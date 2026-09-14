package core

import "testing"

func TestNamedHandKeepsItsActorAfterRosterReorders(t *testing.T) {
	w := withCrew(t, 80)
	w.Player.Crew = append(w.Player.Crew, Crew{ID: "mara", Name: "Mara", Loyalty: 55})
	n := w.NPC("mara")
	if n == nil {
		t.Fatal("missing actor")
	}
	n.Weapon = 3
	hand, ok := w.NamedHands(n.ID)
	if !ok {
		t.Fatal("missing hand")
	}
	edge := w.HandEdge(hand)
	w.Player.Crew[0], w.Player.Crew[1] = w.Player.Crew[1], w.Player.Crew[0]
	if w.HandEdge(hand) != edge || w.strikeAttacker(hand).ID != n.ID || w.strikeAttacker(hand).Weapon != 3 {
		t.Fatal("roster order changed actor or equipment")
	}
	w.HandHurt(hand, 20, "test job")
	if w.Player.Crew[0].Loyalty != 30 || w.Player.Crew[1].Loyalty != 80 || w.Player.Health != 100 {
		t.Fatal("injury charged the wrong person")
	}
}

func TestNamedHandCanBeASignedFamilyMember(t *testing.T) {
	w := withCrew(t, 80)
	n := w.NPC("mara")
	n.Faction = w.PlayerOrganizationID()
	n.Trust = 70
	n.Heading = ""
	n.Arrives = 0
	n.Held = 0
	hand, ok := w.NamedHands(n.ID)
	if !ok || w.HandReadiness(hand) != "" {
		t.Fatal("signed member unavailable", w.HandReadiness(hand))
	}
	w.HandHurt(hand, 20, "test job")
	if n.Trust != 45 || w.Player.Crew[0].Loyalty != 80 {
		t.Fatal("wrong person's loyalty charged")
	}
	n.Held = w.Minute + 60
	if w.HandReadiness(hand) == "" {
		t.Fatal("sent a held member")
	}
	n.Held = 0
	n.Faction = ""
	if w.HandReadiness(hand) == "" {
		t.Fatal("sent somebody after they left")
	}
}

func TestNamedAssassinationLossDoesNotClearOtherAssociates(t *testing.T) {
	deaths, captures := 0, 0
	for seed := uint32(1); seed <= 100; seed++ {
		w := withCrew(t, 80)
		w.RNG = scatter(seed)
		w.Player.Crew = append(w.Player.Crew, Crew{ID: "mara", Name: "Mara", Loyalty: 80})
		hand, _ := w.NamedHands("mara")
		victim := w.NPC("leo")
		w.itWentWrong(victim, hand, "test premises", nil)
		found := false
		for _, c := range w.Player.Crew {
			if c.ID == "leo" {
				found = true
			}
		}
		if !found || w.Player.Health != 100 {
			t.Fatal("another associate or player paid for failure")
		}
		actor := w.NPC("mara")
		if actor.Dead {
			deaths++
		}
		if w.Inside(actor) {
			captures++
		}
	}
	if deaths == 0 || captures == 0 {
		t.Fatal("did not exercise both death and capture", deaths, captures)
	}
}
