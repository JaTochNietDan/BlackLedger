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
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type is how a room is drawn and what kind of place it is to stand in.
	// Kind is what business is run there, and it is what a trade belongs to:
	// how many hands it needs, what it runs on, and what goes wrong in it.
	//
	// They are separate because a trade used to be keyed by street address, so
	// the city could hold exactly one laundry and one casino — a second would
	// have needed its own copy of the same rules, and the two would have
	// drifted apart the first time either was touched. Empty on anything
	// nobody runs a business out of.
	Type     string `json:"type"`
	Kind     string `json:"kind,omitempty"`
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
	// Explosives in hand. Absent in saves from before anybody could buy any.
	Charges int `json:"charges,omitempty"`
	// Standing arrangements with people in the building. Absent in saves from
	// before there was a building.
	Retainers []string `json:"retainers,omitempty"`
	// What the player has asked about lately, and when what they were told
	// stops being current. Absent in saves from before anybody could ask.
	Enquiries map[string]int `json:"enquiries,omitempty"`
	// Whose organization the player answers to, and how much work they have
	// done for them. Absent for somebody who answers to nobody.
	// What is in the tank, and when it was last touched. The minute is what
	// tells a car nobody has ever filled from a car somebody ran dry — without
	// it, settling the world would hand a free tank to anybody who had just
	// used theirs up.
	Fuel    int    `json:"fuel,omitempty"`
	Fuelled int    `json:"fuelled,omitempty"`
	Serves  string `json:"serves,omitempty"`
	Service int    `json:"service,omitempty"`
	// The face the player picked for themselves, one-based into the cast, and
	// nothing at all if they never picked. Nobody in this city is gendered by
	// the rules, so a face is otherwise dealt out by a hash of a name — which
	// is fair to the cast and can still hand somebody a portrait they do not
	// recognise as themselves. This is the one place a player is allowed to
	// overrule the city about who they are looking at.
	Face int `json:"face,omitempty"`
	// When the police let go, and what they took them in for. Absent for
	// anybody who is not inside.
	HeldUntil int    `json:"held_until,omitempty"`
	HeldFor   string `json:"held_for,omitempty"`
	// The minute the door shut, so the morning it opens can say what the city
	// did in between. A sentence runs from two days to twelve — three times
	// the longest trip out of Bellwether — and the release record used to say
	// only that whatever it cost "happened while you were not there to watch
	// it", which is the same empty claim the trips made.
	HeldFrom int `json:"held_from,omitempty"`
	// When the paper last carried something the player put there. Absent for
	// anybody who has never been able to.
	LastPress int `json:"last_press,omitempty"`
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
	// When they died, and whether the paper has carried an obituary yet. The
	// city recorded that somebody was dead and never when, so nothing could ask
	// "who died yesterday" — which is the question a paper asks every morning.
	DiedAt     int  `json:"died_at,omitempty"`
	Remembered bool `json:"remembered,omitempty"`
	// When the police let this one go. Absent for anybody who is not inside,
	// which is everybody in a save written before anybody could be taken in.
	Held int `json:"held,omitempty"`
	// What this person holds against the player personally. Grudges are
	// between people in the city and deliberately do not reach the protagonist;
	// this is the one number that does.
	Sore   int    `json:"sore,omitempty"`
	SoreAt string `json:"sore_at,omitempty"`
	// Where they are walking to, when they get there, and why. Location stays
	// where they set off from until they arrive, so everything that reasons
	// about where a person belongs keeps working while they are on the street.
	Heading string `json:"heading,omitempty"`
	Arrives int    `json:"arrives,omitempty"`
	Errand  string `json:"errand,omitempty"`
	// Where this person's day is, so the evening has somewhere to send them
	// back from. Taken the first time they are seen standing somewhere in the
	// daytime, which is how saves written before the city had a shift acquire
	// one.
	Post string `json:"post,omitempty"`
	// What this person actually has. It used to be worked out at the moment
	// somebody robbed them, from their rank and their family's money, so the
	// same man could be robbed every day of the year and be carrying the same
	// amount every time. Absent in saves written before people had money, which
	// reads as nothing — Pockets settles a person who has never been paid onto
	// what their standing would carry, so an old save behaves as it did.
	Purse int `json:"purse,omitempty"`
	// What they drive, as a tier, and nothing if they walk. The player was the
	// only person in this city who owned a car, which meant the forecourt had
	// one customer, nothing could be stolen off anybody, and a garage had
	// nothing to repair. Absent in saves written before the city drove.
	Car int `json:"car,omitempty"`
	// The minute they got it, so somebody who has never had a car can be told
	// from somebody whose car was taken. Without it, settling a city would
	// quietly hand a replacement to anybody who had just lost one.
	Drove int `json:"drove,omitempty"`
	// Whether what they drive is standing there with the glass out of it. A
	// thief working a street does not take every car on it; the rest are left
	// needing a garage, and that work is what a garage actually lives on.
	Hurt bool `json:"hurt,omitempty"`
	// Whether they are running on fumes. A car that needs petrol is a reason to
	// be at a forecourt, which is what makes one worth holding.
	Dry bool `json:"dry,omitempty"`
	// The minute this person was last paid. It exists so that having nothing
	// can be told from never having been given anything: a man robbed down to
	// nothing must not be quietly refilled by the same pass that settles a
	// newcomer.
	Paid int `json:"paid,omitempty"`
	// When they actually set off, which is not when they decided to. A shift
	// change that empties every building at the same minute puts most of the
	// city on the street at once; people leave over a few hours instead. Zero
	// once they have gone, so it doubles as the record of having left.
	Sets int `json:"sets,omitempty"`
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
	// The strength the paper last reported. Families rise and fall in the
	// simulation constantly and the player had no way to feel any of it without
	// opening a screen and comparing numbers to numbers they did not write down.
	Reported int `json:"reported,omitempty"`
	// Whether this family failed to pay its people yesterday. Absent in saves
	// written before a family had a bill at all, which reads as a family that
	// met it — and for those saves that is true, because there was none.
	Short int `json:"short,omitempty"`
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
	// When the till was last taken, and when a car was last taken apart in the
	// street outside. A place remembers being robbed: it is what stops the same
	// drawer being emptied every forty-five minutes.
	Robbed   int `json:"robbed,omitempty"`
	Stripped int `json:"stripped,omitempty"`
	// Limit is the most this room will take on one bet, set by whoever holds
	// it. Nothing means nobody has said, and the room runs on what it is worth.
	Limit int `json:"limit,omitempty"`
	// What has been fitted to a residence. It belongs to the building rather
	// than to whoever lives there. Absent in saves written before that was
	// possible, which is a place with nothing in it.
	Comforts []string `json:"comforts,omitempty"`
	// A room under the floor with crates in it, and how many. Absent
	// everywhere else and in saves from before there was such a room.
	Armoury bool `json:"armoury,omitempty"`
	Crates  int  `json:"crates,omitempty"`
	// The legitimate trade a place has built up, and whether somebody
	// respectable has standing work with it. Absent in saves written before a
	// business had an outside as well as an inside.
	Custom int  `json:"custom,omitempty"`
	Order  bool `json:"order,omitempty"`
	// Whoever answers to the player and is standing on this door. Absent
	// everywhere nobody is.
	Posted string `json:"posted,omitempty"`
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

