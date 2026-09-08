package core

import "fmt"

// The police in this city took money, stock and premises, and never once took a
// person. Their own newspaper printed ARRESTS EXPECTED after every raid and no
// arrest ever followed. Attention could end a career and could not end a week.
//
// Being taken in is the failure this game was missing. It is not death: you come
// back to everything you had, minus whatever ran itself into the ground while
// you were not there to watch it. Wages come out, supplies run down, rivals move
// on ground you cannot defend, and the city hands you the bill on the way out.
// That is the whole point of it — the world already runs without the player, and
// this is the mechanic that makes them feel it.

const (
	// ChargeMinimum is the weight of evidence below which there is a fine and
	// nothing more. A search that turns up a till is not a case.
	ChargeMinimum = 3
	// SentenceFloor and SentenceCeiling bound how long anybody is held, in
	// days. Nobody is in for an afternoon and nobody disappears for a season.
	SentenceFloor, SentenceCeiling = 2, 12
	// LawyerDaily is what somebody who knows the clerks charges for each day
	// they take off the end of it.
	LawyerDaily = 260
	// ServedRespect is what the street gives somebody for each day they did
	// without saying anything.
	ServedRespect = 4
	// TalkedRespect is what it costs to walk out by explaining who else was
	// involved, and TalkedGoodwill is what everybody named or not takes off you.
	TalkedRespect, TalkedGoodwill = 45, 30
	// FallGuyTrust is what the rest of your people take off you for handing one
	// of them over, and FallGuyDays is how long he is gone.
	FallGuyTrust = 25
	// ReleasedHeat is what attention falls to on the way out. They had their
	// pound of flesh; the file is closed and a new one has to be opened.
	ReleasedHeat = 12
)

// Charge is the weight of what a search found, in the terms a prosecutor would
// use. Money is not evidence of anything; a room full of rifles is a case that
// makes itself.
func (w *World) Charge() (int, string) {
	weight, worst := 0, ""
	if _, ok := w.TheArmoury(); ok && w.Stocked() > 0 {
		weight += 5
		worst = "a room full of crates"
	}
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && prop.Still && w.Own(l.ID) {
			weight += 3
			if worst == "" {
				worst = "a still in the back of " + l.Name
			}
		}
	}
	if w.Player.Charges > 0 {
		weight += 3
		if worst == "" {
			worst = "explosives in a bag under a bed"
		}
	}
	if w.Carrying() >= 10 {
		weight += 2
		if worst == "" {
			worst = fmt.Sprintf("%d units of stock nobody had papers for", w.Carrying())
		}
	}
	if w.Player.Weapon > 0 {
		weight += 2
		if worst == "" {
			worst = "a gun with your name on the licence and none on the paperwork"
		}
	}
	return weight, worst
}

// Sentence is how long the charge is worth, before anybody pays anybody. A man
// the police have been watching for months does not get what a first offender
// gets, so attention lengthens it.
func (w *World) Sentence(weight int) int {
	days := weight + w.Player.Heat/25
	return min(SentenceCeiling, max(SentenceFloor, days))
}

// Held reports whether the player is inside, which is true of nothing else in
// this game: alive, solvent, holding everything they held, and unable to do any
// of it.
func (w *World) Held() bool { return w.Player.Alive && w.Player.HeldUntil > w.Minute }

// DaysLeft is what is left of the sentence, rounded up, because a day inside
// with an hour left in it is still a day.
func (w *World) DaysLeft() int {
	if !w.Held() {
		return 0
	}
	return (w.Player.HeldUntil - w.Minute + 1439) / 1440
}

// Take is the arrest itself. The scene is the decision: go with them, hand them
// somebody else, or find out what the arresting officer costs.
func (w *World) Take(weight int, because string) {
	days := w.Sentence(weight)
	price := 400 + weight*220 + w.Player.Heat*10
	choices := []Choice{
		{ID: "serve", Label: "Go with them", Detail: fmt.Sprintf("%d days. Your businesses run without you and nobody defends your ground. The street respects a man who does his time quietly.", days)},
	}
	if fall := w.FallGuy(); fall != nil {
		choices = append(choices, Choice{ID: "fall:" + fall.ID, Label: "Give them " + fall.Name,
			Detail: fmt.Sprintf("%s answers for it and is gone %d days. Everybody else who works for you finds out what you do when it is one of them or you.", fall.Name, days)})
	}
	choices = append(choices, Choice{ID: "pay", Label: fmt.Sprintf("Find out what the officer wants ($%d)", price), Cost: price,
		Detail: "The evidence goes back in the van and the paperwork is never filed. It does nothing about the attention that brought them here."})
	w.Event = &Scene{
		ID: "arrest-" + ID(), Kind: "arrest", Source: "authored", Minute: w.Minute,
		Speaker: w.HolderID("detective"),
		Title:   "They are not here for the money",
		Body:    fmt.Sprintf("They found %s. This is not a fine and it is not a warning; somebody is going to be charged with it tonight.", because),
		Effect:  Effect{Reward: days}, // the sentence, carried to whoever resolves it
		Target:  because,
		Choices: choices,
	}
}

