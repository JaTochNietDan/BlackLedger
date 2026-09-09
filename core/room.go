package core

import "fmt"

// A location was a set of premises with a list of verbs attached. But half of
// what a player does anywhere is done to somebody who happens to be standing
// there — lending them money, collecting it, putting them on, putting them out,
// taking what they are carrying — and rendering "Lend Perla Mraz money" as one
// more card in a column of verbs throws away the thing the decision is actually
// about, which is Perla Mraz.
//
// This is who is in the room. The core already knows all of it; it has simply
// never been asked the question in this shape.

// Presence is one person where the player is standing, and everything the
// player has earned the right to know about them.
type Presence struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role,omitempty"`
	Faction string `json:"faction,omitempty"`
	// Standing is the one line that says who this person is to the player.
	Standing string `json:"standing"`
	// Trust is what they think of the player, and Sore what they hold against
	// them. Both are only shown for people the player has reason to know.
	Trust int `json:"trust,omitempty"`
	Sore  int `json:"sore,omitempty"`
	// Owes is what they owe, if anything, and Overdue whether it is late.
	Owes    int  `json:"owes,omitempty"`
	Overdue bool `json:"overdue,omitempty"`
	// Yours is whether they answer to the player, Known whether the player has
	// any reason to know their name at all.
	Yours bool `json:"yours,omitempty"`
	// Walking is somebody who is between two addresses right now. Where and
	// WhereID then name where they are *going*, because that is the only place
	// they could be met: reaching a man in the street is not something this
	// game models, and pointing the player at the door he walked out of sends
	// them across the city to an empty room.
	Walking bool `json:"walking,omitempty"`
	Minutes int  `json:"minutes,omitempty"`
	Known   bool `json:"known,omitempty"`
	// Temperament and Manner are what the player has learned of their
	// character; empty for a stranger.
	Temperament string `json:"temperament,omitempty"`
	// Doing is what this person is actually doing right now, in a few words.
	// The city has always known; it has never been asked.
	Doing string `json:"doing,omitempty"`
	// Where they are, for anything looking at the city rather than at one room.
	Where   string `json:"where,omitempty"`
	WhereID string `json:"where_id,omitempty"`
	// Why this person is worth the player's attention, in the order the
	// interface should group them: "yours", "crew", "owes", "sore", "job",
	// "organization", or "street".
	Because string `json:"because,omitempty"`
}

// standingOf is the single line under somebody's name: the most important true
// thing about them from where the player is standing.
func (w *World) standingOf(n *NPC) string {
	switch {
	case w.Inside(n):
		return "In a cell at Ward Street Station"
	case n.Faction == w.PlayerOrganizationID():
		if posted := w.postedWhere(n.ID); posted != "" {
			return "Yours · on the door at " + posted
		}
		return "Yours · " + fmt.Sprintf("$%d a day", MemberWage)
	case len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID:
		return fmt.Sprintf("Your crew · %d loyalty", w.Player.Crew[0].Loyalty)
	case IsOfficial(n.ID):
		if w.Retained(n.ID) {
			o, _ := OfficialByID(n.ID)
			return fmt.Sprintf("%s · paid $%d a day", n.Role, o.Retainer)
		}
		return n.Role
	case w.isRoleHolder(n):
		return n.Role
	case n.Faction != "":
		if n.Rank >= RankLeader {
			return "Head of " + w.factionName(n.Faction)
		}
		return w.factionName(n.Faction)
	case n.Role != "":
		return n.Role
	}
	return "Nobody in particular, yet"
}

// postedWhere is the place somebody of the player's is standing on the door of.
func (w *World) postedWhere(id string) string {
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Posted == id {
			return l.Name
		}
	}
	return ""
}

