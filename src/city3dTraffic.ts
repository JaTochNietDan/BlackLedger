import {onRoute, PITCH} from './city3dPlan.js';
import type {Point} from './city3dPlan.js';

export type TrafficPose = Point & {heading: number};
export type TrafficRequest = {id: string; model: string; points: Point[]; progress: number};
export type TrafficPlacement = {pose: TrafficPose; progress: number; waiting: boolean};
const lengths: Record<string, number> = {ford: 4.7, hudson: 5.1, packard: 5.8, police: 4.7};
export function trafficSize(model: string) {
  return {length: lengths[model] || 0.7, width: lengths[model] ? 2.15 : 0.7};
}
export function trafficOverlap(a: TrafficPose, am: string, b: TrafficPose, bm: string) {
  const as = trafficSize(am),
    bs = trafficSize(bm);
  const af = {x: Math.sin(a.heading), z: Math.cos(a.heading)};
  const ar = {x: af.z, z: -af.x};
  const bf = {x: Math.sin(b.heading), z: Math.cos(b.heading)};
  const br = {x: bf.z, z: -bf.x};
  const dot = (p: Point, q: Point) => p.x * q.x + p.z * q.z;
  const delta = {x: b.x - a.x, z: b.z - a.z};
  for (const axis of [af, ar, bf, br]) {
    const extent =
      (Math.abs(dot(axis, af)) * as.length +
        Math.abs(dot(axis, ar)) * as.width +
        Math.abs(dot(axis, bf)) * bs.length +
        Math.abs(dot(axis, br)) * bs.width) /
        2 +
      0.2;
    if (Math.abs(dot(delta, axis)) >= extent) return false;
  }
  return true;
}
const routeLength = (points: Point[]) =>
  points.slice(1).reduce((sum, p, i) => sum + Math.hypot(p.x - points[i].x, p.z - points[i].z), 0);
function junction(p: Point) {
  const x = Math.round(p.x / PITCH),
    z = Math.round(p.z / PITCH);
  return Math.abs(p.x - x * PITCH) <= 8 && Math.abs(p.z - z * PITCH) <= 8 ? `${x}:${z}` : null;
}
function throughAxis(request: TrafficRequest, progress: number, length: number) {
  if (!length) return null;
  const headings = [-10, 0, 10].map(
    offset => onRoute(request.points, progress + offset / length).heading,
  );
  if (headings.every(h => Math.abs(Math.sin(h)) > 0.999)) return 'horizontal';
  if (headings.every(h => Math.abs(Math.cos(h)) > 0.999)) return 'vertical';
  return null;
}
/** Presentation occupancy only. Saved journey progress remains the upper bound.
 * Small spatial steps prevent fast snapshot interpolation tunnelling through a car.
 * A departure waits inside its source until a physical space opens on its route. */
export class StreetTraffic {
  private entries = new Map<
    string,
    {key: string; progress: number; pose: TrafficPose; model: string; waiting: boolean}
  >();
  clear() {
    this.entries.clear();
  }
  update(requests: TrafficRequest[], seconds: number): Map<string, TrafficPlacement> {
    const wanted = new Set(requests.map(r => r.id));
    for (const id of this.entries.keys()) if (!wanted.has(id)) this.entries.delete(id);
    const lengths = new Map(requests.map(r => [r.id, routeLength(r.points)]));
    const byID = new Map(requests.map(r => [r.id, r]));
    // Preserve occupants first; new journeys yield to cars already on the road.
    const sorted = [...requests].sort(
      (a, b) =>
        Number(this.entries.has(b.id)) - Number(this.entries.has(a.id)) ||
        b.progress * lengths.get(b.id)! - a.progress * lengths.get(a.id)! ||
        a.id.localeCompare(b.id),
    );
    const free = (pose: TrafficPose, model: string, id: string, progress: number) => {
      const crossing = junction(pose);
      const axis = crossing ? throughAxis(byID.get(id)!, progress, lengths.get(id)!) : null;
      return ![...this.entries].some(([other, e]) => {
        if (other === id || e.waiting) return false;
        if (trafficOverlap(pose, model, e.pose, e.model)) return true;
        if (!crossing || crossing !== junction(e.pose)) return false;
        // Reserve the crossing before bodies enter it. Straight parallel lanes
        // can share it; turning or perpendicular traffic waits outside the box.
        return (
          axis === null || axis !== throughAxis(byID.get(other)!, e.progress, lengths.get(other)!)
        );
      });
    };
    for (const r of sorted) {
      const key = JSON.stringify([r.model, r.points]);
      let e = this.entries.get(r.id);
      if (e?.key !== key || e.progress > r.progress + 1e-8) {
        this.entries.delete(r.id);
        let progress = Math.max(0, Math.min(1, r.progress));
        const length = lengths.get(r.id)!;
        let pose = onRoute(r.points, progress);
        while (!free(pose, r.model, r.id, progress) && progress > 0) {
          progress = Math.max(0, progress - 0.35 / length);
          pose = onRoute(r.points, progress);
        }
        e = {key, model: r.model, progress, pose, waiting: !free(pose, r.model, r.id, progress)};
        this.entries.set(r.id, e);
      }
      if (e.waiting) {
        if (!free(e.pose, e.model, r.id, e.progress)) continue;
        e.waiting = false;
      }
      const length = lengths.get(r.id)!;
      if (!length) continue;
      // A snapshot may jump minutes. Compress travel, but retain safe occupancy.
      const target = Math.max(
        e.progress,
        Math.min(r.progress, e.progress + (Math.min(0.1, seconds) * 80) / length),
      );
      while (e.progress < target) {
        const next = Math.min(target, e.progress + 0.25 / length),
          pose = onRoute(r.points, next);
        if (!free(pose, e.model, r.id, next)) break;
        e.progress = next;
        e.pose = pose;
      }
    }
    return new Map(
      [...this.entries].map(([id, e]) => [
        id,
        {pose: e.pose, progress: e.progress, waiting: e.waiting},
      ]),
    );
  }
}
