package core

import (
	"strings"
	"testing"
)

// What a card says it costs, against what pressing it takes.
//
// Every action in this game runs `w.Pay(a.Cost)` in the command layer, and one
// that pays its own fee has to declare `Cost: 0` and use `asks(...)` so the
// panel still shows a price. Get that wrong and the money goes out twice. There
// is a guard for it and it checks exactly one action — a still at a laundry —
// chosen by hand, and the last time this project trusted a hand-written list of
// actions it cost five separate measures before anybody noticed what was not on
// it.
//
// So: every action, in every room, that names a price at all. Press it in a copy
// of the world and see what leaves the pocket. The clock moves while the work is
// done, so rent, wages and a day's income land on top of the fee and the reading
// is never exact; what it can say for certain is that nobody was charged twice.
// Both directions are checked. Across every priced card in the city not one
// came in under half what it said, so "names a price and does not take it" is a
// fault rather than a tolerance — and neither half needs a list of which actions
// pay out, which is the point.
func TestNoCardTakesItsPriceTwice(t *testing.T) {
	t.Parallel()
	base := New(53)
	base.Event, base.District = nil, 9
	base.Player.Cash, base.Player.Respect, base.Player.Health = 400000, 200, 100
	base.Player.Contacts = 5
	base.Player.Car, base.Player.CarWear = 2, 20
	for _, id := range []string{"laundry", "garage", "casino", "poolhall"} {
		if prop := base.Properties[id]; prop != nil {
			prop.Owner = "player:1"
		}
	}
	base.Properties["garage"].Trouble = true

	pressed, wrong := 0, 0
	for _, l := range Locations {
		if l.District > base.District {
			continue
		}
		w := base.Clone()
		w.Player.Location, w.Event = l.ID, nil
		for _, a := range w.Actions(l.ID) {
			price := a.Cost
			if a.Asks > price {
				price = a.Asks
			}
			if price <= 0 || a.Disabled {
				continue
			}
			// A card that opens a conversation rather than doing a thing is
			// answered somewhere else, and what it takes is taken there.
			try := w.Clone()
			try.Player.Location, try.Event = l.ID, nil
			before := try.Player.Cash
			next, err := Execute(try, Command{RequestID: ID(), Revision: try.Revision,
				Kind: a.ID, Target: l.ID})
			if err != nil {
				continue
			}
			pressed++
			// The clock's own money, measured rather than allowed for. Rent,
			// wages and a day's takings all move while the work is done, and
			// they can hide a second charge: a still is $450, a day at the
			// laundry brings in enough that being charged twice read as $859,
			// which is under twice $450. Advancing a copy by the same minutes
			// with nothing pressed says exactly how much of the difference was
			// the clock.
			clock := w.Clone()
			clock.Event = nil
			was := clock.Player.Cash
			clock.Advance(a.Minutes)
			passing := was - clock.Player.Cash

			// Both directions, because with the clock measured rather than
			// allowed for the reading is close enough to say either. Across
			// every priced card in the city not one came in under half what it
			// said, so a card naming a price it does not take is a fault and
			// not a tolerance.
			// A few dollars of slack, because the control and the run are two
			// campaigns and a day's income can land in one and not the other:
			// a thirty-minute $25 card read as $52 that way. Half the price,
			// capped at twenty, is narrower than the fault this looks for.
			slack := min(20, price/2)
			spent := before - next.Player.Cash - passing
			switch {
			case spent >= 2*price+slack:
				wrong++
				t.Errorf("%s at %s says $%d and pressing it took $%d once the "+
					"clock's own $%d is taken out — charged twice",
					a.ID, l.ID, price, spent, passing)
			case spent < price/2-slack:
				wrong++
				t.Errorf("%s at %s says $%d and pressing it took $%d once the "+
					"clock's own $%d is taken out — the price is not taken",
					a.ID, l.ID, price, spent, passing)
			}
		}
	}
	t.Logf("%d priced cards pressed, %d took something other than what they said", pressed, wrong)
	if pressed < 40 {
		t.Fatalf("only %d priced cards were pressed, so this measures nothing", pressed)
	}
}

