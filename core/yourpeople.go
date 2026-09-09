package core

import "fmt"

// The player's organization has holdings, a name, a strength and quarrels, and
// one man. Every other organization in this city has nine people with standing,
// ambition and their own opinion of whoever is above them.
//
// This is the other half of being filed with the others: people who answer to
// you, who cost money every day, who make you harder to move on, who die when
// somebody comes for your premises, and who leave when they stop thinking it is
// worth it — sometimes taking a business with them.

const (
	// SigningCost is what it takes to put somebody on. It is a week up front,
	// the same arrangement the businesses use.
	SigningCost = 140
	// ShareCost is what a share out of your own pocket costs. It was written
	// as a bare 60 in three places: the gate, the payment and the button.
	ShareCost = 60
	// MemberWage is what one of your own costs a day.
	MemberWage = 14
	// MemberTrust is what somebody starts at when they sign on.
	MemberTrust = 45
	// TrustDrift is what a day of being paid and not being asked for anything
	// impossible is worth.
	TrustDrift = 1
	// DefectionTrust is the trust below which somebody starts thinking about
	// where else they could be.
	DefectionTrust = 20
	// MaxOwnPeople is how many will answer to one name at this level. Past this
	// the player is running a family, which is a different game.
	MaxOwnPeople = 6
)

// OwnPeople is everybody who answers to the player.
func (w *World) OwnPeople() []*NPC { return w.Members(w.PlayerOrganizationID()) }

// MemberWages is what your own people cost a day, which joins the rest of the
// bill alongside the wages of the people who work your businesses.
func (w *World) MemberWages() int { return len(w.OwnPeople()) * MemberWage }

// SignOnReadiness explains why somebody cannot be taken on, or returns "".
func (w *World) SignOnReadiness(id string) string {
	if !w.Incorporated() {
		return "Nobody signs on with one person. They sign on with something that has a name"
	}
	n := w.NPC(id)
	if n == nil || n.Dead {
		return "There is nobody of that description"
	}
	if n.Location != w.Player.Location {
		return "You would have to be talking to them"
	}
	if IsOfficial(n.ID) || w.isRoleHolder(n) {
		return n.Name + " has a job already"
	}
	if n.Faction == w.PlayerOrganizationID() {
		return n.Name + " already answers to you"
	}
	if n.Faction != "" {
		return n.Name + " answers to somebody else"
	}
	// Knowing somebody is not required to offer them work: walking up to a man
	// in a bar and asking is the whole of it. An earlier build required it, and
	// since nothing in the game makes a civilian known, nobody could ever be
	// signed on at all.
	if len(w.OwnPeople()) >= MaxOwnPeople {
		return fmt.Sprintf("%d is as many as answer to one name at this level", MaxOwnPeople)
	}
	if w.Player.Cash < SigningCost {
		return fmt.Sprintf("It takes $%d to put somebody on", SigningCost)
	}
	return ""
}

// SignOn puts somebody on. They are a person in the city who now answers to the
// player, with everything that follows from that.
func (w *World) SignOn(id string) error {
	if reason := w.SignOnReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(SigningCost); err != nil {
		return err
	}
	n := w.NPC(id)
	n.Faction, n.Rank, n.Role = w.PlayerOrganizationID(), RankSoldier, "Yours"
	n.Trust = MemberTrust
	w.MeetPerson(n.ID)
	w.Log(n.Name+" signs on", fmt.Sprintf("$%d up front and $%d a day. %s answers to you now, which means they stand in front of what comes at you and can decide one morning that it is not worth it.", SigningCost, MemberWage, n.Name), "politics")
	return nil
}

// LetGoReadiness explains why somebody cannot be let go, or returns "".
func (w *World) LetGoReadiness(id string) string {
	n := w.NPC(id)
	if n == nil || n.Dead || n.Faction != w.PlayerOrganizationID() {
		return "They do not answer to you"
	}
	return ""
}

