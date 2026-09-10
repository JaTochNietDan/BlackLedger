export interface VoiceHandle {
  play(): Promise<void>;
  pause(): void;
  dispose(): void;
  onEnded(fn: () => void): void;
}
export interface VoiceDependencies {
  load(event: string, signal: AbortSignal): Promise<Blob>;
  audio(blob: Blob): VoiceHandle;
  status(text: string): void;
  unavailable(): void;
}
/**
 * Whether a reading is in flight, from the same status the player reports.
 *
 * The scene offered "Stop voice" at all times, beside a button reading "Read
 * aloud" — a control for stopping something that was not happening. The states
 * below are the player's own, so this cannot drift from them without the tests
 * that assert those states failing too.
 */
export function speaking(status: string) {
  return status === PREPARING || status === SPEAKING;
}

export const PREPARING = 'Preparing voice…',
  SPEAKING = 'Speaking…',
  IDLE = 'Read aloud',
  REPLAY = 'Replay voice',
  RETRY = 'Retry voice';

/** Owns cancellation and media lifetime, independently of React or model latency. */
export class VoicePlayer {
  private epoch = 0;
  private request: AbortController | null = null;
  private playing: VoiceHandle | null = null;
  constructor(private deps: VoiceDependencies) {}
  stop() {
    this.epoch++;
    this.request?.abort();
    this.request = null;
    if (this.playing) {
      this.playing.pause();
      this.playing.dispose();
      this.playing = null;
    }
    this.deps.status(IDLE);
  }
  async speak(event: string) {
    this.stop();
    const token = this.epoch,
      abort = new AbortController();
    this.request = abort;
    this.deps.status(PREPARING);
    try {
      const blob = await this.deps.load(event, abort.signal);
      if (token !== this.epoch) return;
      const player = this.deps.audio(blob);
      this.playing = player;
      player.onEnded(() => {
        if (token === this.epoch) {
          player.dispose();
          this.playing = null;
          this.deps.status(REPLAY);
        }
      });
      await player.play();
      if (token !== this.epoch) {
        player.pause();
        return;
      }
      this.deps.status(SPEAKING);
    } catch (err) {
      if (token !== this.epoch) return;
      if (this.playing) {
        this.playing.pause();
        this.playing.dispose();
        this.playing = null;
      }
      if ((err as Error).name !== 'AbortError') {
        this.deps.status(RETRY);
        this.deps.unavailable();
      }
    }
  }
}

/**
 * Who a scene is attributed to, or null when the city has nobody by that name.
 *
 * The scene panel resolved its speaker with `npcs.find(...)||npcs[0]`, so an id
 * the city did not know became whoever happened to be first in the list — a
 * real person, with their portrait, name and job, saying something they never
 * said. Core is the only source of truth for who spoke; when it does not name
 * somebody the view has to say nothing rather than name somebody itself.
 */
export function speakerOf<T extends {id: string}>(
  people: readonly T[],
  id: string | undefined,
): T | null {
  if (!id) return null;
  return people.find(n => n.id === id) ?? null;
}
