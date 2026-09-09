package core

import "testing"

// The People screen opens with "52 people live in this city and you know 33 of
// them. 29 answer to an organization and 17 to nobody." Twenty-nine and
// seventeen make forty-six, and six people are in neither figure: the four
// officials, the fixer, and the player's own crew before they are an
// organization. Somebody doing one of the city's jobs answers to the city, and
// the sentence had no room for them.
//
// It also disagreed with the screen underneath it, which files each person in
// exactly one place and whose chips add to the whole city: the header said
// twenty-nine answer to an organization while the filter beside it said
// "Organizations 23", because that filter does not hold the family heads it
// files under names everybody knows.

func TestThePopulationAddsUp(t *testing.T) {
	check := func(what string, w *World) {
		p := w.PopulationSummary()
		living, organized := p["living"].(int), p["organized"].(int)
		street, jobs := p["street"].(int), p["jobs"].(int)
		if organized+street+jobs != living {
			t.Errorf("%s: %d in an organization + %d on the street + %d holding a job = %d, and %d people live here",
				what, organized, street, jobs, organized+street+jobs, living)
		}
		if known := p["known"].(int); known > living {
			t.Errorf("%s: the player knows %d of %d people", what, known, living)
		}
	}
	w := New(41)
	check("a new city", w)

	// A city that has run, with families made and lost.
	aged := withFamily(41, "estate:kohl", "Vera Kohl's people", "Vera Kohl")
	live(t, aged, 4000)
	check("an aged city", aged)

	// And one where the player is an organization of their own, so their crew
	// answers to somebody.
	own := New(41)
	own.District = 2
	own.Player.Cash, own.Player.Respect = 20000, 200
	for _, id := range []string{"laundry", "garage", "casino"} {
		own.Properties[id].Owner = "player:1"
	}
	own.Incorporate()
	own.Player.Crew = []Crew{{"leo", "Leo Carver", 90}}
	check("a city where the player is somebody", own)
}
