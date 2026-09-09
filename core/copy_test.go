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

// A splinter's name used to carry its own article, so every sentence it began
// started in lower case. Three separate sites were fixed one at a time before
// it was obvious that the name was the fault.

func TestASplinterIsNamedLikeEveryOtherFamily(t *testing.T) {
	for _, form := range splinterForms {
		if strings.HasPrefix(form, "the ") {
			t.Fatalf("a splinter would be named %q, which no sentence can begin with", form)
		}
	}
	// Saves written before tonight still hold the old names, so the repair
	// has to stay.
	if Leads("the Rizzo Crew") != "The Rizzo Crew" {
		t.Fatal("an old save's name can still begin a sentence in lower case")
	}
}

// Seen in the live game's result banner: "The a professional you paid $2501 to
// reach Elena Russo did not finish it." The contract tiers are labelled for
// buttons — "A professional", "Someone who needs the money" — and a label
// carrying its own article cannot be given another one.

func TestAHiredHandIsNamedInASentence(t *testing.T) {
	for _, tier := range contractTiers {
		if tier.Noun == "" {
			t.Fatalf("%q has no form fit for a sentence", tier.Label)
		}
		line := upper1(tier.Noun) + " you paid $200 for reached them."
		for _, bad := range []string{"The a ", "The A ", "The someone ", "The Someone "} {
			if strings.Contains(line, bad) {
				t.Fatalf("a ledger line reads %q", line)
			}
		}
		if strings.ToUpper(line[:1]) != line[:1] {
			t.Fatalf("a ledger line begins in lower case: %q", line)
		}
	}
}

// Read out of a real save while driving the arms loop: "5 crates of Crated arms
// for $1100, at $220 each." The market lists a good by a name fit for a price
// board, and that name cannot follow a count of units.
func TestAGoodReadsCorrectlyAfterACount(t *testing.T) {
	w := New(17)
	for _, g := range w.Goods {
		line := "5 " + g.Unit + "s of " + g.InBulk()
		if strings.Contains(line, "of Crated") || strings.Contains(line, "of Untaxed") ||
			strings.Contains(line, "of Moonshine") {
			t.Fatalf("the ledger reads %q", line)
		}
		if g.InBulk() == "" {
			t.Fatalf("%q has no form fit to follow a count", g.Name)
		}
	}
	if got := w.Good("arms").InBulk(); got != "arms" {
		t.Fatalf("crated arms read as %q after a count", got)
	}
}

// Found by reading a second campaign's paper end to end, after eight copy
// fixes had landed. Two of these three faults were introduced by those fixes.

func TestATitledOfficeTakesNoArticle(t *testing.T) {
	w := New(18)
	f := &w.Factions[0]
	// The fix for "They were Lieutenant" gave every role an article, and this
	// is what that did to a role that already carries its own complement.
	n := &NPC{ID: "z", Name: "Elena Russo", Role: "Head of the " + f.Name, Faction: f.ID}
	if got := describeStanding(n, w); strings.Contains(got, "a head of") {
		t.Fatalf("the city says %q", got)
	}
	if got := describeStanding(n, w); !strings.Contains(got, "They were head of") {
		t.Fatalf("the city says %q", got)
	}
	// An ordinary job still takes one.
	n.Role = "Lieutenant"
	if got := describeStanding(n, w); !strings.Contains(got, "a lieutenant") {
		t.Fatalf("an ordinary role lost its article: %q", got)
	}
}

