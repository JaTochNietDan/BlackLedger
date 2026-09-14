package core

import (
	"fmt"
	"hash/fnv"
)

// Disposition remains stable across reloads and consumes no simulation RNG.
// This is the default arrangement until individual funeral wishes are modeled.
func deathDisposition(personID string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(personID))
	if h.Sum32()%3 == 0 {
		return "crematorium"
	}
	return "cemetery"
}

// A service goes to an operating provider, preferring the deceased's district.
// Ownership does not determine who gets the work.
func (w *World) deathServiceProvider(kind string, person *NPC) string {
	district := -1
	if person != nil {
		if p, ok := PlaceByID(person.Home); ok {
			district = p.District
		}
	}
	fallback := ""
	for _, p := range Locations {
		prop := w.Properties[p.ID]
		if p.Kind != kind || prop == nil || prop.Condition < 60 || prop.Trouble || prop.Staff == 0 || prop.Supply == 0 {
			continue
		}
		if p.District == district {
			return p.ID
		}
		if fallback == "" {
			fallback = p.ID
		}
	}
	return fallback
}

// Paid from the funeral's existing service allowance, never an extra bill.
// One third pays for care, two thirds for disposition. Half covers consumables
// and outside expenses; only the remainder is the provider's margin.
func (w *World) deathServicePayments(person *NPC, allowance int) {
	if person == nil || allowance <= 0 {
		return
	}
	care := allowance / 3
	for _, service := range []struct {
		kind string
		paid int
	}{{"mortuary", care}, {deathDisposition(person.ID), allowance - care}} {
		id := w.deathServiceProvider(service.kind, person)
		if id == "" || service.paid == 0 {
			continue
		}
		margin := service.paid / 2
		if w.Own(id) {
			w.Earn(margin)
			place, _ := PlaceByID(id)
			w.Log("Service accounts at "+place.Name, fmt.Sprintf("Work for %s brought in $%d. Supplies and outside costs left $%d for your accounts.", person.Name, service.paid, margin), "business")
		} else {
			w.changeBusinessFunds(id, margin)
		}
		w.ShiftCustom(id, "a quiet week", BurialTrade)
	}
}
