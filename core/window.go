package core

import "fmt"

// The pawnbroker's window.
//
// A ticket running out used to mean the thing was gone: "Ackerman sold it to
// somebody else", and that was the end of it. That is true of a pawnbroker and
// it is also a wall — the player is told a decision has been made for them and
// there is nothing on the other side of it. The thing did not evaporate. It
// went in the window, and a window is a thing you can stand in front of.
//
// So the window is real. What this city could not redeem is on the shelf with a
// price on it, and anybody with the money can have it, the player included, and
// the thing that was theirs last week included. That is the difference between
// a loss and a bad trade.
//
// It is also the half of this trade that was missing. Four businesses reach
// past their own income by changing what the player already pays: a garage
// halves the car's upkeep, a haulier takes a third off stocking, a yard buys
// the wreck. The pawnbroker only ever pushed a number on the city's side. Hold
// the shop and the window is your own stock: you take what is on the shelf at
// what the counter lent on it rather than at what it asks, which is the whole
// of a pawnbroker's margin, handed to whoever owns the counter.

const (
	// WindowAsk is what the shop wants for a thing on the shelf, as a
	// percentage of what it cost new. The counter lends PawnLend and asks
	// WindowAsk, and the difference between the two is the trade.
	WindowAsk = 60
	// WindowHolds is how much shelf there is. The oldest thing goes to a
	// customer who is not the player when a new one comes in — a window is not
	// a warehouse, and a shelf that only grew would end the campaign with every
	// car the city ever lost on it.
	WindowHolds = 6
	// WindowMinutes is how long buying something out of a window takes.
	WindowMinutes = 30
	// WindowBroke is the purse below which somebody is short enough to pawn
	// what they own and not get it back. The city stocks this shelf; the player
	// pawning their own suit is not the only way anything reaches it.
	WindowBroke = 25
)

// Shelf is one thing in the window, with what the shop wants for it.
type Shelf struct {
	// ID is stable for as long as the thing is on the shelf, because a card is
	// matched to a command by its id and a shelf reorders.
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Tier int    `json:"tier"`
	Wear int    `json:"wear"`
	// Ask is what the shop wants, fixed when it went in the window so the price
	// does not move while the player is looking at it.
	Ask int `json:"ask"`
	// Lent is what the counter gave whoever brought it in, which is what an
	// owner pays instead of the asking price.
	Lent int `json:"lent"`
	// Whose is who could not redeem it, for the card to say. Empty is somebody
	// the city does not name.
	Whose string `json:"whose,omitempty"`
	// Yours is whether the player is looking at their own thing.
	Yours bool `json:"yours,omitempty"`
}

// What the thing on the shelf is called.
func (s Shelf) What() string {
	if s.Kind == "dress" {
		return AttireByTier(s.Tier).Label
	}
	return VehicleByTier(s.Tier).Label
}

// Condition is how hard a thing has been used, in a word. Two pressed suits at
// different prices read as the same card without it, and two cards in a room
// under one name is the oldest fault in this project.
func (s Shelf) Condition() string {
	switch {
	case s.Wear >= 65:
		return "a hard-worn"
	case s.Wear >= 40:
		return "a worn"
	case s.Wear >= 20:
		return "a second-hand"
	}
	return "a barely used"
}

// OnTheShelf is how the thing reads on a card: the condition and the thing,
// with the article the label carries taken off so the two do not collide.
func (s Shelf) OnTheShelf() string {
	what := lowerFirst(s.What())
	for _, article := range []string{"a ", "an ", "the "} {
		if len(what) > len(article) && what[:len(article)] == article {
			what = what[len(article):]
			break
		}
	}
	return s.Condition() + " " + what
}

// worthNew is what the kind and tier cost over a counter that sells them new.
func worthNew(kind string, tier int) int {
	if kind == "dress" {
		return AttireByTier(tier).Cost
	}
	return VehicleByTier(tier).Cost
}

// askFor is what the shop wants for something in the condition it came in.
// Wear comes off the price, because a window is full of things that have been
// used and the player can see that they have.
func askFor(kind string, tier, wear int) int {
	return max(1, worthNew(kind, tier)*WindowAsk/100*max(20, 100-wear)/100)
}

// shelve puts a thing in the window and lets the oldest go if the shelf is
// full.
func (w *World) shelve(s Shelf) {
	s.ID = fmt.Sprintf("%s-%d", s.Kind, w.Minute)
	for i := range w.Window {
		if w.Window[i].ID == s.ID {
			s.ID = fmt.Sprintf("%s-%d-%d", s.Kind, w.Minute, i+1)
		}
	}
	w.Window = append(w.Window, s)
	for len(w.Window) > WindowHolds {
		gone := w.Window[0]
		w.Window = w.Window[1:]
		if gone.Yours {
			w.Log("Somebody bought it",
				fmt.Sprintf("%s is out of the window at %s. Somebody else is driving it, wearing it, or has it in a room you will never see.",
					gone.What(), placeName(w.thePawnshop())), "personal")
		}
	}
}

