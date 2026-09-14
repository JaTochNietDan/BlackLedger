package billiards

import (
	"math"
	"reflect"
	"testing"
)

func ball(id int, x, y float64) Ball {
	return Ball{ID: id, Position: Vec{x, y, Radius}, Orientation: [4]float64{0, 0, 0, 1}}
}
func near(t *testing.T, got, want, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Fatalf("got %.10f want %.10f ± %.10f", got, want, tolerance)
	}
}
func kinetic(b Ball) float64 {
	return b.Velocity.dot(b.Velocity)/2 + Radius*Radius*b.Spin.dot(b.Spin)/5
}

func TestSlidingReachesAnalyticRollingState(t *testing.T) {
	e := New()
	b := ball(0, .6, 1)
	b.Velocity.X = 1
	transition := 2 / (7 * e.Config.Slide * gravity)
	e.cloth(&b, transition)
	near(t, b.Velocity.X, 5.0/7, 1e-12)
	near(t, b.Spin.Y, 5.0/(7*Radius), 1e-10)
	near(t, b.Velocity.X-Radius*b.Spin.Y, 0, 1e-12)
	// Pure rolling has constant deceleration and comes to rest without reversal.
	e.cloth(&b, 1)
	near(t, b.Velocity.X, 5.0/7-e.Config.Roll*gravity, 1e-12)
	e.cloth(&b, 100)
	if b.Velocity.norm() != 0 || b.Spin.norm() != 0 {
		t.Fatal("rolling did not stop", b)
	}
}
func TestElasticBallCollisionConservesMomentumAndEnergy(t *testing.T) {
	e := New()
	e.Config.BallRestitution = 1
	e.Config.BallFriction = 0
	balls := []Ball{ball(0, .4, 1), ball(1, .4+2*Radius, 1)}
	balls[0].Velocity = Vec{2, .5, 0}
	balls[1].Velocity = Vec{-.4, -.2, 0}
	before := balls[0].Velocity.add(balls[1].Velocity)
	energy := kinetic(balls[0]) + kinetic(balls[1])
	e.resolve(balls, contact{kind: "ball", a: 0, b: 1, normal: Vec{1, 0, 0}})
	near(t, balls[0].Velocity.X, -.4, 1e-12)
	near(t, balls[1].Velocity.X, 2, 1e-12)
	near(t, balls[0].Velocity.add(balls[1].Velocity).sub(before).norm(), 0, 1e-12)
	near(t, kinetic(balls[0])+kinetic(balls[1]), energy, 1e-12)
}
func TestFrictionalImpactCannotCreateEnergy(t *testing.T) {
	e := New()
	for _, spin := range []float64{-150, 0, 150} {
		for _, angle := range []float64{-.8, 0, .8} {
			balls := []Ball{ball(0, .4, 1), ball(1, .4+2*Radius, 1)}
			balls[0].Velocity = Vec{2, angle, 0}
			balls[0].Spin.Z = spin
			before := kinetic(balls[0]) + kinetic(balls[1])
			e.resolve(balls, contact{kind: "ball", a: 0, b: 1, normal: Vec{1, 0, 0}})
			if after := kinetic(balls[0]) + kinetic(balls[1]); after > before+1e-10 {
				t.Fatalf("collision created energy %f -> %f", before, after)
			}
		}
	}
}
func TestSideSpinChangesCushionRebound(t *testing.T) {
	e := New()
	var tangential []float64
	for _, spin := range []float64{-70, 0, 70} {
		balls := []Ball{ball(0, Radius, 1)}
		balls[0].Velocity = Vec{-2, 0, 0}
		balls[0].Spin.Z = spin
		before := kinetic(balls[0])
		e.resolve(balls, contact{kind: "rail", a: 0, normal: Vec{1, 0, 0}})
		near(t, balls[0].Velocity.X, 2*e.Config.RailRestitution, 1e-12)
		if kinetic(balls[0]) > before {
			t.Fatal("cushion created energy")
		}
		tangential = append(tangential, balls[0].Velocity.Y)
	}
	if !(tangential[0] < 0 && tangential[1] == 0 && tangential[2] > 0) {
		t.Fatal("english did not change rebound", tangential)
	}
}
func TestContinuousSweepCatchesGrazingCollision(t *testing.T) {
	// Both endpoints miss the target; the swept trajectory must still collide.
	d := Vec{-.1, 2*Radius - .00001, 0}
	v := Vec{16, 0, 0}
	at, ok := circleTime(d, v, 2*Radius, .0125)
	if !ok || at <= 0 || at >= .0125 {
		t.Fatal("missed high-speed grazing contact", at, ok)
	}
	near(t, d.add(v.mul(at)).norm(), 2*Radius, 1e-12)
}
func TestEveryPocketAcceptsDirectShot(t *testing.T) {
	e := New()
	for i, p := range e.pockets {
		direction := p.center.sub(Vec{Width / 2, Length / 2, 0}).unit()
		start := p.center.sub(direction.mul(.30))
		b := ball(0, start.X, start.Y)
		b.Velocity = direction.mul(1.5)
		got, err := e.Roll([]Ball{b})
		if err != nil {
			t.Fatal(i, err)
		}
		if !got.Balls[0].Pocketed || got.Balls[0].Pocket != i {
			t.Fatal("direct shot missed pocket", i, got.Balls[0], got.Events)
		}
	}
}
func TestDrawAndFollowComeFromAngularMomentum(t *testing.T) {
	e := New()
	var travel []float64
	for _, top := range []float64{-.5, 0, .5} {
		got, err := e.Shoot([]Ball{ball(0, .40, 1), ball(1, .40+2*Radius+.015, 1)}, Shot{Angle: 0, Speed: 1.2, Top: top})
		if err != nil {
			t.Fatal(err)
		}
		travel = append(travel, got.Balls[0].Position.X)
	}
	if !(travel[0] < .40 && travel[2] > travel[1] && travel[1] > travel[0]) {
		t.Fatal("draw/stop/follow order is wrong", travel)
	}
}
func TestRackBreakSettlesDeterministicallyWithoutTunnelling(t *testing.T) {
	e := New()
	initial := Rack()
	saved := append([]Ball(nil), initial...)
	a, err := e.Shoot(initial, Shot{Angle: math.Pi / 2, Speed: 7})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(initial, saved) {
		t.Fatal("input was mutated")
	}
	reverse := append([]Ball(nil), initial...)
	for i, j := 0, len(reverse)-1; i < j; i, j = i+1, j-1 {
		reverse[i], reverse[j] = reverse[j], reverse[i]
	}
	b, err := e.Shoot(reverse, Shot{Angle: math.Pi / 2, Speed: 7})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("shot depends on input ordering")
	}
	if len(a.Events) < 10 || a.Duration < 1 || moving(a.Balls) {
		t.Fatal("break did not simulate", len(a.Events), a.Duration)
	}
	for _, f := range a.Frames {
		for i, b := range f.Balls {
			if b.Pocketed {
				continue
			}
			if b.Position.X < -.10 || b.Position.X > Width+.10 || b.Position.Y < -.10 || b.Position.Y > Length+.10 {
				t.Fatal("ball escaped the table", b)
			}
			for j := 0; j < i; j++ {
				if !f.Balls[j].Pocketed && b.Position.sub(f.Balls[j].Position).norm() < 2*Radius-1e-6 {
					t.Fatal("balls overlapped", b.ID, f.Balls[j].ID)
				}
			}
			q := b.Orientation
			near(t, q[0]*q[0]+q[1]*q[1]+q[2]*q[2]+q[3]*q[3], 1, 1e-10)
		}
	}
}
func TestRestingSpinDecaysAndReplayHasRotation(t *testing.T) {
	e := New()
	b := ball(0, .6, 1)
	b.Spin.Z = 16
	got, err := e.Roll([]Ball{b})
	if err != nil {
		t.Fatal(err)
	}
	near(t, got.Duration, 2, 1e-7)
	near(t, got.Balls[0].Position.sub(b.Position).norm(), 0, 1e-12)
	if got.Frames[10].Balls[0].Orientation == b.Orientation {
		t.Fatal("spinning ball has no replay rotation")
	}
}
func TestInvalidShotsAndPlacementsAreRejected(t *testing.T) {
	e := New()
	rack := Rack()
	for _, shot := range []Shot{{Speed: math.NaN()}, {Speed: 9}, {Speed: 1, Top: .61}, {Speed: 1, Angle: math.Inf(1)}} {
		if _, err := e.Shoot(rack, shot); err == nil {
			t.Fatal("accepted invalid shot", shot)
		}
	}
	if err := e.ValidatePlacement(rack, rack[1].Position.X, rack[1].Position.Y); err == nil {
		t.Fatal("accepted overlapping cue placement")
	}
	if err := e.ValidatePlacement(rack, Width/2, Length/4); err != nil {
		t.Fatal(err)
	}
	rack[1].Position = rack[0].Position
	if _, err := e.Shoot(rack, Shot{Speed: 1}); err == nil {
		t.Fatal("accepted overlapping initial state")
	}
}

