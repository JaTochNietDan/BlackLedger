package main

import (
	"blackledger/core"
	"strings"
	"testing"
)

// The director was shipped every family's cash as a bare integer inside the
// raw Faction record, and nothing in either prompt ever mentioned money. So a
// number that says nothing on its own was being sent with no way to read it.

func TestTheDirectorIsToldHowEachFamilyIsPlaced(t *testing.T) {
	w := core.New(41)
	money := organizationMoney(w)
	if len(money) == 0 {
		t.Fatal("no organization has a money picture")
	}
	for name, entry := range money {
		placed, ok := entry["placed"].(string)
		if !ok || placed == "" {
			t.Errorf("%s is not described as placed any way at all", name)
			continue
		}
		switch placed {
		case "comfortable", "getting by", "struggling", "cannot pay its people":
		default:
			t.Errorf("%s is %q, which is not one of the four states the prompt explains", name, placed)
		}
		if _, ok := entry["day_costs_them"]; !ok {
			t.Errorf("%s is described without what a day costs them", name)
		}
	}
}

// A family that is running down says how long it has. One living within its
// income is not counting days, and must not be given a number that invites a
// speaker to talk about three months of runway nobody is thinking about.
func TestOnlyAFamilyRunningDownCountsItsDays(t *testing.T) {
	w := core.New(41)
	var f *core.Faction
	for i := range w.Factions {
		if w.Factions[i].ID == "bellandi" {
			f = &w.Factions[i]
		}
	}
	if f == nil {
		t.Fatal("no bellandi")
	}
	for _, id := range w.FamilyHoldings("bellandi") {
		w.Properties[id].Owner = ""
	}
	f.Cash = 1000
	if _, counting := organizationMoney(w)[f.Name]["days_before_it_cannot_pay"]; !counting {
		t.Error("a family with no income and $1000 is not counting its days")
	}

	w2 := core.New(41)
	var g *core.Faction
	for i := range w2.Factions {
		if w2.Factions[i].ID == "bellandi" {
			g = &w2.Factions[i]
		}
	}
	g.Power, g.Cash = 4, 500
	if w2.FamilyIncome(g) <= w2.FamilyBill(g) {
		t.Skip("this family does not live within its income, so there is nothing to check")
	}
	if _, counting := organizationMoney(w2)[g.Name]["days_before_it_cannot_pay"]; counting {
		t.Error("a family earning more than it spends was given a countdown")
	}
}

// Both briefs have to explain the four words, or the field is sent to a model
// that has never been told what it means.
func TestBothBriefsExplainWhatTheMoneyStatesMean(t *testing.T) {
	for name, text := range map[string]string{"full": prompt, "focused": focusedPrompt} {
		if !strings.Contains(text, "organization_money") {
			t.Errorf("the %s brief never mentions organization_money", name)
		}
		for _, state := range []string{"comfortable", "getting by", "struggling", "cannot pay its people"} {
			if !strings.Contains(text, state) {
				t.Errorf("the %s brief never explains what %q means", name, state)
			}
		}
	}
}
