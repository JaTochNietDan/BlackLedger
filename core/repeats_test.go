package core

import (
	"strings"
	"testing"
)

// Read end to end, a campaign's paper carried seven identical paragraphs about
// a robbery at Saint Agnes in one edition. No newspaper does that, and none of
// the tests could see it because each individual story was correct.

func TestThePaperDoesNotPrintTheSameStorySevenTimes(t *testing.T) {
	w := New(141)
	for i := 0; i < 7; i++ {
		w.Report("robbery", "ROBBERY IN SAINT AGNES", "A man was robbed near Saint Agnes.")
	}
	filed := 0
	for _, s := range w.News {
		if s.Kind == "robbery" {
			filed++
		}
	}
	if filed != 1 {
		t.Fatalf("seven robberies produced %d stories", filed)
	}
	if !strings.Contains(w.News[0].Body, "seven times") {
		t.Errorf("the story does not say how often it happened: %q", w.News[0].Body)
	}
	if w.News[0].Count != 7 {
		t.Errorf("the story was counted %d times", w.News[0].Count)
	}
}

func TestTwoDifferentRobberiesAreTwoStories(t *testing.T) {
	w := New(142)
	w.Report("robbery", "ROBBERY IN SAINT AGNES", "A man was robbed near Saint Agnes.")
	w.Report("robbery", "ROBBERY AT THE MONARCH", "The day's takings went out of the back.")
	if len(w.News) != 2 {
		t.Errorf("two robberies in different places produced %d stories", len(w.News))
	}
}

func TestYesterdaysStoryIsNotTodaysStory(t *testing.T) {
	// The paper collapses a run within one issue. It does not reach back into
	// yesterday's edition, which has gone out.
	w := New(143)
	w.Report("robbery", "ROBBERY IN SAINT AGNES", "A man was robbed near Saint Agnes.")
	w.Minute += 1440
	w.Report("robbery", "ROBBERY IN SAINT AGNES", "A man was robbed near Saint Agnes.")
	if len(w.News) != 2 {
		t.Errorf("a robbery on two separate days produced %d stories", len(w.News))
	}
}

func TestThePoliceLineIsNotPrintedTwice(t *testing.T) {
	w := New(144)
	doubled := w.unattributed("Saint Agnes",
		"A man was robbed near Saint Agnes. Police have asked anybody who saw it to come forward.")
	if strings.Count(doubled, "Police") > 1 {
		t.Errorf("the police appear twice in one paragraph: %q", doubled)
	}
	// And copy that has not mentioned them still gets the line.
	plain := w.unattributed("Saint Agnes", "A sum was taken from Saint Agnes.")
	if !strings.Contains(plain, "no arrest") {
		t.Errorf("copy with no police line did not get one: %q", plain)
	}
}
