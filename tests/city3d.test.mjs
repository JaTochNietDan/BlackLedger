import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
// Node's ES module loader needs the explicit extension on the emitted import.
const places=JSON.parse(readFileSync(new URL('../core/locations.json',import.meta.url)));
const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
const {cityPlan,route,onRoute,intersectsLot,MODEL_LIMIT,PITCH}=await import('../.runtime/frontend-test/city3dPlan.js');
const plan=cityPlan(places);
test('all exported building geometry, including escapes and cornices, stays inside the reserved footprint',()=>{
 for(const lot of plan.lots){const m=manifest[lot.model];assert.ok(m,lot.model);const [a,b]=m.bounds_blender;for(const axis of [0,1]){assert.ok(a[axis]>=-MODEL_LIMIT/2,`${lot.model} min ${axis}`);assert.ok(b[axis]<=MODEL_LIMIT/2,`${lot.model} max ${axis}`);}}
});
test('each current building has a separate reserved footprint',()=>{
 for(const a of plan.lots)for(const b of plan.lots)if(a.id!==b.id)assert.ok(Math.abs(a.x-b.x)>=MODEL_LIMIT||Math.abs(a.z-b.z)>=MODEL_LIMIT,`${a.id}/${b.id}`);
});
test('every pedestrian and vehicle route avoids every building with actor clearance',()=>{
 for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){const points=route(from,to,driving);for(let step=0;step<=200;step++){const p=onRoute(points,step/200);for(const building of plan.lots)assert.equal(intersectsLot(p,building,driving?2.9:.7),false,`${from.id} to ${to.id} intersects ${building.id} at ${JSON.stringify(p)}`);}}
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
   assert.equal(Math.sign(a.z-from.row*PITCH),Math.sign(to.x-from.x));
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
test('authored street furniture stays on pavements and outside all buildings and travel clearances',async()=>{
 const {streetsidePosition,PITCH}=await import('../.runtime/frontend-test/city3dPlan.js');
 const [min,max]=manifest.streetside.bounds_blender;
 const boxes=plan.lots.map(lot=>{
  const p=streetsidePosition(lot),box={x0:p.x+min[0],x1:p.x+max[0],z0:p.z-max[1],z1:p.z-min[1]};
  assert.ok(box.x0>=lot.col*PITCH+4&&box.x1<=(lot.col+1)*PITCH-4);
  assert.ok(box.z0>=lot.row*PITCH+4&&box.z1<=(lot.row+1)*PITCH-4);
  for(const b of plan.lots)assert.ok(box.x1<=b.x-8.5||box.x0>=b.x+8.5||box.z1<=b.z-8.5||box.z0>=b.z+8.5,`furniture overlaps ${b.id}`);
  return box;
 });
 for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){
  const points=route(from,to,driving),padding=driving?2.9:.7;
  for(let step=0;step<=200;step++){
   const p=onRoute(points,step/200);
   for(const box of boxes)assert.ok(p.x<=box.x0-padding||p.x>=box.x1+padding||p.z<=box.z0-padding||p.z>=box.z1+padding,`${from.id}/${to.id}: street furniture obstructs route`);
  }
 }
});
test('parking positions clear buildings, furniture and every travel lane',async()=>{
 const {parkingSpot,streetsidePosition}=await import('../.runtime/frontend-test/city3dPlan.js');
 const {trafficOverlap}=await import('../.runtime/frontend-test/city3dTraffic.js');
 for(const lot of plan.lots){const parked={...parkingSpot(lot),heading:0};
  for(const b of plan.lots)assert.ok(Math.abs(parked.x-b.x)>=8.5+1.075||Math.abs(parked.z-b.z)>=8.5+2.9);
  for(const b of plan.lots){const f=streetsidePosition(b);assert.ok(Math.abs(parked.z-f.z)>3.3||Math.abs(parked.x-f.x)>5.3);}
  for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){
   const points=route(from,to,driving);
   for(let step=0;step<=40;step++)assert.equal(trafficOverlap(parked,'parked-packard',onRoute(points,step/40),driving?'packard':'person'),false,`parked at ${lot.id} blocks ${from.id}/${to.id}`);
  }
 }
});
test('lamp poles stay clear of the articulated pedestrian walking envelope',async()=>{
 const {lampPositions}=await import('../.runtime/frontend-test/city3dPlan.js');
 const lamps=plan.lots.flatMap(lampPositions);
 for(const from of plan.lots)for(const to of plan.lots){const points=route(from,to);
  for(let step=0;step<=100;step++){const p=onRoute(points,step/100);
   for(const lamp of lamps){const dx=lamp.x-p.x,dz=lamp.z-p.z;
    const side=dx*Math.cos(p.heading)-dz*Math.sin(p.heading),forward=dx*Math.sin(p.heading)+dz*Math.cos(p.heading);
    assert.ok(Math.abs(side)>=.85/2+.16||Math.abs(forward)>=1.4/2+.16,`${from.id}/${to.id} walks through a lamp`);
   }
  }
 }
});

