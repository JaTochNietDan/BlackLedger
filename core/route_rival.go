package core

import "fmt"

// Somebody else wants the route.
//
// Layer 6 of docs/LIVING_WORLD.md asks for three risks: "seizure, informants, a
// rival who wants the route". Measured over a hundred campaigns, one existed. A
// policy that ran the trade carried thirty-seven attention and lost the goods
// fourteen times and was never once hurt, because nobody in the city had any
// opinion about a man carrying crates through it.
//
// Now that the waterfront and the exchange disagree about what a crate is
// worth, running between them is a trade — and a trade is a thing worth taking
// off somebody. This is the half of contraband that reaches the faction system,
// which is what that layer means by it being "a common cause of war".

const (
	// RouteNotice is how many units have to cross the city before the word is
	// out. Somebody moving a few crates is nobody; somebody moving forty is a
	// business that other businesses have an opinion about.
	RouteNotice = 30
	// RouteInterest is how often, per half-day, a family that has noticed
	// decides to do something about it.
	RouteInterest = .09
	// RouteCools is how much of the word fades a day. A trade you have stopped
	// running stops being your trade.
	RouteCools = 2
	// RouteWait is how long between somebody deciding and somebody arriving.
	RouteWait = 300
)

// routePlotted reports whether somebody is already coming for the trade.
func (w *World) routePlotted() bool {
	for _, p := range w.Plots {
		if p.Life == w.Life && p.Kind == "route" {
			return true
		}
	}
	return false
}

// RouteRun records a load that crossed the city, which is what gets a man
// noticed. Called from the sale rather than the purchase: buying is a man with
// money, selling is a man with a trade.
func (w *World) RouteRun(units int) {
	if units <= 0 {
		return
	}
	w.Player.Runs = min(200, w.Player.Runs+units)
}

// RouteDay lets the word fade. Nobody is remembered for a trade they gave up.
func (w *World) RouteDay() {
	w.Player.Runs = max(0, w.Player.Runs-RouteCools)
}

// ConsiderRoute is a family deciding that a trade somebody else is running
// would be better run by them. Off the world's own stream: the player is not
// party to the conversation.
func (w *World) ConsiderRoute() {
	if !w.Player.Alive || w.Player.Runs < RouteNotice || w.routePlotted() {
		return
	}
	if w.WorldRandom() >= RouteInterest {
		return
	}
	// Whoever is strongest and is not the player's own organization. A trade is
	// taken by somebody with the people to take it.
	var want *Faction
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() || len(w.Members(f.ID)) == 0 {
			continue
		}
		if want == nil || f.Power > want.Power {
			want = f
		}
	}
	if want == nil {
		return
	}
	w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "route", Life: w.Life,
		Due: w.Minute + RouteWait, Actor: want.ID, Strength: 4})
	want.Goodwill = max(-100, want.Goodwill-10)
	w.Log("Somebody has been asking about the crates",
		fmt.Sprintf("%s %s who has been moving what, and how often. A trade worth running is a trade worth taking.",
			Leads(want.Name), Agree(want.Name, "knows", "know")), "danger")
}

// TakeTheRoute is them arriving. They came for the load, not for the man: if he
// is carrying, they take it and leave him standing, and if he is not, they say
// what they came to say.
func (w *World) TakeTheRoute(plot Plot) {
	actor := w.factionName(plot.Actor)
	carrying := w.Carrying()
	if carrying == 0 {
		w.Player.Runs = max(0, w.Player.Runs-RouteNotice/2)
		w.Log("A word about the trade",
			fmt.Sprintf("Two of %s came to find you carrying and found you empty-handed. The message was that the crates coming off that waterfront are spoken for.", actor), "danger")
		return
	}
	lost := w.Seize(fmt.Sprintf("%s took the load off you in the street.", actor))
	w.Player.Runs = max(0, w.Player.Runs-RouteNotice)
	w.Player.Respect = max(0, w.Player.Respect-2)
	if f := w.faction(plot.Actor); f != nil {
		// They are better off for it, and the city can see who is running that
		// trade now.
		f.Cash += lost * 20
		f.Power = min(peak(f), f.Power+2)
	}
	w.Report("theft", "GOODS TAKEN IN THE STREET",
		w.unattributed(w.whereItHappens(), fmt.Sprintf("A quantity of contraband was taken from somebody on foot near %s. No arrest has been made.", w.whereItHappens())))
	w.Log("They took the load", fmt.Sprintf("%s wanted the trade and now %s it. %s gone, and you walked away, which was the arrangement they had in mind.",
		actor, Agree(actor, "has", "have"), counted(lost, "unit is", "units are")), "danger")
}
