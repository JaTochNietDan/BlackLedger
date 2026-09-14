package billiards

import "fmt"

// Bracket holds actual eight-ball matches. Campaign code owns entry fees and
// eligibility; a client must never submit bracket winners or replacement racks.
type Bracket struct {
	Entrants []string       `json:"entrants"`
	Matches  []BracketMatch `json:"matches"`
	Winner   string         `json:"winner"`
}
type BracketMatch struct {
	Round    int       `json:"round"`
	Table    int       `json:"table"`
	Players  [2]string `json:"players"`
	Rack     *Match    `json:"rack,omitempty"`
	Winner   string    `json:"winner"`
	Next     int       `json:"next"`
	NextSeat int       `json:"next_seat"`
}

// Initial entrant order is seeding. The first entrant remains seat zero in
// successive rounds if they win, allowing the campaign's player-first UI.
func NewBracket(entrants []string) (*Bracket, error) {
	n := len(entrants)
	if n < 2 || n > 8 || n&(n-1) != 0 {
		return nil, fmt.Errorf("a billiards bracket needs two, four or eight entrants")
	}
	seen := map[string]bool{}
	for _, id := range entrants {
		if id == "" || seen[id] {
			return nil, fmt.Errorf("bracket entrants must have distinct identities")
		}
		seen[id] = true
	}
	b := &Bracket{Entrants: append([]string{}, entrants...), Matches: make([]BracketMatch, 0, n-1)}
	for size, round := n/2, 0; size > 0; size, round = size/2, round+1 {
		start := len(b.Matches)
		for i := 0; i < size; i++ {
			cell := BracketMatch{Round: round, Next: -1}
			if size > 1 {
				cell.Next = start + size + i/2
				cell.NextSeat = i % 2
			}
			if round == 0 {
				cell.Table = i + 1
				cell.Players = [2]string{entrants[i*2], entrants[i*2+1]}
				cell.Rack, _ = NewMatch(0)
			}
			b.Matches = append(b.Matches, cell)
		}
	}
	return b, nil
}

// Advance records only winners adjudicated by the embedded physical match (or
// its explicit concession rule). It is idempotent and never plays unseen shots.
func (b *Bracket) Advance() {
	if b.Winner != "" {
		return
	}
	for i := range b.Matches {
		cell := &b.Matches[i]
		if cell.Winner != "" || cell.Rack == nil || cell.Rack.Winner < 0 || cell.Rack.Winner > 1 {
			continue
		}
		winner := cell.Players[cell.Rack.Winner]
		if winner == "" {
			continue
		}
		cell.Winner = winner
		if cell.Next < 0 {
			b.Winner = winner
			return
		}
		next := &b.Matches[cell.Next]
		next.Players[cell.NextSeat] = winner
		if next.Players[0] != "" && next.Players[1] != "" && next.Rack == nil {
			next.Table = b.freeTable()
			next.Rack, _ = NewMatch(0)
		}
	}
}

// ActiveFor returns a playable match, never a speculative future pairing or an
// already resolved rack. Advance should be called after each committed shot.
func (b *Bracket) ActiveFor(id string) (index, seat int, ok bool) {
	if id == "" || b.Winner != "" {
		return 0, 0, false
	}
	for i := range b.Matches {
		cell := &b.Matches[i]
		if cell.Rack == nil || cell.Winner != "" || cell.Rack.Winner >= 0 {
			continue
		}
		for seat, player := range cell.Players {
			if player == id {
				return i, seat, true
			}
		}
	}
	return 0, 0, false
}

// Later rounds may become ready while another opening rack is still playing.
// Allocate an actually free table rather than reusing a round-relative number.
func (b *Bracket) freeTable() int {
	occupied := map[int]bool{}
	for _, cell := range b.Matches {
		if cell.Rack != nil && cell.Rack.Winner < 0 && cell.Winner == "" {
			occupied[cell.Table] = true
		}
	}
	for table := 1; table <= 6; table++ {
		if !occupied[table] {
			return table
		}
	}
	return 0 // Unreachable for a validated bracket of at most eight entrants.
}
