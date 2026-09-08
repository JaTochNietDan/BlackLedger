package core

import (
	"fmt"
	"strings"
)

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
	// What the player had before they decided, so the result can say what the
	// decision actually cost rather than leaving them to diff two screens.
	wasCash, wasRespect, wasHeat, wasHealth := p.Cash, p.Respect, p.Heat, p.Health
	chosen := ""
	previousRecords := make(map[string]bool, len(w.History))
	for _, record := range w.History {
		previousRecords[record.ID] = true
	}
	if c.Kind == "new_life" {
		if p.Alive {
			return fmt.Errorf("this life is still in progress")
		}
		// Whatever the city had started calling them dies with them.
		w.Dissolve(w.PlayerOrganizationID())
		w.Life++
		w.Player = newPerson(w.Life)
		w.Event = nil
		w.Offers = []Offer{}
		w.Plots = []Plot{}
		// Arrangements the dead protagonist paid for die with them. Without
		// this they linger in the save forever, filtered out but never removed.
		w.Contracts = nil
		w.Commissions = nil
		w.Pacts = nil
		w.Player.Serves, w.Player.Service = "", 0
		w.BusinessTruces = nil
		w.SuspendedJob = nil
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
		case "sitdown":
			if err := w.ResolveSitdown(e, c.Choice); err != nil {
				return err
			}
		case "arrest":
			if err := w.ResolveArrest(e, c.Choice); err != nil {
				return err
			}
		case "contract":
			if c.Choice == "leave" {
				w.Log("No name given", "You let the conversation end without saying anything worth repeating.", "personal")
				break
			}
			if err := w.openContractTerms(strings.TrimPrefix(c.Choice, "mark:")); err != nil {
				return err
			}
			// The terms are the next decision, not a committed outcome.
			return nil
		case "contract_terms":
			if c.Choice == "leave" {
				w.Log("Nothing was agreed", "You leave the name where it was.", "personal")
				break
			}
			if err := w.Commission(e.Target, strings.TrimPrefix(c.Choice, "hire:")); err != nil {
				return err
			}
		case "business_pressure":
			if err := w.ResolvePressure(e, c.Choice); err != nil {
				return err
			}
			w.OfferResume()
		case "resume_job":
			saved := w.SuspendedJob
			if saved == nil || saved.Scene == nil {
				return fmt.Errorf("no suspended arrangement")
			}
			w.SuspendedJob = nil
			if c.Choice == "resume" {
				w.RunArrangement(saved.Scene, saved.Remaining)
			} else {
				w.RememberArrangement(saved.Scene, "abandoned")
				w.Log("The arrangement abandoned", "You abandon the remaining work. No reward was paid.", "story")
			}
		case "police_stop":
			// A search finds whatever is being carried, whichever way the stop
			// is settled. This is what makes moving goods quickly matter.
			w.Seize("A detective searched you during the stop.")
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
				w.RunArrangement(e, e.Effect.Minutes)
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
		chosen = a.Label
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
				// What the city was doing while the player was crossing it.
				w.PassThrough(oldLoc, target)
			} else {
				p.Location = oldLoc
				if p.Alive {
					from, _ := PlaceByID(oldLoc)
					to, _ := PlaceByID(target)
					w.Log("Journey interrupted", fmt.Sprintf("The journey to %s was interrupted after %d minutes. You remain based at %s; choose your next destination after resolving the situation.", to.Name, w.Minute-oldTime, from.Name), "travel")
				}
			}
		} else if id, ok := strings.CutPrefix(c.Kind, "serve:"); ok {
			if err := w.Serve(id); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if with, ok := strings.CutPrefix(c.Kind, "pact:"); ok {
			if err := w.MakePact(with); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if with, ok := strings.CutPrefix(c.Kind, "break:"); ok {
			if !w.Allied(with) {
				return fmt.Errorf("you have no understanding with them")
			}
			w.BreakPact(with, "You ended it. Nobody forgets which side did that.")
			if f := w.faction(with); f != nil {
				f.Goodwill = max(-100, f.Goodwill-20)
			}
			w.Advance(a.Minutes)
		} else if about, ok := strings.CutPrefix(c.Kind, "enquire:"); ok {
			if err := w.AskAround(about); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "sign:"); ok {
			if err := w.SignOn(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "share:"); ok {
			if err := w.PayShare(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if org, ok := strings.CutPrefix(c.Kind, "smear:"); ok {
			if err := w.Smear(org); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "lend:"); ok {
			if err := w.Lend(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "lean:"); ok {
			// Resolved before the clock moves, like sabotage: a collection that
			// goes wrong must not also collect the hours it never survived.
			if err := w.Lean(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "extend:"); ok {
			if err := w.Extend(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "forgive:"); ok {
			if err := w.Forgive(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "bail:"); ok {
			if err := w.Bail(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
		} else if person, ok := strings.CutPrefix(c.Kind, "dismiss:"); ok {
			if err := w.LetGo(person); err != nil {
				return err
			}
			w.Advance(a.Minutes)
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
				holder := w.PropertyHolder(target)
				if f := w.FactionByID(holder); f != nil {
					f.Goodwill = max(-100, f.Goodwill-35)
					w.RetaliationFrom(f.ID)
				} else {
					w.Retaliation()
				}
				w.Log("A demand nobody forgets", "The manager refuses. A Bellandi man watches you leave. You have challenged a powerful family on its own ground.", "politics")
			case "play:small", "play:high":
				if err := w.Deal(target, strings.TrimPrefix(c.Kind, "play:")); err != nil {
					return err
				}
			case "hit":
				if err := w.DrawCard(); err != nil {
					return err
				}
			case "stand":
				if err := w.Stand(); err != nil {
					return err
				}
			case "arms:weapon", "arms:armour":
				if err := w.BuyArms(strings.TrimPrefix(c.Kind, "arms:")); err != nil {
					return err
				}
			case "trip:rockridge", "trip:kingsport", "trip:halloway":
				if err := w.Trip(strings.TrimPrefix(c.Kind, "trip:")); err != nil {
					return err
				}
			case "commission":
				if err := w.TakeCommission(target); err != nil {
					return err
				}
			case "fit:door", "fit:telephone", "fit:safe", "fit:cellar":
				if err := w.Fit(target, strings.TrimPrefix(c.Kind, "fit:")); err != nil {
					return err
				}
			case "sign", "share", "dismiss":
				// Never reached: these carry a person's id after the colon and
				// are handled by prefix below. Listed so the switch reads as
				// the full set of what the player can do.
				return fmt.Errorf("who?")
			case "spike":
				if err := w.Spike(); err != nil {
					return err
				}
			case "puff":
				if err := w.Puff(); err != nil {
					return err
				}
			case "retain:commissioner", "retain:mayor", "retain:editor":
				if err := w.Retain(strings.TrimPrefix(c.Kind, "retain:")); err != nil {
					return err
				}
			case "release:commissioner", "release:mayor", "release:editor":
				if err := w.EndRetainer(strings.TrimPrefix(c.Kind, "release:")); err != nil {
					return err
				}
			case "charge":
				if err := w.BuyCharge(); err != nil {
					return err
				}
			case "plant":
				// Resolved before the clock moves, so a charge that kills the
				// player cannot also collect the hours it never survived.
				if err := w.Plant(target); err != nil {
					return err
				}
			case "car":
				if err := w.BuyVehicle(); err != nil {
					return err
				}
			case "service":
				if err := w.Service(target); err != nil {
					return err
				}
			case "bankroll":
				if err := w.Bankroll(target); err != nil {
					return err
				}
			case "draw":
				if err := w.Draw(target); err != nil {
					return err
				}
			case "dress":
				if err := w.BuyAttire(); err != nil {
					return err
				}
			case "press":
				if err := w.Press(target); err != nil {
					return err
				}
			case "armoury":
				if err := w.BuildArmoury(target); err != nil {
					return err
				}
			case "stock_arms":
				if err := w.StockArmoury(); err != nil {
					return err
				}
			case "still":
				if err := w.BuildStill(target); err != nil {
					return err
				}
			case "dismantle":
				if err := w.Dismantle(target); err != nil {
					return err
				}
			case "hire":
				if err := w.Hire(target); err != nil {
					return err
				}
			case "layoff":
				if err := w.LayOff(target); err != nil {
					return err
				}
			case "takeover":
				// Resolved before the clock moves: a man who does not survive
				// it does not collect the evening.
				if err := w.TakeOver(); err != nil {
					return err
				}
			case "leave_service":
				if err := w.LeaveService(); err != nil {
					return err
				}
			case "order":
				if err := w.TakeOrder(target); err != nil {
					return err
				}
			case "restock":
				if err := w.Restock(target); err != nil {
					return err
				}
			case "remedy":
				if err := w.Remedy(target); err != nil {
					return err
				}
			case "deposit":
				if err := w.Deposit(); err != nil {
					return err
				}
			case "offshore_access":
				if err := w.EstablishAccess(); err != nil {
					return err
				}
			case "withdraw":
				if err := w.Withdraw(); err != nil {
					return err
				}
			case "bribe":
				if err := w.Bribe(); err != nil {
					return err
				}
			case "launder":
				if err := w.Launder(target); err != nil {
					return err
				}
			case "mug":
				if err := w.Mug(target, w.OwnHands()); err != nil {
					return err
				}
			case "mug:crew":
				hand, ok := w.CrewHands()
				if !ok {
					return fmt.Errorf("you have nobody to send")
				}
				if err := w.Mug(target, hand); err != nil {
					return err
				}
			case "rob:crew":
				hand, ok := w.CrewHands()
				if !ok {
					return fmt.Errorf("you have nobody to send")
				}
				if err := w.RobBy(target, hand); err != nil {
					return err
				}
			case "sabotage:crew":
				hand, ok := w.CrewHands()
				if !ok {
					return fmt.Errorf("you have nobody to send")
				}
				// Resolved before the clock moves, like the version the player
				// carries out themselves.
				if err := w.SabotageBy(target, hand); err != nil {
					return err
				}
			case "rob":
				if err := w.Rob(target); err != nil {
					return err
				}
			case "contract":
				w.OpenContract()
			case "buy:moonshine", "buy:cigarettes", "buy:arms":
				if err := w.Buy(strings.TrimPrefix(c.Kind, "buy:")); err != nil {
					return err
				}
			case "sell:moonshine", "sell:cigarettes", "sell:arms":
				if err := w.Sell(strings.TrimPrefix(c.Kind, "sell:")); err != nil {
					return err
				}
			case "operate:clean", "operate:standard", "operate:hard":
				if err := w.SetMode(target, strings.TrimPrefix(c.Kind, "operate:")); err != nil {
					return err
				}
			case "incite":
				if err := w.Incite(target); err != nil {
					return err
				}
			case "move":
				// Resolved before the clock moves, like sabotage, so a move
				// that kills the player cannot also collect the hours it never
				// survived.
				if err := w.MoveOn(target); err != nil {
					return err
				}
			case "sabotage":
				// Resolved before the clock moves, so a fatal attempt cannot also
				// collect the time and income of the hours it never survived.
				if err := w.Sabotage(target); err != nil {
					return err
				}
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
					if fixer := w.Holder("fixer"); fixer != nil {
						fixer.Trust += 5
					}
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
				case "sitdown":
					// Opened after the clock moves, like an audience, so the
					// evening actually costs the evening.
					if err := w.CallSitdown(); err != nil {
						return err
					}
				case "audience":
					w.OpenAudience(target)
				case "acquire":
					if w.NextPressure == 0 {
						w.NextPressure = w.Minute + 180
					}
					priorOwner := w.Properties[target].Owner
					w.Properties[target].Owner = fmt.Sprintf("player:%d", w.Life)
					// A business is bought as a going concern: the people
					// working it and what it runs on come with it. Keeping them
					// is the player's problem from here.
					if trade, running := TradeOf(target); running {
						prop := w.Properties[target]
						prop.Staff = max(prop.Staff, trade.Hands)
						prop.Supply = max(prop.Supply, trade.RestockAmount)
					}
					p.Respect += 4
					l, _ := PlaceByID(target)
					w.Log("A foothold in the city", l.Name+" now produces income for you. Earnings accrue as game time passes.", "business")
					// Taking premises in a family's district is noticed by them.
					if previous := priorOwner; previous != "" && previous != "independent" {
						if f := w.FactionByID(previous); f != nil {
							f.Goodwill = max(-100, f.Goodwill-10)
						}
					}
				case "inspect":
					l, _ := PlaceByID(target)
					w.Log("The books are open", fmt.Sprintf("%s: %d%% condition, earning $%d/hour of a possible $%d/hour. Repairs cost $50 and restore up to 40 condition.", l.Name, w.Properties[target].Condition, w.Properties[target].Income*w.Properties[target].Condition/100, w.Properties[target].Income), "business")
				case "sit_out":
					if err := w.SitOut(); err != nil {
						return err
					}
				case "lawyer":
					if err := w.Lawyer(); err != nil {
						return err
					}
				case "talk":
					if err := w.Talk(); err != nil {
						return err
					}
				case "post":
					if err := w.Post(target); err != nil {
						return err
					}
				case "unpost":
					if err := w.Unpost(target); err != nil {
						return err
					}
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
		// Completed travel is an encounter boundary too. Otherwise a prepared
		// contact stays silent until the player performs an unrelated local action.
		// OfferIfReady preserves any urgent incident raised during the journey.
		if c.Kind != "provoke" && c.Kind != "audience" {
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
	w.LastResult = &Result{
		Action: chosen, Kind: c.Kind,
		From: oldLoc, To: w.Player.Location, Elapsed: w.Minute - oldTime,
		Records: newRecords, Cues: w.VisualCues,
		Cash:    w.Player.Cash - wasCash,
		Respect: w.Player.Respect - wasRespect,
		Heat:    w.Player.Heat - wasHeat,
		Health:  w.Player.Health - wasHealth,
	}
	return nil
}
