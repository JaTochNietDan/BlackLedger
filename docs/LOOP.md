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
   with no assets. All of it is built: the case, the payline, the coin tray, the
   room's own hum, and drums that travel through a run of the strip and come to
   rest on the faces the core dealt. Two things caught the drums out and are
   written in the fault shapes — a transition driven in the wrong frame order,
   and a second `.drum` rule inheriting `display:grid` from an earlier one.
   Nothing here is left open.
2. **Cars.** Speed is stated honestly already: the forecourt quotes a real
   journey from where the player is standing, in minutes on foot against
   minutes in the car being sold, and says when plate is weighing the figure
   down. Pictures are in too — three painted offline by `tools/cars.py` and
   shipped as files, drawn by `paintedCar` in `src/cityAssets.ts`, so that line
   was stale. Plate is built for the player's own car and for your people's
   (`core/plate.go`, `core/theirs.go`). Nothing here is left open.
3. **A back-room card game** whose other players are people from the city.
   Built: Texas hold'em behind the poolhall, no house and no edge,
   seats filled from whoever is in the room, four streets of betting with one
   raise, folding, and bluffing off `Ambition` (`core/backroom.go`,
   `core/holdem.go`).
   Screen in `src/BackRoomScene.tsx`, which takes the whole screen the way the
   casino does; what a night costs somebody is remembered: a heavy loser is
   sore, somebody you paid thinks better of you, and a room can be emptied.
   A grudge changes how somebody plays the player: they call light, raise on
   less and bluff more. Sound is in: the room hum while you are in there, a
   card for every card the core deals, and chips when the pot grows, one for
   each ante's worth. Nothing here is left open.
4. **The rest of the open inbox.** Checked this tick and nearly all of it is
   built: the wording, the roulette cloth and multiple chips, the travel bar,
   cursing and threats, the bad blood box, the family lead, and Leo Carver's
   collections are all answered in `docs/LIVING_WORLD.md`. What was genuinely
   left from the user's words was one question — **whether the map is the right
   way to move around the city at all** — and it is measured and part-answered
   in `docs/LIVING_WORLD.md`: the geography is load-bearing, so the map stays,
   and the plain list of addresses behind it is now nearest-first with the
   journey on each row.
5. **The city playing its own games.** Built: one hand a night in every room
   with a game behind it, a seat charge to whoever holds that room, and a
   falling-out for whoever is cleaned out (`core/citygame.go`). There are two
   of those rooms now, the poolhall and the bar — a game with no house belongs
   where there are people of an evening and nobody holding a float.
   An unwatched game runs on a third RNG stream of its own — shuffling off
   `WorldRNG` moved everything else the city does off-screen. The city already
   gambles at every room that runs a float (`core/floor.go`). Nothing here is
   left open: a game with no house belongs where there are people of an evening
   and nobody holding one, which is the poolhall and the bar, and both have it.
6. **The people behind the counter.** `Property.Hands` names the staff of every
   business at an address anybody holds, and somebody standing in front of you
   who works for a rival can be offered a place at one of yours
   (`core/hands.go`). A place a week behind on its wages is offered nobody and
   its position count does not reset, so emptying a counter by not paying it
   is a thing that stays done until the money does. Checked and not built:
   running the city dry of people to hire — sixty days of a city playing itself
   leaves eight held addresses, twenty-nine hands and forty-seven people still
   free, so a rule for it would be a rule that never fires. Asking one of them what they have seen was
   built and wired and answering — and answering with the state of the
   premises, which is on "Review the books" for nothing and no time, so the
   whole of what the conversation added was a footfall count the game already
   had. A counter now says three things a ledger never records: who is walking
   over and whose man he is, who has been asking after you and whether they are
   past asking, and the one thing this trade can see — the glass out of cars for
   a garage, who was short for a pawnbroker, where the cabs took somebody, who
   on this counter is thinking about leaving. Every one of the seventeen trades has a
   line of its own and no two say the same thing, which is guarded over the
   whole city rather than a sample — the first version of that guard asked five
   kinds and passed while eleven counters said nothing at all. Nothing here is
   left open.

