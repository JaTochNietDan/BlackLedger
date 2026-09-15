package core

import "fmt"

// CrewOrder is a saved assignment. Its clock and actor belong to the simulation.
// Targets and return addresses are captured at acceptance, never inferred by UI.
type CrewOrder struct {
	Estate   string `json:"estate,omitempty"`
	Loot     int    `json:"loot"`
	ID       string `json:"id"`
	Life     int    `json:"life"`
	Actor    string `json:"actor"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Target   string `json:"target"`
	Place    string `json:"place"`
	Base     string `json:"base"`
	Stage    string `json:"stage"`
	Due      int    `json:"due"`
	Reserved int    `json:"reserved"`
	Charges  int    `json:"charges"`
	Recall   bool   `json:"recall"`
	Result   string `json:"result"`
}

func (o CrewOrder) active() bool {
	return o.Stage == "outbound" || o.Stage == "working" || o.Stage == "returning" || o.Stage == "guarding"
}
func (w *World) CrewOrderFor(id string) *CrewOrder {
	for i := range w.CrewOrders {
		if w.CrewOrders[i].Actor == id && w.CrewOrders[i].active() {
			return &w.CrewOrders[i]
		}
	}
	return nil
}

// Orders use what the player learned, not the target's hidden current address.
func (w *World) crewTargetAddress(id string) string {
	n := w.NPC(id)
	if n == nil {
		return ""
	}
	if n.Location == w.Player.Location && !w.Travelling(n) {
		return n.Location
	}
	if seen, ok := w.LastSeen(id); ok && w.Minute-seen.When <= SightingLasts {
		return seen.Where
	}
	return ""
}

func (w *World) CrewOrderReadiness(kind, actor, target string) string {
	base := w.Headquarters(w.PlayerOrganizationID())
	if !w.Player.Alive || !w.Incorporated() || base == "" {
		return "Establish a family headquarters first"
	}
	if why := w.HeadquartersReadiness(base, false); why != "" {
		return why
	}
	if w.Event != nil {
		return "Finish the conversation first"
	}
	hand, ok := w.NamedHands(actor)
	if !ok {
		return "They do not work for you"
	}
	if why := w.HandReadiness(hand); why != "" {
		return why
	}
	if w.CrewOrderFor(actor) != nil {
		return hand.Name + " is already on assignment"
	}
	for _, p := range w.Properties {
		if p.Posted == actor {
			return "Relieve them from guard duty first"
		}
	}
	switch kind {
	case "collections":
		return w.crewCollectionTarget(target, "")
	case "guard":
		return w.crewGuardTarget(target, "")
	case "repair", "remedy":
		cost, why := w.crewPropertyWork(kind, target)
		if why != "" {
			return why
		}
		if w.Player.Cash < cost {
			return "Not enough cash"
		}
		return ""
	case "rob":
		return w.robberyReadinessBy(target, hand)
	case "bomb":
		if w.Player.Charges < 1 {
			return "Acquire a charge before issuing this order"
		}
		if n := w.NPC(actor); n == nil || w.Poise(n) < PlantStanding {
			return "The operative needs more experience for demolition"
		}
		return w.crewBombTarget(target)
	case "restock":
		return w.RestockReadiness(target)
	case "assassinate":
		n := w.NPC(target)
		if n == nil || !w.Known(n) {
			return "Locate somebody you know first"
		}
		if _, ok := w.StrikeTargetAt(target, n.Location); !ok {
			return "That is not a valid target"
		}
		if w.crewTargetAddress(target) == "" {
			return "Find a recent address for the target first"
		}
		if _, ok := PlaceByID(w.crewTargetAddress(target)); !ok {
			return "There is no usable address"
		}
		return ""
	}
	return "That operation is not available yet"
}

func (w *World) StartCrewOrder(kind, actor, target string) error {
	if why := w.CrewOrderReadiness(kind, actor, target); why != "" {
		return fmt.Errorf("%s", why)
	}
	hand, _ := w.NamedHands(actor)
	at := target
	reserved := 0
	if kind == "assassinate" {
		at = w.crewTargetAddress(target)
	} else if kind == "restock" {
		reserved = w.RestockCost(target)
	} else if kind == "repair" || kind == "remedy" {
		reserved, _ = w.crewPropertyWork(kind, target)
	}
	if err := w.Pay(reserved); err != nil {
		return err
	}
	o := CrewOrder{ID: ID(), Life: w.Life, Actor: actor, Name: hand.Name, Kind: kind, Target: target, Place: at, Base: w.Headquarters(w.PlayerOrganizationID()), Stage: "outbound", Reserved: reserved}
	if kind == "bomb" {
		w.Player.Charges--
		o.Charges = 1
	}
	w.CrewOrders = append(w.CrewOrders, o)
	w.crewOrderJourney(&w.CrewOrders[len(w.CrewOrders)-1], at)
	place, _ := PlaceByID(at)
	work := kind + " at " + place.Name
	if kind == "bomb" {
		work = "plant a charge at " + place.Name
	}
	if kind == "assassinate" {
		work = "go after " + w.NPC(target).Name + " at " + place.Name
	}
	w.Log("An order from headquarters", fmt.Sprintf("%s sets out to %s. The assignment continues while you attend to other business.", hand.Name, work), "work")
	return nil
}
func (w *World) crewOrderJourney(o *CrewOrder, to string) {
	n := w.NPC(o.Actor)
	if n == nil {
		return
	}
	duration := max(1, TravelMinutes(n.Location, to))
	if n.Location == to {
		o.Due = w.Minute + 1
		n.Heading, n.Arrives, n.Sets, n.Errand = "", 0, 0, ""
		return
	}
	o.Due = w.Minute + duration
	n.Heading, n.Arrives, n.Sets, n.Errand = to, o.Due, 0, "on a headquarters assignment"
	w.noticed(n, true)
}
func (w *World) RecallCrewOrder(id string) error {
	for i := range w.CrewOrders {
		o := &w.CrewOrders[i]
		if o.ID != id || o.Life != w.Life || !o.active() {
			continue
		}
		if w.Event != nil || w.HeadquartersReadiness(w.Headquarters(w.PlayerOrganizationID()), false) != "" {
			return fmt.Errorf("Issue the recall from your usable headquarters")
		}
		if o.Stage == "returning" || o.Recall {
			return fmt.Errorf("They are already coming back")
		}
		o.Recall = true
		// Finish the street leg before returning; no teleport to its starting point.
		if o.Stage == "working" || o.Stage == "guarding" {
			w.returnCrewOrder(o, "Recalled to headquarters")
		}
		return nil
	}
	return fmt.Errorf("That assignment is no longer active")
}
func (w *World) refundCrewOrder(o *CrewOrder) {
	// Reserved job funds travel with the operative, just like loot and charges.
	// Losing the carrier cannot remotely restore their wallet to headquarters.
	carrier := w.NPC(o.Actor)
	if carrier == nil || carrier.Dead || w.Inside(carrier) {
		o.Reserved, o.Loot, o.Charges = 0, 0, 0
		return
	}
	if o.Estate != "" {
		// These funds belong to the successor organization, never a new protagonist.
		// Cash carried by a dead or captured operative is lost with them.
		n := w.NPC(o.Actor)
		if n != nil && !n.Dead && !w.Inside(n) {
			if f := w.faction(o.Estate); f != nil && n.Faction == f.ID {
				f.Cash += o.Reserved + o.Loot
			} else {
				n.Purse += o.Reserved + o.Loot
			}
		}
		if o.Charges > 0 {
			o.Result += "; unused demolition charges removed from circulation after succession"
		}
		o.Reserved, o.Loot, o.Charges = 0, 0, 0
		return
	}

	if o.Loot > 0 {
		n := w.NPC(o.Actor)
		if n != nil && !n.Dead && !w.Inside(n) {
			if o.Life == w.Life && w.Player.Alive && o.Stage == "returning" {
				w.Earn(o.Loot)
			} else {
				n.Purse += o.Loot
			}
		}
		o.Loot = 0
	}
	if o.Charges > 0 {
		n := w.NPC(o.Actor)
		if o.Life == w.Life && w.Player.Alive && n != nil && !n.Dead && !w.Inside(n) {
			w.Player.Charges += o.Charges
		}
		o.Charges = 0
	}
	if o.Reserved <= 0 {
		return
	}
	if o.Life == w.Life && w.Player.Alive {
		w.Player.Cash += o.Reserved
	} else if n := w.NPC(o.Actor); n != nil && !n.Dead {
		n.Purse += o.Reserved
	}
	o.Reserved = 0
}
func (w *World) releaseCrewGuard(o *CrewOrder) {
	if o.Kind == "guard" {
		if p := w.Properties[o.Target]; p != nil && p.Posted == o.Actor {
			p.Posted = ""
		}
	}
}
func (w *World) returnCrewOrder(o *CrewOrder, result string) {
	w.releaseCrewGuard(o)
	o.Result = result
	o.Stage = "returning"
	w.crewOrderJourney(o, o.Base)
}
func (w *World) SettleCrewOrders() {
	for i := range w.CrewOrders {
		o := &w.CrewOrders[i]
		if !o.active() {
			continue
		}
		n := w.NPC(o.Actor)
		_, hired := w.NamedHands(o.Actor)
		issuerAvailable := o.Life == w.Life && w.Player.Alive
		if o.Estate != "" {
			hired = n != nil && n.Faction == o.Estate
			issuerAvailable = w.faction(o.Estate) != nil
		}
		if n == nil || n.Dead || w.Inside(n) || !hired || !issuerAvailable {
			w.releaseCrewGuard(o)
			w.refundCrewOrder(o)
			o.Stage, o.Result = "cancelled", "Assignment ended: operative or issuing family unavailable"
			if n != nil && (n.Dead || w.Inside(n)) && n.Errand == "on a headquarters assignment" {
				n.Heading, n.Arrives, n.Sets, n.Errand = "", 0, 0, ""
			}
			continue
		}
		if o.Stage == "guarding" {
			owner := w.PlayerOrganizationID()
			if o.Estate != "" {
				owner = o.Estate
			}
			p := w.Properties[o.Target]
			if p == nil || p.Owner != owner || p.Posted != o.Actor || p.Condition <= 0 || n.Location != o.Place || w.Travelling(n) {
				w.returnCrewOrder(o, "The guard post is no longer available")
			}
			continue
		}
		if o.Due > w.Minute {
			continue
		}
		switch o.Stage {
		case "outbound":
			if n.Location != o.Place || w.Travelling(n) {
				w.returnCrewOrder(o, "The operative could not reach the address")
				continue
			}
			if o.Recall {
				w.returnCrewOrder(o, "Recalled before the work was committed")
				continue
			}
			if o.Kind == "assassinate" {
				victim, ok := w.StrikeTargetAt(o.Target, o.Place)
				if !ok || w.Travelling(victim) || w.Inside(victim) {
					w.returnCrewOrder(o, "The target was not at the address")
					continue
				}
			}
			if o.Kind == "restock" && !w.Own(o.Target) {
				w.returnCrewOrder(o, "The business changed hands before arrival")
				continue
			}
			if o.Kind == "bomb" && w.crewBombTarget(o.Target) != "" {
				w.returnCrewOrder(o, "The demolition target is no longer available")
				continue
			}
			if o.Kind == "rob" {
				hand, _ := w.NamedHands(o.Actor)
				if w.robberyReadinessBy(o.Target, hand) != "" {
					w.returnCrewOrder(o, "The robbery target is no longer available")
					continue
				}
			}
			if o.Kind == "repair" || o.Kind == "remedy" {
				if _, why := w.crewPropertyWork(o.Kind, o.Target); why != "" {
					w.returnCrewOrder(o, why)
					continue
				}
			}
			if o.Kind == "collections" && w.crewCollectionTarget(o.Target, o.ID) != "" {
				w.returnCrewOrder(o, "The collection round is no longer available")
				continue
			}
			if o.Kind == "guard" && w.crewGuardTarget(o.Target, "") != "" {
				w.returnCrewOrder(o, "The guard post is no longer available")
				continue
			}
			o.Stage = "working"
			o.Due = w.Minute + 30
			if o.Kind == "collections" {
				o.Due = w.Minute + CollectionMinutes
			}
			if o.Kind == "assassinate" {
				o.Due = w.Minute + StrikeMinutes
			}
			if o.Kind == "bomb" {
				o.Due = w.Minute + PlantMinutes
			}
			if o.Kind == "rob" {
				o.Due = w.Minute + 60
			}
			if o.Kind == "repair" || o.Kind == "remedy" {
				o.Due = w.Minute + PropertyWorkMinutes
			}
			if o.Kind == "guard" {
				o.Due = w.Minute + PostingMinutes
			}
		case "working":
			result := "The target is no longer available"
			if n.Location == o.Place && !w.Travelling(n) {
				switch o.Kind {
				case "collections":
					if w.crewCollectionTarget(o.Target, o.ID) == "" {
						o.Loot += CollectionPay
						result = fmt.Sprintf("Collection round finished; $%d carried back", CollectionPay)
					}
				case "guard":
					if w.crewGuardTarget(o.Target, "") == "" {
						w.Properties[o.Target].Posted = o.Actor
						o.Stage, o.Due, o.Result = "guarding", 0, "On guard until recalled or relieved"
						w.Log(o.Name+" takes the guard post", o.Result, "work")
						continue
					}

				case "repair", "remedy":
					cost, why := w.crewPropertyWork(o.Kind, o.Target)
					if why == "" && cost <= o.Reserved {
						o.Reserved -= cost
						if o.Kind == "repair" {
							w.applyRepair(o.Target)
							result = "Property repairs completed"
						} else {
							w.applyRemedy(o.Target, cost)
							result = "Business trouble resolved"
						}
					} else if why != "" {
						result = why
					}
				case "rob":
					hand, _ := w.NamedHands(o.Actor)
					if err := w.robWithOrder(o.Target, hand, o); err == nil {
						result = "The robbery brought back no money"
						if o.Loot > 0 {
							result = fmt.Sprintf("Robbery resolved; $%d carried back", o.Loot)
						}
					} else {
						result = err.Error()
					}

				case "bomb":
					if o.Charges == 1 && w.crewBombTarget(o.Target) == "" {
						o.Charges = 0
						hand, _ := w.NamedHands(o.Actor)
						if w.resolvePlant(o.Target, hand) {
							result = "The charge detonated at the target"
						} else {
							result = "The charge went off prematurely"
						}
					}
				case "restock":
					trade, ok := TradeOf(o.Target)
					p := w.Properties[o.Target]
					if ok && p != nil && w.Own(o.Target) && p.Supply < trade.RestockAmount && w.RestockCost(o.Target) <= o.Reserved {
						cost := w.RestockCost(o.Target)
						o.Reserved -= cost
						w.applyRestock(o.Target, cost)
						result = "Supplies delivered"
					}
				case "assassinate":
					if victim, ok := w.StrikeTargetAt(o.Target, o.Place); ok && !w.Inside(victim) && !w.Travelling(victim) {
						hand, _ := w.NamedHands(o.Actor)
						place, _ := PlaceByID(o.Place)
						family := w.faction(victim.Faction)
						if w.Random() < w.strikeChance(victim, hand) {
							w.finishThem(victim, hand, place.Name, len(w.PeopleHere(o.Place))-1, family)
							result = "The target was killed"
						} else {
							w.itWentWrong(victim, hand, place.Name, family)
							result = "The target survived the attempt"
						}
					}
				}
			}
			if n.Dead || w.Inside(n) {
				w.refundCrewOrder(o)
				o.Stage, o.Result = "cancelled", result+"; operative lost"
			} else {
				w.returnCrewOrder(o, result)
			}
		case "returning":
			w.refundCrewOrder(o)
			o.Stage = "done"
			w.Log(o.Name+" reports back", o.Result, "work")
		}
	}
}
func (w *World) PublicCrewOrders() []CrewOrder {
	out := []CrewOrder{}
	for _, o := range w.CrewOrders {
		if o.Life == w.Life {
			out = append(out, o)
		}
	}
	return out
}

type CrewOrderOffer struct {
	Charges int    `json:"charges"`
	Actor   string `json:"actor"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Target  string `json:"target"`
	Label   string `json:"label"`
	Cost    int    `json:"cost"`
	Minutes int    `json:"minutes"`
	Reason  string `json:"reason"`
}

