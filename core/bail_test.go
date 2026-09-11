package core

import (
	"strings"
	"testing"
)

// Auditing the actions that had never had their label, gate and effect read
// together. Bail is offered at the precinct for anybody of the player's who is
// inside.
//
// Two things. `BailReadiness` exists, refuses when the player cannot afford it,
// and is called by `Bail` — but the button was added with an empty readiness,
// so it was never shown as unavailable. A player with no money saw a live
// button and got an error back instead of a reason on the button, which is the
// opposite of how every other action in the game behaves.
//
// And the description counts days while the ledger entry beside it does not:
// the effect already uses `plainly` to write "for the day still on them", and
// the description said "for the 1 days still on them".

func heldMan(t *testing.T, days int) (*World, *NPC) {
	t.Helper()
	w := New(23)
	w.Player.Location = "precinct"
	w.Player.Cash = 5000
	w.Player.Respect = 120
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Fatal("the player has no organization, so nobody answers to them")
	}
	w.NPCs = append(w.NPCs, NPC{ID: "held", Name: "Otto Reiss", Location: "precinct",
		Faction: w.PlayerOrganizationID(), Rank: RankSoldier,
		Held: w.Minute + days*1440})
	n := w.NPC("held")
	if !w.Inside(n) {
		t.Fatal("he is not inside")
	}
	return w, n
}

func bailAction(t *testing.T, w *World) Action {
	t.Helper()
	for _, a := range w.Actions("precinct") {
		if strings.HasPrefix(a.ID, "bail:") {
			return a
		}
	}
	t.Fatal("the precinct does not offer to bail anybody out")
	return Action{}
}

func TestBailIsRefusedOnTheButtonWhenItCannotBePaid(t *testing.T) {
	t.Parallel()
	w, _ := heldMan(t, 2)
	if a := bailAction(t, w); a.Disabled {
		t.Fatalf("a player who can afford it is refused: %q", a.Reason)
	}
	w.Player.Cash = 0
	a := bailAction(t, w)
	if !a.Disabled {
		t.Fatal("a player with no money is offered a live bail button")
	}
	if !contains(a.Reason, "cash") && !contains(a.Reason, "money") {
		t.Fatalf("the refusal does not say the money is the problem: %q", a.Reason)
	}
	// And the command still refuses, as it always did.
	if err := w.Bail("held"); err == nil {
		t.Fatal("bail was paid with no money")
	}
}

func TestBailCountsOneDayAsOneDay(t *testing.T) {
	t.Parallel()
	w, _ := heldMan(t, 1)
	a := bailAction(t, w)
	if contains(a.Detail, "1 days") {
		t.Errorf("the description counts one day as many: %q", a.Detail)
	}
	if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: a.ID, Target: "precinct"}); err != nil {
		t.Fatalf("bail failed: %v", err)
	}
	// The ledger entry has always got this right; the description should agree.
	last := w.History[len(w.History)-1]
	if contains(last.Text, "1 days") {
		t.Errorf("the ledger counts one day as many: %q", last.Text)
	}
}

// The other half, and the more important one: the exception must be exactly one
// action wide. Everything else aimed at somebody in a cell stays refused, which
// is the rule the sweep was written for.
func TestOnlyBailReachesIntoACell(t *testing.T) {
	t.Parallel()
	w, n := heldMan(t, 2)
	w.Player.Crew = []Crew{{n.ID, n.Name, 60}}
	w.Player.Location = "bar"
	n.Location = "bar"
	offered, refused := 0, 0
	for _, a := range w.Actions("bar") {
		if a.Subject != n.ID {
			continue
		}
		offered++
		if a.Disabled && contains(a.Reason, "held at Ward Street Station") {
			refused++
			continue
		}
		t.Errorf("%q is offered for a man the police are holding: disabled=%v %q",
			a.Label, a.Disabled, a.Reason)
	}
	if offered == 0 {
		t.Fatal("nothing was aimed at him at all, so this proved nothing")
	}
	t.Logf("%d of %d actions aimed at him are refused because he is inside", refused, offered)
	// And bail itself, at the precinct, is not.
	w.Player.Location = "precinct"
	if a := bailAction(t, w); a.Disabled {
		t.Fatalf("bail is still refused: %q", a.Reason)
	}
}

