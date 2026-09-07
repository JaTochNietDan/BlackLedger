package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func act(t *testing.T, w **World, kind, target string) {
	t.Helper()
	n, e := Execute(*w, Command{Revision: (*w).Revision, Kind: kind, Target: target})
	if e != nil {
		t.Fatal(e)
	}
	*w = n
}
func choice(t *testing.T, w **World, c string) {
	t.Helper()
	n, e := Execute(*w, Command{Revision: (*w).Revision, Kind: "choice", Event: (*w).Event.ID, Choice: c})
	if e != nil {
		t.Fatal(e)
	}
	*w = n
}
func TestReadingDoesNotAdvance(t *testing.T) {
	w := New(27)
	before, _ := json.Marshal(w)
	for i := 0; i < 10; i++ {
		w.Public()
	}
	after, _ := json.Marshal(w)
	if string(before) != string(after) {
		t.Fatal("reading mutated world")
	}
}
func TestTravel(t *testing.T) {
	w := New(27)
	act(t, &w, "travel", "bar")
	if w.Player.Location != "bar" || w.Minute != 480+TravelMinutes("room", "bar") {
		t.Fatal("bad travel")
	}
}

func TestPresentationRecordsSurviveHistoryRollover(t *testing.T) {
	w := New(27)
	for i := 0; i < 200; i++ {
		w.Log("Old record", "Earlier in the campaign", "personal")
	}
	act(t, &w, "travel", "bar")
	if len(w.History) != 180 || len(w.LastResult.Records) == 0 {
		t.Fatal("history rollover discarded the new command's presentation")
	}
	for _, record := range w.LastResult.Records {
		if record.Title == "Old record" {
			t.Fatal("old history replayed as a new outcome")
		}
	}
	if w.LastResult.Records[len(w.LastResult.Records)-1].ID != w.History[len(w.History)-1].ID {
		t.Fatal("latest outcome missing from presentation")
	}
}
func TestUnavailableActions(t *testing.T) {
	for _, c := range []Command{{Kind: "courier", Target: "bar"}, {Kind: "acquire", Target: "laundry"}, {Kind: "security", Target: "room"}, {Kind: "new_life"}} {
		w := New(27)
		_, e := Execute(w, c)
		if e == nil || w.Revision != 0 || w.Player.Cash != 90 {
			t.Fatal("invalid command accepted or mutated state")
		}
	}
}
func TestEarlyProgression(t *testing.T) {
	w := New(27)
	act(t, &w, "travel", "bar")
	for i := 0; i < 3; i++ {
		act(t, &w, "courier", "bar")
		if w.Event != nil {
			choice(t, &w, "accept")
		}
	}
	act(t, &w, "travel", "laundry")
	act(t, &w, "acquire", "laundry")
	if !w.Own("laundry") {
		t.Fatal("not acquired")
	}
	cash := w.Player.Cash
	act(t, &w, "wait", "laundry")
	if w.Player.Cash <= cash {
		t.Fatal("no income")
	}
}
func TestHiddenHit(t *testing.T) {
	w := New(27)
	act(t, &w, "travel", "club")
	act(t, &w, "provoke", "club")
	if len(w.Plots) != 1 {
		t.Fatal("no plot")
	}
	data, _ := json.Marshal(w.Public())
	if strings.Contains(string(data), w.Plots[0].ID) || strings.Contains(string(data), `"rng"`) || strings.Contains(string(data), `"plots"`) {
		t.Fatal("private state exposed")
	}
}
func TestWarningAndInterruption(t *testing.T) {
	w := New(27)
	w.Player.Contacts = 2
	w.Retaliation()
	w.Advance(151)
	if !w.Plots[0].Known {
		t.Fatal("no warning")
	}
	w.Advance(300)
	if w.Minute != 720 || w.Event == nil {
		t.Fatal("did not pause at attack")
	}
}
func TestAbsentPlayer(t *testing.T) {
	w := New(27)
	w.Player.Location = "bar"
	w.Retaliation()
	w.Advance(300)
	if !w.Player.Alive || w.Properties["room"].Condition != 55 {
		t.Fatal("absent target killed at home")
	}
}
func TestTribute(t *testing.T) {
	w := New(27)
	w.Player.Cash = 300
	act(t, &w, "travel", "club")
	act(t, &w, "provoke", "club")
	act(t, &w, "audience", "club")
	choice(t, &w, "tribute")
	if len(w.Plots) != 0 {
		t.Fatal("plot not canceled")
	}
}
func TestStaleChoice(t *testing.T) {
	w := New(27)
	act(t, &w, "travel", "club")
	act(t, &w, "audience", "club")
	_, e := Execute(w, Command{Revision: w.Revision, Kind: "choice", Event: "stale", Choice: "leave"})
	if e == nil || w.Event == nil {
		t.Fatal("stale decision accepted")
	}
}
func TestSecurityMatters(t *testing.T) {
	for _, g := range []int{0, 3} {
		w := New(50)
		w.Player.Security = g
		w.Player.Contacts = 2
		w.Retaliation()
		w.Advance(240)
		choice(t, &w, "defend")
		if w.Player.Alive != (g == 3) {
			t.Fatalf("security %d outcome wrong", g)
		}
	}
}
func TestDeathAndNewLife(t *testing.T) {
	w := New(27)
	id := w.ID
	w.Properties["laundry"].Owner = "player:1"
	w.Die("Test death")
	act(t, &w, "new_life", "")
	if w.ID != id || w.Life != 2 || w.Own("laundry") || w.Player.Cash != 90 || len(w.Dead) != 1 {
		t.Fatal("legacy reset broken")
	}
}
func TestDelegate(t *testing.T) {
	w := New(27)
	w.Player.Crew = []Crew{{"leo", "Leo", 65}}
	act(t, &w, "delegate", "room")
	before := w.Player.Cash
	w.Advance(130)
	if w.Player.Cash != before+65 {
		t.Fatal("task payment")
	}
	w.Advance(50)
	if w.Player.Cash != before+65 {
		t.Fatal("task paid twice")
	}
}
func TestBills(t *testing.T) {
	w := New(27)
	w.Minute = 1439
	w.Advance(1)
	if w.Player.Cash != 75 {
		t.Fatal("bad rent")
	}
}
func TestValidation(t *testing.T) {
	w := New(27)
	for _, p := range []Proposal{{}, {Speaker: "unknown"}, {Speaker: "mara", Operation: "kill_player"}} {
		if _, e := w.ValidateProposal(p); e == nil {
			t.Fatal("bad proposal accepted")
		}
	}
	e, err := w.ValidateProposal(Proposal{"Sealed envelope", "Would you deliver this?", "mara", "courier", "Secret future sentence.", ""})
	if err != nil {
		t.Fatal(err)
	}
	if e.Effect.Reward != 75 {
		t.Fatal("model owns economics")
	}
	w.Event = e
	b, _ := json.Marshal(w.Public())
	if strings.Contains(string(b), "Secret future sentence") {
		t.Fatal("future exposed")
	}
}
func TestInterruptedRestNoFreeHealing(t *testing.T) {
	w := New(27)
	w.Player.Health = 25
	w.Player.Security = 1
	w.Retaliation()
	w.Plots[0].Due = w.Minute + 15
	act(t, &w, "rest", "room")
	if w.Player.Health != 25 || w.Event == nil || w.Minute != 495 {
		t.Fatal("rest completed despite interruption")
	}
}
func TestInspectNotInfiniteRespect(t *testing.T) {
	w := New(27)
	w.Player.Location = "laundry"
	w.Properties["laundry"].Owner = "player:1"
	act(t, &w, "inspect", "laundry")
	if w.Player.Respect != 0 {
		t.Fatal("inspection farming exploit")
	}
}
func BenchmarkAdvanceCityDay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := New(27)
		w.Properties["laundry"].Owner = "player:1"
		w.Advance(1440)
	}
}