func (w *World) CrewOrderOffers() []CrewOrderOffer {
	out := []CrewOrderOffer{}
	if !w.Incorporated() {
		return out
	}
	seen := map[string]bool{}
	ids := []string{}
	for _, c := range w.Player.Crew {
		if !seen[c.ID] {
			seen[c.ID] = true
			ids = append(ids, c.ID)
		}
	}
	for _, n := range w.OwnPeople() {
		if !seen[n.ID] {
			seen[n.ID] = true
			ids = append(ids, n.ID)
		}
	}
	base := w.Headquarters(w.PlayerOrganizationID())
	for _, id := range ids {
		n := w.NPC(id)
		if n == nil || n.Dead {
			continue
		}
		offer := func(kind, target, at, label string, cost, work int) {
			outbound, back := max(1, TravelMinutes(n.Location, at)), max(1, TravelMinutes(at, base))
			if n.Location == at {
				outbound = 1
			}
			if at == base {
				back = 1
			}
			duration := outbound + work + back
			if kind == "guard" {
				duration = outbound + work
			}
			if at == "" {
				duration = 0
			}
			charges := 0
			if kind == "bomb" {
				charges = 1
			}
			out = append(out, CrewOrderOffer{Charges: charges, Actor: id, Name: n.Name, Kind: kind, Target: target, Label: label, Cost: cost, Minutes: duration, Reason: w.CrewOrderReadiness(kind, id, target)})
		}
		for _, l := range Locations {
			if p := w.Properties[l.ID]; p != nil && w.Own(l.ID) {
				if p.Income > 0 {
					offer("collections", l.ID, l.ID, fmt.Sprintf("Collections at %s · $%d on return", l.Name, CollectionPay), 0, CollectionMinutes)
					offer("guard", l.ID, l.ID, "Guard "+l.Name, 0, PostingMinutes)
				}
				offer("repair", l.ID, l.ID, "Repair "+l.Name, RepairCost, PropertyWorkMinutes)
				if trade, ok := TradeOf(l.ID); ok {
					offer("remedy", l.ID, l.ID, trade.Remedy+" at "+l.Name, trade.RemedyCost, PropertyWorkMinutes)
				}
			}
			if p := w.Properties[l.ID]; p != nil && p.Income > 0 && !w.Own(l.ID) {
				offer("bomb", l.ID, l.ID, "Bomb "+l.Name, 0, PlantMinutes)
				offer("rob", l.ID, l.ID, "Rob "+l.Name, 0, 60)
			}
			if _, ok := TradeOf(l.ID); ok && w.Own(l.ID) {
				offer("restock", l.ID, l.ID, "Restock "+l.Name, w.RestockCost(l.ID), 30)
			}
		}
		for i := range w.NPCs {
			victim := &w.NPCs[i]
			if !victim.Dead && w.Known(victim) {
				if _, ok := w.StrikeTargetAt(victim.ID, victim.Location); ok {
					offer("assassinate", victim.ID, w.crewTargetAddress(victim.ID), "Assassinate "+victim.Name, 0, StrikeMinutes)
				}
			}
		}
	}
	return out
}

