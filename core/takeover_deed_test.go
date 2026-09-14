package core

import "testing"

func TestAmbitiousTakeoverRecordsOnePersonalDeed(t *testing.T) {
	w := New(404)
	for _, p := range w.Properties {
		p.Owner = "bellandi"
	}
	w.Properties["laundry"].Owner = "independent"
	first, second := w.NPC("leo"), w.NPC("mara")
	for _, n := range []*NPC{first, second} {
		n.Faction = ""
		n.Skill = 80
		n.Ambition = 80
		n.Held = 0
	}
	first.Heading = "club"
	first.Arrives = w.Minute + 60
	if !w.claimPremises(first) {
		t.Fatal("could not claim unheld business")
	}
	if w.Properties["laundry"].Owner != first.ID || first.Heading != "" || first.Location != "laundry" {
		t.Fatal("takeover did not establish its real owner and location")
	}
	if w.claimPremises(second) {
		t.Fatal("second person claimed the same business")
	}
	if w.claimPremises(first) {
		t.Fatal("already established proprietor claimed again")
	}
	w.Kill(first.ID, "test succession")
	if w.Properties["laundry"].Owner != "independent" {
		t.Fatal("personal deed did not return to market")
	}
	if !w.claimPremises(second) {
		t.Fatal("vacant business never became available again")
	}
}
