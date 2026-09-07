package core

import "fmt"

// Execute validates and applies a command to a copy: rejected commands never partially mutate state.
func Execute(original *World, c Command) (*World, error) {
	w := original.Clone()
	if c.Revision != w.Revision {
		return nil, fmt.Errorf("the city has changed; refresh before deciding")
	}
	if err := w.apply(c); err != nil {
		return nil, err
	}
	return w, nil
}
func (w *World) apply(c Command) error {
	p := &w.Player
	oldTime := w.Minute
	oldLoc := p.Location
	start := len(w.History)
	if c.Kind == "new_life" {
		if p.Alive {
			return fmt.Errorf("this life is still in progress")
		}
		w.Life++
		w.Player = newPerson(w.Life)
		w.Event = nil
		w.Offers = []Offer{}
		w.Plots = []Plot{}
		w.NextPressure = 0
		for i := range w.NPCs {
			w.NPCs[i].Trust = 0
		}
		w.Director = Director{"available", "A new life begins. Ready to prepare encounters.", -9999}
		w.Minute += 480
		for i := range w.Factions {
			w.Factions[i].Goodwill = 0
		}
		w.Log("A stranger arrives", w.Player.Name+" rents a room at the Mariner. The previous life left its mark on the city.", "personal")
	} else if !p.Alive {
		return fmt.Errorf("this life has ended")
	} else if c.Kind == "choice" {
		e := w.Event
		if e == nil || e.ID != c.Event {
			return fmt.Errorf("that conversation is no longer current")
		}
		var choice *Choice
		for i := range e.Choices {
			if e.Choices[i].ID == c.Choice {
				choice = &e.Choices[i]
			}
		}
		if choice == nil {
			return fmt.Errorf("unknown decision")
		}
		if err := w.Pay(choice.Cost); err != nil {
			return err
		}
		w.Event = nil
		switch e.Kind {
		case "attack":
			if c.Choice == "bargain" {
				p.Respect = max(0, p.Respect-5)
				w.Log("A costly reprieve", "The men accept your money and withdraw. This does not make Bellandi your friend.", "danger")
			} else {
				base := .12
				if c.Choice == "escape" {
					base = .5
				}
				chance := base + float64(w.Guard())*.16 + float64(p.Contacts)*.04
				if chance > .9 {
					chance = .9
				}
				if w.Random() > chance {
					w.Die("You did not survive the attack at your residence.")
				} else {
					lost := min(p.Security, 1+int(w.Random()*2))
					p.Security -= lost
					damage := 40
					if c.Choice == "escape" {
						damage = 20
						p.Location = "bar"
					}
					p.Health = max(1, p.Health-damage)
					w.Log("Alive, at a price", fmt.Sprintf("You survive wounded. %d security details are lost. The attackers withdraw.", lost), "danger")
				}
			}
		case "audience":
			if c.Choice == "tribute" {
				w.Factions[0].Goodwill += 8
				remaining := []Plot{}
				for _, t := range w.Plots {
					if t.Actor != "bellandi" {
						remaining = append(remaining, t)
					}
				}
				w.Plots = remaining
				w.Log("A temporary understanding", "Bellandi accepts the tribute and calls off his current operation against you.", "politics")
			} else {
				w.Log("You leave the Monarch", "Nothing was agreed. Existing threats remain.", "personal")
			}
		case "business_pressure":
			if err := w.ResolvePressure(e, c.Choice); err != nil {
				return err
			}
		case "proposal":
			if c.Choice == "accept" {
				w.Advance(e.Effect.Minutes)
				if p.Alive && w.Event == nil {
					w.Earn(e.Effect.Reward)
					p.Respect += e.Effect.Respect
					p.Heat = min(100, p.Heat+e.Effect.Heat)
					w.NPC(e.Speaker).Trust += 3
					w.Log(e.Title, e.Outcome+fmt.Sprintf(" ($%d, respect +%d)", e.Effect.Reward, e.Effect.Respect), "story")
				} else {
					w.Log("An interrupted arrangement", "The operation could not be completed. No reward was paid.", "story")
				}
			} else {
				w.Log("An offer declined", "You decline "+w.NPC(e.Speaker).Name+"'s proposal. No payment changes hands.", "story")
			}
		}
	} else {
		if w.Event != nil {
			return fmt.Errorf("resolve the current situation first")
		}
		target := c.Target
		if target == "" {
			target = p.Location
		}
		var a *Action
		actions := w.Actions(target)
		for i := range actions {
			if actions[i].ID == c.Kind {
				a = &actions[i]
			}
		}
		if a == nil {
			return fmt.Errorf("that action is not available here")
		}
		if a.Disabled {
			return fmt.Errorf("%s", a.Reason)
		}
		if err := w.Pay(a.Cost); err != nil {
			return err
		}
		if c.Kind == "travel" {
			p.Location = "transit"
			w.Advance(a.Minutes)
			if p.Alive && w.Event == nil {
				p.Location = target
				l, _ := PlaceByID(target)
				w.Log("Arrived at "+l.Name, fmt.Sprintf("The journey took %d minutes.", w.Minute-oldTime), "travel")
			} else {
				p.Location = oldLoc
			}
		} else {
			// Hiring/delegating commits arrangements immediately; work rewards require reaching completion.
			switch c.Kind {
			case "security":
				p.Security++
				w.Log("Someone at the door", "Another security detail is assigned to your residence. It adds $10 a day to your expenses.", "personal")
			case "delegate":
				w.Tasks = append(w.Tasks, Task{ID(), "Leo · collections", w.Minute + 120})
				w.Log("Leo heads out", "Collections should be completed in two hours.", "work")
			case "provoke":
				p.Respect++
				w.Factions[0].Goodwill -= 35
				w.Retaliation()
				w.Log("A demand nobody forgets", "The manager refuses. A Bellandi man watches you leave. You have challenged a powerful family on its own ground.", "politics")
			}
			w.Advance(a.Minutes)
			if p.Alive && w.Event == nil {
				switch c.Kind {
				case "courier":
					w.Earn(45)
					p.Respect += 2
					p.JobCount++
					w.Log("Envelope delivered", "Mara pays $45. A small favor, completed without questions.", "work")
				case "dockwork":
					w.Earn(75)
					p.Respect++
					if w.Random() < .15 {
						p.Health = max(1, p.Health-10)
					}
					w.Log("Cargo shifted", "$75 for a long shift on Pier 14.", "work")
				case "contact":
					p.Contacts = min(5, p.Contacts+1)
					w.NPC("mara").Trust += 5
					w.Log("A useful conversation", "Mara will keep an ear open. Your information network improves.", "personal")
				case "recruit":
					p.Crew = append(p.Crew, Crew{"leo", "Leo Carver", 65})
					p.Respect += 2
					w.Log("Your first associate", "Leo Carver joins you. He expects $12 a day and a boss who keeps their word.", "personal")
				case "investigate":
					found := false
					for i := range w.Plots {
						if w.Plots[i].Life == w.Life {
							w.Plots[i].Known = true
							found = true
						}
					}
					text := "Your sources have no evidence of an active operation against you. This is not a guarantee of safety."
					if found {
						text = "Bellandi has commissioned an operation against you. Avoid home, negotiate, or arrange protection."
					}
					w.Log("Word on the street", text, "intel")
				case "lie_low":
					p.Heat = max(0, p.Heat-10)
					w.Log("Out of the spotlight", "You avoid attention for a while.", "personal")
				case "audience":
					w.Event = &Scene{ID: ID(), Title: "A seat across from Bellandi", Body: "“People mistake an open door for an invitation. Tell me you understand the difference.”", Speaker: "vittorio", Kind: "audience", Source: "authored", Minute: w.Minute, Choices: []Choice{{ID: "tribute", Label: "Offer $150 in tribute", Cost: 150, Detail: "Improves relations and cancels his current operation against you."}, {ID: "leave", Label: "Leave without an agreement", Detail: "No payment. Existing threats remain."}}}
				case "acquire":
					if w.NextPressure == 0 {
						w.NextPressure = w.Minute + 180
					}
					w.Properties[target].Owner = fmt.Sprintf("player:%d", w.Life)
					p.Respect += 4
					l, _ := PlaceByID(target)
					w.Log("A foothold in the city", l.Name+" now produces income for you. Earnings accrue as game time passes.", "business")
					if target == "garage" {
						w.Factions[1].Goodwill -= 10
					}
				case "inspect":
					l, _ := PlaceByID(target)
					w.Log("The books are open", fmt.Sprintf("%s: %d%% condition. Income depends on its condition.", l.Name, w.Properties[target].Condition), "business")
				case "repair":
					w.Properties[target].Condition = min(100, w.Properties[target].Condition+40)
					w.Log("Repairs arranged", "The property is restored by 40 condition.", "business")
				case "move_home":
					p.Home = target
					p.Security = 0
					p.Respect += 3
					l, _ := PlaceByID(target)
					w.Log("A different view", l.Name+" is now your residence. Hired security must be arranged here.", "personal")
				case "rest":
					p.Health = min(100, p.Health+25)
					p.Heat = max(0, p.Heat-5)
					w.Log("Time to recover", "Rest restored your health. The city did not stop while you slept.", "personal")
				case "expand":
					w.District = min(2, w.District+1)
					w.Log("The city opens up", "Your contacts introduce you to another district. More properties are accessible.", "city")
				}
			} else if a.Cost > 0 && c.Kind != "security" {
				w.Player.Cash += a.Cost
				w.Log("A commitment interrupted", "Unspent funds were returned. The arrangement was not completed.", "personal")
			}
		}
		if c.Kind != "travel" && c.Kind != "provoke" && c.Kind != "audience" {
			w.OfferIfReady()
		}
	}
	w.Revision++
	start = min(start, len(w.History))
	w.LastResult = &Result{oldLoc, w.Player.Location, w.Minute - oldTime, append([]Record{}, w.History[start:]...)}
	return nil
}
