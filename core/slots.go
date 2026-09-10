package core

import "fmt"

// The machines. A casino is a room with tables in it, and tables need a dealer,
// a floor and somebody to watch the floor. A machine needs a wall. That is why
// the mob put them everywhere there was a wall to spare — and why a poolhall or
// a bar can take money off people all day without being a casino at all.
//
// Everything here is decided by the core, the same as every other game: the
// reels are rolled from the player's own stream and the interface is told where
// they stopped.

// MachineStops is how many positions each reel has. The three reels carry the
// same strip, which is what a real bandit does and what makes the odds
// something a person could work out with a pencil.
const MachineStops = 20

// A Symbol is one face on the strip, how many of the twenty stops it takes, and
// what three of them return as a multiple of what was put in.
type Symbol struct {
	ID    string
	Face  string
	Stops int
	Pays  int
}

// The strip. Twenty stops, in the order a person would read them off the drum.
var reelStrip = []Symbol{
	{"seven", "7", 1, 100},
	{"bar", "BAR", 2, 50},
	{"bell", "BELL", 3, 25},
	{"plum", "PLUM", 4, 14},
	{"orange", "ORANGE", 4, 12},
	{"lemon", "LEMON", 4, 10},
	{"cherry", "CHERRY", 2, 25},
}

const (
	// TwoCherries and OneCherry are what the fruit pays on its own, which is
	// the whole reason a machine feels like it is nearly paying: most of what
	// comes back comes back a nickel at a time.
	TwoCherries = 5
	OneCherry   = 1
)

// slotStakes are what a machine takes. A bandit is not a table: it eats small
// money all day, which is exactly why there is one in every bar.
var slotStakes = []Stake{
	{ID: "nickel", Label: "Play the nickel machine", Amount: 5},
	{ID: "dollar", Label: "Play the dollar machine", Amount: 25},
}

func slotStake(id string) (Stake, bool) {
	for _, s := range slotStakes {
		if s.ID == id {
			return s, true
		}
	}
	return Stake{}, false
}

// SlotStakes is the list, for anything that offers them.
func SlotStakes() []Stake { return slotStakes }

// ReelStrip is the strip, for the interface to draw the drums with. It is the
// core's own list: a machine showing faces the core does not have is a machine
// showing somebody a lie.
func ReelStrip() []Symbol { return reelStrip }

// HasMachines reports whether a room has a bandit against the wall. Every
// casino has them, and so does anywhere people stand around with loose change
// in their pockets.
func HasMachines(id string) bool {
	place, ok := PlaceByID(id)
	if !ok {
		return false
	}
	return place.Type == "casino" || place.Kind == "poolhall" || place.Kind == "burlesque" || id == "bar"
}

// symbolAt turns a stop on the drum into the face standing in the window.
func symbolAt(stop int) Symbol {
	at := ((stop % MachineStops) + MachineStops) % MachineStops
	for _, s := range reelStrip {
		if at < s.Stops {
			return s
		}
		at -= s.Stops
	}
	return reelStrip[len(reelStrip)-1]
}

// MachinePays is what a line of three faces returns, as a multiple of the
// stake. Three of anything pays what the strip says it pays; cherries pay on
// their own, which is the whole of why a machine feels like it is nearly
// paying. Nothing else returns anything.
func MachinePays(line [3]Symbol) int {
	if line[0].ID == line[1].ID && line[1].ID == line[2].ID {
		return line[0].Pays
	}
	cherries := 0
	for _, s := range line {
		if s.ID == "cherry" {
			cherries++
		}
	}
	switch cherries {
	case 2:
		return TwoCherries
	case 1:
		return OneCherry
	}
	return 0
}

