package core

import (
	"strings"
	"testing"
)

func boss(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := proprietor(t)
	own(w, "laundry", "garage")
	w.Player.Respect, w.Player.Contacts = OrganizationStanding, 3
	w.OrganizationDay()
	if !w.Incorporated() {
		t.Fatal("the fixture did not become an organization")
	}
	var candidate *NPC
	for _, n := range w.Civilians() {
		if !IsOfficial(n.ID) && !w.isRoleHolder(n) {
			candidate = n
			break
		}
	}
	if candidate == nil {
		t.Fatal("nobody in the city was available")
	}
	w.Player.Location = candidate.Location
	candidate.Trust = 1 // somebody the player has dealt with
	return w, candidate
}

func TestNobodySignsOnWithAMan(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	own(w, "laundry")
	w.Player.Respect, w.Player.Contacts = 5, 3
	w.OrganizationDay()
	for _, n := range w.Civilians() {
		w.Player.Location = n.Location
		if w.SignOnReadiness(n.ID) == "" {
			t.Fatal("somebody signed on with a man who was not an organization")
		}
		break
	}
}

func TestPuttingSomebodyOnCostsAndBinds(t *testing.T) {
	t.Parallel()
	w, candidate := boss(t)
	if w.SignOnReadiness(candidate.ID) != "" {
		t.Fatal("could not put anybody on:", w.SignOnReadiness(candidate.ID))
	}
	cash, bill := w.Player.Cash, w.DailyCost()
	if err := w.SignOn(candidate.ID); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-SigningCost {
		t.Fatalf("signing cost $%d", cash-w.Player.Cash)
	}
	if w.DailyCost() != bill+MemberWage {
		t.Fatalf("the bill went %d to %d", bill, w.DailyCost())
	}
	if len(w.OwnPeople()) != 1 || w.OwnPeople()[0].ID != candidate.ID {
		t.Fatal("they did not answer to the player afterwards")
	}
	if candidate.Trust != MemberTrust {
		t.Fatalf("they arrived at %d trust", candidate.Trust)
	}
	if w.SignOnReadiness(candidate.ID) == "" {
		t.Fatal("signed the same person on twice")
	}
}

func TestYouCannotTakeOnSomebodyElsesPeople(t *testing.T) {
	t.Parallel()
	w, _ := boss(t)
	for _, n := range w.Members("bellandi") {
		w.Player.Location = n.Location
		if w.SignOnReadiness(n.ID) == "" {
			t.Fatalf("%s signed on while answering to Bellandi", n.Name)
		}
		break
	}
	// Nor anybody with a job the city needs doing.
	for _, r := range roles {
		holder := w.Holder(r.ID)
		w.Player.Location = holder.Location
		if w.SignOnReadiness(holder.ID) == "" {
			t.Fatalf("the %s signed on", r.Title)
		}
	}
	if w.SignOnReadiness("commissioner") == "" {
		t.Fatal("the commissioner signed on")
	}
}

func TestYourPeopleMakeYouHarderToMoveOn(t *testing.T) {
	t.Parallel()
	w, candidate := boss(t)
	before := w.PlayerStrength()
	w.SignOn(candidate.ID)
	if w.PlayerStrength() <= before {
		t.Fatalf("somebody answering to you was worth nothing: %d against %d", w.PlayerStrength(), before)
	}
	weak := w.PlayerStrength()
	candidate.Trust = 100
	if w.PlayerStrength() <= weak {
		t.Fatal("somebody who means it was worth no more than somebody who does not")
	}
}

func TestARaidReachesYourOwnPeople(t *testing.T) {
	t.Parallel()
	found := false
	for seed := uint32(1); seed <= 300 && !found; seed++ {
		w, candidate := boss(t)
		w.WorldRNG = seed * 2654435761
		w.SignOn(candidate.ID)
		attacker := &w.Factions[0]
		attacker.Power = 100
		for i := 0; i < 20 && !candidate.Dead; i++ {
			w.contest(attacker, w.PlayerOrganization())
		}
		if candidate.Dead {
			found = true
		}
	}
	if !found {
		t.Fatal("no raid in 300 ever reached anybody who answered to the player")
	}
}

func TestNeglectIsWhatCostsYouPeople(t *testing.T) {
	t.Parallel()
	w, candidate := boss(t)
	w.SignOn(candidate.ID)
	// Paid, and they settle.
	w.Player.Cash = 50000
	for i := 0; i < 30; i++ {
		w.OwnPeopleDay()
	}
	if candidate.Trust <= MemberTrust {
		t.Fatalf("a month of being paid left them at %d", candidate.Trust)
	}
	// Not paid, and they do not.
	w.Player.Cash = 0
	for i := 0; i < 20; i++ {
		w.OwnPeopleDay()
	}
	if candidate.Trust >= MemberTrust {
		t.Fatalf("three weeks of unpaid bills left them at %d", candidate.Trust)
	}
}

