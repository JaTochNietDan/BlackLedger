// Package sim exercises the production command boundary using public information only.
package sim

import (
	"blackledger/core"
	"encoding/json"
	"fmt"
)

type Place struct {
	ID        string `json:"id"`
	Owned     bool   `json:"owned"`
	Locked    bool   `json:"locked"`
	Cost      int    `json:"cost"`
	Condition int    `json:"condition"`
	Income    int    `json:"income"`
	// How much longer this place is keeping its money somewhere else after
	// being robbed. The city knows it whether or not anybody is standing in the
	// room, which is what lets a policy walk to a different one instead.
	Shy int `json:"shy"`
	// What is behind the counter: how many positions the trade has, how many
	// are filled, who fills them, what they are paid against the rate, and who
	// has the keys. A policy that only ever buys a place and walks away cannot
	// see any of it, which is why twelve ticks of work on running a business
	// reached no simulated campaign.
	// What it has left to trade on. A business out of stock barely trades, and
	// no policy in this harness ever bought any: measured across eight hundred
	// campaigns, zero restock commands, so every rule about stock, and the
	// price of carrying it, was invisible to the balance.
	Supply    int    `json:"supply"`
	Staff     int    `json:"staff"`
	Positions int    `json:"positions"`
	Wage      int    `json:"wage"`
	Rate      int    `json:"rate"`
	Runs      string `json:"runs"`
	Hands     []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"hands"`
	Actions []core.Action `json:"actions"`
}
type Event struct {
	ID      string        `json:"id"`
	Kind    string        `json:"kind"`
	Choices []core.Choice `json:"choices"`
}

// View deliberately excludes plots, RNG, queued offers and private director memory.
type View struct {
	// The market. A price that moves is the whole of the underground trade, and
	// no policy in here could see one: the view carried premises and people and
	// not the one number the trade is decided on.
	Goods          []core.Good    `json:"goods"`
	BusinessTruces map[string]int `json:"business_truces"`
	Revision       int            `json:"revision"`
	Minute         int            `json:"minute"`
	Player         core.Person    `json:"player"`
	Locations      []Place        `json:"locations"`
	Event          *Event         `json:"event"`
	District       int            `json:"district"`
	// What this campaign has already done, so a policy can prefer what it has
	// not. Only the magpie reads them; they are filled in the run loop rather
	// than by Public, because they are a fact about the run and not about the
	// world.
	Tried map[string]int `json:"-"`
	Stood map[string]int `json:"-"`
}
type Step struct {
	Number  int          `json:"number"`
	Minute  int          `json:"minute"`
	Cash    int          `json:"cash"`
	Health  int          `json:"health"`
	Command core.Command `json:"command"`
	// Collections that came home during this command, and what they paid.
	//
	// Money that arrives late lands on whatever command happens to be running
	// when it does. Sending somebody on a round pays two hours after the
	// decision, so `delegate` read as free and whichever command the crew came
	// back during — resting, mostly — read as generous. A third of what a
	// publican appeared to earn by sitting still was collections landing while
	// it was asleep, and any table built on the trace will say so unless the
	// trace says which money was which.
	Settled int `json:"settled,omitempty"`
	PerTask int `json:"per_task,omitempty"`
}
type Report struct {
	Seed     uint32 `json:"seed"`
	Strategy string `json:"strategy"`
	Director string `json:"director"`
	Commands int    `json:"commands"`
	Minutes  int    `json:"game_minutes"`
	Alive    bool   `json:"alive"`
	Cash     int    `json:"cash"`
	Respect  int    `json:"respect"`
	// What the city took out of them. Attention is the whole risk of the
	// underground trade and the harness could not see it: a report of deaths
	// and cash says nothing about a policy whose money is taken rather than
	// whose life is.
	Heat   int `json:"heat"`
	Health int `json:"health"`
	// Times the goods were taken, by a search or in the street.
	Seizures int `json:"seizures"`
	// Times somebody in the city told the police about them.
	Informed int `json:"informed"`
	// The lowest the player's health ever got. "Ended below full health" turned
	// out to mean almost nothing: a policy that works the docks takes a ten
	// point knock about one shift in seven and ends every campaign bruised,
	// which read as a hundred runs in a hundred "hurt" and put it alongside a
	// policy that was being shot at. How close somebody came to dying is the
	// thing worth reporting.
	Lowest       int            `json:"lowest_health"`
	Milestones   map[string]int `json:"milestone_commands"`
	Actions      map[string]int `json:"action_counts"`
	Events       map[string]int `json:"event_counts"`
	Error        string         `json:"error,omitempty"`
	ReplayQueued int            `json:"replay_queued,omitempty"`
	// What the city did on its own while this run was happening.
	//
	// docs/LIVING_WORLD.md asks for exactly these and they did not exist: a
	// simulation reporting only deaths and cash cannot say whether the living
	// world is alive or whether every campaign plays out the same way. A system
	// is not finished because it runs; it is finished when a long simulation
	// shows it producing varied outcomes.
	World WorldMeasures `json:"world"`
	Trace []Step        `json:"trace,omitempty"`
}

// WorldMeasures counts the things the city did that the player did not do.
type WorldMeasures struct {
	// Conflicts that reached the state "war" during this run.
	WarsStarted int `json:"wars_started"`
	// Of those, the ones the player's organization was no party to.
	WarsElsewhere int `json:"wars_elsewhere"`
	// Premises whose owner changed. Between other organizations only: the
	// player buying a laundry is the player playing, not the city moving.
	HoldingsChangedHands int `json:"holdings_changed_hands"`
	// Organizations that came into existence and that stopped existing. Not
	// counting the player's own: naming your outfit is you playing.
	FactionsCreated   int `json:"factions_created"`
	FactionsDestroyed int `json:"factions_destroyed"`
	// Commands after which the player was worse off — health or a holding —
	// while a war they were no party to was running. This is the closest thing
	// to "affected by a conflict it had no part in" that can be counted
	// honestly: it does not prove the war caused the harm, only that the player
	// was taking damage while somebody else's war was on.
	HurtDuringOthersWar int `json:"hurt_during_others_war"`
}

// watcher remembers enough of the city to notice what changed between two
// commands. Nothing here reads the player's own actions: the question is what
// the city did while the player was doing something else.
type watcher struct {
	wars     map[string]bool
	owners   map[string]string
	factions map[string]bool
	health   int
	// lowest is the low-water mark of the player's health across the whole run.
	lowest   int
	holdings int
}

func watch(w *core.World) *watcher {
	m := &watcher{wars: map[string]bool{}, owners: map[string]string{}, factions: map[string]bool{}}
	m.note(w)
	m.health = w.Player.Health
	m.lowest = w.Player.Health
	m.holdings = owned(w)
	return m
}

func (m *watcher) note(w *core.World) {
	for _, c := range w.Conflicts {
		if c.State == "war" {
			m.wars[c.A+"|"+c.B] = true
		}
	}
	for _, l := range core.Locations {
		if prop := w.Properties[l.ID]; prop != nil {
			m.owners[l.ID] = prop.Owner
		}
	}
	for i := range w.Factions {
		m.factions[w.Factions[i].ID] = true
	}
}

func owned(w *core.World) int {
	n := 0
	for _, l := range core.Locations {
		if w.Own(l.ID) {
			n++
		}
	}
	return n
}

// changed folds one command's worth of city movement into the measures.
func (m *watcher) changed(w *core.World, into *WorldMeasures) {
	player := w.PlayerOrganizationID()
	elsewhere := false
	for _, c := range w.Conflicts {
		if c.State != "war" {
			continue
		}
		key := c.A + "|" + c.B
		if !m.wars[key] {
			into.WarsStarted++
			if c.A != player && c.B != player {
				into.WarsElsewhere++
			}
		}
		if c.A != player && c.B != player {
			elsewhere = true
		}
	}
	for _, l := range core.Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		was, had := m.owners[l.ID]
		// Only movement between other organizations counts. The player taking
		// premises is the player playing.
		if had && was != prop.Owner && was != player && prop.Owner != player {
			into.HoldingsChangedHands++
		}
	}
	// The player naming their own organization is the player playing, not the
	// city making a new family. Counting it made this read as one new
	// organization per campaign for the strategies that form one, and almost
	// none for the strategies that do not — which says something about the
	// strategies and nothing at all about the city.
	live := map[string]bool{}
	for i := range w.Factions {
		id := w.Factions[i].ID
		live[id] = true
		if !m.factions[id] && id != player {
			into.FactionsCreated++
		}
	}
	for id := range m.factions {
		if !live[id] && id != player {
			into.FactionsDestroyed++
		}
	}
	if elsewhere && (w.Player.Health < m.health || owned(w) < m.holdings) {
		into.HurtDuringOthersWar++
	}
	m.wars = map[string]bool{}
	m.owners = map[string]string{}
	m.factions = map[string]bool{}
	m.note(w)
	m.health = w.Player.Health
	if w.Player.Health < m.lowest {
		m.lowest = w.Player.Health
	}
	m.holdings = owned(w)
}

func Public(w *core.World) View {
	b, _ := json.Marshal(w.Public())
	var v View
	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	return v
}
func (v View) place(id string) Place {
	for _, p := range v.Locations {
		if p.ID == id {
			return p
		}
	}
	return Place{}
}
func (v View) action(target, kind string) (core.Command, bool) {
	for _, a := range v.place(target).Actions {
		if a.ID == kind && !a.Disabled {
			return core.Command{Kind: kind, Target: target, Revision: v.Revision}, true
		}
	}
	return core.Command{}, false
}

// The two ends of a route. Fixed facts about this city rather than something to
// discover: what comes ashore comes ashore at the waterfront, and the floor
// where the buyers are is the exchange.
func (v View) cheapEnd(where, good string) bool { return where == "docks" }

func (v View) dearEnd(good string) string {
	if good == "arms" {
		return "docks" // nowhere else deals in them at all
	}
	return "market"
}

func (v View) cheapestFloor(good string) string {
	if good == "cigarettes" {
		return "market" // the only floor that takes them
	}
	return "docks"
}

// most is the largest figure this card will take, which the core works out
// from what the player can carry and what they can pay for. A policy asks the
// card rather than doing that arithmetic again.
func (v View) most(target, kind string) int {
	for _, a := range v.place(target).Actions {
		if a.ID == kind && a.Sum != nil {
			return a.Sum.Most
		}
	}
	return 0
}

func (v View) at(target, kind string) (core.Command, bool) {
	if v.Player.Location != target {
		return v.action(target, "travel")
	}
	return v.action(target, kind)
}

// earn is where the money comes from when the next thing on the ladder costs
// more than the player has.
//
// Every policy here reached for the same thing — carry an envelope at the bar —
// and treated it as inexhaustible. It is a favour rather than a job, and the
// moment there were only so many a day, seven of the eight ended their
// campaigns on the first afternoon with nothing to do.
//
// The pier is the answer and the trick is not to commute to it. Sending them
// back to the bar between every shift spent most of the day walking, which read
// as the economy collapsing and was travel time. Somebody already standing on
// the waterfront stays there.
func (v View) earn() (string, string) {
	switch v.Player.Location {
	case "docks":
		return "docks", "dockwork"
	case "bar":
		if _, ok := v.action("bar", "courier"); ok {
			return "bar", "courier"
		}
		return "docks", "dockwork"
	}
	return "bar", "courier"
}

func Choose(v View, strategy string) (core.Command, error) {
	if v.Event != nil {
		priorities := []string{"approach:careful", "accept", "pay", "escape", "acknowledge", "leave", "decline"}
		if strategy == "defiant" && v.Event.Kind == "business_pressure" {
			priorities = []string{"resist"}
		}
		if strategy == "diplomat" && v.Event.Kind == "audience" {
			priorities = []string{"business_truce", "leave"}
		}
		if strategy == "reckless" {
			priorities = []string{"approach:press", "accept", "resist", "defend", "leave", "decline"}
		}
		for _, id := range priorities {
			for _, c := range v.Event.Choices {
				if c.ID == id && !c.Disabled {
					return core.Command{Kind: "choice", Choice: id, Event: v.Event.ID, Revision: v.Revision}, nil
				}
			}
		}
		for _, c := range v.Event.Choices {
			if !c.Disabled {
				return core.Command{Kind: "choice", Choice: c.ID, Event: v.Event.ID, Revision: v.Revision}, nil
			}
		}
		return core.Command{}, fmt.Errorf("no affordable event choice")
	}
	if strategy == "reckless" {
		// One explicit early challenge, followed by staying home without security.
		// Respect is a public result of that challenge; no hidden schedule is read.
		target, kind := "club", "provoke"
		if v.Player.Respect > 0 {
			target, kind = v.Player.Home, "rest"
		}
		if c, ok := v.at(target, kind); ok {
			return c, nil
		}
	}

	if v.Player.Health < 65 {
		if c, ok := v.at(v.Player.Home, "rest"); ok {
			return c, nil
		}
	}
	// The thief builds the ordinary way first — a name, a crew, somebody on the
	// door — and only then starts taking tills. A policy that robs from the
	// first minute is dead inside six commands and measures nothing.
	// The smuggler. Buys when a good is cheap against its own base and sells
	// when it is dear, which is the only edge there is: prices move toward base
	// with noise and a war disrupts supply. It exists because nothing in this
	// harness had ever traded a crate of anything, so the main high-variance
	// income path in the game was entirely unmeasured.
	if strategy == "smuggler" && v.Player.Respect >= 6 {
		// Attention is what kills this policy, and lying low is the only cure.
		if v.Player.Heat >= 55 {
			if c, ok := v.at(v.Player.Home, "lie_low"); ok {
				return c, nil
			}
		}
		here := v.Player.Location
		trading := here == "market" || here == "docks"
		carrying := 0
		for _, g := range v.Goods {
			carrying += v.Player.Stock[g.ID]
		}
		if trading {
			// Sell whatever is being carried, wherever this floor pays more for
			// it than the one it came from. The floors are not the same price
			// any more, so a policy carries rather than waits.
			for _, g := range v.Goods {
				held := v.Player.Stock[g.ID]
				if held == 0 {
					continue
				}
				if c, ok := v.action(here, "sell:"+g.ID); ok {
					c.Amount = min(held, v.most(here, c.Kind))
					if c.Amount > 0 {
						return c, nil
					}
				}
			}
			// Buy anything cheap, as much as the card will take. How much is
			// asked of the card rather than worked out again: the first version
			// did its own arithmetic and asked for four crates with room for
			// one, and the campaign ended on the refusal.
			if v.Player.Heat < 40 && v.Player.Cash > 600 {
				for _, g := range v.Goods {
					// Only where this floor is the cheap end of a route.
					if !v.cheapEnd(here, g.ID) {
						continue
					}
					c, ok := v.action(here, "buy:"+g.ID)
					if !ok {
						continue
					}
					c.Amount = min(v.most(here, c.Kind), max(1, (v.Player.Cash-500)/max(1, g.Price)))
					if c.Amount > 0 {
						return c, nil
					}
				}
			}
		}
		// Carrying something means walking it to the other end of the route,
		// which is the risk this whole trade is made of.
		if carrying > 0 && v.Player.Heat < 55 {
			for _, g := range v.Goods {
				if v.Player.Stock[g.ID] == 0 {
					continue
				}
				if to := v.dearEnd(g.ID); to != "" && to != here {
					if c, ok := v.at(to, "wait"); ok && c.Kind == "travel" {
						return c, nil
					}
				}
			}
		}
		if !trading && carrying == 0 && v.Player.Cash > 900 && v.Player.Heat < 40 {
			for _, g := range v.Goods {
				if from := v.cheapestFloor(g.ID); from != "" {
					if c, ok := v.at(from, "wait"); ok && c.Kind == "travel" {
						return c, nil
					}
				}
			}
		}
	}
	// violent is the handful of things that get somebody killed. Named here rather
	// than guessed at, because a policy that avoids anything that sounds dangerous
	// would avoid most of the game.
	var violent = map[string]bool{
		"mug": true, "strike": true, "charge": true, "plant": true,
		"sabotage": true, "provoke": true, "takeover": true, "rob": true,
		"rob:crew": true, "dockwork": false,
	}

	// The magpie takes whatever it has taken least.
	//
	// Every other policy here is a person with a plan, and between the eight of
	// them they exercise 78 of the game's 116 kinds of action. Ninety-four kinds
	// were never taken by any of eight hundred campaigns — which does not mean
	// they are unreachable, it means no policy was written to want them, and
	// every balance figure this harness has ever printed was silent about all
	// of them. A window at a pawnbroker, a boat at a pier, a night at a room you
	// host and a meeting between two families had all been added and priced
	// without one campaign ever touching them.
	//
	// So this one has no plan. It looks at everything offered in the room it is
	// standing in, takes whichever it has taken fewest times, and walks
	// somewhere else when the room is exhausted. It plays badly on purpose: the
	// point is coverage, not a score, and its cash column should be read as
	// what happens to somebody who does everything once rather than as a
	// strategy anybody would follow.
	if strategy == "magpie" {
		// Poor first, curious second.
		//
		// The first version of this only ever took whatever it had taken least
		// in the room it was standing in, and it starved: with no income
		// everything with a price on it is refused, so the only cards left
		// offered were the free ones and eight runs reached six kinds of thing
		// — fewer than the eight policies with plans. A policy that means to
		// see the whole game has to be able to afford the whole game.
		// Upkeep before curiosity, and resting is upkeep.
		//
		// Taking only what it has never taken meant every kind of thing
		// happened exactly once, resting included: nine times in twelve hundred
		// commands, and a median health of seventy that never climbed. Which in
		// turn meant it could never be well enough to take an attempt on
		// anybody — it was above ninety for ten steps out of twelve hundred —
		// so eight kinds of action stayed out of reach for a reason that had
		// nothing to do with the game.
		if v.Player.Health < 70 {
			if c, ok := v.at(v.Player.Home, "rest"); ok {
				return c, nil
			}
		}
		// And knowing people is upkeep too, which is the question this was
		// built to ask. Six of eleven deaths came while walking between two
		// addresses — somebody the policy had robbed or hit, coming for it in
		// the street — and the odds on an attack in the core are written as
		// the odds of an *unwarned* one. A coffee costs ten dollars and the
		// city has no other way of telling a stranger you have wronged from a
		// stranger.
		if v.Player.Contacts < 5 {
			if c, ok := v.at("bar", "contact"); ok {
				return c, nil
			}
		}
		if v.Player.Cash < 1500 {
			if c, ok := v.at("docks", "dockwork"); ok {
				return c, nil
			}
		}
		// A room of its own, first, before anything else is interesting.
		//
		// Everything a business can do — hiring, the wage, restocking, putting
		// the trouble right, laundering, a night, the window, the boat — needs
		// the player to hold the place, and none of it had ever been played by
		// any policy in this harness. Left to wander and take whatever it had
		// taken least, this one never bought anything at all in eight campaigns
		// of forty days: it earns at the docks, spends what it earns on the
		// first priced card it walks past, and is poor again by the time it is
		// standing somewhere with a deed for sale. A laundry wants about $1,500
		// in hand and it was carrying four hundred.
		//
		// So the first business is a goal rather than an accident. Work until
		// it can afford one, walk to the cheapest thing for sale, buy it. After
		// that the whole branch is offered wherever it stands and the ordinary
		// rule can have it back.
		holds := 0
		for _, p := range v.Locations {
			if p.Owned {
				holds++
			}
		}
		if holds == 0 {
			if v.Player.Cash < 1500 {
				if c, ok := v.at("docks", "dockwork"); ok {
					return c, nil
				}
			}
			cheapest, price := "", 0
			for _, p := range v.Locations {
				if p.Owned || p.Locked || p.Cost <= 0 || p.Income <= 0 {
					continue
				}
				if cheapest == "" || p.Cost < price {
					cheapest, price = p.ID, p.Cost
				}
			}
			if cheapest != "" {
				if c, ok := v.at(cheapest, "acquire"); ok {
					return c, nil
				}
			}
		}
		// And any other room of its own, whenever one is offered.
		//
		// Everything a business can do — hiring, the wage, restocking, putting
		// the trouble right, laundering, a night, the window, the boat — needs
		// the player to hold the place, and none of it had ever been played by
		// any policy in this harness. Left to "take whatever you have taken
		// least", this one never bought anything at all: it earns at the docks,
		// spends what it earns on the first priced card it meets, and is poor
		// again by the time it is standing somewhere with a deed for sale.
		//
		// So buying is not one card among many. It is the thing that opens the
		// rest of the game, and a policy meant to see the game takes it.
		if c, ok := v.action(v.Player.Location, "acquire"); ok {
			return c, nil
		}
		best, fewest := "", 0
		here := v.place(v.Player.Location)
		for _, a := range here.Actions {
			if a.Disabled || a.ID == "travel" {
				continue
			}
			// Not the ones that end the run or undo the point of it.
			if a.ID == "rest" && v.Player.Health > 92 {
				continue
			}
			// Curious, not suicidal. Taking every attempt on a person the
			// moment it was offered gave this a median life of a day and a
			// half, which explores nothing: it reached nine kinds of action
			// nobody else reached and then died before it could reach a tenth.
			// It will still do all of them — it simply waits until it is in a
			// condition to survive them.
			// Fit enough to survive it, which has to be a figure resting can
			// actually reach. Requiring full health while resting stopped at
			// ninety-two meant the two rules contradicted each other and the
			// policy could never take an attempt on anybody at all: mugging,
			// striking, robbing, charging, planting, sabotage, provoking and
			// moving on a family — eight kinds of action — were locked out by
			// arithmetic rather than by any decision.
			if violent[root(a.ID)] && v.Player.Health < 90 {
				continue
			}
			// By kind, not by person.
			n := 0
			for id, count := range v.Tried {
				if root(id) == root(a.ID) {
					n += count
				}
			}
			if best == "" || n < fewest {
				best, fewest = a.ID, n
			}
		}
		// Only if it is something this campaign has never done. Taking the
		// least-done card instead meant there was always something to take
		// here, so the walking-on rule below never ran once and the policy
		// spent forty days in whichever room it happened to be standing in.
		if best != "" && fewest == 0 {
			if c, ok := v.action(v.Player.Location, best); ok {
				// A card with a field on it wants a number, and sending none
				// is refused. Seven runs in eight ended on exactly that —
				// "a house limit runs from $20 to $5000", "nobody stands
				// behind a counter for less than $3 a day" — which is why the
				// policy stopped at a bit under eight days however many steps
				// it was given. The card states its own range; the smallest
				// figure in it is always a legal answer.
				for _, a := range here.Actions {
					if a.ID == best && a.Sum != nil {
						c.Amount = a.Sum.Least
					}
				}
				return c, nil
			}
		}
		// Nothing left here that it has not done. Somewhere else — and a room
		// of its own first, because hiring, restocking, putting the trouble
		// right and reading the books are only offered to somebody standing in
		// the place they hold. Buying one and never going back reached two of
		// that branch out of a dozen.
		where, stood := "", 0
		for _, p := range v.Locations {
			if p.ID == v.Player.Location || p.Locked || !p.Owned {
				continue
			}
			if n := v.Stood[p.ID]; where == "" || n < stood {
				where, stood = p.ID, n
			}
		}
		// And anywhere at all, when it has stood in all of its own more than
		// it has stood in somewhere it has never been.
		for _, p := range v.Locations {
			if p.ID == v.Player.Location || p.Locked {
				continue
			}
			if n := v.Stood[p.ID]; where == "" || n < stood {
				where, stood = p.ID, n
			}
		}
		if where != "" {
			if c, ok := v.at(where, "travel"); ok {
				return c, nil
			}
		}
	}
	// The publican runs the businesses rather than only buying them. Twelve
	// ticks of work — hiring, the wage, putting somebody in charge, restocking,
	// the people who walk out and the families who come for them — reached no
	// simulated campaign at all, because every policy here buys a place and
	// then never thinks about it again. Every "baseline unmoved: no campaign
	// policy does this" in the development log is this gap.
	if strategy == "publican" {
		// It buys the way the investor does — that ladder is below — and this
		// is what it does with what it has bought.
		if v.Player.Health < 85 {
			if c, ok := v.at(v.Player.Home, "rest"); ok {
				return c, nil
			}
		}
		for _, p := range v.Locations {
			if !p.Owned || p.Income <= 0 {
				continue
			}
			// Somebody in charge of it, so it stocks itself.
			if p.Runs == "" {
				for _, h := range p.Hands {
					if c, ok := v.at(p.ID, "incharge:"+h.ID); ok {
						return c, nil
					}
				}
			}
			// And paid over the rate, so nobody listens to a rival.
			if p.Wage > 0 && p.Wage <= p.Rate {
				if c, ok := v.at(p.ID, "wage"); ok {
					c.Amount = p.Rate + 2
					return c, nil
				}
			}
			if p.Staff < p.Positions {
				if c, ok := v.at(p.ID, "hire"); ok {
					return c, nil
				}
			}
			// And something to trade with. Somebody in charge stocks the place
			// themselves, so this is for the ones nobody is minding.
			if p.Supply <= 0 {
				if c, ok := v.at(p.ID, "restock"); ok {
					return c, nil
				}
			}
		}
	}
	if strategy == "publican" {
		// Everything else the investor does: the crew, the home, the security
		// and the ladder of premises. A publican is an investor who reads the
		// books afterwards.
		strategy = "investor"
	}
	// The racketeer does both, which is the only way to exercise an informant:
	// somebody has to hate them and they have to be doing something worth
	// telling the police about, and no policy here satisfied both halves at
	// once. It trades like the smuggler and takes tills like the thief.
	if strategy == "racketeer" && v.Player.Respect >= 6 && v.Player.Security >= 1 {
		if v.Player.Health < 85 {
			if c, ok := v.at(v.Player.Home, "rest"); ok {
				return c, nil
			}
		}
		here := v.Player.Location
		if here == "market" || here == "docks" {
			for _, g := range v.Goods {
				if held := v.Player.Stock[g.ID]; held > 0 {
					if c, ok := v.action(here, "sell:"+g.ID); ok {
						c.Amount = min(held, v.most(here, c.Kind))
						if c.Amount > 0 {
							return c, nil
						}
					}
				}
			}
			if v.Player.Heat < 40 && v.Player.Cash > 600 && here == "docks" {
				for _, g := range v.Goods {
					c, ok := v.action(here, "buy:"+g.ID)
					if !ok {
						continue
					}
					c.Amount = min(v.most(here, c.Kind), max(1, (v.Player.Cash-500)/max(1, 40)))
					if c.Amount > 0 {
						return c, nil
					}
				}
			}
		}
		carrying := 0
		for _, g := range v.Goods {
			carrying += v.Player.Stock[g.ID]
		}
		if carrying > 0 && here != "market" {
			if c, ok := v.at("market", "wait"); ok && c.Kind == "travel" {
				return c, nil
			}
		}
		// And on the way, take whatever is going. This is what makes enemies.
		for _, l := range v.Locations {
			if l.Owned || l.Locked || l.Income <= 0 || l.Shy > 0 {
				continue
			}
			if c, ok := v.action(here, "rob"); ok && l.ID == here {
				return c, nil
			}
		}
		if c, ok := v.action(here, "mug"); ok {
			return c, nil
		}
		if carrying == 0 && v.Player.Cash > 900 && here != "docks" {
			if c, ok := v.at("docks", "wait"); ok && c.Kind == "travel" {
				return c, nil
			}
		}
	}
	if strategy == "thief" && v.Player.Respect >= 6 && v.Player.Security >= 1 {
		// A robbery that goes wrong is a beating, and two in a row is a death,
		// so this one waits until it is whole before trying another.
		if v.Player.Health < 85 {
			if c, ok := v.at(v.Player.Home, "rest"); ok {
				return c, nil
			}
		}
		// Takes tills, and nothing else while there is one to take. No other
		// policy here ever robs anything, so a rule about robbing the same
		// place twice could be added or deleted and every number in the report
		// would stay where it was.
		for _, l := range v.Locations {
			if l.Owned || l.Locked || l.Income <= 0 || l.Shy > 0 {
				continue
			}
			if c, ok := v.at(l.ID, "rob"); ok {
				return c, nil
			}
		}
	}
	if strategy == "worker" {
		if c, ok := v.at("docks", "dockwork"); ok {
			return c, nil
		}
	}
	if strategy == "diplomat" && v.Player.Cash >= 200 {
		for _, p := range v.Locations {
			if !p.Owned || p.Income <= 0 {
				continue
			}
			actor, venue := "bellandi", "club"
			if p.ID == "garage" || p.ID == "casino" {
				actor, venue = "russo", "garage"
			}
			if v.BusinessTruces[actor] <= v.Minute {
				if c, ok := v.at(venue, "audience"); ok {
					return c, nil
				}
			}
		}
	}
	if len(v.Player.Crew) > 0 && v.Player.Crew[0].Loyalty < 50 {
		if c, ok := v.action(v.Player.Location, "crew_bonus"); ok {
			return c, nil
		}
	}
	if c, ok := v.action(v.Player.Location, "delegate"); ok {
		return c, nil
	}
	target, kind, cost := "bar", "courier", 0
	switch {
	case v.Player.Respect < 6:
	case !v.place("laundry").Owned:
		target, kind, cost = "laundry", "acquire", v.place("laundry").Cost
	case len(v.Player.Crew) == 0:
		target, kind, cost = "bar", "recruit", 90
	case v.Player.Contacts < 2:
		target, kind, cost = "bar", "contact", 10
	case v.District == 0:
		target, kind, cost = "apartment", "expand", 100
	case v.Player.Home == "room":
		target, kind, cost = "apartment", "move_home", 180
	case v.Player.Security < 1:
		target, kind, cost = v.Player.Home, "security", 100
	case !v.place("garage").Owned:
		target, kind, cost = "garage", "acquire", v.place("garage").Cost
	case !v.place("casino").Owned:
		target, kind, cost = "casino", "acquire", v.place("casino").Cost
	default:
		target, kind = v.Player.Home, "rest"
	}
	for _, p := range v.Locations {
		if p.Owned && p.Income > 0 && p.Condition < 70 {
			target, kind, cost = p.ID, "repair", 50
			break
		}
	}
	if v.Player.Cash < cost {
		target, kind = v.earn()
	}
	if kind == "expand" {
		if c, ok := v.action(target, kind); ok {
			return c, nil
		}
	}
	if c, ok := v.at(target, kind); ok {
		return c, nil
	}
	if c, ok := v.at(v.earn()); ok {
		return c, nil
	}
	// Whatever this room offers, rather than nothing at all.
	//
	// One campaign in nine hundred ended with "no policy action at precinct": a
	// policy taken in and held has no ladder to climb and cannot walk to the
	// bar or the pier, so every rule above this had nothing to say and the run
	// stopped. A cell has things to do in it — sit the sentence out, talk, send
	// for a lawyer — and a policy that will not do any of them is a hole in the
	// harness rather than a fact about the game.
	for _, a := range v.place(v.Player.Location).Actions {
		if a.Disabled || a.ID == "travel" {
			continue
		}
		if c, ok := v.action(v.Player.Location, a.ID); ok {
			return c, nil
		}
	}
	return core.Command{}, fmt.Errorf("no policy action at %s", v.Player.Location)
}
func Run(seed uint32, strategy, director string, limit int, trace bool) Report {
	return RunRecorded(seed, strategy, director, limit, trace, nil)
}

// RunRecorded consumes each recorded proposal once; exhaustion falls back to authored play.
func RunRecorded(seed uint32, strategy, director string, limit int, trace bool, corpus []core.Proposal) Report {
	r := Report{Seed: seed, Strategy: strategy, Director: director, Milestones: map[string]int{}, Actions: map[string]int{}, Events: map[string]int{}}
	w := core.New(seed)
	eyes := watch(w)
	start := w.Minute
	nextOffer := start + 240
	cursor := 0
	stood := map[string]int{}
	for i := 0; i < limit && w.Player.Alive; i++ {
		// This is a deterministic test provider, not the real AI director. It uses the same validator.
		if (director == "fixture" || (director == "replay" && cursor < len(corpus))) && w.Minute >= nextOffer && w.Event == nil && len(w.Offers) == 0 {
			operation := w.NextDirectorOperation()
			p := core.Proposal{Title: "Simulation arrangement", Body: "I need help with a discreet neighborhood job.", Speaker: "mara", Operation: operation, Outcome: "The agreed job is complete.", Approaches: []core.Approach{{Method: "careful", Label: "Prepare carefully"}, {Method: "press", Label: "Push the schedule"}}}
			if director == "replay" {
				p = corpus[cursor]
			}
			scene, err := w.ValidateProposal(p)
			if err != nil {
				r.Error = err.Error()
				break
			}
			if director == "replay" {
				cursor++
			}
			w.Offers = append(w.Offers, core.Offer{Ready: w.Minute + 30, Event: scene})
			nextOffer = w.Minute + 240
		}
		v := Public(w)
		v.Tried, v.Stood = r.Actions, stood
		c, err := Choose(v, strategy)
		if err != nil {
			r.Error = err.Error()
			break
		}
		outstanding := len(w.Tasks)
		if trace {
			r.Trace = append(r.Trace, Step{Number: i + 1, Minute: w.Minute,
				Cash: w.Player.Cash, Health: w.Player.Health, Command: c})
		}
		stood[w.Player.Location]++
		n, err := core.Execute(w, c)
		if err != nil {
			r.Error = fmt.Sprintf("command %d %s: %v", i+1, c.Kind, err)
			break
		}
		if n.Minute < w.Minute || n.Revision != w.Revision+1 || n.Player.Cash < 0 || n.Player.Health < 0 || n.Player.Health > 100 {
			r.Error = "command invariant failed"
			break
		}
		if v.Event != nil {
			r.Events[v.Event.Kind]++
		}
		r.Actions[c.Kind]++
		// What came home during this command, so a table built on the trace can
		// put late money against the decision that sent for it rather than
		// against whatever was happening when it arrived.
		if trace && len(r.Trace) > 0 {
			if settled := outstanding - len(n.Tasks); settled > 0 {
				r.Trace[len(r.Trace)-1].Settled = settled
				r.Trace[len(r.Trace)-1].PerTask = core.CollectionPay
			}
		}
		w = n
		eyes.changed(w, &r.World)
		r.Commands++
		marks := map[string]bool{"crew": len(w.Player.Crew) > 0, "laundry": w.Own("laundry"), "garage": w.Own("garage"), "casino": w.Own("casino"), "housing": w.Player.Home != "room", "security": w.Player.Security > 0, "death": !w.Player.Alive}
		for name, done := range marks {
			if done && r.Milestones[name] == 0 {
				r.Milestones[name] = r.Commands
			}
		}
	}
	r.ReplayQueued = cursor
	r.Minutes = w.Minute - start
	r.Cash = w.Player.Cash
	r.Respect = w.Player.Respect
	r.Heat, r.Health = w.Player.Heat, w.Player.Health
	r.Lowest = eyes.lowest
	// What the city took. The record is the only place a seizure is written
	// down, which is right — it is a thing that happened, not a counter.
	for _, entry := range w.History {
		if entry.Title == "The goods are gone" {
			r.Seizures += max(1, entry.Count)
		}
		if entry.Title == "Somebody talked" {
			r.Informed += max(1, entry.Count)
		}
	}
	r.Alive = w.Player.Alive
	return r
}

// root is an action id without whoever it is about.
//
// Half the ids in this game carry a person after a colon — `strike:person-8`,
// `about:leo:vittorio` — and both of the magpie's rules were reading the whole
// id. It cost the policy its life and most of its purpose. The list of things
// that get somebody killed never matched `strike:person-8`, so the one rule
// meant to keep it alive did nothing at all; and "taken fewest times" counted
// each person as a separate thing to try, so a room with twelve people in it
// was twelve untried cards and it never left. Every trace ended the same way:
// dockwork, strike, strike, strike, and health from seventy to twenty-two in
// three commands.
func root(id string) string {
	for i := 0; i < len(id); i++ {
		if id[i] == ':' {
			return id[:i]
		}
	}
	return id
}