func TestThePaperCountsInWords(t *testing.T) {
	// "There were twice such incidents before the day was out." A frequency
	// where a count belongs.
	w := New(19)
	w.News = nil
	for i := 0; i < 3; i++ {
		w.Report("attack", "DAMAGE AT SAINT AGNES", "Saint Agnes was attacked overnight.")
	}
	body := w.News[len(w.News)-1].Body
	if strings.Contains(body, "twice such") || strings.Contains(body, "three times such") {
		t.Fatalf("the paper reads %q", body)
	}
	if !strings.Contains(body, "three such incidents") {
		t.Fatalf("the paper reads %q", body)
	}
	// And a robbery still takes the frequency, because that sentence wants one.
	w.News = nil
	for i := 0; i < 2; i++ {
		w.Report("robbery", "ROBBERY IN SAINT AGNES", "A man was robbed.")
	}
	if body := w.News[0].Body; !strings.Contains(body, "twice in the same day") {
		t.Fatalf("the paper reads %q", body)
	}
}

func TestAnObituaryCountsInWords(t *testing.T) {
	if spelled(2) != "two" || spelled(12) != "twelve" {
		t.Fatal("the paper is writing small numbers as figures")
	}
	if spelled(41) != "41" {
		t.Fatalf("the paper spelled out a large number as %q", spelled(41))
	}
}

// "DETECTIVE HARLOW TAKES OVER THE BLUE HOUR" — the city detective had walked
// off his beat and seized a casino. The guard on that excluded the officials
// and the heads of organizations, but not the people holding the city's
// standing jobs.
func TestSomebodyWithAJobDoesNotSeizeACasino(t *testing.T) {
	w := New(20)
	for _, r := range roles {
		n := w.NPC(r.Seed)
		if n == nil {
			continue
		}
		n.Role, n.Location = r.Title, r.Where
		n.Skill, n.Ambition = 100, 100
		if !w.keepsPost(n) {
			t.Fatalf("the %s is free to walk off and take over a business", r.Title)
		}
	}
	// Somebody with no job still can, or nothing in the city ever changes.
	nobody := &NPC{ID: "nobody", Name: "Nobody", Skill: 100, Ambition: 100, Post: "bar"}
	if w.keepsPost(nobody) {
		t.Fatal("a person with no job and no rank is pinned to the spot")
	}
}

// "You know Mayor Ellis Crane now: They answered to nobody." An official has no
// faction, so the city described the mayor by what he did not belong to rather
// than by the office he held.
func TestTheCityNamesAnOfficeWhenSomebodyHoldsOne(t *testing.T) {
	w := New(24)
	for _, o := range Officials() {
		n := w.NPC(o.ID)
		if n == nil {
			t.Fatalf("%s is not in the city", o.Name)
		}
		got := describeStanding(n, w)
		if strings.Contains(got, "answered to nobody") {
			t.Fatalf("the city says the %s %q", o.Role, got)
		}
		if !strings.Contains(strings.ToLower(got), strings.ToLower(o.Role)) {
			t.Fatalf("the city describes the %s as %q", o.Role, got)
		}
		// Exactly one person holds each of these, so it is "the".
		if strings.Contains(got, "a police") || strings.Contains(got, "a mayor") {
			t.Fatalf("a unique office took an indefinite article: %q", got)
		}
	}
}

// Read out of a real save after taking an organization: "You hold 4 of its
// premises, 8 of its people stayed and 0 would not." A zero written as a figure
// reads like a report from a machine.
func TestTakingAnOrganizationReadsLikeProse(t *testing.T) {
	if spelled(0) != "no" {
		t.Fatalf("zero reads as %q", spelled(0))
	}
	for _, n := range []int{0, 1, 4, 8, 12} {
		if s := spelled(n); strings.ContainsAny(s, "0123456789") {
			t.Fatalf("%d reads as %q in prose", n, s)
		}
	}
	// And the first fix for it produced "You hold no of its premises", because
	// "no" reads in "no people" and never after a preposition.
	if countOf(0) != "none" {
		t.Fatalf("a count of nothing followed by \"of\" reads as %q", countOf(0))
	}
	for _, n := range []int{0, 1, 4, 12} {
		if line := "You hold " + countOf(n) + " of its premises."; strings.Contains(line, " no of ") {
			t.Fatalf("the record reads %q", line)
		}
	}

}
