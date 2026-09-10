package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func buyer(t *testing.T) *World {
	t.Helper()
	w := New(17)
	w.Player.Location = "market"
	w.Player.Contacts = 2
	w.Player.Cash = 20000
	return w
}

func TestWhatItCostsDependsOnWhoTheyAre(t *testing.T) {
	t.Parallel()
	w := buyer(t)
	professional, _ := contractTier("professional")
	leader := w.ContractPrice("vittorio", professional)
	var soldier int
	for _, n := range w.Members("bellandi") {
		if n.Rank == RankSoldier {
			soldier = w.ContractPrice(n.ID, professional)
		}
	}
	if soldier == 0 {
		t.Fatal("no soldier to price")
	}
	if leader <= soldier {
		t.Fatalf("killing the head of a family costs no more than killing a soldier: %d against %d", leader, soldier)
	}
	// Better people cost more than desperate ones.
	cheap, _ := contractTier("cheap")
	specialist, _ := contractTier("specialist")
	if w.ContractPrice("vittorio", specialist) <= w.ContractPrice("vittorio", cheap) {
		t.Fatal("a specialist costs no more than someone who needs the money")
	}
}

func TestCommissioningTakesTheMoneyAndCommitsNothingYet(t *testing.T) {
	t.Parallel()
	w := buyer(t)
	professional, _ := contractTier("professional")
	fee := w.ContractPrice("vittorio", professional)
	cash := w.Player.Cash
	if err := w.Commission("vittorio", "professional"); err != nil {
		t.Fatal(err)
	}
	// The fee is charged by the choice that reaches Commission, carrying it as
	// its cost, so booking one here leaves cash alone. Charging in both places
	// took the money twice and refused any arrangement the player could only
	// just afford.
	if w.Player.Cash != cash {
		t.Fatalf("Commission charged the fee a second time: cash %d, expected %d", w.Player.Cash, cash)
	}
	if w.Contracts[0].Fee != fee {
		t.Fatalf("the contract recorded a fee of %d, expected %d", w.Contracts[0].Fee, fee)
	}
	if len(w.Contracts) != 1 {
		t.Fatal("no contract was booked")
	}
	if w.Contracts[0].Due <= w.Minute {
		t.Fatal("a hit resolved the moment it was paid for")
	}
	if w.NPC("vittorio").Dead {
		t.Fatal("the target died before the contract was due")
	}
	// A name that cannot be reached is refused, and costs nothing.
	w.Kill("elena", "Something else.")
	if err := w.Commission("elena", "professional"); err == nil {
		t.Fatal("commissioned a killing of somebody already dead")
	}
	if err := w.Commission("vittorio", "nobody-like-that"); err == nil {
		t.Fatal("hired a class of person who does not exist")
	}
}

func TestAContractIsNeverVisibleToThePlayer(t *testing.T) {
	t.Parallel()
	w := buyer(t)
	if err := w.Commission("vittorio", "specialist"); err != nil {
		t.Fatal(err)
	}
	public := w.Public()
	if _, listed := public["contracts"]; listed {
		t.Fatal("the public projection lists commissioned killings")
	}
	body, err := json.Marshal(public)
	if err != nil {
		t.Fatal(err)
	}
	// The offer to arrange one is public; who is marked, by whom and when is not.
	if strings.Contains(string(body), w.Contracts[0].ID) {
		t.Fatal("a contract id reached the public projection")
	}
	if strings.Contains(string(body), fmt.Sprintf("\"due\":%d", w.Contracts[0].Due)) {
		t.Fatal("a contract's timing reached the public projection")
	}
	// The target is not flagged as marked anywhere the player can read.
	for _, n := range public["npcs"].([]*NPC) {
		if n.ID == "vittorio" && n.Dead {
			t.Fatal("the target is already dead in the projection")
		}
	}
}

