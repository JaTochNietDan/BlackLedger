package core

import "fmt"

// Everybody in this city stood still. A person had an address and kept it for
// life: the woman who ran the laundry was inside the laundry at every hour of
// every day, and the man who had spent six weeks deciding he hated somebody
// across town never once walked over to look at him. The only thing that moved
// anybody was a takeover, and it moved them instantly.
//
// This gives the city's people errands. Twice a day, somebody with a reason to
// be elsewhere sets off, and until they get there they are on the street rather
// than in either building — which is the first thing in this game that makes
// the city look like it is being lived in rather than staffed.
//
// Journeys are on foot. Nobody in this city but the player has a car.

// Coming is somebody arriving in or leaving the room the player is standing in.
// Everything else in this city happens at a distance and is read about; this is
// the one thing that happens in front of them.
type Coming struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Where   string `json:"where"`
	To      string `json:"to,omitempty"`
	Leaving bool   `json:"leaving,omitempty"`
	Note    string `json:"note"`
}

// noticed records a coming or going, but only where the player can see it. A
// room they are not in changes without remark, which is what a room they are
// not in does.
func (w *World) noticed(n *NPC, leaving bool) {
	if !w.Player.Alive || n.Dead {
		return
	}
	// Either way this is where they now stand: SetOut leaves Location at the
	// door they walked out of, and Arrivals has already put them down.
	at := n.Location
	if at != w.Player.Location {
		return
	}
	here, _ := PlaceByID(at)
	c := Coming{ID: n.ID, Name: n.Name, Where: at, Leaving: leaving}
	if leaving {
		to, ok := PlaceByID(n.Heading)
		if !ok {
			return
		}
		c.To = to.Name
		c.Note = n.Name + " leaves " + here.Name
		if n.Errand != "" {
			c.Note += ", " + n.Errand
		}
		c.Note += "."
	} else {
		c.Note = n.Name + " comes in."
	}
	w.Comings = append(w.Comings, c)
}

// Journeying is one person between two addresses, for anybody who needs to
// draw the street or say what is happening on it.
type Journeying struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	FromID  string `json:"from_id"`
	From    string `json:"from"`
	ToID    string `json:"to_id"`
	To      string `json:"to"`
	Because string `json:"because"`
	// Minutes is how long they still have to walk, and Progress how far along
	// they are, so a view can put them somewhere sensible between the two.
	Minutes  int     `json:"minutes"`
	Progress float64 `json:"progress"`
	Yours    bool    `json:"yours"`
}

// Travelling reports whether somebody is out on the street rather than
// standing somewhere. Their Location stays where they set off from, so
// everything that reasons about where a person belongs keeps working; only
// what is in a room, and what is on the street, changes.
func (w *World) Travelling(n *NPC) bool {
	return n != nil && n.Heading != "" && n.Sets == 0 && n.Arrives > w.Minute && !n.Dead
}

// setsOff spreads a shift change over the hours a shift change takes. Fixed
// per person, so the one who always leaves early always leaves early.
func setsOff(id string) int {
	sum := 0
	for i := 0; i < len(id); i++ {
		sum = sum*131 + int(id[i])
	}
	if sum < 0 {
		sum = -sum
	}
	return 1 + sum%300
}

// errand is somewhere a person has reason to be, and the reason.
type errand struct {
	where   string
	because string
}

// wanted is where this person ought to be, in priority order. The first reason
// that fits is the one they act on, and the order is the order a person would
// weigh them: the work that feeds you, the office you hold, the man you have
// not forgiven, the people you answer to.
func (w *World) wanted(n *NPC) (errand, bool) {
	// Whoever runs a place is supposed to be in it.
	for _, l := range Locations {
		if n.Role == "Runs "+l.Name && n.Location != l.ID {
			return errand{l.ID, "opening up at " + l.Name}, true
		}
	}
	// A title comes with somewhere to hold it.
	for _, r := range roles {
		if n.Role == r.Title && r.Where != "" && n.Location != r.Where {
			place, ok := PlaceByID(r.Where)
			if ok {
				return errand{r.Where, "due at " + place.Name}, true
			}
		}
	}
	// Something carried long enough to act on takes a man across town. This is
	// the first time a grudge makes anybody do anything of their own accord.
	for _, g := range w.Grudges {
		if g.Holder != n.ID || g.Weight < GrudgeActs {
			continue
		}
		target := w.NPC(g.Against)
		if target == nil || target.Dead || target.Location == "" || target.Location == n.Location || w.Travelling(target) {
			continue
		}
		return errand{target.Location, "looking for " + target.Name + " over " + g.Because}, true
	}
	// A car with the glass out of it. Somebody who can find the fee takes it in
	// rather than driving it broken, which is why a garage has anybody at the
	// counter at all: the trade is people standing there, not a figure moving
	// once a day. It outranks minding a holding for the morning — a family's
	// ground is still there in the afternoon.
	if n.Hurt && n.Car > 0 && n.Purse >= GlassCost {
		if id := w.theGarage(); id != "" && id != n.Location {
			if place, ok := PlaceByID(id); ok {
				return errand{id, "taking the car in to " + place.Name}, true
			}
		}
	}
	// An organization's ground has to be minded by somebody. A family holding
	// with nobody standing in it pulls one of that family's people across the
	// city — which is also what makes a war visible on the street: take a
	// holding and their man walks out of it, and one of yours walks in.
	if n.Faction != "" && n.Rank < RankLeader && !w.soleMinder(n) {
		for _, id := range w.FamilyHoldings(n.Faction) {
			if id == n.Location || w.minded(id, n) {
				continue
			}
			place, ok := PlaceByID(id)
			if !ok {
				continue
			}
			return errand{id, "minding " + place.Name}, true
		}
	}
	// The day ends and people go out. This is last because it is the weakest
	// reason anybody has to be anywhere, and it must never pull somebody off
	// work they are needed for.
	if e, ok := w.routine(n); ok {
		return e, true
	}
	// A summons from whoever is above them was tried here and removed. It
	// pulled every member to their leader, so a man who had just walked across
	// the city to mind a holding became its only minder and was immediately
	// sent for again: the same four people walked back and forth for the whole
	// game. Nobody in this city moves without a reason that ends.
	return errand{}, false
}

