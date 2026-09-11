package core

import (
	"strings"
	"testing"
)

// "Maybe only some people know where someone is and we would have to be able to
// convince us to give them their location, especially if it's a rival family
// lead... Some people may not want to give us that information and even asking
// them could affect our reputation with them."

func askers(t *testing.T) (*World, *NPC, *NPC) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	w.Player.Location = "bar"
	var here, away *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) {
			continue
		}
		if here == nil && n.Location == "bar" {
			here = n
			continue
		}
		if away == nil && n.Location != "bar" && n.Faction != "" {
			away = n
		}
	}
	if here == nil || away == nil {
		t.Fatal("this city cannot stage a question")
	}
	return w, here, away
}

func TestSomebodyWhoHasNotSeenThemSaysSo(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	here.Trust = 80
	if w.Knows(here.ID, away.ID) {
		t.Fatalf("%s is across the city and %s is credited with knowing where", away.Name, here.Name)
	}
	if err := w.AskAbout(here.ID, away.ID); err != nil {
		t.Fatal(err)
	}
	if w.KnowsWhere(away.ID) {
		t.Fatalf("somebody who had not seen %s told the player where they are", away.Name)
	}
	if !strings.Contains(w.History[len(w.History)-1].Text, "has not seen") {
		t.Fatalf("the answer reads %q", w.History[len(w.History)-1].Text)
	}
}

func TestSomebodyWhoThinksNothingOfYouTellsYouNothing(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	// They work at the same counter, which is how somebody in the bar knows
	// where a colleague is without the player being able to see either of them.
	// A rival's butcher, so the player has no claim on the answer.
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Hands = []string{here.ID, away.ID}
	away.Faction, away.Location = "", "butcher"
	here.Trust = 0
	if !w.Knows(here.ID, away.ID) {
		t.Fatal("two people in the same room and one does not know where the other is")
	}
	if err := w.AskAbout(here.ID, away.ID); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(w.History[len(w.History)-1].Text, "rather not be asked") {
		t.Fatalf("somebody who thinks nothing of you answered: %q", w.History[len(w.History)-1].Text)
	}
}

func TestAskingSomebodyToGiveUpTheirOwnIsRemembered(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	here.Faction, here.Trust = away.Faction, 40
	away.Location = "docks"
	if !w.Knows(here.ID, away.ID) {
		t.Fatal("two people in the same family and one cannot say where the other is")
	}
	sore, trust := here.Sore, here.Trust
	if err := w.AskAbout(here.ID, away.ID); err != nil {
		t.Fatal(err)
	}
	if w.KnowsWhere(away.ID) {
		t.Fatalf("%s gave up one of their own family", here.Name)
	}
	if here.Sore <= sore {
		t.Fatalf("%s was asked to give up their own and holds nothing against you", here.Name)
	}
	if here.Trust >= trust {
		t.Fatalf("%s thinks the same of you after being asked", here.Name)
	}
	// And somebody who thinks a great deal of you will.
	loyal, mate, mark := askers(t)
	mate.Faction, mate.Trust = mark.Faction, 90
	mark.Location = "docks"
	if err := loyal.AskAbout(mate.ID, mark.ID); err != nil {
		t.Fatal(err)
	}
	if !loyal.KnowsWhere(mark.ID) {
		t.Fatalf("%s thinks the world of you and would not say where %s is", mate.Name, mark.Name)
	}
}

// What you are told is where they were, not where they are: a secondhand
// address is worth acting on and is not the same as having seen it.
func TestWhatYouAreToldIsAlreadyOld(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Hands = []string{here.ID, away.ID}
	away.Faction, here.Trust = "", 80
	away.Location = "butcher"
	if err := w.AskAbout(here.ID, away.ID); err != nil {
		t.Fatal(err)
	}
	seen, ok := w.LastSeen(away.ID)
	if !ok {
		t.Fatal("they said where and the player wrote nothing down")
	}
	if seen.When >= w.Minute {
		t.Fatalf("a secondhand address is as fresh as standing there: %d against %d", seen.When, w.Minute)
	}
	if !w.KnowsWhere(away.ID) {
		t.Fatal("an address four hours old is not worth acting on")
	}
}

// And there is no point asking about somebody you can already see.
func TestThereIsNoPointAskingAboutSomebodyYouCanSee(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	away.Location = w.Player.Location
	w.SeeTheRoom()
	if reason := w.AskAboutReadiness(here.ID, away.ID); reason == "" {
		t.Fatalf("the player is looking at %s and is offered a question about where they are", away.Name)
	}
}

// An action nobody can press is a feature in a file: this walks the same path
// the interface walks, including the second name the command has to carry.
func TestAskingAfterSomebodyIsOfferedAndCarriesBothNames(t *testing.T) {
	t.Parallel()
	w, here, away := askers(t)
	here.Trust, away.Rank = 80, RankLeader
	here.Faction, away.Location = "", "docks"
	// They share a counter, so the person in front of the player knows.
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Hands = []string{here.ID, away.ID}

	var offered *Action
	for _, a := range w.Actions(w.Player.Location) {
		// The id carries both names now, so two cards in one room cannot
		// share it: "about:leo:vittorio" is Leo asked where Vittorio is.
		if a.ID == "about:"+here.ID+":"+away.ID && a.Choice == away.ID {
			offered = &a
		}
	}
	if offered == nil {
		t.Fatalf("%s cannot be placed and there is no way to ask %s about them", away.Name, here.Name)
	}
	if offered.Disabled {
		t.Fatalf("asking is refused where it should work: %s", offered.Reason)
	}
	if offered.Subject != away.ID {
		t.Fatalf("a question about %s is filed under %q", away.Name, offered.Subject)
	}
	if err := w.apply(Command{Kind: "about:" + here.ID + ":" + away.ID, Choice: away.ID,
		Target: w.Player.Location, RequestID: "askafterthem1"}); err != nil {
		t.Fatalf("asking through the same path as everything else failed: %v", err)
	}
	if !w.KnowsWhere(away.ID) {
		t.Fatalf("%s told the player where %s is and the player did not write it down", here.Name, away.Name)
	}
}
