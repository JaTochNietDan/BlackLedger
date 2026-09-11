package core

import (
	"fmt"
	"strings"
	"testing"
)

func armed(t *testing.T) *World {
	t.Helper()
	w := New(89)
	w.Player.Location = "docks"
	w.Player.Cash = 5000
	return w
}

func TestArmsComeOffABoatAndNowhereElse(t *testing.T) {
	t.Parallel()
	w := armed(t)
	if w.ArmsReadiness("weapon", 1) != "" {
		t.Fatal("the waterfront refused to sell:", w.ArmsReadiness("weapon", 1))
	}
	for _, elsewhere := range []string{"bar", "market", "club", "room"} {
		w.Player.Location = elsewhere
		if w.ArmsReadiness("weapon", 1) == "" {
			t.Fatalf("%s was selling guns", elsewhere)
		}
		if err := w.BuyArms("weapon", 1); err == nil {
			t.Fatalf("bought a gun at %s", elsewhere)
		}
	}
}

// Everything is on the counter, each at its own price, and none of it has to be
// climbed to. The dock offered exactly one of each — the next one up — so a
// Thompson meant buying a revolver and a shotgun first and throwing both away:
// "when buying guns/armor and whatnot I don't think you should have to progress
// through them, you should be able to buy any of them at any time, you don't
// need to go through some sort of upgrade cycle."
func TestEveryGunIsOnTheCounterAtOnce(t *testing.T) {
	t.Parallel()
	w := armed(t)
	w.Player.Cash = 100000
	// Every one of them is offered, and the best of them can be had first.
	offered := map[string]bool{}
	for _, a := range w.Actions("docks") {
		if strings.HasPrefix(a.ID, "arms:") {
			offered[a.ID] = !a.Disabled
		}
	}
	for _, kind := range []string{"weapon", "armour"} {
		for _, arm := range Armaments(kind) {
			id := fmt.Sprintf("arms:%s:%d", kind, arm.Tier)
			ready, ok := offered[id]
			if !ok {
				t.Fatalf("%s is not on the counter at all", arm.Label)
			}
			if !ready {
				t.Fatalf("%s is on the counter and refused to somebody with $100,000", arm.Label)
			}
		}
	}
	top := Armaments("weapon")[len(Armaments("weapon"))-1]
	if err := w.BuyArms("weapon", top.Tier); err != nil {
		t.Fatalf("could not buy %s without buying everything under it first: %v", top.Label, err)
	}
	if w.Player.Weapon != top.Tier {
		t.Fatalf("bought %s and are carrying tier %d", top.Label, w.Player.Weapon)
	}
	// And paid for that one rather than for the climb.
	if spent := 100000 - w.Player.Cash; spent != top.Cost {
		t.Fatalf("%s cost $%d and $%d was spent", top.Label, top.Cost, spent)
	}
}

// What you cannot do is pay for something worse than what you are carrying.
// Nothing in this city rewards carrying less gun, so that would be a trap
// rather than a choice — and the card says so rather than going quiet.
func TestYouCannotPayForALesserGunThanYouCarry(t *testing.T) {
	t.Parallel()
	w := armed(t)
	w.Player.Cash = 100000
	guns := Armaments("weapon")
	top, lesser := guns[len(guns)-1], guns[0]
	if err := w.BuyArms("weapon", top.Tier); err != nil {
		t.Fatal(err)
	}
	cash := w.Player.Cash
	if err := w.BuyArms("weapon", lesser.Tier); err == nil {
		t.Fatalf("bought %s while carrying %s", lesser.Label, top.Label)
	}
	if w.ArmsReadiness("weapon", lesser.Tier) == "" {
		t.Fatal("the card for a lesser gun says nothing about why it is refused")
	}
	if w.ArmsReadiness("weapon", top.Tier) != "You are carrying it" {
		t.Fatalf("the one being carried says %q", w.ArmsReadiness("weapon", top.Tier))
	}
	if w.Player.Cash != cash || w.Player.Weapon != top.Tier {
		t.Fatal("the refusal still took the money or the gun")
	}
	// Being armed is itself a reason to be looked at.
	poor := armed(t)
	heat := poor.Player.Heat
	if err := poor.BuyArms("armour", 1); err != nil {
		t.Fatal(err)
	}
	if poor.Player.Heat <= heat {
		t.Fatal("carrying hardware drew no attention at all")
	}
}

func TestArmourReducesWhatViolenceCostsWithoutMakingItFree(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// A gun bought for one of your own was the only thing in this city that carried
// no risk at all. The search took what was in your coat and left what was in
// theirs, so arming the man you send was strictly better than arming yourself
// in the one way that matters — and that was a hole opened the same night the
// counter started selling for them.
//
// They are yours, they are in the room, and the room is being turned out.
func TestASearchReachesWhoeverIsStandingWithYou(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 40, 20000
	w.Player.Weapon, w.Player.Armour = 2, 1

	here := w.Holder("driver")
	away := w.Holder("fixer")
	if here == nil || away == nil {
		t.Fatal("this city has nobody to sign on")
	}
	here.Faction, here.Location, here.Weapon = w.PlayerOrganizationID(), w.Player.Location, 3
	away.Faction, away.Location, away.Weapon = w.PlayerOrganizationID(), "docks", 3
	if w.Player.Location == "docks" {
		t.Fatal("both of them are standing in the same room, so this measures nothing")
	}

	if !w.SeizeArms() {
		t.Fatal("a search of somebody carrying a shotgun found nothing")
	}
	if w.Player.Weapon != 0 || w.Player.Armour != 0 {
		t.Fatal("the search left the player armed")
	}
	if here.Weapon != 0 {
		t.Fatalf("%s was standing in the room and kept their gun", here.Name)
	}
	// And it does not reach across the city. A man at the docks was not in the
	// room being turned out.
	if away.Weapon == 0 {
		t.Fatalf("%s was at the docks and the search took their gun anyway", away.Name)
	}
	said := w.History[len(w.History)-1].Text
	if !strings.Contains(said, here.Name) {
		t.Fatalf("%s lost their gun and nobody said so: %q", here.Name, said)
	}
}

// And a search of somebody carrying nothing, with nobody armed beside them,
// finds nothing — which is what stops the line being printed every raid.
func TestASearchOfNobodyArmedFindsNothing(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Weapon, w.Player.Armour = 0, 0
	if w.SeizeArms() {
		t.Fatal("a search took hardware off somebody carrying none")
	}
}
