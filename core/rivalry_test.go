package core

import (
	"fmt"
	"testing"
)

// readyAttacker is a player who meets every published requirement for moving
// against a family, so tests exercise the outcome rather than the gate.
func readyAttacker(t *testing.T) *World {
	t.Helper()
	w := New(7)
	w.Player.Respect = 20
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
	w.Player.Location = "club"
	return w
}

func TestFamiliesHoldPropertyBothSidesCanLose(t *testing.T) {
	t.Parallel()
	w := New(1)
	bellandi := w.FamilyHoldings("bellandi")
	russo := w.FamilyHoldings("russo")
	if len(bellandi) == 0 || len(russo) == 0 {
		t.Fatalf("each family needs holdings to be pressured: bellandi=%v russo=%v", bellandi, russo)
	}
	for _, id := range append(append([]string{}, bellandi...), russo...) {
		if w.Properties[id].Income <= 0 {
			t.Fatalf("%s is a holding with no income, so damaging it would cost its family nothing", id)
		}
	}
	// The player's own property is never a sabotage target.
	w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
	if _, ok := w.SabotageTarget("laundry"); ok {
		t.Fatal("the player's own business was offered as a sabotage target")
	}
	if _, ok := w.SabotageTarget("garage"); ok {
		t.Fatal("an independent property was offered as a sabotage target")
	}
}

func TestSabotageRequirementsAreStatedBeforeTheyAreEnforced(t *testing.T) {
	t.Parallel()
	w := New(2)
	w.Player.Location = "club"
	if w.SabotageReadiness("club") == "" {
		t.Fatal("a new arrival should not be able to move against a family")
	}
	if err := w.Sabotage("club"); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	w.Player.Respect = 20
	if reason := w.SabotageReadiness("club"); reason == "" {
		t.Fatal("respect alone should not be enough without crew")
	}
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 20}}
	if reason := w.SabotageReadiness("club"); reason == "" {
		t.Fatal("disloyal crew should refuse this work")
	}
	w.Player.Crew[0].Loyalty = 80
	if reason := w.SabotageReadiness("club"); reason != "" {
		t.Fatal("a prepared player was still refused:", reason)
	}
	// The offered action carries the same reason the command would give.
	w.Player.Crew[0].Loyalty = 20
	for _, a := range w.Actions("club") {
		if a.ID == "sabotage" && (!a.Disabled || a.Reason == "") {
			t.Fatal("the interface offered an attack the rules would refuse")
		}
	}
}

func TestSuccessfulSabotageCostsTheFamilyRealStanding(t *testing.T) {
	t.Parallel()
	w := readyAttacker(t)
	before := *w.faction("bellandi")
	condition := w.Properties["club"].Condition
	// Force the success branch regardless of seed.
	w.RNG = 1
	for tries := 0; tries < 200; tries++ {
		probe := readyAttacker(t)
		probe.RNG = uint32(tries + 1)
		if probe.Random() < probe.sabotageChance(probe.faction("bellandi"), probe.OwnHands()) {
			w = readyAttacker(t)
			w.RNG = uint32(tries + 1)
			break
		}
	}
	if err := w.Sabotage("club"); err != nil {
		t.Fatal(err)
	}
	after := w.faction("bellandi")
	if w.Properties["club"].Condition >= condition {
		t.Fatal("a landed attack left the property undamaged")
	}
	if after.Power >= before.Power {
		t.Fatalf("family power did not fall: %d -> %d", before.Power, after.Power)
	}
	if after.Cash >= before.Cash {
		t.Fatalf("family cash did not fall: %d -> %d", before.Cash, after.Cash)
	}
	if after.Goodwill >= before.Goodwill {
		t.Fatal("attacking a family did not worsen standing with them")
	}
	if w.Player.Respect <= 20 {
		t.Fatal("a landed attack earned no respect")
	}
	// The family answers. A hit is scheduled but never exposed publicly here.
	answered := false
	for _, plot := range w.Plots {
		if plot.Actor == "bellandi" && plot.Life == w.Life {
			answered = true
		}
	}
	if !answered {
		t.Fatal("the family did not answer an attack on its own holding")
	}
}

