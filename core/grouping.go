package core

import "strings"

// The sidebar renders whatever the core offers, in the order the core happens
// to build it. That was the right decision while there were eight actions and
// the wrong one at ninety: standing in a bar now produces an undifferentiated
// column of cards where taking a job, robbing the till, lending a stranger
// money and paying your own man a share all look identical and sit in whatever
// order the switch statement was written in.
//
// Grouping belongs here rather than in the interface. The core knows what an
// action is for; the interface only knows its id. Putting the answer in one
// function means a new action is grouped the day it is written, and the test
// below fails if it is not.

// Group is a family of actions, in the order they should be offered. The order
// is a judgement about what a player is usually looking for: what earns, then
// what they own, then the people, then the risky work, then the slow work of
// becoming somebody, then money, then leaving.
type Group struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Blurb string `json:"blurb"`
}

var groups = []Group{
	{"work", "Work", "Jobs that pay today."},
	{"business", "Your premises", "Keeping what you own earning."},
	{"people", "People", "Who works for you, who owes you, who you know."},
	{"street", "The street", "Work that can go wrong, and hurt somebody."},
	{"standing", "Standing", "Who you are to this city, and who owes you a favour."},
	{"tables", "The tables", "Cards, the wheel, the dice and the machines."},
	{"money", "Money", "Moving it, hiding it, or putting it somewhere else."},
	{"travel", "Elsewhere", "Leaving where you are standing."},
}

// Groups is the ordered list, for anything that renders them.
func Groups() []Group { return groups }

// actionGroup is the exception table: everything not covered by the prefix
// rules below. Kept as one map so a new action is one line rather than a
// decision spread across the file that offers it.
var actionGroup = map[string]string{
	"form_family": "business", "set_headquarters": "business",
	// Work that pays on the day.
	"courier": "work", "dockwork": "work", "rushorder": "work", "householdwork": "work", "delegate": "work", "commission": "work",
	"contract": "work", "order": "work",

	// Premises.
	"sell_property": "business", "buy_residence": "business", "acquire": "business", "repair": "business", "hire": "business", "layoff": "business", "wage": "business",
	"restock": "business", "remedy": "business", "inspect": "business", "still": "business",
	"dismantle": "business", "armoury": "business", "stock_arms": "business",
	"bankroll": "business", "post": "business", "unpost": "business",
	// Taking money back off your own tables is running the place, not a night
	// out at it. This was filed under work that can go wrong.
	"landing": "business", "draw": "business", "limit": "business", "night": "business",

	// People.
	"recruit": "people", "crew_bonus": "people", "contact": "people",

	// Work that can go wrong.
	"rob": "street", "mug": "street", "sabotage": "street", "move": "street", "frighten": "street",
	"demand": "street", "push": "street", "incite": "street", "provoke": "street", "takeover": "street", "charge": "street",
	"plant": "street", "incendiary": "street", "driveby-building": "street", "strip": "street",

	// Becoming somebody.
	"expand": "standing", "audience": "standing", "sitdown": "standing", "bribe": "standing",
	"investigate": "standing", "lie_low": "standing", "press": "standing",
	"security": "standing", "move_home": "standing", "car": "standing", "leave_service": "standing",
	// Everything about the car you own, in one place. Fuelling it and plating
	// it were nowhere, and having it worked on was filed as a job that pays.
	"fill": "standing", "plate": "standing", "service": "standing",
	"spike": "standing", "puff": "standing", "lawyer": "standing", "talk": "standing",
	"sit_out": "standing",

	// The tables. Sitting down, getting up, and every verb of every game.
	"scrap": "business",
	"sit":   "tables", "rise": "tables", "play": "tables", "pull": "tables",
	"deal": "tables", "cashout": "tables",
	"wheel": "tables", "dice": "tables", "roll": "tables", "hit": "tables",
	"stand": "tables", "cards": "tables", "bet": "tables", "call": "tables", "fold": "tables",

	// Money.
	"launder": "money", "deposit": "money", "withdraw": "money", "offshore_access": "money",

	// Leaving.
	"travel": "travel", "rest": "travel", "wait": "travel",
}

// prefixGroup covers the actions that carry an id after a colon. The prefix is
// the verb and the suffix is who or what it is done to, so the verb decides.
var prefixGroup = [][2]string{
	{"home_strike:", "street"},
	{"burgle:", "street"},
	{"buy_apartment:", "business"}, {"sell_apartment:", "business"},
	{"sign:", "people"}, {"share:", "people"}, {"dismiss:", "people"}, {"poach:", "business"}, {"ask:", "people"}, {"about:", "people"}, {"incharge:", "business"},
	{"lend:", "people"}, {"lean:", "people"}, {"extend:", "people"},
	{"word:", "standing"}, {"forgive:", "people"}, {"bail:", "people"}, {"break:", "people"}, {"funeral:", "people"},
	{"rob:", "street"}, {"mug:", "street"}, {"sabotage:", "street"},
	{"arms:", "street"},
	// Going after somebody, or sending one of your own to. The whole of the
	// violence in this game sat under "jobs that pay today".
	{"strike:", "street"}, {"send:", "street"},
	// A car for one of your own, and plate on it: about them, not about you.
	{"lot:", "standing"}, {"attire:", "standing"}, {"car:", "people"}, {"plate:", "people"}, {"give:", "people"},
	// A counter is where money comes from when there is none.
	{"pawn:", "money"}, {"redeem:", "money"}, {"window:", "money"},
	{"buy:", "money"}, {"sell:", "money"}, {"play:", "tables"},
	// Which of the two things in the room you are sitting down to.
	{"sit:", "tables"},
	{"retain:", "standing"}, {"release:", "standing"}, {"smear:", "standing"},
	{"enquire:", "standing"}, {"pact:", "standing"}, {"serve:", "standing"},
	{"operate:", "business"}, {"fit:", "business"},
	{"trip:", "travel"},
}

// Classified reports whether an action's group was chosen for it rather than
// arrived at by the fallback. The fallback exists so that a new action still
// reaches the player; it is not a place for an action to live.
func Classified(id string) bool {
	if _, ok := actionGroup[id]; ok {
		return true
	}
	for _, rule := range prefixGroup {
		if strings.HasPrefix(id, rule[0]) {
			return true
		}
	}
	return false
}

// GroupOf is what an action is for. Everything the core can offer has an
// answer; TestEveryActionBelongsSomewhere is what keeps that true.
func GroupOf(id string) string {
	if group, ok := actionGroup[id]; ok {
		return group
	}
	for _, rule := range prefixGroup {
		if strings.HasPrefix(id, rule[0]) {
			return rule[1]
		}
	}
	// An action nobody has classified is still an action. Offering it under
	// the work heading is better than not offering it at all.
	return "work"
}
