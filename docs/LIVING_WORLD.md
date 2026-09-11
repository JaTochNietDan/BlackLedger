# The living world

The city has to be worth losing. That means it must keep moving when the player
is not looking, and the things that happen to it must be things the player could
also have done. This document is the working plan for that; it records intent and
order, not completed work. `docs/DEVELOPMENT.md` records what is actually built
and what evidence exists for it.

## The principle everything else follows from

Every actor in the city runs on the same rules. A family takes a rival's
warehouse the same way the player takes one, loses standing for it the same way,
and can be answered the same way. Nothing is special-cased for the player, and
nothing is special-cased for the story. When the director narrates a betrayal, it
is describing a transaction the simulation actually performed.

This is what makes surprise trustworthy. A scenario nobody designed can still be
coherent, because it is assembled from committed state rather than invented at
the moment of telling.

## Layers

**1. Actors and holdings.** Organizations own property, collect from it, repair
it, and lose strength when it is damaged. *Built:* families hold income-earning
property; power tracks holding condition and recovers toward a peak; the player
can attack a rival holding and is answered for it.

**2. Faction agency.** Families act on their own schedule against each other, not
only against the player: pressure, sabotage, seizure of holdings, truces. The
player may be a bystander, a beneficiary, or collateral. Requires removing the
assumption that there are exactly two families and that the player is party to
every conflict.

**3. Dynamic organizations.** Factions are created and destroyed during play. A
lieutenant whose family is weak and whose standing is poor splits off with a
holding and becomes a new organization. A family reduced to nothing is absorbed
or disappears, and its property changes hands. New names, new leaders, new
grudges, none of them authored in advance.

**4. Territory and war.** Conflicts have a state: cold, feuding, at war,
settled. War changes what actors do — more seizures, more retaliation, weaker
defence of outlying holdings — and creates the openings a player can exploit.
Wars end in a settlement, a conquest, or exhaustion.

**5. Enterprise.** Businesses are run, not merely owned. Operating choices trade
profit against heat, condition, and attention: skim harder for cash and risk, run
clean for stability, invest for growth. A badly run business loses money. Some
premises support their own activities — a casino can be played as well as owned.

**6. Contraband.** Goods with a price that moves, sources, routes, and buyers.
Profit comes from risk taken knowingly: seizure, informants, a rival who wants
the route. This is the main high-variance income path and a common cause of war.

**7. Succession.** Deferred. Death currently ends a run completely, which is the
intended hardcore default. Ranks within an organization and a named successor are
a later expansion, and must not soften the base game: a coordinated killing of an
entire organization still ends everything.

## What the director does with this

The director is given the world's live conflict state — who is at war, who lost a
holding this week, who is newly weak, who just defected — as ordinary generation
context. It proposes situations grounded in that state and speaks through
characters who have a real position in it. It never invents a war, a betrayal or
a transfer of property; it narrates ones the simulation committed, and it
proposes actions the rules will adjudicate.

The existing validation applies unchanged: the model may not invent terms,
relationships, ranks, or attributions that saved state does not support. As the
world gains more real structure, those guards get *more* to check against, not less.

## How this gets verified

Systems land with rules and tests first, then are exercised headlessly at scale
before they are trusted.

- `go test -race ./...` for rules and invariants.
- `cmd/simulate` for many campaigns against the rules in process, reporting
  milestones, deaths and errors per strategy. Living-world systems need their own
  measures here: wars started, holdings changed hands, factions created and
  destroyed, and how often a run is affected by a conflict it had no part in.
- `cmd/apicheck` for the HTTP stack, gating, idempotency and stale revisions.
- Targeted play through the API for anything that needs a human read.

A system is not finished because it runs. It is finished when a long simulation
shows it producing varied, non-degenerate outcomes, and when the results still
make sense read back as a story.

## Interface

- Intuitive for player to find stuff
- Avoid overcrowding interface with tons of information
- Avoid scrolling if possible by focusing on better design principles
- Don't allow main viewport to scroll, only sub areas like the actions bar when necessary
- Use icons where possible to save text space and have more recognizable icons for common actions that you can start to recognize efficiently

## What the user sends, and how to tell it apart

This section belongs to the user. Add anything here in any form: a mechanic, a
scenario you want to be possible, a moment you want to be surprised by, a
complaint about how something plays. Nothing here is a commitment or a promise
of order.

Three different things arrive here and they are not the same kind of thing:

**Standing instructions** are about how the work is done rather than what to
build. They are never finished, they never move to `docs/DEVELOPMENT.md`, and
nothing should ever be ticked off against them.

**Open** entries are ideas, requests and complaints the game does not answer
yet. They stay in the user's exact words until it does.

**Answered** entries have been built or answered. The account of what was built
lives in `docs/DEVELOPMENT.md`; what stays here is the original wording and a
line saying where it went.

Every line Claude writes in these sections begins with an em dash. Anything
without one is the user's own words and is never edited, summarised or split.
One rule to an entry: each message sits between horizontal rules, so where one
ends and the next begins is never a matter of reading carefully.

## Standing instructions

Ultiamtely the expectation is that I don't have to think of all of these concepts, you, the agent will think of these too and build them out and playtest and simulate them yourself.

— This is the loop itself. Nothing finishes it.

---

Remember we can also have multiple casinos, and multiples of businesses. We can
also do strip clubs, as is very typical of mafia life. Keep fleshing out
businesses and expanding them and adding more and how they interact with the
city simulation and remember we want to really have a lot of characters living
in this city

— Standing: there is always another business, and always more people to put in
  the city. Burlesque, a poolhall, a restaurant, a cab company, a haulier, a
  butcher, two garages, a scrapyard, a forecourt and two filling stations exist
  because of this line, and it is still not finished.

---

Ensure efficiency of development loops by increasing efficiency of your workflow in any way that you can accomplish.

