package core

import (
	"fmt"
	"strings"
)

// A business's staff was a number. You hired a pair of hands, the wage bill went
// up, and nobody in Bellwether had a job: the person behind the counter of a
// place the player owned did not exist, could not be talked to, killed, robbed,
// poached or arrested, and the city's own people had no reason to be anywhere in
// the daytime except the one the routine invented for them.
//
// The count stays, because everything that reads it — what a place can handle,
// what the wages are, whether it is short-handed — is right to read a count. It
// is now the length of a list of people who live here. Somebody who dies, or is
// taken in and held, is not behind that counter, and the position is empty until
// it is filled again.

// EmployerOf is where this person works, and nothing if they do not.
func (w *World) EmployerOf(id string) string {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		for _, who := range prop.Hands {
			if who == id {
				return l.ID
			}
		}
	}
	return ""
}

// AtWork is everybody the player employs, which is a thing worth being able to
// ask now that it has an answer.
func (w *World) AtWork() []*NPC {
	out := []*NPC{}
	for _, l := range Locations {
		if !w.Own(l.ID) || w.Properties[l.ID] == nil {
			continue
		}
		for _, who := range w.Properties[l.ID].Hands {
			if n := w.NPC(who); n != nil && !n.Dead {
				out = append(out, n)
			}
		}
	}
	return out
}

// HolderFaction is whose organization holds this address, if any.
func (w *World) HolderFaction(id string) string {
	if prop := w.Properties[id]; prop != nil {
		return prop.Owner
	}
	return ""
}

// takeOn finds somebody in this city to stand behind that counter. Nobody with
// a job already, nobody a family owns, nobody who is dead or inside. Returns
// the empty string when there is nobody, which is a real answer: a city can run
// out of people willing to work for you.
func (w *World) takeOn(id string) string {
	best, bestScore := "", -1
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Held > w.Minute || IsOfficial(n.ID) || n.Rank >= RankLeader {
			continue
		}
		// Somebody who already holds a position keeps it, and the city keeps
		// them where it is: a driver taken on at a laundry had the laundry
		// written down as their post and went on standing in the bar for the
		// rest of the week, because a role holder is exempt from the errand
		// that would walk them there.
		if w.keepsPost(n) {
			continue
		}
		if w.EmployerOf(n.ID) != "" {
			continue
		}
		// Somebody another family already has is not looking for a counter job:
		// they have one, and the city sends them off to mind their own
		// holdings, so a laundry that hired one had a position filled by
		// somebody who was never behind the counter. Measured over a week of
		// daytimes: one of three hands was never once in the room.
		// A family's own people work for their family, not behind somebody
		// else's counter.
		if n.Faction != "" && n.Faction != w.HolderFaction(id) {
			continue
		}
		// Somebody already spending their day there is the obvious hire, then
		// somebody who trusts you, then anybody at all.
		score := n.Trust
		if n.Post == id || n.Location == id {
			score += 50
		}
		if score > bestScore {
			best, bestScore = n.ID, score
		}
	}
	return best
}

// putToWork sets somebody behind a counter: the day is spent there, which is
// what having a job means to everything else in this city.
func (w *World) putToWork(who, id string) {
	n := w.NPC(who)
	if n == nil {
		return
	}
	prop := w.Properties[id]
	prop.Hands = append(prop.Hands, who)
	n.Post = id
	// They start today. Writing down where somebody works and leaving them
	// across the city is how a counter ends up staffed by nobody.
	if !w.Travelling(n) && !Evening(w.Minute) {
		n.Location = id
	}
}

// letGo takes the last one on and gives them their day back.
func (w *World) letGo(id string) {
	prop := w.Properties[id]
	if len(prop.Hands) == 0 {
		return
	}
	who := prop.Hands[len(prop.Hands)-1]
	prop.Hands = prop.Hands[:len(prop.Hands)-1]
	if n := w.NPC(who); n != nil && n.Post == id {
		n.Post = ""
	}
}