func TestSomebodyFarEnoughDownLeaves(t *testing.T) {
	t.Parallel()
	const runs = 300
	left, tookSomething := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, candidate := boss(t)
		w.WorldRNG = seed * 2654435761
		// The fixture picks the same person every time, so their temperament is
		// set deliberately rather than left to whatever that name produces.
		candidate.Name = nameWith(t, "hot")
		w.SignOn(candidate.ID)
		candidate.Trust, candidate.Ambition = 0, 90
		w.Player.Cash = 0
		for i := 0; i < 40 && len(w.OwnPeople()) > 0; i++ {
			w.OwnPeopleDay()
		}
		if len(w.OwnPeople()) == 0 {
			left++
			if len(w.FamilyHoldings(w.PlayerOrganizationID())) < 2 {
				tookSomething++
			}
		}
	}
	if left == 0 {
		t.Fatal("nobody at zero trust ever left")
	}
	if tookSomething == 0 {
		t.Fatal("nobody ambitious ever walked out with a business")
	}
	t.Logf("%d of %d campaigns lost somebody who had stopped being paid, and %d of those walked out with a business", left, runs, tookSomething)
}

func TestALoyalManDoesNotTakeYourBusiness(t *testing.T) {
	t.Parallel()
	w, candidate := boss(t)
	candidate.Name = nameWith(t, "loyal")
	w.SignOn(candidate.ID)
	candidate.Trust, candidate.Ambition = 0, 100
	held := len(w.FamilyHoldings(w.PlayerOrganizationID()))
	for i := 0; i < 400 && len(w.OwnPeople()) > 0; i++ {
		w.Player.Cash = 0
		w.OwnPeopleDay()
	}
	if len(w.FamilyHoldings(w.PlayerOrganizationID())) != held {
		t.Fatal("a loyal man walked out with a business")
	}
}

func TestPuttingSomebodyOutAndPayingAShare(t *testing.T) {
	t.Parallel()
	w, candidate := boss(t)
	w.SignOn(candidate.ID)
	candidate.Trust = 30
	cash := w.Player.Cash
	if err := w.PayShare(candidate.ID); err != nil {
		t.Fatal(err)
	}
	if candidate.Trust != 55 || w.Player.Cash != cash-60 {
		t.Fatalf("trust %d, cash %d", candidate.Trust, w.Player.Cash)
	}
	if err := w.LetGo(candidate.ID); err != nil {
		t.Fatal(err)
	}
	if len(w.OwnPeople()) != 0 || candidate.Faction != "" {
		t.Fatal("they still answered to the player")
	}
	if w.PayShareReadiness(candidate.ID) == "" || w.LetGoReadiness(candidate.ID) == "" {
		t.Fatal("somebody who does not answer to you could still be paid or put out")
	}
}

// Somebody of yours can walk out with one of your businesses, and the only
// warning was a number on a card measured against a line the player has never
// been told. The game even says "it had been coming" the morning after; this is
// the part that comes before.
func TestSomebodyOnTheirWayOutIsSaidSo(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect = 100, 40
	w.Player.Cash = 0 // nothing covered, so they come down rather than up
	hand := w.Holder("driver")
	if hand == nil {
		t.Fatal("nobody to sign on")
	}
	hand.Faction, hand.Rank, hand.Trust = w.PlayerOrganizationID(), RankAssociate, DefectionTrust+3
	if hand.Trust < DefectionTrust {
		t.Fatal("they are already past it, so the crossing cannot be seen")
	}

	said, marked := false, false
	for day := 0; day < 6 && !said; day++ {
		before := len(w.History)
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		for _, r := range w.History[min(before, len(w.History)):] {
			if strings.Contains(r.Title, "thinking about it") {
				said = true
			}
		}
	}
	for _, p := range w.PeopleHere(hand.Location) {
		if p.ID == hand.ID {
			marked = p.Restless
		}
	}
	t.Logf("%s is at %d trust: told=%v, marked on the card=%v", hand.Name, hand.Trust, said, marked)
	if !said {
		t.Fatal("somebody got far enough down to walk and nobody said anything")
	}
	if !marked {
		t.Fatal("their card says nothing about it")
	}

	// And it is said once rather than every morning. Somebody with no ambition
	// at all never decides to go, so they stay below the line for as long as
	// this needs — otherwise they walk out before the second morning and the
	// repeat can never be seen.
	w.NPC(hand.ID).Ambition = 0
	again := 0
	for day := 0; day < 6; day++ {
		before := len(w.History)
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		for _, r := range w.History[min(before, len(w.History)):] {
			if strings.Contains(r.Title, "thinking about it") {
				again++
			}
		}
		if w.NPC(hand.ID).Faction != w.PlayerOrganizationID() {
			t.Fatal("somebody with no ambition at all walked out anyway")
		}
	}
	if again > 0 {
		t.Fatalf("the player was told %d more times about the same man", again)
	}
}

// Somebody who is paid and content is not marked, which is the ordinary case.
func TestSomebodyContentIsNotMarked(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 100000
	hand := w.Holder("driver")
	hand.Faction, hand.Rank, hand.Trust = w.PlayerOrganizationID(), RankAssociate, 65
	for day := 0; day < 5; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	for _, p := range w.PeopleHere(hand.Location) {
		if p.ID == hand.ID && p.Restless {
			t.Fatalf("%s is paid and at %d trust and is marked as leaving", hand.Name, hand.Trust)
		}
	}
	for _, r := range w.History {
		if strings.Contains(r.Title, "thinking about it") {
			t.Fatal("a contented man was reported as thinking about leaving")
		}
	}
}
