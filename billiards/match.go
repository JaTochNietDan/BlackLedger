package billiards

import (
	"errors"
	"fmt"
	"math"
)

const HeadString = Length / 4

type Call struct {
	Ball, Pocket int
	Safety       bool
}
type BreakDecision struct {
	Kind    string
	Breaker int
	Foul    bool
}
type Match struct {
	Balls                      []Ball
	Turn                       int
	Groups                     [2]int // 0 open, 1 solids, 2 stripes
	Breaking, InHand, HeadOnly bool
	Winner                     int // -1 until finished
	Shots                      int
	Pending                    *BreakDecision
	Last                       Outcome
}
type Outcome struct {
	Shooter        int
	Foul, KeptTurn bool
	Winner         int
	Reason         string
	Pocketed       []int
}
type TurnResult struct {
	Physics Result
	Outcome Outcome
}

func NewMatch(breaker int) (*Match, error) {
	if breaker < 0 || breaker > 1 {
		return nil, errors.New("invalid breaker")
	}
	return &Match{Balls: Rack(), Turn: breaker, Breaking: true, InHand: true, HeadOnly: true, Winner: -1, Last: Outcome{Winner: -1}}, nil
}
func group(id int) int {
	if id >= 1 && id <= 7 {
		return 1
	}
	if id >= 9 && id <= 15 {
		return 2
	}
	return 0
}
func (m *Match) Remaining(g int) int {
	n := 0
	for _, b := range m.Balls {
		if !b.Pocketed && group(b.ID) == g {
			n++
		}
	}
	return n
}
func (m *Match) ball(id int) *Ball {
	for i := range m.Balls {
		if m.Balls[i].ID == id {
			return &m.Balls[i]
		}
	}
	return nil
}
func (m *Match) LegalFirst(seat, id int) bool {
	if id <= 0 || id > 15 || seat < 0 || seat > 1 {
		return false
	}
	if m.Breaking {
		return true
	}
	g := m.Groups[seat]
	if g == 0 {
		return id != 8 || m.Remaining(1) == 0 || m.Remaining(2) == 0
	}
	if m.Remaining(g) == 0 {
		return id == 8
	}
	return group(id) == g
}
func (m *Match) checkSeat(seat int) error {
	if seat < 0 || seat > 1 || seat != m.Turn {
		return errors.New("it is not your turn")
	}
	if m.Winner != -1 {
		return errors.New("the rack is finished")
	}
	return nil
}

// PlaceCue changes only a ball-in-hand position. HeadOnly remains in force
// until the shot, including the restriction on first contact behind the line.
func (m *Match) PlaceCue(e *Engine, seat int, x, y float64) error {
	if err := m.checkSeat(seat); err != nil {
		return err
	}
	if m.Pending != nil {
		return errors.New("resolve the break decision first")
	}
	if !m.InHand {
		return errors.New("the cue ball is not in hand")
	}
	if m.HeadOnly && y >= HeadString {
		return errors.New("place the cue ball behind the head string")
	}
	if err := e.ValidatePlacement(m.Balls, x, y); err != nil {
		return err
	}
	b := m.ball(0)
	if b == nil {
		return errors.New("cue ball is missing")
	}
	*b = Ball{ID: 0, Position: Vec{x, y, Radius}, Orientation: [4]float64{0, 0, 0, 1}}
	m.InHand = false
	return nil
}

