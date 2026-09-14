package core

import "testing"

func TestShowdownSeatDescriptionUsesSharedBoard(t *testing.T) {
	w := New(61)
	w.Game = &CardGame{Place: "bar", Mine: hand("Th", "5h"), Board: hand("6h", "Js", "Jd", "4h", "2c"), Seats: []Seat{{Who: "mara", Cards: hand("2h", "Jc")}}}
	hidden := w.CardsDescription()["seats"].([]map[string]any)[0]
	if _, ok := hidden["hand"]; ok {
		t.Fatal("opponent rank leaked before showdown")
	}
	if _, ok := hidden["cards"]; ok {
		t.Fatal("opponent cards leaked before showdown")
	}
	w.Game.Done = true
	shown := w.CardsDescription()["seats"].([]map[string]any)[0]
	if shown["hand"] != "a full house" {
		t.Fatalf("shared-board full house described as %v", shown["hand"])
	}
}
