package core

import "testing"

// Measured on a real fifty-eight day campaign: the ledger holds sixty records
// and forty-nine of them were the same sentence — "Mara pays $45. A small
// favor, completed without questions." Repetition had not merely made the log
// unreadable, it had pushed every notable thing that ever happened out of the
// archive. The paper learned this within a day hours earlier; the ledger had
// it across a whole campaign.

func TestTheLedgerDoesNotSayTheSameThingFiftyTimes(t *testing.T) {
	t.Parallel()
	w := New(31)
	w.History = nil
	for i := 0; i < 49; i++ {
		w.Log("Envelope delivered", "Mara pays $45. A small favor, completed without questions.", "work")
	}
	if len(w.History) != 1 {
		t.Fatalf("the same errand forty-nine times filled %d records", len(w.History))
	}
	if w.History[0].Count != 49 {
		t.Fatalf("the record says it happened %d times", w.History[0].Count)
	}
}

// A collapsed repeat must still read as something that just happened: the
// result panel after an action is built by diffing record ids.
func TestACollapsedRepeatIsStillANewOutcome(t *testing.T) {
	t.Parallel()
	w := New(32)
	w.History = nil
	w.Log("Envelope delivered", "Mara pays $45.", "work")
	first := w.History[0].ID
	w.Log("Envelope delivered", "Mara pays $45.", "work")
	if len(w.History) != 1 {
		t.Fatalf("two identical errands made %d records", len(w.History))
	}
	if w.History[0].ID == first {
		t.Fatal("a collapsed repeat kept its old id, so the result panel would show nothing happened")
	}
	// And it is the newest thing in the log, not stranded where it started.
	w.Log("Something else", "A different line.", "personal")
	w.Log("Envelope delivered", "Mara pays $45.", "work")
	if w.History[len(w.History)-1].Title != "Envelope delivered" {
		t.Fatalf("the collapsed record is not the newest: %q", w.History[len(w.History)-1].Title)
	}
}

// Yesterday's errand is not today's, and a different amount is a different
// event.
func TestTheLedgerKeepsWhatIsActuallyDifferent(t *testing.T) {
	t.Parallel()
	w := New(33)
	w.History = nil
	w.Log("Envelope delivered", "Mara pays $45.", "work")
	w.Minute += 1440
	w.Log("Envelope delivered", "Mara pays $45.", "work")
	if len(w.History) != 2 {
		t.Fatalf("the same errand on two days made %d records", len(w.History))
	}
	w.Log("Envelope delivered", "Mara pays $60.", "work")
	if len(w.History) != 3 {
		t.Fatalf("a different amount collapsed into the same record")
	}
	w.Life++
	w.Log("Envelope delivered", "Mara pays $60.", "work")
	if len(w.History) != 4 {
		t.Fatal("a new life inherited the last one's errand")
	}
}
