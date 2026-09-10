package core

import "testing"

// A player with money, standing, and somebody in the room to lend it to.
func lender(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := proprietor(t)
	w.Player.Cash, w.Player.Respect = 6000, 40
	w.District = 2 // people live all over this city and money travels
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) {
			continue
		}
		w.Player.Location = n.Location
		if w.LendReadiness(n.ID) == "" {
			return w, n
		}
	}
	t.Fatal("nobody in this city would take money")
	return nil, nil
}

// defaulter is a loan that genuinely goes bad: somebody worth lending to whose
// position collapses while they are carrying it. The floor loan is serviceable
// by anybody alive, so a default in this game is a man whose circumstances
// changed rather than a man who was always poor.
func defaulter(t *testing.T) (*World, *NPC) {
	t.Helper()
	w, n := lender(t)
	n.Rank, n.Faction, n.Ambition = RankLieutenant, "bellandi", 10
	if f := w.faction("bellandi"); f != nil {
		f.Cash = 200000
	}
	w.Player.Cash = 40000
	if err := w.Lend(n.ID); err != nil {
		t.Fatal(err)
	}
	n.Rank, n.Faction, n.Skill = 0, "", 5
	w.Minute = w.LoanTo(n.ID).Due
	w.LoanDay()
	if l := w.LoanTo(n.ID); l == nil || l.Missed == 0 {
		t.Fatal("the loan did not go bad")
	}
	return w, n
}

func TestNobodyBorrowsFromAStranger(t *testing.T) {
	t.Parallel()
	w, n := lender(t)
	w.Player.Respect = LendStanding - 1
	if w.LendReadiness(n.ID) == "" {
		t.Fatal("a man nobody has heard of was lending money")
	}
	w.Player.Respect = 40

	// Not across town, and not to a man with a title.
	here := n.Location
	n.Location = "docks"
	if here == "docks" {
		n.Location = "club"
	}
	if w.LendReadiness(n.ID) == "" {
		t.Fatal("money changed hands across the city")
	}
	n.Location = w.Player.Location
	if reason := w.LendReadiness(n.ID); reason != "" {
		t.Fatal("nobody would take it:", reason)
	}
	// And never twice to the same man at once.
	if err := w.Lend(n.ID); err != nil {
		t.Fatal(err)
	}
	if w.LendReadiness(n.ID) == "" {
		t.Fatal("the same man borrowed twice")
	}
}

func TestLendingIsBoundedByWhatYouCanAffordToBeOwed(t *testing.T) {
	t.Parallel()
	w, _ := lender(t)
	w.Player.Cash = 6000
	lent, borrowers := 0, 0
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) {
			continue
		}
		w.Player.Location = n.Location
		if w.LendReadiness(n.ID) != "" {
			continue
		}
		before := w.Player.Cash
		if err := w.Lend(n.ID); err != nil {
			t.Fatal(err)
		}
		lent += before - w.Player.Cash
		borrowers++
	}
	if borrowers < 2 {
		t.Fatalf("only %d people in the city would borrow anything", borrowers)
	}
	if out := w.OutOnLoan(); out != lent {
		t.Fatalf("the book says $%d and $%d left the pocket", out, lent)
	}
	// Never more than the share of what there was to begin with.
	if float64(lent) > float64(6000)*LoanShare+1 {
		t.Fatalf("$%d of $6000 went out on loan", lent)
	}
	if w.Player.Cash <= 0 {
		t.Fatal("a lender lent themselves to nothing")
	}
}

func TestSomebodyWhoCanPayPaysAndSomebodyWhoCannotIsADecision(t *testing.T) {
	t.Parallel()
	w, n := lender(t)
	if err := w.Lend(n.ID); err != nil {
		t.Fatal(err)
	}
	l := w.LoanTo(n.ID)
	principal, owed := l.Principal, l.Owed

	// Not before the day.
	cash := w.Player.Cash
	w.LoanDay()
	if w.Player.Cash != cash || w.LoanTo(n.ID) == nil {
		t.Fatal("a loan settled before it was due")
	}

	// A man carrying enough pays, and thinks better of you for it.
	w.Minute = l.Due
	// Somebody who can find it and has no reason not to: nerve is willingness,
	// and a man with no ambition simply pays.
	n.Ambition = 0
	n.Rank, n.Faction = RankLieutenant, "bellandi"
	if f := w.faction("bellandi"); f != nil {
		f.Cash = 400000
	}
	if !w.canPay(n, owed) {
		t.Skip("even a lieutenant of a rich family could not carry it")
	}
	trust := n.Trust
	w.LoanDay()
	if w.LoanTo(n.ID) != nil {
		t.Fatal("a man who paid still owes it")
	}
	if w.Player.Cash != cash+owed {
		t.Fatalf("$%d came back on a $%d loan owing $%d", w.Player.Cash-cash, principal, owed)
	}
	if n.Trust <= trust {
		t.Fatal("paying you back was worth nothing")
	}
}

