package core

import (
	"fmt"
	"sort"
)

type PoolHostInput struct {
	Fee        int  `json:"fee"`
	CutPercent int  `json:"cut_percent"`
	Enter      bool `json:"enter"`
}

func (w *World) HostPoolTournament(settings *PoolHostInput) error {
	if !w.Own(PoolPlace) {
		return fmt.Errorf("you must own The Green Baize to arrange its tournaments")
	}
	if settings == nil {
		return fmt.Errorf("set the entry fee, house cut and whether you will play")
	}
	ids := []string{}
	for i := range w.NPCs {
		n := &w.NPCs[i]
		if !n.Dead && n.Held <= w.Minute && n.Location == PoolPlace && !w.Travelling(n) && w.Pockets(n) >= settings.Fee {
			ids = append(ids, n.ID)
		}
	}
	sort.Strings(ids)
	needed := 4
	if settings.Enter {
		needed = 3
	}
	if len(ids) >= needed+4 {
		needed += 4
	}
	if len(ids) < needed {
		return fmt.Errorf("need %d funded local opponents at this entry fee; %d are available", needed, len(ids))
	}
	if err := w.startPoolTournament(ids[:needed], settings.Fee, settings.CutPercent, settings.Enter); err != nil {
		return err
	}
	opens, _ := w.poolTournamentWindow()
	if w.Minute >= opens {
		w.LastPoolTournamentSlot = opens
	}
	return nil
}
