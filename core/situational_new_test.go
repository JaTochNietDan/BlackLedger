package core

import (
	"fmt"
	"testing"
)

func offers(w *World) map[string]SituationalOperation {
	out := map[string]SituationalOperation{}
	for _, s := range w.SituationalOperations() {
		out[s.ID] = s
	}
	return out
}

func TestAQuietLifeOffersNothing(t *testing.T) {
	w := New(151)
	w.MigrateLivingWorld()
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "cold", 10
	}
	w.Grudges, w.Commissions = nil, nil
	for i := range w.Factions {
		w.Factions[i].Goodwill = 0
	}
	for _, id := range []string{"consignment", "grievance", "obligation", "warning_off"} {
		if _, ok := offers(w)[id]; ok {
			t.Fatalf("a player with nothing going on was offered %s", id)
		}
	}
}

func TestCratesAndAWarIsASituation(t *testing.T) {
	w := New(151)
	w.MigrateLivingWorld()
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Player.Location, w.Player.Cash = "laundry", 20000
	if err := w.BuildArmoury("laundry"); err != nil {
		t.Fatal(err)
	}
	w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
	if _, ok := offers(w)["consignment"]; ok {
		t.Fatal("an empty room was a situation")
	}
	w.Properties["laundry"].Crates = 30
	s, ok := offers(w)["consignment"]
	if !ok {
		t.Fatal("thirty crates and a war was not a situation")
	}
	if !containsName(s.Because, "30 crates") {
		t.Fatalf("the fact was %q", s.Because)
	}
}

func TestAQuarrelYouKnowAboutIsASituation(t *testing.T) {
	w := New(151)
	w.MigrateLivingWorld()
	people := w.People()
	w.Resent(people[3].ID, people[4].ID, GrudgeCap, "an old debt")
	w.Player.Contacts = 0
	if _, ok := offers(w)["grievance"]; ok {
		t.Fatal("a player with no contacts was told about somebody else's quarrel")
	}
	w.Player.Contacts = 3
	s, ok := offers(w)["grievance"]
	if !ok {
		t.Fatal("a live quarrel the player could hear about was not a situation")
	}
	if !containsName(s.Because, people[3].Name) || !containsName(s.Because, "old debt") {
		t.Fatalf("the fact was %q", s.Because)
	}
}

func TestWorkYouPromisedIsASituation(t *testing.T) {
	w, giver := petitioner(t)
	if _, ok := offers(w)["obligation"]; ok {
		t.Fatal("a player who owed nobody anything was reminded of it")
	}
	if err := w.TakeCommission("bar"); err != nil {
		t.Fatal(err)
	}
	s, ok := offers(w)["obligation"]
	if !ok {
		t.Fatal("outstanding work was not a situation")
	}
	if !containsName(s.Because, giver.Name) {
		t.Fatalf("the fact was %q", s.Because)
	}
}

func TestAnEnemyPastTalkingIsASituation(t *testing.T) {
	w := New(151)
	w.MigrateLivingWorld()
	w.Factions[0].Goodwill = -54
	if _, ok := offers(w)["warning_off"]; ok {
		t.Fatal("an organization still complaining was treated as past talking")
	}
	w.Factions[0].Goodwill = -55
	s, ok := offers(w)["warning_off"]
	if !ok {
		t.Fatal("an organization at -55 was not a situation")
	}
	if !containsName(s.Because, w.Factions[0].Name) {
		t.Fatalf("the fact was %q", s.Because)
	}
}

func TestEveryNewSituationHasAPriceALabelAndAnOutcome(t *testing.T) {
	effects := SituationalEffects()
	for _, id := range []string{"consignment", "grievance", "obligation", "warning_off"} {
		effect, ok := effects[id]
		if !ok || effect.Reward <= 0 || effect.Minutes <= 0 {
			t.Fatalf("%s had no price: %+v", id, effect)
		}
		if operationLabel(id) == "Take the work" {
			t.Fatalf("%s had no label of its own", id)
		}
		if operationOutcome(id) == "You completed the arrangement." {
			t.Fatalf("%s had no outcome of its own", id)
		}
	}
}

func TestASituationsPriceMatchesWhatItPays(t *testing.T) {
	// The catalog and the offer have to agree, or a proposal validates against
	// one number and pays another.
	w := New(151)
	w.MigrateLivingWorld()
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	w.Player.Location, w.Player.Cash = "laundry", 20000
	w.BuildArmoury("laundry")
	w.Properties["laundry"].Crates = 30
	w.Conflicts[0].State, w.Conflicts[0].Hostility = "war", 80
	w.Player.Contacts = 3
	people := w.People()
	w.Resent(people[3].ID, people[4].ID, GrudgeCap, "an old debt")
	w.Factions[0].Goodwill = -70

	effects := SituationalEffects()
	for _, s := range w.SituationalOperations() {
		if effects[s.ID] != s.Effect {
			t.Fatalf("%s is offered as %+v and paid as %+v", s.ID, s.Effect, effects[s.ID])
		}
		if s.Because == "" {
			t.Fatalf("%s was offered without a reason", s.ID)
		}
	}
}