func TestABetterHitmanSucceedsMoreOften(t *testing.T) {
	t.Parallel()
	rates := map[string]int{}
	for _, tier := range []string{"cheap", "specialist"} {
		for i := uint32(1); i <= 300; i++ {
			w := New(i * 2654435761)
			w.Player.Cash = 100000
			w.Contracts = []Contract{{ID: ID(), Life: w.Life, Target: "vittorio", Tier: tier, Payer: "player", Due: w.Minute}}
			w.ResolveContracts()
			if w.NPC("vittorio").Dead {
				rates[tier]++
			}
		}
	}
	t.Logf("of 300 attempts on a family head: cheap killed %d, a specialist killed %d", rates["cheap"], rates["specialist"])
	if rates["specialist"] <= rates["cheap"] {
		t.Fatal("paying for a better hitman bought nothing")
	}
	if rates["cheap"] == 0 {
		t.Fatal("the cheapest option never works, so it is not a choice")
	}
	if rates["specialist"] == 300 {
		t.Fatal("a specialist never fails, so there is no risk to weigh")
	}
}

func TestAFailedAttemptCanBeTracedBack(t *testing.T) {
	t.Parallel()
	traced, escaped := 0, 0
	for i := uint32(1); i <= 400; i++ {
		w := New(i * 2654435761)
		w.Contracts = []Contract{{ID: ID(), Life: w.Life, Target: "vittorio", Tier: "cheap", Payer: "player", Due: w.Minute}}
		heat, goodwill := w.Player.Heat, w.faction("bellandi").Goodwill
		w.ResolveContracts()
		if w.NPC("vittorio").Dead {
			continue
		}
		if w.Player.Heat > heat+10 && w.faction("bellandi").Goodwill < goodwill {
			traced++
		} else {
			escaped++
		}
	}
	t.Logf("of failed cheap attempts: %d were traced back, %d got away clean", traced, escaped)
	if traced == 0 {
		t.Fatal("a captured hitman never gives up who paid them")
	}
	if escaped == 0 {
		t.Fatal("a failed attempt is always traced, so the tier's capture rate means nothing")
	}
}

func TestAKillingIsDescribedByWhereItHappened(t *testing.T) {
	t.Parallel()
	w := New(19)
	seen := map[string]bool{}
	for i := 0; i < 40; i++ {
		tier, _ := contractTier("professional")
		seen[w.killingMethod(w.NPC("vittorio"), tier)] = true
	}
	if len(seen) < 2 {
		t.Fatal("every killing at the same place reads identically")
	}
	for text := range seen {
		if !strings.Contains(text, "The Monarch") {
			t.Fatalf("a killing did not say where it happened: %q", text)
		}
	}
}

func TestKillingALeaderThroughAContractStillPromotesSomebody(t *testing.T) {
	t.Parallel()
	w := New(23)
	deputy := w.Members("bellandi")[1].Name
	w.Contracts = []Contract{{ID: ID(), Life: w.Life, Target: "vittorio", Tier: "specialist", Payer: "player", Due: w.Minute}}
	for attempt := 0; attempt < 60 && !w.NPC("vittorio").Dead; attempt++ {
		w.Contracts = []Contract{{ID: ID(), Life: w.Life, Target: "vittorio", Tier: "specialist", Payer: "player", Due: w.Minute}}
		w.ResolveContracts()
	}
	if !w.NPC("vittorio").Dead {
		t.Skip("no successful attempt in range")
	}
	if w.faction("bellandi").Leader != deputy {
		t.Fatalf("a family whose head was killed did not promote its deputy: leader is %q", w.faction("bellandi").Leader)
	}
}

