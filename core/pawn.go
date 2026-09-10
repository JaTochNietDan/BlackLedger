package core

import "fmt"

// The pawnbroker.
//
// Three brass balls over the door, and it sits at the join between everything
// else in this city. What gets taken off the street has to turn into money
// somewhere, and somebody who is short has to turn what they own into money and
// hope to get it back. Both of those already existed and had nowhere to happen.
//
// It is deliberately a bad deal. They lend a fraction of what a thing is worth
// and want more back than they gave, and if the ticket runs out the thing is
// sold to somebody else. That is the trade, and it is why anybody uses one only
// when they have to.

const (
	// PawnLend is the percentage of a thing's price the counter will advance.
	PawnLend = 35
	// PawnBack is the percentage of the advance it takes to get it back, over
	// and above the advance itself. A month's interest, taken up front.
	PawnBack = 130
	// PawnDays is how long the ticket runs before the thing goes in the window.
	PawnDays = 7
	// PawnMinutes is how long the business of it takes.
	PawnMinutes = 45
	// FenceTrade is what a piece of work coming through the shop is worth to
	// its takings, in trade.
	FenceTrade = 3
)

// Ticket is something of the player's sitting behind the counter.
type Ticket struct {
	// What it is: "car" or "dress", and which one, so the same thing comes back
	// rather than something of the same kind.
	Kind string `json:"kind"`
	Tier int    `json:"tier"`
	Wear int    `json:"wear"`
	// Lent is what they gave for it, and Due is when the ticket runs out.
	Lent int `json:"lent"`
	Due  int `json:"due"`
	// Sold is what happened when nobody came back for it. The ticket is kept
	// rather than thrown away, because the counter should be able to tell you
	// what became of the thing rather than shrugging at you.
	Sold bool `json:"sold,omitempty"`
	// Life ties it to one protagonist. Nobody inherits somebody else's ticket.
	Life int `json:"life"`
}

// thePawnshop is the address, asked for once so nothing has to know its name.
func (w *World) thePawnshop() string {
	for _, l := range Locations {
		if l.Kind == "pawn" {
			return l.ID
		}
	}
	return ""
}

// FenceAbout is what a piece of work coming off the street is worth to the shop
// that buys it. The same shape as PartsAbout: a city with more being taken in
// it is a city where a pawnbroker has a better week, which is the whole reason
// anybody would hold one.
func (w *World) FenceAbout(delta int) {
	for _, l := range Locations {
		if l.Kind == "pawn" {
			w.ShiftCustom(l.ID, "things coming in over the counter with no history", delta)
		}
	}
}

// pawnable is what the player has that a counter would take, with what it is
// worth new.
func (w *World) pawnable(kind string) (label string, price int, held bool) {
	switch kind {
	case "car":
		car := VehicleByTier(w.Player.Car)
		return car.Label, car.Cost, w.Player.Car > 0
	case "dress":
		suit := AttireByTier(w.Player.Dress)
		return suit.Label, suit.Cost, w.Player.Dress > 0
	}
	return "", 0, false
}

// PawnValue is what the counter will advance on it: a fraction of what it cost
// new, less what wear has taken off it.
func (w *World) PawnValue(kind string) int {
	_, price, held := w.pawnable(kind)
	if !held || price <= 0 {
		return 0
	}
	condition := 100
	switch kind {
	case "car":
		condition = w.CarCondition()
	case "dress":
		condition = w.DressCondition()
	}
	return max(1, price*PawnLend/100*max(20, condition)/100)
}

// Ticketed finds the player's ticket of this kind, if there is one.
func (w *World) Ticketed(kind string) *Ticket {
	for i := range w.Tickets {
		if t := &w.Tickets[i]; t.Kind == kind && t.Life == w.Life {
			return t
		}
	}
	return nil
}

// PawnReadiness explains why something cannot be put behind the counter, or
// returns "".
func (w *World) PawnReadiness(kind string) string {
	if w.Player.Location != w.thePawnshop() {
		return "This is done over a counter"
	}
	label, _, held := w.pawnable(kind)
	if !held {
		return "You have nothing of that kind to leave"
	}
	if t := w.Ticketed(kind); t != nil && !t.Sold {
		return "There is already a ticket out on one"
	}
	if w.PawnValue(kind) <= 0 {
		return label + " is not worth anything over this counter"
	}
	return ""
}

