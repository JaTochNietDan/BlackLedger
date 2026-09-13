/** Deterministic presentation of a committed blast; no gameplay physics. */
const smooth = (a: number, b: number, value: number) => {
  const t = Math.max(0, Math.min(1, (value - a) / (b - a)));
  return t * t * (3 - 2 * t);
};
export function blastParticle(index: number, seconds: number) {
  const smoke = index >= 12;
  const delay = smoke ? .08 + (index % 5) * .055 : (index % 4) * .018;
  const age = Math.max(0, seconds - delay);
  const angle = index * 2.399963;
  const fade = smoke ? (seconds >= 3 ? 0 : 1) : 1 - smooth(.25, .85, age);
  const radius = smoke ? .4 + age * (.45 + index % 3 * .1) : age * (2.5 + index % 4 * .5);
  return {
    x: Math.cos(angle) * radius,
    y: smoke ? .8 + age * (1.8 + index % 3 * .2) : .4 + age * (1.2 + index % 3 * .4),
    // Keep the blast on the exposed facade side, rather than inside the model.
    z: -(Math.abs(Math.sin(angle)) * radius + .1),
    size: seconds < delay ? 0 : (smoke ? 1.1 + age * .85 : .3 + smooth(0, .12, age) * 2.5) * Math.sqrt(fade),
    smoke,
    color: smoke ? 0x77736c : index % 3 ? 0xff8e29 : 0xffe4a1,
  };
}
export function blastOpacity(seconds: number) {
  return .8 * (1 - smooth(2.3, 3, seconds));
}
export function blastLight(seconds: number) {
  return 160 * (1 - smooth(.025, .45, seconds));
}
/** Soft, irregular opacity for one billow, generated once at scene startup. */
export function billowAlpha(x: number, y: number) {
  const radius = Math.hypot(x, y);
  const lumps = Math.sin(x * 9 + Math.sin(y * 7)) * .075 + Math.cos(y * 11 - x * 4) * .055;
  const edge = 1 - smooth(.35, .94 + lumps, radius);
  const grain = .86 + .09 * Math.sin(x * 31 + y * 17) * Math.cos(y * 29 - x * 13);
  return Math.max(0, Math.min(1, edge * grain));
}
