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
	w := buyer(t)
	professional, _ := contractTier("professional")
	fee := w.ContractPrice("vittorio", professional)
	cash := w.Player.Cash
	if err := w.Commission("vittorio", "professional"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Cash != cash-fee {
		t.Fatalf("fee was not taken: cash %d, expected %d", w.Player.Cash, cash-fee)
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
	before := w.Player.Cash
	if err := w.Commission("elena", "professional"); err == nil {
		t.Fatal("commissioned a killing of somebody already dead")
	}
	if w.Player.Cash != before {
		t.Fatal("a refused commission still took the money")
	}
	if err := w.Commission("vittorio", "nobody-like-that"); err == nil {
		t.Fatal("hired a class of person who does not exist")
	}
}

func TestAContractIsNeverVisibleToThePlayer(t *testing.T) {
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
	rates := map[string]int{}
	for _, tier := range []string{"cheap", "specialist"} {
		for seed := uint32(1); seed <= 300; seed++ {
			w := New(seed)
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
	traced, escaped := 0, 0
	for seed := uint32(1); seed <= 400; seed++ {
		w := New(seed)
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
	w := New(19)
	seen := map[string]bool{}
	for i := 0; i < 40; i++ {
		seen[w.killingMethod(w.NPC("vittorio"))] = true
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
	booked, killed := 0, 0
	for seed := uint32(1); seed <= 300; seed++ {
		w := New(seed)
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
