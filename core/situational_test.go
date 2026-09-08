package core

import (
	"strings"
	"testing"
)

func TestAQuietCityOffersNoConflictWork(t *testing.T) {
	w := New(101)
	if len(w.SituationalOperations()) != 0 {
		t.Fatalf("a city with no war, no seizure and no succession offered %v", w.SituationalOperations())
	}
	// And the director keeps to the ordinary rotation.
	op := w.NextDirectorOperation()
	if _, situational := SituationalEffects()[op]; situational {
		t.Fatalf("a quiet city selected conflict work: %s", op)
	}
}

func TestWarCreatesWorkThatOnlyAWarCouldCreate(t *testing.T) {
	w := New(103)
	w.Antagonize("bellandi", "russo", 100)
	if w.Conflict("bellandi", "russo").State != "war" {
		t.Fatal("the test world is not at war")
	}
	ids := map[string]string{}
	for _, s := range w.SituationalOperations() {
		ids[s.ID] = s.Because
	}
	for _, want := range []string{"escort", "warning"} {
		because, ok := ids[want]
		if !ok {
			t.Fatalf("a war produced no %s work: %v", want, ids)
		}
		if !strings.Contains(because, "Bellandi") || !strings.Contains(because, "Russo") {
			t.Fatalf("%s does not say which war it exists because of: %q", want, because)
		}
	}
	// Every situational operation pays something the catalog knows about.
	for id := range ids {
		if _, ok := SituationalEffects()[id]; !ok {
			t.Fatalf("%s has no reward table", id)
		}
	}
}

func TestGroundChangingHandsCreatesSomethingToRecover(t *testing.T) {
	w := New(107)
	place, _ := PlaceByID("club")
	w.Report("seizure", upper(place.Name)+" CHANGES HANDS", "It changed hands.")
	w.Properties["club"].Owner = "russo"
	found := ""
	for _, s := range w.SituationalOperations() {
		if s.ID == "recovery" {
			found = s.Because
		}
	}
	if found == "" {
		t.Fatal("a seizure created nothing to recover")
	}
	if !strings.Contains(found, "The Monarch") || !strings.Contains(found, "Russo") {
		t.Fatalf("the recovery does not name the place or who holds it now: %q", found)
	}
	// It stops being current once it is old news.
	w.Minute += 10000
	for _, s := range w.SituationalOperations() {
		if s.ID == "recovery" {
			t.Fatal("a seizure from a week ago is still generating work")
		}
	}
}

func TestANewLeaderLeavesArrangementsUnsettled(t *testing.T) {
	w := New(109)
	deputy := w.Members("bellandi")[1].Name
	w.Kill("vittorio", "Shot at the club.")
	if w.faction("bellandi").Leader != deputy {
		t.Fatal("succession did not happen")
	}
	found := ""
	for _, s := range w.SituationalOperations() {
		if s.ID == "settlement" {
			found = s.Because
		}
	}
	if found == "" {
		t.Fatal("a change at the top of a family created no work")
	}
	if !strings.Contains(found, deputy) {
		t.Fatalf("the settlement does not name who now leads: %q", found)
	}
}

func TestConflictWorkDoesNotCrowdOutEverythingElse(t *testing.T) {
	w := New(113)
	w.Antagonize("bellandi", "russo", 100)
	seen := map[string]int{}
	for i := 0; i < 12; i++ {
		op := w.NextDirectorOperation()
		seen[op]++
		// Record it as told, the way a completed arrangement would.
		w.Arrangements = append(w.Arrangements, ArrangementMemory{
			ID: ID(), Life: w.Life, Operation: op, Status: "completed", Speaker: "mara",
		})
	}
	if len(seen) < 3 {
		t.Fatalf("a city at war told only %d kinds of story: %v", len(seen), seen)
	}
	situational, ordinary := 0, 0
	for op, n := range seen {
		if _, is := SituationalEffects()[op]; is {
			situational += n
		} else {
			ordinary += n
		}
	}
	if situational == 0 {
		t.Fatal("a city at war never told a war story")
	}
	if ordinary == 0 {
		t.Fatal("a war stopped every ordinary job in the city")
	}
	t.Logf("over 12 offers in a city at war: %d conflict jobs, %d ordinary, kinds %v", situational, ordinary, seen)
}

func TestEveryOperationHasWordsForItsButtonAndItsLedger(t *testing.T) {
	// The conflict-derived operations were added without labels, which rendered
	// as an unpressable blank choice in a live offer.
	operations := []string{"courier", "mediation", "collection"}
	for id := range SituationalEffects() {
		operations = append(operations, id)
	}
	for _, op := range operations {
		if label := operationLabel(op); len(strings.TrimSpace(label)) < 4 {
			t.Fatalf("%s has no usable accept label: %q", op, label)
		}
		if outcome := operationOutcome(op); len(strings.TrimSpace(outcome)) < 10 {
			t.Fatalf("%s has no outcome for the ledger: %q", op, outcome)
		}
	}
	// An unknown operation still produces something pressable rather than blank.
	if len(strings.TrimSpace(operationLabel("nonsense"))) < 4 {
		t.Fatal("an unknown operation rendered a blank button")
	}
}

func TestABusinessInTroubleIsAReasonForWork(t *testing.T) {
	w := New(811)
	w.Properties["laundry"].Owner = "player:1"
	if w.hasSituational("supply") {
		t.Fatal("a business running properly generated work for itself")
	}
	w.Properties["laundry"].Supply = 0
	if !w.hasSituational("supply") {
		t.Fatal("a business out of supplies generated no work")
	}
	because := w.becauseOf("supply")
	if !strings.Contains(because, "Bluebird Laundry") || !strings.Contains(because, "soap") {
		t.Fatalf("the reason does not say what is wrong or where: %q", because)
	}
	// A business somebody else owns is not the player's problem to solve.
	other := New(813)
	other.Properties["laundry"].Supply = 0
	if other.hasSituational("supply") {
		t.Fatal("a business the player does not own generated work for them")
	}
}

func TestStockThatHasToMoveIsAReasonForWork(t *testing.T) {
	w := New(817)
	if w.hasSituational("distribution") {
		t.Fatal("an empty pocket generated distribution work")
	}
	w.Player.Stock = map[string]int{"moonshine": 20}
	if !w.hasSituational("distribution") {
		t.Fatal("twenty crates generated no reason to move anything")
	}
	if !strings.Contains(w.becauseOf("distribution"), "20") {
		t.Fatalf("the reason does not say how much: %q", w.becauseOf("distribution"))
	}
}

func (w *World) hasSituational(id string) bool {
	for _, s := range w.SituationalOperations() {
		if s.ID == id {
			return true
		}
	}
	return false
}

func (w *World) becauseOf(id string) string {
	for _, s := range w.SituationalOperations() {
		if s.ID == id {
			return s.Because
		}
	}
	return ""
}
