package core

import (
	"strings"
	"testing"
)

// Offering to do a family a favour, and doing it.
//
// The audience across a table has four answers. Two of them — paying tribute and
// walking out — the harness plays. Two it has never touched, and one of those is
// the only one that goes anywhere: offering to carry something for them opens an
// arrangement, and an arrangement finished for a family is work done for them,
// which is the one road anybody in this city comes up by.
//
// So it is four steps rather than a branch: reach the room, take the seat, offer
// the favour, and finish it. Nothing had walked any of it.
func TestOfferingAFamilyAFavourOpensWorkYouCanFinish(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash, w.Player.Respect = 100, 40000, 200

	// Wherever somebody who can agree to something is standing.
	seat, actor := "", ""
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location = l.ID
		if who := w.AudienceActor(l.ID); who != "" && w.AudienceReadiness(l.ID) == "" {
			seat, actor = l.ID, who
			break
		}
	}
	if seat == "" {
		t.Skip("nobody in this city will sit down with the player")
	}

	w.Player.Location, w.Event = seat, nil
	var card Action
	for _, a := range w.Actions(seat) {
		if a.ID == "audience" {
			card = a
		}
	}
	if card.ID == "" || card.Disabled {
		t.Fatalf("no seat is offered at %s: %q", seat, card.Reason)
	}
	after, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "audience", Target: seat})
	if err != nil {
		t.Fatal(err)
	}
	w = after
	if w.Event == nil || w.Event.Kind != "audience" {
		t.Fatal("requesting an audience did not open one")
	}

	offer := ""
	for _, c := range w.Event.Choices {
		if c.ID == "work" {
			offer = c.ID
		}
	}
	if offer == "" {
		t.Fatalf("a seat across from %s offers no favour to do", actor)
	}
	taken, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "choice", Event: w.Event.ID, Choice: offer})
	if err != nil {
		t.Fatal(err)
	}
	w = taken

	// What it opens is an arrangement, with the family as the beneficiary —
	// which is what makes finishing it work done for them rather than a job.
	if w.Event == nil {
		t.Fatal("offering to carry something opened nothing")
	}
	job := w.Event
	if job.Beneficiary != actor {
		t.Fatalf("the favour is for %q and the seat was across from %q", job.Beneficiary, actor)
	}
	if len(job.Choices) == 0 {
		t.Fatal("the arrangement offers no way to take it on")
	}

	// And it can be taken on. A road with a door at the end of it that does not
	// open is not a road.
	// What the family thinks of the player before the favour is done, because
	// the arrangement runs its own minutes inside the command that accepts it.
	// Reading it afterwards and then advancing the clock reads +6 against +6 —
	// the rise had already happened, and it looked like a promise the card did
	// not keep.
	standing := w.faction(actor).Goodwill
	before := w.Player.Cash
	way := job.Choices[0].ID
	done, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision,
		Kind: "choice", Event: job.ID, Choice: way})
	if err != nil {
		t.Fatalf("the favour cannot be taken on with %q: %v", way, err)
	}
	w = done
	if w.Player.Cash < before {
		t.Logf("taking it on cost $%d", before-w.Player.Cash)
	}
	// Somebody has to have heard about it, or nothing happened at all.
	told := false
	for _, r := range w.History {
		told = told || strings.Contains(r.Title, "message") || strings.Contains(r.Text, "message")
	}
	if !told && w.Event == nil {
		t.Fatal("the favour was taken on and the city has nothing to say about it")
	}

	// And the card's own promise: "Completion improves their standing."
	if got := w.faction(actor).Goodwill; got <= standing {
		t.Fatalf("the favour was done and %s thinks of the player at %+d, having thought %+d",
			w.faction(actor).Name, got, standing)
	}
}
