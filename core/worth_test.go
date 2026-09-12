package core

import (
	"strings"
	"testing"
)

// "It would also be good if we generated nice car images to show what you're
// buying and stats information about the speed of the car relevant to what it
// gives to you."
//
// The forecourt said "journeys take 74% of the time they take on foot", which is
// true and tells nobody anything. What a car is worth is minutes on a road this
// player actually walks, and the honest number has to include what the car is
// carrying: plate is weight, and a plated car is slower than the table says.

func motorist(t *testing.T) *World {
	t.Helper()
	w := New(47)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 40000
	w.Player.Location = w.theForecourt()
	if w.Player.Location == "" {
		t.Skip("this city has no forecourt")
	}
	return w
}

func TestWhatACarIsWorthIsSaidInMinutesOnARealRoad(t *testing.T) {
	t.Parallel()
	w := motorist(t)
	worth := w.CarWorth(1)
	if worth.To == "" {
		t.Fatal("the lot named no road at all")
	}
	if worth.Walking <= 0 || worth.Driving <= 0 {
		t.Fatalf("a journey with no time in it: %+v", worth)
	}
	if worth.Driving >= worth.Walking {
		t.Fatalf("a car that is no quicker than walking: %d against %d", worth.Driving, worth.Walking)
	}
	// The road named is a real address, and one worth naming: the furthest
	// thing from here, because that is where a car is worth most.
	if _, ok := PlaceByID(worth.ToID); !ok {
		t.Fatalf("the lot named somewhere that is not in this city: %q", worth.ToID)
	}
	// And the button says it.
	a := actionByID(w.Actions(w.Player.Location), "lot:1")
	if a == nil {
		t.Fatal("nothing for sale on the forecourt")
	}
	if !strings.Contains(a.Detail, worth.To) {
		t.Fatalf("the description does not name the road it is quoting: %q", a.Detail)
	}
	if strings.Contains(a.Detail, "% of the time") {
		t.Fatalf("the description is still quoting a percentage: %q", a.Detail)
	}
	// And it says which car it is, so the panel can show that one rather than
	// working it out from the label.
	if a.Tier == 0 {
		t.Fatal("the card does not say which car is on the lot")
	}
}

func TestAPlatedCarIsQuotedAtWhatItActuallyDoes(t *testing.T) {
	t.Parallel()
	w := motorist(t)
	w.Player.Car, w.Player.CarWear = 1, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	bare := w.CarWorth(w.Player.Car)
	w.Player.Plate = PlateStages
	plated := w.CarWorth(w.Player.Car)
	if plated.Driving <= bare.Driving {
		t.Fatalf("plate cost nothing in the quote: %d against %d", plated.Driving, bare.Driving)
	}
	if plated.Driving >= plated.Walking {
		t.Fatal("a plated car was quoted as slower than walking, which is not what plate does")
	}
}
