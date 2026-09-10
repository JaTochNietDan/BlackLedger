package core

import "testing"

func house(t *testing.T) (*World, *Faction) {
	t.Helper()
	w := New(43)
	w.MigrateLivingWorld()
	f := &w.Factions[0]
	// A failing organization, which is when people start looking upward.
	f.Peak, f.Power = 100, 40
	return w, f
}

func TestALoyalLieutenantNeverMovesOnTheirOwnLeader(t *testing.T) {
	t.Parallel()
	w, f := house(t)
	members := w.Members(f.ID)
	if len(members) < 2 {
		t.Fatal("the organization had one person in it")
	}
	leader, challenger := members[0], members[1]
	challenger.Name, challenger.Ambition, challenger.Rank = nameWith(t, "loyal"), 100, RankLieutenant
	for _, other := range members[2:] {
		other.Rank = RankSoldier
	}
	if w.coupNerve(challenger, leader) != 0 {
		t.Fatalf("a loyal person had %f nerve for it", w.coupNerve(challenger, leader))
	}
	for i := 0; i < 300; i++ {
		w.InternalMove(f)
	}
	if leader.Dead {
		t.Fatal("a loyal lieutenant took the top by force")
	}
}

func TestAGrievanceAgainstYourOwnLeaderIsWorthMoreThanAmbition(t *testing.T) {
	t.Parallel()
	w, f := house(t)
	members := w.Members(f.ID)
	leader, challenger := members[0], members[1]
	challenger.Name, challenger.Ambition, challenger.Rank = nameWith(t, "careful"), 40, RankLieutenant
	plain := w.coupNerve(challenger, leader)
	w.Resent(challenger.ID, leader.ID, 60, "being passed over")
	if owed := w.coupNerve(challenger, leader); owed <= plain {
		t.Fatalf("a grievance was worth nothing: %f against %f", owed, plain)
	}
}

func TestNobodyMovesOnALeaderWhoIsWinningWithoutAReason(t *testing.T) {
	t.Parallel()
	w, f := house(t)
	f.Power, f.Peak = 100, 100 // winning
	members := w.Members(f.ID)
	leader, challenger := members[0], members[1]
	challenger.Name, challenger.Ambition, challenger.Rank = nameWith(t, "hot"), 100, RankLieutenant
	for i := 0; i < 300; i++ {
		w.InternalMove(f)
	}
	if leader.Dead || challenger.Dead {
		t.Fatal("somebody moved on a leader who was winning, with nothing owed")
	}
	// Give them a specific reason that has been growing, and it is different.
	w.Resent(challenger.ID, leader.ID, GrudgeCap, "an old debt")
	moved := false
	for i := 0; i < 300 && !moved; i++ {
		moved = w.InternalMove(f)
	}
	if !moved {
		t.Fatal("a heavy grievance against a winning leader was never acted on")
	}
}

func TestTheWeekAfterIsWhereTheDramaIs(t *testing.T) {
	t.Parallel()
	const runs = 400
	purged, walkedOut, resented := 0, 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, f := house(t)
		w.WorldRNG = seed * 2654435761
		members := w.Members(f.ID)
		winner, loser := members[0], members[1]
		before := 0
		for i := range w.NPCs {
			if w.NPCs[i].Dead {
				before++
			}
		}
		organizations := len(w.Factions)
		w.CoupAftermath(f, winner, loser)
		after := 0
		for i := range w.NPCs {
			if w.NPCs[i].Dead {
				after++
			}
		}
		if after > before {
			purged++
		}
		if len(w.Factions) > organizations {
			walkedOut++
		}
		for _, g := range w.Grudges {
			if g.Against == winner.ID {
				resented++
				break
			}
		}
	}
	if purged == 0 || purged == runs {
		t.Fatalf("somebody was dealt with in %d of %d aftermaths", purged, runs)
	}
	if resented == 0 {
		t.Fatal("nobody left in the organization had a view on the man who did it")
	}
	t.Logf("across %d aftermaths: somebody was dealt with %d times, somebody walked out and started their own %d times, and the man who did it was resented by somebody left behind %d times",
		runs, purged, walkedOut, resented)
}

func TestACoupLeavesTheOrganizationWorseOff(t *testing.T) {
	t.Parallel()
	const runs = 300
	weakened, upheaval := 0, 0
	for seed := uint32(1); seed <= runs; seed++ {
		w, f := house(t)
		w.WorldRNG = seed * 2654435761
		members := w.Members(f.ID)
		if len(members) < 2 {
			continue
		}
		members[1].Name, members[1].Ambition, members[1].Rank = nameWith(t, "hot"), 100, RankLieutenant
		power := f.Power
		moved := false
		for i := 0; i < 100 && !moved; i++ {
			moved = w.InternalMove(f)
		}
		if !moved {
			continue
		}
		if f.Power < power {
			weakened++
		}
		if w.hasNewsKind("politics") || w.hasRecord("A move against "+members[0].Name) || w.hasRecord("A move that failed") {
			upheaval++
		}
	}
	if weakened == 0 {
		t.Fatal("a coup never cost the organization anything")
	}
	t.Logf("%d of %d organizations that ate one of their own came out weaker, with %d of them making the record", weakened, runs, upheaval)
}
