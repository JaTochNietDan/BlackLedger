package core

// A business's staff was a number. You hired a pair of hands, the wage bill went
// up, and nobody in Bellwether had a job: the person behind the counter of a
// place the player owned did not exist, could not be talked to, killed, robbed,
// poached or arrested, and the city's own people had no reason to be anywhere in
// the daytime except the one the routine invented for them.
//
// The count stays, because everything that reads it — what a place can handle,
// what the wages are, whether it is short-handed — is right to read a count. It
// is now the length of a list of people who live here. Somebody who dies, or is
// taken in and held, is not behind that counter, and the position is empty until
// it is filled again.

// EmployerOf is where this person works, and nothing if they do not.
func (w *World) EmployerOf(id string) string {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		for _, who := range prop.Hands {
			if who == id {
				return l.ID
			}
		}
	}
	return ""
}

// AtWork is everybody the player employs, which is a thing worth being able to
// ask now that it has an answer.
func (w *World) AtWork() []*NPC {
	out := []*NPC{}
	for _, l := range Locations {
		if !w.Own(l.ID) || w.Properties[l.ID] == nil {
			continue
		}
		for _, who := range w.Properties[l.ID].Hands {
			if n := w.NPC(who); n != nil && !n.Dead {
				out = append(out, n)
			}
		}
	}
	return out
}

// takeOn finds somebody in this city to stand behind that counter. Nobody with
// a job already, nobody a family owns, nobody who is dead or inside. Returns
// the empty string when there is nobody, which is a real answer: a city can run
// out of people willing to work for you.
func (w *World) takeOn(id string) string {
	best, bestScore := "", -1
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Held > w.Minute || IsOfficial(n.ID) || n.Rank >= RankLeader {
			continue
		}
		// Somebody who already holds a position keeps it, and the city keeps
		// them where it is: a driver taken on at a laundry had the laundry
		// written down as their post and went on standing in the bar for the
		// rest of the week, because a role holder is exempt from the errand
		// that would walk them there.
		if w.keepsPost(n) {
			continue
		}
		if w.EmployerOf(n.ID) != "" {
			continue
		}
		// Somebody another family already has is not looking for a counter job:
		// they have one, and the city sends them off to mind their own
		// holdings, so a laundry that hired one had a position filled by
		// somebody who was never behind the counter. Measured over a week of
		// daytimes: one of three hands was never once in the room.
		if n.Faction != "" && n.Faction != w.PlayerOrganizationID() {
			continue
		}
		// Somebody already spending their day there is the obvious hire, then
		// somebody who trusts you, then anybody at all.
		score := n.Trust
		if n.Post == id || n.Location == id {
			score += 50
		}
		if score > bestScore {
			best, bestScore = n.ID, score
		}
	}
	return best
}

// putToWork sets somebody behind a counter: the day is spent there, which is
// what having a job means to everything else in this city.
func (w *World) putToWork(who, id string) {
	n := w.NPC(who)
	if n == nil {
		return
	}
	prop := w.Properties[id]
	prop.Hands = append(prop.Hands, who)
	n.Post = id
	// They start today. Writing down where somebody works and leaving them
	// across the city is how a counter ends up staffed by nobody.
	if !w.Travelling(n) && !Evening(w.Minute) {
		n.Location = id
	}
}

// letGo takes the last one on and gives them their day back.
func (w *World) letGo(id string) {
	prop := w.Properties[id]
	if len(prop.Hands) == 0 {
		return
	}
	who := prop.Hands[len(prop.Hands)-1]
	prop.Hands = prop.Hands[:len(prop.Hands)-1]
	if n := w.NPC(who); n != nil && n.Post == id {
		n.Post = ""
	}
}

// EmptyChairs takes the dead off the books and puts a name to any position that
// has not got one. A position held by somebody who is not coming in is not a
// position that is filled, and the count has to say so or the place goes on
// handling work nobody is there to do.
//
// The naming half is also how a save written before anybody had a job acquires
// one, and how a business the player has just taken over gets the people who
// were already working in it.
func (w *World) EmptyChairs() {
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		if prop == nil {
			continue
		}
		kept := prop.Hands[:0]
		for _, who := range prop.Hands {
			if n := w.NPC(who); n != nil && !n.Dead {
				kept = append(kept, who)
			}
		}
		if len(kept) != len(prop.Hands) {
			prop.Hands = kept
			prop.Staff = min(prop.Staff, len(kept))
		}
		if !w.Own(l.ID) {
			continue
		}
		for len(prop.Hands) < prop.Staff {
			who := w.takeOn(l.ID)
			if who == "" {
				// Nobody left in this city willing to stand behind it, which is
				// a real answer rather than a reason to invent somebody.
				prop.Staff = len(prop.Hands)
				break
			}
			w.putToWork(who, l.ID)
		}
	}
}

// HandsDescription is who stands behind that counter, for the room to draw.
// Nothing at all where the player has no business knowing, which is everywhere
// they do not hold: a rival's payroll is not public.
func (w *World) HandsDescription(id string) []map[string]any {
	prop := w.Properties[id]
	if prop == nil || !w.Own(id) || len(prop.Hands) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(prop.Hands))
	for _, who := range prop.Hands {
		n := w.NPC(who)
		if n == nil {
			continue
		}
		out = append(out, map[string]any{
			"id": n.ID, "name": n.Name, "role": n.Role,
			"here": n.Location == id && !w.Travelling(n),
		})
	}
	return out
}
