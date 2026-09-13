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
  if (kind === 'killing' || kind === 'gunfight')
    return (kind === 'gunfight' ? [-3, -6, 0, 3, 6] : [0, -3, 3, -6, 6]).map(offset => {
      const root = {x: lot.x + offset, z: lot.row * PITCH + 6.35};
      return {root, pose: {x: root.x + 0.8, z: root.z, heading: 0}, model: 'casualty'};
    });
  if (kind === 'raid' || kind === 'arrest')
    return [1, -1].flatMap(side => [-6, 0, 6].map(offset => {
      const root = {x: lot.x + side * 9.6, z: lot.z + offset};
      return {root, pose: {...root, heading: 0}, model: 'police'};
    }));
  return [];
}
export function availableSceneSlot(
  lot: Lot, kind: string, occupied: {pose: TrafficPose; model: string}[],
): SceneSlot | undefined {
  return sceneSlots(lot, kind).find(slot =>
    occupied.every(other => !trafficOverlap(slot.pose, slot.model, other.pose, other.model)));
}
export function casualtyFall(progress: number) {
  const angle = Math.min(1, Math.max(0, progress) * 2.5) * Math.PI / 2;
  // Rotating around the feet would push the jacket/arms through the pavement.
  return {rotation: -angle, height: 0.2 + 0.4 * Math.sin(angle)};
}

export const GUNFIRE_SHOTS = [0.7, 1.05, 1.5, 1.9] as const;

/** Three-second schematic gunfire sequence; no inferred target or damage. */
export function gunfightPose(seconds: number) {
  const ease = (t: number) => { const x = Math.max(0, Math.min(1, t)); return x*x*(3-2*x); };
  const aim = ease((seconds - 0.1) / 0.45) * (1 - ease((seconds - 2.3) / 0.6));
  const index = GUNFIRE_SHOTS.filter(at => at <= seconds).length;
  const age = index ? seconds - GUNFIRE_SHOTS[index - 1] : Infinity;
  const recoil = Math.max(0, 1 - age / 0.16) * 0.14;
  return {index, arm: -Math.PI / 2 * aim - recoil, flash: age < 0.065,
    smoke: age < 0.35 ? 1 - age / 0.35 : 0};
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
  private stop?: () => void;
  private closed = false;
  constructor(private fire: () => (() => void) | undefined) {}
  update(seconds: number, enabled = true) {
    if (this.closed) return;
    if (!enabled) {
      this.stop?.();
      this.stop = undefined;
    }
    const pose = gunfightPose(seconds);
    if (pose.index <= this.consumed) return;
    this.consumed = pose.index;
    if (!pose.flash || !enabled) return;
    this.stop?.();
    this.stop = this.fire();
    if (this.stop) this.started++;
  }
  dispose() {
    this.closed = true;
    this.stop?.();
    this.stop = undefined;
  }
}

/** The rendered cause supplies audio for co-located casualties in that moment. */
export function cityOwnsAudio(cue: VisualCue, batch: VisualCue[]) {
  if (cue.kind === 'gunfight' || cue.kind === 'explosion') return true;
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
  constructor(private fire: () => (() => void) | undefined) {}
  update(seconds: number, enabled = true) {
    if (this.closed) return;
    if (!enabled) { this.stop?.(); this.stop = undefined; }
    if (this.consumed) return;
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
