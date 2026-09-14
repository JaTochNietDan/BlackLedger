/** Cosmetic twilight follows game minutes, including interpolated travel time. */
export function cityNightAmount(minute: number) {
  const time = ((minute % 1440) + 1440) % 1440;
  const smooth = (v: number) => { const t = Math.max(0, Math.min(1, v)); return t * t * (3 - 2 * t); };
  return time < 720 ? 1 - smooth((time - 330) / 60) : smooth((time - 1170) / 60);
}
const blendColor = (day: number, night: number, amount: number) => {
  let color = 0;
  for (const shift of [16, 8, 0]) {
    const a = (day >> shift) & 255, b = (night >> shift) & 255;
    color |= Math.round(a + (b - a) * amount) << shift;
  }
  return color;
};
export function cityWeather(sky: {kind: string; wet: number} | undefined, night: boolean | number) {
  const amount = typeof night === 'boolean' ? Number(night) : Math.max(0, Math.min(1, night));
  const blend = (day: number, dark: number) => day + (dark - day) * amount;
  const wet = Number.isFinite(sky?.wet) ? Math.max(0, Math.min(1, sky!.wet)) : 0;
  const rain = sky?.kind === 'rain';
  const dull = rain || sky?.kind === 'overcast' || sky?.kind === 'fog';
  return {wet, rain, sun: blend(dull ? 1.05 : 3.2, .35),
    ambient: blend(dull ? 1.8 : 2.1, .9),
    background: blendColor(dull ? 0x5f6970 : 0x657477, 0x17232c, amount),
    fogFar: sky?.kind === 'fog' ? 440 : rain ? 560 : 850,
    roadRoughness: .92 - wet * .48, pavementRoughness: .9 - wet * .25,
    roadTone: 1 - wet * .28, pavementTone: 1 - wet * .12};
}
const hash = (n: number) => { const v = Math.sin(n * 127.1 + 311.7) * 43758.5453; return v - Math.floor(v); };
/** One reusable vertex buffer, with no gameplay clock or per-frame allocations. */
export function rainVertices(vertices: Float32Array, seconds: number, width: number, depth: number) {
  for (let i = 0; i < vertices.length / 6; i++) {
    const phase = (hash(i + 9) + seconds * (12 + hash(i + 5) * 7) / 32) % 1;
    const y = .2 + 32 * (1 - phase), length = .6 + hash(i + 3) * .8;
    const x = hash(i + 1) * width + (32 - y) * .08, z = hash(i + 2) * depth;
    vertices[i*6] = x; vertices[i*6+1] = y; vertices[i*6+2] = z;
    vertices[i*6+3] = x - length * .08; vertices[i*6+4] = Math.min(32.2, y + length); vertices[i*6+5] = z;
  }
}
