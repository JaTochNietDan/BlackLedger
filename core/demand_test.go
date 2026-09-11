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

// And when asking is no longer the point.
//
// A family can come for your deed once you have refused them long enough. The
// player could not do the same: a takeover reaches only a family you serve, and
// everything else aimed at a rival was a war. A weak family that already hated
// you could hold a shop on your own street for ever.

func TestAFamilyWalksOutWhenItHasCountedWhatYouAreWorth(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 200)
	f.Goodwill, f.Power = DemandFloor, 20
	// Somewhere else to go, so this is not the last thing they have.
	for _, l := range Locations {
		if l.ID != id && w.Properties[l.ID] != nil && w.Properties[l.ID].Income > 0 && !w.Own(l.ID) {
			w.Properties[l.ID].Owner = f.ID
			break
		}
	}
	if len(w.FamilyHoldings(f.ID)) <= PushLeft {
		t.Skipf("%s holds only %d places", f.Name, len(w.FamilyHoldings(f.ID)))
	}
	if reason := w.PushReadiness(id); reason != "" {
		t.Fatalf("refused: %s", reason)
	}
	power, respect, heat := f.Power, w.Player.Respect, w.Player.Heat
	if err := w.TakeItFromThem(id); err != nil {
		t.Fatal(err)
	}
	if !w.Own(id) {
		t.Fatal("they stayed")
	}
	if f.Power >= power {
		t.Fatalf("their power went from %d to %d", power, f.Power)
	}
	if w.Player.Respect <= respect {
		t.Fatal("taking a room off a family was worth no standing")
	}
	if w.Player.Heat <= heat {
		t.Fatal("Ward Street did not notice")
	}
	// And everybody else in the city did.
	for i := range w.Factions {
		if w.Factions[i].ID == f.ID {
			continue
		}
		if w.Factions[i].Goodwill >= 0 {
			t.Errorf("%s watched it happen and thinks the same of you", w.Factions[i].Name)
		}
	}
}

func TestTheyFightForItUntilTheyHaveStoppedPretending(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 200)
	f.Power = 20
	// On any sort of terms with you, and they would rather fight.
	f.Goodwill = DemandFloor + 1
	if w.PushReadiness(id) == "" {
		t.Fatalf("they walked out at a standing of %+d where it takes %+d", f.Goodwill, DemandFloor)
	}
	// And a family worth more in this city than you are does not walk out of
	// anywhere, however it feels about you. It is the margin that decides it
	// rather than either figure on its own: a family at the top of its power
	// still leaves a room to somebody whose name is worth far more, which is
	// the rule and not what the first version of this test assumed.
	f.Goodwill = -100
	f.Power = w.Presence() - PushMargin + 1
	if w.PushReadiness(id) == "" {
		t.Fatalf("a family at power %d walked out for a presence of %d, where it takes %d clear",
			f.Power, w.Presence(), PushMargin)
	}
	f.Power = w.Presence() - PushMargin
	if w.PushReadiness(id) != "" {
		t.Fatalf("exactly %d clear and they stayed: %s", PushMargin, w.PushReadiness(id))
	}
}

func TestNobodyTakesTheLastThingAFamilyHas(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 200)
	f.Goodwill, f.Power = -100, 10
	// Everything of theirs but this one.
	for _, l := range Locations {
		if l.ID != id && w.Properties[l.ID] != nil && w.Properties[l.ID].Owner == f.ID {
			w.Properties[l.ID].Owner = "independent"
		}
	}
	if len(w.FamilyHoldings(f.ID)) != 1 {
		t.Skipf("%s holds %d places", f.Name, len(w.FamilyHoldings(f.ID)))
	}
	if w.PushReadiness(id) == "" {
		t.Fatal("a family was put out of this city in an afternoon")
	}
	if err := w.TakeItFromThem(id); err == nil {
		t.Fatal("and it went through anyway")
	}
}

// The whole way down, in one test.
//
// Asking for a share and taking the room are two halves of one idea and were
// built a tick apart, which is how this project has repeatedly ended up with a
// thing that works in pieces and is unreachable as a path. The take needs a
// family at −60; refusing their demands reaches about −40 in a campaign, so if
// leaning on them did not also drive them down the whole feature would sit
// behind a number nothing produces.
//
// It does. Four demands take a family past the floor and six put them at the
// bottom of it, and then they walk out of the room.
func TestLeaningOnAFamilyIsThePathToTakingItsRoom(t *testing.T) {
	t.Parallel()
	w, id, f := leaning(t, 200)
	f.Power = 20
	// Somewhere else for them to go, so this is not the last thing they hold.
	for _, l := range Locations {
		if l.ID != id && w.Properties[l.ID] != nil && w.Properties[l.ID].Income > 0 && !w.Own(l.ID) {
			w.Properties[l.ID].Owner = f.ID
			break
		}
	}
	if len(w.FamilyHoldings(f.ID)) <= PushLeft {
		t.Skipf("%s holds only %d places", f.Name, len(w.FamilyHoldings(f.ID)))
	}
	if w.PushReadiness(id) == "" {
		t.Fatal("they would walk out before anybody had asked them for anything")
	}
	asked := 0
	for i := 0; i < 12 && w.PushReadiness(id) != ""; i++ {
		if reason := w.DemandReadiness(id); reason != "" {
			t.Fatalf("after %d demands they will not even be asked: %s", asked, reason)
		}
		if err := w.DemandAShare(id); err != nil {
			t.Fatal(err)
		}
		w.Event = nil
		asked++
	}
	if w.PushReadiness(id) != "" {
		t.Fatalf("twelve demands took them to %+d and they still will not go: %s",
			f.Goodwill, w.PushReadiness(id))
	}
	t.Logf("%d demands took them from nothing to %+d, and then they walked out", asked, f.Goodwill)
	if err := w.TakeItFromThem(id); err != nil {
		t.Fatal(err)
	}
	if !w.Own(id) {
		t.Fatal("they stayed after all")
	}
}
