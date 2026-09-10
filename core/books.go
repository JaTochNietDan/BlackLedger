package core

import "fmt"

// The Ledger screen showed the underground market, two figures and then sixty
// rows of undifferentiated history — five screens of scrolling for a page whose
// whole job is to answer "what am I worth and what is this costing me". The
// arithmetic existed: `DailyCost` added nine things together and returned a
// single number, so a player who wanted to know why they were losing money had
// to guess which of the nine it was.
//
// The books are truth the core already holds. Adding them up is the core's job;
// deciding how to draw them is not.

// Line is one entry in the books: what it is, what it costs or earns a day, and
// what it is made of in words.
type Line struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
	Detail string `json:"detail,omitempty"`
}

// Books is everything the player's position is made of, in the order a person
// would read it: what comes in, what goes out, what they are holding, and what
// is out of reach.
func (w *World) Books() map[string]any {
	p := &w.Player

	// What comes in. Income is per hour everywhere else in the interface, so
	// this states the day, which is the unit the costs are in.
	income, holdings := 0.0, 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) {
			continue
		}
		holdings++
		income += float64(prop.Income*prop.Condition) / 100 *
			operatingMode(prop.Mode).Take * w.Capacity(l.ID) * w.TradeMultiplier(l.ID) *
			(1 + w.LicenceTake())
	}
	daily := int(income * 24)

	out := []Line{}
	add := func(label string, amount int, detail string) {
		if amount != 0 {
			out = append(out, Line{label, amount, detail})
		}
	}
	add("Rent", HomeRent(p.Home), placeName(p.Home))
	add("Security", 10*p.Security, plural(p.Security, "detail", "details"))
	add("Crew", 12*len(p.Crew), plural(len(p.Crew), "on the payroll", "on the payroll"))
	add("Staff", w.Wages(), plural(w.staffed(), "hand", "hands")+" across your premises")
	add("Your own people", w.MemberWages(), plural(len(w.OwnPeople()), "on the payroll", "on the payroll"))
	add("The car", w.CarUpkeep(), w.carName())
	add("The house", w.ComfortUpkeep(), "what is fitted at "+placeName(p.Home))
	add("Retainers", w.RetainerCost(), w.retainerNames())
	add("Understandings", w.PactCost(), w.pactNames())

	costs := 0
	for _, l := range out {
		costs += l.Amount
	}

	// And what is behind. The day's costs are what the player owes; this is
	// where they have already failed to pay it. Wages are the last thing the
	// night gives up and the only one with people on the other side of it: a
	// week of this and somebody stops coming in, and nobody new will take the
	// job until it is paid. The page whose whole job is "what is this costing
	// me" was the one place in the game that did not say so.
	behind := []map[string]any{}
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || !w.Own(l.ID) || prop.Unpaid == 0 {
			continue
		}
		place, _ := PlaceByID(l.ID)
		behind = append(behind, map[string]any{
			"id": l.ID, "place": place.Name, "nights": prop.Unpaid,
			"hands": len(prop.Hands), "positions": tradeHands(l.ID),
			"shut": w.wordIsOut(l.ID),
		})
	}

	return map[string]any{
		"behind":    behind,
		"income":    daily,
		"costs":     costs,
		"net":       daily - costs,
		"lines":     out,
		"cash":      p.Cash,
		"sheltered": w.Sheltered(),
		"lent":      w.OutOnLoan(),
		"owed":      w.owedToPlayer(),
		"offshore":  w.Offshore,
		"holdings":  holdings,
		"earned":    p.Earned,
	}
}

// owedToPlayer is everything out on the street with a name on it, at the price
// it comes back at rather than what went out.
func (w *World) owedToPlayer() int {
	total := 0
	for _, l := range w.Book() {
		total += l.Owed
	}
	return total
}

// staffed is how many people work the player's premises.
func (w *World) staffed() int {
	total := 0
	for id, prop := range w.Properties {
		if w.Own(id) {
			total += prop.Staff
		}
	}
	return total
}

func (w *World) retainerNames() string {
	names := []string{}
	for _, id := range w.Player.Retainers {
		if w.Retained(id) {
			if o, ok := OfficialByID(id); ok {
				names = append(names, o.Name)
			}
		}
	}
	return join(names)
}

func (w *World) pactNames() string {
	names := []string{}
	for _, p := range w.Pacts {
		if p.Life == w.Life {
			names = append(names, w.factionName(p.With))
		}
	}
	return join(names)
}

// carName is what the player drives, for the books.
func (w *World) carName() string {
	if v := w.VehicleDescription(); v != nil {
		if name, ok := v["car"].(string); ok {
			return name
		}
	}
	return ""
}

func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

func join(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	}
	out := names[0]
	for _, n := range names[1 : len(names)-1] {
		out += ", " + n
	}
	return out + " and " + names[len(names)-1]
}