func TestAManWhoCannotPayIsAManStandingInFrontOfYou(t *testing.T) {
	t.Parallel()
	w, n := defaulter(t)
	if l := w.LoanTo(n.ID); l.Missed != 1 {
		t.Fatalf("missed %d times", l.Missed)
	}
	// Now there are three answers, and all three are in the room.
	offered := map[string]bool{}
	for _, a := range w.Actions(w.Player.Location) {
		if !a.Disabled {
			offered[a.ID] = true
		}
	}
	for _, want := range []string{"lean:" + n.ID, "extend:" + n.ID, "forgive:" + n.ID} {
		if !offered[want] {
			t.Fatalf("%s was not offered: %v", want, offered)
		}
	}
	// And leaving it standing costs standing every day.
	respect := w.Player.Respect
	w.LoanDay()
	if w.Player.Respect >= respect {
		t.Fatal("a man owed money who did nothing about it lost nothing")
	}
}

func TestCollectingWorksAndCostsSomethingEveryTime(t *testing.T) {
	t.Parallel()
	w, n := defaulter(t)
	w.Player.Respect, w.Player.Heat = 60, 0
	respect, heat, cash := w.Player.Respect, w.Player.Heat, w.Player.Cash
	w.RNG = 1
	if err := w.Lean(n.ID); err != nil {
		t.Fatal(err)
	}
	if w.Player.Heat <= heat {
		t.Fatal("a collection in the street drew no attention")
	}
	settled := w.LoanTo(n.ID) == nil
	if settled {
		if w.Player.Cash <= cash {
			t.Fatal("a collection that worked collected nothing")
		}
		if w.Player.Respect <= respect {
			t.Fatal("the street thought nothing of it")
		}
		if n.Trust != 0 {
			t.Fatal("he thinks as well of you as he did")
		}
	}
	// Either way he remembers it, and it is the one thing in this city
	// somebody can hold against the player personally.
	if n.Sore == 0 {
		t.Fatal("nobody held anything against a man who came for their money")
	}
	if n.SoreAt == "" {
		t.Fatal("he could not say what it was about")
	}
}

func TestAnotherWeekIsHowThisTradeMakesItsMoney(t *testing.T) {
	t.Parallel()
	w, n := defaulter(t)
	owed, trust := w.LoanTo(n.ID).Owed, n.Trust
	if err := w.Extend(n.ID); err != nil {
		t.Fatal(err)
	}
	l := w.LoanTo(n.ID)
	if l.Owed <= owed {
		t.Fatalf("another week cost him nothing: $%d against $%d", l.Owed, owed)
	}
	if l.Missed != 0 || l.Due <= w.Minute {
		t.Fatal("the clock did not start again")
	}
	if n.Trust <= trust {
		t.Fatal("he was not grateful for it")
	}
	if w.ExtendReadiness(n.ID) == "" {
		t.Fatal("it could be extended again the same afternoon")
	}
}

func TestWritingItOffBuysTheOneThingMoneyCannot(t *testing.T) {
	t.Parallel()
	w, n := defaulter(t)
	w.Aggrieve(n.ID, 50, "something older")
	w.Player.Respect = 60
	trust, cash := n.Trust, w.Player.Cash
	if err := w.Forgive(n.ID); err != nil {
		t.Fatal(err)
	}
	if w.LoanTo(n.ID) != nil {
		t.Fatal("it is still on the books")
	}
	if w.Player.Cash != cash {
		t.Fatal("writing it off returned money")
	}
	if w.Player.Respect != 60-ForgiveRespect {
		t.Fatalf("the street charged %d for it", 60-w.Player.Respect)
	}
	if n.Trust <= trust+20 {
		t.Fatalf("he thinks of you at %d", n.Trust)
	}
	if n.Sore != 0 {
		t.Fatalf("he still holds %d against you", n.Sore)
	}
}

func TestADebtDiesWithTheManCarryingIt(t *testing.T) {
	t.Parallel()
	w, n := lender(t)
	w.Lend(n.ID)
	l := w.LoanTo(n.ID)
	w.Minute = l.Due
	n.Dead = true
	cash := w.Player.Cash
	w.LoanDay()
	if w.LoanTo(n.ID) != nil {
		t.Fatal("a dead man still owes it")
	}
	if w.Player.Cash != cash {
		t.Fatal("a dead man paid")
	}
}

func TestTheBookSurvivesBeingWrittenDown(t *testing.T) {
	t.Parallel()
	w, n := lender(t)
	w.Lend(n.ID)
	book := w.LoanDescription()
	if len(book) != 1 || book[0]["name"] != n.Name {
		t.Fatalf("the book reads %v", book)
	}
	if book[0]["owed"].(int) <= book[0]["principal"].(int) {
		t.Fatal("nothing was owed on top")
	}
}