// Play is the only shot entry point. The caller supplies aim and a declared
// ball/pocket, never pocket results or the winner. Rejecting input or a physics
// failure leaves the match untouched; money settlement belongs to the campaign.
func (m *Match) Play(e *Engine, seat int, shot Shot, call Call) (TurnResult, error) {
	if err := m.checkSeat(seat); err != nil {
		return TurnResult{}, err
	}
	if m.Pending != nil {
		return TurnResult{}, errors.New("resolve the break decision first")
	}
	if m.InHand {
		return TurnResult{}, errors.New("place the cue ball before shooting")
	}
	if !m.Breaking && !call.Safety {
		b := m.ball(call.Ball)
		if b == nil || b.Pocketed || !m.LegalFirst(seat, call.Ball) || call.Pocket < 0 || call.Pocket >= 6 {
			return TurnResult{}, errors.New("call a legal ball and pocket, or a safety")
		}
	}
	result, err := e.Shoot(m.Balls, shot)
	if err != nil {
		return TurnResult{}, err
	}
	// Judge on a private copy so spotting failures cannot half-commit a shot.
	next := *m
	next.Balls = append([]Ball(nil), m.Balls...)
	out, err := next.adjudicate(e, result, call)
	if err != nil {
		return TurnResult{}, err
	}
	*m = next
	return TurnResult{result, out}, nil
}

func (m *Match) adjudicate(e *Engine, result Result, call Call) (Outcome, error) {
	shooter := m.Turn
	out := Outcome{Shooter: shooter, Winner: -1}
	first := -1
	firstTime := math.Inf(1)
	railAfter := false
	rails := map[int]bool{}
	pots := map[int]int{}
	for _, event := range result.Events {
		if event.Kind == "ball" && first < 0 && (event.Ball == 0 || event.Other == 0) {
			first = event.Other
			if first == 0 {
				first = event.Ball
			}
			firstTime = event.Time
		}
		if event.Kind == "ball" && first >= 0 && math.Abs(event.Time-firstTime) < 1e-9 && (event.Ball == 0 || event.Other == 0) {
			id := event.Other
			if id == 0 {
				id = event.Ball
			}
			if !m.LegalFirst(shooter, first) && m.LegalFirst(shooter, id) {
				first = id
			}
		}
		if event.Kind == "rail" && event.Time >= firstTime && m.newRailContact(e, result, event) {
			railAfter = true
			if event.Ball != 0 {
				rails[event.Ball] = true
			}
		}
		if event.Kind == "pocket" {
			pots[event.Ball] = event.Other
			out.Pocketed = append(out.Pocketed, event.Ball)
		}
	}
	_, scratch := pots[0]
	eightPocket, eight := pots[8]
	objectPot := false
	for id := range pots {
		if id != 0 {
			objectPot = true
		}
	}
	foul := ""
	if first < 0 {
		foul = "The cue ball did not contact an object ball."
	} else if !m.LegalFirst(shooter, first) {
		foul = "The cue ball struck the wrong group first."
	}
	if scratch {
		foul = "The cue ball was pocketed."
	}
	if m.HeadOnly && first >= 0 {
		crossed := false
		for _, f := range result.Frames {
			if f.Time > firstTime+1e-9 {
				break
			}
			for _, b := range f.Balls {
				if b.ID == 0 && b.Position.Y >= HeadString {
					crossed = true
				}
			}
		}
		firstBall := m.ball(first)
		if !crossed && firstBall != nil && firstBall.Position.Y < HeadString {
			foul = "The cue ball hit a ball behind the head string before crossing it."
		}
	}
	if !m.Breaking && first >= 0 && !objectPot && !railAfter {
		foul = "No ball reached a cushion or pocket after contact."
	}
	wasBreak := m.Breaking
	eligibleEight := m.LegalFirst(shooter, 8)
	m.Balls = append([]Ball(nil), result.Balls...)
	m.Shots++
	m.Breaking = false
	m.InHand = false
	m.HeadOnly = false
	out.Foul = foul != ""
	out.Reason = foul
	if wasBreak {
		switch {
		case eight:
			if out.Foul {
				m.Turn = 1 - shooter
			}
			m.Pending = &BreakDecision{Kind: "eight", Breaker: shooter, Foul: out.Foul}
			out.Reason = "The eight fell on the break. Choose spotting it or a new rack."
			if out.Foul {
				out.Reason = "The eight fell on a foul break. The other player chooses spotting it with ball in hand or a new rack."
			}
		case !objectPot && len(rails) < 4:
			m.Turn = 1 - shooter
			m.Pending = &BreakDecision{Kind: "illegal", Breaker: shooter, Foul: out.Foul}
			out.Reason = "The break did not pocket an object ball or drive four object balls to cushions."
		case out.Foul:
			m.Turn = 1 - shooter
			m.Pending = &BreakDecision{Kind: "foul", Breaker: shooter, Foul: true}
		case objectPot:
			out.KeptTurn = true
			out.Reason = "A ball fell on the break. The table remains open."
		default:
			m.Turn = 1 - shooter
			out.Reason = "The break is legal. The other player shoots at an open table."
		}
	} else if eight {
		if out.Foul || !eligibleEight || call.Safety || call.Ball != 8 || call.Pocket != eightPocket {
			m.Winner = 1 - shooter
			out.Reason = "The eight fell before it was legally pocketed in its called pocket."
		} else {
			m.Winner = shooter
			out.Reason = "The eight fell in its called pocket. The rack is won."
		}
		out.Winner = m.Winner
	} else if out.Foul {
		m.Turn = 1 - shooter
		m.InHand = true
	} else {
		pocket, calledPot := pots[call.Ball]
		made := !call.Safety && calledPot && pocket == call.Pocket
		if made {
			if m.Groups[shooter] == 0 {
				m.Groups[shooter] = group(call.Ball)
				m.Groups[1-shooter] = 3 - m.Groups[shooter]
			}
			out.KeptTurn = true
			out.Reason = "The called ball fell in its pocket. Keep shooting."
		} else {
			m.Turn = 1 - shooter
			out.Reason = "The other player shoots."
			if call.Safety {
				out.Reason = "Safety called. The other player shoots."
			}
		}
	}
	m.Last = out
	return out, nil
}

