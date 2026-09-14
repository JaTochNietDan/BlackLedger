package billiards

import (
	"errors"
	"math"
	"sort"
)

type Intent struct {
	Shot Shot
	Call Call
}
type OpponentTurn struct {
	Match     Match
	Result    TurnResult
	Intent    Intent
	Placement *Vec
	Decision  string
}
type candidate struct {
	Intent
	cost float64
}

// Opponent plays precisely one stroke. It prepares required break decisions or
// ball-in-hand placement, plans with visible geometry, and executes the chosen
// cue input through Match.Play. Nothing here can grant a pot or a win directly.
func Opponent(e *Engine, source *Match, seat int, skill float64, seed uint64) (OpponentTurn, error) {
	var turn OpponentTurn
	if source == nil || !finite(skill) || skill < 0 || skill > 1 {
		return turn, errors.New("invalid opponent state")
	}
	if err := source.checkSeat(seat); err != nil {
		return turn, err
	}
	m := *source
	m.Balls = append([]Ball(nil), source.Balls...)
	if m.Pending != nil {
		choices := m.BreakChoices()
		choice := choices[0]
	chooseDecision:
		for _, want := range []string{"spot-eight", "hand-behind", "rerack-self"} {
			for _, c := range choices {
				if c == want {
					choice = want
					break chooseDecision
				}
			}
		}
		if err := m.DecideBreak(e, seat, choice); err != nil {
			return turn, err
		}
		turn.Decision = choice
	}
	if m.InHand {
		position, err := opponentPlacement(e, &m, seat)
		if err != nil {
			return turn, err
		}
		if err = m.PlaceCue(e, seat, position.X, position.Y); err != nil {
			return turn, err
		}
		turn.Placement = &position
	}
	plan, err := planOpponent(e, &m, seat)
	if err != nil {
		return turn, err
	}
	// Skill changes execution accuracy, never the pocketed balls or settlement.
	// This private local stream cannot consume campaign RNG or change on retries.
	random := func() float64 {
		seed += 0x9e3779b97f4a7c15
		z := seed
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		return float64(z>>11) / float64(uint64(1)<<53)
	}
	plan.Shot.Angle += (random()*2 - 1) * (1 - skill) * .015
	plan.Shot.Speed = clamp(plan.Shot.Speed*(1+(random()*2-1)*(1-skill)*.12), .05, 8)
	result, err := m.Play(e, seat, plan.Shot, plan.Call)
	if err != nil {
		return turn, err
	}
	turn.Match = m
	turn.Result = result
	turn.Intent = plan
	return turn, nil
}

func opponentPlacement(e *Engine, m *Match, seat int) (Vec, error) {
	candidates := []Vec{}
	if m.Breaking {
		candidates = append(candidates, Vec{Width * .35, HeadString - .04, Radius})
	}
	for _, b := range m.Balls {
		if b.Pocketed || !m.LegalFirst(seat, b.ID) {
			continue
		}
		for _, p := range e.pockets {
			direction := p.center.sub(Vec{b.Position.X, b.Position.Y, 0}).unit()
			for _, distance := range []float64{.20, .35, .55} {
				candidates = append(candidates, b.Position.sub(direction.mul(2*Radius+distance)))
			}
		}
	}
	for _, y := range []float64{.15, .35, .55, .85, 1.2, 1.6, 2, 2.3} {
		for _, x := range []float64{.15, .4, .635, .87, 1.12} {
			candidates = append(candidates, Vec{x, y, Radius})
		}
	}
	best := math.Inf(1)
	var chosen Vec
	found := false
	for _, p := range candidates {
		if m.HeadOnly && p.Y >= HeadString {
			continue
		}
		if err := e.ValidatePlacement(m.Balls, p.X, p.Y); err != nil {
			continue
		}
		if m.Breaking {
			return p, nil
		}
		copy := *m
		copy.Balls = append([]Ball(nil), m.Balls...)
		cue := copy.ball(0)
		if cue == nil {
			return Vec{}, errors.New("cue ball missing")
		}
		cue.Position = p
		cue.Pocketed = false
		plans := potCandidates(e, &copy, seat)
		cost := 100.0
		if len(plans) > 0 {
			cost = plans[0].cost
		}
		if cost < best {
			best = cost
			chosen = p
			found = true
		}
	}
	if !found {
		return Vec{}, errors.New("no legal cue-ball placement")
	}
	return chosen, nil
}