// And the reader is worth no more than what it can see. A card whose price the
// engine takes on top of its own fee is the fault this looks for, so one is
// built here and the same reading is taken of it.
func TestPressingACardThatChargesTwiceIsSeen(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Location = "laundry"

	var still Action
	for _, a := range w.Actions("laundry") {
		if a.ID == "still" {
			still = a
		}
	}
	if still.ID == "" || still.Disabled {
		t.Skipf("a still cannot be built here: %s", still.Reason)
	}
	price := still.Cost
	if still.Asks > price {
		price = still.Asks
	}
	if price <= 0 {
		t.Fatal("a still names no price, so there is nothing to read")
	}

	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "still", Target: "laundry"})
	if err != nil {
		t.Fatal(err)
	}
	clock := w.Clone()
	clock.Event = nil
	was := clock.Player.Cash
	clock.Advance(still.Minutes)
	passing := was - clock.Player.Cash

	once := before - next.Player.Cash - passing
	if once >= 2*price {
		t.Fatalf("a still is charged twice already: $%d for a $%d room", once, price)
	}
	// Charged a second time by hand, which is what the fault looks like.
	twice := next.Clone()
	if err := twice.Pay(price); err != nil {
		t.Fatal(err)
	}
	if spent := before - twice.Player.Cash - passing; spent < 2*price {
		t.Fatalf("taking the fee a second time reads as $%d, under twice $%d, "+
			"so the sweep above cannot see a double charge", spent, price)
	}
	if !strings.Contains(strings.ToLower(still.Detail), "$") {
		t.Errorf("a card that costs money says nothing about money: %q", still.Detail)
	}
}

// And the other number on every card: how long it takes.
//
// A price is one of two things a card promises before it is pressed. The other
// is the hour, and in this game the hour is the real currency — the clock is
// what brings rent, wages, a rival's move and the police, and a card that says
// thirty minutes and spends a day is a worse lie than one that overcharges.
// Nothing had ever read it.
//
// Same shape as the sweep above and for the same reason: every card in every
// room, pressed in a copy of the world, with no list of which ones are special.
func TestNoCardSpendsMoreOfTheDayThanItSays(t *testing.T) {
	t.Parallel()
	base := New(53)
	base.Event, base.District = nil, 9
	base.Player.Cash, base.Player.Respect, base.Player.Health = 400000, 200, 100
	base.Player.Contacts = 5
	base.Player.Car, base.Player.CarWear = 2, 20
	for _, id := range []string{"laundry", "garage", "casino", "poolhall"} {
		if prop := base.Properties[id]; prop != nil {
			prop.Owner = "player:1"
		}
	}
	base.Properties["garage"].Trouble = true

	pressed, wrong := 0, 0
	for _, l := range Locations {
		if l.District > base.District {
			continue
		}
		w := base.Clone()
		w.Player.Location, w.Event = l.ID, nil
		for _, a := range w.Actions(l.ID) {
			if a.Disabled {
				continue
			}
			try := w.Clone()
			try.Player.Location, try.Event = l.ID, nil
			was := try.Minute
			next, err := Execute(try, Command{RequestID: ID(), Revision: try.Revision,
				Kind: a.ID, Target: l.ID})
			if err != nil {
				continue
			}
			pressed++
			// What the card declares, which is `Minutes` for work done inside
			// the day and `Away` for work that is days out of the city. The
			// second is the same trick `Asks` plays with money: the trips run
			// their own days, so they carry `Minutes: 0` and the panel prints
			// "4 days away" out of `Away` instead. Reading only `Minutes` made
			// all three trips look like a card that says nothing and spends
			// four days, which is the fault this looks for rather than the
			// convention that avoids it.
			declared := a.Minutes
			if a.Away > declared {
				declared = a.Away
			}
			spent := next.Minute - was
			// A card may cost less than it says — being turned away at a door
			// takes the walk and not the evening — but it must never quietly
			// take more of the day than it named. A card that names no time at
			// all is a card that says it takes none, and the same rule holds,
			// which is what makes this the guard for `Away` as well.
			if spent > declared {
				wrong++
				t.Errorf("%s at %s says %d minutes and took %d",
					a.ID, l.ID, declared, spent)
			}
		}
	}
	t.Logf("%d cards pressed, %d took more of the day than they said", pressed, wrong)
	if pressed < 120 {
		t.Fatalf("only %d cards were pressed, so this measures nothing", pressed)
	}
}

