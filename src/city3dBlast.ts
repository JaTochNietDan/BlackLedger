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

const fragmentVariation = (index: number) => {
  let h = Math.imul(index + 1, 0x45d9f3b);
  h = Math.imul(h ^ (h >>> 16), 0x45d9f3b);
  return ((h ^ (h >>> 16)) >>> 0) / 4294967296;
};
/** Small facade fragments settle independently; paths never enter the road. */
export function debrisPose(index: number, seconds: number) {
  const delay = (index % 3) * .025;
  const duration = .62 + (index % 4) * .09;
  const u = Math.max(0, Math.min(1, (seconds - delay) / duration));
  const lane = (index - 5.5) * .5;
  const spin = (1 - smooth(.65, 1, u));
  return {
    x: lane * (.75 + .25 * u) + (fragmentVariation(index + 12) - .5) * .1,
    z: -.25 - u * (.18 + fragmentVariation(index) * 1.7),
    height: u === 1 ? 0 : .4 * (1 - u) + Math.sin(Math.PI * u) * (.6 + index % 4 * .15),
    rx: u * (5 + index % 3) * spin, ry: index * 1.7 + u * 3,
    rz: spin ? u * (index % 2 ? -4 : 4) * spin : 0,
    scale: seconds < delay || seconds >= 3 ? 0 : .85 + index % 3 * .1,
  };
}
export function fragmentBlocked(x: number, z: number, body: {x: number; z: number; heading: number}, width: number, length: number) {
  const dx=x-body.x, dz=z-body.z, c=Math.cos(body.heading), s=Math.sin(body.heading);
  return Math.abs(dx*c-dz*s)<width/2+.15 && Math.abs(dx*s+dz*c)<length/2+.15;
}
