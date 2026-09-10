package core

import "testing"

// Ninety actions in one undifferentiated column is not a menu, it is a wall.
// These guard the grouping that makes it navigable: everything the core can
// offer has a home, and every home the core names is one the interface knows.

func TestEveryActionBelongsSomewhere(t *testing.T) {
	known := map[string]bool{}
	for _, g := range Groups() {
		known[g.ID] = true
	}
	// A campaign rich enough to offer most of what exists: premises, people,
	// money, a car, an organization, a debt, and somebody to fight.
	w := proprietor(t)
	own(w, "laundry", "garage", "casino")
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Contacts = 90000, 200, 4
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
	w.ensureOfficials()
	w.OrganizationDay()
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) || w.isRoleHolder(n) || len(w.OwnPeople()) >= 2 {
			continue
		}
		w.Player.Location = n.Location
		w.SignOn(n.ID)
	}

	seen, ungrouped := 0, []string{}
	for _, l := range Locations {
		w.Player.Location = l.ID
		for _, a := range w.Actions(l.ID) {
			seen++
			if a.Group == "" || !known[a.Group] {
				ungrouped = append(ungrouped, a.ID+" -> "+a.Group)
			}
		}
	}
	if seen < 40 {
		t.Fatalf("only %d actions were offered anywhere in the city; this is not measuring what it claims to", seen)
	}
	if len(ungrouped) > 0 {
		t.Fatalf("actions with no home the interface knows: %v", ungrouped)
	}
	t.Logf("%d actions offered across the city, every one of them in one of %d groups", seen, len(known))
}

func TestAnActionNobodyClassifiedIsStillOffered(t *testing.T) {
	// A group table is a thing somebody forgets to update. Whatever happens, an
	// action must still reach the player rather than vanishing into a group the
	// interface does not render.
	group := GroupOf("something-nobody-has-written-yet")
	known := false
	for _, g := range Groups() {
		if g.ID == group {
			known = true
		}
	}
	if !known {
		t.Fatalf("an unclassified action landed in %q, which nothing renders", group)
	}
}

func TestThingsThatBelongTogetherAreTogether(t *testing.T) {
	// The grouping is a judgement, so these are the judgements, written down.
	for _, c := range []struct{ id, group string }{
		{"courier", "work"}, {"dockwork", "work"},
		{"repair", "business"}, {"restock", "business"}, {"post", "business"},
		{"operate:hard", "business"}, {"fit:safe", "business"},
		{"recruit", "people"}, {"sign:street-1", "people"}, {"lend:street-1", "people"},
		{"lean:street-1", "people"}, {"forgive:street-1", "people"}, {"bail:street-1", "people"},
		{"rob", "street"}, {"rob:crew", "street"}, {"mug", "street"}, {"sabotage:crew", "street"},
		{"takeover", "street"}, {"plant", "street"},
		{"bribe", "standing"}, {"retain:editor", "standing"}, {"smear:bellandi", "standing"},
		{"spike", "standing"}, {"puff", "standing"}, {"dress", "standing"},
		{"sit_out", "standing"}, {"lawyer", "standing"}, {"talk", "standing"},
		{"launder", "money"}, {"buy:moonshine", "money"}, {"deposit", "money"},
		{"travel", "travel"}, {"trip:rockridge", "travel"}, {"wait", "travel"},
	} {
		if got := GroupOf(c.id); got != c.group {
			t.Errorf("%s is offered under %q, expected %q", c.id, got, c.group)
		}
	}
}

// The guard above only asked whether an action's group is one the interface
// renders, and "work" is. So nine actions that nobody had classified — the
// house limit, plating a car, and every verb at the tables — sat quietly under
// "Jobs that pay today" and the test went on passing. This asks the question
// that was meant: is this action's group a decision somebody made, or the
// fallback?
func TestNoOfferedActionIsThereByDefault(t *testing.T) {
	w := proprietor(t)
	own(w, "laundry", "garage", "casino")
	w.District = 2
	w.Player.Cash, w.Player.Respect, w.Player.Contacts = 90000, 200, 4
	w.Player.Crew = []Crew{{ID: "leo", Name: "Leo Carver", Loyalty: 70}}
	w.Player.Car, w.Player.CarWear = 1, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	w.ensureOfficials()
	w.OrganizationDay()

	unclassified := map[string]bool{}
	for _, l := range Locations {
		w.Player.Location = l.ID
		for _, a := range w.Actions(l.ID) {
			if !Classified(a.ID) {
				unclassified[a.ID] = true
			}
		}
	}
	if len(unclassified) > 0 {
		names := []string{}
		for id := range unclassified {
			names = append(names, id)
		}
		t.Fatalf("actions filed under the fallback rather than anywhere chosen: %v", names)
	}
}
