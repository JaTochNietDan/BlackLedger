package core

import "testing"

// A business's trade — how many hands it needs, what it runs on, what goes
// wrong in it — was keyed by street address. So there could only ever be one
// laundry and one casino in the city: a second would need its own copy of the
// same rules under a different key, and the two would drift apart the first
// time either was touched.
//
// A trade belongs to a KIND of business. An address has a kind, and any number
// of addresses can share one.

func TestEveryEarningAddressDeclaresWhatKindOfBusinessItIs(t *testing.T) {
	for _, l := range Locations {
		earns := PlaceIncome[l.ID] > 0
		_, runs := TradeOf(l.ID)
		if runs && l.Kind == "" {
			t.Errorf("%s has a trade and does not say what kind of business it is", l.ID)
		}
		if l.Kind != "" && !earns {
			t.Errorf("%s is a %s and earns nothing", l.ID, l.Kind)
		}
		if l.Kind != "" {
			if _, known := trades[l.Kind]; !known {
				t.Errorf("%s is a %q, which is not a kind of business the game knows how to run", l.ID, l.Kind)
			}
		}
	}
}

// The point of the change: two addresses of the same kind run by the same
// rules, and are not two copies of them.
func TestTwoAddressesOfOneKindShareTheirTrade(t *testing.T) {
	byKind := map[string][]string{}
	for _, l := range Locations {
		if l.Kind != "" {
			byKind[l.Kind] = append(byKind[l.Kind], l.ID)
		}
	}
	shared := 0
	for kind, ids := range byKind {
		if len(ids) < 2 {
			continue
		}
		shared++
		first, _ := TradeOf(ids[0])
		for _, id := range ids[1:] {
			other, ok := TradeOf(id)
			if !ok {
				t.Errorf("%s is a %s and has no trade", id, kind)
				continue
			}
			if other != first {
				t.Errorf("%s and %s are both %ss and run by different rules", ids[0], id, kind)
			}
		}
	}
	if shared == 0 {
		t.Fatal("no kind of business has two addresses, so nothing here is tested")
	}
	t.Logf("%d kinds of business have more than one address in the city", shared)
}

// And a world made of that city gives every one of them an inside.
func TestEveryAddressOfAKindOpensAsAGoingConcern(t *testing.T) {
	w := New(97)
	for _, l := range Locations {
		trade, runs := TradeOf(l.ID)
		if !runs {
			continue
		}
		prop := w.Properties[l.ID]
		if prop == nil {
			t.Errorf("%s has no premises", l.ID)
			continue
		}
		if prop.Staff != trade.Hands || prop.Supply != trade.RestockAmount {
			t.Errorf("%s opens with %d of %d hands and %d of %d supplies",
				l.ID, prop.Staff, trade.Hands, prop.Supply, trade.RestockAmount)
		}
	}
}

// Anything the game offers "at a laundry" has to be offered at every laundry.
// These were written when the city could only hold one of each, so they name an
// address where they mean a kind of business, and owning the second laundry in
// town would have got the player none of it.
func TestWhatAKindOfBusinessLetsYouDoIsTrueOfEveryOneOfThem(t *testing.T) {
	kindOf := map[string]string{}
	for _, l := range Locations {
		if l.Kind != "" {
			kindOf[l.ID] = l.Kind
		}
	}
	laundries := []string{}
	for id, kind := range kindOf {
		if kind == "laundry" {
			laundries = append(laundries, id)
		}
	}
	if len(laundries) < 2 {
		t.Fatal("the city has fewer than two laundries, so this proves nothing")
	}

	for _, id := range laundries {
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 40000, 200
		w.Properties[id].Owner = "player:1"
		if !ArmourySite(id) {
			t.Errorf("a crate room can be built under one laundry and not under %s", id)
		}
		// Having your own clothes done is a thing a laundry does, not a thing
		// one particular address does.
		w.Player.Location = id
		if !w.PressAtOwnPlace(id) {
			t.Errorf("%s is a laundry the player owns and will not press their clothes", id)
		}
	}
}
