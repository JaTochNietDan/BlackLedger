// Package core is the authoritative headless simulation. It has no rendering or HTTP dependencies.
package core

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

//go:embed locations.json
var placeJSON []byte
var Locations []Place

func init() {
	if err := json.Unmarshal(placeJSON, &Locations); err != nil {
		panic(err)
	}
}

type Place struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	District int    `json:"district"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Cost     int    `json:"cost"`
	Blurb    string `json:"blurb"`
}
type Person struct {
	Name     string `json:"name"`
	Cash     int    `json:"cash"`
	Health   int    `json:"health"`
	Respect  int    `json:"respect"`
	Heat     int    `json:"heat"`
	Location string `json:"location"`
	Home     string `json:"home"`
	BestHome int    `json:"best_home,omitempty"`
	Security int    `json:"security"`
	Contacts int    `json:"contacts"`
	Crew     []Crew `json:"crew"`
	Alive    bool   `json:"alive"`
	Earned   int    `json:"earned"`
	JobCount int    `json:"job_count"`
}
type Crew struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Loyalty int    `json:"loyalty"`
}
type NPC struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Trust int    `json:"trust"`
	Voice string `json:"voice"`
	Color string `json:"color"`
}
type Faction struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Leader   string `json:"leader"`
	Power    int    `json:"power"`
	Goodwill int    `json:"goodwill"`
	Cash     int    `json:"cash"`
}
type Property struct {
	Owner     string  `json:"owner"`
	Condition int     `json:"condition"`
	Income    int     `json:"income"`
	Carry     float64 `json:"carry"`
}
type Plot struct {
	Target   string `json:"target,omitempty"`
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Life     int    `json:"life"`
	Due      int    `json:"due"`
	Actor    string `json:"actor"`
	Strength int    `json:"strength"`
	Known    bool   `json:"known"`
}
type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Due  int    `json:"due"`
}
type Record struct {
	ID     string `json:"id"`
	Minute int    `json:"minute"`
	Life   int    `json:"life"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Kind   string `json:"kind"`
}
type Death struct {
	Name   string `json:"name"`
	Minute int    `json:"minute"`
	Life   int    `json:"life"`
	Cause  string `json:"cause"`
}
type Choice struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Detail   string `json:"detail"`
	Cost     int    `json:"cost"`
	Disabled bool   `json:"disabled"`
}
type Effect struct {
	Reward  int `json:"reward"`
	Respect int `json:"respect"`
	Heat    int `json:"heat"`
	Minutes int `json:"minutes"`
}
type Scene struct {
	Connection   *StoryConnection  `json:"connection,omitempty"`
	Operation    string            `json:"operation,omitempty"`
	JobID        string            `json:"job_id,omitempty"`
	Alternatives map[string]Effect `json:"alternatives,omitempty"`
	Beneficiary  string            `json:"beneficiary,omitempty"`
	Target       string            `json:"target,omitempty"`
	Actor        string            `json:"actor,omitempty"`
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Speaker      string            `json:"speaker"`
	Choices      []Choice          `json:"choices"`
	Kind         string            `json:"kind"`
	Source       string            `json:"source"`
	Minute       int               `json:"minute"`
	Effect       Effect            `json:"effect"`
	Outcome      string            `json:"outcome"`
}
type Offer struct {
	Ready int    `json:"ready"`
	Event *Scene `json:"event"`
}
type Director struct {
	Status      string `json:"status"`
	Detail      string `json:"detail"`
	LastRequest int    `json:"last_request"`
}
type VisualCue struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Target  string `json:"target"`
	Caption string `json:"caption"`
}
type Result struct {
	Cues    []VisualCue `json:"cues,omitempty"`
	From    string      `json:"from_location"`
	To      string      `json:"to_location"`
	Elapsed int         `json:"elapsed"`
	Records []Record    `json:"records"`
}
type World struct {
	VisualCues     []VisualCue          `json:"-"`
	Arrangements   []ArrangementMemory  `json:"arrangements,omitempty"`
	NextPressure   int                  `json:"next_pressure,omitempty"`
	Version        int                  `json:"version"`
	ID             string               `json:"id"`
	Revision       int                  `json:"revision"`
	Life           int                  `json:"life"`
	Minute         int                  `json:"minute"`
	RNG            uint32               `json:"rng"`
	Player         Person               `json:"player"`
	District       int                  `json:"district"`
	BusinessTruces map[string]int       `json:"business_truces,omitempty"`
	Factions       []Faction            `json:"factions"`
	NPCs           []NPC                `json:"npcs"`
	Properties     map[string]*Property `json:"properties"`
	Plots          []Plot               `json:"plots"`
	Tasks          []Task               `json:"tasks"`
	Event          *Scene               `json:"event"`
	History        []Record             `json:"history"`
	Dead           []Death              `json:"dead"`
	Director       Director             `json:"director"`
	Offers         []Offer              `json:"offers"`
	LastResult     *Result              `json:"last_result"`
}
type Command struct {
	RequestID string `json:"request_id"`
	Revision  int    `json:"revision"`
	Kind      string `json:"kind"`
	Target    string `json:"target"`
	Event     string `json:"event"`
	Choice    string `json:"choice"`
}
type Action struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Minutes  int    `json:"minutes"`
	Cost     int    `json:"cost"`
	Disabled bool   `json:"disabled"`
	Reason   string `json:"reason"`
	Detail   string `json:"detail"`
	Target   string `json:"target"`
}

