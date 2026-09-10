package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	css := rawSource(t, "src/style.css")
	rule := regexp.MustCompile(`\.theatre\{[^}]*\}`)
	found := rule.FindString(css)
	if found == "" {
		t.Fatal("the theatre has no styling at all")
	}
	// An opaque wash over the whole city view is the thing being prevented.
	if strings.Contains(found, "inset:0") && !strings.Contains(found, "pointer-events:none") {
		t.Fatalf("the theatre covers the whole city and swallows its clicks: %s", found)
	}
	street := source(t, "src/CityStreet.tsx")
	for _, want := range []string{"spotlight", "lit"} {
		if !holds(street, want) {
			t.Fatalf("the street cannot single out one address: no %q", want)
		}
	}
}

// People cross the city over real time, and the street only listed them in a
// band: names and minutes, in a box, above a picture of the city they were
// supposedly walking through. The last piece of the living city is seeing them
// on it — placed between the two fronts according to how far along they are.
func TestWalkersAreDrawnOnTheStreetAndNotOnlyListed(t *testing.T) {
	t.Parallel()
	street := source(t, "src/CityStreet.tsx")
	body := street
	for _, want := range []string{"walker-figure", "getBoundingClientRect", "progress"} {
		if !strings.Contains(body, want) {
			t.Fatalf("the street cannot place a walker between two addresses: no %q", want)
		}
	}
	css := rawSource(t, "src/style.css")
	rule := regexp.MustCompile(`\.walker-figure\{[^}]*\}`)
	found := rule.FindString(css)
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
	t.Parallel()
	css := rawSource(t, "src/style.css")
	sheet := css
	if !regexp.MustCompile(`\.workspace\{[^}]*overflow:hidden`).MatchString(sheet) {
		t.Skip("the workspace no longer clips, so its columns need not scroll")
	}
	for _, column := range []string{".city-pane", ".sidebar"} {
		rules := regexp.MustCompile(regexp.QuoteMeta(column)+`\{[^}]*\}`).FindAllString(sheet, -1)
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
	t.Parallel()
	css := rawSource(t, "src/style.css")
	sheet := css
	stage := regexp.MustCompile(`\.interior-stage\{[^}]*\}`).FindString(sheet)
	if stage == "" {
		t.Fatal("the room has no layout at all")
	}
	if !strings.Contains(stage, "grid-template-columns:minmax(0,1fr);") {
		t.Errorf("the room still puts the work in a column beside the picture: %s", stage)
	}
	// Below the picture and in the same column, asked as an order rather than
	// as a row number: putting the head of the room in row one moved all three
	// of these down and failed a guard that was only ever about which comes
	// first.
	row := func(selector string) int {
		rule := regexp.MustCompile(regexp.QuoteMeta(selector) + `\{[^}]*\}`).FindString(sheet)
		found := regexp.MustCompile(`grid-row:(\d+)`).FindStringSubmatch(rule)
		if found == nil {
			t.Fatalf("%s has no row of its own, so it falls wherever the grid puts it", selector)
		}
		n, _ := strconv.Atoi(found[1])
		if !strings.Contains(rule, "grid-column:1") {
			t.Errorf("%s is not full width: %s", selector, rule)
		}
		return n
	}
	if room, work := row(".room"), row(".room-work"); work <= room {
		t.Errorf("the work is on row %d and the picture on row %d, so it is not under it", work, room)
	}
	compact := regexp.MustCompile(`\.actions\.compact\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(compact, "repeat(auto-fill") {
		t.Errorf("the cards do not lay out across the room: %s", compact)
	}
}

// And the column beside the map stops repeating the room. Two copies of the
// same twenty-six cards is how the list got long enough to complain about.
func TestTheColumnBesideTheMapDoesNotRepeatTheRoom(t *testing.T) {
	t.Parallel()
	source := source(t, "src/main.tsx")
	main := source
	if !holds(main, "here-instead") {
		t.Error("standing in a place, the column beside the map offers no way into the room")
	}
	// The full list is still what a place you are NOT standing in gets.
	if !holds(main, "<ActionList actions={l.actions}") {
		t.Error("a place across the city lost its own panel")
	}
}

// "Feels like the action boxes should be consistently sized." They were sized by
// their own words: a card with a long sentence in it was taller than the one
// beside it, and the two grids in a section (what you can do, and what you
// cannot) settled on two different heights. Every card in a room is one size.
func TestEveryActionCardInARoomIsTheSameSize(t *testing.T) {
	t.Parallel()
	css := rawSource(t, "src/style.css")
	sheet := css
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

// Work that belongs to the player rather than to a room has three homes, and
// each one is where the thing it is about already lives: names on the People
// screen, understandings and enquiries on the family they concern, and what the
// police think on the ledger beside the rest of the accounts. A flag that
// removes work from the room panel without giving it somewhere else to be is
// how a button disappears.
func TestPlayerWorkHasSomewhereToBe(t *testing.T) {
	t.Parallel()
	source := source(t, "src/main.tsx")
	main := source
	for _, home := range []string{
		"<PeopleScreen world={w} actions={anywhere}",
		"<FamiliesScreen world={w} actions={anywhere}",
		"<LedgerScreen world={w} render={actionButton} actions={anywhere.filter(",
	} {
		if !holds(main, home) {
			t.Errorf("no screen is given the player's own work: %q is not in the interface", home)
		}
	}
	// And the room panel is still the thing dropping it, or it would be in two
	// places at once.
	if !holds(main, "&& !a.anywhere") {
		t.Error("the room panel no longer leaves the player's own work out")
	}
}

// The screen that lists everybody in the city was a list you read and then went
// hunting on a map for: a man owes you money and is overdue, and the card said
// so and offered nothing. A name you can deal with now carries the work on its
// own card, and a name you cannot carries the address they are standing at.
func TestANameYouKnowCanBeDealtWithOrFound(t *testing.T) {
	t.Parallel()
	screen := source(t, "src/PeopleScreen.tsx")
	s := screen
	if !strings.Contains(s, "a.subject === who.id") {
		t.Error("a person's card does not carry the work the core says is about them")
	}
	if !strings.Contains(s, "find-them") {
		t.Error("a person standing somewhere else cannot be reached from their card")
	}
	main := source(t, "src/main.tsx")
	// The work has to come from where the player actually is, or the card is
	// offering something the command layer will refuse.
	if !holds(main, "at={p.location} here={(w.locations.find(l => l.id === p.location)?.actions || []).filter( a => !!a.subject, )}") {
		t.Error("the people screen is not given the work available where the player is standing")
	}
}

// The wheel now takes its time arriving at the pocket, which is the one place a
// view could quietly start deciding things. Nothing that draws a game is
// allowed to roll for anything: the core spins the wheel and deals the cards,
// and the felt only arranges what it was told. Camera shake and noise elsewhere
// are presentation; a number on a table is a fact.
func TestNothingThatDrawsAGameRollsForAnything(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Tables.tsx", "Casino.tsx", "cards.ts"} {
		drawn := source(t, "src/"+name)
		if holds(drawn, "Math.random") {
			t.Errorf("%s rolls for something; the core owns every number on a table", name)
		}
	}
	tables := source(t, "src/Tables.tsx")
	s := tables
	if !strings.Contains(s, "wheel.pocket ?? 0") {
		t.Error("the pocket the ball lands in is not the one the core spun")
	}
	// And the ball has to travel to it rather than appear in it. The angle it
	// travels from used to be React state and is now a ref, because a render
	// must not be able to restart a flight — the fact being guarded is the
	// same one.
	if !strings.Contains(s, "ballAngle(was.ball, pocket") {
		t.Error("the ball no longer goes round to the pocket")
	}
}

// The table is drawn as a table: a bowl with a head that turns under the ball,
// and a cloth laid out the way a cloth is laid out with a chip on it. The one
// thing that must hold whatever the browser is doing is where the ball comes to
// rest — the pocket is written on the element, and the flight is only the
// journey to it. An animation is not allowed to be the only thing that puts the
// ball in the pocket: on a frozen clock it never finishes and holds its first
// frame for ever, which is a ball that never lands.
func TestTheBallsRestingPlaceDoesNotDependOnAnimation(t *testing.T) {
	t.Parallel()
	tables := source(t, "src/Tables.tsx")
	s := tables
	if !strings.Contains(s, "el.style.transform = to;") {
		t.Error("the ball's resting place is not written on the element")
	}
	if !strings.Contains(s, "f.cancel()") {
		t.Error("a flight that outlives its time is never cancelled, so it hides where the ball landed")
	}
	if !strings.Contains(s, "ballAngle(was.ball, pocket") {
		t.Error("the ball no longer travels to the pocket the core spun")
	}
	css := rawSource(t, "src/style.css")
	sheet := css
	for _, part := range []string{".wheel-bowl{", ".wheel-head{", ".wheel-cone{", ".chip{", ".cloth-cell{"} {
		if !strings.Contains(sheet, part) {
			t.Errorf("the table has no %s", part)
		}
	}
	// One ball, one head: a second copy of either rule left over from an older
	// wheel silently overrides the new one, which is how the ball spent an
	// afternoon sitting at the top of the bowl.
	// Counted at the start of a line, because ".wheel-bowl.falling .wheel-ball{"
	// contains the same text and is a different rule — the first version of this
	// guard counted that one too and reported a duplicate that was not there.
	for _, one := range []string{"\n.wheel-ball{", "\n.wheel-head{"} {
		if n := strings.Count(sheet, one); n != 1 {
			t.Errorf("%s is declared %d times; a leftover copy overrides the live one", strings.TrimPrefix(one, "\n"), n)
		}
	}
}

// A chip has a ring drawn inside it with an absolute inset. If the chip itself
// is not positioned, that ring escapes to the nearest positioned ancestor and
// draws a circle the width of the whole screen across the game — which is
// exactly what happened the first time a chip was laid in the flow of a line
// rather than pinned to the corner of a cell.
func TestTheChipsRingCannotEscapeTheChip(t *testing.T) {
	t.Parallel()
	css := rawSource(t, "src/style.css")
	sheet := css
	chip := regexp.MustCompile(`\n\.chip\{[^}]*\}`).FindString(sheet)
	if chip == "" {
		t.Fatal("there is no chip")
	}
	if !strings.Contains(chip, "position:relative") {
		t.Errorf("a chip is not positioned, so the ring inside it will escape: %s", chip)
	}
	ring := regexp.MustCompile(`\.chip i\{[^}]*\}`).FindString(sheet)
	if !strings.Contains(ring, "position:absolute") {
		t.Skip("the ring is no longer drawn with an inset")
	}
}

// "It should again be a separate scene that takes up the screen when you're
// playing it and you have to leave it rather than right now it just lives in a
// small box above the action bar. That's silly stuff. We need to stop doing that
// in future and always dedicate these games to their own screen."
//
// So this is the rule, written down where it can fail: a game is a screen you
// go into and come out of, never a panel drawn beside the staffing figures and
// the supply count.
func TestEveryGameGetsAScreenOfItsOwn(t *testing.T) {
	t.Parallel()
	main := source(t, "src/main.tsx")
	// Both takeovers hang off the same fact — the player has taken a seat the
	// world knows about — and both are dismissed by getting up.
	for _, screen := range []string{"<BackRoomScene", "<Casino"} {
		if !holds(main, screen) {
			t.Errorf("%s is not mounted, so that game has no screen of its own", screen)
		}
	}
	if !holds(main, "atTable && inTheBackRoom") || !holds(main, "atTable && !inTheBackRoom") {
		t.Error("the two takeovers do not divide the seat between them, so one of them draws over the other")
	}
	// And the room behind it must not draw a felt of its own. The interior
	// panel used to be handed the table as a prop, which is exactly the shape
	// the note was about.
	room := source(t, "src/Interior.tsx")
	if holds(room, "backroom") || holds(room, "<BackRoom") {
		t.Error("the room panel draws a card table beside the premises work again")
	}
	// The verbs of a hand belong to the felt, not to the room's list of work.
	if !holds(main, "!isBackRoomAction(a.id)") {
		t.Error("the room's action list offers the verbs of a hand between hiring and restocking")
	}
}

// "It seems to swap the results on the rollers at the end which is odd, they
// just flip around at random mid-end game. For example it shows 7-7- as it
// progresses then at the very end it flips to bell, lemon, cherry."
//
// The case asked "has everything stopped, or is this drum still going?" and
// showed the result only then, so a drum that had stopped while the others
// turned fell back to the strip's first symbol and all three jumped when the
// last one came down. The rule now is the brief's own: the resting state is
// correct without animation, and the blur is the only thing the animation does.
func TestADrumIsAlreadyShowingWhereItWillStop(t *testing.T) {
	t.Parallel()
	tables := source(t, "src/Tables.tsx")
	if !holds(tables, "drumFaces(strip, line, !!machine.pulled)") {
		t.Error("the case works out its own faces again instead of asking where the drums land")
	}
	if holds(tables, "settled || rolling[i]") {
		t.Error("a drum's face depends on whether its neighbours have stopped")
	}
	// The tray is allowed to wait for the last drum, because what a pull paid
	// is not a thing to announce while they are still going.
	if !holds(tables, "const settled = machine.pulled && !rolling.some(Boolean)") {
		t.Error("the tray no longer waits for the drums")
	}
}

// "We should probably show options like 'buy kerrigan haulage' before you can
// afford it instead of having it hidden. We probably should just show all
// hidden options tbh. Not sure if hiding them is productive."
//
// The core was already offering every one of them, refused with a sentence
// saying what would change it — 251 refusals across the city, all of them
// explained. What hid them was the panel, which folded them behind a "Show N
// you cannot do yet" and started closed. They start open. The toggle stays,
// because a way to tidy a long list is not the same thing as a wall.
func TestNothingYouCannotDoYetIsHiddenByDefault(t *testing.T) {
	t.Parallel()
	list := source(t, "src/ActionList.tsx")
	if !holds(list, "const [open, setOpen] = useState(true)") {
		t.Error("the work you cannot do with somebody standing here is folded away again")
	}
	if !holds(list, "const showBlocked = !hidBlocked[s.id] || !!needle") {
		t.Error("a room's refusals are folded away again")
	}
	// And the toggle is still there, reading the right way round.
	if !holds(list, "{showBlocked ? 'Hide' : 'Show'} {s.blocked.length} you cannot do yet") {
		t.Error("there is no way to tidy a long list of refusals away")
	}
}

// "When inside a building it shows this info at the bottom of the action list
// which is wrong. It'd probably be better to have a more fleshed out display of
// current building your in with the name and stuff up higher in the fold in a
// consistent place when you're inside."
//
// The interior is a grid with explicit rows for the picture, the people and the
// work, and the strip naming the room had no row of its own — so it auto-placed
// into an implicit row after all three and came out under the action list. It
// is row one now, and the three explicit rows moved down to make space. This
// guard is about the placement, because that is the whole of what went wrong.
func TestTheRoomSaysWhereYouAreAtTheTopOfIt(t *testing.T) {
	t.Parallel()
	css := rawSource(t, "src/style.css")
	for _, rule := range []string{
		".room-holder{grid-column:1;grid-row:1",
		".room{grid-column:1;grid-row:2",
		".room-people{grid-column:1;grid-row:3",
		".room-work{grid-column:1;grid-row:4",
	} {
		if !strings.Contains(css, rule) {
			t.Errorf("the interior no longer places %q, so it falls wherever the grid puts it", rule)
		}
	}
	// Four rows to put them in, or the last one auto-places again.
	if !strings.Contains(css, "grid-template-rows:auto auto auto 1fr") {
		t.Error("the interior grid has no row for the head of the room")
	}
	// And it says more than a name: who holds it, and the figures a player
	// standing in a business wants in the same place in every room.
	room := source(t, "src/Interior.tsx")
	for _, fact := range []string{"Condition", "Working at", "On the books", "Earns"} {
		if !holds(room, fact) {
			t.Errorf("the head of the room does not say %q", fact)
		}
	}
}

// An action can be about two people — asking somebody where a third person is —
// and the id only has room for one. The second travels on the action's choice,
// which the panel has to send back or the command asks about nobody.
func TestAnActionAboutTwoPeopleCarriesBothNames(t *testing.T) {
	main := source(t, "src/main.tsx")
	if !holds(main, "commit({kind: a.id, target: a.target, choice: a.choice})") {
		t.Error("an ordinary action card drops the second name it was given")
	}
	sum := source(t, "src/SumAction.tsx")
	if !holds(sum, "commit({kind: a.id, target: a.target, choice: a.choice, amount})") {
		t.Error("an action with a figure in it drops the second name it was given")
	}
}

// "When playing the slot machine we should show actual images for the stuff on
// the rollers."
//
// The drums showed the core's own words set in a row, which is a list of
// symbols rather than a machine. Every symbol on the core's strip is painted
// now, and a symbol the view has never heard of falls back to its word rather
// than to a blank drum — so adding one to the strip shows up as itself.
func TestEverySymbolOnTheStripIsPainted(t *testing.T) {
	// Raw and anchored to the start of a line. Asking `holds` for "cherry:"
	// matched "notcherry:", and leading the needle with a space did not help
	// because flattening trims it — so renaming a symbol passed the guard.
	reels := rawSource(t, "src/reels.ts")
	for _, s := range core.ReelStrip() {
		key := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(s.ID) + `:`)
		if !key.MatchString(reels) {
			t.Errorf("the strip has a %q on it and nothing paints one", s.ID)
		}
	}
	// The drums ask what to paint rather than working it out, and the fallback
	// is the word.
	tables := source(t, "src/Tables.tsx")
	if !holds(tables, "const art = reelArt(painted[i][at])") {
		t.Error("the case decides for itself what a symbol looks like")
	}
	// The fallback, named: `") : ("` matched every ternary in the file and
	// could not have failed.
	if !holds(tables, "art ? ( <svg") || !holds(tables, "/> ) : ( face )") {
		t.Error("a symbol nothing paints leaves the drum blank")
	}
}

// A field the player types a figure into says what is wrong with the figure,
// and it is used for a stake, a float, a house limit and a wage. "More than the
// $18 there is" is a statement about the player's pocket: true of a stake and
// false of a wage, where the ceiling is what the trade will carry.
func TestATypedFigureDoesNotClaimToKnowWhatIsInYourPocket(t *testing.T) {
	sum := source(t, "src/SumAction.tsx")
	if holds(sum, "More than the ${money(sum.most)} there is") {
		t.Error("the figure field tells a business owner what is in their pocket")
	}
	if !holds(sum, "`${money(sum.most)} is the most`") {
		t.Error("the figure field no longer says what the most is")
	}
}

// And a business of the player's says what it pays, now that the wage is a
// decision rather than the trade's own rate.
func TestTheRoomSaysWhatItPays(t *testing.T) {
	room := source(t, "src/Interior.tsx")
	if !holds(room, "{what: 'Pays', is: '$' + place.wage + '/day'}") {
		t.Error("a business of yours does not say what it pays")
	}
}
