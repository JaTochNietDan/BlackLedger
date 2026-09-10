package sim

import "blackledger/core"

// Watching the city rather than the player.
//
// docs/LIVING_WORLD.md sets the standard for these systems: "a system is not
// finished because it runs. It is finished when a long simulation shows it
// producing varied, non-degenerate outcomes." The campaign runner cannot show
// that. It follows one protagonist, and the median campaign lasts between three
// and fourteen days — the reckless one lasts five hours — so `factions_destroyed`
// read zero in every report the harness has ever printed. Not because nothing
// ever falls: left alone for sixty days, ten organizations fell across
// twenty-five cities. The measurement was watching the wrong thing.
//
// So this runs the city with nobody in it. No policy, no commands: the clock
// turns and the families do whatever the rules make them do. What it reports is
// the shape of the world after a season, which is the only way to tell a living
// city from one that quietly collapses into a single owner or a graveyard.

// CityReport is what a season did to one city.
type CityReport struct {
	Seed uint32 `json:"seed"`
	Days int    `json:"days"`
	// Organizations at the start and at the end, and the churn between.
	Started  int `json:"organizations_started"`
	Ended    int `json:"organizations_ended"`
	Formed   int `json:"organizations_formed"`
	Fell     int `json:"organizations_fell"`
	Wars     int `json:"wars_started"`
	Settled  int `json:"wars_settled"`
	Changed  int `json:"holdings_changed_hands"`
	Killed   int `json:"people_killed"`
	Living   int `json:"people_living"`
	Homeless int `json:"holdings_unheld"`
	// Biggest is the share of held property in one organization's hands at the
	// end, as a percentage. A city that ends with one owner has degenerated
	// however busy it looked getting there.
	Biggest int `json:"biggest_share"`
}

// City runs one city for a season with nobody playing it.
func City(seed uint32, days int) CityReport {
	w := core.New(seed)
	// Two streams, and the world's own is what everything off-screen draws
	// from. Striding it keeps two cities from having the same weather.
	w.WorldRNG = seed ^ 0x9e3779b9
	// Nobody is playing. The player is the one thing here that would need a
	// policy, so there is not one: they sit at home and the city goes on.
	r := CityReport{Seed: seed, Days: days}

	was := map[string]bool{}
	for _, f := range w.Factions {
		if f.ID != w.PlayerOrganizationID() {
			was[f.ID] = true
		}
	}
	r.Started = len(was)
	owners := map[string]string{}
	for _, l := range core.Locations {
		owners[l.ID] = w.Properties[l.ID].Owner
	}
	wars := map[string]bool{}
	for _, c := range w.PublicConflicts() {
		if c.State == "war" {
			wars[c.Between[0]+"|"+c.Between[1]] = true
		}
	}

	for day := 0; day < days; day++ {
		// A scene waits for somebody to answer it, and nobody is here. Clear it
		// and let the clock go on: this is a measurement of the city, not of
		// what the director would have said about it.
		w.Event = nil
		w.Advance(1440)
		w.Event = nil

		live := map[string]bool{}
		for _, f := range w.Factions {
			if f.ID == w.PlayerOrganizationID() {
				continue
			}
			live[f.ID] = true
			if !was[f.ID] {
				r.Formed++
				was[f.ID] = true
			}
		}
		for id := range was {
			if !live[id] {
				r.Fell++
				delete(was, id)
			}
		}
		now := map[string]bool{}
		for _, c := range w.PublicConflicts() {
			if c.State != "war" {
				continue
			}
			key := c.Between[0] + "|" + c.Between[1]
			now[key] = true
			if !wars[key] {
				r.Wars++
			}
		}
		for key := range wars {
			if !now[key] {
				r.Settled++
			}
		}
		wars = now
		for _, l := range core.Locations {
			if owner := w.Properties[l.ID].Owner; owner != owners[l.ID] {
				r.Changed++
				owners[l.ID] = owner
			}
		}
	}

	r.Ended = 0
	held := map[string]int{}
	total := 0
	for _, f := range w.Factions {
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		r.Ended++
	}
	for _, l := range core.Locations {
		if w.Properties[l.ID].Income <= 0 {
			continue
		}
		owner := w.Properties[l.ID].Owner
		if owner == "" {
			r.Homeless++
			continue
		}
		total++
		held[owner]++
	}
	for _, n := range held {
		if total > 0 && n*100/total > r.Biggest {
			r.Biggest = n * 100 / total
		}
	}
	for _, n := range w.NPCs {
		if n.Dead {
			r.Killed++
		} else {
			r.Living++
		}
	}
	return r
}
