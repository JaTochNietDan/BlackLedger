package core

import "testing"

// "They could have knowledge of where you live but not where you currently are,
// and when you move house they will no longer know where you live until they
// find out via some contact."

func hunted(t *testing.T) (*World, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 5000, 100
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	// Nobody of theirs anywhere near the player.
	for i := range w.NPCs {
		if w.NPCs[i].Faction == rival {
			w.NPCs[i].Location = "estate"
		}
	}
	w.Player.Location = "bar"
	return w, rival
}

func TestTheyKnowYourAddressAndNotYourEvening(t *testing.T) {
	t.Parallel()
	w, rival := hunted(t)
	if w.TheyKnowWhereYouAre(rival) {
		t.Fatal("a family with nobody near the player knows exactly where they are drinking")
	}
	// At home they always can find you: it is an address, and addresses do not
	// move.
	w.Player.Location = w.Player.Home
	if !w.TheyKnowWhereYouAre(rival) {
		t.Fatal("a family cannot find the player at the address with their name on it")
	}
}

func TestBeingSeenIsWhatCostsYouTheCover(t *testing.T) {
	t.Parallel()
	w, rival := hunted(t)
	// One of theirs walks into the bar.
	for i := range w.NPCs {
		if w.NPCs[i].Faction == rival {
			w.NPCs[i].Location = "bar"
			break
		}
	}
	w.noticedByTheCity()
	if !w.TheyKnowWhereYouAre(rival) {
		t.Fatal("one of theirs stood in the same room and they still cannot find the player")
	}
	// And it goes stale.
	w.Minute += TheyKnowYou + 60
	if w.TheyKnowWhereYouAre(rival) {
		t.Fatal("a day later they still know which bar the player was in")
	}
}

// The whole point of it: an attempt on somebody they cannot place goes to the
// address they do have.
func TestAnAttemptOnSomebodyTheyCannotFindGoesToTheHouse(t *testing.T) {
	t.Parallel()
	w, rival := hunted(t)
	w.Properties[w.Player.Home].Condition = 100
	health := w.Player.Health
	w.Attack(Plot{ID: ID(), Kind: "hit", Life: w.Life, Actor: rival, Strength: 60})
	if w.Player.Health != health {
		t.Fatalf("they could not find the player and hurt them anyway: %d against %d", w.Player.Health, health)
	}
	if w.Properties[w.Player.Home].Condition >= 100 {
		t.Fatal("they could not find the player and left the house alone as well")
	}
}

// Moving house takes back what they knew.
func TestMovingHouseTakesBackWhatTheyKnew(t *testing.T) {
	t.Parallel()
	w, rival := hunted(t)
	for i := range w.NPCs {
		if w.NPCs[i].Faction == rival {
			w.NPCs[i].Location = "bar"
			break
		}
	}
	w.noticedByTheCity()
	if !w.TheyKnowWhereYouAre(rival) {
		t.Fatal("they saw the player and know nothing")
	}
	w.MovedHouse()
	if w.TheyKnowWhereYouAre(rival) {
		t.Fatal("the player moved and everybody followed them")
	}
}
