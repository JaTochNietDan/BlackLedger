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
	ID, Title, Blurb string
}

var groups = []Group{
	{"work", "Work", "Jobs that pay today."},
	{"business", "Your premises", "Keeping what you own earning."},
	{"people", "People", "Who works for you, who owes you, who you know."},
	{"street", "The street", "Work that can go wrong, and hurt somebody."},
	{"standing", "Standing", "Who you are to this city, and who owes you a favour."},
	{"money", "Money", "Moving it, hiding it, or putting it somewhere else."},
	{"travel", "Elsewhere", "Leaving where you are standing."},
}

// Groups is the ordered list, for anything that renders them.
func Groups() []Group { return groups }

// actionGroup is the exception table: everything not covered by the prefix
// rules below. Kept as one map so a new action is one line rather than a
// decision spread across the file that offers it.
var actionGroup = map[string]string{
	// Work that pays on the day.
	"courier": "work", "dockwork": "work", "delegate": "work", "commission": "work",
	"contract": "work", "service": "work", "order": "work",

	// Premises.
	"acquire": "business", "repair": "business", "hire": "business", "layoff": "business",
	"restock": "business", "remedy": "business", "inspect": "business", "still": "business",
	"dismantle": "business", "armoury": "business", "stock_arms": "business",
	"bankroll": "business", "post": "business", "unpost": "business",

	// People.
	"recruit": "people", "crew_bonus": "people", "contact": "people",

	// Work that can go wrong.
	"rob": "street", "mug": "street", "sabotage": "street", "move": "street",
	"incite": "street", "provoke": "street", "takeover": "street", "charge": "street",
	"plant": "street", "hit": "street", "stand": "street", "draw": "street",

	// Becoming somebody.
	"expand": "standing", "audience": "standing", "sitdown": "standing", "bribe": "standing",
	"investigate": "standing", "lie_low": "standing", "dress": "standing", "press": "standing",
	"security": "standing", "move_home": "standing", "car": "standing", "leave_service": "standing",
	"spike": "standing", "puff": "standing", "lawyer": "standing", "talk": "standing",
	"sit_out": "standing",

	// Money.
	"launder": "money", "deposit": "money", "withdraw": "money", "offshore_access": "money",

	// Leaving.
	"travel": "travel", "rest": "travel", "wait": "travel",
}

// prefixGroup covers the actions that carry an id after a colon. The prefix is
// the verb and the suffix is who or what it is done to, so the verb decides.
var prefixGroup = [][2]string{
	{"sign:", "people"}, {"share:", "people"}, {"dismiss:", "people"},
	{"lend:", "people"}, {"lean:", "people"}, {"extend:", "people"},
	{"forgive:", "people"}, {"bail:", "people"}, {"break:", "people"},
	{"rob:", "street"}, {"mug:", "street"}, {"sabotage:", "street"},
	{"arms:", "street"},
	{"buy:", "money"}, {"sell:", "money"}, {"play:", "money"},
	{"retain:", "standing"}, {"release:", "standing"}, {"smear:", "standing"},
	{"enquire:", "standing"}, {"pact:", "standing"}, {"serve:", "standing"},
	{"operate:", "business"}, {"fit:", "business"},
	{"trip:", "travel"},
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
