package core

import "fmt"

// Pier 14, and the boat that comes in.
//
// The waterfront is where it comes ashore — the city's own price less a fifth
// for moonshine, less a sixth for guns — and that spread is a fact about the
// floor, true for anybody standing on it. Holding the pier changed nothing at
// all, which is the same fault The Monarch had: an address you fight a war for
// that pays a number and offers no decision.
//
// A pier of your own is not a better price. It is knowing when a boat is in.
// Somebody who holds the wharf hears that a quantity is coming ashore tonight
// and wants to be gone by morning, at a price no floor in this city offers —
// and then has to decide how much of it to carry, which is the decision the
// whole underground trade is built on and the one thing a discount can never
// be.

const (
	// LandingUnder is the price of a landing as a percentage of what the same
	// good costs on the dock floor. The floor is already cheap; this is the
	// price of somebody who has to be at sea before it is light.
	LandingUnder = 62
	// LandingLot is roughly how much comes off a boat.
	LandingLot = 22
	// LandingChance is the chance in a hundred that a boat is in on an
	// ordinary night, at a pier that is working.
	LandingChance = 30
	// LandingLasts is how long it is there. One night. That is the whole
	// point of it.
	LandingLasts = 1440
	// LandingMinutes is how long it takes to get it off the boat.
	LandingMinutes = 60
)

// Landing is a quantity ashore tonight at a price, and when it goes.
type Landing struct {
	Good  string `json:"good"`
	Units int    `json:"units"`
	Price int    `json:"price"`
	Until int    `json:"until"`
	Where string `json:"where"`
}

// Wharf reports whether this address is a waterfront a boat can come into.
func Wharf(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == "wharf"
}

// BoatIsIn reports whether there is something ashore the player can still take.
func (w *World) BoatIsIn() bool {
	return w.Landed.Units > 0 && w.Minute < w.Landed.Until
}

// LandingNight is the boat coming in, once a day from the clock.
//
// Derived from the seed and the day rather than drawn from either random
// stream, for the reason written on the pawnbroker's window: a daily feature
// that spends the city's randomness moves everything downstream of it, and two
// runs of the same seed should land the same boat.
func (w *World) LandingNight() {
	for _, l := range Locations {
		if !Wharf(l.ID) || !w.Own(l.ID) {
			continue
		}
		prop := w.Properties[l.ID]
		trade, runs := TradeOf(l.ID)
		// A pier with nobody working it and a pier with something wrong on it
		// are both piers nobody unloads at in the dark.
		if prop == nil || !runs || prop.Trouble || prop.Staff < trade.Hands {
			return
		}
		day := w.Minute / 1440
		if roll(w.Seed, day, 11)%100 >= LandingChance {
			return
		}
		landed := []string{}
		for _, g := range w.Goods {
			if TradesAt(l.ID, g.ID) {
				landed = append(landed, g.ID)
			}
		}
		if len(landed) == 0 {
			return
		}
		good := landed[roll(w.Seed, day, 12)%len(landed)]
		units := LandingLot/2 + roll(w.Seed, day, 13)%LandingLot
		price := max(1, w.PriceAt(l.ID, good)*LandingUnder/100)
		w.Landed = Landing{Good: good, Units: units, Price: price, Until: w.Minute + LandingLasts, Where: l.ID}
		g := w.Good(good)
		w.Log("A boat is in at "+l.Name,
			fmt.Sprintf("%d %ss of %s at $%d, against $%d on your own floor. They want to be gone before it is light.",
				units, g.Unit, g.InBulk(), price, w.PriceAt(l.ID, good)), "business")
		return
	}
}

// LandingReadiness explains why the boat cannot be unloaded, or returns "".
func (w *World) LandingReadiness(units int) string {
	if !w.BoatIsIn() {
		return "There is nothing ashore tonight"
	}
	if w.Player.Location != w.Landed.Where {
		return "This is done on the wharf"
	}
	if units <= 0 || units > w.Landed.Units {
		units = w.Landed.Units
	}
	if w.Player.Cash < w.Landed.Price*units {
		return fmt.Sprintf("The lot is $%d and one crate is $%d", w.Landed.Price*w.Landed.Units, w.Landed.Price)
	}
	return ""
}

// TakeTheLanding buys what came ashore, or as much of it as was asked for.
func (w *World) TakeTheLanding(units int) error {
	if reason := w.LandingReadiness(units); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if units <= 0 || units > w.Landed.Units {
		units = w.Landed.Units
	}
	cost := w.Landed.Price * units
	if err := w.Pay(cost); err != nil {
		return err
	}
	if w.Player.Stock == nil {
		w.Player.Stock = map[string]int{}
	}
	good := w.Landed.Good
	w.Player.Stock[good] += units
	w.Landed.Units -= units
	g := w.Good(good)
	// What is carried is what draws attention, which is the decision this is
	// for. A cheap lot you cannot hide is not a cheap lot.
	w.Log("Off the boat",
		fmt.Sprintf("%s of %s for $%d at $%d each. You are holding %s, %d of them where anybody can find them.",
			counted(units, g.Unit, g.Unit+"s"), g.InBulk(), cost, w.Landed.Price,
			counted(w.Carrying(), "unit", "units"), w.Exposed()), "business")
	return nil
}
