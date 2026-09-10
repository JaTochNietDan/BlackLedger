package core

import "testing"

func mugger(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(101)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Health, w.Player.Cash = "club", 100, 500
	w.Player.Contacts = 3 // knows the city
	mark, ok := w.MuggingTarget("club")
	if !ok {
		t.Fatal("there was nobody at the casino to take anything off")
	}
	return w, mark
}

func TestYouDoNotWalkUpToAStranger(t *testing.T) {
	t.Parallel()
	w := New(101)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Contacts = "club", 0
	// Only leaders are names everybody knows, so move them out of the way.
	for i := range w.NPCs {
		if w.NPCs[i].Rank >= RankLeader {
			w.NPCs[i].Location = "estate"
		}
	}
	if _, ok := w.MuggingTarget("club"); ok {
		t.Fatal("a stranger with no contacts picked somebody out of a room")
	}
	if w.MuggingReadiness("club") == "" {
		t.Fatal("took something off nobody")
	}
}

func TestNotYourOwnCrew(t *testing.T) {
	t.Parallel()
	w, _ := mugger(t)
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
	if n := w.NPC("leo"); n != nil {
		n.Location = "club"
	}
	for i := range w.NPCs {
		if w.NPCs[i].ID != "leo" {
			w.NPCs[i].Location = "estate"
		}
	}
	if mark, ok := w.MuggingTarget("club"); ok {
		t.Fatalf("the player was offered %s, who works for them", mark.Name)
	}
}

func TestWhatSomebodyIsCarryingComesFromWhoTheyAre(t *testing.T) {
	t.Parallel()
	w, mark := mugger(t)
	soldier := *mark
	soldier.Rank, soldier.Faction = RankSoldier, ""
	boss := *mark
	boss.Rank = RankLeader
	// What somebody is carrying is now their own money, so Pockets reports
	// what they have rather than working it out from their rank on the spot.
	// The rule this test was written for moved to StandingPurse: that is what
	// a person of this standing would have on an ordinary day, and it is what
	// the city settles a newcomer onto and pays them toward.
	if w.StandingPurse(&boss) <= w.StandingPurse(&soldier) {
		t.Fatalf("a boss would carry $%d and a soldier $%d", w.StandingPurse(&boss), w.StandingPurse(&soldier))
	}
	dead := *mark
	dead.Dead = true
	if w.Pockets(&dead) != 0 {
		t.Fatal("a dead man was carrying money")
	}
	// A rich organization puts more in its people's pockets — which is now a
	// statement about what they are paid toward rather than about what is in
	// their hand this second. Somebody's own money does not change because
	// their family had a good morning.
	poor := w.StandingPurse(mark)
	if f := w.faction(mark.Faction); f != nil {
		f.Cash += 40000
	}
	if w.StandingPurse(mark) <= poor {
		t.Fatal("a richer organization did not put more in anybody's pocket")
	}
}

func TestStandingMakesYouBetterAtThisAndWorseAtGettingAway(t *testing.T) {
	t.Parallel()
	w, _ := mugger(t)
	w.Player.Respect = RecognisedAt - 1
	if w.Recognised(w.OwnHands()) {
		t.Fatal("a man nobody had heard of was recognised")
	}
	w.Player.Respect = RecognisedAt
	if !w.Recognised(w.OwnHands()) {
		t.Fatal("a known face was not recognised")
	}
	// And somebody sent is the man being described instead.
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
	hand, _ := w.CrewHands()
	if w.Recognised(hand) {
		t.Fatal("the player was named for something somebody else did")
	}
}

