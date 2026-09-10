package core

import "testing"

// A trade settled at midnight that depends on somebody being there. Two of them
// were found this way — a garage and a forecourt, both billing whoever happened
// to be standing on the spot at the exact minute the day turned over, which is
// the minute a counter is dark. These are the others in the day loop that read
// where people are, measured against the city's own clock rather than reasoned
// about.

// nightsWithAHand counts the days on which the city's own gambling actually
// found somebody on a floor to play, running the clock the way the game does.
func TestTheTablesFindSomebodyAtMidnight(t *testing.T) {
	t.Parallel()
	played, days := 0, 0
	for _, seed := range []uint32{5, 23, 61, 97, 181} {
		w := New(seed)
		w.District = 2
		w.Properties["casino"].Owner = "bellandi"
		if f := w.faction("bellandi"); f != nil {
			f.Cash = 20000
		}
		for day := 0; day < 7; day++ {
			w.RNG, w.WorldRNG = seed+uint32(day), seed+uint32(day)
			aDay(w)
			days++
			before := w.cityPurses()
			w.TableNight()
			if w.cityPurses() != before {
				played++
			}
		}
	}
	t.Logf("the tables found somebody on %d of %d nights", played, days)
	if played == 0 {
		t.Error("a week of midnights in five cities and the tables never found a soul")
	}
}

// A family takes somebody on at its own address, and only when it is short of
// people: a family that starts full never recruits, so the first version of
// this measured a fortnight of peace and called it a fault. It is short of
// people after a war, which is when the rule is supposed to matter — and the
// question that matters is whether anybody is standing at their door when the
// day turns over.
func TestAFamilyShortOfPeopleTakesSomebodyOn(t *testing.T) {
	t.Parallel()
	took, cities := 0, 0
	for _, seed := range []uint32{5, 23, 61, 97, 181} {
		w := New(seed)
		w.District = 2
		cities++
		f := w.faction("bellandi")
		if f == nil {
			t.Fatal("nobody to be short of people")
		}
		// A bad month: three of theirs are gone.
		lost := 0
		for _, n := range w.Members(f.ID) {
			if lost >= 3 || n.Rank >= RankLeader {
				continue
			}
			w.Kill(n.ID, "a bad month")
			lost++
		}
		if lost == 0 {
			t.Fatal("nobody to lose")
		}
		was := len(w.Members(f.ID))
		for day := 0; day < 14; day++ {
			w.RNG, w.WorldRNG = seed+uint32(day), seed+uint32(day)
			aDay(w)
			w.RecruitDay()
		}
		if len(w.Members(f.ID)) > was {
			took++
		}
	}
	t.Logf("%d of %d families short of people took somebody on inside a fortnight", took, cities)
	if took == 0 {
		t.Error("five families lost three people each and none of them replaced anybody in a fortnight")
	}
}
