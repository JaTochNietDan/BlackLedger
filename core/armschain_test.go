package core

import (
	"strings"
	"testing"
)

// The whole arms trade, walked end to end with the game's own actions.
//
// `stock_arms` — putting the crates under the floor — was on the list of
// actions the harness has never played, and it is the middle of a chain rather
// than a thing on its own: own premises of the right kind, build a room under
// the floor, buy crates where crates are sold, carry them there, put them down,
// and then a family at war comes to the door with money. Every link existed and
// nothing had ever joined them.
func TestTheArmsChainFromTheWaterfrontToSomebodyAtWar(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event, w.District, w.Player.Cash = nil, 9, 900000

	// Premises of a kind that has a floor worth lifting.
	site := ""
	for _, l := range Locations {
		if ArmourySite(l.ID) {
			site = l.ID
			break
		}
	}
	if site == "" {
		t.Skip("this city has nowhere for a room like that")
	}
	w.Properties[site].Owner = "player:1"
	w.Player.Location = site

	if reason := w.ArmouryReadiness(site); reason != "" {
		t.Fatalf("the room cannot be built: %s", reason)
	}
	built, err := Execute(w, Command{Kind: "armoury", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = built
	if _, ok := w.TheArmoury(); !ok {
		t.Fatal("the room was paid for and is not there")
	}

	// Crates, bought where crates are bought.
	market := ""
	for _, l := range Locations {
		if TradesAt(l.ID, "arms") {
			market = l.ID
			break
		}
	}
	if market == "" {
		t.Skip("nowhere in this city sells arms")
	}
	w.Player.Location, w.Event = market, nil
	if reason := w.TradeReadiness("arms", "buy", 4); reason != "" {
		t.Fatalf("four crates cannot be bought at the waterfront: %s", reason)
	}
	bought, err := Execute(w, Command{Kind: "buy:arms", Amount: 4, Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = bought
	if w.Holding("arms") != 4 {
		t.Fatalf("carrying %d crates after buying four", w.Holding("arms"))
	}

	// Carried back and put down. Read off the room's own card.
	w.Player.Location, w.Event = site, nil
	var card Action
	for _, a := range w.Actions(site) {
		if a.ID == "stock_arms" {
			card = a
		}
	}
	if card.ID == "" {
		t.Fatal("the room does not offer to take the crates")
	}
	if card.Disabled {
		t.Fatalf("the crates cannot be put down: %s", card.Reason)
	}
	stocked, err := Execute(w, Command{Kind: "stock_arms", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = stocked
	if w.Stocked() != 4 {
		t.Fatalf("%d crates are under the floor after putting four down", w.Stocked())
	}
	if w.Holding("arms") != 0 {
		t.Fatalf("%d crates are still being carried", w.Holding("arms"))
	}

	// And the war comes to the door. Somebody has to be fighting, which is the
	// city's business rather than the player's, so it is arranged here.
	if len(w.Factions) < 2 {
		t.Skip("this city has nobody to fight")
	}
	a, b := &w.Factions[0], &w.Factions[1]
	if c := w.Conflict(a.ID, b.ID); c != nil {
		c.Hostility, c.State = 100, "war"
	}
	a.Cash, b.Cash = 900000, 900000
	if len(w.buyers()) == 0 {
		t.Fatal("two organizations are at war and neither is buying")
	}

	cash, crates := w.Player.Cash, w.Stocked()
	w.ArmouryDay()
	if w.Stocked() >= crates {
		t.Fatalf("%d crates under the floor and a war on, and %d are still there", crates, w.Stocked())
	}
	if w.Player.Cash <= cash {
		t.Fatalf("crates left the room and the player is no better off: $%d then, $%d now", cash, w.Player.Cash)
	}
}

// And what a search does to the room, which is the risk the whole chain is
// priced against.
//
// The room is always found — that is the same rule the still is under, and it
// is what the card is now allowed to say. What was wrong was the reporting: an
// empty room came back as found with nought crates in it, so the log said
// "0 crates of arms out through the front door" and the paper carried it as the
// largest find of its kind this year.
func TestAnEmptyRoomIsNotTheYearsBiggestFind(t *testing.T) {
	t.Parallel()
	site, w := "", New(23)
	w.Event, w.District, w.Player.Cash = nil, 9, 900000
	for _, l := range Locations {
		if ArmourySite(l.ID) {
			site = l.ID
			break
		}
	}
	if site == "" {
		t.Skip("this city has nowhere for a room like that")
	}
	w.Properties[site].Owner = "player:1"
	w.Properties[site].Income = 40
	w.Player.Location = site
	if err := w.BuildArmoury(site); err != nil {
		t.Fatal(err)
	}
	if w.Stocked() != 0 {
		t.Fatal("a room just built already has something in it")
	}

	before := w.Properties[site].Condition
	w.search()

	if w.Properties[site].Armoury {
		t.Fatal("the floor came up and the room is still there")
	}
	// At least what pulling the floor up costs. A search does other things to a
	// building on the way through, so this is a floor and not an equality.
	if after := w.Properties[site].Condition; after > before-ArmouryRuin {
		t.Fatalf("the building went from %d%% to %d%%, and pulling a floor up alone costs %d",
			before, after, ArmouryRuin)
	}
	told := false
	for _, entry := range w.History {
		if !strings.HasPrefix(entry.Title, "They found the room at") {
			continue
		}
		told = true
		if strings.Contains(entry.Text, " 0 crates") {
			t.Fatalf("the search reports nothing as a number of crates: %q", entry.Text)
		}
	}
	if !told {
		t.Fatal("the floor came up and the player was never told")
	}
	for _, s := range w.News {
		if strings.Contains(s.Headline, "ARMS CACHE SEIZED") {
			t.Fatalf("the paper carries a cache that was not there: %q", s.Headline)
		}
	}
}
