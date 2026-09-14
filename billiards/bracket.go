package billiards

import "fmt"

// Bracket holds actual eight-ball matches. Campaign code owns entry fees and
// eligibility; a client must never submit bracket winners or replacement racks.
type Bracket struct {
	Withdrawn map[string]bool `json:"withdrawn"`
	Finished  bool            `json:"finished"`
	Entrants  []string        `json:"entrants"`
	Matches   []BracketMatch  `json:"matches"`
	Winner    string          `json:"winner"`
}
type BracketMatch struct {
	Resolved bool      `json:"resolved"`
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

// Advance records physical match winners, explicit concessions and withdrawal
// walkovers. It is idempotent and never plays unseen shots.
func (b *Bracket) Advance() {
	if b.Finished || b.Winner != "" {
		b.Finished = true
		return
	}
	for i := range b.Matches {
		cell := &b.Matches[i]
		if cell.Resolved {
			continue
		}
		if cell.Winner != "" {
			cell.Resolved = true
			continue
		}
		if cell.Round > 0 {
			ready := 0
			for j := 0; j < i; j++ {
				source := &b.Matches[j]
				if source.Next == i && source.Resolved {
					cell.Players[source.NextSeat] = source.Winner
					ready++
				}
			}
			if ready != 2 {
				continue
			}
		}
		present := [2]bool{}
		for seat, id := range cell.Players {
			present[seat] = id != "" && !b.Withdrawn[id]
		}
		winner := ""
		switch {
		case !present[0] && !present[1]:
			// A resolved empty branch supplies a bye, not a fabricated champion.
		case !present[0] || !present[1]:
			seat := 0
			if present[1] {
				seat = 1
			}
			winner = cell.Players[seat]
			if cell.Rack != nil && cell.Rack.Winner < 0 {
				_ = cell.Rack.Concede(1 - seat)
			}
		default:
			if cell.Rack == nil {
				cell.Table = b.freeTable()
				cell.Rack, _ = NewMatch(0)
			}
			if cell.Rack.Winner < 0 {
				continue
			}
			winner = cell.Players[cell.Rack.Winner]
		}
		cell.Winner = winner
		cell.Resolved = true
		if cell.Next < 0 {
			b.Winner = winner
			b.Finished = true
			return
		}
	}
}

// Withdraw removes eligibility, including while waiting for another table. It
// cannot rewrite a finished event. Campaign code supplies the actual cause.
func (b *Bracket) Withdraw(id string) error { return b.WithdrawMany([]string{id}) }

// Simultaneous departures must be marked together before advancing, otherwise
// the first withdrawal could prematurely award the event to another casualty.
func (b *Bracket) WithdrawMany(ids []string) error {
	if b.Finished || b.Winner != "" {
		return fmt.Errorf("the tournament is already finished")
	}
	entrants := map[string]bool{}
	for _, id := range b.Entrants {
		entrants[id] = true
	}
	for _, id := range ids {
		if !entrants[id] {
			return fmt.Errorf("that person did not enter this tournament")
		}
	}
	if b.Withdrawn == nil {
		b.Withdrawn = map[string]bool{}
	}
	for _, id := range ids {
		b.Withdrawn[id] = true
	}
	b.Advance()
	return nil
}

// ActiveFor returns a playable match, never a speculative future pairing or an
// already resolved rack. Advance should be called after each committed shot.
func (b *Bracket) ActiveFor(id string) (index, seat int, ok bool) {
	if id == "" || b.Finished || b.Winner != "" || b.Withdrawn[id] {
		return 0, 0, false
	}
	for i := range b.Matches {
		cell := &b.Matches[i]
		if cell.Rack == nil || cell.Resolved || cell.Winner != "" || cell.Rack.Winner >= 0 {
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
		if cell.Rack != nil && cell.Rack.Winner < 0 && !cell.Resolved && cell.Winner == "" {
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
