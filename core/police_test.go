package core

import "testing"

func watched(t *testing.T, heat int) *World {
	t.Helper()
	w := New(79)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Heat = heat
	w.Player.Cash = 5000
	return w
}

func TestAttentionFadesWhenNothingIsAdded(t *testing.T) {
	t.Parallel()
	w := watched(t, 20)
	w.PoliceDay()
	// The constant was a flat CoolOff for everybody. How much a person is
	// forgotten now depends on what they are known to own, so this asks the
	// world rather than a number: a player holding nothing watched gets the
	// base rate.
	if w.Player.Heat != 20-w.CoolOff() {
		t.Fatalf("attention did not fade: %d", w.Player.Heat)
	}
	// It never goes below nothing.
	w.Player.Heat = 0
	w.PoliceDay()
	if w.Player.Heat != 0 {
		t.Fatal("attention went negative")
	}
	// Fading is slower than a hard-run business generates, or the operating
	// decision would cost nothing.
	if BaseCool >= operatingMode("hard").Heat {
		t.Fatal("a skimmed business generates less attention than fades in a day")
	}
}

func TestARaidTakesStockCashAndCondition(t *testing.T) {
	t.Parallel()
	w := watched(t, 60)
	w.Player.Stock = map[string]int{"moonshine": 15}
	cash, condition := w.Player.Cash, w.Properties["laundry"].Condition
	w.Raid()
	if w.Carrying() != 0 {
		t.Fatal("a raid left the stock alone")
	}
	if w.Player.Cash >= cash {
		t.Fatal("a raid cost nothing in fines")
	}
	if w.Properties["laundry"].Condition >= condition {
		t.Fatal("a raid left the premises untouched")
	}
	if w.Player.Heat >= 60 {
		t.Fatal("attention did not fall after they got what they came for")
	}
	if !w.Own("laundry") {
		t.Fatal("a business was forfeited below the forfeiture threshold")
	}
	// The city reads about it.
	found := false
	for _, s := range w.Edition() {
		if s["kind"] == "police" {
			found = true
		}
	}
	if !found {
		t.Fatal("a raid was never reported")
	}
}

func TestPastAPointTheyTakeTheBusiness(t *testing.T) {
	t.Parallel()
	w := watched(t, ForfeitThreshold+5)
	w.Raid()
	if w.Own("laundry") {
		t.Fatal("a business survived a raid above the forfeiture threshold")
	}
	if w.Properties["laundry"].Owner != "independent" {
		t.Fatalf("the forfeited business went to %q", w.Properties["laundry"].Owner)
	}
}

func TestTheyOnlyComeWhenThereIsEnoughToComeFor(t *testing.T) {
	t.Parallel()
	quiet, loud := 0, 0
	for i := uint32(1); i <= 400; i++ {
		// Consecutive seeds give this generator nearly identical first draws,
		// so stride them or the measured rate is an artifact of the seed.
		seed := i * 2654435761
		low := watched(t, RaidThreshold-10)
		low.WorldRNG = seed
		before := low.Player.Cash
		low.considerRaid()
		if low.Player.Cash != before {
			quiet++
		}
		high := watched(t, 75)
		high.WorldRNG = seed
		cash := high.Player.Cash
		high.considerRaid()
		if high.Player.Cash != cash {
			loud++
		}
	}
	t.Logf("of 400 days: %d raids below the threshold, %d at high attention", quiet, loud)
	if quiet != 0 {
		t.Fatal("the police raided somebody they were not watching")
	}
	if loud == 0 {
		t.Fatal("the police never arrive however well known you become")
	}
	if loud == 400 {
		t.Fatal("a raid every single day is not a risk, it is a schedule")
	}
}

func TestADetectiveWillNotBeSeenWithYouForever(t *testing.T) {
	t.Parallel()
	w := watched(t, 30)
	cost := w.BribeCost()
	cash := w.Player.Cash
	if err := w.Bribe(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Heat >= 30 {
		t.Fatal("the bribe cleared nothing")
	}
	if w.Player.Cash != cash-cost {
		t.Fatal("the bribe was not paid for")
	}
	// Past the ceiling, money stops working.
	hot := watched(t, BribeCeiling+5)
	if hot.BribeReadiness() == "" {
		t.Fatal("a detective took money from somebody far too well known")
	}
	if err := hot.Bribe(); err == nil {
		t.Fatal("the command ignored the ceiling the action reports")
	}
	// And there is nothing to buy when nobody is looking.
	clean := watched(t, 0)
	if clean.BribeReadiness() == "" {
		t.Fatal("a bribe was offered with no attention to clear")
	}
}

func TestWarCostsOrganizationsMoneyToo(t *testing.T) {
	t.Parallel()
	w := New(83)
	w.Antagonize("bellandi", "russo", 100)
	if w.Conflict("bellandi", "russo").State != "war" {
		t.Fatal("the test world is not at war")
	}
	before := w.faction("bellandi").Cash
	w.PoliceDay()
	if w.faction("bellandi").Cash >= before {
		t.Fatal("a war drew no police attention to the organizations fighting it")
	}
	// An organization at peace is left alone.
	peace := New(83)
	quiet := peace.faction("bellandi").Cash
	peace.PoliceDay()
	if peace.faction("bellandi").Cash != quiet {
		t.Fatal("the police fined an organization that was not fighting anyone")
	}
}
