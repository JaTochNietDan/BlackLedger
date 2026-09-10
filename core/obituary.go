package core

import (
	"fmt"
	"strings"
)

// The morning after.
//
// A killing is reported the same day, as it should be: a body was found, police
// say enquiries are continuing. That is the crime desk, and it is about what
// happened. It is not about who it happened to.
//
// A paper carries the other thing the next morning. Somebody who ran a business
// on the harbour road for two years, or stood behind a family boss for longer
// than most of them last, gets three sentences about what they were. This is
// the only place in the game where a person is described as a life rather than
// as a threat, an asset or an obstacle — and in a city that remembers, the day
// after somebody dies is exactly when it should say so.
//
// It knows nothing it should not. An obituary never says who arranged it; it
// says what the city could see of a person while they were alive.

// mourned is whether the city would notice this death enough to print anything.
// A soldier nobody had heard of does not get a column, and neither does
// somebody who was never given a name to be known by.
func (w *World) mourned(n *NPC) bool {
	if n == nil || n.Name == "" {
		return false
	}
	if _, official := OfficialByID(n.ID); official {
		return true // a man with a title always gets one
	}
	if n.Rank >= RankLieutenant {
		return true // somebody who stood over other people
	}
	// Or somebody the city had reason to know: they ran premises, or the
	// player had actually dealt with them.
	for _, l := range Locations {
		if prop := w.Properties[l.ID]; prop != nil && n.Location == l.ID && n.Rank >= RankSoldier {
			return true
		}
	}
	return n.Trust > 0
}

// obituary writes what the city could say about a person, and nothing else.
func (w *World) obituary(n *NPC) (string, string) {
	headline := strings.ToUpper(n.Name)

	// What they were. Not describeStanding: that appends the family's name to
	// the role, and a boss's role already carries the family in it, so the
	// paper was printing "Russo boss of Russo Outfit".
	opening := "They answered to nobody."
	if role := n.Role; role != "" {
		job := lowerFirst(role)
		opening = "They were " + withArticle(job) + "."
		if f := w.faction(n.Faction); f != nil && !sharesAName(role, f.Name) {
			opening = "They were " + withArticle(job) + " of " + f.Name + "."
		}
	} else if f := w.faction(n.Faction); f != nil {
		opening = "They were one of " + f.Name + "."
	}

	// Where they were usually found, which is how most of the city knew them.
	where := ""
	if place, ok := PlaceByID(n.Location); ok {
		where = fmt.Sprintf(" Those who had business with them had it at %s.", place.Name)
	}

	// What they were holding when they died, which anybody walking past can see
	// has changed hands or has not.
	// Their own place first: naming a holding across the city when the man was
	// found at the door of one he actually ran reads as though the paper had
	// picked a building at random, which is exactly what it was doing.
	left := ""
	if n.Rank >= RankLieutenant {
		for _, id := range []string{n.Location, ""} {
			for _, l := range Locations {
				if id != "" && l.ID != id {
					continue
				}
				if prop := w.Properties[l.ID]; prop != nil && prop.Owner == n.Faction {
					left = fmt.Sprintf(" %s carries on under whoever is left to carry it on.", l.Name)
					break
				}
			}
			if left != "" {
				break
			}
		}
	}

	// And whether anybody is left who answered to them, because that is the
	// part of a death the rest of the city actually feels.
	under := 0
	for i := range w.NPCs {
		other := &w.NPCs[i]
		if !other.Dead && other.ID != n.ID && other.Faction == n.Faction && other.Faction != "" &&
			other.Rank < n.Rank {
			under++
		}
	}
	after := ""
	if under == 1 {
		after = " One person in the same organization stood below them, and will be told by somebody."
	} else if under > 1 {
		after = fmt.Sprintf(" %s people in the same organization stood below them, and will be told by somebody.", upper1(spelled(under)))
	}

	return headline, opening + where + left + after
}

// ObituaryDay files the previous day's obituaries. Called from the clock, so
// the paper carries them the morning after rather than alongside the killing.
func (w *World) ObituaryDay() {
	today := w.Minute / 1440
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead || n.Remembered || n.DiedAt/1440 != today-1 {
			continue
		}
		// Marked whether or not a column is written, so somebody the city would
		// not have noticed is passed over once rather than considered forever.
		n.Remembered = true
		if !w.mourned(n) {
			continue
		}
		headline, body := w.obituary(n)
		w.ReportAbout("obituary", headline, body, n.ID)
	}
}

// sharesAName reports whether a role already carries the family's name in it,
// so the paper does not print "Russo boss of Russo Outfit".
func sharesAName(role, family string) bool {
	for _, word := range strings.Fields(family) {
		if len(word) > 3 && strings.Contains(role, word) {
			return true
		}
	}
	return false
}
