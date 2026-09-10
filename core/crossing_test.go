package core

import (
	"strings"
	"testing"
)

// "The walking between buildings simulation is not that great right now because
// we don't have the map working properly and the little bar that explains that
// you're traveling between buildings is at the bottom of the page often below
// the fold."
//
// The bar was in the wrong place and said almost nothing. It matters more than
// it did: the street between two addresses is now the one stretch of the city
// where a hit can catch you cold, so what you are crossing it in is a fact the
// player should have before they set off, not after.

func crosser(t *testing.T) *World {
	t.Helper()
	w := New(37)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 5000
	w.Player.Location = "bar"
	return w
}

func TestTheCrossingSaysWhatYouAreCrossingItIn(t *testing.T) {
	w := crosser(t)
	on := w.Crossing("bar", "market")
	if on["minutes"].(int) <= 0 {
		t.Fatalf("a journey that takes no time: %+v", on)
	}
	if on["driving"].(bool) {
		t.Fatal("somebody with no car was driving")
	}
	foot := on["minutes"].(int)

	w.Player.Car, w.Player.CarWear = 2, 100
	w.Player.Fuel, w.Player.Fuelled = FuelFull, max(1, w.Minute)
	on = w.Crossing("bar", "market")
	if !on["driving"].(bool) {
		t.Fatal("a car in the road and nobody driving it")
	}
	if on["minutes"].(int) >= foot {
		t.Fatalf("driving took %d against %d on foot", on["minutes"].(int), foot)
	}
	if on["plate"].(int) != 0 {
		t.Fatalf("a bare car carried plate: %+v", on)
	}
	w.Player.Plate = PlateStages
	if w.Crossing("bar", "market")["plate"].(int) == 0 {
		t.Fatal("plate on the car did not reach the crossing")
	}
}

func TestTheCrossingSaysWhoIsLookingForYou(t *testing.T) {
	w := crosser(t)
	quiet := w.Crossing("bar", "market")
	if quiet["warned"].(bool) {
		t.Fatalf("a quiet city warned about nothing: %+v", quiet)
	}
	w.RetaliationFrom("bellandi")
	for i := range w.Plots {
		w.Plots[i].Known = true
	}
	loud := w.Crossing("bar", "market")
	if !loud["warned"].(bool) {
		t.Fatal("somebody had a price on the player and the street said nothing")
	}
	if note, _ := loud["note"].(string); note == "" || !strings.Contains(strings.ToLower(note), "street") {
		t.Fatalf("the warning does not say where the risk is: %q", note)
	}
}

// The holder sets the house limit, and the figure has to be typed. The action
// was offered with no field on it, so it sent nothing and was refused every
// time: "the max limit should be defined by the casino owner dynamically."
func TestTheHolderCanActuallyNameTheLimit(t *testing.T) {
	w := houseKeeper(t)
	w.Event, w.District = nil, 9
	a := actionByID(w.Actions("casino"), "limit")
	if a == nil {
		t.Fatal("the holder of a casino cannot set its limit")
	}
	if a.Sum == nil {
		t.Fatal("there is nowhere to type the figure")
	}
	if a.Sum.Least != HouseLimitFloor || a.Sum.Most != HouseLimitCeiling {
		t.Fatalf("the field runs %d..%d and the rule runs %d..%d",
			a.Sum.Least, a.Sum.Most, HouseLimitFloor, HouseLimitCeiling)
	}
	if err := w.SetLimit("casino", a.Sum.Preset); err != nil {
		t.Fatalf("the figure the field starts on is refused: %v", err)
	}
}
