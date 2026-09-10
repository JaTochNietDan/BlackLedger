# The overnight loop

Paste the whole of this file as the argument to `/loop` in a fresh session. It is
the working brief: what the loop is for, what is already built, what is queued,
and the mistakes that keep costing time.

---

## What this is

Work autonomously on Black Ledger in `~/Workarea/black-ledger`, one real slice
per tick, continuously. Schedule the next wakeup at **60s** every tick — that is
the scheduler's floor, and the user has said four minutes is too long. Never
idle. Never ask questions; make the call and write down what you decided.

The user's standing words: *"I don't have to think of all of these concepts,
you, the agent will think of these too and build them out and playtest and
simulate them yourself"* and *"gameplay expansion, adding businesses, finding
ways to link them together."*

Read `docs/LIVING_WORLD.md` first every tick. It has three sections:

- **Standing instructions** — never finished, never ticked off.
- **Inbox — open** — the user's exact words, until the game does the thing.
- **Inbox — answered** — kept, with a line saying where the work went.

Every line you write in those sections starts with an em dash. Anything without
one is the user's own words: never edited, summarised or split across sections.
One entry per horizontal rule. A message carrying both an idea and a standing
instruction stays whole, filed by what it mostly asks for, with your line saying
which half is standing.

---

## The queue

1. **Sound, and a better bandit.** `src/sound.ts` synthesises noises in-browser
   with no assets. The slot machine wants a real case, a payline, a coin tray
   and strips that roll rather than one face per drum, plus ambient room sound.
2. **Cars.** Still open: pictures of what you are buying, and what a car is
   worth in speed stated honestly. Plate is built for the player's own car and
   for your people's (`core/plate.go`, `core/theirs.go`).
3. **A back-room card game** whose other players are people from the city, which
   is the ambitious one. Craps is built.
4. **The rest of the open inbox.** Wording ("establish protection" for buying a
   business), the roulette table beside the wheel rather than under it and
   multiple chips down at once, newspaper pictures, the travel bar below the
   fold, cursing and threats from characters, a "bad blood" box that never
   clears, a family lead who sits you down somewhere he is not, and why Leo
   Carver can be sent on collections from every building in the city.

---

## What exists (do not rebuild)

The city is a Go core that owns every fact, a React view that invents nothing,
and a local model that writes encounters and can be switched off.

- **Work placement.** `Action.Anywhere` marks work that belongs to the player
  rather than a room; it is offered wherever they stand and filed on the People,
  Families and Ledger screens. Orders aimed at a building (`core/orders.go`)
  appear on that building's panel whether or not you are in it; going in
  yourself stays present-only.
- **The room.** Work sits below the picture, full width, in cards of one size,
  grouped by the core's own groups, with refusals folded away.
- **The tables.** A takeover screen: blackjack on baize, a roulette bowl with a
  turning head and a ball that lands in the core's pocket, and slot machines
  (`core/slots.go`, 17 in every hundred worked out over all 8,000 lines), and
  craps with the pass line, the don't and the field (`core/dice.go`, 1.414%
  measured over two million decisions). The player types the stake; the holder
  sets the house limit (`core/limits.go`). Sitting down is a decision the world
  knows about and it clears the felt.
- **Cars.** Bought on a forecourt, repaired at a garage, wrecked to a scrapyard,
  and fuelled at two filling stations. Petrol burns as you drive and a dry car
  is standing wherever it stopped. The city's own drivers do all of this too,
  and the trade happens when somebody walks in, not at midnight.
- **Violence.** A price on a name, going after somebody yourself, or sending one
  of your own (`core/strike.go`) — with capture, interrogation and the family
  learning who sent them.
- **Identity.** One decision picks a person's face and their voice
  (`core/voices.go`); the server sends the voice with the line.
- **Businesses.** Every address that earns can be bought and run. A place with
  no price does not change hands. Twenty-six addresses, twelve kinds; the
  pawnbroker is where what is taken off the street turns into money and where
  somebody short pawns the suit off their back.

---

## Fault shapes that keep biting

- **A content table written for the smaller city.** Three ids by name gated the
  entire business block; the same shape hid the car button and the respect
  requirement. Ask what table was written for a world that has since grown.
