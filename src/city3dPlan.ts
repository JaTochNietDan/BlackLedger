import {grid, bounds} from './iso.js';
import type {Place} from './types';

export type Point = {x: number; z: number};
export type Lot = Point & {id: string; col: number; row: number; model: string};
export const PITCH = 28;
export const STREET_WIDTH = 8;
export const FOOTWAY = 4.8;
export const LANE = 1.6;
export const MODEL_LIMIT = 17;
export function streetsidePosition(lot: Lot): Point {
  return {x: lot.x, z: lot.z + 9};
}
export function parkingSpot(lot: Lot): Point {
  return {x: lot.x + 9.65, z: lot.z};
}
export function cityPlan(places: Pick<Place, 'id' | 'x' | 'y' | 'type'>[]) {
  const cells = grid(places);
  const lots: Lot[] = places.map(p => {
    const c = cells.get(p.id)!;
    const specialized: Record<string, string> = {
      docks: 'docks',
      filling: 'filling',
      pumps: 'filling',
      garage: 'garage',
      archway: 'garage',
      dealer: 'dealer',
      cabstand: 'dealer',
      chapel: 'chapel',
      haulage: 'haulage',
      scrapyard: 'haulage',
      estate: 'villa',
      market: 'warehouse',
      tailor: 'shop',
      burlesque: 'casino',
      goldenlily: 'casino',
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
  return {lots, width: extent.cols * PITCH, depth: extent.rows * PITCH, ...extent};
}
export function entrance(lot: Lot, driving = false): Point {
  return {x: lot.x, z: lot.row * PITCH + (driving ? LANE : FOOTWAY)};
}
// Orthogonal street routes share the same pitch as the building footprints.
// Paths stay in the carriageway/footway, including at intermediate blocks.
export function route(from: Lot, to: Lot, driving = false): Point[] {
  if (driving) return drivingRoute(from, to);
  const a = entrance(from, driving),
    b = entrance(to, driving);
  const offset = driving ? LANE : FOOTWAY;
  const turn = from.col * PITCH + offset;
  const points = [a, {x: turn, z: a.z}, {x: turn, z: b.z}, b];
  return points.filter((p, i) => i === 0 || p.x !== points[i - 1].x || p.z !== points[i - 1].z);
}
// Offset each directed road segment onto its right-hand lane. Mitered
// intersections then receive a short curve so cars steer through the junction.
function drivingRoute(from: Lot, to: Lot): Point[] {
  if (from.id === to.id) return [entrance(from, true)];
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
      a: {x: q.x - dz * LANE, z: q.z + dx * LANE},
      b: {x: p.x - dz * LANE, z: p.z + dx * LANE},
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
  return {...points[0], heading: 0};
}
export function intersectsLot(point: Point, lot: Lot, clearance = 0) {
  return (
    Math.abs(point.x - lot.x) < MODEL_LIMIT / 2 + clearance &&
    Math.abs(point.z - lot.z) < MODEL_LIMIT / 2 + clearance
  );
}
