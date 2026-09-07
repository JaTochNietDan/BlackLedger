package main

import (
	"blackledger/sim"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	runs := flag.Int("runs", 100, "campaigns per strategy")
	steps := flag.Int("steps", 200, "maximum commands per campaign")
	first := flag.Uint("seed", 1, "first simulation seed; later runs use a Weyl stride")
	profiles := flag.String("strategies", "worker,investor,defiant,reckless", "comma-separated player policies")
	director := flag.String("director", "authored", "authored or fixture (no model calls)")
	trace := flag.Bool("trace", false, "include each pre-command public state and command")
	flag.Parse()
	if *runs < 1 || *runs > 10000 || *steps < 1 || *steps > 10000 || *first > 4294967295 {
		fmt.Fprintln(os.Stderr, "invalid run/step/seed bounds")
		os.Exit(2)
	}
	if *director != "authored" && *director != "fixture" {
		fmt.Fprintln(os.Stderr, "director must be authored or fixture")
		os.Exit(2)
	}
	strategies := strings.Split(*profiles, ",")
	for _, p := range strategies {
		if p != "worker" && p != "investor" && p != "reckless" && p != "defiant" {
			fmt.Fprintln(os.Stderr, "unknown strategy", p)
			os.Exit(2)
		}
	}
	start := time.Now()
	reports := []sim.Report{}
	failed := false
	for _, p := range strategies {
		for i := 0; i < *runs; i++ {
			r := sim.Run(uint32(*first)+uint32(i)*0x9e3779b9, p, *director, *steps, *trace)
			reports = append(reports, r)
			failed = failed || r.Error != ""
		}
	}
	summaries := map[string]any{}
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
		ms := map[string]any{}
		for k, values := range milestones {
			sort.Ints(values)
			ms[k] = map[string]any{"reached": len(values), "median_commands_when_reached": values[len(values)/2]}
		}
		summaries[p] = map[string]any{"runs": *runs, "deaths": deaths, "errors": errors, "median_final_cash": cash[len(cash)/2], "median_game_minutes": minutes[len(minutes)/2], "milestones": ms}
	}
	out := map[string]any{"elapsed_seconds": time.Since(start).Seconds(), "director": *director, "max_commands": *steps, "summary": summaries, "campaigns": reports}
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