func TestTakingItOffSomebodyIsMoneyAndAGrievance(t *testing.T) {
	t.Parallel()
	const runs = 400
	took, madeEnemy, hurt := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, mark := mugger(t)
		w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
		for i := uint32(0); i < seed%11; i++ {
			w.Random()
		}
		w.Player.Respect = 60
		w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
		cash := w.Player.Cash
		if err := w.Mug("club", w.OwnHands()); err != nil {
			t.Fatal(err)
		}
		if w.Player.Cash > cash {
			took++
		}
		if w.Player.Health < 100 {
			hurt++
		}
		// The player is not a person in the city's records, so what the mark
		// answers with is their organization's standing.
		if f := w.faction(mark.Faction); f != nil && f.Goodwill < 0 {
			madeEnemy++
		}
	}
	if took == 0 || took == runs {
		t.Fatalf("it worked %d times in %d", took, runs)
	}
	if hurt == 0 {
		t.Fatal("nobody ever fought back")
	}
	if madeEnemy != runs {
		t.Fatalf("only %d of %d marks' organizations took any notice", madeEnemy, runs)
	}
	t.Logf("%d attempts on a man in a casino: took the money %d times, came off worse %d times, and every single one of their organizations took notice", runs, took, hurt)
}

func TestBeingRecognisedIsWhatItCosts(t *testing.T) {
	t.Parallel()
	const runs = 300
	quiet, known := 0, 0
	for _, respect := range []int{0, 80} {
		for seed := uint32(1); seed <= runs; seed++ {
			w, mark := mugger(t)
			w.RNG, w.WorldRNG = seed*2654435761, seed*2654435761
			w.Player.Respect = respect
			f := w.faction(mark.Faction)
			before := f.Goodwill
			w.Mug("club", w.OwnHands())
			if f.Goodwill < before {
				if respect == 0 {
					quiet++
				} else {
					known++
				}
			}
		}
	}
	if known <= quiet {
		t.Fatalf("a known face cost %d in standing against %d for a nobody, across %d attempts each", known, quiet, runs)
	}
	t.Logf("across %d attempts each: a nobody's organization took notice %d times, a known face's %d times", runs, quiet, known)
}

func TestRobbingAManWithATitleIsNotRobbingAMan(t *testing.T) {
	t.Parallel()
	w := New(101)
	w.MigrateLivingWorld()
	w.Player.Location, w.Player.Health, w.Player.Cash, w.Player.Respect = CityHall, 100, 20000, 60
	if err := w.Retain("commissioner"); err != nil {
		t.Fatal(err)
	}
	for i := range w.NPCs {
		if w.NPCs[i].ID != "commissioner" {
			w.NPCs[i].Location = "estate"
		}
	}
	mark, ok := w.MuggingTarget(CityHall)
	if !ok || mark.ID != "commissioner" {
		t.Fatal("the commissioner was not there to be robbed")
	}
	if w.Pockets(mark) <= 100 {
		t.Fatalf("a man with a salary and an arrangement carried $%d", w.Pockets(mark))
	}
	// Force the outcome that matters rather than waiting for it.
	found := false
	for seed := uint32(1); seed <= 400 && !found; seed++ {
		probe := New(101)
		probe.MigrateLivingWorld()
		probe.Player.Location, probe.Player.Health, probe.Player.Cash, probe.Player.Respect = CityHall, 100, 20000, 60
		probe.RNG = seed * 2654435761
		probe.Retain("commissioner")
		for i := range probe.NPCs {
			if probe.NPCs[i].ID != "commissioner" {
				probe.NPCs[i].Location = "estate"
			}
		}
		cash := probe.Player.Cash
		probe.Mug(CityHall, probe.OwnHands())
		if probe.Player.Cash > cash {
			found = true
			if probe.Retained("commissioner") {
				t.Fatal("a commissioner kept taking money from the man who robbed him")
			}
			if !probe.hasRecord("Of all the people in this city") {
				t.Fatal("nobody mentioned who it had been")
			}
		}
	}
	if !found {
		t.Fatal("nobody ever successfully robbed the commissioner in 400 attempts")
	}
}

func TestSendingSomebodyKeepsYourNameOutOfIt(t *testing.T) {
	t.Parallel()
	w, mark := mugger(t)
	w.Player.Respect = 80
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 90}}
	hand, _ := w.CrewHands()
	if err := w.Mug("club", hand); err != nil {
		t.Fatal(err)
	}
	if w.Player.Health != 100 {
		t.Fatal("the player was hurt by something they were not at")
	}
	for _, g := range w.Grudges {
		if g.Holder == mark.ID && g.Against == "leo" {
			return
		}
	}
	t.Fatal("the mark held it against nobody")
}