// A card the game offers must work when it is pressed.
//
// Every action in this game answers `Disabled` and a `Reason` before anybody
// touches it, and that is the whole contract of the panel: what you can do reads
// live, what you cannot reads dim with the reason on it. A live card that
// returns an error is the contract broken — the player gets a toast instead of a
// reason, which the bail card did for a long time and which was found by
// auditing one action by hand.
//
// A card that takes a figure the player types is pressed with a figure inside
// the bounds the core itself published. Sending nothing to a card that wants a
// number is a fault in the test rather than in the game: the first run of this
// refused six cards, all of them a wage or a house limit, all of them because it
// typed a zero.
//
// Seven situations rather than one rich player, because the interesting cards
// are the ones a city only offers to somebody broke, hurt, wanted, nameless or
// thirty days in, and one comfortable world never sees them.
func TestEveryLiveCardWorksWhenPressed(t *testing.T) {
	t.Parallel()
	situations := []struct {
		name string
		fit  func(w *World)
	}{
		{"comfortable", func(w *World) {}},
		{"broke", func(w *World) { w.Player.Cash = 3 }},
		{"hurt", func(w *World) { w.Player.Health = 12 }},
		{"wanted", func(w *World) { w.Player.Heat = 95 }},
		{"nameless", func(w *World) { w.Player.Respect = 0 }},
		{"with somebody", func(w *World) {
			if driver := w.Holder("driver"); driver != nil {
				w.Player.Crew = append(w.Player.Crew,
					Crew{ID: driver.ID, Name: driver.Name, Loyalty: 70})
			}
		}},
		{"a month in", func(w *World) { w.Advance(1440 * 30) }},
	}

	// Fifteen seconds of wall clock at six seeds a situation. The fast half of
	// the gate takes one seed each, which still reaches every situation.
	seeds, floor := uint32(6), 12000
	if testing.Short() {
		seeds, floor = 1, 2000
	}
	live, refused := 0, 0
	for _, situation := range situations {
		for seed := uint32(1); seed <= seeds; seed++ {
			base := New(spread(seed))
			base.Event, base.District = nil, 9
			base.Player.Cash, base.Player.Respect, base.Player.Health = 400000, 200, 100
			base.Player.Contacts = 5
			base.Player.Car, base.Player.CarWear = 2, 20
			for _, id := range []string{"laundry", "garage", "casino", "poolhall"} {
				if prop := base.Properties[id]; prop != nil {
					prop.Owner = "player:1"
				}
			}
			base.Properties["garage"].Trouble = true
			situation.fit(base)

			for _, l := range Locations {
				if l.District > base.District {
					continue
				}
				w := base.Clone()
				w.Player.Location, w.Event = l.ID, nil
				for _, a := range w.Actions(l.ID) {
					if a.Disabled {
						continue
					}
					live++
					amount := 0
					if a.Sum != nil {
						amount = a.Sum.Preset
						if amount < a.Sum.Least || amount > a.Sum.Most {
							amount = a.Sum.Least
						}
					}
					try := w.Clone()
					try.Player.Location, try.Event = l.ID, nil
					if _, err := Execute(try, Command{RequestID: ID(), Revision: try.Revision,
						Kind: a.ID, Target: l.ID, Amount: amount}); err != nil {
						refused++
						t.Errorf("%s at %s reads live to a %s player and refuses when pressed: %v",
							a.ID, l.ID, situation.name, err)
					}
				}
			}
		}
	}
	t.Logf("%d live cards pressed across seven situations, %d refused", live, refused)
	if live < floor {
		t.Fatalf("only %d live cards were pressed, so this measures nothing", live)
	}
}

// "It is the cheapest respect in this city and the only kind nobody had to be
// hurt for."
//
// That is written above `Puff` in core/newsroom.go, and it is the game making a
// claim about itself across every other card it offers. Nothing could check a
// sentence like that: the prose sweeps read whether a line is well formed, and
// the price and clock sweeps read one card against its own declaration. This
// reads one card against all the others.
//
// Same machinery as the sweeps above — press everything in a copy of the world,
// with the clock's own money measured rather than allowed for — and then divide.
// A card that raises standing and takes money has a price a point, and the paper
// has to be the lowest of them or the sentence comes out.
// harmed counts what this city is carrying that somebody will pay for: the dead,
// the injured, and the quarrels. A card can then be asked whether anybody was
// the worse for it without anybody writing down which cards those are.
//
// The quarrels are the half this did not have at first, and leaving them out
// certified a sentence that is not true. Putting a word in the right ear costs
// $25 and buys a point of standing — cheaper a point than the paper — and
// nobody is bleeding when the card resolves. What it does is leave two families
// believing the other moved on them, which is a fight somebody else has later.
// Reading only the dead and the hurt at the instant of the press calls that a
// clean purchase.
func harmed(w *World) int {
	n := 0
	for i := range w.NPCs {
		if who := &w.NPCs[i]; who.Dead || who.Hurt {
			n++
		}
	}
	if w.Player.Health < 100 {
		n++
	}
	for _, c := range w.Conflicts {
		n += c.Hostility
	}
	return n
}

