package core

import "fmt"

// The top bar grew to eight figures and never once said what any of them meant.
// A new player reads "PRESENCE 253" and "HEAT 10" and has no way to learn what
// either does short of dying of it. Worse, the city's temperature was drawn as
// a number above the words "The city is ordinary", which parses as two ordinary
// cities.
//
// The explanations belong here rather than in the interface, because they are
// about rules and they contain the actual thresholds. Written as copy in the
// front end they would rot the way the Guide rotted; written here they change
// when the constant changes.

// Stat is one figure in the top bar: what it is, what it reads, and what it
// means in a sentence a new player can act on.
type Stat struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
	// Note is the qualifier under the number, when there is one.
	Note string `json:"note,omitempty"`
	// Meaning is the sentence, with the real numbers in it.
	Meaning string `json:"meaning"`
	// Warn is whether this figure is currently something to worry about.
	Warn bool `json:"warn,omitempty"`
}

// Dashboard is the top bar, explained.
func (w *World) Dashboard() []Stat {
	p := &w.Player
	out := []Stat{
		{ID: "cash", Label: "Cash on hand", Value: cash(p.Cash),
			Meaning: fmt.Sprintf("What you can spend today. The day costs $%d whether you earn anything or not.", w.DailyCost()),
			Warn:    p.Cash < w.DailyCost()},
		{ID: "respect", Label: "Respect", Value: itoa(p.Respect),
			Meaning: fmt.Sprintf("What this city thinks you are worth. It opens doors: %d to be filed with the families, and it is what somebody weighs before deciding whether to move on you.", OrganizationStanding)},
		{ID: "crew", Label: "Presence", Value: itoa(w.Presence()),
			Meaning: "Respect plus how you are dressed and what you drive. It decides the odds when you lean on somebody, and whether they can name you afterwards."},
		{ID: "health", Label: "Health", Value: itoa(p.Health),
			Meaning: "How much you can take. Rest restores it; being caught without it is how a life ends.",
			Warn:    p.Health < 40},
		{ID: "heat", Label: "Attention", Value: itoa(p.Heat),
			Meaning: fmt.Sprintf("How interested the police are in you specifically. Past %d they come to the door; past %d they take the premises.", RaidThreshold, ForfeitThreshold),
			Warn:    p.Heat >= RaidThreshold},
	}
	if w.Incorporated() {
		f := w.PlayerOrganization()
		out = append(out, Stat{ID: "families", Label: "Your organization", Value: itoa(f.Power), Note: f.Name,
			Meaning: "What your organization is worth in a fight: your people, your ground and your standing. It decides who dares move on you."})
	}
	out = append(out, Stat{ID: "city", Label: "The city", Value: itoa(w.Scrutiny()), Note: w.scrutinyWord(),
		Meaning: fmt.Sprintf("How hard the whole city is being looked at, from everything the Herald carried. It is not about you. Past %d the police stop waiting and raids begin up to %d attention sooner, for everybody.", ScrutinyCrackdown, ScrutinyRaidRelief),
		Warn:    w.UnderCrackdown()})
	return out
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }

// cash writes a sum the way a person reads one. The top bar used to format its
// own money and stopped when the value moved here; four figures ran together.
func cash(n int) string {
	sign := ""
	if n < 0 {
		sign, n = "-", -n
	}
	out, digits := "", itoa(n)
	for i, c := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			out += ","
		}
		out += string(c)
	}
	return sign + "$" + out
}
