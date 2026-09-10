package core

import "testing"

// A suit has to be a running cost rather than a purchase, or standing becomes
// something you buy once on day one and never think about again. These measure
// what keeping one actually costs across a campaign rather than asserting it.

// TestASuitIsAConsumable measures how long a tailored suit lasts if nobody
// looks after it. It has to stop being worth anything inside a season, or the
// pressing action and the laundry tie-in are decoration.
func TestASuitIsAConsumable(t *testing.T) {
	heavy(t)
	total, runs := 0, 200
	for i := 0; i < runs; i++ {
		w := New(uint32(i*2654435761 + 1))
		w.Player.Dress, w.Player.DressWear = 2, 100
		days := 0
		for w.Standing() > 0 && days < 200 {
			w.DressDay()
			days++
		}
		if days >= 200 {
			t.Fatalf("a suit survived %d days untouched", days)
		}
		total += days
	}
	average := total / runs
	if average < 20 || average > 45 {
		t.Fatalf("a neglected suit lasts %d days on average, which is not a season", average)
	}
	t.Logf("a tailored suit goes shabby after %d days of ordinary wear", average)
}

// TestKeepingUpAppearancesCostsRealMoney measures a year of a player who keeps
// a tailored suit presentable at home, against what the same year earns them at
// the high tables the suit unlocks. The suit should cost enough to notice and
// not enough to be a trap.
func TestKeepingUpAppearancesCostsRealMoney(t *testing.T) {
	heavy(t)
	pressings, heat := 0, 0
	const days = 180
	w := New(7)
	w.Player.Dress, w.Player.DressWear = 2, 100
	w.Player.Cash = 100000
	w.Player.Location = w.Player.Home
	for day := 0; day < days; day++ {
		before := w.Player.Heat
		w.DressDay()
		heat += w.Player.Heat - before
		w.Player.Heat = before // measured, not accumulated, so bribery is not in play
		if w.DressCondition() < 60 {
			if err := w.Press(w.Player.Home); err != nil {
				t.Fatal(err)
			}
			pressings++
		}
	}
	spent := pressings * PressCost
	if spent < 100 || spent > 600 {
		t.Fatalf("half a year of looking presentable cost $%d", spent)
	}
	if heat < days/2 {
		t.Fatalf("a tailored suit drew %d attention over %d days", heat, days)
	}
	t.Logf("half a year: %d pressings, $%d, %d police attention drawn", pressings, spent, heat)
}

// TestTheDoorIsOpenedByEitherHalfOfStanding checks that clothes and reputation
// are genuinely interchangeable at the door, so a new arrival with money has a
// route in and a known figure in working clothes is not shut out.
func TestTheDoorIsOpenedByEitherHalfOfStanding(t *testing.T) {
	heavy(t)
	high, _ := tableStake("high")

	reputation := New(3)
	reputation.Player.Location, reputation.Player.Cash = "club", 5000
	reputation.Player.Respect = HighTableStanding
	if reputation.TableReadiness("club", high) != "" {
		t.Fatal("a respected man in working clothes was turned away")
	}

	money := New(3)
	money.Player.Location, money.Player.Cash = "club", 5000
	money.Player.Dress, money.Player.DressWear = 3, 100
	if money.Player.Respect != 0 {
		t.Fatal("the fixture was not a nobody")
	}
	if money.TableReadiness("club", high) != "" {
		t.Fatal("a nobody in bespoke was turned away, so clothes buy nothing")
	}

	neither := New(3)
	neither.Player.Location, neither.Player.Cash = "club", 5000
	neither.Player.Dress, neither.Player.DressWear = 1, 100
	if neither.TableReadiness("club", high) == "" {
		t.Fatal("a nobody in an off-the-rack suit walked into the high tables")
	}
}
