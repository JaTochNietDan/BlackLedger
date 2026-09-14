// Sampled gunshots and procedural Web Audio effects, with cancellable scene voices.

let context: AudioContext | null = null;

const sceneSounds = new Set<() => void>();
class SceneSound {
  private voices = new Set<() => void>();
  constructor() { sceneSounds.add(this.cancel); }
  add(stop: () => void) { this.voices.add(stop); }
  release(stop: () => void) {
    this.voices.delete(stop);
    if (!this.voices.size) sceneSounds.delete(this.cancel);
  }
  cancel = () => {
    for (const stop of [...this.voices]) stop();
    sceneSounds.delete(this.cancel);
  };
}

function startVoice(source: AudioScheduledSourceNode, nodes: AudioNode[], at: number, end: number, scene?: SceneSound, onFinish?:()=>void) {
  let finished = false;
  const cleanup = () => {
    if (finished) return;
    finished = true;
    onFinish?.();
    source.disconnect();
    nodes.forEach(node => node.disconnect());
    scene?.release(stop);
  };
  const stop = () => {
    if (finished) return;
    try { source.stop(); } catch { /* already ended or not started */ }
    cleanup();
  };
  source.onended = cleanup;
  scene?.add(stop);
  try {
    source.start(at);
    if(Number.isFinite(end))source.stop(end);
  } catch (error) {
    stop();
    throw error;
  }
  return stop;
}

// A browser will not let a page make a noise before somebody has touched it,
// so the context is built on the first sound and reused after that.
function audio(): AudioContext | null {
  if (context) return context;
  try {
    context = new (window.AudioContext || (window as any).webkitAudioContext)();
    return context;
  } catch {
    return null;
  }
}

export function soundOn(): boolean {
  try {
    return localStorage.getItem('black-ledger-sound') !== 'off';
  } catch {
    return true;
  }
}

export function setSound(on: boolean) {
  if (!on) for (const cancel of [...sceneSounds]) cancel();
  try {
    localStorage.setItem('black-ledger-sound', on ? 'on' : 'off');
  } catch {
    /* a private window */
  }
}

// noise is the raw material for anything that is not a tone: a shot, a blast,
// the crackle under a fire.
function noise(ctx: AudioContext, seconds: number): AudioBufferSourceNode {
  const frames = Math.max(1, Math.floor(ctx.sampleRate * seconds));
  const buffer = ctx.createBuffer(1, frames, ctx.sampleRate);
  const data = buffer.getChannelData(0);
  for (let i = 0; i < frames; i++) data[i] = Math.random() * 2 - 1;
  const source = ctx.createBufferSource();
  source.buffer = buffer;
  return source;
}

// One shot: a crack with almost no attack and a short tail.
function shot(ctx: AudioContext, at: number, level = 0.5, scene?: SceneSound) {
  const source = noise(ctx, 0.3);
  const band = ctx.createBiquadFilter();
  band.type = 'bandpass';
  band.frequency.value = 1800;
  band.Q.value = 0.7;
  const gain = ctx.createGain();
  gain.gain.setValueAtTime(0, at);
  gain.gain.linearRampToValueAtTime(level, at + 0.004);
  gain.gain.exponentialRampToValueAtTime(0.0001, at + 0.22);
  source.connect(band).connect(gain).connect(ctx.destination);
  return startVoice(source, [band, gain], at, at + 0.32, scene);
}