---

7. **Businesses that link to each other. Done, and guarded.** All thirteen
   trades reach past their own income, and `TestEveryTradeReachesPastItsOwnIncome`
   in `core/reach_test.go` holds each one and does not hold it and measures the
   number it is supposed to change. It fails if a link stops working and it
   fails if a fourteenth trade is added with nobody asking about it. Do not
   rewrite this list from a reading of the code — run the test, it prints the
   figures:

   | trade | what holding it changes | not held | held |
   |---|---|---|---|
   | garage | half off what the car costs to keep | 7 | 4 |
   | filling | your own petrol at what it cost the pumps | 25 | 8 |
   | dealer | a car without the forecourt's margin | 4400 | 3300 |
   | haulage | a third off stocking everything else | 200 | 132 |
   | cabs | a ride when your own car cannot take you | no | yes |
   | scrapyard | more for the wreck when the yard is yours | 105 | 157 |
   | laundry | attention the books absorb, which nothing else takes off you | 0 | 18 |
   | burlesque | a night on, which draws (guarded in `night_test.go`) | no | yes |
   | pawn | the window at what the counter lent, not what it asks | 400 | 150 |
   | butcher | a cold room things sit in without being looked at | 20 | 24 |
   | restaurant | a dining room two families will sit down in | no | yes |
   | poolhall | the seat money off the city's own game | 0 | 32 |
   | casino | a bankroll the house plays out of, and a night of its own | no | yes |

   This list was wrong four times in one night — the burlesque, the butcher, the
   poolhall and the casino were all down as unlinked and were not. Two rules
   came out of that. **Contributing cover to laundering is not reaching**; every
   trade with a front does that, and counting it would have finished the
   question before it was asked. And **a link can be wired up, offered on a card,
   correct to read, and move nothing** — the burlesque's night did exactly that
   for as long as it existed. So the sweep asks the game rather than the code,
   and where a row can only ask whether something is offered it says so and
   points at the test that asks whether anybody comes.
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
- **A night at cards.** The back room is a sitting, not a hand: money goes on
  the table, everybody plays out of what is in front of them, the ante is a
  twentieth of the buy-in but never more than a quarter of the shortest stack,
  the cleaned-out go home and whoever is in the room takes the empty chair.
  Getting up or walking out picks your money up.
- **Arms.** Every gun and vest is on the counter at the docks at its own price,
  in any order. What you carry shifts an attempt on somebody — 17 in 200 with
  empty hands against 58 with a Thompson — and so does what you put in the hand
  of whoever you send. A search reaches whoever is standing with you.
- **Businesses. They pay, and they pay late.** Over two hundred commands a
  publican who hires, pays over the rate and restocks ends nearly the poorest
  policy in the game; over four hundred it ends the richest at $27,516, ahead of
  the worker and the investor both. The balance baseline ran at two hundred for
  the whole life of this project and warned on every single run that its own
  city measures need twenty game days, which two hundred does not reach. It runs
  at four hundred now. **Read no balance figure taken at two hundred commands.**
- **Businesses.** Every address that earns can be bought or taken and run. A
  family's seat has no price on purpose — it changes hands by force — and used
  to be the only kind of address in the city that ran on nothing: you fought a
  war for the busiest room in the game and won an income figure. The club is a
  trade now, with hands, drink, its own trouble and the best front after a
  casino. Pier 14 is one too: six dockers, rope and
  fuel, a crane that goes down, the best hiding place in the city and the worst
  explanation for cash. Holding it does not buy a better price on the
  waterfront — the floor's spread is true for anybody standing on it — it buys
  knowing when a boat is in. The Mercer Exchange is one as well: four on
  the floor, ledgers and scales, weights that get condemned. It does not pay
  better — the exchange already pays over the odds and that is a fact about the
  floor — it tells you what the numbers mean. What a good is normally worth
  lives in the core and is on no card, so a trader can see two prices and still
  not know whether this is a week to buy. A floor of yours says how far over or
  under a price is and which way the drift is pulling, about the whole market
  rather than the room you are in, which is the point of holding the sell side.
  Saint Agnes is the last of the four and is one
  too: four behind the bar, the cellar and the glasses, a cellar that floods.
  What a public house is that nothing else in the city is, is somewhere
  everything gets said out loud — so holding one is a line into the city worth
  the same as a telephone in the hall, and reach is what decides whether
  anybody warns you that one side came to a sitdown to finish it. **All four of
  the rooms a family holds are businesses now. Fifteen trades, every one of
  them reaching, all measured by `TestEveryTradeReachesPastItsOwnIncome`.** Every address that earns can be bought and run. A place with
  no price does not change hands. Twenty-six addresses, twelve kinds; the
  pawnbroker is where what is taken off the street turns into money and where
  somebody short pawns the suit off their back. Its window holds what the city
  could not redeem, at 60% of new less wear; hold the counter and you take that
  stock at what it lent instead.

