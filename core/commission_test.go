package core

import (
	"encoding/json"
	"testing"
)

func petitioner(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(31)
	w.MigrateLivingWorld()
	var giver *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Faction != "" && n.Rank >= RankLieutenant && !n.Dead {
			n.Location = "bar"
			giver = n
			break
		}
	}
	if giver == nil {
		t.Fatal("no organization in this city had anybody who could ask for anything")
	}
	w.Player.Location = "bar"
	w.Player.Cash = 3000
	return w, giver
}

func TestNobodyAsksForAnythingWhereNobodyIs(t *testing.T) {
	w, giver := petitioner(t)
	giver.Location = "estate"
	if _, ok := w.AvailableCommission("bar"); ok {
		t.Fatal("an empty room offered work")
	}
	if w.CommissionReadiness("bar") == "" {
		t.Fatal("an empty room was ready to hand out work")
	}
}

func TestPeopleWhoWantYouDeadDoNotHandYouWork(t *testing.T) {
	w, giver := petitioner(t)
	if _, ok := w.AvailableCommission("bar"); !ok {
		t.Fatal("nobody offered anything to begin with")
	}
	w.faction(giver.Faction).Goodwill = -60
	if _, ok := w.AvailableCommission("bar"); ok {
		t.Fatal("a family that wants you dead offered you work")
	}
}

func TestWorkComesFromWhatTheOrganizationIsLivingThrough(t *testing.T) {
	w, giver := petitioner(t)
	f := w.faction(giver.Faction)
	rival := w.Rival(f.ID)
	if rival == nil {
		t.Fatal("the fixture had no quarrel in it")
	}

	// At war with somebody who holds ground: they want the ground hurt.
	c, ok := w.AvailableCommission("bar")
	if !ok || c.Kind != ObjectiveDamage {
		t.Fatalf("a family at war asked for %q", c.Kind)
	}
	if w.Properties[c.Target] == nil || w.HolderName(c.Target) != rival.Name {
		t.Fatalf("they asked about %q, which is not their rival's", c.Target)
	}

	// With that ground already wrecked, they want the money instead.
	w.Properties[c.Target].Condition = 20
	for _, id := range w.FamilyHoldings(rival.ID) {
		w.Properties[id].Condition = 20
	}
	c, ok = w.AvailableCommission("bar")
	if !ok || c.Kind != ObjectiveDrain || c.Target != rival.ID {
		t.Fatalf("a family with nothing left to break asked for %q against %q", c.Kind, c.Target)
	}
	if c.Baseline != rival.Cash {
		t.Fatalf("the baseline was $%d against $%d standing", c.Baseline, rival.Cash)
	}
}

func TestADrainCreditsOnlyWhatThePlayerTook(t *testing.T) {
	w, giver := petitioner(t)
	f := w.faction(giver.Faction)
	rival := w.Rival(f.ID)
	for _, id := range w.FamilyHoldings(rival.ID) {
		w.Properties[id].Condition = 20
	}
	if err := w.TakeCommission("bar"); err != nil {
		t.Fatal(err)
	}
	c := w.Live()[0]
	if c.Kind != ObjectiveDrain {
		t.Fatalf("expected a drain, got %q", c.Kind)
	}

	// Money that was never there is not money the player took.
	rival.Cash += 5000
	if w.Met(c) {
		t.Fatal("a rival getting richer counted as draining them")
	}
	rival.Cash = c.Baseline - c.Amount
	if !w.Met(w.Live()[0]) {
		t.Fatal("taking exactly what was asked did not count")
	}
}