func ID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
func PlaceByID(id string) (Place, bool) {
	for _, p := range Locations {
		if p.ID == id {
			return p, true
		}
	}
	return Place{}, false
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (w *World) Random() float64 {
	w.RNG = 1664525*w.RNG + 1013904223
	return float64(w.RNG) / 4294967296
}
func (w *World) Log(title, text, kind string) {
	w.History = append(w.History, Record{ID(), w.Minute, w.Life, title, text, kind})
	if len(w.History) > 180 {
		w.History = w.History[len(w.History)-180:]
	}
}
func (w *World) CanAcquire(id string) bool {
	p := w.Properties[id]
	return p != nil && (p.Owner == "independent" || strings.HasPrefix(p.Owner, "former:"))
}
func (w *World) Own(id string) bool {
	return w.Properties[id] != nil && w.Properties[id].Owner == fmt.Sprintf("player:%d", w.Life)
}
func newPerson(life int) Person {
	names := []string{"Alex Varga", "Nico Ward", "Frankie Vale", "Sam Costa", "Jamie Moretti", "Robin Hale"}
	return Person{Name: names[(life-1)%len(names)], Cash: 90, Health: 100, Location: "room", Home: "room", Alive: true, Crew: []Crew{}}
}
func New(seed uint32) *World {
	w := &World{Version: 2, ID: ID(), Life: 1, Minute: 480, RNG: seed, Player: newPerson(1), Properties: map[string]*Property{}, Tasks: []Task{}, Plots: []Plot{}, History: []Record{}, Dead: []Death{}, Offers: []Offer{}, Director: Director{"authored", "Authored opening. Local AI can prepare additional encounters.", -9999}}
	w.Factions = []Faction{{"bellandi", "Bellandi Family", "Vittorio Bellandi", 90, 0, 8000}, {"russo", "Russo Outfit", "Elena Russo", 58, 0, 4500}}
	w.NPCs = []NPC{{"mara", "Mara Bell", "Fixer", 10, "af_heart", "#a48761"}, {"leo", "Leo Carver", "Driver", 20, "am_michael", "#9ca795"}, {"vittorio", "Vittorio Bellandi", "Bellandi boss", 0, "bm_george", "#ad7970"}, {"elena", "Elena Russo", "Russo boss", 0, "bf_emma", "#83989b"}}
	for _, p := range Locations {
		owner := "independent"
		if p.ID == "club" {
			owner = "bellandi"
		}
		income := 0
		switch p.ID {
		case "laundry":
			income = 14
		case "garage":
			income = 24
		case "casino":
			income = 48
		}
		w.Properties[p.ID] = &Property{owner, 100, income, 0}
	}
	w.Log("A room. A name. No protection.", "Mara Bell left word at Saint Agnes: there is work, if you can be discreet. Your room costs $15 each midnight.", "personal")
	return w
}
func (w *World) Clone() *World {
	b, _ := json.Marshal(w)
	var n World
	_ = json.Unmarshal(b, &n)
	return &n
}
func (w *World) Earn(n int) { w.Player.Cash += n; w.Player.Earned += n }
func (w *World) Pay(n int) error {
	if n < 0 || w.Player.Cash < n {
		return fmt.Errorf("you cannot afford this commitment")
	}
	w.Player.Cash -= n
	return nil
}
func (w *World) NPC(id string) *NPC {
	for i := range w.NPCs {
		if w.NPCs[i].ID == id {
			return &w.NPCs[i]
		}
	}
	return nil
}
func HomeRent(id string) int {
	switch id {
	case "apartment":
		return 35
	case "estate":
		return 90
	}
	return 15
}
func HomeRank(id string) int {
	switch id {
	case "estate":
		return 2
	case "apartment":
		return 1
	}
	return 0
}
func (w *World) Guard() int {
	n := w.Player.Security
	switch w.Player.Home {
	case "apartment":
		n++
	case "estate":
		n += 2
	}
	return n
}
func (w *World) DailyCost() int {
	return HomeRent(w.Player.Home) + 10*w.Player.Security + 12*len(w.Player.Crew)
}
func TravelMinutes(a, b string) int {
	x, _ := PlaceByID(a)
	y, _ := PlaceByID(b)
	return max(10, int(math.Round((math.Abs(float64(x.X-y.X))+math.Abs(float64(x.Y-y.Y)))/60))*5)
}
func (w *World) Actions(id string) []Action {
	out := []Action{}
	p := &w.Player
	l, ok := PlaceByID(id)
	if !ok || !p.Alive || w.Event != nil {
		return out
	}
	add := func(id, label string, minutes, cost int, reason, detail string) {
		if reason == "" && p.Cash < cost {
			reason = "Not enough cash"
		}
		out = append(out, Action{id, label, minutes, cost, reason != "", reason, detail, l.ID})
	}
	need := func(b bool, s string) string {
		if b {
			return s
		}
		return ""
	}
	if l.District > w.District {
		add("expand", "Establish contacts across town", 90, 100, need(p.Respect < 10, "Earn 10 respect first"), "Opens the next district. Existing businesses and rivals remain.")
		return out
	}
	if l.ID != p.Location {
		add("travel", "Visit "+l.Name, TravelMinutes(p.Location, l.ID), 0, "", "Travel advances the city clock. Known threats may interrupt you.")
		return out
	}
	switch id {
	case "bar":
		add("courier", "Carry a discreet envelope", 45, 0, "", "Earn $45 and 2 respect. A reliable introduction to the neighborhood.")
		add("contact", "Buy Mara a coffee", 30, 10, need(p.Contacts >= 5, "Your information network is fully developed"), "Build trust and an information network. Contacts may warn you of trouble.")
		reason := need(p.Respect < 6, "Earn 6 respect first")
		if len(p.Crew) > 0 {
			reason = "Leo is already in your crew"
		}
		add("recruit", "Recruit Leo Carver", 30, 90, reason, "A driver and collector. $12 daily wages; loyalty matters.")
	case "garage":
		add("audience", "Request an audience with Russo", 45, 0, "", "Discuss your standing with the Russo Outfit.")
	case "docks":
		add("dockwork", "Work the night cargo", 90, 0, "", "Earn $75 and 1 respect. Small chance of a work injury.")
	case "market":
		add("investigate", "Ask about threats", 45, 30, "", "Investigate existing threats. Evidence is not a guarantee of safety.")
		add("lie_low", "Keep a low profile", 120, 15, "", "Lose 10 heat. Time still passes for rivals and businesses.")
	case "club":
		add("audience", "Request an audience", 45, 0, "", "Discuss your standing with the Bellandi family.")
		add("provoke", "Demand protection money", 30, 0, "", "EXTREME RISK. Bellandi owns this casino. Challenging him can bring lethal retaliation.")
	}
	if id == "laundry" || id == "garage" || id == "casino" {
		if w.Own(id) {
			add("inspect", "Review the books", 0, 0, "", "Read current income and repair needs without advancing time.")
			add("repair", "Repair the property", 60, 50, need(w.Properties[id].Condition >= 100, "Already in good condition"), "Restore 40 condition.")
		} else {
			req := 6
			if id == "garage" {
				req = 10
			}
			if id == "casino" {
				req = 20
			}
			reason := need(p.Respect < req, fmt.Sprintf("Earn %d respect first", req))
			if !w.CanAcquire(id) {
				reason = "This property belongs to another organization"
			}
			label := "Establish protection"
			if id == "casino" {
				label = "Reopen the casino"
			}
			cost := l.Cost
			if strings.HasPrefix(w.Properties[id].Owner, "former:") {
				cost *= 2
				label = "Buy out the former organization"
			}
			add("acquire", label, 60, cost, reason, fmt.Sprintf("Earn up to $%d/hour. Income accrues automatically; rivals may take notice.", w.Properties[id].Income))
		}
	}
	if l.Type == "home" {
		if p.Home != id {
			label := "Rent this apartment"
			cost, reason := l.Cost, ""
			if id == "room" {
				label = "Return to a rented room"
			}
			if id == "estate" {
				label = "Buy this residence"
				if w.Own(id) {
					label, cost = "Return to your residence", 0
				} else if w.Properties[id].Owner != "independent" {
					reason = "This residence belongs to another organization"
				}
			}
			add("move_home", label, 60, cost, reason, fmt.Sprintf("$%d/day upkeep. Moving resets hired security. Respect is earned only for a new housing tier.", HomeRent(id)))
		} else {
			add("rest", "Rest for four hours", 240, 0, "", "Recover up to 25 health as you rest. Rivals can act while you sleep.")
			add("security", "Hire another security detail", 30, 100, need(p.Security >= 3, "Maximum security hired"), "Improves detection and survival at home. Adds $10/day upkeep.")
			if w.Properties[id].Condition < 100 {
				add("repair", "Repair the residence", 60, 50, "", "Restore 40 condition.")
			}
		}
	}
	if len(p.Crew) > 0 {
		reason := need(p.Crew[0].Loyalty < 30, "Leo refuses assignments below 30 loyalty. Pay a bonus to rebuild trust.")
		if len(w.Tasks) > 0 {
			reason = "Leo is already on assignment"
		}
		add("delegate", "Send Leo on collections", 15, 0, reason, "Completes after 120 game minutes: $65. Requires 30 loyalty.")
		add("crew_bonus", "Pay Leo a bonus", 15, 40, need(p.Crew[0].Loyalty >= 100, "Loyalty is already at its maximum"), "Restore up to 25 loyalty. Below 30 he refuses collections; at 50 he can help protect businesses when available.")
	}
	add("wait", "Let an hour pass", 60, 0, "", "Income, rent, operations and rival plans continue.")
	return out
}
func (w *World) Retaliation() { w.RetaliationFrom("bellandi") }

// RetaliationFrom commits one active personal operation per family and life.
func (w *World) RetaliationFrom(actor string) {
	valid := false
	for _, f := range w.Factions {
		if f.ID == actor {
			valid = true
		}
	}
	if !valid {
		return
	}
	for _, p := range w.Plots {
		if p.Life == w.Life && p.Kind == "hit" && p.Actor == actor {
			return
		}
	}
	w.Plots = append(w.Plots, Plot{ID: ID(), Kind: "hit", Life: w.Life, Due: w.Minute + 240, Actor: actor, Strength: 5})
}
func (w *World) factionName(id string) string {
	for _, f := range w.Factions {
		if f.ID == id {
			return f.Name
		}
	}
	return "An unidentified family"
}
func (w *World) Die(cause string) {
	p := &w.Player
	p.Alive = false
	p.Health = 0
	w.Dead = append(w.Dead, Death{p.Name, w.Minute, w.Life, cause})
	for id, prop := range w.Properties {
		if w.Own(id) {
			prop.Owner = "former:" + p.Name
			if prop.Income > 0 {
				prop.Income = max(5, prop.Income-2)
			}
		}
	}
	w.Tasks = []Task{}
	w.Event = nil
	w.Log(p.Name+" is dead", cause+" Your life ends here. The city continues.", "death")
}
func (w *World) Attack(plot Plot) {
	p := &w.Player
	if p.Location != p.Home {
		w.Properties[p.Home].Condition = max(10, w.Properties[p.Home].Condition-45)
		w.VisualCues = append(w.VisualCues, VisualCue{ID(), "attack", p.Home, "Armed men damaged your residence while you were away."})
		w.Log("Someone came looking", "You were away. Armed men damaged your residence and left before anyone could identify them.", "danger")
		return
	}
	if plot.Known || w.Guard() > 0 || p.Contacts >= 2 {
		home, _ := PlaceByID(p.Home)
		body := "A car stops outside " + home.Name + ". "
		if w.Guard() > 0 {
			body += "Your security raises the alarm."
		} else {
			body += "A contact calls: leave by the back, now."
		}
		w.Event = &Scene{ID: "attack-" + plot.ID, Title: "Headlights outside", Body: body + " You have moments to act.", Speaker: "mara", Kind: "attack", Source: "authored", Minute: w.Minute, Choices: []Choice{{ID: "escape", Label: "Leave through the rear", Detail: "A chance to escape. Security and contacts help; injuries reduce your odds."}, {ID: "defend", Label: "Hold the entrance", Detail: "Rely on your security. Injuries and a weak defense can be fatal."}, {ID: "bargain", Label: "Offer $180 to stand down", Cost: 180, Detail: "Money may settle this incident, but your standing suffers."}}}
	} else if w.Random() < .78 {
		w.Die("An attack at your residence caught you without warning or protection.")
	} else {
		p.Health = max(1, p.Health-65)
		w.Log("You survived by inches", "The attackers leave you wounded. Nobody warned you. You need rest and protection.", "danger")
	}
}
func (w *World) Advance(minutes int) {
	p := &w.Player
	end := w.Minute + max(0, minutes)
	for w.Minute < end {
		if !p.Alive || w.Event != nil {
			return
		}
		// Jump to the next meaningful boundary; presentation never drives this clock.
		next := min(end, (w.Minute/1440+1)*1440)
		if w.NextPressure > 0 {
			next = min(next, max(w.Minute+1, w.NextPressure))
		}
		for _, task := range w.Tasks {
			next = min(next, max(w.Minute+1, task.Due))
		}
		for _, plot := range w.Plots {
			if plot.Life == w.Life {
				next = min(next, max(w.Minute+1, plot.Due))
				if plot.Kind == "hit" && !plot.Known && p.Contacts >= 2 {
					next = min(next, max(w.Minute+1, plot.Due-90))
				}
			}
		}
		elapsed := next - w.Minute
		w.Minute = next
		for id, prop := range w.Properties {
			if w.Own(id) {
				prop.Carry += float64(prop.Income*prop.Condition*elapsed) / 6000
				n := int(prop.Carry + 1e-9)
				prop.Carry -= float64(n)
				w.Earn(n)
			}
		}
		for j := 0; j < len(w.Tasks); {
			if w.Tasks[j].Due <= w.Minute {
				w.Earn(65)
				w.Log("Leo returns", "$65 from collections. He is available again.", "business")
				w.Tasks = append(w.Tasks[:j], w.Tasks[j+1:]...)
			} else {
				j++
			}
		}
		if w.Minute%1440 == 0 {
			bill := w.DailyCost()
			if p.Cash >= bill {
				p.Cash -= bill
				w.Log("Accounts settled", fmt.Sprintf("$%d paid for housing, security and crew.", bill), "business")
			} else {
				p.Cash = max(0, p.Cash-15)
				p.Home = "room"
				p.Security = 0
				if len(p.Crew) > 0 {
					p.Crew[0].Loyalty = max(0, p.Crew[0].Loyalty-20)
				}
				w.Log("Your arrangements unravel", "You could not cover the bills. Security leaves; your residence is now a rented room. Unpaid crew lose loyalty.", "danger")
			}
		}
		for j := 0; j < len(w.Plots); j++ {
			plot := &w.Plots[j]
			if plot.Life != w.Life {
				continue
			}
			if plot.Kind == "hit" && !plot.Known && p.Contacts >= 2 && plot.Due-w.Minute <= 90 {
				plot.Known = true
				warning := w.factionName(plot.Actor) + " has people asking where you sleep."
				w.Log("Mara has heard something", warning+" You may have very little time.", "danger")
				if plot.Due > w.Minute {
					w.Event = &Scene{ID: "warning-" + plot.ID, Title: "A call worth answering", Body: warning + " I cannot tell you exactly when they will come. Stop what you are doing and think about where you want to be tonight.", Speaker: "mara", Kind: "warning", Source: "authored", Minute: w.Minute, Choices: []Choice{{ID: "acknowledge", Label: "Put down the phone and prepare", Detail: "Clock stays paused. You can leave, arrange security or seek an audience. The threat remains."}}}
					return
				}
			}
			if plot.Due <= w.Minute {
				copy := *plot
				w.Plots = append(w.Plots[:j], w.Plots[j+1:]...)
				if copy.Kind == "sabotage" {
					w.ResolveSabotage(copy)
				} else {
					w.Attack(copy)
				}
				break
			}
		}
		if w.Event == nil && p.Alive && w.NextPressure > 0 && w.Minute >= w.NextPressure {
			w.BusinessPressure()
		}
	}
}

// KnownThreats exposes discovered intelligence, never schedules or undiscovered plans.
func (w *World) KnownThreats() []string {
	out := []string{}
	if !w.Player.Alive {
		return out
	}
	for _, p := range w.Plots {
		if p.Life != w.Life || !p.Known {
			continue
		}
		text := w.factionName(p.Actor) + " has commissioned an attack against you. Consider leaving home, arranging security, or negotiating."
		if p.Kind == "sabotage" {
			l, ok := PlaceByID(p.Target)
			if !ok {
				continue
			}
			text = w.factionName(p.Actor) + " is targeting " + l.Name + ". Available loyal crew can limit damage; negotiation can stop this family's current operations."
		}
		out = append(out, text)
	}
	return out
}
func (w *World) Public() map[string]any {
	locs := []map[string]any{}
	income := 0.0
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if w.Own(l.ID) {
			income += float64(prop.Income*prop.Condition) / 100
		}
		locs = append(locs, map[string]any{"id": l.ID, "name": l.Name, "type": l.Type, "district": l.District, "x": l.X, "y": l.Y, "cost": l.Cost, "blurb": l.Blurb, "owner": prop.Owner, "condition": prop.Condition, "income": prop.Income, "owned": w.Own(l.ID), "locked": l.District > w.District, "actions": w.Actions(l.ID)})
	}
	var scene any = nil
	if e := w.Event; e != nil {
		choices := append([]Choice{}, e.Choices...)
		for i := range choices {
			choices[i].Disabled = w.Player.Cash < choices[i].Cost
		}
		scene = map[string]any{"id": e.ID, "title": e.Title, "body": e.Body, "speaker": e.Speaker, "kind": e.Kind, "source": e.Source, "minute": e.Minute, "choices": choices, "connection": e.Connection}
	}
	history := w.History
	if len(history) > 60 {
		history = history[len(history)-60:]
	}
	return map[string]any{"id": w.ID, "version": w.Version, "revision": w.Revision, "life": w.Life, "minute": w.Minute, "player": w.Player, "district": w.District, "factions": w.Factions, "npcs": w.NPCs, "locations": locs, "event": scene, "history": history, "dead": w.Dead, "tasks": w.Tasks, "director": w.Director, "last_result": w.LastResult, "daily_cost": w.DailyCost(), "income": income, "security": w.Guard(), "opportunity": w.NextOpportunity(), "known_threats": w.KnownThreats(), "business_truces": w.ActiveBusinessTruces()}
}
func (w *World) hasRecord(title string) bool {
	for _, r := range w.History {
		if r.Life == w.Life && r.Title == title {
			return true
		}
	}
	return false
}
func (w *World) OfferIfReady() {
	if !w.Player.Alive || w.Event != nil {
		return
	}
	if len(w.Offers) > 0 && w.Minute >= w.Offers[0].Ready {
		w.Event = w.Offers[0].Event
		w.Event.Minute = w.Minute
		w.RememberArrangement(w.Event, "offered")
		w.Log(w.Event.Title, w.Event.Body, "story")
		w.Offers = w.Offers[1:]
		w.Director.Status = "available"
		w.Director.Detail = "Ready to prepare another situation."
		return
	}
	if w.Player.JobCount == 2 && !w.hasRecord("A favor with a price") {
		e, _ := w.ValidateProposal(Proposal{"", "A favor with a price", "“A merchant wants a sealed ledger moved before his partners arrive. I would understand if you preferred the ordinary work.”", "mara", "courier", "You moved the ledger. Mara now knows you can handle sensitive work.", "", []Approach{{Method: "careful", Label: "Wait for a quiet route"}, {Method: "press", Label: "Move it before the partners arrive"}}})
		e.Source = "authored"
		w.Event = e
		w.RememberArrangement(e, "offered")
		w.Log("A favor with a price", "Mara offers more sensitive work.", "story")
	}
}

