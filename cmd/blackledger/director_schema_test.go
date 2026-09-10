package main

import (
	"blackledger/core"
	"testing"
)

func TestProposalSchemaRestrictsFollowUpIdentity(t *testing.T) {
	t.Parallel()
	w := core.New(27)
	c := &core.ArrangementMemory{Speaker: "mara", Beneficiary: ""}
	s := proposalSchema(w, "collection", c)
	p := s["properties"].(map[string]any)
	for key, want := range map[string]string{"speaker": "mara", "beneficiary": "", "operation": "collection"} {
		values := p[key].(map[string]any)["enum"].([]string)
		if len(values) != 1 || values[0] != want {
			t.Fatalf("%s constraint wrong", key)
		}
	}
	if s["additionalProperties"] != false {
		t.Fatal("arbitrary model fields allowed")
	}
	fresh := proposalSchema(w, "mediation", nil)["properties"].(map[string]any)
	if values := fresh["speaker"].(map[string]any)["enum"].([]string); len(values) != 1 || values[0] != "mara" {
		t.Fatal("newcomer was offered unestablished contacts")
	}
}
