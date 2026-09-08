package core

import "testing"

func armed(t *testing.T) *World {
	t.Helper()
	w := New(89)
	w.Player.Location = "docks"
	w.Player.Cash = 5000
	return w
}

func TestArmsComeOffABoatAndNowhereElse(t *testing.T) {
	w := armed(t)
	if w.ArmsReadiness("weapon") != "" {
		t.Fatal("the waterfront refused to sell:", w.ArmsReadiness("weapon"))
	}
	for _, elsewhere := range []string{"bar", "market", "club", "room"} {
		w.Player.Location = elsewhere
		if w.ArmsReadiness("weapon") == "" {
			t.Fatalf("%s was selling guns", elsewhere)
		}
		if err := w.BuyArms("weapon"); err == nil {
			t.Fatalf("bought a gun at %s", elsewhere)
		}
	}
}

func TestBuyingArmsWorksUpwardAndRunsOut(t *testing.T) {
	w := armed(t)
	w.Player.Cash = 100000
	for tier := 1; tier < len(weapons); tier++ {
		if err := w.BuyArms("weapon"); err != nil {
			t.Fatal(err)
		}
		if w.Player.Weapon != tier {
			t.Fatalf("expected tier %d, got %d", tier, w.Player.Weapon)
		}
	}
	if w.ArmsReadiness("weapon") == "" {
		t.Fatal("there was something better than the best")
	}
	if err := w.BuyArms("weapon"); err == nil {
		t.Fatal("bought past the top of the range")
	}
	// Being armed is itself a reason to be looked at.
	poor := armed(t)
	heat := poor.Player.Heat
	if err := poor.BuyArms("armour"); err != nil {
		t.Fatal(err)
	}
	if poor.Player.Heat <= heat {
		t.Fatal("carrying hardware drew no attention at all")
	}
}

func TestArmourReducesWhatViolenceCostsWithoutMakingItFree(t *testing.T) {
	bare, vested := New(91), New(91)
	vested.Player.Armour = len(armour) - 1
	for _, injury := range []int{10, 30, 65} {
		plain, protected := bare.Absorb(injury), vested.Absorb(injury)
		if protected >= plain {
			t.Fatalf("armour did not reduce a %d injury: %d against %d", injury, protected, plain)
		}
		if protected <= 0 {
			t.Fatalf("armour made a %d injury free", injury)
		}
	}
	// It never fully absorbs a serious one.
	if vested.Absorb(65) < 65/3 {
		t.Fatal("armour absorbed more than its floor allows")
	}
}

func TestAWeaponShiftsOddsWithoutGuaranteeingAnything(t *testing.T) {
	unarmed, carrying := New(93), New(93)
	carrying.Player.Weapon = len(weapons) - 1
	for _, w := range []*World{unarmed, carrying} {
		w.Player.Respect = 30
		w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 80}}
		w.Player.Location = "club"
	}
	if carrying.robberyOdds("club", carrying.OwnHands()) <= unarmed.robberyOdds("club", unarmed.OwnHands()) {
		t.Fatal("being armed did not improve the odds of a robbery")
	}
	if carrying.sabotageChance(carrying.faction("bellandi"), carrying.OwnHands()) <= unarmed.sabotageChance(unarmed.faction("bellandi"), unarmed.OwnHands()) {
		t.Fatal("being armed did not improve the odds of an attack")
	}
	if carrying.robberyOdds("club", carrying.OwnHands()) >= 1 {
		t.Fatal("a weapon made a robbery certain")
	}
}

func TestASearchTakesTheHardware(t *testing.T) {
	w := armed(t)
	w.Player.Weapon, w.Player.Armour = 2, 1
	if !w.SeizeArms() {
		t.Fatal("a search found nothing on an armed player")
	}
	if w.Player.Weapon != 0 || w.Player.Armour != 0 {
		t.Fatal("the player kept their hardware through a search")
	}
	if w.SeizeArms() {
		t.Fatal("a search of an unarmed player took something")
	}
	// A police raid takes it along with everything else.
	raided := armed(t)
	raided.Properties["laundry"].Owner = "player:1"
	raided.Player.Weapon, raided.Player.Armour = 3, 2
	raided.Player.Heat = 60
	raided.Raid()
	if raided.Player.Weapon != 0 || raided.Player.Armour != 0 {
		t.Fatal("a raid left the player armed")
	}
}

func TestArmourHelpsSurviveAnAttackAtHome(t *testing.T) {
	died := map[int]int{}
	for _, plate := range []int{0, 2} {
		for i := uint32(1); i <= 400; i++ {
			// Consecutive seeds give this generator nearly identical first
			// draws, so stride them the way cmd/simulate does.
			w := New(i * 2654435761)
			w.Player.Armour = plate
			w.Player.Contacts = 0
			w.Player.Security = 0
			w.Player.Location = w.Player.Home
			w.Attack(Plot{ID: ID(), Kind: "hit", Life: w.Life, Actor: "bellandi"})
			if !w.Player.Alive {
				died[plate]++
			}
		}
	}
	t.Logf("of 400 unwarned attacks: unarmoured died %d, armoured died %d", died[0], died[2])
	if died[2] >= died[0] {
		t.Fatal("armour did not help survive an attack at home")
	}
	if died[2] == 0 {
		t.Fatal("armour made an unwarned attack survivable every time")
	}
}
