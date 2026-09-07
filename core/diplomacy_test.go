package core

import "testing"

func TestRussoAudienceCancelsOnlyRussoPlans(t *testing.T) {
	w := New(27)
	w.District = 1
	w.Player.Location = "garage"
	w.Player.Cash = 300
	w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: 1, Actor: "russo", Target: "garage", Due: 1000}, {ID: ID(), Kind: "hit", Life: 1, Actor: "bellandi", Due: 1100}}
	var err error
	w, err = Execute(w, Command{Kind: "audience", Target: "garage", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if w.Event == nil || w.Event.Actor != "russo" || w.Event.Speaker != "elena" || w.Minute != 525 {
		t.Fatal("audience routed incorrectly")
	}
	oldCash := w.Factions[1].Cash
	oldEvent := w.Event.ID
	choice(t, &w, "tribute")
	if w.Player.Cash != 150 || w.Factions[1].Cash != oldCash+150 || w.Factions[1].Goodwill != 8 || w.Factions[0].Goodwill != 0 {
		t.Fatal("tribute charged or credited wrong party")
	}
	if len(w.Plots) != 1 || w.Plots[0].Actor != "bellandi" {
		t.Fatal("other family's plan canceled")
	}
	if _, err = Execute(w, Command{Kind: "choice", Event: oldEvent, Choice: "tribute", Revision: w.Revision}); err == nil {
		t.Fatal("tribute replay accepted")
	}
}
func TestAudienceFavorIsOptionalAndUsesNormalJobRules(t *testing.T) {
	for _, decision := range []string{"accept", "decline"} {
		w := New(27)
		w.OpenAudience("garage")
		start := w.Minute
		w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: 1, Actor: "russo", Target: "garage", Due: 1000}}
		choice(t, &w, "work")
		if w.Minute != start || w.Player.Cash != 90 || w.Event == nil || w.Event.Kind != "proposal" || w.Event.Beneficiary != "russo" || w.Event.Speaker != "elena" || len(w.Plots) != 1 {
			t.Fatal("hearing offer changed world or lost affiliation")
		}
		choice(t, &w, decision)
		if decision == "accept" {
			if w.Minute != start+45 || w.Player.Cash != 165 || w.Factions[1].Goodwill != 6 || w.Factions[0].Goodwill != -3 {
				t.Fatal("favor bypassed normal rewards/standing")
			}
		} else if w.Minute != start || w.Player.Cash != 90 || w.Factions[1].Goodwill != 0 {
			t.Fatal("declined favor changed economy")
		}
	}
}
func TestLegacyAudienceAndUnaffordableTribute(t *testing.T) {
	w := New(27)
	w.OpenAudience("club")
	w.Event.Actor = ""
	before := w.Clone()
	if _, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "tribute", Revision: w.Revision}); err == nil {
		t.Fatal("unaffordable tribute accepted")
	}
	if w.Player.Cash != before.Player.Cash || w.Event == nil {
		t.Fatal("failed payment mutated original")
	}
	w.Player.Cash = 150
	w.Factions[0].Goodwill = 99
	choice(t, &w, "tribute")
	if w.Factions[0].Goodwill != 100 || w.Player.Cash != 0 {
		t.Fatal("legacy audience or standing cap failed")
	}
}
