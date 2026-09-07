package main

import (
	"blackledger/core"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestOrdinaryLeaderJobsServeTheirOwnFamily(t *testing.T) {
	w := core.New(27)
	for _, tc := range []struct {
		speaker, beneficiary string
		valid                bool
	}{{"elena", "russo", true}, {"elena", "bellandi", false}, {"elena", "", false}, {"vittorio", "bellandi", true}, {"vittorio", "russo", false}, {"mara", "russo", true}, {"leo", "bellandi", true}, {"mara", "", true}} {
		if (validateSpeakerBeneficiary(w, tc.speaker, tc.beneficiary) == nil) != tc.valid {
			t.Fatal(tc)
		}
	}
	w.Factions[1].Goodwill = 9
	w.Arrangements = []core.ArrangementMemory{{ID: "old", Life: 1, Speaker: "mara", Status: "declined"}}
	p := proposalSchema(w, "mediation", nil)["properties"].(map[string]any)
	if !reflect.DeepEqual(p["speaker"].(map[string]any)["enum"], []string{"elena"}) || !reflect.DeepEqual(p["beneficiary"].(map[string]any)["enum"], []string{"russo"}) {
		t.Fatal("schema allows contradictory single-speaker pairing")
	}
	w.Arrangements = []core.ArrangementMemory{{ID: "legacy", Life: 1, Speaker: "elena", Beneficiary: "bellandi", Status: "completed"}}
	if directorConnection(w) != nil || w.Arrangements[0].Beneficiary != "bellandi" {
		t.Fatal("invalid legacy connection propagated or rewritten")
	}
	w.Arrangements[0].Beneficiary = "russo"
	if directorConnection(w) == nil {
		t.Fatal("valid follow-up lost")
	}
}

func TestWrongLeaderBeneficiaryIsCorrectedBeforeQueue(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Factions[1].Goodwill = 9
		w.Arrangements = []core.ArrangementMemory{{Life: 1, Speaker: "mara", Status: "declined"}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	calls := 0
	model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		beneficiary := "bellandi"
		if calls > 1 {
			beneficiary = "russo"
		}
		b, _ := json.Marshal(core.Proposal{Location: "bar", Title: "A scheduling dispute", Body: "Mediate shared access at Saint Agnes for " + beneficiary + ".", Speaker: "elena", Beneficiary: beneficiary, Operation: "mediation", Outcome: "Settled."})
		json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
	}))
	defer model.Close()
	t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
	snapshot, _ := a.s.Read()
	if err := a.generate(snapshot); err != nil {
		t.Fatal(err)
	}
	saved, _ := a.s.Read()
	if calls != 2 || len(saved.Offers) != 1 || saved.Offers[0].Event.Beneficiary != "russo" {
		t.Fatal("contradictory job queued")
	}
}

func TestSpokenOfferNamesTheFactionReceivingCredit(t *testing.T) {
	w := core.New(27)
	for _, tc := range []struct {
		body, beneficiary string
		valid             bool
	}{
		{"A vendor refuses his contract with Bellandi. Settle the dispute.", "russo", false},
		{"Settle shared access for the Russo Outfit.", "russo", true},
		{"This is Russo's business.", "russo", true},
		{"This helps Russoness.", "russo", false},
		{"Bellandi disagrees, but Russo needs shared access.", "russo", true},
		{"Settle a disagreement between two staff.", "", true},
	} {
		if (validateBeneficiaryMention(w, core.Proposal{Body: tc.body, Beneficiary: tc.beneficiary}) == nil) != tc.valid {
			t.Fatal(tc)
		}
	}
}

func TestRivalOnlyDialogueIsCorrectedBeforeQueue(t *testing.T) {
	a := testApp(t)
	if err := a.s.Change(func(w *core.World) error {
		w.Factions[1].Goodwill = 9
		w.Arrangements = []core.ArrangementMemory{{Life: w.Life, Speaker: "mara", Operation: "courier", Status: "declined"}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	calls := 0
	model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		body := "At Mercer Exchange, settle the vendor's contract with Bellandi."
		if calls > 1 {
			body = "At Mercer Exchange, settle loading access for Russo suppliers."
		}
		b, _ := json.Marshal(core.Proposal{Location: "market", Title: "Shared access", Body: body, Speaker: "elena", Beneficiary: "russo", Operation: "mediation", Outcome: "Settled."})
		json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
	}))
	defer model.Close()
	t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
	snapshot, _ := a.s.Read()
	if err := a.generate(snapshot); err != nil {
		t.Fatal(err)
	}
	saved, _ := a.s.Read()
	if calls != 2 || len(saved.Offers) != 1 || saved.Offers[0].Event.Body != "At Mercer Exchange, settle loading access for Russo suppliers." {
		t.Fatal("rival-only dialogue entered the queue")
	}
}