— Standing: how the work is done, so it never finishes. What it has changed so
  far: the ~140s core suite and the ~4m simulation run in the background while
  the next piece of work goes on, rather than the tick sitting and watching
  them; the simulation is parsed rather than grepped; and a slice is measured by
  breaking it once rather than by running the whole gate twice.

## Inbox — open

Send word ahead to a family, so a meeting can be arranged rather than stumbled
into.

— built, and the thing that stopped it last time was the instrument.

  The feature: a seat across from a family used to be offered wherever somebody
  who could speak for them happened to be standing that minute, so the only way
  to get one was to walk the city until a lieutenant turned up in a room — the
  diplomat spends seven hundred and ninety-one journeys buying twenty-six seats.
  People had telephones in 1930 and used them for exactly this. Thirty minutes
  and $25 at the place you live, a telephone or two people who know your name,
  and somebody who can agree to something waits at their own hall for eight
  hours. A family that will not take your call says so.

  Last time this read as coverage falling from 97 kinds of action to 77, and
  four innocent parts of the change were bisected one at a time looking for it.
  It was none of them. Of the thirteen policies the harness runs, the exploring
  one plays about seventy kinds on its own and the other twelve add a handful
  each — and the explorer dies early in most seeds. At four seeds a policy the
  whole reading hung on whether one magpie campaign survived its seven hundred
  commands: seed three alive plays ninety-one kinds, seed three dead plays
  twenty-seven. Pressing one new card once moved that campaign's turn order and
  killed it, and the pin reported twenty kinds lost by a change that took
  nothing away. Sixteen seeds of the explorer rather than four; the union stops
  moving by about the twelfth, no single campaign is worth more than a kind or
  two, and it costs thirty seconds. The figure with the telephone in is 105.

  Two real faults came out of it. The command branch never advanced the clock,
  so a card saying thirty minutes took none. And the price sweep needed a few
  dollars of slack — the control and the run are two campaigns, a day's income
  can land in one and not the other, and a $25 card read as $52 that way.

  The twelfth time a measurement was the instrument rather than the game.

---

Teach the diplomat to telephone. It is the policy the feature was written for
and it does not use it.

— done, and it took two corrections to the feature itself.

  The first was a crash. A seat was being conjured into a room nobody was
  standing in: the audience card asks who is across the table from you, and
  with word sent there was a family expecting you and nobody of theirs there to
  be it. Nil, every time the diplomat pressed the card. Nobody is conjured now
  — whoever can agree to something walks to their own hall on their own feet,
  takes the minutes it takes and can be seen going, and the seat opens when
  they arrive. They have to be set off at the moment the call is made, because
  the city reconsiders where people should be twice a day and a message sent at
  nine in the morning is about this morning.

  The second was the room. The card was offered only where the player lives,
  and the view only carries a room's cards when the player is standing in it,
  so no policy could see it from the branch that wanted it — the diplomat used
  it twenty-three times in eight campaigns. Asking around about a family and
  reaching an understanding with one are both already offered wherever the
  player stands; this belongs with them. A telephone at home is what makes it
  cheap and certain, and without one it is a runner, and neither is a thing you
  do at a counter.

  Eight campaigns, before and after:

      travel    3127 -> 2392   (56% of commands -> 43%)
      audience   309 -> 620
      word         0 -> 126

  Coverage 105 -> 107.

---

Most of what the harness measures is walking.

— answered as far as it goes. The soldier is down from **99% of its commands to
  12%**: `at(place, card)` walks to an address on the strength of the address
  alone, so a branch naming a card that room is not offering walks over, finds
  it refused, falls through to the branch that sends it back, and does that all
  campaign. Its crew branch was asking the bar for a recruit the bar was not
  offering.

  The respectable and the diplomat stay at 65%, and that is the cost of
  diplomacy rather than a fault. Both do their business at one address — the
  paper, a seat at a table — and both of those are refused everywhere else.
  Asking whether the card is ready before setting out means never setting out:
  the paper stopped selling this policy anything at all, and the diplomat took
  nought seats in three campaigns where it had taken twenty-six. Both guards
  were written, measured and taken back out in the same tick.

  The racketeer's 55% is its trade route: buy at the docks, carry, sell at the
  market, walk back. Half of that is walking by construction.

---

A name has no ceiling, and the highest thing it explicitly gates is forty.

— answered, and it was not the ceiling that was wrong. A shift on the pier paid
  a point of standing, every time, for ever. Four hundred shifts at seventy-five
  dollars is a job rather than a reputation, and it is where three hundred and
  thirty-eight came from. It carries a docker to fifteen now — past the six a
  premises takes and the eight it takes to go to work for a family, nowhere near
  the forty for a chair — and after that the pier pays money and nothing else.

  Every policy now ends inside the range the game's own gates use: the worker on
  15, the magpie on 35, the diplomat on 46, the soldier on 62, the publican on
  64. Cash did not move. What moved was a number the panel was showing the
  player as though it were a score.

  A commission that asked the player to "be worth thirty more than you are
  today" went with it — an ask that moves up every time anybody reaches it. It
  is the family's own strength plus ten now, which is real early and gettable
  late.

---

Spread the campaign number across the world's stream inside `New`, and re-measure
what moves.

— done, over three nights. The stream is a plain linear congruential generator
  and it used to start from the campaign's own number, so seeds 1 to 200 opened
  on 0.236 through 0.313 — an eighth of the range — and a hundred and twenty-five
  loops in the test suite start at one. One multiply scatters all of them.

  It turned thirteen tests red and every one was a guard that was wrong rather
  than a number that needed nudging: a poker hand measured in the wrong pocket,
  a slot payout read as money when it is odds, security decided by one roll, a
  raid asked to empty a cellar it cannot find, a refusal asked to draw a
  three-in-four every time, a takeover assumed to have landed, a mayor judged on
  a two per cent gap in one campaign, a country run that varied the world's
  stream while the risk was drawn from the player's, and a door that was
  guarding nothing at all.

  The city itself was quieter than it really is. Left alone for a long campaign
  it now sees three organizations fall where none ever did, seventeen form
  against five, nineteen wars start against seven, and thirty-five people killed
  against twenty-eight. Every measure this project has ever taken of a city on
  its own was taken through that eighth.

