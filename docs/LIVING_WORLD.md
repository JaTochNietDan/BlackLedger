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

## Inbox — open

Ok this is much better. It may be worth adding a couple more casino games. It
could be a reason to have other casinos, to have other casino games like slot
machines and whatnot, any other common games you can think of adding.

— slot machines are built (56de529); another game or two is not

---

When setting up the funding on a casino that you own, you should be able to set
that to an actual number by entering it, not having to use a pre-set amount. We
definitely want more level of detail with business management like that. It
should be as dynamic and user settable as possible. Like deposit or withdraw
money.

— not started

---

When gambling you should be able to actually choose how much to gamble, not use
set amounts, up to a max limit. The max limit should be defined by the casino
owner dynamically, whether by you the owner by or by someone else who owns it.

— not started

---

We should also add ambient sounds and sounds to the slot machines and whatnot. I
also want you to flesh out the slot machine a lot more, make it much nicer like
you did for blackjack and roulette. Right now it looks scraggy.

— not started

---

Can people only make attempts on your life while you're at home? They always
seem to hit my home when I'm not there and they are coming after me.

— answered: only at home, by the rule in w.Attack. The rule itself is not fixed yet

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
  are on the forecourt now. Car pictures, what a car is worth in speed, armour at a garage,
  and cars for your own people are all still to do.

---

It seems like you can rob places or take from people's cars multiple times in a row, that should
probably be tracked and time limited etc. Or in the case of the car - until that person repairs
their car or gets a new car, which is a dynamic living NPC thing they could do when their car is
in a damaged state.

— not started

## Inbox — answered

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
