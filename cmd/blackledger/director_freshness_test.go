package main

import (
	"blackledger/core"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirectorDiscardsOwnershipChangeDuringModelRequest(t *testing.T) {
	a := testApp(t)
	snapshot, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	model := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		// Simulate the player buying property while the local model is working.
		if err := a.s.Change(func(current *core.World) error {
			property := current.Properties["laundry"]
			property.Owner = "player:1"
			current.Properties["laundry"] = property
			return nil
		}); err != nil {
			t.Error(err)
		}
		p := core.Proposal{Location: "bar", Title: "A discreet delivery", Body: "Deliver a sealed parcel to Saint Agnes.", Speaker: "mara", Operation: snapshot.NextDirectorOperation(), Outcome: "Delivered."}
		b, _ := json.Marshal(p)
		json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": string(b)}})
	}))
	defer model.Close()
	t.Setenv("BLACK_LEDGER_OLLAMA", model.URL)
	if err := a.generate(snapshot); !errors.Is(err, errDirectorContextChanged) {
		t.Fatalf("expected stale context, got %v", err)
	}
	saved, err := a.s.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(saved.Offers) != 0 || saved.Properties["laundry"].Owner != "player:1" || calls != 1 {
		t.Fatal("stale offer queued, player change lost, or obsolete context retried")
	}
}

func TestDirectorFreshnessAllowsProgressButRequiresAvailableNewContact(t *testing.T) {
	t.Parallel()
	before := core.New(27)
	before.Player.Crew = []core.Crew{{ID: "leo", Loyalty: 65}}
	proposal := core.Proposal{Speaker: "leo"}
	now := before.Clone()
	now.Minute += 90
	now.Player.Cash += 75
	p := now.Properties["laundry"]
	p.Condition -= 15
	now.Properties["laundry"] = p
	if err := validateDirectorFreshness(before, now, proposal, nil); err != nil {
		t.Fatal(err)
	}
	now.Player.Crew[0].Loyalty = 29
	if err := validateDirectorFreshness(before, now, proposal, nil); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("disloyal new contact still allowed")
	}
	now = before.Clone()
	now.Player.Crew = nil
	if err := validateDirectorFreshness(before, now, proposal, nil); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("departed new contact still allowed")
	}
}

func TestDirectorFreshnessPreservesCanonicalFollowupAfterStandingChange(t *testing.T) {
	t.Parallel()
	before := core.New(27)
	before.Factions[1].Goodwill = 9
	before.Arrangements = []core.ArrangementMemory{{ID: "finished", Life: before.Life, Speaker: "elena", Beneficiary: "russo", Status: "completed", Result: "Payment collected."}}
	connection := &before.Arrangements[0]
	proposal := core.Proposal{Speaker: "elena", Beneficiary: "russo"}
	now := before.Clone()
	now.Factions[1].Goodwill = -10
	if err := validateDirectorFreshness(before, now, proposal, connection); err != nil {
		t.Fatal(err)
	}
	if err := validateDirectorFreshness(before, now, proposal, nil); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("new offer ignores lost access")
	}
	now.Arrangements[0].Status = "interrupted"
	if err := validateDirectorFreshness(before, now, proposal, connection); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("invalid completion reused")
	}
	now = before.Clone()
	now.Factions[1].Leader = "Another boss"
	if err := validateDirectorFreshness(before, now, proposal, connection); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("changed leadership reused")
	}
	now = before.Clone()
	now.NPC("elena").Name = "Another person"
	if err := validateDirectorFreshness(before, now, proposal, connection); !errors.Is(err, errDirectorContextChanged) {
		t.Fatal("changed identity reused")
	}
}