func TestOrganizationsBuyTheSameServiceThePlayerCan(t *testing.T) {
	t.Parallel()
	booked, killed := 0, 0
	for i := uint32(1); i <= 300; i++ {
		w := New(i * 2654435761)
		w.Antagonize("bellandi", "russo", 100)
		if w.Conflict("bellandi", "russo").State != "war" {
			t.Fatal("the test world is not at war")
		}
		before := len(w.Contracts)
		for turn := 0; turn < 20; turn++ {
			w.ConsiderFactionContracts()
		}
		if len(w.Contracts) > before {
			booked++
			for _, c := range w.Contracts {
				if c.Payer == "player" {
					t.Fatal("a faction contract was recorded as the player's")
				}
			}
			w.Minute += 2000
			w.ResolveContracts()
			for _, n := range w.NPCs {
				if n.Dead {
					killed++
					break
				}
			}
		}
	}
	t.Logf("of 300 cities at war: %d saw an organization commission a killing, %d of those landed", booked, killed)
	if booked == 0 {
		t.Fatal("organizations never buy what the player can buy")
	}
	if killed == 0 {
		t.Fatal("no commissioned killing by an organization ever succeeded")
	}
}

func TestAnOrganizationWillNotSpendMoneyItDoesNotHave(t *testing.T) {
	t.Parallel()
	w := New(29)
	w.Antagonize("bellandi", "russo", 100)
	for i := range w.Factions {
		w.Factions[i].Cash = 0
	}
	for turn := 0; turn < 200; turn++ {
		w.ConsiderFactionContracts()
	}
	if len(w.Contracts) != 0 {
		t.Fatal("a penniless organization commissioned a killing")
	}
	for _, f := range w.Factions {
		if f.Cash < 0 {
			t.Fatal("an organization went into debt to pay a hitman")
		}
	}
}

func TestThePlayerCanSeeWhatTheyPaidFor(t *testing.T) {
	t.Parallel()
	w := buyer(t)
	if len(w.PendingArrangements()) != 0 {
		t.Fatal("arrangements existed before any were made")
	}
	if err := w.Commission("vittorio", "professional"); err != nil {
		t.Fatal(err)
	}
	pending := w.PendingArrangements()
	if len(pending) != 1 {
		t.Fatalf("the player cannot see the arrangement they paid for: %d listed", len(pending))
	}
	if pending[0]["target"] != "Vittorio Bellandi" {
		t.Fatal("the arrangement does not name who it is for")
	}
	if pending[0]["paid"].(int) <= 0 {
		t.Fatal("the arrangement does not say what it cost")
	}
	// When it is due is never revealed.
	for key := range pending[0] {
		if key == "due" || key == "minute" {
			t.Fatal("the player was told exactly when it happens")
		}
	}
	// Somebody else's arrangement is not the player's business.
	w.Contracts = append(w.Contracts, Contract{ID: ID(), Life: w.Life, Target: "elena", Tier: "cheap", Payer: "bellandi", Due: w.Minute + 100})
	if len(w.PendingArrangements()) != 1 {
		t.Fatal("a faction's arrangement was shown to the player")
	}
	// It clears once it resolves.
	w.Contracts[0].Due = w.Minute
	w.ResolveContracts()
	for _, a := range w.PendingArrangements() {
		if a["target"] == "Vittorio Bellandi" {
			t.Fatal("a settled arrangement is still listed as pending")
		}
	}
}

func TestAResolutionThePlayerPaidForSaysSo(t *testing.T) {
	t.Parallel()
	for _, want := range []struct{ tier, phrase string }{{"specialist", "Your arrangement is settled"}, {"cheap", "Your arrangement failed"}} {
		found := false
		for seed := uint32(1); seed <= 400 && !found; seed++ {
			w := New(seed)
			w.Contracts = []Contract{{ID: ID(), Life: w.Life, Target: "vittorio", Tier: want.tier, Payer: "player", Fee: 900, Due: w.Minute}}
			w.ResolveContracts()
			for _, r := range w.History {
				if r.Title == want.phrase {
					found = true
					if !contains(r.Text, "900") {
						t.Fatalf("%q did not say what it cost: %q", want.phrase, r.Text)
					}
				}
			}
		}
		if !found {
			t.Fatalf("no outcome ever reported %q to the player", want.phrase)
		}
	}
}

