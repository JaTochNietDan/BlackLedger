package core

import "fmt"

// The player could read a rival's exact strength, their exact standing and the
// exact money in their books, off a screen, having never spoken to anybody
// inside them. Nothing in the fiction accounted for it: a man who has never met
// anybody in an organization does not know what is in its accounts.
//
// What you can find out about somebody now depends on who you know. It is the
// first thing in this game that rewards dealing with people rather than
// spending money, and the first that can be wrong.

const (
	// EnquiryCost is what asking around costs.
	EnquiryCost = 90
	// EnquiryMinutes is the afternoon it takes.
	EnquiryMinutes = 90
	// EnquiryLasts is how long what you learn stays current. A week, because
	// an organization is a different shape a week later.
	EnquiryLasts = 10080
)

// Intelligence is how much the player can find out about an organization, from
// nothing to everything.
func (w *World) Intelligence(id string) int {
	if id == w.PlayerOrganizationID() {
		return 3 // your own books
	}
	level := 0
	// A network hears the ordinary things about everybody.
	if w.Reach() >= 2 {
		level++
	}
	// Somebody inside it you have actually dealt with is worth more than any
	// number of people who have heard something.
	for _, n := range w.Members(id) {
		if n.Trust > 0 {
			level++
			break
		}
	}
	// And asking around directly, while it is still current.
	if until, ok := w.Player.Enquiries[id]; ok && until > w.Minute {
		level++
	}
	return min(3, level)
}

// EnquiryReadiness explains why nobody can be asked about, or returns "".
func (w *World) EnquiryReadiness(id string) string {
	if w.faction(id) == nil || id == w.PlayerOrganizationID() {
		return "There is nobody of that description to ask about"
	}
	if w.Reach() < 1 {
		return "You have nobody to ask"
	}
	if until, ok := w.Player.Enquiries[id]; ok && until > w.Minute {
		return "What you were told is still current"
	}
	if w.Player.Cash < EnquiryCost {
		return "Not enough cash"
	}
	return ""
}

// AskAround buys a week of knowing more about somebody than the street does.
func (w *World) AskAround(id string) error {
	if reason := w.EnquiryReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(EnquiryCost); err != nil {
		return err
	}
	if w.Player.Enquiries == nil {
		w.Player.Enquiries = map[string]int{}
	}
	w.Player.Enquiries[id] = w.Minute + EnquiryLasts
	f := w.faction(id)
	w.Log("You ask about "+f.Name, fmt.Sprintf("$%d in the right pockets. What comes back is current for about a week, and a week is a long time in this city.", EnquiryCost), "politics")
	return nil
}

// roughly describes a number the way somebody passing it on would, which is
// what the player gets before they know anybody worth knowing.
// handful describes a headcount without counting it. roughly's words are made
// for money and strength — "as much as anybody" is a strange thing to say about
// eight people — so the impression of a crowd gets its own.
func handful(n int) string {
	switch {
	case n == 0:
		return "nobody worth naming"
	case n < 4:
		return "a handful"
	case n < 8:
		return "a fair few"
	}
	return "a great many"
}

func roughly(n int) string {
	switch {
	case n < 25:
		return "not much"
	case n < 50:
		return "something"
	case n < 75:
		return "a good deal"
	}
	return "as much as anybody"
}

