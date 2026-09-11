package core

import "testing"

// Walking into a room that is not yours and saying that it is now.
//
// A family that holds ground on your streets sends for a share of what your
// places take. The player could do none of that back: everything aimed at a
// family was violence, or politics at a distance. This is the mirror — the same
// share of the same takings, paid or refused for the same reason, which is
// whether the people in the room think you can make it stick.

func leaning(t *testing.T, presence int) (*World, string, *Faction) {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 5000
	w.Player.Respect = presence
	id := ""
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop != nil && prop.Income > 0 && w.faction(prop.Owner) != nil {
			id = l.ID
			break
		}
	}
	if id == "" {
		t.Skip("no family holds anything that earns")
	}
	w.Player.Location = id
	return w, id, w.faction(w.Properties[id].Owner)
}

func TestAFamilyPaysWhenYourNameIsWorthMoreThanTheirs(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 200)
	f.Power = 10
	f.Cash = 50000
	if reason := w.DemandReadiness(id); reason != "" {
		t.Fatalf("refused: %s", reason)
	}
	if !w.WouldPay(id) {
		t.Fatalf("presence %d against power %d and they would not pay", w.Presence(), f.Power)
	}
	cash, theirs, standing, respect := w.Player.Cash, f.Cash, f.Goodwill, w.Player.Respect
	share := w.TheirShare(id)
	if err := w.DemandAShare(id); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash+share {
		t.Fatalf("took $%d where a share is $%d", w.Player.Cash-cash, share)
	}
	if f.Cash != theirs-share {
		t.Fatalf("their till went from $%d to $%d", theirs, f.Cash)
	}
	if f.Goodwill >= standing {
		t.Fatal("asking cost them nothing of their opinion of you")
	}
	if w.Player.Respect <= respect {
		t.Fatal("being paid in front of the room was worth no standing")
	}
}

func TestAStrongerFamilyTellsYouToLeave(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 40)
	f.Power = 100
	f.Cash = 50000
	if w.WouldPay(id) {
		t.Fatalf("presence %d against power %d and they opened the till", w.Presence(), f.Power)
	}
	cash, theirs, standing := w.Player.Cash, f.Cash, f.Goodwill
	if err := w.DemandAShare(id); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash {
		t.Fatalf("they refused and $%d changed hands", w.Player.Cash-cash)
	}
	if f.Cash != theirs {
		t.Fatal("their till moved on a refusal")
	}
	// Asking is the insult. Being refused costs the same standing as being paid.
	if f.Goodwill != standing-DemandGrudge {
		t.Fatalf("their opinion went from %+d to %+d where asking costs %d", standing, f.Goodwill, DemandGrudge)
	}
}

func TestNobodyHereHasHeardOfYou(t *testing.T) {
	t.Parallel()
	w, id, _ := leaning(t, 0)
	w.Player.Respect = 0
	if w.Presence() >= DemandStanding {
		t.Skipf("this player starts at a presence of %d", w.Presence())
	}
	if reason := w.DemandReadiness(id); reason == "" {
		t.Fatal("a nobody walked in and was taken seriously")
	}
	if err := w.DemandAShare(id); err == nil {
		t.Fatal("a nobody was paid anyway")
	}
}

func TestYouCannotLeanOnYourOwnRoomOrAnEmptyOne(t *testing.T) {
	t.Parallel()
	w, id, _ := leaning(t, 200)
	own(w, id)
	// The player's own outfit put into the family list, which is the state a
	// takeover leaves behind: the deed is yours and the owner is a family id.
	// Without it, "is this somebody's" is enough on its own to keep you out of
	// your own till and the rule that checks the deed cannot be shown to
	// matter — the first version of this test proved nothing for that reason.
	mine := w.PlayerOrganizationID()
	w.Factions = append(w.Factions, Faction{ID: mine, Name: "Your outfit", Power: 20, Peak: 20})
	w.Properties[id].Owner = mine
	if w.faction(mine) == nil || !w.Own(id) {
		t.Fatal("the deed did not end up with the player's own outfit")
	}
	if Leaning(w, id) {
		t.Fatal("a room of your own is somebody else's to lean on")
	}
	if w.DemandReadiness(id) == "" {
		t.Fatal("you asked yourself for a share")
	}
	// And somewhere nobody holds.
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop != nil && w.faction(prop.Owner) == nil && !w.Own(l.ID) {
			w.Player.Location = l.ID
			if Leaning(w, l.ID) {
				t.Fatalf("%s answers to a family", l.ID)
			}
			if w.DemandReadiness(l.ID) == "" {
				t.Fatalf("asked %s for a share and there is nobody to ask", l.ID)
			}
			return
		}
	}
}
