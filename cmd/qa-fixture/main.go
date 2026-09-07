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
		log.Fatal("usage: go run ./cmd/qa-fixture <new-qa.sqlite3> [police|damage]")
	}
	scenario := "police"
	if len(os.Args) == 3 {
		scenario = os.Args[2]
	}
	if scenario != "police" && scenario != "damage" {
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
		w.Event = scene
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created isolated", scenario, "QA save:", path)
}