// The path nothing had ever walked: somebody of yours answers for what the
// police found, and later you buy them back.
//
// Both halves were on the never-played list — `bail:` as an action, `fall:` as
// a scene branch — and every test of bail above puts a man in a cell by writing
// the minute he comes out straight onto him. That is the state, not the road to
// it. The only thing in the game that puts one of the player's own people
// inside is the arrest scene's third answer, and nothing checked that what it
// leaves behind is what the precinct will sell back.
func TestGivingThemSomebodyAndBuyingThemBack(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event, w.Player.Cash, w.Player.Respect = nil, 40000, 200
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Fatal("the player has no organization, so nobody answers to them")
	}

	// Somebody signed on the way the game signs people on: standing in front of
	// you, and paid for.
	var signed *NPC
	for _, n := range w.NPCs {
		who := w.NPC(n.ID)
		w.Player.Location = who.Location
		if w.SignOnReadiness(who.ID) != "" {
			continue
		}
		if err := w.SignOn(who.ID); err != nil {
			t.Fatal(err)
		}
		signed = who
		break
	}
	if signed == nil {
		t.Skip("nobody in this city would sign on")
	}

	// The door. The arrest is raised by the world, and the answer goes through
	// the command layer like any other.
	w.Take(2, "what was in the back room")
	if w.Event == nil || w.Event.Kind != "arrest" {
		t.Fatal("the police did not come to the door")
	}
	fall := ""
	for _, c := range w.Event.Choices {
		if strings.HasPrefix(c.ID, "fall:") {
			fall = c.ID
		}
	}
	if fall == "" {
		t.Fatal("a player with somebody of their own is offered nobody to give them")
	}
	after, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: fall, Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = after

	inside := w.NPC(strings.TrimPrefix(fall, "fall:"))
	if inside == nil || !w.Inside(inside) {
		t.Fatal("the man handed over at the door is not inside")
	}
	if inside.Faction != w.PlayerOrganizationID() {
		t.Fatalf("%s still answers to %q, so the precinct will not sell them back",
			inside.Name, inside.Faction)
	}

	// And he is at the station, not standing where he was arrested. The line
	// under his name said "In a cell at Ward Street Station" while the room he
	// was taken from went on listing him among the people standing in it.
	if inside.Location != "precinct" {
		t.Fatalf("%s is in a cell and standing at %q", inside.Name, inside.Location)
	}
	for _, p := range w.PeopleHere(signed.Location) {
		if p.ID == inside.ID && signed.Location != "precinct" {
			t.Fatalf("%s is in a cell and in the room at %s", p.Name, signed.Location)
		}
	}
	here := false
	for _, p := range w.PeopleHere("precinct") {
		here = here || p.ID == inside.ID
	}
	if !here {
		t.Fatalf("%s is being held and is nowhere in the station", inside.Name)
	}

	// And the precinct sells them back. Read off the room's own card, because a
	// card nobody is offered is a road nobody can walk.
	w.Player.Location = "precinct"
	card := bailAction(t, w)
	if card.ID != "bail:"+inside.ID {
		t.Fatalf("the precinct offers %q rather than the man who was handed over", card.ID)
	}
	if card.Disabled {
		t.Fatalf("a player with $%d cannot buy them back: %s", w.Player.Cash, card.Reason)
	}
	bought, err := Execute(w, Command{Kind: card.ID, Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = bought
	if out := w.NPC(inside.ID); w.Inside(out) {
		t.Fatalf("%s is still inside after being bailed", out.Name)
	}
}
