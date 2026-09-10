package core

import (
	"strings"
	"testing"
)

// The paper says what the street has noticed about the families. It must not
// print a number, must not report the same slide every morning, and must not
// report the player to the player.

func civicAbout(w *World, name string) []Story {
	out := []Story{}
	for _, s := range w.News {
		if s.Kind == "civic" && strings.Contains(s.Headline, upper(name)) {
			out = append(out, s)
		}
	}
	return out
}

func TestAFamilyLosingGroundIsNoticed(t *testing.T) {
	t.Parallel()
	w := New(111)
	f := &w.Factions[0]
	f.Power -= FortuneShift + 4
	w.FortunesDay()
	got := civicAbout(w, f.Name)
	if len(got) != 1 {
		t.Fatalf("a family losing %d strength produced %d stories", FortuneShift+4, len(got))
	}
	if !strings.Contains(got[0].Headline, "STRUGGLING") {
		t.Errorf("a family that lost ground was reported as %q", got[0].Headline)
	}
}

func TestTheSameSlideIsNotReportedEveryMorning(t *testing.T) {
	t.Parallel()
	w := New(112)
	f := &w.Factions[0]
	f.Power -= FortuneShift + 2
	w.FortunesDay()
	for i := 0; i < 6; i++ {
		w.Minute += 1440
		w.FortunesDay()
	}
	if got := civicAbout(w, f.Name); len(got) != 1 {
		t.Errorf("one fall was reported %d times", len(got))
	}
}

func TestASmallDriftIsNotNews(t *testing.T) {
	t.Parallel()
	w := New(113)
	f := &w.Factions[0]
	f.Power -= FortuneShift - 1
	w.FortunesDay()
	if got := civicAbout(w, f.Name); len(got) != 0 {
		t.Errorf("a drift of %d was reported as news: %q", FortuneShift-1, got[0].Headline)
	}
}

func TestThePaperNeverPrintsTheNumber(t *testing.T) {
	t.Parallel()
	// "Power 41" is a statistic. A paper writes about what people have noticed.
	w := New(114)
	for i := range w.Factions {
		w.Factions[i].Power -= FortuneShift + 9
	}
	w.FortunesDay()
	for _, s := range w.News {
		if s.Kind != "civic" {
			continue
		}
		for _, digit := range "0123456789" {
			if strings.ContainsRune(s.Body+s.Headline, digit) {
				t.Errorf("the paper printed a figure out of the simulation: %q", s.Body)
				break
			}
		}
	}
}

func TestThePlayersOwnOutfitIsNotNewsToThePlayer(t *testing.T) {
	t.Parallel()
	w := New(115)
	w.Factions = append(w.Factions, Faction{
		ID: w.PlayerOrganizationID(), Name: "Ward's people", Power: 60, Reported: 60,
	})
	w.Factions[len(w.Factions)-1].Power -= FortuneShift * 3
	w.FortunesDay()
	if got := civicAbout(w, "Ward's people"); len(got) != 0 {
		t.Errorf("the paper told the player about their own outfit: %q", got[0].Headline)
	}
}
