package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"blackledger/core"
)

// Twelve addresses, and for months only five of them had a picture. Two of the
// seven that did not were places this project added itself and never went back
// to. Art is now generated for all of them, and this is what stops the next
// location being added without any: a building nobody painted shows a wireframe
// box, which makes the whole city look unfinished.

func artFile(t *testing.T, candidates ...string) bool {
	t.Helper()
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func TestEveryAddressInTheCityHasAPicture(t *testing.T) {
	root := "../.."
	if _, err := os.Stat(filepath.Join(root, "public", "art")); err != nil {
		t.Skip("no art tree beside this build")
	}
	// The hand-painted five are listed in the street manifest; the rest are
	// generated fronts named after the place.
	manifest, err := os.ReadFile(filepath.Join(root, "public", "art", "buildings.json"))
	if err != nil {
		t.Fatal(err)
	}
	var painted []struct{ ID, File string }
	if err := json.Unmarshal(manifest, &painted); err != nil {
		t.Fatal(err)
	}
	byID := map[string]string{}
	for _, b := range painted {
		byID[b.ID] = b.File
	}

	missing := []string{}
	for _, l := range core.Locations {
		if file, ok := byID[l.ID]; ok {
			if artFile(t, filepath.Join(root, "public", "art", file)) {
				continue
			}
			missing = append(missing, l.ID+" (manifest names "+file+", which is not there)")
			continue
		}
		if artFile(t,
			filepath.Join(root, "public", "art", "fronts", "front-"+l.ID+"-v1.jpg"),
			filepath.Join(root, "public", "art", "previews.json")) {
			// previews.json is checked separately below; a front is enough.
			if artFile(t, filepath.Join(root, "public", "art", "fronts", "front-"+l.ID+"-v1.jpg")) {
				continue
			}
		}
		missing = append(missing, l.ID+" has no painted front")
	}
	if len(missing) > 0 {
		t.Fatalf("addresses with no picture: %v. Run `mise run exteriors`.", missing)
	}
	t.Logf("all %d addresses in the city have a picture", len(core.Locations))
}

func TestEveryAddressHasAnInsideToStandIn(t *testing.T) {
	root := "../.."
	if _, err := os.Stat(filepath.Join(root, "public", "art", "rooms")); err != nil {
		t.Skip("no room art beside this build")
	}
	missing := []string{}
	for _, l := range core.Locations {
		if !artFile(t, filepath.Join(root, "public", "art", "rooms", "room-"+l.ID+"-v1.jpg")) {
			missing = append(missing, l.ID)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("addresses with no interior: %v. Run `mise run interiors`.", missing)
	}
}

func TestEveryMomentTheTheatreCanPlayHasAPlate(t *testing.T) {
	root := "../.."
	if _, err := os.Stat(filepath.Join(root, "public", "art", "scenes")); err != nil {
		t.Skip("no scene art beside this build")
	}
	missing := []string{}
	for _, kind := range []string{"killing", "explosion", "gunfight", "raid", "seizure", "arrest", "attack", "robbery"} {
		if core.Gravity(kind) == 0 {
			t.Fatalf("%q is not a moment the city rates", kind)
		}
		if !artFile(t, filepath.Join(root, "public", "art", "scenes", "scene-"+kind+"-v1.jpg")) {
			missing = append(missing, kind)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("moments with no plate: %v. Run `mise run scenes`.", missing)
	}
}

// The stated goal for a loud moment is that "the camera is taken there" — to
// the building it happened in, on the city view, with the headline afterwards.
// The theatre instead washed the whole city out to near-black and drew its own
// picture on top, which is the camera being taken *away* from the city.
//
// These are text guards on the interface sources, the same cheap kind that
// caught the hook below an early return. They fail if the theatre goes back to
// covering the city, or if the street stops being able to spotlight one address.
func TestTheCameraGoesToTheBuildingRatherThanOverTheCity(t *testing.T) {
	css, err := os.ReadFile("../../src/style.css")
	if err != nil {
		t.Skip("no stylesheet beside this build")
	}
	rule := regexp.MustCompile(`\.theatre\{[^}]*\}`)
	found := rule.FindString(string(css))
	if found == "" {
		t.Fatal("the theatre has no styling at all")
	}
	// An opaque wash over the whole city view is the thing being prevented.
	if strings.Contains(found, "inset:0") && !strings.Contains(found, "pointer-events:none") {
		t.Fatalf("the theatre covers the whole city and swallows its clicks: %s", found)
	}
	street, err := os.ReadFile("../../src/CityStreet.tsx")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"spotlight", "lit"} {
		if !strings.Contains(string(street), want) {
			t.Fatalf("the street cannot single out one address: no %q", want)
		}
	}
}

// People cross the city over real time, and the street only listed them in a
// band: names and minutes, in a box, above a picture of the city they were
// supposedly walking through. The last piece of the living city is seeing them
// on it — placed between the two fronts according to how far along they are.
func TestWalkersAreDrawnOnTheStreetAndNotOnlyListed(t *testing.T) {
	street, err := os.ReadFile("../../src/CityStreet.tsx")
	if err != nil {
		t.Skip("no interface sources beside this build")
	}
	body := string(street)
	for _, want := range []string{"walker-figure", "getBoundingClientRect", "progress"} {
		if !strings.Contains(body, want) {
			t.Fatalf("the street cannot place a walker between two addresses: no %q", want)
		}
	}
	css, err := os.ReadFile("../../src/style.css")
	if err != nil {
		t.Fatal(err)
	}
	rule := regexp.MustCompile(`\.walker-figure\{[^}]*\}`)
	found := rule.FindString(string(css))
	if found == "" {
		t.Fatal("a walker on the street has no styling")
	}
	// Positioned against the street rather than sitting in the flow, or it is
	// a list item again with a different name.
	if !strings.Contains(found, "position:absolute") {
		t.Fatalf("a walker is not placed on the street: %s", found)
	}
}

// The workspace clips what it cannot fit — it is overflow:hidden, which is what
// keeps the map and the sidebar in their own columns. That makes every column
// inside it a scroller in its own right, or its content is simply cut off with
// no scrollbar to reach it. The city column was not one: on a short screen the
// bottom of the room, the addresses and the band that says where something
// happened all ran below the fold and could not be reached at all.
func TestEveryColumnInsideTheWorkspaceCanBeScrolled(t *testing.T) {
	css, err := os.ReadFile("../../src/style.css")
	if err != nil {
		t.Skip("no stylesheet beside this build")
	}
	sheet := string(css)
	if !regexp.MustCompile(`\.workspace\{[^}]*overflow:hidden`).MatchString(sheet) {
		t.Skip("the workspace no longer clips, so its columns need not scroll")
	}
	for _, column := range []string{".city-pane", ".sidebar"} {
		rules := regexp.MustCompile(regexp.QuoteMeta(column) + `\{[^}]*\}`).FindAllString(sheet, -1)
		if len(rules) == 0 {
			t.Errorf("%s has no styling at all", column)
			continue
		}
		scrolls := false
		for _, rule := range rules {
			if strings.Contains(rule, "overflow:auto") || strings.Contains(rule, "overflow-y:auto") {
				scrolls = true
			}
		}
		if !scrolls {
			t.Errorf("%s sits in a workspace that clips and cannot be scrolled: %v", column, rules)
		}
	}
}

// "It'd be better to display actions below the interior render when inside a
// building." The work used to be a 330px column beside the picture, so a room
// with twenty-six things to do in it was ten feet of scrolling in a slot the
// width of a receipt. Below the picture it has the whole width, and the cards
// lay out across it instead of down it.
func TestTheWorkInARoomSitsUnderThePictureAndAcross(t *testing.T) {
	css, err := os.ReadFile("../../src/style.css")
	if err != nil {
		t.Skip("no stylesheet beside this build")
	}
	sheet := string(css)
	stage := regexp.MustCompile(`\.interior-stage\{[^}]*\}`).FindString(sheet)
	if stage == "" {
		t.Fatal("the room has no layout at all")
	}
	if !strings.Contains(stage, "grid-template-columns:minmax(0,1fr);") {
		t.Errorf("the room still puts the work in a column beside the picture: %s", stage)
	}
	work := regexp.MustCompile(`\.room-work\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(work, "grid-row:3") || !strings.Contains(work, "grid-column:1") {
		t.Errorf("the work is not the row under the picture: %s", work)
	}
	compact := regexp.MustCompile(`\.actions\.compact\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(compact, "repeat(auto-fill") {
		t.Errorf("the cards do not lay out across the room: %s", compact)
	}
}

// And the column beside the map stops repeating the room. Two copies of the
// same twenty-six cards is how the list got long enough to complain about.
func TestTheColumnBesideTheMapDoesNotRepeatTheRoom(t *testing.T) {
	source, err := os.ReadFile("../../src/main.tsx")
	if err != nil {
		t.Skip("no interface beside this build")
	}
	main := string(source)
	if !strings.Contains(main, "here-instead") {
		t.Error("standing in a place, the column beside the map offers no way into the room")
	}
	// The full list is still what a place you are NOT standing in gets.
	if !strings.Contains(main, "<ActionList actions={l.actions}") {
		t.Error("a place across the city lost its own panel")
	}
}

// "Feels like the action boxes should be consistently sized." They were sized by
// their own words: a card with a long sentence in it was taller than the one
// beside it, and the two grids in a section (what you can do, and what you
// cannot) settled on two different heights. Every card in a room is one size.
func TestEveryActionCardInARoomIsTheSameSize(t *testing.T) {
	css, err := os.ReadFile("../../src/style.css")
	if err != nil {
		t.Skip("no stylesheet beside this build")
	}
	sheet := string(css)
	grid := regexp.MustCompile(`\.actions\.compact\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(grid, "grid-auto-rows:1fr") {
		t.Errorf("rows are sized by their own contents, so cards differ between rows: %s", grid)
	}
	if !strings.Contains(grid, "align-items:stretch") {
		t.Errorf("cards do not fill the height of their row: %s", grid)
	}
	card := regexp.MustCompile(`\.actions\.compact \.action\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(card, "min-height:") {
		t.Errorf("a card has no floor to its height, so a short one and a long one differ: %s", card)
	}
	// The clamp is what keeps one wordy card from setting the height of every
	// card in the room.
	desc := regexp.MustCompile(`\.actions\.compact \.action \.desc\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(desc, "line-clamp") {
		t.Errorf("the detail is unbounded: %s", desc)
	}
}
