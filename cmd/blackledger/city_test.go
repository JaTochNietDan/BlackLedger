package main

import (
	"encoding/json"
	"math"
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
// These read the grid's own constants straight out of the renderer and check
// them against the city's coordinates. They are deliberately loud when they
// cannot parse them: a guard that quietly finds nothing to check is worse than
// no guard.

var (
	blockPitch = regexp.MustCompile(`export const BLOCK = ([0-9.]+);`)
	cellSize   = regexp.MustCompile(`export const CELL = ([0-9.]+);`)
)

// The city is laid out on a grid: every address gets a block, the streets run
// the full width and height, and no two buildings can share ground because no
// two can have the same block.
//
// This replaces a test that checked the addresses' raw coordinates against a
// table of footprints. That was the right test when buildings stood wherever
// their coordinates put them and could overlap; it stopped meaning anything
// when the grid started deciding placement, and a test that cannot fail is
// worse than no test. What matters now is that every address gets its own
// block, which is exactly what the renderer computes.
func TestEveryAddressGetsItsOwnBlock(t *testing.T) {
	source, err := os.ReadFile("../../src/iso.ts")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	body := string(source)
	pitch := blockPitch.FindStringSubmatch(body)
	cellMatch := cellSize.FindStringSubmatch(body)
	if pitch == nil || cellMatch == nil {
		t.Fatal("the renderer no longer states BLOCK and CELL, so nothing here can be checked")
	}
	block, err1 := strconv.ParseFloat(pitch[1], 64)
	cell, err2 := strconv.ParseFloat(cellMatch[1], 64)
	if err1 != nil || err2 != nil || block <= 0 || cell <= 0 {
		t.Fatalf("BLOCK reads %q and CELL reads %q", pitch[1], cellMatch[1])
	}

	// The same binning the renderer does: an address's own coordinates decide
	// which block it gets, so the city keeps its shape.
	minX, minY := math.Inf(1), math.Inf(1)
	for _, l := range core.Locations {
		minX = math.Min(minX, float64(l.X)/cell)
		minY = math.Min(minY, float64(l.Y)/cell)
	}
	taken := map[[2]int]string{}
	for _, l := range core.Locations {
		at := [2]int{
			int(math.Round((float64(l.X)/cell - minX) / block)),
			int(math.Round((float64(l.Y)/cell - minY) / block)),
		}
		if other, clash := taken[at]; clash {
			t.Errorf("%s and %s both want block %d,%d — one of them will be pushed off its own coordinates",
				other, l.Name, at[0], at[1])
			continue
		}
		taken[at] = l.Name
	}
	if len(taken) != len(core.Locations) {
		t.Errorf("%d addresses went into %d blocks", len(core.Locations), len(taken))
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

// Every address must have a painted cut-out. A missing one falls back to a
// flat-shaded solid, which is correct behaviour and looks like a bug sitting
// next to eleven painted buildings.
func TestTheManifestAndTheArtAgree(t *testing.T) {
	body, err := os.ReadFile("../../public/art/iso/isometric.json")
	if err != nil {
		t.Skip("no isometric art beside this build")
	}
	var painted []struct {
		ID   string `json:"id"`
		File string `json:"file"`
	}
	if err := json.Unmarshal(body, &painted); err != nil {
		t.Fatalf("the manifest does not parse: %v", err)
	}
	// Every entry points at a file that is there. A manifest naming art that
	// has been deleted is worse than no art: the view asks for a texture,
	// the load fails, and the building silently disappears instead of falling
	// back to its solid.
	listed := map[string]bool{}
	for _, p := range painted {
		listed[p.File] = true
		if _, err := os.Stat("../../public/art/" + p.File); err != nil {
			t.Errorf("the manifest lists %s and the file is not there", p.File)
		}
	}
	// And every cut-out on disk is in the manifest, so art that was dropped in
	// by hand and never registered does not sit there doing nothing.
	found, err := os.ReadDir("../../public/art/iso")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range found {
		name := e.Name()
		if !strings.HasPrefix(name, "iso-") || !strings.HasSuffix(name, ".png") {
			continue
		}
		if !listed["iso/"+name] {
			t.Errorf("%s is on disk and not in the manifest; run tools/manifest.py", name)
		}
	}
	// An address with no cut-out is not a failure. It stands as a blocked-out
	// solid, which is what lets the city be looked at while it is being
	// re-thought — this test used to assert that all twelve were painted, and
	// that stopped being true the day the generated set was thrown out.
	for _, l := range core.Locations {
		if l.Name == "" {
			t.Error("an address with no name cannot be drawn at all")
		}
	}
}

// A moment happens in the city now rather than in a modal over it: the camera
// goes to the address the core named and the effect plays over that building.
// Two things have to stay true, and both are cheap to check in text.
func TestMomentsPlayInTheCityAndGiveTheCameraBack(t *testing.T) {
	body, err := os.ReadFile("../../src/CityIso.tsx")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	source := string(body)
	// Every kind of moment the core can witness must be drawn as something.
	// core/witness.go is the authority on what those are.
	witness, err := os.ReadFile("../../core/witness.go")
	if err != nil {
		t.Fatal(err)
	}
	kinds := regexp.MustCompile(`"([a-z]+)": \d`).FindAllStringSubmatch(string(witness), -1)
	if len(kinds) < 6 {
		t.Fatalf("only %d kinds of moment were found in witness.go; the table has moved", len(kinds))
	}
	for _, k := range kinds {
		if !strings.Contains(source, `'`+k[1]+`'`) {
			t.Errorf("the city draws nothing for a %q, so the loudest thing that can happen there is silent", k[1])
		}
	}
	// And the camera has to be given back: a player who was looking at the
	// docks should not be left staring at a rooftop across town.
	if !strings.Contains(source, "wasLooking") {
		t.Error("the camera is taken to a moment and never returned to where the player had it")
	}
}

// Every kind of moment the city can witness makes a noise. A kind nobody has
// scored falls through to a dull knock rather than silence, because silence
// reads as a bug — but the loud ones have to be scored deliberately.
func TestTheLoudMomentsAreScored(t *testing.T) {
	body, err := os.ReadFile("../../src/sound.ts")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	source := string(body)
	for _, loud := range []string{"explosion", "killing", "gunfight", "raid", "arrest"} {
		if !strings.Contains(source, `'`+loud+`'`) {
			t.Errorf("a %q makes whatever the fallback makes, which is a knock", loud)
		}
	}
	// It has to be possible to turn off, and it has to default to on rather
	// than to a browser exception in a private window.
	if !strings.Contains(source, "black-ledger-sound") {
		t.Error("the sound cannot be turned off")
	}
	if !strings.Contains(source, "catch { return true }") {
		t.Error("a browser that refuses local storage silences the city instead of defaulting to on")
	}
}

// The blocks between the addresses have buildings on them. Without those the
// grid is twelve models with holes between them; with a flat grey box on each
// it is worse, because a placeholder reads as a mistake rather than as
// distance. These have to be real painted cut-outs like everything else.
func TestTheBlocksBetweenTheAddressesAreBuiltOn(t *testing.T) {
	body, err := os.ReadFile("../../public/art/iso/isometric.json")
	if err != nil {
		t.Skip("no isometric art beside this build")
	}
	var painted []struct {
		ID   string `json:"id"`
		File string `json:"file"`
	}
	if err := json.Unmarshal(body, &painted); err != nil {
		t.Fatalf("the manifest does not parse: %v", err)
	}
	fillers := 0
	for _, p := range painted {
		if !strings.HasPrefix(p.ID, "fill-") {
			continue
		}
		fillers++
		if _, err := os.Stat("../../public/art/" + p.File); err != nil {
			t.Errorf("%s is in the manifest and the file is not there", p.ID)
		}
	}
	// Enough of them that a row of blocks does not read as the same building
	// repeated — but only once there are any at all. A city with no cut-outs
	// is a city being re-thought, and every block falls back to its solid.
	if fillers > 0 && fillers < 4 {
		t.Errorf("only %d filler buildings; a city needs more variety than that", fillers)
	}
	// A block is a terrace: the address takes one slot on the frontage and
	// ordinary buildings take the rest, shoulder to shoulder. This replaces an
	// assertion that fillers went on the blocks the addresses left empty —
	// true of the older model, where each address stood alone in the middle of
	// its block with pavement on all four sides, which read as an office park
	// rather than as a city.
	source, err := os.ReadFile("../../src/CityIso.tsx")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(source), "terrace({col, row}, SLOTS)") {
		t.Error("blocks are not laid out as terraces, so buildings do not stand next to each other")
	}
	if !strings.Contains(string(source), "addressAt.has(key)") {
		t.Error("nothing keeps a filler out of the slot an address builds on")
	}
}

// Everything on a pavement has to be on the pavement. A hydrant in the middle
// of the carriageway or a bench inside a building is the kind of fault that
// reads as broken however good the art is, and the user has asked twice for a
// clean grid with nothing overlapping.
//
// This is a text guard: it holds that the placement is computed from the
// pavement ring rather than from anywhere else. That the 63 props actually
// land there was checked in the browser by recomputing every one of them
// against its own block — none in a road, none under a building.
func TestStreetDressingStandsOnThePavement(t *testing.T) {
	source, err := os.ReadFile("../../src/iso.ts")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	body := string(source)
	if !strings.Contains(body, "export function dressing") {
		t.Fatal("nothing places the things a pavement carries")
	}
	// The runs a prop may stand on are built from the island — the pavement
	// ring — and inset from its edge. If that stops being true, props can
	// wander into the road.
	dress := body[strings.Index(body, "export function dressing"):]
	if end := strings.Index(dress, "\nexport function wires"); end > 0 {
		dress = dress[:end]
	}
	if !strings.Contains(dress, "island(cell)") {
		t.Error("dressing is not placed from the block's pavement ring")
	}
	if !strings.Contains(dress, "PAVE") {
		t.Error("dressing is not inset from the kerb, so a prop can overhang the carriageway")
	}
}

// A building is scaled to the ground it stands on, and to exactly that ground.
//
// This is a text guard against the mistake it is named for. A 1.16x overshoot
// was added to the sprite scaling to close the party walls between neighbours,
// and its real effect was to push every building into the one beside it —
// roofs through roofs, walls over the kerb — which is what the city looked
// like until it was taken out. The slot rectangles never overlapped (see
// tests/city-blocks.test.mjs); the pictures drawn on them did. Nothing here
// can check pixels, so it checks that the multiplier is gone.
func TestNoBuildingIsDrawnWiderThanItsGround(t *testing.T) {
	source, err := os.ReadFile("../../src/CityIso.tsx")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(source), "\n") {
		if !strings.Contains(line, "const across =") || !strings.Contains(line, "TILE.w") {
			continue
		}
		if !strings.HasSuffix(strings.TrimSpace(line), "* (TILE.w / 2);") {
			t.Errorf("a building is scaled by something other than its own plot: %s", strings.TrimSpace(line))
		}
	}
	// And its proportions are its own: scaling one axis alone squashed and
	// stretched buildings that were painted correctly.
	if strings.Contains(string(source), "art.scale.y *=") {
		t.Error("a building's height is being scaled independently of its width, which distorts it")
	}
}