---

Asking somebody where to find somebody else publishes a shorter id over HTTP
than the core builds, and refuses about once in a thousand commands.

— answered, and it was never the game. `cmd/blackledger` had no flags at all.
  Every command written down in this project passed `-addr` and `-db` to it —
  the tick ritual, the API playtest's task, the comment at the top of
  `cmd/apicheck` — and Go says nothing about arguments nobody reads, so all of
  it was decoration. A server told to open a scratch database on port 8862
  opened the live campaign on 8791, found something already there, and exited;
  a server somebody had started by hand an hour before went on answering every
  run, on a build from an hour before that. The short id and the refusal were
  both a stale binary talking.

  The flags are real now, with the environment variables as their defaults, and
  the playtest task checks that the server it started is still alive before it
  talks to anything. Run properly: the ids are `about:mara:elena` and the like,
  seventy-one kinds of command instead of forty-nine, and three runs of two
  hundred and fifty commands with no invariant failures.

---

You should be able to buy any car at any time instead of having to go through an upgrade process.

Also adding armor to your car should be different than just buying a car.

---

Nobody has taken over the Bellweather after I killed them

— answered. An office is not a person. The city hall's offices were seeded once
  and asked whether the person existed rather than whether the desk was filled,
  and a dead man exists — so the editor's chair stayed empty for the rest of the
  campaign. Somebody else is behind it within a day now, the newspaper reports
  the change, and whatever you had arranged with the last one lapses, because an
  understanding is with a person and not with a building. Checked against a copy
  of your own save: the mayor and the editor were both empty and both filled on
  the next day.

---

When buying guns/armor and whatnot I don't think you should have to progress through them, you should be able to buy any of them at any time, you don't need to go through some sort of upgrade cycle

— answered. Everything is on the counter at the docks now, each at its own
  price, and a Thompson can be the first thing you buy. The one you are carrying
  says so, and one worse than it is refused with a reason rather than going
  quiet: nothing in this city rewards carrying less gun, so buying down would be
  a trap rather than a choice.

---

When playing slots, it should scroll through the items before displaying the final result. Right now they just shake but the result is shown already. You should really try to animate it smooth and nicely.

— answered, and the second half of it was still wrong until this tick. The
  drums stopped shaking a while ago: the column travels now and the face it
  lands on is decided before a pixel moves. What it travelled through was one
  symbol repeated twenty times. The run behind the window was an arithmetic
  sequence with a stride of seven, and this machine carries seven faces, so
  every step landed on the same one. The drum moved and nothing scrolled past.
  The guard that existed asked whether the three drums differed from each other,
  which they did, and never asked whether one drum differed from itself.

  The run is the strip now, walked backwards in the order it is painted on the
  drum, started a face apart on each one so the three are out of step. Every
  face on the strip goes by. Guarded at the machine's own seven and at every
  strip length from two to twenty-four, because the fault only appears when the
  stride and the length share a factor. What is still not checked by anything
  but eyes is how the travel looks — speed, easing, and whether it reads as a
  machine rather than a list going past.

---

"Nobody has told you where to find them" on the people screen messes up the panel, it pushes the other text to the right.

— answered. The card is a two-column grid: the portrait, then everything about
  them. The sentence was landing in the portrait's column, and that column was
  `auto`, which in a grid means as wide as its widest child. So one line of
  prose in the wrong column made the column as wide as the line and shoved the
  rest of the card's text across.

  The first fix was to tell every block at the foot of the card to span both
  columns, with a guard that reads the card itself and checks each one. That is
  the right rule and it is a rule somebody has to remember every time a block is
  added. So the column has a width of its own now — the small portrait's, which
  is one number the card and the picture both read — and nothing can widen it.
  The span rule stays as the second line rather than the only one. Both are
  guarded, and both guards fail when the change is undone.

  Read out of the stylesheet rather than off a screen: there is no browser on
  this loop. What a fixed column cannot do is a fact about the grid, not about
  the rendering, so this one does not need eyes the way the drums do.

---