---

## Fault shapes that keep biting

- **A test that skips rather than fails.** A door test looked for one of the
  player's own people to post, found none, and skipped — reporting a pass in
  every summary while proving nothing, and passing with the rule it guarded
  removed. Twenty-odd tests here stand down like that: "no lieutenant in this
  world", "this city has no forecourt", "no breakaway for this seed".
  `TestTheCityStillMakesWhatItsTestsLookFor` asserts those preconditions
  directly, so a change to the city that empties them is one loud failure rather
  than a dozen silent passes.
- **A new rule that punishes the ordinary case.** Twice in two ticks: "below
  twenty trust" is every employee in the game, because everybody here starts at
  nothing and thinks nothing of a stranger; and paying exactly the going rate
  was worth half of the floor's temptation. Both emptied a laundry over a month
  with nothing having happened. `TestAPlayerWhoDoesNothingLosesNothing` is
  written for the city rather than for a rule: a player who owns two businesses,
  pays the rate and wrongs nobody still has them and their people 45 days later.
  Both mistakes, re-made, fail it.

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
- **Two cards in one room under one name.** A command is matched to a card by
  its id and the match takes the last one that fits, so two cards with the same
  id means the wrong card's price, minutes and label are used whatever was
  pressed. It put somebody who asked for the machines into a hand of cards, and
  asking two people where two marks were was offered twice as `about:<who>`.
  `TestNoTwoCardsInARoomShareAName` rules it out everywhere: 615 cards, each
  named once in its room. An id can carry two names — `arms:weapon:3`,
  `sit:back`, `about:leo:vittorio` — and guards that parse the subject out of an
  id have to expect that.
- **A branch in a switch that can never be reached.** A cab's line was written
  after the bare "on foot" case and being warned matched both, so the first one
  won for ever. Caught by reading the code back, not by a test. When a case is
  added to a switch whose earlier cases are broader, it goes above them.
- **Something named by a person where a role was meant.** Three faults in one
  night: the second job crashed the request when the fixer was dead, ten strings
  went on buying coffee for a woman the player had killed, and an office stayed
  empty for the rest of the campaign because the check asked whether the man
  existed rather than whether the desk was filled — and a dead man exists.
  `TestTheCityOutlivesEverybodyItWasWrittenWith` buries every seeded holder and
  runs a fortnight; `TestNothingTheCityShowsNamesTheFirstHolderByName` reads
  every card in the city and fails on any that names one of them.
- **Teeth with no page.** A rule with real consequences that the player has no
  way to see coming. The wage rules that emptied a counter, the gun in somebody's
  coat, the man about to walk out with one of your businesses, and how a
  business is being run. Two sweeps in the log compare every field on a business
  and on a person against what the room and the card actually say; the answer is
  judgement, but the diff is mechanical and worth re-running after anything is
  added.
- **A measurement that reads the wrong thing and says nothing is there.** The
  police looked as though they never came, because the probe counted log lines by
  titles that do not exist — they say "They came to the door" and "Turned over
  at", never "raid". Before believing a system does nothing, check the
  measurement can see it doing something.