// A ball frozen to a cushion must leave it before that same cushion can
// satisfy the after-contact rule. Replay preserves every impact/turning point.
func (m *Match) newRailContact(e *Engine, result Result, event Event) bool {
	if event.Other < 0 || event.Other >= len(e.rails) {
		return false
	}
	rail := e.rails[event.Other]
	distance := func(p Vec) float64 {
		p.Z = 0
		d := rail.b.sub(rail.a)
		u := clamp(p.sub(rail.a).dot(d)/d.dot(d), 0, 1)
		return p.sub(rail.a.add(d.mul(u))).norm()
	}
	initial := m.ball(event.Ball)
	if initial == nil {
		return false
	}
	if distance(initial.Position) > Radius+1e-6 {
		return true
	}
	for _, frame := range result.Frames {
		if frame.Time >= event.Time-1e-10 {
			break
		}
		for _, b := range frame.Balls {
			if b.ID == event.Ball && distance(b.Position) > Radius+1e-5 {
				return true
			}
		}
	}
	return false
}

func (m *Match) Concede(seat int) error {
	if seat < 0 || seat > 1 {
		return errors.New("invalid player")
	}
	if m.Winner != -1 {
		return errors.New("the rack is finished")
	}
	m.Winner = 1 - seat
	m.Pending = nil
	m.InHand = false
	m.Last = Outcome{Shooter: seat, Winner: m.Winner, Reason: "The rack was conceded."}
	return nil
}

func (m *Match) BreakChoices() []string {
	if m.Pending == nil {
		return nil
	}
	p := m.Pending
	if p.Kind == "eight" {
		return []string{"spot-eight", "rerack-self"}
	}
	if p.Kind == "illegal" {
		choices := []string{"accept", "rerack-self", "rerack-opponent"}
		if p.Foul {
			choices = append(choices, "hand-behind")
		}
		return choices
	}
	choices := []string{"hand-behind"}
	if cue := m.ball(0); cue != nil && !cue.Pocketed {
		choices = append(choices, "accept")
	}
	return choices
}