What we talked about before, making it so that hidden actions are not hidden anymore, just showed as lower priority in the list (not changing location of sub sections, just putting unavailable actions at the end of the list in each subsection

— answered. There is no toggle in front of them any more. Each subsection is in
  the place it was, and inside it what you can do comes first and what you
  cannot follows, in the same list. A refused card reads as lower priority
  rather than as broken — dimmer, with a quieter border — and the reason keeps
  its colour, because the reason is the point of the card. Three places had the
  fold: the people standing in the room, the room's own sections, and the panel
  inside a building, where the toggle started closed and hid half of what the
  room had.

---

It still says buy Mara a coffee even though now it's Ivo Costa for me since I killed Mara.

— answered, and checked rather than taken on trust. Every line that named a
  role-holder reads the holder now. The check is a sweep: kill whoever holds
  each of the five named roles, let the city fill the desks, then stand in every
  room in turn and look at every card's label, reason and detail, and at the
  guide, for any of the dead names. Nothing says them. The first version of that
  sweep asked each room for its cards without standing in the room, and a room
  you are not in offers two cards — so it found nothing because it looked at
  almost nothing. Standing in the room first, the coffee card reads the living
  fixer's name.

---

When playing poker the game should continue until you stop playing, right now it just requires you to leave the table and rejoin. Realistically it feels like it should be more like actual poker, where you have a buy in and whatnot and you play until people go bust or you can leave.

— answered. There is money on the table now: you put a stake in front of you and
  play out of it, and so does everybody else. The ante is a twentieth of what
  you brought but never more than a quarter of the shortest stack at the table,
  because the people in these rooms carry fifty dollars and a game pitched at
  your bankroll would be a game nobody could sit in. Between hands the table is
  still there: the next hand is one decision, whoever cannot cover the ante
  takes what is left and goes home, and whoever else is in the room takes the
  empty chair. It ends when you have nothing left, when the room has nobody
  left, or when you pick your money up — and getting up or walking out both
  pick it up for you. Measured at the bar: twenty hands on a $1,000 buy-in at
  $12 a hand before the room ran out of people. What is not built: a bigger
  table than four, rebuying after you are cleaned out, and a stake you set
  yourself rather than one the room sets.

---

Now when gambling for some reason you changed it back to "play the nickle machine" which doesn't even show our slots page it just seems to run some background simulation and it doesn't allow you to set your own bet as we fleshed out prior.

We used to have "sit down at the tables" when you were at a casino which is much nicer. No idea why you changed to this weird action button thing again that doesn't even show our fancy interface.

Also a bug that existed before/after that change, merely sitting down to play the games already started playing a game, which is wrong.

This all feels a bit haphazard. Make sure you are keeping track of what you are doing and not regressing. We need to be fixing out UI and making it better.

— part answered: sitting down is a decision the world knows about now, not a
  screen the interface opened by itself. Taking a seat clears the last hand, the
  last spin and the drums, so walking up to a table no longer shows a game
  somebody already played — that was the "already started playing a game" bug,
  and its cause was that the felt is saved state. Getting up is refused in the
  middle of a hand and walking out of the room ends the sitting. "Sit down at
  the tables" is a real action the core offers rather than a button the panel
  drew for itself. The nickel and dollar machines went in 5f5e30d, before this
  message, and the label no longer exists anywhere the game can offer it — a
  stale bundle is the only way to still see it, which is worth knowing. The
  standing half of this, keeping track and not regressing, stays open forever.

— on whether the map is the right way to move around the city: measured before
  answering. Six hundred and fifty journeys between the city's addresses run
  from ten minutes to a hundred and twenty-five over twenty-three distinct
  lengths, so where a place is decides most of what going there costs and the
  map is carrying a real fact rather than decorating one. What was wrong was
  the other way in: the plain list of addresses behind the drawing was
  twenty-six names in no order with nothing on them. It is nearest first now,
  with the journey in minutes on each row, what is yours marked, and what is
  shut to you at the bottom. If the map should still go, it is a bigger
  question than the one this answered — say so and it goes back on the list.

---

## Inbox — answered

When playing the slot machine we should show actual images for the stuff on the rollers. Also it seems to swap the results on the rollers at the end which is odd, they just flip around at random mid-end game. For example it shows 7-7- as it progresses then at the very end it flips to "bell", "lemon", "cherry". Sort that out.

— half built: the flipping is fixed, and it was a real fault rather than an
  animation quirk. The case asked "has everything stopped, or is this drum still
  going?" and showed the result only then — so a drum that had stopped while the
  others were still turning satisfied neither half and fell back to the first
  symbol on the strip, and when the last drum came down all three jumped to the
  real result at once. Every drum shows where it is going to stop from the
  moment the handle goes down; the blur is all the animation does, which is the
  rule this project already had written down about never relying on watching
  motion. The tray still waits for the last drum, because what a pull paid is
  not a thing to announce while they are going. Actual images on the rollers are
  built. Every symbol on the core's strip is painted now — a red seven, a gold
  bar, a bell with its clapper, a plum with a leaf, an orange, a lemon and two
  cherries on one stem — drawn the way the city and the top bar are rather than
  photographed, so they cost nothing to load and match the rest of the game. The
  core still owns which symbol is on the payline and what it pays; this only
  decides what that symbol looks like, and one it has never heard of falls back
  to the word rather than to a blank drum.

---

Playing the game in the back room at the Green Baize is weird. It should again be a separate scene that takes up the screen when you're playing it and you have to leave it rather than right now it just lives in a small box above the action bar. That's silly stuff. We need to stop doing that in future and always dedicate these games to their own screen.

I'd also prefer if this game was the Texas Hold Em version as it's better to play so we can fix that maybe.

— half built: the back room takes the screen. Sitting in on the game is a seat
  the world knows about, exactly as it is at the tables — you go through, the
  city waits, you cannot get up with money in the middle of the table, and
  walking out of the poolhall ends the sitting. The felt is gone from the room's
  panel and the verbs of a hand are gone from its action list. The standing rule
  in this message, that a game always gets its own screen, is written down as a
  guard that fails if either takeover stops being mounted or the room is handed
  a table again.
— and built: the game is Texas hold'em. Two cards each, three on the table, then
  one, then one, with a round of betting on every street — four times the
  decision the draw gave. The best five of the seven a player can see wins, and
  the board being most of everybody's hand is what makes reading the other seats
  the game. Measured over 3,000 hands at a $50 ante: calling everything is
  -$67,112, folding what is beaten is -$33,812, and folding and betting the good
  ones is +$48,338 — a wider gap than the draw's, which is the reason to prefer
  it.

---

You seem to be taking a long time to develop each iteration on the loop. Much longer than your progress before.

Please investigate anything you can do to increase your ability to iterate and develop efficiently. Find bottlenecks and improve them.

— built, in two rounds, and measured both times. First: the balance suites are
  54 files of tens of thousands of simulated days answering "did I move the
  balance", which is not the question while a change is still being written, so
  they skip under `-short` and `mise run quick` was born. Second, and this was
  the whole of it: exactly one test in the project called `t.Parallel()`, so a
  thousand tests summing 62 seconds of work ran one after another and took 69 on
  a machine with cores to spare. Every test has the line now except eighteen
  that reach for an environment variable, which is process-wide. Iterating went
  from 120 seconds to 25, and the full gate from 118 to 63; a tick that used to
  cost 210 seconds of waiting costs 63. `-race` is clean across every package,
  which is the evidence for a mechanical change of that size rather than my
  reading of it.

---

When inside a building it shows this info at the bottom of the action list which is wrong. It'd probably be better to have a more fleshed out display of current building your in with the name and stuff up higher in the fold in a consistent place when you're inside.

The Golden Lily
INDEPENDENT
95% condition

— built, and the cause was worth knowing: the interior is a grid with explicit
  rows for the picture, the people and the work, and the strip naming the room
  had no row of its own. CSS auto-placed it into an implicit row after all three,
  which is why it came out under the action list. It is row one now, above the
  picture, in the same place in every room. It says more than a name, too: who
  holds it, and the figures a player standing in a business wants — condition,
  what it is working at, how many are on the books and what it earns an hour —
  laid out as figures rather than joined with middots at the end of a sentence.

---

When I went after Tila myself, the newspaper info on the attempt showed the wrong portait.

— built. The paper worked out who a story was about by reading its own headline
  back and returning the first person in the city whose name appeared in it, so
  a headline naming two people — which is most of what a paper about people
  going after each other prints — was illustrated with whichever of them the
  city happened to list first. The core knew whose story it was when it filed it
  and threw that away. A story carries its subject now, and everything that
  files one about a person says who: the attempt, the killing, the arrest, the
  obituary and the city hall death. Headlines written before this still get read
  back, because those stories are still in the archive.

---

We should probably show options like "buy kerrigan haulage" before you can afford it instead of having it hidden. We probably should just show all hidden options tbh. Not sure if hiding them is productive.

— built, and the core was already doing its half. Measured: every one of the 17
  businesses for sale is offered to a player with $40 in their pocket, refused
  rather than absent, with the price on the card; and all 251 refusals across
  the city say what would change them, which is what makes a refused card worth
  more than an absent one. What hid them was the panel, which folded refusals
  behind a "Show N you cannot do yet" and started closed. They start open now,
  in the room and beside a person. The toggle stays, because a way to tidy a
  long list is not the same thing as a wall to get past. A home you cannot
  afford is offered the same way, since wanting somewhere better to live is the
  same kind of wanting.

---

Let's verify, can you be attacked while traversing the map? Can your car affect this, whether it's armored, etc? This should all be the case.

— verified, and the answer is yes on every count. It is in numbers now rather
  than in an assurance, because "this should all be the case" is a thing that
  should stay true: an attack that finds you kills on 53% in a bar, 78% in an
  empty room and 88% out on the road. On the road it is 88% on foot, 88% in an
  ordinary car and 62% in a plated one — a car is not cover and a plated one is,
  which is the whole of what paying for the plate buys, and it is worth nothing
  when you are standing in a bar because you are not in it. What you are wearing
  helps in both places. Nobody warns you out there: the people who watch your
  door are at your door, so the street is the one place an attack arrives
  without a moment to decide first. A car also shortens how long you are
  findable — five crossings out of the bar are 205 minutes on foot and 123
  driving — and a wreck or a dry tank is worth neither, because it is not a car
  you are driving. The crossing says all of this before you set off and reads
  differently on foot and in a plated car.

---

While managing the Green Baize I don't see its current funds or how to add to the funds or withdraw from the funds dynamically like we talked about.

— built. The Green Baize is a racket rather than a casino, so none of the float
  machinery reached it — while being the one room in this city that runs a card
  game and charges for the seat. It has a till now: the room panel says what is
  in it, and the room offers putting money in and taking it out in whatever
  figure you type rather than a lot the room decided for you. What the table
  takes for the seat goes into the till rather than straight into your pocket,
  which is what makes it worth looking at — measured over thirty days, $632
  taken and $632 in the till. It is a till and not a float: nothing is covered
  out of it, because the money across that table belongs to the people sitting
  at it.

---

It probably would make sense that we don't always know everyone's location and that's a service that we have to pay for to gain that information from. Maybe only some people know where someone is and we would have to be able to convince us to give them their location, especially if it's a rival family lead.

This would also act as a way of making it harder to make attempts on people's lives in the game as it would be a drawn out task more than anything.

The same should be true of other living NPCs who are working on taking out other people, it should not always be easy to find out where someone is at any time.

They could have knowledge of where you live but not where you currently are, and when you move house they will no longer know where you live until they find out via some contact.

We don't want it to be a situation where we can just talk to one singular character and find out people's locations by paying a fee, that's too simplistic. We want a dynamic knowledge situation and bargainning situation. Some people may not want to give us that information and even asking them could affect our reputation with them. 

I want you to think through this and consider what makes the most sense.

— thought through, and the foundation is built. What was actually easy was not
  the killing: going after somebody, or sending one of your own, already needs
  them standing in front of you. What was easy was *finding* them — the People
  screen listed every living soul with their current address and a button that
  walked you to it, so a name off a list was an address.
— so the player's knowledge is now a thing of its own. You know where somebody
  is because you saw them there, sightings are written down as the clock moves
  through a room you are standing in, and they go stale after a day. The list
  says "Last seen at The Monarch, three hours ago" or "Nobody has told you where
  to find them", and there is no button to walk to. Your own family and the
  people you employ are always findable, because they work for you.
— and the asking is built. There is no broker and no fee: you ask whoever is
  standing in front of you, about somebody you cannot place. Whether they know
  is one thing — the same room, the same family, the same counter, all of them
  facts the city already holds — and whether they say is another. Below 25 trust
  they tell you they have not seen them, the way people say it when they would
  rather not be asked again. Asking somebody to give up their own family takes
  70 trust, and being asked at all is a thing they hold against you whatever
  they answer. What you are told is four hours old, because it is where they
  last saw them and not where they are.
— and the city's own killers play by it too. A family knows your address,
  because it is an address and addresses do not move; where you are standing
  tonight is a different question, and the answer is whether any of their people
  have laid eyes on you in the last day. An attempt on somebody they cannot
  place goes to the house instead. Being somewhere they have not looked is
  cover, being seen is what costs it, and moving house takes back everything
  they had learned. The campaign death rates did not move — the strategies spend
  their lives in public — which is the honest reading: this is cover for a player
  who chooses to use it rather than a general softening.

---

Taking over businesses should be a lot more expensive and high level stuff that you build up to over time.

— built, and measured against what a campaign actually earns rather than
  guessed. A freehold is four times its old listing and every premises already
  held raises the price of the next by nearly half, so the cheapest door in the
  city is $720 against a careful fortnight's $12,500, and the haulage yard goes
  $3,920, then $5,684, then $7,448 as you take the city. The investor still
  holds a casino in 99 runs of a hundred; it just ends them with $2,367 in hand
  rather than $14,129, which is what building up to something costs.
— the measurement found a real fault on the way. A location published its listed
  price while the command charged the real one, so a policy with enough for the
  listing tried, was told "not enough cash", and tried again forever: at four
  times the price six of the seven strategies stalled at $195 and never did
  anything else. The price a place publishes is now the price it costs. Without
  that, this change was impossible rather than merely expensive.

---

We should also add ambient sounds and sounds to the slot machines and whatnot. I
also want you to flesh out the slot machine a lot more, make it much nicer like
you did for blackjack and roulette. Right now it looks scraggy.

— built. The machine was three letters in three boxes, which is a picture of a
  result rather than a machine. It is a cabinet now: a crown, a window with the
  payline across the middle of it, three drums showing three faces each off the
  core's own twenty-stop strip, the handle down the side and a tray at the
  bottom that says what fell into it. And it makes a noise — the handle going
  over, each drum knocking as it stops, and a run of coins into the tray that is
  longer the more it paid. There is a floor tone under the room while you are at
  the tables, a card for a hand and a rattle for the dice, all of it synthesised
  the way the rest of the city's noise is: nothing downloaded, nothing licensed.

---

I thought we added car sales so you could buy new cars that would increase you traversal speed
but when I go to the car dealership I don't see any option to buy cars.

It would also be good if we generated nice car images to show what you're buying and stats
information about the speed of the car relevant to what it gives to you.

You should probably also be able to outfit your car with protection like armor etc at a garage
which will help you survive attacks when traversing out in the streets.

We could probably also extend to be able to provide cars, and armor the cars for people in our
family to keep them more protected from attacks.

— part built: the missing button is fixed. Cars were offered in the garage's own case in
  the room switch and refused there by a rule that asks for a forecourt, and the forecourt
  never offered them at all — a button and its rule asking two different questions. They
  are on the forecourt now.
— and armour is built. A garage plates a car in two stages, fitted to the car
  rather than to you, and the Packard that has always been called armoured now
  actually is. It is worth something only on the street between two addresses,
  which is the one stretch of the city where nothing else protects anybody: of
  600 unwarned hits out there, 71 survived bare, 149 with the doors plated and
  225 with the glass in too. It costs speed, because plate is weight.
— and cars for your own people are built. Buy one of yours a car on the
  forecourt and have a garage plate it, and it is what gets them out when a job
  goes wrong: of 500 jobs that went the other way, 120 walked away from it on
  foot, 175 with a car at the kerb and 231 with a plated one. It is the same
  three endings as before — killed, taken alive with your name coming out of it,
  or out with nothing — and a car moves them toward the last.
— and both of the last two are built. The lot shows a painting of the car it is
  selling (tools/cars.py, the same offline contract as the buildings), and it
  quotes what the car is worth in minutes on a road this player actually walks
  rather than as a percentage nobody walks. The quote includes the weight of any
  plate on it, because a figure that does not agree with what happens is a lie
  told slowly.

---

When inside a building you own the top buttons should probably be for owner management and under a separate subtitle for management actions.

— built, and it turned up something worse. The block was already at the top; it
  had no subtitle saying what it was, and inside a place of yours it now reads
  "Running The Blue Bird" over "staff, stock, repairs and what the house takes".
  But going looking for what belonged in it found that the grouping guard was
  only checking that an action's group was one the interface renders — and the
  fallback group is one it renders. So the house limit, plating a car, fuelling
  it, taking money off your own tables, and every strike and every order to send
  one of your own after somebody were all filed under "Jobs that pay today". A
  hundred and sixty actions, quietly in the wrong place. They are classified now
  and the guard asks the question it meant to ask.

---

We should try to improve the images being displayed on the newspaper. Having a portrait of an affected person or building would be great. Some other black and white dramatization of something would also be great if plausible.

— built: the cut is a picture of the person or the building the story is about.
  The city already had a painted face for everybody in it and a painted front
  for its addresses; the paper screens and inks them so they read as something a
  press ran. The drawn plate is still there as the fallback for a subject with
  no picture — and it is a fallback rather than an underlay, because a
  silhouette is solid black and a photograph laid over one leaves the shape
  showing through whatever the blending mode.

---

When gambling you should be able to actually choose how much to gamble, not use
set amounts, up to a max limit. The max limit should be defined by the casino
owner dynamically, whether by you the owner by or by someone else who owns it.

— built: the stake is typed at every game in the room (5f5e30d), and the holder
  sets the limit. This entry said "not started" for longer than it was true,
  which is my bookkeeping and not the game's. Half of it was also not quite
  finished when I looked: the button to set the limit was offered with no field
  on it, so it sent an amount of nothing and was refused every time. It has a
  field now.

---

The walking between buildings simulation is not that great right now because we don't have the map working properly and the little bar that explains that you're traveling between buildings is at the bottom of the page often below the fold.

— half built: the bar is at the top of the city pane now and says something
  worth reading — where you are going, how long, what you are crossing in, and
  whether anybody is known to be looking for you while you are out in it. That
  matters more than it did: the street is now the one stretch of the city where
  a hit can catch you cold.
— and the map reads now. Fourteen of the twenty-five addresses were drawn as
  two silhouettes — seven identical sheds for the docks, the haulage yard, the
  cab stand, the forecourt, the scrapyard and two filling stations, and seven
  identical shopfronts for the rest — so the city was a grid of copies you could
  not find anything on. The ones that do different things look different: pumps
  under a canopy, a crane over a scrapyard, plate glass on the forecourt, a line
  of cabs, a flatbed in the haulage yard, roller doors on the garages, an awning
  on the butcher, upstairs windows on the poolhall, vent stacks on the steam
  laundry, and a lit marquee on the burlesque, which was being drawn as a house
  because nobody had given its kind a shape.
— still open: whether the map is the right way to move around the city at all.

---

I don't see a lot of cursing from characters in this game, we should increase that since it's with the mafia style. Characters should be able to make threats to you too, I have not seen that yet.

— built: somebody carrying something against you says so to your face, and the
  closer they are to acting on it the less polite the saying gets. Three levels
  of it, shaped by their temperament — a careful man does not make speeches and
  a vain one wants the room to hear. Of 498 lines from people past the weight at
  which people move, 289 had an oath in them.

---

I got some gossip by getting that girl a coffee and then I see "bad blood" red box at the top of the page but it never goes away and it's annoying.

— built: it was a standing state drawn as news. Bad blood fades a point a day,
  so a serious grudge sat at the top of the page for a month and a half. What is
  worth hearing is that two people have just fallen out, so it says so for three
  days and then stops. The grudge itself goes on being true underneath, and the
  person carrying it now tells you about it themselves.

---

The roulette graphics look better but I think that the betting table part should be to the right of the wheel, not below it, like on a real table. You should be able to pick your specific bet amount, up to the maximum (as we talked about in another inbox item, maximum can be set by the casino owner). You can also place multiple bets in roulette, on different numbers, combinations etc, like the real game by putting down chips on each one you want to bet on.

— built, all three. The cloth sits to the right of the wheel. The chip is worth
  whatever you type, up to the house limit the holder sets. And the table takes
  as many chips as you want to put on it: click a spot to lay one down, click
  again to stack, right-click to take one off, and every chip on the cloth is
  settled against the same pocket.

---

Ok this is much better. It may be worth adding a couple more casino games. It
could be a reason to have other casinos, to have other casino games like slot
machines and whatnot, any other common games you can think of adding.

— built: slot machines (56de529) and now craps, which is a third shape of
  decision rather than a third set of numbers. The cards ask you something on
  every card, the wheel asks once and then there is nothing to do, and the dice
  ask once and then make you sit through a run of throws nobody can affect. Pass
  and don't pass both keep about 1.4 in every hundred, measured; the field keeps
  nearly twice that and says so.

---

Vitor Bellendi is always at the kessler filling station for some reason. That seems odd?

I go to The Monarch and request a sit down with the controlling family and it sits me down with the family lead but the lead is not present at this location. That seems unusual. That lead is at Kesseler Filling Station...

— built, and they were the same fault. A lead never had a reason to walk
  anywhere — the evening exempts them, so wherever an errand once left them was
  where they stayed for the rest of the game. A lead holds court at the family's
  own seat now, and goes back to it. The sit-down happens where somebody who can
  settle terms is standing rather than at an address that meant a family by
  name: the club meant Bellandi and the garage meant Russo whatever either was
  doing. If their lead is not there but a lieutenant of theirs is, you sit with
  the lieutenant and the scene says so.

---

Why does it seem like you can send Leo Carver on collections in practically every single building's
action menu?

— built: because it was offered in every one. An order to your own man is work
  that belongs to you rather than to whatever counter you are standing at, and
  the mechanism for saying so already existed. Sending him on collections and
  paying him a bonus both follow you now, filed under the person.

---

Is buying a business called "establish protection"? That's not super clear what that means, we should fix
the wording on that to explain what that actually entails.

— built: it reads "Buy The Green Baize" now, and the description says what you
  are taking on — the positions to keep filled, the wage, the supplies and the
  repairs — rather than only what it earns.

---

When inside a building you can't see who the family that owns it (if any) is anymore.

— built: the room says whose it is at the top, with its condition and how it is
  trading.

---

I don't think "moving against X business yourself" should required respect, that doesn't make sense.

— built: going in yourself asks about your crew, their loyalty and whether they
  are already out, and nothing about your name. Sending them in on your account
  still asks, because that is what a name is for.

---

It seems like you can rob places or take from people's cars multiple times in a row, that should
probably be tracked and time limited etc. Or in the case of the car - until that person repairs
their car or gets a new car, which is a dynamic living NPC thing they could do when their car is
in a damaged state.

— built: a place remembers being robbed and keeps its money somewhere else for
  two days, however the attempt went; a street that lost a car is watched from
  the windows for a night. Both wear off, because a place that can never be
  robbed again has been deleted rather than defended. The car half was already
  true and I had not noticed: taking one apart leaves the owner with none, and
  they buy another off the forecourt when they next walk onto it with the money.
  Mugging was already the same shape — an emptied purse refills on its own.

---

Can people only make attempts on your life while you're at home? They always
seem to hit my home when I'm not there and they are coming after me.

— built: a hit comes to where you are standing. The room decides how much of a
  chance you get — people in it make them careful and give somebody time to see
  them come in, your own people make them slower, and the street between two
  addresses has neither. Two places they cannot walk into: a police cell, and
  anywhere outside the city. A warning is still worth what it always was: they
  go to the house they were told to find you at, and being somewhere else costs
  you the door instead of your life.

---

When setting up the funding on a casino that you own, you should be able to set
that to an actual number by entering it, not having to use a pre-set amount. We
definitely want more level of detail with business management like that. It
should be as dynamic and user settable as possible. Like deposit or withdraw
money.

— built: money into and out of a business now moves in whatever figure is typed.
  Putting money behind the tables, taking it off them, wiring it out of the city
  and bringing it home each carry a number field bounded by what is actually
  there, and the old lot survives only as what the field starts on. The core
  owns the bounds (Action.Sum) so the interface invents none of them.

---

Another idea is that we can have a car dealership that actually acts as the
place you buy your cars from and other living NPCs buy their cars from. We can
have all the city link together to be internally consistent that way. The
business could get more business if cars are destroyed during operations and
whatnot. That's the kind of intertwined and linked up world we need to be
creating and I want you to be continuously autonomously creating

— built: the forecourt, and the city buying its own cars (f888f62, 878a0da, cd97b71).
— The last line of it, "I want you to be continuously autonomously creating", is a standing
  instruction and is not finished by any of that.

---

Somebody who works for you can be reached instead of you (standing
instruction; the inbox is empty).

— built: a family that comes for you and cannot find you used to wreck the
  house, every time, whichever of the four ways you were out of reach — in a
  cell, out of the city, warned and standing elsewhere, or simply not seen
  lately. About a third of the time now they find somebody who works for you
  instead. Everybody else on the books takes ten points off you for it, the
  family's quarrel with you gets worse, and there is a funeral to arrange. Until
  this, signing somebody on carried no risk the city could deliver, and the only
  thing that ever killed one of your own was you walking them into a rival's
  holding yourself.

---

Burying your own (standing instruction; the inbox is empty).

— built: a man who died working for you used to leave a line in the log and
  nothing else — the wage stopped, his name left the books, and everybody still
  on the payroll carried on as though the week had been ordinary. There is a
  funeral now, arranged at a funeral director's within six days of the death. It
  costs $260, or $90 if the parlour is yours because the cars and the box
  already are, and everybody still on your books thinks twelve points better of
  you for it. That is the undertaker's other half: the trade earns on strangers,
  and this is what it costs the player and what it buys.

---

An undertaker (standing instruction: add businesses; the inbox is empty).

— built: Thorne & Sons, a funeral director in the second district. It is the one
  trade in this city whose custom is made entirely by everybody else's work, and
  the counter counts the week's dead the way a garage counts broken glass. A
  casket is the best place in the district for a thing to sit, so holding it
  hides more than a butcher's cold room, and the books explain cash about as
  well as a laundry's because a funeral is paid for in notes by people nobody
  wants to press. Three people work there and there is a way of dying in the
  back yard.

---

A car in pieces and a bench of yours, joined (standing instruction; the inbox is
empty).

— built: what comes off a car taken apart in the street now goes onto the bench
  at a garage of yours, the thinnest one first, rather than only paying cash and
  nudging every garage in the city — rivals' included. The card says which
  garage and how full it is before the crowbar comes out. Somebody with no
  garage loses nothing: the parts still pay, the trade still lifts, and no
  shelf fills by itself. Three cars keeps a bench running on nothing but what
  the city was parking in the street.

---

A still and a bar, joined (taken from the standing instruction rather than the
inbox, which is empty).

— built: a room that sells drink can be restocked out of your own crates instead
  of cash. Both halves had been in the game a long time with no road between
  them: the still made moonshine, the only buyer in this city was the market,
  and the bar was restocked with money. Somebody who owned both was carrying
  crates past his own cellar to sell them to a stranger. A crate on your hands
  is a crate the police can find, so a cellar behind a bar is the one place in
  this city where stock stops being contraband and starts being stock. Four
  trades take it — a pool hall one crate, a saloon and a club three, a revue bar
  five — and a room nobody drinks in is not offered it.

---

Another thing is that we could have the ability for people to steal car parts
and sell them to garages and that's how garages make money, they make more if
there's more car theft or more car repairs to be made from broken windows from
theft etc. Maybe a scrapyard would also be a good idea that links in similarly.
Flesh that out in the loop. Add it to the inbox. Think about more interlinked
stuff too.

— built: parts, the scrapyard, the glass and the garage's trade (878a0da, 0b330cc,
  522cd9b, 832b5f1). "Think about more interlinked stuff too" is standing and stays open.

