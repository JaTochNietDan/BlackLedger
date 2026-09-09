package core

import "fmt"

// What the street says about how the families are doing.
//
// The simulation moves families constantly — they take each other's premises,
// lose strength when a holding is wrecked, recover it when it is repaired, go
// to war and come out of it smaller. Measured over a long campaign a war starts
// in most runs and organizations are destroyed in a quarter of them.
//
// None of which the player could feel. It was all true and all invisible: to
// know that the Russo Outfit had been losing for a fortnight you had to open a
// screen and compare a number to a number you had not written down. A city that
// is busy behind glass is not a living city.
//
// So the paper says so, the way a paper would: not "power 41" but that people
// who deal with them have noticed. It reports a move only when the move is big
// enough to be worth a paragraph, and it remembers what it last said so it does
// not report the same slide every morning for a week.

// FortuneShift is how far a family's strength has to move before anybody
// outside it would notice. Small enough that a bad week shows; large enough
// that the ordinary rise and fall of a quiet month does not fill the paper.
const FortuneShift = 12

// FortunesDay files what the street has noticed about the families. Called from
// the clock, so it lands in the morning's edition with everything else.
func (w *World) FortunesDay() {
	for i := range w.Factions {
		f := &w.Factions[i]
		// The player's own organization is not news to the player.
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		if f.Reported == 0 {
			f.Reported = f.Power
			continue
		}
		move := f.Power - f.Reported
		if move > -FortuneShift && move < FortuneShift {
			continue
		}
		f.Reported = f.Power
		if move > 0 {
			w.Report("civic", upper(f.Name)+" GAINING GROUND", fmt.Sprintf(
				"%s is spoken of more confidently than it was. People who deal with them say there is more of them to deal with, and those who owe them are said to be paying on time.",
				f.Name))
			continue
		}
		// A family that has lost a lot is the more interesting story, and the
		// paper says it the way a paper does: carefully, and about other people.
		holdings := len(w.FamilyHoldings(f.ID))
		lost := "Nobody at the family would say what has gone wrong."
		if holdings == 0 {
			lost = "They are not known to be holding anything at present, which for an organization of that name is new."
		} else if holdings == 1 {
			lost = "One premises is still spoken of as theirs. It was more than that."
		}
		w.Report("civic", upper(f.Name)+" SAID TO BE STRUGGLING", fmt.Sprintf(
			"%s has had a poor few weeks by the reckoning of people who watch such things. %s", f.Name, lost))
	}
}