func TestWorkIsPaidWhenItIsDoneAndChargedWhenItIsNot(t *testing.T) {
	w, giver := petitioner(t)
	f := w.faction(giver.Faction)
	f.Goodwill = 0
	c := Commission{ID: ID(), PatronID: f.ID, GiverName: giver.Name, Kind: ObjectiveStanding,
		Amount: 5, Brief: "Be somebody.", Due: w.Minute + CommissionWindow, Life: w.Life,
		Pay: 300, Respect: 2, Goodwill: 10, Penalty: 6}
	w.Commissions = append(w.Commissions, c)

	w.SettleCommissions()
	if len(w.Live()) != 1 {
		t.Fatal("unfinished work settled itself")
	}
	cash, respect := w.Player.Cash, w.Player.Respect
	w.Player.Respect = 20
	w.SettleCommissions()
	if len(w.Live()) != 0 {
		t.Fatal("finished work stayed open")
	}
	if w.Player.Cash != cash+300 || w.Player.Respect != 20+2 || f.Goodwill != 10 {
		t.Fatalf("cash %d, respect %d, goodwill %d", w.Player.Cash, w.Player.Respect, f.Goodwill)
	}
	if respect == 0 && w.Player.Respect == 0 {
		t.Fatal("nothing moved at all")
	}

	// And the other way: a deadline that passes costs standing, once.
	f.Goodwill = 0
	late := c
	late.ID, late.Amount = ID(), 9999
	w.Commissions = append(w.Commissions, late)
	w.Minute += CommissionWindow + 1
	w.SettleCommissions()
	if f.Goodwill != -6 {
		t.Fatalf("failing cost %d standing", -f.Goodwill)
	}
	w.SettleCommissions()
	if f.Goodwill != -6 {
		t.Fatal("failing charged twice")
	}
}

func TestOnlySoMuchWorkAtOnce(t *testing.T) {
	w, giver := petitioner(t)
	for i := 0; i < MaxCommissions; i++ {
		w.Commissions = append(w.Commissions, Commission{ID: ID(), PatronID: "nobody" + string(rune('a'+i)),
			GiverName: giver.Name, Kind: ObjectiveStanding, Amount: 9999, Due: w.Minute + CommissionWindow, Life: w.Life})
	}
	if w.CommissionReadiness("bar") == "" {
		t.Fatal("took on a fourth piece of work")
	}
	if err := w.TakeCommission("bar"); err == nil {
		t.Fatal("took on a fourth piece of work")
	}
}

func TestOneOrganizationAsksForOneThingAtATime(t *testing.T) {
	w, _ := petitioner(t)
	if err := w.TakeCommission("bar"); err != nil {
		t.Fatal(err)
	}
	if _, ok := w.AvailableCommission("bar"); ok {
		t.Fatal("the same family asked for a second thing while the first was outstanding")
	}
}

func TestNobodyInheritsAnObligation(t *testing.T) {
	w, _ := petitioner(t)
	w.TakeCommission("bar")
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if len(next.Live()) != 0 || len(next.Commissions) != 0 {
		t.Fatal("the next arrival owed somebody a favour")
	}
}

func TestTheInterfaceIsToldSomethingTrue(t *testing.T) {
	w, _ := petitioner(t)
	w.TakeCommission("bar")
	public := w.PublicCommissions()
	if len(public) != 1 {
		t.Fatalf("%d pieces of work shown", len(public))
	}
	entry := public[0]
	if entry["brief"] == "" || entry["progress"] == "" || entry["met"] != false {
		t.Fatalf("the interface was told %v", entry)
	}
	if entry["minutes_left"].(int) != CommissionWindow {
		t.Fatalf("%v minutes left on a commission just taken", entry["minutes_left"])
	}
}

// A commission has to survive being written to a save and read back, or a
// deadline stops meaning anything the moment the player closes the game.
func TestWorkSurvivesBeingWrittenDown(t *testing.T) {
	w, _ := petitioner(t)
	if err := w.TakeCommission("bar"); err != nil {
		t.Fatal(err)
	}
	before := w.Live()[0]
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var back World
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if len(back.Live()) != 1 {
		t.Fatalf("%d pieces of work came back", len(back.Live()))
	}
	after := back.Live()[0]
	if after.Kind != before.Kind || after.Target != before.Target || after.Amount != before.Amount ||
		after.Due != before.Due || after.Life != before.Life || after.Pay != before.Pay ||
		after.Respect != before.Respect || after.Goodwill != before.Goodwill || after.Penalty != before.Penalty ||
		after.Baseline != before.Baseline || after.GoodID != before.GoodID {
		t.Fatalf("what came back was not what went in:\n%+v\n%+v", before, after)
	}
}
