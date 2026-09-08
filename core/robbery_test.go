package core

import "testing"

func robberWorld(t *testing.T) *World {
	t.Helper()
	w := New(31)
	w.Player.Location = "club" // Bellandi premises, and they trade
	w.Player.Respect = 30
	return w
}

func TestRobberyNeedsSomethingWorthTakingThatIsNotYours(t *testing.T) {
	w := robberWorld(t)
	if w.RobberyReadiness("club") != "" {
		t.Fatal("premises with takings were refused:", w.RobberyReadiness("club"))
	}
	if w.RobberyReadiness("room") == "" {
		t.Fatal("a rented room with no takings was offered as a robbery")
	}
	w.Properties["laundry"].Owner = "player:1"
	if w.RobberyReadiness("laundry") == "" {
		t.Fatal("the player was offered a robbery of their own business")
	}
	if err := w.Rob("laundry"); err == nil {
		t.Fatal("the command ignored a requirement the action reports")
	}
	w.Player.Health = 20
	if w.RobberyReadiness("club") == "" {
		t.Fatal("a badly hurt player was sent to rob premises")
	}
}

func TestASuccessfulRobberyPaysAndIsNoticed(t *testing.T) {
	var w *World
	for seed := uint32(1); seed <= 200; seed++ {
		probe := robberWorld(t)
		probe.RNG = seed
		if probe.Random() < probe.robberyOdds("club") {
			w = robberWorld(t)
			w.RNG = seed
			break
		}
	}
	if w == nil {
		t.Skip("no succeeding seed")
	}
	cash, heat := w.Player.Cash, w.Player.Heat
	owner := *w.faction("bellandi")
	if err := w.Rob("club"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash <= cash {
		t.Fatal("a successful robbery paid nothing")
	}
	if w.Player.Heat <= heat {
		t.Fatal("a robbery drew no police attention")
	}
	after := w.faction("bellandi")
	if after.Goodwill >= owner.Goodwill {
		t.Fatal("robbing a family did not worsen standing with them")
	}
	if after.Cash >= owner.Cash {
		t.Fatal("the money came from nowhere; the owner lost nothing")
	}
	answered := false
	for _, p := range w.Plots {
		if p.Actor == "bellandi" && p.Life == w.Life {
			answered = true
		}
	}
	if !answered {
		t.Fatal("a family robbed of its takings did not answer")
	}
}

func TestAFailedRobberyHurts(t *testing.T) {
	var w *World
	for seed := uint32(1); seed <= 400; seed++ {
		probe := robberWorld(t)
		probe.RNG = seed
		if probe.Random() >= probe.robberyOdds("club") {
			w = robberWorld(t)
			w.RNG = seed
			break
		}
	}
	if w == nil {
		t.Skip("no failing seed")
	}
	cash, health := w.Player.Cash, w.Player.Health
	if err := w.Rob("club"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash {
		t.Fatal("a failed robbery still paid out")
	}
	if w.Player.Health >= health {
		t.Fatal("a failed robbery cost nothing")
	}
}

func TestTheCityRobsThePlayerBack(t *testing.T) {
	robbed, named, anonymous := 0, 0, 0
	for i := uint32(1); i <= 400; i++ {
		w := New(i * 2654435761)
		w.Player.Cash = 4000
		w.Player.Stock = map[string]int{"moonshine": 20}
		w.Player.Contacts = 0
		before := w.Carrying()
		w.ConsiderRobbery()
		if w.Carrying() < before {
			robbed++
			anonymous++
		}
	}
	for i := uint32(1); i <= 400; i++ {
		w := New(i * 2654435761)
		w.Player.Cash = 4000
		w.Player.Stock = map[string]int{"moonshine": 20}
		w.Player.Contacts = 3
		before := w.Carrying()
		w.ConsiderRobbery()
		if w.Carrying() < before {
			for _, r := range w.History {
				if r.Kind == "danger" {
					for _, n := range w.NPCs {
						if n.Name != "" && contains(r.Text, n.Name) {
							named++
						}
					}
				}
			}
		}
	}
	t.Logf("of 400 exposed players: %d were robbed; with contacts, %d learned a name", robbed, named)
	if robbed == 0 {
		t.Fatal("nobody in the city ever robs the player")
	}
	if robbed == 400 {
		t.Fatal("carrying goods means being robbed every single time")
	}
	if named == 0 {
		t.Fatal("an information network never produces the name of who robbed you")
	}
	if anonymous == 0 {
		t.Fatal("a player with no contacts always learns who did it")
	}
}

func contains(haystack, needle string) bool {
	if len(needle) == 0 || len(haystack) < len(needle) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func TestCarryingNothingIsSafer(t *testing.T) {
	empty, loaded := New(41), New(41)
	empty.Player.Cash = 100
	loaded.Player.Cash = 100
	loaded.Player.Stock = map[string]int{"moonshine": 30}
	if empty.robberyExposure() >= loaded.robberyExposure() {
		t.Fatalf("carrying stock is no more dangerous than carrying none: %.3f against %.3f",
			empty.robberyExposure(), loaded.robberyExposure())
	}
	// A loyal, available crew member makes the player a worse mark.
	guarded := New(41)
	guarded.Player.Cash = 100
	guarded.Player.Stock = map[string]int{"moonshine": 30}
	guarded.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
	if guarded.robberyExposure() >= loaded.robberyExposure() {
		t.Fatal("having somebody with you did not make you a harder mark")
	}
}

func TestNobodyIsRobbedByANobody(t *testing.T) {
	w := New(43)
	for i := 0; i < 200; i++ {
		thief := w.robber()
		if thief == nil {
			t.Fatal("no living person was available to do it")
		}
		if thief.Name == "" {
			t.Fatal("an unnamed person robbed the player")
		}
		if thief.Rank >= RankLeader {
			t.Fatal("the head of an organization did this personally")
		}
		if thief.Dead {
			t.Fatal("a dead person robbed the player")
		}
	}
}
