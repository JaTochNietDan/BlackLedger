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

// DepositReadiness explains why money cannot be sent out, or returns "".
func (w *World) DepositReadiness() string {
	if !BankReachable(w.Player.Location) {
		return "This is not arranged here"
	}
	if w.Player.Cash < DepositLot {
		return fmt.Sprintf("You need $%d to send out at once", DepositLot)
	}
	return ""
}

// Deposit sends a lot out of the city, less the cut.
func (w *World) Deposit() error {
	if reason := w.DepositReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	if err := w.Pay(DepositLot); err != nil {
		return err
	}
	kept := DepositLot * (100 - DepositCut) / 100
	w.Offshore += kept
	w.Log("Money leaves the city", fmt.Sprintf("$%d sent out, $%d of it arrives. The account holds $%d and answers to nobody here, including you if anything happens.",
		DepositLot, kept, w.Offshore), "business")
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

// WithdrawReadiness explains why nothing can be drawn, or returns "".
func (w *World) WithdrawReadiness() string {
	if !BankReachable(w.Player.Location) {
		return "This is not arranged here"
	}
	if !w.Player.Offshore {
		return "The account does not answer to you yet"
	}
	if w.Offshore == 0 {
		return "The account is empty"
	}
	return ""
}

// Withdraw brings everything back into the city at once. Money in hand is money
// that can be taken, which is the risk of bringing it home.
func (w *World) Withdraw() error {
	if reason := w.WithdrawReadiness(); reason != "" {
		return fmt.Errorf("%s", reason)
	}
	amount := w.Offshore
	w.Offshore = 0
	w.Earn(amount)
	w.Log("It comes home", fmt.Sprintf("$%d back in the city and in your hands, where anybody can take it from you.", amount), "business")
	return nil
}