// shelfByID finds one thing in the window.
func (w *World) shelfByID(id string) *Shelf {
	for i := range w.Window {
		if w.Window[i].ID == id {
			return &w.Window[i]
		}
	}
	return nil
}

// WindowLabels is what each thing on the shelf is called on its own card,
// guaranteed to differ from every other card on the shelf.
//
// A window of second-hand suits is exactly where two cards in a room under one
// name happens, and it did: two "second-hand pressed suit" cards at different
// prices. The condition word separates most of them, the price separates the
// rest, and where a shelf holds two things that are genuinely the same thing at
// the same price in the same state, they are counted — because a player
// clicking one of two identical cards should still be able to tell which one
// they clicked.
func (w *World) WindowLabels() map[string]string {
	name := map[string]string{}
	count := map[string]int{}
	for _, s := range w.Window {
		label := "Take " + s.OnTheShelf() + " out of the window"
		count[label]++
	}
	twice := map[string]int{}
	for _, s := range w.Window {
		label := "Take " + s.OnTheShelf() + " out of the window"
		if count[label] > 1 {
			label = fmt.Sprintf("Take %s out of the window at $%d", s.OnTheShelf(), w.WindowPrice(s.ID))
		}
		name[s.ID] = label
	}
	// Anything still sharing a name is the same thing at the same price, so it
	// is numbered.
	again := map[string]int{}
	for _, label := range name {
		again[label]++
	}
	for _, s := range w.Window {
		if again[name[s.ID]] > 1 {
			twice[name[s.ID]]++
			name[s.ID] = fmt.Sprintf("%s (%s)", name[s.ID], counted(twice[name[s.ID]], "one", "ones"))
		}
	}
	return name
}

// WindowPrice is what this costs the player. The asking price, unless the
// counter is theirs — an owner takes stock at what the shop paid for it, which
// is a pawnbroker's whole margin and the reason to hold one.
func (w *World) WindowPrice(id string) int {
	s := w.shelfByID(id)
	if s == nil {
		return 0
	}
	if w.Own(w.thePawnshop()) {
		return max(1, s.Lent)
	}
	return s.Ask
}

// WindowReadiness explains why something in the window cannot be bought, or
// returns "".
func (w *World) WindowReadiness(id string) string {
	if w.Player.Location != w.thePawnshop() {
		return "This is done over a counter"
	}
	s := w.shelfByID(id)
	if s == nil {
		return "That is not in the window"
	}
	// Nowhere to put it. A man already in a suit cannot also be in this one,
	// and the game has one car.
	switch s.Kind {
	case "car":
		if w.Player.Car > 0 {
			return "You have a car already. Sell or pawn the one you have first"
		}
	case "dress":
		if w.Player.Dress >= s.Tier {
			return "What you have on is no worse than that"
		}
	}
	if w.Player.Cash < w.WindowPrice(id) {
		return fmt.Sprintf("It is $%d", w.WindowPrice(id))
	}
	return ""
}

// BuyFromWindow takes a thing off the shelf.
func (w *World) BuyFromWindow(id string) error {
	if reason := w.WindowReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	s := *w.shelfByID(id)
	price := w.WindowPrice(id)
	if err := w.Pay(price); err != nil {
		return err
	}
	switch s.Kind {
	case "car":
		w.Player.Car, w.Player.CarWear, w.Player.Plate = s.Tier, s.Wear, 0
	case "dress":
		w.Player.Dress, w.Player.DressWear = s.Tier, s.Wear
	}
	kept := w.Window[:0]
	for _, x := range w.Window {
		if x.ID != id {
			kept = append(kept, x)
		}
	}
	w.Window = kept
	// A thing coming off the shelf is the shop's day, whoever bought it.
	w.FenceAbout(FenceTrade)
	back := ""
	if s.Yours {
		back = " It was yours until the ticket ran out, and now it is yours again."
	}
	w.Log("Out of the window", fmt.Sprintf("%s for $%d at %s.%s",
		s.What(), price, placeName(w.thePawnshop()), back), "personal")
	return nil
}