type Approach struct {
	Method string `json:"method"`
	Label  string `json:"label"`
}

type Proposal struct {
	Location    string     `json:"location,omitempty"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Speaker     string     `json:"speaker"`
	Operation   string     `json:"operation"`
	Outcome     string     `json:"outcome"`
	Beneficiary string     `json:"beneficiary,omitempty"`
	Approaches  []Approach `json:"approaches,omitempty"`
}

func (w *World) ValidateProposal(p Proposal) (*Scene, error) {
	if p.Location != "" {
		location, ok := PlaceByID(p.Location)
		if !ok || location.District > w.District {
			return nil, fmt.Errorf("job location is not accessible")
		}
	}
	catalog := map[string]Effect{"courier": {75, 3, 3, 45}, "mediation": {55, 5, 0, 60}, "collection": {130, 3, 5, 90}}
	fx, ok := catalog[p.Operation]
	if !ok || w.NPC(p.Speaker) == nil {
		return nil, fmt.Errorf("unsupported director operation or speaker")
	}
	if len(strings.TrimSpace(p.Title)) < 3 || len(p.Title) > 70 || len(p.Body) < 3 || len(p.Body) > 1200 || len(p.Outcome) < 3 || len(p.Outcome) > 700 {
		return nil, fmt.Errorf("invalid director text")
	}
	politicalDetail := ""
	if p.Beneficiary != "" {
		found := false
		for _, f := range w.Factions {
			if f.ID == p.Beneficiary {
				found = true
				politicalDetail = " · " + f.Name + " standing +6; rival standing −3."
			}
		}
		if !found {
			allowed := []string{""}
			for _, f := range w.Factions {
				allowed = append(allowed, f.ID)
			}
			return nil, fmt.Errorf("unknown beneficiary faction %q; use exactly one of %q", p.Beneficiary, allowed)
		}
	}
	scene := &Scene{Target: p.Location, Operation: p.Operation, ID: ID(), Title: p.Title, Beneficiary: p.Beneficiary, Body: p.Body, Speaker: p.Speaker, Kind: "proposal", Source: "local-ai", Minute: w.Minute, Effect: fx, Outcome: map[string]string{"courier": "You delivered the sealed package and reported back.", "mediation": "You completed the requested mediation without violence.", "collection": "You collected the agreed payment and reported back."}[p.Operation], Choices: []Choice{{ID: "accept", Label: map[string]string{"courier": "Deliver the package", "mediation": "Mediate the dispute", "collection": "Collect the payment"}[p.Operation], Detail: fmt.Sprintf("$%d · %d minutes · +%d respect · +%d heat", fx.Reward, fx.Minutes, fx.Respect, fx.Heat) + " · At 15 heat, police may stop completion." + politicalDetail}, {ID: "decline", Label: "Decline the arrangement", Detail: "No cost or time."}}}
	if err := addApproaches(scene, p.Approaches, politicalDetail); err != nil {
		return nil, err
	}
	if p.Location != "" {
		place, _ := PlaceByID(p.Location)
		for i := range scene.Choices {
			if scene.Choices[i].ID != "decline" {
				scene.Choices[i].Detail = place.Name + " · " + scene.Choices[i].Detail
			}
		}
	}
	return scene, nil
}
