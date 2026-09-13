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
