package core

import (
	"sort"
	"strings"
	"testing"
)

// Bail could never be pressed. One rule disabled it with the very reason it was
// offered, and no test asked the obvious question of it: is there any state at
// all in which this button is live?
//
// This asks that of every button the game can produce. It drives many cities in
// many conditions, records every action id ever offered and every id ever
// offered ENABLED, and reports anything that appears only in the first list. A
// button that is offered but never live is either dead like bail was, or needs
// a state this sweep does not reach — and either way it is worth knowing which.

// kindOf strips the subject from an id so "lend:mara" and "lend:leo" count as
// one button rather than forty.
func kindOf(id string) string {
	at := strings.IndexByte(id, ':')
	if at < 0 {
		return id
	}
	head, tail := id[:at], id[at+1:]
	// Some ids carry a real variant after the colon (buy:arms, operate:hard,
	// trip:kingsport) and some carry a person or a place. Keep the variant.
	switch head {
	case "buy", "sell", "operate", "trip", "fit", "arms", "retain", "play", "pact", "enquire", "serve", "smear":
		// These carry a real variant, except when the variant is a family the
		// living world invented — one of those is the same button as the one
		// aimed at a family the campaign started with.
		if strings.HasPrefix(tail, "splinter-") || strings.HasPrefix(tail, "estate:") || strings.HasPrefix(tail, "player:") {
			return head + ":<a family>"
		}
		return id
	}
	if len(tail) > 0 && (tail == "crew" || tail == "weapon" || tail == "armour") {
		return id
	}
	return head + ":*"
}

func sweepCities(t *testing.T) (offered, enabled map[string]int) {
	t.Helper()
	offered, enabled = map[string]int{}, map[string]int{}
	note := func(w *World) {
		for i := range Locations {
			id := Locations[i].ID
			w.Player.Location = id
			for _, a := range w.Actions(id) {
				k := kindOf(a.ID)
				offered[k]++
				if !a.Disabled {
					enabled[k]++
				}
			}
		}
	}
	for _, seed := range []uint32{31, 47, 88, 103, 219} {
		// A rich, established player with a crew, an organization, goods, a
		// car, an account, people on the books and somebody in a cell.
		w := New(seed)
		w.District = 2
		w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
		w.Player.Contacts = 3
		w.Offshore = 4000
		w.Player.Offshore = true
		w.Player.Car, w.Player.CarWear = 1, 90
		w.Player.Dress, w.Player.DressWear = 1, 90
		for _, id := range []string{"laundry", "garage", "casino"} {
			w.Properties[id].Owner = "player:1"
		}
		w.Incorporate()
		w.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
		w.NPCs = append(w.NPCs, NPC{ID: "inside", Name: "Otto Reiss", Location: "precinct",
			Faction: w.PlayerOrganizationID(), Rank: RankSoldier, Held: w.Minute + 2*1440})
		note(w)
		// The same city, aged, so late-game buttons appear.
		live(t, w, 2000)
		w.Player.Cash, w.Player.Respect, w.Player.Health = 20000, 200, 100
		note(w)

		// A beginner with enough standing to open the next district.
		b := New(seed)
		b.Player.Cash, b.Player.Respect = 900, 20
		note(b)

		// Somebody serving a family, which unlocks a different set, and who is
		// not yet an organization of their own.
		s := New(seed)
		s.District = 2
		s.Player.Cash, s.Player.Respect, s.Player.Health = 8000, 60, 100
		s.Player.Serves, s.Player.Service = "bellandi", 10
		for i := range s.Factions {
			s.Factions[i].Goodwill = 60
		}
		s.Properties["laundry"].Owner = "player:1"
		note(s)

		// Nobody's man, but well thought of everywhere, which is what it takes
		// to be taken on or to reach an understanding.
		f := New(seed)
		f.District = 2
		f.Player.Cash, f.Player.Respect, f.Player.Health = 8000, 60, 100
		for i := range f.Factions {
			f.Factions[i].Goodwill = 80
		}
		note(f)

		// A proprietor with the press and the police on a retainer, short-handed
		// premises, a clean pair of hands and nothing offshore yet.
		r := New(seed)
		r.District = 2
		r.Player.Cash, r.Player.Respect, r.Player.Health = 20000, 200, 100
		r.Player.Heat = 0
		for _, id := range []string{"laundry", "garage", "casino"} {
			r.Properties[id].Owner = "player:1"
			r.Properties[id].Staff = 0
			r.Properties[id].Supply = 40
			r.Properties[id].Mode = "hard"
		}
		r.Incorporate()
		// Retainers are held by the ROLE, not by the person in it — an
		// arrangement is with the editor's chair, and it survives the editor.
		r.ensureOfficials()
		r.Player.Retainers = append(r.Player.Retainers, "editor", "commissioner", "mayor")
		r.News = append(r.News, Story{ID: ID(), Minute: r.Minute, Life: r.Life,
			Headline: "QUESTIONED AT THE DOCKS", Body: "Police called again.", Kind: "police"})
		note(r)

		// The same proprietor with the police interested in him, premises busy
		// enough to be offered a standing order, money sitting abroad he has
		// not yet reached, and families who think well of him.
		r.Player.Heat = 60
		for _, id := range []string{"laundry", "garage", "casino"} {
			r.Properties[id].Supply = 200
			r.Properties[id].Staff = 9
			r.Properties[id].Condition = 100
		}
		r.Offshore = 5000
		r.Player.Offshore = false
		for i := range r.Factions {
			r.Factions[i].Goodwill = 60
		}
		live(t, r, 200)
		r.Player.Cash, r.Player.Health = 20000, 100
		// A hundred hours of a hard-run house leaves the police very interested
		// indeed: heat comes out at 77 to 95 across these seeds, and a detective
		// will not be seen taking money from anybody past BribeCeiling. This
		// sweep was reaching the bribe by luck — one seed happening to land just
		// under the line — so a tick that changed nothing about the police could
		// take the button dark. The state the comment above describes is being
		// built rather than hoped for.
		r.Player.Heat = BribeCeiling - 10
		note(r)
	}
	return offered, enabled
}

