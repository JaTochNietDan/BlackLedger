/** Soft pavement projection, not an additional gameplay light or visibility rule. */
export function headlightAlpha(x: number, distance: number) {
  if (distance <= 0 || distance >= 1) return 0;
  const spread = .10 + distance * .88;
  const across = Math.max(0, 1 - (x / spread) ** 2);
  return across ** 2 * Math.min(1, distance / .08) * (1 - distance) ** 1.5;
}
export function headlightCentre(x: number, z: number, heading: number, length: number, side: number) {
  const forward = length / 2 + 3.5, lateral = side * .61;
  return {x: x + Math.sin(heading) * forward + Math.cos(heading) * lateral,
    z: z + Math.cos(heading) * forward - Math.sin(heading) * lateral};
}