// doingNow is what somebody is actually doing at this moment, read out of the
// same state everything else is read out of. The ultimate shape of this game is
// a city the player can watch, and a city you can watch is one where every
// person on the screen is visibly occupied with something true.
func (w *World) doingNow(n *NPC) string {
	place, known := PlaceByID(n.Location)
	where := "the district"
	if known {
		where = place.Name
	}
	switch {
	case w.Travelling(n):
		// Said before anything else: a man on the street is not on a door, not
		// on duty, and not behind a desk, whatever his job is.
		if to, ok := PlaceByID(n.Heading); ok {
			out := "Walking to " + to.Name + ", " + counted(max(1, n.Arrives-w.Minute), "minute", "minutes") + " out"
			// The errand usually names the same building, and saying it twice
			// in one line reads as a stutter rather than as a reason.
			if n.Errand != "" && !containsName(n.Errand, to.Name) {
				out += " — " + n.Errand
			}
			return out
		}
		return "Somewhere on the street"
	case w.Inside(n):
		return "Held at Ward Street Station"
	case w.onARound(n.ID):
		// What he is actually doing outranks what his job is called. He read as
		// "Driver, on duty at Bluebird Laundry" while he was standing in it
		// collecting, which is his title rather than his afternoon.
		return "Collecting at " + where
	case w.postedWhere(n.ID) != "":
		return "Standing on the door at " + w.postedWhere(n.ID)
	case IsOfficial(n.ID):
		if w.Retained(n.ID) {
			return "Taking your money and answering your calls"
		}
		return "Working, and expensive to interrupt"
	case w.isRoleHolder(n):
		return n.Role + ", on duty at " + where
	}
	// Somebody who runs premises is at work in them.
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && l.ID == n.Location && prop.Income > 0 {
			if prop.Owner == n.Faction && n.Rank >= RankSoldier {
				return "Running " + l.Name
			}
			if n.Faction != "" && prop.Owner == n.Faction {
				return "Working for " + w.factionName(n.Faction) + " at " + l.Name
			}
		}
	}
	if l := w.LoanTo(n.ID); l != nil {
		if l.Missed > 0 {
			return "Avoiding you, and not doing it well"
		}
		return "Carrying money that is yours"
	}
	if n.Sore >= 30 {
		return "Nursing something they hold against you"
	}
	switch {
	case n.Faction == w.PlayerOrganizationID():
		return "Waiting to be given something to do"
	case n.Rank >= RankLeader:
		return "Holding court at " + where
	case n.Faction != "":
		if w.fighting(n.Faction) != nil {
			return "Armed, and expecting trouble"
		}
		return "About " + w.factionName(n.Faction) + " business"
	case n.Ambition >= 70:
		return "Looking for a way up"
	}
	return "Getting on with the day at " + where
}

// because is the reason this person is on the player's screen at all, which is
// also the order they should be read in. A city of fifty people rendered as
// fifty identical cards is a register; grouped by why they matter, it is a
// list of the people in your life and then everybody else.
func (w *World) because(n *NPC) string {
	switch {
	case n.Faction == w.PlayerOrganizationID():
		return "yours"
	case len(w.Player.Crew) > 0 && w.Player.Crew[0].ID == n.ID:
		return "crew"
	case w.LoanTo(n.ID) != nil:
		return "owes"
	case n.Sore > 0:
		return "sore"
	case IsOfficial(n.ID) || w.isRoleHolder(n) || n.Rank >= RankLeader:
		return "job"
	case n.Faction != "":
		return "organization"
	}
	return "street"
}

var becauseOrder = map[string]int{
	"yours": 0, "crew": 1, "owes": 2, "sore": 3, "job": 4, "organization": 5, "street": 6,
}

// see builds what the player is allowed to know about somebody, wherever they
// are being looked at from.
func (w *World) see(n *NPC) Presence {
	known := w.Known(n)
	p := Presence{
		ID: n.ID, Name: n.Name, Role: n.Role,
		Standing: w.standingOf(n), Doing: w.doingNow(n),
		Yours:   n.Faction == w.PlayerOrganizationID(),
		Known:   known,
		Because: w.because(n),
		WhereID: n.Location,
	}
	if place, ok := PlaceByID(n.Location); ok {
		p.Where = place.Name
	}
	if w.Travelling(n) {
		if to, ok := PlaceByID(n.Heading); ok {
			p.Walking, p.WhereID, p.Where = true, n.Heading, "On the way to "+to.Name
			p.Minutes = max(1, n.Arrives-w.Minute)
		}
	}
	if known {
		p.Faction = w.factionName(n.Faction)
		p.Trust, p.Sore = n.Trust, n.Sore
		p.Temperament = TemperamentOf(n).Label
	}
	if l := w.LoanTo(n.ID); l != nil {
		p.Owes, p.Overdue = l.Owed, l.Missed > 0
	}
	return p
}

// Everyone is the whole city, ordered by why each person matters to the player.
// The People screen used to render every living soul as an identical card in
// whatever order the save happened to hold them — fifty of them, four screens
// of scrolling, with the man who works for you indistinguishable from a docker
// he has never met.
func (w *World) Everyone() []Presence {
	out := []Presence{}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead {
			out = append(out, w.see(n))
		}
	}
	rank := func(p Presence) int {
		at, ok := becauseOrder[p.Because]
		if !ok {
			return len(becauseOrder)
		}
		return at
	}
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if a, b := rank(out[j]), rank(out[i]); a < b || (a == b && out[j].Owes > out[i].Owes) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// PeopleHere is everybody standing where the player is, in the order a person
// would notice them: their own first, then the people they know, then whoever
// else is in the room.
func (w *World) PeopleHere(id string) []Presence {
	out := []Presence{}
	for i := range w.NPCs {
		// Somebody on the street between two addresses is in neither of them.
		if n := &w.NPCs[i]; !n.Dead && n.Location == id && !w.Travelling(n) {
			out = append(out, w.see(n))
		}
	}
	// Yours first, then anybody who owes you, then people you know.
	rank := func(p Presence) int {
		switch {
		case p.Yours:
			return 0
		case p.Owes > 0:
			return 1
		case p.Known:
			return 2
		}
		return 3
	}
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if rank(out[j]) < rank(out[i]) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
