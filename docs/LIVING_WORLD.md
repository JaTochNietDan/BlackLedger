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

## Inbox — unsorted ideas

This section belongs to the user. Add anything here in any form: a mechanic, a
scenario you want to be possible, a moment you want to be surprised by, a
complaint about how something plays. Nothing here is a commitment or a promise
of order.

Claude folds these into the layers above as they are designed, and moves them
into `docs/DEVELOPMENT.md` once they are actually built and tested. An entry
stays here, in the original wording, until it is genuinely implemented — so if
it is still in this list, it does not exist in the game yet.

Remember we can also have multiple casinos, and multiples of businesses. We can
also do strip clubs, as is very typical of mafia life. Keep fleshing out
businesses and expanding them and adding more and how they interact with the
city simulation and remember we want to really have a lot of characters living
in this city

Another idea is that we can have a car dealership that actually acts as the
place you buy your cars from and other living NPCs buy their cars from. We can
have all the city link together to be internally consistent that way. The
business could get more business if cars are destroyed during operations and
whatnot. That's the kind of intertwined and linked up world we need to be
creating and I want you to be continuously autonomously creating

Another thing is that we could have the ability for people to steal car parts
and sell them to garages and that's how garages make money, they make more if
there's more car theft or more car repairs to be made from broken windows from
theft etc. Maybe a scrapyard would also be a good idea that links in similarly.
Flesh that out in the loop. Add it to the inbox. Think about more interlinked
stuff too.

Built: parts and the trade that buys them (core/parts.go), a scrapyard for what
is left of a car, and the last line of it — the glass. A car taken apart in a
street leaves the rest of the row broken rather than gone, a raid that does not
burn a car goes through it instead, and one owner a day pays a garage out of
their own purse to be put right (core/repairs.go). Somebody with nothing keeps
driving it broken, so a poor district gives its garage the same crimes and less
work. "Think about more interlinked stuff too" stays open; it is standing.

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

Ultiamtely the expectation is that I don't have to think of all of these concepts, you, the agent will think of these too and build them out and playtest and simulate them yourself.
