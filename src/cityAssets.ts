import buildings from '../public/art/buildings.json';
import fronts from '../public/art/fronts.json';
import previews from '../public/art/previews.json';

export const paintedLocations = buildings.map(b => b.id);
export const paintedBuildings = buildings;

// Five addresses were painted by hand for the street study. The rest showed a
// wireframe box in the address book, which made half the city look unfinished,
// so `tools/exteriors.py` painted a front for every one of them.
//
// Seven of those twenty-one were named in a list here and the other fourteen
// were never shown: the pictures were on disk the whole time and the address
// book went on drawing a wireframe box over them — "a lot of the building
// previews are empty in the addresses view." A list of ids written by hand is
// the same fault as a content table written for the smaller city, so the list
// is the files now (`public/art/fronts.json`, written by the same tool), and
// the wireframe stays as the last resort so a building added tomorrow still has
// a card rather than a hole.
const painted = new Set(fronts.map(f => f.id));
const frontFile = (id: string) => fronts.find(f => f.id === id)?.file ?? '';

export const paintedAsset = (id: string, condition = 100) => {
  const b = buildings.find(b => b.id === id);
  if (b) return '/art/' + (b.damage && condition < b.damage.below ? b.damage.file : b.file);
  const p = previews.find(p => p.id === id);
  if (p) return '/art/' + p.file;
  return painted.has(id) ? '/art/' + frontFile(id) : null;
};

// Whether this address is shown as a painted street view rather than as one of
// the isometric cut-outs. The two are drawn differently on purpose: a cut-out
// is a model of a building and fits its frame, a street view is a picture taken
// across the road and fills it.
export const paintedFront = (id: string) => painted.has(id);

export const paintedMask = (id: string) => {
  const p = previews.find(p => p.id === id);
  if (p) return '/art/' + p.mask;
  // A photographic front has no silhouette to cut it out with, and should not
  // be masked by itself: that punches the picture out of its own frame.
  return painted.has(id) ? null : paintedAsset(id);
};

// The car on the forecourt. Three cars were three lines of text that looked
// identical on the way past; these are painted offline the same way the
// buildings are (tools/cars.py) and shipped as files.
export const paintedCar = (tier: number) =>
  tier >= 1 && tier <= 3 ? `/art/cars/car-${tier}-v1.jpg` : null;
