package core

import "testing"

// A seat across from a family used to be luck: it was offered wherever somebody
// who could speak for them happened to be standing that minute, so the only way
// to get one was to walk the city until a lieutenant turned up in a room. The
// diplomat policy spent seven hundred and ninety-one journeys buying twenty-six
// seats. People had telephones in 1930 and used them for exactly this.
func TestSendingWordGetsYouASeat(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event = nil
	w.Player.Cash, w.Player.Contacts = 5000, WordContacts
	f := &w.Factions[0]
	if f.ID == w.PlayerOrganizationID() {
		f = &w.Factions[1]
	}
	f.Goodwill = max(f.Goodwill, 0)
	where := w.homeOf(f.ID)
	if where == "" {
		t.Fatalf("%s has nowhere to receive anybody", f.Name)
	}

	// Their own hall, with nobody of theirs standing in it, offers nothing.
	for i := range w.NPCs {
		if w.NPCs[i].Location == where {
			w.NPCs[i].Location = "precinct"
		}
	}
	if who := w.AudienceActor(where); who == f.ID {
		t.Fatalf("a seat with %s was already going for nothing", f.Name)
	}

	if reason := w.SendWordReadiness(f.ID); reason != "" {
		t.Fatalf("word could not be sent: %s", reason)
	}
	before := w.Player.Cash
	if err := w.SendWord(f.ID); err != nil {
		t.Fatalf("sending word: %v", err)
	}
	if spent := before - w.Player.Cash; spent != WordCost {
		t.Fatalf("sending word cost $%d and says $%d", spent, WordCost)
	}
	if who := w.AudienceActor(where); who != f.ID {
		t.Fatalf("nobody was waiting at %s after word was sent: got %q", where, who)
	}

	// It is the rest of the day and not for ever.
	w.Advance(WordHolds + 1)
	if w.Expecting(f.ID) {
		t.Fatalf("%s was still expecting you a day later", f.Name)
	}
}

func TestNobodyCarriesAMessageForAStranger(t *testing.T) {
	t.Parallel()
	w := New(53)
	w.Event = nil
	w.Player.Cash, w.Player.Contacts = 5000, 0
	f := &w.Factions[0]
	if f.ID == w.PlayerOrganizationID() {
		f = &w.Factions[1]
	}
	f.Goodwill = max(f.Goodwill, 0)
	if w.Fitted("telephone") {
		t.Skip("the player starts with a telephone, so there is no stranger to test")
	}
	if reason := w.SendWordReadiness(f.ID); reason == "" {
		t.Fatalf("a stranger with no telephone and nobody who knows them sent word anyway")
	}
	if err := w.SendWord(f.ID); err == nil {
		t.Fatalf("sending word succeeded with nobody to carry it")
	}
	if w.Expecting(f.ID) {
		t.Fatalf("%s is expecting somebody who never sent word", f.Name)
	}
}
