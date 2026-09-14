import {isPedestrian} from './city3dCast.js';
import {onRoute, PITCH, FOOTWAY} from './city3dPlan.js';
import type {Point} from './city3dPlan.js';

export type TrafficPose = Point & {heading: number};
export type TrafficRequest = {id: string; model: string; points: Point[]; progress: number};
export type TrafficPlacement = {pose: TrafficPose; progress: number; waiting: boolean; yielding?: boolean};
const lengths: Record<string, number> = {ford: 4.7, hudson: 5.1, packard: 5.8, police: 4.7, 'fire-engine':5.8};
export function trafficSpeed(model: string) {
  return isPedestrian(model) ? 1.8 : 11;
}
export function trafficSize(model: string) {
  if (model === 'planter') return {length:6.4,width:3.4};
  if (model === 'incendiary') return {length:1.8,width:6.2};
  if (model === 'custody') return {length:2,width:4.4};
  if (model === 'assassination') return {length:1.4,width:8.2};
  if (model === 'police-approach') return {length:6.6,width:1.4};
  if (model === 'casualty') return {length: 1.4, width: 2.6};
  const parked = model.startsWith('parked-');
  const vehicle = parked ? model.slice(7) : model;
  if(vehicle==='fire-engine')return {length:5.8,width:parked?2.35:2.55};
  return {length: lengths[vehicle] || 1.4, width: lengths[vehicle] ? (parked ? 2.15 : 2.35) : 0.85};
}
export function trafficModel(model: string, stationary: boolean) {
  return stationary && lengths[model] ? `parked-${model}` : model;
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
  const crossing = junction(onRoute(request.points, progress));
  const headings: number[] = [];
  for (let offset = 0; offset <= 32; offset++) {
    const pose = onRoute(request.points, progress + offset / length);
    if (junction(pose) !== crossing) break;
    headings.push(pose.heading);
  }
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
    {key: string; progress: number; pose: TrafficPose; model: string; waiting: boolean; yieldTo?: TrafficPose; yielding?: boolean; returnTo?: Point; joining?: Point}
  >();
  clear() {
    this.entries.clear();
  }
  placement(id: string): TrafficPlacement | undefined {
    const e=this.entries.get(id);
    return e ? {pose:{...e.pose},progress:e.progress,waiting:e.waiting,yielding:e.yielding} : undefined;
  }
  update(
    requests: TrafficRequest[],
    seconds: number,
    playbackRate = 1,
    pending: readonly {pose:TrafficPose;model:string}[] = [],
  ): Map<string, TrafficPlacement> {
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
      // A waiting scene closes its approach to new moving traffic. Existing
      // occupants may finish leaving; the scene itself still waits for clearance.
      const current=this.entries.get(id);
      if((lengths.get(id)!>0||(current&&!current.waiting&&isPedestrian(model)))&&pending.some(space=>trafficOverlap(pose,model,space.pose,space.model)&&
        !(current&&!current.waiting&&trafficOverlap(current.pose,current.model,space.pose,space.model))))return false;
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
      if(e?.key!==key&&e?.returnTo&&!e.waiting&&e.model===r.model&&lengths.get(r.id)!>0&&
        Math.hypot(r.points[0].x-e.returnTo.x,r.points[0].z-e.returnTo.z)<1e-8){
        // Rejoin a newly committed journey from the stance we actually rendered.
        // Only a known same-frontage sidestep can use this short connector.
        e={key,model:r.model,progress:0,pose:e.pose,waiting:false,joining:{...r.points[0]}};
        this.entries.set(r.id,e);
      }
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
      e.yielding=false;
      if(e.joining){
        let budget=Math.min(.1,Math.max(0,seconds))*playbackRate*trafficSpeed(r.model);
        while(budget>1e-8){
          const delta=e.joining.x-e.pose.x,distance=Math.abs(delta);
          if(distance<1e-8){e.joining=undefined;break;}
          const step=Math.min(distance,.125,budget);
          const next={x:e.pose.x+Math.sign(delta)*step,z:e.pose.z,heading:Math.sign(delta)*Math.PI/2};
          if(!free(next,r.model,r.id,0))break;
          e.pose=next;budget-=step;
        }
        // Journey progress stays at zero until its physical start is reached.
        continue;
      }
      if (!length) {
        // A stationary pedestrian may make room along this frontage, without
        // crossing a road or changing their authoritative location. Keep the
        // resulting stance until a real route replaces this request.
        const departures=requests.flatMap(other=>{
          if(!lengths.get(other.id)||!isPedestrian(other.model))return [];
          const entry=this.entries.get(other.id);
          if(entry?.waiting)return [{pose:entry.pose,model:other.model}];
          if(!entry&&other.progress===0)return [{pose:onRoute(other.points,0),model:other.model}];
          return [];
        });
        const spaces=[...pending,...departures];
        const inScene=(pose:TrafficPose)=>spaces.some(s=>trafficOverlap(pose,r.model,s.pose,s.model));
        const row=Math.floor(e.pose.z/PITCH),col=Math.floor(e.pose.x/PITCH);
        if(isPedestrian(r.model)&&Math.abs(e.pose.z-(row*PITCH+FOOTWAY))<.001){
          if(!spaces.length)e.yieldTo=undefined;
          if(!e.yieldTo&&inScene(e.pose)){
            for(let distance=.25;distance<=6&&!e.yieldTo;distance+=.25)for(const side of [-1,1]){
              const target={x:e.pose.x+side*distance,z:e.pose.z,heading:side*Math.PI/2};
              if(target.x<col*PITCH+8||target.x>col*PITCH+24||inScene(target))continue;
              let clear=true;
              for(let step=0;step<=distance;step+=.125){
                if(!free({...target,x:e.pose.x+side*step},r.model,r.id,e.progress)){clear=false;break;}
              }
              if(clear){e.yieldTo=target;break;}
            }
          }
          if(e.yieldTo){
            const delta=e.yieldTo.x-e.pose.x;
            const step=Math.min(Math.abs(delta),Math.min(.1,Math.max(0,seconds))*playbackRate*trafficSpeed(r.model));
            const next={...e.yieldTo,x:e.pose.x+Math.sign(delta)*step};
            if(free(next,r.model,r.id,e.progress)){
              if(step>0)e.returnTo={...r.points[0]};
              e.pose=next;e.yielding=step>0;
            }
            else e.yieldTo=undefined;
            if(e.yieldTo&&Math.abs(e.pose.x-e.yieldTo.x)<1e-8)e.yieldTo=undefined;
          }
        }
        continue;
      }
      // A snapshot may jump minutes. Compress travel, but retain safe occupancy.
      const target = Math.max(
        e.progress,
        Math.min(
          r.progress,
          e.progress +
            (Math.min(0.1, Math.max(0, seconds)) * playbackRate * trafficSpeed(r.model)) / length,
        ),
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
        {pose: e.pose, progress: e.progress, waiting: e.waiting, yielding:e.yielding},
      ]),
    );
  }
}

