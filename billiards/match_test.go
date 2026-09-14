package billiards

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func testMatch() *Match {
	m, _ := NewMatch(0)
	m.Breaking = false
	m.InHand = false
	m.HeadOnly = false
	return m
}
func contactEvent(first int) Event  { return Event{Time: .1, Kind: "ball", Ball: 0, Other: first} }
func potEvent(id, pocket int) Event { return Event{Time: .3, Kind: "pocket", Ball: id, Other: pocket} }
func resultFor(m *Match, events ...Event) Result {
	r := Result{Balls: append([]Ball(nil), m.Balls...), Events: events}
	for _, ev := range events {
		if ev.Kind == "pocket" {
			for i := range r.Balls {
				if r.Balls[i].ID == ev.Ball {
					r.Balls[i].Pocketed = true
					r.Balls[i].Pocket = ev.Other
				}
			}
		}
	}
	return r
}
func judge(t *testing.T, m *Match, call Call, events ...Event) Outcome {
	t.Helper()
	out, err := m.adjudicate(New(), resultFor(m, events...), call)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestMatchOnlyCalledBallAssignsGroups(t *testing.T) {
	for _, test := range []struct {
		name      string
		call      Call
		events    []Event
		wantGroup int
		keep      bool
	}{
		{"called solid", Call{Ball: 1, Pocket: 2}, []Event{contactEvent(9), potEvent(1, 2)}, 1, true},
		{"called stripe despite solid also falling", Call{Ball: 9, Pocket: 2}, []Event{contactEvent(1), potEvent(1, 0), potEvent(9, 2)}, 2, true},
		{"wrong pocket", Call{Ball: 1, Pocket: 2}, []Event{contactEvent(1), potEvent(1, 3)}, 0, false},
		{"uncalled ball", Call{Ball: 1, Pocket: 2}, []Event{contactEvent(1), potEvent(9, 2)}, 0, false},
		{"safety with pocketed ball", Call{Safety: true}, []Event{contactEvent(1), potEvent(1, 2)}, 0, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := testMatch()
			out := judge(t, m, test.call, test.events...)
			if out.Foul || m.Groups[0] != test.wantGroup || out.KeptTurn != test.keep {
				t.Fatal(m.Groups, out)
			}
			if test.keep && m.Turn != 0 || !test.keep && m.Turn != 1 {
				t.Fatal("wrong next player")
			}
		})
	}
}
func TestMatchFoulsGrantBallInHandWithoutAssigningOpenGroups(t *testing.T) {
	for _, events := range [][]Event{{}, {contactEvent(8), potEvent(1, 0)}, {contactEvent(1)}, {contactEvent(1), potEvent(1, 0), potEvent(0, 1)}} {
		m := testMatch()
		out := judge(t, m, Call{Ball: 1, Pocket: 0}, events...)
		if !out.Foul || !m.InHand || m.Turn != 1 || m.Groups != [2]int{} || m.Winner != -1 {
			t.Fatal(m, out)
		}
	}
	m := testMatch()
	m.Groups = [2]int{1, 2}
	out := judge(t, m, Call{Ball: 1, Pocket: 0}, contactEvent(9), potEvent(1, 0))
	if !out.Foul || !m.InHand {
		t.Fatal("wrong group wasn't a foul")
	}
}
func TestRailBeforeContactDoesNotAvoidFoul(t *testing.T) {
	m := testMatch()
	out := judge(t, m, Call{Safety: true}, Event{Time: .01, Kind: "rail", Ball: 0, Other: 0}, contactEvent(1))
	if !out.Foul {
		t.Fatal("pre-contact rail was counted")
	}
	m = testMatch()
	out = judge(t, m, Call{Safety: true}, contactEvent(1), Event{Time: .2, Kind: "rail", Ball: 1, Other: 0})
	if out.Foul {
		t.Fatal("post-contact object rail was not counted")
	}
}
func TestFrozenRailNeedsDepartureBeforeItCounts(t *testing.T) {
	m := testMatch()
	m.ball(1).Position = Vec{.6, Radius, Radius}
	impact := Event{Time: .2, Kind: "rail", Ball: 1, Other: 0}
	r := resultFor(m, contactEvent(1), impact)
	out, err := m.adjudicate(New(), r, Call{Safety: true})
	if err != nil || !out.Foul {
		t.Fatal("frozen rail counted", out, err)
	}
	m = testMatch()
	m.ball(1).Position = Vec{.6, Radius, Radius}
	r = resultFor(m, contactEvent(1), impact)
	away := append([]Ball(nil), m.Balls...)
	for i := range away {
		if away[i].ID == 1 {
			away[i].Position.Y = .2
		}
	}
	r.Frames = []Frame{{Time: .15, Balls: away}}
	out, err = m.adjudicate(New(), r, Call{Safety: true})
	if err != nil || out.Foul {
		t.Fatal("return to rail did not count", out, err)
	}
}
func TestEightBallVictoryRequiresGroupClearedBeforeShot(t *testing.T) {
	for _, test := range []struct {
		name   string
		clear  bool
		call   Call
		events []Event
		winner int
	}{
		{"called win", true, Call{Ball: 8, Pocket: 2}, []Event{contactEvent(8), potEvent(8, 2)}, 0},
		{"wrong pocket", true, Call{Ball: 8, Pocket: 1}, []Event{contactEvent(8), potEvent(8, 2)}, 1},
		{"scratch while winning", true, Call{Ball: 8, Pocket: 2}, []Event{contactEvent(8), potEvent(8, 2), potEvent(0, 0)}, 1},
		{"safety eight", true, Call{Safety: true}, []Event{contactEvent(8), potEvent(8, 2)}, 1},
		{"early eight", false, Call{Ball: 1, Pocket: 0}, []Event{contactEvent(1), potEvent(8, 2)}, 1},
		{"last group ball and eight same shot", false, Call{Ball: 1, Pocket: 0}, []Event{contactEvent(1), potEvent(1, 0), potEvent(8, 2)}, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := testMatch()
			m.Groups = [2]int{1, 2}
			for i := range m.Balls {
				if group(m.Balls[i].ID) == 1 {
					m.Balls[i].Pocketed = true
				}
			}
			if !test.clear {
				m.ball(1).Pocketed = false
			}
			out := judge(t, m, test.call, test.events...)
			if m.Winner != test.winner || out.Winner != test.winner {
				t.Fatal(m.Winner, out)
			}
		})
	}
}
func TestScratchWhileOnEightDoesNotLoseWithoutPocketingEight(t *testing.T) {
	m := testMatch()
	m.Groups = [2]int{1, 2}
	for i := range m.Balls {
		if group(m.Balls[i].ID) == 1 {
			m.Balls[i].Pocketed = true
		}
	}
	out := judge(t, m, Call{Ball: 8, Pocket: 0}, contactEvent(8), potEvent(0, 1))
	if !out.Foul || m.Winner != -1 || !m.InHand {
		t.Fatal(out)
	}
}
func TestOpenTableCanClaimAlreadyClearedGroupToPlayEight(t *testing.T) {
	m := testMatch()
	for i := range m.Balls {
		if group(m.Balls[i].ID) == 1 {
			m.Balls[i].Pocketed = true
		}
	}
	out := judge(t, m, Call{Ball: 8, Pocket: 0}, contactEvent(8), potEvent(8, 0))
	if out.Winner != 0 {
		t.Fatal(out)
	}
}
func TestBreakPocketKeepsOpenTable(t *testing.T) {
	m, _ := NewMatch(0)
	m.InHand = false
	m.HeadOnly = false
	out := judge(t, m, Call{}, contactEvent(1), potEvent(9, 0))
	if out.Foul || !out.KeptTurn || m.Groups != [2]int{} || m.Breaking {
		t.Fatal(out, m)
	}
}
func TestBreakRequiresFourDifferentObjectBallsAtRails(t *testing.T) {
	for _, distinct := range []bool{false, true} {
		m, _ := NewMatch(0)
		m.InHand = false
		m.HeadOnly = false
		events := []Event{contactEvent(1)}
		for i := 0; i < 4; i++ {
			id := 1
			if distinct {
				id = i + 1
			}
			events = append(events, Event{Time: .2 + float64(i)/10, Kind: "rail", Ball: id, Other: 0})
		}
		judge(t, m, Call{}, events...)
		if distinct && m.Pending != nil || !distinct && (m.Pending == nil || m.Pending.Kind != "illegal") {
			t.Fatal("break rail count incorrect", distinct, m.Pending)
		}
	}
}
func TestIllegalBreakChoicesAndReracks(t *testing.T) {
	for _, choice := range []string{"accept", "rerack-self", "rerack-opponent"} {
		m, _ := NewMatch(0)
		m.InHand = false
		m.HeadOnly = false
		judge(t, m, Call{}, contactEvent(1))
		if err := m.DecideBreak(New(), 0, choice); err == nil {
			t.Fatal("breaker stole incoming player's decision")
		}
		if err := m.DecideBreak(New(), 1, choice); err != nil {
			t.Fatal(err)
		}
		if m.Pending != nil {
			t.Fatal("decision remains pending")
		}
		if choice == "accept" {
			if m.Breaking || m.Turn != 1 {
				t.Fatal(m)
			}
		} else if !m.Breaking || !m.InHand || !m.HeadOnly {
			t.Fatal("rerack not ready", m)
		}
		if choice == "rerack-opponent" && m.Turn != 0 {
			t.Fatal("wrong breaker")
		}
	}
}
func TestEightOnBreakOffersSpotOrRerackToCorrectPlayer(t *testing.T) {
	for _, scratch := range []bool{false, true} {
		m, _ := NewMatch(0)
		m.InHand = false
		m.HeadOnly = false
		events := []Event{contactEvent(1), potEvent(8, 0)}
		if scratch {
			events = append(events, potEvent(0, 1))
		}
		judge(t, m, Call{}, events...)
		seat := 0
		if scratch {
			seat = 1
		}
		if m.Winner != -1 || m.Pending == nil || m.Turn != seat {
			t.Fatal("eight on break incorrectly settled", m)
		}
		if err := m.DecideBreak(New(), seat, "spot-eight"); err != nil {
			t.Fatal(err)
		}
		if m.ball(8).Pocketed || m.Breaking || m.InHand != scratch || m.Groups != [2]int{} {
			t.Fatal("bad spot decision", m)
		}
		if err := New().validate(m.Balls); err != nil {
			t.Fatal("spotting overlaps other balls", err)
		}
	}
}
func TestCuePlacementRestrictionsAndRejectedActionsAreAtomic(t *testing.T) {
	m, _ := NewMatch(0)
	before, _ := json.Marshal(m)
	e := New()
	if err := m.PlaceCue(e, 0, .6, HeadString+.1); err == nil {
		t.Fatal("accepted forward break placement")
	}
	if _, err := m.Play(e, 0, Shot{Speed: 1}, Call{}); err == nil {
		t.Fatal("shot before placement")
	}
	after, _ := json.Marshal(m)
	if string(before) != string(after) {
		t.Fatal("invalid input changed match")
	}
	if err := m.PlaceCue(e, 0, .6, .4); err != nil {
		t.Fatal(err)
	}
	before, _ = json.Marshal(m)
	if _, err := m.Play(e, 0, Shot{Speed: math.NaN()}, Call{}); err == nil {
		t.Fatal("accepted NaN")
	}
	after, _ = json.Marshal(m)
	if string(before) != string(after) {
		t.Fatal("failed physics changed match")
	}
}
func TestRealPhysicalCalledPocketAndSaveRoundTrip(t *testing.T) {
	m := testMatch()
	m.Groups = [2]int{1, 2}
	for i := range m.Balls {
		m.Balls[i].Pocketed = true
	}
	*m.ball(0) = ball(0, .34, Length/2)
	*m.ball(8) = ball(8, .14, Length/2)
	*m.ball(9) = ball(9, .9, 2)
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var restored Match
	if err = json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	result, err := restored.Play(New(), 0, Shot{Angle: math.Pi, Speed: .6}, Call{Ball: 8, Pocket: 4})
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome.Winner != 0 || !restored.ball(8).Pocketed {
		t.Fatal("actual pocket did not win", result.Outcome)
	}
	if m.Winner != -1 {
		t.Fatal("restored match aliased source")
	}
	if err := restored.Concede(0); err == nil {
		t.Fatal("finished match could be conceded again")
	}
}
func TestConcessionAndDecisionsSurviveSerialization(t *testing.T) {
	m, _ := NewMatch(1)
	m.InHand = false
	m.HeadOnly = false
	judge(t, m, Call{}, contactEvent(1))
	raw, _ := json.Marshal(m)
	var restored Match
	if err := json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m.BreakChoices(), restored.BreakChoices()) {
		t.Fatal("lost pending choices")
	}
	if err := restored.Concede(0); err != nil {
		t.Fatal(err)
	}
	if restored.Winner != 1 || restored.Pending != nil {
		t.Fatal("concession did not settle")
	}
}

