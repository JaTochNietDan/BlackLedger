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
	w.VisualCues = nil
	oldTime := w.Minute
	oldLoc := p.Location
	previousRecords := make(map[string]bool, len(w.History))
	for _, record := range w.History {
		previousRecords[record.ID] = true
	}
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
		case "warning":
			w.Log("Time to prepare", "You heed Mara's warning. The threat remains; your next action is yours to choose.", "intel")
		case "attack":
			if c.Choice == "bargain" {
				p.Respect = max(0, p.Respect-5)
				w.Log("A costly reprieve", "The men accept your money and withdraw. Their withdrawal does not repair your relationship with the family that sent them.", "danger")
			} else {
				base := .12
				if c.Choice == "escape" {
					base = .5
				}
				chance := base + float64(w.Guard())*.16 + float64(p.Contacts)*.04 - float64(100-p.Health)*.0025
				if chance < .05 {
					chance = .05
				}
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
			caption := "The attackers withdraw. You survived the encounter."
			if !p.Alive {
				caption = "The attack ended your life."
			} else if c.Choice == "bargain" {
				caption = "The attackers accepted payment and withdrew."
			}
			w.VisualCues = append(w.VisualCues, VisualCue{ID(), "attack", p.Home, caption})
		case "audience":
			if err := w.ResolveAudience(e, c.Choice); err != nil {
				return err
			}
		case "business_pressure":
			if err := w.ResolvePressure(e, c.Choice); err != nil {
				return err
			}
		case "police_stop":
			if c.Choice == "pay" {
				p.Heat = max(0, p.Heat-10)
				w.CompleteArrangement(e)
			} else {
				p.Heat = max(0, p.Heat-6)
				w.RememberArrangement(e, "abandoned")
				w.Log("The arrangement abandoned", "You surrender the package or paperwork and leave without completing the job. No reward was paid. Police attention eases.", "story")
			}
		case "proposal":
			w.RememberArrangement(e, "offered")
			alternative, hasAlternative := e.Alternatives[c.Choice]
			if hasAlternative {
				e.Effect = alternative
				if c.Choice == "approach:careful" {
					e.Outcome += " You took extra time to keep the work discreet."
				} else {
					e.Outcome += " You pushed the schedule and drew more attention."
				}
			}
			if c.Choice == "accept" || hasAlternative {
				w.RememberArrangement(e, "in_progress")
				w.Advance(e.Effect.Minutes)
				if p.Alive && w.Event == nil {
					if p.Heat+e.Effect.Heat >= 15 {
						w.PoliceStop(e)
					} else {
						w.CompleteArrangement(e)
					}
				} else {
					w.RememberArrangement(e, "interrupted")
					w.Log("An interrupted arrangement", "The operation could not be completed. No reward was paid.", "story")
				}
			} else {
				w.RememberArrangement(e, "declined")
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
			if p.Alive && w.Minute-oldTime >= a.Minutes {
				p.Location = target
				l, _ := PlaceByID(target)
				w.Log("Arrived at "+l.Name, fmt.Sprintf("The journey took %d minutes.", w.Minute-oldTime), "travel")
			} else {
				p.Location = oldLoc
				if p.Alive {
					from, _ := PlaceByID(oldLoc)
					to, _ := PlaceByID(target)
					w.Log("Journey interrupted", fmt.Sprintf("The journey to %s was interrupted after %d minutes. You remain based at %s; choose your next destination after resolving the situation.", to.Name, w.Minute-oldTime, from.Name), "travel")
				}
			}
		} else {
			// Hiring/delegating commits arrangements immediately; work rewards require reaching completion.
			switch c.Kind {
			case "security":
				p.Security++
				w.Log("Someone at the door", "Another security detail is assigned to your residence. It adds $10 a day to your expenses.", "personal")
			case "crew_bonus":
				before := p.Crew[0].Loyalty
				p.Crew[0].Loyalty = min(100, before+25)
				w.Log("A share for Leo", fmt.Sprintf("You paid a $40 bonus. Loyalty rose from %d to %d.", before, p.Crew[0].Loyalty), "personal")
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
					w.Investigate()
				case "lie_low":
					p.Heat = max(0, p.Heat-10)
					w.Log("Out of the spotlight", "You avoid attention for a while.", "personal")
				case "audience":
					w.OpenAudience(target)
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
					w.Log("The books are open", fmt.Sprintf("%s: %d%% condition, earning $%d/hour of a possible $%d/hour. Repairs cost $50 and restore up to 40 condition.", l.Name, w.Properties[target].Condition, w.Properties[target].Income*w.Properties[target].Condition/100, w.Properties[target].Income), "business")
				case "repair":
					restored := min(40, 100-w.Properties[target].Condition)
					w.Properties[target].Condition += restored
					w.Log("Repairs arranged", fmt.Sprintf("The property is restored by %d condition, to %d%%.", restored, w.Properties[target].Condition), "business")
				case "move_home":
					p.BestHome = max(p.BestHome, HomeRank(p.Home))
					if HomeRank(target) > p.BestHome {
						p.Respect += 3
						p.BestHome = HomeRank(target)
					}
					if target == "estate" {
						w.Properties[target].Owner = fmt.Sprintf("player:%d", w.Life)
					}
					p.Home = target
					p.Security = 0
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
	// History is capped. Its old length is not a stable cursor once new entries
	// evict old ones; identify this command's records by their persistent IDs.
	newRecords := []Record{}
	for _, record := range w.History {
		if !previousRecords[record.ID] {
			newRecords = append(newRecords, record)
		}
	}
	w.LastResult = &Result{From: oldLoc, To: w.Player.Location, Elapsed: w.Minute - oldTime, Records: newRecords, Cues: w.VisualCues}
	return nil
}
