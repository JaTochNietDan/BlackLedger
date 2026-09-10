package core

import (
	"fmt"
	"strings"
	"testing"
)

// Measured over forty campaigns: four days away left the city different on
// return every time, and the player was told "4 days gone. The city did not
// wait" — a claim with nothing in it. Everything that had happened was in the
// Herald, filed under days the player had no reason to go back and read.

func TestComingHomeSaysWhatHappened(t *testing.T) {
	t.Parallel()
	w := New(31)
	left := w.Minute
	w.Minute += 1440
	w.Report("killing", "SOMEBODY IS DEAD", "b")
	w.Report("war", "TWO FAMILIES AT WAR", "b")
	// The record is written when the player gets back, which is after.
	w.Minute += 1440
	line := w.WhatYouMissed(left)
	if !strings.Contains(line, "SOMEBODY IS DEAD") || !strings.Contains(line, "TWO FAMILIES AT WAR") {
		t.Fatalf("a returning player reads %q", line)
	}
	// The biggest story leads, the way it led the front page.
	if strings.Index(line, "SOMEBODY IS DEAD") > strings.Index(line, "TWO FAMILIES AT WAR") {
		t.Fatalf("a war led a killing: %q", line)
	}
}

// Somebody who has been away four days does not need to be told it rained.
func TestComingHomeLeavesOutTheWeather(t *testing.T) {
	t.Parallel()
	w := New(32)
	left := w.Minute
	w.Minute += 1440
	w.Report("civic", "RAIN ACROSS THE DISTRICT", "b")
	w.Report("civic", "THE CITY COUNTED", "b")
	w.Minute += 1440
	if line := w.WhatYouMissed(left); !strings.HasPrefix(line, "Nothing") {
		t.Fatalf("four days of weather came back as %q", line)
	}
}

// Nothing before the trip belongs in it. The player was here for that.
func TestComingHomeDoesNotReportWhatYouWereThereFor(t *testing.T) {
	t.Parallel()
	w := New(33)
	w.Report("killing", "BEFORE YOU LEFT", "b")
	w.Minute += 1440
	left := w.Minute
	w.Minute += 1440
	w.Report("killing", "WHILE YOU WERE GONE", "b")
	w.Minute += 1440
	line := w.WhatYouMissed(left)
	if strings.Contains(line, "BEFORE YOU LEFT") {
		t.Fatalf("the city told a returning player something they watched happen: %q", line)
	}
	if !strings.Contains(line, "WHILE YOU WERE GONE") {
		t.Fatalf("the city left out what it did in their absence: %q", line)
	}
}

// A returning player gets a summary, not the whole paper.
func TestComingHomeIsASummaryNotAnArchive(t *testing.T) {
	t.Parallel()
	w := New(34)
	left := w.Minute
	w.Minute += 1440
	for i := 0; i < 9; i++ {
		w.Report("robbery", "ROBBERY NUMBER "+string(rune('A'+i)), "b")
	}
	w.Minute += 1440
	line := w.WhatYouMissed(left)
	if strings.Count(line, "ROBBERY NUMBER") > MissedHeadlines {
		t.Fatalf("the whole paper came back with them: %q", line)
	}
	if !strings.Contains(line, "6 other stories") {
		t.Fatalf("the rest went unmentioned: %q", line)
	}
}

// A previous life's news is not this one's.
func TestComingHomeDoesNotReadAPreviousLifesPaper(t *testing.T) {
	t.Parallel()
	w := New(35)
	left := w.Minute
	w.Minute += 1440
	w.Report("killing", "IN A LIFE BEFORE THIS ONE", "b")
	w.Minute += 1440
	w.Life++
	if line := w.WhatYouMissed(left); strings.Contains(line, "IN A LIFE BEFORE") {
		t.Fatalf("a new protagonist was handed the last one's headlines: %q", line)
	}
}