/** Wheel roll follows distance actually rendered, including traffic waits. */
export const WHEEL_RADIUS = .37;
export function advanceWheel(angle: number, distance: number) {
  return (angle + Math.max(0, distance) / WHEEL_RADIUS) % (Math.PI * 2);
}

const HALF_TRACK = .87;
function wheelbase(model: string) { return (lengths[model] - .2) * .59; }
function steeringLimit(model: string) {
  // Keep the inside tyre within the authored/tested half-radian sweep.
  const base = wheelbase(model), tangent = Math.tan(.5);
  return Math.atan(base * tangent / (base + HALF_TRACK * tangent));
}
/** Ackermann steering: both front axles intersect the same rear-axle turn centre. */
export function frontWheelSteering(angle: number, model: string, lateral: number) {
  if (!lengths[model]) return 0;
  const base = wheelbase(model), limit = steeringLimit(model);
  const tangent = Math.tan(Math.max(-limit, Math.min(limit, angle)));
  return Math.atan(base * tangent / (base - lateral * tangent));
}

/** Average route curvature across the wheelbase avoids faceted-curve steering jitter. */
export function wheelSteering(points: Point[], progress: number, model: string) {
  const length = routeLength(points);
  if (length < .001 || !lengths[model]) return 0;
  const base = wheelbase(model);
  const radius = base / 2 / length;
  const before = Math.max(0, progress - radius), after = Math.min(1, progress + radius);
  if (after <= before) return 0;
  const a = onRoute(points,before).heading, b = onRoute(points,after).heading;
  const turn = Math.atan2(Math.sin(b-a),Math.cos(b-a));
  const limit = steeringLimit(model);
  return Math.max(-limit,Math.min(limit,Math.atan(base * turn / ((after-before)*length))));
}
export function advanceSteering(angle: number, target: number, distance: number) {
  return angle + (target-angle) * (1-Math.exp(-Math.max(0,distance)/.45));
}