const gunSamples: Record<string,{url:string;gain:number}> = {
  revolver:{url:'/audio/guns/revolver.wav',gain:.75},
  shotgun:{url:'/audio/guns/shotgun.wav',gain:.65},
  thompson:{url:'/audio/guns/thompson.wav',gain:.8},
};
const gunBuffers=new Map<string,AudioBuffer>();
const sampledShots:Record<string,number>={};
const activeSamples=new Set<AudioBufferSourceNode>();
let loadingGuns:Promise<void>|undefined,gunBus:DynamicsCompressorNode|undefined,fallbackShots=0;
/** Decode once ahead of playback. Completion never starts a late sound. */
export function preloadCityGunshots():Promise<void> {
  if(loadingGuns)return loadingGuns;
  const ctx=audio();if(!ctx)return Promise.resolve();
  loadingGuns=Promise.all(Object.entries(gunSamples).map(async([name,sample])=>{
    try{
      const response=await fetch(sample.url);if(!response.ok)return;
      const buffer=await ctx.decodeAudioData(await response.arrayBuffer());gunBuffers.set(name,buffer);
    }catch{ /* A failed asset retains the immediate synthesized fallback. */ }
  })).then(()=>{});
  return loadingGuns;
}
export function cityGunshotStatus(){return {loaded:[...gunBuffers.keys()],played:{...sampledShots},fallbacks:fallbackShots,active:activeSamples.size,context:context?.state};}
/** Sample tails can overlap naturally; each voice remains cancellable on Skip/mute. */
export function playCityGunshot(weapon='revolver'): (() => void) | undefined {
  if (!soundOn()) return;
  const ctx = audio();
  if (!ctx) return;
  const scene=new SceneSound();
  try {
    if (ctx.state === 'suspended') ctx.resume().catch(() => {});
    const buffer=gunBuffers.get(weapon),sample=gunSamples[weapon];
    if(buffer&&sample){
      if(!gunBus){
        gunBus=ctx.createDynamicsCompressor();
        gunBus.threshold.value=-4;gunBus.knee.value=6;gunBus.ratio.value=12;
        gunBus.attack.value=.001;gunBus.release.value=.12;gunBus.connect(ctx.destination);
      }
      const source=ctx.createBufferSource(),gain=ctx.createGain();source.buffer=buffer;
      gain.gain.value=sample.gain;source.connect(gain).connect(gunBus);
      const at=ctx.currentTime+.005;activeSamples.add(source);
      startVoice(source,[gain],at,at+buffer.duration,scene,()=>activeSamples.delete(source));
      sampledShots[weapon]=(sampledShots[weapon]||0)+1;
    }else{
      shot(ctx,ctx.currentTime+.02,.4,scene);fallbackShots++;
    }
    return scene.cancel;
  } catch {
    scene.cancel();
    // Sound failure must never interrupt the city animation loop.
    return undefined;
  }
}

const effectSamples:Record<string,{gain:number;loop?:boolean}>={
  newspaper:{gain:.5},explosion:{gain:.65},raid:{gain:.65},siren:{gain:.35},'door-kick':{gain:.7},
  pain:{gain:.4},panic:{gain:.22},fire:{gain:.16,loop:true},
  'engine-idle':{gain:.12,loop:true},'vehicle-approach':{gain:.3},'drive-away':{gain:.4},
};
const effectBuffers=new Map<string,AudioBuffer>(),activeEffects=new Map<AudioBufferSourceNode,string>();
const playedEffects:Record<string,number>={};let loadingEffects:Promise<void>|undefined;
/** Crossfade a short end/start overlap instead of clicking at every loop seam. */
export function seamlessLoop(ctx:AudioContext,source:AudioBuffer){
  const fade=Math.min(Math.floor(source.sampleRate*.05),Math.floor(source.length/4));
  if(fade<2)return source;
  const result=ctx.createBuffer(source.numberOfChannels,source.length-fade,source.sampleRate);
  for(let ch=0;ch<source.numberOfChannels;ch++){
    const input=source.getChannelData(ch),out=result.getChannelData(ch);out.set(input.subarray(fade));
    for(let i=0;i<fade;i++){
      const t=i/(fade-1),at=out.length-fade+i;
      out[at]=input[source.length-fade+i]*(1-t)+input[i]*t;
    }
  }
  return result;
}
export function preloadCityEffects():Promise<void>{
  if(loadingEffects)return loadingEffects;const ctx=audio();if(!ctx)return Promise.resolve();
  loadingEffects=Promise.all(Object.entries(effectSamples).map(async([name,sample])=>{
    try{const response=await fetch(`/audio/effects/${name}.wav`);if(!response.ok)return;
      const decoded=await ctx.decodeAudioData(await response.arrayBuffer());
      effectBuffers.set(name,sample.loop?seamlessLoop(ctx,decoded):decoded);
    }catch{ /* Missing effects remain silent or use the existing procedural fallback. */ }
  })).then(()=>{});return loadingEffects;
}
export function cityEffectStatus(){return {loaded:[...effectBuffers.keys()],played:{...playedEffects},active:[...activeEffects.values()]};}
export function playRecordedEffect(name:string):(()=>void)|undefined{
  if(!soundOn())return;const ctx=audio(),buffer=effectBuffers.get(name),sample=effectSamples[name];
  if(!ctx||!buffer||!sample)return;
  const scene=new SceneSound();
  try{
    if(ctx.state==='suspended')ctx.resume().catch(()=>{});
    const source=ctx.createBufferSource(),gain=ctx.createGain();source.buffer=buffer;source.loop=!!sample.loop;
    gain.gain.value=sample.gain;source.connect(gain).connect(ctx.destination);
    activeEffects.set(source,name);const at=ctx.currentTime+.005;
    startVoice(source,[gain],at,sample.loop?Infinity:at+buffer.duration,scene,()=>activeEffects.delete(source));
    playedEffects[name]=(playedEffects[name]||0)+1;return scene.cancel;
  }catch{scene.cancel();return;}
}