// EmptyChairs takes the dead off the books and puts a name to any position that
// has not got one. A position held by somebody who is not coming in is not a
// position that is filled, and the count has to say so or the place goes on
// handling work nobody is there to do.
//
// The naming half is also how a save written before anybody had a job acquires
// one, and how a business the player has just taken over gets the people who
// were already working in it.
func (w *World) EmptyChairs() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		kept := make([]string, 0, len(prop.Hands))
		for _, who := range prop.Hands {
			if n := w.NPC(who); n != nil && !n.Dead {
				kept = append(kept, who)
			}
		}
		if len(kept) != len(prop.Hands) {
			prop.Hands = kept
			prop.Staff = min(prop.Staff, len(kept))
		}
		// Every address somebody holds, not only the player's. A rival with a
		// laundry has somebody in it, or there is nobody to poach, nobody to
		// lean on, and nobody to lose when the place is taken.
		if prop.Owner == "" || prop.Owner == "independent" {
			continue
		}
		if trade, ok := TradeOf(l.ID); ok && prop.Staff == 0 && prop.Income > 0 && !w.wordIsOut(l.ID) {
			prop.Staff = trade.Hands
		}
		// A counter that somebody has just walked out of stays short. Without
		// this the city handed the position straight back to the person who
		// had had enough of it, on the same morning they left.
		if w.Minute < prop.Shorthanded {
			continue
		}
		// And nobody comes to work at a place that is not paying. Without
		// this, a player who had paid nobody for twenty-five days emptied
		// their laundry and the city handed them three fresh hands two days
		// later, for nothing, over and over — nineteen unpaid nights and
		// strangers still turning up.
		if w.wordIsOut(l.ID) {
			w.nobodyWillWork(l.ID)
			continue
		}
		for len(prop.Hands) < prop.Staff {
			who := w.takeOn(l.ID)
			if who == "" {
				// Nobody left in this city willing to stand behind it, which is
				// a real answer rather than a reason to invent somebody.
				prop.Staff = len(prop.Hands)
				break
			}
			w.putToWork(who, l.ID)
		}
	}
}

// HandsDescription is who stands behind that counter, for the room to draw.
// Nothing at all where the player has no business knowing, which is everywhere
// they do not hold: a rival's payroll is not public.
func (w *World) HandsDescription(id string) []map[string]any {
	prop := w.Properties[id]
	out := []map[string]any{}
	if prop == nil || !w.Own(id) || len(prop.Hands) == 0 {
		return out
	}
	for _, who := range prop.Hands {
		n := w.NPC(who)
		if n == nil {
			continue
		}
		out = append(out, map[string]any{
			"id": n.ID, "name": n.Name, "role": n.Role,
			"here": n.Location == id && !w.Travelling(n),
		})
	}
	return out
}

// Poaching. Now that the people behind a counter are people, they can be taken.
// Somebody good is somebody a rival is already paying, and money is the whole
// of the argument — which is also the cheapest way for a business war to be
// fought without anybody being shot.

// PoachPays is the signing money, as weeks of the wage the position carries.
const PoachPays = 3

// PoachCost is what it takes to walk somebody off a rival's counter and onto
// yours. The wage is the one the new position carries, because that is what the
// player is committing to pay.
func (w *World) PoachCost(into string) int {
	trade, ok := TradeOf(into)
	if !ok {
		return 0
	}
	return trade.Wage * 7 * PoachPays
}

// PoachReadiness explains why somebody cannot be taken on, or returns "".
func (w *World) PoachReadiness(who, into string) string {
	n := w.NPC(who)
	if n == nil || n.Dead {
		return "There is nobody here by that name"
	}
	trade, ok := TradeOf(into)
	if !ok || !w.Own(into) {
		return "That is not a business of yours"
	}
	from := w.EmployerOf(who)
	if from == "" {
		return n.Name + " works for nobody, and can simply be hired"
	}
	if w.Own(from) {
		return n.Name + " already works for you"
	}
	prop := w.Properties[into]
	if prop.Staff >= trade.Hands {
		place, _ := PlaceByID(into)
		return "Every position at " + place.Name + " is filled"
	}
	if w.Player.Cash < w.PoachCost(into) {
		return "Not enough cash to make it worth their while"
	}
	return ""
}

