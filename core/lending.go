package core

import "fmt"

// Money in this city only ever went one way. The families lean on the player
// for a share, the player pays tribute, the police fine, the street robs — and
// in all of it there was never a single person who owed the player anything.
// Cash sat in a pocket doing nothing between purchases, and the oldest business
// this trade has did not exist.
//
// A loan is the whole game in miniature. Lending is easy and costs nothing.
// Collecting is where it becomes a decision: somebody who cannot pay is not a
// number, they are a man standing in front of you, and what you do about it is
// the difference between a lender and a moneylender.

const (
	// LoanTermDays is how long somebody has before it is due.
	LoanTermDays = 7
	// LoanRate is what is owed on top, over the term.
	LoanRate = .35
	// LoanFloor is the smallest sum anybody bothers borrowing. Somebody whose
	// position cannot service even this is where the bad debts come from.
	LoanFloor = 120
	// PayReach is how many days of what somebody carries they could put their
	// hands on if the alternative was unpleasant.
	PayReach = 9
	// LoanShare is the most of the player's cash that can be out at once, as a
	// fraction. Lending everything you have is not a strategy, it is a way of
	// having nothing.
	LoanShare = .4
	// LendMinutes is the conversation.
	LendMinutes = 30
	// LeanMinutes is the other conversation.
	LeanMinutes = 60
	// ExtendRate is the extra owed for another week, which is where this trade
	// actually makes its money.
	ExtendRate = .25
	// LeanHeat is the attention a collection draws. Somebody always sees it.
	LeanHeat = 5
	// LeanRespect is what the street gives a man who collects what he is owed,
	// and ForgiveRespect is what it takes off one who does not.
	LeanRespect, ForgiveRespect = 5, 6
	// SlowRespect is what standing costs each day a bad debt is left standing.
	// A man who is owed money and does nothing about it is a man everybody
	// borrows from.
	SlowRespect = 1
	// LendOffers is how many new borrowers are put in front of the player in
	// one room at once. The people who already owe you are never hidden.
	LendOffers = 4
	// LendStanding is the respect it takes before anybody in this city would
	// take money from you rather than laugh.
	LendStanding = 12
)

// Loan is money out with somebody's name on it.
type Loan struct {
	ID     string `json:"id"`
	Debtor string `json:"debtor"`
	Life   int    `json:"life"`
	// Principal is what was handed over; Owed is what comes back.
	Principal int `json:"principal"`
	Owed      int `json:"owed"`
	// Due is when, and Missed is how many times it has not been.
	Due    int `json:"due"`
	Missed int `json:"missed"`
	// Since is when it was made, for anything that wants to know how long
	// somebody has been carrying it.
	Since int `json:"since"`
}

// Book is every loan of this life still outstanding.
func (w *World) Book() []Loan {
	out := []Loan{}
	for _, l := range w.Loans {
		if l.Life == w.Life {
			out = append(out, l)
		}
	}
	return out
}

// LoanTo is what somebody owes, or nil.
func (w *World) LoanTo(id string) *Loan {
	for i := range w.Loans {
		if l := &w.Loans[i]; l.Life == w.Life && l.Debtor == id {
			return l
		}
	}
	return nil
}

// OutOnLoan is everything the player has out at once.
func (w *World) OutOnLoan() int {
	total := 0
	for _, l := range w.Book() {
		total += l.Principal
	}
	return total
}

// Repayable is what somebody could find over a week if they had to: not what is
// in their pocket tonight, but what their position in this city is worth. A
// lieutenant of a rich family can find a great deal more than a docker, and the
// whole skill of this trade is knowing the difference.
func (w *World) Repayable(n *NPC) int { return w.Pockets(n) * PayReach }

// LoanSize is what this person could be lent. Never more than they could
// plausibly bring back with the interest on it — a lender who hands somebody
// more than they can service is not running a business, he is buying a problem
// — and never more than the player can afford to have out at once.
func (w *World) LoanSize(n *NPC) int {
	if n == nil {
		return 0
	}
	serviceable := int(float64(w.Repayable(n))/(1+LoanRate)) * 3 / 4
	want := max(LoanFloor, serviceable)
	ceiling := int(float64(w.Player.Cash+w.OutOnLoan())*LoanShare) - w.OutOnLoan()
	return min(want, max(0, ceiling))
}

// Overdue reports whether a loan has passed its day without being paid.
func (w *World) Overdue(l *Loan) bool { return l != nil && w.Minute >= l.Due }

// BadDebt is everybody who has missed at least once, which is who the player
// has to make a decision about.
func (w *World) BadDebt() []*Loan {
	out := []*Loan{}
	for i := range w.Loans {
		if l := &w.Loans[i]; l.Life == w.Life && l.Missed > 0 {
			out = append(out, l)
		}
	}
	return out
}