- **A guard whose sample cannot show the fault.** "Against your a pair" was
  guarded with a hand that has no article, so the sentence read the same either
  way and putting the fault back left it passing. Pick the sample that can
  fail.

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

**`mise run quick` (~25s) while you are working.** The core, API, sim and store
suites under `-short`, plus gofmt and tsc. The balance measurements skip: they
are 54 files running tens of thousands of simulated days and they answer "did I
move the balance" rather than "did I break something".

Every test in the project calls `t.Parallel()`. One did before this, so a
thousand tests summing 62 seconds ran one after another and took 69. The
exceptions are eighteen director and speech tests that reach for an environment
variable, which is process-wide: Go's own `t.Setenv` refuses to be called in a
parallel test for exactly that reason.

**`mise run gate` (~63s) before committing.** All of it, balance included: the
core suite and the API suites start first and the format check, vet, prettier,
the build and the 56 node tests run while they go.

**`mise run simulate` (~85s) only when the tick could have moved the balance.**
It prints the baseline itself now.

The campaigns run at once. Every one builds its own world and draws from its own
streams; nothing in the core writes package-level state after `init`, there is
no global randomness, and no clock or environment read anywhere in the campaign
path — so sequence was buying nothing but the order the reports landed in, and
an index buys that more cheaply. Eight hundred campaigns went from 414 seconds
to 81 on sixteen cores, with the output identical to the byte. `-workers 1` runs
them one after another, which is how that was checked.

This file said ninety seconds for a long time and that was wrong by a factor of
four before the change and is about right again after it. Do not reach for
`-runs 25` any more: a quarter-size baseline saves forty seconds now and hides
small moves in noisier medians.

Running it inside the gate still costs more than running it after, because three
saturating jobs at once made the whole gate 4m49 against 1m58 — and the
campaigns now saturate the machine on their own.

It also runs **six long campaigns** under `a_long_campaign`, because every
strategy in the baseline ends at the command limit after about six days and the
rules about running a business are about weeks. Across eight hundred standard
campaigns there were 1,446 acquisitions and zero restocks: nothing about stock,
missed payroll or somebody walking out was visible to any number here. The long
ones run a publican for twelve hundred commands, about ninety days, and report
what they actually did — and what they still never do, by name. Today that is
`remedy`, `bankroll` and `poach`, which no policy in this harness has ever
tried.

One tick is: `quick` while iterating, `gate` once, `simulate` only if the
balance could have moved, then the live game and the commit.

Balance baseline, seven strategies
(worker/investor/defiant/reckless/thief/smuggler/racketeer/publican):
deaths 0/0/0/82/75/0/33/0, median cash 12585/2407/5817/90/1090/3564/1658/1732.
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
- **Ask of every change: is the new thing strictly better than the old?** A gun
  bought for one of your own carried no risk at all for an hour, because the
  search took what was in the player's coat and left what was in theirs. This
  file keeps asking that question of the game and it has to be asked of the
  change too, in the tick that makes it.
- Restart the live game on 8791 every tick, backing the save up first, and curl
  it afterwards: build to `.runtime/blackledger-ops.new`, `pkill -f
  blackledger-ops`, `mv` into place, `nohup` it. **Then `grep -c panic
  .runtime/ops.log`.** A 200 from `/api/state` says nothing about the command
  path: a nil dereference in `OfferIfReady` panicked every command the browser
  sent for an unknown number of ticks while the state endpoint answered
  perfectly, and the loop's own check reported the game healthy each time.
- `w.Pay(a.Cost)` runs for every action: one that pays its own fee declares
  `Cost: 0` and uses `asks(...)`; one that pays out uses `add(...)`.
- Two RNG streams — set both `w.RNG` and `w.WorldRNG` when comparing runs.
- Role names must not assume a gender.
- Adding an address: `core/locations.json`, `PlaceIncome`, a trade in
  `core/operations.go` if it is a new kind, `mannerByPlace`, `streetTrades`,
  prompts in `tools/exteriors.py` and `tools/interiors.py` then run both, and a
  free block on the map (`BLOCK=3.6`, `CELL=48` in `src/iso.ts`).
