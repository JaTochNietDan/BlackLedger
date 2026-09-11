package core

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// An action that pays its own fee has to declare Cost of nothing, or the
// command layer takes the money a second time. The panel prints the price from
// Cost, so those actions showed no price at all: calling two families to a room
// takes $220 and its button said nothing about money.
//
// Asks is that fee for the panel only. Two things have to hold: it has to be
// there, and declaring it must not charge anybody twice.

func priced(t *testing.T, w *World, place, kind string) Action {
	t.Helper()
	w.Player.Location = place
	for _, a := range w.Actions(place) {
		if a.ID == kind {
			return a
		}
	}
	t.Fatalf("%s is not offered at %s", kind, place)
	return Action{}
}

func TestWorkThatPaysItsOwnFeeStillShowsAPrice(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Player.Contacts = 5
	w.Offshore = 5000
	w.Player.Offshore = false
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	// Something wrong in the back of one of them, so the remedy is offered;
	// something of the player's behind the pawnbroker's counter and something
	// on his shelf, so redeeming and the window are offered too. Each of these
	// is an action that pays its own fee and was going unchecked because the
	// city this test builds had no reason to offer it.
	w.Properties["garage"].Trouble = true
	w.Player.Car, w.Player.CarWear = 2, 20
	w.Tickets = append(w.Tickets, Ticket{
		Kind: "dress", Tier: 2, Wear: 10, Lent: 90,
		Due: w.Minute + PawnDays*1440, Life: w.Life})
	w.Player.Dress = 0
	w.Window = append(w.Window, Shelf{
		ID: "shelf-probe", Kind: "dress", Tier: 2, Wear: 25, Ask: 300, Lent: 120})
	w.Incorporate()
	w.ensureOfficials()
	w.Player.Retainers = append(w.Player.Retainers, "editor")
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})

	for _, c := range []struct {
		place, kind string
		fee         int
	}{
		{"herald", "spike", SpikeCost},
		{"herald", "puff", PuffCost},
		{"market", "offshore_access", AccessCost},
		{"laundry", "still", StillCost},
	} {
		a := priced(t, w, c.place, c.kind)
		if a.Cost != 0 {
			t.Errorf("%s declares a cost of %d, and the command layer would charge it on top of its own fee", c.kind, a.Cost)
		}
		if a.Asks != c.fee {
			t.Errorf("%s takes $%d and the button offers %d", c.kind, c.fee, a.Asks)
		}
	}
}

// The other half, and the reason Cost has to stay nothing: declaring a price
// must not take the money twice.
func TestDeclaringAPriceDoesNotChargeItTwice(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 40000, 200, 100
	w.Properties["laundry"].Owner = "player:1"
	a := priced(t, w, "laundry", "still")
	if a.Disabled {
		t.Skipf("a still cannot be built here: %s", a.Reason)
	}
	before := w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "still", Target: "laundry"})
	if err != nil {
		t.Fatalf("building it failed: %v", err)
	}
	spent := before - next.Player.Cash
	// The clock moves, so rent and income land too; the fee must be in there
	// once, not twice.
	// The clock moves while it is built, so a day's income and rent land on top
	// of the fee. What matters is that the fee itself went out once: comfortably
	// more than nothing, and nowhere near twice.
	if spent >= 2*StillCost {
		t.Errorf("a still costs $%d and pressing it took $%d — charged twice", StillCost, spent)
	}
	if spent < StillCost/2 {
		t.Errorf("a still costs $%d and pressing it took only $%d — not charged at all", StillCost, spent)
	}
	t.Logf("declared $%d, took $%d once the hours' other money had moved", a.Asks, spent)
}