func TestFreeSlideAndRollTravelMatchesAnalyticDistance(t *testing.T) {
	e := New()
	b := ball(0, .6, .1)
	b.Velocity.Y = .75
	got, err := e.Roll([]Ball{b})
	if err != nil {
		t.Fatal(err)
	}
	slidingTime := 2 * .75 / (7 * e.Config.Slide * gravity)
	slidingDistance := .75*slidingTime - .5*e.Config.Slide*gravity*slidingTime*slidingTime
	rollingSpeed := .75 * 5 / 7
	rollingDistance := rollingSpeed * rollingSpeed / (2 * e.Config.Roll * gravity)
	near(t, got.Balls[0].Position.Y, .1+slidingDistance+rollingDistance, 2e-6)
}
func TestBreaksWithPowerAimAndSpinStayPhysical(t *testing.T) {
	e := New()
	for _, speed := range []float64{.2, 2, 5, 8} {
		for _, aim := range []float64{-.08, 0, .08} {
			for _, side := range []float64{-.4, 0, .4} {
				initial := Rack()
				initial[0].Position.X += aim
				result, err := e.Shoot(initial, Shot{Angle: math.Pi/2 + aim/5, Speed: speed, Top: .3, Side: side})
				if err != nil {
					t.Fatalf("speed=%f aim=%f side=%f: %v", speed, aim, side, err)
				}
				previous := -1.0
				for _, frame := range result.Frames {
					if frame.Time <= previous {
						t.Fatal("replay times are not strictly ordered")
					}
					previous = frame.Time
					for i, b := range frame.Balls {
						if b.Pocketed {
							continue
						}
						if !validVec(b.Position) || !validVec(b.Velocity) || !validVec(b.Spin) {
							t.Fatal("non-finite motion")
						}
						if b.Position.X < -.10 || b.Position.X > Width+.10 || b.Position.Y < -.10 || b.Position.Y > Length+.10 {
							t.Fatalf("ball %d escaped table at %f", b.ID, frame.Time)
						}
						for j := 0; j < i; j++ {
							if !frame.Balls[j].Pocketed && b.Position.sub(frame.Balls[j].Position).norm() < 2*Radius-1e-6 {
								t.Fatal("interpenetration", b.ID, frame.Balls[j].ID, frame.Time)
							}
						}
					}
				}
				for _, impact := range result.Events {
					found := false
					for _, frame := range result.Frames {
						if math.Abs(frame.Time-impact.Time) < 1e-9 {
							found = true
							break
						}
					}
					if !found {
						t.Fatal("replay omits an impact", impact)
					}
				}
			}
		}
	}
}
func TestClothIntegrationConvergesWithHalfStep(t *testing.T) {
	e := New()
	initial := []Ball{ball(0, .30, .5), ball(1, .70, .7)}
	shot := Shot{Angle: .40, Speed: 1.6, Top: .25, Side: .1}
	coarse, err := e.Shoot(initial, shot)
	if err != nil {
		t.Fatal(err)
	}
	e.Config.Step /= 2
	fine, err := e.Shoot(initial, shot)
	if err != nil {
		t.Fatal(err)
	}
	if len(coarse.Events) != len(fine.Events) {
		t.Fatal("halving the step changed event count", len(coarse.Events), len(fine.Events))
	}
	for i, b := range coarse.Balls {
		if b.Pocketed != fine.Balls[i].Pocketed {
			t.Fatal("halving step changed pocket result")
		}
		near(t, b.Position.sub(fine.Balls[i].Position).norm(), 0, .0001)
	}
}

