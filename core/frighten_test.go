package core

import "testing"

// Moving against a business meant damaging the property: go in with your crew,
// take the condition off it, and wait for the answer. That is the loud option
// and it was the only one. A business is people now, and people can be put off
// coming in — which costs the holder the same trade for a couple of days
// without a broken window to point at.

func rivalShop(t *testing.T) (*World, string, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Respect = 5000, 100, 30
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	w.Properties["butcher"].Owner = rival
	w.EmptyChairs()
	w.Player.Location = "butcher"
	return w, "butcher", rival
}

func TestYouCanPutTheWindUpARivalsCounter(t *testing.T) {
	t.Parallel()
	w, id, rival := rivalShop(t)
	prop := w.Properties[id]
	filled := prop.Staff
	if filled == 0 {
		t.Fatal("a rival's butcher has nobody in it")
	}
	goodwill := w.faction(rival).Goodwill
	if reason := w.FrightenReadiness(id); reason != "" {
		t.Fatalf("putting the wind up a rival's people was refused: %s", reason)
	}
	if err := w.Frighten(id); err != nil {
		t.Fatalf("it failed: %v", err)
	}
	if prop.Staff >= filled {
		t.Fatalf("everybody came in as usual: %d of %d", prop.Staff, filled)
	}
	if len(prop.Hands) != prop.Staff {
		t.Fatalf("the butcher says %d filled and has %d names", prop.Staff, len(prop.Hands))
	}
	if prop.Shorthanded <= w.Minute {
		t.Fatal("somebody stopped turning up and the counter is ready to hire the same morning")
	}
	if w.faction(rival).Goodwill >= goodwill {
		t.Fatalf("%s lost a pair of hands to you and thinks as much of you as before", w.faction(rival).Name)
	}
}

// It is not free. Somebody you frighten remembers your face, and it is not
// something you can do to your own people or to a shop nobody holds.
func TestFrighteningPeopleIsRemembered(t *testing.T) {
	t.Parallel()
	w, id, _ := rivalShop(t)
	who := w.Properties[id].Hands[0]
	if err := w.Frighten(id); err != nil {
		t.Fatal(err)
	}
	sore := 0
	for i := range w.NPCs {
		if w.NPCs[i].Sore > 0 {
			sore++
		}
	}
	if sore == 0 {
		t.Fatalf("nobody at %s minded, including %s", id, w.NPC(who).Name)
	}
}

func TestYouCannotFrightenYourOwnPeopleOrNobodys(t *testing.T) {
	t.Parallel()
	w, _, _ := rivalShop(t)
	own(w, "laundry")
	w.EmptyChairs()
	w.Player.Location = "laundry"
	if w.FrightenReadiness("laundry") == "" {
		t.Fatal("a player can frighten their own staff off their own counter")
	}
	// And somewhere nobody holds has nobody to lean on.
	w.Properties["restaurant"].Owner = "independent"
	w.Player.Location = "restaurant"
	if w.FrightenReadiness("restaurant") == "" {
		t.Fatal("a shop nobody holds can be leaned on for somebody's benefit")
	}
}

// And it is offered where the shop is, not from across the city.
func TestItIsOfferedWhereTheShopIs(t *testing.T) {
	t.Parallel()
	w, id, _ := rivalShop(t)
	var lean *Action
	for _, a := range w.Actions(id) {
		if a.ID == "frighten" {
			lean = &a
		}
	}
	if lean == nil {
		t.Fatal("standing in a rival's shop there is no way to put the wind up anybody")
	}
	if lean.Disabled {
		t.Fatalf("it is refused where it should work: %s", lean.Reason)
	}
	if err := w.apply(Command{Kind: "frighten", Target: id, RequestID: "windupthecounter1"}); err != nil {
		t.Fatalf("through the same path as everything else it failed: %v", err)
	}
}
