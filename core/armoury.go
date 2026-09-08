package core

import "fmt"

// Arms could be bought and never sold. The player could carry a gun and nobody
// else in the city could buy one, which meant a war between two organizations
// was fought with whatever they already had and cost them nothing but people.
//
// An armoury is a room under a business the player owns with crates in it, and
// the customers are the wars. An organization fighting one will buy from
// whoever has stock, at a price nobody would pay in peacetime, and the strength
// it buys is real. It is the most profitable thing in this game and the most
// difficult to explain to anybody.

const (
	// ArmouryCost is what fitting one out takes.
	ArmouryCost = 1200
	// ArmouryMinutes is how long that takes.
	ArmouryMinutes = 180
	// ArmouryHeat is the attention the room itself draws each day once there
	// is anything in it, before what is in it is counted.
	ArmouryHeat = 3
	// ArmouryCrateHeat is how many crates it takes to add another point. A
	// room with sixty in it is a louder room than one with five.
	ArmouryCrateHeat = 12
	// ArmouryHold is how many crates a room will take.
	ArmouryHold = 60
	// WarPremium is what an organization at war pays, as a percentage of the
	// waterfront price. Nobody haggles when they are losing.
	WarPremium = 235
	// ArmsStrength is the strength a crate is worth to whoever buys it.
	ArmsStrength = 2
	// ArmouryGoodwill is what selling to somebody is worth with them.
	ArmouryGoodwill = 3
)

// ArmourySite reports whether a business can hide a room like this. A casino
// floor cannot: too many people, and the wrong ones.
func ArmourySite(id string) bool { return id == "laundry" || id == "garage" }

// TheArmoury finds the player's armoury, if they have one.
func (w *World) TheArmoury() (string, bool) {
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Armoury && w.Own(l.ID) {
			return l.ID, true
		}
	}
	return "", false
}

// Stocked is how many crates are in the armoury.
func (w *World) Stocked() int {
	id, ok := w.TheArmoury()
	if !ok {
		return 0
	}
	return w.Properties[id].Crates
}

// ArmouryReadiness explains why one cannot be fitted out, or returns "".
func (w *World) ArmouryReadiness(id string) string {
	if !ArmourySite(id) || !w.Own(id) {
		return "There is nowhere here for a room like that"
	}
	if w.Properties[id].Armoury {
		return "There is already one under this floor"
	}
	if _, ok := w.TheArmoury(); ok {
		return "Everything in one place is bad enough. Two places is asking for it"
	}
	if w.Player.Cash < ArmouryCost {
		return "Not enough cash"
	}
	return ""
}

// BuildArmoury puts one in.
func (w *World) BuildArmoury(id string) error {
	if reason := w.ArmouryReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(ArmouryCost); err != nil {
		return err
	}
	w.Properties[id].Armoury = true
	place, _ := PlaceByID(id)
	w.Log("A room under "+place.Name, fmt.Sprintf("$%d of brick, board and a door that locks from the outside. It holds %d crates, draws %d attention a day plus one for every %d in it, and there is exactly one kind of customer for what goes in there.", ArmouryCost, ArmouryHold, ArmouryHeat, ArmouryCrateHeat), "business")
	return nil
}

// StockReadiness explains why crates cannot be put away, or returns "".
func (w *World) StockReadiness() string {
	id, ok := w.TheArmoury()
	if !ok {
		return "You have nowhere to put them"
	}
	if w.Player.Location != id {
		place, _ := PlaceByID(id)
		return "They would have to be carried to " + place.Name
	}
	if w.Holding("arms") == 0 {
		return "You are not carrying any"
	}
	if w.Properties[id].Crates >= ArmouryHold {
		return "The room is full"
	}
	return ""
}

// StockArmoury moves what the player is carrying into the room, where it stops
// being their problem to carry and starts being merchandise.
func (w *World) StockArmoury() error {
	if reason := w.StockReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	id, _ := w.TheArmoury()
	prop := w.Properties[id]
	moved := min(w.Holding("arms"), ArmouryHold-prop.Crates)
	w.Player.Stock["arms"] -= moved
	prop.Crates += moved
	place, _ := PlaceByID(id)
	w.Log("Crates under "+place.Name, fmt.Sprintf("%d crates off your back and into the room. There are %d down there now.", moved, prop.Crates), "business")
	return nil
}

