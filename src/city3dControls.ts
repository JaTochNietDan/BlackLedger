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
