package core

import "fmt"

// The Bellwether Herald prints what the city can see. It is assembled at the
// moment events are committed, not summarised afterwards, so a headline always
// refers to something that actually happened.
//
// The paper knows only public things. It never prints a plot that has not
// happened yet, never names who paid for a killing, and reports the player's
// own crimes the way a newspaper would: as an event with no suspect. Reading it
// is free and advances nothing.

type Story struct {
	ID       string `json:"id"`
	Minute   int    `json:"minute"`
	Life     int    `json:"life"`
	Headline string `json:"headline"`
	Body     string `json:"body"`
	Kind     string `json:"kind"`
}

// newsCapacity bounds the archive. A campaign runs for weeks of game time and
// the paper is a recent record, not a complete one.
const newsCapacity = 60

// Report files a story. Callers pass what the city could observe, which is why
// the wording differs from the private record written to the player's history.
func (w *World) Report(kind, headline, body string) {
	w.News = append(w.News, Story{
		ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: headline, Body: body, Kind: kind,
	})
	if len(w.News) > newsCapacity {
		w.News = w.News[len(w.News)-newsCapacity:]
	}
}

// Edition is the paper as it stands, most recent first, with the day each story
// ran. Only the current life's stories are printed: a new person picks up a
// paper that has been running the whole time, but the archive they can read
// starts when they arrive.
func (w *World) Edition() []map[string]any {
	out := []map[string]any{}
	for i := len(w.News) - 1; i >= 0; i-- {
		s := w.News[i]
		if s.Life != w.Life {
			continue
		}
		out = append(out, map[string]any{
			"id": s.ID, "headline": s.Headline, "body": s.Body,
			"kind": s.Kind, "minute": s.Minute, "subject": w.SubjectOf(s),
			"standfirst": w.standfirst(s), "byline": deskFor(s.Kind), "dateline": Dateline(s.Minute),
			"day": s.Minute/1440 + 1, "time": fmt.Sprintf("%02d:%02d", s.Minute%1440/60, s.Minute%60),
		})
	}
	return out
}

// unattributed is how the paper describes a crime nobody has been charged with.
// The player reads their own work here with no name attached, which is exactly
// what the rest of the city sees.
func (w *World) unattributed(place string, what string) string {
	return fmt.Sprintf("%s Police have made no arrest and are appealing for anyone who saw the incident at %s.", what, place)
}

// hasNewsKind reports whether the paper has carried a story of a kind this
// life, which is what a test can ask without reaching into the wording.
func (w *World) hasNewsKind(kind string) bool {
	for _, s := range w.News {
		if s.Life == w.Life && s.Kind == kind {
			return true
		}
	}
	return false
}

// A headline on its own reads like a log line. A paper of the period gave every
// story a standfirst under the headline, a byline that told you which desk it
// came from, and a date at the top of the page. None of that is new truth — it
// is the truth already committed, set the way the city would have set it.

// pressEpoch is the weekday the first day of a campaign falls on, and the date
// it carries. Bellwether keeps its own calendar; what matters is that it reads
// like a paper somebody left on a table.
var (
	weekdays   = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	monthNames = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
)

// Dateline is the line under the masthead: what day of the week it is, and the
// date, derived from the same clock everything else uses.
func Dateline(minute int) string {
	day := minute / 1440
	// The third of March, a Tuesday, is day one.
	weekday := weekdays[(day+1)%7]
	month, date := 2, 3+day
	for month < len(monthNames) && date > daysIn(month) {
		date -= daysIn(month)
		month++
	}
	return fmt.Sprintf("%s, %s %d, 1953", weekday, monthNames[min(month, len(monthNames)-1)], date)
}

func daysIn(month int) int {
	switch month {
	case 1:
		return 28
	case 3, 5, 8, 10:
		return 30
	}
	return 31
}

// deskFor is which desk filed it, in the way a paper of the period signed its
// copy: rarely a name, usually a department.
func deskFor(kind string) string {
	switch kind {
	case "killing", "attack", "robbery":
		return "By our crime reporter"
	case "police", "seizure":
		return "By our police correspondent"
	case "politics", "war", "split":
		return "By our city hall correspondent"
	case "arrest", "attempt", "recovery":
		return "By our courts correspondent"
	case "business":
		return "By our commercial editor"
	case "collapse":
		return "By our municipal correspondent"
	}
	return "Staff report"
}

// standfirst is the italic line under a headline: what the story is really
// about, in the register of a paper that has to sell itself on a newsstand and
// cannot say what everybody knows.
func (w *World) standfirst(s Story) string {
	subject := w.SubjectOf(s)
	where := "the district"
	if subject.Kind != "city" && subject.Name != "" {
		where = subject.Name
	}
	switch s.Kind {
	case "killing":
		return fmt.Sprintf("Police at %s say they are pursuing several lines of inquiry. Neighbours describe a quiet street.", where)
	case "attack":
		return fmt.Sprintf("Damage at %s is described as extensive. No arrests have been made and none appear imminent.", where)
	case "robbery":
		return fmt.Sprintf("The proprietor of %s declined to be photographed. Officers have asked witnesses to come forward.", where)
	case "police":
		return "The department describes the operation as part of a continuing investigation. It would not be drawn on what it expects to find."
	case "seizure":
		return fmt.Sprintf("Ownership of %s has been in question for some time. Nobody at the premises would give a name.", where)
	case "war":
		return "Interests on both sides deny that anything is happening at all. Two funerals this month suggest otherwise."
	case "politics":
		return "Nothing was announced and nothing was denied. The arrangement is understood to have been settled some days ago."
	case "business":
		return fmt.Sprintf("Trade at %s is said to be steady. Rivals in the district were not available for comment.", where)
	case "collapse":
		return "The city says the matter is closed. Those who worked there have not been told where to apply."
	case "arrest":
		return "The accused was not represented. The department would not say how the name came to it."
	case "attempt":
		return fmt.Sprintf("Whoever went to %s left in a hurry and left nothing behind. Police are not treating it as random.", where)
	case "split":
		return "Those who left were not asked to explain themselves. Those who stayed were not asked either."
	case "recovery":
		return fmt.Sprintf("Property has been returned to %s. How it left is not a question anybody is putting in writing.", where)
	}
	return "The Herald understands the position remains unchanged."
}
