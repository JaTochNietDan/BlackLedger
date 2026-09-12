package core

import "fmt"

// Respect is earned and cannot be bought. How you look can be bought, and in
// this city it is read first: a man is judged at the door long before anyone
// asks what he has done.
//
// Attire is therefore a second, purchasable half of standing, and it is
// deliberately fragile. It wears out with the days, it is ruined by every
// beating and search, and a good suit on somebody with no visible income is
// exactly the thing a detective remembers. It opens doors that respect alone
// does not, and it costs money to keep open.

// Attire is a level of dress and what it is worth in a room.
type Attire struct {
	Tier int
	// Label is what the player is wearing, in words.
	Label string
	// Detail is why anybody would pay for it.
	Detail string
	// Cost is the price of moving up to it.
	Cost int
	// Presence is what it is worth when somebody is sizing you up, at full
	// condition. Wear takes it away proportionally.
	Presence int
	// Notice is the police attention it draws each day, because dressing above
	// your visible means is a question waiting to be asked.
	Notice int
}

var attires = []Attire{
	{Tier: 0, Label: "Working clothes", Detail: "Nobody looks twice, which is its own kind of useful."},
	{Tier: 1, Label: "A pressed suit", Detail: "Off a rack, but clean, and it says you are not nobody.", Cost: 190, Presence: 5},
	{Tier: 2, Label: "A tailored suit", Detail: "Cut for you. Doors open a little further and stay open a little longer.", Cost: 680, Presence: 11, Notice: 1},
	{Tier: 3, Label: "Bespoke, from a house that knows your name", Detail: "The kind of clothes people describe to each other afterwards.", Cost: 1900, Presence: 19, Notice: 2},
}

const (
	// DressUpkeep is the condition a suit loses each day simply from being worn
	// in this city.
	DressUpkeep = 2
	// PressCost is what having it cleaned and put right costs, unless you own
	// somewhere that does that sort of thing.
	PressCost = 25
	// PressMinutes is how long that takes.
	PressMinutes = 45
	// Shabby is the condition below which nobody is impressed by any of it.
	Shabby = 40
)

// AttireByTier is the dress at a tier, clamped so an unknown save cannot panic.
func AttireByTier(tier int) Attire { return attires[max(0, min(tier, len(attires)-1))] }

// DressCondition is how well kept what the player is wearing currently is.
// Saves written before clothes existed carry a zero here and no suit, which is
// the same as working clothes in perfect order.
func (w *World) DressCondition() int {
	if w.Player.Dress == 0 {
		return 100
	}
	return max(0, min(100, w.Player.DressWear))
}

// Standing is what attire is worth right now: nothing at all once it is shabby,
// because a ruined good suit reads worse than honest working clothes.
func (w *World) Standing() int {
	condition := w.DressCondition()
	if condition < Shabby {
		return 0
	}
	return AttireByTier(w.Player.Dress).Presence * condition / 100
}

// Presence is how the player reads to somebody who does not know them: what
// they have done, plus what they are wearing. This is the number rooms and
// doormen judge, rather than respect alone.
func (w *World) Presence() int { return w.Player.Respect + w.Standing() }

// Ruin is what violence, a search or a bad night does to good clothes. Called
// wherever the player takes something physical, so the suit is a consumable and
// not a permanent purchase.
func (w *World) Ruin(amount int) {
	if w.Player.Dress == 0 || amount <= 0 {
		return
	}
	before := w.DressCondition()
	w.Player.DressWear = max(0, before-amount)
	if before >= Shabby && w.Player.DressWear < Shabby {
		w.Log("The suit is finished", fmt.Sprintf("%s, torn and marked past anything a brush will fix. Until it is put right you are just another face in the street.", AttireByTier(w.Player.Dress).Label), "personal")
	}
}

// DressDay is the ordinary wear of walking around this city.
func (w *World) DressDay() {
	if w.Player.Dress == 0 {
		return
	}
	w.Ruin(DressUpkeep)
	if notice := AttireByTier(w.Player.Dress).Notice; notice > 0 && w.DressCondition() >= Shabby {
		w.Player.Heat = min(100, w.Player.Heat+notice)
	}
}

// Attires is everything on the rail, working clothes included. All of it is for
// sale at any time: a ladder you had to climb a rung at a time meant the only
// way to be cut for is to buy off a rack first and throw it away, and there is
// no way down at all — which matters here, because a good suit on somebody with
// no visible income is exactly what a detective remembers. Going back into
// working clothes is a decision somebody in this trade makes on purpose.
func Attires() []Attire { return attires }

// TailorReachable reports where clothes are bought. It used to be the exchange:
// the market sold everything else, so it sold this too, and the one purchase in
// this game that is about how you are read happened at a counter between the
// fish and the cloth. There is a tailor's now.
func TailorReachable(location string) bool {
	place, ok := PlaceByID(location)
	return ok && place.Kind == "tailor"
}

