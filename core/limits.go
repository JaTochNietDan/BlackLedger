package core

import "fmt"

// The house limit. Gambling ran on two fixed lots — fifty dollars or three
// hundred at a table, a nickel or a dollar at a machine — so what a player put
// down was a button rather than a decision, and a room could not be run
// differently from the room next door.
//
// What you put down is now any amount you like up to the limit, and the limit
// belongs to whoever holds the room. A house that raises it takes bigger action
// and carries bigger nights; one that lowers it is a room serious money walks
// out of. Somebody else's house is somebody else's decision, which is the whole
// point of it being a business.

const (
	// LeastStake is the smallest bet a room will take. Below this a table is
	// not running a game, it is keeping somebody company.
	LeastStake = 1
	// HouseLimitFloor and HouseLimitCeiling bound what any holder can set. A
	// limit of nothing closes the room and a limit of everything is not a
	// limit; both are ways to make the rules say nothing.
	HouseLimitFloor   = 20
	HouseLimitCeiling = 5000
	// MachineShare is the part of a room's limit a machine will take. A bandit
	// is small money by design, whatever the tables are doing.
	MachineShare = 10
	// UsualStake and UsualPull are what goes down when the player presses the
	// button without naming a figure. A button that is offered as available and
	// then refused because nobody typed anything is a button that lies, so the
	// readiness above is asked about exactly this amount and the command falls
	// back to exactly this amount.
	UsualStake = 50
	UsualPull  = 5
)

// TableLimit is the most this room will take on one bet. A holder's own figure
// if they have set one, and otherwise what the room is worth: an address that
// earns more, and has more behind the tables, takes more.
func (w *World) TableLimit(id string) int {
	prop := w.Properties[id]
	if prop == nil {
		return 0
	}
	if prop.Limit > 0 {
		return min(HouseLimitCeiling, max(HouseLimitFloor, prop.Limit))
	}
	limit := prop.Income * 10
	if prop.Bankroll > 0 {
		limit = max(limit, prop.Bankroll/2)
	}
	return min(HouseLimitCeiling, max(HouseLimitFloor, limit))
}

// MachineLimit is the most a machine in this room will take, which is a tenth
// of what the tables will take and never less than a dollar.
func (w *World) MachineLimit(id string) int {
	return max(LeastStake, w.TableLimit(id)/MachineShare)
}

// StakeReadiness explains why an amount cannot be put down here, or returns "".
// One place, so a table, a wheel and a machine cannot come to different answers
// about the same money.
func (w *World) StakeReadiness(id string, amount, limit int) string {
	if amount < LeastStake {
		return "That is not a bet"
	}
	if amount > limit {
		return fmt.Sprintf("The house takes no more than $%d on one bet", limit)
	}
	if w.Player.Cash < amount {
		return "Not enough cash"
	}
	return ""
}

// LimitReadiness explains why the player cannot set the house limit here, or
// returns "".
func (w *World) LimitReadiness(id string) string {
	if !HasTables(id) && !HasMachines(id) {
		return "There is nothing to set a limit on here"
	}
	if !w.Own(id) {
		return "This is not your house to run"
	}
	return ""
}

// SetLimit is the holder deciding what their room will take.
func (w *World) SetLimit(id string, amount int) error {
	if reason := w.LimitReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if amount < HouseLimitFloor || amount > HouseLimitCeiling {
		return fmt.Errorf("a house limit runs from $%d to $%d", HouseLimitFloor, HouseLimitCeiling)
	}
	w.Properties[id].Limit = amount
	place, _ := PlaceByID(id)
	w.Log("The limit at "+place.Name, fmt.Sprintf("$%d on one bet from tonight. Bigger action pays better and loses worse; the floor knows which rooms will cover them.", amount), "business")
	return nil
}

// LimitDescription is what a room takes, for the interface.
func (w *World) LimitDescription(id string) map[string]any {
	return map[string]any{
		"limit": w.TableLimit(id), "machine": w.MachineLimit(id),
		"least": LeastStake, "floor": HouseLimitFloor, "ceiling": HouseLimitCeiling,
		"yours": w.Own(id),
	}
}

// TableStakeOrUsual is what a hand or a spin costs when the player has not said:
// the ordinary stake, never more than the house takes and never more than they
// are carrying.
func (w *World) TableStakeOrUsual(id string, amount int) int {
	if amount > 0 {
		return amount
	}
	return max(LeastStake, min(min(UsualStake, w.TableLimit(id)), w.Player.Cash))
}

// MachineStakeOrUsual is the same for a machine, which takes small money.
func (w *World) MachineStakeOrUsual(id string, amount int) int {
	if amount > 0 {
		return amount
	}
	return max(LeastStake, min(min(UsualPull, w.MachineLimit(id)), w.Player.Cash))
}

// HouseDescription is what the room the player is standing in will take, so the
// interface can offer a figure to type into rather than inventing a range of
// its own.
func (w *World) HouseDescription() map[string]any {
	id := w.Player.Location
	if !HasTables(id) && !HasMachines(id) {
		return map[string]any{"games": false}
	}
	return map[string]any{
		"games": true, "place": id,
		"limit": w.TableLimit(id), "machine": w.MachineLimit(id),
		"least": LeastStake, "usual": w.TableStakeOrUsual(id, 0), "pull": w.MachineStakeOrUsual(id, 0),
		"floor": HouseLimitFloor, "ceiling": HouseLimitCeiling,
		"yours": w.Own(id), "high": HighTableMoney,
	}
}
