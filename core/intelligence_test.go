package core

import (
	"strings"
	"testing"
)

func enquirer(t *testing.T) *World {
	t.Helper()
	w := New(263)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Cash, w.Player.Contacts = "market", 5000, 0
	for i := range w.NPCs {
		w.NPCs[i].Trust = 0
	}
	return w
}

func TestAStrangerKnowsNothingAboutAnybody(t *testing.T) {
	w := enquirer(t)
	for _, f := range w.PublicFactions() {
		if w.Intelligence(f.ID) != 0 {
			t.Fatalf("a stranger was at %d on %s", w.Intelligence(f.ID), f.Name)
		}
		if f.Name == "" {
			t.Fatal("a name is common knowledge and was withheld")
		}
		if f.Leader != "" || f.Power != 0 || f.Cash != 0 {
			t.Fatalf("a stranger was told %+v", f)
		}
		if f.Strength != "nobody will say" || f.Money != "nobody will say" {
			t.Fatalf("a stranger was told %q and %q", f.Strength, f.Money)
		}
	}
}

func TestANetworkHearsTheOrdinaryThings(t *testing.T) {
	w := enquirer(t)
	w.Player.Contacts = 2
	f := w.PublicFactions()[0]
	if w.Intelligence(f.ID) != 1 {
		t.Fatalf("two contacts were worth %d", w.Intelligence(f.ID))
	}
	if f.Leader == "" {
		t.Fatal("a network did not know who ran anything")
	}
	if f.Power != 0 || f.Cash != 0 {
		t.Fatal("a network read exact numbers off a screen")
	}
	if f.Strength == "nobody will say" || strings.Contains(f.Strength, "hundred") {
		t.Fatalf("a network was told %q", f.Strength)
	}
}

func TestSomebodyInsideIsWorthMoreThanAnyNetwork(t *testing.T) {
	w := enquirer(t)
	w.Player.Contacts = 2
	id := w.Factions[0].ID
	members := w.Members(id)
	if len(members) == 0 {
		t.Fatal("the fixture's organization had nobody in it")
	}
	members[0].Trust = 1
	if w.Intelligence(id) != 2 {
		t.Fatalf("somebody inside was worth %d", w.Intelligence(id))
	}
	for _, f := range w.PublicFactions() {
		if f.ID != id {
			continue
		}
		if f.Power == 0 || !strings.Contains(f.Strength, "hundred") {
			t.Fatalf("somebody inside did not know their strength: %+v", f)
		}
		if f.Cash != 0 {
			t.Fatal("somebody inside read the books")
		}
	}
}

func TestAskingAroundBuysAWeekOfKnowingMore(t *testing.T) {
	w := enquirer(t)
	w.Player.Contacts = 2
	id := w.Factions[0].ID
	w.Members(id)[0].Trust = 1
	before := w.Intelligence(id)
	cash := w.Player.Cash
	if err := w.AskAround(id); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-EnquiryCost {
		t.Fatalf("asking cost $%d", cash-w.Player.Cash)
	}
	if w.Intelligence(id) != before+1 {
		t.Fatalf("asking was worth %d", w.Intelligence(id)-before)
	}
	for _, f := range w.PublicFactions() {
		if f.ID == id && (f.Cash == 0 || !strings.HasPrefix(f.Money, "$")) {
			t.Fatalf("everything known did not include the books: %+v", f)
		}
	}
	// It is still current, so it cannot be bought again.
	if w.EnquiryReadiness(id) == "" {
		t.Fatal("asked the same question twice in a week")
	}
	// And it lapses.
	w.Minute += EnquiryLasts + 1
	if w.Intelligence(id) != before {
		t.Fatalf("a week later they were still at %d", w.Intelligence(id))
	}
	if w.EnquiryReadiness(id) != "" {
		t.Fatal("could not ask again after it lapsed:", w.EnquiryReadiness(id))
	}
}

