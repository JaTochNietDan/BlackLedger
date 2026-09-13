import {grid, bounds} from './iso.js';
import type {Place} from './types';

export type Point = {x: number; z: number};
export type Lot = Point & {id: string; col: number; row: number; model: string};
export const PITCH = 28;
export const STREET_WIDTH = 8;
export const FOOTWAY = 4.8;
export const LANE = 1.6;
export const MODEL_LIMIT = 17;
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
  const a = entrance(from, driving),
    b = entrance(to, driving);
  const offset = driving ? LANE : FOOTWAY;
  const turn = from.col * PITCH + offset;
  const points = [a, {x: turn, z: a.z}, {x: turn, z: b.z}, b];
  return points.filter((p, i) => i === 0 || p.x !== points[i - 1].x || p.z !== points[i - 1].z);
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
