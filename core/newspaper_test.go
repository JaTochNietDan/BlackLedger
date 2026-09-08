package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestThePaperPrintsWhatTheCityCanSee(t *testing.T) {
	w := New(51)
	w.Player.Location = "club"
	w.Player.Respect = 40
	// A war, a killing and a robbery are all public events.
	w.Antagonize("bellandi", "russo", 100)
	w.Minute += 720
	w.FactionTurn()
	w.Kill("elena", "Shot on the steps of the exchange.")
	for seed := uint32(1); seed <= 200; seed++ {
		probe := New(51)
		probe.Player.Location, probe.Player.Respect, probe.RNG = "club", 40, seed
		if probe.Random() < probe.robberyOdds("club", probe.OwnHands()) {
			w.RNG = seed
			break
		}
	}
	_ = w.Rob("club")

	edition := w.Edition()
	if len(edition) == 0 {
		t.Fatal("a war, a killing and a robbery produced no news at all")
	}
	kinds := map[string]bool{}
	for _, s := range edition {
		kinds[s["kind"].(string)] = true
		if s["headline"] == "" || s["body"] == "" {
			t.Fatal("a story ran with no headline or no body")
		}
		if s["day"].(int) < 1 {
			t.Fatal("a story ran before the campaign began")
		}
	}
	for _, want := range []string{"killing", "robbery"} {
		if !kinds[want] {
			t.Fatalf("the paper never reported a %s: saw %v", want, kinds)
		}
	}
	// Most recent first.
	if len(edition) > 1 && edition[0]["minute"].(int) < edition[len(edition)-1]["minute"].(int) {
		t.Fatal("the paper is printed oldest first")
	}
}

func TestThePaperNeverPrintsWhatIsHidden(t *testing.T) {
	w := New(53)
	w.Player.Cash = 50000
	w.Player.Location = "market"
	w.Player.Contacts = 2
	if err := w.Commission("vittorio", "specialist"); err != nil {
		t.Fatal(err)
	}
	w.RetaliationFrom("bellandi")
	if len(w.Contracts) == 0 || len(w.Plots) == 0 {
		t.Fatal("the test did not create anything hidden")
	}
	body, err := json.Marshal(w.Edition())
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, c := range w.Contracts {
		if strings.Contains(text, c.ID) || strings.Contains(text, c.Target) {
			t.Fatal("the paper printed a commissioned killing that has not happened")
		}
	}
	for _, plot := range w.Plots {
		if strings.Contains(text, plot.ID) {
			t.Fatal("the paper printed a plot that has not happened")
		}
	}
	if strings.Contains(strings.ToLower(text), "contract") {
		t.Fatal("the paper referred to a contract")
	}
}

func TestThePlayersOwnCrimesAreReportedWithoutTheirName(t *testing.T) {
	var w *World
	for seed := uint32(1); seed <= 300; seed++ {
		probe := New(57)
		probe.Player.Location, probe.Player.Respect, probe.RNG = "club", 40, seed
		if probe.Random() < probe.robberyOdds("club", probe.OwnHands()) {
			w = New(57)
			w.Player.Location, w.Player.Respect, w.RNG = "club", 40, seed
			break
		}
	}
	if w == nil {
		t.Skip("no succeeding seed")
	}
	name := w.Player.Name
	if err := w.Rob("club"); err != nil {
		t.Fatal(err)
	}
	printed := false
	for _, s := range w.Edition() {
		if s["kind"] != "robbery" {
			continue
		}
		printed = true
		text := s["headline"].(string) + " " + s["body"].(string)
		if strings.Contains(text, name) {
			t.Fatal("the paper named the player as the robber")
		}
		if !strings.Contains(text, "no arrest") {
			t.Fatalf("the report did not read as an unsolved crime: %q", text)
		}
	}
	if !printed {
		t.Fatal("a robbery was never reported")
	}
	// The player's own private record still tells them what they did.
	own := false
	for _, r := range w.History {
		if strings.Contains(r.Title, "Taken from") {
			own = true
		}
	}
	if !own {
		t.Fatal("the player's own history lost the robbery")
	}
}

