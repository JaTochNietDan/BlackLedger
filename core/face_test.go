package core

import (
	"os"
	"regexp"
	"strconv"
	"testing"
)

// The city hands out faces by a hash of a name, because nobody here is gendered
// by the rules and any face can belong to any name. That is fair to the cast
// and it can still hand a player somebody they do not recognise as themselves.
// Picking your own is the one place a player overrules the city about what they
// are looking at, so the core owns the number and refuses a bad one.

func TestAPlayerCanChooseTheirOwnFace(t *testing.T) {
	w := New(7)
	if w.Player.Face != 0 {
		t.Fatalf("a new life arrives with a face already picked: %d", w.Player.Face)
	}
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "face", Choice: "9"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Player.Face != 9 {
		t.Errorf("picked face 9 and got %d", next.Player.Face)
	}
	// And back to whatever the city would have given them.
	back, err := Execute(next, Command{Revision: next.Revision, Kind: "face", Choice: "0"})
	if err != nil {
		t.Fatal(err)
	}
	if back.Player.Face != 0 {
		t.Errorf("gave it back and kept %d", back.Player.Face)
	}
}

func TestTheCoreRefusesAFaceThatIsNotThere(t *testing.T) {
	w := New(7)
	for _, bad := range []string{"25", "-1", "", "four", "999"} {
		if _, err := Execute(w, Command{Revision: w.Revision, Kind: "face", Choice: bad}); err == nil {
			t.Errorf("%q was accepted as a face", bad)
		}
	}
	// The last one there is the point: the view knows how big the sheet is, and
	// so must the core, or the core cannot refuse anything.
	if _, err := Execute(w, Command{Revision: w.Revision, Kind: "face", Choice: "24"}); err != nil {
		t.Errorf("the last face in the sheet was refused: %v", err)
	}
}

// Choosing a face is not something you do in a room, and it must not cost a
// minute of anybody's evening.
func TestChoosingAFaceCostsNothingAndTakesNoTime(t *testing.T) {
	w := New(7)
	w.Player.Cash = 500
	next, err := Execute(w, Command{Revision: w.Revision, Kind: "face", Choice: "3"})
	if err != nil {
		t.Fatal(err)
	}
	if next.Minute != w.Minute || next.Player.Cash != w.Player.Cash {
		t.Errorf("it cost %d minutes and $%d", next.Minute-w.Minute, w.Player.Cash-next.Player.Cash)
	}
}

// The core refuses a face the sheet does not have, which only works while the
// core's count and the sheet's are the same number. Two constants in two
// languages agreeing by hand is exactly the kind of thing that quietly stops
// being true, so it is asked here.
func TestTheCoreAndTheSheetAgreeOnHowManyFacesThereAre(t *testing.T) {
	source, err := os.ReadFile("../src/Portrait.tsx")
	if err != nil {
		t.Skipf("no view to compare against: %v", err)
	}
	// Written across three lines since the view was formatted, so the two
	// numbers are found separately rather than as one phrase: a guard that
	// reads another language's source must not also be a guard on its layout.
	cols := regexp.MustCompile(`CAST_COLS\s*=\s*(\d+)`).FindSubmatch(source)
	rows := regexp.MustCompile(`CAST_ROWS\s*=\s*(\d+)`).FindSubmatch(source)
	if cols == nil || rows == nil {
		t.Fatal("the sheet no longer says how big it is; this guard cannot see it")
	}
	found := [][]byte{nil, cols[1], rows[1]}
	wide, _ := strconv.Atoi(string(found[1]))
	high, _ := strconv.Atoi(string(found[2]))
	if wide*high != CastFaces {
		t.Errorf("the sheet holds %d faces and the core will accept %d", wide*high, CastFaces)
	}
}
