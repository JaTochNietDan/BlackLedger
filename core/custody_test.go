package core

import "testing"

// A campaign with something in it worth charging somebody with.
func chargeable(t *testing.T) *World {
	t.Helper()
	w, _ := testator(t)
	w.Player.Cash = 20000
	w.Player.Heat = 60
	if err := w.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	return w
}

func TestAFineIsNotACaseAndACaseIsNotAFine(t *testing.T) {
	t.Parallel()
	w, _ := testator(t)
	if weight, _ := w.Charge(); weight != 0 {
		t.Fatalf("a search of a laundry with nothing in it was worth %d", weight)
	}
	w.Player.Cash = 20000
	if err := w.BuildStill("laundry"); err != nil {
		t.Fatal(err)
	}
	weight, because := w.Charge()
	if weight < ChargeMinimum {
		t.Fatalf("a still was worth %d, below the %d it takes to charge anybody", weight, ChargeMinimum)
	}
	if because == "" {
		t.Fatal("nobody could say what they found")
	}
	// Attention lengthens it: the same evidence against a man they have been
	// watching is a longer sentence.
	w.Player.Heat = 0
	quiet := w.Sentence(weight)
	w.Player.Heat = 95
	if watched := w.Sentence(weight); watched <= quiet {
		t.Fatalf("a watched man got %d days for what a quiet one got %d for", watched, quiet)
	}
	if w.Sentence(100) > SentenceCeiling || w.Sentence(0) < SentenceFloor {
		t.Fatal("a sentence escaped its bounds")
	}
}

func TestARaidThatFindsACaseTakesAPerson(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Raid()
	if w.Event == nil || w.Event.Kind != "arrest" {
		t.Fatal("a raid that found a still did not produce an arrest")
	}
	if w.Held() {
		t.Fatal("nobody was given the choice; they were simply taken")
	}
	ids := map[string]bool{}
	for _, c := range w.Event.Choices {
		ids[c.ID] = true
	}
	if !ids["serve"] || !ids["pay"] {
		t.Fatalf("the arrest offered %v", ids)
	}
}

func TestSomebodyRetainingACommissionerIsNotCharged(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Player.Retainers = []string{"commissioner"}
	if w.RaidRelief() == 0 {
		t.Skip("a retained commissioner is worth nothing here")
	}
	w.Raid()
	if w.Event != nil && w.Event.Kind == "arrest" {
		t.Fatal("a man paying the commissioner was charged anyway")
	}
}

func TestGoingWithThemCostsTheDaysAndPaysTheStreet(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Raid()
	e := w.Event
	days := e.Effect.Reward
	respect, minute := w.Player.Respect, w.Minute
	w.Event = nil
	if err := w.ResolveArrest(e, "serve"); err != nil {
		t.Fatal(err)
	}
	if !w.Held() || w.Player.Location != "precinct" {
		t.Fatal("they went with them and nothing happened")
	}
	if w.Player.Weapon != 0 || w.Player.Charges != 0 {
		t.Fatal("they walked into a cell carrying")
	}
	if w.DaysLeft() != days {
		t.Fatalf("%d days left of a %d day sentence", w.DaysLeft(), days)
	}
	// While inside, the map offers nothing anywhere else.
	for _, l := range Locations {
		if l.ID == "precinct" {
			continue
		}
		if len(w.Actions(l.ID)) > 0 {
			t.Fatalf("a man in a cell could still act at %s", l.ID)
		}
	}
	inside := map[string]bool{}
	for _, a := range w.Actions("precinct") {
		inside[a.ID] = true
	}
	if !inside["sit_out"] || !inside["talk"] || !inside["lawyer"] {
		t.Fatalf("the cell offered %v", inside)
	}

	if err := w.SitOut(); err != nil {
		t.Fatal(err)
	}
	if w.Held() {
		t.Fatal("they did the time and are still in there")
	}
	if w.Minute-minute < days*1440 {
		t.Fatalf("%d days passed for a %d day sentence", (w.Minute-minute)/1440, days)
	}
	if w.Player.Alive && w.Player.Respect <= respect {
		t.Fatalf("doing %d days quietly was worth %d respect", days, w.Player.Respect-respect)
	}
	if w.Player.Alive && w.Player.Heat > ReleasedHeat {
		t.Fatalf("they came out at %d attention", w.Player.Heat)
	}
}

func TestTalkingIsTheFastWayOutAndTheOnlyOneNobodyForgets(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Raid()
	e := w.Event
	w.Event = nil
	if err := w.ResolveArrest(e, "serve"); err != nil {
		t.Fatal(err)
	}
	w.Player.Respect = 100
	goodwill := w.Factions[0].Goodwill
	minute := w.Minute
	if err := w.Talk(); err != nil {
		t.Fatal(err)
	}
	if w.Held() {
		t.Fatal("they talked and stayed in")
	}
	if w.Minute != minute {
		t.Fatal("talking took days")
	}
	if w.Player.Respect != 100-TalkedRespect {
		t.Fatalf("talking cost %d respect", 100-w.Player.Respect)
	}
	if w.Factions[0].Goodwill >= goodwill {
		t.Fatal("no organization in the city minded")
	}
	for _, n := range w.OwnPeople() {
		if n.Trust != 0 {
			t.Fatalf("%s still trusts a man who did that", n.Name)
		}
	}
}

