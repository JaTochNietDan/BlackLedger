package core

import "testing"

// Every action is put through `w.Pay(a.Cost)` before its effect runs — the
// prefixed ones included, because that line sits above the prefix chain and
// not inside the final switch. So an action whose effect pays a fee of its own
// AND whose button declares that fee as a Cost takes the money twice.
//
// This walks the seam for bail, where the button declares `days*BailDaily` and
// `Bail` pays `days*BailDaily` again. The pocket is the only witness worth
// trusting: a second world is put through the same hour without the bail, so
// what the hour itself costs is subtracted out and the fee stands alone.

func TestBailTakesTheMoneyOnce(t *testing.T) {
	t.Parallel()
	w, n := heldMan(t, 3)
	fee := 3 * BailDaily

	var bail *Action
	for i, a := range w.Actions("precinct") {
		if a.ID == "bail:"+n.ID {
			bail = &w.Actions("precinct")[i]
		}
	}
	if bail == nil {
		t.Fatal("no bail button for somebody of yours inside")
	}
	if bail.Cost != fee && bail.Asks != fee {
		t.Fatalf("the button names neither a cost nor a price of $%d (cost=%d asks=%d)", fee, bail.Cost, bail.Asks)
	}

	before := w.Player.Cash
	if err := w.apply(Command{Kind: "bail:" + n.ID, Target: "precinct", RequestID: "bailonce1"}); err != nil {
		t.Fatalf("bail: %v", err)
	}
	spent := before - w.Player.Cash

	// The same world, the same hour, no bail.
	quiet, _ := heldMan(t, 3)
	quietBefore := quiet.Player.Cash
	quiet.Advance(60)
	hour := quietBefore - quiet.Player.Cash

	if paid := spent - hour; paid != fee {
		t.Fatalf("bail of $%d took $%d out of the pocket (spent %d, the hour itself costs %d)", fee, paid, spent, hour)
	}
}

// The rule behind that one case. `w.Pay(a.Cost)` sits above the prefix chain,
// so it runs for prefixed actions too — and every effect the prefix chain
// reaches pays whatever it needs from inside itself. So none of these may
// declare a Cost. Any that does is charged for twice.
//
// The list is written out rather than derived: deriving it from the same file
// the rule is about would make the test agree with the code by construction.

var prefixChain = []string{
	"serve:", "pact:", "break:", "enquire:", "sign:", "share:", "smear:",
	"lend:", "lean:", "extend:", "forgive:", "bail:", "dismiss:",
}

func TestNothingInThePrefixChainDeclaresACost(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for _, a := range everyAction(t) {
		for _, prefix := range prefixChain {
			if len(a.ID) < len(prefix) || a.ID[:len(prefix)] != prefix {
				continue
			}
			seen[prefix] = true
			if a.Cost != 0 {
				t.Errorf("%s declares a cost of %d, and its effect pays as well, so the money goes out twice", a.ID, a.Cost)
			}
		}
	}
	if len(seen) < 6 {
		t.Fatalf("the sweep only reached %d of the %d prefixes, so it is not evidence about the rest", len(seen), len(prefixChain))
	}
}

// everyAction walks a well-supplied city and a family servant's city, and
// returns every button either of them is ever offered.
func everyAction(t *testing.T) []Action {
	t.Helper()
	var all []Action
	for _, seed := range []uint32{31, 47, 88, 103, 219} {
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
		w.Player.Contacts = 3
		w.Offshore, w.Player.Offshore = 4000, true
		w.Player.Car, w.Player.CarWear = 1, 90
		w.Player.Dress, w.Player.DressWear = 1, 90
		for _, id := range []string{"laundry", "garage", "casino"} {
			w.Properties[id].Owner = "player:1"
		}
		w.Incorporate()
		w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
		w.NPCs = append(w.NPCs, NPC{ID: "inside", Name: "Otto Reiss", Location: "precinct",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Held: w.Minute + 2*1440})
		w.Loans = append(w.Loans, Loan{})
		for i := range Locations {
			id := Locations[i].ID
			w.Player.Location = id
			all = append(all, w.Actions(id)...)
		}
	}
	return all
}
