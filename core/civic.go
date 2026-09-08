package core

import "fmt"

// The Herald prints what the city can see, and it turned out the city could
// only see violence. Measured on a real campaign: twenty-one days of play
// produced one story, because every route into the paper is a killing, a raid,
// a robbery or a war, and that campaign had none. Run the same save forward
// into a war and it files fourteen stories in thirty days.
//
// A paper that prints nothing for three weeks is not a paper. Real ones fill
// the space with trade, civic business and whatever the town is talking about,
// and they do it every single day. So does this one now — out of state the city
// already holds, never invented, and worth nothing to the city's temperature,
// because a column about the price of coal is not a reason for the police to
// look harder at anybody.

// CivicDay files the ordinary edition: what the paper would carry on a day when
// nothing happened. It runs after the day's real news, and only fills the space
// that day left empty, so a morning with a killing in it is not padded with the
// price of coal.
func (w *World) CivicDay() {
	filed := 0
	for _, s := range w.News {
		if s.Life == w.Life && w.Minute-s.Minute < 1440 {
			filed++
		}
	}
	if filed >= 2 {
		return // a day with real news in it does not need filling
	}
	for _, item := range w.civicItems() {
		w.Report("civic", item[0], item[1])
		filed++
		if filed >= 2 {
			return
		}
	}
}

// civicItems is everything true about this city today that a local paper would
// print, in the order it would lead with them. Every one is read out of state
// rather than made up, so the paper can be trusted the way the rest of the game
// is trusted.
func (w *World) civicItems() [][2]string {
	out := [][2]string{}
	day := w.Minute / 1440

	// Trade: the market moves every day and a commercial paper leads on it.
	if len(w.Goods) > 0 {
		g := w.Goods[day%len(w.Goods)]
		direction, verb := "STEADY", "held its level"
		if g.Price > g.Base*11/10 {
			direction, verb = "RISES", "climbed above what it usually fetches"
		} else if g.Price < g.Base*9/10 {
			direction, verb = "FALLS", "slipped below its usual price"
		}
		out = append(out, [2]string{
			fmt.Sprintf("PRICE OF %s %s", upper(g.Name), direction),
			fmt.Sprintf("The price of %s %s this week, at $%d the %s against $%d in an ordinary month. Traders at Mercer Exchange offered no explanation.", g.Name, verb, g.Price, g.Unit, g.Base),
		})
	}

	// Whoever holds a job in this city holds it in public.
	if roles := w.RoleDescription(); len(roles) > 0 {
		r := roles[day%len(roles)]
		out = append(out, [2]string{
			fmt.Sprintf("%s CONTINUES IN POST", upper(fmt.Sprint(r["title"]))),
			fmt.Sprintf("%s remains %s of Bellwether. The appointment has drawn no comment from the council, which is itself a kind of comment.", r["name"], r["title"]),
		})
	}

	// The town's own size and shape, which changes as people arrive and die.
	if p := w.PopulationSummary(); p != nil {
		out = append(out, [2]string{
			"DISTRICT POPULATION HOLDS AT " + fmt.Sprint(p["living"]),
			fmt.Sprintf("The district counts %v people at work or looking for it, of whom %v answer to one of the organizations and %v to nobody at all. The figure is compiled from the rating rolls.", p["living"], p["organized"], p["street"]),
		})
	}

	// How hard everybody is being looked at, in the words a paper would use.
	temperature := "quiet"
	switch {
	case w.UnderCrackdown():
		temperature = "under a crackdown nobody will call one"
	case w.Scrutiny() > 30:
		temperature = "watchful"
	}
	out = append(out, [2]string{
		"POLICE REPORT DISTRICT " + upper(temperature),
		fmt.Sprintf("The department describes the past week in the district as %s. Officers were reassigned from other duties, and the commissioner declined to say from where.", temperature),
	})

	// A business the city can see, doing better or worse than it was.
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || prop.Income <= 0 {
			continue
		}
		if prop.Condition < 60 {
			out = append(out, [2]string{
				"REPAIRS AWAITED AT " + upper(l.Name),
				fmt.Sprintf("%s is trading in visibly poor repair, at %d%% of the condition its rateable value assumes. Its proprietor was not available.", l.Name, prop.Condition),
			})
			break
		}
	}
	return out
}
