package core

// Campaigns begun before the city had a life of its own carry factions with no
// holdings, no people and no quarrels. Loading such a save on the current build
// would leave the world inert: nothing to fight over, nobody to promote, and no
// rivalry to escalate. This migration gives an existing campaign the same
// structure a new one starts with, without taking anything that already belongs
// to somebody.

// SaveVersion is the shape the current build writes.
const SaveVersion = 19

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
	defer w.SettleHeadquarters()
	for id, p := range w.Properties {
		if p != nil && w.Own(id) && p.BoughtLife == 0 {
			p.BoughtLife = w.Life
		}
	}
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

	// The casino's hourly number used to include the tables. Now the tables are
	// run every night off their own float, so the hourly number is the floor
	// take alone and an existing campaign is corrected to match.
	if prop := w.Properties["casino"]; prop != nil && prop.Income > 18 {
		prop.Income = 18
	}

	// A casino somebody already owned always had money behind its tables; the
	// campaign simply could not see it. Enough to run, not enough to relax.
	for _, l := range Locations {
		if HasBankroll(l.ID) && w.Own(l.ID) && w.Properties[l.ID].Bankroll == 0 {
			w.Properties[l.ID].Bankroll = BankrollLot * 2
		}
	}

	// A campaign begun before there was anybody in the building gets them now.
	w.ensureOfficials()

	// And a city that was eight people gets the rest of its population, and
	// somebody doing every job the game speaks through.
	w.Populate()
	w.SettleHousing()
	w.SettleApartments()

	// A campaign that already had a car is driving one in working order.
	if w.Player.Car > 0 && w.Player.CarWear == 0 {
		w.Player.CarWear = 100
	}

	// The top of the range used to come plated, as a sentence in its
	// description that the rules had started reading. Plate is fitted at a
	// garage now and no car carries any off the lot, so a campaign already
	// driving one keeps the cover it had rather than losing two stages
	// overnight to a change in how they are counted.
	if w.Player.Car >= 3 && w.Player.Plate == 0 {
		w.Player.Plate = PlateStages
	}
	for i := range w.NPCs {
		if n := &w.NPCs[i]; n.Car >= 3 && n.Plate == 0 {
			n.Plate = PlateStages
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
	// A good added after a campaign began is missing from it entirely, and a
	// market the player cannot see is a market that does not exist. Anything
	// new starts at its reference price so nobody inherits a windfall.
	for _, good := range newGoods() {
		if w.Good(good.ID) == nil {
			w.Goods = append(w.Goods, good)
		}
	}

	w.SettleNewPlaces()
	marinerAdded := false
	if p := w.Properties["room"]; p != nil {
		marinerAdded = p.Income == 0
		p.Income = PlaceIncome["room"]
	}

	// Businesses acquired before they had an inside were working concerns all
	// along. Without this they would read as unstaffed and unstocked, and start
	// earning a fraction of what the campaign had come to expect.
	for _, l := range Locations {
		// This upgrade only equips the newly introduced lodging trade. Do not
		// refill existing businesses whose owner ran out of stock or staff.
		if w.Version >= 16 && (l.ID != "room" || !marinerAdded) {
			continue
		}
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

// SettleNewPlaces gives a campaign a record for any address added to the city
// since it began. It runs on EVERY load, not only when the save version has
// moved: adding four businesses to the city crashed every existing save,
// because those saves were already at the current version and the repair was
// gated behind a version bump nobody had made.
//
// Anything new starts unowned, in good order, earning what that address earns,
// and — if it is a trading business — already staffed and stocked, because a
// business is a going concern before anybody buys it.
func (w *World) SettleNewPlaces() {
	for _, l := range Locations {
		if w.Properties[l.ID] != nil {
			continue
		}
		prop := &Property{Owner: "independent", Condition: 100, Income: PlaceIncome[l.ID]}
		if trade, running := TradeOf(l.ID); running {
			prop.Staff, prop.Supply = trade.Hands, trade.RestockAmount
		}
		w.Properties[l.ID] = prop
	}
}
