// Package billiards owns deterministic shot physics, independently of rendering,
// campaign time and money. Coordinates are metres: cloth in XY, Z upwards.
package billiards

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

const Radius = .028575
const Width = 1.27
const Length = 2.54
const gravity = 9.81

type Vec struct{ X, Y, Z float64 }

func (a Vec) add(b Vec) Vec     { return Vec{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a Vec) sub(b Vec) Vec     { return Vec{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a Vec) mul(s float64) Vec { return Vec{a.X * s, a.Y * s, a.Z * s} }
func (a Vec) dot(b Vec) float64 { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a Vec) norm() float64     { return math.Sqrt(a.dot(a)) }
func (a Vec) unit() Vec {
	if n := a.norm(); n > 0 {
		return a.mul(1 / n)
	}
	return Vec{}
}
func cross(a, b Vec) Vec    { return Vec{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X} }
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func validVec(v Vec) bool   { return finite(v.X) && finite(v.Y) && finite(v.Z) }

// Orientation is a unit quaternion in XYZW order. Numbered/striped balls can
// replay their true rotation, including a ball spinning without translating.
type Ball struct {
	ID                       int
	Position, Velocity, Spin Vec
	Orientation              [4]float64
	Pocketed                 bool
	Pocket                   int
}

// Shot uses radians, metres/second and metre-based tip offsets, respectively.
type Shot struct{ Angle, Speed, Top, Side float64 }
type Event struct {
	Time        float64
	Kind        string
	Ball, Other int
	Speed       float64
}
type Frame struct {
	Time  float64
	Balls []Ball
}
type Result struct {
	Balls    []Ball
	Events   []Event
	Frames   []Frame
	Duration float64
}
type segment struct{ a, b Vec }
type pocket struct {
	center Vec
	radius float64
}

// Config coefficients are explicit so calibration and convergence tests do not
// need to alter the collision solver. Defaults are a starting cloth/rail setup,
// not a claim to reproduce a measured individual table.
type Config struct {
	Step, Slide, Roll, SpinDrag, BallRestitution, RailRestitution, BallFriction, RailFriction float64
	FrameRate, MaxSeconds                                                                     float64
}

func DefaultConfig() Config {
	return Config{Step: 1.0 / 2400, Slide: .20, Roll: .012, SpinDrag: 8, BallRestitution: .96, RailRestitution: .82, BallFriction: .04, RailFriction: .16, FrameRate: 30, MaxSeconds: 60}
}

type Engine struct {
	Config  Config
	rails   []segment
	pockets []pocket
}

func New() *Engine {
	e := &Engine{Config: DefaultConfig()}
	// Nose segments and short angled facings form actual pocket jaws. The
	// endpoint-circle sweep gives each facing a rounded, ball-radius contact.
	const corner = .085
	const side = .075
	e.rails = []segment{
		{Vec{corner, 0, 0}, Vec{Width - corner, 0, 0}}, {Vec{corner, Length, 0}, Vec{Width - corner, Length, 0}},
	}
	for _, x := range []float64{0, Width} {
		e.rails = append(e.rails, segment{Vec{x, corner, 0}, Vec{x, Length/2 - side, 0}}, segment{Vec{x, Length/2 + side, 0}, Vec{x, Length - corner, 0}})
		outside := -.045
		if x > 0 {
			outside = Width + .045
		}
		e.rails = append(e.rails, segment{Vec{x, Length/2 - side, 0}, Vec{outside, Length/2 - .065, 0}}, segment{Vec{x, Length/2 + side, 0}, Vec{outside, Length/2 + .065, 0}})
	}
	for _, x := range []float64{0, Width} {
		for _, y := range []float64{0, Length} {
			sx, sy := 1.0, 1.0
			if x > 0 {
				sx = -1
			}
			if y > 0 {
				sy = -1
			}
			e.rails = append(e.rails, segment{Vec{x + sx*corner, y, 0}, Vec{x + sx*.049, y - sy*.035, 0}}, segment{Vec{x, y + sy*corner, 0}, Vec{x - sx*.035, y + sy*.049, 0}})
			e.pockets = append(e.pockets, pocket{Vec{x - sx*.026, y - sy*.026, 0}, .067})
		}
	}
	e.pockets = append(e.pockets, pocket{Vec{-.045, Length / 2, 0}, .065}, pocket{Vec{Width + .045, Length / 2, 0}, .065})
	return e
}

// Rack makes a tight eight-ball rack: apex at the foot spot, eight in the
// centre, opposite groups on the rear corners. Randomizing the other numbered
// positions, when desired, belongs to the seeded match rather than physics.
func Rack() []Ball {
	balls := []Ball{{ID: 0, Position: Vec{Width / 2, Length * .25, Radius}, Orientation: [4]float64{0, 0, 0, 1}}}
	order := []int{1, 9, 2, 3, 8, 10, 11, 4, 12, 5, 6, 13, 7, 14, 15}
	k := 0
	d := 2*Radius + 1e-6
	for row := 0; row < 5; row++ {
		for col := 0; col <= row; col++ {
			balls = append(balls, Ball{ID: order[k], Position: Vec{Width/2 + (float64(col)-float64(row)/2)*d, Length*.75 + float64(row)*d*math.Sqrt(3)/2, Radius}, Orientation: [4]float64{0, 0, 0, 1}})
			k++
		}
	}
	return balls
}

func (e *Engine) validate(balls []Ball) error {
	c := e.Config
	for _, v := range []float64{c.Step, c.Slide, c.Roll, c.SpinDrag, c.BallRestitution, c.RailRestitution, c.BallFriction, c.RailFriction, c.FrameRate, c.MaxSeconds} {
		if !finite(v) {
			return errors.New("non-finite physics configuration")
		}
	}
	if c.Step < 1.0/9600 || c.Step > 1.0/120 || c.Slide <= 0 || c.Slide > 1 || c.Roll <= 0 || c.Roll > 1 || c.SpinDrag <= 0 || c.SpinDrag > 100 || c.FrameRate < 1 || c.FrameRate > 60 || c.MaxSeconds <= 0 || c.MaxSeconds > 120 || c.BallRestitution < 0 || c.BallRestitution > 1 || c.RailRestitution < 0 || c.RailRestitution > 1 || c.BallFriction < 0 || c.BallFriction > 1 || c.RailFriction < 0 || c.RailFriction > 1 {
		return errors.New("invalid physics configuration")
	}
	if len(balls) == 0 || len(balls) > 16 {
		return errors.New("a table needs one to sixteen balls")
	}
	ids := map[int]bool{}
	for i, b := range balls {
		if b.ID < 0 || b.ID > 15 || ids[b.ID] {
			return errors.New("invalid or duplicate ball number")
		}
		ids[b.ID] = true
		if !validVec(b.Position) || !validVec(b.Velocity) || !validVec(b.Spin) {
			return errors.New("non-finite ball state")
		}
		for _, q := range b.Orientation {
			if !finite(q) {
				return errors.New("invalid ball orientation")
			}
		}
		if b.Pocketed {
			if b.Pocket < 0 || b.Pocket >= 6 {
				return errors.New("invalid pocket")
			}
			continue
		}
		if math.Abs(b.Position.Z-Radius) > 1e-7 || b.Velocity.Z != 0 || b.Position.X < -.12 || b.Position.X > Width+.12 || b.Position.Y < -.12 || b.Position.Y > Length+.12 || b.Velocity.norm() > 8.001 || b.Spin.norm() > 750 {
			return errors.New("ball outside supported table motion")
		}
		for j := 0; j < i; j++ {
			if !balls[j].Pocketed && b.Position.sub(balls[j].Position).norm() < 2*Radius-1e-7 {
				return errors.New("overlapping balls")
			}
		}
	}
	return nil
}

// Shoot specifies actual initial cue-ball speed and a contact offset in ball
// radii. A horizontal impulse gives I=2mr²/5 angular momentum. The tip is kept
// within 0.6 radii; jump/masse strokes are not implemented by this ground solver.
func (e *Engine) Shoot(initial []Ball, shot Shot) (Result, error) {
	if err := e.validate(initial); err != nil {
		return Result{}, err
	}
	for _, v := range []float64{shot.Angle, shot.Speed, shot.Top, shot.Side} {
		if !finite(v) {
			return Result{}, errors.New("non-finite shot")
		}
	}
	if shot.Speed < .05 || shot.Speed > 8 || math.Hypot(shot.Top, shot.Side) > .6*Radius {
		return Result{}, errors.New("shot outside cue limits")
	}
	balls := append([]Ball(nil), initial...)
	cue := -1
	for i, b := range balls {
		if !b.Pocketed && (b.Velocity.norm() > 1e-8 || b.Spin.norm() > 1e-8) {
			return Result{}, errors.New("balls are still moving")
		}
		if b.ID == 0 && !b.Pocketed {
			cue = i
		}
	}
	if cue < 0 {
		return Result{}, errors.New("place the cue ball before shooting")
	}
	d := Vec{math.Cos(shot.Angle), math.Sin(shot.Angle), 0}
	balls[cue].Velocity = d.mul(shot.Speed)
	balls[cue].Spin = Vec{-d.Y * shot.Top, d.X * shot.Top, -shot.Side}.mul(2.5 * shot.Speed / (Radius * Radius))
	return e.Roll(balls)
}

// Roll evolves an already moving, validated table. Its input is never mutated.
// Motion is swept continuously between friction steps, so even grazing balls
// and thin pocket jaws cannot be skipped by a high-speed shot.
func (e *Engine) Roll(initial []Ball) (Result, error) {
	if err := e.validate(initial); err != nil {
		return Result{}, err
	}
	balls := append([]Ball(nil), initial...)
	sort.Slice(balls, func(i, j int) bool { return balls[i].ID < balls[j].ID })
	for i := range balls {
		normalize(&balls[i].Orientation)
	}
	result := Result{}
	frame := func(t float64) {
		f := Frame{t, append([]Ball(nil), balls...)}
		if n := len(result.Frames); n > 0 && math.Abs(result.Frames[n-1].Time-t) < 1e-10 {
			result.Frames[n-1] = f
		} else {
			result.Frames = append(result.Frames, f)
		}
	}
	frame(0)
	nextFrame := 1 / e.Config.FrameRate
	t := 0.0
	for moving(balls) {
		if t >= e.Config.MaxSeconds {
			return Result{}, errors.New("shot did not settle within physics limit")
		}
		dt := math.Min(e.Config.Step, e.Config.MaxSeconds-t)
		for i := range balls {
			if !balls[i].Pocketed {
				e.cloth(&balls[i], dt/2)
			}
		}
		remaining := dt
		for events := 0; remaining > 1e-10; events++ {
			if events > 128 {
				return Result{}, errors.New("contact solver failed to progress")
			}
			hit := e.next(balls, remaining)
			travel := remaining
			if hit.kind != "" {
				travel = hit.time
			}
			for i := range balls {
				if !balls[i].Pocketed {
					balls[i].Position = balls[i].Position.add(balls[i].Velocity.mul(travel))
					rotate(&balls[i], travel)
				}
			}
			remaining -= travel
			if hit.kind == "" {
				break
			}
			impact := e.resolve(balls, hit)
			result.Events = append(result.Events, Event{t + dt - remaining, hit.kind, balls[hit.a].ID, hit.b, impact})
			if hit.kind == "ball" {
				result.Events[len(result.Events)-1].Other = balls[hit.b].ID
			}
			// Preserve each velocity discontinuity for faithful piecewise replay.
			frame(t + dt - remaining)
		}
		for i := range balls {
			if !balls[i].Pocketed {
				e.cloth(&balls[i], dt/2)
			}
		}
		t += dt
		if t+1e-9 >= nextFrame {
			frame(t)
			nextFrame += 1 / e.Config.FrameRate
		}
	}
	if t > 0 && (len(result.Frames) == 1 || result.Frames[len(result.Frames)-1].Time < t-1e-9) {
		frame(t)
	}
	result.Balls = balls
	result.Duration = t
	return result, nil
}
func moving(balls []Ball) bool {
	for _, b := range balls {
		if !b.Pocketed && (b.Velocity.norm() > 1e-9 || b.Spin.norm() > 1e-9) {
			return true
		}
	}
	return false
}

func (e *Engine) cloth(b *Ball, dt float64) {
	// Independent torsional drag applies during sliding as well as rolling.
	b.Spin.Z = math.Copysign(math.Max(0, math.Abs(b.Spin.Z)-e.Config.SpinDrag*dt), b.Spin.Z)
	// Slip at the cloth is v + omega × (0,0,-R). Friction acts on
	// both translation and rotation; slip decays at 7/2 mu*g, not mu*g.
	slip := b.Velocity.add(cross(b.Spin, Vec{0, 0, -Radius}))
	slip.Z = 0
	if s := slip.norm(); s > 1e-10 {
		slideTime := 2 * s / (7 * e.Config.Slide * gravity)
		h := math.Min(dt, slideTime)
		dv := slip.mul(-e.Config.Slide * gravity * h / s)
		b.Velocity = b.Velocity.add(dv)
		b.Spin = b.Spin.add(cross(Vec{0, 0, -Radius}, dv).mul(2.5 / (Radius * Radius)))
		dt -= h
		if h >= slideTime-1e-12 {
			b.Spin.X = -b.Velocity.Y / Radius
			b.Spin.Y = b.Velocity.X / Radius
		}
	}
	if dt > 0 {
		speed := b.Velocity.norm()
		next := math.Max(0, speed-e.Config.Roll*gravity*dt)
		if speed > 0 {
			b.Velocity = b.Velocity.mul(next / speed)
		}
		b.Spin.X = -b.Velocity.Y / Radius
		b.Spin.Y = b.Velocity.X / Radius
	}
}

func normalize(q *[4]float64) {
	n := math.Sqrt(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])
	if n < 1e-12 {
		*q = [4]float64{0, 0, 0, 1}
		return
	}
	for i := range q {
		q[i] /= n
	}
}
func rotate(b *Ball, dt float64) {
	w := b.Spin.norm()
	if w == 0 {
		return
	}
	s := math.Sin(w*dt/2) / w
	c := math.Cos(w * dt / 2)
	a := b.Spin.mul(s)
	q := b.Orientation
	b.Orientation = [4]float64{c*q[0] + a.X*q[3] + a.Y*q[2] - a.Z*q[1], c*q[1] - a.X*q[2] + a.Y*q[3] + a.Z*q[0], c*q[2] + a.X*q[1] - a.Y*q[0] + a.Z*q[3], c*q[3] - a.X*q[0] - a.Y*q[1] - a.Z*q[2]}
	normalize(&b.Orientation)
}

type contact struct {
	kind   string
	time   float64
	a, b   int
	normal Vec
}

func circleTime(delta, velocity Vec, radius, limit float64) (float64, bool) {
	delta.Z = 0
	velocity.Z = 0
	a := velocity.dot(velocity)
	b := delta.dot(velocity)
	c := delta.dot(delta) - radius*radius
	if a < 1e-18 || b >= -1e-12 {
		return 0, false
	}
	if c <= 1e-12 {
		return 0, true
	}
	disc := b*b - a*c
	if disc < 0 {
		return 0, false
	}
	t := c / (-b + math.Sqrt(disc))
	return t, t >= 0 && t <= limit
}
func (e *Engine) next(balls []Ball, limit float64) contact {
	best := contact{time: limit + 1}
	offer := func(h contact) {
		if h.time < best.time-1e-12 {
			best = h
		}
	}
	for i, a := range balls {
		if a.Pocketed {
			continue
		}
		p := a.Position
		p.Z = 0
		for j, pocket := range e.pockets {
			delta := p.sub(pocket.center)
			if delta.norm() < pocket.radius {
				offer(contact{"pocket", 0, i, j, Vec{}})
			} else if t, ok := circleTime(delta, a.Velocity, pocket.radius, limit); ok {
				offer(contact{"pocket", t, i, j, Vec{}})
			}
		}
		for j := i + 1; j < len(balls); j++ {
			b := balls[j]
			if b.Pocketed {
				continue
			}
			delta := b.Position.sub(a.Position)
			relative := b.Velocity.sub(a.Velocity)
			if t, ok := circleTime(delta, relative, 2*Radius, limit); ok {
				offer(contact{"ball", t, i, j, delta.add(relative.mul(t)).unit()})
			}
		}
		for j, rail := range e.rails {
			tangent := rail.b.sub(rail.a).unit()
			normal := Vec{-tangent.Y, tangent.X, 0}
			distance := p.sub(rail.a).dot(normal)
			vn := a.Velocity.dot(normal)
			for _, sign := range []float64{-1, 1} {
				if vn*sign >= -1e-12 {
					continue
				}
				t := (sign*Radius - distance) / vn
				if t < -1e-10 || t > limit {
					continue
				}
				t = math.Max(0, t)
				at := p.add(a.Velocity.mul(t))
				projection := at.sub(rail.a).dot(tangent)
				if projection >= 0 && projection <= rail.b.sub(rail.a).norm() {
					offer(contact{"rail", t, i, j, normal.mul(sign)})
				}
			}
			for _, end := range []Vec{rail.a, rail.b} {
				if t, ok := circleTime(p.sub(end), a.Velocity, Radius, limit); ok {
					offer(contact{"rail", t, i, j, p.add(a.Velocity.mul(t)).sub(end).unit()})
				}
			}
		}
	}
	if best.time > limit {
		return contact{}
	}
	return best
}
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func (e *Engine) resolve(balls []Ball, h contact) float64 {
	a := &balls[h.a]
	if h.kind == "pocket" {
		speed := a.Velocity.norm()
		a.Pocketed = true
		a.Pocket = h.b
		a.Velocity = Vec{}
		a.Spin = Vec{}
		return speed
	}
	n := h.normal
	t := Vec{-n.Y, n.X, 0}
	if h.kind == "ball" {
		b := &balls[h.b]
		relative := b.Velocity.sub(a.Velocity)
		closing := relative.dot(n)
		jn := -(1 + e.Config.BallRestitution) * closing / 2
		slip := relative.dot(t) - Radius*(a.Spin.Z+b.Spin.Z)
		jt := clamp(-slip/7, -e.Config.BallFriction*jn, e.Config.BallFriction*jn)
		impulse := n.mul(jn).add(t.mul(jt))
		a.Velocity = a.Velocity.sub(impulse)
		b.Velocity = b.Velocity.add(impulse)
		a.Spin.Z -= 2.5 * jt / Radius
		b.Spin.Z -= 2.5 * jt / Radius
		return -closing
	}
	closing := a.Velocity.dot(n)
	jn := -(1 + e.Config.RailRestitution) * closing
	slip := a.Velocity.dot(t) - Radius*a.Spin.Z
	jt := clamp(-slip/3.5, -e.Config.RailFriction*jn, e.Config.RailFriction*jn)
	a.Velocity = a.Velocity.add(n.mul(jn)).add(t.mul(jt))
	a.Spin.Z -= 2.5 * jt / Radius
	return -closing
}

// ValidatePlacement checks a stationary ball-in-hand position on the cloth,
// including pocket mouths, jaw faces and every unpocketed object ball.
func (e *Engine) ValidatePlacement(balls []Ball, x, y float64) error {
	if !finite(x) || !finite(y) || x < Radius || x > Width-Radius || y < Radius || y > Length-Radius {
		return errors.New("place the cue ball on the cloth")
	}
	p := Vec{x, y, Radius}
	for _, b := range balls {
		if b.ID != 0 && !b.Pocketed && p.sub(b.Position).norm() < 2*Radius+1e-5 {
			return fmt.Errorf("cue ball overlaps ball %d", b.ID)
		}
	}
	for _, k := range e.pockets {
		d := p.sub(k.center)
		d.Z = 0
		if d.norm() < k.radius {
			return errors.New("cue ball is in a pocket mouth")
		}
	}
	return nil
}
