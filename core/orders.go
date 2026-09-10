package core

import "fmt"

// Work aimed at a place that does not put the player in it. Both of these were
// offered only while standing in the building they were about, which meant
// crossing the city to give an order somebody else was going to carry out —
// and neither rule ever asked where the player was. What they ask for is a
// crew, their loyalty, standing, and what you know.
//
// Going in yourself is not here. That one genuinely means going there, and it
// is the whole difference between the two buttons.
//
// These are NOT marked Anywhere. That flag is for work that belongs to the
// player rather than to any place — putting a price on a name, reaching an
// understanding — and the interface files it away from the premises panel
// entirely. An order aimed at a building belongs in that building's panel; what
// changed is that the panel is readable without walking there.
func (w *World) ordersAbout(l Place) []Action {
	out := []Action{}
	f, ok := w.SabotageTarget(l.ID)
	if !ok {
		return out
	}
	order := func(id, label string, minutes, cost int, reason, detail string) {
		if reason == "" && w.Player.Cash < cost {
			reason = "Not enough cash"
		}
		out = append(out, Action{
			Group: GroupOf(id), ID: id, Label: label, Minutes: minutes, Cost: cost,
			Disabled: reason != "", Reason: reason, Detail: detail, Target: l.ID,
		})
	}
	if hand, ok := w.CrewHands(); ok {
		reason := w.SendAgainstReadiness(l.ID)
		if reason == "" {
			reason = w.DelegateReadiness()
		}
		order("sabotage:crew", "Send "+hand.Name+" against "+l.Name, 90, 0, reason,
			fmt.Sprintf("The same damage to %s and worse odds. %d less attention on you and a fifth of the standing. Turned away, they take the beating and %d loyalty, and sometimes they do not come back.", f.Name, HandHeatRelief, HandLoyaltyCost))
	}
	if rival := w.Rival(f.ID); rival != nil {
		order("incite", "Point "+f.Name+" at "+rival.Name, 45, 25, w.InciteReadiness(l.ID),
			fmt.Sprintf("Spend $25 on the right conversations so %s %s %s moved against them. Hardens their quarrel and can start a war you are not part of. A story that does not hold up costs you standing with %s.", f.Name, Agree(f.Name, "believes", "believe"), rival.Name, f.Name))
	}
	return out
}
