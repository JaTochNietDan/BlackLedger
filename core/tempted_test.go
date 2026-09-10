package core

import "testing"

// A rival with an empty counter and money in the bank should come for your
// people, and what you pay them should be the reason they stay. Until now
// somebody only ever left over a grudge, so a well-run business was safe from
// the city entirely and the wage was a lever with nothing pulling against it.

func tempted(t *testing.T) (*World, string, string) {
	t.Helper()
	w := New(61)
	w.Event, w.District = nil, 9
	w.Player.Cash, w.Player.Health = 20000, 100
	own(w, "laundry")
	w.EmptyChairs()
	rival := w.Factions[0].ID
	if rival == w.PlayerOrganizationID() {
		rival = w.Factions[1].ID
	}
	// A rival with a counter to fill and the money to fill it.
	w.Properties["butcher"].Owner = rival
	w.Properties["butcher"].Staff, w.Properties["butcher"].Hands = 0, nil
	if f := w.faction(rival); f != nil {
		f.Cash = 40000
	}
	return w, "laundry", rival
}

func stayed(w *World, id string, days int) int {
	for day := 0; day < days; day++ {
		w.Event = nil
		w.Advance(1440)
		w.Event = nil
	}
	return len(w.Properties[id].Hands)
}

func TestARivalComesForYourPeople(t *testing.T) {
	t.Parallel()
	w, id, _ := tempted(t)
	trade, _ := TradeOf(id)
	// Paid the least anybody stands there for, and nobody carrying anything
	// against the player: the only reason to go is the money.
	least, _ := WageBounds(trade.Wage)
	if err := w.SetWage(id, least); err != nil {
		t.Fatal(err)
	}
	filled := len(w.Properties[id].Hands)
	if left := stayed(w, id, 60); left == filled {
		t.Fatalf("a rival with an empty counter and $40,000 never took one of the %d you pay $%d a day", filled, least)
	}
}

func TestPayingWellKeepsThem(t *testing.T) {
	t.Parallel()
	w, id, _ := tempted(t)
	trade, _ := TradeOf(id)
	_, most := WageBounds(trade.Wage)
	if err := w.SetWage(id, most); err != nil {
		t.Fatal(err)
	}
	filled := len(w.Properties[id].Hands)
	if left := stayed(w, id, 60); left < filled {
		t.Fatalf("people paid $%d a day walked out for a rival anyway: %d of %d", most, left, filled)
	}
}

// Where they land is a different question from whether they go. A family that
// can pay takes them; a city where nobody can simply loses them, because being
// underpaid is a reason to stop turning up whether or not anybody is hiring.
func TestWhoTakesThemIsADifferentQuestionFromWhetherTheyGo(t *testing.T) {
	t.Parallel()
	w, id, _ := tempted(t)
	trade, _ := TradeOf(id)
	least, _ := WageBounds(trade.Wage)
	if err := w.SetWage(id, least); err != nil {
		t.Fatal(err)
	}
	gone := []string{}
	was := append([]string{}, w.Properties[id].Hands...)
	stayed(w, id, 60)
	for _, who := range was {
		if w.EmployerOf(who) != id {
			gone = append(gone, who)
		}
	}
	if len(gone) == 0 {
		t.Fatalf("nobody paid $%d a day left in two months", least)
	}
	somewhere := 0
	for _, who := range gone {
		if at := w.EmployerOf(who); at != "" {
			somewhere++
		}
	}
	t.Logf("of %d who left a counter paying $%d, %d went straight to somebody else's", len(gone), least, somewhere)
	// Nobody who left is on two sets of books, which is the thing that would
	// actually be wrong.
	for _, who := range gone {
		at := w.EmployerOf(who)
		if at == id {
			t.Fatalf("%s left and is still on your books", w.NPC(who).Name)
		}
	}
}

// The rate is the safe wage. Paying exactly what the work is worth in this city
// is not a reason for anybody to go anywhere, and the first version of this had
// it worth half of the floor's temptation — a business paying the going rate
// bled people for no reason anybody could name.
func TestPayingTheRateTemptsNobody(t *testing.T) {
	t.Parallel()
	w, id, _ := tempted(t)
	trade, _ := TradeOf(id)
	if w.WageAt(id) != trade.Wage {
		t.Fatalf("a business nobody has touched pays $%d against a rate of $%d", w.WageAt(id), trade.Wage)
	}
	if odds := w.tempted(id); odds != 0 {
		t.Fatalf("paying the going rate tempts people away at %.3f a day", odds)
	}
	filled := len(w.Properties[id].Hands)
	if left := stayed(w, id, 60); left < filled {
		t.Fatalf("people paid the going rate left anyway: %d of %d", left, filled)
	}
	// And a penny under it is worth something, or the lever has no travel.
	if err := w.SetWage(id, trade.Wage-1); err != nil {
		t.Fatal(err)
	}
	if odds := w.tempted(id); odds <= 0 {
		t.Fatal("a wage under the rate tempts nobody at all")
	}
}
