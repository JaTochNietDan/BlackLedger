import test from 'node:test';
import assert from 'node:assert/strict';
import {StreetTraffic,trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';
const line=[{x:0,z:0},{x:200,z:0}];
function clear(requests,poses){
 for(const a of requests)for(const b of requests){if(a.id>=b.id)continue;
 const ap=poses.get(a.id),bp=poses.get(b.id);
 if(!ap.waiting&&!bp.waiting)assert.equal(trafficOverlap(ap.pose,a.model,bp.pose,b.model),false,`${a.id} intersects ${b.id}`);
 }
}
test('twelve cars sharing committed progress queue without overlap or backwards movement',()=>{
 const traffic=new StreetTraffic();
 const requests=Array.from({length:12},(_,i)=>({id:`car-${i}`,model:i%2?'packard':'ford',points:line,progress:.7}));
 let previous=traffic.update(requests,1/60);clear(requests,previous);
 assert.equal([...previous.values()].filter(p=>p.waiting).length,0);
 for(let frame=0;frame<180;frame++){
  for(const r of requests)r.progress=.7+.2*Math.min(1,frame/100);
  const poses=traffic.update(requests,1/60);clear(requests,poses);
  for(const r of requests){assert.ok(poses.get(r.id).progress>=previous.get(r.id).progress);assert.ok(poses.get(r.id).progress<=r.progress+1e-8);}
  previous=poses;
 }
});
test('crowded departures remain represented in the waiting set and enter when space opens',()=>{
 const traffic=new StreetTraffic();const requests=[0,1].map(i=>({id:`car-${i}`,model:'ford',points:line,progress:0}));
 let poses=traffic.update(requests,1/60);assert.equal([...poses.values()].filter(p=>p.waiting).length,1);
 for(let frame=0;frame<90;frame++){requests.forEach(r=>r.progress=.5);poses=traffic.update(requests,1/60);clear(requests,poses);}
 assert.ok([...poses.values()].every(p=>!p.waiting&&p.progress>0));
});
test('opposing lanes do not block each other',()=>{
 const traffic=new StreetTraffic();const requests=[{id:'east',model:'packard',points:[{x:0,z:1.6},{x:100,z:1.6}],progress:.5},{id:'west',model:'packard',points:[{x:100,z:-1.6},{x:0,z:-1.6}],progress:.5}];
 const poses=traffic.update(requests,1/60);clear(requests,poses);assert.ok([...poses.values()].every(p=>p.progress===.5&&!p.waiting));
});
test('crossing traffic cannot tunnel through occupied vehicles between frames',()=>{
 const traffic=new StreetTraffic();const requests=[{id:'east',model:'ford',points:[{x:0,z:0},{x:100,z:0}],progress:.4},{id:'north',model:'packard',points:[{x:50,z:-50},{x:50,z:50}],progress:.3}];
 traffic.update(requests,1/60);
 for(let frame=0;frame<450;frame++){requests.forEach(r=>r.progress=.9);clear(requests,traffic.update(requests,1/30));}
 for(const p of traffic.update(requests,1/30).values())assert.ok(p.progress>.89,'crossing traffic must finish');
});
test('a newly restarted journey on the same route never inherits old progress',()=>{
 const traffic=new StreetTraffic();const request={id:'returning',model:'ford',points:line,progress:.8};
 traffic.update([request],1/60);request.progress=.1;
 assert.equal(traffic.update([request],1/60).get(request.id).progress,.1);
 assert.equal(traffic.update([],1/60).size,0);
});
test('occupancy dimensions enclose the authored car meshes',async()=>{
 const {readFileSync}=await import('node:fs');const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
 for(const model of ['ford','hudson','packard','police']){
  const [min,max]=manifest[model].bounds_blender,size=trafficSize(model);
  assert.ok(size.width>=max[0]-min[0],`${model} width`);assert.ok(size.length>=max[1]-min[1],`${model} length`);
 }
});
test('long journeys remain in transit past the old 2.4-second cutoff',()=>{
 const traffic=new StreetTraffic();const request={id:'player',model:'hudson',points:[{x:0,z:0},{x:300,z:0}],progress:0};
 traffic.update([request],1/60);let pose;
 for(let frame=1;frame<=144;frame++){request.progress=frame/144;pose=traffic.update([request],1/60).get('player');}
 assert.ok(pose.progress<1,'completion must use actual rendered arrival');
 for(let frame=0;frame<2400;frame++)pose=traffic.update([request],1/60).get('player');
 assert.equal(pose.progress,1);
});
test('four-way arrivals clear the junction at different frame rates without deadlocking',()=>{
 for(const fps of [30,60,144]){
  const requests=[{id:'east',points:[{x:-50,z:1.6},{x:50,z:1.6}]},{id:'west',points:[{x:50,z:-1.6},{x:-50,z:-1.6}]},{id:'north',points:[{x:1.6,z:50},{x:1.6,z:-50}]},{id:'south',points:[{x:-1.6,z:-50},{x:-1.6,z:50}]}].map(r=>({...r,model:'packard',progress:.3}));
  const traffic=new StreetTraffic();traffic.update(requests,1/fps);let poses;
  for(let frame=0;frame<fps*20;frame++){requests.forEach(r=>r.progress=.9);poses=traffic.update(requests,1/fps);clear(requests,poses);}
  for(const [id,p]of poses)assert.ok(p.progress>=.89,`${id} deadlocked at ${fps} FPS`);
 }
});
test('pedestrian occupancy encloses the exported animated stride',async()=>{
 const {readFileSync}=await import('node:fs');const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
 const [min,max]=manifest.person.motion_bounds_blender,size=trafficSize('person');
 assert.ok(size.width/2>=Math.max(-min[0],max[0]));assert.ok(size.length/2>=Math.max(-min[1],max[1]));
 assert.ok(min[2]+.2+.05>=.17,'walking toes clear the pavement top');
});
test('a walking player can pass a snapshot-held opposing pedestrian',async()=>{
 const {readFileSync}=await import('node:fs');const {cityPlan,route}=await import('../.runtime/frontend-test/city3dPlan.js');
 const lots=cityPlan(JSON.parse(readFileSync(new URL('../core/locations.json',import.meta.url)))).lots;
 const at=id=>lots.find(l=>l.id===id);
 const requests=[{id:'npc',model:'person',points:route(at('room'),at('apartment')),progress:.25},{id:'player',model:'person',points:route(at('bar'),at('room')),progress:0}];
 const traffic=new StreetTraffic();traffic.update(requests,1/60);let poses;
 for(let frame=0;frame<3600;frame++){requests[1].progress=Math.min(1,frame/144);poses=traffic.update(requests,1/60);clear(requests,poses);}
 assert.equal(poses.get('player').progress,1);assert.equal(poses.get('npc').progress,.25);
});
test('walking and driving respect physical speed caps at normal and accelerated playback',()=>{
 for(const [model,speed]of [['person',1.8],['hudson',11]])for(const rate of [1,4])for(const fps of [30,144]){
  const traffic=new StreetTraffic(),request={id:'traveller',model,points:line,progress:0};
  traffic.update([request],1/fps,rate);request.progress=1;let placement;
  for(let frame=0;frame<fps;frame++)placement=traffic.update([request],1/fps,rate).get(request.id);
  assert.ok(Math.abs(placement.pose.x-speed*rate)<1e-6,`${model}, ${rate}x, ${fps} FPS`);
 }
});
test('event reservations stop later arrivals and release them after playback',()=>{
 const traffic=new StreetTraffic();
 const scene={id:'scene:kill',model:'casualty',points:[{x:10,z:0}],progress:0};
 const walker={id:'npc',model:'person',points:line,progress:0};
 traffic.update([scene,walker],1/60);walker.progress=1;
 let poses;
 for(let frame=0;frame<600;frame++){
  poses=traffic.update([scene,walker],1/60);clear([scene,walker],poses);
 }
 assert.ok(poses.get('npc').pose.x<9,'walker yields to the complete fall envelope');
 for(let frame=0;frame<600;frame++)poses=traffic.update([walker],1/60);
 assert.ok(poses.get('npc').pose.x>15,'scene removal releases the waiting traveller');
});

test('wheel rotation advances by actual distance and is independent of frame rate',async()=>{
 const {advanceWheel,WHEEL_RADIUS}=await import('../.runtime/frontend-test/city3dTraffic.js');
 assert.ok(Math.abs(advanceWheel(0,2*Math.PI*WHEEL_RADIUS))<1e-12);
 for(const fps of [30,60,144]){
  let angle=0;
  for(let frame=0;frame<fps*3;frame++)angle=advanceWheel(angle,11/fps);
  assert.ok(Math.abs(angle-advanceWheel(0,33))<1e-10);
  assert.equal(advanceWheel(angle,0),angle,'waiting traffic must not spin its wheels');
 }
});

test('front-wheel steering follows turn direction, stays bounded and eases by travelled distance',async()=>{
 const {wheelSteering,advanceSteering}=await import('../.runtime/frontend-test/city3dTraffic.js');
 for(const model of ['ford','hudson','packard','police']){
  assert.equal(wheelSteering([{x:0,z:0},{x:0,z:100}],.5,model),0);
  assert.equal(wheelSteering([{x:0,z:0}],.5,model),0);
  for(const side of [-1,1]){
   const curve=Array.from({length:81},(_,i)=>({x:side*10*(1-Math.cos(i*Math.PI/160)),z:10*Math.sin(i*Math.PI/160)}));
   const angle=wheelSteering(curve,.5,model);
   assert.ok(angle*side>.2&&Math.abs(angle)<=.5);
  }
 }
 for(const fps of [30,60,144]){
  let angle=0;for(let i=0;i<fps;i++)angle=advanceSteering(angle,.5,3/fps);
  assert.ok(Math.abs(angle-advanceSteering(0,.5,3))<1e-10);
  assert.equal(advanceSteering(angle,0,0),angle);
 }
});

test('parked vehicles reserve straight-wheel width while moving vehicles reserve the steering sweep',async()=>{
 const {trafficModel,trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 for(const model of ['ford','hudson','packard','police']){
  assert.equal(trafficModel(model,false),model);
  assert.equal(trafficSize(trafficModel(model,true)).width,2.15);
  assert.equal(trafficSize(trafficModel(model,false)).width,2.35);
  assert.equal(trafficSize(trafficModel(model,true)).length,trafficSize(model).length);
 }
 assert.equal(trafficModel('person',true),'person');
});

test('front tyres share a turn centre with a tighter inside angle in either direction',async()=>{
 const {frontWheelSteering}=await import('../.runtime/frontend-test/city3dTraffic.js');
 for(const [model,length] of Object.entries({ford:4.7,hudson:5.1,packard:5.8,police:4.7})){
  const base=(length-.2)*.59;
  for(const angle of [-.5,-.3,-.01,.01,.3,.5]){
   const left=frontWheelSteering(angle,model,-.87),right=frontWheelSteering(angle,model,.87);
   assert.ok(Math.abs(left)<=.50000001&&Math.abs(right)<=.50000001);
   assert.ok(Math.sign(left)===Math.sign(angle)&&Math.sign(right)===Math.sign(angle));
   assert.ok(angle>0?right>left:Math.abs(left)>Math.abs(right),'inside wheel must turn more sharply');
   assert.ok(Math.abs((-.87+base/Math.tan(left))-(.87+base/Math.tan(right)))<1e-9,'front wheel axes must meet at one rear-axle turn centre');
   assert.ok(Math.abs(left+frontWheelSteering(-angle,model,.87))<1e-12,'opposite turns must mirror');
  }
  assert.equal(frontWheelSteering(0,model,-.87),0);
  assert.equal(frontWheelSteering(0,model,.87),0);
 }
});

test('pending scenes let existing traffic leave but stop new arrivals entering',()=>{
 const traffic=new StreetTraffic(),points=[{x:40,z:16},{x:70,z:16}];
 const occupied={pose:{x:50,z:16,heading:0},model:'planter'};
 const inside={id:'inside',model:'person',points,progress:.3};
 traffic.update([inside],1/60);
 const behind={id:'behind',model:'person',points,progress:0};
 traffic.update([inside,behind],1/60,1,[occupied]);
 inside.progress=1;behind.progress=1;
 let poses;
 for(let i=0;i<900;i++){
  poses=traffic.update([inside,behind],1/60,1,[occupied]);
  clear([inside,behind],poses);
  const arrival=poses.get('behind');
  if(!arrival.waiting)assert.equal(trafficOverlap(arrival.pose,'person',occupied.pose,occupied.model),false);
 }
 assert.ok(poses.get('inside').progress>.8,'occupant could not leave');
 assert.ok(poses.get('behind').progress<.3,'new arrival crossed pending scene');
 for(let i=0;i<1200;i++)poses=traffic.update([inside,behind],1/60);
 assert.ok(poses.get('behind').progress>.9,'traffic failed to resume after scene release');
});

test('stationary pavement pedestrian walks aside continuously and retains their stance after release',()=>{
 for(const fps of [30,60,144]){
  const traffic=new StreetTraffic();
  const person={id:'player',model:'person',points:[{x:48,z:4.65}],progress:0};
  const space={pose:{x:48,z:6,heading:0},model:'planter'};
  traffic.update([person],0);
  let last=traffic.placement('player').pose, moved=false;
  for(let i=0;i<fps*5;i++){
   const p=traffic.update([person],1/fps,1,[space]).get('player');
   assert.equal(p.waiting,false);
   assert.equal(p.progress,0);
   assert.ok(Math.hypot(p.pose.x-last.x,p.pose.z-last.z)<=1.8/fps+1e-8);
   assert.equal(p.pose.z,4.65);
   moved ||= !!p.yielding;last=p.pose;
  }
  assert.ok(moved);
  assert.equal(trafficOverlap(last,'person',space.pose,space.model),false);
  assert.deepEqual(traffic.update([person],1/fps).get('player').pose,last);
  assert.deepEqual(person.points,[{x:48,z:4.65}]);
 }
});

test('stationary exit recovery respects occupied pavement and leaves parked cars in place',()=>{
 const traffic=new StreetTraffic();
 const person={id:'player',model:'person',points:[{x:48,z:4.65}],progress:0};
 const blocker={id:'blocker',model:'parked-ford',points:[{x:44,z:4.65}],progress:0};
 const space={pose:{x:48,z:6,heading:0},model:'planter'};
 let poses=traffic.update([blocker,person],0);
 const car=poses.get('blocker').pose;
 for(let i=0;i<300;i++){
  poses=traffic.update([blocker,person],1/60,1,[space]);
  clear([blocker,person],poses);
  assert.deepEqual(poses.get('blocker').pose,car);
 }
 assert.ok(poses.get('player').pose.x>48,'person must choose the unoccupied side');
 assert.equal(trafficOverlap(poses.get('player').pose,'person',space.pose,space.model),false);
});

test('a new journey rejoins its start from the visible sidestep without teleporting',()=>{
 for(const fps of [30,60,144]){
  const traffic=new StreetTraffic();
  const request={id:'player',model:'person',points:[{x:48,z:4.65}],progress:0};
  const space={pose:{x:48,z:6,heading:0},model:'planter'};
  for(let i=0;i<fps*3;i++)traffic.update([request],1/fps,1,[space]);
  let last=traffic.placement('player').pose;
  assert.ok(last.x<46);
  request.points=[{x:48,z:4.65},{x:60,z:4.65}];request.progress=1;
  let reachedStart=false;
  for(let i=0;i<fps*12;i++){
   const p=traffic.update([request],1/fps).get('player');
   assert.equal(p.waiting,false);
   assert.ok(Math.hypot(p.pose.x-last.x,p.pose.z-last.z)<=1.8/fps+1e-8,'handoff teleported');
   if(p.pose.x<48-1e-8)assert.equal(p.progress,0);
   if(Math.abs(p.pose.x-48)<1e-8)reachedStart=true;
   last=p.pose;
  }
  assert.ok(reachedStart);
  assert.ok(last.x>59);
 }
});

test('the connector waits for an occupied route start and resumes after clearance',()=>{
 const traffic=new StreetTraffic();
 const person={id:'player',model:'person',points:[{x:48,z:4.65}],progress:0};
 const space={pose:{x:48,z:6,heading:0},model:'planter'};
 for(let i=0;i<180;i++)traffic.update([person],1/60,1,[space]);
 const obstacle={id:'blocker',model:'person',points:[{x:48,z:4.65}],progress:0};
 traffic.update([person,obstacle],1/60);
 person.points=[{x:48,z:4.65},{x:60,z:4.65}];person.progress=1;
 for(let i=0;i<180;i++){
  const poses=traffic.update([person,obstacle],1/60);
  clear([person,obstacle],poses);
  assert.equal(poses.get('player').progress,0);
 }
 for(let i=0;i<600;i++)traffic.update([person],1/60);
 assert.ok(traffic.placement('player').progress>.9);
});

test('a stationary person clears a waiting pedestrian departure without scene reservations',()=>{
 for(const fps of [30,60,144]){
  const traffic=new StreetTraffic();
  const player={id:'player',model:'person',points:[{x:48,z:4.65}],progress:0};
  traffic.update([player],0);
  const departing={id:'departing',model:'person',points:[{x:48,z:4.65},{x:60,z:4.65}],progress:0};
  let poses;
  for(let i=0;i<fps*3;i++){
   poses=traffic.update([player,departing],1/fps);
   clear([player,departing],poses);
   assert.equal(poses.get('player').waiting,false);
  }
  assert.equal(poses.get('departing').waiting,false);
  assert.ok(poses.get('player').pose.x<48);
  assert.equal(poses.get('departing').progress,0,'clearance must not advance a committed journey');
  departing.progress=1;
  for(let i=0;i<fps*8;i++){
   poses=traffic.update([player,departing],1/fps);clear([player,departing],poses);
  }
  assert.ok(poses.get('departing').progress>.9);
 }
});
