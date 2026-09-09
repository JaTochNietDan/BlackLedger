package core

import (
	"strings"
	"testing"
)

// An action that pays its own fee has to declare Cost of nothing, or the
// command layer takes the money a second time. The panel prints the price from
// Cost, so those actions showed no price at all: calling two families to a room
// takes $220 and its button said nothing about money.
//
// Asks is that fee for the panel only. Two things have to hold: it has to be
// there, and declaring it must not charge anybody twice.

func priced(t *testing.T, w *World, place, kind string) Action {
	t.Helper()
	w.Player.Location = place
	for _, a := range w.Actions(place) {
		if a.ID == kind {
			return a
		}
	}
	t.Fatalf("%s is not offered at %s", kind, place)
	return Action{}
}

func TestWorkThatPaysItsOwnFeeStillShowsAPrice(t *testing.T) {
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Player.Contacts = 5
	w.Offshore = 5000
	w.Player.Offshore = false
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	w.ensureOfficials()
	w.Player.Retainers = append(w.Player.Retainers, "editor")
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})

	for _, c := range []struct{ place, kind string; fee int }{
		{"herald", "spike", SpikeCost},
		{"herald", "puff", PuffCost},
		{"market", "offshore_access", AccessCost},
		{"laundry", "still", StillCost},
	} {
		a := priced(t, w, c.place, c.kind)
		if a.Cost != 0 {
			t.Errorf("%s declares a cost of %d, and the command layer would charge it on top of its own fee", c.kind, a.Cost)
		}
		if a.Asks != c.fee {
			t.Errorf("%s takes $%d and the button offers %d", c.kind, c.fee, a.Asks)
		}
	}
}

// The other half, and the reason Cost has to stay nothing: declaring a price
// must not take the money twice.
func TestDeclaringAPriceDoesNotChargeItTwice(t *testing.T) {
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Properties["laundry"].Owner = "player:1"
	a := priced(t, w, "laundry", "still")
	if a.Disabled {
		t.Skipf("a still cannot be built here: %s", a.Reason)
	}
	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "still", Target: "laundry"})
	if err != nil {
		t.Fatalf("building it failed: %v", err)
	}
	spent := before - next.Player.Cash
	// The clock moves, so rent and income land too; the fee must be in there
	// once, not twice.
	// The clock moves while it is built, so a day's income and rent land on top
	// of the fee. What matters is that the fee itself went out once: comfortably
	// more than nothing, and nowhere near twice.
	if spent >= 2*StillCost {
		t.Errorf("a still costs $%d and pressing it took $%d — charged twice", StillCost, spent)
	}
	if spent < StillCost/2 {
		t.Errorf("a still costs $%d and pressing it took only $%d — not charged at all", StillCost, spent)
	}
	t.Logf("declared $%d, took $%d once the hours' other money had moved", a.Asks, spent)
}

// The regression guard, by name rather than by count. A first version asserted
// only that "at least twelve" actions named a fee, and silencing one still left
// twenty-four — a test that passed when I broke the code. These are the actions
// that pay their own way, so each must name a price and must declare no cost,
// or the engine would take the money on top of the fee.
var paysItsOwnWay = []string{
	"sitdown", "spike", "puff", "offshore_access", "still", "armoury",
	"launder", "car", "dress", "arms:weapon", "arms:armour", "charge",
	"bribe", "bankroll", "hire", "restock", "remedy",
	"retain:commissioner", "retain:mayor", "retain:editor",
	"share:", "sign:", "bail:",
	"pact:bellandi", "enquire:bellandi", "smear:bellandi",
}

func TestEveryPricedActionKeepsItsPrice(t *testing.T) {
	w := New(59)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 60000, 200, 100
	w.Player.Contacts = 5
	w.Player.Heat = 10
	w.Offshore = 5000
	w.Player.Offshore = false
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
		w.Properties[id].Staff = 0
		w.Properties[id].Supply = 0
		w.Properties[id].Condition = 60
	}
	w.Incorporate()
	w.ensureOfficials()
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})
	for i := range w.Factions {
		w.Factions[i].Goodwill = 60
	}
	// Somebody of the player's standing in a room, so a share can be paid, and
	// somebody of theirs in a cell, so bail is offered.
	w.NPCs = append(w.NPCs,
		NPC{ID: "ownman", Name: "Otto Reiss", Location: "bar",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Role: "Yours", Trust: 40},
		NPC{ID: "cellman", Name: "Bruno Sala", Location: "precinct",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Role: "Yours",
			Trust: 40, Held: w.Minute + 2*1440})

	seen := map[string]Action{}
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			if _, had := seen[a.ID]; !had {
				seen[a.ID] = a
			}
		}
	}
	missing, checked := []string{}, 0
	for _, kind := range paysItsOwnWay {
		// Three of these carry a person's id after the colon, so they are
		// matched on the prefix and the first one found stands for the rest.
		a, offered := seen[kind]
		if !offered && strings.HasSuffix(kind, ":") {
			for id, found := range seen {
				if strings.HasPrefix(id, kind) {
					a, offered = found, true
					break
				}
			}
		}
		if !offered {
			missing = append(missing, kind)
			continue
		}
		checked++
		if a.Asks <= 0 {
			t.Errorf("%s pays its own fee and names no price", kind)
		}
		if a.Cost != 0 {
			t.Errorf("%s pays its own fee and also declares a cost of %d, so the money would go out twice",
				kind, a.Cost)
		}
	}
	if checked < len(paysItsOwnWay)-4 {
		t.Errorf("only %d of %d priced actions were offered anywhere; this city is not exercising them: %v",
			checked, len(paysItsOwnWay), missing)
	}
	if len(missing) > 0 {
		t.Logf("not offered in this city, so not checked: %v", missing)
	}
	t.Logf("checked %d actions that pay their own way", checked)
}
