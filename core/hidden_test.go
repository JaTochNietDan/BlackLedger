package core

import "testing"

// "We should probably show options like 'buy kerrigan haulage' before you can
// afford it instead of having it hidden. We probably should just show all
// hidden options tbh. Not sure if hiding them is productive."
//
// An action the player cannot take yet is a thing to want, and a game that
// hides what you cannot afford is a game that never tells you what to save for.
// Every refusal is a sentence saying what would change it, so a refused card is
// worth more than an absent one.

func broke(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 40, 100
	return w
}

func TestWhatYouCannotAffordIsStillOffered(t *testing.T) {
	w := broke(t)
	// Every business that can be bought, from a pocket with $40 in it. A home
	// is not one: it is rented rather than taken over, and it offers its own
	// refusal.
	checked := 0
	for _, l := range Locations {
		if l.Cost <= 0 || w.Own(l.ID) || l.District > w.District || l.Type == "home" {
			continue
		}
		checked++
		w.Player.Location = l.ID
		var buy *Action
		for _, a := range w.Actions(l.ID) {
			if a.ID == "acquire" {
				buy = &a
			}
		}
		if buy == nil {
			t.Fatalf("%s is for sale at $%d and a player with $40 is not shown it at all", l.Name, l.Cost)
		}
		if !buy.Disabled {
			t.Fatalf("%s costs $%d and a player with $40 is offered it as though they could", l.Name, l.Cost)
		}
		if buy.Reason == "" {
			t.Fatalf("buying %s is refused and does not say why", l.Name)
		}
		if buy.Cost <= 0 {
			t.Fatalf("buying %s is offered at $%d", l.Name, buy.Cost)
		}
	}
	if checked == 0 {
		t.Fatal("nothing in this city is for sale, so this proves nothing")
	}
	t.Logf("checked %d businesses a player with $40 cannot afford", checked)
	// And the home a player cannot afford is offered too, refused rather than
	// absent, because wanting somewhere better to live is the same kind of
	// wanting.
	w.Player.Location = "apartment"
	var rent *Action
	for _, a := range w.Actions("apartment") {
		if a.ID == "move_home" {
			rent = &a
		}
	}
	if rent == nil || !rent.Disabled || rent.Reason == "" {
		t.Fatal("a home a player cannot afford is hidden rather than refused")
	}
}

// And every refusal anywhere says what would change it, because a card that
// says "no" and nothing else is worse than no card.
func TestEveryRefusalSaysWhatWouldChangeIt(t *testing.T) {
	w := broke(t)
	silent, checked := []string{}, 0
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location = l.ID
		for _, a := range w.Actions(l.ID) {
			if !a.Disabled {
				continue
			}
			checked++
			if a.Reason == "" {
				silent = append(silent, l.ID+"/"+a.ID)
			}
		}
	}
	t.Logf("checked %d refusals across the city", checked)
	if checked == 0 {
		t.Fatal("a player with $40 was refused nothing anywhere, so this proves nothing")
	}
	if len(silent) > 0 {
		t.Fatalf("%d refusals say nothing about what would change them: %v", len(silent), silent)
	}
}
