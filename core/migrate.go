package core

// Campaigns begun before the city had a life of its own carry factions with no
// holdings, no people and no quarrels. Loading such a save on the current build
// would leave the world inert: nothing to fight over, nobody to promote, and no
// rivalry to escalate. This migration gives an existing campaign the same
// structure a new one starts with, without taking anything that already belongs
// to somebody.

// SaveVersion is the shape the current build writes.
const SaveVersion = 6

// seedHoldings is the property each established family holds in a new city.
var seedHoldings = map[string][]string{
	"bellandi": {"club", "docks"},
	"russo":    {"market", "bar"},
}

// holdingIncome is what those premises earn their owner.
var holdingIncome = map[string]int{"club": 30, "market": 18, "docks": 22, "bar": 12}

// MigrateLivingWorld is idempotent: running it on an already-migrated save
// changes nothing.
func (w *World) MigrateLivingWorld() {
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.Peak == 0 {
			f.Peak = f.Power
		}
		// Claim only premises nobody holds. A property the player owns, or that
		// a dead protagonist's estate still holds, is never taken by this.
		for _, id := range seedHoldings[f.ID] {
			if prop := w.Properties[id]; prop != nil && prop.Owner == "independent" {
				prop.Owner = f.ID
			}
		}
		for _, id := range w.FamilyHoldings(f.ID) {
			if prop := w.Properties[id]; prop != nil && prop.Income == 0 {
				prop.Income = holdingIncome[id]
			}
		}
	}

	// Attach each leader to the organization they lead, then make sure somebody
	// could replace them.
	for i := range w.Factions {
		f := &w.Factions[i]
		for j := range w.NPCs {
			n := &w.NPCs[j]
			if n.Dead || n.Name != f.Leader || n.Faction != "" {
				continue
			}
			n.Faction, n.Rank = f.ID, RankLeader
			if n.Location == "" {
				n.Location = w.homeOf(f.ID)
			}
			if n.Ambition == 0 {
				n.Ambition = 70
			}
			if n.Skill == 0 {
				n.Skill = 75
			}
		}
		for len(w.Members(f.ID)) < 3 {
			role, rank := "Soldier", RankSoldier
			if len(w.Members(f.ID)) < 2 {
				role, rank = "Lieutenant", RankLieutenant
			}
			if w.AddMember(f.ID, role, rank, w.homeOf(f.ID)) == nil {
				break
			}
		}
	}

	// Everyone is somewhere, and everyone sounds like themselves.
	for j := range w.NPCs {
		n := &w.NPCs[j]
		if n.Location == "" {
			n.Location = "bar"
		}
		if n.Voice == "" {
			n.Voice = w.voiceFor(n.Name)
		}
	}

	// A campaign that already had a suit before clothes could wear out is
	// wearing it in good order, not in rags.
	if w.Player.Dress > 0 && w.Player.DressWear == 0 {
		w.Player.DressWear = 100
	}

	// The established families were already rivals before any of this existed.
	if len(w.Conflicts) == 0 && len(w.Factions) >= 2 {
		w.Antagonize(w.Factions[0].ID, w.Factions[1].ID, 50)
	}

	// The underground trade has always been there; the campaign just could not
	// see it. Prices start at their reference so nobody inherits a windfall.
	if len(w.Goods) == 0 {
		w.Goods = newGoods()
	}

	// Businesses acquired before they had an inside were working concerns all
	// along. Without this they would read as unstaffed and unstocked, and start
	// earning a fraction of what the campaign had come to expect.
	for _, l := range Locations {
		prop := w.Properties[l.ID]
		trade, running := TradeOf(l.ID)
		if !running || prop == nil {
			continue
		}
		if prop.Staff == 0 {
			prop.Staff = trade.Hands
		}
		if prop.Supply == 0 {
			prop.Supply = trade.RestockAmount
		}
	}
}
