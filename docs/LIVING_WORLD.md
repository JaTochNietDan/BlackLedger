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

---

We should also add ambient sounds and sounds to the slot machines and whatnot. I
also want you to flesh out the slot machine a lot more, make it much nicer like
you did for blackjack and roulette. Right now it looks scraggy.

— not started

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

## Inbox — answered

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
  a hit can catch you cold. The map itself is still open.

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
