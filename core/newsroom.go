package core

import "fmt"

// The Bellwether Herald was the one thing in this city nobody could touch. It
// reported everything and answered to no one: the player read their own work in
// it with no name attached and could do nothing about the fact that it was
// there at all. Meanwhile the city's whole temperature is derived from what the
// paper carried — so the one institution nobody could reach was the one that
// decided how hard everything else was.
//
// A man on the top floor changes that. Stories about you can be pulled before
// they run. Stories about anybody else can be arranged. Both are permanent
// facts about you held by somebody who is paid weekly, which is the least
// secure way there is to hold a secret.

const (
	// HeraldPlace is where this is arranged.
	HeraldPlace = "herald"
	// SpikeCost is what pulling one story costs, and SpikeMinutes the evening
	// it takes.
	SpikeCost, SpikeMinutes = 260, 45
	// SpikeExposure is the chance, each time, that somebody in that building
	// notices what is being done and the arrangement ends.
	SpikeExposure = .14
	// SmearCost is what running a story about somebody else costs.
	SmearCost, SmearMinutes = 340, 60
	// SmearCustom is the trade each of their premises loses, and SmearPower
	// what the organization loses.
	SmearCustom, SmearPower = 14, 4
	// SmearTrace is the chance they work out who paid for it, and
	// SmearGoodwill what it costs when they do.
	SmearTrace, SmearGoodwill = .22, 25
	// PuffCost is what a favourable piece about yourself costs, and PuffRespect
	// and PuffHeat what it is worth.
	PuffCost, PuffMinutes = 180, 45
	PuffRespect, PuffHeat = 6, 5
	// SmearFloor is as low as anything printed can push a place's trade. Below
	// it, whoever is still going there is going for reasons a story will not
	// change.
	SmearFloor = 25
	// PressCooldown is how long between anything the player puts in the paper.
	// A newspaper that carried the same man's arrangements twice in a week
	// would not be a newspaper.
	PressCooldown = 2880
)

// TheEditor reports whether there is a man at the Herald taking the player's
// money, and whether he is currently in a position to be worth it.
func (w *World) TheEditor() bool { return w.Retained("editor") && !w.Outbid("editor") }

// Spikeable is what could be pulled: stories from the last day about premises
// of the player's, or about the police, which in this city usually means them.
// The paper is the city's memory, so a story that never runs is a thing that,
// as far as everybody else is concerned, did not happen.
func (w *World) Spikeable() []Story {
	out := []Story{}
	for _, s := range w.News {
		if s.Life != w.Life || w.Minute-s.Minute > 1440 {
			continue
		}
		subject := w.SubjectOf(s)
		if subject.Kind == "place" && w.Own(subject.ID) {
			out = append(out, s)
			continue
		}
		if s.Kind == "police" || s.Kind == "seizure" {
			out = append(out, s)
		}
	}
	return out
}

// SpikeReadiness explains why nothing can be pulled, or returns "".
func (w *World) SpikeReadiness() string {
	if !w.TheEditor() {
		return "Nobody at that paper takes your calls"
	}
	if w.Player.Location != HeraldPlace {
		return "That conversation happens at the paper"
	}
	if len(w.Spikeable()) == 0 {
		return "There is nothing in today's paper worth pulling"
	}
	if w.Player.Cash < SpikeCost {
		return "Not enough cash"
	}
	return ""
}

// Spike pulls the worst of today's stories. The city's interest in it goes with
// it, because the city's interest was never anything but what it read.
func (w *World) Spike() error {
	if reason := w.SpikeReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(SpikeCost); err != nil {
		return err
	}
	// The worst of them, by what it is doing to the city's temperature.
	worst := w.Spikeable()[0]
	for _, s := range w.Spikeable() {
		if scrutinyWeight[s.Kind] > scrutinyWeight[worst.Kind] {
			worst = s
		}
	}
	w.dropStory(worst.ID)
	w.Attention = max(0, w.Attention-scrutinyWeight[worst.Kind])
	w.Player.Heat = max(0, w.Player.Heat-4)
	w.Log("It does not run", fmt.Sprintf("%q was set and is not in tomorrow's paper. As far as this city is concerned it did not happen, and one more person knows it did.", worst.Headline), "personal")

	// Every one of these is a thing somebody in that building knows.
	if w.Random() < SpikeExposure {
		w.EndRetainer("editor")
		w.Log("Somebody at the paper talked", "The arrangement at the Herald is over. Nobody says why, and nobody at that desk will take a call from you again.", "danger")
	}
	return nil
}

