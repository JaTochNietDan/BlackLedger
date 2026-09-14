package core

import (
	"fmt"
	"strings"
)

// What the day costs, and what happens when it cannot be met.
//
// It was all or nothing. A player a dollar short lost their security, their
// address and their crew's loyalty in a single night — and then the bill was
// $33 instead of $2,033, so it cleared every night after and nothing was ever
// unpaid again. Being briefly short was a cliff, and being persistently short
// was impossible: measured over twenty days of a player with nothing, the
// longest run of unpaid wages was one night.
//
// It gives things up in order now, one at a time, and each only while what is
// left of the bill is still out of reach. The security goes first, because it
// is the thing a person stops paying for first. Then the crew feel it, then the
// address, and the wages last: the people behind a counter are the last thing
// anybody stops paying and the first thing they lose by not.

// Settle pays the day's bill, giving up whatever the player cannot cover in the
// order somebody would give it up, and reports what it cost them.
func (w *World) Settle() {
	p := &w.Player
	lost := []string{}
	bill := w.DailyCost()

	// The security. Ten a day each and the first thing to go.
	if p.Cash < bill && p.Security > 0 {
		lost = append(lost, "Security leaves")
		bill -= 10 * p.Security
		p.Security = 0
	}
	// The crew. Their day goes unpaid, and they are people, so they notice.
	if p.Cash < bill && len(p.Crew) > 0 {
		if p.Crew[0].Loyalty > 0 {
			lost = append(lost, p.Crew[0].Name+" is not being paid and knows it")
			p.Crew[0].Loyalty = max(0, p.Crew[0].Loyalty-20)
		}
		bill -= 12 * len(p.Crew)
	}
	// The roof over their head, down to a rented room.
	if p.Cash < bill && p.Home != "room" {
		lost = append(lost, "your residence is now a rented room")
		bill -= w.HomeCost(p.Home) - w.HomeCost("room")
		p.Home = "room"
	}

	// And the wages, last, because they are the thing worth keeping.
	if p.Cash >= bill {
		p.Cash -= bill
		w.EverybodyGotPaid()
	} else {
		// There is nothing left to pay anybody with. What the player was
		// carrying goes towards the day and the counters go unpaid.
		p.Cash = 0
		if short := w.NobodyGotPaid(); short > 0 {
			lost = append(lost, plural(short, "hand", "hands")+" went unpaid")
		}
	}

	if len(lost) == 0 {
		w.Log("Accounts settled",
			fmt.Sprintf("$%d paid for housing, security, crew and wages.", bill), "business")
		return
	}
	w.Log("Your arrangements unravel",
		"You could not cover the bills. "+upper1(strings.Join(lost, "; "))+".", "danger")
}
