package core

import "testing"

// The undertaker's trade is the city's dead.
//
// The room went in with a counter that counts the week's funerals and a line in
// the brief calling it the one trade whose custom is made by everybody else's
// work. Nothing did it. The parlour earned the same whether the district buried
// two people that week or twenty, which is a sentence describing a consequence
// with no code behind it — the fault this project keeps finding, written this
// time by the tick before.
func theParlour(t *testing.T) (*World, string) {
	t.Helper()
	w := New(53)
	w.Event, w.District, w.Player.Cash = nil, 9, 40000
	id := w.TheParlour()
	if id == "" {
		t.Skip("this city has no undertaker")
	}
	return w, id
}

func TestEveryFuneralIsWorkForSomebody(t *testing.T) {
	t.Parallel()
	w, parlour := theParlour(t)
	before := w.Custom(parlour)

	// Somebody with a name, a wage and people to answer to. Found rather than
	// built, so this measures the city's own dead.
	buried := 0
	for _, n := range w.NPCs {
		who := w.NPC(n.ID)
		if who.Dead || who.Faction == "" {
			continue
		}
		if !w.Kill(who.ID, "Shot on the steps.") {
			continue
		}
		buried++
		if buried >= 4 {
			break
		}
	}
	if buried == 0 {
		t.Skip("nobody in this city answers to anybody")
	}
	after := w.Custom(parlour)
	if after <= before {
		t.Fatalf("%d funerals and the parlour went from %d%% to %d%%", buried, before, after)
	}
	if want := before + buried*BurialTrade; after != want {
		t.Fatalf("%d funerals at %d each should read %d%% and reads %d%%",
			buried, BurialTrade, want, after)
	}
}

// And a pauper is not a customer. Somebody with no name in this city, nobody to
// send a card to and nothing in their pockets is buried out of the parish's
// money, which is not money: the room does the work and is not better off.
func TestBuryingSomebodyNobodyIsPayingForIsNotTrade(t *testing.T) {
	t.Parallel()
	w, parlour := theParlour(t)
	before := w.Custom(parlour)

	pauper := w.AddCivilian()
	if pauper == nil {
		t.Skip("nobody can be added to this city")
	}
	pauper.Faction, pauper.Rank, pauper.Purse = "", RankAssociate, 0
	if !w.Kill(pauper.ID, "Found in the road.") {
		t.Fatal("they could not be killed")
	}
	if after := w.Custom(parlour); after != before {
		t.Fatalf("a pauper's funeral took the parlour from %d%% to %d%%", before, after)
	}
}

// A parlour of the player's is where their own city's business ends up, the
// same rule the garage is under.
func TestYourOwnParlourTakesTheCitysFunerals(t *testing.T) {
	t.Parallel()
	w, first := theParlour(t)
	mine := ""
	for _, l := range Locations {
		if l.Kind == "undertaker" && l.ID != first {
			mine = l.ID
		}
	}
	if mine == "" {
		// One parlour in this city, so the rule cannot be told apart from the
		// fallback. What can be checked is that owning it does not change which
		// room takes the work.
		w.Properties[first].Owner = "player:1"
		if w.TheParlour() != first {
			t.Fatal("the city's only parlour stopped being the one that buries people")
		}
		return
	}
	w.Properties[mine].Owner = "player:1"
	if got := w.TheParlour(); got != mine {
		t.Fatalf("the funerals go to %s and %s belongs to the player", got, mine)
	}
}
