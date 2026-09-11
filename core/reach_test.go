package core

import "testing"

// Does every trade reach past its own income?
//
// The brief has carried a list of which do and which do not, and the list has
// been wrong four times in one night. The burlesque was down as unlinked and
// had been a valid host for putting a night on since the day that was built
// — though the feature was dead, which is why the list being wrong is not
// harmless. The butcher was down as unlinked and hides four units of
// contraband. The poolhall was down as unlinked and takes the seat money from
// the city's own card game into its bankroll. The casino was down as unlinked
// and has a bankroll and a house edge of its own.
//
// A list in a document cannot be trusted about this and neither can a reading
// of the code, because a link can be wired up and still move nothing. So the
// question is asked of the game: for every trade in the city, hold it and do
// not hold it, and measure the number it is supposed to change.
//
// A trade that reaches has something of its own. Contributing cover to
// laundering is not reaching — every trade with a front does that, and if it
// counted then the question would have been finished before it was asked.
//
// What this sweep proves is narrow and worth being exact about: holding the
// place changes the number. Whether the number means anything in play is each
// feature's own guard's job. The burlesque is the reason that line matters —
// its night was wired up, offered on a card and correct to read, and moved not
// one person in the city for as long as it existed. A row here would have
// passed on it the whole time. So where a row can only ask whether something is
// offered, it says so, and the feature carries its own test of whether anybody
// comes.

// reaching is one trade, what holding it changes, and how to measure it.
type reaching struct {
	kind string
	// what the link is, for the failure to say.
	what string
	// better is the direction holding it should move the number.
	better string
	// measure is the number, taken from a world that holds the place and one
	// that does not.
	measure func(w *World, id string) int
}

func addressOf(kind string) string {
	for _, l := range Locations {
		if l.Kind == kind {
			return l.ID
		}
	}
	return ""
}

// standing is a campaign able to exercise any of these: money, a car, stock,
// a quarrel to stand between, and three days of a city behind it.
func holdingIt(t *testing.T, id string, hold bool) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 400000, 90
	w.Player.Car, w.Player.CarWear, w.Player.Dress = 2, 30, 1
	w.Player.Stock = map[string]int{"moonshine": 30}
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "war", 90
	}
	if hold {
		own(w, id)
		if trade, runs := TradeOf(id); runs {
			w.Properties[id].Staff = trade.Hands
			w.Properties[id].Supply = trade.RestockAmount
		}
	}
	w.Player.Location = id
	w.Event = nil
	return w
}