// Read out of a real save after a four-day trip: "Whatever was arranged for you
// happened 1 times to a locked door." A count of one takes a singular noun.
func TestOneIsNotPlural(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		n    int
		want string
	}{{1, "once"}, {2, "2 times"}, {5, "5 times"}} {
		got := plainly(c.n, "once", fmt.Sprintf("%d times", c.n))
		if got != c.want {
			t.Fatalf("%d reads as %q, wanted %q", c.n, got, c.want)
		}
	}
}

// Read out of a real save: "POLICE PRESSURE ON RUSSO OUTFIT · POLICE PRESSURE
// ON RUSSO OUTFIT". The paper collapses repeats within a day; a trip spans
// several, and two days of the same headline is one thing to be told.
func TestComingHomeDoesNotSayTheSameThingTwice(t *testing.T) {
	t.Parallel()
	w := New(36)
	left := w.Minute
	for day := 1; day <= 3; day++ {
		w.Minute = left + day*1440
		w.Report("police", "POLICE PRESSURE ON RUSSO OUTFIT", "b")
	}
	w.Minute += 1440
	line := w.WhatYouMissed(left)
	if strings.Count(line, "POLICE PRESSURE") != 1 {
		t.Fatalf("a returning player is told the same thing more than once: %q", line)
	}
	if strings.Contains(line, "other stor") {
		t.Fatalf("the duplicates were counted as further news: %q", line)
	}
}

// The paper files the story of the player's own arrest at the exact minute the
// door shuts, so a man walking out of a cell was being told that the paper had
// reported him walking into it. The minute you leave is a minute you were there
// for.
func TestComingOutIsNotToldAboutYourOwnArrest(t *testing.T) {
	t.Parallel()
	w := New(41)
	w.Player.Home = "room"
	w.Confine(4, "what was found at the laundry")
	arrest := w.News[len(w.News)-1].Headline
	if err := w.SitOut(); err != nil {
		t.Fatal(err)
	}
	var out string
	for _, r := range w.History {
		if r.Title == "Out" {
			out = r.Text
		}
	}
	if out == "" {
		t.Fatal("nothing was filed on release")
	}
	if strings.Contains(out, arrest) {
		t.Fatalf("a man leaving a cell was told the paper reported him entering it: %q", out)
	}
}

// Twelve days inside is three times the longest trip out of the city, and the
// release record used to say only that whatever it cost happened while the
// player was not there to watch it.
func TestComingOutSaysWhatHappened(t *testing.T) {
	t.Parallel()
	w := New(42)
	w.Player.Home = "room"
	w.Confine(6, "what was found at the laundry")
	w.Minute += 1440
	w.Report("killing", "SOMEBODY IS DEAD", "b")
	w.Minute = w.Player.HeldUntil - 1440
	if err := w.SitOut(); err != nil {
		t.Fatal(err)
	}
	var out string
	for _, r := range w.History {
		if r.Title == "Out" {
			out = r.Text
		}
	}
	if !strings.Contains(out, "SOMEBODY IS DEAD") {
		t.Fatalf("six days inside and the city said nothing: %q", out)
	}
}

// And the other end of the same fault: after talking his way out, the player
// was told the paper had reported him talking his way out.
func TestComingOutIsNotToldAboutYourOwnRelease(t *testing.T) {
	t.Parallel()
	w := New(43)
	w.Player.Home = "room"
	w.Confine(6, "what was found at the laundry")
	since := w.Player.HeldFrom
	w.Minute += 2 * 1440
	w.Report("killing", "SOMEBODY IS DEAD", "b")
	// The release is filed at the moment the player walks out, which is now.
	w.Minute += 1440
	w.Report("police", "CHARGES DROPPED AFTER COOPERATION", "b")
	line := w.WhatYouMissed(since)
	if strings.Contains(line, "CHARGES DROPPED") {
		t.Fatalf("the player was told the paper reported his own release: %q", line)
	}
	if !strings.Contains(line, "SOMEBODY IS DEAD") {
		t.Fatalf("real news went with it: %q", line)
	}
}
