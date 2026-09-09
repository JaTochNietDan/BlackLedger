package core

import "fmt"

// The city page.
//
// The Herald only ever printed things that had just happened to somebody: a
// killing, a robbery, a business changing hands. A real paper is mostly not
// that. It is weather, and prices, and the number of people in the place, and
// a paragraph about premises standing empty on a street nobody walks down —
// copy that exists because a paper comes out whether or not anything happened.
// Without it the city only speaks when it is being violent.
//
// Every one of these is composed from something the world already knows, and
// only from things the city could see. The paper does not learn anything here
// that it could not learn by looking out of the window: what the sky is doing,
// what a crate of moonshine costs, how hard the police are looking, how many
// people are on the streets. Nothing here reveals who arranged anything.
//
// They are filed on the day turn like any other story, which means they are
// archived with that day's issue and a back number still carries the weather it
// was printed with. Composing them on demand instead would have yesterday's
// paper reporting today's rain.

// brief is one piece of city-page copy and the fact that justified it.
type brief struct {
	headline string
	body     string
}

// cityPage is everything the paper could run today, in a fixed order. Which of
// them actually runs is chosen below.
func (w *World) cityPage() []brief {
	out := []brief{}
	day := w.Minute / 1440
	sky := w.Sky()

	// The weather, which is the one thing every newspaper has always printed.
	switch sky.Kind {
	case "rain":
		out = append(out, brief{"RAIN ACROSS THE DISTRICT",
			"Rain through the day and into the evening. The gutters on the lower streets are carrying more than they were built for, and shopkeepers have been advised that the drains will be looked at when there is money to look at them."})
	case "fog":
		out = append(out, brief{"FOG OFF THE WATER",
			"Fog again off the water, thick enough by evening that the far side of a street is a rumour. Drivers are asked to go carefully. Those with business after dark are asked to consider whether they have business after dark."})
	case "overcast":
		out = append(out, brief{"NO CHANGE IN THE WEATHER",
			"Cloud over the city and no sign of it lifting. The forecast, for what it has been worth lately, says the same again tomorrow."})
	default:
		out = append(out, brief{"A CLEAR DAY",
			"Clear over Bellwether, and warm enough by afternoon that the benches on the front were taken by eleven. It will not last."})
	}
	if sky.Kind != "rain" && sky.Wet > 0 {
		out = append(out, brief{"THE STREETS STILL WET",
			"The rain has stopped and the streets have not dried. Standing water at the low end of the harbour road again, which residents there point out is where it always is."})
	}

	// The market, which the paper reports as prices rather than as contraband.
	if g, ok := w.dearest(); ok {
		if g.Price > g.Base {
			out = append(out, brief{"PRICES UP AT THE DOCKS",
				fmt.Sprintf("%s is fetching more than it did, at around $%d the %s against $%d not long ago. Those who deal in it say supply; those who buy it say something else.",
					g.Name, g.Price, g.Unit, g.Base)})
		} else {
			out = append(out, brief{"A GLUT ON THE WATERFRONT",
				fmt.Sprintf("%s has fallen to about $%d the %s, down from $%d. More of it has come into the city than the city has use for, and it is being sold accordingly.",
					g.Name, g.Price, g.Unit, g.Base)})
		}
	}

	// How hard the police are looking, which is public: the city can see the
	// cars on the corners.
	switch w.scrutinyWord() {
	case "quiet":
		out = append(out, brief{"A QUIET WEEK FOR THE POLICE",
			"The department reports a week without incident of any note, and has said so at some length. Officers were seen on the harbour road on Tuesday, which is where officers are usually seen."})
	default:
		out = append(out, brief{"MORE OFFICERS ON THE STREETS",
			"Additional officers have been assigned to the district. The department declines to say for how long, or in response to what, and would not be drawn on either."})
	}

	// The size of the place, and how much of it answers to somebody.
	living, organized := len(w.People()), 0
	for i := range w.Factions {
		organized += len(w.Members(w.Factions[i].ID))
	}
	if living > 0 {
		out = append(out, brief{"THE CITY COUNTED",
			fmt.Sprintf("Some %d people are living in the district by the latest count, of whom %d are understood to be in the employ of one of the families. The remainder work for a living or are between arrangements.",
				living, organized)})
	}

	// Premises standing without anybody answering for them, which anybody
	// walking past can see for themselves.
	if empty, ok := w.standingEmpty(); ok {
		out = append(out, brief{"PREMISES STANDING EMPTY",
			fmt.Sprintf("%s has been standing without anybody answering for it. The building is sound. What it wants is somebody to open it, and the terms are said to be reasonable to anybody who asks about them.", empty)})
	}

	// And the day of the week, because a paper that never mentions Sunday is
	// not a paper.
	if day%7 == 6 {
		out = append(out, brief{"SUNDAY IN BELLWETHER",
			"Saint Agnes will hold its usual services. The market on the harbour road will not open. Those with business that cannot wait for Monday will find the usual people in the usual places, as they always have."})
	}
	return out
}

// dearest is whichever good has moved furthest from what it usually costs,
// which is the one a commercial editor would write about.
func (w *World) dearest() (Good, bool) {
	best, gap := Good{}, 0
	for _, g := range w.Goods {
		if g.Base == 0 {
			continue
		}
		if d := abs(g.Price - g.Base); d > gap {
			best, gap = g, d
		}
	}
	// A few dollars either way is not a story.
	return best, gap*100/max(1, best.Base) >= 12
}

// standingEmpty is premises nobody is answering for: visible from the street,
// and not the player's business to explain.
func (w *World) standingEmpty() (string, bool) {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil || w.Own(l.ID) || w.faction(prop.Owner) != nil {
			continue
		}
		if prop.Owner == "independent" || hasPrefix(prop.Owner, "former:") {
			if prop.Staff == 0 {
				return l.Name, true
			}
		}
	}
	return "", false
}

// CityPageDay files the day's filler. Two pieces: enough that the paper has
// something in it on a day when nothing happened, few enough that the news is
// still the news.
func (w *World) CityPageDay() {
	page := w.cityPage()
	if len(page) == 0 {
		return
	}
	// Deterministic from the day, so the same day of the same campaign always
	// prints the same page, and consecutive days do not run the same two.
	h := skySeed(w.ID) ^ uint32(w.Minute/1440)*2246822519
	h ^= h >> 13
	for n := 0; n < 2 && n < len(page); n++ {
		pick := page[int((h>>uint(n*8))%uint32(len(page)))]
		// Never the same brief twice in one issue.
		if n == 1 && pick.headline == page[int(h%uint32(len(page)))].headline {
			pick = page[(int(h%uint32(len(page)))+1)%len(page)]
		}
		w.Report("civic", pick.headline, pick.body)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
