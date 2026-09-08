package core

import (
	"fmt"
	"strings"
)

// Nobody in this city is scenery. Everyone with a name has an organization, a
// place they are usually found, a standing inside that organization, and a life
// that can end. Whatever can be done to the player can be done to them, and by
// them. When someone at the top dies, the people below them do not wait.

// Ranks. Higher is closer to the top; the leader holds RankLeader.
const (
	RankAssociate  = 10
	RankSoldier    = 25
	RankLieutenant = 55
	RankLeader     = 100
)

// firstNames and surnames generate the city's ordinary people. A name is drawn
// once and then belongs to that person for as long as they live.
// Kept apart so a character is not given a voice that contradicts how the rest
// of the city refers to them.
// A city of hundreds needs enough names for hundreds. Fourteen first names and
// twelve surnames gave 168 combinations, which is a hard ceiling on how many
// people this city could ever contain — and newPersonName gives up after sixty
// attempts, so it started failing long before that.
var mensFirstNames = []string{"Gio", "Aldo", "Emil", "Ivo", "Luca", "Anton", "Piet",
	"Bruno", "Cesare", "Dante", "Ennio", "Fausto", "Gustav", "Hugo", "Igor",
	"Janos", "Karel", "Lorenz", "Marek", "Nico", "Otto", "Pavel", "Rudi",
	"Sandor", "Tomas", "Ugo", "Valter", "Wim", "Zoltan", "Bela"}

var womensFirstNames = []string{"Nina", "Perla", "Rosa", "Greta", "Mirela", "Sofia", "Dora",
	"Alma", "Bianca", "Clara", "Dita", "Elsa", "Franca", "Gina", "Hedda",
	"Ilona", "Jelena", "Katia", "Lidia", "Magda", "Nadia", "Olga", "Pia",
	"Renata", "Stella", "Tilda", "Ursa", "Vera", "Wanda", "Zora"}

var peopleSurnames = []string{"Costa", "Varga", "Lenz", "Moreau", "Sabbatini",
	"Novak", "Hale", "Duarte", "Weiss", "Petrov", "Ferro", "Blum",
	"Aldini", "Berger", "Corvi", "Draga", "Esposito", "Falk", "Gruber",
	"Havel", "Iordan", "Janssen", "Kovac", "Lombardi", "Mraz", "Nagy",
	"Olsen", "Palma", "Quintero", "Rossi", "Steiner", "Toth", "Ulmann",
	"Vance", "Wolf", "Zanetti", "Bassi", "Cerny", "Doyle", "Erdos"}

// A voice belongs to a person for as long as they live, so a character the
// player has heard before sounds the same the next time they speak. Drawn from
// the voices the local synthesis service already provides.
var womensVoices = []string{"af_alloy", "af_aoede", "af_bella", "af_jessica",
	"af_kore", "af_nicole", "af_nova", "af_river", "af_sarah", "bf_alice", "bf_isabella"}

var mensVoices = []string{"am_echo", "am_eric", "am_fenrir", "am_liam",
	"am_onyx", "am_puck", "bm_daniel", "bm_fable"}

// voiceFor picks a voice no living person is already using, so two characters
// in a scene are never the same voice. It falls back to a stable choice once
// every voice is spoken for.
func (w *World) voiceFor(name string) string {
	taken := map[string]bool{}
	for _, n := range w.NPCs {
		if !n.Dead {
			taken[n.Voice] = true
		}
	}
	pool := mensVoices
	first, _, _ := strings.Cut(name, " ")
	for _, womans := range womensFirstNames {
		if strings.EqualFold(first, womans) {
			pool = womensVoices
		}
	}
	sum := 0
	for _, r := range name {
		sum = sum*31 + int(r)
	}
	if sum < 0 {
		sum = -sum
	}
	for offset := 0; offset < len(pool); offset++ {
		candidate := pool[(sum+offset)%len(pool)]
		if !taken[candidate] {
			return candidate
		}
	}
	return pool[sum%len(pool)]
}

// Living reports whether a person is still in the city.
func (n NPC) Living() bool { return !n.Dead }

// People returns everyone still alive, in stable order.
func (w *World) People() []*NPC {
	out := []*NPC{}
	for i := range w.NPCs {
		if w.NPCs[i].Living() {
			out = append(out, &w.NPCs[i])
		}
	}
	return out
}