// Poach takes somebody off a rival's books and puts them behind your counter.
func (w *World) Poach(who, into string) error {
	if reason := w.PoachReadiness(who, into); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	from := w.EmployerOf(who)
	if err := w.Pay(w.PoachCost(into)); err != nil {
		return err
	}
	old := w.Properties[from]
	kept := make([]string, 0, len(old.Hands))
	for _, id := range old.Hands {
		if id != who {
			kept = append(kept, id)
		}
	}
	old.Hands = kept
	old.Staff = min(old.Staff, len(kept))
	n := w.NPC(who)
	w.Properties[into].Staff++
	w.putToWork(who, into)
	// Money is not friendship, but it is not nothing either.
	n.Trust = min(100, n.Trust+8)
	// The family they were working for takes it as what it is.
	if f := w.faction(old.Owner); f != nil {
		f.Goodwill = max(-100, f.Goodwill-PoachGalls)
	}
	here, _ := PlaceByID(into)
	there, _ := PlaceByID(from)
	w.Log(n.Name+" comes to work for you",
		fmt.Sprintf("$%d to walk out of %s and behind the counter at %s. %s will have noticed.",
			w.PoachCost(into), there.Name, here.Name, w.HolderName(from)), "business")
	return nil
}

// PoachGalls is what a family thinks of somebody taking their people. Small
// beside a burnt-out holding and large beside nothing at all, which is what
// this used to cost.
const PoachGalls = 9

// shortHanded is a business of the player's with a position going, and nothing
// if every counter they hold is covered.
func (w *World) shortHanded() string {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, ok := TradeOf(l.ID)
		if prop == nil || !ok || !w.Own(l.ID) {
			continue
		}
		if prop.Staff < trade.Hands {
			return l.ID
		}
	}
	return ""
}

// What the counter sees. The people behind your counters watch the street all
// day, and somebody standing across the road from your laundry for three
// afternoons is the sort of thing they would mention. Until this, the only way
// to learn that a family had commissioned an attack on a business of yours was
// to go and investigate it — so a business the player owned was a number that
// could be broken, and the people in it were not looking out of the window.
//
// It also gives the trust of somebody you employ a job. A man who thinks well
// of you tells you what he saw; a man who does not keeps his head down, which
// is the difference between the two kinds of employee a player can have.

// NoticeAt is the best odds anybody has of spotting it in a day, at full trust.
// Kept low: a business that always sees it coming has nothing to fear, and
// investigating would stop being worth the trip.
const NoticeAt = .14

// ReadyFor is how much of an attack each pair of hands turns away when the
// place is expecting it.
const ReadyFor = 9

// WordFromTheCounter is one day of your people watching their own street. Called
// every business day.
func (w *World) WordFromTheCounter() {
	for i := range w.Plots {
		p := &w.Plots[i]
		if p.Life != w.Life || p.Known || p.Target == "" {
			continue
		}
		prop := w.Properties[p.Target]
		if prop == nil || !w.Own(p.Target) || len(prop.Hands) == 0 {
			continue
		}
		for _, who := range prop.Hands {
			n := w.NPC(who)
			if n == nil || n.Dead {
				continue
			}
			// Somebody who thinks nothing of you saw the same street and said
			// nothing about it.
			if w.WorldRandom() >= NoticeAt*float64(n.Trust)/100 {
				continue
			}
			p.Known = true
			place, _ := PlaceByID(p.Target)
			w.Log("A word from "+place.Name,
				fmt.Sprintf("%s says the same car has been parked across from %s three afternoons running, and the same two people in it who never get out. Somebody is looking the place over.",
					n.Name, place.Name), "danger")
			break
		}
	}
}

// HandNames is who is behind that counter, for a sentence.
func (w *World) HandNames(id string) string {
	prop := w.Properties[id]
	if prop == nil {
		return "nobody"
	}
	names := make([]string, 0, len(prop.Hands))
	for _, who := range prop.Hands {
		if n := w.NPC(who); n != nil {
			names = append(names, n.Name)
		}
	}
	if len(names) == 0 {
		return "nobody"
	}
	return joinNames(names)
}

// Asking the counter. The people behind yours have names, a wage and a view of
// the street, and until now nothing to say: they noticed a car parked across
// the road on their own schedule and the player could not ask. Everything in
// the answer is a fact the city already holds — who has been in, what the trade
// is doing, whether somebody is looking the place over — so this is a way of
// reading the world through somebody who lives in it rather than a new source
// of anything.

