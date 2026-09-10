package core

import "fmt"

// Everything a person owns dies with them. That is the point of this game, and
// nothing here changes it: a new arrival still starts with ninety dollars, a
// rented room and no protection.
//
// What an account abroad changes is what a careful person can leave behind for
// whoever comes next. Money sent out of the city survives its owner, but it
// costs a cut going in, it earns nothing sitting there, and reaching it as a
// stranger costs money a stranger does not have. It is a decision to give up
// something now against something later, taken while you are still alive to
// take it.

const (
	// DepositLot is how much is wired at a time, so the interface can offer a
	// plain action rather than an amount field.
	DepositLot = 500
	// DepositCut is the percentage the arrangement takes on the way out of the
	// city. Banking money is a loss taken deliberately.
	DepositCut = 18
	// AccessCost is what a new arrival pays to establish that the account is
	// theirs: papers, a journey, and somebody vouching. It is deliberately more
	// than a new person is given.
	AccessCost = 400
	// AccessMinutes is the time that takes.
	AccessMinutes = 240
)

// BankReachable reports whether arrangements of this kind can be made here.
func BankReachable(location string) bool { return location == "market" }

// DepositLeast is the smallest sum the arrangement will carry: below it the
// cut is not worth anybody's trouble.
const DepositLeast = 100

// DepositSum is the figure a request actually wires out: what was typed, or the
// lot when nothing was.
func (w *World) DepositSum(amount int) int {
	if amount > 0 {
		return amount
	}
	return DepositLot
}

// DepositReadiness explains why the typed figure cannot be sent out, or
// returns "".
func (w *World) DepositReadiness(amount int) string {
	if !BankReachable(w.Player.Location) {
		return "This is not arranged here"
	}
	amount = w.DepositSum(amount)
	if w.Player.Cash < DepositLeast {
		return fmt.Sprintf("You need $%d to send anything out", DepositLeast)
	}
	if amount < DepositLeast {
		return fmt.Sprintf("Nothing under $%d is carried", DepositLeast)
	}
	if amount > w.Player.Cash {
		return fmt.Sprintf("You have $%d", w.Player.Cash)
	}
	return ""
}

// Deposit sends the typed figure out of the city, less the cut.
func (w *World) Deposit(amount int) error {
	if reason := w.DepositReadiness(amount); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	amount = w.DepositSum(amount)
	if err := w.Pay(amount); err != nil {
		return err
	}
	kept := amount * (100 - DepositCut) / 100
	w.Offshore += kept
	w.Log("Money leaves the city", fmt.Sprintf("$%d sent out, $%d of it arrives. The account holds $%d and answers to nobody here, including you if anything happens.",
		amount, kept, w.Offshore), "business")
	return nil
}

// AccessReadiness explains why the account cannot be reached, or returns "".
func (w *World) AccessReadiness() string {
	if !BankReachable(w.Player.Location) {
		return "This is not arranged here"
	}
	if w.Player.Offshore {
		return "The account already answers to you"
	}
	if w.Offshore == 0 {
		return "There is nothing out there to reach"
	}
	if w.Player.Cash < AccessCost {
		return fmt.Sprintf("You need $%d to establish that it is yours", AccessCost)
	}
	return ""
}

// EstablishAccess is what a new arrival does to prove the account is theirs.
// The money exists whether or not anyone can reach it; this is the reaching.
func (w *World) EstablishAccess() error {
	if reason := w.AccessReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(AccessCost); err != nil {
		return err
	}
	w.Player.Offshore = true
	w.Log("The account answers to you", fmt.Sprintf("$%d in papers, a journey and somebody willing to vouch. There is $%d out there and you can now reach it.", AccessCost, w.Offshore), "personal")
	return nil
}

// WithdrawLeast is the smallest sum worth a journey to fetch.
const WithdrawLeast = 25

// WithdrawSum is the figure a request actually brings home. Nothing typed means
// all of it, which is what the plain action has always done.
func (w *World) WithdrawSum(amount int) int {
	if amount > 0 {
		return amount
	}
	return w.Offshore
}

// WithdrawReadiness explains why the typed figure cannot be drawn, or returns "".
func (w *World) WithdrawReadiness(amount int) string {
	if !BankReachable(w.Player.Location) {
		return "This is not arranged here"
	}
	if !w.Player.Offshore {
		return "The account does not answer to you yet"
	}
	if w.Offshore == 0 {
		return "The account is empty"
	}
	amount = w.WithdrawSum(amount)
	if amount > w.Offshore {
		return fmt.Sprintf("There is $%d out there", w.Offshore)
	}
	if amount < min(WithdrawLeast, w.Offshore) {
		return fmt.Sprintf("Bring $%d or more", min(WithdrawLeast, w.Offshore))
	}
	return ""
}

// Withdraw brings the typed figure back into the city. Money in hand is money
// that can be taken, which is the risk of bringing it home.
func (w *World) Withdraw(amount int) error {
	if reason := w.WithdrawReadiness(amount); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	amount = w.WithdrawSum(amount)
	w.Offshore -= amount
	w.Earn(amount)
	// "$0 is still out there" is a figure that says nothing while looking like
	// one, the same fault the panel's own wording had.
	rest := "The account is empty now."
	if w.Offshore > 0 {
		rest = fmt.Sprintf("$%d is still out there.", w.Offshore)
	}
	w.Log("It comes home", fmt.Sprintf("$%d back in the city and in your hands, where anybody can take it from you. %s", amount, rest), "business")
	return nil
}
