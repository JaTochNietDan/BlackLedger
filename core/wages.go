package core

import "fmt"

// What you pay them.
//
// The people behind your counters can walk out over something you did, and
// there was nothing to be done about it except not do it. A wage is the oldest
// answer to that, and the one thing a business owner actually decides every
// week: pay over the rate and they think better of you, pay under it and they
// think less, and the books say what it costs either way.
//
// The trade's own wage is the rate — what the work is worth in this city — and
// what a place pays is a decision on top of it.

const (
	// WageFloor is the least anybody works for, as a share of the rate. Nobody
	// stands behind a counter for nothing.
	WageFloor = .5
	// WageCeiling is the most a business will promise, as a multiple of the
	// rate. Past it you are not paying a wage, you are giving money away, and
	// the trade cannot carry it.
	WageCeiling = 3
	// Generous and Mean are what a day at either end is worth to somebody's
	// opinion of you. Small: it is a wage, not a gift, and it takes weeks.
	Generous = 1
	Mean     = 1
)

// WageAt is what this business pays a hand a day.
func (w *World) WageAt(id string) int {
	trade, ok := TradeOf(id)
	if !ok {
		return 0
	}
	if prop := w.Properties[id]; prop != nil && prop.Wage > 0 {
		return prop.Wage
	}
	return trade.Wage
}

// WageBounds is what a business may pay, from the least anybody stands there
// for to the most the trade can carry.
func WageBounds(rate int) (int, int) {
	return max(1, int(float64(rate)*WageFloor)), rate * WageCeiling
}

// PayReadiness explains why a figure cannot be paid, or returns "".
func (w *World) PayReadiness(id string, amount int) string {
	trade, ok := TradeOf(id)
	if !ok || !w.Own(id) {
		return "This is not a business of yours"
	}
	least, most := WageBounds(trade.Wage)
	if amount < least {
		return fmt.Sprintf("Nobody stands behind a counter for less than $%d a day", least)
	}
	if amount > most {
		return fmt.Sprintf("The trade will not carry more than $%d a day a head", most)
	}
	return ""
}

// SetWage is the decision. It changes nothing today: what it buys is what the
// people behind the counter make of you over the weeks they are paid it.
func (w *World) SetWage(id string, amount int) error {
	if reason := w.PayReadiness(id, amount); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	trade, _ := TradeOf(id)
	prop := w.Properties[id]
	was := w.WageAt(id)
	prop.Wage = amount
	place, _ := PlaceByID(id)
	how := "which is the rate"
	switch {
	case amount > trade.Wage:
		how = fmt.Sprintf("which is $%d over the rate", amount-trade.Wage)
	case amount < trade.Wage:
		how = fmt.Sprintf("which is $%d under it", trade.Wage-amount)
	}
	w.Log("The wage at "+place.Name,
		fmt.Sprintf("$%d a day a head, %s. It was $%d. Nobody says anything about it today.",
			amount, how, was), "business")
	return nil
}

// PayDay is what a week of being paid over or under the rate does to what the
// people behind a counter think of the person paying it. Called every business
// day, which is where the wage bill is already settled.
func (w *World) PayDay() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, ok := TradeOf(l.ID)
		if prop == nil || !ok || !w.Own(l.ID) || len(prop.Hands) == 0 {
			continue
		}
		paid := w.WageAt(l.ID)
		for _, who := range prop.Hands {
			n := w.NPC(who)
			if n == nil || n.Dead {
				continue
			}
			switch {
			case paid > trade.Wage:
				n.Trust = min(100, n.Trust+Generous)
			case paid < trade.Wage:
				n.Trust = max(0, n.Trust-Mean)
			}
		}
	}
}

// Unpaid is what a day nobody covered costs the regard of the people who were
// not paid. Larger than a day over the rate is worth, because not being paid is
// not the same kind of thing as being paid a little less.
const Unpaid = 4

const (
	// PatienceRunsOut is how many unpaid nights somebody stands behind a
	// counter for before it stops being a job. A week: long enough that a bad
	// stretch is not the end of anybody's employment, short enough that a
	// player who has stopped paying finds out by losing people.
	PatienceRunsOut = 7
	// NotBeingPaid is the daily chance one of them does not come in once that
	// week has passed. Higher than a rival's offer, because there is no money
	// in this at all, and still a decision somebody makes on a day of their
	// own choosing rather than a number crossing a line.
	NotBeingPaid = .2
)

// NobodyGotPaid takes the day off what the people behind the player's counters
// think of them, counts the night against each place, and reports how many
// hands went unpaid. Called from the night the bills do not clear.
func (w *World) NobodyGotPaid() int {
	short := 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) || len(prop.Hands) == 0 {
			continue
		}
		prop.Unpaid++
		for _, who := range prop.Hands {
			n := w.NPC(who)
			if n == nil || n.Dead {
				continue
			}
			n.Trust = max(0, n.Trust-Unpaid)
			short++
		}
		// A week of this is a reason to stop coming in, and Notice reads the
		// count for exactly that. It was unreachable while the day's bill was
		// all or nothing: one short night stripped the security and the
		// address, the cost fell from $2,033 to $33, and every night after
		// cleared out of what the business earned. Settle gives things up in
		// order now, so a place earning less than its counter costs goes
		// unpaid night after night and the count climbs.
	}
	return short
}

// EverybodyGotPaid forgets the unpaid nights, because they were paid. Called
// from the night the bills clear.
func (w *World) EverybodyGotPaid() {
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Unpaid > 0 && w.Own(l.ID) {
			prop.Unpaid = 0
		}
	}
}
