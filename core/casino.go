package core

import "fmt"

// A casino was premises with an hourly number on it, which is the one business
// in this city that number describes least well. A casino is a float: the money
// behind the tables is what the house is willing to lose in a night, and it is
// the only reason anybody with real money walks in.
//
// So an owned casino now runs a night every day. The house keeps its edge over
// the long run and loses on plenty of individual nights, the size of the action
// it can attract is set by what is behind the tables, and a house that cannot
// pay a winner is finished as a room worth visiting. Bankrolling it too thin is
// safe and earns nothing; bankrolling it heavily earns well and puts real money
// where one bad night can reach it.

const (
	// BankrollLot is what goes behind the tables at a time.
	BankrollLot = 250
	// BankrollFull is the float at which a house runs the biggest action it
	// will ever see. Beyond this, more money is only more exposure.
	BankrollFull = 1500
	// NightHandle is what crosses the tables on a full night at a room in good
	// order, before any of the things that reduce it.
	NightHandle = 2400
	// HouseEdge is the percentage of the handle the house keeps in expectation.
	// It matches the edge a player faces at somebody else's tables.
	HouseEdge = 6
	// HighRollerChance is how often somebody with real money sits down at a
	// fully funded house.
	HighRollerChance = .15
	// RuinCondition is what a house loses when it cannot pay a winner. The room
	// keeps trading; it stops being a room where serious money plays.
	RuinCondition = 25
)

// HasBankroll reports whether a property is a room that runs a float behind
// tables of its own — a casino, and only a casino, because that is what the
// nightly handle is worked out for.
func HasBankroll(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Type == "casino"
}

// RunsAGame reports whether a room has money of its own that the player can
// look at, add to and take out. Every casino, and the poolhall, which is a
// racket rather than a casino and so had none of this — while being the one
// room in the city that runs a card game and charges for the seat. The money it
// took went straight into the holder's pocket without ever being anywhere they
// could see it, which is what the note was about: "I don't see its current
// funds or how to add to the funds or withdraw from the funds dynamically."
func RunsAGame(id string) bool { return HasBankroll(id) || id == BackRoom }

// Confidence is how much of its potential action a house attracts, which is
// decided by what is behind the tables. A thin float is not a secret: the
// people who bet seriously know exactly which rooms can cover them.
func (w *World) Confidence(id string) float64 {
	prop := w.Properties[id]
	if prop == nil || !HasBankroll(id) {
		return 0
	}
	return min64(1, float64(prop.Bankroll)/BankrollFull)
}

func min64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// NightHandleAt is what will cross the tables tonight: the room's condition,
// its staffing and stock, how hard it is being run, and what it can cover.
func (w *World) NightHandleAt(id string) int {
	prop := w.Properties[id]
	if prop == nil || !HasBankroll(id) || !w.Own(id) {
		return 0
	}
	// And who is actually in the room. A room's action is the people in it, and
	// the one business in this city that people go out to was the one whose
	// takings ignored them entirely — so a war that emptied every bar in the
	// city left a casino covering exactly as much as it had the night before.
	handle := float64(NightHandle) * w.Capacity(id) * float64(prop.Condition) / 100 * w.Confidence(id) * operatingMode(prop.Mode).Take * w.roomAction(id)
	return int(handle)
}

// roomAction is what the floor is worth tonight. Same shape as the trade a shop
// does, and bounded the same way: a full room is a good night rather than a
// different house.
func (w *World) roomAction(id string) float64 {
	factor := 1 + float64(w.Footfall(id)-TypicalRoom)*PerHead
	if factor < QuietRoom {
		return QuietRoom
	}
	if factor > BusyRoom {
		return BusyRoom
	}
	return factor
}

// CasinoDay runs one night at every casino the player owns. Money won and lost
// stays behind the tables; taking it out is a decision of its own.
func (w *World) CasinoDay() {
	for _, l := range Locations {
		if !HasBankroll(l.ID) || !w.Own(l.ID) {
			continue
		}
		w.night(l)
	}
}