func planOpponent(e *Engine, m *Match, seat int) (Intent, error) {
	if m.Breaking {
		cue := m.ball(0)
		if cue == nil {
			return Intent{}, errors.New("cue ball missing")
		}
		delta := Vec{Width / 2, Length * .75, Radius}.sub(cue.Position)
		base := math.Atan2(delta.Y, delta.X)
		best := Intent{Shot: Shot{Angle: base, Speed: 7.5}}
		bestScore := math.Inf(-1)
		for _, offset := range []float64{0, -.015, .015, -.03, .03, -.06, .06} {
			shot := Shot{Angle: base + offset, Speed: 7.5}
			copy := *m
			copy.Balls = append([]Ball(nil), m.Balls...)
			result, err := copy.Play(e, seat, shot, Call{})
			if err != nil {
				continue
			}
			rails := map[int]bool{}
			pots := 0
			for _, ev := range result.Physics.Events {
				if ev.Kind == "rail" && ev.Ball != 0 {
					rails[ev.Ball] = true
				}
				if ev.Kind == "pocket" && ev.Ball != 0 {
					pots++
				}
			}
			score := float64(len(rails)*2 + pots*10)
			if result.Outcome.Foul {
				score -= 50
			}
			if copy.Pending != nil && copy.Pending.Kind == "illegal" {
				score -= 100
			}
			if score > bestScore {
				bestScore = score
				best = Intent{Shot: shot}
			}
		}
		return best, nil
	}

	plans := potCandidates(e, m, seat)
	if len(plans) > 8 {
		plans = plans[:8]
	}
	// Preview draw and follow on the strongest pot lines. All offsets are metres.
	// Keep centre hits first so a needless spin stroke loses a tied evaluation.
	potLines := append([]candidate(nil), plans...)
	for i, line := range potLines {
		if i >= 4 {
			break
		}
		for _, top := range []float64{-.012, -.008, .008, .012} {
			variation := line
			variation.Shot.Top = top
			variation.cost += .03
			plans = append(plans, variation)
		}
	}
	cue := m.ball(0)
	if cue == nil || cue.Pocketed {
		return Intent{}, errors.New("cue ball not on table")
	}
	// A direct safety and one-cushion escapes remain candidates when no pot has
	// a clear line. These are physically evaluated, not assumed to make contact.
	targets := []Ball{}
	for _, b := range m.Balls {
		if !b.Pocketed && m.LegalFirst(seat, b.ID) {
			targets = append(targets, b)
		}
	}
	sort.SliceStable(targets, func(i, j int) bool {
		return targets[i].Position.sub(cue.Position).norm() < targets[j].Position.sub(cue.Position).norm()
	})
	if len(targets) > 3 {
		targets = targets[:3]
	}
	for _, b := range targets {
		d := b.Position.sub(cue.Position)
		plans = append(plans, candidate{Intent{Shot: Shot{Angle: math.Atan2(d.Y, d.X), Speed: clamp(2.4+d.norm()*.5, 2.4, 3.8)}, Call: Call{Safety: true}}, 100 + d.norm()})
		// Reflect the target across the ball-centre rail line to aim a kick.
		for axis := 0; axis < 4; axis++ {
			mirrored := b.Position
			switch axis {
			case 0:
				mirrored.X = 2*Radius - b.Position.X
			case 1:
				mirrored.X = 2*(Width-Radius) - b.Position.X
			case 2:
				mirrored.Y = 2*Radius - b.Position.Y
			case 3:
				mirrored.Y = 2*(Length-Radius) - b.Position.Y
			}
			delta := mirrored.sub(cue.Position)
			for _, offset := range []float64{0, -.12, -.06, .06, .12} {
				plans = append(plans, candidate{Intent{Shot: Shot{Angle: math.Atan2(delta.Y, delta.X) + offset, Speed: clamp(3.6+delta.norm()*.5, 3.6, 5)}, Call: Call{Safety: true}}, 110 + delta.norm() + math.Abs(offset)})
			}
		}
	}
	if len(plans) == 0 {
		return Intent{}, errors.New("no object ball to play")
	}
	preview := *e
	preview.Config.Step = 1.0 / 600
	bestScore := math.Inf(-1)
	best := plans[0].Intent
	for _, plan := range plans {
		// A declared safety cannot retain the turn or win. Once a better result
		// exists, its maximum possible score cannot justify another simulation.
		if plan.Call.Safety && bestScore >= -plan.cost {
			continue
		}
		copy := *m
		copy.Balls = append([]Ball(nil), m.Balls...)
		result, err := copy.Play(&preview, seat, plan.Shot, plan.Call)
		if err != nil {
			continue
		}
		out := result.Outcome
		score := -plan.cost
		if out.Foul {
			score -= 10000
		}
		if out.KeptTurn {
			score += 1000
			if copy.Winner < 0 {
				// Evaluate the actual settled cue ball: a clear next pot is worth
				// more than stranding it behind traffic or against a cushion.
				next := potCandidates(&preview, &copy, seat)
				if len(next) == 0 {
					score -= 80
				} else {
					score -= math.Min(80, next[0].cost*10)
				}
			}
		}
		if copy.Winner == seat {
			score += 100000
		}
		if copy.Winner == 1-seat {
			score -= 1000000
		}
		if score > bestScore {
			bestScore = score
			best = plan.Intent
		}
	}
	if math.IsInf(bestScore, -1) {
		return Intent{}, errors.New("opponent could not simulate a legal input")
	}
	return best, nil
}