// The regression guard, by name rather than by count. A first version asserted
// only that "at least twelve" actions named a fee, and silencing one still left
// twenty-four — a test that passed when I broke the code. These are the actions
// that pay their own way, so each must name a price and must declare no cost,
// or the engine would take the money on top of the fee.
// paysItsOwnWay is every action built with `asks(...)`, read out of the core's
// own source rather than written down here.
//
// It was a list of twenty-five, and the core had forty-two. Seventeen actions
// that pay their own fee were unguarded, including every one added in a night
// of work — the window at the pawnbroker, the boat at the pier, the night at a
// room you host. A missing entry is not a gap in coverage, it is an action that
// could declare a cost as well as paying its own fee and send the money out
// twice with nothing to say so.
//
// This is the fifth measure in this project caught sampling where it could have
// swept, and the shape is the same every time: a list a person wrote where the
// thing itself could have been asked.
func paysItsOwnWay(t *testing.T) []string {
	t.Helper()
	source, err := os.ReadFile("world.go")
	if err != nil {
		t.Fatalf("cannot read the core's own source: %v", err)
	}
	found, odd := map[string]bool{}, 0
	for _, m := range asksCall.FindAllStringSubmatch(string(source), -1) {
		arg := strings.TrimSpace(m[1])
		switch {
		case literalID.MatchString(arg):
			found[literalID.FindStringSubmatch(arg)[1]] = true
		case builtID.MatchString(arg):
			// `fmt.Sprintf("arms:weapon:%d", …)` and `"give:" + who` both name
			// a prefix, and the first card found under it stands for the rest.
			found[builtID.FindStringSubmatch(arg)[1]] = true
		default:
			odd++
		}
	}
	// The extraction has to be shown to be reading the file rather than
	// quietly matching nothing.
	if len(found) < 30 {
		t.Fatalf("only %d asks() call sites found in world.go, which cannot be right", len(found))
	}
	// And the number is pinned, because reading the asks() sites cannot see an
	// action that stops being one.
	//
	// Demonstrated: change `asks("window:"…)` to `add("window:"…, price, …)`
	// and the sweep loses sight of it entirely — while the command layer now
	// pays the price and `BuyFromWindow` pays it again, which is the exact
	// double charge this guard exists to prevent. A count that only goes up on
	// purpose catches the conversion the id list cannot.
	if len(found) != asksCallSites {
		t.Fatalf("world.go has %d actions that pay their own fee and this expects %d.\n"+
			"    If one was added, raise asksCallSites. If one was removed, check its handler\n"+
			"    does not still call Pay — an action that pays itself and declares a cost\n"+
			"    sends the money out twice.", len(found), asksCallSites)
	}
	if odd > 1 {
		t.Fatalf("%d asks() call sites whose id could not be read", odd)
	}
	out := []string{}
	for id := range found {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// asksCallSites is how many distinct actions in world.go pay their own fee.
const asksCallSites = 41

var (
	asksCall  = regexp.MustCompile(`(?s)asks\(\s*(.+?),`)
	literalID = regexp.MustCompile(`^"([^"]*)"\s*$`)
	builtID   = regexp.MustCompile(`^(?:fmt\.Sprintf\()?"([^"%]*)`)
	// A range said in words, for the cards whose stake is set on the felt.
	staked = regexp.MustCompile(`from \$\d+ to \$\d+`)
)

func TestEveryPricedActionKeepsItsPrice(t *testing.T) {
	t.Parallel()
	w := New(59)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 60000, 200, 100
	w.Player.Contacts = 5
	w.Player.Heat = 10
	w.Offshore = 5000
	w.Player.Offshore = false
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
		w.Properties[id].Staff = 0
		w.Properties[id].Supply = 0
		w.Properties[id].Condition = 60
	}
	// Something wrong in the back of one of them, so the remedy is offered;
	// something of the player's behind the pawnbroker's counter and something
	// on his shelf, so redeeming and the window are offered too. Each of these
	// is an action that pays its own fee and was going unchecked because the
	// city this test builds had no reason to offer it.
	w.Properties["garage"].Trouble = true
	w.Player.Car, w.Player.CarWear = 2, 20
	w.Tickets = append(w.Tickets, Ticket{
		Kind: "dress", Tier: 2, Wear: 10, Lent: 90,
		Due: w.Minute + PawnDays*1440, Life: w.Life})
	w.Player.Dress = 0
	w.Window = append(w.Window, Shelf{
		ID: "shelf-probe", Kind: "dress", Tier: 2, Wear: 25, Ask: 300, Lent: 120})
	w.Incorporate()
	w.ensureOfficials()
	w.News = append(w.News, Story{ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})
	for i := range w.Factions {
		w.Factions[i].Goodwill = 60
	}
	// Somebody of the player's standing in a room, so a share can be paid, and
	// somebody of theirs in a cell, so bail is offered.
	w.NPCs = append(w.NPCs,
		NPC{ID: "ownman", Name: "Otto Reiss", Location: "bar",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Role: "Yours", Trust: 40},
		NPC{ID: "cellman", Name: "Bruno Sala", Location: "precinct",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Role: "Yours",
			Trust: 40, Held: w.Minute + 2*1440})

	seen := map[string]Action{}
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			// The offered one, where the city offers the same card in more
			// than one room and refuses it in some of them.
			if had, seenIt := seen[a.ID]; !seenIt || (had.Disabled && !a.Disabled) {
				seen[a.ID] = a
			}
		}
	}
	pays := paysItsOwnWay(t)
	missing, checked := []string{}, 0
	for _, kind := range pays {
		// Three of these carry a person's id after the colon, so they are
		// matched on the prefix and the first one found stands for the rest.
		a, offered := seen[kind]
		if !offered && strings.HasSuffix(kind, ":") {
			for id, found := range seen {
				if strings.HasPrefix(id, kind) {
					a, offered = found, true
					break
				}
			}
		}
		if !offered {
			missing = append(missing, kind)
			continue
		}
		checked++
		// An action whose figure the player types states its price in the
		// field, not on the card: there is no one number to print before they
		// have said how much. Cost still has to be zero or the money leaves
		// twice, which is what this guard is actually for.
		// Three ways an action can name its price, and a fourth would be a
		// card that does not.
		//
		// Asks is the figure on the card. Sum is a field the player types into,
		// where there is no one number to print before they have said how much.
		// And the floor games name a range in words — "anything from $2 to
		// $500 a pull" — because the stake is set on the felt rather than on
		// the card, and the felt is the screen the player is looking at when
		// they choose it. That third one was not expressed here, so four
		// correct cards read as faults the first time this swept the whole
		// core instead of a list of twenty-five.
		named := a.Asks > 0 || a.Sum != nil || staked.MatchString(a.Detail)
		// A card that is refused has nothing to charge for. The pumps ask
		// nothing of a full tank and say so — "it reads 100 of 100, there is
		// nowhere for it to go" — and demanding a figure of that is demanding
		// a price for something nobody is selling. The cost check below still
		// applies to it, because a refused card carrying a cost is a double
		// charge waiting for the day it is offered.
		if !named && !a.Disabled {
			t.Errorf("%s pays its own fee and names no price: %q", kind, a.Detail)
		}
		if a.Sum != nil && (a.Sum.Least <= 0 || a.Sum.Most < a.Sum.Least) {
			t.Errorf("%s asks for a figure between %d and %d", kind, a.Sum.Least, a.Sum.Most)
		}
		if a.Cost != 0 {
			t.Errorf("%s pays its own fee and also declares a cost of %d, so the money would go out twice",
				kind, a.Cost)
		}
	}
	// Most of them have to be reachable in one city. The ones that are not are
	// named, so a whole area falling out of reach is visible rather than
	// absorbed into a threshold.
	if checked < len(pays)*2/3 {
		t.Errorf("only %d of %d priced actions were offered anywhere; this city is not exercising them: %v",
			checked, len(pays), missing)
	}
	if len(missing) > 0 {
		t.Logf("not offered in this city, so not checked: %v", missing)
	}
	t.Logf("checked %d actions that pay their own way", checked)
}
