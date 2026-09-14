package billiards

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func botEightPosition() *Match {
	m := testMatch()
	m.Groups = [2]int{1, 2}
	for i := range m.Balls {
		m.Balls[i].Pocketed = true
	}
	*m.ball(0) = ball(0, .34, Length/2)
	*m.ball(8) = ball(8, .14, Length/2)
	*m.ball(9) = ball(9, .9, 2)
	return m
}
func TestOpponentPotsCalledEightThroughPhysics(t *testing.T) {
	m := botEightPosition()
	before, _ := json.Marshal(m)
	turn, err := Opponent(New(), m, 0, 1, 17)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Match.Winner != 0 || turn.Intent.Call.Ball != 8 || len(turn.Result.Physics.Events) == 0 {
		t.Fatal("opponent did not play physical winning shot", turn.Intent, turn.Result.Outcome)
	}
	after, _ := json.Marshal(m)
	if string(before) != string(after) {
		t.Fatal("planning mutated original match")
	}
	// Re-execute precisely the advertised stroke; the result must be identical.
	copy := *m
	copy.Balls = append([]Ball(nil), m.Balls...)
	result, err := copy.Play(New(), 0, turn.Intent.Shot, turn.Intent.Call)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, turn.Result) || !reflect.DeepEqual(copy, turn.Match) {
		t.Fatal("bot result does not follow its cue input")
	}
}
func TestOpponentPlanningIsDeterministicAcrossRestore(t *testing.T) {
	m := botEightPosition()
	raw, _ := json.Marshal(m)
	var saved Match
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatal(err)
	}
	a, err := Opponent(New(), m, 0, .72, 1234)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Opponent(New(), &saved, 0, .72, 1234)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("retry/restore changed shot")
	}
	c, err := Opponent(New(), m, 0, .72, 4321)
	if err != nil {
		t.Fatal(err)
	}
	if a.Intent.Shot == c.Intent.Shot {
		t.Fatal("execution accuracy ignores its seed")
	}
}
func TestOpponentPlacesCueAndResolvesBreakChoices(t *testing.T) {
	m, _ := NewMatch(1)
	turn, err := Opponent(New(), m, 1, .9, 77)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Placement == nil || turn.Placement.Y >= HeadString || turn.Match.Shots != 1 || len(turn.Result.Physics.Events) < 10 {
		t.Fatal("opponent did not place and break")
	}
	m = botEightPosition()
	m.Pending = &BreakDecision{Kind: "eight", Breaker: 0}
	m.ball(8).Pocketed = true
	turn, err = Opponent(New(), m, 0, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Decision != "spot-eight" {
		t.Fatal("eight repeatedly reracked", turn.Decision)
	}
}
func TestOpponentNeverActsForWrongSeatOrFinishedRack(t *testing.T) {
	m := botEightPosition()
	for _, skill := range []float64{-1, 2, math.NaN()} {
		if _, err := Opponent(New(), m, 0, skill, 1); err == nil {
			t.Fatal("invalid skill accepted")
		}
	}
	if _, err := Opponent(New(), m, 1, 1, 1); err == nil {
		t.Fatal("wrong player acted")
	}
	m.Winner = 0
	if _, err := Opponent(New(), m, 0, 1, 1); err == nil {
		t.Fatal("finished rack replayed")
	}
}
func TestOpponentEscapesBlockedStraightLineThroughACushion(t *testing.T) {
	m := testMatch()
	m.Groups = [2]int{1, 2}
	for i := range m.Balls {
		m.Balls[i].Pocketed = true
	}
	*m.ball(0) = ball(0, .8, .8)
	*m.ball(1) = ball(1, .6, 1.3)
	*m.ball(8) = ball(8, .9, 2)
	*m.ball(9) = ball(9, .7, 1.05)
	if clearLine(m.Balls, m.ball(0).Position, m.ball(1).Position, 1) {
		t.Fatal("planner ignored blocker")
	}

	turn, err := Opponent(New(), m, 0, 1, 12)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Result.Outcome.Foul {
		t.Fatal("planner selected blocked foul despite alternatives", turn.Intent, turn.Result.Outcome)
	}
}
func TestPhysicalOpponentsCanFinishCompleteRacks(t *testing.T) {
	legalWins := 0
	for _, seed := range []uint64{7, 41} {
		m, _ := NewMatch(0)
		for shots := 0; shots < 160 && m.Winner < 0; shots++ {
			turn, err := Opponent(New(), m, m.Turn, .9, seed+uint64(shots)*7919)
			if err != nil {
				t.Fatalf("seed %d shot %d: %v", seed, shots, err)
			}
			m = &turn.Match
			if shots%20 == 19 {
				t.Logf("seed %d shot %d groups %v left %d/%d last %s", seed, shots+1, m.Groups, m.Remaining(1), m.Remaining(2), m.Last.Reason)
			}
		}
		if m.Winner < 0 {
			t.Fatalf("seed %d did not finish after %d physical shots", seed, m.Shots)
		}
		if m.Last.Shooter == m.Winner {
			legalWins++
		}
		t.Logf("seed %d: winner %d after %d shots; %s", seed, m.Winner, m.Shots, m.Last.Reason)
	}
	if legalWins == 0 {
		t.Fatal("racks ended only through early-eight losses")
	}
}

func TestOpponentBreakDoesNotDependOnRepeatedReracks(t *testing.T) {
	legal := 0
	for seed := uint64(1); seed <= 8; seed++ {
		m, _ := NewMatch(0)
		turn, err := Opponent(New(), m, 0, .9, seed)
		if err != nil {
			t.Fatal(err)
		}
		rails := map[int]int{}
		pots := []int{}
		for _, ev := range turn.Result.Physics.Events {
			if ev.Kind == "rail" {
				rails[ev.Ball]++
			}
			if ev.Kind == "pocket" {
				pots = append(pots, ev.Ball)
			}
		}
		t.Logf("seed %d shot %+v rails %v pots %v reason %s", seed, turn.Intent.Shot, rails, pots, turn.Match.Last.Reason)
		if turn.Match.Pending == nil || turn.Match.Pending.Kind != "illegal" {
			legal++
		}
	}
	if legal < 6 {
		t.Fatalf("only %d of 8 first breaks were legal", legal)
	}
}

func BenchmarkOpponentMidRack(b *testing.B) {
	m, _ := NewMatch(0)
	opening, err := Opponent(New(), m, 0, 1, 7)
	if err != nil {
		b.Fatal(err)
	}
	m = &opening.Match
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Opponent(New(), m, m.Turn, .9, 41); err != nil {
			b.Fatal(err)
		}
	}
}
