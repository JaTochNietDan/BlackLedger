package core

import "testing"

func quietCity(t *testing.T) *World {
	t.Helper()
	w := New(307)
	w.MigrateLivingWorld()
	w.News = nil
	w.Attention = 0
	return w
}

func TestAQuietCityCoolsAndALoudOneDoesNot(t *testing.T) {
	t.Parallel()
	w := quietCity(t)
	w.Attention = 40
	for i := 0; i < 10; i++ {
		w.Minute += 1440
		w.ScrutinyDay()
	}
	if w.Scrutiny() != 30 {
		t.Fatalf("ten quiet days left the city at %d", w.Scrutiny())
	}
	// It cannot fall below nothing.
	for i := 0; i < 200; i++ {
		w.Minute += 1440
		w.ScrutinyDay()
	}
	if w.Scrutiny() != 0 {
		t.Fatalf("the city settled at %d", w.Scrutiny())
	}
}

func TestWhatThePaperCarriesIsWhatTheCityCounts(t *testing.T) {
	t.Parallel()
	w := quietCity(t)
	w.Minute = 5000
	w.Report("killing", "SOMEBODY KILLED", "It happened.")
	w.Report("business", "A SHOP CHANGES HANDS", "It happened.")
	w.ScrutinyDay()
	if w.Scrutiny() != scrutinyWeight["killing"] {
		t.Fatalf("a killing and a shop were worth %d", w.Scrutiny())
	}
	// Yesterday's news is not today's.
	w.Minute += 2880
	before := w.Scrutiny()
	w.ScrutinyDay()
	if w.Scrutiny() >= before {
		t.Fatal("an old story kept raising the temperature")
	}
	// And it cannot go past the ceiling.
	w.Attention = ScrutinyCeiling
	for i := 0; i < 5; i++ {
		w.Minute += 1440
		w.Report("attack", "AN EXPLOSION", "It happened.")
		w.ScrutinyDay()
	}
	if w.Scrutiny() != ScrutinyCeiling {
		t.Fatalf("the city went to %d", w.Scrutiny())
	}
}

func TestACrackdownIsAnnouncedOnceAndLiftsOnce(t *testing.T) {
	t.Parallel()
	w := quietCity(t)
	w.Attention = ScrutinyCrackdown - 1
	w.Minute = 5000
	w.Report("killing", "SOMEBODY KILLED", "It happened.")
	w.ScrutinyDay()
	if !w.UnderCrackdown() {
		t.Fatal("the city did not act")
	}
	if !w.hasRecord("The city has had enough") || !w.hasNewsKind("police") {
		t.Fatal("nobody was told")
	}
	// It is not announced again while it holds.
	count := 0
	for _, r := range w.History {
		if r.Title == "The city has had enough" {
			count++
		}
	}
	w.Report("killing", "SOMEBODY ELSE KILLED", "It happened.")
	w.ScrutinyDay()
	again := 0
	for _, r := range w.History {
		if r.Title == "The city has had enough" {
			again++
		}
	}
	if again != count {
		t.Fatal("the crackdown was announced twice")
	}
	// And it lifts.
	w.News = nil
	for i := 0; i < 200 && w.UnderCrackdown(); i++ {
		w.Minute += 1440
		w.ScrutinyDay()
	}
	if w.UnderCrackdown() || !w.hasRecord("It has gone quiet") {
		t.Fatal("the crackdown never lifted")
	}
}

func TestACrackdownBringsThePoliceSooner(t *testing.T) {
	t.Parallel()
	w := quietCity(t)
	if w.ScrutinyRaidShift() != 0 {
		t.Fatal("a quiet city brought anybody sooner")
	}
	w.Attention = ScrutinyCeiling
	if w.ScrutinyRaidShift() != ScrutinyRaidRelief {
		t.Fatalf("the worst the city gets brings them %d sooner", w.ScrutinyRaidShift())
	}

	const runs = 400
	raided := func(scrutiny int) int {
		hit := 0
		for seed := uint32(1); seed <= runs; seed++ {
			probe := quietCity(t)
			probe.WorldRNG = seed * 2654435761
			probe.Attention = scrutiny
			probe.Player.Heat = RaidThreshold - 10
			probe.Properties["laundry"].Owner = probe.PlayerOrganizationID()
			probe.PoliceDay()
			if probe.hasRecord("Turned over at Bluebird Laundry") || probe.hasRecord("They came to the door") {
				hit++
			}
		}
		return hit
	}
	if quiet, loud := raided(0), raided(ScrutinyCeiling); loud <= quiet {
		t.Fatalf("below the ordinary threshold, a quiet city raided %d times and a crackdown %d", quiet, loud)
	} else {
		t.Logf("held ten below the ordinary raid threshold: a quiet city comes to the door %d times in %d, a city at its worst %d", quiet, runs, loud)
	}
}

func TestACrackdownMakesTheBuildingDearerAndFussier(t *testing.T) {
	t.Parallel()
	w := quietCity(t)
	o, _ := OfficialByID("mayor")
	if w.ScrutinyPremium() != 0 || w.OfficialOpening(o) != o.Opening || w.OfficialCeiling(o) != o.Ceiling {
		t.Fatal("a quiet city charged a premium")
	}
	w.Attention = ScrutinyCeiling
	if w.OfficialOpening(o) <= o.Opening {
		t.Fatalf("an arrangement cost $%d at the worst of it against $%d", w.OfficialOpening(o), o.Opening)
	}
	if w.OfficialCeiling(o) >= o.Ceiling {
		t.Fatalf("the ceiling was %d against %d", w.OfficialCeiling(o), o.Ceiling)
	}
	if w.OfficialCeiling(o) < 20 {
		t.Fatal("the ceiling fell below anything workable")
	}
}
