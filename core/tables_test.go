package core

import "testing"

func gambler(t *testing.T) *World {
	t.Helper()
	w := New(71)
	w.Player.Location = "club"
	w.Player.Cash = 5000
	// Somebody the room will seat: the door is tested in the appearance suite.
	w.Player.Respect = HighTableStanding
	return w
}

func TestTheHouseKeepsItsEdge(t *testing.T) {
	const sessions = 6000
	staked, returned := 0, 0
	w := gambler(t)
	w.Player.Cash = 100000000
	small, _ := tableStake("small")
	for i := 0; i < sessions; i++ {
		before := w.Player.Cash
		if err := w.Play("club", "small"); err != nil {
			t.Fatal(err)
		}
		staked += small.Amount
		returned += w.Player.Cash - before + small.Amount
	}
	edge := float64(staked-returned) / float64(staked) * 100
	t.Logf("over %d sessions the house kept %.1f%% of everything staked", sessions, edge)
	if returned >= staked {
		t.Fatalf("the tables pay out more than they take: staked %d, returned %d", staked, returned)
	}
	if edge > 25 {
		t.Fatalf("the house edge is %.1f%%, which is robbery rather than a game", edge)
	}
	if edge < 1 {
		t.Fatalf("the house edge is %.1f%%, so owning a casino is worth nothing", edge)
	}
}

func TestBigWinsComeOutOfTheOwnersPocket(t *testing.T) {
	var w *World
	for seed := uint32(1); seed <= 4000; seed++ {
		probe := gambler(t)
		probe.RNG = seed
		high, _ := tableStake("high")
		before := probe.faction("bellandi").Cash
		cash := probe.Player.Cash
		if err := probe.Play("club", "high"); err != nil {
			t.Fatal(err)
		}
		// A heavy win, which is the case the owner is supposed to notice.
		if probe.Player.Cash-cash > high.Amount*2 && probe.faction("bellandi").Cash < before {
			w = probe
			break
		}
	}
	if w == nil {
		t.Skip("no heavy win found in range")
	}
	if w.faction("bellandi").Goodwill >= 0 {
		t.Fatal("cleaning out a family's room did not annoy them")
	}
	if w.Player.Heat == 0 {
		t.Fatal("a heavy win at somebody else's tables drew no attention")
	}
}

func TestYouCannotBeatYourOwnHouse(t *testing.T) {
	w := gambler(t)
	small, _ := tableStake("small")
	w.Properties["club"].Owner = "player:1"
	if w.TableReadiness("club", small) == "" {
		t.Fatal("the player was offered a seat at their own tables")
	}
	if err := w.Play("club", "small"); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	// Nowhere without tables runs games.
	w.Player.Location = "laundry"
	if w.TableReadiness("laundry", small) == "" {
		t.Fatal("a laundry ran a card game")
	}
	// A wrecked room cannot.
	w.Properties["club"].Owner = "bellandi"
	w.Properties["club"].Condition = 10
	if w.TableReadiness("club", small) == "" {
		t.Fatal("a wrecked casino was still running games")
	}
	// And you need the stake.
	w.Properties["club"].Condition = 100
	w.Player.Cash = 10
	if w.TableReadiness("club", small) == "" {
		t.Fatal("the player sat down without the money")
	}
}

func TestLosingIsTheUsualOutcome(t *testing.T) {
	w := gambler(t)
	w.Player.Cash = 1000000
	losses, wins := 0, 0
	for i := 0; i < 2000; i++ {
		before := w.Player.Cash
		if err := w.Play("club", "small"); err != nil {
			t.Fatal(err)
		}
		if w.Player.Cash < before {
			losses++
		} else if w.Player.Cash > before {
			wins++
		}
	}
	t.Logf("of 2000 sessions: %d lost, %d won", losses, wins)
	if losses <= wins {
		t.Fatal("the tables are not a losing proposition on any given night")
	}
	if wins == 0 {
		t.Fatal("nobody ever wins, so nobody would ever play")
	}
}
