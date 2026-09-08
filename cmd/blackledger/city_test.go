package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"blackledger/core"
)

// The isometric city places every address from the coordinates the core keeps
// for it. Two buildings standing on the same ground is the one fault that
// cannot be styled away — it reads as broken however good the art gets — and
// the first version did it twice, putting The Mariner and Saint Agnes through
// the Bellwether Herald.
//
// This reads the footprint table and the spacing straight out of the renderer
// and checks them against the city's own coordinates. It is deliberately loud
// when it cannot parse them: a guard that quietly finds nothing to check is
// worse than no guard.

var (
	cellSize  = regexp.MustCompile(`export const CELL = ([0-9.]+);`)
	footprint = regexp.MustCompile(`(?m)^\s+([a-z]+): \[([0-9.]+), ([0-9.]+)\],`)
)

func TestNoTwoBuildingsStandOnTheSameGround(t *testing.T) {
	source, err := os.ReadFile("../../src/iso.ts")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	body := string(source)

	cellMatch := cellSize.FindStringSubmatch(body)
	if cellMatch == nil {
		t.Fatal("the renderer no longer states a CELL, so nothing here can be checked")
	}
	cell, err := strconv.ParseFloat(cellMatch[1], 64)
	if err != nil || cell <= 0 {
		t.Fatalf("CELL reads %q", cellMatch[1])
	}

	table := map[string][2]float64{}
	for _, m := range footprint.FindAllStringSubmatch(body, -1) {
		w, _ := strconv.ParseFloat(m[2], 64)
		d, _ := strconv.ParseFloat(m[3], 64)
		table[m[1]] = [2]float64{w, d}
	}
	if len(table) < 7 {
		t.Fatalf("only %d footprints were found in the renderer; the table has moved or been renamed", len(table))
	}

	// Every kind of place the city actually contains must be drawn deliberately
	// rather than falling back to a house.
	for _, l := range core.Locations {
		if _, ok := table[l.Type]; !ok {
			t.Errorf("%s is a %q and nothing in the renderer draws one", l.Name, l.Type)
		}
	}

	sizeOf := func(kind string) [2]float64 {
		if s, ok := table[kind]; ok {
			return s
		}
		return table["home"]
	}
	for i, a := range core.Locations {
		for _, b := range core.Locations[i+1:] {
			as, bs := sizeOf(a.Type), sizeOf(b.Type)
			ax, ay := float64(a.X)/cell, float64(a.Y)/cell
			bx, by := float64(b.X)/cell, float64(b.Y)/cell
			if ax < bx+bs[0] && bx < ax+as[0] && ay < by+bs[1] && by < ay+as[1] {
				t.Errorf("%s and %s stand on the same ground at CELL %g", a.Name, b.Name, cell)
			}
		}
	}
}

// And the switch has to keep offering the card view while the city is built,
// so a broken renderer never leaves the game unplayable.
func TestTheCardViewIsStillReachable(t *testing.T) {
	body, err := os.ReadFile("../../src/main.tsx")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	if !strings.Contains(string(body), "CityStreet") || !strings.Contains(string(body), "CityIso") {
		t.Fatal("the city view and the card view are not both reachable")
	}
}

// The city is drawn on a GPU now, and a renderer that fails to start would
// leave a blank pane where the city was. Two things have to stay true: the
// card view remains reachable (checked above), and the addresses remain
// reachable without a mouse or WebGL at all.
func TestTheCityCanBeReadWithoutWebGL(t *testing.T) {
	body, err := os.ReadFile("../../src/CityIso.tsx")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	source := string(body)
	if !strings.Contains(source, "iso-reader") {
		t.Error("the city has no text alternative, so a browser without WebGL shows an empty pane")
	}
	// Every address, not a selection of them.
	if !strings.Contains(source, "state.locations.map") {
		t.Error("the text alternative does not list the city's own addresses")
	}
	// The camera must not be reset by an ordinary update: a player who has
	// zoomed in on the docks should stay there when an hour passes.
	if strings.Contains(source, "useEffect(frame") {
		t.Error("the camera is re-framed on every update, which throws away where the player was looking")
	}
}
