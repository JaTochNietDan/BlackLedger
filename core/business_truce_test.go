package core

import "testing"

func TestBusinessTruceHasLimitedScopeAndExpires(t *testing.T) {
	w := New(27)
	w.Player.Cash = 300
	w.District = 1
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["garage"].Owner = "player:1"
	w.Plots = []Plot{{ID: "russo-business", Life: 1, Actor: "russo", Kind: "sabotage", Target: "garage"}, {ID: "russo-personal", Life: 1, Actor: "russo", Kind: "hit"}, {ID: "bellandi-business", Life: 1, Actor: "bellandi", Kind: "sabotage", Target: "laundry"}}
	w.OpenAudience("garage")
	start := w.Minute
	cash := w.Factions[1].Cash
	choice(t, &w, "business_truce")
	if w.Player.Cash != 200 || w.Factions[1].Cash != cash+100 || w.Minute != start || w.BusinessTruces["russo"] != start+1440 {
		t.Fatal("wrong payment or duration")
	}
	if len(w.Plots) != 2 || w.Plots[0].ID != "russo-personal" || w.Plots[1].ID != "bellandi-business" {
		t.Fatal("ceasefire altered personal/other-family plans")
	}
	w = w.Clone()
	if w.ActiveBusinessTruces()["russo"] != start+1440 {
		t.Fatal("saved agreement lost")
	}
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "bellandi" || w.Event.Target != "laundry" {
		t.Fatal("wrong family suppressed")
	}
	w.Event = nil
	before := w.Properties["garage"].Condition
	w.ResolveSabotage(Plot{Actor: "russo", Target: "garage", Strength: 35})
	if w.Properties["garage"].Condition != before {
		t.Fatal("sabotage violated ceasefire")
	}
	w.Minute = start + 1440
	if len(w.ActiveBusinessTruces()) != 0 {
		t.Fatal("expired agreement still presented as active")
	}
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "russo" {
		t.Fatal("demands never resumed")
	}
	w.ResolveSabotage(Plot{Actor: "russo", Target: "garage", Strength: 35})
	if w.Properties["garage"].Condition >= before {
		t.Fatal("sabotage permanently suppressed")
	}
}

func TestBusinessTrucePaymentAndNewLife(t *testing.T) {
	w := New(27)
	w.OpenAudience("club")
	if _, err := Execute(w, Command{Revision: w.Revision, Kind: "choice", Event: w.Event.ID, Choice: "business_truce"}); err == nil {
		t.Fatal("unaffordable ceasefire accepted")
	}
	if w.Player.Cash != 90 || len(w.BusinessTruces) != 0 {
		t.Fatal("failed agreement mutated original")
	}
	w.Player.Cash = 300
	choice(t, &w, "business_truce")
	first := w.BusinessTruces["bellandi"]
	w.OpenAudience("club")
	choice(t, &w, "business_truce")
	if w.BusinessTruces["bellandi"] != first {
		t.Fatal("renewal stacked duration")
	}
	w.Player.Alive = false
	w.Event = nil
	if len(w.ActiveBusinessTruces()) != 0 {
		t.Fatal("dead player's agreement active")
	}
	act(t, &w, "new_life", "")
	if len(w.BusinessTruces) != 0 {
		t.Fatal("new stranger inherited old agreement")
	}
}