// AskReadiness explains why this person will not answer, or returns "".
func (w *World) AskReadiness(who string) string {
	n := w.NPC(who)
	if n == nil || n.Dead {
		return "There is nobody here by that name"
	}
	at := w.EmployerOf(who)
	if at == "" {
		return n.Name + " does not work for anybody"
	}
	if !w.Own(at) {
		return n.Name + " works for " + w.HolderName(at)
	}
	if n.Location != at || w.Travelling(n) {
		return n.Name + " is not behind the counter"
	}
	return ""
}

// AskTheCounter is what they saw from behind it.
func (w *World) AskTheCounter(who string) error {
	if reason := w.AskReadiness(who); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n := w.NPC(who)
	at := w.EmployerOf(who)
	place, _ := PlaceByID(at)
	prop := w.Properties[at]

	said := []string{}
	// The thing worth saying first, if there is one. By the index, because a
	// range over the plans hands out copies of them: the player was told about
	// the car across the road and the world went on not knowing, so a place
	// that had been warned by its own people took the full damage anyway.
	for i := range w.Plots {
		p := &w.Plots[i]
		if p.Life != w.Life || p.Target != at {
			continue
		}
		p.Known = true
		said = append(said, "There has been a car across the road three afternoons running, and nobody gets out of it.")
		break
	}
	// Then the money, before anything about the work. Somebody who has not
	// been paid for a week is not going to open with the footfall.
	if prop.Unpaid > 0 {
		if w.wordIsOut(at) {
			said = append(said, fmt.Sprintf("Nobody here has been paid in %s. They say it plainly, and they say it first.",
				plural(prop.Unpaid, "night", "nights")))
		} else {
			said = append(said, fmt.Sprintf("The wages are %s behind.", plural(prop.Unpaid, "night", "nights")))
		}
	}
	switch {
	case prop.Trouble:
		said = append(said, "Something has gone wrong in the back and nobody has put it right.")
	case prop.Staff < handsWanted(at):
		said = append(said, fmt.Sprintf("There are %d of us doing the work of %d.", prop.Staff, handsWanted(at)))
	case prop.Supply == 0:
		said = append(said, "We are out of nearly everything.")
	}
	// Then the things the books cannot hold.
	//
	// Everything above this line is the state of the premises, and the player
	// can read all of it off "Review the books" without talking to anybody. A
	// counter is only worth talking to for what a person standing in a room all
	// day sees and a ledger never records: who is walking over, who has been
	// asking, and the one thing this particular trade is in a position to
	// notice. Without those the answer was "three people through the door
	// today", which is a number the game already had.
	if coming := w.walkingOver(at); coming != "" {
		said = append(said, coming)
	}
	if asking := w.beenAsking(at); asking != "" {
		said = append(said, asking)
	}
	if seen := w.fromBehindThisCounter(at); seen != "" {
		said = append(said, seen)
	}
	if n := w.Footfall(at); n > 0 {
		said = append(said, fmt.Sprintf("%s through the door today.", counted(n, "person", "people")))
	} else {
		said = append(said, "Nobody at all today.")
	}
	if n.Sore > 0 {
		said = append(said, "They say it without looking at you.")
	} else if n.Trust >= 40 {
		said = append(said, "They seem glad you asked.")
	}
	w.Log(n.Name+" behind the counter",
		fmt.Sprintf("At %s. %s", place.Name, strings.Join(said, " ")), "personal")
	n.Trust = min(100, n.Trust+2)
	return nil
}

// handsWanted is how many positions this business has, and none where it is not
// a business that runs on people.
func handsWanted(id string) int {
	if trade, ok := TradeOf(id); ok {
		return trade.Hands
	}
	return 0
}

// Leaving. You can walk somebody off a rival's counter and nothing walks
// anybody off yours, so employment only ever happened in one direction. A
// business the city can take people out of is a business worth defending, and
// it is the reason to care what the people behind your counters think of you.

