package core

import (
	"strings"
	"testing"
)

// The city view showed who was standing in each place and nothing about the
// place. A player looking at twelve addresses could not see that one was out of
// soap, one had a press broken and one had a still running in the back.

func TestAPlaceSaysTheMostImportantTrueThingAboutItself(t *testing.T) {
	t.Parallel()
	w, member := testator(t)
	prop := w.Properties["laundry"]
	trade, _ := TradeOf("laundry")

	// A working business says how it is doing and is not a worry.
	prop.Staff, prop.Supply, prop.Condition, prop.Trouble = trade.Hands, trade.RestockAmount, 100, false
	if note := w.PlaceNote("laundry"); note == "" {
		t.Fatal("a working business says nothing about itself")
	}
	if w.PlaceWarn("laundry") {
		t.Fatalf("a healthy business is a worry: %q", w.PlaceNote("laundry"))
	}

	// Each thing that can be wrong outranks the ones below it.
	prop.Condition = 40
	if !strings.Contains(w.PlaceNote("laundry"), "repair") || !w.PlaceWarn("laundry") {
		t.Fatalf("a run-down business reads %q", w.PlaceNote("laundry"))
	}
	prop.Staff = 0
	if !strings.Contains(w.PlaceNote("laundry"), "Short-handed") {
		t.Fatalf("an empty business reads %q", w.PlaceNote("laundry"))
	}
	prop.Supply = 0
	if !strings.Contains(w.PlaceNote("laundry"), trade.Supplies) {
		t.Fatalf("a business out of supplies reads %q", w.PlaceNote("laundry"))
	}
	prop.Trouble = true
	if w.PlaceNote("laundry") != trade.Trouble {
		t.Fatalf("a business in trouble reads %q", w.PlaceNote("laundry"))
	}
	// And the police outrank all of it.
	w.Player.Heat = ForfeitThreshold
	if !strings.Contains(w.PlaceNote("laundry"), "take it") {
		t.Fatalf("at %d attention the laundry reads %q", w.Player.Heat, w.PlaceNote("laundry"))
	}

	// A man on the door is worth saying when nothing is wrong.
	w.Player.Heat = 0
	prop.Staff, prop.Supply, prop.Condition, prop.Trouble = trade.Hands, trade.RestockAmount, 100, false
	w.Player.Location = "laundry"
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(w.PlaceNote("laundry"), member.Name) {
		t.Fatalf("with a man on the door it reads %q", w.PlaceNote("laundry"))
	}
}

func TestAPasserByOnlySeesWhatIsVisibleFromTheStreet(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	// Somebody else's premises, out of supplies and in trouble.
	prop := w.Properties["club"]
	prop.Trouble, prop.Supply, prop.Staff, prop.Still = true, 0, 0, true
	prop.Condition = 100
	if note := w.PlaceNote("club"); note != "" {
		t.Fatalf("a passer-by can see inside somebody else's casino: %q", note)
	}
	// But a boarded window is not a secret.
	prop.Condition = 35
	if note := w.PlaceNote("club"); note == "" {
		t.Fatal("a wrecked building looks perfectly fine from the street")
	}
	if !w.PlaceWarn("club") {
		t.Fatal("a wrecked building is not worth noticing")
	}
}

func TestEveryAddressCanBeAskedAboutItself(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	for _, l := range Locations {
		// It must never panic and never contradict itself.
		note := w.PlaceNote(l.ID)
		if w.PlaceWarn(l.ID) && note == "" {
			t.Fatalf("%s is a worry with nothing to say about it", l.Name)
		}
	}
}

// "Trading at 100% of what it could" was measured against staffing, supply and
// custom — and not against the condition of the building, which the clock uses
// to scale every dollar the place earns. So a laundry knocked down to 60%
// earned 169 a day where a sound one earned 305, and told the player it was
// trading at everything it could. The sentence promises a share of what the
// place could earn; the number has to be that share.
func TestWhatAPlaceIsTradingAtIsWhatItActuallyEarns(t *testing.T) {
	t.Parallel()
	// Gross, not net: the day's rent and wages are the same whatever state the
	// building is in, so measuring the cash left over would compare the wrong
	// thing and make a wrecked laundry look worse than it trades.
	earnings := func(condition int) int {
		w := New(4)
		w.Properties["laundry"].Owner = "player:1"
		w.Properties["laundry"].Condition = condition
		w.Player.Location = "laundry"
		before := w.Player.Earned
		w.Advance(1440)
		return w.Player.Earned - before
	}
	best := earnings(100)
	if best <= 0 {
		t.Fatal("a sound laundry earned nothing")
	}
	for _, condition := range []int{100, 80, 60, 30} {
		w := New(4)
		w.Properties["laundry"].Owner = "player:1"
		w.Properties["laundry"].Condition = condition
		w.Player.Location = "laundry"

		claimed := w.Trading("laundry")
		actual := float64(earnings(condition)) / float64(best)
		if diff := claimed - actual; diff > .06 || diff < -.06 {
			t.Errorf("at %d%% condition the place says it trades at %.0f%% and earns %.0f%% of what a sound one does",
				condition, claimed*100, actual*100)
		}
	}
	// And the words the player reads follow the same number.
	w := New(4)
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 80
	w.Player.Location = "laundry"
	if note := w.PlaceNote("laundry"); contains(note, "100%") {
		t.Fatalf("a laundry at 80%% condition says %q", note)
	}
	// And the sentence must not describe the number as a share of a maximum,
	// because strong custom carries a place past an ordinary day.
	w.Properties["laundry"].Condition = 100
	w.ShiftCustom("laundry", "a good week", 40)
	if w.Trading("laundry") <= 1 {
		t.Skip("custom no longer lifts a place past an ordinary day")
	}
	if note := w.PlaceNote("laundry"); contains(note, "of what it could") {
		t.Fatalf("a place doing better than usual reports a percentage of a ceiling, above the ceiling: %q", note)
	}
}