func TestFailedSabotageInjuresWithoutDamagingTheHolding(t *testing.T) {
	t.Parallel()
	var w *World
	for tries := 0; tries < 400; tries++ {
		probe := readyAttacker(t)
		probe.RNG = uint32(tries + 1)
		if probe.Random() >= probe.sabotageChance(probe.faction("bellandi"), probe.OwnHands()) {
			w = readyAttacker(t)
			w.RNG = uint32(tries + 1)
			break
		}
	}
	if w == nil {
		t.Skip("no failing seed found in range")
	}
	condition, health := w.Properties["club"].Condition, w.Player.Health
	goodwill := w.faction("bellandi").Goodwill
	if err := w.Sabotage("club"); err != nil {
		t.Fatal(err)
	}
	if w.Properties["club"].Condition != condition {
		t.Fatal("a turned-away attack still damaged the property")
	}
	if w.Player.Health >= health {
		t.Fatal("a turned-away attack cost the player nothing")
	}
	if w.faction("bellandi").Goodwill >= goodwill {
		t.Fatal("a failed attack left the family's opinion unchanged")
	}
	// Being seen trying is enough to be answered.
	answered := false
	for _, plot := range w.Plots {
		if plot.Actor == "bellandi" && plot.Life == w.Life {
			answered = true
		}
	}
	if !answered {
		t.Fatal("the family did not answer an attempted attack")
	}
}

func TestFamiliesRebuildTheirHoldingsOverDays(t *testing.T) {
	t.Parallel()
	w := New(3)
	f := w.faction("russo")
	f.Peak = 58
	holding := w.FamilyHoldings("russo")[0]
	w.Properties[holding].Condition = 20
	f.Power = 20
	cash := f.Cash

	w.FamilyDay()
	if w.Properties[holding].Condition <= 20 {
		t.Fatal("a damaged holding was never repaired by its owner")
	}
	if f.Cash <= cash {
		t.Fatal("a family collected nothing from its holdings")
	}

	for day := 0; day < 40; day++ {
		w.FamilyDay()
	}
	if w.Properties[holding].Condition != 100 {
		t.Fatalf("holding never returned to full condition: %d", w.Properties[holding].Condition)
	}
	if f.Power != 58 {
		t.Fatalf("power did not recover to its peak: %d", f.Power)
	}
}

func TestDamagedHoldingsEarnTheirFamilyLess(t *testing.T) {
	t.Parallel()
	whole, damaged := New(4), New(4)
	holding := whole.FamilyHoldings("bellandi")[0]
	damaged.Properties[holding].Condition = 40
	before, damagedBefore := whole.faction("bellandi").Cash, damaged.faction("bellandi").Cash
	whole.FamilyDay()
	damaged.FamilyDay()
	if whole.faction("bellandi").Cash-before <= damaged.faction("bellandi").Cash-damagedBefore {
		t.Fatal("condition did not affect what a family collects")
	}
}

func TestIncitingARivalryNeedsContactsAndStanding(t *testing.T) {
	t.Parallel()
	w := New(11)
	w.Player.Location = "club"
	if w.InciteReadiness("club") == "" {
		t.Fatal("a newcomer with no network could set two families against each other")
	}
	if err := w.Incite("club"); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	w.Player.Contacts = 2
	w.Player.Respect = 8
	if reason := w.InciteReadiness("club"); reason != "" {
		t.Fatal("a connected player was still refused:", reason)
	}
	// The action carries the same reason the command would give.
	w.Player.Contacts = 0
	for _, a := range w.Actions("club") {
		if a.ID == "incite" && (!a.Disabled || a.Reason == "") {
			t.Fatal("the interface offered an incitement the rules would refuse")
		}
	}
}

func TestSuccessfulIncitementHardensTheOtherQuarrel(t *testing.T) {
	t.Parallel()
	var w *World
	for seed := 1; seed <= 200; seed++ {
		probe := New(uint32(seed))
		probe.Player.Contacts, probe.Player.Respect = 2, 8
		if probe.Random() >= .25 { // the branch where the story holds
			w = New(uint32(seed))
			w.Player.Contacts, w.Player.Respect = 2, 8
			break
		}
	}
	if w == nil {
		t.Skip("no succeeding seed found")
	}
	w.Player.Location = "club"
	before := w.Conflict("bellandi", "russo").Hostility
	if err := w.Incite("club"); err != nil {
		t.Fatal(err)
	}
	if got := w.Conflict("bellandi", "russo").Hostility; got <= before {
		t.Fatalf("the quarrel did not harden: %d -> %d", before, got)
	}
	// Pointing one family at another is not free of attention.
	if w.Player.Heat == 0 {
		t.Fatal("incitement drew no attention at all")
	}
}

