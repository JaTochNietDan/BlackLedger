package core

import (
	"strings"
	"testing"
)

// "I don't see a lot of cursing from characters in this game, we should increase
// that since it's with the mafia style. Characters should be able to make
// threats to you too, I have not seen that yet."
//
// And, from the same inbox: "I see 'bad blood' red box at the top of the page
// but it never goes away and it's annoying." The two are the same system. What
// somebody holds against you should be said to your face by the person holding
// it, and should stop being a banner that sits there for a month.

func sore(t *testing.T) (*World, *NPC) {
	t.Helper()
	w := New(29)
	w.Event, w.District = nil, 9
	w.Player.Health, w.Player.Cash = 100, 500
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) {
			continue
		}
		n.Trust = 30 // known well enough to be spoken to
		n.Location = "bar"
		w.Player.Location = "bar"
		return w, n
	}
	t.Fatal("nobody in this city")
	return nil, nil
}

func TestSomebodyWithSomethingAgainstYouSaysSo(t *testing.T) {
	w, n := sore(t)
	if said := w.Threat(n); said != "" {
		t.Fatalf("somebody with nothing against you threatened you: %q", said)
	}
	w.Aggrieve(n.ID, 12, "the money you took off them")
	n.Sore = 12
	mild := w.Threat(n)
	if mild == "" {
		t.Fatal("somebody sore at you said nothing at all")
	}
	// Three weights, three different things said. Two would not tell them
	// apart: with the worst branch deleted the top of the scale falls through
	// to the middle one, which still reads as an escalation and is not.
	n.Sore = Plain
	warned := w.Threat(n)
	n.Sore = GrudgeActs
	worse := w.Threat(n)
	said := map[string]bool{mild: true, warned: true, worse: true}
	if len(said) != 3 {
		t.Fatalf("the scale has only %d voices on it: %q / %q / %q", len(said), mild, warned, worse)
	}
	// It reaches the room, on the person carrying it.
	found := false
	for _, p := range w.PeopleHere("bar") {
		if p.ID == n.ID && p.Says != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("nothing they said reached the room")
	}
}

func TestTheCityDoesSwear(t *testing.T) {
	// Not one line in five hundred, which is what "increase the cursing" means.
	seen, salty := 0, 0
	words := []string{"hell", "damn", "christ", "bastard", "son of a", "goddamn"}
	// Every person in a few cities, not the first person in many: the first
	// name a seed produces is often the same name, which would measure one
	// man's manners three hundred times.
	for seed := uint32(1); seed <= 6; seed++ {
		w := New(seed * 2654435761)
		for i := range w.NPCs {
			n := &w.NPCs[i]
			if n.Dead || IsOfficial(n.ID) {
				continue
			}
			n.Trust = 30
			w.Aggrieve(n.ID, GrudgeActs+10, "what you did")
			said := strings.ToLower(w.Threat(n))
			if said == "" {
				continue
			}
			seen++
			for _, word := range words {
				if strings.Contains(said, word) {
					salty++
					break
				}
			}
		}
	}
	t.Logf("of %d lines from people who want you dead, %d had heat in them", seen, salty)
	if seen == 0 {
		t.Fatal("nobody said anything")
	}
	if salty*3 < seen {
		t.Fatalf("only %d of %d lines from people who want you dead had any heat in them", salty, seen)
	}
}

func TestBadBloodIsNewsAndStopsBeingNews(t *testing.T) {
	w := New(29)
	w.Event, w.District = nil, 9
	w.Player.Contacts = 5 // reach enough to hear gossip at all
	var a, b *NPC
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if n.Dead || IsOfficial(n.ID) {
			continue
		}
		if a == nil {
			a = n
			continue
		}
		b = n
		break
	}
	w.Resent(a.ID, b.ID, 40, "a car that went missing")
	if len(w.GrudgeSummary()) == 0 {
		t.Fatal("fresh bad blood was not worth hearing about")
	}
	// A week later nobody is still talking about it. Asked about this quarrel
	// rather than about the summary being empty: the city now falls out over
	// its own card games while the week passes, so an empty page would mean
	// nothing had happened in seven days rather than that this had stopped
	// being news.
	w.Advance(GrudgeNews + 1440)
	for _, line := range w.GrudgeSummary() {
		if line["holder"] == a.Name && line["against"] == b.Name {
			t.Fatalf("the same bad blood is still at the top of the page: %+v", line)
		}
	}
	// And it is still there in the world, because it did not stop being true.
	held := false
	for _, g := range w.Grudges {
		if g.Holder == a.ID && g.Against == b.ID && g.Weight > 0 {
			held = true
		}
	}
	if !held {
		t.Fatal("the grudge itself was thrown away with the headline")
	}
}
