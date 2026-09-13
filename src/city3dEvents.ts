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