// Pawn leaves something behind the counter and takes the money.
func (w *World) Pawn(kind string) error {
	if reason := w.PawnReadiness(kind); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	label, _, _ := w.pawnable(kind)
	lent := w.PawnValue(kind)
	ticket := Ticket{Kind: kind, Lent: lent, Due: w.Minute + PawnDays*1440, Life: w.Life}
	switch kind {
	case "car":
		ticket.Tier, ticket.Wear = w.Player.Car, w.Player.CarWear
		w.Player.Car, w.Player.CarWear, w.Player.Plate = 0, 0, 0
	case "dress":
		ticket.Tier, ticket.Wear = w.Player.Dress, w.Player.DressWear
		w.Player.Dress, w.Player.DressWear = 0, 0
	}
	w.Tickets = append(w.Tickets, ticket)
	w.Earn(lent)
	place, _ := PlaceByID(w.Player.Location)
	w.Log("Left at "+place.Name, fmt.Sprintf("%s, and $%d over the counter for it. The ticket runs %d days; after that it is in the window.",
		label, lent, PawnDays), "personal")
	return nil
}

// what is the thing itself, named from the ticket rather than from what the
// player is holding now — they are not holding it, that is the point.
func (t *Ticket) what(w *World) string {
	if t.Kind == "dress" {
		return AttireByTier(t.Tier).Label
	}
	return VehicleByTier(t.Tier).Label
}

// RedeemPrice is what getting it back costs.
func (w *World) RedeemPrice(kind string) int {
	t := w.Ticketed(kind)
	if t == nil {
		return 0
	}
	return max(t.Lent+1, t.Lent*PawnBack/100)
}

// RedeemReadiness explains why something cannot be got back, or returns "".
func (w *World) RedeemReadiness(kind string) string {
	if w.Player.Location != w.thePawnshop() {
		return "This is done over a counter"
	}
	t := w.Ticketed(kind)
	if t == nil {
		return "There is nothing of yours behind this counter"
	}
	if t.Sold || w.Minute > t.Due {
		return "The ticket ran out and it was sold. It is not coming back"
	}
	if w.Player.Cash < w.RedeemPrice(kind) {
		return fmt.Sprintf("Getting it back costs $%d", w.RedeemPrice(kind))
	}
	return ""
}

// Redeem buys it back.
func (w *World) Redeem(kind string) error {
	if reason := w.RedeemReadiness(kind); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	t := w.Ticketed(kind)
	price := w.RedeemPrice(kind)
	if err := w.Pay(price); err != nil {
		return err
	}
	switch kind {
	case "car":
		w.Player.Car, w.Player.CarWear = t.Tier, t.Wear
	case "dress":
		w.Player.Dress, w.Player.DressWear = t.Tier, t.Wear
	}
	w.dropTicket(kind)
	label, _, _ := w.pawnable(kind)
	place, _ := PlaceByID(w.Player.Location)
	w.Log("Back out of "+place.Name, fmt.Sprintf("%s, for $%d against the $%d they gave you. That is what a counter is for.",
		label, price, t.Lent), "personal")
	return nil
}

func (w *World) dropTicket(kind string) {
	kept := w.Tickets[:0]
	for _, t := range w.Tickets {
		if t.Kind == kind && t.Life == w.Life {
			continue
		}
		kept = append(kept, t)
	}
	w.Tickets = kept
}

// PawnDay sells what nobody came back for. The shop has a better day for it,
// which is the other half of why anybody holds one.
func (w *World) PawnDay() {
	for i := range w.Tickets {
		t := &w.Tickets[i]
		if t.Life != w.Life || t.Sold || w.Minute <= t.Due {
			continue
		}
		t.Sold = true
		what := VehicleByTier(t.Tier).Label
		if t.Kind == "dress" {
			what = AttireByTier(t.Tier).Label
		}
		w.FenceAbout(FenceTrade)
		w.Log("It went in the window", fmt.Sprintf("%s. The ticket ran out and Ackerman sold it to somebody else, which is what a ticket running out means.", what), "personal")
	}
}

// TicketDescription is what is behind the counter, for the interface.
func (w *World) TicketDescription() []map[string]any {
	out := []map[string]any{}
	for _, t := range w.Tickets {
		if t.Life != w.Life || t.Sold {
			continue
		}
		what := VehicleByTier(t.Tier).Label
		if t.Kind == "dress" {
			what = AttireByTier(t.Tier).Label
		}
		out = append(out, map[string]any{
			"kind": t.Kind, "what": what, "lent": t.Lent,
			"back": w.RedeemPrice(t.Kind), "days": max(0, (t.Due-w.Minute)/1440),
		})
	}
	return out
}
