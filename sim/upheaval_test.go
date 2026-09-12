package sim

import "testing"

// The living world is five layers deep on paper — families act on each other,
// organizations are created and destroyed, wars start and are settled, holdings
// change hands, people die for it — and the only measure of whether any of it
// happens is the city run with nobody playing it.
//
// That measure said nought. Not for one layer: organizations fell nought times
// in every run this project has ever taken, and wars started three times across
// twelve cities. It read as a layer that had been built and did not work.
//
// It was the window. Sixty days is shorter than the thing being looked at: a
// family has to lose every holding it has and be beaten down to nothing before
// it falls, and that takes about a hundred days of the city grinding at itself.
// The same twelve cities over two hundred days start a hundred wars and bury
// eighteen organizations.
//
// So this exists to stop the measure quietly going back to reporting nothing.
// It asks one question of each layer — does it ever happen — and it is the
// season it asks over that is doing the work.
func TestASeasonOfTheCityProducesEveryKindOfUpheaval(t *testing.T) {
	if testing.Short() {
		t.Skip("a season of twelve cities: run without -short")
	}
	t.Parallel()
	const cities, season = 12, 200
	total := CityReport{}
	for i := 0; i < cities; i++ {
		r := City(uint32(i+1)*2654435761, season)
		total.Formed += r.Formed
		total.Fell += r.Fell
		total.Wars += r.Wars
		total.Settled += r.Settled
		total.Changed += r.Changed
		total.Killed += r.Killed
	}
	t.Logf("%d cities over %d days each: %d organizations formed, %d fell, %d wars started, "+
		"%d settled, %d holdings changed hands, %d people killed",
		cities, season, total.Formed, total.Fell, total.Wars, total.Settled, total.Changed, total.Killed)
	for _, layer := range []struct {
		what  string
		count int
	}{
		{"no organization was ever created", total.Formed},
		{"no organization ever fell", total.Fell},
		{"no war ever started", total.Wars},
		{"no war was ever settled", total.Settled},
		{"no holding ever changed hands", total.Changed},
		{"nobody was ever killed over any of it", total.Killed},
	} {
		if layer.count == 0 {
			t.Errorf("%d cities ran for %d days each and %s", cities, season, layer.what)
		}
	}
}
