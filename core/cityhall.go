package core

import "fmt"

// A detective who will lose a file for money already existed. What did not is
// the rest of the building: a mayor, a commissioner, and the fact that a man
// with a salary and a mortgage is the easiest thing in this city to buy and the
// hardest thing to keep bought.
//
// An official on a retainer is a standing arrangement rather than a favour. It
// costs every day, it is worth different things depending on who it is, it is
// cut the moment holding it would cost them more than it pays, and somebody
// with more money than the player can outbid them for it. They are also people:
// they can be killed, by the player or by anybody else, and the city notices
// that more than it notices anything else that happens in it.

// Official is somebody in the building who can be reached.
type Official struct {
	ID, Name, Role string
	// Detail is what an arrangement with them is worth, in words.
	Detail string
	// Retainer is what they cost a day, and Opening what it takes to start.
	Retainer, Opening int
	// Ceiling is the police attention past which they will not be seen with
	// the player at any price.
	Ceiling int
}

var officials = []Official{
	{ID: "commissioner", Name: "Commissioner Vance", Role: "Police commissioner",
		Detail:   "Files go to the bottom of piles. Raids begin 20 attention later than they otherwise would, and nothing is ever forfeited while he is paid.",
		Retainer: 45, Opening: 900, Ceiling: 78},
	{ID: "mayor", Name: "Mayor Ellis Crane", Role: "Mayor of Bellwether",
		Detail:   "Licences, inspections and the right words in the right rooms. Every business of yours earns a fifth more while he is paid.",
		Retainer: 60, Opening: 1400, Ceiling: 65},
}

const (
	// RetainerRelief is the extra attention a commissioner absorbs before
	// anybody comes to the door.
	RetainerRelief = 20
	// MayorTake is the share a mayor's licences add to what a business earns.
	MayorTake = .2
	// OfficialMinutes is how long an arrangement takes to make.
	OfficialMinutes = 90
	// CityHall is where these arrangements are made. The exchange, because
	// nobody makes them in the building itself.
	CityHall = "market"
	// OutbidBy is how much richer than the player an organization has to be
	// before an official quietly starts taking their money instead.
	OutbidBy = 4000
)

// OfficialByID is one of them, and whether they exist.
func OfficialByID(id string) (Official, bool) {
	for _, o := range officials {
		if o.ID == id {
			return o, true
		}
	}
	return Official{}, false
}

// IsOfficial reports whether somebody holds a title rather than a position in
// an organization. They are people like anybody else — they can be resented,
// and they can be killed — but they do not go out at night taking tills.
func IsOfficial(id string) bool {
	_, ok := OfficialByID(id)
	return ok
}

// Officials is everybody who can be reached, for the interface.
func Officials() []Official { return officials }

// ensureOfficials puts them in the city as people, so everything that can
// happen to anybody can happen to them. Idempotent, and called wherever they
// might first be needed.
func (w *World) ensureOfficials() {
	for _, o := range officials {
		if w.NPC(o.ID) != nil {
			continue
		}
		w.NPCs = append(w.NPCs, NPC{
			ID: o.ID, Name: o.Name, Role: o.Role, Voice: w.voiceFor(o.Name),
			Color: "#7c8791", Location: CityHall, Rank: RankLieutenant,
			Ambition: 55, Skill: 30,
		})
	}
}

// Retained reports whether an official is currently taking the player's money
// and is still alive to be worth it.
func (w *World) Retained(id string) bool {
	for _, held := range w.Player.Retainers {
		if held != id {
			continue
		}
		if n := w.NPC(id); n == nil || n.Dead {
			return false
		}
		return true
	}
	return false
}

// RetainerCost is what the arrangements cost a day, joining the rest of the
// bill. An official who is dead costs nothing, which is one way to end one.
func (w *World) RetainerCost() int {
	total := 0
	for _, id := range w.Player.Retainers {
		if !w.Retained(id) {
			continue
		}
		o, _ := OfficialByID(id)
		total += o.Retainer
	}
	return total
}

// Outbid reports whether an organization with more money than the player is
// paying the same official, which is what makes an arrangement something to
// defend rather than something to buy once.
func (w *World) Outbid(id string) bool {
	if !w.Retained(id) {
		return false
	}
	// Money alone is not enough: every organization in this city has more of it
	// than a man starting out. Somebody outbids the player when they are richer
	// *and* have a reason to want the player without friends in that building.
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.Cash > w.Player.Cash+OutbidBy && f.Goodwill <= -20 {
			return true
		}
	}
	return false
}

// RaidRelief is the attention a commissioner absorbs before anybody arrives,
// and nothing at all while somebody richer is paying him too.
func (w *World) RaidRelief() int {
	if w.Retained("commissioner") && !w.Outbid("commissioner") {
		return RetainerRelief
	}
	return 0
}

// LicenceTake is what a mayor's arrangement adds to what a business earns.
func (w *World) LicenceTake() float64 {
	if w.Retained("mayor") && !w.Outbid("mayor") {
		return MayorTake
	}
	return 0
}