func TestNobodyToAskMeansNobodyToAsk(t *testing.T) {
	w := enquirer(t)
	if w.EnquiryReadiness(w.Factions[0].ID) == "" {
		t.Fatal("somebody with no contacts asked around")
	}
	w.Player.Contacts = 1
	if w.EnquiryReadiness(w.Factions[0].ID) != "" {
		t.Fatal("one contact was not enough:", w.EnquiryReadiness(w.Factions[0].ID))
	}
	if w.EnquiryReadiness("nobody") == "" {
		t.Fatal("asked about an organization that does not exist")
	}
}

func TestYouAlwaysKnowWhatTheyThinkOfYouAndYourOwnBooks(t *testing.T) {
	w := enquirer(t)
	w.Factions[0].Goodwill = -40
	for _, f := range w.PublicFactions() {
		if f.ID == w.Factions[0].ID && f.Goodwill != -40 {
			t.Fatal("the player could not tell how they were being treated")
		}
	}
	// And their own organization is not something they have to ask about.
	own(w, "laundry", "garage")
	w.Player.Respect = OrganizationStanding
	w.OrganizationDay()
	if w.Intelligence(w.PlayerOrganizationID()) != 3 {
		t.Fatal("the player had to ask around about themselves")
	}
}

// The Families screen showed a strength, a money word and a bare number for
// standing — "+45", "-63" — and left the player to work out what any of it
// meant. What an organization holds, who it is fighting and what its number
// means are all things the city already knew.

func TestAnOrganizationSaysWhatItHoldsAndWhoItIsFighting(t *testing.T) {
	w, _ := testator(t)
	w.Antagonize("bellandi", "russo", 90)
	for i := range w.Conflicts {
		w.Conflicts[i].State = "war"
	}
	seen := 0
	for _, f := range w.PublicFactions() {
		if f.Standing == "" {
			t.Fatalf("%s has no standing anybody can read", f.Name)
		}
		if len(w.FamilyHoldings(f.ID)) > 0 && len(f.Holdings) == 0 {
			t.Fatalf("%s holds ground and the screen shows none of it", f.Name)
		}
		if f.ID == "bellandi" || f.ID == "russo" {
			if len(f.Fighting) == 0 {
				t.Fatalf("%s is at war and the screen does not say so", f.Name)
			}
			seen++
		}
	}
	if seen != 2 {
		t.Fatalf("%d of the two families at war were reported", seen)
	}

	// The player's own organization says so rather than reporting a number
	// about how much it likes itself.
	own := false
	for _, f := range w.PublicFactions() {
		if f.ID == w.PlayerOrganizationID() {
			own = true
			if !f.Yours || f.Standing != "Yours" {
				t.Fatalf("the player's own organization reads %q", f.Standing)
			}
		}
	}
	if !own {
		t.Fatal("the player's own organization is not on the screen")
	}
}

func TestStandingIsSaidInWordsNotNumbers(t *testing.T) {
	for _, c := range []struct {
		goodwill int
		want     string
	}{{-90, "They have decided about you"}, {-50, "Hostile"}, {-20, "They do not like you"},
		{0, "They have no opinion of you"}, {20, "Cordial"}, {50, "They think well of you"},
		{90, "You are as good as one of theirs"}} {
		if got := standingWith(c.goodwill); got != c.want {
			t.Errorf("%d reads %q, expected %q", c.goodwill, got, c.want)
		}
	}
}

func TestWhatYouCannotCountYouAreNotTold(t *testing.T) {
	// Counting somebody's people needs an informant. Their premises do not.
	w := proprietor(t)
	w.Player.Enquiries = nil
	for _, f := range w.PublicFactions() {
		if f.Knowledge >= 2 {
			continue
		}
		if f.People != 0 {
			t.Fatalf("%s reports %d people without anybody inside", f.Name, f.People)
		}
		if len(w.FamilyHoldings(f.ID)) > 0 && len(f.Holdings) == 0 {
			t.Fatalf("%s holds premises anybody could walk past and they are hidden", f.Name)
		}
	}
}