// MachineEdge is what the machine keeps out of every hundred put through it,
// worked out from the strip and the paytable rather than written down beside
// them. A number stated on a wall that does not follow the thing it describes
// is the interface making a promise the core has not agreed to.
func MachineEdge() int {
	total, paid := 0, 0
	for a := 0; a < MachineStops; a++ {
		for b := 0; b < MachineStops; b++ {
			for c := 0; c < MachineStops; c++ {
				total++
				paid += MachinePays([3]Symbol{symbolAt(a), symbolAt(b), symbolAt(c)})
			}
		}
	}
	return 100 - paid*100/total
}

// PullReadiness explains why the handle cannot be pulled, or returns "".
func (w *World) PullReadiness(id string, stake Stake) string {
	if !HasMachines(id) {
		return "There is no machine in here"
	}
	if w.Own(id) {
		return "You would be playing your own machine"
	}
	if w.Player.Cash < stake.Amount {
		return "Not enough cash"
	}
	return ""
}

// A Pull is where the reels stopped, kept so the interface can show the
// machine rather than a sentence about it.
type Pull struct {
	Place string
	Stake string
	Stops [3]int
	Pays  int
}

// PullHandle plays one line on a machine.
func (w *World) PullHandle(id, stakeID string) error {
	stake, ok := slotStake(stakeID)
	if !ok {
		return fmt.Errorf("no such machine")
	}
	if reason := w.PullReadiness(id, stake); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(stake.Amount); err != nil {
		return err
	}
	// The player is standing at this machine, so this is the player's stream.
	stops := [3]int{}
	for i := range stops {
		stops[i] = int(w.Random() * MachineStops)
		if stops[i] >= MachineStops {
			stops[i] = MachineStops - 1
		}
	}
	line := [3]Symbol{symbolAt(stops[0]), symbolAt(stops[1]), symbolAt(stops[2])}
	pays := MachinePays(line)
	w.Reels = &Pull{Place: id, Stake: stakeID, Stops: stops, Pays: pays}

	place, _ := PlaceByID(id)
	house := w.faction(w.Properties[id].Owner)
	returned := stake.Amount * pays
	if returned > 0 {
		w.Earn(returned)
	}
	net := returned - stake.Amount
	if house != nil {
		house.Cash = max(0, house.Cash-net)
	}
	w.tableAftermath(place.Name, house, net, stake.Amount)
	faces := line[0].Face + " · " + line[1].Face + " · " + line[2].Face
	if pays > 0 {
		w.Log("The machine at "+place.Name, fmt.Sprintf("%s. It pays %d to 1 and $%d drops into the tray.", faces, pays, returned), "business")
	} else {
		w.Log("The machine at "+place.Name, fmt.Sprintf("%s. Nothing, and the $%d is the machine's.", faces, stake.Amount), "business")
	}
	return nil
}

// MachineDescription is where the reels stopped, for the interface.
func (w *World) MachineDescription() map[string]any {
	faces := []map[string]any{}
	for _, s := range reelStrip {
		faces = append(faces, map[string]any{"id": s.ID, "face": s.Face, "stops": s.Stops, "pays": s.Pays})
	}
	out := map[string]any{
		"pulled": false, "strip": faces, "stops": MachineStops, "edge": MachineEdge(),
		"two_cherries": TwoCherries, "one_cherry": OneCherry,
	}
	if w.Reels == nil {
		return out
	}
	place, _ := PlaceByID(w.Reels.Place)
	stake, _ := slotStake(w.Reels.Stake)
	line := []string{}
	for _, stop := range w.Reels.Stops {
		line = append(line, symbolAt(stop).ID)
	}
	out["pulled"], out["place"], out["stake"] = true, place.Name, stake.Amount
	out["line"], out["pays"], out["won"] = line, w.Reels.Pays, w.Reels.Pays > 0
	return out
}

// sevenPays is the top line on the machine, for the description to quote
// without a second copy of the number living in a sentence.
func sevenPays() int {
	for _, s := range reelStrip {
		if s.ID == "seven" {
			return s.Pays
		}
	}
	return 0
}
