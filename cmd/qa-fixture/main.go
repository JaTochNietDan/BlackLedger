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
		log.Fatal("usage: go run ./cmd/qa-fixture <new-qa.sqlite3> [police|damage|warning|russo-warning|attack|voice|contact|paused-job|leader|doorman|arrest|debt|herald|killing|dead|offer|audience|street|room|gone|post|round]")
	}
	scenario := "police"
	if len(os.Args) == 3 {
		scenario = os.Args[2]
	}
	if scenario != "police" && scenario != "damage" && scenario != "warning" && scenario != "russo-warning" && scenario != "attack" && scenario != "voice" && scenario != "contact" && scenario != "paused-job" && scenario != "leader" && scenario != "doorman" && scenario != "arrest" && scenario != "debt" && scenario != "herald" && scenario != "killing" && scenario != "dead" && scenario != "offer" && scenario != "audience" && scenario != "street" && scenario != "room" && scenario != "gone" && scenario != "post" && scenario != "round" {
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
