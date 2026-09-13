import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
// Node's ES module loader needs the explicit extension on the emitted import.
const places=JSON.parse(readFileSync(new URL('../core/locations.json',import.meta.url)));
const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
const {cityPlan,route,onRoute,intersectsLot,MODEL_LIMIT}=await import('../.runtime/frontend-test/city3dPlan.js');
const plan=cityPlan(places);
test('all exported building geometry, including escapes and cornices, stays inside the reserved footprint',()=>{
 for(const lot of plan.lots){const m=manifest[lot.model];assert.ok(m,lot.model);const [a,b]=m.bounds_blender;for(const axis of [0,1]){assert.ok(a[axis]>=-MODEL_LIMIT/2,`${lot.model} min ${axis}`);assert.ok(b[axis]<=MODEL_LIMIT/2,`${lot.model} max ${axis}`);}}
});
test('each current building has a separate reserved footprint',()=>{
 for(const a of plan.lots)for(const b of plan.lots)if(a.id!==b.id)assert.ok(Math.abs(a.x-b.x)>=MODEL_LIMIT||Math.abs(a.z-b.z)>=MODEL_LIMIT,`${a.id}/${b.id}`);
});
test('every pedestrian and vehicle route avoids every building with actor clearance',()=>{
 for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){const points=route(from,to,driving);for(let step=0;step<=200;step++){const p=onRoute(points,step/200);for(const building of plan.lots)assert.equal(intersectsLot(p,building,driving?2.9:.4),false,`${from.id} to ${to.id} intersects ${building.id} at ${JSON.stringify(p)}`);}}
});
test('route progress clamps and travels by path distance',()=>{const path=[{x:0,z:0},{x:10,z:0},{x:10,z:30}];assert.deepEqual(onRoute(path,.5),{x:10,z:10,heading:0});assert.equal(onRoute(path,-1).x,0);assert.equal(onRoute(path,2).z,30);});
const {CityCueQueue}=await import('../.runtime/frontend-test/city3dEvents.js');
const cue={id:'one',kind:'explosion',target:'bar',caption:'A committed explosion'};
test('old save cues stay silent; repeats and same-result revisions never replay',()=>{const q=new CityCueQueue();assert.deepEqual(q.take('city',[cue],null,true),[]);assert.deepEqual(q.take('city',[cue],null,false),[]);const next={...cue,id:'two'};assert.deepEqual(q.take('city',[next],next,false),[next]);assert.deepEqual(q.take('city',[next],next,false),[]);});
test('active playback survives mounting the city and IDs can be reused by a fresh world',()=>{const q=new CityCueQueue();assert.deepEqual(q.take('city',[cue],cue,true),[cue]);assert.deepEqual(q.take('new-city',[cue],cue,true),[cue]);});
test('entering the city for a committed result preserves its simultaneous explosion and casualty',()=>{const q=new CityCueQueue();const death={...cue,id:'death',kind:'killing'};assert.deepEqual(q.take('city',[cue,death],death,true),[cue,death]);assert.deepEqual(q.take('city',[cue,death],death,false),[]);});
test('explicit replay presents each saved cue once without enabling refresh replay',()=>{const q=new CityCueQueue();q.take('city',[cue],null,true);assert.deepEqual(q.take('city',[cue],cue,false,true),[cue]);assert.deepEqual(q.take('city',[cue],cue,false),[]);});
test('opposing cars use separate right-hand lanes without a same-street detour',()=>{
 for(const from of plan.lots)for(const to of plan.lots){
  if(from.id===to.id||from.row!==to.row)continue;
  const forward=route(from,to,true),reverse=route(to,from,true);
  assert.equal(forward.length,2);
  for(let step=0;step<=100;step++){
   const a=onRoute(forward,step/100),b=onRoute(reverse,1-step/100);
   assert.ok(Math.abs(a.x-b.x)<1e-8);
   assert.ok(Math.abs(a.z-b.z)>3.1);
   assert.equal(Math.sign(a.z-from.row*28),Math.sign(to.x-from.x));
  }
 }
});
test('vehicle junction curves have continuous heading and retain building clearance',()=>{
 for(const from of plan.lots)for(const to of plan.lots){
  const points=route(from,to,true);let last;
  for(let i=1;i<points.length;i++){
   const a=points[i-1],b=points[i],heading=Math.atan2(b.x-a.x,b.z-a.z);
   if(last!==undefined){const delta=Math.atan2(Math.sin(heading-last),Math.cos(heading-last));assert.ok(Math.abs(delta)<.18,`${from.id}/${to.id}: abrupt steering ${delta}`);}
   last=heading;
   for(const lot of plan.lots)assert.equal(intersectsLot(b,lot,2.9),false,`${from.id}/${to.id}: curve clips ${lot.id}`);
  }
 }
});
