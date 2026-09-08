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

Focus heavily on building out and hydrating the interface, cleaning it up and
making finding stuff way more intuitive and user friendly. We want to be able
to find all of the actions easily and understand how everything links together.

We also want some effort put into styling it appropriately and fleshing out
all of the pictures of the people etc. Maybe we can use a local image gen model
to flesh out people's portraits. We want to work on stuff like the newspaper to
make it look like an actual 1950s or whatnot newspaper that suits the time and
maybe hydrate the text on it a bit more to make it sound like a real news headline
with some subtext.

I expect you to be able to intuit what the right thing to do here is for the user
experience and makes the game better to play as the interface that we use is crucial.

I would also love to start seeing some of the in-game effects simulated on screen,
like gang hits, assassinations, explosions etc displaying in the simulation when
these events occur would be the ultimate goal and something we should start working
toward if we can.

Maybe we can hook into and utilize another model if we want image generation that
you cannot do or something like that to flesh out the city. The ultimate goal
would be to have the entire city living, with all of the living NPCs being displayed
on screen doing the actions that they are actually doing and when a major event 
occurs, the camera is taken there and we see it played out with a news headline 
appearing just after it plays out.

It's also ok if we want to start using another model, local or online to generate
images as long as they fit our asthetic. If we don't have a local model you can
download one and experiment with using it to generate images as necessary to flesh
out the UI.

As we did with the entering a business thing and cleaning up the view there I want 
you to focus on user friendliness issues like this and really clean it up.

You have to have the mindset of a player in terms of trying to make it easy to understand
and not make any one interface too busy and overwhelming while still having all that
complexity built into it.

Remember that a big goal here is to flesh out the city and have it fully simulated
and animated. We want to see a scene of someone sabotaging a building with an explosion
or committing arson (when we add that mechanic) with fire etc and we want sounds to go
with it. We also want to have all the buildings represented in the city at any one time.
For now we can try our best to create the buildings with our image gen but if it's not
ultra up to scratch that's ok as long as we can regen them to have more detail. In order
to ensure that they are regeneratable we probably need to account for accurately regenerating
and keeping the same size. Maybe reskinning more so than regenerating so the new skin can
add higher textured detail if that makes sense.

We want that map view of all the buildings and city to actually show the living NPCs and
their actions throughout the day as well. That's the ultimate goal, so we should be working
toward that.

- When playing theater the events should be somewhat gruesome and bloody and intense, this is an adult game
- City should have ambience, some effects like smoke and lighting etc to make it feel alive
- Interiors of buildings should have the same treatment, lots of ambiance

## Inbox — unsorted ideas

This section belongs to the user. Add anything here in any form: a mechanic, a
scenario you want to be possible, a moment you want to be surprised by, a
complaint about how something plays. Nothing here is a commitment or a promise
of order.

Claude folds these into the layers above as they are designed, and moves them
into `docs/DEVELOPMENT.md` once they are actually built and tested. An entry
stays here, in the original wording, until it is genuinely implemented — so if
it is still in this list, it does not exist in the game yet.


Ultiamtely the expectation is that I don't have to think of all of these concepts, you, the agent will think of these too and build them out and playtest and simulate them yourself.