func TestGivingThemSomebodyElseIsSomethingEverybodyElseFindsOut(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	// A second man, so there is somebody left to have an opinion about it.
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) || len(w.OwnPeople()) >= 2 {
			continue
		}
		w.Player.Location = n.Location
		w.SignOn(n.ID)
	}
	if len(w.OwnPeople()) < 2 {
		t.Skip("only one man would sign on")
	}
	for _, n := range w.OwnPeople() {
		n.Trust = 60
	}
	w.Player.Location = "laundry"
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	w.Raid()
	e := w.Event
	w.Event = nil
	fall := w.FallGuy()
	if fall == nil {
		t.Fatal("nobody could be given to them")
	}
	if err := w.ResolveArrest(e, "fall:"+fall.ID); err != nil {
		t.Fatal(err)
	}
	if w.Held() {
		t.Fatal("they gave them somebody and went in anyway")
	}
	if !w.Inside(fall) {
		t.Fatal("the man they were given is not inside")
	}
	if w.PostedAt("laundry") != nil {
		t.Fatal("a man in a cell is still on the door")
	}
	for _, n := range w.OwnPeople() {
		if n.ID != fall.ID && n.Trust >= 60 {
			t.Fatalf("%s thought no less of it", n.Name)
		}
	}
	// He is not available for anything until they let him go.
	if len(w.Unposted()) != 0 && w.Unposted()[0].ID == fall.ID {
		t.Fatal("a man inside was offered for a door")
	}

	// And he can be bought back out, which he notices.
	w.Player.Cash, w.Player.Location = 20000, "precinct"
	trust := fall.Trust
	bail := map[string]bool{}
	for _, a := range w.Actions("precinct") {
		bail[a.ID] = true
	}
	if !bail["bail:"+fall.ID] {
		t.Fatalf("the station offered %v", bail)
	}
	if err := w.Bail(fall.ID); err != nil {
		t.Fatal(err)
	}
	if w.Inside(fall) || fall.Trust <= trust {
		t.Fatal("he came out no better disposed")
	}
}

func TestTheDoorOpensByItselfIfNobodyOpensIt(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Confine(3, "a still")
	if !w.Held() {
		t.Fatal("they are not inside")
	}
	// The clock passing the morning is enough; nothing has to be pressed.
	w.Minute = w.Player.HeldUntil
	w.CustodyDay()
	if w.Held() || w.Player.HeldFor != "" {
		t.Fatal("the door never opened")
	}
	if w.Player.Location != w.Player.Home {
		t.Fatalf("they came out at %s", w.Player.Location)
	}
}

func TestNobodyBringsWorkToAManInACell(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	w.Confine(4, "a still")
	w.Offers = []Offer{{Ready: w.Minute, Event: &Scene{ID: ID(), Kind: "job", Title: "Something to do", Body: "b", Choices: []Choice{{ID: "accept", Label: "Accept"}}}}}
	w.OfferIfReady()
	if w.Event != nil {
		t.Fatal("somebody brought a job to a cell")
	}
	// And the police do not come to the door of a man they are already holding.
	w.Player.Heat = 100
	for i := 0; i < 50; i++ {
		w.considerRaid()
	}
	if w.Event != nil {
		t.Fatal("they raided a man they had in a cell")
	}
}

func TestNobodyBringsTheTakingsToAStationExceptYourOwnMan(t *testing.T) {
	t.Parallel()
	w := chargeable(t)
	if w.CollectionShare("laundry") != 1 {
		t.Fatal("a man standing in his own business did not get his own takings")
	}
	w.Confine(5, "a still")
	if w.CollectionShare("laundry") != 0 {
		t.Fatalf("a man in a cell collected %v of the till", w.CollectionShare("laundry"))
	}
	// Somebody of yours on the door is the reason any of it survives the week.
	w.Player.HeldUntil = 0
	w.Player.Location = "laundry"
	if err := w.Post("laundry"); err != nil {
		t.Fatal(err)
	}
	w.Confine(5, "a still")
	if share := w.CollectionShare("laundry"); share != HeldCollection {
		t.Fatalf("a manned door was worth %v of the till while its owner was inside", share)
	}
	if w.CollectionShare("garage") != 0 {
		t.Fatal("a door with nobody on it paid out anyway")
	}
}

func TestACellIsTheOneThingItIsGoodFor(t *testing.T) {
	t.Parallel()
	// Nobody is robbed in the street, mugged, or visited at home while the
	// police have them. It is the only protection in this game that is free.
	w := chargeable(t)
	w.Confine(6, "a still")
	cash := w.Player.Cash
	for i := 0; i < 200; i++ {
		w.ConsiderRobbery()
	}
	if w.Player.Cash != cash {
		t.Fatalf("a man in a cell was robbed in the street of $%d", cash-w.Player.Cash)
	}
}
