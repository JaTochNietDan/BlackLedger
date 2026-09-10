package core

import (
	"strings"
	"testing"
)

// The last of the three risks layer 6 asks for. Seizure was there; a rival who
// wants the route went in last slice; an informant is somebody in this city
// deciding that what you are carrying is worth telling the police about.
//
// It has to be a person. "The people who rob the player are named people with a
// place in the city, not anonymous thieves, and they can be answered" — the
// same holds here, or it is a dice roll wearing a hat.

func marked2(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(89)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 3000
	w.Player.Runs = RouteNotice * 2
	w.Player.Contacts = 4 // reach enough to hear a name back
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) || n.Faction == w.PlayerOrganizationID() {
			continue
		}
		n.Trust = 25
		w.Aggrieve(n.ID, GrudgeActs, "the money you took off them")
		return w, n
	}
	t.Fatal("nobody in this city")
	return nil, nil
}

func TestSomebodyWithAGrudgeTalksToThePolice(t *testing.T) {
	w, n := marked2(t)
	heat := w.Player.Heat
	spoke := false
	for i := 0; i < 400 && !spoke; i++ {
		w.WorldRNG = uint32(i*2654435761 + 11)
		spoke = w.ConsiderInformant()
	}
	if !spoke {
		t.Fatal("a man who wants you finished never once picked up a telephone")
	}
	if w.Player.Heat <= heat {
		t.Fatalf("somebody talked and the police took no more interest: %d then %d", heat, w.Player.Heat)
	}
	// And it is spent. Telling them once is the thing they had to say.
	if w.NPC(n.ID).Sore >= GrudgeActs {
		t.Fatalf("they said their piece and are as sore as ever: %d", w.NPC(n.ID).Sore)
	}
}

func TestTheNameComesBackIfYouCanReach(t *testing.T) {
	w, n := marked2(t)
	for i := 0; i < 400; i++ {
		w.WorldRNG = uint32(i*2654435761 + 11)
		if w.ConsiderInformant() {
			break
		}
	}
	said := ""
	for _, r := range w.History {
		if strings.Contains(r.Title, "Somebody talked") {
			said = r.Text
		}
	}
	if said == "" {
		t.Fatal("nothing about it reached the record")
	}
	if !strings.Contains(said, n.Name) {
		t.Fatalf("four contacts and the name never came back: %q", said)
	}
	// Without the reach, it is a fact about the city and not about a person.
	w2, _ := marked2(t)
	w2.Player.Contacts = 0
	for i := 0; i < 400; i++ {
		w2.WorldRNG = uint32(i*2654435761 + 11)
		if w2.ConsiderInformant() {
			break
		}
	}
	for _, r := range w2.History {
		if strings.Contains(r.Title, "Somebody talked") && strings.Contains(r.Text, "Contacts") {
			t.Fatal("a name came back to somebody with nobody to ask")
		}
	}
}

func TestNobodyWithNothingAgainstYouTalks(t *testing.T) {
	w := New(89)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Runs = 100, RouteNotice*2
	for i := range w.NPCs {
		w.NPCs[i].Sore, w.NPCs[i].SoreAt = 0, ""
	}
	for i := 0; i < 200; i++ {
		w.WorldRNG = uint32(i*2654435761 + 3)
		if w.ConsiderInformant() {
			t.Fatal("somebody with no reason to said something anyway")
		}
	}
}

// Taking somebody's wallet in the street had gone through answerFor, which
// reaches their family, and never through Aggrieve, which reaches them. So a
// policy that mugged five hundred and seventy people made not one enemy who
// would say anything about it, and the informant path could not be exercised at
// all.
func TestAMuggedManHoldsItAgainstYou(t *testing.T) {
	w := New(97)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health, w.Player.Respect = 500, 100, 30
	mark := ""
	for _, l := range Locations {
		w.Player.Location = l.ID
		if n, ok := w.MuggingTarget(l.ID); ok && w.MuggingReadiness(l.ID) == "" {
			mark = n.ID
			break
		}
	}
	if mark == "" {
		t.Skip("nobody in this city worth walking up to")
	}
	before := w.NPC(mark).Sore
	// However it goes, they know it happened to them.
	if err := w.Mug(w.Player.Location, w.OwnHands()); err != nil {
		t.Fatal(err)
	}
	if w.NPC(mark).Sore <= before {
		t.Fatalf("somebody was robbed in the street and holds %d against anybody", w.NPC(mark).Sore)
	}
}
