package main

import (
	"blackledger/core"
	"testing"
)

func attributionWorld(barOwner string) *core.World {
	w := &core.World{
		Factions: []core.Faction{
			{ID: "bellandi", Name: "Bellandi Family", Leader: "Vittorio Bellandi"},
			{ID: "russo", Name: "Russo Outfit", Leader: "Elena Russo"},
		},
		Properties: map[string]*core.Property{
			"bar":  {Owner: barOwner, Condition: 100},
			"club": {Owner: "bellandi", Condition: 100},
		},
	}
	w.Player.Location = "bar"
	w.Player.Name = "Alex Varga"
	return w
}

func TestNeutralWorkDoesNotHandItsPeopleToAFamily(t *testing.T) {
	// Observed live on qwen3.5:35b-a3b at Saint Agnes, which is independently
	// owned: a neutral mediation described "two Bellandi staff", implying family
	// standing that the neutral reward never moves.
	w := attributionWorld("independent")
	for _, text := range []string{
		"Two Bellandi staff are at odds over access to a storeroom at Saint Agnes.",
		"Two of the Bellandi's men are arguing over the storeroom.",
		"Bellandi soldiers are blocking the door.",
		"Vittorio's people want this settled.",
		"The Russo crew cannot agree on the loading bay.",
	} {
		t.Run(text, func(t *testing.T) {
			if validateFactionAttribution(w, core.Proposal{Location: "bar", Body: text}) == nil {
				t.Fatal("neutral job attributed to a family was accepted")
			}
		})
	}
}

func TestFamilyAttributionIsAllowedWhenEarnedOrOwned(t *testing.T) {
	w := attributionWorld("independent")
	// The family receiving the credit may own the people in its own job.
	if err := validateFactionAttribution(w, core.Proposal{Location: "bar", Beneficiary: "bellandi", Body: "Two Bellandi staff are at odds over the storeroom."}); err != nil {
		t.Fatal("beneficiary's own people rejected:", err)
	}
	// A family that owns the premises staffs them, whoever the work serves.
	if err := validateFactionAttribution(w, core.Proposal{Location: "club", Body: "Two Bellandi staff are at odds over the storeroom."}); err != nil {
		t.Fatal("owner's own people rejected:", err)
	}
	// A property that changed hands is judged on the saved owner, not the name.
	owned := attributionWorld("bellandi")
	if err := validateFactionAttribution(owned, core.Proposal{Location: "bar", Body: "Two Bellandi staff are at odds over the storeroom."}); err != nil {
		t.Fatal("current owner's people rejected:", err)
	}
}

func TestOrdinaryMentionsOfAFamilyRemainLegal(t *testing.T) {
	w := attributionWorld("independent")
	for _, text := range []string{
		"Keep the Bellandi Family out of this.",
		"The Russo Outfit will hear about it either way.",
		"Two of the floor staff are deadlocking over the loading bay.",
		"Vittorio Bellandi drinks here, so keep it quiet.",
		"Do this well and the Russo Outfit may notice you.",
	} {
		t.Run(text, func(t *testing.T) {
			if err := validateFactionAttribution(w, core.Proposal{Location: "bar", Body: text}); err != nil {
				t.Fatal("ordinary family mention rejected:", err)
			}
		})
	}
}

func TestTheLiveStoreroomOfferIsRejected(t *testing.T) {
	// Generated live on qwen3.5:35b-a3b during a fresh playtest at 20d1c76:
	// mechanically neutral mediation, independently owned Saint Agnes.
	w := attributionWorld("independent")
	body := "You handled the envelope quietly, Alex. Now, two Bellandi staff are at odds over access to a storeroom at Saint Agnes. They won’t listen to me, but they’ll hear you. Talk to each one separately, find a compromise, and make sure the deal sticks. No need for a show, but don’t let it drag on."
	p := core.Proposal{Location: "bar", Title: "Storeroom Dispute", Body: body, Speaker: "mara", Operation: "mediation"}
	if err := validateFactionAttribution(w, p); err == nil {
		t.Fatal("the observed neutral-but-Bellandi offer was accepted")
	} else {
		t.Log("rejected:", err)
	}
	// Everything else about that offer was acceptable and must stay so.
	if err := validateInWorldVoice(p); err != nil {
		t.Fatal("in-world voice rejected a good offer:", err)
	}
	if err := validateSceneTitle(p); err != nil {
		t.Fatal("scene title rejected a good offer:", err)
	}
	if err := validateSpokenTerms(p); err != nil {
		t.Fatal("spoken terms rejected a good offer:", err)
	}
	if err := validatePlayerRole(w, p); err != nil {
		t.Fatal("player role rejected a good offer:", err)
	}
}
