package main

import (
	"blackledger/core"
	"blackledger/sim"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

// livingWorldHorizon is how many game days a campaign has to cover before the
// city measures say anything. Families escalate over weeks; at a fortnight the
// answer is always "the city did nothing", and that is a fact about the run
// rather than about the city.
const livingWorldHorizon = 20

func main() {
	runs := flag.Int("runs", 100, "campaigns per strategy")
	steps := flag.Int("steps", 200, "maximum commands per campaign")
	first := flag.Uint("seed", 1, "first simulation seed; later runs use a Weyl stride")
	profiles := flag.String("strategies", "worker,investor,defiant,reckless,thief,smuggler,racketeer,publican", "comma-separated player policies")
	director := flag.String("director", "authored", "authored, fixture or replay (no model calls)")
	corpusPath := flag.String("corpus", "", "JSON proposal array required for replay")
	trace := flag.Bool("trace", false, "include each pre-command public state and command")
	// The city with nobody in it. A campaign follows one protagonist for a
	// median of a few days, which is why no report this harness has ever
	// printed showed an organization ending. They end; nobody was watching
	// long enough. These runs have no player policy at all.
	cities := flag.Int("cities", 12, "cities to run with nobody playing them")
	season := flag.Int("season", 60, "days to run each of those cities for")
	flag.Parse()
	if *runs < 1 || *runs > 10000 || *steps < 1 || *steps > 10000 || *first > 4294967295 {
		fmt.Fprintln(os.Stderr, "invalid run/step/seed bounds")
		os.Exit(2)
	}
	if *director != "authored" && *director != "fixture" && *director != "replay" {
		fmt.Fprintln(os.Stderr, "director must be authored, fixture or replay")
		os.Exit(2)
	}
	var corpus []core.Proposal
	corpusHash := ""
	if (*director == "replay") != (*corpusPath != "") {
		fmt.Fprintln(os.Stderr, "replay requires -corpus; other modes do not accept it")
		os.Exit(2)
	}
	if *director == "replay" {
		f, err := os.Open(*corpusPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		corpus, corpusHash, err = sim.ReadCorpus(f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	strategies := strings.Split(*profiles, ",")
	for _, p := range strategies {
		if p != "worker" && p != "investor" && p != "reckless" && p != "defiant" && p != "diplomat" && p != "thief" && p != "smuggler" && p != "racketeer" && p != "publican" {
			fmt.Fprintln(os.Stderr, "unknown strategy", p)
			os.Exit(2)
		}
	}
	start := time.Now()
	reports := []sim.Report{}
	failed := false
	for _, p := range strategies {
		for i := 0; i < *runs; i++ {
			r := sim.RunRecorded(uint32(*first)+uint32(i)*0x9e3779b9, p, *director, *steps, *trace, corpus)
			reports = append(reports, r)
			failed = failed || r.Error != ""
		}
	}
	summaries := map[string]any{}
	// Which strategies did not last long enough for the city block to mean
	// anything, named in the output so nobody reads a short campaign as a quiet
	// city again.
	short := []string{}
	for _, p := range strategies {
		deaths, errors := 0, 0
		cash := []int{}
		minutes := []int{}
		milestones := map[string][]int{}
		for _, r := range reports {
			if r.Strategy != p {
				continue
			}
			if !r.Alive {
				deaths++
			}
			if r.Error != "" {
				errors++
			}
			cash = append(cash, r.Cash)
			minutes = append(minutes, r.Minutes)
			for k, v := range r.Milestones {
				milestones[k] = append(milestones[k], v)
			}
		}
		sort.Ints(cash)
		sort.Ints(minutes)
		// What the city did on its own. Reported as totals across the runs and
		// as how many runs saw any of it at all, because a living world that
		// produces one war in a hundred campaigns is not a living world — the
		// total alone would hide that.
		city := map[string]any{}
		{
			type measure struct {
				total, runsWith int
			}
			m := map[string]*measure{}
			add := func(name string, n int) {
				if m[name] == nil {
					m[name] = &measure{}
				}
				m[name].total += n
				if n > 0 {
					m[name].runsWith++
				}
			}
			for _, r := range reports {
				if r.Strategy != p {
					continue
				}
				add("wars_started", r.World.WarsStarted)
				add("wars_elsewhere", r.World.WarsElsewhere)
				add("holdings_changed_hands", r.World.HoldingsChangedHands)
				add("factions_created", r.World.FactionsCreated)
				add("factions_destroyed", r.World.FactionsDestroyed)
				add("hurt_during_others_war", r.World.HurtDuringOthersWar)
			}
			for name, v := range m {
				city[name] = map[string]any{"total": v.total, "runs_with_any": v.runsWith}
			}
		}
		ms := map[string]any{}
		for k, values := range milestones {
			sort.Ints(values)
			ms[k] = map[string]any{"reached": len(values), "median_commands_when_reached": values[len(values)/2]}
		}
		// Days, not minutes. The living-world measures are about what the city
		// does over weeks, and a summary reported in minutes let a 12-day
		// median pass for a campaign — which produced a confident and wrong
		// conclusion that the city was dormant. It was not; the runs were short.
		span := minutes[len(minutes)/2]
		days := float64(span) / 1440
		// What the city took, as against what the player kept. A policy whose
		// money is taken rather than whose life is looked identical to a safe
		// one in this report.
		heat, seizures, hurt, informed := 0, 0, 0, 0
		lows := []int{}
		for _, r := range reports {
			if r.Strategy != p {
				continue
			}
			heat += r.Heat
			seizures += r.Seizures
			informed += r.Informed
			lows = append(lows, r.Lowest)
			if r.Health < 100 {
				hurt++
			}
		}
		summaries[p] = map[string]any{"runs": *runs, "deaths": deaths, "errors": errors,
			"mean_heat": heat / max(1, *runs), "seizures": seizures,
			"informed": informed, "median_lowest_health": median(lows),
			"median_final_cash": cash[len(cash)/2], "median_game_minutes": span,
			"median_game_days": math.Round(days*10) / 10, "milestones": ms, "city": city,
			"city_measures_meaningful": days >= livingWorldHorizon,
		}
		if days < livingWorldHorizon {
			short = append(short, fmt.Sprintf("%s (%.1f days)", p, days))
		}
	}
	// What a season does to a city that nobody is playing.
	out := map[string]any{"city_alone": seasonReport(*cities, *season), "corpus_sha256": corpusHash, "corpus_proposals": len(corpus), "elapsed_seconds": time.Since(start).Seconds(), "director": *director, "max_commands": *steps, "summary": summaries, "campaigns": reports}
	if len(short) > 0 {
		out["warning"] = fmt.Sprintf(
			"the city measures are not meaningful for %s: a campaign has to run past about %d days "+
				"before families have time to escalate, split or fall, and these ended sooner. "+
				"Raise -steps, or read only the strategies that survived.",
			strings.Join(short, ", "), livingWorldHorizon)
		fmt.Fprintln(os.Stderr, out["warning"])
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}

// seasonReport runs a handful of cities with nobody in them and totals what
// happened. It is the only measure in this harness that can see a family fall,
// because a campaign does not last long enough for one to.
func seasonReport(cities, days int) map[string]any {
	if cities <= 0 || days <= 0 {
		return map[string]any{"cities": 0}
	}
	total := map[string]int{}
	worst, runs := 0, []sim.CityReport{}
	for i := 0; i < cities; i++ {
		r := sim.City(uint32(i+1)*2654435761, days)
		runs = append(runs, r)
		total["organizations_formed"] += r.Formed
		total["organizations_fell"] += r.Fell
		total["wars_started"] += r.Wars
		total["wars_settled"] += r.Settled
		total["holdings_changed_hands"] += r.Changed
		total["people_killed"] += r.Killed
		if r.Biggest > worst {
			worst = r.Biggest
		}
	}
	return map[string]any{
		"cities": cities, "days": days, "totals": total,
		"biggest_share_anywhere": worst, "runs": runs,
	}
}

// median is the middle of a set of run figures, which is what a summary should
// report: one campaign that went badly is not the shape of a hundred.
func median(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	sort.Ints(xs)
	return xs[len(xs)/2]
}
