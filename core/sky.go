package core

// The weather over Bellwether.
//
// This lives in the core rather than in the view for one reason: what the sky
// is doing is a fact about the world, and the view is not allowed to invent
// facts. The city view could perfectly well have decided for itself that it
// was raining, and then the rain would have existed only where somebody was
// looking — a different sky in the street view than in the room, and none at
// all in the newspaper.
//
// Nothing in the rules turns on it yet. That is deliberate rather than
// unfinished: the fact is established first and cheaply, so that when working
// a door in the rain is worth a modifier, the modifier has somewhere to live
// and every screen already agrees about what day it got wet.

// Sky is what it is doing outside, and what the streets are like because of it.
type Sky struct {
	Kind string  `json:"kind"` // clear, overcast, rain or fog
	Wet  float64 `json:"wet"`  // 0 dry to 1 running with water
}

// skySeed is a stable per-campaign number, so two campaigns do not get the
// same fortnight of weather and one campaign's Tuesday is always the same
// Tuesday however many times it is loaded.
func skySeed(id string) uint32 {
	var h uint32 = 2166136261
	for _, c := range id {
		h ^= uint32(c)
		h *= 16777619
	}
	return h
}

// skyOnDay is the weather for one day, with no memory of any other day.
func skyOnDay(day int, seed uint32) string {
	h := seed ^ uint32(day)*2654435761
	h ^= h >> 15
	h *= 2246822519
	h ^= h >> 13
	switch roll := h % 100; {
	case roll < 46:
		return "clear"
	case roll < 74:
		return "overcast"
	case roll < 92:
		return "rain"
	default:
		return "fog"
	}
}

// SkyOn is the sky on a given day of a given campaign. A street does not dry
// the moment the rain stops, so the day after rain is still wet underfoot,
// which is the whole reason this takes yesterday into account at all.
func SkyOn(day int, seed uint32) Sky {
	kind := skyOnDay(day, seed)
	wet := 0.0
	switch kind {
	case "rain":
		wet = 1
	case "fog":
		wet = .3
	}
	if kind != "rain" && day > 0 && skyOnDay(day-1, seed) == "rain" {
		if wet < .45 {
			wet = .45 // still drying
		}
	}
	return Sky{Kind: kind, Wet: wet}
}

// Sky is today's weather in this campaign.
func (w *World) Sky() Sky { return SkyOn(w.Minute/1440, skySeed(w.ID)) }
