package core

import "fmt"

// Attention is a number attached to one man. Nothing in this city gets harder
// because of what everybody in it has been doing — a campaign of bombings,
// killings and daylight robberies leaves the police exactly as interested in
// week twenty as they were in week one.
//
// Scrutiny is the city's own temperature. It rises with everything the Herald
// reports, whoever did it, and falls slowly through a quiet week. High enough
// and the police stop waiting: raids begin sooner, an understanding costs more,
// and the men in that building will not be seen with anybody.

const (
	// ScrutinyCeiling is as interested as the city ever gets.
	ScrutinyCeiling = 100
	// ScrutinyCool is what a quiet day takes off it.
	ScrutinyCool = 1
	// ScrutinyCrackdown is where the police stop waiting.
	ScrutinyCrackdown = 60
	// ScrutinyRaidRelief is how much sooner raids begin at the ceiling: the
	// threshold falls by this much, scaled by how bad it is.
	ScrutinyRaidRelief = 20
)

// scrutinyWeight is what a story in the paper is worth to the city's
// temperature. Killings and explosions are what a city notices.
// The first set of weights produced no crackdown in a hundred and fifty cities
// at open war over four months: the city files a story every few days and a
// point a day of cooling ate all of it. These are what a paper full of killings
// is actually worth.
var scrutinyWeight = map[string]int{
	"killing": 14, "attack": 12, "war": 8, "police": 6,
	"seizure": 4, "robbery": 3, "collapse": 3, "politics": 2, "business": 0,
	// A column about the price of coal is not a reason for anybody to look
	// harder at anybody. The ordinary edition costs the city nothing.
	"civic": 0,
}

// Scrutiny is how hard the city is looking, from nothing to everything.
func (w *World) Scrutiny() int { return max(0, min(ScrutinyCeiling, w.Attention)) }

// ScrutinyDay moves the city's temperature: up by what the paper carried today,
// down by a point when it carried nothing worth noticing.
func (w *World) ScrutinyDay() {
	gained := 0
	for _, s := range w.News {
		if s.Minute <= w.Minute-1440 || s.Minute > w.Minute {
			continue
		}
		gained += scrutinyWeight[s.Kind]
	}
	before := w.Scrutiny()
	if gained == 0 {
		w.Attention = max(0, w.Attention-ScrutinyCool)
	} else {
		w.Attention = min(ScrutinyCeiling, w.Attention+gained)
	}
	if before < ScrutinyCrackdown && w.Scrutiny() >= ScrutinyCrackdown {
		w.Log("The city has had enough", fmt.Sprintf("It is not about you. Between the killings, the explosions and whatever was in the paper this morning, the police have stopped waiting for a reason. Everything is harder now, for everybody."), "danger")
		w.Report("police", "CITY ORDERS CRACKDOWN ON ORGANIZED CRIME",
			"The commissioner has announced what he describes as an end to tolerance. Officers have been reassigned from other duties and the courts have been asked to sit longer.")
	}
	if before >= ScrutinyCrackdown && w.Scrutiny() < ScrutinyCrackdown {
		w.Log("It has gone quiet", "Whatever the police were doing, they are doing less of it. The city has found something else to worry about.", "personal")
	}
}

// UnderCrackdown reports whether the city has stopped waiting.
func (w *World) UnderCrackdown() bool { return w.Scrutiny() >= ScrutinyCrackdown }

// ScrutinyRaidShift is how much sooner the police come for anybody, because of
// what everybody has been doing.
func (w *World) ScrutinyRaidShift() int {
	return w.Scrutiny() * ScrutinyRaidRelief / ScrutinyCeiling
}

// ScrutinyPremium is what a crackdown adds to the price of anything bought from
// somebody with a career to protect, as a percentage.
func (w *World) ScrutinyPremium() int {
	if !w.UnderCrackdown() {
		return 0
	}
	return 25 + (w.Scrutiny() - ScrutinyCrackdown)
}

// ScrutinyDescription is the city's temperature, for the interface.
func (w *World) ScrutinyDescription() map[string]any {
	state := "ordinary"
	switch {
	case w.Scrutiny() >= 85:
		state = "every door in the city"
	case w.UnderCrackdown():
		state = "a crackdown"
	case w.Scrutiny() >= 30:
		state = "watchful"
	}
	return map[string]any{
		"scrutiny": w.Scrutiny(), "state": state,
		"crackdown": w.UnderCrackdown(), "premium": w.ScrutinyPremium(),
		"raids_sooner": w.ScrutinyRaidShift(),
	}
}
