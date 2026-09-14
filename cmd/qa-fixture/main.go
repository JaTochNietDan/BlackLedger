// qa-fixture creates a new, isolated save for a repeatable browser edge-case test.
// It refuses to overwrite any existing file, including the normal campaign.
package main

import (
	"blackledger/core"
	"blackledger/store"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatal("usage: go run ./cmd/qa-fixture <new-qa.sqlite3> [butcher|ashbury|cypress|burglary|apartments|poker|building-driveby|city3d-dusk|fatal-charge-car|survivor-charge|fatal-charge|faction-planter|strike-unarmed|strike-revolver|strike-shotgun|strike-thompson|raid-presence|city3d-rain|gunfight-killing|gunfight|city3d-walk|city3d-junction|city3d-traffic|city3d|city3d-night|city3d-blast|police|damage|warning|russo-warning|attack|voice|contact|paused-job|leader|doorman|arrest|debt|herald|killing|dead|offer|audience|street|room|gone|post|round|bereaved|inside|writeoff|writeoff-dead|worn|tables|bench|petrol]")
	}
	scenario := "police"
	if len(os.Args) == 3 {
		scenario = os.Args[2]
	}
	weaponTier, weaponFixture := map[string]int{"strike-unarmed": 0, "strike-revolver": 1, "strike-shotgun": 2, "strike-thompson": 3}[scenario]
	if !weaponFixture && scenario != "butcher" && scenario != "ashbury" && scenario != "cypress" && scenario != "burglary" && scenario != "apartments" && scenario != "poker" && scenario != "building-driveby" && scenario != "city3d-dusk" && scenario != "fatal-charge-car" && scenario != "survivor-charge" && scenario != "fatal-charge" && scenario != "faction-planter" && scenario != "raid-presence" && scenario != "city3d-rain" && scenario != "gunfight-killing" && scenario != "gunfight" && scenario != "city3d-walk" && scenario != "city3d-junction" && scenario != "city3d-traffic" && scenario != "city3d-blast" && scenario != "city3d" && scenario != "city3d-night" && scenario != "police" && scenario != "damage" && scenario != "warning" && scenario != "russo-warning" && scenario != "attack" && scenario != "voice" && scenario != "contact" && scenario != "paused-job" && scenario != "leader" && scenario != "doorman" && scenario != "arrest" && scenario != "debt" && scenario != "herald" && scenario != "killing" && scenario != "dead" && scenario != "offer" && scenario != "audience" && scenario != "street" && scenario != "room" && scenario != "gone" && scenario != "post" && scenario != "round" && scenario != "bereaved" && scenario != "inside" && scenario != "writeoff" && scenario != "writeoff-dead" && scenario != "worn" && scenario != "tables" && scenario != "bench" && scenario != "petrol" {
		log.Fatal("unsupported QA scenario")
	}
	path := os.Args[1]
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Fatal("QA output must be a new file; existing saves are never overwritten")
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	s, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer s.DB.Close()
	err = s.Change(func(w *core.World) error {
		if scenario == "building-driveby" {
			candidate := core.New(61)
			candidate.MigrateLivingWorld()
			candidate.Minute = 600
			candidate.Player.Location = "club"
			candidate.Player.Cash, candidate.Player.Respect = 20000, 40
			candidate.Player.Weapon, candidate.Player.Car, candidate.Player.CarWear = 3, 3, 100
			candidate.Player.Fuel, candidate.Player.Fuelled = 30, 600
			candidate.Player.Crew = []core.Crew{{ID: "leo", Name: candidate.NPC("leo").Name, Loyalty: 90}}
			driver := candidate.NPC("leo")
			driver.Location, driver.Heading, driver.Arrives, driver.Sets = "club", "", 0, 0
			candidate.Tasks, candidate.Plots = nil, nil
			candidate.Properties["club"].Condition = 100
			*w = *candidate
			return nil
		}
		if scenario == "fatal-charge" || scenario == "survivor-charge" || scenario == "fatal-charge-car" {
			// Probe disposable worlds, then save the untouched pre-command setup.
			// The browser must commit the actual action to exercise the death flow.
			setup := func(seed uint32) *core.World {
				candidate := core.New(61)
				candidate.MigrateLivingWorld()
				candidate.Minute = 600
				candidate.Player.Location, candidate.Player.Car = "club", 0
				if scenario == "fatal-charge-car" {
					candidate.Player.Car, candidate.Player.CarWear = 1, 100
					candidate.Player.Fuel, candidate.Player.Fuelled = 30, candidate.Minute
				}
				candidate.Player.Cash, candidate.Player.Respect = 20000, 40
				candidate.Player.Health, candidate.Player.Charges = 40, 1
				candidate.RNG = seed * 2654435761
				return candidate
			}
			for seed := uint32(1); seed <= 128; seed++ {
				probe := setup(seed)
				if err := probe.Plant("club"); err != nil {
					return err
				}
				if len(probe.VisualCues) > 0 && probe.VisualCues[len(probe.VisualCues)-1].Accident != nil && probe.Player.Alive == (scenario == "survivor-charge") {
					*w = *setup(seed)
					return nil
				}
			}
			return fmt.Errorf("no matching charge accident seed found")
		}
		if scenario == "faction-planter" {
			w.Minute, w.Player.Location, w.Player.Car = 600, "laundry", 0
			w.Properties["laundry"].Owner = fmt.Sprintf("player:%d", w.Life)
			w.Plots = nil
			for i := range w.Conflicts {
				w.Conflicts[i].State = "cold"
			}
			for i := range w.Factions {
				f := &w.Factions[i]
				f.Cash = 0
				if f.ID == "bellandi" {
					f.Cash, f.Goodwill = 30000, -80
				}
			}
			var planter *core.NPC
			for i := range w.NPCs {
				if w.NPCs[i].Faction == "bellandi" && !w.NPCs[i].Dead {
					planter = &w.NPCs[i]
					break
				}
			}
			if planter == nil {
				return fmt.Errorf("fixture has no Bellandi member")
			}
			planter.Location, planter.Home, planter.Heading = "laundry", "room", ""
			planter.Held, planter.Sets, planter.Arrives, planter.Car = 0, 0, 0, 0
			w.NPCs = []core.NPC{*planter}
			w.VisualCues = nil
			for i := 0; i < 2000 && len(w.VisualCues) == 0; i++ {
				w.DemolitionDay()
			}
			if len(w.VisualCues) == 0 {
				return fmt.Errorf("no faction blast resolved")
			}
			w.LastResult = &core.Result{Kind: "qa-faction-planter", From: "laundry", To: "laundry", Cues: w.VisualCues}
			return nil
		}
		if scenario == "city3d-dusk" || scenario == "city3d-rain" || scenario == "city3d-walk" || scenario == "city3d-junction" || scenario == "city3d-traffic" || scenario == "city3d" || scenario == "city3d-night" || scenario == "city3d-blast" {
			w.Player.Cash, w.Player.Respect, w.District = 12000, 60, 2
			w.Player.Location, w.Player.Car, w.Player.CarWear = "bar", 2, 95
			w.SettleFuel()
			w.Minute = 600
			if scenario == "city3d-rain" {
				for i := 0; ; i++ {
					w.ID = fmt.Sprintf("isolated-city-rain-%d", i)
					if w.Sky().Kind == "rain" {
						break
					}
				}
			}
			if scenario == "city3d-night" {
				w.Minute = 1260
			}
			if scenario == "city3d-dusk" {
				w.Minute, w.Player.Car = 1195, 0
			}
			w.Properties["laundry"].Owner = "player:1"
			w.Plots = []core.Plot{{ID: "city3d-visible-test", Kind: "sabotage", Life: w.Life, Due: w.Minute + 5, Actor: "bellandi", Target: "laundry", Strength: 35}}
			if scenario == "city3d-dusk" {
				w.Plots = nil
			}
			if scenario == "city3d-blast" {
				w.Player.Location, w.Player.Charges, w.RNG = "club", 2, 1
				w.Plots = nil
			}
			for i := 0; i < 12 && i < len(w.NPCs); i++ {
				n := &w.NPCs[i]
				n.Location = core.Locations[i%len(core.Locations)].ID
				n.Heading = core.Locations[(i+3)%len(core.Locations)].ID
				n.Sets, n.Car, n.Dry, n.Hurt = 0, i%4, false, false
				if scenario == "city3d-walk" {
					w.Player.Car, n.Car = 0, 0
					w.Plots = nil
				}
				n.Errand = "Crossing town in an isolated city presentation fixture"
				n.Arrives = w.Minute + max(1, core.TravelMinutes(n.Location, n.Heading)*2/3)
				if scenario == "city3d-junction" {
					pairs := [][2]string{{"laundry", "pawn"}, {"pawn", "laundry"}, {"casino", "market"}, {"market", "casino"}}
					pair := pairs[i%4]
					n.Location, n.Heading, n.Car = pair[0], pair[1], 1+i%3
					progress := []float64{.2, .45, .2, .45}[i%4] - float64(i/4)*.07
					n.Arrives = w.Minute + max(1, int(float64(core.TravelMinutes(n.Location, n.Heading))*(1-progress)))
					w.Plots = nil
				}
				if scenario == "city3d-traffic" {
					w.Player.Location = "tailor"
					n.Location, n.Heading, n.Car = "tailor", "dealer", 1+i%3
					n.Arrives = w.Minute + max(1, core.TravelMinutes(n.Location, n.Heading)/2)
					w.Plots = nil
				}
			}
			return nil
		}
		if scenario == "worn" {
			// A business of the player's knocked about but still standing, so the
			// screen has to say what it is really earning rather than what a sound
			// one would.
			w.Player.Cash, w.Player.Respect = 4000, 20
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["laundry"].Condition = 80
			w.Player.Location = "laundry"
			return nil
		}
		if scenario == "writeoff" || scenario == "writeoff-dead" {
			// Money out with a man the city then kills, so the ledger has an
			// asset that stops existing while the player is looking at it.
			w.Player.Cash, w.Player.Respect = 5000, core.OrganizationStanding
			w.Player.Location = "bar"
			for _, n := range w.PeopleHere("bar") {
				if w.LendReadiness(n.ID) == "" {
					if err := w.Lend(n.ID); err != nil {
						return err
					}
					if scenario == "writeoff-dead" {
						w.Kill(n.ID, "Shot over something that had nothing to do with the money.")
					}
					return nil
				}
			}
			return fmt.Errorf("nobody at the bar would borrow")
		}
		if scenario == "inside" {
			// The player's crew member in a police cell, which used to be no
			// obstacle to sending him out to rob somebody.
			w.Player.Cash, w.Player.Respect = 2500, 20
			w.Player.Location = "bar"
			next, err := core.Execute(w, core.Command{RequestID: core.ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
			if err != nil {
				return err
			}
			*w = *next
			w.Player.Location = "bar"
			if n := w.NPC(w.Player.Crew[0].ID); n != nil {
				n.Held = w.Minute + 3*1440
				n.Location = "precinct"
			}
			return nil
		}
		if scenario == "bereaved" {
			// A crew member the player recruited and the city then killed, which
			// used to leave him on the books, offered work and paid for.
			w.Player.Cash, w.Player.Respect = 2500, 20
			w.Player.Location = "bar"
			next, err := core.Execute(w, core.Command{RequestID: core.ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
			if err != nil {
				return err
			}
			*w = *next
			w.Player.Location = "bar"
			w.Kill(w.Player.Crew[0].ID, "Shot twice outside the Mariner, over nothing anybody will name.")
			return nil
		}
		if scenario == "round" {
			// A crew member at the bar and a laundry of the player's across the
			// district, so sending him on collections is a walk to watch.
			w.Player.Cash, w.Player.Respect = 2500, 20
			w.Player.Location = "bar"
			w.Properties["laundry"].Owner = "player:1"
			next, err := core.Execute(w, core.Command{RequestID: core.ID(), Revision: w.Revision, Kind: "recruit", Target: "bar"})
			if err != nil {
				return err
			}
			*w = *next
			w.Player.Location = "bar"
			if n := w.NPC(w.Player.Crew[0].ID); n != nil {
				n.Location = "bar"
			}
			return nil
		}
		if scenario == "post" {
			// The player standing in their own laundry with their one man across
			// the city, so putting him on the door is a journey he has to make.
			w.Player.Cash, w.Player.Respect, w.Player.Contacts = 9000, core.OrganizationStanding, 3
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["garage"].Owner = "player:1"
			w.OrganizationDay()
			for _, n := range w.Civilians() {
				if core.IsOfficial(n.ID) {
					continue
				}
				w.Player.Location = n.Location
				if w.SignOn(n.ID) == nil {
					break
				}
			}
			for _, n := range w.OwnPeople() {
				n.Location = "club" // the far side of the district
			}
			w.Player.Location = "laundry"
			return nil
		}
		if scenario == "gone" {
			// The player standing in the bar with somebody they deal with, five
			// minutes short of the half-day when she is due somewhere else. One
			// action and she is out on the street rather than across the table.
			w.Player.Cash, w.Player.Respect = 4000, 30
			w.Player.Location = "bar"
			if m := w.NPC("mara"); m != nil {
				m.Location, m.Role = "bar", "Runs Russo Motor Works"
			}
			w.Minute = 715
			return nil
		}
		if scenario == "petrol" {
			// A car most of the way through a tank, standing on a forecourt, so
			// the pumps have something to sell and the gauge has something to
			// say.
			w.Player.Cash, w.Player.Respect = 2500, 30
			w.Player.Car, w.Player.CarWear = 2, 90
			w.SettleFuel()
			w.Burn(600)
			w.Player.Location = "filling"
			return nil
		}
		if scenario == "bench" {
			// A city where the glass has been going: people with broken cars and
			// the fee in their pockets, so a garage has somebody at the counter
			// who is there for a reason rather than an empty room with a trade
			// figure attached to it.
			w.Player.Cash, w.Player.Respect = 3000, 40
			w.District = 2
			w.Player.Location = "garage"
			w.Properties["garage"].Owner = "player:1"
			broken := 0
			for i := range w.NPCs {
				n := &w.NPCs[i]
				if n.Dead || n.Car == 0 || broken >= 4 {
					continue
				}
				n.Hurt, n.Purse = true, 400
				broken++
			}
			if broken == 0 {
				return fmt.Errorf("nobody in this city drives")
			}
			// Two half-days, so they set off and arrive.
			for half := 0; half < 3; half++ {
				w.SetOut()
				w.Minute += 720
				w.Arrivals()
			}
			return nil
		}
		if scenario == "cypress" {
			w.District = 2
			w.Player.Cash, w.Player.Respect = 12000, 60
			w.Player.Home, w.Player.Location = "estate", "estate"
			w.Properties["estate"].Owner = "player:1"
			w.Event, w.Plots, w.Tasks = nil, nil, nil
			w.Minute = 1020
			for i := 0; i < 7 && i < len(w.NPCs); i++ {
				n := &w.NPCs[i]
				n.Location, n.Heading, n.Arrives, n.Sets = "estate", "", 0, 0
				w.MeetPerson(n.ID)
			}
			w.SettleApartments()
			return nil
		}
		if scenario == "burglary" {
			w.Player.Cash, w.Player.Respect = 1000, 40
			w.Player.Location = "mercercourt"
			w.District = 2
			w.Event, w.Plots, w.Tasks = nil, nil, nil
			n := w.NPC("mara")
			n.Home, n.Location, n.Heading = "mercercourt", "bar", ""
			n.Purse = 90
			n.Faction = ""
			w.MeetPerson(n.ID)
			w.HouseholdSavings = map[string]core.HouseholdAccount{n.ID: {Cash: 180, Day: 1}}
			w.SettleApartments()
			w.RNG = 1
			return nil
		}
		if scenario == "butcher" {
			w.Player.Location, w.Player.Cash = "butcher", 6000
			w.District, w.Minute = 9, 600
			w.Event, w.Plots, w.Tasks = nil, nil, nil
			for i := 0; i < 6 && i < len(w.NPCs); i++ {
				n := &w.NPCs[i]
				n.Location, n.Heading, n.Arrives, n.Sets = "butcher", "", 0, 0
				if i == 0 {
					n.Role = "Butcher"
				}
				w.MeetPerson(n.ID)
			}
			return nil
		}
		if scenario == "apartments" || scenario == "ashbury" {
			w.Player.Cash, w.Player.Respect = 6000, 30
			w.Player.Home, w.Player.Location = "mercercourt", "mercercourt"
			w.District = 2
			w.Event, w.Plots, w.Tasks = nil, nil, nil
			if scenario == "ashbury" {
				w.Player.Home, w.Player.Location = "apartment", "apartment"
				for i := 0; i < 5 && i < len(w.NPCs); i++ {
					n := &w.NPCs[i]
					n.Location, n.Heading, n.Arrives, n.Sets = "apartment", "", 0, 0
					w.MeetPerson(n.ID)
				}
			}
			w.SettleApartments()
			return nil
		}
		if scenario == "poker" {
			candidate := core.New(61)
			candidate.Event, candidate.District = nil, 9
			candidate.Player.Health, candidate.Player.Respect = 100, 30
			candidate.Player.Cash, candidate.Player.Location = 5000, "bar"
			for candidate.Minute%1440 < 1200 {
				candidate.Event = nil
				candidate.Advance(60)
				candidate.Event = nil
			}
			if err := candidate.Sit("bar", core.Backroom); err != nil {
				return err
			}
			if err := candidate.SitInTheBackRoom("bar", 600); err != nil {
				return err
			}
			*w = *candidate
			return nil
		}
		if scenario == "tables" {
			// The player sitting down at a casino with money on them, for looking
			// at the tables themselves: the cards and the wheel take the screen
			// and hold it until the player gets up, so this is the one fixture
			// where the point is what the game looks like while you play it.
			w.Player.Cash, w.Player.Respect = 4000, 60
			w.District = 2
			w.Player.Location = "casino"
			return nil
		}
		if scenario == "room" {
			// The player standing inside a business while somebody walks out of it
			// on an errand, so the room has traffic to remark on.
			w.Player.Cash, w.Player.Respect = 4000, 30
			w.Properties["laundry"].Owner = "player:1"
			w.Player.Location = "laundry"
			// Somebody standing in the laundry who is due somewhere else.
			n := w.NPC("leo")
			if n == nil {
				return fmt.Errorf("nobody to move")
			}
			n.Location = "laundry"
			// Mara stands at the bar with an errand of her own, so the player can
			// watch somebody they were dealing with become unreachable.
			if m := w.NPC("mara"); m != nil {
				m.Role = "Runs Russo Motor Works"
			}
			// Five minutes short of midday: the city's people set off on the
			// half-day, so any ordinary action reaches it.
			w.Minute = 715
			return nil
		}
		if scenario == "street" {
			// Ground changing hands is what puts people on the street, so this
			// takes a holding off one family and gives it to the other, then lets
			// the city set off and catches it part way through the walk.
			w.Player.Cash, w.Player.Respect = 3000, 30
			w.Player.Location = "bar"
			w.SetOut()
			for i := 0; i < 6; i++ {
				w.Minute += 720
				w.Arrivals()
				w.SetOut()
			}
			w.Properties["club"].Owner = "russo"
			w.Properties["laundry"].Owner = "bellandi"
			w.Minute += 720
			w.SetOut()
			w.Minute += 10 // caught part way across the city
			w.Arrivals()
			return nil
		}
		if scenario == "audience" {
			// A family across the table asking for money the player does not have,
			// so the refused choices have to say why themselves.
			w.Player.Cash, w.Player.Respect = 40, 20
			w.Player.Location = "laundry"
			w.OpenAudience("laundry")
			return nil
		}
		if scenario == "offer" {
			// A job on the table with two ways of doing it, which is the whole
			// reason the scene modal exists: a comparison.
			w.Player.Cash, w.Player.Respect = 900, 25
			w.Player.Location = "bar"
			e, err := w.ValidateProposal(core.Proposal{Location: "bar", Title: "A ledger before dawn",
				Body:    "Move the ledger out of the office before the merchant's partners come looking for it. Nobody wants a conversation about where it went.",
				Speaker: "mara", Operation: "courier", Outcome: "Delivered.", Beneficiary: "russo",
				Approaches: []core.Approach{{Method: "careful", Label: "Wait for the street to clear"}, {Method: "press", Label: "Make the delivery before closing"}}})
			if err != nil {
				return err
			}
			w.Event = e
			return nil
		}
		if scenario == "paused-job" {
			w.Player.Cash = 200
			w.Properties["laundry"].Owner = "player:1"
			w.NextPressure = w.Minute + 30
			e, err := w.ValidateProposal(core.Proposal{Location: "bar", Title: "A sealed message", Body: "Deliver a message at Saint Agnes.", Speaker: "mara", Operation: "courier", Outcome: "Delivered."})
			if err != nil {
				return err
			}
			w.Event = e
			next, err := core.Execute(w, core.Command{Revision: w.Revision, Kind: "choice", Event: e.ID, Choice: "accept"})
			if err != nil {
				return err
			}
			*w = *next
			return nil
		}
		if scenario == "dead" {
			// A protagonist who built something, with people to leave it to, so
			// the death screen has an estate to report rather than a rule.
			w.Player.Cash, w.Player.Respect, w.Player.Contacts = 9000, core.OrganizationStanding+40, 3
			w.Player.Earned = 44000
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["garage"].Owner = "player:1"
			w.OrganizationDay()
			for _, n := range w.Civilians() {
				if core.IsOfficial(n.ID) {
					continue
				}
				w.Player.Location = n.Location
				if w.SignOn(n.ID) == nil {
					break
				}
			}
			w.Minute += 9 * 1440
			w.Report("killing", "A MAN IS FOUND AT PIER 14", "Police say enquiries are continuing.")
			w.Report("police", "RAID AT BLUEBIRD LAUNDRY", "Officers searched the premises.")
			w.Die("Shot outside the Bluebird, in front of two people who will not say so.")
			return nil
		}
		if scenario == "raid-presence" {
			w.Player.Location = "bar"
			w.Witness("raid", "bar", "Police search Saint Agnes.", "RAID AT SAINT AGNES")
			w.LastResult = &core.Result{Kind: "qa-raid", From: "bar", To: "bar", Cues: w.VisualCues}
			return nil
		}
		if weaponFixture {
			w.Player.Location = "bar"
			w.Player.Cash, w.Player.Respect, w.Player.Weapon = 6000, 60, weaponTier
			victim := w.NPC("mara")
			victim.Location = "bar"
			w.RNG, w.WorldRNG = 1, 1
			if err := w.Strike(victim.ID, w.OwnHands()); err != nil {
				return err
			}
			if !w.NPC(victim.ID).Dead {
				return fmt.Errorf("weapon fixture strike did not succeed")
			}
			w.LastResult = &core.Result{Kind: "qa-weapon-strike", From: "bar", To: "bar", Cues: w.VisualCues}
			return nil
		}
		if scenario == "gunfight" || scenario == "gunfight-killing" {
			// Isolated presentation fixture for the same anonymous cue used by
			// a meeting that ends in gunfire. The paired variant explicitly kills
			// its victim through core; neither variant identifies a shooter.
			w.Player.Location = "bar"
			if scenario == "gunfight-killing" {
				victim := w.NPC("mara")
				victim.Location = "bar"
				w.Kill(victim.ID, "Shot at Saint Agnes.")
			}
			w.Witness("gunfight", "bar", "A meeting at Saint Agnes ended in gunfire.", "")
			w.LastResult = &core.Result{Kind: "qa-gunfight", From: "bar", To: "bar", Cues: w.VisualCues}
			return nil
		}
		if scenario == "killing" {
			// Somebody the player knows, dead where they stood, so the theatre
			// has a face to show and a headline to follow it.
			w.Player.Cash, w.Player.Respect = 6000, 40
			w.Properties["laundry"].Owner = "player:1"
			w.Player.Location = "bar"
			victim := w.NPC("mara")
			if victim == nil {
				return fmt.Errorf("nobody to lose")
			}
			victim.Location = "bar"
			w.VisualCues = nil
			w.Kill(victim.ID, "Shot twice at the counter, in front of everyone and nobody.")
			// Preserve the actual core-produced cue for explicit browser replay.
			w.LastResult = &core.Result{Kind: "qa-killing", From: "bar", To: "bar", Cues: w.VisualCues}
			return nil
		}
		if scenario == "herald" {
			// A campaign standing at the paper with a raid in this morning's
			// edition and the money to do something about it.
			w.Player.Cash, w.Player.Respect = 25000, 60
			w.Properties["laundry"].Owner = "player:1"
			w.Player.Location = core.HeraldPlace
			w.Player.Heat = 35
			w.Attention = 45
			w.Report("police", "RAID AT BLUEBIRD LAUNDRY", "Officers searched the premises this morning. No charges have yet been brought.")
			return nil
		}
		if scenario == "debt" {
			// Money already out with a name on it, overdue, and the man who owes
			// it standing in front of you.
			w.Player.Cash, w.Player.Respect = 40000, 80
			w.District = 2
			for _, n := range w.Civilians() {
				if core.IsOfficial(n.ID) {
					continue
				}
				w.Player.Location = n.Location
				n.Rank, n.Faction = core.RankLieutenant, "bellandi"
				if w.Lend(n.ID) != nil {
					n.Rank, n.Faction = 0, ""
					continue
				}
				// His position collapses while he is carrying it.
				n.Rank, n.Faction, n.Skill = 0, "", 5
				w.Minute = w.LoanTo(n.ID).Due
				w.LoanDay()
				break
			}
			return nil
		}
		if scenario == "arrest" {
			// A campaign with a still in the back of the laundry, people of its
			// own, and enough attention that the police are already coming.
			w.Player.Cash = 9000
			w.Player.Respect, w.Player.Contacts = core.OrganizationStanding, 3
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["garage"].Owner = "player:1"
			w.Player.Location = "laundry"
			w.OrganizationDay()
			for _, n := range w.Civilians() {
				if core.IsOfficial(n.ID) {
					continue
				}
				w.Player.Location = n.Location
				if w.SignOn(n.ID) == nil {
					break
				}
			}
			w.Player.Location = "laundry"
			if err := w.BuildStill("laundry"); err != nil {
				return err
			}
			// Enough attention for them to come, short of the point where they
			// take the premises instead of the person.
			w.Player.Heat = 60
			// The visit itself, so the fixture opens on the decision rather than
			// on a wait for the dice to produce one.
			w.Raid()
			return nil
		}
		if scenario == "doorman" {
			// A campaign that owns premises and has people of its own, which is
			// the only state in which anybody can be put on a door.
			w.Player.Cash = 5000
			w.Player.Respect, w.Player.Contacts = core.OrganizationStanding, 3
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["garage"].Owner = "player:1"
			w.Player.Location = "laundry"
			w.OrganizationDay()
			for _, n := range w.Civilians() {
				if core.IsOfficial(n.ID) {
					continue
				}
				w.Player.Location = n.Location
				if w.SignOn(n.ID) == nil {
					break
				}
			}
			w.Player.Location = "laundry"
			return nil
		}
		if scenario == "leader" {
			w.Factions[1].Goodwill = 9
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Status: "declined", Title: "A refused introduction"}}
			return nil
		}
		if scenario == "contact" {
			w.Player.Crew = []core.Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 65}}
			w.Player.Respect = 6
			w.Arrangements = []core.ArrangementMemory{{ID: core.ID(), Life: w.Life, Speaker: "mara", Operation: "courier", Status: "declined", Title: "A courier offer declined"}}
			return nil
		}
		if scenario == "attack" {
			w.Player.Location = "laundry"
			w.Properties["laundry"].Owner = "player:1"
			w.Plots = append(w.Plots, core.Plot{ID: core.ID(), Kind: "sabotage", Life: w.Life, Due: w.Minute + 15, Actor: "bellandi", Target: "laundry", Strength: 35})
			return nil
		}
		if scenario == "russo-warning" {
			w.Player.Contacts = 2
			w.Player.Cash = 300
			w.Player.Home = "apartment"
			w.Player.Location = "apartment"
			w.District = 1
			w.Factions[1].Goodwill = -40
			w.RetaliationFrom("russo")
			w.Advance(240)
			return nil
		}
		if scenario == "warning" {
			w.Player.Contacts = 2
			w.Player.Health = 60
			w.Retaliation()
			w.Advance(240)
			return nil
		}
		if scenario == "damage" {
			w.Player.Location = "laundry"
			w.Player.Respect = 6
			w.Properties["laundry"].Owner = "player:1"
			w.Properties["laundry"].Condition = 45
			w.Log("A damaged storefront", "An isolated repair and visual-state QA scenario.", "danger")
			return nil
		}
		w.Player.Heat = 14
		w.Player.Location = "bar"
		scene, err := w.ValidateProposal(core.Proposal{Title: "A Russo delivery", Body: "Take these sealed papers to our contact. With police watching your movements, the arrangement may become expensive.", Speaker: "mara", Operation: "courier", Outcome: "Delivered the sealed papers.", Beneficiary: "russo", Approaches: []core.Approach{{Method: "careful", Label: "Wait until the street clears"}, {Method: "press", Label: "Deliver before the doors close"}}})
		if err != nil {
			return err
		}
		scene.Source = "authored"
		if scenario == "voice" {
			w.Offers = append(w.Offers, core.Offer{Ready: w.Minute + 30, Event: scene})
			w.Director.Status = "ready"
		} else {
			w.Event = scene
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created isolated", scenario, "QA save:", path)
}
