package core

import (
	"strings"
	"testing"
)

// "Death is permanent. What you built passes to the strongest of your own
// people and becomes an organization you can deal with, or fight."
//
// The passing was verified once. The second half was not: an organization the
// next protagonist can DEAL WITH, or FIGHT. That means it has to be a real
// family in the city — one the new player can enquire about, reach an
// understanding with, go to work for, move against, and read about in the
// paper — and not a bookkeeping entry with the old holdings attached.

func builtSomething(t *testing.T, seed uint32) *World {
	t.Helper()
	w := New(seed)
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Health = 6000, 140, 100
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = w.PlayerOrganizationID()
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Fatal("the player never became an organization")
	}
	// People of their own, so there is somebody to inherit.
	for i, n := range w.Civilians() {
		if i >= 3 {
			break
		}
		n.Faction = w.PlayerOrganizationID()
		n.Rank = RankSoldier
		n.Trust = 60
	}
	if len(w.OwnPeople()) == 0 {
		t.Fatal("nobody answers to the player")
	}
	return w
}

func TestWhatYouBuiltBecomesSomethingTheNextLifeCanDealWith(t *testing.T) {
	t.Parallel()
	w := builtSomething(t, 37)
	mine := w.PlayerOrganizationID()
	held := len(w.FamilyHoldings(mine))
	if held == 0 {
		t.Fatal("the organization holds nothing")
	}
	w.Die("Shot on the steps of the Mariner.")
	heir := ""
	for _, f := range w.Factions {
		if strings.HasPrefix(f.ID, "estate:") {
			heir = f.ID
		}
	}
	if heir == "" {
		t.Fatal("nothing was inherited")
	}
	// The next life arrives.
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life", Target: ""})
	if err != nil {
		t.Fatalf("could not begin again: %v", err)
	}
	w = next

	f := w.faction(heir)
	if f == nil {
		t.Fatal("the organization did not survive into the next life")
	}
	if got := len(w.FamilyHoldings(heir)); got != held {
		t.Errorf("it held %d premises and now holds %d", held, got)
	}
	if len(w.Members(heir)) == 0 {
		t.Error("nobody answers to it")
	}
	if w.Leader(heir) == nil {
		t.Error("nobody leads it")
	}

	// "Deal with, or fight." Both have to be reachable, which means the city
	// has to treat it like any other family.
	w.Player.Cash, w.Player.Respect, w.Player.Health = 8000, 120, 100
	for i := range w.Factions {
		w.Factions[i].Goodwill = 40
	}
	kinds := map[string]bool{}
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		for _, a := range w.Actions(id) {
			if strings.HasSuffix(a.ID, ":"+heir) {
				kinds[strings.SplitN(a.ID, ":", 2)[0]] = true
			}
		}
	}
	for _, want := range []string{"enquire", "pact", "serve", "smear"} {
		if !kinds[want] {
			t.Errorf("the city offers no %q about what the last life built", want)
		}
	}
	// And fighting it: its premises can be moved on like anybody else's.
	fought := false
	for _, id := range w.FamilyHoldings(heir) {
		if _, ok := w.SabotageTarget(id); ok {
			fought = true
		}
	}
	if !fought {
		t.Error("none of its premises can be moved on")
	}
	t.Logf("the heir holds %d premises, %d people, and the city offers %v about it",
		len(w.FamilyHoldings(heir)), len(w.Members(heir)), keysOf(kinds))
}

func keysOf(m map[string]bool) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}

// "The city does not scale to you. An organization at ninety strength will kill
// you on your first day if you give it a reason."
//
// Bellandi opens at ninety. Provoking them at the club is the reason. This
// counts what happens to a new arrival who does it, over many cities.
func TestAStrongFamilyWillKillANewArrivalWhoProvokesThem(t *testing.T) {
	t.Parallel()
	died, survived, refused := 0, 0, 0
	for i := 0; i < 120; i++ {
		w := New(uint32(4000 + i*13))
		if f := w.faction("bellandi"); f == nil || f.Power < 90 {
			t.Fatalf("the opening family is not at ninety: %v", f)
		}
		w.Player.Location = "club"
		found := false
		for _, a := range w.Actions("club") {
			if a.ID == "provoke" && !a.Disabled {
				found = true
			}
		}
		if !found {
			refused++
			continue
		}
		next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "provoke", Target: "club"})
		if err != nil {
			refused++
			continue
		}
		w = next
		// Home, because that is where a person is at night and where an attack
		// looks for them. Standing in the nightclub they provoked for two days
		// running means the car pulls up outside an empty room — which is what
		// a first version of this test measured.
		w.Player.Location = w.Player.Home
		// The rest of the first day, and the night after it. When they come,
		// meet it the way somebody with nothing meets it — by standing in the
		// doorway. Clearing the scene instead of answering it discards the
		// attack entirely, which is how a first version of this test found
		// nobody ever died.
		for j := 0; j < 96 && w.Player.Alive; j++ {
			w.Advance(30)
			if w.Event == nil {
				continue
			}
			answer := ""
			for _, c := range w.Event.Choices {
				if c.ID == "defend" || c.ID == "acknowledge" {
					answer = c.ID
				}
			}
			if answer == "" && len(w.Event.Choices) > 0 {
				answer = w.Event.Choices[0].ID
			}
			after, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
				Kind: "choice", Event: w.Event.ID, Choice: answer})
			if err != nil {
				w.Event = nil
				continue
			}
			w = after
		}
		if w.Player.Alive {
			survived++
		} else {
			died++
		}
	}
	if refused > 0 {
		t.Logf("%d cities would not offer the provocation at all", refused)
	}
	t.Logf("of %d new arrivals who provoked a family at ninety: %d died, %d were still standing two days later",
		died+survived, died, survived)
	// "Will kill you" is not "might". A majority has to die, or the sentence is
	// a warning the city does not mean.
	if died*2 <= died+survived {
		t.Errorf("the rules say a family at ninety will kill you on your first day; %d of %d survived",
			survived, died+survived)
	}
	// And the other half: it has to be survivable, or the sentence is not a
	// warning but a death notice, and the choice it describes is not a choice.
	if survived == 0 {
		t.Error("nobody survived it at all, which makes the provocation a way of ending the game rather than a risk")
	}
}
