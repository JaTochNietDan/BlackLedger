package core

import "testing"

// legacyWorld reproduces the shape of a campaign saved before the city had
// holdings, people or quarrels: the player's own estate, a dead protagonist's
// businesses, one family holding a single property and the other holding none.
func legacyWorld() *World {
	w := New(27)
	w.Version = 2
	w.Life = 2
	w.Conflicts = nil
	w.Properties["estate"].Owner = "player:2"
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "former:Alex Varga"
	}
	for _, id := range []string{"docks", "market", "bar", "room", "apartment"} {
		w.Properties[id].Owner = "independent"
		w.Properties[id].Income = 0
	}
	w.Properties["club"].Owner = "bellandi"
	for i := range w.Factions {
		w.Factions[i].Peak = 0
	}
	// Nobody had an organization, a rank or a place to be found.
	for i := range w.NPCs {
		w.NPCs[i].Faction, w.NPCs[i].Rank, w.NPCs[i].Location = "", 0, ""
		w.NPCs[i].Ambition, w.NPCs[i].Skill = 0, 0
	}
	w.NPCs = w.NPCs[:4] // the four the old build seeded
	return w
}

func TestMigrationNeverTakesPropertyThatIsHeld(t *testing.T) {
	t.Parallel()
	w := legacyWorld()
	w.MigrateLivingWorld()
	if !w.Own("estate") {
		t.Fatal("the player's residence was taken away")
	}
	for _, id := range []string{"laundry", "garage", "casino"} {
		if w.Properties[id].Owner != "former:Alex Varga" {
			t.Fatalf("%s was taken from the previous protagonist's estate: now %q", id, w.Properties[id].Owner)
		}
	}
}

func TestMigrationGivesEveryOrganizationSomethingToLose(t *testing.T) {
	t.Parallel()
	w := legacyWorld()
	w.MigrateLivingWorld()
	for _, f := range w.Factions {
		holdings := w.FamilyHoldings(f.ID)
		if len(holdings) < 2 {
			t.Fatalf("%s holds %v; it cannot be fought over or split", f.Name, holdings)
		}
		for _, id := range holdings {
			if w.Properties[id].Income <= 0 {
				t.Fatalf("%s holds %s, which earns it nothing", f.Name, id)
			}
		}
		if f.Peak <= 0 {
			t.Fatalf("%s has no strength to recover toward", f.Name)
		}
	}
}

func TestMigrationGivesOrganizationsPeople(t *testing.T) {
	t.Parallel()
	w := legacyWorld()
	w.MigrateLivingWorld()
	for _, f := range w.Factions {
		members := w.Members(f.ID)
		if len(members) < 3 {
			t.Fatalf("%s has %d people; nobody could replace its leader", f.Name, len(members))
		}
		if members[0].Name != f.Leader || members[0].Rank != RankLeader {
			t.Fatalf("%s: the leader is not the strongest standing", f.Name)
		}
	}
	for _, n := range w.People() {
		if n.Location == "" || n.Voice == "" {
			t.Fatalf("%s is not really anywhere or has no voice", n.Name)
		}
	}
	// Established rivals already dislike each other.
	if len(w.PublicConflicts()) == 0 {
		t.Fatal("the families have no history with each other after migrating")
	}
}

func TestMigrationIsIdempotent(t *testing.T) {
	t.Parallel()
	once := legacyWorld()
	once.MigrateLivingWorld()
	twice := legacyWorld()
	twice.MigrateLivingWorld()
	twice.MigrateLivingWorld()
	if len(twice.NPCs) != len(once.NPCs) {
		t.Fatalf("running the migration twice added people: %d then %d", len(once.NPCs), len(twice.NPCs))
	}
	for _, f := range twice.Factions {
		if len(twice.FamilyHoldings(f.ID)) != len(once.FamilyHoldings(f.ID)) {
			t.Fatal("running the migration twice changed holdings")
		}
	}
	if len(twice.Conflicts) != len(once.Conflicts) {
		t.Fatal("running the migration twice added quarrels")
	}
}

func TestMigrationDoesNotChargeOrAdvanceThePlayer(t *testing.T) {
	t.Parallel()
	w := legacyWorld()
	cash, minute, respect, revision := w.Player.Cash, w.Minute, w.Player.Respect, w.Revision
	w.MigrateLivingWorld()
	if w.Player.Cash != cash || w.Minute != minute || w.Player.Respect != respect || w.Revision != revision {
		t.Fatal("migration changed the player's position")
	}
}

func TestAMigratedCampaignCanActuallyFightAWar(t *testing.T) {
	t.Parallel()
	// Measured across seeds rather than one: a single campaign may legitimately
	// stalemate, and tying the assertion to one random stream makes the test
	// fail whenever an unrelated system draws from it.
	seizures, kept := 0, 0
	for seed := 0; seed < 40; seed++ {
		w := legacyWorld()
		w.RNG, w.WorldRNG = uint32(seed*7919+1), uint32(seed*104729+1)
		w.MigrateLivingWorld()
		w.Antagonize("bellandi", "russo", 100)
		for day := 0; day < 120; day++ {
			before := map[string]string{}
			for id, prop := range w.Properties {
				before[id] = prop.Owner
			}
			w.Minute += 720
			w.FactionTurn()
			w.Minute += 720
			w.FactionTurn()
			w.FamilyDay()
			for id, prop := range w.Properties {
				if before[id] != prop.Owner {
					seizures++
				}
			}
		}
		// The player's own holdings are never spoils of a war between families.
		if w.Own("estate") {
			kept++
		}
	}
	if seizures == 0 {
		t.Fatal("forty migrated cities fought for four months and nothing ever changed hands")
	}
	if kept != 40 {
		t.Fatalf("a war between families took the player's residence in %d of 40 cities", 40-kept)
	}
	t.Logf("across 40 migrated cities at war: %d holdings changed hands", seizures)
}
