package core

import "fmt"

// What you are carrying decides how a bad night ends. Arms come off a boat and
// are never legal to hold: they make violence go better and a police search go
// much worse, which is the whole trade.
//
// This is equipment, not a stat line. A better weapon shifts the odds of
// something the player chose to start; armour reduces what it costs when it
// goes wrong. Neither makes anyone safe.

type Armament struct {
	Tier   int
	Label  string
	Detail string
	Cost   int
}

var weapons = []Armament{
	{Tier: 0, Label: "Nothing but your hands"},
	{Tier: 1, Label: "A revolver", Detail: "Old, reliable enough, and easy to lose in a river.", Cost: 220},
	{Tier: 2, Label: "A pump shotgun", Detail: "Nobody argues with it at close range. Hard to explain if you are searched.", Cost: 650},
	{Tier: 3, Label: "A Thompson", Detail: "The kind of thing that ends an argument and starts an investigation.", Cost: 1800},
}

var armour = []Armament{
	{Tier: 0, Label: "Just a coat"},
	{Tier: 1, Label: "A padded vest", Detail: "Takes some of the sting out of a beating.", Cost: 180},
	{Tier: 2, Label: "A steel plate vest", Detail: "Heavy, obvious under a jacket, and it has saved people.", Cost: 700},
}

// ArmsSource is where arms change hands: off a boat, not over a counter.
func ArmsSource(id string) bool { return id == "docks" }

// Armaments is everything of this kind that can be had, without the empty
// hands at the top of the list. The dock offered exactly one of each — the next
// one up — so getting a Thompson meant buying a revolver and a shotgun first
// and throwing both away: "when buying guns/armor and whatnot I don't think you
// should have to progress through them, you should be able to buy any of them
// at any time, you don't need to go through some sort of upgrade cycle."
//
// Everything is on the counter now, each at its own price. What you cannot do
// is pay for something worse than what you are already carrying, and the card
// says so rather than going quiet: a revolver bought over a Thompson would be a
// trap, not a choice, because nothing in this city rewards carrying less gun.
func Armaments(kind string) []Armament {
	list := weapons
	if kind == "armour" {
		list = armour
	}
	return append([]Armament{}, list[1:]...)
}

// carrying is what the player has of this kind.
func (w *World) carrying(kind string) int {
	if kind == "armour" {
		return w.Player.Armour
	}
	return w.Player.Weapon
}

// WeaponEdge is how much a weapon shifts the odds of violence the player
// starts. It is deliberately modest: being armed is an advantage, not a promise.
func (w *World) WeaponEdge() float64 { return float64(w.Player.Weapon) * .07 }

// Absorb reduces injury by what is being worn, never below a real cost. Armour
// makes a beating survivable; it does not make it free.
func (w *World) Absorb(injury int) int {
	reduced := injury - w.Player.Armour*7
	return max(injury/3, max(1, reduced))
}

// ArmsReadiness explains why an upgrade cannot be bought, or returns "".
func (w *World) ArmsReadiness(kind string, tier int) string {
	if !ArmsSource(w.Player.Location) {
		return "Nobody sells this here"
	}
	list := weapons
	if kind == "armour" {
		list = armour
	}
	if tier <= 0 || tier >= len(list) {
		return "Nobody sells this here"
	}
	have := w.carrying(kind)
	if tier == have {
		return "You are carrying it"
	}
	if tier < have {
		return "You are already carrying something better"
	}
	if w.Player.Cash < list[tier].Cost {
		return "Not enough cash"
	}
	return ""
}

// BuyArms buys the one that was asked for, whichever it is. Holding arms is
// itself a reason for the police to take an interest.
func (w *World) BuyArms(kind string, tier int) error {
	if reason := w.ArmsReadiness(kind, tier); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	list := weapons
	if kind == "armour" {
		list = armour
	}
	next := list[tier]
	if err := w.Pay(next.Cost); err != nil {
		return err
	}
	if kind == "armour" {
		w.Player.Armour = next.Tier
	} else {
		w.Player.Weapon = next.Tier
	}
	w.Player.Heat = min(100, w.Player.Heat+2)
	w.Log("Off a boat at Pier 14", fmt.Sprintf("%s, $%d. Carrying it is its own risk if anyone searches you.", next.Label, next.Cost), "personal")
	return nil
}

// SeizeArms is what a search costs somebody who is armed. Handled separately
// from stock because arms are evidence rather than merchandise.
func (w *World) SeizeArms() bool {
	if w.Player.Weapon == 0 && w.Player.Armour == 0 {
		return false
	}
	w.Player.Weapon, w.Player.Armour = 0, 0
	w.Log("They took the hardware", "Everything you were carrying is evidence now. You are unarmed.", "danger")
	return true
}

// ArmsDescription is what the player is currently carrying, for the interface.
func (w *World) ArmsDescription() map[string]any {
	return map[string]any{
		"weapon":  weapons[min(w.Player.Weapon, len(weapons)-1)].Label,
		"armour":  armour[min(w.Player.Armour, len(armour)-1)].Label,
		"charges": w.Player.Charges,
	}
}