// FallGuy is whoever would answer for it instead: somebody standing on a door,
// or failing that anybody who works for you. Never a crewman — Leo is a
// character rather than a body.
func (w *World) FallGuy() *NPC {
	for _, l := range Locations {
		if n := w.PostedAt(l.ID); n != nil {
			return n
		}
	}
	for _, n := range w.OwnPeople() {
		if !w.Inside(n) {
			return n
		}
	}
	return nil
}

// ResolveArrest applies whichever way the player answered the door.
func (w *World) ResolveArrest(e *Scene, choice string) error {
	days := max(SentenceFloor, e.Effect.Reward)
	p := &w.Player
	switch {
	case choice == "pay":
		w.Player.Heat = min(100, w.Player.Heat+6)
		w.Log("It never happened", "The evidence goes back in the van. Nothing is filed and nothing is forgotten either; the officer knows exactly what you are worth now.", "danger")
		return nil
	case choice == "serve":
		w.Confine(days, e.Target)
		return nil
	case len(choice) > 5 && choice[:5] == "fall:":
		n := w.NPC(choice[5:])
		if n == nil || n.Dead {
			return fmt.Errorf("he is not there to give them")
		}
		for _, l := range Locations {
			if prop := w.Properties[l.ID]; prop != nil && prop.Posted == n.ID {
				prop.Posted = ""
			}
		}
		n.Held = w.Minute + days*1440
		// Everybody who works for you learns something about you tonight.
		for _, other := range w.OwnPeople() {
			if other.ID != n.ID {
				other.Trust = max(0, other.Trust-FallGuyTrust)
			}
		}
		w.Resent(n.ID, "", 60, "being handed to the police at the door")
		p.Respect = max(0, p.Respect-8)
		w.Player.Heat = max(0, w.Player.Heat-20)
		w.Log(n.Name+" goes instead", fmt.Sprintf("%s answers for what was in the building and is gone %d days. Everybody who works for you heard about it before morning.", n.Name, days), "danger")
		w.Report("police", "ONE MAN CHARGED AFTER "+upper(w.cityName())+" SEARCHES",
			fmt.Sprintf("%s has been charged following searches across the district. Police said no other arrests were expected.", n.Name))
		return nil
	}
	return fmt.Errorf("unknown decision")
}

// Confine puts the player inside and takes off them everything a man walking
// into a cell does not keep.
func (w *World) Confine(days int, because string) {
	p := &w.Player
	p.HeldUntil = w.Minute + days*1440
	p.HeldFor = because
	p.Location = "precinct"
	p.Weapon, p.Charges = 0, 0
	w.Log("Taken in", fmt.Sprintf("%d days for %s. Your businesses keep their hours and your ground keeps nobody on it. You will hear about all of it afterwards.", days, because), "danger")
	w.Report("police", "MAN CHARGED AFTER DISTRICT SEARCHES",
		fmt.Sprintf("%s has been remanded following searches across the district. Police said the investigation was continuing.", p.Name))
}

// Release is the morning it ends.
func (w *World) Release() {
	p := &w.Player
	if p.HeldUntil == 0 {
		return
	}
	served := p.HeldFor
	p.HeldUntil, p.HeldFor = 0, ""
	p.Location = p.Home
	p.Heat = min(p.Heat, ReleasedHeat)
	w.Log("Out", fmt.Sprintf("Nobody meets you. Whatever %s cost you happened while you were not there to watch it, and it is on the books now.", served), "personal")
}

// SitOut is doing the time. It advances the clock to the morning they open the
// door, and the street pays for it in the only currency it has.
func (w *World) SitOut() error {
	if !w.Held() {
		return fmt.Errorf("nobody is holding you")
	}
	days := w.DaysLeft()
	w.Advance(w.Player.HeldUntil - w.Minute)
	// Something in the city can stop the clock mid-sentence. If it did, the
	// player is still inside and the rest of the days are still theirs to do.
	if !w.Player.Alive || w.Minute < w.Player.HeldUntil {
		return nil
	}
	w.Player.Respect += days * ServedRespect
	w.Log("You did not say anything", fmt.Sprintf("%d days and not a name out of you. That is worth %d respect on this street and it is the only thing you earned in there.", days, days*ServedRespect), "personal")
	w.Release()
	return nil
}

