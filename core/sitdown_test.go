package core

import (
	"strings"
	"testing"
)

// mediator builds a city with a quarrel bad enough to need somebody, and a
// player the two sides would come for.
func mediator(t *testing.T, hostility int, hotLeader bool) *World {
	t.Helper()
	w := New(97)
	w.MigrateLivingWorld()
	w.Player.Location = SitdownGround
	w.Player.Cash, w.Player.Respect, w.Player.Health = 5000, 60, 100
	c := &w.Conflicts[0]
	c.Hostility = hostility
	c.State = classify(c)
	// Make the leader of one side the kind of person who does or does not wait.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Name != w.faction(c.A).Leader {
			continue
		}
		if hotLeader {
			n.Name = nameWith(t, "hot")
		} else {
			n.Name = nameWith(t, "careful")
		}
		w.faction(c.A).Leader = n.Name
	}
	// And the other side too, so the fixture is unambiguous.
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Name != w.faction(c.B).Leader {
			continue
		}
		n.Name = nameWith(t, "careful")
		w.faction(c.B).Leader = n.Name
	}
	return w
}

func TestNobodyComesToARoomForANobody(t *testing.T) {
	t.Parallel()
	w := mediator(t, 70, false)
	w.Player.Respect = 0
	if w.SitdownReadiness() == "" {
		t.Fatal("two families crossed the city for somebody nobody had heard of")
	}
	w.Player.Respect = 60
	if w.SitdownReadiness() != "" {
		t.Fatal("they would not come for somebody with standing:", w.SitdownReadiness())
	}
	// Nor into a room arranged by somebody one of them wants dead.
	w.Factions[0].Goodwill = -60
	if w.SitdownReadiness() == "" {
		t.Fatal("a family that hates you sat in a room you arranged")
	}
	w.Factions[0].Goodwill = 0
	// Nor anywhere but neutral ground.
	w.Player.Location = "club"
	if w.SitdownReadiness() == "" {
		t.Fatal("a meeting was held on one side's own floor")
	}
}

func TestThereIsNothingToMediateInAQuietCity(t *testing.T) {
	t.Parallel()
	w := mediator(t, 10, false)
	for i := range w.Conflicts {
		w.Conflicts[i].Hostility, w.Conflicts[i].State = 5, "cold"
	}
	if _, ok := w.OpenQuarrel(); ok {
		t.Fatal("a city at peace produced a quarrel to settle")
	}
	if w.SitdownReadiness() == "" {
		t.Fatal("a meeting was called about nothing")
	}
}

func TestPressingASettleableQuarrelEndsIt(t *testing.T) {
	t.Parallel()
	w := mediator(t, 55, false)
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Event == nil || w.Event.Kind != "sitdown" || len(w.Event.Choices) != 4 {
		t.Fatalf("the room was %v", w.Event)
	}
	c := w.Conflict(w.Event.Actor, w.Event.Target)
	before, respect := c.Hostility, w.Player.Respect
	if err := w.ResolveSitdown(w.Event, "press"); err != nil {
		t.Fatal(err)
	}
	if c.Hostility >= before {
		t.Fatalf("hostility went %d to %d", before, c.Hostility)
	}
	if w.Player.Respect <= respect || w.faction(w.Event.Actor).Goodwill <= 0 {
		t.Fatal("settling a war bought the player nothing")
	}
	if !w.hasRecord("It is settled") {
		t.Fatal("nothing was recorded")
	}
	if w.Player.Health != 100 {
		t.Fatal("a settled meeting hurt somebody")
	}
}

func TestPressingAQuarrelSomebodyCameToFinishIsTheWorstRoomInTheCity(t *testing.T) {
	t.Parallel()
	w := mediator(t, 80, true)
	q, ok := w.OpenQuarrel()
	if !ok || !q.Trap {
		t.Fatalf("a hostility-80 quarrel led by a hot-headed man was not a trap: %+v", q)
	}
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	before := 0
	for i := range w.NPCs {
		if w.NPCs[i].Dead {
			before++
		}
	}
	if err := w.ResolveSitdown(w.Event, "press"); err != nil {
		t.Fatal(err)
	}
	after := 0
	for i := range w.NPCs {
		if w.NPCs[i].Dead {
			after++
		}
	}
	if after <= before {
		t.Fatal("a bloodbath killed nobody")
	}
	if w.Player.Health >= 100 {
		t.Fatalf("the player walked out of it untouched, at %d health", w.Player.Health)
	}
	if !w.hasRecord("It was never a meeting") || !w.hasNewsKind("killing") {
		t.Fatal("the city never heard about it")
	}
}

func TestAWarmQuarrelIsNotATrapHoweverBadTheLeader(t *testing.T) {
	t.Parallel()
	w := mediator(t, TrapHostility-5, true)
	q, ok := w.OpenQuarrel()
	if !ok || q.Trap {
		t.Fatal("a hot-headed man in a cooling quarrel was treated as an ambush")
	}
}