// RetainerReadiness explains why an arrangement cannot be made, or returns "".
func (w *World) RetainerReadiness(id string) string {
	o, ok := OfficialByID(id)
	if !ok {
		return "There is nobody of that description"
	}
	if w.Player.Location != CityHall {
		return "This is not arranged here"
	}
	if n := w.NPC(id); n != nil && n.Dead {
		return "They are dead. Whoever replaces them does not know you"
	}
	if w.Retained(id) {
		return "That arrangement already stands"
	}
	if w.Player.Heat > w.OfficialCeiling(o) {
		return fmt.Sprintf("Nobody in that building will be seen with you above %d attention", w.OfficialCeiling(o))
	}
	if w.Presence() < 25 {
		return "They would not take a call from you"
	}
	if w.Player.Cash < w.OfficialOpening(o) {
		return fmt.Sprintf("It takes $%d to open the conversation", w.OfficialOpening(o))
	}
	return ""
}

// Retain opens an arrangement. The opening payment is the introduction; the
// daily cost is what keeps it.
func (w *World) Retain(id string) error {
	if reason := w.RetainerReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	w.ensureOfficials()
	o, _ := OfficialByID(id)
	if err := w.Pay(w.OfficialOpening(o)); err != nil {
		return err
	}
	w.Player.Retainers = append(w.Player.Retainers, id)
	w.MeetPerson(id)
	w.Log("An arrangement with "+o.Name, fmt.Sprintf("$%d to open it and $%d a day to keep it. %s It ends the day he decides you are worth more trouble than money.", o.Opening, o.Retainer, o.Detail), "politics")
	return nil
}

// EndRetainer is the player letting one go, which is the only way to stop
// paying without somebody dying.
func (w *World) EndRetainer(id string) error {
	if !w.Retained(id) {
		return fmt.Errorf("there is no such arrangement")
	}
	kept := []string{}
	for _, held := range w.Player.Retainers {
		if held != id {
			kept = append(kept, held)
		}
	}
	w.Player.Retainers = kept
	o, _ := OfficialByID(id)
	w.Log("The arrangement with "+o.Name+" ends", "You stop paying. He does not argue, which tells you what it was worth to him.", "politics")
	return nil
}

// CityHallDay is what these arrangements cost when they stop being worth it to
// the man taking the money. Nobody in that building goes down with the player.
func (w *World) CityHallDay() {
	for _, id := range append([]string{}, w.Player.Retainers...) {
		if !w.Retained(id) {
			continue
		}
		o, _ := OfficialByID(id)
		if w.Player.Heat <= w.OfficialCeiling(o) {
			continue
		}
		w.EndRetainerQuietly(id)
		w.Log(o.Name+" is not taking calls", fmt.Sprintf("Your attention is at %d and he has a pension. The arrangement is over and the money you paid to open it is not coming back.", w.Player.Heat), "danger")
	}
}

// EndRetainerQuietly is an official cutting the player loose, which needs no
// permission from anybody.
func (w *World) EndRetainerQuietly(id string) {
	kept := []string{}
	for _, held := range w.Player.Retainers {
		if held != id {
			kept = append(kept, held)
		}
	}
	w.Player.Retainers = kept
}

// OfficialKilled is what the city does when somebody kills a man with a title.
// It is the loudest thing that can happen here, and it happens to everybody
// rather than only to whoever did it.
func (w *World) OfficialKilled(id string) {
	o, ok := OfficialByID(id)
	if !ok {
		return
	}
	w.EndRetainerQuietly(id)
	w.Player.Heat = min(100, w.Player.Heat+35)
	// Everybody pays for it, not only whoever did it. An operation on this
	// scale costs every organization in the city people and money.
	for i := range w.Factions {
		f := &w.Factions[i]
		f.Power = max(10, f.Power-8)
		f.Cash = max(0, f.Cash-1200)
	}
	w.Log("They will turn the city over", fmt.Sprintf("%s is dead. Every man in this city with a name is going to spend the next month explaining where he was, and that includes you.", o.Name), "danger")
	w.Report("police", "CITY REELS AS "+upper(o.Name)+" IS KILLED",
		fmt.Sprintf("%s, %s, was killed today. The police have announced what they describe as an unprecedented operation against organized crime in the city. No arrests have been made.", o.Name, o.Role))
}

// RetainerDescription is what the player is paying for, for the interface.
func (w *World) RetainerDescription() []map[string]any {
	out := []map[string]any{}
	for _, o := range officials {
		if !w.Retained(o.ID) {
			continue
		}
		out = append(out, map[string]any{
			"id": o.ID, "name": o.Name, "role": o.Role,
			"retainer": o.Retainer, "outbid": w.Outbid(o.ID), "detail": o.Detail,
		})
	}
	return out
}

// OfficialOpening is what an arrangement costs to make, which is more while the
// city is looking: a man with a career to protect wants more for the risk.
func (w *World) OfficialOpening(o Official) int {
	return o.Opening * (100 + w.ScrutinyPremium()) / 100
}

// OfficialCeiling is the attention past which they will not be seen with
// anybody, which falls while the city is looking.
func (w *World) OfficialCeiling(o Official) int {
	if !w.UnderCrackdown() {
		return o.Ceiling
	}
	return max(20, o.Ceiling-(w.Scrutiny()-ScrutinyCrackdown)-15)
}