/** One owned loop per sound; global mute cancellation can safely restart on a later frame. */
export class CityAmbientAudio {
  private loops=new Map<string,()=>void>();
  update(names:string[]){
    const wanted=new Set(soundOn()?names:[]);
    for(const [name,stop] of this.loops)if(!wanted.has(name)||!cityEffectStatus().active.includes(name)){stop();this.loops.delete(name);}
    for(const name of wanted)if(!this.loops.has(name)){
      const stop=playRecordedEffect(name);if(stop)this.loops.set(name,stop);
    }
  }
  dispose(){for(const stop of this.loops.values())stop();this.loops.clear();}
}

// Layered boot/wood impact, latch crack and short hinge scrape.
function doorBreach(ctx: AudioContext, at: number, scene: SceneSound) {
  for(const [delay,seconds,frequency,level] of [[0,.22,230,.55],[.015,.12,2600,.28],[.09,.38,750,.12]]) {
    const source=noise(ctx,seconds), band=ctx.createBiquadFilter(), gain=ctx.createGain();
    band.type='bandpass';band.frequency.value=frequency;band.Q.value=1.4;
    const begin=at+delay;
    gain.gain.setValueAtTime(.0001,begin);
    gain.gain.linearRampToValueAtTime(level,begin+.003);
    gain.gain.exponentialRampToValueAtTime(.0001,begin+seconds);
    source.connect(band).connect(gain).connect(ctx.destination);
    startVoice(source,[band,gain],begin,begin+seconds+.01,scene);
  }
}

// Initial glass crack followed by scattered high-frequency fragments.
function glassBreak(ctx:AudioContext,at:number,scene:SceneSound){
 for(let i=0;i<6;i++){
  const duration=i===0?.13:.07+i*.018,begin=at+(i===0?0:.07+i*.046);
  const source=noise(ctx,duration),band=ctx.createBiquadFilter(),gain=ctx.createGain();
  band.type='bandpass';band.frequency.value=2200+i*730;band.Q.value=i===0?.7:5;
  gain.gain.setValueAtTime(.0001,begin);gain.gain.linearRampToValueAtTime(i===0?.32:.09,begin+.002);
  gain.gain.exponentialRampToValueAtTime(.0001,begin+duration);
  source.connect(band).connect(gain).connect(ctx.destination);startVoice(source,[band,gain],begin,begin+duration+.01,scene);
 }
}

// A blast: low, long, and with a body you feel rather than hear.
function blast(ctx: AudioContext, at: number, scene: SceneSound) {
  const source = noise(ctx, 1.6);
  const low = ctx.createBiquadFilter();
  low.type = 'lowpass';
  low.frequency.setValueAtTime(900, at);
  low.frequency.exponentialRampToValueAtTime(90, at + 1.1);
  const gain = ctx.createGain();
  gain.gain.setValueAtTime(0, at);
  gain.gain.linearRampToValueAtTime(0.7, at + 0.02);
  gain.gain.exponentialRampToValueAtTime(0.0001, at + 1.5);
  source.connect(low).connect(gain).connect(ctx.destination);
  startVoice(source, [low, gain], at, at + 1.6, scene);

  // The thump under it.
  const body = ctx.createOscillator();
  body.type = 'sine';
  body.frequency.setValueAtTime(78, at);
  body.frequency.exponentialRampToValueAtTime(26, at + 0.8);
  const thump = ctx.createGain();
  thump.gain.setValueAtTime(0.6, at);
  thump.gain.exponentialRampToValueAtTime(0.0001, at + 0.9);
  body.connect(thump).connect(ctx.destination);
  startVoice(body, [thump], at, at + 1, scene);
}

