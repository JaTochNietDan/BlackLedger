package core

import "fmt"

// A city with eight people in it is a cast, not a city. Every system built so
// far — grudges, coups, muggings, wars that kill people — was operating on a
// population small enough that one bad month emptied it.
//
// This fills the streets. Organizations are the size organizations are, there
// are people in this city who answer to nobody and work for a living, and both
// populations are maintained: a family below strength recruits, and somebody
// who was nobody last month is somebody's soldier now.

const (
	// FamilySize is how many people an established organization keeps. Below
	// this it recruits; it never grows past it on its own.
	FamilySize = 9
	// SplinterSize is what a new organization starts with and works up from.
	SplinterSize = 4
	// StreetCount is how many people in this city answer to nobody.
	StreetCount = 26
	// MaxPeople bounds the save. The city recruits from the street rather than
	// inventing people, so this is a ceiling nothing normally approaches.
	MaxPeople = 200
	// RecruitChance is how often an understrength organization takes somebody
	// on, per day.
	RecruitChance = .25
)

// streetTrades are the people who are not in this business and live here
// anyway. They have somewhere to be and something to lose, which is all it
// takes to be worth robbing, killing or recruiting.
var streetTrades = []struct{ role, place string }{
	{"Barman", "bar"}, {"Bookmaker", "bar"}, {"Cab driver", "bar"},
	{"Docker", "docks"}, {"Ship's clerk", "docks"}, {"Net mender", "docks"},
	{"Laundress", "laundry"}, {"Presser", "laundry"},
	{"Mechanic", "garage"}, {"Panel beater", "garage"},
	{"Croupier", "club"}, {"Doorman", "club"}, {"Cigarette girl", "club"},
	{"Stallholder", "market"}, {"Butcher", "market"}, {"Clerk", "market"},
	{"Dealer", "casino"}, {"Floor manager", "casino"},
	{"Landlady", "room"}, {"Boarder", "room"},
	{"Caretaker", "apartment"}, {"Nurse", "apartment"},
	{"Groundsman", "estate"}, {"Housekeeper", "estate"},
	{"Newspaperman", "market"}, {"Photographer", "bar"},
}

// AddCivilian puts somebody in the city who answers to nobody.
func (w *World) AddCivilian() *NPC {
	if len(w.NPCs) >= MaxPeople {
		return nil
	}
	name, ok := w.newPersonName()
	if !ok {
		return nil
	}
	trade := streetTrades[int(w.WorldRandom()*float64(len(streetTrades)))%len(streetTrades)]
	id := fmt.Sprintf("street-%d", len(w.NPCs)+1)
	for w.NPC(id) != nil {
		id += "x"
	}
	w.NPCs = append(w.NPCs, NPC{
		ID: id, Name: name, Role: trade.role, Voice: w.voiceFor(name), Color: "#7f7a6c",
		Location: trade.place, Rank: RankAssociate,
		Ambition: 15 + int(w.WorldRandom()*60), Skill: 20 + int(w.WorldRandom()*50),
	})
	return &w.NPCs[len(w.NPCs)-1]
}

// Civilians is everybody alive who answers to nobody and holds no title.
func (w *World) Civilians() []*NPC {
	out := []*NPC{}
	for _, n := range w.People() {
		// Somebody doing one of the city's jobs is not loose on the street:
		// no family recruits the detective, and nothing replaces the fixer
		// with the fixer.
		if n.Faction == "" && !IsOfficial(n.ID) && !w.isCrew(n.ID) && !w.isRoleHolder(n) {
			out = append(out, n)
		}
	}
	return out
}

func (w *World) isCrew(id string) bool {
	for _, c := range w.Player.Crew {
		if c.ID == id {
			return true
		}
	}
	return false
}

