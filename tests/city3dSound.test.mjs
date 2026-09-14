import test from 'node:test';
import assert from 'node:assert/strict';

test('city shot respects mute and releases its actual audio graph on end or cancellation',async()=>{
 const previousWindow=globalThis.window,previousStorage=globalThis.localStorage;
 let muted=true,created=0,failed=false;
 const sources=[],nodes=[];
 const parameter=()=>({value:0,setValueAtTime(){},linearRampToValueAtTime(){},exponentialRampToValueAtTime(){}});
 const node=()=>{const n={disconnected:0,connect(next){return next;},disconnect(){this.disconnected++;}};nodes.push(n);return n;};
 class Context {
  constructor(){created++;}
  state='running';currentTime=4;sampleRate=100;destination={};
  createBuffer(_,frames){if(failed)throw new Error('audio unavailable');return{getChannelData:()=>new Float32Array(frames)};}
  createBufferSource(){const n=Object.assign(node(),{starts:[],stops:[],start(at){this.starts.push(at);},stop(at){this.stops.push(at);}});sources.push(n);return n;}
  createBiquadFilter(){return Object.assign(node(),{frequency:parameter(),Q:parameter()});}
  createGain(){return Object.assign(node(),{gain:parameter()});}
 }
 globalThis.window={AudioContext:Context};
 globalThis.localStorage={getItem:()=>muted?'off':'on'};
 try{
  const {playCityGunshot}=await import('../.runtime/frontend-test/sound.js');
  assert.equal(playCityGunshot(),undefined);assert.equal(created,0);
  muted=false;const stop=playCityGunshot();assert.equal(created,1);
  assert.deepEqual(sources[0].starts,[4.02]);assert.ok(sources[0].stops[0]>4.3);
  stop();stop();assert.equal(sources[0].stops.length,2);
  assert.ok(nodes.every(n=>n.disconnected===1));
  const next=playCityGunshot();sources[1].onended();next();
  assert.ok(nodes.every(n=>n.disconnected===1));assert.equal(created,1);
  failed=true;assert.equal(playCityGunshot(),undefined);
 }finally{globalThis.window=previousWindow;globalThis.localStorage=previousStorage;}
});

test('scene sound cancellation stops scheduled explosions and sirens, mute cancels every scene, and natural endings release graphs',async()=>{
 const previousWindow=globalThis.window,previousStorage=globalThis.localStorage;
 let muted=false,failOscillator=false;
 const sources=[],nodes=[];
 const parameter=()=>({value:0,setValueAtTime(){},linearRampToValueAtTime(){},exponentialRampToValueAtTime(){}});
 const node=()=>{const n={disconnected:0,connect(next){return next;},disconnect(){this.disconnected++;}};nodes.push(n);return n;};
 const source=()=>{const n=Object.assign(node(),{frequency:parameter(),starts:[],stops:[],start(at){this.starts.push(at);},stop(at){this.stops.push(at);}});sources.push(n);return n;};
 class Context {
  state='running';currentTime=10;sampleRate=100;destination={};
  createBuffer(_,frames){return{getChannelData:()=>new Float32Array(frames)};}
  createBufferSource(){return source();}
  createOscillator(){if(failOscillator)throw new Error('oscillator unavailable');return source();}
  createBiquadFilter(){return Object.assign(node(),{frequency:parameter(),Q:parameter()});}
  createGain(){return Object.assign(node(),{gain:parameter()});}
 }
 globalThis.window={AudioContext:Context};
 globalThis.localStorage={getItem:()=>muted?'off':'on',setItem:(_,v)=>{muted=v==='off';}};
 try{
  const {playMoment,setSound}=await import('../.runtime/frontend-test/sound.js?scene-lifecycle');
  for(const [kind,count] of [['glass-break',6],['door-breach',3],['explosion',2],['arrest',4],['raid',4],['killing',2],['gunfight',5],['robbery',1]]){
   const begin=sources.length;const cancel=playMoment(kind);assert.equal(sources.length-begin,count);
   assert.ok(sources.slice(begin).every(s=>s.starts[0]>=10.02));
   if(kind==='arrest')assert.ok(sources.at(-1).starts[0]>11,'future siren pulses are actually scheduled');
   cancel();cancel();
   assert.ok(sources.slice(begin).every(s=>s.stops.length===2&&s.stops[1]===undefined));
   assert.ok(nodes.every(n=>n.disconnected===1));
  }
  let begin=sources.length;
  const finished=playMoment('explosion');sources.slice(begin).forEach(s=>s.onended());finished();
  assert.ok(sources.slice(begin).every(s=>s.stops.length===1),'natural ending must not stop a source twice');
  begin=sources.length;
  const blast=playMoment('explosion');playMoment('arrest');
  blast();assert.equal(sources.at(-1).disconnected,0,'cancelling one scene must preserve another');
  setSound(false);assert.ok(nodes.every(n=>n.disconnected===1));
  assert.equal(playMoment('raid'),undefined);assert.equal(sources.length,begin+6);
  setSound(true);assert.equal(sources.length,begin+6,'unmuting must not replay cancelled sound');
  failOscillator=true;begin=sources.length;
  assert.equal(playMoment('explosion'),undefined);
  assert.equal(sources.length,begin+1);assert.equal(sources.at(-1).stops.length,2,'partial audio failure cancels the blast already started');
  assert.ok(nodes.every(n=>n.disconnected===1));
 }finally{globalThis.window=previousWindow;globalThis.localStorage=previousStorage;}
});

