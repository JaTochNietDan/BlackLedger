package core

// Ordinary residents return home overnight, commute from six in the morning,
// and visit their usual evening venues after noon. Work and urgent errands
// take precedence in wanted(). Destinations remain stable across saved games.

// Haunts are the places somebody goes when the day's work is done. All three
// sell a drink; two of them will also take a bet.
// Where the city goes of an evening. The poolhall is on this list because the
// back room is behind it and a game needs three people in the room: measured
// before it was added, the most anybody ever saw in the poolhall in a week was
// one person, at nine in the morning, which made a whole card game unplayable
// in a city that had one. The other fault shape in the brief, a content table
// written for a smaller city, wearing the clothes of a feature that works.
var haunts = []string{"bar", "club", "casino", "poolhall"}

// How the evening divides. A poolhall is not a bar: a few people go, not a
// quarter of the city. Splitting the evening four equal ways moved twenty
// people a night out of the drinking places, and the drop in their takings was
// enough that a diplomat who used to buy the casino inside 220 commands could
// no longer afford one on any seed. One share in twelve is enough to get a
// three handed game up most evenings and small enough that the rest of the
// city's night is where it was.
var haunted = map[string]int{"bar": 4, "club": 4, "casino": 3, "poolhall": 1}

// Evening reports which half of the day it is. Everything that cares asks here.
// EveningFrom is the minute of the day the city stops working and goes out.
// It is also when the clock sends everybody who is going anywhere, so it is
// the hour a night either catches or misses.
const EveningFrom = 720

func Evening(minute int) bool { return minute%1440 >= EveningFrom }

const HomeUntil = 6 * 60

func overnight(minute int) bool { return minute%1440 < HomeUntil }

// haunt is where this person drinks, fixed for life. It is derived from their
// id rather than stored, so it survives every save ever written and cannot
// drift: the whole point is that they are in the same place every evening.
func haunt(id string) string {
	sum := 0
	for i := 0; i < len(id); i++ {
		sum = sum*31 + int(id[i])
	}
	if sum < 0 {
		sum = -sum
	}
	total := 0
	for _, id := range haunts {
		total += haunted[id]
	}
	at := sum % total
	for _, id := range haunts {
		if at -= haunted[id]; at < 0 {
			return id
		}
	}
	return haunts[0]
}

// keepsPost is true for the people whose position is not theirs to choose. A
// role holder is on duty, somebody the police are holding is not going
// anywhere, and anybody with a family holding to mind has already been given a
// better reason by wanted().
func (w *World) keepsPost(n *NPC) bool {
	if n.Held > w.Minute || IsOfficial(n.ID) {
		return true
	}
	for _, office := range officials {
		if n.Role == office.Role {
			return true
		}
	}
	for _, r := range roles {
		if n.Role == r.Title {
			return true
		}
	}
	for _, l := range Locations {
		if n.Role == "Runs "+l.Name {
			return true
		}
	}
	// The leader of a family is not found propping up a bar.
	return n.Rank >= RankLeader
}

// routine is where the hour says this person ought to be. It is the last
// reason wanted() considers, so anybody with work to do does the work.
func (w *World) routine(n *NPC) (errand, bool) {
	if w.keepsPost(n) || n.Post == "" {
		return errand{}, false
	}
	if overnight(w.Minute) && n.Home != "" {
		if place, ok := PlaceByID(n.Home); ok && n.Home != n.Location {
			return errand{n.Home, "heading home to " + place.Name}, true
		}
		return errand{}, false
	}
	if Evening(w.Minute) {
		// Nobody walks to a bar on the night two families are shooting at each
		// other. Most people stay where they are, which is what a war costs an
		// owner who is not in it.
		if w.CityAtWar() && w.staysIn(n) {
			return errand{}, false
		}
		where := haunt(n.ID)
		// Unless somebody has put a night on within walking distance, which is
		// the whole of what paying for a band buys. One night: tomorrow they
		// are back where they always are.
		if on := w.theNight(where, n.ID); on != "" {
			where = on
		}
		if where == n.Post || where == n.Location {
			return errand{}, false
		}
		place, ok := PlaceByID(where)
		if !ok {
			return errand{}, false
		}
		return errand{where, "going to " + place.Name}, true
	}
	if n.Post == n.Location {
		return errand{}, false
	}
	place, ok := PlaceByID(n.Post)
	if !ok {
		return errand{}, false
	}
	return errand{n.Post, "back to " + place.Name}, true
}

// keepPost records where a person's day is. It is taken the first time they
// are seen standing anywhere in the daytime, which bootstraps every save
// written before the city had a shift, and updated whenever work moves them.
func (w *World) keepPost(n *NPC, where string) {
	if _, ok := PlaceByID(where); ok {
		n.Post = where
	}
}
