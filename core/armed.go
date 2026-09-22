package core

// Somebody in the room had a gun.
//
// Eleven rules in this game end a life and two of them ever fired: ninety-eight
// per cent of every death this project has measured is somebody coming for the
// player at an address they were standing at. Everything the player does to
// other people — a till, a pocket, a charge, a move on a holding — is a health
// tax and never a risk. A failed robbery costs ten to thirty health and the
// policies rest below seventy, so the arithmetic simply never reaches zero.
//
// That is not what this city is supposed to be. A man who goes into a room to
// take the till is betting on nobody in there being armed, and sometimes
// somebody is. So a job the player does themselves, and only when it goes
// wrong, can be the end of them.
//
// It is deliberately the same shape as an attempt on somebody's life
// (`itWentWrong`), which has always worked this way and is the one piece of
// violence in this game with a real price on it: a flat roll, taken once, that
// ends the campaign rather than adding damage. Damage the player can rest off.

const (
	// TillGun is the chance a business that was robbed and fought back had
	// somebody armed behind the counter. Lower than a mugging because a shop
	// is a shop; what raises it is who owns it.
	TillGun = .05
	// PocketGun is the same for somebody stopped in the street. Higher,
	// because a man with something worth taking off him in this city is a man
	// who has thought about that.
	PocketGun = .06
	// GunByPower is what a family behind the room or the mark adds, per point
	// of their strength.
	GunByPower = .0012
	// GunByArmour is what each step of what the player is wearing takes off it.
	// The other half of why anybody buys the stuff: it has only ever reduced a
	// beating, which is the one kind of harm that did not matter.
	GunByArmour = .022
)

// armedResistance reports whether the person on the other side of a job that
// has already gone wrong was carrying, and the player is not walking away.
//
// Only for work the player did themselves. Sending somebody has its own price
// and it is already written: they are hurt, they are taken, or they talk.
func (w *World) armedResistance(base float64, power int, hand Hand) bool {
	if hand.Crew {
		return false
	}
	odds := base + float64(power)*GunByPower - float64(w.Player.Armour)*GunByArmour
	if odds <= 0 {
		return false
	}
	return w.Random() < odds
}