// DressMargin is the share of a suit's price that stays with the shop. Holding
// the tailor's is a workroom of your own, and what your own cutters make you
// costs what the cloth cost.
const DressMargin = 30

// DressPrice is what this one costs today. A shop of the player's own keeps no
// margin from them.
func (w *World) DressPrice(tier int) int {
	attire := AttireByTier(tier)
	if w.Own(w.Player.Location) {
		return attire.Cost - attire.Cost*DressMargin/100
	}
	return attire.Cost
}

// DressReadiness explains why this one cannot be had, or returns "".
func (w *World) DressReadiness(tier int) string {
	if !TailorReachable(w.Player.Location) {
		return "Nobody sells this here"
	}
	if tier < 0 || tier >= len(attires) {
		return "Nobody sells this here"
	}
	if tier == w.Player.Dress && w.DressCondition() >= 100 {
		return "It is what you are standing in"
	}
	if w.Player.Cash < w.DressPrice(tier) {
		return "Not enough cash"
	}
	return ""
}

// BuyAttire puts the player in the one they asked for, whichever it is,
// arriving in perfect condition. Nobody takes the old one: clothes off a man's
// back are worth nothing to a shop that makes them, which is the difference
// between this counter and a forecourt.
func (w *World) BuyAttire(tier int) error {
	if reason := w.DressReadiness(tier); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	next := AttireByTier(tier)
	price := w.DressPrice(tier)
	if err := w.Pay(price); err != nil {
		return err
	}
	// The shop keeps its margin, and if the player holds it they are buying
	// from themselves — the money never leaves their pocket.
	if shop := w.Properties[w.Player.Location]; shop != nil && !w.Own(w.Player.Location) {
		if house := w.faction(shop.Owner); house != nil {
			house.Cash += next.Cost * DressMargin / 100
		}
	}
	w.Player.Dress, w.Player.DressWear = next.Tier, 100
	place, _ := PlaceByID(w.Player.Location)
	if next.Tier == 0 {
		w.Log("Out of the suit at "+place.Name,
			fmt.Sprintf("%s. Nobody looks at you twice now, which is the point of it.", next.Label), "personal")
		return nil
	}
	w.Log("Fitted at "+place.Name, fmt.Sprintf("%s, $%d. %s", next.Label, price, next.Detail), "personal")
	return nil
}

// PressPlace reports whether clothes can be put right here: your own home, or a
// laundry, which is the one thing a laundry is actually for.
func (w *World) PressPlace(id string) bool {
	return id == w.Player.Home || w.PressAtOwnPlace(id)
}

// PressAtOwnPlace reports whether this is a laundry of the player's. It asks
// what KIND of business the address is, because the city can hold more than one
// laundry and a rule written about "the laundry" would have been a rule about
// one street corner.
func (w *World) PressAtOwnPlace(id string) bool {
	place, ok := PlaceByID(id)
	return ok && place.Kind == "laundry" && w.Own(id)
}

// PressFee is nothing at a laundry of your own. It is the only return the
// player ever gets on that business that is not money.
func (w *World) PressFee(id string) int {
	if w.PressAtOwnPlace(id) {
		return 0
	}
	return PressCost
}

// PressReadiness explains why clothes cannot be put right, or returns "".
func (w *World) PressReadiness(id string) string {
	if !w.PressPlace(id) {
		return "This is not somewhere that is done"
	}
	if w.Player.Dress == 0 {
		return "There is nothing here worth pressing"
	}
	if w.DressCondition() >= 100 {
		return "It is already immaculate"
	}
	if w.Player.Cash < w.PressFee(id) {
		return "Not enough cash"
	}
	return ""
}

// Press restores what a brush and a press can restore. A suit taken far enough
// down never comes all the way back, which is what makes replacing it a
// decision rather than a formality.
func (w *World) Press(id string) error {
	if reason := w.PressReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.PressFee(id)
	if err := w.Pay(fee); err != nil {
		return err
	}
	before := w.DressCondition()
	ceiling := 100
	if before < Shabby {
		ceiling = 80 // it is never quite the same suit again
	}
	w.Player.DressWear = min(ceiling, before+45)
	place, _ := PlaceByID(id)
	where := "at " + place.Name
	if fee == 0 {
		where = "by your own people at " + place.Name
	}
	w.Log("Cleaned and pressed", fmt.Sprintf("%s put right %s. Condition %d of 100.", AttireByTier(w.Player.Dress).Label, where, w.Player.DressWear), "personal")
	return nil
}

// AppearanceDescription is what the player is wearing, for the interface.
func (w *World) AppearanceDescription() map[string]any {
	return map[string]any{
		"attire":    AttireByTier(w.Player.Dress).Label,
		"condition": w.DressCondition(),
		"standing":  w.Standing(),
		"presence":  w.Presence(),
	}
}

// lowerFirst is for reading a label into the middle of a sentence.
func lowerFirst(s string) string {
	if s == "" || s[0] < 'A' || s[0] > 'Z' {
		return s
	}
	return string(s[0]+32) + s[1:]
}
