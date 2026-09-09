package core

import (
	"fmt"
	"strings"
)

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
	// How many times this happened on the day it was filed. One story, with a
	// number on it, rather than the same paragraph printed over and over.
	Count int `json:"count,omitempty"`
	// Set once a brief has been through the director, whether or not the
	// rewrite was accepted. A refusal is as final as an acceptance: asking the
	// same model the same question again gets the same answer, and the paper is
	// not going to sit there re-asking about Tuesday's weather forever.
	Polished bool `json:"polished,omitempty"`
}

// newsCapacity bounds the archive. Measured: a city at war files a little over
// half a story a day, so sixty was about three months — enough for one life and
// nothing at all across several. A player who can read back through the paper
// is reading the history of the city rather than a rolling status, and this
// game is explicitly about a city that remembers, so the archive is worth the
// bytes: at roughly a quarter of a kilobyte a story this is some sixty of them.
const newsCapacity = 240

// Report files a story. Callers pass what the city could observe, which is why
// the wording differs from the private record written to the player's history.
//
// The same story is never printed twice in one issue. Read end to end, a
// campaign's paper carried seven identical paragraphs about a robbery at Saint
// Agnes in a single edition, because seven people were robbed there that day and
// each one filed its own copy. No newspaper does that; it writes one piece
// saying it happened seven times, and that piece is a better story than any one
// of them.
func (w *World) Report(kind, headline, body string) {
	if prior := w.todayStory(kind, headline); prior != nil {
		prior.Count++
		prior.Minute = w.Minute
		prior.Body = runOf(kind, prior.Count, body)
		return
	}
	w.News = append(w.News, Story{
		ID: ID(), Minute: w.Minute, Life: w.Life,
		Headline: headline, Body: body, Kind: kind, Count: 1,
	})
	w.trimNews()
}

// todayStory finds a story already filed today under the same headline.
func (w *World) todayStory(kind, headline string) *Story {
	day := w.Minute / 1440
	for i := len(w.News) - 1; i >= 0; i-- {
		s := &w.News[i]
		if s.Minute/1440 != day || s.Life != w.Life {
			break // the archive is in order; nothing older can be today's
		}
		if s.Kind == kind && s.Headline == headline {
			return s
		}
	}
	return nil
}

// runOf is how a paper writes the same thing happening repeatedly: as one
// incident with a number on it, which reads as a district in trouble rather
// than as a stuck press.
func runOf(kind string, count int, latest string) string {
	if count < 2 {
		return latest
	}
	times := map[int]string{2: "twice", 3: "three times", 4: "four times",
		5: "five times", 6: "six times", 7: "seven times"}[count]
	if times == "" {
		times = fmt.Sprintf("%d times", count)
	}
	switch kind {
	case "robbery":
		return fmt.Sprintf("%s It happened %s in the same day, which residents say is not the usual run of things.", latest, times)
	case "attack":
		return fmt.Sprintf("%s There were %s such incidents before the day was out.", latest, times)
	case "police":
		return fmt.Sprintf("%s Officers were back %s before the day was out.", latest, times)
	}
	return fmt.Sprintf("%s It happened %s that day.", latest, times)
}

// trimNews keeps the archive inside its bound, and drops filler before news.
//
// The city page files two briefs a day whether or not anything happened, so
// without this the weather would quietly evict the killings: at two a day the
// filler alone would fill a 240-story archive in four months and a player
// reading back through the paper would find nothing but prices. What the
// archive is for is the history of the city, so a brief about a clear day gives
// up its place before a man who was shot does.
func (w *World) trimNews() {
	for len(w.News) > newsCapacity {
		oldest := -1
		for i, s := range w.News {
			if s.Kind == "civic" {
				oldest = i
				break
			}
		}
		if oldest < 0 {
			w.News = w.News[len(w.News)-newsCapacity:]
			return
		}
		w.News = append(w.News[:oldest], w.News[oldest+1:]...)
	}
}

// story is one piece of copy, set the way the paper would set it.
func (w *World) story(s Story) map[string]any {
	return map[string]any{
		"id": s.ID, "headline": s.Headline, "body": s.Body,
		"kind": s.Kind, "minute": s.Minute, "subject": w.SubjectOf(s),
		"standfirst": w.standfirst(s), "byline": deskFor(s.Kind), "dateline": Dateline(s.Minute),
		"day": s.Minute/1440 + 1, "time": fmt.Sprintf("%02d:%02d", s.Minute%1440/60, s.Minute%60),
		"weight": Newsworthiness(s.Kind),
	}
}

// Edition is today's paper, most recent first. Kept because plenty of things
// only ever want the front of the paper.
func (w *World) Edition() []map[string]any {
	out := []map[string]any{}
	for i := len(w.News) - 1; i >= 0; i-- {
		if s := w.News[i]; s.Life == w.Life {
			out = append(out, w.story(s))
		}
	}
	return out
}

// Editions is the paper as an archive rather than a rolling list: one issue a
// day, newest first, every life the save still remembers. A previous
// protagonist's era is readable, which is the whole premise of a city that
// remembers — the player who comes next inherits the news as well as the
// streets.
//
// Within an issue the lead is the biggest story of that day rather than the
// latest. A shop changing hands does not lead over a killing because it
// happened at four in the afternoon.
//
// Ranked by Newsworthiness rather than by what a story costs in police
// attention. Those were the same map, and it does not cover half the kinds the
// paper files: an arrest, an attempt on somebody's life and a family splitting
// all scored zero, so none of them could ever lead an edition.
func (w *World) Editions() []map[string]any {
	type issue struct {
		day, life int
		stories   []Story
	}
	order := []*issue{}
	byDay := map[int]*issue{}
	for _, s := range w.News {
		day := s.Minute/1440 + 1
		key := s.Life*100000 + day
		if byDay[key] == nil {
			byDay[key] = &issue{day: day, life: s.Life}
			order = append(order, byDay[key])
		}
		byDay[key].stories = append(byDay[key].stories, s)
	}
	out := []map[string]any{}
	for i := len(order) - 1; i >= 0; i-- {
		in := order[i]
		// The biggest story of the day leads, then the rest as they came.
		sorted := append([]Story{}, in.stories...)
		for a := range sorted {
			for b := a + 1; b < len(sorted); b++ {
				if Newsworthiness(sorted[b].Kind) > Newsworthiness(sorted[a].Kind) {
					sorted[a], sorted[b] = sorted[b], sorted[a]
				}
			}
		}
		stories := []map[string]any{}
		for _, s := range sorted {
			stories = append(stories, w.story(s))
		}
		out = append(out, map[string]any{
			"day": in.day, "life": in.life, "dateline": Dateline((in.day - 1) * 1440),
			"stories": stories, "count": len(stories),
			"current": in.life == w.Life && in.day == w.Minute/1440+1,
			"mine":    in.life == w.Life,
		})
	}
	return out
}

// unattributed is how the paper describes a crime nobody has been charged with.
// The player reads their own work here with no name attached, which is exactly
// what the rest of the city sees.
func (w *World) unattributed(place string, what string) string {
	// Callers that already said the police have nothing get left alone. The
	// robbery copy read "Police have asked anybody who saw it to come forward.
	// Police have made no arrest and are appealing for anyone who saw the
	// incident at Saint Agnes." — the same sentence twice, in every edition.
	if strings.Contains(what, "Police") || strings.Contains(what, "police") {
		return what
	}
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
	case "civic":
		return "By our municipal correspondent"
	case "obituary":
		return "Obituaries"
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
	case "civic":
		return "Compiled from the district returns. The Herald prints them whether or not anything else happened."
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
