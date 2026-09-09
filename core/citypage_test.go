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

// Read out of a real campaign's paper: THE CITY COUNTED ran on four days in
// seven with the same forty-eight people and the same eighteen of them
// employed, and the Mariner was announced as standing empty five times in the
// same week. The pick was randomised per day, and nothing ever looked at what
// had already been printed.
func TestTheCityPageDoesNotRepeatItselfAllWeek(t *testing.T) {
	w := New(88)
	w.News = nil
	for day := 0; day < 21; day++ {
		w.Minute = day * 1440
		w.CityPageDay()
	}
	// No headline twice inside its rest period.
	seen := map[string]int{}
	for _, s := range w.News {
		if s.Kind != "civic" {
			continue
		}
		day := s.Minute / 1440
		if last, ok := seen[s.Headline]; ok && day-last < CityPageRest {
			t.Fatalf("%q ran on day %d and again on day %d", s.Headline, last+1, day+1)
		}
		seen[s.Headline] = day
	}
	// And no issue runs the same brief twice.
	byDay := map[int]map[string]bool{}
	for _, s := range w.News {
		if s.Kind != "civic" {
			continue
		}
		day := s.Minute / 1440
		if byDay[day] == nil {
			byDay[day] = map[string]bool{}
		}
		if byDay[day][s.Headline] {
			t.Fatalf("day %d ran %q twice", day+1, s.Headline)
		}
		byDay[day][s.Headline] = true
	}
	// The page still comes out. Three weeks of empty issues is not a fix.
	if len(w.News) < 21 {
		t.Fatalf("three weeks produced %d civic briefs", len(w.News))
	}
}

// The city keeps hours — people are at their posts through the morning and in
// the bars and clubs after midday — and the paper had never once mentioned it.
// Counted over a fifty-eight day campaign, the Herald was 59% civic filler and
// twenty-six of its fifty-eight issues carried no news at all; a paper printed
// in this city with nothing to say about its evenings is missing something it
// can see.
func TestThePaperCanSeeTheEvening(t *testing.T) {
	w := New(41)
	// Everybody out.
	for i := range w.NPCs {
		w.NPCs[i].Location, w.NPCs[i].Heading, w.NPCs[i].Dead = haunts[i%len(haunts)], "", false
	}
	full := ""
	for _, b := range w.cityPage() {
		if strings.Contains(b.headline, "FULL NIGHT") {
			full = b.body
		}
	}
	if full == "" {
		t.Fatal("the whole district was in the bars and the paper had nothing to say")
	}
	// Nobody out.
	for i := range w.NPCs {
		w.NPCs[i].Location = "market"
	}
	empty := ""
	for _, b := range w.cityPage() {
		if strings.Contains(b.headline, "HOUSES WERE EMPTY") {
			empty = b.body
		}
	}
	if empty == "" {
		t.Fatal("nobody was out and the paper had nothing to say")
	}
	// An ordinary night is not news.
	live := w.living()
	for i := range w.NPCs {
		if i < live/6 {
			w.NPCs[i].Location = haunts[i%len(haunts)]
		} else {
			w.NPCs[i].Location = "market"
		}
	}
	for _, b := range w.cityPage() {
		if strings.Contains(b.headline, "NIGHT ON THE FRONT") || strings.Contains(b.headline, "HOUSES WERE EMPTY") {
			t.Fatalf("an ordinary evening was reported as news: %q", b.headline)
		}
	}
	// And it counts in words, like the rest of the paper.
	if strings.ContainsAny(full+empty, "0123456789") {
		t.Fatalf("the evening brief prints figures in prose: %q / %q", full, empty)
	}
}
