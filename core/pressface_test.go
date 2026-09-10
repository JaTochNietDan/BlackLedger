package core

import (
	"strings"
	"testing"
)

// "When I went after Tila myself, the newspaper info on the attempt showed the
// wrong portrait."
//
// The paper worked out who a story was about by reading its own headline back
// and returning the first person in the city whose name appeared in it. A
// headline that names two people — which is most of what a paper about people
// going after each other prints — showed whichever of them the city happened to
// list first. The core knew who the story was about when it filed it and threw
// that away.

func newsroom(t *testing.T) (*World, *NPC, *NPC) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	var a, b *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) {
			continue
		}
		if a == nil {
			a = n
			continue
		}
		b = n
		break
	}
	if a == nil || b == nil {
		t.Fatal("this city does not have two people in it")
	}
	return w, a, b
}

func TestAStoryIsAboutWhoItIsAboutRatherThanWhoIsListedFirst(t *testing.T) {
	w, first, second := newsroom(t)
	// A headline naming both of them, filed about the second one — which is
	// exactly the shape of "you went after somebody".
	headline := strings.ToUpper(first.Name) + " GOES AFTER " + strings.ToUpper(second.Name)
	w.ReportAbout("attack", headline, "Two cars and a lot of noise.", second.ID)
	story := w.News[len(w.News)-1]
	if got := w.SubjectOf(story); got.ID != second.ID {
		t.Fatalf("a story about %s is illustrated with %s", second.Name, got.Name)
	}
	// And the other way round, on the same words.
	w.ReportAbout("attack", headline+" AGAIN", "Two cars and a lot of noise.", first.ID)
	if got := w.SubjectOf(w.News[len(w.News)-1]); got.ID != first.ID {
		t.Fatalf("a story about %s is illustrated with %s", first.Name, got.Name)
	}
}

// A story filed with nobody named still reads its own headline, because that is
// how every story written before this worked and they are still in the archive.
func TestAStoryWithNobodyNamedStillFindsAFace(t *testing.T) {
	w, who, _ := newsroom(t)
	w.Report("attack", "SOMEBODY WENT AFTER "+strings.ToUpper(who.Name), "It was over quickly.")
	if got := w.SubjectOf(w.News[len(w.News)-1]); got.ID != who.ID {
		t.Fatalf("a headline naming %s alone is illustrated with %q", who.Name, got.Name)
	}
	// And a story about nowhere and nobody is a picture of the city.
	w.Report("politics", "THE CITY IS QUIET", "Nothing happened.")
	if got := w.SubjectOf(w.News[len(w.News)-1]); got.Kind != "city" {
		t.Fatalf("a story about nothing in particular is illustrated with %q", got.Name)
	}
}

// The premises rule stays: a headline that names a room is a picture of the
// room, unless the story says otherwise, because "KILLED AT THE MONARCH" is a
// photograph of The Monarch.
func TestAStoryAboutSomebodyBeatsThePlaceInItsHeadline(t *testing.T) {
	w, who, _ := newsroom(t)
	headline := strings.ToUpper(who.Name) + " SHOT AT THE MONARCH"
	w.Report("attack", headline, "Outside, in the road.")
	if got := w.SubjectOf(w.News[len(w.News)-1]); got.Kind != "place" {
		t.Fatalf("a headline naming a room was illustrated with %q rather than the room", got.Name)
	}
	w.ReportAbout("attack", headline+" LAST NIGHT", "Outside, in the road.", who.ID)
	if got := w.SubjectOf(w.News[len(w.News)-1]); got.ID != who.ID {
		t.Fatalf("a story filed about %s was illustrated with %q", who.Name, got.Name)
	}
}

// The stories a paper about people going after each other actually prints must
// all say whose they are. This is the list, checked against the core rather
// than against the prose: an attempt, a killing, an arrest and an obituary.
func TestTheStoriesAboutPeopleAllSayWhoTheyAreAbout(t *testing.T) {
	w, mark, _ := newsroom(t)
	w.Player.Cash, w.Player.Health = 5000, 100
	w.Player.Location = mark.Location
	before := len(w.News)

	// An attempt that fails, which is the one the note was about: going after
	// somebody yourself and reading about it afterwards.
	w.itWentWrong(mark, Hand{}, "The Monarch", nil)
	// And a death, which files a killing and later an obituary.
	w.Kill(mark.ID, "shot in the street")

	named, checked := 0, 0
	for _, s := range w.News[before:] {
		checked++
		if s.About == "" {
			t.Errorf("a %q story about people does not say whose it is: %q", s.Kind, s.Headline)
			continue
		}
		named++
		if got := w.SubjectOf(s); got.Kind != "person" {
			t.Errorf("a %q story naming somebody is illustrated with a %s", s.Kind, got.Kind)
		}
	}
	if checked == 0 {
		t.Fatal("going after somebody and killing them filed no stories at all")
	}
	t.Logf("%d of %d stories filed about a person say who", named, checked)
}
