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
		log.Fatal("usage: go run ./cmd/qa-fixture <new-qa.sqlite3> [police|damage|warning|russo-warning|attack|voice|contact|paused-job]")
	}
	scenario := "police"
	if len(os.Args) == 3 {
		scenario = os.Args[2]
	}
	if scenario != "police" && scenario != "damage" && scenario != "warning" && scenario != "russo-warning" && scenario != "attack" && scenario != "voice" && scenario != "contact" && scenario != "paused-job" {
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