const (
	// WalksOut is the chance a day that somebody who has had enough actually
	// goes. Low: they leave on a day of their own choosing rather than the
	// moment a number crosses a line, which is what makes it feel like a
	// decision somebody made rather than a threshold.
	WalksOut = .12
	// FindingSomebody is how long a counter stays short after somebody walks
	// off it. Two days: long enough that losing a pair of hands is felt, short
	// enough that a business is not permanently crippled by one bad week.
	FindingSomebody = 2 * 1440
	// TemptedAway is the day's chance that somebody paid the least anybody
	// stands there for listens to a rival who is hiring. Low, because it is a
	// day: over two months it is most of a counter, which is what a business
	// run on the floor wage should lose.
	TemptedAway = .035
)

// Notice is the day's chance that anybody behind your counters has had enough
// of it, and that a rival with a position going takes one of them. Called every
// business day, after the dead have come off the books.
func (w *World) Notice() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || len(prop.Hands) == 0 || !w.Own(l.ID) {
			continue
		}
		to := w.hiringElsewhere(l.ID)
		for _, who := range append([]string{}, prop.Hands...) {
			n := w.NPC(who)
			if n == nil || n.Dead {
				continue
			}
			odds := 0.
			if w.hadEnough(n) {
				odds = WalksOut
			}
			// And a week of not being paid. This was written once and taken
			// out again for being unreachable: the day's bill was all or
			// nothing, so one short night stripped the security and the
			// address and every night after cleared out of what the business
			// earned. The bill gives things up in order now, so a place that
			// earns less than its counter costs goes unpaid night after night,
			// and this is a week somebody actually lives through.
			if prop.Unpaid >= PatienceRunsOut {
				odds = max64(odds, NotBeingPaid)
			}
			// And what the player pays against what the work is worth in this
			// city. Somebody paid over the rate is not listening to anybody;
			// somebody paid the floor is, whatever they think of you. The city
			// is bigger than any one family, so this does not wait for a
			// particular rival to be hiring — where they land does.
			odds = max64(odds, w.tempted(l.ID))
			if odds <= 0 || w.WorldRandom() >= odds {
				continue
			}
			w.walkOut(who, l.ID, to)
			break // one a day at each address, or a bad week empties the place
		}
	}
}

// hadEnough is why somebody stops standing behind your counter: they are
// carrying something against you.
//
// Not a trust score. Everybody in this city starts at nothing and thinks
// nothing of a stranger, so "below twenty" is every employee in the game and a
// business would bleed people for no reason anybody could name — measured, a
// laundry went from three hands to none in a month with nothing having happened.
// A second cause was tried, a place in trouble under somebody who dislikes you,
// and it failed the same way for the same reason. What is left is one rule that
// names something the player did.
func (w *World) hadEnough(n *NPC) bool {
	return n.Sore >= TakesItPersonally
}

// tempted is the day's chance that somebody standing behind this counter
// listens to somebody else who is hiring. It is about paying under the rate and
// nothing else: at the rate, which is what the work is worth in this city, they
// have no reason to move, and the whole of TemptedAway is what somebody paid
// the floor is worth to a rival.
//
// The first version made the rate itself worth half of it, and a business
// paying exactly what the work is worth bled people for no reason anybody could
// name — which is the same mistake as the first version of hadEnough, made
// again a hundred lines further down.
func (w *World) tempted(id string) float64 {
	trade, ok := TradeOf(id)
	if !ok {
		return 0
	}
	paid := w.WageAt(id)
	least, _ := WageBounds(trade.Wage)
	if paid >= trade.Wage || least >= trade.Wage {
		return 0
	}
	short := float64(trade.Wage-paid) / float64(trade.Wage-least)
	return TemptedAway * short
}

// hiringElsewhere is a rival's address with a position going and the money to
// fill it, and nothing when nobody is hiring.
func (w *World) hiringElsewhere(from string) string {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, ok := TradeOf(l.ID)
		if prop == nil || !ok || l.ID == from || w.Own(l.ID) {
			continue
		}
		if prop.Owner == "" || prop.Owner == "independent" || prop.Staff >= trade.Hands {
			continue
		}
		if f := w.faction(prop.Owner); f == nil || f.Cash < trade.Wage*7 {
			continue
		}
		return l.ID
	}
	return ""
}

