package core

import (
	"strings"
	"testing"
)

// Burying one of your own.
//
// A man who died working for the player left a line in the log and nothing
// else: the wage stopped, his name left the books, and everybody still on the
// payroll carried on as though the week had been ordinary. This is the
// undertaker's other half — the trade earns on strangers, and this is what it
// costs the player and what it buys.
func aDeadHandOfYours(t *testing.T) (*World, string, string) {
	t.Helper()
	w := New(53)
	w.Event, w.District, w.Player.Cash, w.Player.Respect = nil, 9, 40000, 200
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Skip("the player has no organization, so nobody answers to them")
	}
	// Two signed on: one to bury and one to watch it happen.
	signed := []string{}
	for _, n := range w.NPCs {
		who := w.NPC(n.ID)
		w.Player.Location = who.Location
		if w.SignOnReadiness(who.ID) != "" {
			continue
		}
		if err := w.SignOn(who.ID); err != nil {
			t.Fatal(err)
		}
		signed = append(signed, who.ID)
		if len(signed) == 2 {
			break
		}
	}
	if len(signed) < 2 {
		t.Skip("not enough people in this city would sign on")
	}
	parlour := w.TheParlour()
	if parlour == "" {
		t.Skip("this city has no undertaker")
	}
	if !w.Kill(signed[0], "Shot on the steps.") {
		t.Fatal("they could not be killed")
	}
	return w, parlour, signed[0]
}

func TestBuryingOneOfYourOwnIsSomethingEverybodyElseSees(t *testing.T) {
	t.Parallel()
	w, parlour, dead := aDeadHandOfYours(t)
	w.Player.Location, w.Event = parlour, nil

	var card Action
	for _, a := range w.Actions(parlour) {
		if a.ID == "funeral:"+dead {
			card = a
		}
	}
	if card.ID == "" {
		t.Fatal("a funeral director's does not offer to bury somebody of yours")
	}
	if card.Disabled {
		t.Fatalf("the funeral is refused: %s", card.Reason)
	}
	if !strings.Contains(card.Detail, "$") {
		t.Errorf("a card that costs money says nothing about money: %q", card.Detail)
	}

	watching := map[string]int{}
	for _, n := range w.OwnPeople() {
		watching[n.ID] = n.Trust
	}
	standing, cash := w.Player.Respect, w.Player.Cash
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: card.ID})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash >= cash {
		t.Fatalf("the funeral cost nothing: $%d then, $%d now", cash, next.Player.Cash)
	}
	if next.Player.Respect <= standing {
		t.Fatalf("standing went from %d to %d", standing, next.Player.Respect)
	}
	lifted := 0
	for _, n := range next.OwnPeople() {
		if was, knew := watching[n.ID]; knew && n.Trust > was {
			lifted++
		}
	}
	if lifted == 0 {
		t.Fatal("nobody still on the books thought any better of it")
	}
	// And it is done. A man is not buried twice.
	for _, a := range next.Actions(parlour) {
		if a.ID == card.ID {
			t.Fatal("the same funeral is still on offer")
		}
	}
}

// And a week later nobody is arranging anything.
func TestAFuneralNobodyArrangedInTimeIsNotArrangedAtAll(t *testing.T) {
	t.Parallel()
	w, parlour, dead := aDeadHandOfYours(t)
	w.Advance(FuneralWindow + 1440)
	w.Player.Location, w.Event = parlour, nil
	if reason := w.FuneralReadiness(parlour, dead); reason == "" {
		t.Fatal("somebody dead a week is still being buried by the people who knew them")
	}
	for _, a := range w.Actions(parlour) {
		if a.ID == "funeral:"+dead {
			t.Fatal("the card is still offered for somebody nobody is burying")
		}
	}
}

// A parlour of the player's own does it for what the plot and the notices cost,
// because the cars and the box are already theirs.
func TestYourOwnParlourBuriesYourOwnForLess(t *testing.T) {
	t.Parallel()
	w, parlour, _ := aDeadHandOfYours(t)
	full := w.FuneralFee(parlour)
	w.Properties[parlour].Owner = "player:1"
	if mine := w.FuneralFee(parlour); mine >= full {
		t.Fatalf("a funeral costs $%d at somebody else's parlour and $%d at your own", full, mine)
	}
}
