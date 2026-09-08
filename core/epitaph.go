package core

import "fmt"

// The death screen is the most dramatic moment this game has, and it was
// telling the player something that stopped being true. "Your properties pass
// to your former organization. Another life begins without your money, rank, or
// authority." — but what a protagonist leaves has passed to the strongest of
// their own people for a long time now, as an organization under a name they
// never chose, which the next protagonist can deal with or fight. The screen
// had rotted exactly the way the Guide had.
//
// So it reports rather than describes. Everything here is read out of what the
// city actually recorded when they died.

// Epitaph is what a life amounted to and what became of it.
func (w *World) Epitaph() map[string]any {
	if len(w.Dead) == 0 {
		return nil
	}
	last := w.Dead[len(w.Dead)-1]
	p := &w.Player

	// What still stands with their name on it, which is different from what
	// somebody inherited.
	standing := []string{}
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Owner == "former:"+last.Name {
			standing = append(standing, l.Name)
		}
	}

	// What became of what they built, in a sentence rather than a rule.
	became := "They had nobody. Whatever they held stands with no one to answer for it."
	if last.Estate != "" {
		became = fmt.Sprintf("%s took it over. It holds their premises, keeps their quarrels, and is worth rather less than it was, because most of what it was worth was them.", last.Estate)
	}

	// What the next person actually inherits, which is nothing but the city and
	// whatever was put somewhere the city cannot reach.
	inherits := "Nothing but the city, and whatever it remembers."
	if w.Offshore > 0 {
		inherits = fmt.Sprintf("Nothing here. There is $%d outside the city that somebody could prove is theirs, if they could afford to.", w.Offshore)
	}

	return map[string]any{
		"name": last.Name, "cause": last.Cause,
		"day":       last.Minute/1440 + 1,
		"life":      last.Life,
		"respect":   p.Respect,
		"earned":    p.Earned,
		"estate":    last.Estate,
		"became":    became,
		"standing":  standing,
		"inherits":  inherits,
		"headlines": w.headlinesFor(last),
	}
}

// headlinesFor is what the paper carried about this person, which is the only
// obituary this city writes.
func (w *World) headlinesFor(d Death) []string {
	out := []string{}
	for i := len(w.News) - 1; i >= 0 && len(out) < 4; i-- {
		if s := w.News[i]; s.Life == d.Life && scrutinyWeight[s.Kind] > 0 {
			out = append(out, s.Headline)
		}
	}
	return out
}
