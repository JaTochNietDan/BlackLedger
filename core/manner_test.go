package core

import (
	"strings"
	"testing"
)

func morgue(t *testing.T) *World {
	t.Helper()
	w := New(113)
	w.MigrateLivingWorld()
	return w
}

func TestEveryPlaceInTheCityHasItsOwnWaysOfDying(t *testing.T) {
	t.Parallel()
	for _, l := range Locations {
		options := mannerByPlace[l.ID]
		if len(options) < 3 {
			t.Fatalf("%s had %d ways to die in it", l.Name, len(options))
		}
		seen := map[string]bool{}
		for _, o := range options {
			if seen[o] {
				t.Fatalf("%s repeated a way of dying", l.Name)
			}
			seen[o] = true
		}
	}
}

func TestTheSamePlaceDoesNotReadTheSameTwice(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	victim := w.NPC("vittorio")
	seen := map[string]bool{}
	for i := 0; i < 60; i++ {
		seen[w.Manner(victim, nil)] = true
	}
	if len(seen) < 3 {
		t.Fatalf("sixty killings at the same place produced %d descriptions", len(seen))
	}
}

func TestTheHourIsInIt(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	victim := w.NPC("vittorio")
	hours := map[string]bool{}
	for _, minute := range []int{60, 400, 700, 800, 1000, 1300, 1400} {
		w.Minute = minute
		hours[hourOf(w.Minute)] = true
	}
	if len(hours) < 5 {
		t.Fatalf("a whole day produced %d different hours", len(hours))
	}
	w.Minute = 120
	if !strings.Contains(w.Manner(victim, nil), "small hours") {
		t.Fatalf("a killing at two in the morning read as %q", w.Manner(victim, nil))
	}
}

func TestWhoDidItLeavesAMark(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	victim := w.NPC("vittorio")
	marks := map[string]bool{}
	for _, id := range []string{"hot", "careful", "greedy", "vain", "loyal"} {
		killer := &NPC{ID: "k", Name: nameWith(t, id)}
		mark := signature(killer)
		if mark == "" {
			t.Fatalf("a %s killer left nothing behind", id)
		}
		if marks[mark] {
			t.Fatalf("two kinds of person left the same mark")
		}
		marks[mark] = true
		if !strings.HasSuffix(w.Manner(victim, killer), mark) {
			t.Fatalf("the mark did not reach the description")
		}
	}
	if signature(nil) != "" {
		t.Fatal("nobody left a mark")
	}
}

func TestSomebodyKilledNowhereStillGetsADescription(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	stray := &NPC{ID: "stray", Name: "A Stranger", Location: "transit"}
	w.NPCs = append(w.NPCs, *stray)
	manner := w.Manner(w.NPC("stray"), nil)
	if manner == "" || !strings.Contains(manner, "A Stranger") || !strings.Contains(manner, "In the street") {
		t.Fatalf("somebody killed nowhere in particular read as %q", manner)
	}
}

func TestKillByPutsTheMannerInTheRecordAndThePaper(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	killer := w.NPC("elena")
	if !w.KillBy("vittorio", killer, "It was over an old debt.") {
		t.Fatal("nobody died")
	}
	found := false
	for _, r := range w.History {
		if r.Title == "Vittorio Bellandi is dead" {
			found = true
			if !strings.Contains(r.Text, "At The Monarch") || !strings.Contains(r.Text, "old debt") {
				t.Fatalf("the record read %q", r.Text)
			}
		}
	}
	if !found {
		t.Fatal("nothing was recorded")
	}
	if len(w.News) == 0 {
		t.Fatal("the paper did not carry it")
	}
	story := w.News[len(w.News)-1]
	if !strings.Contains(story.Body, "Vittorio Bellandi") {
		t.Fatalf("the paper read %q", story.Body)
	}
	// And nobody dies twice.
	if w.KillBy("vittorio", killer, "Again.") {
		t.Fatal("a dead man was killed a second time")
	}
}

func TestACommissionedKillingSaysWhatTheFeeBought(t *testing.T) {
	t.Parallel()
	w := morgue(t)
	victim := w.NPC("vittorio")
	seen := map[string]bool{}
	for _, id := range []string{"cheap", "professional", "specialist"} {
		tier, ok := contractTier(id)
		if !ok {
			t.Fatalf("no tier called %s", id)
		}
		method := w.killingMethod(victim, tier)
		if seen[method] {
			t.Fatal("two tiers read identically")
		}
		seen[method] = true
		if !strings.Contains(method, victim.Name) {
			t.Fatalf("a %s hit did not name who it was: %q", id, method)
		}
	}
}
