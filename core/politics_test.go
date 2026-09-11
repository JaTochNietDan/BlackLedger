package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func pressureWorld() *World { return pressureWorldFor(27) }

func pressureWorldFor(seed uint32) *World {
	w := New(seed)
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Cash = 200
	w.NextPressure = w.Minute + 30
	// What a room takes now depends on who is standing in it. Every test built
	// on this world asserts money to the dollar and is about something else —
	// pressure, a paused job, a resumed one — so the laundry is held at the
	// city's ordinary three, where the room is worth exactly what it says.
	clearRoom(w, "laundry")
	fill(w, "laundry", TypicalRoom)
	return w
}
func TestPressureInterruptsAtScheduledMinute(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Advance(120)
	if w.Minute != 510 || w.Event == nil || w.Event.Kind != "business_pressure" || w.Event.Target != "laundry" {
		t.Fatal("demand missed boundary or wrong business")
	}
	data, _ := json.Marshal(w.Public())
	if strings.Contains(string(data), "next_pressure") {
		t.Fatal("private schedule disclosed")
	}
}
func TestPressurePaymentIsSpecificNotBlanketProtection(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Retaliation()
	w.Advance(30)
	before := w.Player.Cash
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "pay", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	// A week of what the place takes, rather than a flat sixty for a laundry
	// and a casino alike. The rest of the rule is unchanged and this still
	// asserts it: paying buys standing with that family and a day without
	// another demand, and the plot against the player personally is untouched,
	// because an understanding about a business is not protection for a man.
	share := w.TheirShare(w.Event.Target)
	if share <= 0 {
		t.Fatalf("the family asked for $%d", share)
	}
	if next.Player.Cash != before-share {
		t.Fatalf("paid $%d where the share is $%d", before-next.Player.Cash, share)
	}
	if next.Factions[0].Goodwill != 12 || next.NextPressure != next.Minute+1440 || len(next.Plots) != 1 {
		t.Fatal("payment rules or personal threat changed")
	}
}
func TestRefusingCreatesHiddenConsequences(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Advance(30)
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "resist", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Factions[0].Goodwill != -20 || next.Player.Respect != 2 {
		t.Fatalf("refusing left standing at %+d and respect at %d",
			next.Factions[0].Goodwill, next.Player.Respect)
	}
	// Whether they come for the place is three times in four, and this asked
	// for it every time on one campaign number. Measured across a spread of
	// them: the plan is laid in most cities and not in all, which is what the
	// card says — "the family may retaliate against your business".
	laid, cities := 0, 24
	for n := uint32(1); n <= uint32(cities); n++ {
		v := pressureWorldFor(spread(n))
		v.Advance(30)
		if v.Event == nil {
			cities--
			continue
		}
		after, err := Execute(v, Command{Kind: "choice", Event: v.Event.ID, Choice: "resist", Revision: v.Revision})
		if err != nil {
			t.Fatal(err)
		}
		if len(after.Plots) > 0 {
			laid++
		}
	}
	if cities < 12 {
		t.Fatalf("only %d cities put a demand at all, so this measures nothing", cities)
	}
	if laid*2 <= cities || laid == cities {
		t.Fatalf("a refused family laid a plan in %d of %d cities, which is not "+
			"\"may retaliate\"", laid, cities)
	}
	// And what a laid plan does, read in a city that laid one. The tail of this
	// used to run on whichever world the top of it happened to make, and a
	// refusal that did not draw the three-in-four left nothing to resolve.
	var plotted *World
	for n := uint32(1); n <= 24 && plotted == nil; n++ {
		v := pressureWorldFor(spread(n))
		v.Advance(30)
		if v.Event == nil {
			continue
		}
		after, err := Execute(v, Command{Kind: "choice", Event: v.Event.ID, Choice: "resist", Revision: v.Revision})
		if err != nil {
			t.Fatal(err)
		}
		if len(after.Plots) > 0 {
			plotted = after
		}
	}
	if plotted == nil {
		t.Fatal("no city laid a plan at all, so nothing can be watched resolving")
	}
	data, _ := json.Marshal(plotted.Public())
	if strings.Contains(string(data), "sabotage") {
		t.Fatal("hidden retaliation leaked")
	}
	condition := plotted.Properties["laundry"].Condition
	plotted.Advance(120)
	if got := plotted.Properties["laundry"].Condition; got >= condition {
		t.Fatalf("the plan resolved and the laundry went from %d%% to %d%%", condition, got)
	}
	if len(plotted.Plots) != 0 {
		t.Fatalf("%d plans are still waiting after the hour they were due", len(plotted.Plots))
	}
}
func TestBusinessDamageUsesAvailableCrewNotHomeGuards(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		crew bool
		busy bool
		want int
	}{{"unprotected", false, false, 65}, {"available crew", true, false, 85}, {"crew on job", true, true, 65}} {
		t.Run(tc.name, func(t *testing.T) {
			w := pressureWorld()
			w.Player.Security = 3
			if tc.crew {
				w.Player.Crew = []Crew{{"leo", "Leo", 65}}
			}
			if tc.busy {
				w.Tasks = []Task{{ID(), "Collections", 999}}
			}
			w.ResolveSabotage(Plot{Target: "laundry", Strength: 35})
			if w.Properties["laundry"].Condition != tc.want {
				t.Fatal("wrong protection source")
			}
		})
	}
}
func TestNoDemandWithoutBusiness(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.NextPressure = 500
	w.Advance(60)
	if w.Event != nil || w.NextPressure != 0 || w.Minute != 540 {
		t.Fatal("business demand without business")
	}
}

