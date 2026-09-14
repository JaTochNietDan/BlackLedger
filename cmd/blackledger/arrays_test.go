package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The view blanked on a hand of cards because the core sent it nothing where it
// had been promised a list. `Board []Card` carried `omitempty`, so before the
// flop the key was simply not there, and `cards.board.length` is a blank screen
// rather than an empty table.
//
// That is a shape rather than an incident. Go has two ways of sending a list
// that is empty and neither of them is a list: `omitempty` leaves the key out,
// and a slice nothing has been put in is written as `null`. The view has to
// know which fields do that, remember to guard each one, and go on remembering
// as the core grows — and it only takes one to blank the screen.
//
// So the core does not do either any more. No list omits itself, and
// `FillLists` gives every empty one a body on the way out. These two hold that:
// one reads the core's own field tags, and one reads what actually comes down
// the wire.

var omitted = regexp.MustCompile(`^\s*\w+\s+(\[\]\S+)\s+.*json:"([a-z_]+),omitempty"`)

func TestNoListInTheCoreLeavesItselfOut(t *testing.T) {
	t.Parallel()
	names, err := filepath.Glob(filepath.Join("..", "..", "core", "*.go"))
	if err != nil || len(names) == 0 {
		t.Fatal("no core to read", err)
	}
	lists, wrong := 0, []string{}
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		body, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(body), "\n") {
			if !strings.Contains(line, `json:"`) || !strings.Contains(line, "[]") {
				continue
			}
			lists++
			if omitted.MatchString(line) {
				wrong = append(wrong, filepath.Base(name)+": "+strings.TrimSpace(line))
			}
		}
	}
	t.Logf("%d list fields carry a name the view reads", lists)
	if lists < 30 {
		t.Fatalf("only %d list fields found in the core, so this is reading it wrong", lists)
	}
	if len(wrong) > 0 {
		t.Fatalf("these lists are left out of the payload when they are empty, and a view that "+
			"counts them sees a blank screen:\n%s", strings.Join(wrong, "\n"))
	}
}

// The things that really are absent rather than empty: a scene nobody is in,
// a result from a hand nobody has played, a page nobody has printed. Each is
// one thing or no thing, never a list, and the view asks whether it is there.
var mayBeNothing = map[string]bool{
	"pool_tournament": true, // No current-life tournament at this location.
	"pool":            true, // No current-life billiards rack at the player's location.
	"stroke":          true, // A billiards rack has no previous stroke before its break.
	"rent_register":   true, // Non-residential premises have no lodging register.
	"event":           true, "last_result": true, "press": true, "cards": true, "epitaph": true,
	"trade": true, "posted": true, "note": true, "room": true, "travel_note": true,
	"crossing": true, "away": true, "runs": true, "handle": true, "curtains": true,
	"organization": true, "residence": true, "scene": true, "speech": true, "table": true,
}

func nothings(v any, path string, out map[string]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, held := range t {
			if held == nil {
				if !mayBeNothing[k] {
					out[k] = path + "/" + k
				}
				continue
			}
			nothings(held, path+"/"+k, out)
		}
	case []any:
		for _, held := range t {
			nothings(held, path+"/*", out)
		}
	}
}

func TestNothingTheViewIsSentIsNothing(t *testing.T) {
	t.Parallel()
	a := testApp(t)
	out := request(a, "GET", "/api/state", "")
	if out.Code != 200 {
		t.Fatal(out.Body.String())
	}
	var payload any
	if err := json.Unmarshal(out.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	found := map[string]string{}
	nothings(payload, "", found)
	if len(found) == 0 {
		return
	}
	lines := []string{}
	for key, path := range found {
		lines = append(lines, key+" at "+path)
	}
	t.Fatalf("the view is sent nothing where it expects something. Either the core should send an "+
		"empty list, or the name belongs in mayBeNothing with a reason:\n%s", strings.Join(lines, "\n"))
}
