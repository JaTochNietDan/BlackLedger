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
