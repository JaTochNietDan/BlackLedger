package core

import "fmt"

// What a journey across this city actually is.
//
// The bar that said you were walking somewhere sat at the bottom of the city
// pane, below the fold on a short window, and said the name of the place and a
// number of minutes. That was fine while the street was scenery. It is not
// scenery any more: it is the one stretch of the city with no walls, no door
// and nobody who knows you, which is why a car can be plated and why plate is
// worth anything.
//
// So the crossing is a fact the core states before the player sets off, rather
// than a caption the interface writes afterwards.

// Crossing describes the journey between two addresses: how long it takes, what
// the player is making it in, and whether anybody is known to be looking for
// them while they are out in it.
func (w *World) Crossing(from, to string) map[string]any {
	minutes := w.Journey(from, to)
	driving := w.Driving()
	plate := 0
	if driving {
		plate = w.Plating()
	}
	warned := false
	for _, plot := range w.Plots {
		if plot.Life == w.Life && plot.Kind == "hit" && plot.Known {
			warned = true
		}
	}
	place, _ := PlaceByID(to)
	out := map[string]any{
		"to": place.Name, "to_id": to, "minutes": minutes,
		"driving": driving, "plate": plate, "plate_max": PlateStages,
		"warned": warned, "note": "",
	}
	switch {
	case warned && driving && plate > 0:
		out["note"] = fmt.Sprintf("Somebody is looking for you. The street is the one place nothing covers you, and %s of plate is what you have out there.", counted(plate, "stage", "stages"))
	case warned && driving:
		out["note"] = "Somebody is looking for you, and the street is the one place nothing covers you. There is no plate on this car."
	// Before the bare "on foot", or it can never be reached: a switch takes the
	// first case that matches and being warned matches both.
	case warned && w.RidingWithTheCabs():
		out["note"] = "Somebody is looking for you, and one of your own drivers is taking you. A cab is not cover, but it is not the pavement either."
	case warned:
		out["note"] = "Somebody is looking for you and you are crossing the street on foot, where nothing covers anybody."
	case driving && plate > 0:
		out["note"] = fmt.Sprintf("%s of plate between you and the street.", capitalise(counted(plate, "stage", "stages")))
	case driving:
		out["note"] = "Nothing on this car but its own doors."
	case w.RidingWithTheCabs():
		out["note"] = "One of your own drivers has you, which is faster than the pavement and no safer."
	default:
		out["note"] = "On foot, and the street is the one place in the city nothing covers anybody."
	}
	return out
}