// Two tones, alternating: a car at the kerb with its lamp turning.
function siren(ctx: AudioContext, at: number, scene: SceneSound, times = 4) {
  for (let i = 0; i < times; i++) {
    const when = at + i * 0.42;
    const tone = ctx.createOscillator();
    tone.type = 'square';
    tone.frequency.setValueAtTime(i % 2 ? 460 : 610, when);
    const gain = ctx.createGain();
    gain.gain.setValueAtTime(0, when);
    gain.gain.linearRampToValueAtTime(0.09, when + 0.05);
    gain.gain.setValueAtTime(0.09, when + 0.3);
    gain.gain.exponentialRampToValueAtTime(0.0001, when + 0.4);
    tone.connect(gain).connect(ctx.destination);
    startVoice(tone, [gain], when, when + 0.42, scene);
  }
}

// Something happened, and the city noticed: a dull knock, no drama.
function knock(ctx: AudioContext, at: number, scene: SceneSound) {
  const tone = ctx.createOscillator();
  tone.type = 'triangle';
  tone.frequency.setValueAtTime(220, at);
  tone.frequency.exponentialRampToValueAtTime(90, at + 0.18);
  const gain = ctx.createGain();
  gain.gain.setValueAtTime(0.22, at);
  gain.gain.exponentialRampToValueAtTime(0.0001, at + 0.3);
  tone.connect(gain).connect(ctx.destination);
  startVoice(tone, [gain], at, at + 0.32, scene);
}

// What each kind of moment sounds like. The kinds are the core's own, from
// core/witness.go, and a kind nobody has scored still makes the knock rather
// than nothing — silence reads as a bug.
export function playMoment(kind: string) {
  if (!soundOn()) return;
  const sample=({newspaper:'newspaper',explosion:'explosion',raid:'raid',arrest:'siren','door-breach':'door-kick'} as Record<string,string>)[kind];
  if(sample){const cancel=playRecordedEffect(sample);if(cancel)return cancel;}
  const ctx = audio();
  if (!ctx) return;
  const scene = new SceneSound();
  try {
    if (ctx.state === 'suspended') ctx.resume().catch(() => {});
    const at = ctx.currentTime + 0.02;
    switch (kind) {
      case 'newspaper':
        for(let i=0;i<3;i++){
          const source=noise(ctx,.22),band=ctx.createBiquadFilter(),gain=ctx.createGain(),begin=at+i*.07;
          band.type='bandpass';band.frequency.value=900+i*600;band.Q.value=.7;
          gain.gain.setValueAtTime(.0001,begin);gain.gain.linearRampToValueAtTime(.08,begin+.025);gain.gain.exponentialRampToValueAtTime(.0001,begin+.22);
          source.connect(band).connect(gain).connect(ctx.destination);startVoice(source,[band,gain],begin,begin+.23,scene);
        }
        break;
      case 'glass-break': glassBreak(ctx,at,scene); break;
      case 'door-breach': doorBreach(ctx, at, scene); break;
      case 'explosion': blast(ctx, at, scene); break;
      case 'killing':
        shot(ctx, at, 0.55, scene);
        shot(ctx, at + 0.17, 0.4, scene);
        break;
      case 'gunfight':
        for (let i = 0; i < 5; i++)
          shot(ctx, at + i * 0.13 + Math.random() * 0.04, 0.3 + Math.random() * 0.2, scene);
        break;
      case 'raid':
      case 'arrest': siren(ctx, at, scene); break;
      default: knock(ctx, at, scene);
    }
    return scene.cancel;
  } catch {
    scene.cancel();
    return undefined;
  }
}

// The room, and the machines in it.
//
// "We should also add ambient sounds and sounds to the slot machines and
// whatnot." Same contract as the rest of this file: synthesised, nothing
// downloaded, nothing licensed. A handle is a spring and a clunk, a drum
// stopping is a wooden knock, and a payout is a run of coins into a metal tray.

