import {grid, bounds} from './iso.js';
import type {Place} from './types';

export type Point = {x: number; z: number; heading?: number};
export type Lot = Point & {id: string; col: number; row: number; model: string};
export const PITCH = 32;
export const STREET_WIDTH = 8;
export const FOOTWAY = 4.65;
export const LANE = 1.6;
export const MODEL_LIMIT = 17;
/** Paired crossing lines between existing pavement islands, aligned to walking lanes. */
export function crossingPaint(cols: number, rows: number) {
  const paint: {x: number; z: number; width: number; depth: number}[] = [];
  for (let col = 0; col <= cols; col++) for (let row = 0; row <= rows; row++) {
    for (const side of [-1, 1]) {
      if (col > 0 && col < cols && row + (side < 0 ? -1 : 0) >= 0 && row + (side > 0 ? 1 : 0) <= rows)
        for (const edge of [-.55, .55]) paint.push({x: col * PITCH, z: row * PITCH + side * FOOTWAY + edge, width: STREET_WIDTH - .2, depth: .11});
      if (row > 0 && row < rows && col + (side < 0 ? -1 : 0) >= 0 && col + (side > 0 ? 1 : 0) <= cols)
        for (const edge of [-.55, .55]) paint.push({x: col * PITCH + side * FOOTWAY + edge, z: row * PITCH, width: .11, depth: STREET_WIDTH - .2});
    }
  }
  return paint;
}
export function surfaceHeight(at: Point) {
  const toGrid = (value: number) => Math.abs(value - Math.round(value / PITCH) * PITCH);
  return Math.min(toGrid(at.x), toGrid(at.z)) <= STREET_WIDTH / 2 ? -0.1 : 0.17;
}
/** Conservative stride support: finish rising before either shoe reaches a kerb.
 * The 68cm radius encloses both authored animated pedestrian footprints.
 * Ease across the road side of the edge, retaining the existing gait clearance.
 */