func TestEveryTradeReachesPastItsOwnIncome(t *testing.T) {
	t.Parallel()
	links := []reaching{
		{"garage", "half off what the car costs to keep", "lower",
			func(w *World, id string) int { return w.CarUpkeep() }},
		{"filling", "your own petrol at what it cost the pumps", "lower",
			func(w *World, id string) int {
				w.Player.Fuel, w.Player.Fuelled = 1, w.Minute
				return w.FuelFee(id)
			}},
		{"dealer", "a car without the forecourt's margin", "lower",
			func(w *World, id string) int {
				next, _ := nextVehicle(w.Player.Car)
				cash := w.Player.Cash
				if err := w.BuyVehicle(); err != nil {
					return next.Cost
				}
				return cash - w.Player.Cash
			}},
		{"haulage", "a third off stocking everything else", "lower",
			func(w *World, id string) int { return w.RestockCost(addressOf("butcher")) }},
		{"cabs", "a ride when your own car cannot take you", "higher",
			func(w *World, id string) int {
				w.Player.Car = 0
				if w.RidingWithTheCabs() {
					return 1
				}
				return 0
			}},
		{"scrapyard", "more for the wreck when the yard is yours", "higher",
			func(w *World, id string) int { return w.ScrapWorth(id) }},
		{"laundry", "attention the books can absorb that nothing else can", "higher",
			func(w *World, id string) int {
				// The capacity is a property of the trade whoever holds it, so
				// measuring that measured nothing. What ownership decides is
				// whether the books will take it at all — and they are the only
				// thing in this city that takes attention off you.
				w.Player.Heat = 60
				before := w.Player.Heat
				if w.LaunderReadiness(id) != "" {
					return 0
				}
				if err := w.Launder(id); err != nil {
					return 0
				}
				return before - w.Player.Heat
			}},
		// Offered only. That a night actually moves people is guarded in
		// night_test.go, at five seeds, because it did not for a long time.
		{"burlesque", "a night on is offered, and it draws (see night_test.go)", "higher",
			func(w *World, id string) int {
				if w.NightReadiness(id) != "" {
					return 0
				}
				return 1
			}},
		{"pawn", "the window at what the counter lent, not what it asks", "lower",
			func(w *World, id string) int {
				// shelve gives the thing its own id, so asking for one made up
				// here priced nothing and both sides read zero.
				w.shelve(Shelf{Kind: "dress", Tier: 2, Wear: 30, Ask: 400, Lent: 150})
				return w.WindowPrice(w.Window[len(w.Window)-1].ID)
			}},
		// How well you hear things, which nothing else but coffee and a
		// telephone touches.
		{"saloon", "a line into the city, the size of a telephone in the hall", "higher",
			func(w *World, id string) int { return w.Reach() }},
		// Not a better price either — the exchange already pays over the odds
		// and that is a fact about the floor. What a floor knows is what the
		// numbers mean, and it knows it about the whole market rather than
		// about the room you are standing in.
		{"exchange", "a word on where every price stands against what it is worth", "higher",
			func(w *World, id string) int {
				read := 0
				for _, g := range w.Goods {
					if w.MarketWord(g.ID) != "" {
						read++
					}
				}
				return read
			}},
		// Not a better price — the dock floor's spread is a fact about the
		// floor and true for anybody standing on it. Knowing when a boat is in.
		{"wharf", "a boat in tonight at a price no floor offers", "higher",
			func(w *World, id string) int {
				for day := 0; day < 40; day++ {
					w.Advance(1440)
					w.Event = nil
					if w.BoatIsIn() {
						return w.PriceAt(id, w.Landed.Good) - w.Landed.Price
					}
				}
				return 0
			}},
		// Standing rather than money, which nothing else in this city pays.
		{"club", "a night's standing off a room with your name over the door", "higher",
			func(w *World, id string) int { return w.StandingFromTheDoor(id) }},
		{"butcher", "a cold room things sit in without being looked at", "higher",
			func(w *World, id string) int { return w.Concealed() }},
		{"restaurant", "a dining room two families will sit down in", "higher",
			func(w *World, id string) int {
				if w.SitdownWhere(id) {
					return 1
				}
				return 0
			}},
		{"poolhall", "the seat money off the city's own game", "higher",
			func(w *World, id string) int {
				w.BackRoomNight()
				if prop := w.Properties[id]; prop != nil {
					return prop.Bankroll
				}
				return 0
			}},
		// Whether the night's arithmetic is right is casino_test.go's job; this
		// asks only whether the room has a night at all when it is yours.
		{"casino", "a bankroll the house plays out of, and a night of its own", "higher",
			func(w *World, id string) int {
				prop := w.Properties[id]
				if prop == nil {
					return 0
				}
				prop.Bankroll = 5000
				before := prop.Bankroll
				w.CasinoDay()
				if prop.Bankroll != before {
					return 1
				}
				return 0
			}},
	}
	if len(links) != len(trades) {
		t.Fatalf("the city has %d trades and this asks about %d", len(trades), len(links))
	}
	for _, link := range links {
		id := addressOf(link.kind)
		if id == "" {
			t.Errorf("%s: the city has no such address", link.kind)
			continue
		}
		without := link.measure(holdingIt(t, id, false), id)
		with := link.measure(holdingIt(t, id, true), id)
		moved := with < without
		if link.better == "higher" {
			moved = with > without
		}
		if !moved {
			t.Errorf("%s does not reach: %s. Holding it gives %d where not holding it gives %d, and %s would be %s",
				link.kind, link.what, with, without, "reaching", link.better)
			continue
		}
		t.Logf("%-11s %-52s %6d -> %-6d", link.kind, link.what, without, with)
	}
}