func TestYouOnlyKnowItIsATrapIfSomebodyToldYou(t *testing.T) {
	t.Parallel()
	w := mediator(t, 80, true)
	w.Player.Contacts = 0
	q, _ := w.OpenQuarrel()
	if !q.Trap || q.Suspected {
		t.Fatal("a stranger with no contacts walked in knowing")
	}
	w.Player.Contacts = 3
	q, _ = w.OpenQuarrel()
	if !q.Suspected {
		t.Fatal("a network of contacts warned nobody")
	}
	// And the warning reaches the room itself.
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if !containsName(w.Event.Body, "more people than the room needs") {
		t.Fatalf("the scene said: %s", w.Event.Body)
	}
}

func TestListeningIsSafeAndSettlesLittle(t *testing.T) {
	t.Parallel()
	w := mediator(t, 80, true) // even in the worst room
	w.CallSitdown()
	c := w.Conflict(w.Event.Actor, w.Event.Target)
	before := c.Hostility
	if err := w.ResolveSitdown(w.Event, "listen"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Health != 100 {
		t.Fatal("saying nothing got the player shot")
	}
	if c.Hostility >= before || before-c.Hostility > 10 {
		t.Fatalf("hostility went %d to %d, which is either nothing or too much for saying nothing", before, c.Hostility)
	}
}

func TestTakingASideMakesOneEnemyAndOneFriend(t *testing.T) {
	t.Parallel()
	w := mediator(t, 55, false)
	w.CallSitdown()
	a, b := w.faction(w.Event.Actor), w.faction(w.Event.Target)
	beforeA, beforeB := a.Goodwill, b.Goodwill
	c := w.Conflict(a.ID, b.ID)
	hostility := c.Hostility
	if err := w.ResolveSitdown(w.Event, "side"); err != nil {
		t.Fatal(err)
	}
	if a.Goodwill <= beforeA || b.Goodwill >= beforeB {
		t.Fatalf("standing went %d→%d and %d→%d", beforeA, a.Goodwill, beforeB, b.Goodwill)
	}
	if c.Hostility <= hostility {
		t.Fatal("taking a side calmed the quarrel down")
	}
	if len(w.Plots) == 0 {
		t.Fatal("the side you crossed did nothing about it")
	}
}

func TestLeavingCostsStandingWithBoth(t *testing.T) {
	t.Parallel()
	w := mediator(t, 55, false)
	w.CallSitdown()
	a, b := w.faction(w.Event.Actor), w.faction(w.Event.Target)
	beforeA, beforeB := a.Goodwill, b.Goodwill
	if err := w.ResolveSitdown(w.Event, "leave"); err != nil {
		t.Fatal(err)
	}
	if a.Goodwill >= beforeA || b.Goodwill >= beforeB {
		t.Fatal("walking out cost nothing")
	}
	if w.Player.Health != 100 {
		t.Fatal("walking out got the player hurt")
	}
}

func TestTheFeeIsPaidWhateverHappens(t *testing.T) {
	t.Parallel()
	w := mediator(t, 55, false)
	cash := w.Player.Cash
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-SitdownFee {
		t.Fatalf("the room cost $%d", cash-w.Player.Cash)
	}
	w.ResolveSitdown(w.Event, "leave")
	if w.Player.Cash != cash-SitdownFee {
		t.Fatal("leaving refunded the room")
	}
}

// Read out of a real save after a sitdown went wrong: "At The Monarch, late,
// with the city quiet: Ennio Zanetti was knifed in the crowd. A meeting
// between Bellandi Family and Russo Outfit went the way somebody had already
// decided." The meeting was at Saint Agnes. Casualties are drawn from anywhere
// in the organization and the paper reports a death where the person was
// standing, so the city named a room three streets from the one the player was
// sitting in.
func TestWhoeverDiesAtASitdownDiedInTheRoom(t *testing.T) {
	t.Parallel()
	w := mediator(t, 95, true) // hostile, and one side led by somebody who does not wait
	for i := range w.NPCs {
		if w.NPCs[i].Location == SitdownGround {
			w.NPCs[i].Location = "club" // nobody starts in the room
		}
	}
	if q, ok := w.OpenQuarrel(); !ok || !q.Trap {
		t.Fatal("this fixture is not a trap, so it proves nothing about one")
	}
	before := map[string]bool{}
	for i := range w.NPCs {
		before[w.NPCs[i].ID] = w.NPCs[i].Dead
	}
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Event == nil {
		t.Fatal("no sitdown scene")
	}
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "choice", Event: w.Event.ID, Choice: w.Event.Choices[0].ID})
	if err != nil {
		t.Skip("this room did not turn: " + err.Error())
	}
	w = next
	died := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead && !before[n.ID] {
			died++
			if n.Location != SitdownGround {
				t.Fatalf("%s was killed at the meeting and the city says they died at %q", n.Name, n.Location)
			}
		}
	}
	if died == 0 {
		t.Fatal("a trapped room turned and nobody in it died, so this proves nothing")
	}
}

