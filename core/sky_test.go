package core

import "testing"

// The sky is a fact about the world, so it has to behave like one: the same
// day of the same campaign is always the same day, two campaigns do not share
// a fortnight of weather, and the streets do not dry the instant it stops.

func TestTheSameDayHasTheSameWeatherEveryTimeItIsAsked(t *testing.T) {
	t.Parallel()
	w := &World{ID: "campaign-one", Minute: 5 * 1440}
	first := w.Sky()
	for i := 0; i < 50; i++ {
		if w.Sky() != first {
			t.Fatal("the weather changed while nobody did anything")
		}
	}
	// And it does not depend on the hour, only the day.
	morning := &World{ID: "campaign-one", Minute: 5*1440 + 400}
	if morning.Sky() != first {
		t.Fatal("the weather changed between morning and afternoon of the same day")
	}
}

func TestTwoCampaignsDoNotShareTheirWeather(t *testing.T) {
	t.Parallel()
	same := 0
	for day := 0; day < 40; day++ {
		if SkyOn(day, skySeed("campaign-one")) == SkyOn(day, skySeed("campaign-two")) {
			same++
		}
	}
	if same > 30 {
		t.Fatalf("two campaigns had the same weather on %d of 40 days", same)
	}
}

func TestTheStreetsAreStillWetTheDayAfterRain(t *testing.T) {
	t.Parallel()
	seed := skySeed("campaign-one")
	found := false
	for day := 1; day < 400; day++ {
		if skyOnDay(day-1, seed) == "rain" && skyOnDay(day, seed) == "clear" {
			found = true
			if got := SkyOn(day, seed); got.Wet <= 0 {
				t.Fatalf("day %d came after rain and was bone dry: %+v", day, got)
			}
			if SkyOn(day, seed).Wet >= SkyOn(day-1, seed).Wet {
				t.Fatal("the day after rain is as wet as the rain itself")
			}
		}
	}
	if !found {
		t.Fatal("400 days without a clear morning after rain, which is not weather")
	}
}

func TestEveryKindOfWeatherHappens(t *testing.T) {
	t.Parallel()
	seen := map[string]int{}
	seed := skySeed("campaign-one")
	for day := 0; day < 400; day++ {
		seen[skyOnDay(day, seed)]++
	}
	for _, kind := range []string{"clear", "overcast", "rain", "fog"} {
		if seen[kind] < 10 {
			t.Errorf("%s happened %d times in 400 days", kind, seen[kind])
		}
	}
	// And most days are not raining, because most days are not.
	if seen["rain"]+seen["fog"] > 200 {
		t.Errorf("it was wet on %d of 400 days", seen["rain"]+seen["fog"])
	}
}

func TestRainIsWetterThanFogAndClearIsDry(t *testing.T) {
	t.Parallel()
	seed := skySeed("campaign-one")
	for day := 0; day < 400; day++ {
		s := SkyOn(day, seed)
		if s.Kind == "rain" && s.Wet < 1 {
			t.Fatal("rain that does not wet the streets")
		}
		if s.Kind == "clear" && day > 0 && skyOnDay(day-1, seed) != "rain" && s.Wet != 0 {
			t.Fatal("a clear day after a clear day with wet streets")
		}
	}
}