// LetGo ends it. Somebody put out is somebody who remembers being put out.
func (w *World) LetGo(id string) error {
	if reason := w.LetGoReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n := w.NPC(id)
	n.Faction, n.Rank, n.Role, n.Trust = "", RankAssociate, "Out of work", 0
	w.Log(n.Name+" is let go", "The wage stops today. Nobody takes that as well as they pretend to.", "politics")
	return nil
}

// PayShareReadiness explains why nobody can be paid a share, or returns "".
func (w *World) PayShareReadiness(id string) string {
	if reason := w.LetGoReadiness(id); reason != "" {
		return reason
	}
	if w.NPC(id).Trust >= 100 {
		return "They could not think better of you than they already do"
	}
	if w.Player.Cash < ShareCost {
		return "Not enough cash"
	}
	return ""
}

// PayShare buys back some of what neglect costs.
func (w *World) PayShare(id string) error {
	if reason := w.PayShareReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(ShareCost); err != nil {
		return err
	}
	n := w.NPC(id)
	before := n.Trust
	n.Trust = min(100, before+25)
	w.Log("A share for "+n.Name, fmt.Sprintf("$%d out of your own pocket. They think of you at %d rather than %d.", ShareCost, n.Trust, before), "politics")
	return nil
}

// OwnPeopleDay is what having people costs and what neglecting them does. Trust
// drifts up while they are paid and down when the bill is not met, and somebody
// far enough down goes looking for somewhere else to be.
func (w *World) OwnPeopleDay() {
	people := w.OwnPeople()
	if len(people) == 0 {
		return
	}
	paid := w.Player.Cash >= w.DailyCost()
	for _, n := range people {
		if paid {
			n.Trust = min(100, n.Trust+TrustDrift)
			continue
		}
		n.Trust = max(0, n.Trust-6)
	}
	// One of them a day, at most, decides.
	for _, n := range people {
		if n.Trust >= DefectionTrust {
			continue
		}
		if w.WorldRandom() >= float64(n.Ambition)/300 {
			continue
		}
		w.defect(n)
		return
	}
}

// defect is somebody deciding it is not worth it. What they do about it depends
// on what kind of person they are and whether there is anything to take.
func (w *World) defect(n *NPC) {
	holdings := w.FamilyHoldings(w.PlayerOrganizationID())
	// Somebody ambitious enough, with somewhere to walk out to, takes it.
	if len(holdings) > 1 && n.Ambition >= 65 && !TemperamentOf(n).Loyal {
		taken := holdings[0]
		for _, id := range holdings {
			if w.Properties[id].Income > w.Properties[taken].Income {
				taken = id
			}
		}
		place, _ := PlaceByID(taken)
		w.Properties[taken].Owner = "independent"
		n.Faction, n.Rank, n.Role = "", RankSoldier, "Runs "+place.Name
		n.Location, n.Trust = taken, 0
		w.Log(n.Name+" is running "+place.Name+" now", fmt.Sprintf("They stopped being paid properly and stopped waiting. %s is not yours any more, and whoever took it knows every arrangement you have.", place.Name), "danger")
		w.Report("business", upper(n.Name)+" TAKES OVER "+upper(place.Name),
			fmt.Sprintf("%s is now run by %s, who is understood to have previously worked for its former proprietor.", place.Name, n.Name))
		return
	}
	n.Faction, n.Rank, n.Role, n.Trust = "", RankAssociate, "Out of work", 0
	w.Log(n.Name+" is gone", "They were not there this morning and nobody has seen them. It had been coming.", "politics")
}

// OwnPeopleDescription is who answers to the player, for the interface.
func (w *World) OwnPeopleDescription() []map[string]any {
	out := []map[string]any{}
	for _, n := range w.OwnPeople() {
		out = append(out, map[string]any{
			"id": n.ID, "name": n.Name, "trust": n.Trust,
			"temperament": TemperamentOf(n).Label, "wage": MemberWage,
			"restless": n.Trust < DefectionTrust,
		})
	}
	return out
}
