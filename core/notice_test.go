package core

import "testing"

// A business was pure upside once bought: takings, and now cover. What a place
// IS never cost the owner anything. A revue bar with a late licence and a
// gambling house are watched in a way a laundry is not, and that ought to be
// part of choosing what to own.
//
// It is built as a DRAG ON FORGETTING rather than as attention added each day,
// and the reason is arithmetic. Attention fades by one a day. Anything that
// adds more than one a day therefore climbs without limit: past forty-five the
// police arrive, past eighty they take the premises, and a player who bought a
// burlesque and did nothing else would lose it in under a month with no way to
// stop it. A drag cannot run away, because it never adds anything — at worst
// what you have earned simply does not fade.

func TestWhatYouOwnDecidesHowFastTheCityForgetsYou(t *testing.T) {
	t.Parallel()
	quiet, loud := 0, 0
	for kind, trade := range trades {
		if trade.Watched < 0 || trade.Watched > 3 {
			t.Errorf("a %s carries a notice of %d, outside 0..3", kind, trade.Watched)
		}
		if trade.Watched == 0 {
			quiet++
		}
		if trade.Watched >= 3 {
			loud++
		}
	}
	if quiet == 0 {
		t.Error("every business in the city is watched, so the choice is not a choice")
	}
	if loud == 0 {
		t.Error("no business is watched enough to stop the city forgetting you")
	}
	// A revue bar and a gambling house are watched; a laundry is not.
	laundry, _ := TradeOfKind("laundry")
	burlesque, _ := TradeOfKind("burlesque")
	if burlesque.Watched <= laundry.Watched {
		t.Errorf("a burlesque draws %d and a laundry %d", burlesque.Watched, laundry.Watched)
	}
}

// The measurement, with a control that shares the run of luck: the same player
// with the same attention, differing only in which business they hold.
func TestAWatchedBusinessKeepsYouInView(t *testing.T) {
	t.Parallel()
	cool := func(id string) int {
		w := New(53)
		w.District = 2
		w.Player.Cash, w.Player.Respect = 40000, 200
		w.Player.Heat = 40
		w.RNG, w.WorldRNG = 999, 999
		if id != "" {
			w.Properties[id].Owner = "player:1"
		}
		// The clock has to move. Cooling happens every Nth NIGHT, so calling the
		// day ten times without advancing it is the same night ten times — the
		// first version of this measured exactly that and reported no
		// difference at all between a laundry and a revue bar.
		for day := 0; day < 20; day++ {
			w.Minute = day * 1440
			w.RNG, w.WorldRNG = 999, 999 // the same night's luck in every run
			w.PoliceDay()
		}
		return 40 - w.Player.Heat
	}
	nothing, laundry, burlesque := cool(""), cool("laundry"), cool("burlesque")
	if laundry <= burlesque {
		t.Errorf("twenty nights: a laundry let %d points fade and a burlesque %d", laundry, burlesque)
	}
	if nothing <= 0 {
		t.Fatalf("owning nothing at all faded %d points, so this measures nothing", nothing)
	}
	t.Logf("attention faded over twenty nights: owning nothing %d, a laundry %d, a burlesque %d",
		nothing, laundry, burlesque)
}

// And it can never make attention climb. A drag on forgetting is bounded by
// construction and this is the check that keeps it so.
func TestBeingWatchedNeverRaisesAttentionByItself(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.District = 2
	w.Player.Cash, w.Player.Respect = 40000, 200
	w.Player.Heat = 30
	for _, id := range []string{"burlesque", "casino", "goldenlily", "poolhall"} {
		w.Properties[id].Owner = "player:1"
	}
	for day := 0; day < 60; day++ {
		before := w.Player.Heat
		w.PoliceDay()
		if w.Player.Heat > before {
			t.Fatalf("day %d: owning four watched businesses raised attention from %d to %d",
				day, before, w.Player.Heat)
		}
	}
}
