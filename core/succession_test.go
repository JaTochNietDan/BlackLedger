package core

import (
	"strings"
	"testing"
)

// Three faults in one night were the same shape: something reaching for a
// person where a role was meant.
//
//   - The second job crashed the request when Mara was dead, because the
//     proposal named her by role, the role was empty, and the nil scene was
//     read anyway.
//   - Ten strings said "Buy Mara a coffee" while the action was already aimed
//     at whoever holds the job, so the player was told one thing and the city
//     did another.
//   - The editor's desk stayed empty for the rest of the campaign, because the
//     city asked whether the man existed rather than whether the office was
//     filled, and a dead man exists.
//
// So this kills everybody the campaign was written with — every seeded role and
// every office — and then runs the city for a fortnight. Nothing may crash, and
// every job has to be done by somebody at the end of it.
func TestTheCityOutlivesEverybodyItWasWrittenWith(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 30, 20000

	buried := 0
	for _, r := range roles {
		if n := w.NPC(r.Seed); n != nil && !n.Dead {
			n.Dead = true
			buried++
		}
	}
	for _, o := range Officials() {
		if n := w.NPC(o.ID); n != nil && !n.Dead {
			n.Dead = true
			buried++
		}
	}
	if buried < 5 {
		t.Fatalf("only %d of the people this campaign was written with were buried", buried)
	}

	// A fortnight of the city, which is where the crash was: the second job is
	// offered on the day the player has done two, and it read a scene that was
	// never made.
	w.Player.JobCount = 2
	for day := 0; day < 14; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
		w.OfferIfReady()
	}

	for _, r := range roles {
		if w.Holder(r.ID) == nil {
			t.Errorf("nobody is the %s a fortnight after the first one died", r.Title)
		}
	}
	for _, o := range Officials() {
		who := w.OfficeHolder(o.ID)
		if who == nil {
			t.Errorf("nobody is the %s a fortnight after the first one died", o.Role)
			continue
		}
		if who.Dead {
			t.Errorf("the %s is a dead man", o.Role)
		}
	}
	t.Logf("buried %d, and every job is being done a fortnight later", buried)
}

// And nothing the player can read names one of them by name. A role is filled
// by whoever holds it, so a line that names the first holder is a line that
// goes on naming a dead woman.
func TestNothingTheCityShowsNamesTheFirstHolderByName(t *testing.T) {
	t.Parallel()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Respect, w.Player.Cash = 100, 30, 20000

	// Who the campaign starts with, by name.
	seeded := map[string]string{}
	for _, r := range roles {
		if r.SeedName != "" {
			seeded[r.SeedName] = r.Title
		}
	}
	for _, o := range Officials() {
		seeded[o.Name] = o.Role
	}
	if len(seeded) < 5 {
		t.Fatal("no seeded names to check against")
	}
	// Bury them all, so anybody still named is named wrongly.
	for name := range seeded {
		for i := range w.NPCs {
			if w.NPCs[i].Name == name {
				w.NPCs[i].Dead = true
			}
		}
	}
	w.FillRoles()
	w.ensureOfficials()

	// Every action the city offers, at every address the player can reach.
	said := 0
	for _, l := range Locations {
		if l.District > w.District {
			continue
		}
		w.Player.Location, w.Event = l.ID, nil
		for _, a := range w.Actions(l.ID) {
			said++
			for name, role := range seeded {
				for what, text := range map[string]string{"label": a.Label, "detail": a.Detail, "reason": a.Reason} {
					if strings.Contains(text, name) {
						t.Errorf("%s at %s names %s, who is dead, where the %s is meant: %q",
							what, l.ID, name, role, text)
					}
				}
			}
		}
	}
	if said < 100 {
		t.Fatalf("only %d actions read across the city, so this measures nothing", said)
	}
	t.Logf("%d actions read, none of them naming somebody the city has buried", said)
}