// WindowDay is the city stocking the shelf.
//
// The first version of this asked a modelled person to be broke and to own a
// car, and measured: fifteen people in this city drive, none of them is short,
// and the thinnest purse among them is $164. That is not a bug in the city —
// a man with a car is not the man pawning one — it just means the rule could
// never fire, and a window that is only ever stocked by the player pawning
// their own suit is not a window.
//
// So most of what is on the shelf comes from the city at large, unnamed. There
// are a few hundred people in this town the game does not model and a
// pawnbroker's window is full of their things. When a person the game does
// model really is that short and really does own a car, they are named, and
// that is worth something precisely because it is rare.
const (
	// WindowComes is the chance in a hundred that something reaches the shelf
	// on an ordinary day, while there is room for it.
	WindowComes = 35
	// WindowPoor is the highest tier the city's poor are pawning. Nobody
	// hands a Packard over this counter and fails to come back for it.
	WindowPoor = 2
)

func (w *World) WindowDay() {
	// Derived from the campaign and the day rather than drawn from the world's
	// random stream, which is the same trick the evening haunt and the shift
	// change use: fixed for this city on this day, identical on a replay, and
	// costing the stream nothing.
	//
	// It started as four stream draws a day and broke two tests that had nothing
	// to do with pawnbrokers — a manager keeping a bar stocked, and what a
	// table pays into a till. Both read the world stream in absolute terms, so
	// anything that consumes it moves them. Taking the draws unconditionally
	// was not enough: the stream still shifts for everybody. A daily feature
	// should not be spending the city's randomness at all.
	day := w.Minute / 1440
	chance := roll(w.ID, day, 1) % 100
	sort := roll(w.ID, day, 2) % 100
	pick := roll(w.ID, day, 3) % 100
	worn := roll(w.ID, day, 4) % 100
	if w.thePawnshop() == "" || len(w.Window) >= WindowHolds {
		return
	}
	// Somebody the game knows, short enough to give up the car. Rare.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Purse > WindowBroke || n.Car <= 0 {
			continue
		}
		wear := min(90, 20+worn*60/100)
		lent := max(1, worthNew("car", n.Car)*PawnLend/100*max(20, 100-wear)/100)
		w.shelve(Shelf{
			Kind: "car", Tier: n.Car, Wear: wear,
			Ask: askFor("car", n.Car, wear), Lent: lent, Whose: n.Name,
		})
		if w.Own(w.thePawnshop()) {
			w.Log("Somebody was short", fmt.Sprintf("%s left a %s over your counter for $%d and is not coming back for it. It is in the window.",
				n.Name, lowerFirst(VehicleByTier(n.Car).Label), lent), "business")
		}
		n.Car, n.Purse = 0, n.Purse+lent
		return
	}
	if chance >= WindowComes {
		return
	}
	kind := "dress"
	if sort < 30 {
		kind = "car"
	}
	tier := 1 + pick*WindowPoor/100
	wear := min(90, 25+worn*55/100)
	lent := max(1, worthNew(kind, tier)*PawnLend/100*max(20, 100-wear)/100)
	w.shelve(Shelf{
		Kind: kind, Tier: tier, Wear: wear,
		Ask: askFor(kind, tier, wear), Lent: lent,
	})
}

// roll is a number from 0 up, fixed for this campaign, this day and this
// question. Nothing about it touches either random stream.
//
// The mixing step is the whole of it. The first version ran the campaign id
// through `sum*131 + c` over a starting value made of the day and the salt,
// which is an affine map of that start — so every answer marched in step with
// the day and the shelf filled with four of the same suit at almost the same
// price. A hash used as a substitute for randomness has to avalanche or it is
// a counter wearing a disguise.
func roll(id string, day, salt int) int {
	sum := uint64(salt)*0x9E3779B97F4A7C15 + uint64(day)*0xBF58476D1CE4E5B9
	for i := 0; i < len(id); i++ {
		sum = (sum ^ uint64(id[i])) * 0x100000001B3
	}
	sum ^= sum >> 33
	sum *= 0xFF51AFD7ED558CCD
	sum ^= sum >> 29
	return int(sum >> 33)
}

// theBroker is whoever is behind the counter, by name. "Ackerman" was written
// into the log the day the shop was built and stayed there through every
// campaign in which Ackerman was somebody else or nobody — the same fault as
// the coffee bought for a fixer who had been dead a week.
func (w *World) theBroker() string {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Role == "Pawnbroker" {
			return n.Name
		}
	}
	return "the counter"
}

// tierOf and wearOf are what the player is holding of a kind, so a card can say
// what the thing would be worth in the window before it goes there.
func tierOf(w *World, kind string) int {
	if kind == "dress" {
		return w.Player.Dress
	}
	return w.Player.Car
}

func wearOf(w *World, kind string) int {
	if kind == "dress" {
		return w.Player.DressWear
	}
	return w.Player.CarWear
}
