package core

import "testing"

// No two cards in a room share a name. A command is matched to the card by its
// id, and the match takes the last one that fits — so two cards under one name
// means the wrong card's price, minutes and label are used, whatever the right
// one was. Asking two different people where two different marks are was
// offered as "about:leo" twice: the mark rode on the choice so the right thing
// happened, and the receipt named the wrong person, and the day one of them
// costs more than another the wrong one is charged.
//
// This is the shape that put somebody who asked for the machines into a hand of
// cards, and it is cheap to rule out everywhere at once.
func TestNoTwoCardsInARoomShareAName(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 50000
	// Something held, so the work a proprietor is offered is in the list too.
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["poolhall"].Owner = "player:1"

	read := 0
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location, w.Event = l.ID, nil
		seen := map[string]string{}
		for _, a := range w.Actions(l.ID) {
			read++
			if first, ok := seen[a.ID]; ok {
				t.Errorf("%s offers %q twice: %q and %q — a command matched to it takes the second",
					l.ID, a.ID, first, a.Label)
			}
			seen[a.ID] = a.Label
		}
	}
	if read < 200 {
		t.Fatalf("only %d cards read across the city, so this measures nothing", read)
	}
	t.Logf("%d cards read across the city, every one of them named once in its room", read)
}