func TestBusinessPlotDoesNotPreventPersonalHit(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: w.Life, Actor: "russo", Target: "laundry", Due: 900}}
	w.Retaliation()
	w.Retaliation()
	if len(w.Plots) != 2 || w.Plots[1].Kind != "hit" {
		t.Fatal("different plots incorrectly deduplicated")
	}
}
func TestInvestigationNamesRealActorAndTarget(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Plots = []Plot{{ID: ID(), Kind: "sabotage", Life: w.Life, Actor: "russo", Target: "laundry", Due: 900}}
	w.Investigate()
	record := w.History[len(w.History)-1]
	if !strings.Contains(record.Text, "Russo Outfit") || !strings.Contains(record.Text, "Bluebird Laundry") || strings.Contains(record.Text, "Bellandi") || !w.Plots[0].Known {
		t.Fatal("investigation invented actor or target")
	}
}
func TestOpportunityDoesNotRevealHiddenPlans(t *testing.T) {
	t.Parallel()
	w := New(27)
	before, _ := json.Marshal(w.NextOpportunity())
	w.Retaliation()
	after, _ := json.Marshal(w.NextOpportunity())
	if string(before) != string(after) {
		t.Fatal("guidance leaked private threat")
	}
	w.Player.JobCount = 3
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 65
	if w.NextOpportunity().Target != "laundry" {
		t.Fatal("damage has no actionable guidance")
	}
}

func TestPoliceStopDefersJobRewardUntilDecision(t *testing.T) {
	t.Parallel()
	for _, decision := range []string{"pay", "abandon"} {
		t.Run(decision, func(t *testing.T) {
			w := New(27)
			w.Player.Heat = 14
			var err error
			w.Event, err = w.ValidateProposal(Proposal{"", "A risky delivery", "Please deliver this sealed package.", "mara", "courier", "Delivered.", "", nil})
			if err != nil {
				t.Fatal(err)
			}
			choice(t, &w, "accept")
			if w.Event == nil || w.Event.Kind != "police_stop" || w.Player.Cash != 90 || w.Minute != 525 || w.Player.Respect != 0 {
				t.Fatal("police stop failed to defer reward")
			}
			oldID := w.Event.ID
			w = w.Clone() // Resume the saved police decision after a reload.
			choice(t, &w, decision)
			if decision == "pay" && (w.Player.Cash != 125 || w.Player.Respect != 3 || w.Player.Heat != 7) {
				t.Fatal("paid completion incorrect")
			}
			if decision == "pay" && w.History[len(w.History)-1].Title != "A risky delivery" {
				t.Fatal("completion lost the original arrangement title")
			}
			if decision == "abandon" && (w.Player.Cash != 90 || w.Player.Respect != 0 || w.Player.Heat != 8) {
				t.Fatal("abandoned job paid reward")
			}
			if _, err = Execute(w, Command{Kind: "choice", Event: oldID, Choice: decision, Revision: w.Revision}); err == nil {
				t.Fatal("police choice applied twice")
			}
		})
	}
}

