package core

import "fmt"

// The Mercer Exchange, and knowing what a price means.
//
// A trader in this city can see two numbers: what a good costs on the floor
// they are standing on, and what the other floor pays. That is enough to know
// which way to walk and nothing else. Whether $52 for a crate of moonshine is
// cheap, dear or ordinary — whether this is the week to buy or the week to sit
// on your hands — is not on any card, because the price a thing is normally
// worth is in the core and has never been shown to anybody.
//
// That is what a floor knows. Prices drift a quarter of the way back toward
// what the thing is worth every time the market moves, and while two families
// are shooting at each other scarcity holds them up. None of that is a secret
// to the people who stand on the exchange all day; it is only a secret to
// somebody reading two numbers off a card.
//
// So the Mercer Exchange reaches past its own income by telling you what the
// numbers mean — everywhere, not only in the room. A floor of your own is a
// floor whose clerks take your calls, and the whole point of holding the sell
// side is that it informs the buy side across town.

const (
	// ExchangeHigh and ExchangeLow are how far from what a thing is worth the
	// price has to be, in percent, before the floor will call it anything.
	ExchangeHigh = 115
	ExchangeLow  = 85
)

// TheExchange is the floor the player holds, and nothing if they hold none.
func (w *World) TheExchange() string {
	for _, l := range Locations {
		if l.Kind != "exchange" || !w.Own(l.ID) {
			continue
		}
		prop := w.Properties[l.ID]
		trade, runs := TradeOf(l.ID)
		// A floor with nobody on it and a floor with its weights condemned are
		// both floors that tell you nothing.
		if prop == nil || !runs || prop.Trouble || prop.Staff < trade.Hands {
			continue
		}
		return l.ID
	}
	return ""
}

// MarketWord is what the floor says about where a price stands, and nothing at
// all to somebody who does not hold it.
func (w *World) MarketWord(good string) string {
	if w.TheExchange() == "" {
		return ""
	}
	g := w.Good(good)
	if g == nil || g.Base <= 0 {
		return ""
	}
	against := g.Price * 100 / g.Base
	// Which way it is leaning. The drift is a quarter of the gap and it is not
	// a guess: it happens every time the market moves.
	switch {
	case against >= ExchangeHigh:
		return fmt.Sprintf("Your floor says %d%% over what it is worth, and it comes back down a quarter of the way every time the market moves.", against-100)
	case against <= ExchangeLow:
		return fmt.Sprintf("Your floor says %d%% under what it is worth, and a quarter of that comes back every time the market moves.", 100-against)
	default:
		return "Your floor says this is about what it is worth."
	}
}

// TheWarPremium is the floor's other piece of knowledge: while two
// organizations are fighting, scarcity holds every price up, and the week the
// shooting stops is the week the bottom goes out of it.
func (w *World) TheWarPremium() string {
	if w.TheExchange() == "" || !w.CityAtWar() {
		return ""
	}
	return " Everything is up while the shooting is on, and it will not stay up."
}