func (m *Match) DecideBreak(e *Engine, seat int, choice string) error {
	if err := m.checkSeat(seat); err != nil {
		return err
	}
	legal := false
	for _, c := range m.BreakChoices() {
		if c == choice {
			legal = true
		}
	}
	if !legal {
		return errors.New("that break decision is not available")
	}
	next := *m
	next.Balls = append([]Ball(nil), m.Balls...)
	p := *m.Pending
	switch choice {
	case "rerack-self", "rerack-opponent":
		breaker := seat
		if choice == "rerack-opponent" {
			breaker = p.Breaker
		}
		next.Balls = Rack()
		next.Turn = breaker
		next.Groups = [2]int{}
		next.Breaking = true
		next.InHand = true
		next.HeadOnly = true
	case "spot-eight":
		if err := next.spot(8); err != nil {
			return err
		}
		if p.Foul {
			next.InHand = true
			next.HeadOnly = true
		}
	case "hand-behind":
		next.InHand = true
		next.HeadOnly = true
	case "accept":
		if cue := next.ball(0); cue == nil || cue.Pocketed {
			next.InHand = true
			next.HeadOnly = true
		}
	}
	next.Pending = nil
	next.Last.Reason = "The break decision is resolved."
	if next.Breaking {
		next.Last.Reason = "A new rack is ready. Place the cue ball behind the head string."
	} else if next.InHand {
		next.Last.Reason = "Place the cue ball behind the head string."
	} else if choice == "spot-eight" {
		next.Last.Reason = "The eight is spotted. Continue at the open table."
	} else {
		next.Last.Reason = "The table is accepted in position."
	}
	if next.InHand && next.HeadOnly && !next.Breaking {
		// If every legal object ball is behind the head string, spot the closest
		// one so a legal first contact is possible without an artificial stalemate.
		candidate := -1
		closest := -1.0
		allBehind := true
		for _, b := range next.Balls {
			if b.Pocketed || !next.LegalFirst(seat, b.ID) {
				continue
			}
			if b.Position.Y >= HeadString {
				allBehind = false
			}
			if b.Position.Y > closest {
				closest = b.Position.Y
				candidate = b.ID
			}
		}
		if allBehind && candidate >= 0 {
			if err := next.spot(candidate); err != nil {
				return err
			}
		}
	}
	*m = next
	return nil
}

// Place a spotted ball on the long string, nearest the foot spot toward the
// foot rail, then toward the head rail if that half is occupied. Candidate
// interval endpoints avoid a coarse grid moving a ball gratuitously far away.
func (m *Match) spot(id int) error {
	target := m.ball(id)
	if target == nil {
		return errors.New("spotted ball is missing")
	}
	candidates := []float64{Length * .75}
	for _, b := range m.Balls {
		if b.ID == id || b.Pocketed {
			continue
		}
		dx := math.Abs(b.Position.X - Width/2)
		d := 2*Radius + 1e-6
		if dx < d {
			dy := math.Sqrt(d*d - dx*dx)
			candidates = append(candidates, b.Position.Y+dy, b.Position.Y-dy)
		}
	}
	best := math.Inf(1)
	chosen := 0.0
	for _, y := range candidates {
		if y < Radius || y > Length-Radius {
			continue
		}
		p := Vec{Width / 2, y, Radius}
		free := true
		for _, b := range m.Balls {
			if b.ID != id && !b.Pocketed && p.sub(b.Position).norm() < 2*Radius+5e-7 {
				free = false
				break
			}
		}
		if !free {
			continue
		}
		score := y - Length*.75
		if score < 0 {
			score = Length - score
		}
		if score < best {
			best = score
			chosen = y
		}
	}
	if math.IsInf(best, 1) {
		return fmt.Errorf("no space to spot ball %d", id)
	}
	*target = Ball{ID: id, Position: Vec{Width / 2, chosen, Radius}, Orientation: [4]float64{0, 0, 0, 1}}
	return nil
}