func TestANewPersonReadsTheirOwnEdition(t *testing.T) {
	w := New(59)
	w.Kill("elena", "Shot at the exchange.")
	if len(w.Edition()) == 0 {
		t.Fatal("nothing was printed in the first life")
	}
	w.Die("Test death")
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	w = next
	for _, s := range w.Edition() {
		if s["minute"].(int) < 0 {
			t.Fatal("impossible story")
		}
	}
	if len(w.Edition()) != 0 {
		t.Fatal("a new person inherited the previous life's newspaper archive")
	}
	// The paper keeps running for them.
	w.Report("war", "SOMETHING HAPPENS", "It happened.")
	if len(w.Edition()) != 1 {
		t.Fatal("the paper stopped printing for the new person")
	}
}

func TestTheArchiveIsBounded(t *testing.T) {
	w := New(61)
	for i := 0; i < newsCapacity*3; i++ {
		w.Report("war", "HEADLINE", "Body.")
	}
	if len(w.News) > newsCapacity {
		t.Fatalf("the archive grew to %d stories", len(w.News))
	}
	if len(w.Edition()) != newsCapacity {
		t.Fatalf("the edition prints %d stories", len(w.Edition()))
	}
}

func TestADeadProtagonistsArrangementsDieWithThem(t *testing.T) {
	w := New(63)
	w.Player.Cash = 50000
	w.Player.Location = "market"
	w.Player.Contacts = 2
	if err := w.Commission("vittorio", "cheap"); err != nil {
		t.Fatal(err)
	}
	if len(w.Contracts) == 0 {
		t.Fatal("no contract was booked")
	}
	w.Die("Test death")
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if len(next.Contracts) != 0 {
		t.Fatalf("%d contracts survived the person who paid for them", len(next.Contracts))
	}
}

// A headline on its own reads like a log line. Every kind of story the city can
// file gets a standfirst and a desk that filed it, or the paper prints a line
// that says nothing.
func TestEveryKindOfStoryIsSetLikeANewspaper(t *testing.T) {
	kinds := []string{"police", "politics", "robbery", "business", "killing", "attack",
		"war", "seizure", "split", "recovery", "collapse", "attempt", "arrest"}
	w := proprietor(t)
	generic := 0
	for _, kind := range kinds {
		w.Report(kind, "SOMETHING HAPPENED AT BLUEBIRD LAUNDRY", "A thing occurred.")
		s := w.News[len(w.News)-1]
		stand, desk := w.standfirst(s), deskFor(kind)
		if stand == "" || desk == "" {
			t.Fatalf("%q was filed with standfirst %q and byline %q", kind, stand, desk)
		}
		if stand == "The Herald understands the position remains unchanged." {
			generic++
			t.Errorf("%q falls through to the generic standfirst", kind)
		}
		if desk == "Staff report" {
			t.Errorf("%q was filed by nobody in particular", kind)
		}
	}
	if generic > 0 {
		t.Fatalf("%d of %d story kinds have nothing to say under the headline", generic, len(kinds))
	}
}

func TestThePaperCarriesADate(t *testing.T) {
	if got := Dateline(480); got != "Tuesday, March 3, 1953" {
		t.Fatalf("day one is %q", got)
	}
	// It advances a day at a time and rolls over the end of a month.
	if got := Dateline(29 * 1440); got != "Wednesday, April 1, 1953" {
		t.Fatalf("day thirty is %q", got)
	}
	seen := map[string]bool{}
	for day := 0; day < 200; day++ {
		d := Dateline(day * 1440)
		if d == "" || seen[d] {
			t.Fatalf("day %d reads %q", day, d)
		}
		seen[d] = true
	}
}