- **A rule and its button asking different questions.** The car button lived in
  the garage's case while its rule asked for a forecourt: offered where it was
  always refused, absent where it would have worked.
- **A trade settled at midnight that needs somebody standing there.** By
  midnight the counter is dark. Do it when they walk in.
- **A guard that cannot fail.** Four family-held places were refused for a
  reason other than the one under test. Put them in nobody's hands first.
- **A harness that jumps the clock.** Advancing twelve hours before looking for
  arrivals cancels every journey shorter than the jump. Use `aDay(w)`.
- **The interface stating a rule the core does not have.** "Insurance pays 2 to
  1" on a cloth with no insurance.

---

## How to work

Write the failing test first. **If it passes before you build anything, the
harness is wrong, not the world** — that has happened four times and each time
the measurement was asking the wrong question.

A test that passes when you break the code is not evidence. Break one fix at a
time with `sed`, run `go test -count=1`, and check the break is in the *right
direction*: disabling a rule so that nothing happens will pass a test that
asserts nothing happens.

Measure a change by breaking it: run the measurement, `sed` the new rule out
with `if false && ...`, run it again, report both numbers.

When you change a test rather than the code it guards, say so with the reason.
Do not claim an effect you have not measured.

**Playtest in the browser at least every other tick** and say what looked wrong.

### Gates

`mise run gate` (~90s) is all of it: the core suite and the API suites start
first and the format check, vet, tsc, the build and the 45 node tests run while
they go. `go test ./core` alone is ~90s; the balance tests run in parallel with
each other and the run is race-clean.

Balance baseline, six strategies (worker/investor/defiant/reckless/thief/smuggler):
deaths 0/0/51/82/77/0, median cash 12405/14156/7230/90/547/5500.
`mise run simulate > <scratchpad>/sim.json` then **parse** the JSON; grepping it
is useless. Each strategy also reports `mean_heat`, `seizures` and `runs_hurt`,
because deaths and cash cannot tell a safe policy from one whose money is taken
rather than whose life is. It also runs twelve cities for sixty days with nobody playing them
and reports that under `city_alone` — 11 organizations formed, 7 fell, 28 wars,
104 holdings changed hands. Organizations hold 52% of the city after 60 days and
65% after 240, the largest of them 28% and 57%. Families grow into unheld
premises (`core/expansion.go`) and a family that has won too much of the city
fractures from the inside.
That is the only measure here that can see a family fall: a campaign follows one
protagonist and does not last long enough for one to.

### Browser QA

```
go run ./cmd/qa-fixture <new.sqlite3> room|tables|bench|petrol|killing|leader|police
# kill every non-ops blackledger process, delete the sqlite3 AND -wal AND -shm
BLACK_LEDGER_DB=<db> BLACK_LEDGER_PORT=8862 nohup go run ./cmd/blackledger &
```

Drive it with the Chrome tools and **measure** — `getBoundingClientRect`, and
`/api/state` before and after. `w.Actions` returns nothing while `w.Event != nil`,
so answer the scene first. The tab's animation clock is frozen where you are
looking from: never rely on watching motion, make the resting state correct
without animation and measure the geometry.

### Rules that are not negotiable

- Never touch `.runtime/campaign.sqlite3` for testing.
- Use mise, never brew.
- Restart the live game on 8791 every tick, backing the save up first, and curl
  it afterwards: build to `.runtime/blackledger-ops.new`, `pkill -f
  blackledger-ops`, `mv` into place, `nohup` it.
- `w.Pay(a.Cost)` runs for every action: one that pays its own fee declares
  `Cost: 0` and uses `asks(...)`; one that pays out uses `add(...)`.
- Two RNG streams — set both `w.RNG` and `w.WorldRNG` when comparing runs.
- Role names must not assume a gender.
- Adding an address: `core/locations.json`, `PlaceIncome`, a trade in
  `core/operations.go` if it is a new kind, `mannerByPlace`, `streetTrades`,
  prompts in `tools/exteriors.py` and `tools/interiors.py` then run both, and a
  free block on the map (`BLOCK=3.6`, `CELL=48` in `src/iso.ts`).