// SetOut sends the people who have somewhere to be. Called twice a day from
// the clock, so the city is not permanently in motion: a person makes at most
// one journey in a half-day, and most people make none.
func (w *World) SetOut() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || n.Location == "" || w.Travelling(n) || n.Held > w.Minute {
			continue
		}
		n.Heading, n.Arrives, n.Errand, n.Sets = "", 0, "", 0
		// Where a person is standing in the daytime is where their day is,
		// unless something better has already claimed them.
		if n.Post == "" && !Evening(w.Minute) {
			w.keepPost(n, n.Location)
		}
		where, ok := w.wanted(n)
		if !ok {
			continue
		}
		if _, exists := PlaceByID(where.where); !exists || where.where == n.Location {
			continue
		}
		n.Heading = where.where
		n.Errand = where.because
		n.Sets = w.Minute + setsOff(n.ID)
		n.Arrives = n.Sets + TravelMinutes(n.Location, where.where)
	}
}

// Arrivals puts down everybody whose walk is over. Called on every step of the
// clock rather than at a boundary, because a journey is shorter than a day and
// the player should be able to watch one finish.
func (w *World) Arrivals() {
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Heading == "" {
			continue
		}
		if n.Dead {
			n.Heading, n.Arrives, n.Errand, n.Sets = "", 0, "", 0
			continue
		}
		// Still deciding to go. The room they are in is the room they are in
		// until they walk out of it, and the player watching sees them go at
		// the moment they go.
		if n.Sets > w.Minute {
			continue
		}
		if n.Sets > 0 {
			n.Sets = 0
			w.noticed(n, true)
		}
		if n.Arrives > w.Minute {
			continue
		}
		n.Location = n.Heading
		n.Heading, n.Arrives, n.Errand, n.Sets = "", 0, "", 0
		// Arriving anywhere in the daytime is arriving at work: whatever
		// reason brought them, this is now where their day is, and the evening
		// has somewhere to send them back from. An evening arrival is a drink
		// and changes nothing.
		if !Evening(w.Minute) {
			w.keepPost(n, n.Location)
		}
		w.noticed(n, false)
	}
}

// OnTheStreet is everybody currently between two addresses.
func (w *World) OnTheStreet() []Journeying {
	out := []Journeying{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !w.Travelling(n) {
			continue
		}
		from, okFrom := PlaceByID(n.Location)
		to, okTo := PlaceByID(n.Heading)
		if !okFrom || !okTo {
			continue
		}
		total := max(1, TravelMinutes(n.Location, n.Heading))
		left := max(0, n.Arrives-w.Minute)
		out = append(out, Journeying{
			ID: n.ID, Name: n.Name,
			FromID: n.Location, From: from.Name,
			ToID: n.Heading, To: to.Name,
			Because:  n.Errand,
			Minutes:  left,
			Progress: progress(total, left),
			Yours:    n.Faction != "" && n.Faction == w.PlayerOrganizationID(),
		})
	}
	return out
}

// StreetNote is one line for anybody who wants to say what the street looks
// like without drawing it.
func (w *World) StreetNote() string {
	street := w.OnTheStreet()
	switch len(street) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%s is out on the street, %s.", street[0].Name, street[0].Because)
	}
	if len(street) == 2 {
		return fmt.Sprintf("%s and %s are out on the street.", street[0].Name, street[1].Name)
	}
	return fmt.Sprintf("%s and %d others are out on the street.", street[0].Name, len(street)-1)
}

// progress is how far along a walk is, kept between nothing and all of it.
func progress(total, left int) float64 {
	done := float64(total-left) / float64(total)
	if done < 0 {
		return 0
	}
	if done > 1 {
		return 1
	}
	return done
}

// minded reports whether somebody of this person's own organization is already
// standing in a place, or on their way to it. Two people do not walk to the
// same empty room.
func (w *World) minded(id string, except *NPC) bool {
	for i := range w.NPCs {
		other := &w.NPCs[i]
		if other.Dead || other.ID == except.ID || other.Faction != except.Faction {
			continue
		}
		if other.Heading == id || (other.Location == id && !w.Travelling(other)) {
			return true
		}
	}
	return false
}

// soleMinder reports whether this person is the only one of their own
// organization standing on ground it holds. Somebody minding a place on their
// own has a reason to stay; somebody standing in a room with three of their own
// people, while another of their family's addresses has nobody in it at all,
// does not. Without this the city either shuffles its people between its own
// buildings forever or never moves them at all.
func (w *World) soleMinder(n *NPC) bool {
	prop := w.Properties[n.Location]
	if prop == nil || prop.Owner != n.Faction {
		return false
	}
	for i := range w.NPCs {
		other := &w.NPCs[i]
		if other.Dead || other.ID == n.ID || other.Faction != n.Faction {
			continue
		}
		if other.Location == n.Location && !w.Travelling(other) {
			return false
		}
	}
	return true
}
