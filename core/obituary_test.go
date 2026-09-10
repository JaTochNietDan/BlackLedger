package core

import (
	"strings"
	"testing"
)

// An obituary is the only place in this game where a person is described as a
// life rather than as a threat. It has to arrive the morning after, once, and
// it must never say who arranged anything.

func obituariesIn(w *World) []Story {
	out := []Story{}
	for _, s := range w.News {
		if s.Kind == "obituary" {
			out = append(out, s)
		}
	}
	return out
}

func anyLieutenant(w *World) *NPC {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Rank >= RankLieutenant && n.Faction != "" {
			return n
		}
	}
	return nil
}

func TestTheObituaryComesTheMorningAfter(t *testing.T) {
	t.Parallel()
	w := New(91)
	n := anyLieutenant(w)
	if n == nil {
		t.Skip("no lieutenant in this world")
	}
	w.Kill(n.ID, "Shot in the street.")

	// Nothing on the day itself: the killing was today's news.
	w.ObituaryDay()
	if got := obituariesIn(w); len(got) != 0 {
		t.Fatalf("an obituary ran the same day as the killing: %q", got[0].Headline)
	}

	w.Minute += 1440
	w.ObituaryDay()
	got := obituariesIn(w)
	if len(got) != 1 {
		t.Fatalf("expected one obituary the morning after, got %d", len(got))
	}
	if !strings.Contains(got[0].Headline, strings.ToUpper(n.Name)) {
		t.Errorf("the obituary was for %q, not for %s", got[0].Headline, n.Name)
	}

	// And never again.
	w.Minute += 1440
	w.ObituaryDay()
	if len(obituariesIn(w)) != 1 {
		t.Error("the paper ran the same obituary twice")
	}
}

func TestAnObituaryNeverSaysWhoArrangedIt(t *testing.T) {
	t.Parallel()
	w := New(92)
	victim, killer := anyLieutenant(w), (*NPC)(nil)
	if victim == nil {
		t.Skip("no lieutenant in this world")
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.ID != victim.ID && n.Name != "" {
			killer = n
			break
		}
	}
	w.KillBy(victim.ID, killer, "a debt nobody wrote down")
	w.Minute += 1440
	w.ObituaryDay()
	for _, s := range obituariesIn(w) {
		if killer != nil && strings.Contains(s.Body, killer.Name) {
			t.Errorf("the obituary named %s, who did it: %q", killer.Name, s.Body)
		}
		for _, word := range []string{"murder", "arranged", "ordered by", "on behalf of"} {
			if strings.Contains(strings.ToLower(s.Body), word) {
				t.Errorf("the obituary said %q, which the paper cannot know: %q", word, s.Body)
			}
		}
	}
}

func TestNobodyGetsAColumnForBeingNobody(t *testing.T) {
	t.Parallel()
	w := New(93)
	// A soldier with no standing, nobody has met, running nothing.
	unknown := &NPC{ID: "nobody-at-all", Name: "A Nobody", Rank: RankSoldier}
	if w.mourned(unknown) {
		t.Error("the paper would run an obituary for somebody nobody had heard of")
	}
	// And somebody with a title always gets one.
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Rank >= RankLieutenant {
			if !w.mourned(n) {
				t.Errorf("%s stood over other people and got no obituary", n.Name)
			}
			break
		}
	}
}

func TestAnObituaryCostsTheCityNothing(t *testing.T) {
	t.Parallel()
	// The killing was counted yesterday. Counting it again would have the city
	// look hardest at the people who are mourned most.
	if scrutinyWeight["obituary"] != 0 {
		t.Fatalf("an obituary carries %d scrutiny", scrutinyWeight["obituary"])
	}
}

func TestAnObituarySaysWhatThePersonWas(t *testing.T) {
	t.Parallel()
	w := New(94)
	n := anyLieutenant(w)
	if n == nil {
		t.Skip("no lieutenant in this world")
	}
	w.Kill(n.ID, "Found in the road.")
	w.Minute += 1440
	w.ObituaryDay()
	got := obituariesIn(w)
	if len(got) == 0 {
		t.Fatal("no obituary was written")
	}
	body := got[0].Body
	if len(body) < 40 {
		t.Errorf("the obituary is too short to be one: %q", body)
	}
	if !strings.Contains(body, "They were") && !strings.Contains(body, "They answered") {
		t.Errorf("the obituary never says what they were: %q", body)
	}
}

func TestTheObituaryDoesNotPrintTheFamilyTwice(t *testing.T) {
	t.Parallel()
	// "They were Russo boss of Russo Outfit" is what came out the first time.
	if sharesAName("Russo boss", "Russo Outfit") != true {
		t.Error("a role carrying the family's name was not recognised")
	}
	if sharesAName("Driver", "Russo Outfit") != false {
		t.Error("a role with nothing to do with the family was treated as carrying it")
	}
}

func TestOnePersonIsNotOnePeople(t *testing.T) {
	t.Parallel()
	w := New(95)
	n := anyLieutenant(w)
	if n == nil {
		t.Skip("no lieutenant in this world")
	}
	w.Kill(n.ID, "Shot.")
	w.Minute += 1440
	w.ObituaryDay()
	for _, s := range obituariesIn(w) {
		if strings.Contains(s.Body, "1 people") {
			t.Errorf("the paper counted %q", "1 people")
		}
	}
}

func TestTheObituaryNamesTheirOwnPlaceFirst(t *testing.T) {
	t.Parallel()
	// It named a holding on the other side of the city while the man was found
	// at the door of one he ran, which reads as the paper picking a building at
	// random — which is what it was doing.
	w := New(96)
	n := anyLieutenant(w)
	if n == nil {
		t.Skip("no lieutenant in this world")
	}
	// Put a holding of their family under them.
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil {
			prop.Owner = n.Faction
			n.Location = l.ID
			w.Kill(n.ID, "Shot.")
			w.Minute += 1440
			w.ObituaryDay()
			for _, s := range obituariesIn(w) {
				if strings.Contains(s.Body, "carries on") && !strings.Contains(s.Body, l.Name) {
					t.Errorf("they were found at %s and the paper named somewhere else: %q", l.Name, s.Body)
				}
			}
			return
		}
	}
}