export function pedestrianRootHeight(at: Point) {
  const toGrid = (value: number) => Math.abs(value - Math.round(value / PITCH) * PITCH);
  const edgeDistance = Math.min(toGrid(at.x), toGrid(at.z));
  const t = Math.max(0, Math.min(1, (edgeDistance + 0.68 - (STREET_WIDTH / 2 - 0.4)) / 0.4));
  return -0.07 + 0.27 * t * t * (3 - 2 * t);
}
/** Exported car tyres start 2cm above their root; retain a 5mm surface gap. */
export function vehicleRootHeight(at: Point) {
  return surfaceHeight(at) - 0.015;
}
export function streetsidePosition(lot: Lot): Point {
  return {x: lot.x, z: lot.z + 9};
}
export function lampPositions(lot: Lot): Point[] {
  return [-1, 1].map(side => ({x: lot.x + side * 8.7, z: lot.row * PITCH + 5.6}));
}
export function parkingSpot(lot: Lot): Point {
  return {x: lot.x + 9.6, z: lot.z};
}
export function cityPlan(places: Pick<Place, 'id' | 'x' | 'y' | 'type'>[]) {
  const cells = grid(places);
  const lots: Lot[] = places.map(p => {
    const c = cells.get(p.id)!;
    const specialized: Record<string, string> = {
      riverside: 'riverside-courts', mercercourt: 'mercer-court', room: 'mariner',
      docks: 'docks',
      filling: 'filling',
      pumps: 'filling',
      garage: 'garage',
      archway: 'garage',
      dealer: 'dealer',
      cabstand: 'dealer',
      chapel: 'undertaker',
      haulage: 'haulage',
      scrapyard: 'haulage',
      estate: 'villa',
      market: 'warehouse',
      tailor: 'shop',
      club: 'monarch',
      casino: 'bluehour',
      burlesque: 'papermoon',
      goldenlily: 'goldenlily',
    };
    const model =
      specialized[p.id] ||
      (p.type === 'casino'
        ? 'casino'
        : p.type === 'work'
          ? 'warehouse'
          : p.type === 'civic'
            ? 'civic'
            : p.type === 'home'
              ? 'tenement'
              : p.type === 'bar'
                ? 'tavern'
                : 'shop');
    return {id: p.id, ...c, x: (c.col + 0.5) * PITCH, z: (c.row + 0.5) * PITCH, model};
  });
  const extent = bounds(cells);
  const occupied = new Set(lots.map(lot => `${lot.col}:${lot.row}`));
  const vacant: (Point & {col: number; row: number})[] = [];
  for (let col = 0; col < extent.cols; col++) for (let row = 0; row < extent.rows; row++)
    if (!occupied.has(`${col}:${row}`)) vacant.push({col, row, x: (col + .5) * PITCH, z: (row + .5) * PITCH});
  return {lots, vacant, width: extent.cols * PITCH, depth: extent.rows * PITCH, ...extent};
}
export function entrance(lot: Lot, driving = false): Point {
  return {x: lot.x, z: lot.row * PITCH + (driving ? LANE : FOOTWAY)};
}
// Orthogonal street routes share the same pitch as the building footprints.
// Paths stay in the carriageway/footway, including at intermediate blocks.
export function route(from: Lot, to: Lot, driving = false): Point[] {
  const points = streetRoute(from, to, driving ? LANE : FOOTWAY, driving);
  if (!driving) {
    const start = entrance(from);
    if (points[0].x !== start.x || points[0].z !== start.z) points.unshift(start);
    const end = entrance(to);
    const last = points.at(-1)!;
    if (last.x !== end.x || last.z !== end.z) points.push(end);
  }
  return points;
}
// Offset each directed road segment onto its right-hand lane. Mitered
// intersections then receive a short curve so cars steer through the junction.
function streetRoute(from: Lot, to: Lot, offset: number, driving: boolean): Point[] {
  if (from.id === to.id) return [entrance(from, driving)];
  const a = {x: from.x, z: from.row * PITCH};
  const b = {x: to.x, z: to.row * PITCH};
  const spine =
    from.row === to.row
      ? [a, b]
      : [a, {x: from.col * PITCH, z: a.z}, {x: from.col * PITCH, z: b.z}, b];
  const segments = spine.slice(1).map((p, i) => {
    const q = spine[i],
      length = Math.hypot(p.x - q.x, p.z - q.z);
    const dx = (p.x - q.x) / length,
      dz = (p.z - q.z) / length;
    return {
      a: {x: q.x - dz * offset, z: q.z + dx * offset},
      b: {x: p.x - dz * offset, z: p.z + dx * offset},
      dx,
      dz,
    };
  });
  const lane = [segments[0].a];
  for (let i = 1; i < segments.length; i++) {
    const before = segments[i - 1],
      after = segments[i];
    lane.push(before.dx ? {x: after.a.x, z: before.b.z} : {x: before.b.x, z: after.a.z});
  }
  lane.push(segments.at(-1)!.b);
  if (!driving) return lane;
  const rounded = [lane[0]];
  for (let i = 1; i < lane.length - 1; i++) {
    const a = lane[i - 1],
      corner = lane[i],
      b = lane[i + 1];
    const incoming = Math.hypot(corner.x - a.x, corner.z - a.z);
    const outgoing = Math.hypot(b.x - corner.x, b.z - corner.z);
    const radius = Math.min(3, incoming / 2, outgoing / 2);
    const start = {
      x: corner.x + ((a.x - corner.x) * radius) / incoming,
      z: corner.z + ((a.z - corner.z) * radius) / incoming,
    };
    const end = {
      x: corner.x + ((b.x - corner.x) * radius) / outgoing,
      z: corner.z + ((b.z - corner.z) * radius) / outgoing,
    };
    rounded.push(start);
    for (let sample = 1; sample <= 12; sample++) {
      const t = sample / 12,
        s = 1 - t;
      rounded.push({
        x: s * s * start.x + 2 * s * t * corner.x + t * t * end.x,
        z: s * s * start.z + 2 * s * t * corner.z + t * t * end.z,
      });
    }
  }
  rounded.push(lane.at(-1)!);
  return rounded;
}
export function onRoute(points: Point[], progress: number): Point & {heading: number} {
  let length = 0;
  const lengths = points.slice(1).map((p, i) => {
    const d = Math.hypot(p.x - points[i].x, p.z - points[i].z);
    length += d;
    return d;
  });
  let left = Math.max(0, Math.min(1, progress)) * length;
  for (let i = 0; i < lengths.length; i++) {
    const d = lengths[i],
      a = points[i],
      b = points[i + 1];
    if (left <= d || i === lengths.length - 1) {
      const t = d ? left / d : 0;
      return {
        x: a.x + (b.x - a.x) * t,
        z: a.z + (b.z - a.z) * t,
        heading: Math.atan2(b.x - a.x, b.z - a.z),
      };
    }
    left -= d;
  }
  // Single-point scene reservations have an authored orientation. Erasing it
  // rotates asymmetric swept footprints when traffic admits the scene.
  return {...points[0], heading: points[0]?.heading ?? 0};
}
export function intersectsLot(point: Point, lot: Lot, clearance = 0) {
  return (
    Math.abs(point.x - lot.x) < MODEL_LIMIT / 2 + clearance &&
    Math.abs(point.z - lot.z) < MODEL_LIMIT / 2 + clearance
  );
}