// LendReadiness explains why nobody would take money from this person, or
// returns "".
func (w *World) LendReadiness(id string) string {
	n := w.NPC(id)
	if n == nil || n.Dead {
		return "There is nobody here by that name"
	}
	if n.Location != w.Player.Location {
		return "That conversation happens in person"
	}
	if IsOfficial(n.ID) || w.isRoleHolder(n) {
		return "Nobody with a title borrows from somebody without one"
	}
	if w.LoanTo(id) != nil {
		return n.Name + " already owes you"
	}
	if w.Player.Respect < LendStanding {
		return fmt.Sprintf("Nobody borrows from a stranger. It takes %d respect", LendStanding)
	}
	if w.Inside(n) {
		return "The police have him"
	}
	if size := w.LoanSize(n); size < LoanFloor {
		return fmt.Sprintf("You have $%d out already, which is as much as you can afford to be owed", w.OutOnLoan())
	}
	return ""
}

// Lend puts money on the street.
func (w *World) Lend(id string) error {
	if reason := w.LendReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	n := w.NPC(id)
	size := w.LoanSize(n)
	if err := w.Pay(size); err != nil {
		return err
	}
	owed := size + int(float64(size)*LoanRate)
	w.Loans = append(w.Loans, Loan{
		ID: ID(), Debtor: id, Life: w.Life,
		Principal: size, Owed: owed,
		Due: w.Minute + LoanTermDays*1440, Since: w.Minute,
	})
	w.MeetPerson(id)
	w.Log(fmt.Sprintf("$%d to %s", size, n.Name), fmt.Sprintf("$%d back inside %d days. %s knows what the alternative is, which is the only security this trade has.", owed, LoanTermDays, n.Name), "business")
	return nil
}

// canPay is capacity and then willingness, in that order. Somebody whose
// position has got worse since they borrowed genuinely cannot find it, and is
// not lying about it. Somebody who can find it still might not: nerve is how
// readily a person acts on what they think they can get away with, and a lender
// with nothing to lean on is somebody a certain kind of man simply stops
// answering.
func (w *World) canPay(n *NPC, owed int) bool {
	if w.Repayable(n) < owed {
		return false
	}
	return w.WorldRandom() >= w.Nerve(n)/3
}

// LoanDay settles whatever is due. A debtor who can pay does; one who cannot is
// a decision the player has to make rather than a number that resolves itself.
func (w *World) LoanDay() {
	kept := w.Loans[:0]
	slow := 0
	for _, l := range w.Loans {
		if l.Life != w.Life {
			kept = append(kept, l)
			continue
		}
		n := w.NPC(l.Debtor)
		if n == nil || n.Dead {
			w.Log("Nothing to collect", fmt.Sprintf("$%d went out and whoever was carrying it is not carrying anything now.", l.Principal), "danger")
			continue
		}
		if w.Minute < l.Due {
			kept = append(kept, l)
			continue
		}
		// The day it is due.
		if l.Missed == 0 && w.canPay(n, l.Owed) {
			w.Earn(l.Owed)
			n.Trust = min(100, n.Trust+10)
			w.Log(n.Name+" pays", fmt.Sprintf("$%d back on $%d. That is the business working exactly as it is supposed to, which is rarer than it sounds.", l.Owed, l.Principal), "business")
			continue
		}
		if l.Missed == 0 {
			l.Missed = 1
			w.Log(n.Name+" does not have it", fmt.Sprintf("$%d due and nothing to put against it. What happens next is yours to decide, and everybody on this street is going to hear which way you went.", l.Owed), "danger")
		} else {
			slow++
		}
		kept = append(kept, l)
	}
	w.Loans = kept
	// Money owed and not collected is a standing everybody can see.
	if slow > 0 {
		w.Player.Respect = max(0, w.Player.Respect-slow*SlowRespect)
	}
}

// LeanReadiness explains why nobody can be leaned on, or returns "".
func (w *World) LeanReadiness(id string) string {
	l := w.LoanTo(id)
	n := w.NPC(id)
	if l == nil || n == nil || n.Dead {
		return "Nobody here owes you anything"
	}
	if l.Missed == 0 {
		return "It is not due yet"
	}
	if n.Location != w.Player.Location {
		return "That conversation happens in person"
	}
	if w.Inside(n) {
		return "The police have him"
	}
	return ""
}

// LeanOdds is whether it works: what the player brings against what somebody
// with nothing left to lose brings back.
func (w *World) LeanOdds(n *NPC) float64 {
	mine := w.Presence() + w.Player.Weapon*20
	odds := float64(mine) / float64(mine+w.Poise(n)+20)
	if odds < .15 {
		odds = .15
	}
	if odds > .92 {
		odds = .92
	}
	return odds
}

