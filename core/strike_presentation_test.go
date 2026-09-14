package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSuccessfulStrikeRecordsActualAttackerAndWeapon(t *testing.T) {
	for _, crew := range []bool{false, true} {
		for tier := 0; tier <= 3; tier++ {
			w, mark := striker(t)
			id := mark.ID
			w.Player.Weapon = tier
			hand := w.OwnHands()
			actorID := "player"
			if crew {
				actorID = "weapon-courier"
				w.NPCs = append(w.NPCs, NPC{ID: actorID, Name: "Test Courier", Weapon: tier})
				w.Player.Crew = []Crew{{actorID, "Test Courier", 90}}
				hand = Hand{Crew: true, Name: "Test Courier"}
				w.Player.Weapon = 3 - tier // Must use the courier's weapon, not the player's.
			}
			mark = w.NPC(id)
			w.finishThem(mark, hand, "Saint Agnes", 0, nil)
			var recorded *VisualCue
			for i := range w.VisualCues {
				if w.VisualCues[i].Attacker != nil {
					recorded = &w.VisualCues[i]
				}
			}
			if recorded == nil || recorded.Attacker.ID != actorID || recorded.Attacker.Weapon != tier {
				t.Fatalf("crew=%v tier=%d cue=%+v", crew, tier, recorded)
			}
			want := "gunfight"
			if tier == 0 {
				want = "attack"
			}
			if recorded.Kind != want {
				t.Fatalf("tier=%d kind=%s", tier, recorded.Kind)
			}
			cause := ""
			for _, cue := range w.VisualCues {
				if cue.Kind == "killing" {
					cause = cue.Caption
					break
				}
			}
			if tier > 0 && !strings.Contains(cause, []string{"", "revolver", "pump shotgun", "Thompson"}[tier]) {
				t.Fatalf("weapon and cause disagree: %s", cause)
			}
			if tier == 0 && strings.Contains(cause, "shot") {
				t.Fatalf("unarmed strike described gunfire: %s", cause)
			}
			w.Player.Weapon = 0
			if crew {
				w.NPC(actorID).Weapon = 0
			}
			data, err := json.Marshal(recorded)
			if err != nil {
				t.Fatal(err)
			}
			var restored VisualCue
			if err = json.Unmarshal(data, &restored); err != nil {
				t.Fatal(err)
			}
			if restored.Attacker == nil || restored.Attacker.Weapon != tier {
				t.Fatal("cue equipment changed after the event")
			}
		}
	}
}

func TestStrikeScenarioUsesPreDeathContextAndOneWorldDraw(t *testing.T) {
	for _, tc := range []struct {
		name              string
		tier, armed, sore int
		travelling        bool
		want              string
	}{
		{"unaware", 1, 0, 0, false, "back-of-head"},
		{"armed target", 1, 1, 0, false, "close-shot"},
		{"personal grievance", 1, 0, 60, false, "close-shot"},
		{"moving target", 1, 0, 0, true, "close-shot"},
		{"shotgun", 2, 0, 0, false, "close-shot"},
		{"automatic", 3, 0, 0, false, "burst"},
		{"unarmed", 0, 0, 0, false, "close-quarters"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, n := striker(t)
			w.WorldRNG = 1
			before := w.RNG
			n.Weapon, n.Sore = tc.armed, tc.sore
			n.Heading = ""
			n.Sets = 0
			if tc.travelling {
				n.Heading = "club"
				n.Arrives = w.Minute + 30
			}
			cue, cause := w.strikePresentation(n, tc.tier)
			if cue.Variant != tc.want || cue.Victim.ID != n.ID || cue.Victim.Name != n.Name {
				t.Fatalf("unexpected scenario: %+v", cue)
			}
			if w.RNG != before || w.WorldRNG != 1015568748 {
				t.Fatal("scenario altered combat RNG or used extra world draws")
			}
			if strings.Contains(cause, "once in the back of the head") != (tc.want == "back-of-head") {
				t.Fatalf("contradictory death description: %s", cause)
			}
		})
	}
}

func TestStrikeScenarioVariesAndSurvivesSavedResult(t *testing.T) {
	seen := map[string]bool{}
	for _, seed := range []uint32{1, 1500} {
		w, n := striker(t)
		w.WorldRNG = seed
		w.Player.Weapon = 1
		n.Weapon, n.Sore = 0, 0
		n.Heading = ""
		w.finishThem(n, w.OwnHands(), "Saint Agnes", 0, nil)
		count := 0
		for _, cue := range w.VisualCues {
			if cue.Kind != "killing" && cue.Kind != "gunfight" {
				continue
			}
			if cue.Strike == nil || cue.Strike.Victim.ID != n.ID {
				t.Fatalf("missing linked victim: %+v", cue)
			}
			seen[cue.Strike.Variant] = true
			count++
		}
		if count != 2 {
			t.Fatalf("expected paired attack/death cues, got %d", count)
		}
		w.LastResult = &Result{Cues: w.VisualCues}
		data, err := json.Marshal(w)
		if err != nil {
			t.Fatal(err)
		}
		var restored World
		if err = json.Unmarshal(data, &restored); err != nil {
			t.Fatal(err)
		}
		n.Sore = 100
		n.Weapon = 3
		w.Player.Weapon = 0
		for i, cue := range restored.LastResult.Cues {
			original := w.LastResult.Cues[i]
			if original.Strike != nil && (cue.Strike == nil || *cue.Strike != *original.Strike) {
				t.Fatal("saved scenario changed")
			}
		}
	}
	if !seen["back-of-head"] || !seen["close-shot"] {
		t.Fatalf("eligible contexts lack variation: %v", seen)
	}
}
