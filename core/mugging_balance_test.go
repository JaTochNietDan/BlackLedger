package core

import "testing"

// Who you pick is the whole decision: a soldier carries little and is easy, a
// boss carries a great deal and is standing next to his organization.

func TestWhoYouPickDecidesWhatItIsWorthAndWhatItCosts(t *testing.T) {
	const runs = 400
	type result struct{ took, money, hurt int }
	measure := func(rank int) result {
		r := result{}
		for seed := uint32(1); seed <= runs; seed++ {
			w := New(101)
			w.MigrateLivingWorld()
			w.Player.Location, w.Player.Health, w.Player.Cash = "club", 100, 0
			w.Player.Contacts, w.Player.Respect = 3, 45
			// Everybody out of the room but one man of the wanted standing.
			var mark *NPC
			for i := range w.NPCs {
				n := &w.NPCs[i]
				if n.Location == "club" && mark == nil && !n.Dead {
					mark = n
					continue
				}
				n.Location = "estate"
			}
			if mark == nil {
				t.Fatal("nobody was in the room")
			}
			mark.Rank, mark.Skill = rank, 40+rank/3
			// People carry their own money, settled when they arrived in the
			// city. This test promotes a man after the fact, so his pocket has
			// to be settled again onto the standing he now has — otherwise
			// both ranks are measured carrying whatever the same man happened
			// to have, and the comparison says nothing.
			mark.Purse = w.StandingPurse(mark)
			// Drive the roll from the seed rather than the campaign, so the two
			// ranks face the same run of luck.
			// The first draw after seeding an LCG is not uniform across a run
			// of strided seeds, so the stream is walked forward a little
			// before the attempt. Both ranks walk it the same distance.
			w.RNG = seed*2654435761 + uint32(rank)
			for i := uint32(0); i < seed%11; i++ {
				w.Random()
			}
			if err := w.Mug("club", w.OwnHands()); err != nil {
				t.Fatal(err)
			}
			if w.Player.Cash > 0 {
				r.took++
				r.money += w.Player.Cash
			}
			if w.Player.Health < 100 {
				r.hurt++
			}
		}
		return r
	}

	soldier, boss := measure(RankSoldier), measure(RankLeader)
	if soldier.took <= boss.took {
		t.Fatalf("a soldier was taken %d times in %d and a boss %d, so standing does not protect anybody", soldier.took, runs, boss.took)
	}
	if boss.took > 0 && soldier.took > 0 {
		perSoldier, perBoss := soldier.money/soldier.took, boss.money/boss.took
		if perBoss <= perSoldier {
			t.Fatalf("a boss carried $%d and a soldier $%d", perBoss, perSoldier)
		}
		t.Logf("%d attempts each: a soldier gives way %d times for $%d a time; a boss %d times for $%d a time. Hurt %d times going after the soldier and %d after the boss",
			runs, soldier.took, perSoldier, boss.took, perBoss, soldier.hurt, boss.hurt)
	}
}
