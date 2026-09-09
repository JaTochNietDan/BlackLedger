package core

import (
	"strings"
	"testing"
)

// Read end to end, a campaign's paper printed "the Rizzo Crew has ceased to
// operate" at the head of a paragraph, "Cesare Ferro's people has had a poor
// few weeks", "Crated arms is fetching more than it did", "They were
// Lieutenant", and SUNDAY IN BELLWETHER under a dateline reading Monday. None
// of it was wrong about the world. All of it was wrong about English.

func TestASentenceDoesNotBeginInLowerCase(t *testing.T) {
	if got := Leads("the Rizzo Crew"); got != "The Rizzo Crew" {
		t.Fatalf("a sentence begins %q", got)
	}
	// Mid-sentence the name is correct as it stands, so nothing else moves.
	for _, name := range []string{"Bellandi Family", "Cesare Ferro's people"} {
		if Leads(name) != name {
			t.Fatalf("%q was changed to %q for no reason", name, Leads(name))
		}
	}
}

func TestAPluralNameTakesAPluralVerb(t *testing.T) {
	if Agree("Cesare Ferro's people", "has", "have") != "have" {
		t.Fatal("the paper says a family has had a poor few weeks when it is a plural name")
	}
	if Agree("Bellandi Family", "has", "have") != "has" {
		t.Fatal("a singular family took a plural verb")
	}
}

func TestGrammarSurvivesAnOldSave(t *testing.T) {
	// Whether a good's name is plural is a fact about the word, not about the
	// market, so it must not be read out of the save. A campaign begun before
	// anybody noticed still has to read correctly.
	w := New(11)
	for i := range w.Goods {
		w.Goods[i].Price = w.Goods[i].Base * 2
	}
	arms := w.Good("arms")
	if arms.Agrees("is", "are") != "are" {
		t.Fatal("crated arms is fetching more than it did")
	}
	if w.Good("moonshine").Agrees("is", "are") != "is" {
		t.Fatal("moonshine are fetching more than they did")
	}
}

func TestTheSundayPageRunsOnASunday(t *testing.T) {
	found := 0
	for day := 0; day < 28; day++ {
		w := New(12)
		w.Minute = day * 1440
		sunday := false
		for _, b := range w.cityPage() {
			if strings.Contains(b.headline, "SUNDAY") {
				sunday = true
			}
		}
		if sunday {
			found++
			if Weekday(w.Minute) != "Sunday" {
				t.Fatalf("the Sunday page ran under a %s dateline", Weekday(w.Minute))
			}
		}
	}
	if found != 4 {
		t.Fatalf("four weeks held %d Sundays", found)
	}
}

func TestAnObituarySaysWhatTheyWere(t *testing.T) {
	w := New(13)
	n := &NPC{ID: "x", Name: "Lorenz Zanetti", Role: "Lieutenant"}
	_, body := w.obituary(n)
	if !strings.Contains(body, "They were a lieutenant") {
		t.Fatalf("the obituary reads %q", body)
	}
	n.Role = "Enforcer"
	if _, body = w.obituary(n); !strings.Contains(body, "They were an enforcer") {
		t.Fatalf("the obituary reads %q", body)
	}
}

// The obituary was not the only place the city said what somebody was. The
// killing story said it too, from a different function, and still read "They
// were Lieutenant" after the obituary was fixed.

func TestEveryPlaceThatSaysWhatSomebodyWasSaysItInEnglish(t *testing.T) {
	w := New(14)
	f := &w.Factions[0]
	n := &NPC{ID: "y", Name: "Zora Erdos", Role: "Lieutenant", Faction: f.ID}
	if got := describeStanding(n, w); !strings.HasPrefix(got, "They were a lieutenant of ") {
		t.Fatalf("the city says %q", got)
	}
	n.Role = ""
	if got := describeStanding(n, w); got != "They were one of "+f.Name {
		t.Fatalf("somebody with no role reads as %q", got)
	}
}
