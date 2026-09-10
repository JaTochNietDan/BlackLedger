package core

import "testing"

// diplomat builds a player who is an organization, on good terms with one
// family, and sharing a problem with them.
func diplomat(t *testing.T) (*World, *Faction, *Faction) {
	t.Helper()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect, w.Player.Contacts, w.Player.Location = OrganizationStanding, 3, "market"
	w.OrganizationDay()
	if len(w.Factions) < 3 {
		t.Fatal("the fixture needs two families and the player")
	}
	friend, enemy := &w.Factions[0], &w.Factions[1]
	friend.Goodwill = 40
	// Somebody they both have a problem with.
	for _, pair := range [][2]string{{friend.ID, enemy.ID}, {w.PlayerOrganizationID(), enemy.ID}} {
		c := w.Conflict(pair[0], pair[1])
		c.Hostility, c.State = 60, "feud"
	}
	// And enough known about the friend to know what is being agreed.
	if members := w.Members(friend.ID); len(members) > 0 {
		members[0].Trust = 1
	}
	return w, friend, enemy
}

func TestNobodyMakesThisArrangementWithAMan(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry")
	w.Player.Respect = 5
	w.OrganizationDay()
	if w.PactReadiness(w.Factions[0].ID) == "" {
		t.Fatal("a proprietor reached an understanding with a family")
	}
}

func TestItTakesStandingAndASharedProblem(t *testing.T) {
	t.Parallel()
	w, friend, enemy := diplomat(t)
	if w.PactReadiness(friend.ID) != "" {
		t.Fatal("could not reach an understanding:", w.PactReadiness(friend.ID))
	}
	friend.Goodwill = PactGoodwill - 1
	if w.PactReadiness(friend.ID) == "" {
		t.Fatal("somebody who barely tolerates the player agreed to stand with them")
	}
	friend.Goodwill = 40
	// Take the shared problem away.
	w.Conflict(friend.ID, enemy.ID).State = "cold"
	if w.PactReadiness(friend.ID) == "" {
		t.Fatal("two people with no common enemy agreed to stand together")
	}
}

func TestYouHaveToKnowWhatYouAreAgreeingTo(t *testing.T) {
	t.Parallel()
	w, friend, _ := diplomat(t)
	for _, n := range w.Members(friend.ID) {
		n.Trust = 0
	}
	w.Player.Contacts = 2 // reputation only
	if w.Intelligence(friend.ID) >= 2 {
		t.Fatal("the fixture already knew too much")
	}
	if w.PactReadiness(friend.ID) == "" {
		t.Fatal("agreed to stand with somebody they knew nothing about")
	}
}

func TestStandingWithSomebodyPutsYouInTheirQuarrels(t *testing.T) {
	t.Parallel()
	w, friend, enemy := diplomat(t)
	before := w.Conflict(enemy.ID, w.PlayerOrganizationID()).Hostility
	cash, bill := w.Player.Cash, w.DailyCost()
	if err := w.MakePact(friend.ID); err != nil {
		t.Fatal(err)
	}
	if !w.Allied(friend.ID) {
		t.Fatal("the understanding did not hold")
	}
	if w.Player.Cash != cash-PactOpening {
		t.Fatalf("it cost $%d to open", cash-w.Player.Cash)
	}
	if w.DailyCost() != bill+PactTribute {
		t.Fatalf("the bill went %d to %d", bill, w.DailyCost())
	}
	if w.Conflict(enemy.ID, w.PlayerOrganizationID()).Hostility <= before {
		t.Fatal("standing with somebody was not noticed by the people they fight")
	}
	// And it goes on being noticed.
	drift := w.Conflict(enemy.ID, w.PlayerOrganizationID()).Hostility
	w.Player.Cash = 50000
	w.PactDay()
	if w.Conflict(enemy.ID, w.PlayerOrganizationID()).Hostility <= drift {
		t.Fatal("a day of standing with them cost nothing")
	}
}

func TestNeitherOfYouMovesOnTheOther(t *testing.T) {
	t.Parallel()
	w, friend, _ := diplomat(t)
	w.MakePact(friend.ID)
	held := w.FamilyHoldings(friend.ID)
	if len(held) == 0 {
		t.Fatal("the ally held nothing")
	}
	w.Player.Location, w.Player.Health = held[0], 100
	if w.MoveOnReadiness(held[0]) == "" {
		t.Fatal("the player moved on somebody they stand with")
	}
	// And their raids on the player do not land.
	mine := w.FamilyHoldings(w.PlayerOrganizationID())
	condition := w.Properties[mine[0]].Condition
	for i := 0; i < 30; i++ {
		w.contestAt(friend, w.PlayerOrganization(), mine[0])
	}
	if w.Properties[mine[0]].Condition != condition {
		t.Fatal("somebody the player stands with raided them anyway")
	}
}

func TestSomebodyWhoStandsWithYouMayAnswerTheDoor(t *testing.T) {
	t.Parallel()
	const runs = 400
	answered := 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, friend, enemy := diplomat(t)
		w.WorldRNG = seed * 2654435761
		w.MakePact(friend.ID)
		friend.Power = 100
		mine := w.FamilyHoldings(w.PlayerOrganizationID())
		condition := w.Properties[mine[0]].Condition
		w.contestAt(enemy, w.PlayerOrganization(), mine[0])
		if w.Properties[mine[0]].Condition == condition && w.hasRecord(friend.Name+" was already there") {
			answered++
		}
	}
	if answered == 0 {
		t.Fatal("nobody who stood with the player ever turned up")
	}
	if answered == runs {
		t.Fatal("an understanding made the player untouchable")
	}
	t.Logf("%d raids on somebody with an understanding: the ally was already there %d times", runs, answered)
}

func TestAnUnderstandingNobodyPaysForLapses(t *testing.T) {
	t.Parallel()
	w, friend, _ := diplomat(t)
	w.MakePact(friend.ID)
	w.Player.Cash = 0
	w.PactDay()
	if w.Allied(friend.ID) {
		t.Fatal("an unpaid understanding held")
	}
	if !w.hasRecord("The understanding with " + friend.Name + " is over") {
		t.Fatal("nobody was told")
	}
}

func TestNobodyInheritsSomebodyElsesFriends(t *testing.T) {
	t.Parallel()
	w, friend, _ := diplomat(t)
	w.MakePact(friend.ID)
	w.Player.Alive = false
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "new_life"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Allied(friend.ID) || len(next.Pacts) != 0 || next.PactCost() != 0 {
		t.Fatal("the next arrival stood with somebody")
	}
}