// mustBeLive is the regression this test exists for. Bail was disabled in every
// campaign by a rule with no exception for it, and nothing failed. These are
// the buttons that were hardest to reach — several were on a list of actions I
// had never seen offered at all — and every one is now confirmed live by a
// state built below. If one goes dark again, this fails.
var mustBeLive = []string{
	"bail:*",          // the person is in a cell, which is why it is offered
	"launder",         // needs the police interested in you
	"spike", "puff",   // need an arrangement with the paper
	"smear:bellandi",  // the same, aimed at a family
	"bribe",           // needs the police and a quiet profile at once
	"charge",          // needs the docks and the money
	"offshore_access", // needs money abroad and none of it reachable yet
	"leave_service", "move", "remedy", "trip:rockridge",
	"hire", "expand", "operate:standard", "restock", "still",
}

func TestEveryButtonTheGameOffersCanBePressedSomewhere(t *testing.T) {
	offered, enabled := sweepCities(t)
	dead := []string{}
	for k := range offered {
		if enabled[k] == 0 {
			dead = append(dead, k)
		}
	}
	sort.Strings(dead)
	t.Logf("%d kinds of button offered, %d of them live somewhere", len(offered), len(offered)-len(dead))
	for _, k := range dead {
		t.Logf("NEVER LIVE: %s (offered %d times)", k, offered[k])
	}
	if len(dead) > 0 {
		t.Logf("a button here needs a state this sweep does not build; each was read and explains itself")
	}

	// The regression this exists for. Bail was disabled in every campaign by a
	// rule that had no exception for it, and nothing failed. These are buttons
	// whose whole point is a state that also refuses them, or that took a
	// deliberately built fixture to reach; if one goes dark again, this fails.
	for _, k := range mustBeLive {
		if enabled[k] == 0 {
			t.Errorf("%s is offered %d times and never live: no state in this sweep can press it", k, offered[k])
		}
	}
	live := make([]string, 0, len(enabled))
	for k, n := range enabled {
		if n > 0 {
			live = append(live, k)
		}
	}
	sort.Strings(live)
	t.Logf("live: %v", live)
	if len(offered) < 90 {
		t.Errorf("only %d kinds of button were offered at all; the sweep has stopped reaching the game", len(offered))
	}
}
