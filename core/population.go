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
	// StreetCount is how many people in this city answer to nobody. It was
	// twenty-six, which was one for every street trade when there were
	// twenty-six of them and ten addresses. There are twenty addresses now and
	// the list is longer, so this is longer with it: a city is as inhabited as
	// the number of people with somewhere to be.
	StreetCount = 62
	// MaxPeople bounds the save. The city recruits from the street rather than
	// inventing people, so this is a ceiling nothing normally approaches — but
	// it approached it once the street was filled out, so it is higher.
	MaxPeople = 400
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
	{"Croupier", "club"}, {"Doorkeeper", "club"}, {"Cigarette seller", "club"},
	{"Stallholder", "market"}, {"Butcher", "market"}, {"Clerk", "market"},
	{"Dealer", "casino"}, {"Floor manager", "casino"},
	{"Landlady", "room"}, {"Boarder", "room"},
	{"Caretaker", "apartment"}, {"Nurse", "apartment"},
	{"Groundsman", "estate"}, {"Housekeeper", "estate"},
	{"Newspaperman", "market"}, {"Photographer", "bar"},
	// The addresses added since this list was written. Half the city had
	// nobody in it: a butcher with no butcher, a cab company with no drivers,
	// a revue bar with nobody on the stage. A city is inhabited or it is a set.
	{"Waiter", "restaurant"}, {"Cook", "restaurant"}, {"Cellarman", "restaurant"},
	{"Marker", "poolhall"}, {"Table hand", "poolhall"},
	{"Boner", "butcher"}, {"Delivery hand", "butcher"}, {"Cold store hand", "butcher"},
	{"Loader", "haulage"}, {"Long-haul driver", "haulage"}, {"Yard clerk", "haulage"},
	{"Croupier", "goldenlily"}, {"Cashier", "goldenlily"}, {"Doorkeeper", "goldenlily"},
	{"Presser", "steamworks"}, {"Van driver", "steamworks"}, {"Sorter", "steamworks"},
	{"Dancer", "burlesque"}, {"Stage hand", "burlesque"}, {"Bandleader", "burlesque"},
	{"Dispatcher", "cabstand"}, {"Night driver", "cabstand"}, {"Fitter", "cabstand"},
	{"Salesperson", "dealer"}, {"Lot hand", "dealer"}, {"Finance clerk", "dealer"},
	{"Mechanic", "archway"}, {"Sprayer", "archway"}, {"Parts keeper", "archway"},
	{"Crane driver", "scrapyard"}, {"Cutter", "scrapyard"}, {"Weighbridge clerk", "scrapyard"},
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
	trade := w.nextStreetTrade()
	id := fmt.Sprintf("street-%d", len(w.NPCs)+1)
	for w.NPC(id) != nil {
		id += "x"
	}
	w.NPCs = append(w.NPCs, NPC{
		ID: id, Name: name, Role: trade.role, Voice: w.voiceFor(name), Color: "#7f7a6c",
		Location: trade.place, Rank: RankAssociate,
		Ambition: 15 + int(w.WorldRandom()*60), Skill: 20 + int(w.WorldRandom()*50),
	})
	w.SettlePurses()
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
		was := lowerFirst(roleOrNobody(pick))
		if pick.Role != "" {
			was = article(was) + " " + was
		}
		w.Log(pick.Name+" signs on with "+f.Name, fmt.Sprintf("They were %s a week ago. %s %s short of people and not asking many questions.",
			was, Leads(f.Name), Agree(f.Name, "is", "are")), "politics")
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
	// The person standing in front of the player right now. A scene names its
	// speaker by id and nothing else, so forgetting them leaves the open
	// conversation pointing at nobody — and declining the offer looked their
	// name up without asking whether they were still there.
	if w.Event != nil && w.Event.Speaker == id {
		return true
	}
	// An arrangement the player walked away from mid-job is still theirs to
	// come back to, and it remembers who it is with the same way.
	if w.SuspendedJob != nil && w.SuspendedJob.Scene != nil && w.SuspendedJob.Scene.Speaker == id {
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
	// Everybody is in exactly one of three places, and the screen said so
	// without leaving room for the third. Somebody doing one of the city's
	// jobs — the four officials, the fixer, the driver — answers neither to a
	// family nor to nobody, and the header read "29 answer to an organization
	// and 17 to nobody" of a city of 52.
	living, street := len(w.People()), len(w.Civilians())
	return map[string]any{
		"living": living, "organized": organized,
		"street": street, "jobs": max(0, living-organized-street),
		"known": len(w.Cast()),
	}
}

// nextStreetTrade picks what somebody does for a living. Trades are DEALT
// rather than drawn: the least-taken jobs are found first and one of those is
// chosen, so every way of earning a living in this city is somebody's before
// any of them is a second person's.
//
// Drawing at random with replacement put three newspapermen in the market at
// once — seen in the browser, not in a test — while a third of the city's jobs
// had nobody doing them at all. A room where everybody is the same thing is a
// room of one person repeated, which is a shorter city than it looks.
func (w *World) nextStreetTrade() struct{ role, place string } {
	taken := map[string]int{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Faction != "" || IsOfficial(n.ID) {
			continue
		}
		taken[n.Role+"@"+n.Location]++
	}
	fewest := -1
	for _, trade := range streetTrades {
		if held := taken[trade.role+"@"+trade.place]; fewest < 0 || held < fewest {
			fewest = held
		}
	}
	open := streetTrades[:0:0]
	for _, trade := range streetTrades {
		if taken[trade.role+"@"+trade.place] == fewest {
			open = append(open, trade)
		}
	}
	// Which of the open ones is still chance, so two cities from two seeds are
	// not the same city with the same people in the same order.
	return open[int(w.WorldRandom()*float64(len(open)))%len(open)]
}