// walkOut takes somebody off one set of books and, where there is somewhere to
// go, puts them on another.
func (w *World) walkOut(who, from, to string) {
	old := w.Properties[from]
	kept := make([]string, 0, len(old.Hands))
	for _, id := range old.Hands {
		if id != who {
			kept = append(kept, id)
		}
	}
	old.Hands = kept
	old.Staff = min(old.Staff, len(kept))
	old.Shorthanded = w.Minute + FindingSomebody
	n := w.NPC(who)
	if n != nil && n.Post == from {
		n.Post = ""
	}
	// And the title goes with the job. Somebody who no longer works here does
	// not run it, and a role saying they do is a lie the rest of the city reads
	// — the routine keeps a manager standing at their own address, so a stale
	// one would have somebody minding a counter they had been put off.
	if here, ok := PlaceByID(from); ok && n != nil && n.Role == "Runs "+here.Name {
		n.Role = "Out of work"
	}
	here, _ := PlaceByID(from)
	if to == "" {
		w.Log(n.Name+" has had enough",
			fmt.Sprintf("They are not behind the counter at %s this morning, and nobody expects them back. There are %d of you now.",
				here.Name, old.Staff), "business")
		return
	}
	w.Properties[to].Staff++
	w.putToWork(who, to)
	there, _ := PlaceByID(to)
	w.Log(n.Name+" has gone to "+there.Name,
		fmt.Sprintf("%s pays better than you do, or asks less. %s is short-handed at %d.",
			w.HolderName(to), here.Name, old.Staff), "business")
}

// wordIsOut reports whether a place has gone long enough without paying
// anybody that the city knows about it. The same week somebody standing behind
// the counter will not work through is the week nobody new will start.
func (w *World) wordIsOut(id string) bool {
	prop := w.Properties[id]
	return prop != nil && prop.Unpaid >= PatienceRunsOut
}

// nobodyWillWork says it once, the morning a position first goes begging,
// because a place nobody will stand in is news and then it is the situation.
func (w *World) nobodyWillWork(id string) {
	prop := w.Properties[id]
	trade, ok := TradeOf(id)
	// The trade's own count rather than the position count: somebody walking
	// out takes their position with them, so by the time the word is out the
	// place wants nobody and nothing would ever be said.
	if !ok || !w.Own(id) || len(prop.Hands) >= trade.Hands || prop.Toldabout {
		return
	}
	prop.Toldabout = true
	place, _ := PlaceByID(id)
	w.Log("Nobody will stand behind the counter at "+place.Name,
		fmt.Sprintf("It has been %d nights since anybody there was paid, and that is the sort of thing "+
			"people in this city tell each other. Pay what is owed and somebody will take the job.",
			prop.Unpaid), "danger")
}

// walkingOver is somebody on the street with this address at the end of it.
// The player walks into a room and sees who is in it; the person who stands
// there all day sees who is coming, which is a quarter of an hour's warning and
// the only place in this game that gives any.
func (w *World) walkingOver(at string) string {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Heading != at || !w.Travelling(n) {
			continue
		}
		// Whose man he is, if he is anybody's, because that is the part worth
		// a warning. A hand would not say a name and stop.
		if f := w.faction(n.Faction); f != nil && f.ID != w.PlayerOrganizationID() {
			return fmt.Sprintf("%s of the %s is walking over, and will be here in %s.",
				n.Name, f.Name, counted(max(1, (n.Arrives-w.Minute)/10*10), "minute", "minutes"))
		}
		because := "on their way over"
		if n.Errand != "" {
			because = n.Errand
		}
		return fmt.Sprintf("%s is %s.", n.Name, because)
	}
	return ""
}

// beenAsking is somebody carrying something against the player, near enough to
// this room that the people in it would have heard about it.
//
// What somebody holds against the protagonist is `Sore`, not a grudge record —
// grudges are between two other people and deliberately stop short of the
// player. It is invisible until it arrives as an ordinary robbery in the street
// with a name attached. This is the half-hour before that.
func (w *World) beenAsking(at string) string {
	for _, n := range w.Aggrieved() {
		if n.Sore < SoreAsks {
			break
		}
		near := n.Location == at || n.Post == at || TravelMinutes(n.Location, at) <= AskingDistance
		if !near {
			continue
		}
		over := "something"
		if n.SoreAt != "" {
			over = n.SoreAt
		}
		if n.Sore >= SoreActs {
			return fmt.Sprintf("%s has been in twice this week asking when you are here, over %s. They were not making conversation.", n.Name, over)
		}
		return fmt.Sprintf("%s has been asking after you, over %s.", n.Name, over)
	}
	return ""
}

