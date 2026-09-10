package sim

import (
	"blackledger/core"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func sampleProposal() core.Proposal {
	return core.Proposal{Title: "An agreed payment", Body: "Collect the agreed payment for Russo Outfit.", Speaker: "elena", Operation: "collection", Outcome: "Payment collected.", Beneficiary: "russo"}
}
func TestCorpusValidatesAndFingerprintsExactInput(t *testing.T) {
	t.Parallel()
	b, _ := json.Marshal([]core.Proposal{sampleProposal()})
	p, hash, err := ReadCorpus(strings.NewReader(string(b)))
	if err != nil {
		t.Fatal(err)
	}
	_, again, _ := ReadCorpus(strings.NewReader(string(b)))
	if len(p) != 1 || len(hash) != 64 || hash != again {
		t.Fatal("corpus identity lost")
	}
	invalid := []string{"[]", "null", string(b) + " {}", `[{"title":"Bad","script":"run code"}]`, strings.Repeat(" ", 2<<20) + string(b)}
	bad := sampleProposal()
	bad.Beneficiary = "not-a-family"
	data, _ := json.Marshal([]core.Proposal{bad})
	invalid = append(invalid, string(data))
	for _, s := range invalid {
		if _, _, err := ReadCorpus(strings.NewReader(s)); err == nil {
			t.Fatal("invalid corpus accepted")
		}
	}
}
func TestReplayConsumesOnceAndIsReproducible(t *testing.T) {
	t.Parallel()
	corpus := []core.Proposal{sampleProposal()}
	a := RunRecorded(27, "investor", "replay", 100, false, corpus)
	b := RunRecorded(27, "investor", "replay", 100, false, corpus)
	if a.Error != "" || a.ReplayQueued != 1 || a.Events["proposal"] == 0 || !reflect.DeepEqual(a, b) {
		t.Fatalf("invalid finite replay: %+v", a)
	}
	if corpus[0].Beneficiary != "russo" {
		t.Fatal("replay mutated input")
	}
}
func TestReplayDoesNotSkipInvalidProposal(t *testing.T) {
	t.Parallel()
	p := sampleProposal()
	p.Operation = "invented-operation"
	r := RunRecorded(27, "investor", "replay", 100, false, []core.Proposal{p})
	if r.Error == "" || r.ReplayQueued != 0 {
		t.Fatal("invalid replay silently accepted")
	}
}
