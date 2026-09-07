package main

import (
	"blackledger/core"
	"fmt"
	"strings"
)

func accessibleJobLocations(w *core.World) []map[string]string {
	result := []map[string]string{}
	for _, place := range core.Locations {
		if place.District <= w.District {
			result = append(result, map[string]string{"id": place.ID, "name": place.Name})
		}
	}
	return result
}

func validateJobLocation(w *core.World, p core.Proposal) error {
	place, ok := core.PlaceByID(p.Location)
	if !ok || place.District > w.District {
		return fmt.Errorf("location must be an ID from accessible_job_locations")
	}
	if !strings.Contains(strings.ToLower(p.Body), strings.ToLower(place.Name)) {
		return fmt.Errorf("body must explicitly name the selected venue %q and describe the task there", place.Name)
	}
	// This rejects explicit locked-place references, not every possible semantic
	// contradiction or alias. The saved location remains the authoritative venue.
	prose := strings.ToLower(p.Title + " " + p.Body)
	for _, approach := range p.Approaches {
		prose += " " + strings.ToLower(approach.Label)
	}
	for _, other := range core.Locations {
		if other.District > w.District && strings.Contains(prose, strings.ToLower(other.Name)) {
			return fmt.Errorf("%s is outside the accessible districts; write this request at %s", other.Name, place.Name)
		}
	}
	return nil
}
