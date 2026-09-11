package core

import (
	"strings"
	"testing"
)

// Pulling a story out of tomorrow's paper, walked from the thing that put it
// there.
//
// `spike` was on the list of actions the harness has never played. What makes
// it worth walking rather than calling is that it depends on three separate
// things being true at once: an arrangement with the editor, a story from the
// last day that is about you or about the police, and standing in the paper's
// own building. Nothing had ever put those three together, and the story has to
// come from the city rather than be written into the list by hand — a test that
// appends its own headline proves the pulling and not the road to it.
//
// The thing that writes the headline here is the arrest from the tick before:
// hand the police somebody of yours at your own door and the paper carries it.
func TestPullingTheStoryTheArrestPutInThePaper(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event, w.Player.Cash, w.Player.Respect, w.Player.Heat = nil, 400000, 300, 0
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Incorporate()
	if w.PlayerOrganization() == nil {
		t.Fatal("the player has no organization, so nobody answers to them")
	}

	// Somebody to hand over, signed on the way the game signs people on.
	signed := false
	for _, n := range w.NPCs {
		who := w.NPC(n.ID)
		w.Player.Location = who.Location
		if w.SignOnReadiness(who.ID) != "" {
			continue
		}
		if err := w.SignOn(who.ID); err != nil {
			t.Fatal(err)
		}
		signed = true
		break
	}
	if !signed {
		t.Skip("nobody in this city would sign on")
	}

	// The arrangement with the editor, made in the editor's own building.
	w.Player.Location = HeraldPlace
	if reason := w.RetainerReadiness("editor"); reason != "" {
		t.Fatalf("the editor cannot be retained: %s", reason)
	}
	if err := w.Retain("editor"); err != nil {
		t.Fatal(err)
	}
	if len(w.Spikeable()) != 0 {
		t.Fatal("there is something to pull before anything has happened")
	}

	// The door, and the answer that puts somebody else's name in the paper.
	w.Take(2, "what was in the back room")
	fall := ""
	for _, c := range w.Event.Choices {
		if strings.HasPrefix(c.ID, "fall:") {
			fall = c.ID
		}
	}
	if fall == "" {
		t.Fatal("a player with somebody of their own is offered nobody to give them")
	}
	after, err := Execute(w, Command{Kind: "choice", Event: w.Event.ID, Choice: fall, Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = after
	if len(w.Spikeable()) == 0 {
		t.Fatal("somebody was charged at your door and the paper has nothing about it")
	}

	// And the paper takes the call. Read off the room's own card: a card nobody
	// is offered is a road nobody can walk.
	w.Player.Location, w.Event = HeraldPlace, nil
	var card Action
	for _, a := range w.Actions(HeraldPlace) {
		if a.ID == "spike" {
			card = a
		}
	}
	if card.ID == "" {
		t.Fatal("the paper does not offer to pull anything")
	}
	if card.Disabled {
		t.Fatalf("the story cannot be pulled: %s", card.Reason)
	}
	heat, before := w.Player.Heat, len(w.Spikeable())
	pulled, err := Execute(w, Command{Kind: "spike", Revision: w.Revision})
	if err != nil {
		t.Fatal(err)
	}
	w = pulled
	if len(w.Spikeable()) >= before {
		t.Fatalf("%d stories could be pulled and %d still can", before, len(w.Spikeable()))
	}
	if w.Player.Heat > heat {
		t.Fatalf("pulling the story raised the city's interest from %d to %d", heat, w.Player.Heat)
	}
}

// And nothing is offered for a story the city never counted.
//
// Pulling takes the city's interest in a story out of the city's temperature.
// Half the kinds the paper files are worth nothing to that temperature on
// purpose — a notice that there was music at your club, a column, an obituary
// — and those were being offered as something to pull. The player paid, was
// told the city's interest went with it, and the city's interest did not move,
// with a one in seven chance of losing the arrangement at the paper on top.
func TestNothingIsPulledForAStoryTheCityNeverCounted(t *testing.T) {
	t.Parallel()
	w := New(23)
	w.Event, w.Player.Cash, w.Player.Respect, w.Player.Heat = nil, 400000, 300, 0
	for _, id := range []string{"laundry", "garage", "casino"} {
		w.Properties[id].Owner = "player:1"
	}
	w.Player.Location = HeraldPlace
	if err := w.Retain("editor"); err != nil {
		t.Fatal(err)
	}

	// A story about a place of the player's, of a kind the temperature scores
	// at nothing. Written here rather than played, because what is under test
	// is the filter and not the road to the headline.
	weightless := ""
	for kind, weight := range scrutinyWeight {
		if weight == 0 {
			weightless = kind
			break
		}
	}
	if weightless == "" {
		t.Skip("every kind of story costs the city something")
	}
	place, _ := PlaceByID("casino")
	w.Report(weightless, "MUSIC AT "+upper(place.Name), "There was music at "+place.Name+".")
	filed := Story{}
	for _, s := range w.News {
		if strings.HasPrefix(s.Headline, "MUSIC AT ") {
			filed = s
		}
	}
	if filed.ID == "" {
		t.Fatal("the paper did not carry it, so this measures nothing")
	}
	// And the paper says it is about a place of the player's, or the weight is
	// not what is keeping it off the list.
	if subject := w.SubjectOf(filed); subject.Kind != "place" || !w.Own(subject.ID) {
		t.Fatalf("the story reads as being about %s %q rather than about premises of the "+
			"player's, so this measures something else", subject.Kind, subject.Name)
	}
	if n := len(w.Spikeable()); n != 0 {
		t.Fatalf("%d %s stories are offered as worth pulling, and pulling one takes nothing "+
			"out of what the city thinks", n, weightless)
	}
	for _, a := range w.Actions(HeraldPlace) {
		if a.ID == "spike" && !a.Disabled {
			t.Fatal("the paper offers to pull a story that is costing the player nothing")
		}
	}
}