func BenchmarkRackBreak(b *testing.B) {
	e := New()
	initial := Rack()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := e.Shoot(initial, Shot{Angle: math.Pi / 2, Speed: 7}); err != nil {
			b.Fatal(err)
		}
	}
}

func TestPocketJawDeflectsAnOffCentreEntry(t *testing.T) {
	e := New()
	b := ball(0, .15, Length/2+.05)
	b.Velocity.X = -1
	got, err := e.Roll([]Ball{b})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Events) == 0 || got.Events[0].Kind != "rail" {
		t.Fatal("shot clipped a jaw without a rebound", got.Events)
	}
	// A trajectory parallel to a side-pocket entry must not be pocketed merely
	// because it is near a pocket: it has to pass the actual capture geometry.
	b = ball(0, .15, Length/2+.12)
	b.Velocity.X = -1
	got, err = e.Roll([]Ball{b})
	if err != nil {
		t.Fatal(err)
	}
	if got.Balls[0].Pocketed {
		t.Fatal("pocket swallowed a rail shot")
	}
}
func TestPathologicalConfigurationsAreBounded(t *testing.T) {
	for _, edit := range []func(*Config){func(c *Config) { c.Step = 1e-300 }, func(c *Config) { c.Slide = math.MaxFloat64 }, func(c *Config) { c.FrameRate = 1e9 }, func(c *Config) { c.MaxSeconds = 1e9 }} {
		e := New()
		edit(&e.Config)
		if _, err := e.Shoot(Rack(), Shot{Speed: 1}); err == nil {
			t.Fatal("accepted unbounded physics configuration")
		}
	}
}