// Members lists the living people who answer to an organization, strongest
// standing first, so succession and internal politics are deterministic.
func (w *World) Members(faction string) []*NPC {
	out := []*NPC{}
	for _, n := range w.People() {
		if n.Faction == faction {
			out = append(out, n)
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0; j-- {
			a, b := out[j-1], out[j]
			if b.Rank > a.Rank || (b.Rank == a.Rank && (b.Ambition > a.Ambition || (b.Ambition == a.Ambition && b.ID < a.ID))) {
				out[j-1], out[j] = b, a
				continue
			}
			break
		}
	}
	return out
}

func (w *World) personName(name string) bool {
	for _, n := range w.NPCs {
		if n.Name == name {
			return true
		}
	}
	for _, f := range w.Factions {
		if f.Leader == name {
			return true
		}
	}
	return name == w.Player.Name
}

// newPersonName draws an unused name from the reserved world stream.
func (w *World) newPersonName() (string, bool) {
	for attempt := 0; attempt < 60; attempt++ {
		names := mensFirstNames
		if w.WorldRandom() < .5 {
			names = womensFirstNames
		}
		name := names[int(w.WorldRandom()*float64(len(names)))%len(names)] +
			" " + peopleSurnames[int(w.WorldRandom()*float64(len(peopleSurnames)))%len(peopleSurnames)]
		if !w.personName(name) {
			return name, true
		}
	}
	return "", false
}

// AddMember brings a new person into an organization at a given standing. They
// are a full participant from the moment they exist.
func (w *World) AddMember(faction, role string, rank int, base string) *NPC {
	name, ok := w.newPersonName()
	if !ok {
		return nil
	}
	id := fmt.Sprintf("person-%d", len(w.NPCs)+1)
	for w.NPC(id) != nil {
		id += "x"
	}
	w.NPCs = append(w.NPCs, NPC{
		ID: id, Name: name, Role: role, Voice: w.voiceFor(name), Color: "#8d7f6a",
		Faction: faction, Location: base, Rank: rank,
		Ambition: 20 + int(w.WorldRandom()*70), Skill: 25 + int(w.WorldRandom()*60),
	})
	return &w.NPCs[len(w.NPCs)-1]
}

// homeOf is where an organization's people are usually found: its best holding,
// or the street if it holds nothing.
func (w *World) homeOf(faction string) string {
	holdings := w.FamilyHoldings(faction)
	if len(holdings) == 0 {
		return "bar"
	}
	best := holdings[0]
	for _, id := range holdings {
		if w.Properties[id].Income > w.Properties[best].Income {
			best = id
		}
	}
	return best
}

// Kill ends a person's life for a stated reason. If they led an organization,
// the people below them decide what happens next. This is the same function
// whoever did it and whoever it was.
func (w *World) Kill(id, cause string) bool {
	person := w.NPC(id)
	if person == nil || person.Dead {
		return false
	}
	person.Dead = true
	faction := person.Faction
	led := ""
	for i := range w.Factions {
		if w.Factions[i].Leader == person.Name {
			led = w.Factions[i].ID
		}
	}
	w.Log(person.Name+" is dead", cause+" "+describeStanding(person, w)+".", "danger")
	// A man with a title is not a soldier, and the city does not treat him
	// like one.
	if _, official := OfficialByID(person.ID); official {
		w.OfficialKilled(person.ID)
		return true
	}
	// The paper reports a killing without knowing who arranged it.
	headline := strings.ToUpper(person.Name) + " FOUND DEAD"
	if person.Rank >= RankLieutenant {
		headline = strings.ToUpper(person.Name) + " KILLED"
	}
	w.Report("killing", headline, cause+" "+describeStanding(person, w)+". Police say enquiries are continuing.")
	w.witnessKilling(person, cause, headline)
	if led != "" {
		w.Succeed(led)
	} else if faction != "" {
		// Losing experienced people costs an organization strength.
		if f := w.faction(faction); f != nil {
			f.Power = max(10, f.Power-max(1, person.Rank/20))
		}
	}
	return true
}

func describeStanding(n *NPC, w *World) string {
	if n.Faction == "" {
		return "They answered to nobody"
	}
	if f := w.faction(n.Faction); f != nil {
		// A successor's role already carries the organization's name, so
		// appending it again produced "Head of the Russo Outfit of Russo
		// Outfit" in the Herald.
		if strings.Contains(n.Role, f.Name) {
			return "They were " + n.Role
		}
		return "They were " + n.Role + " of " + f.Name
	}
	return "They were " + n.Role
}

// Succeed promotes the strongest surviving member of an organization to lead it.
// An ambitious successor with weak support is exactly the kind of arrangement
// that does not hold.
func (w *World) Succeed(faction string) {
	f := w.faction(faction)
	if f == nil {
		return
	}
	members := w.Members(faction)
	if len(members) == 0 {
		// Nobody is left to hold it together. The organization is a shell.
		f.Power = max(10, f.Power/2)
		w.Log("Nobody left to lead "+f.Name, "Everyone who could have taken over is dead. What remains of the organization holds together on habit alone.", "politics")
		return
	}
	successor := members[0]
	previous := f.Leader
	f.Leader = successor.Name
	successor.Rank = RankLeader
	successor.Role = "Head of " + f.Name
	// A change at the top is disruptive even when it is orderly.
	f.Power = max(10, f.Power-8)
	w.Log(successor.Name+" takes over "+f.Name,
		fmt.Sprintf("With %s gone, %s now leads %s. The organization is weaker while the change settles.", previous, successor.Name, f.Name),
		"politics")
	// Somebody who thought it should have been them now has a reason of their
	// own, which is how an orderly succession stops being orderly.
	for _, peer := range members[1:] {
		if peer.Dead || peer.Rank < RankLieutenant {
			continue
		}
		w.Resent(peer.ID, successor.ID, 28, "being passed over when "+previous+" died")
		break
	}
}

// casualty picks who dies when violence reaches an organization. The people
// sent to do the work are hit most often, but nobody is exempt, and a leader
// caught in the wrong place dies like anyone else.
func (w *World) casualty(faction string) *NPC {
	members := w.Members(faction)
	if len(members) == 0 {
		return nil
	}
	// Weight toward the bottom of the organization without excluding the top.
	weights := make([]int, len(members))
	total := 0
	for i, m := range members {
		weights[i] = max(1, RankLeader-m.Rank+10)
		total += weights[i]
	}
	roll := int(w.WorldRandom() * float64(total))
	for i, weight := range weights {
		roll -= weight
		if roll < 0 {
			return members[i]
		}
	}
	return members[len(members)-1]
}

// ConsiderInternalMove gives an ambitious deputy in a failing organization a
// reason to take it for themselves. It is the same violence as any other, done
// by people the city already knows, and it can fail.
// ConsiderInternalMove is the older name for what InternalMove now does. Kept
// so that anything still calling it gets the current behaviour rather than the
// version without an aftermath.
func (w *World) ConsiderInternalMove(f *Faction) bool { return w.InternalMove(f) }