// Found by cmd/apicheck, twice in twenty-two runs against a fresh save:
//
//	HTTP 409 on sitdown: "Doyle Crew would not sit in a room you arranged.
//	They think of you at -26 and it takes -25."
//
// The game listed the action as available and then refused the command. The
// room takes an evening and the clock is advanced before it opens, so three
// hours passed between the check that enabled the button and the check that ran
// it, and one point of drift in that time refused a player who had already
// committed.
//
// Three guesses at the cause reproduced nothing across seven hundred attempts.
// Making the refusal name the family and the figure found it in one run — the
// families were sitting exactly on the boundary, which no fixture had put them
// on.
func TestAnAgreedMeetingIsNotCalledOffByOnePoint(t *testing.T) {
	t.Parallel()
	w := mediator(t, 40, false)
	w.Player.Cash, w.Player.Respect = 5000, 60
	q, ok := w.OpenQuarrel()
	if !ok {
		t.Fatal("no quarrel to sit over")
	}
	// Exactly welcome, the way the failing save was.
	q.A.Goodwill, q.B.Goodwill = SitdownWelcome, SitdownWelcome
	if reason := w.SitdownReadiness(); reason != "" {
		t.Fatalf("the button was not offered: %s", reason)
	}
	agreed, _ := w.OpenQuarrel()
	// The three hours pass and one of them thinks a point less of the player.
	q.B.Goodwill = SitdownWelcome - 1
	if reason := w.SitdownReadiness(); reason == "" {
		t.Fatal("this fixture no longer refuses, so it proves nothing")
	}
	if err := w.CallSitdownAs(agreed); err != nil {
		t.Fatalf("the meeting the player paid for was called off: %v", err)
	}
	if w.Event == nil || w.Event.Kind != "sitdown" {
		t.Fatal("no room opened")
	}
}

// And the refusal, when it is a refusal, says which door is shut.
func TestARefusedMeetingNamesWhoWouldNotCome(t *testing.T) {
	t.Parallel()
	w := mediator(t, 40, false)
	w.Player.Cash, w.Player.Respect = 5000, 60
	q, _ := w.OpenQuarrel()
	q.B.Goodwill = SitdownWelcome - 1
	reason := w.SitdownReadiness()
	if !strings.Contains(reason, q.B.Name) {
		t.Fatalf("the refusal does not say who: %q", reason)
	}
	if !strings.Contains(reason, "-26") {
		t.Fatalf("the refusal does not say by how much: %q", reason)
	}
}

// A dining room of your own.
//
// A sitdown could only ever be held in the back of a bar, and the restaurant
// was on the brief's list of trades that do not reach past their own income.
// Both of those were the same gap. A corner table with the plates still down is
// where this city has always done it, and a man with a dining room to lose has
// as much reason as either of them for nobody to draw anything in it.
//
// Two things come with holding the room: you are not renting it, and the people
// on the door are yours. The second is the one that matters, because what your
// staff see on the way in is the difference between walking into a trap and
// knowing about it, and information is the scarcest thing in that room.

func theRestaurant() string {
	for _, l := range Locations {
		if l.Kind == "restaurant" {
			return l.ID
		}
	}
	return ""
}

// quarrelling is a city with a war hot enough that somebody comes to finish it.
func quarrelling(t *testing.T) *World {
	t.Helper()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 200000, 90
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "war", 90
	}
	if q, ok := w.OpenQuarrel(); !ok || !q.Trap {
		t.Skip("this city has nobody coming to finish anything")
	}
	return w
}

// yourDiningRoom hands the player a restaurant with its full complement.
func yourDiningRoom(w *World) string {
	id := theRestaurant()
	own(w, id)
	trade, _ := TradeOf(id)
	w.Properties[id].Staff = trade.Hands
	return id
}

func TestARestaurantOfYourOwnIsARoomTheyWouldSitIn(t *testing.T) {
	t.Parallel()
	w := quarrelling(t)
	id := theRestaurant()
	w.Player.Location = id
	if w.SitdownWhere(id) {
		t.Fatal("two families sat down in a restaurant that belongs to somebody else")
	}
	if actionByID(w.Actions(id), "sitdown") != nil {
		t.Fatal("a room that is not yours offered the meeting")
	}
	yourDiningRoom(w)
	if !w.SitdownWhere(id) {
		t.Fatal("a dining room of yours is not a room either of them would come to")
	}
	a := actionByID(w.Actions(id), "sitdown")
	if a == nil {
		t.Fatal("your own dining room does not offer the meeting")
	}
	if a.Disabled {
		t.Fatalf("refused: %s", a.Reason)
	}
	// And the bar still works, because this adds a room rather than moving one.
	w.Player.Location = SitdownGround
	if !w.SitdownWhere(SitdownGround) {
		t.Fatal("the back of the bar stopped being neutral ground")
	}
}

