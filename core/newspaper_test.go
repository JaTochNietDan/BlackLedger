package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestThePaperPrintsWhatTheCityCanSee(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	w := New(61)
	for i := 0; i < newsCapacity*3; i++ {
		w.Report("war", fmt.Sprintf("HEADLINE %d", i), "Body.")
	}
	if len(w.News) > newsCapacity {
		t.Fatalf("the archive grew to %d stories", len(w.News))
	}
	if len(w.Edition()) != newsCapacity {
		t.Fatalf("the edition prints %d stories", len(w.Edition()))
	}
}

func TestADeadProtagonistsArrangementsDieWithThem(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// The Herald was a rolling list with no back issues: one flat run of the
// current life's stories, oldest silently evicted, and a previous
// protagonist's era unreachable. A city that remembers should let somebody
// read what it remembered.

func TestThePaperKeepsItsBackIssues(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	w.News = nil
	// Three days of a life, then a death, then two days of the next.
	for day := 0; day < 3; day++ {
		w.Minute = day * 1440
		w.Report("business", "SOMETHING ON DAY "+string(rune('1'+day)), "b")
	}
	w.Life++
	for day := 3; day < 5; day++ {
		w.Minute = day * 1440
		w.Report("politics", "A NEW NAME ON DAY "+string(rune('1'+day)), "b")
	}
	issues := w.Editions()
	if len(issues) != 5 {
		t.Fatalf("five days of news came back as %d issues", len(issues))
	}
	// Newest first, and each issue knows whose era it belongs to.
	if issues[0]["day"].(int) != 5 || issues[len(issues)-1]["day"].(int) != 1 {
		t.Fatalf("the archive runs %v to %v", issues[0]["day"], issues[len(issues)-1]["day"])
	}
	mine, theirs := 0, 0
	for _, in := range issues {
		if in["mine"].(bool) {
			mine++
		} else {
			theirs++
		}
	}
	if mine != 2 || theirs != 3 {
		t.Fatalf("%d issues from this life and %d from before it", mine, theirs)
	}
	// A predecessor's paper is readable, which is the point.
	if issues[len(issues)-1]["mine"].(bool) {
		t.Fatal("the oldest issue is attributed to the current life")
	}
	for _, in := range issues {
		if in["dateline"].(string) == "" || in["count"].(int) == 0 {
			t.Fatalf("an issue came back as %v", in)
		}
	}
}

func TestTheBiggestStoryOfTheDayLeads(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	w.News = nil
	w.Minute = 3 * 1440
	w.Report("business", "TRADE STEADY AT THE MARKET", "b")
	w.Report("killing", "A MAN IS FOUND DEAD", "b")
	w.Report("robbery", "A TILL IS EMPTIED", "b")
	// The shop story came first and the last story is the smallest; neither
	// should lead over a killing.
	issues := w.Editions()
	if len(issues) != 1 {
		t.Fatalf("one day came back as %d issues", len(issues))
	}
	stories := issues[0]["stories"].([]map[string]any)
	if stories[0]["headline"] != "A MAN IS FOUND DEAD" {
		t.Fatalf("the day is led by %q", stories[0]["headline"])
	}
	if len(stories) != 3 {
		t.Fatalf("%d stories ran that day", len(stories))
	}
}

func TestTheArchiveIsWorthKeeping(t *testing.T) {
	t.Parallel()
	// Sixty stories was about three months, which is one life and nothing
	// across several. Whatever the cap is, it must hold a campaign.
	if newsCapacity < 200 {
		t.Fatalf("the paper keeps %d stories, which is not an archive", newsCapacity)
	}
	w := proprietor(t)
	w.News = nil
	for i := 0; i < newsCapacity+40; i++ {
		w.Minute = i * 720
		w.Report("business", fmt.Sprintf("STORY %d", i), "b")
	}
	if len(w.News) != newsCapacity {
		t.Fatalf("the paper is holding %d stories against a cap of %d", len(w.News), newsCapacity)
	}
	if len(w.Editions()) == 0 {
		t.Fatal("a full archive produced no issues")
	}
}

// Measured on a real campaign: twenty-one days of play produced one story,
// because every route into the paper was violence and that campaign had none.
// A paper that prints nothing for three weeks is not a paper.

func TestThePaperComesOutEveryDay(t *testing.T) {
	t.Parallel()
	const days = 40
	w := proprietor(t)
	w.News = nil
	for d := 0; d < days; d++ {
		w.Minute = d * 1440
		w.CivicDay()
	}
	issues := w.Editions()
	if len(issues) < days-1 {
		t.Fatalf("%d days of a quiet city produced %d issues", days, len(issues))
	}
	for _, in := range issues {
		if in["count"].(int) == 0 {
			t.Fatalf("day %v went to press empty", in["day"])
		}
	}
	// And it is never the same paper twice running.
	seen := ""
	same := 0
	for _, in := range issues {
		lead := in["stories"].([]map[string]any)[0]["headline"].(string)
		if lead == seen {
			same++
		}
		seen = lead
	}
	if same > len(issues)/2 {
		t.Fatalf("%d of %d issues led with the same headline as the day before", same, len(issues))
	}
	t.Logf("%d quiet days produced %d issues, none of them empty", days, len(issues))
}

func TestTheOrdinaryEditionCostsTheCityNothing(t *testing.T) {
	t.Parallel()
	// Filling the paper must not make the police look harder at anybody.
	w := proprietor(t)
	w.News = nil
	w.Attention = 0
	for d := 1; d <= 25; d++ {
		w.Minute = d * 1440
		w.CivicDay()
		w.ScrutinyDay()
	}
	if w.Scrutiny() != 0 {
		t.Fatalf("a month of ordinary editions put the city at %d scrutiny", w.Scrutiny())
	}
	if len(w.News) == 0 {
		t.Fatal("nothing was printed at all")
	}
}

func TestARealDayIsNotPaddedWithThePriceOfCoal(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	w.News = nil
	w.Minute = 5 * 1440
	w.Report("killing", "A MAN IS FOUND DEAD", "b")
	w.Report("robbery", "A TILL IS EMPTIED", "b")
	w.CivicDay()
	for _, s := range w.News {
		if s.Kind == "civic" {
			t.Fatalf("a day with a killing in it was padded with %q", s.Headline)
		}
	}
	// But a day with only one thing in it gets the rest of its page filled.
	w.News = nil
	w.Report("robbery", "A TILL IS EMPTIED", "b")
	w.CivicDay()
	if len(w.News) != 2 {
		t.Fatalf("a thin day went to press with %d stories", len(w.News))
	}
}

func TestTheOrdinaryEditionIsTrue(t *testing.T) {
	t.Parallel()
	// Everything in it is read out of state the city already holds. If the
	// price is $40 the paper says $40.
	w := proprietor(t)
	w.News = nil
	if len(w.Goods) == 0 {
		t.Skip("no market to report on")
	}
	items := w.civicItems()
	if len(items) < 3 {
		t.Fatalf("the ordinary edition could only find %d things to say", len(items))
	}
	g := w.Goods[0]
	found := false
	for _, item := range items {
		if item[0] == "PRICE OF "+upper(g.Name)+" STEADY" || item[0] == "PRICE OF "+upper(g.Name)+" RISES" || item[0] == "PRICE OF "+upper(g.Name)+" FALLS" {
			found = true
			if !strings.Contains(item[1], fmt.Sprintf("%d", g.Price)) {
				t.Fatalf("the paper reports a price the market does not hold: %q against $%d", item[1], g.Price)
			}
		}
	}
	if !found {
		t.Fatalf("the market went unreported: %v", items[0])
	}
}
