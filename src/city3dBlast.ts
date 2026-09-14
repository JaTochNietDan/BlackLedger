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
    // Local outward travel; windowBurst adds the internal-to-external transition.
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

/** A pressure burst begins behind glazing and escapes through the authored window. */
export function windowBurst(index:number,seconds:number,window:{x:number;y:number;z:number}){
 const p=blastParticle(index,seconds);
 const escape=smooth(0,.18,seconds);
 return {...p,x:window.x+p.x*.35,y:window.y+p.y*.65,z:window.z+.45*(1-escape)+p.z*escape,size:p.size*.72};
}

/** Only a confirmed building fire (or explicit debug scene) identifies an internal detonation. */
export function internalDetonation(cue:{id:string;target:string;minute?:number;detonation?:string},fires:readonly {target:string;minute:number}[]){
 if(cue.detonation)return cue.detonation==='planted';
 return cue.id.startsWith('preview:')||fires.some(f=>f.target===cue.target&&f.minute===cue.minute);
}

/** Window-ejected masonry: fast outward impulse clears the canopy before falling. */
export function windowDebris(index:number,seconds:number,window:{x:number;y:number;z:number},ground:number,count=4){
 const delay=.06+(index%3)*.025,age=Math.max(0,seconds-delay);
 const height=Math.max(0,window.y-ground),up=.4,gravity=9.8;
 const duration=(up+Math.sqrt(up*up+2*gravity*height))/gravity;
 const u=Math.min(1,age/duration),lane=Math.floor(index/count)-(Math.ceil(12/count)-1)/2;
 const settle=1-smooth(.8,1,u);
 return {
  x:window.x+lane*(.22+.08*u),z:window.z-.08-(2.8+(index%3)*.15)*(2*u-u*u),
  height:Math.max(0,height+up*Math.min(age,duration)-gravity*Math.min(age,duration)**2/2),
  rx:age*(5+index%3)*settle,ry:index*1.7+u*3,rz:settle?age*(index%2?-4:4)*settle:0,
  scale:seconds<delay||seconds>=3?0:.85+(index%3)*.1,
 };
}
