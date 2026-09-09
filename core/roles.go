package core

import "fmt"

// Mara is the fixer at Saint Agnes, Leo is a driver, Harlow is the detective who
// keeps turning up. The game refers to all three by name, which means all three
// were immortal by accident: killing Mara would have left every scene that
// speaks through her without a speaker.
//
// A role is a job in this city rather than a person. When whoever holds one is
// gone, somebody else is doing it by the end of the week — a different name, a
// different voice, and no memory of anything the player did for the last one.

// Role is a job somebody in this city does, and what happens when nobody does
// it.
type Role struct {
	ID string
	// Title is the job, Where it is done, and Announce what the city says when
	// somebody new takes it up.
	Title, Where, Announce string
	// Seed is who holds it on the first morning.
	Seed, SeedName, SeedVoice string
}

var roles = []Role{
	{ID: "fixer", Title: "Fixer", Where: "bar", Seed: "mara", SeedName: "Mara Bell", SeedVoice: "af_heart",
		Announce: "Somebody has to be the one people ask, and at Saint Agnes it is %s now."},
	{ID: "driver", Title: "Driver", Where: "bar", Seed: "leo", SeedName: "Leo Carver", SeedVoice: "am_michael",
		Announce: "%s drives now, and takes the work that comes with it."},
	{ID: "detective", Title: "City detective", Where: "market", Seed: "harlow", SeedName: "Detective Harlow", SeedVoice: "bm_lewis",
		Announce: "%s has the caseload that was somebody else's last week."},
}

// RoleByID is one of them.
func RoleByID(id string) (Role, bool) {
	for _, r := range roles {
		if r.ID == id {
			return r, true
		}
	}
	return Role{}, false
}

// Holder is whoever is doing a job now. It is the only way anything should
// reach these people: their ids belong to the person, and the person changes.
func (w *World) Holder(role string) *NPC {
	r, ok := RoleByID(role)
	if !ok {
		return nil
	}
	// Whoever the campaign started with, while they are alive.
	if seed := w.NPC(r.Seed); seed != nil && !seed.Dead {
		return seed
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; !n.Dead && n.Role == r.Title {
			return n
		}
	}
	return nil
}

// HolderID is who to put in a scene as the speaker, or "" when nobody is doing
// that job at this moment — which lasts until the end of the day.
func (w *World) HolderID(role string) string {
	if n := w.Holder(role); n != nil {
		return n.ID
	}
	return ""
}

// FillRoles gives every empty job to somebody. Called on the daily tick, so a
// killing on Tuesday is answered by a stranger on Wednesday.
func (w *World) FillRoles() {
	for _, r := range roles {
		if w.Holder(r.ID) != nil {
			continue
		}
		// The first person to hold a job is the one the campaign was written
		// with. Only once they are gone does the city find somebody itself —
		// the first build promoted a stranger on the first morning and lost
		// Detective Harlow entirely.
		if r.Seed != "" && w.NPC(r.Seed) == nil {
			w.NPCs = append(w.NPCs, NPC{
				ID: r.Seed, Name: r.SeedName, Role: r.Title, Voice: r.SeedVoice,
				Color: "#7c8791", Location: r.Where, Rank: RankAssociate,
				Ambition: 40, Skill: 50,
			})
			continue
		}
		successor := w.nearestTo(r.Where)
		if successor == nil {
			successor = w.AddCivilian()
		}
		if successor == nil {
			continue
		}
		successor.Role, successor.Location = r.Title, r.Where
		successor.Faction, successor.Rank = "", RankAssociate
		successor.Trust = 0 // a stranger is a stranger, whatever the last one knew
		w.Log("Somebody else is doing that job now", fmt.Sprintf(r.Announce, successor.Name), "personal")
		w.Report("politics", upper(successor.Name)+" TAKES OVER AS "+upper(r.Title),
			fmt.Sprintf("%s is doing the work of the %s now. The change was not explained and nobody has said where the last one went.", successor.Name, lowerFirst(r.Title)))
	}
}

// nearestTo finds somebody unaffiliated who could take a job on, preferring
// whoever is already standing where it is done.
func (w *World) nearestTo(place string) *NPC {
	var fallback *NPC
	// Civilians already excludes anybody doing one of these jobs.
	for _, n := range w.Civilians() {
		if IsOfficial(n.ID) {
			continue
		}
		if n.Location == place {
			return n
		}
		if fallback == nil {
			fallback = n
		}
	}
	return fallback
}

// isRoleHolder reports whether somebody is already doing one of these jobs, so
// the fixer is never also the detective.
func (w *World) isRoleHolder(n *NPC) bool {
	for _, r := range roles {
		if n.Role == r.Title {
			return true
		}
	}
	return false
}

// RoleDescription is who is doing what, for the interface.
func (w *World) RoleDescription() []map[string]any {
	out := []map[string]any{}
	for _, r := range roles {
		entry := map[string]any{"role": r.ID, "title": r.Title, "name": "Nobody"}
		if n := w.Holder(r.ID); n != nil {
			entry["name"], entry["id"] = n.Name, n.ID
		}
		out = append(out, entry)
	}
	return out
}
