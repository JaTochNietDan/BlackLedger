import buildings from '../public/art/buildings.json';
import previews from '../public/art/previews.json';

export const paintedLocations = buildings.map(b => b.id);
export const paintedBuildings = buildings;

// Five addresses were painted by hand for the street study. The other seven
// showed a wireframe box in the address book, which made half the city look
// unfinished. Those have a painted front now (tools/exteriors.py); the
// wireframe stays as the last resort so a building added tomorrow still has a
// card rather than a hole.
const fronts = ['docks', 'apartment', 'garage', 'casino', 'estate', 'precinct', 'herald'];

export const paintedAsset = (id: string, condition = 100) => {
  const b = buildings.find(b => b.id === id);
  if (b) return '/art/' + (b.damage && condition < b.damage.below ? b.damage.file : b.file);
  const p = previews.find(p => p.id === id);
  if (p) return '/art/' + p.file;
  return fronts.includes(id) ? `/art/fronts/front-${id}-v1.jpg` : null;
};

// Whether this address is shown as a painted street view rather than as one of
// the isometric cut-outs. The two are drawn differently on purpose: a cut-out
// is a model of a building and fits its frame, a street view is a picture taken
// across the road and fills it.
export const paintedFront = (id: string) => fronts.includes(id);

export const paintedMask = (id: string) => {
  const p = previews.find(p => p.id === id);
  if (p) return '/art/' + p.mask;
  // A photographic front has no silhouette to cut it out with, and should not
  // be masked by itself: that punches the picture out of its own frame.
  return fronts.includes(id) ? null : paintedAsset(id);
};

// The car on the forecourt. Three cars were three lines of text that looked
// identical on the way past; these are painted offline the same way the
// buildings are (tools/cars.py) and shipped as files.
export const paintedCar = (tier: number) =>
  tier >= 1 && tier <= 3 ? `/art/cars/car-${tier}-v1.jpg` : null;
