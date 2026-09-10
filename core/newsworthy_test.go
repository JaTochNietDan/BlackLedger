package core

import "testing"

// What leads the paper and what the police notice are different questions, and
// they were the same map. An arrest, an attempt on somebody's life and a family
// splitting all scored zero for police attention — correctly, they are not why
// a city looks harder at anybody — and that same zero was being used to decide
// what leads an edition, so none of them ever could.

func TestEveryKindThePaperFilesHasAPlaceOnThePage(t *testing.T) {
	t.Parallel()
	// Every kind passed to Report anywhere in the core.
	filed := []string{"arrest", "attack", "attempt", "business", "civic", "collapse",
		"killing", "obituary", "police", "politics", "recovery", "robbery",
		"seizure", "split", "war"}
	for _, kind := range filed {
		if Newsworthiness(kind) <= 0 {
			t.Errorf("%s is filed by the paper and has no place on the page", kind)
		}
	}
}

func TestAKillingLeadsOverAnArrestOverTheWeather(t *testing.T) {
	t.Parallel()
	if !(Newsworthiness("killing") > Newsworthiness("arrest")) {
		t.Error("an arrest leads over a killing")
	}
	if !(Newsworthiness("arrest") > Newsworthiness("civic")) {
		t.Error("the weather leads over an arrest")
	}
	if !(Newsworthiness("civic") > 0) {
		t.Error("the city page cannot be shown at all")
	}
	// And an obituary never leads: it is set apart from the news anyway.
	if Newsworthiness("obituary") >= Newsworthiness("robbery") {
		t.Error("an obituary would lead over a robbery")
	}
}

func TestRankingThePaperDoesNotChangeWhatThePoliceNotice(t *testing.T) {
	t.Parallel()
	// The balance of the game depends on scrutiny. Ranking the front page must
	// not touch it: a column about the price of coal still costs nothing, and a
	// killing still costs what it always did.
	if scrutinyWeight["civic"] != 0 || scrutinyWeight["obituary"] != 0 {
		t.Error("filler started raising the city's temperature")
	}
	if scrutinyWeight["killing"] != 14 {
		t.Errorf("a killing now costs %d attention", scrutinyWeight["killing"])
	}
	if scrutinyWeight["arrest"] != 0 {
		t.Errorf("an arrest now costs %d attention, which is a balance change", scrutinyWeight["arrest"])
	}
}

func TestTheLeadIsTheBiggestStoryNotTheLatest(t *testing.T) {
	t.Parallel()
	w := New(131)
	w.Report("civic", "A CLEAR DAY", "Clear over Bellwether.")
	w.Report("killing", "A MAN IS FOUND", "Somebody was killed at noon.")
	w.Report("civic", "NO CHANGE IN THE WEATHER", "Cloud, and more of it tomorrow.")
	issues := w.Editions()
	stories := issues[0]["stories"].([]map[string]any)
	if stories[0]["kind"] != "killing" {
		t.Errorf("the paper led with %q over a killing", stories[0]["kind"])
	}
}
