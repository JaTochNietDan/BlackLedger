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