func clearLine(balls []Ball, a, b Vec, ignore int) bool {
	a.Z = 0
	b.Z = 0
	d := b.sub(a)
	length2 := d.dot(d)
	if length2 < 1e-12 {
		return false
	}
	for _, ball := range balls {
		if ball.Pocketed || ball.ID == 0 || ball.ID == ignore {
			continue
		}
		p := ball.Position
		p.Z = 0
		along := p.sub(a).dot(d) / length2
		if along < 0 || along > 1 {
			continue
		}
		if p.sub(a.add(d.mul(along))).norm() < 2*Radius+.001 {
			return false
		}
	}
	return true
}

func potCandidates(e *Engine, m *Match, seat int) []candidate {
	cue := m.ball(0)
	if cue == nil {
		return nil
	}
	plans := []candidate{}
	for _, ball := range m.Balls {
		if ball.Pocketed || !m.LegalFirst(seat, ball.ID) {
			continue
		}
		object := ball.Position
		object.Z = 0
		from := cue.Position
		from.Z = 0
		for pocketID, pocket := range e.pockets {
			// Direct pots and single-cushion banks use the same declared pocket.
			targets := []Vec{pocket.center}
			banks := []bool{false}
			for axis := 0; axis < 4; axis++ {
				mirrored := pocket.center
				line := Radius
				if axis == 1 {
					line = Width - Radius
				}
				if axis == 3 {
					line = Length - Radius
				}
				if axis < 2 {
					mirrored.X = 2*line - pocket.center.X
				} else {
					mirrored.Y = 2*line - pocket.center.Y
				}
				delta := mirrored.sub(object)
				den := delta.X
				origin := object.X
				if axis >= 2 {
					den = delta.Y
					origin = object.Y
				}
				if math.Abs(den) < 1e-9 {
					continue
				}
				fraction := (line - origin) / den
				if fraction <= 0 || fraction >= 1 {
					continue
				}
				bounce := object.add(delta.mul(fraction))
				if bounce.X < Radius-1e-6 || bounce.X > Width-Radius+1e-6 || bounce.Y < Radius-1e-6 || bounce.Y > Length-Radius+1e-6 {
					continue
				}
				if !clearLine(m.Balls, bounce, pocket.center, ball.ID) {
					continue
				}
				targets = append(targets, bounce)
				banks = append(banks, true)
			}
			for i, target := range targets {
				direction := target.sub(object).unit()
				ghost := object.sub(direction.mul(2 * Radius))
				aim := ghost.sub(from)
				distance := aim.norm()
				if distance < .001 || ghost.X < Radius || ghost.X > Width-Radius || ghost.Y < Radius || ghost.Y > Length-Radius {
					continue
				}
				cut := aim.unit().dot(direction)
				if cut < .25 {
					continue
				}
				if !clearLine(m.Balls, from, ghost, ball.ID) || !clearLine(m.Balls, object, target, ball.ID) {
					continue
				}
				travel := target.sub(object).norm()
				if banks[i] {
					travel += pocket.center.sub(target).norm()
				}
				speed := clamp((math.Sqrt(2*e.Config.Roll*gravity*travel)+.35)/(cut*.75)+distance*.4, .6, 4.5)
				cost := travel + distance*.5 + (1-cut)*2
				if banks[i] {
					cost += 2
					speed = clamp(speed*1.2, .6, 5)
				}
				plans = append(plans, candidate{Intent{Shot: Shot{Angle: math.Atan2(aim.Y, aim.X), Speed: speed}, Call: Call{Ball: ball.ID, Pocket: pocketID}}, cost})
			}
		}
	}
	sort.SliceStable(plans, func(i, j int) bool { return plans[i].cost < plans[j].cost })
	return plans
}