// PublicFaction is an organization as the player can see it. The name and
// whether you are at odds are common knowledge; everything else is not.
type PublicFaction struct {
	Headquarters string `json:"headquarters,omitempty"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	// Leader is known once anybody will talk to you about them.
	Leader string `json:"leader"`
	// Goodwill is what they think of the player, which the player can always
	// tell from how they are treated.
	Goodwill int `json:"goodwill"`
	// Strength and Money are exact only when the player has somebody inside.
	// Below that they are described rather than counted.
	Power    int    `json:"power"`
	Cash     int    `json:"cash"`
	Strength string `json:"strength"`
	Money    string `json:"money"`
	// Hands is the people figure said the way strength and money are said. It
	// used to be only the number, dropped when it was zero, so a card lost a
	// column rather than admitting it did not know — and a family with nobody
	// left read exactly like a family nobody would talk about.
	Hands string `json:"hands"`
	// Knowledge is how much of this is worth trusting, from nothing to
	// everything, so the interface can say so.
	Knowledge int `json:"knowledge"`
	// LeaderID is who to draw, once anybody will name them.
	LeaderID string `json:"leader_id,omitempty"`
	// Holdings is the ground they are standing on, by name. Premises are the
	// most public thing an organization has: anybody can walk past them.
	Holdings []string `json:"holdings"`
	// Members is how many people answer to them, and Seats the number of those
	// the player could actually put a name to.
	Members, Known int `json:"-"`
	People         int `json:"people,omitempty"`
	// Fighting is who they are at war or at odds with, in words.
	Fighting []string `json:"fighting"`
	// Standing is what their goodwill means, said plainly.
	Standing string `json:"standing,omitempty"`
	// Yours is whether this is the player's own organization.
	Yours bool `json:"yours,omitempty"`
}

// standingWith puts a number nobody can read into words anybody can. The screen
// said "+45" and "-63" and left the player to work out what either meant.
func standingWith(goodwill int) string {
	switch {
	case goodwill <= -70:
		return "They have decided about you"
	case goodwill <= -35:
		return "Hostile"
	case goodwill < -10:
		return "They do not like you"
	case goodwill <= 10:
		return "They have no opinion of you"
	case goodwill < 35:
		return "Cordial"
	case goodwill < 70:
		return "They think well of you"
	}
	return "You are as good as one of theirs"
}

// PublicFactions is every organization as the player can see it.
func (w *World) PublicFactions() []PublicFaction {
	out := []PublicFaction{}
	for i := range w.Factions {
		f := &w.Factions[i]
		level := w.Intelligence(f.ID)
		entry := PublicFaction{
			ID: f.ID, Name: f.Name, Goodwill: f.Goodwill, Knowledge: level,
			Strength: "nobody will say", Money: "nobody will say",
		}
		if level >= 1 || f.ID == w.PlayerOrganizationID() {
			entry.Headquarters = w.Headquarters(f.ID)
		}
		if level >= 1 {
			entry.Leader = f.Leader
			entry.Strength = roughly(f.Power)
			entry.Money = roughly(min(100, f.Cash/300))
		}
		if level >= 2 {
			entry.Power = f.Power
			entry.Strength = fmt.Sprintf("%d of a hundred", f.Power)
		}
		if level >= 3 {
			entry.Cash = f.Cash
			entry.Money = fmt.Sprintf("$%d", f.Cash)
		}
		// Premises are the most public thing an organization has: anybody can
		// walk past them, so they need no informant.
		for _, id := range w.FamilyHoldings(f.ID) {
			if place, ok := PlaceByID(id); ok {
				entry.Holdings = append(entry.Holdings, place.Name)
			}
		}
		entry.Yours = f.ID == w.PlayerOrganizationID()
		entry.Standing = standingWith(f.Goodwill)
		if entry.Yours {
			entry.Standing = "Yours"
		}
		if level >= 1 {
			for _, n := range w.People() {
				if n.Name == f.Leader {
					entry.LeaderID = n.ID
				}
			}
		}
		// How many people answer to them is a thing you need somebody inside
		// to count, but that there is a quarrel at all is public.
		entry.Hands = "nobody will say"
		if level >= 1 {
			entry.Hands = handful(len(w.Members(f.ID)))
		}
		if level >= 2 {
			entry.People = len(w.Members(f.ID))
			entry.Hands = counted(entry.People, "person", "people")
			if entry.People == 0 {
				entry.Hands = "nobody left"
			}
		}
		for _, c := range w.Conflicts {
			if c.State != "war" && c.State != "feud" {
				continue
			}
			other := ""
			if c.A == f.ID {
				other = c.B
			} else if c.B == f.ID {
				other = c.A
			}
			if other == "" {
				continue
			}
			word := "at war with "
			if c.State == "feud" {
				word = "at odds with "
			}
			entry.Fighting = append(entry.Fighting, word+w.factionName(other))
		}
		out = append(out, entry)
	}
	return out
}