// Populate brings the city up to strength: every organization to its size, and
// enough people on the street for it to be a street. Idempotent, so it can be
// called on a new world and on a save that predates any of this.
func (w *World) Populate() {
	for i := range w.Factions {
		w.fillOut(&w.Factions[i], FamilySize)
	}
	// Jobs are filled before the street is counted, because filling one takes
	// somebody off the street — doing it the other way round meant the city
	// grew by one every time it was counted.
	w.FillRoles()
	for len(w.Civilians()) < StreetCount && len(w.NPCs) < MaxPeople {
		if w.AddCivilian() == nil {
			break
		}
	}
}

// fillOut takes an organization up to a size, giving each new person a standing
// that makes sense: a couple of lieutenants and the rest below them.
func (w *World) fillOut(f *Faction, size int) {
	for len(w.Members(f.ID)) < size && len(w.NPCs) < MaxPeople {
		role, rank := "Soldier", RankSoldier
		if lieutenants(w.Members(f.ID)) < 2 {
			role, rank = "Lieutenant", RankLieutenant
		} else if w.WorldRandom() < .35 {
			role, rank = "Associate", RankAssociate
		}
		if w.AddMember(f.ID, role, rank, w.homeOf(f.ID)) == nil {
			return
		}
	}
}

func lieutenants(members []*NPC) int {
	n := 0
	for _, m := range members {
		if m.Rank >= RankLieutenant && m.Rank < RankLeader {
			n++
		}
	}
	return n
}

// RecruitDay is how an organization that has lost people gets more. It takes
// them off the street rather than inventing them, so the city's population moves
// between the two rather than only growing.
func (w *World) RecruitDay() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if len(w.Members(f.ID)) >= FamilySize || w.WorldRandom() >= RecruitChance {
			continue
		}
		// Whoever wants it most, out of the people who are near them.
		var pick *NPC
		for _, n := range w.Civilians() {
			if n.Location != w.homeOf(f.ID) {
				continue
			}
			if pick == nil || n.Ambition > pick.Ambition {
				pick = n
			}
		}
		if pick == nil {
			// Nobody local. The city makes another one if there is room.
			if made := w.AddCivilian(); made != nil {
				continue
			}
			continue
		}
		pick.Faction, pick.Rank, pick.Role = f.ID, RankSoldier, "Soldier"
		pick.Location = w.homeOf(f.ID)
		w.Log(pick.Name+" signs on with "+f.Name, fmt.Sprintf("They were %s a week ago. %s is short of people and not asking many questions.", lowerFirst(roleOrNobody(pick)), f.Name), "politics")
	}
}

func roleOrNobody(n *NPC) string {
	if n.Role == "" {
		return "nobody in particular"
	}
	return n.Role
}

// PrunePeople keeps the save bounded by forgetting the dead once nothing refers
// to them. The record of the death survives in the city's history and in the
// Herald; what goes is the empty person.
func (w *World) PrunePeople() {
	if len(w.NPCs) <= MaxPeople*3/4 {
		return
	}
	kept := w.NPCs[:0]
	for _, n := range w.NPCs {
		if n.Dead && !w.referenced(n.ID) {
			continue
		}
		kept = append(kept, n)
	}
	w.NPCs = kept
}

// referenced reports whether anything still needs this person to exist.
func (w *World) referenced(id string) bool {
	if w.isCrew(id) || IsOfficial(id) {
		return true
	}
	if n := w.NPC(id); n != nil && w.isRoleHolder(n) {
		return true
	}
	for _, c := range w.Contracts {
		if c.Target == id {
			return true
		}
	}
	for _, g := range w.Grudges {
		if g.Holder == id || g.Against == id {
			return true
		}
	}
	for _, c := range w.Commissions {
		if c.Target == id && !c.Done && !c.Failed {
			return true
		}
	}
	return false
}

// PopulationSummary is the shape of the city, for the interface.
func (w *World) PopulationSummary() map[string]any {
	organized := 0
	for i := range w.Factions {
		organized += len(w.Members(w.Factions[i].ID))
	}
	return map[string]any{
		"living": len(w.People()), "organized": organized,
		"street": len(w.Civilians()), "known": len(w.Cast()),
	}
}