func TestHeadStringRestrictionUsesActualPathBeforeFirstContact(t *testing.T) {
	for _, crossed := range []bool{false, true} {
		m := testMatch()
		m.HeadOnly = true
		m.ball(1).Position.Y = .3
		r := resultFor(m, contactEvent(1), potEvent(1, 0))
		cue := *m.ball(0)
		cue.Position.Y = .4
		if crossed {
			cue.Position.Y = .8
		}
		r.Frames = []Frame{{Time: .05, Balls: []Ball{cue}}}
		out, err := m.adjudicate(New(), r, Call{Ball: 1, Pocket: 0})
		if err != nil || out.Foul == crossed {
			t.Fatal("head string restriction did not follow shot path", crossed, out, err)
		}
	}
}
func TestSimultaneousFirstContactDoesNotPenalizeLegalBall(t *testing.T) {
	m := testMatch()
	m.Groups = [2]int{1, 2}
	out := judge(t, m, Call{Ball: 1, Pocket: 0}, contactEvent(9), contactEvent(1), potEvent(1, 0))
	if out.Foul {
		t.Fatal("simultaneous legal contact was ignored")
	}
}
func TestAllTargetsBehindHeadStringAreSpottedForBallInHand(t *testing.T) {
	m := testMatch()
	for i := range m.Balls {
		m.Balls[i].Pocketed = true
	}
	*m.ball(1) = ball(1, .3, .2)
	*m.ball(9) = ball(9, .9, .5)
	m.Pending = &BreakDecision{Kind: "foul", Breaker: 1, Foul: true}
	m.Turn = 0
	if err := m.DecideBreak(New(), 0, "hand-behind"); err != nil {
		t.Fatal(err)
	}
	near(t, m.ball(9).Position.Y, Length*.75, 1e-9)
	if !m.HeadOnly || !m.InHand || m.ball(1).Position.Y != .2 {
		t.Fatal("wrong head-string remedy", m)
	}
}
func TestRealPhysicsBreakReachesPlayableMatchState(t *testing.T) {
	m, _ := NewMatch(0)
	e := New()
	if err := m.PlaceCue(e, 0, Width/2, HeadString-.04); err != nil {
		t.Fatal(err)
	}
	result, err := m.Play(e, 0, Shot{Angle: math.Pi / 2, Speed: 7}, Call{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Physics.Events) < 10 || m.Shots != 1 || m.Winner != -1 {
		t.Fatal("break did not become a match", m.Last)
	}
	if m.Pending != nil {
		choices := m.BreakChoices()
		if len(choices) == 0 {
			t.Fatal("break has no way forward")
		}
		if err := m.DecideBreak(e, m.Turn, choices[0]); err != nil {
			t.Fatal(err)
		}
	}
	if m.Breaking || m.Groups != [2]int{} {
		t.Fatal("opening break did not leave open table")
	}
}
