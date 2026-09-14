import {PITCH} from './city3dPlan.js';
import type {Lot, Point} from './city3dPlan.js';
import {trafficOverlap} from './city3dTraffic.js';
import type {TrafficPose} from './city3dTraffic.js';
import type {VisualCue} from './types';
/** Each committed cue is staged once per mounted city. Loading an old save is
 * silent; an explicit currently playing cue survives navigation into the city. */
export class CityCueQueue {
  private seen = new Set<string>();
  private world = '';
  take(
    world: string,
    cues: VisualCue[],
    active: VisualCue | null,
    initial: boolean,
    replay = false,
  ): VisualCue[] {
    if (world !== this.world) {
      this.world = world;
      this.seen.clear();
    }
    const out: VisualCue[] = [];
    const batch = new Set<string>();
    for (const cue of [...cues, ...(active ? [active] : [])]) {
      if (batch.has(cue.id) || (!replay && this.seen.has(cue.id))) continue;
      batch.add(cue.id);
      this.seen.add(cue.id);
      if (!initial || active !== null) out.push(cue);
    }
    // Keep memory bounded. Cue IDs are stable; last_result has bounded history.
    while (this.seen.size > 2048) this.seen.delete(this.seen.values().next().value!);
    return out;
  }
}

export type SceneSlot = {root: Point; pose: TrafficPose; model: string};
/** Dedicated forecourt/side bays keep reenactments out of public travel lanes.
 * The casualty reservation encloses the whole fall, including the standing pose. */
export function sceneSlots(lot: Lot, kind: string): SceneSlot[] {
  if(kind==='incendiary')return [3,0,-3].map(offset=>{
    const root={x:lot.x+offset,z:lot.row*PITCH+6.35};
    return {root,pose:{x:root.x-1.5,z:root.z,heading:0},model:'incendiary'};
  });
  if(kind==='custody')return [-3,-7,1].map(offset=>{
    const root={x:lot.x+offset,z:lot.row*PITCH+6.35};
    return {root,pose:{x:root.x+1.5,z:root.z+.3,heading:0},model:'custody'};
  });
  if(kind==='assassination')return [-4,-8.6,-1].map(offset=>{
    const root={x:lot.x+offset,z:lot.row*PITCH+6.35};
    return {root,pose:{x:root.x+3.4,z:root.z,heading:0},model:'assassination'};
  });
  if (['killing','gunfight','officer','detainee','raid-officer'].includes(kind))
    return (kind === 'gunfight' ? [-3, -6, 0, 3, 6] : [0, -3, 3, -6, 6]).map(offset => {
      const root = {x: lot.x + offset, z: lot.row * PITCH + 6.35};
      return kind==='raid-officer' ? {root,pose:{x:root.x,z:root.z+2.6,heading:0},model:'police-approach'}
        : {root, pose: {x: root.x + 0.8, z: root.z, heading: 0}, model: 'casualty'};
    });
  if (['raid','arrest','police-unit','raid-unit','fire-engine'].includes(kind))
    return [1, -1].flatMap(side => [-6, 0, 6].map(offset => {
      const root = {x: lot.x + side * (kind==='fire-engine'?10:9.6), z: lot.z + offset};
      return {root, pose: {...root, heading: 0}, model: kind==='fire-engine'?'parked-fire-engine':'parked-police'};
    }));
  return [];
}
export function availableSceneSlot(
  lot: Lot, kind: string, occupied: {pose: TrafficPose; model: string}[], entry?:Point,
): SceneSlot | undefined {
  const candidates=sceneSlots(lot,kind);
  if(kind==='raid-officer'&&entry){
    // Keep a door-aligned approach available when a body occupies the forecourt.
    // Its conservative swept reservation also protects the entry after contact.
    const root={x:entry.x,z:Math.max(lot.row*PITCH+6.35,entry.z-1.55)};
    if(entry.z-root.z>.65&&entry.z-root.z<=3.65)
      candidates.splice(entry.z-candidates[0].root.z>3.65?0:1,0,{root,pose:{x:root.x,z:root.z+2.6,heading:0},model:'police-approach'});
  }
  return candidates.find(slot =>
    occupied.every(other => !trafficOverlap(slot.pose, slot.model, other.pose, other.model)));
}
export function casualtyFall(progress: number) {
  const angle = Math.min(1, Math.max(0, progress) * 2.5) * Math.PI / 2;
  // Rotating around the feet would push the jacket/arms through the pavement.
  return {rotation: -angle, height: 0.2 + 0.4 * Math.sin(angle)};
}

export const GUNFIRE_SHOTS = [0.7, 1.05, 1.5, 1.9] as const;

/** Three-second schematic gunfire sequence; no inferred target or damage. */
export function gunfightPose(seconds: number,beats:readonly number[]=GUNFIRE_SHOTS) {
  const ease = (t: number) => { const x = Math.max(0, Math.min(1, t)); return x*x*(3-2*x); };
  const aim = ease((seconds - 0.1) / 0.45) * (1 - ease((seconds - 2.3) / 0.6));
  const index = beats.filter(at => at <= seconds).length;
  const age = index ? seconds - beats[index - 1] : Infinity;
  const recoil = Math.max(0, 1 - age / 0.16) * 0.14;
  return {index, arm: -Math.PI / 2 * aim - recoil, flash: age < 0.065,
    smoke: age < 0.35 ? 1 - age / 0.35 : 0};
}