func TestGeneratedFactionWorkHasValidatedPoliticalEffect(t *testing.T) {
	t.Parallel()
	w := New(27)
	if _, err := w.ValidateProposal(Proposal{Title: "A job", Body: "Help them.", Speaker: "mara", Operation: "courier", Outcome: "Done.", Beneficiary: "invented-family"}); err == nil {
		t.Fatal("unknown faction accepted")
	}
	scene, err := w.ValidateProposal(Proposal{Title: "A Russo favor", Body: "Deliver these papers for Russo.", Speaker: "mara", Operation: "courier", Outcome: "Done.", Beneficiary: "russo"})
	if err != nil {
		t.Fatal(err)
	}
	// Who gains by the job holds however it is done, so it is stated with the
	// scene's other conditions rather than repeated on every button.
	if !strings.Contains(scene.Conditions, "Russo Outfit standing +6") {
		t.Fatalf("political stakes hidden: %q", scene.Conditions)
	}
	w.Event = scene
	choice(t, &w, "accept")
	if w.Factions[1].Goodwill != 6 || w.Factions[0].Goodwill != -3 {
		t.Fatal("completed faction job has no politics")
	}
	w.Player.Heat = 14
	w.Event = scene
	w.Event.ID = ID()
	choice(t, &w, "accept")
	if w.Event == nil || w.Event.Kind != "police_stop" || w.Event.Beneficiary != "russo" {
		t.Fatal("political context lost on interruption")
	}
	before := w.Factions[1].Goodwill
	choice(t, &w, "abandon")
	if w.Factions[1].Goodwill != before {
		t.Fatal("abandoned work granted favor")
	}
}

func TestAllyPressureReportsBothRelationships(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Advance(30)
	before := w.Player.Cash
	next, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: "ally", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Cash != before-35 || next.Factions[0].Goodwill != -12 || next.Factions[1].Goodwill != 12 {
		t.Fatal("incorrect alliance cost or standing")
	}
	record := next.History[len(next.History)-1]
	for _, phrase := range []string{"Russo Outfit", "Bellandi Family", "+12", "-12", "does not guarantee protection"} {
		if !strings.Contains(record.Text, phrase) {
			t.Fatalf("missing political feedback %q: %s", phrase, record.Text)
		}
	}
	if strings.Contains(record.Text, "sabotage") {
		t.Fatal("private retaliation revealed")
	}
}

func TestBookReviewDoesNotAdvanceCity(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Player.Location = "laundry"
	w.Properties["laundry"].Condition = 85
	next, err := Execute(w, Command{Kind: "inspect", Target: "laundry", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	if next.Minute != w.Minute || next.Player.Cash != w.Player.Cash || next.Event != nil {
		t.Fatal("inspection advanced economics or pressure")
	}
	if !strings.Contains(next.History[len(next.History)-1].Text, "$11/hour") {
		t.Fatal("missing actual income")
	}
}

func TestFriendlyHighIncomeDistrictDoesNotSilenceHostileFamily(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[1].Goodwill = 40
	w.Factions[0].Goodwill = -20
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "bellandi" || w.Event.Target != "laundry" {
		t.Fatal("friendly casino concealed hostile laundry claim")
	}
}
func TestHostileHighIncomeDistrictStillMakesItsOwnClaim(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[0].Goodwill = 40
	w.Factions[1].Goodwill = -20
	w.BusinessPressure()
	if w.Event == nil || w.Event.Actor != "russo" || w.Event.Target != "casino" {
		t.Fatal("wrong territorial claimant")
	}
}
func TestGoodRelationsInAllOwnedDistrictsPreservePeace(t *testing.T) {
	t.Parallel()
	w := pressureWorld()
	w.Properties["casino"].Owner = "player:1"
	w.Factions[0].Goodwill, w.Factions[1].Goodwill = 25, 25
	w.BusinessPressure()
	if w.Event != nil || len(w.Plots) != 0 || w.NextPressure <= w.Minute {
		t.Fatal("good relationships must still prevent demands")
	}
	w.Factions[0].Goodwill = 24
	w.Advance(720)
	if w.Event == nil || w.Event.Actor != "bellandi" {
		t.Fatal("claim did not resume after standing changed")
	}
}

func TestSustainedDefianceEscalatesForEitherFamily(t *testing.T) {
	t.Parallel()
	for _, actor := range []string{"bellandi", "russo"} {
		w := pressureWorld()
		e := &Scene{Actor: actor, Target: "laundry"}
		if err := w.ResolvePressure(e, "resist"); err != nil {
			t.Fatal(err)
		}
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				t.Fatal("first refusal became a personal hit")
			}
		}
		if err := w.ResolvePressure(e, "resist"); err != nil {
			t.Fatal(err)
		}
		hits := 0
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				hits++
				if p.Actor != actor || p.Due != w.Minute+240 {
					t.Fatal("wrong personal operation")
				}
			}
		}
		if hits != 1 {
			t.Fatal("sustained feud did not escalate")
		}
		w.RetaliationFrom(actor)
		count := 0
		for _, p := range w.Plots {
			if p.Kind == "hit" {
				count++
			}
		}
		if count != 1 {
			t.Fatal("personal operation duplicated")
		}
		data, _ := json.Marshal(w.Public())
		if strings.Contains(string(data), "\"plots\"") {
			t.Fatal("private operation leaked")
		}
	}
}

