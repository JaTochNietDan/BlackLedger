package core

import "testing"

func launderer(t *testing.T) *World {
	t.Helper()
	w := New(73)
	w.Properties["laundry"].Owner = "player:1"
	w.Properties["laundry"].Condition = 100
	w.Player.Location = "laundry"
	w.Player.Cash = 5000
	w.Player.Heat = 30
	return w
}

func TestTheBooksTradeMoneyForAttention(t *testing.T) {
	w := launderer(t)
	heat, cash := w.Player.Heat, w.Player.Cash
	fee := w.LaunderFee("laundry")
	if err := w.Launder("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.Player.Heat >= heat {
		t.Fatal("laundering cleared no attention")
	}
	if w.Player.Cash != cash-fee {
		t.Fatalf("the cut was not taken: cash %d, expected %d", w.Player.Cash, cash-fee)
	}
	if w.Properties["laundry"].Condition >= 100 {
		t.Fatal("the premises took no wear")
	}
	if w.Player.LastLaunder != w.Minute {
		t.Fatal("the books do not remember when they last absorbed a round")
	}
}

func TestTheBooksNeedTimeBetweenRounds(t *testing.T) {
	w := launderer(t)
	if err := w.Launder("laundry"); err != nil {
		t.Fatal(err)
	}
	if w.LaunderReadiness("laundry") == "" {
		t.Fatal("the books absorbed two rounds back to back")
	}
	if err := w.Launder("laundry"); err == nil {
		t.Fatal("the command ignored the cooldown the action reports")
	}
	w.Minute += LaunderCooldown
	w.Player.Heat = 20
	if reason := w.LaunderReadiness("laundry"); reason != "" {
		t.Fatal("a day later the books still refused:", reason)
	}
}

func TestOnlyYourOwnCashBusinessLaunders(t *testing.T) {
	w := launderer(t)
	// Somebody else's premises are not your books.
	w.Properties["laundry"].Owner = "bellandi"
	if w.LaunderReadiness("laundry") == "" {
		t.Fatal("the player used a rival's business to clean their money")
	}
	if err := w.Launder("laundry"); err == nil {
		t.Fatal("the command ignored ownership")
	}
	// This used to say "a casino is not a laundry" and refuse one. That was
	// the room type talking: laundering asked what KIND OF ROOM a place was
	// rather than what trade was run in it, and a casino is the best cash
	// front in any city — it is most of the reason anybody owns one. A casino
	// launders now, and less well than a laundry, which is the real rule.
	//
	// What has to stay refused is a place that runs no business at all.
	w.Properties["casino"].Owner = "player:1"
	if w.LaunderReadiness("casino") != "" {
		t.Fatalf("a casino of the player's own was refused its own books: %q", w.LaunderReadiness("casino"))
	}
	// And it absorbs more than a laundry does, because more loose cash crosses
	// a casino floor in a night than a laundry sees in a week. That is what
	// the eight hundred and fifty dollars buys.
	if w.launderCapacity("casino") <= w.launderCapacity("laundry") {
		t.Errorf("a casino absorbs %d and a laundry %d", w.launderCapacity("casino"), w.launderCapacity("laundry"))
	}
	w.Properties["room"].Owner = "player:1"
	if w.LaunderReadiness("room") == "" {
		t.Fatal("a rented room with no trade in it was treated as a cash-handling front")
	}
	// Nothing to clean is nothing to do.
	w.Properties["laundry"].Owner = "player:1"
	w.Player.Heat = 0
	if w.LaunderReadiness("laundry") == "" {
		t.Fatal("the books ran with no attention to clear")
	}
}

func TestARundownBusinessExplainsLess(t *testing.T) {
	good, poor := launderer(t), launderer(t)
	poor.Properties["laundry"].Condition = 45
	if poor.launderCapacity("laundry") >= good.launderCapacity("laundry") {
		t.Fatal("condition does not affect how much the books can absorb")
	}
	if poor.LaunderFee("laundry") <= 0 {
		t.Fatal("a run-down business launders for nothing")
	}
	// Below the threshold it cannot be used at all.
	poor.Properties["laundry"].Condition = 20
	if poor.LaunderReadiness("laundry") == "" {
		t.Fatal("premises falling apart still explained the money")
	}
}