test('provided gun samples preload once, retain full tails, select the recorded weapon and cancel on mute',async()=>{
 const previousWindow=globalThis.window,previousStorage=globalThis.localStorage,previousFetch=globalThis.fetch;
 let muted=false,fetches=0;const sources=[],gains=[];
 const node=()=>({disconnected:0,connect(next){return next;},disconnect(){this.disconnected++;}});
 class Context {
  state='running';currentTime=10;destination={};
  decodeAudioData(data){return Promise.resolve({duration:new Uint8Array(data)[0]/10});}
  createBufferSource(){const s=Object.assign(node(),{starts:[],stops:[],start(at){this.starts.push(at);},stop(at){this.stops.push(at);}});sources.push(s);return s;}
  createGain(){const g=Object.assign(node(),{gain:{value:0}});gains.push(g);return g;}
  createDynamicsCompressor(){return Object.assign(node(),{threshold:{},knee:{},ratio:{},attack:{},release:{}});}
 }
 globalThis.window={AudioContext:Context};
 globalThis.localStorage={getItem:()=>muted?'off':'on',setItem:(_,v)=>muted=v==='off'};
 globalThis.fetch=async url=>{fetches++;return{ok:true,arrayBuffer:async()=>new Uint8Array([url.includes('revolver')?12:url.includes('shotgun')?10:11]).buffer};};
 try{
  const sound=await import('../.runtime/frontend-test/sound.js?provided-samples');
  await Promise.all([sound.preloadCityGunshots(),sound.preloadCityGunshots()]);assert.equal(fetches,3);assert.equal(sources.length,0,'preload must not make a sound');
  for(const [name,duration,gain] of [['revolver',1.2,.75],['shotgun',1,.65],['thompson',1.1,.8]]){
   sound.playCityGunshot(name);const source=sources.at(-1);
   assert.equal(source.buffer.duration,duration);assert.equal(gains.at(-1).gain.value,gain);
   assert.ok(Math.abs(source.stops[0]-source.starts[0]-duration)<1e-10,'sample tail was truncated');
  }
  assert.equal(sound.cityGunshotStatus().fallbacks,0);assert.equal(sources.length,3);assert.equal(sound.cityGunshotStatus().active,3);
  assert.ok(sources.every(s=>s.disconnected===0),'a later shot cancelled an earlier tail');
  sound.setSound(false);assert.equal(sound.cityGunshotStatus().active,0);assert.ok(sources.every(s=>s.disconnected===1));assert.ok(gains.every(g=>g.disconnected===1));
  assert.equal(sound.playCityGunshot('revolver'),undefined);
  sound.setSound(true);assert.equal(sources.length,3,'unmute replayed old shots');
  const stop=sound.playCityGunshot('thompson');sources.at(-1).onended();stop();
  assert.equal(sources.at(-1).stops.length,1,'completed sample was stopped twice');
 }finally{globalThis.window=previousWindow;globalThis.localStorage=previousStorage;globalThis.fetch=previousFetch;}
});

test('recorded effects preserve tails, crossfade stereo loops, and release ambient voices on mute',async()=>{
 const old={window:globalThis.window,localStorage:globalThis.localStorage,fetch:globalThis.fetch};
 let muted=false,fetches=0;const sources=[];
 const buffer=(channels,length,rate)=>{const data=Array.from({length:channels},(_,ch)=>Float32Array.from({length},(_,i)=>(i/length-.5)*(ch+1)));return{numberOfChannels:channels,length,sampleRate:rate,duration:length/rate,getChannelData:ch=>data[ch]};};
 const node=()=>({connect(next){return next;},disconnect(){this.disconnected=true;}});
 class Context {
  state='running';currentTime=3;destination={};
  createBuffer=buffer;
  decodeAudioData(){return Promise.resolve(buffer(2,1400,100));}
  createGain(){return Object.assign(node(),{gain:{value:0}});}
  createBufferSource(){const s=Object.assign(node(),{stops:[],start(at){this.at=at;},stop(at){this.stops.push(at);}});sources.push(s);return s;}
 }
 globalThis.window={AudioContext:Context};globalThis.localStorage={getItem:()=>muted?'off':'on',setItem:(_,v)=>muted=v==='off'};
 globalThis.fetch=async()=>{fetches++;return{ok:true,arrayBuffer:async()=>new ArrayBuffer(1)};};
 try{
  const sound=await import('../.runtime/frontend-test/sound.js?recorded-effects');
  await Promise.all([sound.preloadCityEffects(),sound.preloadCityEffects()]);assert.equal(fetches,10);assert.equal(sources.length,0);
  sound.playMoment('explosion');assert.equal(sources.length,1);assert.equal(sources[0].stops[0]-sources[0].at,14);
  const ambient=new sound.CityAmbientAudio();ambient.update(['fire','fire','engine-idle']);ambient.update(['fire','engine-idle']);assert.equal(sources.length,3);
  assert.ok(sources.slice(1).every(s=>s.loop&&s.stops.length===0));
  const original=buffer(2,100,100),copy=original.getChannelData(0).slice();const loop=sound.seamlessLoop(new Context(),original);
  assert.equal(loop.length,95);assert.deepEqual(original.getChannelData(0),copy);
  for(let ch=0;ch<2;ch++){const data=loop.getChannelData(ch);assert.ok(Math.abs(data.at(-1)-data[0])<.021);assert.ok(data.every(Number.isFinite));}
  sound.setSound(false);assert.deepEqual(sound.cityEffectStatus().active,[]);
  ambient.update(['fire']);sound.setSound(true);ambient.update(['fire']);assert.equal(sources.length,4);
  ambient.dispose();ambient.dispose();assert.deepEqual(sound.cityEffectStatus().active,[]);assert.ok(sources.every(s=>s.disconnected));
 }finally{Object.assign(globalThis,old);}
});