// LawyerFee is what buying the rest of it costs.
func (w *World) LawyerFee() int { return w.DaysLeft() * LawyerDaily }

// LawyerReadiness explains why nobody can be paid, or returns "".
func (w *World) LawyerReadiness() string {
	if !w.Held() {
		return "Nobody is holding you"
	}
	if w.Reachable() < w.LawyerFee() {
		return "Not enough money anybody can get to"
	}
	return ""
}

// Lawyer is somebody who knows the clerks. It costs money the player cannot
// reach from a cell unless they had the sense to put it somewhere reachable.
func (w *World) Lawyer() error {
	if reason := w.LawyerReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	fee := w.LawyerFee()
	w.Player.Cash -= fee
	w.Log("Somebody who knows the clerks", fmt.Sprintf("$%d and the rest of it goes away. You walk out owing nothing and having earned nothing.", fee), "personal")
	w.Release()
	return nil
}

// Talk is the way out that costs everything except time. Naming the people who
// were involved gets the door open the same afternoon, and there is nobody in
// this city who does not hear about it.
func (w *World) Talk() error {
	if !w.Held() {
		return fmt.Errorf("nobody is holding you")
	}
	p := &w.Player
	p.Respect = max(0, p.Respect-TalkedRespect)
	for i := range w.Factions {
		w.Factions[i].Goodwill = max(-100, w.Factions[i].Goodwill-TalkedGoodwill)
	}
	for _, n := range w.OwnPeople() {
		n.Trust = 0
	}
	w.Log("You explained who else was involved", "The door opens the same afternoon. There is not an organization in this city that does not know why by the end of the week.", "danger")
	w.Report("police", "CHARGES DROPPED AFTER COOPERATION",
		fmt.Sprintf("Charges against %s were dropped after what police described as assistance with a continuing investigation.", p.Name))
	w.Release()
	return nil
}

// CustodyDay lets out anybody of the player's who was inside, and the player
// themselves if the clock passed the morning without them sitting it out.
func (w *World) CustodyDay() {
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Held > 0 && n.Held <= w.Minute {
			n.Held = 0
			w.Log(n.Name+" is out", "Whatever they had was not enough to keep him. He is available again, and he remembers where he was.", "personal")
		}
	}
	if w.Player.HeldUntil > 0 && w.Player.HeldUntil <= w.Minute {
		w.Release()
	}
}

// Inside is everybody of the player's the police are still holding, which is the
// reason they cannot be put on a door or sent anywhere.
func (w *World) Inside(n *NPC) bool { return n != nil && n.Held > w.Minute }

// cityName is what the paper calls where this happens.
func (w *World) cityName() string { return "district" }

// BailDaily is what buying somebody else out costs a day. Dearer than a lawyer
// for yourself, because nobody is doing you a favour.
const BailDaily = 320

// BailReadiness explains why somebody cannot be got out, or returns "".
func (w *World) BailReadiness(id string) string {
	n := w.NPC(id)
	if n == nil || !w.Inside(n) || n.Faction != w.PlayerOrganizationID() {
		return "Nobody of yours is in there under that name"
	}
	if w.Player.Cash < ((n.Held-w.Minute+1439)/1440)*BailDaily {
		return "Not enough cash"
	}
	return ""
}

// Bail gets one of the player's people out early. Somebody who was handed over
// at the door does not forget it, but somebody who was bought back remembers
// that too.
func (w *World) Bail(id string) error {
	if reason := w.BailReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n := w.NPC(id)
	days := (n.Held - w.Minute + 1439) / 1440
	if err := w.Pay(days * BailDaily); err != nil {
		return err
	}
	n.Held = 0
	n.Trust = min(100, n.Trust+18)
	for i := range w.Grudges {
		if g := &w.Grudges[i]; g.Holder == n.ID && g.Against == "" {
			g.Weight = max(0, g.Weight-40)
		}
	}
	w.Log(n.Name+" comes out", fmt.Sprintf("$%d for the %d days still on him. He walks out knowing exactly who paid it.", days*BailDaily, days), "personal")
	return nil
}

// HeldCollection is the share of a business's takings that still reaches
// somebody who is not there to take it, and only because one of their own
// people is standing on the door. The rest of it stays with whoever counted it.
const HeldCollection = .6

// CollectionShare is how much of what a business earns actually reaches the
// player. Standing in the city, all of it. In a cell, none of it — the staff
// keep the books and nobody comes to the station with an envelope. Somebody of
// yours on the door is the only reason any of it survives the week, which is
// the second reason that job exists.
func (w *World) CollectionShare(id string) float64 {
	if !w.Held() {
		return 1
	}
	if w.PostedAt(id) != nil {
		return HeldCollection
	}
	return 0
}