func TestRussoWarningNamesActualFamilyAndCanBeNegotiated(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.Player.Contacts = 2
	w.RetaliationFrom("russo")
	w.Advance(200)
	if w.Event == nil || w.Event.Kind != "warning" || !strings.Contains(w.Event.Body, "Russo") || strings.Contains(w.Event.Body, "Bellandi") {
		t.Fatal("wrong family in warning")
	}
	w.Event = nil
	w.RetaliationFrom("bellandi")
	if err := w.ResolveAudience(&Scene{Actor: "russo"}, "tribute"); err != nil {
		t.Fatal(err)
	}
	if len(w.Plots) != 1 || w.Plots[0].Actor != "bellandi" {
		t.Fatal("audience cancelled wrong family's threats")
	}
}

func TestFavorCanCreateRussoFeud(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.Factions[1].Goodwill = -28
	w.ResolveBeneficiary("bellandi")
	if len(w.Plots) != 1 || w.Plots[0].Actor != "russo" || w.Plots[0].Kind != "hit" {
		t.Fatal("Russo feud has no consequence")
	}
}

func TestKnownThreatsRevealOnlyDiscoveredCurrentLifePlans(t *testing.T) {
	t.Parallel()
	w := New(27)
	w.RetaliationFrom("russo")
	if len(w.KnownThreats()) != 0 {
		t.Fatal("undiscovered plot exposed")
	}
	w.Plots[0].Known = true
	w.Plots = append(w.Plots, Plot{Kind: "hit", Actor: "bellandi", Life: 0, Known: true})
	hints := w.KnownThreats()
	if len(hints) != 1 || !strings.Contains(hints[0], "Russo") {
		t.Fatal("wrong intelligence shown")
	}
	data, _ := json.Marshal(hints)
	if strings.Contains(string(data), "due") || strings.Contains(string(data), w.Plots[0].ID) {
		t.Fatal("private schedule or id exposed")
	}
	if err := w.ResolveAudience(&Scene{Actor: "russo"}, "tribute"); err != nil {
		t.Fatal(err)
	}
	if len(w.KnownThreats()) != 0 {
		t.Fatal("cancelled operation still shown")
	}
}

func TestDefensiveContributionReportsActualDamageAvoided(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		condition, strength, want int
		credited                  bool
	}{{100, 35, 85, true}, {5, 35, 0, false}, {100, 5, 95, false}, {100, 0, 100, false}} {
		w := New(27)
		w.Properties["laundry"].Owner = "player:1"
		w.Properties["laundry"].Condition = tc.condition
		w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
		w.ResolveSabotage(Plot{Actor: "bellandi", Target: "laundry", Strength: tc.strength})
		if w.Properties["laundry"].Condition != tc.want {
			t.Fatal("defense increased or miscounted damage", tc)
		}
		text := w.History[len(w.History)-1].Text
		if strings.Contains(text, "Leo Carver helped defend") != tc.credited {
			t.Fatal("incorrect defensive credit", text)
		}
	}
}

// A share is not a fee.
//
// The demand was sixty dollars flat — the same for a two-room laundry and the
// best casino in the city, and the same whether you held one address or six —
// while the family's own words in the scene are "my people expect a share". A
// publican who ends a campaign with $27,516 paid nine of those over a month,
// which is two per cent of what they made and nothing at all to decide about.
func TestAFamilyAsksForAShareAndNotAFee(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 50000

	// A small place and a big one.
	small, big := "", ""
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 || l.Cost <= 0 {
			continue
		}
		if small == "" || prop.Income < w.Properties[small].Income {
			small = l.ID
		}
		if big == "" || prop.Income > w.Properties[big].Income {
			big = l.ID
		}
	}
	if small == "" || small == big {
		t.Fatal("this city has no two businesses of different sizes")
	}
	cheap, dear := w.TheirShare(small), w.TheirShare(big)
	t.Logf("%s earns %d a day and they want $%d; %s earns %d and they want $%d",
		small, w.Properties[small].Income, cheap, big, w.Properties[big].Income, dear)
	if cheap >= dear {
		t.Fatalf("they want $%d for %s and $%d for %s, which is a fee rather than a share",
			cheap, small, dear, big)
	}
	// Never cheaper than the flat demand used to be.
	if cheap < PressureFloor {
		t.Fatalf("a family settled for $%d where the old demand was $%d", cheap, PressureFloor)
	}
	// And a claim on one shop stays a claim on one shop.
	if dear > PressureCeiling {
		t.Fatalf("a family asked $%d for an understanding about a single address", dear)
	}
	// Somewhere that earns nothing is still worth the floor to be left alone.
	if quiet := w.TheirShare("precinct"); quiet != PressureFloor {
		t.Fatalf("they want $%d about a police station", quiet)
	}
}
