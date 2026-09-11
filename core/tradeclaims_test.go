package core

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// readSource reads a file from this package, for the guards that read the game
// rather than call it.
func readSource(t *testing.T, name string) string {
	t.Helper()
	body, err := os.ReadFile(name)
	if err != nil {
		t.Skip("no sources beside this build")
	}
	return string(body)
}

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
	for _, claim := range tradeClaims() {
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

// tradeClaim is one sentence out of the trade table, the trade whose paragraph
// says it, and whether the numbers agree.
//
// `of` is what makes the sweep below exact. Checking only that a trade's name
// appears somewhere in this file passed a butcher that had started explaining
// itself, because the wharf's own claim mentions a cold room — a trade can be
// named by somebody else's sentence and have none of its own.
type tradeClaim struct {
	of    string
	said  string
	holds bool
}

func tradeClaims() []tradeClaim {
	cover := func(kind string) int { return trades[kind].Cover }
	hides := func(kind string) int { return trades[kind].Hides }
	return []tradeClaim{
		// casino: "The best front there is."
		{"casino", "a casino covers more than anything else", func() bool {
			for kind := range trades {
				if kind != "casino" && cover(kind) > cover("casino") {
					return false
				}
			}
			return true
		}()},
		// club: "The best front there is after a casino."
		{"club", "a club covers less than a casino and as much as anything else",
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
		{"wharf", "a bonded shed hides more than anything else", func() bool {
			for kind := range trades {
				if kind != "wharf" && hides(kind) >= hides("wharf") {
					return false
				}
			}
			return true
		}()},
		{"wharf", "a bonded shed hides more than a cold room", hides("wharf") > hides("butcher")},
		{"wharf", "a bonded shed hides more than a yard of trucks", hides("wharf") > hides("haulage")},
		// exchange: "explains a great deal and hides nothing at all."
		{"exchange", "a trading floor hides next to nothing", hides("exchange") <= 1},
		// saloon: "explains cash better than a garage and nowhere near a
		// laundry, and hides almost nothing."
		{"saloon", "a public house covers more than a garage", cover("saloon") > cover("garage")},
		{"saloon", "a public house covers well short of a laundry", cover("saloon") < cover("laundry")},
		{"saloon", "a public house hides almost nothing", hides("saloon") <= 1},
		// undertaker: "a casket hides more than a butcher's cold room and less
		// than a bonded shed, and the books explain cash better than a garage's
		// and well short of a laundry's."
		{"undertaker", "a casket hides more than a cold room", hides("undertaker") > hides("butcher")},
		{"undertaker", "a casket hides less than a bonded shed", hides("undertaker") < hides("wharf")},
		{"undertaker", "an undertaker covers more than a garage", cover("undertaker") > cover("garage")},
		{"undertaker", "an undertaker covers well short of a laundry", cover("undertaker") < cover("laundry")},
	}
}

// And every trade that makes such a claim has to be in the list above.
//
// The guard above is a written list, and when it was written its own weakness
// was written with it: a nineteenth trade could explain itself in a paragraph
// and nobody would notice that nothing pinned the sentence. This closes that.
//
// It reads the trade table's own source, finds the comment above each trade,
// and asks whether that comment talks about the two things only a comment ever
// explains — how well the books explain cash, and how well the room hides a
// thing. A trade whose paragraph says either has to be named in the claims
// above. A trade whose paragraph is about something else is left alone: the
// restaurant's is about what goes wrong in a kitchen, which pins nothing and
// needs nothing pinned.
func TestEveryTradeThatExplainsItselfIsPinned(t *testing.T) {
	t.Parallel()
	source := readSource(t, "operations.go")
	pinned := map[string]bool{}
	for _, claim := range tradeClaims() {
		pinned[claim.of] = true
	}

	// The words a trade uses when it is placing itself against the others.
	about := []string{
		"explains cash", "launder", "hides", "hiding", "out of sight",
		"front there is", "for a thing to sit", "explains a great deal",
	}
	body := source[strings.Index(source, "var trades = map[string]Trade{"):]
	block := regexp.MustCompile(`((?:\t//[^\n]*\n)*)\t"(\w+)": \{`)
	claiming, covered := 0, 0
	for _, m := range block.FindAllStringSubmatch(body, -1) {
		comment, kind := m[1], m[2]
		if strings.TrimSpace(comment) == "" {
			continue
		}
		said := false
		for _, word := range about {
			said = said || strings.Contains(comment, word)
		}
		if !said {
			continue
		}
		claiming++
		if pinned[kind] {
			covered++
			continue
		}
		t.Errorf("the paragraph above %q places it against the other trades and nothing above pins what it says", kind)
	}
	t.Logf("%d trades explain themselves in a paragraph, %d of them pinned", claiming, covered)
	if claiming < 4 {
		t.Fatalf("only %d trades were read as making a claim, so this measures nothing", claiming)
	}
}