// clunk is a mechanism doing something: the handle going over, a drum stopping.
function clunk(ctx: AudioContext, at: number, pitch: number, level = 0.3) {
  const tone = ctx.createOscillator();
  tone.type = 'triangle';
  tone.frequency.setValueAtTime(pitch, at);
  tone.frequency.exponentialRampToValueAtTime(pitch * 0.45, at + 0.09);
  const body = ctx.createGain();
  body.gain.setValueAtTime(level, at);
  body.gain.exponentialRampToValueAtTime(0.0001, at + 0.13);
  tone.connect(body).connect(ctx.destination);
  tone.start(at);
  tone.stop(at + 0.15);

  // The rattle of the thing it is attached to.
  const rattle = noise(ctx, 0.08);
  const band = ctx.createBiquadFilter();
  band.type = 'bandpass';
  band.frequency.value = pitch * 4;
  band.Q.value = 2;
  const edge = ctx.createGain();
  edge.gain.setValueAtTime(level * 0.5, at);
  edge.gain.exponentialRampToValueAtTime(0.0001, at + 0.07);
  rattle.connect(band).connect(edge).connect(ctx.destination);
  rattle.start(at);
  rattle.stop(at + 0.09);
}

// coin is one piece of metal landing on other metal.
function coin(ctx: AudioContext, at: number) {
  const tone = ctx.createOscillator();
  tone.type = 'square';
  tone.frequency.setValueAtTime(1800 + Math.random() * 900, at);
  const gain = ctx.createGain();
  gain.gain.setValueAtTime(0, at);
  gain.gain.linearRampToValueAtTime(0.06, at + 0.003);
  gain.gain.exponentialRampToValueAtTime(0.0001, at + 0.12);
  tone.connect(gain).connect(ctx.destination);
  tone.start(at);
  tone.stop(at + 0.14);
}

// What each thing at the tables sounds like. Called by the interface at the
// moment the core says the thing happened, never on a timer of its own.
export function playTable(kind: string, count = 1) {
  if (!soundOn()) return;
  const ctx = audio();
  if (!ctx) return;
  if (ctx.state === 'suspended') ctx.resume().catch(() => {});
  const at = ctx.currentTime + 0.02;
  switch (kind) {
    case 'handle':
      clunk(ctx, at, 150, 0.35);
      clunk(ctx, at + 0.11, 110, 0.2);
      break;
    case 'reel':
      clunk(ctx, at, 320, 0.22);
      break;
    case 'coins':
      // Longer for a bigger win: the tray is how a machine tells the room.
      for (let i = 0; i < Math.max(4, Math.min(26, count)); i++)
        coin(ctx, at + i * 0.045 + Math.random() * 0.015);
      break;
    case 'card':
      clunk(ctx, at, 620, 0.12);
      break;
    // Chips going into the middle. Higher and drier than a coin in a tray,
    // because clay on baize is not brass on steel, and one for each rather
    // than one for the lot: a raise sounds like more money than a call.
    case 'chips':
      for (let i = 0; i < Math.max(1, Math.min(8, count)); i++)
        clunk(ctx, at + i * 0.055 + Math.random() * 0.02, 780 + Math.random() * 120, 0.09);
      break;
    case 'dice':
      for (let i = 0; i < 5; i++)
        clunk(ctx, at + i * 0.06 + Math.random() * 0.02, 420 + Math.random() * 200, 0.1);
      break;
    default:
      clunk(ctx, at, 240, 0.18);
  }
}

// The room itself: the hum of a busy floor under everything else. Started when
// the player sits down and stopped when they get up, because a noise that goes
// on after you have left the table is a noise nobody asked for.
let room: {gain: GainNode; nodes: AudioScheduledSourceNode[]} | null = null;

export function roomTone(on: boolean) {
  const ctx = audio();
  if (!ctx) return;
  if (!on || !soundOn()) {
    if (room) {
      const {gain, nodes} = room;
      room = null;
      gain.gain.setTargetAtTime(0, ctx.currentTime, 0.25);
      setTimeout(() => {
        nodes.forEach(n => {
          try {
            n.stop();
          } catch {
            /* already stopped */
          }
        });
      }, 900);
    }
    return;
  }
  if (room) return;
  if (ctx.state === 'suspended') ctx.resume().catch(() => {});
  const gain = ctx.createGain();
  gain.gain.setValueAtTime(0, ctx.currentTime);
  gain.gain.setTargetAtTime(0.035, ctx.currentTime, 0.8);
  gain.connect(ctx.destination);

  // A room full of people is low broadband noise with the top taken off it.
  const hum = noise(ctx, 4);
  hum.loop = true;
  const soft = ctx.createBiquadFilter();
  soft.type = 'lowpass';
  soft.frequency.value = 620;
  soft.Q.value = 0.4;
  hum.connect(soft).connect(gain);
  hum.start();

  room = {gain, nodes: [hum]};
}
