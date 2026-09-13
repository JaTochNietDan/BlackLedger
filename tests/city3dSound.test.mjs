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
