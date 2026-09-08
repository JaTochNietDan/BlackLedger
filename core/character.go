package core

import "fmt"

// Everybody in this city has a name, a job, an ambition and a number for how
// good they are at things. None of it made any two of them behave differently:
// a person with 80 ambition and a person with 80 skill did the same thing in
// the same situation, only more or less often.
//
// A temperament is the part of somebody that decides how they do a thing rather
// than whether. It is derived from their name, so it costs nothing in a save,
// it is the same person every time the campaign is loaded, and a person the
// player has come to expect something from keeps doing it.

// Temperament is how somebody goes about things.
type Temperament struct {
	ID string
	// Label is the word for it, and Detail what it means for what they do.
	Label, Detail string
	// Nerve scales how readily they act on a grievance.
	Nerve float64
	// Poise scales how well it goes when they do.
	Poise float64
	// Greed scales what they take when they help themselves to something.
	Greed float64
	// Loyal people do not move against their own organization.
	Loyal bool
}

var temperaments = []Temperament{
	{ID: "hot", Label: "Hot-headed", Nerve: 1.6, Poise: .8, Greed: 1,
		Detail: "Moves on a grievance long before it is wise to, and moves badly."},
	{ID: "careful", Label: "Careful", Nerve: .55, Poise: 1.35, Greed: .9,
		Detail: "Waits, and is rarely the one who ends up in the river."},
	{ID: "greedy", Label: "Grasping", Nerve: 1.1, Poise: 1, Greed: 1.5,
		Detail: "Takes more than they came for, every time."},
	{ID: "loyal", Label: "Loyal", Nerve: .8, Poise: 1.1, Greed: .8, Loyal: true,
		Detail: "Will not move against their own people, whatever they are owed."},
	{ID: "vain", Label: "Vain", Nerve: 1.2, Poise: .95, Greed: 1.2,
		Detail: "Wants to be seen doing it, which is how anybody ever finds out."},
}

// TemperamentOf is who somebody is, derived from their name so it is stable
// across saves and costs nothing to store. Two people called the same thing
// would be the same kind of person, which is a price worth paying for never
// having to migrate it.
func TemperamentOf(n *NPC) Temperament {
	if n == nil {
		return temperaments[1]
	}
	// A plain weighted sum of the letters landed two of the five temperaments
	// on more than a third of the city each. This mixes properly, which is the
	// difference between five kinds of person and two.
	hash := uint32(2166136261)
	for _, c := range n.Name {
		hash ^= uint32(c)
		hash *= 16777619
	}
	hash ^= hash >> 15
	return temperaments[int(hash%uint32(len(temperaments)))]
}

// Nerve is how readily this person acts on something they are owed.
func (w *World) Nerve(n *NPC) float64 {
	return float64(n.Ambition) / 220 * TemperamentOf(n).Nerve
}

// Poise is what they bring to it when they do act.
func (w *World) Poise(n *NPC) int {
	return max(5, int(float64(n.Skill)*TemperamentOf(n).Poise))
}

// Known reports whether the player has any reason to know who somebody is.
// Trust is built by dealing with them; anybody at the head of an organization
// is a name everybody knows.
func (w *World) Known(n *NPC) bool {
	return n != nil && (n.Trust > 0 || n.Rank >= RankLeader || w.Reach() >= 3)
}

// Dossier is everything the player has actually learned about somebody, built
// out of records the city committed rather than out of anything hidden. A
// stranger is a stranger: this returns nothing for somebody they have never
// had a reason to hear of.
func (w *World) Dossier(id string) map[string]any {
	n := w.NPC(id)
	if n == nil || !w.Known(n) {
		return nil
	}
	temperament := TemperamentOf(n)
	entry := map[string]any{
		"id": n.ID, "name": n.Name, "role": n.Role,
		"temperament": temperament.Label, "manner": temperament.Detail,
		"organization": w.factionName(n.Faction), "known_for": w.knownFor(n),
	}
	if n.Faction == "" {
		entry["organization"] = "Nobody in particular"
	}
	return entry
}

// knownFor is what the city's own record says this person has done, in the
// order it happened. It reads History rather than any private ledger, so the
// player can only ever be told things that were reported.
func (w *World) knownFor(n *NPC) []string {
	out := []string{}
	for _, r := range w.History {
		if len(out) >= 4 {
			break
		}
		if containsName(r.Title, n.Name) || containsName(r.Text, n.Name) {
			out = append(out, r.Title)
		}
	}
	return out
}

func containsName(haystack, name string) bool {
	if name == "" || len(name) > len(haystack) {
		return false
	}
	for i := 0; i+len(name) <= len(haystack); i++ {
		if haystack[i:i+len(name)] == name {
			return true
		}
	}
	return false
}

// Cast is who the player knows, for the interface. Everybody else is a face in
// the street until they have a reason to be more than that.
func (w *World) Cast() []map[string]any {
	out := []map[string]any{}
	for _, n := range w.People() {
		if entry := w.Dossier(n.ID); entry != nil {
			out = append(out, entry)
		}
	}
	return out
}

// MeetPerson is how somebody stops being a stranger: the player dealt with them
// directly. Called wherever a named person and the player actually interact.
func (w *World) MeetPerson(id string) {
	if n := w.NPC(id); n != nil && n.Trust == 0 {
		n.Trust = 1
		w.Log("You know "+n.Name+" now", fmt.Sprintf("%s. %s", describeStanding(n, w), TemperamentOf(n).Detail), "personal")
	}
}
