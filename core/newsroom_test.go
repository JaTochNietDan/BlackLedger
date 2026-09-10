package core

import "testing"

// A campaign standing at the Herald with a man on the top floor.
func editor(t *testing.T) *World {
	t.Helper()
	w := proprietor(t)
	own(w, "laundry")
	w.Player.Cash, w.Player.Respect = 20000, 60
	w.Player.Location = HeraldPlace
	w.ensureOfficials()
	if reason := w.RetainerReadiness("editor"); reason != "" {
		t.Fatal("nobody at the paper would talk:", reason)
	}
	if err := w.Retain("editor"); err != nil {
		t.Fatal(err)
	}
	return w
}

func TestNobodyAtThatPaperTakesYourCalls(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	w.Player.Location = HeraldPlace
	w.ensureOfficials()
	if w.TheEditor() {
		t.Fatal("a man nobody is paying was taking the player's calls")
	}
	for _, reason := range []string{w.SpikeReadiness(), w.PuffReadiness(), w.SmearReadiness("bellandi")} {
		if reason == "" {
			t.Fatal("the paper did something for nothing")
		}
	}
	// And it is arranged at the paper, not in the building with the others.
	w.Player.Location = CityHall
	if w.RetainerReadiness("editor") == "" {
		t.Fatal("the editor was retained from the market")
	}
	if w.RetainerReadiness("commissioner") == "This is not arranged here" {
		t.Fatal("the commissioner is no longer arranged in the building he works in")
	}
}

func TestAStoryThatNeverRunsDidNotHappen(t *testing.T) {
	t.Parallel()
	w := editor(t)
	w.Player.Heat = 40
	w.Attention = 50
	// The city files two things about the player this morning.
	w.Report("police", "RAID AT BLUEBIRD LAUNDRY", "Officers searched the premises.")
	w.Report("business", "TRADE STEADY IN THE DISTRICT", "Nothing much happened.")
	before := len(w.News)
	if len(w.Spikeable()) == 0 {
		t.Fatal("nothing in the paper was about the player")
	}
	attention, heat, cash := w.Attention, w.Player.Heat, w.Player.Cash
	if err := w.Spike(); err != nil {
		t.Fatal(err)
	}
	if len(w.News) != before-1 {
		t.Fatalf("the paper still has %d stories in it", len(w.News))
	}
	for _, s := range w.News {
		if s.Headline == "RAID AT BLUEBIRD LAUNDRY" {
			t.Fatal("the worst of them still ran")
		}
	}
	if w.Attention >= attention {
		t.Fatalf("the city is as interested as it was: %d against %d", w.Attention, attention)
	}
	if w.Player.Heat >= heat {
		t.Fatal("the police are as interested as they were")
	}
	if w.Player.Cash >= cash {
		t.Fatal("pulling a story was free")
	}
	// Yesterday's paper is already printed.
	w.Report("police", "ARRESTS EXPECTED", "b")
	w.Minute += 2000
	if len(w.Spikeable()) != 0 {
		t.Fatal("a story from two days ago was pulled")
	}
}

func TestSomebodyInThatBuildingEventuallyNotices(t *testing.T) {
	t.Parallel()
	// Every story pulled is a permanent fact about the player held by a man on
	// a weekly retainer, and over enough of them it comes apart.
	const runs = 400
	lost := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := editor(t)
		w.RNG = seed * 2654435761
		w.Report("police", "RAID AT BLUEBIRD LAUNDRY", "b")
		if err := w.Spike(); err != nil {
			t.Fatal(err)
		}
		if !w.TheEditor() {
			lost++
		}
	}
	if lost == 0 {
		t.Fatal("400 stories were pulled and nobody at that paper ever noticed")
	}
	if lost > runs/3 {
		t.Fatalf("the arrangement ended %d times in %d, which is not an arrangement", lost, runs)
	}
	t.Logf("of %d stories pulled, the arrangement came apart %d times", runs, lost)
}

func TestAStoryAboutSomebodyElseCostsThemAndTheWholeCity(t *testing.T) {
	t.Parallel()
	w := editor(t)
	f := w.faction("bellandi")
	holdings := w.FamilyHoldings("bellandi")
	if len(holdings) == 0 {
		t.Skip("they hold nothing")
	}
	custom, power, stories := w.Custom(holdings[0]), f.Power, len(w.News)
	w.RNG = 7
	if err := w.Smear("bellandi"); err != nil {
		t.Fatal(err)
	}
	if w.Custom(holdings[0]) >= custom {
		t.Fatalf("their trade is at %d, was %d", w.Custom(holdings[0]), custom)
	}
	if f.Power >= power {
		t.Fatal("the organization lost no strength")
	}
	if len(w.News) != stories+1 {
		t.Fatal("nothing appeared in the paper")
	}
	// The story counts toward the city's temperature like any other, which is
	// what makes this cost something even when it works perfectly.
	before := w.Attention
	w.Minute += 60
	w.ScrutinyDay()
	if w.Attention <= before {
		t.Fatal("a story in the paper about organized crime cooled the city")
	}
	// And not twice in the same week.
	if w.SmearReadiness("bellandi") == "" {
		t.Fatal("the paper carried the same man's arrangements twice")
	}
	if w.PuffReadiness() == "" {
		t.Fatal("the cooldown is not shared")
	}
}

func TestSometimesTheyFindOutWhoPaidForIt(t *testing.T) {
	t.Parallel()
	const runs = 400
	traced := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := editor(t)
		w.RNG = seed * 2654435761
		goodwill := w.faction("bellandi").Goodwill
		if err := w.Smear("bellandi"); err != nil {
			t.Fatal(err)
		}
		if w.faction("bellandi").Goodwill < goodwill {
			traced++
		}
	}
	if traced == 0 || traced > runs/2 {
		t.Fatalf("it was traced back %d times in %d", traced, runs)
	}
	t.Logf("of %d stories arranged about a rival, they worked out who paid for it %d times", runs, traced)
}

func TestTheCheapestStandingInThisCity(t *testing.T) {
	t.Parallel()
	w := editor(t)
	w.Player.Heat, w.Player.Respect = 30, 60
	stories := len(w.News)
	if err := w.Puff(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Respect != 60+PuffRespect || w.Player.Heat != 30-PuffHeat {
		t.Fatalf("respect %d, heat %d", w.Player.Respect, w.Player.Heat)
	}
	if len(w.News) != stories+1 {
		t.Fatal("nothing was printed")
	}
	// And it shares the cooldown with everything else the player puts in.
	if w.PuffReadiness() == "" || w.SmearReadiness("bellandi") == "" {
		t.Fatal("the paper would carry him again the same afternoon")
	}
}

func TestAnEditorIsAPersonLikeTheOthers(t *testing.T) {
	t.Parallel()
	w := editor(t)
	if !w.TheEditor() {
		t.Fatal("nobody is taking the money")
	}
	// Attention above his ceiling and he will not be seen with you.
	o, _ := OfficialByID("editor")
	w.Player.Heat = w.OfficialCeiling(o) + 5
	w.CityHallDay()
	if w.TheEditor() {
		t.Fatal("a man with a career stayed on the payroll of somebody at 80 attention")
	}
	// And he can be killed, like anybody with a name.
	w.Player.Heat = 0
	if err := w.Retain("editor"); err != nil {
		t.Fatal(err)
	}
	w.NPC("editor").Dead = true
	if w.TheEditor() {
		t.Fatal("a dead man was still pulling stories")
	}
}