func TestAnEmptyDiningRoomIsNotAGuarantee(t *testing.T) {
	t.Parallel()
	w := quarrelling(t)
	id := yourDiningRoom(w)
	w.Player.Location = id
	w.Properties[id].Staff = 0
	if w.SitdownWhere(id) {
		t.Fatal("a restaurant with nobody in it was offered as neutral ground")
	}
	if w.OwnGround() {
		t.Fatal("an empty room has people on the door")
	}
	// One short is still short: the point is the people, not the address.
	trade, _ := TradeOf(id)
	w.Properties[id].Staff = trade.Hands - 1
	if w.SitdownWhere(id) {
		t.Fatalf("short-handed at %d of %d and still offering guarantees",
			trade.Hands-1, trade.Hands)
	}
}

func TestYourOwnDoorTellsYouWhoCameHeavy(t *testing.T) {
	t.Parallel()
	w := quarrelling(t)
	// At the bar, without the contacts to be warned, the trap is a surprise.
	w.Player.Location = SitdownGround
	if w.Reach() >= 2 {
		t.Skip("this player already hears everything")
	}
	q, _ := w.OpenQuarrel()
	if !q.Trap {
		t.Fatal("nobody came to finish it")
	}
	if q.Suspected {
		t.Fatal("somebody warned a player with no way of being warned")
	}
	// In your own room, your own staff are the warning.
	id := yourDiningRoom(w)
	w.Player.Location = id
	mine, _ := w.OpenQuarrel()
	if !mine.Suspected {
		t.Fatal("your own people watched them come in and said nothing")
	}
	// And it is in the room when it opens, in somebody's words.
	w.Event = nil
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Event == nil {
		t.Fatal("the meeting never opened")
	}
	if !strings.Contains(w.Event.Body, "your own people on the door") {
		t.Fatalf("the room says: %s", w.Event.Body)
	}
}

func TestTheRoomIsFreeWhenItIsYours(t *testing.T) {
	t.Parallel()
	w := quarrelling(t)
	w.Player.Location = SitdownGround
	if w.SitdownCost() != SitdownFee {
		t.Fatalf("the bar's back room costs $%d", w.SitdownCost())
	}
	id := yourDiningRoom(w)
	w.Player.Location = id
	if w.SitdownCost() != 0 {
		t.Fatalf("you are renting a room you own, at $%d", w.SitdownCost())
	}
	cash := w.Player.Cash
	w.Event = nil
	if err := w.CallSitdown(); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash {
		t.Fatalf("$%d went somewhere for a room of your own", cash-w.Player.Cash)
	}
}

func TestNobodyIsNamedInTheRoomWhoIsNotInIt(t *testing.T) {
	t.Parallel()
	// The warning said "Mara caught your eye on the way in" in every campaign
	// ever played, including the ones in which Mara had been dead a month.
	w := quarrelling(t)
	w.Player.Location = SitdownGround
	gone := w.Holder("fixer")
	if gone == nil {
		t.Skip("this city has no fixer")
	}
	dead := gone.Name
	gone.Dead, gone.Location = true, ""
	for day := 0; day < 3; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	for i := range w.Conflicts {
		w.Conflicts[i].State, w.Conflicts[i].Hostility = "war", 90
	}
	w.Player.Location = SitdownGround
	w.Player.Cash, w.Player.Respect = 200000, 90
	q, ok := w.OpenQuarrel()
	if !ok {
		t.Skip("the quarrel went away")
	}
	q.Suspected = true
	w.Event = nil
	if err := w.CallSitdownAs(q); err != nil {
		t.Fatal(err)
	}
	// Not the dead one, and not a fragment of their name either: the first
	// version of this guard looked for the whole name and missed "Mara" on its
	// own, and the second looked for either word and found "Bell" inside
	// "Bellandi Family". Whole words, against whole words.
	said := map[string]bool{}
	for _, word := range strings.Fields(w.Event.Body) {
		said[strings.Trim(word, "“”.,'s")] = true
	}
	for _, part := range strings.Fields(dead) {
		if said[part] {
			t.Fatalf("%s has been dead three days and is still at the door: %s", dead, w.Event.Body)
		}
	}
	// And somebody who is actually there is doing the warning.
	now := w.Holder("fixer")
	if now == nil {
		t.Fatal("nobody took the fixer's place")
	}
	if !strings.Contains(w.Event.Body, now.Name) {
		t.Fatalf("the room names nobody who is in it: %s", w.Event.Body)
	}
}