// A price is charged exactly once. The command layer pays what an action or a
// choice declares, so a handler that also pays takes it twice and refuses
// anything the player could only just afford.
func TestNothingIsChargedTwice(t *testing.T) {
	t.Parallel()
	for _, probe := range []struct {
		name    string
		place   string
		prepare func(*World)
		kind    string
	}{
		{"arms", "docks", func(w *World) { w.Player.Cash = 240 }, "arms:weapon"},
		{"tables", "club", func(w *World) { w.Player.Cash = 60 }, "play"},
		{"bribe", "market", func(w *World) { w.Player.Cash = 100000; w.Player.Heat = 20 }, "bribe"},
		{"contraband", "market", func(w *World) { w.Player.Cash = 210 }, "buy:moonshine"},
		{"launder", "laundry", func(w *World) {
			w.Properties["laundry"].Owner = "player:1"
			w.Player.Cash, w.Player.Heat = 100000, 30
		}, "launder"},
	} {
		t.Run(probe.name, func(t *testing.T) {
			w := New(307)
			probe.prepare(w)
			w.Player.Location = probe.place
			var offered *Action
			for _, a := range w.Actions(probe.place) {
				if a.ID == probe.kind {
					copy := a
					offered = &copy
				}
			}
			if offered == nil {
				t.Fatalf("%s was not offered at %s", probe.kind, probe.place)
			}
			if offered.Disabled {
				t.Fatalf("%s was refused before it could be tried: %s", probe.kind, offered.Reason)
			}
			next, err := Execute(w, Command{Revision: w.Revision, Kind: probe.kind, Target: probe.place})
			if err != nil {
				t.Fatalf("an affordable %s was refused: %v", probe.kind, err)
			}
			if next.Player.Cash < 0 {
				t.Fatal("the player was charged into debt")
			}
		})
	}
}

// Naming somebody used to offer every living person in the city except the
// player's own crew: in a city of fifty that rendered as fifty identical rows
// in one scene, including a laundress the player had never heard of. A name is
// something you have a reason to say.

func TestYouCanOnlyNameSomebodyYouHaveAReasonToName(t *testing.T) {
	t.Parallel()
	w, member := testator(t)
	w.Populate()
	all := len(w.People())
	if all < 20 {
		t.Fatalf("only %d people in the city; this is not measuring what it claims to", all)
	}
	targets := w.ContractTargets()
	if len(targets) == 0 {
		t.Fatal("there is nobody in this city the player could name")
	}
	if len(targets) >= all/2 {
		t.Fatalf("%d of %d people are on the list, which is still a wall", len(targets), all)
	}

	named := map[string]bool{}
	for _, n := range targets {
		named[n.ID] = true
		if n.Faction == w.PlayerOrganizationID() {
			t.Fatalf("%s works for the player and is on the list", n.Name)
		}
	}
	// Somebody the player has never had a reason to hear of is not on it.
	stranger := ""
	for _, n := range w.Civilians() {
		if !w.Known(n) && n.Sore == 0 && !IsOfficial(n.ID) && !w.isRoleHolder(n) && n.Rank < RankLeader {
			stranger = n.ID
			break
		}
	}
	if stranger == "" {
		t.Skip("everybody in this city is known")
	}
	if named[stranger] {
		t.Fatal("a complete stranger can have a price put on them")
	}

	// But give that stranger a reason and they are the first name on the list.
	w.Aggrieve(stranger, 45, "what was done to them over money")
	targets = w.ContractTargets()
	if len(targets) == 0 || targets[0].ID != stranger {
		got := "nobody"
		if len(targets) > 0 {
			got = targets[0].Name
		}
		t.Fatalf("the first name offered is %s, not the man carrying something against the player", got)
	}
	_ = member
}

func TestTheHeadsOfOrganizationsAreAlwaysNameable(t *testing.T) {
	t.Parallel()
	w := proprietor(t)
	targets := map[string]bool{}
	for _, n := range w.ContractTargets() {
		targets[n.ID] = true
	}
	for _, f := range w.Factions {
		for _, n := range w.People() {
			if n.Name == f.Leader && !targets[n.ID] {
				t.Fatalf("%s runs %s and cannot be named", n.Name, f.Name)
			}
		}
	}
}