// TaskHome remembers where somebody was standing when they were sent to do
// something, so the city can put them back rather than leaving them wherever
// the work happened to be.
type TaskHome struct {
	Task   string `json:"task"`
	Person string `json:"person"`
	Where  string `json:"where"`
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
	// How many times this happened today. Measured on a real fifty-eight day
	// campaign: forty-nine of the sixty records the ledger holds were the same
	// sentence — "Mara pays $45. A small favor, completed without questions."
	// Repetition had not merely made the log unreadable, it had pushed every
	// notable thing that ever happened out of the archive.
	Count int `json:"count,omitempty"`
}
type Death struct {
	Name   string `json:"name"`
	Minute int    `json:"minute"`
	Life   int    `json:"life"`
	Cause  string `json:"cause"`
	// Estate is what the city calls whatever they left, when somebody was left
	// to hold it. What became of what a person built is part of the record of
	// their death; it used to be worked out and thrown away.
	Estate string `json:"estate,omitempty"`
	// Whether the paper has already carried an obituary. One column, the
	// morning after, and never again — a paper that runs the same obituary
	// every day is a paper nobody believes.
	Remembered bool `json:"remembered,omitempty"`
}
type Choice struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Cost   int    `json:"cost"`
	// The terms of this way of doing the job, as figures rather than as prose.
	// A scene exists to make the player compare two or three approaches, and a
	// run-on sentence is the one form those numbers cannot be compared in.
	// Whatever is here must be what the command will actually apply.
	Pay     int `json:"pay,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	Respect int `json:"respect,omitempty"`
	Heat    int `json:"heat,omitempty"`
	// Reason is why this cannot be taken, in the city's words. Only the core
	// knows; the interface used to guess and say "Not enough cash" for every
	// refusal it was shown.
	Reason   string `json:"reason,omitempty"`
	Disabled bool   `json:"disabled"`
}

// terms puts an effect's figures on a choice, so the button quotes the price
// the job will pay rather than a sentence written beside it.
func (c Choice) terms(fx Effect) Choice {
	c.Pay, c.Minutes, c.Respect, c.Heat = fx.Reward, fx.Minutes, fx.Respect, fx.Heat
	return c
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
	// Conditions are what holds however the job is done — where it is, what
	// the police do at 15 heat, who gains standing by it. Said once, above the
	// choices, rather than copied onto each of them.
	Conditions string `json:"conditions,omitempty"`
	Kind       string `json:"kind"`
	Source     string `json:"source"`
	Minute     int    `json:"minute"`
	Effect     Effect `json:"effect"`
	Outcome    string `json:"outcome"`
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

// VisualCue is the core saying what a moment looked like and where. It is not
// choreography — the interface decides how a thing is played — but what
// happened, to whom, and what the paper will say about it belongs here.
// CueActor is somebody who was in a moment.
type CueActor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VisualCue struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Target  string `json:"target"`
	Caption string `json:"caption"`
	// Headline is what the Herald carries about it, so the paper can arrive
	// after the scene rather than instead of it.
	Headline string `json:"headline,omitempty"`
	// Actors are the people who were in it. Their ids travel with their names
	// because a face is drawn from an id, and a scene about somebody that
	// cannot show them is a scene about nobody.
	Actors []CueActor `json:"actors,omitempty"`
	// Gravity is how much it is worth stopping for, so the interface never has
	// to guess which of five things in one command is the one to show.
	Gravity int `json:"gravity,omitempty"`
	Minute  int `json:"minute,omitempty"`
}

// Result is what just happened, which is the most important thing on the
// screen the moment after a player commits to something. It carried the records
// the command wrote and nothing about the command itself — not what was done,
// not what it cost — so the interface could only ever show the city's minutes
// and never the player's decision.
type Result struct {
	Cues []VisualCue `json:"cues,omitempty"`
	// Who came into or left the room while this was being done.
	Comings []Coming `json:"comings,omitempty"`
	// Action is what the player chose, in their own words, and Kind its id.
	Action  string   `json:"action,omitempty"`
	Kind    string   `json:"kind,omitempty"`
	From    string   `json:"from_location"`
	To      string   `json:"to_location"`
	Elapsed int      `json:"elapsed"`
	Records []Record `json:"records"`
	// What it cost or paid, measured across the whole command rather than
	// announced by whichever rule happened to move a number.
	Cash    int `json:"cash"`
	Respect int `json:"respect"`
	Heat    int `json:"heat"`
	Health  int `json:"health"`
}
type World struct {
	SuspendedJob *SuspendedJob `json:"suspended_job,omitempty"`
	VisualCues   []VisualCue   `json:"-"`
	// What walked in or out of the room the player is standing in during this
	// command. Like VisualCues, it belongs to the command rather than the save.
	Comings        []Coming             `json:"-"`
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
	Offshore  int        `json:"offshore,omitempty"`
	Contracts []Contract `json:"contracts,omitempty"`
	// Standing work an organization has asked for. Tied to one protagonist:
	// nobody inherits somebody else's obligations.
	Commissions []Commission `json:"commissions,omitempty"`
	// What people in this city hold against each other. Absent in saves from
	// before anybody remembered anything.
	Grudges []Grudge `json:"grudges,omitempty"`
	// Understandings the player has with organizations. Tied to one
	// protagonist: nobody inherits somebody else's friends.
	Pacts []Pact `json:"pacts,omitempty"`
	// Money out with somebody's name on it. Tied to one protagonist, because a
	// debt is owed to a man and not to an address.
	Loans []Loan `json:"loans,omitempty"`
	// How hard the city as a whole is looking, which is nobody's attention in
	// particular and everybody's problem. Absent in older saves, which is a
	// city that has not been counting.
	Attention int `json:"attention,omitempty"`
	// Where the player is sitting down to play, if anywhere. A sitting is a
	// decision the world knows about rather than a screen the interface opens
	// on its own, because the table's memory — the last hand, the last spin,
	// where the drums stopped — is saved state, and somebody walking up to it
	// should not be shown a game that started without them.
	Seated string `json:"seated,omitempty"`
	// Where the player is, when that is two hundred miles outside the city. A
	// journey between two addresses and a week in Halloway both park the player
	// in "transit", and they are not the same thing at all: one is a walk
	// somebody can follow you down, the other is out of everybody's reach.
	Abroad string `json:"abroad,omitempty"`
	// A hand on the table that has not been settled. Absent whenever nobody is
	// sitting at one, which is nearly always.
	Hand *TableHand `json:"hand,omitempty"`
	// Where the drums stopped on the last pull, so a machine can be drawn
	// rather than described.
	Reels *Pull `json:"reels,omitempty"`
	// The last turn of the wheel. Absent in saves written before the room had
	// one, which reads as a wheel nobody has played.
	Spin  *Spin   `json:"spin,omitempty"`
	News  []Story `json:"news,omitempty"`
	Plots []Plot  `json:"plots"`
	Tasks []Task  `json:"tasks"`
	// Where the people doing those tasks were standing when they were sent.
	Homes      []TaskHome `json:"homes,omitempty"`
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
	// Amount is money the player typed rather than picked: what to put down on
	// a bet, what to put behind the tables, what to wire out of the city. Zero
	// means they did not name one, and whatever the action would have done with
	// a fixed lot still applies.
	Amount int `json:"amount,omitempty"`
}

// Sum is a figure the player types on an action card, bounded by what the core
// knows they can actually move.
type Sum struct {
	Least  int    `json:"least"`
	Most   int    `json:"most"`
	Preset int    `json:"preset"`
	Label  string `json:"label"`
}

type Action struct {
	// Group is what this action is for, so the interface can offer ninety of
	// them in an order a person can navigate. Set by the core, never guessed
	// by the presentation.
	Group string `json:"group"`
	// Anywhere marks work that belongs to the player rather than to the room
	// they happen to be standing in. Putting a price on a name, reaching an
	// understanding with a family and asking around about one are all done
	// through other people; none of them needs the player to walk to a
	// particular counter first, and the city's own screens for people and for
	// families are where they belong. The room's own panel leaves them out.
	Anywhere bool `json:"anywhere,omitempty"`
	// Subject is whose name this action is about, when it is about somebody
	// standing here rather than about the premises. Lending a man money is not
	// the same kind of thing as repairing a roof, and a list that shows them
	// as two identical cards has thrown away what the player needs to decide.
	Subject string `json:"subject,omitempty"`
	ID      string `json:"id"`
	Label   string `json:"label"`
	Minutes int    `json:"minutes"`
	// Away is time the action spends that Minutes does not, because the effect
	// runs the clock itself. Only the trips out of the city do this: they run
	// day by day so that what happens while the player is gone happens to a
	// city they are not standing in. Minutes is what the command layer spends
	// and had to stay zero; the panel prints the cost of an action from it, so
	// the longest actions in the game showed no time at all beside "Rest for
	// four hours · 240 min".
	Away int `json:"away,omitempty"`
	// Cost is what the command layer charges. An effect that pays its own fee
	// has to declare nothing here or the money would be taken twice.
	Cost int `json:"cost"`
	// Asks is that fee, for the panel. The price is printed from Cost, so work
	// that charges its own way showed no price at all: calling two families to
	// a room takes $220 and the button said nothing about money, with the
	// figure only in the description. Same shape as the trips and their days.
	Asks int `json:"asks,omitempty"`
	// Sum, when set, means the player types how much rather than accepting a
	// fixed lot. The panel puts a number field on the card and sends what is in
	// it as Command.Amount; the core still decides what an empty field means.
	Sum      *Sum   `json:"sum,omitempty"`
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
	// The same thing happening again today is not a new thing to say. The
	// record is removed and re-appended rather than updated in place, so it
	// keeps its position in the order of events and takes a fresh id — the
	// result panel after an action is built by diffing record ids, and a
	// collapsed repeat has to still read as something that just happened.
	day := w.Minute / 1440
	for i := len(w.History) - 1; i >= 0; i-- {
		r := w.History[i]
		if r.Minute/1440 != day || r.Life != w.Life {
			break
		}
		if r.Title == title && r.Text == text {
			again := r.Count + 1
			if again < 2 {
				again = 2
			}
			w.History = append(w.History[:i], w.History[i+1:]...)
			w.History = append(w.History, Record{ID(), w.Minute, w.Life, title, text, kind, again})
			return
		}
	}
	w.History = append(w.History, Record{ID(), w.Minute, w.Life, title, text, kind, 1})
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
	w.Factions = []Faction{
		{ID: "bellandi", Name: "Bellandi Family", Leader: "Vittorio Bellandi", Power: 90, Cash: 8000, Peak: 90, Reported: 90},
		{ID: "russo", Name: "Russo Outfit", Leader: "Elena Russo", Power: 58, Cash: 4500, Peak: 58, Reported: 58},
	}
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
		income := PlaceIncome[p.ID]
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
	w.ensureOfficials()
	w.Populate()
	w.Goods = newGoods()
	// The two established families are already rivals when the player arrives.
	w.Antagonize("bellandi", "russo", 50)
	w.SettlePurses()
	w.SettleCars()
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

// Watchers is the people who could raise an alarm: hired security, and whoever
// is around a better address. A door cannot shout, which is why it is not
// counted here.
func (w *World) Watchers() int {
	n := w.Player.Security
	switch w.Player.Home {
	case "apartment":
		n++
	case "estate":
		n += 2
	}
	return n
}

// Guard is everything standing between the player and somebody at the door:
// the people, and what the building itself does.
func (w *World) Guard() int {
	n := w.Player.Security
	if w.Fitted("door") {
		n++
	}
	switch w.Player.Home {
	case "apartment":
		n++
	case "estate":
		n += 2
	}
	return n
}
func (w *World) DailyCost() int {
	return HomeRent(w.Player.Home) + 10*w.Player.Security + 12*len(w.Player.Crew) + w.Wages() + w.CarUpkeep() + w.ComfortUpkeep() + w.RetainerCost() + w.MemberWages() + w.PactCost()
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
		// Most work aimed at a person carries their id after the colon, so the
		// core already knows who it is about without anybody restating it.
		subject := ""
		if at := strings.IndexByte(id, ':'); at >= 0 {
			if who := id[at+1:]; w.NPC(who) != nil {
				subject = who
			}
		}
		out = append(out, Action{Group: GroupOf(id), Subject: subject, ID: id, Label: label, Minutes: minutes, Cost: cost, Disabled: reason != "", Reason: reason, Detail: detail, Target: l.ID})
	}
	// about names the person the action just added is aimed at, for the ones
	// that put a name in the label rather than in the id.
	about := func(who string) {
		if len(out) > 0 {
			out[len(out)-1].Subject = who
		}
	}
	// away offers work whose effect runs the clock itself, so the panel can
	// print how long it takes without the command layer spending it twice.
	away := func(id, label string, minutes int, reason, detail string) {
		add(id, label, 0, 0, reason, detail)
		if len(out) > 0 {
			out[len(out)-1].Away = minutes
		}
	}
	// asks is the same for money: work that pays its own fee, so the panel can
	// print the price without the command layer taking it a second time.
	asks := func(id, label string, minutes, price int, reason, detail string) {
		add(id, label, minutes, 0, reason, detail)
		if len(out) > 0 {
			out[len(out)-1].Asks = price
		}
	}
	// sum puts a typed figure on the action just added: how much goes behind the
	// tables, how much is wired out. Work that moves money the player names
	// rather than a lot somebody chose for them.
	sum := func(least, most, preset int, label string) {
		if len(out) == 0 {
			return
		}
		out[len(out)-1].Sum = &Sum{Least: least, Most: max(least, most), Preset: min(max(preset, least), max(least, most)), Label: label}
	}
	// anywhere marks the action just added as work the room has nothing to do
	// with, so the interface can put it where it belongs rather than in the
	// list of things you can do at a fishmonger's.
	anywhere := func() {
		if len(out) > 0 {
			out[len(out)-1].Anywhere = true
		}
	}
	need := func(b bool, s string) string {
		if b {
			return s
		}
		return ""
	}
	// Inside, there are three ways out and no fourth. Everywhere else on the
	// map offers nothing, because a man in a cell cannot go to any of it.
	if w.Held() {
		if l.ID != "precinct" {
			return out
		}
		days := w.DaysLeft()
		add("sit_out", fmt.Sprintf("Do the %d days", days), 0, 0, "",
			fmt.Sprintf("The city runs without you and hands you the books afterwards. Worth %d respect for not saying anything.", days*ServedRespect))
		add("lawyer", fmt.Sprintf("Pay somebody who knows the clerks ($%d)", w.LawyerFee()), 90, 0, w.LawyerReadiness(),
			fmt.Sprintf("$%d a day for what is left of it. You walk out today having earned nothing and lost nothing but the money.", LawyerDaily))
		add("talk", "Explain who else was involved", 60, 0, "",
			fmt.Sprintf("Out the same afternoon. Costs %d respect, %d standing with every organization in the city, and everybody who works for you.", TalkedRespect, TalkedGoodwill))
		return out
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
		// What can be done about a place without being in it. An order given to
		// somebody else is given wherever you happen to be standing; walking
		// across the city to tell your own crew to go somewhere else is a
		// journey, not a decision.
		out = append(out, w.ordersAbout(l)...)
		return out
	}
	// Work that belongs to the player rather than to a room. All of it used to
	// live at the exchange, which meant crossing the city to tell your own crew
	// to move on somebody, or to put a price on a name the city already knows —
	// none of which happens at a counter in a building. It is offered wherever
	// the player is standing and the interface files it under people and
	// families rather than under the premises.
	add("contract", "Ask about a name", 30, 0,
		need(p.Contacts < 1, "Build a contact who will carry this"),
		"Put a price on somebody. What it costs depends on who they are and who does the work. A failed attempt can be traced back to you.")
	anywhere()
	add("investigate", "Ask about threats", 45, 30, "", "Investigate existing threats. Evidence is not a guarantee of safety.")
	anywhere()
	// What the police think of you follows you around; neither of these happens
	// at a counter. A detective who will lose some paperwork does not take the
	// money across a desk at the exchange, and keeping your head down is not
	// somewhere you go.
	asks("bribe", "An understanding with the detective", 45, w.BribeCost(), w.BribeReadiness(),
		fmt.Sprintf("$%d to Detective Harlow to lose some paperwork. Clears attention now and buys nothing later. Above %d heat nobody will be seen taking it.", w.BribeCost(), BribeCeiling))
	anywhere()
	add("lie_low", "Keep a low profile", 120, 15, "", "Lose 10 heat. Time still passes for rivals and businesses.")
	anywhere()
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() {
			continue
		}
		asks("enquire:"+f.ID, "Ask around about "+f.Name, EnquiryMinutes, EnquiryCost, w.EnquiryReadiness(f.ID),
			fmt.Sprintf("$%d in the right pockets. What comes back is current for about a week. You are at %d of 3 on them as it stands.", EnquiryCost, w.Intelligence(f.ID)))
		anywhere()
		if w.Allied(f.ID) {
			add("break:"+f.ID, "End the understanding with "+f.Name, 30, 0, "",
				fmt.Sprintf("Stops the $%d a day and the quarrels that come with it. They will remember that you did it first.", PactTribute))
			anywhere()
			continue
		}
		asks("pact:"+f.ID, "Reach an understanding with "+f.Name, PactMinutes, PactOpening, w.PactReadiness(f.ID),
			fmt.Sprintf("$%d to open and $%d a day. Neither of you moves on the other, they may answer when somebody comes for you, and every quarrel of theirs becomes yours.", PactOpening, PactTribute))
		anywhere()
	}
	switch id {
	case "precinct":
		for _, n := range w.OwnPeople() {
			if !w.Inside(n) {
				continue
			}
			days := (n.Held - w.Minute + 1439) / 1440
			// The ledger entry this writes has always counted correctly; the
			// description beside it said "for the 1 days still on them".
			asks("bail:"+n.ID, "Bail out "+n.Name, 60, days*BailDaily, w.BailReadiness(n.ID),
				fmt.Sprintf("$%d for the %s still on them. They come out owing you, which is not the same as being grateful.",
					days*BailDaily, plainly(days, "day", counted(days, "day", "days"))))
		}
	case "bar":
		add("courier", "Carry a discreet envelope", CourierMinutes, 0, "",
			fmt.Sprintf("Earn $%d and %d respect. A reliable introduction to the neighborhood.", CourierPay, CourierRespect))
		add("contact", "Buy Mara a coffee", 30, 10, need(p.Contacts >= 5, "Your information network is fully developed"), "Build trust and an information network. Contacts may warn you of trouble.")
		about(w.HolderID("fixer"))
		if q, ok := w.OpenQuarrel(); ok {
			warning := "You would be standing between them."
			if q.Suspected {
				warning = "Somebody has already told you that one of them is not coming to talk."
			}
			asks("sitdown", "Call "+q.A.Name+" and "+q.B.Name+" to a room", SitdownMinutes, SitdownFee, w.SitdownReadiness(),
				fmt.Sprintf("$%d for the room and the guarantees, paid whether or not anybody agrees to anything. The only thing in this city that ends a war without either side losing it. %s", SitdownFee, warning))
		}
		reason := need(p.Respect < PremisesRespect, fmt.Sprintf("Earn %d respect first", PremisesRespect))
		if len(p.Crew) > 0 {
			// Name the man on the books. "Leo is already in your crew" was
			// printed whoever was actually on them, which is how a crew of one
			// dead Leo Carver went unread for fifty-eight days.
			reason = p.Crew[0].Name + " is already in your crew"
		}
		// The button named Leo Carver whatever had happened to him. Nobody
		// holds a job for ever in this city: when the man who drives is dead
		// the role goes to somebody else by the end of the week, and until it
		// does there is nobody to hire. Naming whoever actually drives makes
		// this an action about a person, so the same question is asked of it
		// as of everything else aimed at one — which is what stops it
		// offering to recruit a corpse.
		driver := w.HolderID("driver")
		if driver == "" {
			driver = "leo"
		}
		hand := "Leo Carver"
		if n := w.NPC(driver); n != nil {
			hand = n.Name
		}
		add("recruit", "Recruit "+hand, 30, 90, reason, "A driver and collector. $12 daily wages; loyalty matters.")
		about(driver)
	case "garage":
		add("audience", "Request an audience with Russo", 45, 0, w.AudienceReadiness(id), "Discuss your standing with the Russo Outfit.")
		fee := w.ServiceFee()
		// With no car there is no condition to state, and "currently 0 of 100"
		// reads as a wreck in the yard rather than as nothing at all.
		service := "What a garage does, once there is something of yours in it. Restores up to 55 condition."
		if p.Car > 0 {
			service = fmt.Sprintf("$%d. Restores up to 55 condition, currently %d of 100. Below %d it is worth nothing to you.", fee, w.CarCondition(), Wreck)
			if fee == 0 {
				service = fmt.Sprintf("Your own people, at no charge. Restores up to 55 condition, currently %d of 100.", w.CarCondition())
			}
		}
		add("service", "Have the car worked on", CarServiceMinutes, 0, w.ServiceReadiness(id), service)
	case "docks":
		add("dockwork", "Work the night cargo", 90, 0, "", "Earn $75 and 1 respect. Small chance of a work injury.")
		if next, ok := nextArmament(weapons, p.Weapon); ok {
			asks("arms:weapon", "Buy "+next.Label, 45, next.Cost, w.ArmsReadiness("weapon"),
				fmt.Sprintf("$%d. %s Improves your odds when violence is your idea. A search takes it.", next.Cost, next.Detail))
		}
		asks("charge", "Buy a charge off a boat", ChargeMinutes, ChargeCost, w.ChargeReadiness(),
			fmt.Sprintf("$%d. Not a message: a declaration. Wrecks a business outright, kills whoever was standing in it about a third of the time, and cannot be mistaken for anything else. Draws %d police attention a day while you hold it, and a search that finds it is a prosecution rather than a fine.", ChargeCost, ChargeHeat))
		if next, ok := nextArmament(armour, p.Armour); ok {
			asks("arms:armour", "Buy "+next.Label, 45, next.Cost, w.ArmsReadiness("armour"),
				fmt.Sprintf("$%d. %s Reduces what a beating costs you. A search takes it.", next.Cost, next.Detail))
		}
	case "herald":
		for _, o := range officials {
			if o.Place() != id {
				continue
			}
			if w.Retained(o.ID) {
				add("release:"+o.ID, "Stop paying "+o.Name, 30, 0, "",
					fmt.Sprintf("Ends the arrangement and the $%d a day. Opening it again costs the opening payment over.", o.Retainer))
			} else {
				asks("retain:"+o.ID, "An arrangement with "+o.Name, OfficialMinutes, w.OfficialOpening(o), w.RetainerReadiness(o.ID),
					fmt.Sprintf("$%d to open and $%d a day after. %s They cut you loose above %d attention and keep the opening payment.", w.OfficialOpening(o), o.Retainer, o.Detail, w.OfficialCeiling(o)))
			}
		}
		spikeable := "Nothing in today's paper is about you."
		if n := len(w.Spikeable()); n > 0 {
			// One story is about you; two stories are. The noun and the verb
			// both have to follow the number.
			is := "are"
			if n == 1 {
				is = "is"
			}
			spikeable = fmt.Sprintf("%s in today's paper %s about you or about the police.",
				upper1(counted(n, "story", "stories")), is)
		}
		asks("spike", "Pull a story", SpikeMinutes, SpikeCost, w.SpikeReadiness(),
			fmt.Sprintf("%s The worst of them does not run, and the city's interest in it goes with it. About one time in seven somebody in that building notices and the arrangement is over.", spikeable))
		asks("puff", "A paragraph about a local businessman", PuffMinutes, PuffCost, w.PuffReadiness(),
			fmt.Sprintf("$%d for %d respect and %d off what the police think. The cheapest standing in this city and the only kind nobody was hurt for.", PuffCost, PuffRespect, PuffHeat))
		for i := range w.Factions {
			f := &w.Factions[i]
			if f.ID == w.PlayerOrganizationID() {
				continue
			}
			asks("smear:"+f.ID, "Run something about "+f.Name, SmearMinutes, SmearCost, w.SmearReadiness(f.ID),
				fmt.Sprintf("$%d. Every one of their places loses %d trade and the organization loses %d strength. A paper full of crime is a paper full of crime whoever it is about, so the whole city gets harder — and about one time in five they find out who paid for it.", SmearCost, SmearCustom, SmearPower))
		}
	case "market":
		add("deposit", "Wire money out of the city", 45, 0, w.DepositReadiness(0),
			fmt.Sprintf("The arrangement takes %d%% of whatever you send; $%d of a $%d wire arrives. It survives you, and whoever comes next can reach it if they can afford to.", DepositCut, DepositLot*(100-DepositCut)/100, DepositLot))
		sum(DepositLeast, p.Cash, DepositLot, "Sent out")
		asks("offshore_access", "Establish that the account is yours", AccessMinutes, AccessCost, w.AccessReadiness(),
			fmt.Sprintf("$%d in papers and a journey. Only worth it if there is enough out there to be worth reaching.", AccessCost))
		// An empty account has no sum to bring home. "Brings $0 back into the
		// city" is a figure that says nothing while looking like one.
		bringing := "Brings whatever is out there back into the city, where it can be taken from you."
		if w.Offshore > 0 {
			bringing = fmt.Sprintf("There is $%d out there. Whatever you bring home can be taken from you.", w.Offshore)
		}
		add("withdraw", "Bring money home", 45, 0, w.WithdrawReadiness(0), bringing)
		sum(min(WithdrawLeast, w.Offshore), w.Offshore, w.Offshore, "Brought home")
		if next, ok := nextAttire(p.Dress); ok {
			notice := "Nobody official looks twice at it."
			if next.Notice > 0 {
				notice = fmt.Sprintf("Dressing above your visible means draws %d police attention a day.", next.Notice)
			}
			asks("dress", "Be measured for "+lowerFirst(next.Label), 60, next.Cost, w.DressReadiness(),
				fmt.Sprintf("$%d. %s Worth %d presence while it is kept, and it wears. %s", next.Cost, next.Detail, next.Presence, notice))
		}
		for _, o := range officials {
			if o.Place() != id {
				continue // a man is arranged with where he actually is
			}
			if w.Retained(o.ID) {
				add("release:"+o.ID, "Stop paying "+o.Name, 30, 0, "",
					fmt.Sprintf("Ends the arrangement and the $%d a day. Opening it again costs the opening payment over.", o.Retainer))
				continue
			}
			asks("retain:"+o.ID, "An arrangement with "+o.Name, OfficialMinutes, w.OfficialOpening(o), w.RetainerReadiness(o.ID),
				fmt.Sprintf("$%d to open and $%d a day after. %s They cut you loose above %d attention and keep the opening payment.", w.OfficialOpening(o), o.Retainer, o.Detail, w.OfficialCeiling(o)))
		}
		for _, d := range destinations {
			// Minutes are zero here because Trip runs the days itself, a day at
			// a time, so that what happens in the city while the player is
			// away happens to a city the player is not standing in.
			away("trip:"+d.ID, "Travel to "+d.Name, d.Days*1440, w.TripReadiness(d.ID),
				fmt.Sprintf("%s %s $%d all in and %d days away. The city runs without you: businesses go unwatched, work you promised runs down, and anything arranged for you happens to an empty house. Attention falls %d a day while you are gone.", d.Blurb, d.Purpose, w.TripCost(d.ID), d.Days, d.Relief))
		}
	case "club":
		add("audience", "Request an audience", 45, 0, w.AudienceReadiness(id), "Discuss your standing with the Bellandi family.")
		add("provoke", "Demand protection money", 30, 0, "", "EXTREME RISK. Bellandi owns this casino. Challenging them can bring lethal retaliation.")
	}
	if w.Hand != nil && !w.Hand.Done && w.Hand.Place == id {
		add("hit", "Take another card", 5, 0, "",
			fmt.Sprintf("You are showing %d and the dealer is showing %d. Over %d and it is finished.", w.Hand.Player, w.Hand.Dealer, Bust))
		add("stand", "Stand on "+fmt.Sprint(w.Hand.Player), 5, 0, "",
			fmt.Sprintf("The dealer draws to %d and stands on %d. A tie gives your money back.", DealerStands-1, DealerStands))
	}
	// Cars are sold on a forecourt. This lived in the garage's own case, where
	// it was offered and then permanently refused by a rule that asks for a
	// forecourt — and the forecourt, which is where the rule points, never
	// offered it at all. A button and its rule have to be asking the same
	// question.
	if CarSource(id) {
		if next, ok := nextVehicle(p.Car); ok {
			hides := "Nothing to hide anything in."
			if next.Compartment > 0 {
				hides = fmt.Sprintf("A false floor a search will not find %d units under.", next.Compartment)
			}
			asks("car", "Buy "+lowerFirst(next.Label), 60, next.Cost, w.CarReadiness(),
				fmt.Sprintf("$%d, then $%d a day to keep on the road. %s Journeys take %d%% of the time they take on foot. %s A car outside is a thing witnesses describe.", next.Cost, next.Upkeep, next.Detail, int(next.Pace*100), hides))
		}
	}
	if Pumps(id) {
		short := FuelFull - w.Fuel()
		detail := fmt.Sprintf("$%d for what the tank is short. It reads %d of %d.", w.FuelFee(id), w.Fuel(), FuelFull)
		if w.Own(id) {
			detail = fmt.Sprintf("Your own pumps, at what the petrol cost you: $%d. It reads %d of %d.", w.FuelFee(id), w.Fuel(), FuelFull)
		}
		if p.Car == 0 {
			detail = "What a forecourt is for, once there is something of yours to put it in."
		} else if short == 0 {
			detail = fmt.Sprintf("It reads %d of %d. There is nowhere for it to go.", w.Fuel(), FuelFull)
		}
		asks("fill", "Fill the tank", FillMinutes, w.FuelFee(id), w.FillReadiness(id), detail)
	}
	// A seat, before any game. Sitting down is a decision the world knows
	// about: it clears the last hand, the last spin and the drums, so nobody
	// walks up to a table and finds a game somebody else was playing. Zero
	// minutes and zero money — the games spend both themselves.
	if Playable(id) && id == p.Location {
		if w.Seated == id {
			add("rise", "Get up and leave the tables", 0, 0, w.RiseReadiness(),
				"Ends the sitting. Nothing is left on the table for whoever sits down next.")
		} else {
			seat, what := "Sit down at the tables", "The cards, the wheel and the machines, played out in front of you until you get up."
			if !HasTables(id) {
				seat, what = "Play the machines", "Three drums and a handle, against the wall where they always are."
			}
			add("sit", seat, 0, 0, w.SitReadiness(id), what)
		}
	}
	// The machines. Not a table: no dealer, no floor, no minimum — just a wall
	// and whatever is in your pocket, which is why they are in rooms that are
	// not casinos at all.
	if HasMachines(id) {
		// One machine, and you decide what goes in it. Cost is zero because
		// PullHandle takes the money itself, and the amount rides on the
		// command rather than being one of two buttons.
		asks("pull", "Play the machines", 15, 0, w.PullReadiness(id, Stake{Amount: w.MachineStakeOrUsual(id, 0)}),
			fmt.Sprintf("Anything from $%d to $%d a pull on three drums of twenty. Three of a kind pays what is written on the machine, up to %d to 1 for the sevens, and a cherry on its own gives you your money back. The machine keeps about %d in every hundred that goes through it.",
				LeastStake, w.MachineLimit(id), sevenPays(), MachineEdge()))
	}
	if HasTables(id) && !w.Own(id) {
		// Cost stays nothing on both of these: each takes its own money, and
		// what it takes is the amount the player named rather than one of two
		// lots the room decided for them.
		bet := w.TableStakeOrUsual(id, 0)
		reason := w.TableReadiness(id, Stake{Amount: bet})
		if reason == "" && w.Hand != nil && !w.Hand.Done {
			// Where the hand is, the buttons to play it are right beside this
			// and the refusal needs to say nothing more. Anywhere else, the
			// refusal is all the player gets.
			reason = "You are in the middle of a hand"
			if w.Hand.Place != id {
				where := w.Hand.Place
				if place, ok := PlaceByID(w.Hand.Place); ok {
					where = place.Name
				}
				reason = "There is a hand of yours still on the table at " + where
			}
		}
		asks("play", "Sit in at the cards", 60, 0, reason,
			fmt.Sprintf("Anything from $%d to $%d a hand, played out a card at a time. The dealer draws to %d and stands on %d, a tie gives your money back, and going over is finished before the dealer plays at all. Above $%d the floor wants to know who you are.",
				LeastStake, w.TableLimit(id), DealerStands-1, DealerStands, HighTableMoney))
		asks("wheel", "Play the wheel", 45, 0, w.SpinReadiness(id, "red", Stake{Amount: bet}),
			fmt.Sprintf("Anything from $%d to $%d on one bet. Thirty-seven pockets and one of them is the nought, which belongs to the house and takes every bet on the outside. A number pays %d to 1, red or black and odd or even pay even money, a dozen pays 2 to 1. Every one of those is the true price.",
				LeastStake, w.TableLimit(id), 35))
	}
	// The house limit belongs to whoever holds the room.
	if w.Own(id) && (HasTables(id) || HasMachines(id)) {
		add("limit", "Set the house limit", 30, 0, w.LimitReadiness(id),
			fmt.Sprintf("The most this room will take on one bet, anywhere from $%d to $%d. It reads $%d at the tables and $%d at the machines. Bigger action pays better and loses worse.",
				HouseLimitFloor, HouseLimitCeiling, w.TableLimit(id), w.MachineLimit(id)))
	}
	// The trade, not the kind of room. The readiness function asks what is run
	// here; this asked what the room looked like, so a casino of the player's
	// own was refused a button it would have been allowed to press.
	if trade, runs := TradeOf(id); runs && trade.Cover > 0 && w.Own(id) {
		asks("launder", "Run takings through the books", 90, w.LaunderFee(id), w.LaunderReadiness(id),
			fmt.Sprintf("$%d to clear up to %d police attention through %s. Wears the premises, and the books need a day between rounds.", w.LaunderFee(id), w.launderCapacity(id), l.Name))
	}
	if prop := w.Properties[id]; prop != nil && prop.Income > 0 && !w.Own(id) {
		add("rob", "Take the day's cash yourself", 45, 0, w.RobberyReadiness(id),
			fmt.Sprintf("Walk out with what is in the till at %s. Your standing and whatever you are carrying improve the odds. A haul, police attention, and an owner who will work out who would dare. Going wrong means a beating, and it is yours.", l.Name))
		if hand, ok := w.CrewHands(); ok {
			reason := w.RobberyReadiness(id)
			if reason == "" {
				reason = w.DelegateReadiness()
			}
			add("rob:crew", "Send "+hand.Name+" for the till", 45, 0, reason,
				fmt.Sprintf("The same money and worse odds, because they bring their loyalty to it and not your name. %d less police attention on you and a fifth of the standing. Going wrong costs them %d loyalty, and sometimes more than that.", HandHeatRelief, HandLoyaltyCost))
		}
	}
	for _, g := range w.Goods {
		if !TradesAt(id, g.ID) {
			continue
		}
		// Cost is zero because Buy charges the lot itself; declaring it would
		// have the command layer charge it a second time.
		add("buy:"+g.ID, fmt.Sprintf("Buy %d %ss of %s", Lot, g.Unit, g.InBulk()), 30, 0,
			w.TradeReadiness(g.ID, "buy"),
			fmt.Sprintf("$%d for the lot, at $%d each today. Holding stock draws police attention every day until it is sold, and can be taken from you.", g.Price*Lot, g.Price))
		if held := w.Holding(g.ID); held > 0 {
			add("sell:"+g.ID, fmt.Sprintf("Sell %d %ss of %s", held, g.Unit, g.InBulk()), 30, 0,
				w.TradeReadiness(g.ID, "sell"),
				fmt.Sprintf("$%d each today, for $%d.", g.Price, g.Price*held))
		}
	}
	if mark, ok := w.StripTarget(id); ok {
		add("strip", "Take "+mark.Name+"'s car apart", 60, 0, w.StripReadiness(id),
			fmt.Sprintf("$%d for what comes off it. They will be walking, and they will know by morning that it was somebody. Every garage in the city has more work the more of this there is.", w.PartsWorth(mark)))
	}
	if mark, ok := w.MuggingTarget(id); ok {
		// People carry their own money now, so somebody can genuinely have
		// nothing — a man whose family has missed payday for a month. "About
		// $0 on them" is not a sum, it is a thing to say in words.
		carrying := fmt.Sprintf("About $%d on them.", w.Pockets(mark))
		if w.Broke(mark) {
			carrying = "Nothing worth taking, by the look of them."
		}
		add("mug", "Take what "+mark.Name+" is carrying", MuggingMinutes, 0, w.MuggingReadiness(id),
			fmt.Sprintf("%s Your standing improves the odds and makes you the person they describe afterwards: above %d presence they can name you. They will hold it against you either way, and so will %s.", carrying, RecognisedAt, w.factionName(mark.Faction)))
		about(mark.ID)
		if hand, ok := w.CrewHands(); ok {
			reason := w.MuggingReadiness(id)
			if reason == "" {
				reason = w.DelegateReadiness()
			}
			add("mug:crew", "Send "+hand.Name+" after "+mark.Name, MuggingMinutes, 0, reason,
				fmt.Sprintf("The same $%d and worse odds, and it is their face rather than yours. %s holds it against them instead.", w.Pockets(mark), mark.Name))
			about(mark.ID)
		}
	}
	// Going after somebody who is standing here. Two buttons, because they are
	// two different bets: your own hands and your own risk, or somebody else's
	// face and what happens to them if it goes wrong.
	for _, who := range w.PeopleHere(id) {
		mark, ok := w.StrikeTarget(who.ID)
		if !ok {
			continue
		}
		add("strike:"+mark.ID, "Go after "+mark.Name+" yourself", StrikeMinutes, 0, w.StrikeReadiness(mark.ID),
			fmt.Sprintf("Your own hands, and the best odds you can get: what you are carrying and what you are worth both count. It can kill you, a room with people in it remembers your face, and %s answers for %s either way.",
				w.factionOrStreet(mark), theirOrTheir(mark)))
		about(mark.ID)
		if hand, ok := w.CrewHands(); ok {
			reason := w.SendReadiness(mark.ID)
			add("send:"+mark.ID, "Send "+hand.Name+" after "+mark.Name, StrikeMinutes, 0, reason,
				fmt.Sprintf("Their face at the scene rather than yours, on worse odds than your own. Taken alive they are known to be yours and they give you up; sometimes they do not come back at all."))
			about(mark.ID)
		}
	}
	if prop := w.Properties[id]; prop != nil && prop.Income > 0 && !w.Own(id) && p.Charges > 0 {
		add("plant", "Put the charge under "+l.Name, PlantMinutes, 0, w.PlantReadiness(id),
			fmt.Sprintf("Wrecks %s, empties it of stock and staff, and kills somebody who worked there about a third of the time. The owner will know exactly what it was. Going wrong means it goes off with you under it.", l.Name))
	}
	if holder, ok := w.MoveTarget(id); ok {
		add("move", "Move on "+l.Name, MoveMinutes, 0, w.MoveOnReadiness(id),
			fmt.Sprintf("Commit your organization against %s. Your strength against theirs decides it, and a holding they cannot defend becomes yours. Driven off, it costs you somebody who went with you. Draws %d police attention and hardens the quarrel.", holder.Name, MoveHeat))
	}
	out = append(out, w.ordersAbout(l)...)
	if f, ok := w.SabotageTarget(id); ok {
		// The difference between these two buttons is whether you are there,
		// not whether the crew is: SabotageReadiness requires a crew either
		// way, the odds read their loyalty either way, and the failure text
		// says "You and Leo left without reaching anything" here against
		// "went in without you" beside it. The description said "Send your
		// crew", which is the other button — and that one is offered from
		// anywhere, by ordersAbout, because giving it does not put you in the
		// building. This one does, so it is only here.
		add("sabotage", "Move against "+l.Name+" yourself", 90, 0, w.SabotageReadiness(id),
			fmt.Sprintf("Go in with your crew against %s. Damages the property, weakens %s and costs you standing with them. They will retaliate, and a failed attempt injures you.", l.Name, f.Name))
		_ = f
	}
	// Somewhere that earns is somewhere that can be held, run and taken. This
	// asked for three addresses by name — the laundry, the garage and the
	// casino, which was every business in the city on the day it was written —
	// so a restaurant, a butcher, a cab company, a scrapyard, a forecourt and
	// two filling stations could be walked into and robbed and never bought.
	// The rule underneath was general the whole time. Only the button was not.
	if prop := w.Properties[id]; prop != nil && prop.Income > 0 {
		if w.Own(id) {
			add("inspect", "Review the books", 0, 0, "", "Read current income and repair needs without advancing time.")
			if trade, running := TradeOf(id); running {
				prop := w.Properties[id]
				asks("hire", "Take somebody on", 45, trade.Wage*7, w.HireReadiness(id),
					fmt.Sprintf("%d of %d positions filled. A week's wages up front at $%d a day after. Short-handed, it earns less and attracts trouble.", prop.Staff, trade.Hands, trade.Wage))
				add("layoff", "Let somebody go", 30, 0, w.LayOffReadiness(id),
					fmt.Sprintf("Cuts $%d a day from the wage bill and what the place can handle.", trade.Wage))
				if !prop.Order {
					add("order", "Take on a standing order", 60, 0, w.OrderReadiness(id),
						fmt.Sprintf("$%d a day from somebody respectable, for as long as %s keeps working at %d%%. It needs %d%% trade before anybody offers one, and losing it costs %d trade on top of the money.", OrderBonus, l.Name, int(OrderCapacity*100), OrderCustom, OrderLoss))
				}
				asks("restock", "Buy "+trade.Supplies, 45, trade.Restock, w.RestockReadiness(id),
					fmt.Sprintf("$%d. Currently %d left; a business out of %s barely trades.", trade.Restock, prop.Supply, trade.Supplies))
				if prop.Trouble {
					asks("remedy", trade.Remedy, 60, trade.RemedyCost, w.RemedyReadiness(id),
						fmt.Sprintf("$%d. %s %s", trade.RemedyCost, trade.Trouble, trade.RemedyDetail))
				}
				if ArmourySite(id) {
					if prop.Armoury {
						add("stock_arms", "Put the crates under the floor", 60, 0, w.StockReadiness(),
							fmt.Sprintf("Moves what you are carrying into the room. %d of %d crates down there, and an organization at war pays $%d apiece for them.", prop.Crates, ArmouryHold, w.ArmsPrice()))
					} else {
						asks("armoury", "Build a room under the floor", ArmouryMinutes, ArmouryCost, w.ArmouryReadiness(id),
							fmt.Sprintf("$%d. Holds %d crates of arms and draws %d attention a day plus one for every %d in it. Organizations at war buy at %d%% of the waterfront price and get stronger for it. A search that finds it takes everything and the premises with it.", ArmouryCost, ArmouryHold, ArmouryHeat, ArmouryCrateHeat, WarPremium))
					}
				}
				if StillSite(id) {
					if prop.Still {
						add("dismantle", "Take the still out", 90, 0, w.DismantleReadiness(id),
							fmt.Sprintf("Ends production of about %d crates a day and the attention that comes with it.", w.StillOutput(id)))
					} else {
						asks("still", "Set up a still in the back", StillMinutes, StillCost, w.StillReadiness(id),
							fmt.Sprintf("$%d. Produces moonshine you can sell, draws %d attention a day on top of what the stock draws, and a search that finds it costs far more than one that does not.", StillCost, StillHeat))
					}
				}
			}
			if HasBankroll(id) {
				prop := w.Properties[id]
				asks("bankroll", "Put money behind the tables", 45, 0, w.BankrollReadiness(id, 0),
					fmt.Sprintf("Into the float, currently $%d. %s The house keeps roughly %d%% of what crosses the tables over a season and loses on plenty of single nights. A house that cannot pay a winner is finished as a room worth playing in.", prop.Bankroll, coverage(w.NightHandleAt(id)), HouseEdge))
				sum(BankrollLeast, p.Cash, BankrollLot, "Behind the tables")
				add("draw", "Take money off the tables", 45, 0, w.DrawReadiness(id, 0),
					fmt.Sprintf("Out of the $%d float and into your hands. It is the only way this room's winnings reach you, and everything taken is action it can no longer attract.", prop.Bankroll))
				sum(BankrollLeast, prop.Bankroll, min(BankrollLot, prop.Bankroll), "Taken out")
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
			if w.Properties[id].Income > 0 {
				if posted := w.PostedAt(id); posted != nil {
					add("unpost", "Take "+posted.Name+" off the door", PostingMinutes, 0, "",
						fmt.Sprintf("Worth %d against anybody coming for %s while they are on it.", w.PostingDefenceAt(id), l.Name))
				} else {
					// What it costs is part of what it is: whoever goes has to
					// walk there, and the door is worth nothing until they
					// arrive. The button used to say "standing here" and quote
					// a figure that would not be true for another half hour.
					detail := fmt.Sprintf("Your most reliable person. Worth about %d against anybody coming for it once they are standing in it, and turns away most of what the street tries. They are also the one standing in it when somebody does come.", PostingDefence+15)
					if free := w.Unposted(); len(free) > 0 {
						best := free[0]
						for _, n := range free {
							if n.Trust > best.Trust {
								best = n
							}
						}
						if best.Location != id {
							detail = fmt.Sprintf("%s is at %s, %d minutes away. The door is worth nothing until they get there. %s",
								best.Name, placeName(best.Location), TravelMinutes(best.Location, id), detail)
						}
					}
					add("post", "Put somebody on the door", PostingMinutes, 0, w.PostReadiness(id), detail)
				}
			}
			add("repair", "Repair the property", 60, 50, need(w.Properties[id].Condition >= 100, "Already in good condition"), "Restore 40 condition.")
		} else {
			label := "Establish protection"
			if id == "casino" {
				label = "Reopen the casino"
			}
			cost := AcquisitionCost(w, id)
			if strings.HasPrefix(w.Properties[id].Owner, "former:") {
				label = "Buy out the former organization"
			}
			add("acquire", label, 60, cost, w.AcquireReadiness(id), fmt.Sprintf("Earn up to $%d/hour. Income accrues automatically; rivals may take notice.", w.Properties[id].Income))
		}
	}
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.ID == w.PlayerOrganizationID() || w.homeOf(f.ID) != id {
			continue
		}
		if w.Player.Serves == f.ID {
			add("leave_service", "Stop answering to "+f.Name, ServiceMinutes, 0, "",
				fmt.Sprintf("You are a %s of theirs on $%d a day. Walking out costs %d standing with them and they will answer it.", lowerFirst(w.ServiceTitle()), w.ServicePay(), LeavingCost))
			continue
		}
		add("serve:"+f.ID, "Go to work for "+f.Name, ServiceMinutes, 0, w.ServeReadiness(f.ID),
			fmt.Sprintf("$%d a day and a place to stand, rising to a share of what they take. Their quarrels become yours, and nobody comes up through somebody else's organization while they have one of their own.", SoldierWage))
	}
	if w.Player.Serves != "" {
		if leader := w.Leader(w.Player.Serves); leader != nil && leader.Location == id {
			f := w.faction(w.Player.Serves)
			add("takeover", "Move on "+leader.Name, TakeoverMinutes, 0, w.TakeoverReadiness(),
				fmt.Sprintf("Everything %s has becomes yours: the premises, the people who stay, and every quarrel the name was in. They are not an easy person to be in a room with, and if they are expecting it you are not one of theirs any more, if you are anything.", f.Name))
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
			leaving := ""
			if fitted := w.Comforts(p.Home); len(fitted) > 0 {
				labels := []string{}
				for _, f := range fitted {
					c, _ := ComfortByID(f)
					labels = append(labels, lowerFirst(c.Label))
				}
				leaving = " You would be leaving " + joinNames(labels) + " behind."
			}
			add("move_home", label, 60, cost, reason, fmt.Sprintf("$%d/day upkeep. Moving resets hired security. Respect is earned only for a new housing tier.%s", HomeRent(id), leaving))
		} else {
			add("rest", "Rest for four hours", 240, 0, "", "Recover up to 25 health as you rest. Rivals can act while you sleep.")
			add("security", "Hire another security detail", 30, 100, need(p.Security >= 3, "Maximum security hired"), "Improves detection and survival at home. Adds $10/day upkeep.")
			for _, c := range comforts {
				add("fit:"+c.ID, "Fit "+lowerFirst(c.Label), c.Minutes, 0, w.FitReadiness(id, c.ID),
					fmt.Sprintf("$%d, then $%d a day. %s It stays with the building if you move out.", c.Cost, c.Upkeep, c.Detail))
			}
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
	for _, n := range w.OwnPeople() {
		if n.Location != id {
			continue
		}
		asks("share:"+n.ID, "Pay "+n.Name+" a share", 15, ShareCost, w.PayShareReadiness(n.ID),
			fmt.Sprintf("They think of you at %d of 100. Below %d they start looking for somewhere else to be, and somebody ambitious who leaves takes a business with them.", n.Trust, DefectionTrust))
		add("dismiss:"+n.ID, "Put "+n.Name+" out", 15, 0, w.LetGoReadiness(n.ID),
			fmt.Sprintf("Ends the $%d a day. Nobody takes that as well as they pretend to.", MemberWage))
	}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != id || n.Faction != "" {
			continue
		}
		if IsOfficial(n.ID) || w.isRoleHolder(n) {
			continue
		}
		asks("sign:"+n.ID, "Put "+n.Name+" on", 45, SigningCost, w.SignOnReadiness(n.ID),
			fmt.Sprintf("$%d up front and $%d a day. Adds to what your organization is worth in a fight, stands in front of what comes at you, and can decide one morning that it is not worth it.", SigningCost, MemberWage))
		break
	}
	// Money out with a name on it, and what to do about it when the name stops
	// being able to pay. Anybody who already owes you is always listed; new
	// lending is capped, because a room of fifteen strangers rendered as
	// fifteen identical buttons is not a choice, it is a wall.
	offers := 0
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location != id || IsOfficial(n.ID) {
			continue
		}
		if l := w.LoanTo(n.ID); l != nil {
			days := (l.Due - w.Minute + 1439) / 1440
			state := fmt.Sprintf("$%d due in %s.", l.Owed, counted(days, "day", "days"))
			if l.Missed > 0 {
				state = fmt.Sprintf("$%d overdue, missed %s.", l.Owed, counted(l.Missed, "time", "times"))
			}
			add("lean:"+n.ID, "Collect from "+n.Name, LeanMinutes, 0, w.LeanReadiness(n.ID),
				fmt.Sprintf("%s About %d%% of getting most of it. Costs %d attention, %d respect gained, and somebody who remembers it for as long as they live.",
					state, int(w.LeanOdds(n)*100), LeanHeat, LeanRespect))
			add("extend:"+n.ID, "Give "+n.Name+" another week", 30, 0, w.ExtendReadiness(n.ID),
				fmt.Sprintf("%s Becomes $%d, due in %s. They are grateful today and further from paying it than they were this morning.",
					state, l.Owed+int(float64(l.Owed)*ExtendRate), counted(LoanTermDays, "day", "days")))
			add("forgive:"+n.ID, "Write off what "+n.Name+" owes", 15, 0, w.ForgiveReadiness(n.ID),
				fmt.Sprintf("%s Costs the money and %d respect. Buys the one thing money cannot: somebody who knows exactly what that was worth.", state, ForgiveRespect))
			continue
		}
		if reason := w.LendReadiness(n.ID); (reason == "" || w.Known(n)) && offers < LendOffers {
			offers++
			size := w.LoanSize(n)
			// LoanSize is nothing when the book is full or the money is not
			// there, and the refusal beside this already says which. Terms of
			// "$0 out, $0 back" are not terms.
			terms := fmt.Sprintf("Money out at %d%% back inside %s, when there is room on your book for it.",
				int(LoanRate*100), counted(LoanTermDays, "day", "days"))
			if size > 0 {
				terms = fmt.Sprintf("$%d out, $%d back inside %s.",
					size, size+int(float64(size)*LoanRate), counted(LoanTermDays, "day", "days"))
			}
			add("lend:"+n.ID, "Lend "+n.Name+" money", LendMinutes, 0, reason,
				terms+" If they cannot pay, what you do about it is the decision, and the street will hear which way you went.")
		}
	}
	if len(p.Crew) > 0 {
		// Whoever is actually on the books. Three buttons in this same block
		// named them by asking — rob:crew, mug:crew and sabotage:crew all use
		// CrewHands, and the collections description below already reads
		// n.Name — while these four had "Leo" written into the string. Once
		// hiring started naming whoever drives, one room called one crew Bela
		// Havel three times and Leo three times.
		hand := p.Crew[0].Name
		reason := need(p.Crew[0].Loyalty < 30, hand+" refuses assignments below 30 loyalty. Pay a bonus to rebuild trust.")
		if len(w.Tasks) > 0 {
			reason = hand + " is already on assignment"
		}
		round := w.CollectionRound()
		roundPlace, _ := PlaceByID(round)
		collecting := fmt.Sprintf("Two hours of doors at %s: $%d. Requires 30 loyalty.", roundPlace.Name, CollectionPay)
		if n := w.NPC(p.Crew[0].ID); n != nil && n.Location != round {
			collecting = fmt.Sprintf("%s walks to %s — %d minutes — and is not here while they are doing it. Two hours of doors: $%d. Requires 30 loyalty.",
				n.Name, roundPlace.Name, TravelMinutes(n.Location, round), CollectionPay)
		}
		add("delegate", "Send "+hand+" on collections", 15, 0, reason, collecting)
		about(p.Crew[0].ID)
		add("crew_bonus", "Pay "+hand+" a bonus", 15, 40, need(p.Crew[0].Loyalty >= 100, "Loyalty is already at its maximum"), "Restore up to 25 loyalty. Below 30 they refuse collections; at 50 they can help protect businesses when available.")
		about(p.Crew[0].ID)
	}
	if offer, ok := w.AvailableCommission(id); ok {
		add("commission", "Hear what "+offer.GiverName+" wants", 30, 0, w.CommissionReadiness(id),
			fmt.Sprintf("%s $%d, %d respect and %d standing with %s. Three days. Failing costs %d standing with them.", offer.Brief, offer.Pay, offer.Respect, offer.Goodwill, w.factionName(offer.PatronID), offer.Penalty))
	}
	add("wait", "Let an hour pass", 60, 0, "", "Income, rent, operations and rival plans continue.")
	// One rule, applied once, over everything aimed at a person: somebody out
	// on the street is not here to be dealt with. Every action about a person
	// is written at the place it belongs to and none of them checked whether
	// the person was still standing in it, so the player could buy a coffee for
	// a woman who was halfway across the city. Doing it here rather than at
	// forty call sites means the next action about a person cannot forget.
	for i := range out {
		if out[i].Subject == "" || out[i].Disabled || reachesTheUnreachable(out[i].ID) {
			continue
		}
		if reason := w.OutOfReach(out[i].Subject); reason != "" {
			out[i].Disabled, out[i].Reason = true, reason
		}
	}
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
	w.Dead = append(w.Dead, Death{Name: p.Name, Minute: w.Minute, Life: w.Life, Cause: cause})
	// What they built outlives them if anybody was left to hold it. Whoever
	// takes it over keeps the premises, so this runs before the estate claims
	// anything.
	inherited := ""
	if w.Incorporated() {
		inherited = w.Inherit()
	}
	w.orphan(w.PlayerOrganizationID())
	for id, prop := range w.Properties {
		if w.Own(id) {
			prop.Owner = "former:" + p.Name
			if prop.Income > 0 {
				prop.Income = max(5, prop.Income-2)
			}
		}
	}
	if inherited != "" {
		w.Dead[len(w.Dead)-1].Estate = w.EstateName(inherited)
	}
	w.Tasks = []Task{}
	w.Event = nil
	w.Log(p.Name+" is dead", cause+" Your life ends here. The city continues.", "death")
}
func (w *World) Attack(plot Plot) {
	// Two places in the world they cannot walk into: a police cell, and
	// anywhere that is not this city. Both leave them the house.
	if w.Held() || w.Abroad != "" {
		w.wreckTheHouse()
		return
	}
	// A warning is worth exactly this. Somebody told you they were coming and
	// told you to stay away from home, so that is where they went. Being
	// somewhere else is the whole value of having been warned; it costs the
	// house, which is what makes it a decision rather than a free pass.
	if plot.Known && w.Player.Location != w.Player.Home {
		w.wreckTheHouse()
		return
	}
	if w.Warned(plot) {
		w.ambushScene(plot)
		return
	}
	if w.Random() < w.AmbushOddsHere() {
		w.Die(fmt.Sprintf("An attack at %s caught you without warning or protection.", w.whereItHappens()))
		return
	}
	w.survived()
}
func (w *World) Advance(minutes int) {
	p := &w.Player
	// A seat is taken in one room. Whatever moved the player out of it — a
	// journey, a scene, an errand somebody else ran for them — they are not at
	// those tables any more.
	w.LeaveTable()
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
				if plot.Kind == "hit" && !plot.Known && w.Reach() >= 2 {
					next = min(next, max(w.Minute+1, plot.Due-90))
				}
			}
		}
		elapsed := next - w.Minute
		w.Minute = next
		for id, prop := range w.Properties {
			if w.Own(id) {
				prop.Carry += float64(prop.Income*prop.Condition*elapsed) * operatingMode(prop.Mode).Take * w.Capacity(id) * w.TradeMultiplier(id) * (1 + w.LicenceTake()) * w.CollectionShare(id) / 6000
				n := int(prop.Carry + 1e-9)
				prop.Carry -= float64(n)
				w.Earn(n)
			}
		}
		w.settleTasks()
		// Organizations reconsider each other twice a day; their books settle once.
		if w.Minute%720 == 0 {
			w.FactionTurn()
			w.MarketPrices()
			w.ConsiderRobbery()
			w.SetOut()
		}
		w.Arrivals()
		w.ResolveContracts()
		w.SettleCommissions()
		if w.Minute%1440 == 0 {
			w.FamilyDay()
			w.BusinessDay()
			w.ContrabandDay()
			w.PoliceDay()
			w.CustodyDay()
			w.LoanDay()
			w.PeopleDay()
			w.GrudgeDay()
			w.SettleGrudges()
			w.OperationsDay()
			w.StillDay()
			w.CasinoDay()
			w.TableNight()
			w.CarDay()
			w.CarTrade()
			w.RepairsDay()
			w.DryDay()
			w.PumpDay()
			w.ChargeDay()
			w.DemolitionDay()
			w.CityHallDay()
			w.ArmouryDay()
			w.RecruitDay()
			w.FillRoles()
			w.CustomDay()
			w.OrderDay()
			w.OrganizationDay()
			w.OwnPeopleDay()
			w.PactDay()
			w.ServiceDay()
			w.FortunesDay()
			w.ObituaryDay()
			w.CityPageDay()
			w.ScrutinyDay()
			w.CivicDay()
			w.PrunePeople()
			w.DressDay()
			bill := w.DailyCost()
			if p.Cash >= bill {
				p.Cash -= bill
				w.Log("Accounts settled", fmt.Sprintf("$%d paid for housing, security and crew.", bill), "business")
			} else {
				// Report the night that happened, not the first one. This said
				// "Security leaves; your residence is now a rented room.
				// Unpaid crew lose loyalty" every midnight the player came up
				// short, so the second night and every night after claimed
				// three losses that had already been taken.
				lost := []string{}
				if p.Security > 0 {
					lost = append(lost, "Security leaves")
					p.Security = 0
				}
				if p.Home != "room" {
					lost = append(lost, "your residence is now a rented room")
					p.Home = "room"
				}
				if len(p.Crew) > 0 && p.Crew[0].Loyalty > 0 {
					lost = append(lost, p.Crew[0].Name+" is not being paid and knows it")
					p.Crew[0].Loyalty = max(0, p.Crew[0].Loyalty-20)
				}
				had := p.Cash
				p.Cash = max(0, p.Cash-15)
				text := "You could not cover the bills."
				if len(lost) > 0 {
					text += " " + upper1(strings.Join(lost, "; ")) + "."
				} else if had > 0 {
					text += fmt.Sprintf(" There was nothing left to take but $%d of what you were carrying.", had-p.Cash)
				} else {
					text += " There is nothing left to take, which is its own kind of trouble."
				}
				w.Log("Your arrangements unravel", text, "danger")
			}
		}
		for j := 0; j < len(w.Plots); j++ {
			plot := &w.Plots[j]
			if plot.Life != w.Life {
				continue
			}
			if plot.Kind == "hit" && !plot.Known && w.Reach() >= 2 && plot.Due-w.Minute <= 90 {
				plot.Known = true
				actor := w.factionName(plot.Actor)
				warning := fmt.Sprintf("%s %s people asking where you sleep.", Leads(actor), Agree(actor, "has", "have"))
				w.Log("Mara has heard something", warning+" You may have very little time.", "danger")
				if plot.Due > w.Minute {
					w.Event = &Scene{ID: "warning-" + plot.ID, Title: "A call worth answering", Body: warning + " I cannot tell you exactly when they will come. Stop what you are doing and think about where you want to be tonight.", Speaker: w.HolderID("fixer"), Kind: "warning", Source: "authored", Minute: w.Minute, Choices: []Choice{{ID: "acknowledge", Label: "Put down the phone and prepare", Detail: "Clock stays paused. You can leave, arrange security or seek an audience. The threat remains."}}}
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
		locs = append(locs, map[string]any{"id": l.ID, "name": l.Name, "type": l.Type, "district": l.District, "x": l.X, "y": l.Y, "cost": l.Cost, "blurb": l.Blurb, "owner": prop.Owner, "holder": w.HolderName(l.ID), "staff": prop.Staff, "supply": prop.Supply, "trouble": prop.Trouble, "shy": w.Shy(l.ID), "curtains": w.Curtains(l.ID), "trade": w.CustomDescription(l.ID), "posted": w.PostingDescription(l.ID), "people": w.PeopleHere(l.ID), "note": w.PlaceNote(l.ID), "room": w.RoomNote(l.ID), "note_warn": w.PlaceWarn(l.ID), "away": w.Away(l.ID), "travel_note": w.TravelNote(l.ID), "still": prop.Still, "bankroll": prop.Bankroll, "handle": w.NightHandleAt(l.ID), "capacity": w.Capacity(l.ID), "trading": w.Trading(l.ID), "condition": prop.Condition, "income": prop.Income, "owned": w.Own(l.ID), "locked": l.District > w.District, "actions": w.Actions(l.ID)})
	}
	var scene any = nil
	if e := w.Event; e != nil {
		choices := append([]Choice{}, e.Choices...)
		for i := range choices {
			choices[i].Disabled = w.Player.Cash < choices[i].Cost
			choices[i].Reason = ""
			if choices[i].Disabled {
				choices[i].Reason = "Not enough cash: this takes " + cash(choices[i].Cost) + " and you are holding " + cash(w.Player.Cash)
			}
		}
		scene = map[string]any{"id": e.ID, "title": e.Title, "body": e.Body, "speaker": e.Speaker, "kind": e.Kind, "source": e.Source, "minute": e.Minute, "choices": choices, "connection": e.Connection, "conditions": e.Conditions}
	}
	history := w.History
	if len(history) > 60 {
		history = history[len(history)-60:]
	}
	return map[string]any{"id": w.ID, "version": w.Version, "revision": w.Revision, "life": w.Life, "minute": w.Minute, "sky": w.Sky(), "player": w.Player, "district": w.District, "factions": w.PublicFactions(), "npcs": w.People(), "locations": locs, "event": scene, "history": history, "dead": w.Dead, "tasks": w.Tasks, "director": w.Director, "last_result": w.LastResult, "daily_cost": w.DailyCost(), "books": w.Books(), "guide": w.Guide(), "rules": GuideRules(), "groups": Groups(), "income": income, "security": w.Guard(), "opportunity": w.NextOpportunity(), "known_threats": w.KnownThreats(), "business_truces": w.ActiveBusinessTruces(), "conflicts": w.PublicConflicts(), "goods": w.Goods, "arms": w.ArmsDescription(), "appearance": w.AppearanceDescription(), "vehicle": w.VehicleDescription(), "residence": w.ResidenceDescription(), "offshore": map[string]any{"balance": w.Offshore, "reachable": w.Player.Offshore}, "newspaper": w.Edition(), "editions": w.Editions(), "arrangements": w.PendingArrangements(), "commissions": w.PublicCommissions(), "grudges": w.GrudgeSummary(), "cast": w.Cast(), "everyone": w.Everyone(), "retainers": w.RetainerDescription(), "armoury": w.ArmouryDescription(), "population": w.PopulationSummary(), "seated": w.Seated, "hand": w.HandDescription(), "wheel": w.WheelDescription(), "machine": w.MachineDescription(), "house": w.HouseDescription(), "roles": w.RoleDescription(), "organization": w.PlayerOrganizationDescription(), "own_people": w.OwnPeopleDescription(), "pacts": w.PactDescription(), "book": w.LoanDescription(), "press": w.PressDescription(), "service": w.ServiceDescription(), "city": w.ScrutinyDescription(), "dashboard": w.Dashboard(), "epitaph": w.Epitaph(), "street": w.OnTheStreet(), "street_note": w.StreetNote()}
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
	// Nobody brings work to a man in a cell.
	if !w.Player.Alive || w.Held() || w.Event != nil || w.SuspendedJob != nil || len(w.KnownThreats()) > 0 {
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
		e, _ := w.ValidateProposal(Proposal{"", "A favor with a price", "“A merchant wants a sealed ledger moved before their partners arrive. I would understand if you preferred the ordinary work.”", w.HolderID("fixer"), "courier", "You moved the ledger. Mara now knows you can handle sensitive work.", "", []Approach{{Method: "careful", Label: "Wait for a quiet route"}, {Method: "press", Label: "Move it before the partners arrive"}}})
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
	scene := &Scene{Target: p.Location, Operation: p.Operation, ID: ID(), Title: p.Title, Beneficiary: p.Beneficiary, Body: p.Body, Speaker: p.Speaker, Kind: "proposal", Source: "local-ai", Minute: w.Minute, Effect: fx, Outcome: operationOutcome(p.Operation), Choices: []Choice{Choice{ID: "accept", Label: operationLabel(p.Operation)}.terms(fx), {ID: "decline", Label: "Decline the arrangement", Detail: "No cost or time."}}}
	scene.Conditions = "At 15 heat, police may stop completion." + politicalDetail
	if err := addApproaches(scene, p.Approaches); err != nil {
		return nil, err
	}
	if p.Location != "" {
		place, _ := PlaceByID(p.Location)
		scene.Conditions = place.Name + " · " + scene.Conditions
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
		"consignment": "Hand the crates over", "grievance": "Put yourself between them",
		"obligation": "Get it finished", "warning_off": "Have the conversation",
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
		"consignment":  "You handed the crates over away from your own premises and took the money for them.",
		"grievance":    "You stood between the two of them long enough for it to stop being about tonight.",
		"obligation":   "You finished what you had promised somebody and were seen to finish it.",
		"warning_off":  "You had the conversation, and they know you know.",
	}
	if outcome, ok := outcomes[operation]; ok {
		return outcome
	}
	return "You completed the arrangement."
}

// reachesTheUnreachable is the one exception to the rule above: an action whose
// whole premise is that the person cannot be dealt with in the ordinary way.
//
// Bailing somebody out is not a conversation with them. It is a conversation
// with a clerk about them, and it exists precisely because they are in a cell.
// The sweep disabled it for the reason it is offered — "Otto Reiss is being
// held at Ward Street Station" — so the button could never be pressed, and the
// only way to reach the mechanic was to post the command directly.
func reachesTheUnreachable(id string) bool {
	return strings.HasPrefix(id, "bail:")
}

// AcquisitionCost is what taking a premises costs. Buying out what is left of a
// dead organization costs twice what an unclaimed door does.
func AcquisitionCost(w *World, id string) int {
	place, ok := PlaceByID(id)
	if !ok {
		return 0
	}
	if prop := w.Properties[id]; prop != nil && strings.HasPrefix(prop.Owner, "former:") {
		return place.Cost * 2
	}
	return place.Cost
}

// AcquireReadiness explains why a premises cannot be taken, or returns "".
//
// This used to have no function of its own: the button computed the standing
// requirement inline and let the generic cost check add "Not enough cash", and
// the guide page — which promises it "asks the game the same question the
// buttons ask" — kept its own copy that knew about respect and not about money.
// So a player with $51 and 33 respect was told premises were available now
// while every door in the city said "Not enough cash".
func (w *World) AcquireReadiness(id string) string {
	// Only somewhere that earns is for sale. A rented room, the station and the
	// newspaper are not premises anybody takes over, and the button is not
	// offered there — so neither is this, or the guide would count a door
	// nobody can walk through as an opportunity.
	if prop := w.Properties[id]; prop == nil || prop.Income <= 0 {
		return "There is nothing here to take over"
	}
	// And a place with no price is not for sale. The bar, the docks, the
	// exchange and the club all earn and all carry a cost of nothing, because
	// none of them was ever meant to change hands — opening the business block
	// to everywhere that earns offered them for free.
	if place, ok := PlaceByID(id); !ok || place.Cost <= 0 {
		return "This is not somewhere that changes hands"
	}
	if w.Own(id) {
		return "This is already yours"
	}
	if !w.CanAcquire(id) {
		return "This property belongs to another organization"
	}
	// What it takes to be given the keys, by what the place is rather than by
	// its name: a room that runs games wants somebody the city has heard of, a
	// yard with machinery in it wants more than nothing, and the rest want you
	// to have done something.
	req := 6
	if place, ok := PlaceByID(id); ok {
		switch {
		case place.Type == "casino":
			req = 20
		case place.Kind == "garage" || place.Kind == "scrapyard" || place.Kind == "dealer":
			req = 10
		}
	}
	if w.Player.Respect < req {
		return fmt.Sprintf("Earn %d respect first", req)
	}
	if cost := AcquisitionCost(w, id); w.Player.Cash < cost {
		return "Not enough cash"
	}
	return ""
}

// PlaceIncome is what each address earns its owner by the hour. It is a table
// rather than a switch inside world creation because a save loaded from an
// earlier build has to be given the same figure, and a business added to the
// city after a campaign began was otherwise worth nothing in it forever.
//
// The casino's is the floor take only: the bar, the door and the rooms
// upstairs. What the tables make is decided every night by the float behind
// them, in CasinoDay, rather than accruing by the hour. Haulage is the most an
// ordinary business earns here and the most it costs to run — six men on the
// books before it has moved anything.
var PlaceIncome = map[string]int{
	"laundry": 14, "garage": 24, "casino": 18, "club": 30,
	"market": 18, "docks": 22, "bar": 12,
	"restaurant": 20, "poolhall": 16, "butcher": 28, "haulage": 40,
	// A second house and a second laundry: the same kind of business, their own
	// takings, because what an address earns belongs to the address and not to
	// the kind. The Golden Lily is quieter than the Blue Hour and the Ordway is
	// bigger than the Bluebird.
	"goldenlily": 15, "steamworks": 19, "burlesque": 26, "cabstand": 32,
	// A forecourt earns steadily on servicing and part-exchange; the money in
	// it is the cars, and that is settled when one is sold.
	"dealer": 22, "archway": 21,
	"filling": 19, "pumps": 16,
	// A yard earns on what everybody else is finished with.
	"scrapyard": 18,
}
