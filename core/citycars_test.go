package core

import "testing"

// Nobody in this city owned a car but the player. The forecourt had exactly one
// customer, no car could be stolen from anybody, and there was nothing for a
// garage to repair. Both of the inbox entries about cars stand on this and
// neither could be finished without it.

func TestPeopleInThisCityOwnCars(t *testing.T) {
	w := New(53)
	w.District = 2
	driving, people := 0, 0
	for _, n := range w.People() {
		people++
		if n.Car > 0 {
			driving++
		}
	}
	if people == 0 {
		t.Fatal("nobody lives here")
	}
	if driving == 0 {
		t.Fatal("nobody in the city owns a car")
	}
	if driving == people {
		t.Error("everybody in the city drives, which is not a 1950s city")
	}
	t.Logf("%d of %d people own a car", driving, people)
}

// Who drives is who could afford one. A family head drives; somebody mending
// nets at the docks does not.
func TestWhoDrivesIsWhoCouldAffordTo(t *testing.T) {
	w := New(53)
	w.District = 2
	high, low := 0, 0
	highDriving, lowDriving := 0, 0
	for _, n := range w.People() {
		if n.Rank >= RankLieutenant {
			high++
			if n.Car > 0 {
				highDriving++
			}
			continue
		}
		if n.Faction == "" && !IsOfficial(n.ID) {
			low++
			if n.Car > 0 {
				lowDriving++
			}
		}
	}
	if high == 0 || low == 0 {
		t.Fatalf("the city has %d people of standing and %d on the street", high, low)
	}
	if highDriving*low <= lowDriving*high {
		t.Errorf("%d of %d people of standing drive and %d of %d on the street do",
			highDriving, high, lowDriving, low)
	}
	t.Logf("%d of %d of standing drive; %d of %d on the street do", highDriving, high, lowDriving, low)
}

// And the forecourt has customers. Somebody who can afford a car and has none
// buys one, and the money reaches whoever holds the lot — which is the whole
// point of the dealership existing at all.
func TestTheCityBuysCarsAndTheForecourtTakesTheMargin(t *testing.T) {
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
		}
	}
	w := New(53)
	w.District = 2
	w.Properties[lot].Owner = "bellandi"
	house := w.faction("bellandi")
	if house == nil {
		t.Fatal("no bellandi")
	}
	// Somebody of standing, with money, and nothing to drive.
	var buyer *NPC
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Rank >= RankLieutenant && !n.Dead {
			buyer = n
			break
		}
	}
	if buyer == nil {
		t.Fatal("nobody of standing in the city")
	}
	buyer.Car, buyer.Purse = 0, 4000

	before, purse := house.Cash, buyer.Purse
	sold := 0
	for day := 0; day < 40 && buyer.Car == 0; day++ {
		w.Minute = day * 1440
		w.CarTrade()
		if buyer.Car > 0 {
			sold++
		}
	}
	if buyer.Car == 0 {
		t.Fatal("somebody of standing with four thousand dollars never bought a car in forty days")
	}
	if buyer.Purse >= purse {
		t.Errorf("they bought a car and are carrying $%d, up from $%d", buyer.Purse, purse)
	}
	if house.Cash <= before {
		t.Errorf("a car was sold off Bellandi's forecourt and they took $%d", house.Cash-before)
	}
	t.Logf("one sale: the buyer went from $%d to $%d and the lot took $%d", purse, buyer.Purse, house.Cash-before)
}

// The wiring, at both ends. A city that only settles cars when it is MADE
// leaves every campaign that predates them without a single driver, and a
// trade that is never called sells nothing.
func TestAnOldCampaignGetsTheCityOnTheRoad(t *testing.T) {
	w := New(53)
	w.Version = SaveVersion
	for i := range w.NPCs {
		w.NPCs[i].Car, w.NPCs[i].Drove = 0, 0
	}
	w.SettleCars()
	driving := 0
	for _, n := range w.People() {
		if n.Car > 0 {
			driving++
		}
	}
	if driving == 0 {
		t.Fatal("a campaign that predates cars has nobody driving after the repair")
	}
}

// Somebody whose car was taken is not quietly handed another by the settling.
func TestSettlingDoesNotReplaceACarThatWasTaken(t *testing.T) {
	w := New(53)
	w.District = 2
	var driver *NPC
	for i := range w.NPCs {
		if w.NPCs[i].Car > 0 {
			driver = &w.NPCs[i]
			break
		}
	}
	if driver == nil {
		t.Fatal("nobody drives")
	}
	driver.Car = 0 // taken, burned, whatever happened to it
	w.SettleCars()
	if driver.Car != 0 {
		t.Errorf("%s lost a car and the city handed them another", driver.Name)
	}
}

// And the day actually calls the trade. This is the seam that has caught me
// before: a function that works and is never reached.
func TestTheDayPutsCarsOnTheRoad(t *testing.T) {
	w := New(53)
	w.District = 2
	lot := ""
	for _, l := range Locations {
		if l.Kind == "dealer" {
			lot = l.ID
		}
	}
	w.Properties[lot].Owner = "bellandi"
	// Give everybody of standing money and nothing to drive.
	waiting := 0
	for i := range w.NPCs {
		if n := &w.NPCs[i]; w.WouldDrive(n) {
			n.Car, n.Purse = 0, 5000
			waiting++
		}
	}
	if waiting < 3 {
		t.Fatalf("only %d people would drive, so this proves little", waiting)
	}
	live(t, w, 48*20) // twenty days of the city running itself
	bought := 0
	for _, n := range w.People() {
		if n.Car > 0 {
			bought++
		}
	}
	if bought == 0 {
		t.Fatal("twenty days passed and the city bought no cars at all")
	}
	t.Logf("%d of %d who wanted one bought a car over twenty days", bought, waiting)
}