/** Match only the same recorded moment; explicit victim identity narrows modern strikes. */
export function gunVictim(gun:VisualCue,victim:VisualCue){
  return gun.kind==='gunfight'&&victim.kind==='killing'&&gun.target===victim.target&&gun.minute===victim.minute&&
    (!gun.strike||!!victim.actors?.some(actor=>actor.id===gun.strike!.victim.id));
}
export function gunCastReady(gun:VisualCue,cues:VisualCue[],ready:(id:string)=>boolean){
  return cues.filter(cue=>gunVictim(gun,cue)).every(cue=>ready(cue.id));
}

/** Co-located casualty playback waits until the associated gun scene fires.
 * Repeated holds also cover a gun scene waiting for an available staging slot. */
export function casualtySceneStart(since: number, now: number, gunSince?: number) {
  return gunSince !== undefined && now < gunSince + 700 ? now : since;
}

/** Sounds follow rendered muzzle pulses, with no future Web Audio schedule.
 * A stalled/background frame consumes old beats without replaying a backlog. */
export class GunfireAudio {
  started = 0;
  private consumed = 0;
  private stops=new Set<()=>void>();
  private cancel(){for(const stop of this.stops)stop();this.stops.clear();}
  private closed = false;
  constructor(private fire: () => (() => void) | undefined,private beats:readonly number[]=GUNFIRE_SHOTS,private overlap=false) {}
  update(seconds: number, enabled = true) {
    if (this.closed) return;
    if (!enabled) {
      this.cancel();
    }
    const pose = gunfightPose(seconds,this.beats);
    if (pose.index <= this.consumed) return;
    this.consumed = pose.index;
    if (!pose.flash || !enabled) return;
    if(!this.overlap)this.cancel();
    const stop=this.fire();
    if(stop){this.stops.add(stop);this.started++;}
  }
  dispose() {
    this.closed = true;
    this.cancel();
  }
}

/** The rendered cause supplies audio for co-located casualties in that moment. */
export function cityOwnsAudio(cue: VisualCue, batch: VisualCue[]) {
  if (['gunfight','explosion','raid','arrest'].includes(cue.kind)) return true;
  return cue.kind === 'killing' && batch.some(other =>
    (other.kind === 'gunfight' || other.kind === 'explosion') &&
    other.target === cue.target && other.minute === cue.minute);
}

/** A blast sounds once at the visible onset; a late frame never plays a backlog. */
export class BlastAudio {
  started = 0;
  private consumed = false;
  private closed = false;
  private stop?: () => void;
  constructor(private fire: () => (() => void) | undefined,private beats:readonly number[]=GUNFIRE_SHOTS) {}
  update(seconds: number, enabled = true) {
    if (this.closed) return;
    if (!enabled) { this.stop?.(); this.stop = undefined; }
    if (this.consumed || seconds < 0) return;
    this.consumed = true;
    if (seconds < 0 || seconds > .15 || !enabled) return;
    this.stop = this.fire();
    if (this.stop) this.started++;
  }
  dispose() {
    this.closed = true;
    this.stop?.();
    this.stop = undefined;
  }
}

/** Supporting presentation cast; these are not new backend events. Never infer
 * the prisoner from cue actors, which can name an arresting detective. */
export function policeCast(cue: VisualCue): VisualCue[] {
  if(!['raid','arrest'].includes(cue.kind))return [cue];
  const result=[cue];
  const add=(kind:string,index:number,actors=cue.actors)=>result.push({...cue,id:`${cue.id}:${kind}:${index}`,kind,actors});
  for(let i=0;i<(cue.kind==='raid'?2:1);i++)add(cue.kind==='raid'?'raid-unit':'police-unit',i,[]);
  if(cue.kind==='arrest'&&cue.detainee)add('detainee',0,[cue.detainee]);
  for(let i=0;i<(cue.kind==='raid'?4:cue.detainee?1:2);i++)add(cue.kind==='raid'?'raid-officer':'officer',i,[]);
  return result;
}

/** A staggered purposeful walk; callers reserve its complete swept path. */
export function officerApproach(seconds: number, distance: number, index: number) {
 const start=.35+index*.08;
 const travelled=Math.min(Math.max(0,distance),Math.max(0,seconds-start)*1.4);
 const walking=seconds>start&&travelled<distance;
 return {travelled,walking,phase:travelled/1.15*Math.PI*2};
}

/** Door contact precedes its opening; entry starts only after the leaf clears. */
export function raidEntryPose(seconds: number, distance: number) {
  const approach=officerApproach(seconds,distance,0);
  const arrived=.35+distance/1.4;
  const contact=seconds-arrived;
  const kick=contact>=0 && contact<.65 ? Math.sin(contact/.65*Math.PI) : 0;
  const door=Math.max(0,Math.min(1,(contact-.3)/.35));
  const inside=Math.max(0,Math.min(1.9,(contact-1)*1.4));
  const travelled=approach.travelled+inside;
  const walking=approach.walking || (contact>1 && inside<1.9);
  return {travelled,walking,phase:travelled/1.15*Math.PI*2,kick,door};
}
export function policeSceneSeconds(kind: string) {
  return ['arrest','police-unit','officer','detainee'].includes(kind)?9:kind==='explosion'?14:['raid','raid-unit','raid-officer'].includes(kind)?10:3;
}