---

I also think we have a lot of UI QoL fixes to make here, it's weird that for
example you have to go to Mercer Exchange to tell Leo to try to take out Russo.
Stuff like that should really be doable anywhere right? Or commissioning
attacks, stuff like that, seem like they should be doable from any location? It's
weird that you need to go to where Russo is to commission a hit on her. It feels
like a lot of the game is based around where you are for ALL actions which is not
the right way. You need to really fix all that up, make it more logical, improve
the UX around it, not every action should be located in the action pane of these
buildings. You need to think through what the right place for this stuff is and
right context.

— built: work that is not about a room now follows the player, and the audit is written down (730e6b8, e24ab05, 28758c4)

---

We should also focus on cleaning up that massive action bar, I think it'd be
better to display actions below the interior render when inside a building or
something like that. I want you to experiment, think it through and really clean
it up

— built: the room's work sits under the picture in cards of one size (b8a50fd, 5d3fb66)

---

Add a gas station business to the inbox, probably multiple locations. Sells gas
that cars need and it sells other stuff that usual gas stations sell.

— built: two filling stations, and petrol a car actually burns (6d18d9f)

---

Why can't I attempt to take people out? How does that work? I thought we talked
before about fleshing out the ability to either send a family member after
someone to kill them or to attempt to kill them myself, where doing it myself
comes with much greater risk of my own injury or death based on my skills and
equipment. If you send someone of your own then there's a chance they are killed
or captured and then they could be interrogated and give you up as the assailant
or they would know who they are and who they are working on behalf of and send
them to sleep with the fishies anyways.

— built: going after somebody yourself, and sending one of your own
  (core/strike.go). Yours are the best odds you can buy and the only ones that can kill
  you; sending puts their face at the scene, and taken alive they are known to be yours.