// fromBehindThisCounter is the one thing this particular trade is in a position
// to notice, drawn from that trade's own machinery rather than invented for it.
// A garage sees what came in broken. A pawnbroker sees who was short. A yard
// full of cabs knows where people went. It is the reason to hold more than one
// kind of place: each counter shows you a different part of the same city.
func (w *World) fromBehindThisCounter(at string) string {
	place, ok := PlaceByID(at)
	if !ok {
		return ""
	}
	switch place.Kind {
	case "garage":
		broken := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; !n.Dead && n.Hurt && n.Car > 0 {
				broken++
			}
		}
		if broken > 0 {
			return fmt.Sprintf("%s in this city driving with the glass out, and every one of them has to come to somebody.",
				counted(broken, "car", "cars"))
		}
		return "Nothing on the road is broken this week, which is the wrong kind of quiet for a garage."
	case "undertaker":
		// The one counter in this city whose trade is made entirely by
		// everybody else's. It counts the week's dead the way a garage counts
		// broken glass.
		week := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; n.Dead && n.DiedAt > 0 && w.Minute-n.DiedAt <= 7*1440 {
				week++
			}
		}
		if week > 0 {
			return fmt.Sprintf("%s buried out of this room since last week, and not one family asked a question about the arrangements.",
				counted(week, "person", "people"))
		}
		return "Nobody has died in this district all week. It is restful, and it does not pay for the cars."
	case "pawn":
		for i := len(w.Window) - 1; i >= 0; i-- {
			if s := w.Window[i]; s.Whose != "" {
				return fmt.Sprintf("%s was short enough to leave %s over the counter and has not been back.",
					s.Whose, lowerFirst(s.What()))
			}
		}
		if len(w.Window) > 0 {
			return fmt.Sprintf("%s in the window that nobody came back for.", counted(len(w.Window), "thing", "things"))
		}
		return "Nothing in the window. Either the city is flush or it has stopped trusting us."
	case "cabs":
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if !w.Travelling(n) || n.Heading == "" || n.Car > 0 {
				continue
			}
			to, ok := PlaceByID(n.Heading)
			if !ok {
				continue
			}
			return fmt.Sprintf("One of ours took %s over to %s this morning.", n.Name, to.Name)
		}
		return "The cars have been standing on the rank all day."
	case "filling":
		driving := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; !n.Dead && n.Car > 0 {
				driving++
			}
		}
		if driving == 0 {
			return "Nobody in this city is keeping a car on the road, and pumps with nobody at them are a wall with a roof."
		}
		return fmt.Sprintf("%s in this city keeping a car on the road, and the tank does not fill itself.",
			counted(driving, "person", "people"))
	case "butcher", "restaurant":
		// Somebody on this counter who has got far enough down to be thinking
		// about leaving. The same figure the room panel draws a warning from,
		// said out loud by the person standing next to them.
		for _, who := range w.Properties[at].Hands {
			n := w.NPC(who)
			if n == nil || n.Dead {
				continue
			}
			if n.Faction == w.PlayerOrganizationID() && n.Trust < DefectionTrust {
				return fmt.Sprintf("%s has been talking about going somewhere else.", n.Name)
			}
		}
		// And two rooms that share a case have to part company here, because
		// two counters saying the identical sentence is two cards in a room
		// under one name wearing a different hat. A shop and a dining room do
		// not see the same week.
		if place.Kind == "restaurant" {
			return "The same tables booked by the same people on the same nights, and one of them always wants the corner."
		}
		return "The orders are the orders. Same cuts, same days, and nothing on the hooks by four."

	// The eleven below were nothing at all until they were counted. The guard
	// over this asked five kinds whether any two said the same thing, and
	// passed while two thirds of the city's counters said nothing — a guard
	// over a sample, which is the same fault three other measures had tonight.
	// It asks every trade now.
	case "club":
		crowd := w.TheEveningCrowd(at)
		if crowd == 0 {
			return "Nobody has been in but the staff, and they are drinking the profit."
		}
		return fmt.Sprintf("%s make an evening of it here, and that is who this room is.",
			counted(crowd, "person", "people"))
	case "casino":
		if prop := w.Properties[at]; prop != nil && prop.Bankroll > 0 {
			return fmt.Sprintf("There is $%d behind the tables, which is what decides who will sit down at them.", prop.Bankroll)
		}
		return "There is nothing behind the tables, so nobody serious is going to sit at one."
	case "poolhall":
		if w.BackRoomTake > 0 {
			return fmt.Sprintf("The room has taken $%d in seat money off the game in the back.", w.BackRoomTake)
		}
		return "Nobody has sat down in the back this week. The cloth is the cleanest it has ever been."
	case "saloon":
		if drinkers := w.TheEveningCrowd(at); drinkers > 0 {
			return fmt.Sprintf("%s drink here every evening of their lives, and none of them lowers their voice.",
				counted(drinkers, "person", "people"))
		}
		return "Nobody has been in all week, and an empty bar hears nothing worth hearing."
	case "haulage":
		runs := 0
		for _, l := range Locations {
			if _, ok := TradeOf(l.ID); ok && w.Own(l.ID) && l.ID != at {
				runs++
			}
		}
		if runs == 0 {
			return "The trucks are running for other people, which is a waste of a yard you hold."
		}
		return fmt.Sprintf("The trucks are stocking %s of yours besides this one, and every one of them cheaper for it.",
			counted(runs, "place", "places"))
	case "scrapyard":
		gone := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; !n.Dead && n.Drove > 0 && n.Car == 0 {
				gone++
			}
		}
		if gone == 0 {
			return "Nobody in this city has lost a car this week, and a yard lives on the ones that stop going anywhere."
		}
		return fmt.Sprintf("%s in this city used to drive and do not any more, and a car goes somewhere when it stops going anywhere.",
			counted(gone, "person", "people"))
	case "laundry":
		if w.Player.Heat == 0 {
			return "Quiet books, and nobody asking you anything. That is the week to do the paperwork."
		}
		return fmt.Sprintf("The books will absorb %d of the %d points of attention on you, and they will hold.",
			min(w.launderCapacity(at), w.Player.Heat), w.Player.Heat)
	case "exchange":
		// Not Footfall. That counts whoever is standing here, which includes
		// the person answering the question the moment they step behind the
		// counter — the figure moved between one call and the next and the
		// guard over it caught the wobble. Who buys on this floor is a fact
		// about the city, not about who is in the doorway.
		buyers := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; !n.Dead && n.Purse >= w.PriceAt(at, "moonshine") && n.Post != at {
				buyers++
			}
		}
		if buyers == 0 {
			return "Nobody in this city could cover a crate at today's price, which is a floor with nobody on it."
		}
		return fmt.Sprintf("%s in this city could cover a crate at today's price, and every one of them hears what this floor says.",
			counted(buyers, "person", "people"))
	case "burlesque":
		if w.NightOn(at) {
			return "There is a night on, and the room will be full of people who drink somewhere else."
		}
		return "An ordinary evening, which for a room like this means the chairs outnumber the people."
	case "dealer":
		walking := 0
		for i := range w.NPCs {
			if n := &w.NPCs[i]; !n.Dead && n.Car == 0 && n.Drove == 0 {
				walking++
			}
		}
		if walking == 0 {
			return "Everybody in this city drives something, which is a lot with nothing to sell."
		}
		return fmt.Sprintf("%s in this city have never owned a car, and every one of them walks past a lot like this.",
			counted(walking, "person", "people"))
	case "wharf":
		if w.BoatIsIn() && w.Landed.Where == at {
			g := w.Good(w.Landed.Good)
			return fmt.Sprintf("There is a boat on the quay now with %d %ss of %s on it, and they are not waiting.",
				w.Landed.Units, g.Unit, g.InBulk())
		}
		return "Nothing on the quay tonight. The word goes round when there is."
	}
	return ""
}

// lastHand is whoever letGo would take off the books next, so a caller can
// name them before they are gone.
func (w *World) lastHand(id string) string {
	prop := w.Properties[id]
	if prop == nil || len(prop.Hands) == 0 {
		return ""
	}
	return prop.Hands[len(prop.Hands)-1]
}
