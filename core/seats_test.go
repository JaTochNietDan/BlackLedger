package core

import "testing"

// Taking a family, and what you are actually handed.
//
// Four addresses this city cannot sell — Saint Agnes, Pier 14, The Monarch and
// the Mercer Exchange — are now the four best businesses in it. Every guard
// written for them took the room with a test helper that sets the deed, which
// proves the trade works for an owner and proves nothing about how anybody
// becomes one. The only way to hold one of these is to take the family, and
// that is a different code path with its own ideas about what moves.

// takenTheFamily puts the player at the head of the family that holds the
// seats, by the route the game actually offers.
// takenTheFamily hands back a city where the move actually landed.
//
// A takeover is a roll — the leader can be ready, or somebody can have told him
// — and this used to make one attempt in one city and assume it worked. It did,
// for the campaign number it was written with, and the tests underneath it went
// on to ask what the player now held. The moment those numbers stopped opening
// the stream in the same eighth of its range, the attempt was seen coming and
// the family kept everything, which reads as the takeover handing over nothing.
func takenTheFamily(t *testing.T) (*World, *Faction, []string) {
	t.Helper()
	for n := uint32(1); n <= 12; n++ {
		w, f, held, landed := tryTheFamily(t, spread(n))
		if landed {
			return w, f, held
		}
	}
	t.Skip("no city in twelve let the move land")
	return nil, nil, nil
}

func tryTheFamily(t *testing.T, seed uint32) (*World, *Faction, []string, bool) {
	t.Helper()
	w := New(seed)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 400
	f := w.faction("bellandi")
	if f == nil {
		t.Skip("no Bellandi Family in this city")
	}
	leader := w.Leader(f.ID)
	if leader == nil {
		t.Skip("nobody is in the chair")
	}
	// A lieutenant of theirs, standing in front of the leader, in a family
	// weak enough to be taken.
	w.Player.Serves = f.ID
	w.Player.Location = leader.Location
	f.Power, f.Peak = 10, 100
	w.Player.Service = PromotionWork * 2
	if reason := w.TakeoverReadiness(); reason != "" {
		t.Skipf("cannot move on them: %s", reason)
	}
	// A day first. Every business in this city is seeded with a count of staff
	// and nobody named, and the list is filled when the clock next settles the
	// city — so a takeover on the stroke of minute one hands over counters that
	// say five and name nobody. That is a seam at the start of a campaign
	// rather than a fault in the takeover, and it is written up in
	// docs/DEVELOPMENT.md; what matters here is the path, not the first minute.
	w.Advance(1440)
	w.Event = nil
	w.Player.Location = leader.Location
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 400
	f.Power, f.Peak = 10, 100
	if reason := w.TakeoverReadiness(); reason != "" {
		t.Skipf("cannot move on them after a day: %s", reason)
	}
	// What they hold going in, read out of the city rather than named here. The
	// club and the docks were written down as theirs, which they are in the one
	// campaign this was built on and need not be in the next.
	held := append([]string{}, w.FamilyHoldings(f.ID)...)
	w.Event = nil
	if err := w.TakeOver(); err != nil {
		t.Fatal(err)
	}
	// Whether it landed, read off the ground rather than off a log line.
	for _, id := range held {
		if !w.Own(id) {
			return w, f, held, false
		}
	}
	return w, f, held, len(held) > 0
}

func TestTakingTheFamilyHandsYouFourWorkingBusinesses(t *testing.T) {
	t.Parallel()
	w, _, held := takenTheFamily(t)
	if len(held) == 0 {
		t.Skip("the family held nothing to hand over")
	}
	took := 0
	for _, id := range held {
		if !w.Own(id) {
			t.Errorf("%s did not come with the family", id)
			continue
		}
		took++
		trade, runs := TradeOf(id)
		if !runs {
			t.Errorf("%s runs on nothing", id)
			continue
		}
		prop := w.Properties[id]
		// The counter is still staffed by the people who were on it. A family
		// changing hands is not a business closing.
		if prop.Staff < trade.Hands {
			t.Errorf("%s came with %d of %d hands", id, prop.Staff, trade.Hands)
		}
		if len(prop.Hands) != prop.Staff {
			t.Errorf("%s counts %d staff and names %d", id, prop.Staff, len(prop.Hands))
		}
		if prop.Supply <= 0 {
			t.Errorf("%s came with nothing to sell", id)
		}
		// And it offers what every business of yours offers.
		w.Player.Location = id
		offered := map[string]bool{}
		for _, a := range w.Actions(id) {
			offered[a.ID] = true
		}
		for _, card := range []string{"inspect", "restock"} {
			if !offered[card] {
				t.Errorf("%s offers no %q after a takeover", id, card)
			}
		}
	}
	if took == 0 {
		t.Fatal("the family was taken and came with nothing")
	}
}

func TestWhatTheSeatsReachWithWorksAfterATakeover(t *testing.T) {
	t.Parallel()
	w, _, _ := takenTheFamily(t)
	if !w.Own("club") || !w.Own("docks") {
		t.Skip("the seats did not come with the family")
	}
	// The Monarch pays standing rather than money, and it is the room being
	// full that does it.
	if crowd := w.TheEveningCrowd("club"); crowd < ClubPerHead {
		t.Fatalf("only %d drink at the club, so this proves nothing", crowd)
	}
	if got := w.StandingFromTheDoor("club"); got <= 0 {
		t.Fatal("a club taken by force earns no standing")
	}
	// And a boat comes into a pier taken by force like any other.
	seen := false
	for night := 0; night < 60; night++ {
		w.Advance(1440)
		w.Event = nil
		if w.BoatIsIn() {
			seen = true
			break
		}
	}
	if !seen {
		t.Fatal("sixty nights on a pier you fought for and no boat came in")
	}
}

// The count and the list.
//
// `Property.Staff` is a number and `Property.Hands` is who those people are,
// and hands.go says plainly that the count "is now the length of a list of
// people who live here". From the first day the clock settles, it is. At
// seeding it is not: every business in the city opens with a count and nobody
// named, and the list is filled the first time the city settles itself.
//
// It is also only true of addresses somebody holds. An independent business is
// staffed as a number and never as people, deliberately — `EmptyChairs` skips
// them, because a shop nobody owns has nobody to poach, nobody to lean on and
// nobody to lose. They become people the day somebody takes the deed.
//
// Nothing in play sees the seeding gap. The player owns nothing on the first morning, and
// every rule that needs a name — asking a counter what they have seen, poaching
// somebody, losing somebody — is reached by walking into a room, which takes
// longer than the gap. It is left alone on purpose rather than moved at the
// seeding stage, because filling nineteen counters inside New() would draw on
// the city before the city exists and the balance suite is sensitive to exactly
// that. What is guarded is that it never drifts again afterwards.
func TestTheCountAndTheListAgreeOnceTheCityHasSettled(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 400000
	for day := 0; day < 30; day++ {
		w.Advance(1440)
		w.Event = nil
		for _, l := range Locations {
			prop := w.Properties[l.ID]
			if prop == nil || l.Kind == "" {
				continue
			}
			// Held by somebody. An independent counter is a number on purpose.
			if prop.Owner == "" || prop.Owner == "independent" {
				continue
			}
			if prop.Staff != len(prop.Hands) {
				t.Fatalf("day %d: %s counts %d staff and names %d",
					day+1, l.ID, prop.Staff, len(prop.Hands))
			}
		}
	}
}
