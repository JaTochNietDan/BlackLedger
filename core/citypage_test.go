package core

import "strings"

import "testing"

// The city page is filler, and filler has to be harmless: it must not make the
// police look harder, must not lead over a killing, and must not tell the
// reader anything the city could not see for itself.

func pageOf(w *World) []Story {
	out := []Story{}
	for _, s := range w.News {
		if s.Kind == "civic" {
			out = append(out, s)
		}
	}
	return out
}

func TestTheCityPagePrintsSomethingOnAQuietDay(t *testing.T) {
	w := New(71)
	w.CityPageDay()
	page := pageOf(w)
	if len(page) != 2 {
		t.Fatalf("expected two briefs on a quiet day, got %d", len(page))
	}
	if page[0].Headline == page[1].Headline {
		t.Error("the same brief ran twice in one issue")
	}
	for _, s := range page {
		if strings.TrimSpace(s.Body) == "" {
			t.Error("a brief with no copy in it")
		}
	}
}

func TestTheCityPageCostsTheCityNothing(t *testing.T) {
	// A column about the price of coal is not a reason for anybody to look
	// harder at anybody.
	if scrutinyWeight["civic"] != 0 {
		t.Fatalf("filler carries %d scrutiny; it must carry none", scrutinyWeight["civic"])
	}
	w := New(71)
	before := w.Scrutiny()
	for i := 0; i < 20; i++ {
		w.Minute += 1440
		w.CityPageDay()
	}
	w.ScrutinyDay()
	if w.Scrutiny() > before {
		t.Errorf("twenty days of filler moved scrutiny from %d to %d", before, w.Scrutiny())
	}
}

func TestTheCityPageNeverLeads(t *testing.T) {
	w := New(71)
	w.CityPageDay()
	w.Report("killing", "A MAN IS FOUND", "Somebody was killed.")
	issues := w.Editions()
	if len(issues) == 0 {
		t.Fatal("no issue went to press")
	}
	stories := issues[0]["stories"].([]map[string]any)
	if stories[0]["kind"] != "killing" {
		t.Errorf("the paper led with %q over a killing", stories[0]["kind"])
	}
}

func TestTheSameDayPrintsTheSamePage(t *testing.T) {
	a, b := New(71), New(72)
	b.ID = a.ID
	a.Minute, b.Minute = 9*1440, 9*1440
	a.CityPageDay()
	b.CityPageDay()
	first, second := pageOf(a), pageOf(b)
	if len(first) != len(second) {
		t.Fatal("the same day printed a different number of briefs")
	}
	for i := range first {
		if first[i].Headline != second[i].Headline {
			t.Errorf("day nine printed %q one time and %q the next", first[i].Headline, second[i].Headline)
		}
	}
}

func TestTheCityPageSaysWhatTheSkyIsDoing(t *testing.T) {
	// The weather in the paper has to be the weather the world had, or the
	// paper is making things up — which is the one thing it may not do.
	seed := skySeed("weather-campaign")
	rainy := -1
	for day := 0; day < 400 && rainy < 0; day++ {
		if skyOnDay(day, seed) == "rain" {
			rainy = day
		}
	}
	if rainy < 0 {
		t.Fatal("400 days without rain, which is not weather")
	}
	w := New(71)
	w.ID = "weather-campaign"
	w.Minute = rainy * 1440
	w.CityPageDay()
	printed := false
	for _, s := range w.cityPage() {
		if strings.Contains(s.headline, "RAIN") {
			printed = true
		}
	}
	if !printed {
		t.Errorf("it rained on day %d and the paper had nothing to say about it", rainy)
	}
}

func TestFillerGivesUpItsPlaceBeforeNews(t *testing.T) {
	// Two briefs a day would fill the archive in four months, and a player
	// reading back through the paper would find nothing but the weather.
	w := New(73)
	w.Report("killing", "THE FIRST KILLING", "Somebody was killed on the first day.")
	for i := 0; i < newsCapacity*2; i++ {
		w.Minute += 1440
		w.CityPageDay()
	}
	if len(w.News) > newsCapacity {
		t.Fatalf("the archive grew to %d past its bound of %d", len(w.News), newsCapacity)
	}
	kept := false
	for _, s := range w.News {
		if s.Headline == "THE FIRST KILLING" {
			kept = true
		}
	}
	if !kept {
		t.Error("a killing was evicted from the archive by filler about the weather")
	}
}