func TestScheduledIncomeExact(t *testing.T) {
	w := New(27)
	w.Properties["laundry"].Owner = "player:1"
	cash := w.Player.Cash
	w.Advance(60)
	if w.Player.Cash != cash+14 {
		t.Fatal("hourly income drift")
	}
}
func TestScheduleStopsBeforePostAttackTask(t *testing.T) {
	w := New(27)
	w.Player.Security = 1
	w.Retaliation()
	w.Tasks = append(w.Tasks, Task{ID(), "Leo", w.Minute + 300})
	w.Advance(600)
	if w.Minute != 720 || len(w.Tasks) != 1 {
		t.Fatal("advanced past player decision")
	}
}

func TestNewPersonDoesNotInheritRelationshipsOrPendingDirector(t *testing.T) {
	w := New(27)
	w.NPCs[0].Trust = 85
	name := w.NPCs[0].Name
	w.Director.Status = "writing"
	w.Die("A previous life ended.")
	act(t, &w, "new_life", "")
	if w.NPCs[0].Trust != 0 || w.NPCs[0].Name != name || w.Director.Status == "writing" {
		t.Fatal("new person inherited trust or an obsolete generation job")
	}
}

func TestDirectorCannotInventMechanicalCompletion(t *testing.T) {
	w := New(27)
	scene, err := w.ValidateProposal(Proposal{"A payment", "Please collect this payment.", "mara", "collection", "The rival is killed and you own his casino.", ""})
	if err != nil {
		t.Fatal(err)
	}
	if scene.Outcome != "You collected the agreed payment and reported back." {
		t.Fatal("unverified outcome became canonical")
	}
	if scene.Choices[0].Label != "Collect the payment" {
		t.Fatal("operation hidden behind vague choice")
	}
}

func TestInjuryReducesAttackSurvival(t *testing.T) {
	survived := func(health int) int {
		n := 0
		for seed := uint32(1); seed < 1000; seed++ {
			w := New(seed * 97)
			w.Player.Health = health
			w.Event = &Scene{ID: "attack", Kind: "attack", Choices: []Choice{{ID: "defend"}}}
			next, err := Execute(w, Command{Kind: "choice", Event: "attack", Choice: "defend", Revision: 0})
			if err != nil {
				t.Fatal(err)
			}
			if next.Player.Alive {
				n++
			}
		}
		return n
	}
	healthy, injured := survived(100), survived(25)
	if healthy <= injured || injured == 0 {
		t.Fatalf("health has no meaningful bounded effect: healthy %d injured %d", healthy, injured)
	}
}
