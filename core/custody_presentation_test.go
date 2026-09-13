package core

import "testing"

func TestArrestCueIdentifiesPrisonerSeparatelyFromDetective(t *testing.T) {
	w := proprietor(t)
	w.Confine(2, "an armed robbery")
	cue := w.VisualCues[len(w.VisualCues)-1]
	if cue.Kind != "arrest" || cue.Detainee == nil || cue.Detainee.ID != "player" || cue.Detainee.Name != w.Player.Name {
		t.Fatalf("wrong prisoner: %+v", cue)
	}
	if cue.Target != "precinct" || !w.Held() {
		t.Fatal("presentation lost committed custody")
	}
}