func TestSomebodyPutAgainstAWallComesBackForIt(t *testing.T) {
	t.Parallel()
	// A collection is not free the moment it happens. Over enough days, the man
	// it happened to is the one who picks the player out of a street.
	const runs = 400
	picked := func(sore int) int {
		hits := 0
		for seed := uint32(1); seed <= runs; seed++ {
			w, n := lender(t)
			w.WorldRNG = seed * 2654435761
			n.Sore = sore
			if thief := w.robber(); thief != nil && thief.ID == n.ID {
				hits++
			}
		}
		return hits
	}
	quiet, sore := picked(0), picked(45)
	if sore <= quiet {
		t.Fatalf("a man with a reason was picked %d times in %d against %d for a man without one", sore, runs, quiet)
	}
	t.Logf("of %d robberies, the man who was put against a wall is the thief %d times against %d for a stranger", runs, sore, quiet)

	// And it fades, because nothing in this city is held forever.
	w, n := lender(t)
	n.Sore, n.SoreAt = 3, "something"
	for i := 0; i < 4; i++ {
		w.GrudgeDay()
	}
	if n.Sore != 0 || n.SoreAt != "" {
		t.Fatalf("he is still carrying %d of it", n.Sore)
	}
}

func TestARoomOfStrangersIsNotFifteenIdenticalButtons(t *testing.T) {
	t.Parallel()
	w, _ := lender(t)
	// Put a crowd in one room, all of them plausible borrowers.
	where := w.Player.Location
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && !IsOfficial(n.ID) && !w.isRoleHolder(n) {
			n.Location = where
		}
	}
	w.Player.Cash = 400000 // enough that the book, not the pocket, is the limit
	offers, owing := 0, 0
	for _, a := range w.Actions(where) {
		if len(a.ID) > 5 && a.ID[:5] == "lend:" {
			offers++
		}
	}
	if offers == 0 || offers > LendOffers {
		t.Fatalf("%d people were offered money in one room", offers)
	}
	// But everybody who already owes you is listed, however many that is.
	lent := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != where || w.LendReadiness(n.ID) != "" || lent >= LendOffers+3 {
			continue
		}
		if w.Lend(n.ID) == nil {
			lent++
		}
	}
	if lent <= LendOffers {
		t.Skipf("only %d people could be lent to", lent)
	}
	for _, a := range w.Actions(where) {
		if len(a.ID) > 5 && a.ID[:5] == "lean:" {
			owing++
		}
	}
	if owing != lent {
		t.Fatalf("%d people owe money and %d of them are on the screen", lent, owing)
	}
}

func TestNobodyWithATitleBorrowsFromSomebodyWithout(t *testing.T) {
	t.Parallel()
	w, _ := lender(t)
	for _, role := range []string{"detective", "fixer"} {
		id := w.HolderID(role)
		if id == "" {
			continue
		}
		n := w.NPC(id)
		if n == nil {
			continue
		}
		n.Location = w.Player.Location
		if w.LendReadiness(id) == "" {
			t.Fatalf("%s took money off the player", n.Name)
		}
	}
}

// The books counted money owed by a man who had been shot. LoanDay already
// knew to write it off — "$249 went out and whoever was carrying it is not
// carrying anything now" — but it only runs at midnight, so between the
// killing and the next morning the ledger claimed an asset that was in the
// ground. That window is exactly when a player looks at their books.
func TestADebtDiesWithTheManWhoOwedIt(t *testing.T) {
	t.Parallel()
	w := New(4)
	w.Player.Cash, w.Player.Respect = 5000, OrganizationStanding
	w.Player.Location = "bar"
	var borrower *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && !IsOfficial(n.ID) && w.LendReadiness(n.ID) == "" {
			borrower = n
			break
		}
	}
	if borrower == nil {
		t.Fatal("nobody in this city would borrow")
	}
	if err := w.Lend(borrower.ID); err != nil {
		t.Fatal(err)
	}
	if w.Books()["owed"].(int) <= 0 || len(w.Book()) != 1 {
		t.Fatal("the loan was never on the books")
	}
	w.Kill(borrower.ID, "Shot over something else entirely.")
	if owed := w.Books()["owed"].(int); owed != 0 {
		t.Fatalf("the books still say $%d is owed by a dead man", owed)
	}
	if len(w.Book()) != 0 {
		t.Fatalf("the loan is still in the book: %+v", w.Book())
	}
	// And the player is told, at the moment it stops being true.
	told := false
	for _, r := range w.History {
		if r.Title == "Nothing to collect" {
			told = true
		}
	}
	if !told {
		t.Fatal("the money was written off and nobody said so")
	}
}
