package billiards

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBracketWaitsForRealPairingsAndSurvivesReopen(t *testing.T) {
	b, err := NewBracket([]string{"player", "leo", "mara", "elena", "vance", "ellis", "vittorio", "alma"})
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Matches) != 7 {
		t.Fatal("wrong bracket size")
	}
	for i := 0; i < 4; i++ {
		if b.Matches[i].Rack == nil || b.Matches[i].Table != i+1 {
			t.Fatal("opening tables not allocated")
		}
	}
	if b.Matches[4].Rack != nil {
		t.Fatal("future match started without opponents")
	}
	if err = b.Matches[0].Rack.Concede(1); err != nil {
		t.Fatal(err)
	}
	b.Advance()
	if b.Matches[4].Players[0] != "player" || b.Matches[4].Rack != nil {
		t.Fatal("advanced before other table finished")
	}
	if _, _, ok := b.ActiveFor("player"); ok {
		t.Fatal("waiting entrant has a live rack")
	}
	raw, _ := json.Marshal(b)
	var restored Bracket
	if err = json.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	b = &restored
	for i := 1; i < 4; i++ {
		if err = b.Matches[i].Rack.Concede(1); err != nil {
			t.Fatal(err)
		}
	}
	b.Advance()
	if i, seat, ok := b.ActiveFor("player"); !ok || i != 4 || seat != 0 {
		t.Fatal("player not seated in semifinal", i, seat, ok)
	}
	for i := 4; i < 6; i++ {
		if err = b.Matches[i].Rack.Concede(1); err != nil {
			t.Fatal(err)
		}
	}
	b.Advance()
	if b.Matches[6].Players != [2]string{"player", "vance"} {
		t.Fatal("wrong finalists", b.Matches[6].Players)
	}
	if err = b.Matches[6].Rack.Concede(1); err != nil {
		t.Fatal(err)
	}
	b.Advance()
	if b.Winner != "player" {
		t.Fatal("no champion")
	}
	before, _ := json.Marshal(b)
	for i := 0; i < 5; i++ {
		b.Advance()
	}
	after, _ := json.Marshal(b)
	if string(before) != string(after) {
		t.Fatal("repeated advancement changed completed bracket")
	}
	if _, _, ok := b.ActiveFor("player"); ok {
		t.Fatal("champion still playing")
	}
}
func TestBracketValidatesEntrantsAndOwnsItsSeedList(t *testing.T) {
	for _, ids := range [][]string{nil, {"a"}, {"a", "a"}, {"a", ""}, {"a", "b", "c"}} {
		if _, err := NewBracket(ids); err == nil {
			t.Fatal("accepted invalid bracket", ids)
		}
	}
	ids := []string{"a", "b"}
	b, err := NewBracket(ids)
	if err != nil {
		t.Fatal(err)
	}
	ids[0] = "changed"
	if !reflect.DeepEqual(b.Entrants, []string{"a", "b"}) {
		t.Fatal("caller can change seeds")
	}
	b.Advance()
	if b.Winner != "" {
		t.Fatal("unfinished match declared champion")
	}
}
func TestBracketChampionComesFromPhysicalEightBall(t *testing.T) {
	b, err := NewBracket([]string{"leo", "mara"})
	if err != nil {
		t.Fatal(err)
	}
	for shots := 0; shots < 120 && b.Winner == ""; shots++ {
		m := b.Matches[0].Rack
		turn, err := Opponent(New(), m, m.Turn, .9, uint64(7+shots*193))
		if err != nil {
			t.Fatal(err)
		}
		b.Matches[0].Rack = &turn.Match
		b.Advance()
	}
	if b.Winner == "" {
		t.Fatal("physical match did not finish")
	}
	m := b.Matches[0].Rack
	if b.Winner != b.Matches[0].Players[m.Winner] || m.Shots < 2 {
		t.Fatal("champion disagrees with physical match")
	}
	t.Logf("%s wins after %d physical strokes: %s", b.Winner, m.Shots, m.Last.Reason)
}

func TestBracketLaterRoundNeverTakesAnOccupiedOpeningTable(t *testing.T) {
	b, err := NewBracket([]string{"a", "b", "c", "d", "e", "f", "g", "h"})
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range []int{2, 3} {
		if err = b.Matches[i].Rack.Concede(1); err != nil {
			t.Fatal(err)
		}
	}
	b.Advance()
	if b.Matches[5].Rack == nil {
		t.Fatal("ready semifinal not started")
	}
	used := map[int]bool{}
	for _, cell := range b.Matches {
		if cell.Rack == nil || cell.Rack.Winner >= 0 {
			continue
		}
		if cell.Table < 1 || cell.Table > 6 || used[cell.Table] {
			t.Fatal("active games share a table", cell.Table)
		}
		used[cell.Table] = true
	}
	if b.Matches[5].Table == 1 || b.Matches[5].Table == 2 {
		t.Fatal("semifinal displaced active opening game")
	}
}