func TestThePaperIsTheCheapestStandingInTheCity(t *testing.T) {
	t.Parallel()
	type buy struct {
		id    string
		spent int
		gain  int
	}
	dearer := func(a, b buy) bool { return a.spent*b.gain > b.spent*a.gain }

	best := map[string]buy{}
	for seed := uint32(1); seed <= 3; seed++ {
		base := New(spread(seed))
		base.Event, base.District = nil, 9
		base.Player.Cash, base.Player.Respect, base.Player.Health = 400000, 200, 100
		base.Player.Contacts = 5
		base.Player.Car, base.Player.CarWear = 2, 20
		// Somebody at the paper takes the player's calls, or the card this is
		// about is never offered and the sweep measures the rest of the city
		// against nothing.
		base.ensureOfficials()
		base.Player.Retainers = append(base.Player.Retainers, "editor")
		for _, id := range []string{"laundry", "garage", "casino", "poolhall"} {
			if p := base.Properties[id]; p != nil {
				p.Owner = "player:1"
			}
		}
		for _, l := range Locations {
			if l.District > base.District {
				continue
			}
			w := base.Clone()
			w.Player.Location, w.Event = l.ID, nil
			for _, a := range w.Actions(l.ID) {
				if a.Disabled {
					continue
				}
				try := w.Clone()
				try.Player.Location, try.Event = l.ID, nil
				cash, standing := try.Player.Cash, try.Player.Respect
				next, err := Execute(try, Command{RequestID: ID(), Revision: try.Revision,
					Kind: a.ID, Target: l.ID})
				if err != nil {
					continue
				}
				clock := w.Clone()
				clock.Event = nil
				was := clock.Player.Cash
				clock.Advance(a.Minutes)
				passing := was - clock.Player.Cash

				// And nobody the worse for it. The sentence has two halves —
				// cheapest, and the only kind nobody was hurt for — so a card
				// that leaves somebody dead or hurt is not in the comparison at
				// all. Read off the city rather than off a list of which cards
				// are violent: going after somebody yourself came out at a
				// dollar for nine points of standing, which is true and is not
				// what the paper is being compared against.
				if harmed(try) != harmed(next) {
					continue
				}
				got := buy{a.ID, cash - next.Player.Cash - passing, next.Player.Respect - standing}
				if got.gain <= 0 || got.spent <= 0 {
					continue
				}
				if had, seen := best[got.id]; !seen || dearer(had, got) {
					best[got.id] = got
				}
			}
		}
	}

	paper, ok := best["puff"]
	if !ok {
		t.Fatal("the paper never sold a paragraph, so this measures nothing")
	}
	if len(best) < 3 {
		t.Fatalf("only %d cards in this city buy standing for money, so this compares almost nothing", len(best))
	}
	for id, other := range best {
		if id == "puff" {
			continue
		}
		t.Logf("%-12s $%d for %d", id, other.spent, other.gain)
		if dearer(paper, other) {
			t.Errorf("%s buys standing at $%d for %d and the paper wants $%d for %d, "+
				"so the paper is not the cheapest respect in this city",
				id, other.spent, other.gain, paper.spent, paper.gain)
		}
	}
	t.Logf("paper $%d for %d across %d cards that sell standing", paper.spent, paper.gain, len(best))
}

// No two cards in a room may share an id.
//
// A command is matched to a card by its id, so two cards with one id is a
// receipt that names the wrong person — and the moment one of them costs more
// than the other, the wrong one is charged. The comment above the asking card
// in world.go says exactly that and says it was fixed; nothing has ever checked
// it, and the fault it describes is invisible from inside the core because both
// cards do the right thing.
//
// Written after a probe of a running server published three cards all called
// `about:leo`, with the person being asked about riding on the choice field.
// A freshly built server publishes `about:leo:vittorio` and the two others,
// which is what the source says — so what was seen was a stale binary rather
// than a fault in the game. This is the guard that would have said which, in
// one line, instead of an hour.
func TestNoTwoCardsInARoomShareAnID(t *testing.T) {
	t.Parallel()
	rooms, cards := 0, 0
	for seed := uint32(1); seed <= 6; seed++ {
		w := New(spread(seed))
		w.Event, w.District = nil, 9
		w.Player.Cash, w.Player.Respect, w.Player.Health = 400000, 300, 100
		for _, id := range []string{"laundry", "garage", "casino", "poolhall"} {
			if p := w.Properties[id]; p != nil {
				p.Owner = "player:1"
			}
		}
		// A city that has been running, because the cards that name two people
		// need somebody to owe somebody and somebody to have gone missing.
		w.Advance(1440 * 10)
		for _, l := range Locations {
			if l.District > w.District {
				continue
			}
			w.Player.Location, w.Event = l.ID, nil
			seen := map[string]Action{}
			rooms++
			for _, a := range w.Actions(l.ID) {
				cards++
				if was, twice := seen[a.ID]; twice {
					t.Errorf("%s offers %q twice: %q and %q",
						l.ID, a.ID, was.Label, a.Label)
					continue
				}
				seen[a.ID] = a
			}
		}
	}
	t.Logf("%d cards read across %d rooms", cards, rooms)
	if cards < 2000 {
		t.Fatalf("only %d cards were read, so this measures nothing", cards)
	}
}