func (w *World) night(l Place) {
	prop := w.Properties[l.ID]
	handle := w.NightHandleAt(l.ID)
	if handle < 50 {
		if prop.Bankroll < BankrollLot {
			w.Log("A quiet night at "+l.Name, "There is not enough behind the tables for anybody serious to sit down. The room takes what it takes at the bar and nothing more.", "business")
		}
		return
	}

	// The ordinary run of the night: the house keeps its edge in expectation,
	// and loses on a fair share of individual nights.
	expected := float64(handle) * HouseEdge / 100
	swing := (w.WorldRandom() - .5) * float64(handle) * .5
	result := int(expected + swing)

	// Somebody with real money, drawn by a house that can cover them. They play
	// to the limit the float allows, and they win rather less than half the
	// time, which is the only thing standing between the house and the door.
	roller, rollerNet := false, 0
	if w.WorldRandom() < HighRollerChance*w.Confidence(l.ID) {
		roller = true
		// What they turn over is set by the room, not by the float: there are
		// only so many people in this city who bet like this and they bet the
		// same way wherever they sit down. The house holds its edge on that
		// turnover, and the swing around it is what can empty the tables. A
		// thin float sees one of these rarely and cannot survive it; a deep one
		// sees them often and barely notices.
		turnover := float64(handle) * (.55 + w.WorldRandom()*1.15)
		rollerNet = int(turnover*HouseEdge/100 + (w.WorldRandom()-.5)*turnover*1.5)
		result += rollerNet
	}

	prop.Bankroll += result
	switch {
	case prop.Bankroll < 0:
		// The house cannot pay. This is the thing a thin float actually costs.
		short := -prop.Bankroll
		prop.Bankroll = 0
		prop.Condition = max(0, prop.Condition-RuinCondition)
		w.Player.Respect = max(0, w.Player.Respect-5)
		w.Log("The tables could not cover it at "+l.Name,
			fmt.Sprintf("A winner was owed $%d more than there was behind the tables. Word of that travels faster than anything else in this city, and the serious money will drink somewhere else now.", short), "danger")
		w.Report("business", "HOUSE CANNOT PAY AT "+upper(l.Name),
			fmt.Sprintf("A winner at %s is said to have been turned away without full payment. Patrons described the room afterwards as emptying quickly.", l.Name))
	case roller && rollerNet < 0:
		w.Log("A bad night at "+l.Name,
			fmt.Sprintf("Somebody sat down with intent and left $%d of the house's money up. There is $%d behind the tables now.", -rollerNet, prop.Bankroll), "business")
	case roller:
		w.Log("A good night at "+l.Name,
			fmt.Sprintf("Somebody with real money played long and left $%d of it here. The float is up to $%d.", rollerNet, prop.Bankroll), "business")
	case result < 0:
		w.Log("The tables lost at "+l.Name,
			fmt.Sprintf("$%d of the float went out of the door tonight, leaving $%d behind the tables. The house wins over a season, not over a night.", -result, prop.Bankroll), "business")
	default:
		w.Log("A night at "+l.Name,
			fmt.Sprintf("$%d across the tables, $%d of it kept. The float stands at $%d.", handle, result, prop.Bankroll), "business")
	}
}

// BankrollLeast is the smallest figure worth calling a float. Below it the
// player is not funding a room, they are rounding.
const BankrollLeast = 25

// BankrollReadiness explains why the typed figure cannot be put behind the
// tables, or returns "". A figure of zero is whatever the field would have
// started on, so a caller that names no amount still moves a lot.
func (w *World) BankrollReadiness(id string, amount int) string {
	if !RunsAGame(id) || !w.Own(id) {
		return "This is not a room of yours"
	}
	amount = w.BankrollSum(id, amount)
	if w.Player.Cash < BankrollLeast {
		return "Not enough cash"
	}
	if amount < BankrollLeast {
		return fmt.Sprintf("A float starts at $%d", BankrollLeast)
	}
	if amount > w.Player.Cash {
		return fmt.Sprintf("You have $%d", w.Player.Cash)
	}
	return ""
}

// BankrollSum is the figure a request actually moves: what was typed, or the
// lot when nothing was.
func (w *World) BankrollSum(id string, amount int) int {
	if amount > 0 {
		return amount
	}
	return min(BankrollLot, max(BankrollLeast, w.Player.Cash))
}

// Bankroll puts money of the player's own behind the tables, where the house
// can lose it and where it is the reason anybody worth taking money from comes
// through the door.
func (w *World) Bankroll(id string, amount int) error {
	if reason := w.BankrollReadiness(id, amount); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	amount = w.BankrollSum(id, amount)
	if err := w.Pay(amount); err != nil {
		return err
	}
	prop := w.Properties[id]
	prop.Bankroll += amount
	place, _ := PlaceByID(id)
	w.Log("Money behind the tables at "+place.Name,
		fmt.Sprintf("$%d more, for $%d in the float. The room can now cover about $%d of action a night.", amount, prop.Bankroll, w.NightHandleAt(id)), "business")
	return nil
}

// DrawReadiness explains why the typed figure cannot come off the tables, or
// returns "".
func (w *World) DrawReadiness(id string, amount int) string {
	if !RunsAGame(id) || !w.Own(id) {
		return "This is not a room of yours"
	}
	float := w.Properties[id].Bankroll
	if float < BankrollLeast {
		return "There is not enough behind the tables to take any out"
	}
	amount = w.DrawSum(id, amount)
	if amount < BankrollLeast {
		return fmt.Sprintf("Take $%d or more", BankrollLeast)
	}
	if amount > float {
		return fmt.Sprintf("There is $%d behind the tables", float)
	}
	return ""
}

// DrawSum is the figure a request actually takes off the tables.
func (w *World) DrawSum(id string, amount int) int {
	if amount > 0 {
		return amount
	}
	return min(BankrollLot, w.Properties[id].Bankroll)
}

// Draw takes money back out of the float and into the player's hands. It is the
// only way the money a casino makes ever reaches them, and everything taken is
// action the room can no longer attract.
func (w *World) Draw(id string, amount int) error {
	if reason := w.DrawReadiness(id, amount); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	amount = w.DrawSum(id, amount)
	prop := w.Properties[id]
	prop.Bankroll -= amount
	w.Earn(amount)
	place, _ := PlaceByID(id)
	w.Log("Taken off the tables at "+place.Name,
		fmt.Sprintf("$%d in your hands, $%d left in the float. The room can cover about $%d of action a night now.", amount, prop.Bankroll, w.NightHandleAt(id)), "business")
	return nil
}

// coverage puts the room's nightly capacity into words, including the case
// where there is nothing behind the tables at all.
func coverage(handle int) string {
	if handle <= 0 {
		return "The tables are dark: nobody plays in a room that cannot cover them."
	}
	return fmt.Sprintf("It covers about $%d of action a night.", handle)
}
