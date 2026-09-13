package core

// BuildingFire records an observable fire caused by a building detonation.
// The response runs on game time; presentation cannot extinguish or repair it.
type BuildingFire struct {
	ID             string `json:"id"`
	Target         string `json:"target"`
	Minute         int    `json:"minute"`
	BrigadeAt      int    `json:"brigade_at"`
	ExtinguishedAt int    `json:"extinguished_at"`
	CleanupAt      int    `json:"cleanup_at"`
}

func (w *World) igniteBuilding(target string) {
	active := w.ActiveBuildingFires()
	for i := range active {
		if active[i].Target == target {
			// A second detonation renews the fire, retaining an already arrived brigade.
			if w.Minute > active[i].Minute {
				active[i].Minute = w.Minute
				active[i].ExtinguishedAt = w.Minute + 45
				active[i].CleanupAt = w.Minute + 90
			}
			w.BuildingFires = active
			return
		}
	}
	w.BuildingFires = append(active, BuildingFire{ID: ID(), Target: target, Minute: w.Minute,
		BrigadeAt: w.Minute + 10, ExtinguishedAt: w.Minute + 45, CleanupAt: w.Minute + 90})
}

// ActiveBuildingFires returns a detached projection, including brigade attendance
// after extinguishing. Existing property condition remains authoritative damage.
func (w *World) ActiveBuildingFires() []BuildingFire {
	result := []BuildingFire{}
	for _, fire := range w.BuildingFires {
		if fire.Minute <= w.Minute && w.Minute < fire.CleanupAt {
			result = append(result, fire)
		}
	}
	return result
}
