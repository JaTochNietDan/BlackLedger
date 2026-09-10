package core

import (
	"fmt"
	"testing"
)

func bomber(t *testing.T) *World {
	t.Helper()
	w := New(61)
	w.MigrateLivingWorld()
	w.Player.Location = "docks"
	w.Player.Cash = 20000
	w.Player.Health = 100
	w.Player.Respect = 40
	return w
}

func TestChargesComeOffABoatAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := bomber(t)
	if w.ChargeReadiness() != "" {
		t.Fatal("the waterfront refused:", w.ChargeReadiness())
	}
	for _, elsewhere := range []string{"bar", "market", "club", "room"} {
		w.Player.Location = elsewhere
		if w.ChargeReadiness() == "" {
			t.Fatalf("%s was selling explosives", elsewhere)
		}
		if err := w.BuyCharge(); err == nil {
			t.Fatalf("bought a charge at %s", elsewhere)
		}
	}
}

func TestHoldingOneIsTheMostConspicuousThingInTheCity(t *testing.T) {
	t.Parallel()
	w := bomber(t)
	if err := w.BuyCharge(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Charges != 1 {
		t.Fatalf("holding %d", w.Player.Charges)
	}
	heat := w.Player.Heat
	w.ChargeDay()
	if w.Player.Heat != heat+ChargeHeat {
		t.Fatalf("a day holding a charge drew %d attention", w.Player.Heat-heat)
	}
	// Two of them is twice the problem, and nobody carries three.
	w.BuyCharge()
	if w.ChargeReadiness() == "" {
		t.Fatal("bought a third")
	}
	heat = w.Player.Heat
	w.ChargeDay()
	if w.Player.Heat != heat+ChargeHeat*2 {
		t.Fatalf("two charges drew %d", w.Player.Heat-heat)
	}
	// And a search ends it.
	if !w.SeizeCharges() || w.Player.Charges != 0 {
		t.Fatal("a search did not find them")
	}
	if w.SeizeCharges() {
		t.Fatal("a search found charges that were already gone")
	}
}

func TestYouCannotPlantWhatYouDoNotHaveOrWhereItMakesNoSense(t *testing.T) {
	t.Parallel()
	w := bomber(t)
	if w.PlantReadiness("club") == "" {
		t.Fatal("planted a charge without one")
	}
	w.BuyCharge()
	w.Player.Location = "docks"
	if w.PlantReadiness("club") == "" {
		t.Fatal("planted a charge somewhere the player was not standing")
	}
	w.Player.Location = "club"
	if w.PlantReadiness("club") != "" {
		t.Fatal("could not plant at a rival casino:", w.PlantReadiness("club"))
	}
	w.Properties["club"].Owner = fmt.Sprintf("player:%d", w.Life)
	if w.PlantReadiness("club") == "" {
		t.Fatal("blew up their own business")
	}
	w.Properties["club"].Owner = "bellandi"
	w.Player.Respect = 0
	if w.PlantReadiness("club") == "" {
		t.Fatal("a nobody got under a rival's floor")
	}
}

func TestAChargeIsSpentWhetherItWorksOrNot(t *testing.T) {
	t.Parallel()
	for seed := uint32(1); seed <= 200; seed++ {
		w := bomber(t)
		w.RNG = seed * 2654435761
		w.BuyCharge()
		w.Player.Location = "club"
		if err := w.Plant("club"); err != nil {
			t.Fatal(err)
		}
		if w.Player.Charges != 0 {
			t.Fatalf("still holding %d after using one", w.Player.Charges)
		}
	}
}

func TestABlastTakesTheBuildingAndWhatWasInIt(t *testing.T) {
	t.Parallel()
	w := bomber(t)
	prop := w.Properties["club"]
	prop.Condition, prop.Supply, prop.Staff, prop.Bankroll, prop.Still = 100, 40, 5, 900, true
	before := w.faction("bellandi").Power
	w.detonate("club", "A charge went off.")
	if prop.Condition > 55 {
		t.Fatalf("condition %d after a charge", prop.Condition)
	}
	if prop.Supply != 0 || prop.Staff != 4 || !prop.Trouble {
		t.Fatalf("stock %d staff %d trouble %v", prop.Supply, prop.Staff, prop.Trouble)
	}
	if prop.Bankroll >= 900 || prop.Still {
		t.Fatalf("float %d still %v", prop.Bankroll, prop.Still)
	}
	if w.faction("bellandi").Power >= before {
		t.Fatal("the organization that owned it was not weakened")
	}
	if !w.hasRecord("An explosion at The Monarch") {
		t.Fatal("nothing was recorded")
	}
	if len(w.News) == 0 {
		t.Fatal("the paper did not carry it")
	}
}

func TestSomebodyIsUsuallyStandingThere(t *testing.T) {
	t.Parallel()
	const runs = 400
	killed := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := bomber(t)
		w.WorldRNG = seed * 2654435761
		before := 0
		for i := range w.NPCs {
			if w.NPCs[i].Dead {
				before++
			}
		}
		w.detonate("club", "A charge went off.")
		after := 0
		for i := range w.NPCs {
			if w.NPCs[i].Dead {
				after++
			}
		}
		if after > before {
			killed++
		}
	}
	rate := float64(killed) / runs
	if rate < .2 || rate > .5 {
		t.Fatalf("a charge killed somebody %.0f%% of the time", rate*100)
	}
	t.Logf("across %d explosions, somebody who worked there died in %d of them", runs, killed)
}

func TestTheCityUsesTheSameThingThePlayerDoes(t *testing.T) {
	t.Parallel()
	const runs, days = 300, 60
	bombed, atWar := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.DemolitionDay()
			w.Minute += 1440
		}
		if w.hasNewsKind("attack") {
			bombed++
		}
		for _, c := range w.Conflicts {
			if c.State == "war" {
				atWar++
				break
			}
		}
	}
	if bombed == 0 {
		t.Fatal("no organization in 300 cities ever destroyed anything")
	}
	if bombed == runs {
		t.Fatal("every city was bombed, so it is weather rather than an event")
	}
	t.Logf("%d of %d cities saw an organization destroy premises outright in %d days, with %d ending at war", bombed, runs, days, atWar)
}

func TestAnOrganizationThatHatesYouEnoughComesForYourBusiness(t *testing.T) {
	t.Parallel()
	const runs, days = 300, 60
	hit := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
		w.Properties["laundry"].Income = 20
		f := w.faction("bellandi")
		f.Cash, f.Goodwill = 30000, -80
		// Nobody to be at war with, so the only reason left is the player.
		for i := range w.Conflicts {
			w.Conflicts[i].State, w.Conflicts[i].Hostility = "cold", 0
		}
		condition := w.Properties["laundry"].Condition
		for day := 0; day < days; day++ {
			w.DemolitionDay()
		}
		if w.Properties["laundry"].Condition < condition {
			hit++
		}
	}
	if hit == 0 {
		t.Fatalf("an organization that hated the player never touched their business in %d cities", runs)
	}
	t.Logf("with one organization at -80 standing and no war to fight, the player's laundry was destroyed in %d of %d cities inside %d days", hit, runs, days)
}
