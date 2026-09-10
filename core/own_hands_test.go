package core

import (
	"strings"
	"testing"
)

// Three pieces of work come in pairs: go yourself, or send somebody. The pair
// is the whole point — the crew version trades odds and standing for less
// attention on you — so the two descriptions have to say which is which.
//
// Reading the garage, the docks and the club side by side, sabotage did not.
// Its label read "Move against Bellandi Family yourself" and its own
// description opened "Send your crew against The Monarch", which is the other
// button.
//
// What settles it is the failure text, not the label. Sabotage requires a crew
// on both halves and reads their loyalty on both, so the player never goes
// alone; the difference is whether they go at all. This half says "You and Leo
// left without reaching anything", the other says the crew "went in without
// you". So the property is narrow and exact: the half the player attends must
// not describe itself as dispatching somebody. Saying the crew comes along is
// true, and stays.

func TestWorkYouDoYourselfIsNotDescribedAsSendingSomebody(t *testing.T) {
	t.Parallel()
	w := New(3)
	w.District = 2
	w.Player.Cash = 2000
	w.Player.Respect = 60
	// Somebody to send, so both halves of every pair are offered.
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 65}}
	checked := 0
	for i := range Locations {
		id := Locations[i].ID
		w.Player.Location = id
		byID := map[string]Action{}
		for _, a := range w.Actions(id) {
			byID[a.ID] = a
		}
		for kind, a := range byID {
			if strings.HasSuffix(kind, ":crew") {
				continue
			}
			if _, paired := byID[kind+":crew"]; !paired {
				continue
			}
			checked++
			for _, sending := range []string{"Send ", "send "} {
				if strings.Contains(a.Detail, sending) {
					t.Errorf("%s at %s is done by the player, and its description says %q:\n  label:  %s\n  detail: %s",
						kind, id, sending, a.Label, a.Detail)
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("no paired work was offered anywhere, so this proved nothing")
	}
	t.Logf("checked %d self/crew pairs", checked)
}

// The pair also has to be recognisable as a pair. Robbery and mugging both name
// the same subject on each half — the till, or the person being robbed — so the
// player reads them as one choice about who carries it out. Sabotage named the
// family on one half and the premises on the other, which reads as two
// different acts. The family is still in the description, where robbery and
// mugging keep the same kind of detail.
func TestBothHalvesOfAPairNameTheSameTarget(t *testing.T) {
	t.Parallel()
	w := New(3)
	w.District = 2
	w.Player.Cash = 2000
	w.Player.Crew = []Crew{{"leo", "Leo Carver", 65}}
	checked := 0
	for i := range Locations {
		id, name := Locations[i].ID, Locations[i].Name
		w.Player.Location = id
		byID := map[string]Action{}
		for _, a := range w.Actions(id) {
			byID[a.ID] = a
		}
		crew, ok := byID["sabotage:crew"]
		if !ok {
			continue
		}
		own := byID["sabotage"]
		checked++
		if !strings.Contains(crew.Label, name) {
			t.Fatalf("the crew half stopped naming the premises: %q", crew.Label)
		}
		if !strings.Contains(own.Label, name) {
			t.Errorf("at %s the two halves name different things:\n  own:  %s\n  crew: %s", id, own.Label, crew.Label)
		}
	}
	if checked == 0 {
		t.Fatal("sabotage was never offered, so this proved nothing")
	}
}

// One room, one crew, one name. Reading the market on a driven save, three
// buttons said "Send Bela Havel for the till", "Send Bela Havel after Elena
// Russo", "Send Bela Havel against Mercer Exchange" — and beside them, "Send
// Leo on collections", "Pay Leo a bonus", and a refusal reading "Leo refuses
// assignments below 30 loyalty". Bela Havel was the crew. Leo Carver was a name
// written into three strings, and hiring whoever actually drives is what made
// it visible.
func TestOneCrewIsCalledByOneName(t *testing.T) {
	t.Parallel()
	w := New(3)
	w.District = 2
	w.Player.Cash = 2000
	w.Player.Crew = []Crew{{"bela", "Bela Havel", 20}}
	w.NPCs = append(w.NPCs, NPC{ID: "bela", Name: "Bela Havel", Location: "market"})
	w.Player.Location = "market"
	found := 0
	for _, a := range w.Actions("market") {
		if a.ID != "delegate" && a.ID != "crew_bonus" {
			continue
		}
		found++
		for _, text := range []string{a.Label, a.Reason} {
			if text == "" {
				continue
			}
			if strings.Contains(text, "Leo") {
				t.Errorf("%s names somebody who is not in the crew: %q", a.ID, text)
			}
		}
		if !strings.Contains(a.Label, "Bela") {
			t.Errorf("%s does not name the person in the crew: %q", a.ID, a.Label)
		}
	}
	if found != 2 {
		t.Fatalf("expected delegate and crew_bonus, found %d", found)
	}
}