test('event bays clear buildings, lamps, furniture and every ordinary route',async()=>{
 const {sceneSlots}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {trafficSize,trafficOverlap}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const {lampPositions,streetsidePosition}=await import('../.runtime/frontend-test/city3dPlan.js');
 const slots=plan.lots.flatMap(lot=>['killing','raid'].flatMap(kind=>sceneSlots(lot,kind)));
 for(const slot of slots){
  const size=trafficSize(slot.model);
  for(const lot of plan.lots){
   assert.ok(Math.abs(slot.pose.x-lot.x)>=MODEL_LIMIT/2+size.width/2||Math.abs(slot.pose.z-lot.z)>=MODEL_LIMIT/2+size.length/2,`building ${lot.id}`);
   for(const lamp of lampPositions(lot))assert.ok(Math.abs(slot.pose.x-lamp.x)>size.width/2+.16||Math.abs(slot.pose.z-lamp.z)>size.length/2+.16,'lamp');
   const furniture=streetsidePosition(lot),[min,max]=manifest.streetside.bounds_blender;
   assert.ok(slot.pose.x+size.width/2<furniture.x+min[0]||slot.pose.x-size.width/2>furniture.x+max[0]||slot.pose.z+size.length/2<furniture.z-max[1]||slot.pose.z-size.length/2>furniture.z-min[1],'furniture');
  }
 }
 for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){
  const path=route(from,to,driving);
  for(let i=0;i<=150;i++){
   const pose=onRoute(path,i/150);
   for(const slot of slots){
    // Circumscribed circles reject distant bays without an expensive SAT check.
    const a=trafficSize(driving?'packard':'person'),b=trafficSize(slot.model);
    const radius=(Math.hypot(a.width,a.length)+Math.hypot(b.width,b.length))/2+.2;
    if(Math.abs(pose.x-slot.pose.x)>radius||Math.abs(pose.z-slot.pose.z)>radius)continue;
    assert.equal(trafficOverlap(pose,driving?'packard':'person',slot.pose,slot.model),false,`${from.id}/${to.id}: ${JSON.stringify(slot)}`);
   }
  }
 }
});
test('event overflow waits for a clear slot and respects a parked player car',async()=>{
 const {sceneSlots,availableSceneSlot}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {parkingSpot}=await import('../.runtime/frontend-test/city3dPlan.js');
 const lot=plan.lots[0];
 for(const kind of ['killing','arrest']){
  const occupied=kind==='arrest'?[{model:'parked-packard',pose:{...parkingSpot(lot),heading:0}}]:[];
  const initial=occupied.length;
  for(let i=0;i<sceneSlots(lot,kind).length-initial;i++){
   const slot=availableSceneSlot(lot,kind,occupied);assert.ok(slot);occupied.push(slot);
  }
  assert.equal(availableSceneSlot(lot,kind,occupied),undefined);
  const released=occupied.pop();assert.deepEqual(availableSceneSlot(lot,kind,occupied),released);
 }
});
test('the complete casualty fall remains above pavement and inside its reservation',async()=>{
 const {casualtyFall,sceneSlots}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const slot=sceneSlots(plan.lots[0],'killing')[0],size=trafficSize(slot.model);
 const [min,max]=manifest.person.bounds_blender;
 for(let frame=0;frame<=120;frame++){
  const {rotation,height}=casualtyFall(frame/120);
  for(const x of [min[0],max[0]])for(const y of [min[2],max[2]])for(const z of [-max[1],-min[1]]){
   const wx=slot.root.x+x*Math.cos(rotation)-y*Math.sin(rotation);
   const wy=height+x*Math.sin(rotation)+y*Math.cos(rotation);
   assert.ok(wy>=.17,`ground at ${frame}: ${wy}`);
   assert.ok(Math.abs(wx-slot.pose.x)<=size.width/2,`width at ${frame}`);
   assert.ok(Math.abs(z)<=size.length/2);
  }
 }
});

test('unused grid cells receive scenery exactly once without claiming playable addresses',()=>{
 const occupied=new Set(plan.lots.map(l=>`${l.col}:${l.row}`));
 const vacant=new Set(plan.vacant.map(l=>`${l.col}:${l.row}`));
 assert.equal(vacant.size,plan.vacant.length);
 assert.equal(vacant.size+occupied.size,plan.cols*plan.rows);
 for(const key of vacant)assert.equal(occupied.has(key),false);
 assert.ok(vacant.size>0,'current map contains unused parcels');
 const [min,max]=manifest['vacant-lot'].bounds_blender;
 for(const axis of [0,1])assert.ok(min[axis]>=-8.5&&max[axis]<=8.5);
 for(const driving of [false,true])for(const from of plan.lots)for(const to of plan.lots){
  const points=route(from,to,driving),padding=driving?2.9:.7;
  for(let i=0;i<=200;i++){
   const at=onRoute(points,i/200);
   for(const lot of plan.vacant)assert.ok(at.x<=lot.x+min[0]-padding||at.x>=lot.x+max[0]+padding||at.z<=lot.z-max[1]-padding||at.z>=lot.z-min[1]+padding,'vacant fencing must not obstruct a route');
  }
 }
});

test('Mercer Court occupies a vacant block without shifting existing addresses',()=>{
 const previous=cityPlan(places.filter(p=>p.id!=='mercercourt'));
 for(const old of previous.lots){const current=plan.lots.find(p=>p.id===old.id);assert.equal(current.x,old.x,old.id);assert.equal(current.z,old.z,old.id);}
 const court=plan.lots.find(p=>p.id==='mercercourt');
 assert.equal(court.model,'mercer-court');
 assert.ok(previous.vacant.some(v=>v.col===court.col&&v.row===court.row));
 assert.equal(plan.cols,previous.cols);assert.equal(plan.rows,previous.rows);
});

test('scene clock preserves opening beats across stalls and hidden tabs',async()=>{
 const {ScenePlaybackClock}=await import('../.runtime/frontend-test/city3dEvents.js');
 const clock=new ScenePlaybackClock(1000);
 assert.equal(clock.step(1016,false),1016);
 assert.equal(clock.step(9000,true),1016);
 assert.equal(clock.step(20000,false),1116);
 assert.equal(clock.step(20016,false),1132);
});
