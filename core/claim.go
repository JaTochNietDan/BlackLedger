package core

import "fmt"

// When a family stops asking for a share and asks for the place.
//
// A publican who buys, hires, pays over the rate and restocks ends a campaign
// the richest policy in the game — $26,087 over twenty-seven days — with no
// police attention, no seizures, and every milestone reached in ninety-eight
// campaigns out of a hundred. Measured: nothing in this city ever costs a
// legitimate operator anything except money, and the money is a week of one
// shop's takings at a time.
//
// That is a coherent arrangement and it is also the end of the story. A family
// watching somebody assemble half the district does not go on sending for an
// envelope. At some point the demand stops being a share and becomes the place,
// and what a player does about that is the decision the whole legitimate path
// was missing.
//
// It is deliberately late and deliberately rare. This is not another tax; it is
// the moment a quiet campaign acquires a problem it cannot pay its way out of
// cheaply.

const (
	// ClaimHoldings is how many earning addresses the player has to hold in a
	// family's own districts before that family wants one of them outright.
	ClaimHoldings = 3
	// ClaimStanding is the goodwill below which a family thinks of taking
	// rather than asking. Somebody they are fond of gets the envelope.
	//
	// It has to be under the standing at which a family sends any demand at
	// all, which is 25: above that they are content and `BusinessPressure`
	// passes over them entirely. Setting it higher made a claim that could
	// never fire.
	//
	// Setting it *at* 25 made the opposite mistake and is the more interesting
	// one. Every demand to a player holding a few places became a demand for a
	// place — four hundred and eighty out of four hundred and eighty, measured
	// — which is not an escalation, it is a different tax. Below nothing is the
	// answer: paying an ordinary share lifts a family twelve points and keeps
	// them asking for envelopes, and only somebody who has been refusing them
	// falls far enough for them to come for the deed. The escalation is the
	// relationship rather than the count.
	ClaimStanding = -10
	// ClaimBuyout is what buying them off costs, as a multiple of the ordinary
	// share. It is meant to hurt: this is the price of not deciding.
	ClaimBuyout = 6
)

// TheirClaim reports whether this family has stopped asking for a share of this
// address and wants the address, and what buying them off would cost.
func (w *World) TheirClaim(family, id string) (int, bool) {
	f := w.faction(family)
	prop := w.Properties[id]
	if f == nil || prop == nil || !w.Own(id) || f.Goodwill >= ClaimStanding {
		return 0, false
	}
	if w.HoldingsNear(family) < ClaimHoldings {
		return 0, false
	}
	return w.TheirShare(id) * ClaimBuyout, true
}

// HoldingsNear is how many earning addresses of the player's stand in the
// districts this family has a claim on. A family leans on what it can see.
func (w *World) HoldingsNear(family string) int {
	districts := map[int]bool{}
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop != nil && prop.Owner == family {
			districts[l.District] = true
		}
	}
	held := 0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop != nil && w.Own(l.ID) && prop.Income > 0 && districts[l.District] {
			held++
		}
	}
	return held
}

// HandOver gives a family the deed, which is the one thing paying never had to
// do. The people behind the counter stay where they are — they work there, and
// who is being paid is not their decision.
func (w *World) HandOver(id, family string) error {
	prop := w.Properties[id]
	f := w.faction(family)
	if prop == nil || f == nil || !w.Own(id) {
		return fmt.Errorf("that is not yours to give")
	}
	place, _ := PlaceByID(id)
	prop.Owner = family
	f.Goodwill = min(100, f.Goodwill+30)
	f.Power = min(100, f.Power+ClaimPower)
	w.NextPressure = w.Minute + ClaimQuiet
	w.Log("Somebody else's name over the door",
		fmt.Sprintf("%s is %s now. Whoever was behind the counter is still behind it, and is not the one being paid. %s %s it plainly and %s left you alone about the rest.",
			place.Name, f.Name, Leads(f.Name), Agree(f.Name, "took", "took"), Agree(f.Name, "has", "have")), "politics")
	return nil
}

const (
	// ClaimPower is what a family gains by taking a working business rather
	// than being paid out of one.
	ClaimPower = 8
	// ClaimQuiet is how long they leave you alone afterwards. Long enough to
	// be worth it and not long enough to be a strategy.
	ClaimQuiet = 5 * 1440
)
