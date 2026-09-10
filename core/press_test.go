package core

import "testing"

func pressroom(t *testing.T) *World {
	t.Helper()
	w := New(233)
	w.MigrateLivingWorld()
	return w
}

func TestAStoryAboutPremisesIsAPictureOfThePremises(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	s := Story{Headline: "ROBBERY AT THE MONARCH", Kind: "robbery"}
	subject := w.SubjectOf(s)
	if subject.Kind != "place" || subject.ID != "club" || subject.Name != "The Monarch" {
		t.Fatalf("the picture was of %+v", subject)
	}
}

func TestAStoryAboutSomebodyIsAPictureOfThem(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	person := w.NPC("vittorio")
	s := Story{Headline: upper(person.Name) + " KILLED", Kind: "killing"}
	subject := w.SubjectOf(s)
	if subject.Kind != "person" || subject.ID != person.ID || subject.Name != person.Name {
		t.Fatalf("the picture was of %+v", subject)
	}
}

func TestPremisesWinWhenAHeadlineNamesBoth(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	person := w.NPC("vittorio")
	s := Story{Headline: upper(person.Name) + " KILLED AT THE MONARCH", Kind: "killing"}
	if subject := w.SubjectOf(s); subject.Kind != "place" {
		t.Fatalf("a killing at a named place was a picture of %+v", subject)
	}
}

func TestTheDeadStillGetTheirPicture(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	person := w.NPC("vittorio")
	name := person.Name
	person.Dead = true
	s := Story{Headline: upper(name) + " KILLED", Kind: "killing"}
	if subject := w.SubjectOf(s); subject.Kind != "person" || subject.Name != name {
		t.Fatalf("a dead man's story was a picture of %+v", subject)
	}
}

func TestAnythingElseIsAPictureOfTheCity(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	for _, headline := range []string{"OPEN WAR ON THE WATERFRONT", "POLICE PRESSURE MOUNTS", ""} {
		subject := w.SubjectOf(Story{Headline: headline, Kind: "war"})
		if subject.Kind != "city" || subject.Name != "Bellwether" {
			t.Fatalf("%q was a picture of %+v", headline, subject)
		}
	}
}

func TestEveryStoryTheCityFilesGetsAPicture(t *testing.T) {
	t.Parallel()
	// A city left to run, and every story it printed checked: none of them may
	// come back without something to draw.
	const runs, days = 60, 120
	kinds := map[string]int{}
	stories := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w := New(seed)
		w.MigrateLivingWorld()
		w.WorldRNG = seed * 2654435761
		for day := 0; day < days; day++ {
			w.FactionTurn()
			w.FamilyDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.Minute += 1440
		}
		for _, s := range w.News {
			stories++
			subject := w.SubjectOf(s)
			if subject.Kind == "" || subject.Name == "" {
				t.Fatalf("%q had nothing to draw", s.Headline)
			}
			kinds[subject.Kind]++
		}
	}
	if stories == 0 {
		t.Fatal("no city printed anything")
	}
	if len(kinds) < 3 {
		t.Fatalf("%d stories produced only these subjects: %v", stories, kinds)
	}
	if kinds["city"]*2 > stories {
		t.Fatalf("%d of %d stories fell back to a picture of the city", kinds["city"], stories)
	}
	t.Logf("%d stories across %d cities: %v", stories, runs, kinds)
}

func TestTheSameStoryIsAlwaysTheSamePicture(t *testing.T) {
	t.Parallel()
	w := pressroom(t)
	s := Story{Headline: "ROBBERY AT THE MONARCH", Kind: "robbery"}
	first := w.SubjectOf(s)
	for i := 0; i < 50; i++ {
		if w.SubjectOf(s) != first {
			t.Fatal("the same story produced two different pictures")
		}
	}
}