// ArmouryAttention is what the room draws in a day: the room itself, and how
// much is in it. Selling it down is the only way to make it quieter.
func (w *World) ArmouryAttention() int {
	id, ok := w.TheArmoury()
	if !ok || w.Properties[id].Crates <= 0 {
		return 0
	}
	return ArmouryHeat + w.Properties[id].Crates/ArmouryCrateHeat
}

// ArmsPrice is what an organization at war pays a crate. Nobody haggles when
// they are losing.
func (w *World) ArmsPrice() int {
	g := w.Good("arms")
	if g == nil {
		return 0
	}
	return g.Price * WarPremium / 100
}

// buyers is every organization currently fighting somebody, strongest first so
// the day is deterministic.
func (w *World) buyers() []*Faction {
	out := []*Faction{}
	for i := range w.Factions {
		f := &w.Factions[i]
		if w.fighting(f.ID) != nil {
			out = append(out, f)
		}
	}
	return out
}

// ArmouryDay is the war coming to the player's door with money. It sells to
// whoever is fighting, and it does not ask which side.
func (w *World) ArmouryDay() {
	id, ok := w.TheArmoury()
	if !ok {
		return
	}
	prop := w.Properties[id]
	if prop.Crates <= 0 {
		return
	}
	w.Player.Heat = min(100, w.Player.Heat+w.ArmouryAttention())

	price := w.ArmsPrice()
	place, _ := PlaceByID(id)
	for _, f := range w.buyers() {
		if prop.Crates <= 0 {
			return
		}
		want := 1 + int(w.WorldRandom()*3)
		want = min(want, prop.Crates)
		if f.Cash < price*want {
			want = f.Cash / max(1, price)
		}
		if want <= 0 {
			continue
		}
		paid := price * want
		f.Cash -= paid
		f.Power = min(100, f.Power+want*ArmsStrength)
		f.Goodwill = min(100, f.Goodwill+ArmouryGoodwill)
		prop.Crates -= want
		w.Earn(paid)
		w.Log("Somebody came to "+place.Name, fmt.Sprintf("%s took %d crates and left $%d. They are %d strong now, and there are %d crates left under the floor.", f.Name, want, paid, f.Power, prop.Crates), "business")

		// Arming one side is a thing the other side finds out about, sooner if
		// they have anybody worth having.
		if rival := w.fighting(f.ID); rival != nil && w.WorldRandom() < .45 {
			rival.Goodwill = max(-100, rival.Goodwill-8)
			if w.Reach() >= 2 {
				w.Log("Word gets back", fmt.Sprintf("%s knows where %s got them. Their standing with you falls to %+d.", rival.Name, f.Name, rival.Goodwill), "danger")
			}
		}
	}
}

// ArmouryFound is what a search costs somebody with a room like this. It is the
// worst thing in the game to be caught with, and it is meant to be.
func (w *World) ArmouryFound() (string, int, bool) {
	id, ok := w.TheArmoury()
	if !ok {
		return "", 0, false
	}
	prop := w.Properties[id]
	crates := prop.Crates
	prop.Armoury, prop.Crates = false, 0
	prop.Condition = max(0, prop.Condition-35)
	w.Player.Heat = min(100, w.Player.Heat+20)
	place, _ := PlaceByID(id)
	return place.Name, crates, true
}

// ArmouryDescription is the room and what is in it, for the interface.
func (w *World) ArmouryDescription() map[string]any {
	id, ok := w.TheArmoury()
	if !ok {
		return map[string]any{"held": false}
	}
	place, _ := PlaceByID(id)
	return map[string]any{
		"held": true, "place": place.Name, "crates": w.Properties[id].Crates,
		"capacity": ArmouryHold, "price": w.ArmsPrice(), "buyers": len(w.buyers()),
		"attention": w.ArmouryAttention(),
	}
}
