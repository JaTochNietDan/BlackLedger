package core

// Measured before this was written: over one simulated week, sampled every
// hour, forty-seven of fifty people never moved at all and three moved once.
// The busiest places at three in the morning were the busiest places at three
// in the afternoon, by the same margin. The city had errands — a man crosses
// town to settle a grudge, a family sends somebody to mind a holding it has
// just taken — but every one of them is a reason that ends, and once everybody
// had a reason to be where they already were, nobody ever moved again.
//
// What was missing is the ordinary reason: the day ends and people go out.
// This gives the city a shift. In the first half of the day people are at
// their posts; in the second they are where they drink. Which half it is comes
// from the clock, and where a person drinks never changes, so the player can
// learn that Ivo Costa is at The Blue Hour in the evening and find him there
// every evening for the rest of his life.
//
// The resolution is a half-day because that is the resolution the clock has:
// the world only stops at twelve-hour boundaries, and pretending to finer
// grain than that would be a lie told by this file rather than a fact about
// the city. In practice a person sets off around noon and is drinking by half
// past, which for this city is not early.

// Haunts are the places somebody goes when the day's work is done. All three
// sell a drink; two of them will also take a bet.
var haunts = []string{"bar", "club", "casino"}

// Evening reports which half of the day it is. Everything that cares asks here.
func Evening(minute int) bool { return minute%1440 >= 720 }

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
	return haunts[sum%len(haunts)]
}

// keepsPost is true for the people whose position is not theirs to choose. A
// role holder is on duty, somebody the police are holding is not going
// anywhere, and anybody with a family holding to mind has already been given a
// better reason by wanted().
func (w *World) keepsPost(n *NPC) bool {
	if n.Held > w.Minute || IsOfficial(n.ID) {
		return true
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
	if Evening(w.Minute) {
		where := haunt(n.ID)
		// Unless somebody has put a night on within walking distance, which is
		// the whole of what paying for a band buys. One night: tomorrow they
		// are back where they always are.
		if on := w.theNight(n.Location); on != "" {
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
