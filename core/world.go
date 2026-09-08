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
	// What the player is carrying. Absent in saves written before the trade
	// existed, which is the same as carrying nothing.
	Stock map[string]int `json:"stock,omitempty"`
	// When the books last absorbed a round. Absent in older saves, which is the
	// same as never having done it.
	LastLaunder int `json:"last_launder,omitempty"`
	// What the player is carrying. Absent in older saves, which is the same as
	// carrying nothing.
	Weapon int `json:"weapon,omitempty"`
	Armour int `json:"armour,omitempty"`
	// How the player is dressed, and how well kept it is. Absent in saves
	// written before clothes existed, which is working clothes in good order.
	Dress     int `json:"dress,omitempty"`
	DressWear int `json:"dress_wear,omitempty"`
	// What the player drives, and the state of it. Absent in saves written
	// before there were cars, which is walking.
	Car     int `json:"car,omitempty"`
	CarWear int `json:"car_wear,omitempty"`
	// Whether this person has established that the account abroad is theirs.
	// Reset with every life, which is what makes inheriting it a decision.
	Offshore bool `json:"offshore_access,omitempty"`
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
	// Everyone with a name is a participant in the city, not scenery. Dead
	// rather than Alive so that saves written before people could die read back
	// as living.
	Faction  string `json:"faction,omitempty"`
	Location string `json:"location,omitempty"`
	Rank     int    `json:"rank,omitempty"`
	Ambition int    `json:"ambition,omitempty"`
	Skill    int    `json:"skill,omitempty"`
	Dead     bool   `json:"dead,omitempty"`
}
type Faction struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Leader   string `json:"leader"`
	Power    int    `json:"power"`
	Goodwill int    `json:"goodwill"`
	Cash     int    `json:"cash"`
	// Peak is the strength this family recovers toward once its holdings are
	// repaired. Saves written before families held property carry no peak.
	Peak int `json:"peak,omitempty"`
}
type Property struct {
	Owner     string  `json:"owner"`
	Condition int     `json:"condition"`
	Income    int     `json:"income"`
	Carry     float64 `json:"carry"`
	// How the business is run. Empty means the ordinary way, so saves written
	// before this was a decision keep earning exactly what they earned.
	Mode string `json:"mode,omitempty"`
	// The inside of a business: who works it, what it runs on, and whether
	// something has gone wrong. Absent in older saves, which read as a business
	// nobody has staffed yet.
	Staff   int  `json:"staff,omitempty"`
	Supply  int  `json:"supply,omitempty"`
	Trouble bool `json:"trouble,omitempty"`
	Still   bool `json:"still,omitempty"`
	// What is behind the tables at a casino. Absent everywhere else, and in
	// saves written before a room ran a float of its own.
	Bankroll int `json:"bankroll,omitempty"`
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
	SuspendedJob   *SuspendedJob        `json:"suspended_job,omitempty"`
	VisualCues     []VisualCue          `json:"-"`
	Arrangements   []ArrangementMemory  `json:"arrangements,omitempty"`
	NextPressure   int                  `json:"next_pressure,omitempty"`
	Version        int                  `json:"version"`
	ID             string               `json:"id"`
	Revision       int                  `json:"revision"`
	Life           int                  `json:"life"`
	Minute         int                  `json:"minute"`
	RNG            uint32               `json:"rng"`
	WorldRNG       uint32               `json:"world_rng,omitempty"`
	Player         Person               `json:"player"`
	District       int                  `json:"district"`
	BusinessTruces map[string]int       `json:"business_truces,omitempty"`
	Factions       []Faction            `json:"factions"`
	NPCs           []NPC                `json:"npcs"`
	Properties     map[string]*Property `json:"properties"`
	Conflicts      []Conflict           `json:"conflicts,omitempty"`
	Goods          []Good               `json:"goods,omitempty"`
	// Money sent out of the city. Deliberately on the world rather than the
	// player, because it outlives them; new_life resets the player and leaves
	// this standing.
	Offshore   int        `json:"offshore,omitempty"`
	Contracts  []Contract `json:"contracts,omitempty"`
	News       []Story    `json:"news,omitempty"`
	Plots      []Plot     `json:"plots"`
	Tasks      []Task     `json:"tasks"`
	Event      *Scene     `json:"event"`
	History    []Record   `json:"history"`
	Dead       []Death    `json:"dead"`
	Director   Director   `json:"director"`
	Offers     []Offer    `json:"offers"`
	LastResult *Result    `json:"last_result"`
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

// WorldRandom draws from a stream reserved for events the player is not party
// to, so families quarrelling off-screen cannot shift the odds of a decision the
// player is making. Saves written before this stream existed start it from the
// campaign seed.
func (w *World) WorldRandom() float64 {
	if w.WorldRNG == 0 {
		w.WorldRNG = w.RNG ^ 0x9e3779b9
	}
	w.WorldRNG = 1664525*w.WorldRNG + 1013904223
	return float64(w.WorldRNG) / 4294967296
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
	w := &World{Version: SaveVersion, ID: ID(), Life: 1, Minute: 480, RNG: seed, Player: newPerson(1), Properties: map[string]*Property{}, Tasks: []Task{}, Plots: []Plot{}, History: []Record{}, Dead: []Death{}, Offers: []Offer{}, Director: Director{"authored", "Authored opening. Local AI can prepare additional encounters.", -9999}}
	w.Factions = []Faction{{"bellandi", "Bellandi Family", "Vittorio Bellandi", 90, 0, 8000, 90}, {"russo", "Russo Outfit", "Elena Russo", 58, 0, 4500, 58}}
	w.NPCs = []NPC{
		{ID: "mara", Name: "Mara Bell", Role: "Fixer", Trust: 10, Voice: "af_heart", Color: "#a48761", Location: "bar", Rank: RankAssociate, Ambition: 45, Skill: 60},
		{ID: "leo", Name: "Leo Carver", Role: "Driver", Trust: 20, Voice: "am_michael", Color: "#9ca795", Location: "bar", Rank: RankAssociate, Ambition: 35, Skill: 45},
		{ID: "vittorio", Name: "Vittorio Bellandi", Role: "Head of the Bellandi Family", Voice: "bm_george", Color: "#ad7970", Faction: "bellandi", Location: "club", Rank: RankLeader, Ambition: 70, Skill: 80},
		{ID: "elena", Name: "Elena Russo", Role: "Head of the Russo Outfit", Voice: "bf_emma", Color: "#83989b", Faction: "russo", Location: "market", Rank: RankLeader, Ambition: 75, Skill: 72},
	}
	for _, p := range Locations {
		owner := "independent"
		switch p.ID {
		case "club", "docks":
			owner = "bellandi"
		case "market", "bar":
			// Russo needs holdings of its own, or only one family can be
			// pressured. Two each also means a war costs ground before it costs
			// an organization its existence.
			owner = "russo"
		}
		income := 0
		switch p.ID {
		case "laundry":
			income = 14
		case "garage":
			income = 24
		case "casino":
			// The floor take only: the bar, the door and the rooms upstairs.
			// What the tables make is decided every night by the float behind
			// them, in CasinoDay, rather than accruing by the hour.
			income = 18
		case "club":
			income = 30
		case "market":
			income = 18
		case "docks":
			income = 22
		case "bar":
			income = 12
		}
		property := &Property{Owner: owner, Condition: 100, Income: income}
		// A trading business is already running before anybody buys it: it has
		// people working it and something to work with. Ownership changes who
		// answers for that, not whether it exists.
		if trade, running := TradeOf(p.ID); running {
			property.Staff, property.Supply = trade.Hands, trade.RestockAmount
		}
		w.Properties[p.ID] = property
	}
	// Each organization is people, not a name and a number. These are the ones
	// who would step up if the person above them died.
	for _, f := range []string{"bellandi", "russo"} {
		w.AddMember(f, "Lieutenant", RankLieutenant, w.homeOf(f))
		w.AddMember(f, "Soldier", RankSoldier, w.homeOf(f))
	}
	w.Goods = newGoods()
	// The two established families are already rivals when the player arrives.
	w.Antagonize("bellandi", "russo", 50)
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
	return HomeRent(w.Player.Home) + 10*w.Player.Security + 12*len(w.Player.Crew) + w.Wages() + w.CarUpkeep()
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
		detail := "Travel advances the city clock. Known threats may interrupt you."
		if walk := TravelMinutes(p.Location, l.ID); w.Driving() && w.Journey(p.Location, l.ID) < walk {
			detail = fmt.Sprintf("%d minutes on foot, %d driving. Travel advances the city clock. Known threats may interrupt you.", walk, w.Journey(p.Location, l.ID))
		}
		add("travel", "Visit "+l.Name, w.Journey(p.Location, l.ID), 0, "", detail)
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
		if next, ok := nextVehicle(p.Car); ok {
			hides := "Nothing to hide anything in."
			if next.Compartment > 0 {
				hides = fmt.Sprintf("A false floor a search will not find %d units under.", next.Compartment)
			}
			add("car", "Buy "+lowerFirst(next.Label), 60, 0, w.CarReadiness(),
				fmt.Sprintf("$%d, then $%d a day to keep on the road. %s Journeys take %d%% of the time they take on foot. %s A car outside is a thing witnesses describe.", next.Cost, next.Upkeep, next.Detail, int(next.Pace*100), hides))
		}
		fee := w.ServiceFee()
		service := fmt.Sprintf("$%d. Restores up to 55 condition, currently %d of 100. Below %d it is worth nothing to you.", fee, w.CarCondition(), Wreck)
		if fee == 0 {
			service = fmt.Sprintf("Your own people, at no charge. Restores up to 55 condition, currently %d of 100.", w.CarCondition())
		}
		add("service", "Have the car worked on", CarServiceMinutes, 0, w.ServiceReadiness(id), service)
	case "docks":
		add("dockwork", "Work the night cargo", 90, 0, "", "Earn $75 and 1 respect. Small chance of a work injury.")
		if next, ok := nextArmament(weapons, p.Weapon); ok {
			add("arms:weapon", "Buy "+next.Label, 45, 0, w.ArmsReadiness("weapon"),
				fmt.Sprintf("$%d. %s Improves your odds when violence is your idea. A search takes it.", next.Cost, next.Detail))
		}
		if next, ok := nextArmament(armour, p.Armour); ok {
			add("arms:armour", "Buy "+next.Label, 45, 0, w.ArmsReadiness("armour"),
				fmt.Sprintf("$%d. %s Reduces what a beating costs you. A search takes it.", next.Cost, next.Detail))
		}
	case "market":
		add("investigate", "Ask about threats", 45, 30, "", "Investigate existing threats. Evidence is not a guarantee of safety.")
		add("lie_low", "Keep a low profile", 120, 15, "", "Lose 10 heat. Time still passes for rivals and businesses.")
		add("deposit", fmt.Sprintf("Wire $%d out of the city", DepositLot), 45, 0, w.DepositReadiness(),
			fmt.Sprintf("$%d of it arrives; the arrangement takes %d%%. It survives you, and whoever comes next can reach it if they can afford to.", DepositLot*(100-DepositCut)/100, DepositCut))
		add("offshore_access", "Establish that the account is yours", AccessMinutes, 0, w.AccessReadiness(),
			fmt.Sprintf("$%d in papers and a journey. Only worth it if there is enough out there to be worth reaching.", AccessCost))
		add("withdraw", "Bring it all home", 45, 0, w.WithdrawReadiness(),
			fmt.Sprintf("Brings $%d back into the city, where it can be taken from you.", w.Offshore))
		if next, ok := nextAttire(p.Dress); ok {
			notice := "Nobody official looks twice at it."
			if next.Notice > 0 {
				notice = fmt.Sprintf("Dressing above your visible means draws %d police attention a day.", next.Notice)
			}
			add("dress", "Be measured for "+lowerFirst(next.Label), 60, 0, w.DressReadiness(),
				fmt.Sprintf("$%d. %s Worth %d presence while it is kept, and it wears. %s", next.Cost, next.Detail, next.Presence, notice))
		}
		add("bribe", "An understanding with the detective", 45, 0, w.BribeReadiness(),
			fmt.Sprintf("$%d to Detective Harlow to lose some paperwork. Clears attention now and buys nothing later. Above %d heat nobody will be seen taking it.", w.BribeCost(), BribeCeiling))
		add("contract", "Ask about a name", 30, 0,
			need(p.Contacts < 1, "Build a contact who will carry this"),
			"Put a price on somebody. What it costs depends on who they are and who does the work. A failed attempt can be traced back to you.")
	case "club":
		add("audience", "Request an audience", 45, 0, "", "Discuss your standing with the Bellandi family.")
		add("provoke", "Demand protection money", 30, 0, "", "EXTREME RISK. Bellandi owns this casino. Challenging him can bring lethal retaliation.")
	}
	if HasTables(id) && !w.Own(id) {
		for _, stake := range tableStakes {
			// Cost is zero here because Play charges the stake itself; declaring
			// it would have the command layer charge it a second time.
			add("play:"+stake.ID, stake.Label, 60, 0, w.TableReadiness(id, stake),
				fmt.Sprintf("Stake $%d against the house. The house holds the edge, and winning heavily in somebody else's room is noticed.", stake.Amount))
		}
	}
	if place, ok := PlaceByID(id); ok && place.Type == "racket" && w.Own(id) {
		add("launder", "Run takings through the books", 90, 0, w.LaunderReadiness(id),
			fmt.Sprintf("$%d to clear up to %d police attention through %s. Wears the premises, and the books need a day between rounds.", w.LaunderFee(id), w.launderCapacity(id), l.Name))
	}
	if prop := w.Properties[id]; prop != nil && prop.Income > 0 && !w.Own(id) {
		add("rob", "Take the day's cash", 45, 0, w.RobberyReadiness(id),
			fmt.Sprintf("Walk out with what is in the till at %s. A haul, police attention, and an owner who will work out who would dare. Going wrong means a beating.", l.Name))
	}
	for _, g := range w.Goods {
		if !TradesAt(id, g.ID) {
			continue
		}
		// Cost is zero because Buy charges the lot itself; declaring it would
		// have the command layer charge it a second time.
		add("buy:"+g.ID, fmt.Sprintf("Buy %d %ss of %s", Lot, g.Unit, g.Name), 30, 0,
			w.TradeReadiness(g.ID, "buy"),
			fmt.Sprintf("$%d for the lot, at $%d each today. Holding stock draws police attention every day until it is sold, and can be taken from you.", g.Price*Lot, g.Price))
		if held := w.Holding(g.ID); held > 0 {
			add("sell:"+g.ID, fmt.Sprintf("Sell %d %ss of %s", held, g.Unit, g.Name), 30, 0,
				w.TradeReadiness(g.ID, "sell"),
				fmt.Sprintf("$%d each today, for $%d.", g.Price, g.Price*held))
		}
	}
	if f, ok := w.SabotageTarget(id); ok {
		add("sabotage", "Move against "+f.Name, 90, 0, w.SabotageReadiness(id),
			fmt.Sprintf("Send your crew against %s. Damages the property, weakens %s and costs you standing with them. They will retaliate, and a failed attempt injures you.", l.Name, f.Name))
		if rival := w.Rival(f.ID); rival != nil {
			add("incite", "Point "+f.Name+" at "+rival.Name, 45, 25, w.InciteReadiness(id),
				fmt.Sprintf("Spend $25 on the right conversations so %s believes %s moved against them. Hardens their quarrel and can start a war you are not part of. A story that does not hold up costs you standing with %s.", f.Name, rival.Name, f.Name))
		}
	}
	if id == "laundry" || id == "garage" || id == "casino" {
		if w.Own(id) {
			add("inspect", "Review the books", 0, 0, "", "Read current income and repair needs without advancing time.")
			if trade, running := TradeOf(id); running {
				prop := w.Properties[id]
				add("hire", "Take somebody on", 45, 0, w.HireReadiness(id),
					fmt.Sprintf("%d of %d positions filled. A week's wages up front at $%d a day after. Short-handed, it earns less and attracts trouble.", prop.Staff, trade.Hands, trade.Wage))
				add("layoff", "Let somebody go", 30, 0, w.LayOffReadiness(id),
					fmt.Sprintf("Cuts $%d a day from the wage bill and what the place can handle.", trade.Wage))
				add("restock", "Buy "+trade.Supplies, 45, 0, w.RestockReadiness(id),
					fmt.Sprintf("$%d. Currently %d left; a business out of %s barely trades.", trade.Restock, prop.Supply, trade.Supplies))
				if prop.Trouble {
					add("remedy", trade.Remedy, 60, 0, w.RemedyReadiness(id),
						fmt.Sprintf("$%d. %s %s", trade.RemedyCost, trade.Trouble, trade.RemedyDetail))
				}
				if StillSite(id) {
					if prop.Still {
						add("dismantle", "Take the still out", 90, 0, w.DismantleReadiness(id),
							fmt.Sprintf("Ends production of about %d crates a day and the attention that comes with it.", w.StillOutput(id)))
					} else {
						add("still", "Set up a still in the back", StillMinutes, 0, w.StillReadiness(id),
							fmt.Sprintf("$%d. Produces moonshine you can sell, draws %d attention a day on top of what the stock draws, and a search that finds it costs far more than one that does not.", StillCost, StillHeat))
					}
				}
			}
			if HasBankroll(id) {
				prop := w.Properties[id]
				add("bankroll", "Put money behind the tables", 45, 0, w.BankrollReadiness(id),
					fmt.Sprintf("$%d into the float, currently $%d. %s The house keeps roughly %d%% of what crosses the tables over a season and loses on plenty of single nights. A house that cannot pay a winner is finished as a room worth playing in.", BankrollLot, prop.Bankroll, coverage(w.NightHandleAt(id)), HouseEdge))
				add("draw", "Take money off the tables", 45, 0, w.DrawReadiness(id),
					fmt.Sprintf("$%d out of the $%d float and into your hands. It is the only way this room's winnings reach you, and every lot taken is action it can no longer attract.", BankrollLot, prop.Bankroll))
			}
			if w.Properties[id].Income > 0 {
				current := w.Mode(id)
				for _, m := range operatingModes {
					reason := ""
					if m.ID == current.ID {
						reason = "Already run this way"
					}
					add("operate:"+m.ID, m.Label, 0, 0, reason, m.Detail)
				}
			}
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
	if w.PressPlace(id) {
		fee := w.PressFee(id)
		detail := fmt.Sprintf("$%d. Restores up to 45 condition, currently %d of 100.", fee, w.DressCondition())
		if fee == 0 {
			detail = fmt.Sprintf("Your own people, at no charge. Restores up to 45 condition, currently %d of 100.", w.DressCondition())
		}
		add("press", "Have your clothes cleaned and pressed", PressMinutes, 0, w.PressReadiness(id), detail)
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
	w.SuspendedJob = nil
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
	} else if w.Random() < .78-float64(p.Armour)*.09 {
		w.Die("An attack at your residence caught you without warning or protection.")
	} else {
		p.Health = max(1, p.Health-w.Absorb(65))
		w.Ruin(55)
		w.Damage(40)
		worn := "Nobody warned you."
		if p.Armour > 0 {
			worn = "What you were wearing took the worst of it."
		}
		w.Log("You survived by inches", "The attackers leave you wounded. "+worn+" You need rest and protection.", "danger")
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
		next := min(end, (w.Minute/720+1)*720)
		if w.NextPressure > 0 {
			next = min(next, max(w.Minute+1, w.NextPressure))
		}
		for _, task := range w.Tasks {
			next = min(next, max(w.Minute+1, task.Due))
		}
		for _, contract := range w.Contracts {
			if contract.Life == w.Life {
				next = min(next, max(w.Minute+1, contract.Due))
			}
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
				prop.Carry += float64(prop.Income*prop.Condition*elapsed) * operatingMode(prop.Mode).Take * w.Capacity(id) / 6000
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
		// Organizations reconsider each other twice a day; their books settle once.
		if w.Minute%720 == 0 {
			w.FactionTurn()
			w.MarketPrices()
			w.ConsiderRobbery()
		}
		w.ResolveContracts()
		if w.Minute%1440 == 0 {
			w.FamilyDay()
			w.BusinessDay()
			w.ContrabandDay()
			w.PoliceDay()
			w.PeopleDay()
			w.OperationsDay()
			w.StillDay()
			w.CasinoDay()
			w.CarDay()
			w.DressDay()
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
		locs = append(locs, map[string]any{"id": l.ID, "name": l.Name, "type": l.Type, "district": l.District, "x": l.X, "y": l.Y, "cost": l.Cost, "blurb": l.Blurb, "owner": prop.Owner, "holder": w.HolderName(l.ID), "staff": prop.Staff, "supply": prop.Supply, "trouble": prop.Trouble, "still": prop.Still, "bankroll": prop.Bankroll, "handle": w.NightHandleAt(l.ID), "capacity": w.Capacity(l.ID), "condition": prop.Condition, "income": prop.Income, "owned": w.Own(l.ID), "locked": l.District > w.District, "actions": w.Actions(l.ID)})
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
	return map[string]any{"id": w.ID, "version": w.Version, "revision": w.Revision, "life": w.Life, "minute": w.Minute, "player": w.Player, "district": w.District, "factions": w.Factions, "npcs": w.People(), "locations": locs, "event": scene, "history": history, "dead": w.Dead, "tasks": w.Tasks, "director": w.Director, "last_result": w.LastResult, "daily_cost": w.DailyCost(), "income": income, "security": w.Guard(), "opportunity": w.NextOpportunity(), "known_threats": w.KnownThreats(), "business_truces": w.ActiveBusinessTruces(), "conflicts": w.PublicConflicts(), "goods": w.Goods, "arms": w.ArmsDescription(), "appearance": w.AppearanceDescription(), "vehicle": w.VehicleDescription(), "offshore": map[string]any{"balance": w.Offshore, "reachable": w.Player.Offshore}, "newspaper": w.Edition(), "arrangements": w.PendingArrangements()}
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
	// Leave room to prepare for discovered danger. Hidden plots must not
	// change offer timing and indirectly reveal themselves.
	if !w.Player.Alive || w.Event != nil || w.SuspendedJob != nil || len(w.KnownThreats()) > 0 {
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
	for id, effect := range SituationalEffects() {
		catalog[id] = effect
	}
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
	scene := &Scene{Target: p.Location, Operation: p.Operation, ID: ID(), Title: p.Title, Beneficiary: p.Beneficiary, Body: p.Body, Speaker: p.Speaker, Kind: "proposal", Source: "local-ai", Minute: w.Minute, Effect: fx, Outcome: operationOutcome(p.Operation), Choices: []Choice{{ID: "accept", Label: operationLabel(p.Operation), Detail: fmt.Sprintf("$%d · %d minutes · +%d respect · +%d heat", fx.Reward, fx.Minutes, fx.Respect, fx.Heat) + " · At 15 heat, police may stop completion." + politicalDetail}, {ID: "decline", Label: "Decline the arrangement", Detail: "No cost or time."}}}
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

// Every operation needs words for the button the player presses and for what
// the ledger records afterwards. A missing entry rendered as an unpressable
// blank choice when the conflict-derived operations were added.
func operationLabel(operation string) string {
	labels := map[string]string{
		"courier": "Deliver the package", "mediation": "Mediate the dispute", "collection": "Collect the payment",
		"escort": "Travel with it", "warning": "Deliver the message", "recovery": "Get it out", "settlement": "Settle the matter",
		"supply": "Fetch what it needs", "distribution": "Move the stock on",
	}
	if label, ok := labels[operation]; ok {
		return label
	}
	return "Take the work"
}

func operationOutcome(operation string) string {
	outcomes := map[string]string{
		"courier":      "You delivered the sealed package and reported back.",
		"mediation":    "You completed the requested mediation without violence.",
		"collection":   "You collected the agreed payment and reported back.",
		"escort":       "You saw it across the city and handed it over intact.",
		"warning":      "You delivered the message in person and walked out.",
		"recovery":     "You got it out without a confrontation.",
		"settlement":   "You put the matter in front of the new leadership and settled it.",
		"supply":       "You brought back what the business needed.",
		"distribution": "You moved the stock across the city and handed it over.",
	}
	if outcome, ok := outcomes[operation]; ok {
		return outcome
	}
	return "You completed the arrangement."
}
