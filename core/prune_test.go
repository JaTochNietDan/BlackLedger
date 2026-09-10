package core

import "testing"

// PrunePeople forgets the dead once nothing refers to them, which is what keeps
// a long save bounded. What counts as a reference was a list, and the list was
// missing the person standing in front of the player right now. An open scene
// names its speaker by id; so does a suspended arrangement; so does every loan
// on the player's book. Forget any of them and the id in the save points at
// nobody.

func crowded(t *testing.T) *World {
	t.Helper()
	w := New(11)
	// Past the threshold, so a prune actually runs.
	for i := 0; len(w.NPCs) <= MaxPeople*3/4; i++ {
		w.NPCs = append(w.NPCs, NPC{ID: "filler-" + itoa(i), Name: "Somebody " + itoa(i), Location: "bar"})
	}
	return w
}

func TestThePersonSpeakingIsNotForgotten(t *testing.T) {
	t.Parallel()
	w := crowded(t)
	w.NPCs = append(w.NPCs, NPC{ID: "teller", Name: "Ines Farrow", Location: "bar"})
	w.Event = &Scene{ID: ID(), Kind: "proposal", Speaker: "teller", Title: "A quiet word", Body: "“There is work.”"}
	w.Kill("teller", "Shot between the offer and the answer.")
	w.PrunePeople()
	if w.NPC("teller") == nil {
		t.Fatal("the city forgot the person making the offer the player is looking at")
	}
}

// Negative result, kept so it is not chased again: the book does NOT need
// protecting here. Kill writes off a current-life loan the moment it happens,
// and a loan carried over from an earlier life is already read through a nil
// check in both LoanDay and LoanDescription. The first version of this test
// asserted a dead debtor survived pruning, and it was measuring its own setup.

func TestASuspendedArrangementKeepsItsSpeaker(t *testing.T) {
	t.Parallel()
	w := crowded(t)
	w.NPCs = append(w.NPCs, NPC{ID: "patron", Name: "Vera Kohl", Location: "bar"})
	w.SuspendedJob = &SuspendedJob{Scene: &Scene{ID: ID(), Kind: "proposal", Speaker: "patron"}, Remaining: 90}
	w.Kill("patron", "Shot while the job was half done.")
	w.PrunePeople()
	if w.NPC("patron") == nil {
		t.Fatal("the city forgot who the unfinished arrangement is with")
	}
}

// And the consequence, so this is not filed as a tidiness complaint: declining
// an offer prints the name of whoever made it, straight off the lookup.
func TestDecliningAnOfferFromAForgottenPersonDoesNotCrash(t *testing.T) {
	t.Parallel()
	w := crowded(t)
	w.NPCs = append(w.NPCs, NPC{ID: "teller", Name: "Ines Farrow", Location: "bar"})
	w.Event = &Scene{ID: ID(), Kind: "proposal", Speaker: "teller", Title: "A quiet word",
		Body: "“There is work.”", Choices: []Choice{{ID: "accept", Label: "Take it"}, {ID: "decline", Label: "Pass"}}}
	w.Kill("teller", "Shot between the offer and the answer.")
	w.PrunePeople()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("declining crashed the city: %v", r)
		}
	}()
	if _, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "choice", Event: w.Event.ID, Choice: "decline"}); err != nil {
		t.Logf("declined with an error rather than a crash: %v", err)
	}
}

// An old save can already have lost a speaker, and no migration puts somebody
// back. The refusal has to survive reading it.
func TestDecliningAnOfferWithNoSpeakerLeftReadsWithoutAName(t *testing.T) {
	t.Parallel()
	w := New(11)
	w.NPCs = append(w.NPCs, NPC{ID: "teller", Name: "Ines Farrow", Location: "bar"})
	w.Event = &Scene{ID: ID(), Kind: "proposal", Speaker: "teller", Title: "A quiet word",
		Body: "“There is work.”", Choices: []Choice{{ID: "accept", Label: "Take it"}, {ID: "decline", Label: "Pass"}}}
	// The save the player loaded no longer contains them at all.
	kept := w.NPCs[:0]
	for _, n := range w.NPCs {
		if n.ID != "teller" {
			kept = append(kept, n)
		}
	}
	w.NPCs = kept
	next, err := Execute(w, Command{RequestID: ID(), Revision: w.Revision, Kind: "choice", Event: w.Event.ID, Choice: "decline"})
	if err != nil {
		t.Fatalf("declining failed: %v", err)
	}
	last := next.History[len(next.History)-1]
	if !contains(last.Text, "decline the proposal") {
		t.Fatalf("the refusal does not read without a name: %q", last.Text)
	}
}