// SmearReadiness explains why a story cannot be arranged about somebody, or
// returns "".
func (w *World) SmearReadiness(id string) string {
	if !w.TheEditor() {
		return "Nobody at that paper takes your calls"
	}
	if w.Player.Location != HeraldPlace {
		return "That conversation happens at the paper"
	}
	f := w.faction(id)
	if f == nil || id == w.PlayerOrganizationID() {
		return "There is no such organization"
	}
	if len(w.FamilyHoldings(id)) == 0 {
		return f.Name + " has nothing left for a story to be about"
	}
	if w.pressedRecently() {
		return "The paper carried something of yours too recently"
	}
	if w.Player.Cash < SmearCost {
		return "Not enough cash"
	}
	return ""
}

// Smear runs a story about a rival. It costs them trade and standing, it costs
// the whole city a hotter week — a paper full of crime is a paper full of
// crime, whoever it is about — and it can be traced back.
func (w *World) Smear(id string) error {
	if reason := w.SmearReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(SmearCost); err != nil {
		return err
	}
	f := w.faction(id)
	w.Player.LastPress = w.Minute
	holdings := w.FamilyHoldings(id)
	// Print enough and a place empties out, but only so far. A newspaper can
	// talk a business down and cannot talk it out of existence: past this the
	// people still going are going for reasons no story will change.
	moved := false
	for _, place := range holdings {
		if w.Custom(place) <= SmearFloor {
			continue
		}
		w.ShiftCustom(place, "What the Herald printed about who runs it", -SmearCustom)
		moved = true
	}
	if moved {
		f.Power = max(10, f.Power-SmearPower)
	}
	place, _ := PlaceByID(holdings[0])
	w.Report("politics", "POLICE EXAMINE ACCOUNTS AT "+upper(place.Name),
		fmt.Sprintf("The Herald understands that premises associated with %s are the subject of an inquiry into their books. %s has not responded to the newspaper's questions.", f.Name, f.Leader))
	w.Log("It runs on the third page", fmt.Sprintf("%s wakes up to a story about their accounts. Every one of their places is quieter today, and the city as a whole is a little harder for everybody, including you.", f.Name), "politics")

	if w.Random() < SmearTrace {
		f.Goodwill = max(-100, f.Goodwill-SmearGoodwill)
		w.Antagonize(id, w.PlayerOrganizationID(), 10)
		w.Log("They know where it came from", fmt.Sprintf("Somebody at the Herald says a name and it is yours. %s is not going to write a letter about it.", f.Name), "danger")
	}
	return nil
}

// PuffReadiness explains why nothing favourable can be arranged, or returns "".
func (w *World) PuffReadiness() string {
	if !w.TheEditor() {
		return "Nobody at that paper takes your calls"
	}
	if w.Player.Location != HeraldPlace {
		return "That conversation happens at the paper"
	}
	if w.pressedRecently() {
		return "The paper carried something of yours too recently"
	}
	if w.Player.Cash < PuffCost {
		return "Not enough cash"
	}
	return ""
}

// Puff is a paragraph about a businessman. It is the cheapest respect in this
// city and the only kind nobody had to be hurt for.
func (w *World) Puff() error {
	if reason := w.PuffReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(PuffCost); err != nil {
		return err
	}
	w.Player.LastPress = w.Minute
	w.Player.Respect += PuffRespect
	w.Player.Heat = max(0, w.Player.Heat-PuffHeat)
	w.Report("business", "LOCAL BUSINESSMAN BACKS DISTRICT TRADE",
		fmt.Sprintf("%s, who holds interests across the district, told the Herald that the neighbourhood's future lies in its own businesses. The remarks were welcomed by local traders.", w.Player.Name))
	w.Log("Page five, with a photograph", fmt.Sprintf("A paragraph about a local businessman. Worth %d respect and %d off what the police think, and it cost nobody anything.", PuffRespect, PuffHeat), "personal")
	return nil
}

// pressedRecently is whether the paper has carried something the player put
// there inside the cooldown. A zero is somebody who has never been able to,
// not somebody who did it at the beginning of time.
func (w *World) pressedRecently() bool {
	return w.Player.LastPress > 0 && w.Minute-w.Player.LastPress < PressCooldown
}

// dropStory removes a story from the paper, which is the only way anything ever
// leaves it besides being pushed out by newer news.
func (w *World) dropStory(id string) {
	kept := w.News[:0]
	for _, s := range w.News {
		if s.ID != id {
			kept = append(kept, s)
		}
	}
	w.News = kept
}

// PressDescription is what the paper can currently be made to do, for the
// interface.
func (w *World) PressDescription() map[string]any {
	if !w.TheEditor() {
		return nil
	}
	return map[string]any{
		"spikeable": len(w.Spikeable()),
		"ready":     !w.pressedRecently(),
	}
}
