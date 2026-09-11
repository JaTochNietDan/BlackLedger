package core

import "testing"

// What the trade table says about itself, against what it is.
//
// Every trade's cover and hiding are explained in one place only: the comment
// above it. Nothing else in this project tells a reader why a casino launders
// better than a yard of trucks, so a wrong comment is a wrong fact about the
// city and there is nowhere else to find out.
//
// Three of them were wrong. A saloon "explains cash about as well as a laundry"
// at twelve against eighteen; and the undertaker, written two nights ago,
// claimed a casket was "the best place in the district for a thing to sit" at
// six against a bonded shed's eight, and its books "about as well as a
// laundry" at thirteen.
//
// The ordering each sentence asserts is pinned here. It is a written list and
// that is the point: every line is a sentence somebody wrote in the source, and
// the list is how the sentence stays true when a number moves.
func TestTheTradeTableMeansWhatItSays(t *testing.T) {
	t.Parallel()
	cover := func(kind string) int { return trades[kind].Cover }
	hides := func(kind string) int { return trades[kind].Hides }

	for _, claim := range []struct {
		said  string
		holds bool
	}{
		// casino: "The best front there is."
		{"a casino covers more than anything else", func() bool {
			for kind := range trades {
				if kind != "casino" && cover(kind) > cover("casino") {
					return false
				}
			}
			return true
		}()},
		// club: "The best front there is after a casino."
		{"a club covers less than a casino and as much as anything else",
			cover("club") < cover("casino") && func() bool {
				for kind := range trades {
					if kind != "casino" && cover(kind) > cover("club") {
						return false
					}
				}
				return true
			}()},
		// wharf: "a bonded shed is the best place in the city for a thing to sit
		// without anybody looking at it — more than a cold room, more than a
		// yard of trucks."
		{"a bonded shed hides more than anything else", func() bool {
			for kind := range trades {
				if kind != "wharf" && hides(kind) >= hides("wharf") {
					return false
				}
			}
			return true
		}()},
		{"a bonded shed hides more than a cold room", hides("wharf") > hides("butcher")},
		{"a bonded shed hides more than a yard of trucks", hides("wharf") > hides("haulage")},
		// exchange: "explains a great deal and hides nothing at all."
		{"a trading floor hides next to nothing", hides("exchange") <= 1},
		// saloon: "explains cash better than a garage and nowhere near a
		// laundry, and hides almost nothing."
		{"a public house covers more than a garage", cover("saloon") > cover("garage")},
		{"a public house covers well short of a laundry", cover("saloon") < cover("laundry")},
		{"a public house hides almost nothing", hides("saloon") <= 1},
		// undertaker: "a casket hides more than a butcher's cold room and less
		// than a bonded shed, and the books explain cash better than a garage's
		// and well short of a laundry's."
		{"a casket hides more than a cold room", hides("undertaker") > hides("butcher")},
		{"a casket hides less than a bonded shed", hides("undertaker") < hides("wharf")},
		{"an undertaker covers more than a garage", cover("undertaker") > cover("garage")},
		{"an undertaker covers well short of a laundry", cover("undertaker") < cover("laundry")},
	} {
		if !claim.holds {
			t.Errorf("the source says %s and the table does not", claim.said)
		}
	}

	// And the list is worth no more than its length. Every trade whose comment
	// makes a claim has to be in it, which is five of them; a sixth trade
	// writing a sentence about itself and not adding a line here is the way
	// this quietly stops meaning anything.
	if len(trades) < 15 {
		t.Fatalf("only %d trades in this city, so this pins very little", len(trades))
	}
}