// Lean is the collection. It works more often than not and it costs something
// every time: attention, and a man who will remember it for as long as he
// lives.
func (w *World) Lean(id string) error {
	if reason := w.LeanReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	l, n := w.LoanTo(id), w.NPC(id)
	owed := l.Owed
	w.Player.Heat = min(100, w.Player.Heat+LeanHeat)
	w.Aggrieve(n.ID, 45, "what was done to them over money")
	place, _ := PlaceByID(n.Location)

	if w.Random() >= w.LeanOdds(n) {
		// He has nothing, or he has friends, and either way you are leaving
		// without it.
		l.Missed++
		w.Player.Health = max(1, w.Player.Health-w.Absorb(12))
		w.Log("It did not go your way at "+place.Name, fmt.Sprintf("%s had nothing and was not frightened enough to find any. You leave with a mark on you and $%d still out.", n.Name, owed), "danger")
		return nil
	}

	// Somebody who is squeezed hard enough finds most of it. Never all of it:
	// a man who had it would have paid.
	got := owed * (55 + int(w.Random()*40)) / 100
	w.Earn(got)
	w.Player.Respect += LeanRespect
	n.Trust = 0
	w.dropLoan(l.ID)
	w.Log("Collected at "+place.Name, fmt.Sprintf("$%d of $%d out of %s, and the rest written off because there was no more of it. The street heard about this before you got home.", got, owed, n.Name), "business")
	w.Report("robbery", "ASSAULT REPORTED AT "+upper(place.Name),
		w.unattributed(place.Name, fmt.Sprintf("A man was assaulted near %s. Police say the victim has declined to make a complaint.", place.Name)))
	return nil
}

// ExtendReadiness explains why a debt cannot be rolled over, or returns "".
func (w *World) ExtendReadiness(id string) string {
	l := w.LoanTo(id)
	if l == nil {
		return "Nobody here owes you anything"
	}
	if l.Missed == 0 {
		return "It is not due yet"
	}
	if n := w.NPC(id); n == nil || n.Dead || n.Location != w.Player.Location {
		return "That conversation happens in person"
	}
	return ""
}

// Extend is another week at a worse rate. It is how this trade actually makes
// money, and how a small debt becomes one nobody can ever pay.
func (w *World) Extend(id string) error {
	if reason := w.ExtendReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	l, n := w.LoanTo(id), w.NPC(id)
	before := l.Owed
	l.Owed += int(float64(l.Owed) * ExtendRate)
	l.Due = w.Minute + LoanTermDays*1440
	l.Missed = 0
	n.Trust = min(100, n.Trust+5)
	w.Log("Another week for "+n.Name, fmt.Sprintf("$%d becomes $%d. They are grateful today and further from paying it than they were this morning.", before, l.Owed), "business")
	return nil
}

// ForgiveReadiness explains why a debt cannot be written off, or returns "".
func (w *World) ForgiveReadiness(id string) string {
	l := w.LoanTo(id)
	if l == nil {
		return "Nobody here owes you anything"
	}
	// The three answers are three answers to the same moment. Before it is due
	// there is nothing to be generous about; the money is simply out.
	if l.Missed == 0 {
		return "It is not due yet"
	}
	if n := w.NPC(id); n == nil || n.Dead || n.Location != w.Player.Location {
		return "That conversation happens in person"
	}
	return ""
}

// Forgive writes it off. It costs the money and it costs standing, and it buys
// the only thing in this city money cannot: somebody who owes you nothing and
// would do anything for you.
func (w *World) Forgive(id string) error {
	if reason := w.ForgiveReadiness(id); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	l, n := w.LoanTo(id), w.NPC(id)
	owed := l.Owed
	w.dropLoan(l.ID)
	w.Player.Respect = max(0, w.Player.Respect-ForgiveRespect)
	n.Trust = min(100, n.Trust+40)
	n.Sore, n.SoreAt = 0, ""
	w.Log(n.Name+" owes you nothing", fmt.Sprintf("$%d off the books and a man who knows exactly what that was worth. It is not how this business is done, and everybody will hear that too.", owed), "personal")
	return nil
}

// dropLoan removes a settled debt from the book.
func (w *World) dropLoan(id string) {
	kept := w.Loans[:0]
	for _, l := range w.Loans {
		if l.ID != id {
			kept = append(kept, l)
		}
	}
	w.Loans = kept
}

// LoanDescription is the book, for the interface.
func (w *World) LoanDescription() []map[string]any {
	out := []map[string]any{}
	for _, l := range w.Book() {
		n := w.NPC(l.Debtor)
		if n == nil {
			continue
		}
		place, _ := PlaceByID(n.Location)
		out = append(out, map[string]any{
			"id": l.Debtor, "name": n.Name, "where": place.Name,
			"principal": l.Principal, "owed": l.Owed,
			"days": (l.Due - w.Minute + 1439) / 1440, "missed": l.Missed,
		})
	}
	return out
}
