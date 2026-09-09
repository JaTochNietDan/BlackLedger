package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The naming scans pointed the harness at awkward NAMES. This points it at
// awkward STATES: a family holding nothing, a family with nobody left, a city
// down to one organization, a player with nothing, and a world run five times
// longer than anything measured before.
//
// The checks are the faults already found tonight, turned into a set that runs
// against everything the city writes rather than against one sentence at a
// time.

var (
	oneOfSomethings = regexp.MustCompile(`\b1 [a-z]+s\b`)
	// Go's regexp has no backreferences, so the doubled words are spelled out.
	doubledWord = regexp.MustCompile(`(?i)\b(the the|a a|of of|to to|and and|is is|has has|in in|at at)\b`)
)

// flaws returns what is wrong with a passage, or nothing.
func flaws(text string) []string {
	out := []string{}
	if strings.Contains(text, "  ") {
		out = append(out, "a doubled space")
	}
	if strings.Contains(text, " .") || strings.Contains(text, " ,") {
		out = append(out, "a space before punctuation")
	}
	if m := doubledWord.FindString(text); m != "" {
		out = append(out, "a doubled word "+m)
	}
	if m := oneOfSomethings.FindString(text); m != "" {
		out = append(out, "a plural after one: "+m)
	}
	if strings.Contains(text, "$0 ") || strings.HasSuffix(text, "$0") || strings.Contains(text, "$0.") {
		out = append(out, "a figure of zero")
	}
	if m := disagrees(text); m != "" {
		out = append(out, "a disagreement: "+m)
	}
	for _, sentence := range strings.Split(text, ". ") {
		if strings.HasPrefix(sentence, "the ") {
			out = append(out, "a sentence beginning in lower case")
			break
		}
	}
	return out
}

func scan(t *testing.T, w *World, what string) int {
	t.Helper()
	read := 0
	for _, text := range writings(w) {
		if strings.TrimSpace(text) == "" {
			continue
		}
		read++
		for _, flaw := range flaws(text) {
			t.Errorf("%s: %s\n  %s", what, flaw, text)
		}
	}
	return read
}

// A family that holds nothing, and one with nobody left to hold it.
func TestACityWithHollowedOutFamiliesStillReads(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 47, 88, 103, 219} {
		w := withFamily(seed, "estate:hollow", "Vera Kohl's people", "Vera Kohl")
		// Hollowed out: it owns nothing at all.
		for id := range w.Properties {
			if w.Properties[id].Owner == "estate:hollow" {
				w.Properties[id].Owner = "independent"
			}
		}
		// And the other one loses everybody.
		for i := range w.NPCs {
			if w.NPCs[i].Faction == "russo" {
				w.Kill(w.NPCs[i].ID, "A bad week.")
			}
		}
		live(t, w, 4000)
		read += scan(t, w, fmt.Sprintf("seed %d, hollowed families", seed))
	}
	if read == 0 {
		t.Fatal("nothing was written, so this proved nothing")
	}
	t.Logf("read %d passages", read)
}

// One organization left standing in the whole city.
func TestACityDownToOneOrganizationStillReads(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := withFamily(seed, "estate:last", "Bruno Duarte's people", "Bruno Duarte")
		w.Dissolve("bellandi")
		w.Dissolve("russo")
		if len(w.Factions) != 1 {
			t.Fatalf("seed %d: expected one family, got %d", seed, len(w.Factions))
		}
		live(t, w, 4000)
		read += scan(t, w, fmt.Sprintf("seed %d, one family left", seed))
	}
	t.Logf("read %d passages", read)
}

// A player with nothing at all, which is where every life begins and where a
// bad one ends.
func TestAPlayerWithNothingStillReads(t *testing.T) {
	read := 0
	for _, seed := range []uint32{31, 88, 219} {
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 0, 0, 20
		live(t, w, 2000)
		w.Player.Cash, w.Player.Respect = 0, 0
		read += scan(t, w, fmt.Sprintf("seed %d, a player with nothing", seed))
	}
	t.Logf("read %d passages", read)
}

// Five times longer than anything measured before: does the city repeat itself,
// run out of names, or degrade?
func TestALongCityDoesNotRepeatItself(t *testing.T) {
	w := withFamily(59, "estate:long", "Otto Reiss's people", "Otto Reiss")
	live(t, w, 20000)
	read := scan(t, w, "20000 half-hours")
	tally := func(of []string) (int, string, int) {
		counts := map[string]int{}
		for _, s := range of {
			counts[s]++
		}
		worst, worstOf := 0, ""
		for s, n := range counts {
			if n > worst {
				worst, worstOf = n, s
			}
		}
		return len(counts), worstOf, worst
	}
	ledger := make([]string, 0, len(w.History))
	for _, r := range w.History {
		ledger = append(ledger, r.Title)
	}
	paper := make([]string, 0, len(w.News))
	for _, s := range w.News {
		paper = append(paper, s.Headline)
	}
	sort.Strings(ledger)
	distinctLedger, worstLedger, ledgerRepeats := tally(ledger)
	distinctPaper, worstPaper, paperRepeats := tally(paper)
	t.Logf("read %d passages: %d ledger entries (%d distinct), %d stories (%d distinct), %d families, %d people",
		read, len(w.History), distinctLedger, len(w.News), distinctPaper, len(w.Factions), len(w.NPCs))
	t.Logf("most repeated: ledger %q x%d, paper %q x%d", worstLedger, ledgerRepeats, worstPaper, paperRepeats)

	// The ledger is deliberately not asserted on here. This harness keeps a
	// broke player alive for four hundred game days, so the bills fail every
	// midnight and one entry dominates — that is what I built, not what the
	// game does. A real campaign ends.
	//
	// The paper is the living world's own voice and owes nothing to the
	// player standing still, so it is fair to ask whether the city runs out of
	// things to say over four hundred days.
	if len(w.News) > 20 && paperRepeats*2 > len(w.News) {
		t.Errorf("the paper prints one headline more than half the time: %q x%d of %d",
			worstPaper, paperRepeats, len(w.News))
	}
	if distinctPaper < 8 {
		t.Errorf("the city ran out of things to print: only %d distinct headlines in %d stories",
			distinctPaper, len(w.News))
	}
}