func (w *World) crewBombTarget(id string) string {
	p := w.Properties[id]
	if _, ok := PlaceByID(id); !ok || p == nil || p.Income <= 0 {
		return "There is no business to target"
	}
	if w.Own(id) {
		return "Your family owns this business"
	}
	if p.Condition <= 0 {
		return "The building is already destroyed"
	}
	return ""
}

// inheritCrewOrders winds down the former leader's commitments. The successor
// retains custody of family money while operatives complete their return legs.
func (w *World) inheritCrewOrders(estate string) {
	for i := range w.CrewOrders {
		o := &w.CrewOrders[i]
		if o.Life != w.Life || !o.active() || o.Estate != "" {
			continue
		}
		n := w.NPC(o.Actor)
		// Personally hired associates do not become sworn family members on death.
		if n == nil || n.Dead || n.Faction != estate {
			continue
		}
		o.Estate = estate
		if base := w.Headquarters(estate); base != "" {
			o.Base = base
		}
		if o.Stage == "returning" || o.Stage == "guarding" {
			continue
		}
		o.Recall = true
		o.Result = "Recalled after the change of leadership"
		if o.Stage == "working" {
			w.returnCrewOrder(o, o.Result)
		}
	}
}

func (w *World) crewGuardTarget(id, actor string) string {
	p := w.Properties[id]
	if p == nil || !w.Own(id) || p.Income <= 0 {
		return "Choose a business your family owns"
	}
	if p.Condition <= 0 {
		return "The building is destroyed"
	}
	if n := w.PostedAt(id); n != nil && n.ID != actor {
		return n.Name + " is already on the door"
	}
	return ""
}

// Collections retain the established two-hour round's modeled fee. This is
// supplemental work, not a second withdrawal of automatic business income.
func (w *World) crewCollectionTarget(id, order string) string {
	p := w.Properties[id]
	if p == nil || !w.Own(id) || p.Income <= 0 {
		return "Choose a business your family owns"
	}
	if p.Condition <= 0 || p.Trouble {
		return "Restore the business before sending a collection round"
	}
	if len(w.Tasks) > 0 {
		return "An existing collection round must finish first"
	}
	for _, o := range w.CrewOrders {
		if o.ID != order && o.Kind == "collections" && o.Target == id && o.active() {
			return "Someone is already handling this collection round"
		}
	}
	return ""
}
