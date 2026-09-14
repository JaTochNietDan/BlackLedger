import type {Point} from './city3dPlan';
export type CameraCommand = 'zoom-in' | 'zoom-out' | 'rotate-left' | 'rotate-right' |
  'reset' | 'pan-left' | 'pan-right' | 'pan-up' | 'pan-down';
/** Browser/OS shortcuts and IME input must retain their normal behavior. */
export function cameraCommand(event: {
  key: string; ctrlKey?: boolean; metaKey?: boolean; altKey?: boolean; isComposing?: boolean;
}): CameraCommand | null {
  if (event.ctrlKey || event.metaKey || event.altKey || event.isComposing) return null;
  const commands: Record<string, CameraCommand> = {
    '+': 'zoom-in', '=': 'zoom-in', '-': 'zoom-out',
    w: 'pan-up', a: 'pan-left', s: 'pan-down', d: 'pan-right',
    q: 'rotate-left', e: 'rotate-right', home: 'reset',
    arrowleft: 'pan-left', arrowright: 'pan-right', arrowup: 'pan-up', arrowdown: 'pan-down',
  };
  const key = event.key.toLowerCase();
  return Object.hasOwn(commands, key) ? commands[key] : null;
}
/** Ground-plane movement aligned with the camera's screen, at any azimuth. */
export function screenPan(camera: Point, target: Point, command: CameraCommand, distance = 5): Point {
  const dx = target.x - camera.x, dz = target.z - camera.z;
  const length = Math.hypot(dx, dz);
  const forward = length > 1e-8 ? {x: dx / length, z: dz / length} : {x: 0, z: 1};
  const right = {x: -forward.z, z: forward.x};
  const axis = command === 'pan-up' || command === 'pan-down' ? forward : right;
  const sign = command === 'pan-left' || command === 'pan-down' ? -1 : 1;
  return {x: axis.x * distance * sign, z: axis.z * distance * sign};
}

/** Held keys drive every rendered frame, independently of OS key repeat. */
export class KeyboardPan {
  private held = new Map<string, CameraCommand>();
  press(event: Parameters<typeof cameraCommand>[0]) {
    const command = cameraCommand(event);
    if (!command?.startsWith('pan-')) return false;
    this.held.set(event.key.toLowerCase(), command);
    return true;
  }
  release(key: string) { this.held.delete(key.toLowerCase()); }
  clear() { this.held.clear(); }
  step(camera: Point, target: Point, seconds: number, speed: number): Point {
    const commands = new Set(this.held.values());
    const horizontal = Number(commands.has('pan-right')) - Number(commands.has('pan-left'));
    const vertical = Number(commands.has('pan-up')) - Number(commands.has('pan-down'));
    const length = Math.hypot(horizontal, vertical);
    if (!length) return {x: 0, z: 0};
    // A suspended tab cannot accumulate a large camera jump on resumption.
    const distance = Math.max(0, Math.min(seconds, .05)) * speed / length;
    const right = screenPan(camera, target, 'pan-right', horizontal * distance);
    const up = screenPan(camera, target, 'pan-up', vertical * distance);
    return {x: right.x + up.x, z: right.z + up.z};
  }
}

/** Release even outside the canvas; focus loss must never leave a stuck camera. */
export function bindKeyboardPan(canvas: HTMLElement, pan: KeyboardPan) {
  const release = (event: KeyboardEvent) => pan.release(event.key);
  const clear = () => pan.clear();
  const modifiers = (event: KeyboardEvent) => {
    if (event.metaKey || event.ctrlKey || event.altKey || event.isComposing) clear();
  };
  window.addEventListener('keyup', release);
  window.addEventListener('keydown', modifiers);
  window.addEventListener('blur', clear);
  canvas.addEventListener('blur', clear);
  document.addEventListener('visibilitychange', clear);
  return () => {
    clear();
    window.removeEventListener('keyup', release);
    window.removeEventListener('keydown', modifiers);
    window.removeEventListener('blur', clear);
    canvas.removeEventListener('blur', clear);
    document.removeEventListener('visibilitychange', clear);
  };
}