func TestOnlyVisibleQuarrelsAreReported(t *testing.T) {
	t.Parallel()
	w := New(12)
	// Established families start uneasy but not yet newsworthy at war.
	for _, c := range w.PublicConflicts() {
		if c.State == "cold" {
			t.Fatal("cold relations were reported as news")
		}
		if len(c.Between) != 2 || c.Between[0] == "" {
			t.Fatal("a reported quarrel did not name both organizations")
		}
	}
	w.Antagonize("bellandi", "russo", 100)
	w.Conflict("bellandi", "russo").State = "war"
	found := false
	for _, c := range w.PublicConflicts() {
		if c.State == "war" {
			found = true
		}
	}
	if !found {
		t.Fatal("an open war was not visible to the city")
	}
}

func TestAnOrganizationHoldingNothingLosesItsStrength(t *testing.T) {
	t.Parallel()
	// Observed in a 1500-command campaign: Russo held nothing, had no income
	// and no cash, and still sat at power 58, because a landless organization
	// recovered toward the strength it had when it still owned property.
	w := New(401)
	f := w.faction("russo")
	f.Peak, f.Power = 58, 58
	for _, id := range w.FamilyHoldings("russo") {
		w.Properties[id].Owner = "independent"
	}
	// This used to fade on a timer, four points a day from the morning the last
	// business went. It does not any more, and the change is deliberate: an
	// organization with money in the safe holds together for exactly as long as
	// the money lasts, and comes apart at the first payday it cannot meet. The
	// rule the test was written for still holds — a landless family loses its
	// strength — but the money is now the reason rather than the calendar.
	before, reserves := f.Power, f.Cash
	w.FamilyDay()
	if w.faction("russo").Power != before {
		t.Fatalf("an organization that paid everybody in full lost strength anyway: %d from %d", w.faction("russo").Power, before)
	}
	if w.faction("russo").Cash >= reserves {
		t.Fatalf("an organization with no income did not spend anything on wages: %d from %d", w.faction("russo").Cash, reserves)
	}
	worst := 0
	for day := 0; day < 30; day++ {
		w.FamilyDay()
		worst = max(worst, w.faction("russo").Short)
	}
	if w.faction("russo").Power > 15 {
		t.Fatalf("a landless organization is still strong after a month: %d", w.faction("russo").Power)
	}
	// And it earned nothing all month, because it had nothing to earn from:
	// everything it had is gone on wages it could not keep meeting.
	if w.faction("russo").Cash != 0 {
		t.Fatalf("an organization with no holdings collected money: %d", w.faction("russo").Cash)
	}
	// Read as the worst run rather than as today: an operation that has come
	// apart to nobody costs nothing to run, so it is not short any more.
	if worst == 0 {
		t.Fatal("an organization that ran out of money never once missed a payday")
	}
}

func TestAnOrganizationWithPeopleCanFindSomewhereToStartAgain(t *testing.T) {
	t.Parallel()
	recovered := 0
	for i := uint32(1); i <= 200; i++ {
		w := New(i * 2654435761)
		f := w.faction("russo")
		for _, id := range w.FamilyHoldings("russo") {
			w.Properties[id].Owner = "independent"
		}
		f.Power = 40
		for turn := 0; turn < 40 && len(w.FamilyHoldings("russo")) == 0; turn++ {
			w.considerReestablish()
		}
		if len(w.FamilyHoldings("russo")) > 0 {
			recovered++
		}
	}
	t.Logf("of 200 landless organizations with people and strength, %d found somewhere to start again", recovered)
	if recovered == 0 {
		t.Fatal("an organization with people and strength can never recover any ground")
	}
}

func TestRecoveryNeverTakesWhatSomebodyHolds(t *testing.T) {
	t.Parallel()
	w := New(403)
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["garage"].Owner = "former:Alex Varga"
	w.Properties["garage"].Income = 24
	f := w.faction("russo")
	for _, id := range w.FamilyHoldings("russo") {
		w.Properties[id].Owner = "independent"
	}
	f.Power = 40
	for turn := 0; turn < 400; turn++ {
		w.considerReestablish()
	}
	if !w.Own("laundry") {
		t.Fatal("an organization took the player's business to re-establish itself")
	}
	for _, f := range w.Factions {
		for _, id := range w.FamilyHoldings(f.ID) {
			if w.Properties[id].Owner != f.ID {
				t.Fatal("holdings and ownership disagree")
			}
		}
	}
	// Bellandi's own premises are not available to Russo either.
	if w.Properties["club"].Owner != "bellandi" {
		t.Fatalf("a rival's holding was taken by re-establishment: club is %q", w.Properties["club"].Owner)
	}
}
