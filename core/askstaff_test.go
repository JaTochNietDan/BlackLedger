package core

import (
	"strings"
	"testing"
)

// The people behind your counters have names, a wage and a view of the street,
// and nothing to say. They notice a car parked across the road on their own
// schedule; a player standing in front of one should be able to ask.
//
// Everything the answer contains is a fact the city already holds: who has been
// in, what the trade is doing, and whether somebody has been looking the place
// over. Nothing is invented for it.

func employer(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	w.Player.Location = "laundry"
	return w, "laundry"
}

func TestYouCanAskYourOwnPeopleWhatTheyHaveSeen(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	hand := w.Properties[id].Hands[0]
	var ask *Action
	for _, a := range w.Actions(id) {
		if a.ID == "ask:"+hand {
			ask = &a
		}
	}
	if ask == nil {
		t.Fatalf("%s works behind this counter and there is no way to ask them anything", w.NPC(hand).Name)
	}
	if ask.Disabled {
		t.Fatalf("asking your own employee is refused: %s", ask.Reason)
	}
	if ask.Subject != hand {
		t.Fatalf("a question for %s is filed under %q", w.NPC(hand).Name, ask.Subject)
	}
	if err := w.apply(Command{Kind: "ask:" + hand, Target: id, RequestID: "askthecounter1"}); err != nil {
		t.Fatalf("asking was refused: %v", err)
	}
	last := w.History[len(w.History)-1]
	if !strings.Contains(last.Title, w.NPC(hand).Name) {
		t.Fatalf("the answer is filed under %q", last.Title)
	}
	if last.Text == "" {
		t.Fatal("they were asked and said nothing at all")
	}
}

// Somebody who works for a rival is not going to tell you anything, and nobody
// who works for nobody has a counter to watch.
func TestYouCanOnlyAskThePeopleWhoWorkForYou(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.EmptyChairs()
	theirs := w.Properties["butcher"].Hands[0]
	// Standing in their shop, in front of their employee.
	w.Player.Location = "butcher"
	for _, a := range w.Actions("butcher") {
		if a.ID == "ask:"+theirs && !a.Disabled {
			t.Fatalf("%s works for a rival and answers your questions", w.NPC(theirs).Name)
		}
	}
	// And refused if it is asked for anyway, which is the half the panel does
	// not cover: the question is only offered inside a room you hold, so
	// removing the check on who they work for failed nothing at all.
	if err := w.AskTheCounter(theirs); err == nil {
		t.Fatalf("%s works for a rival and answered", w.NPC(theirs).Name)
	}
	// And the answer comes from the room they actually stand in.
	w.Player.Location = id
	mine := w.Properties[id].Hands[0]
	if err := w.AskTheCounter(mine); err != nil {
		t.Fatalf("asking your own was refused: %v", err)
	}
	place, _ := PlaceByID(id)
	if !strings.Contains(w.History[len(w.History)-1].Text, place.Name) {
		t.Fatalf("the answer does not mention the room it came from: %q", w.History[len(w.History)-1].Text)
	}
}

// And what they say is drawn from what the city holds rather than invented: a
// place somebody is casing reads differently from a quiet one.
func TestWhatTheySayComesFromWhatIsActuallyHappening(t *testing.T) {
	t.Parallel()
	quiet, id := employer(t)
	watched, _ := employer(t)
	rival := watched.Factions[0].ID
	if rival == watched.PlayerOrganizationID() {
		rival = watched.Factions[1].ID
	}
	watched.Plots = append(watched.Plots, Plot{ID: ID(), Kind: "sabotage", Life: watched.Life,
		Due: watched.Minute + 4320, Actor: rival, Target: id, Strength: 35, Known: true})

	if err := quiet.AskTheCounter(quiet.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	if err := watched.AskTheCounter(watched.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	calm := quiet.History[len(quiet.History)-1].Text
	uneasy := watched.History[len(watched.History)-1].Text
	if calm == uneasy {
		t.Fatalf("a laundry somebody is looking over reads exactly like a quiet one: %q", calm)
	}
	if !strings.Contains(uneasy, "car") {
		t.Fatalf("somebody casing the place is not mentioned: %q", uneasy)
	}
}

// And asking is worth something afterwards. The word from your own counter is
// how a business that is about to be attacked finds out in time, and a place
// that saw them coming takes a third of the damage — but the walk that marked
// the plan known ran over copies of the plans rather than the plans, so the
// player was told about the car across the road and the world went on not
// knowing. Investigating recorded it; being told by your own people did not.
func TestAskingTheCounterIsWorthSomethingAfterwards(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "sabotage", Life: w.Life,
		Due: w.Minute + 4320, Actor: rival, Target: id, Strength: 35})
	if err := w.AskTheCounter(w.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	for _, p := range w.Plots {
		if p.Target != id {
			continue
		}
		if !p.Known {
			t.Fatal("the player was told about the car across the road and the world did not write it down")
		}
		return
	}
	t.Fatal("the plan against the laundry is gone from the world entirely")
}

// Somebody who has not been paid says so, and says it first. Three nights of
// rules about wages and the person standing behind the counter would tell you
// about the footfall and the stock and never mention the money.
func TestSomebodyWhoIsNotBeingPaidSaysSo(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	w.Properties[id].Unpaid = PatienceRunsOut + 2
	if err := w.AskTheCounter(w.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	said := w.History[len(w.History)-1].Text
	if !strings.Contains(said, "paid") {
		t.Fatalf("nobody there has been paid for nine nights and they talk about the weather: %q", said)
	}
}

// And somebody who is being paid does not bring it up, which is the ordinary
// case.
func TestSomebodyWhoIsPaidDoesNotMentionIt(t *testing.T) {
	t.Parallel()
	w, id := employer(t)
	if err := w.AskTheCounter(w.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	if said := w.History[len(w.History)-1].Text; strings.Contains(said, "paid") {
		t.Fatalf("somebody who is paid on time complains about the money: %q", said)
	}
}

// What asking is worth, as a number. A place that saw them coming takes a third
// of it, or nine less for every pair of hands standing in it, whichever is
// kinder — and until the walk above was fixed, asking bought none of that.
func TestAskingTheCounterIsWorthConditionWhenTheyCome(t *testing.T) {
	t.Parallel()
	asked, id := employer(t)
	silent, _ := employer(t)
	rival := asked.Factions[0].ID
	if rival == asked.PlayerOrganizationID() {
		rival = asked.Factions[1].ID
	}
	plot := Plot{ID: ID(), Kind: "sabotage", Life: asked.Life,
		Due: asked.Minute + 4320, Actor: rival, Target: id, Strength: 60}
	asked.Plots = append(asked.Plots, plot)
	silent.Plots = append(silent.Plots, plot)
	if err := asked.AskTheCounter(asked.Properties[id].Hands[0]); err != nil {
		t.Fatal(err)
	}
	asked.ResolveSabotage(asked.Plots[len(asked.Plots)-1])
	silent.ResolveSabotage(silent.Plots[len(silent.Plots)-1])
	warned, blind := asked.Properties[id].Condition, silent.Properties[id].Condition
	t.Logf("a laundry that asked its own people is at %d%% after the attack; one that did not is at %d%%",
		warned, blind)
	if warned <= blind {
		t.Fatalf("asking bought nothing when they came: %d%% against %d%%", warned, blind)
	}
}
