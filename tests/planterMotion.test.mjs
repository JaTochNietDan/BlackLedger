import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityPlanter,PLANTER_BLAST} from '../.runtime/frontend-test/city3dPlanter.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('planter exits the authored vestibule without crossing the door or masonry',async()=>{
 for(const name of ['monarch','tavern','tenement','mercer-court','casino','civic','shop','warehouse','bluehour','goldenlily','papermoon','mariner','undertaker','filling','docks','haulage','garage','dealer','villa'])for(const person of ['person','woman']){
  const building=await model(name),actor=await model(person),cast=new CityPlanter(actor);
  if(['filling','docks','haulage','garage','dealer'].includes(name))building.rotation.y=Math.PI;
  building.updateMatrixWorld(true);
  cast.root.position.copy(building.getObjectByName(name==='villa'?'entrance-landing':'entrance-threshold').getWorldPosition(new THREE.Vector3()));
  const door=building.getObjectByName('entrance-door-hinge'),triangle=new THREE.Triangle(),bounds=new THREE.Box3();
  for(let frame=0;frame<=(name==='villa'?75:186);frame++){
   const pose=cast.update(frame/30);door.rotation.y=-Math.PI/2*pose.door;building.updateMatrixWorld(true);
   const actorBounds=new THREE.Box3().setFromObject(actor,true);if(name!=='undertaker')actorBounds.min.y+=.12;
   building.traverse(o=>{
    if(!(o instanceof THREE.Mesh))return;
    o.geometry.computeBoundingBox();bounds.copy(o.geometry.boundingBox).applyMatrix4(o.matrixWorld);
    if(!actorBounds.intersectsBox(bounds))return;
    const positions=o.geometry.attributes.position,index=o.geometry.index,count=index?.count??positions.count;
    for(let i=0;i<count;i+=3){
     triangle.a.fromBufferAttribute(positions,index?index.getX(i):i).applyMatrix4(o.matrixWorld);
     triangle.b.fromBufferAttribute(positions,index?index.getX(i+1):i+1).applyMatrix4(o.matrixWorld);
     triangle.c.fromBufferAttribute(positions,index?index.getX(i+2):i+2).applyMatrix4(o.matrixWorld);
     assert.equal(actorBounds.intersectsTriangle(triangle),false,`${name}/${person} intersects ${o.name} at ${frame/30}`);
    }
   });
  }
  if(name==='villa'){assert.equal(building.getObjectByName('entrance-threshold'),undefined,'stair exit enabled before choreography');continue;}
  const finish=cast.update(PLANTER_BLAST);assert.equal(finish.blast,true);assert.equal(finish.door,0);
  assert.ok(actor.position.z<-2&&actor.position.x<=-2);
  assert.equal(cast.update(PLANTER_BLAST-.001).blast,false);
 }
});

test('planter reservation covers both rigs through translated and rotated exits',async()=>{
 const {planterReservation,availableSceneSlot}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 for(const person of ['person','woman'])for(const heading of [0,Math.PI/2,Math.PI,1.1]){
  const cast=new CityPlanter(await model(person)),slot=planterReservation({x:30,z:50},heading),size=trafficSize(slot.model);
  cast.root.position.set(slot.root.x,.2,slot.root.z);cast.root.rotation.y=heading;
  const inverse=new THREE.Matrix4().compose(new THREE.Vector3(slot.pose.x,0,slot.pose.z),new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,1,0),heading),new THREE.Vector3(1,1,1)).invert();
  for(let frame=0;frame<=186;frame++){
   cast.update(frame/30);const bounds=new THREE.Box3();
   cast.actor.traverse(o=>{if(!(o instanceof THREE.Mesh))return;const p=o.geometry.attributes.position;for(let i=0;i<p.count;i++)bounds.expandByPoint(new THREE.Vector3().fromBufferAttribute(p,i).applyMatrix4(o.matrixWorld).applyMatrix4(inverse));});
   assert.ok(bounds.min.x>=-size.width/2&&bounds.max.x<=size.width/2,`${person} escapes reservation width`);
   assert.ok(bounds.min.z>=-size.length/2&&bounds.max.z<=size.length/2,`${person} escapes reservation depth`);
  }
 }
 const lot={x:0,z:16,row:0},entry={x:0,z:9},occupied=planterReservation(entry);
 assert.equal(availableSceneSlot(lot,'planter',[],undefined),undefined,'invented a doorway');
 assert.equal(availableSceneSlot(lot,'planter',[occupied],entry),undefined,'reserved an occupied exit');
 assert.deepEqual(availableSceneSlot(lot,'planter',[],entry),occupied);
});


test('planter soles stay on the pavement through walking, stopping and turning',async()=>{
 for(const person of ['person','woman'])for(const heading of [0,1.1]){
  const cast=new CityPlanter(await model(person));
  cast.root.position.set(30,.2,50);cast.root.rotation.y=heading;
  let previousY;
  for(let frame=0;frame<=816;frame++){
   cast.update(frame/120);let low=Infinity,shoes=0;
   cast.actor.traverse(o=>{
    if(!(o instanceof THREE.Mesh)||!o.name.startsWith('shoe'))return;
    shoes++;
    const p=o.geometry.attributes.position;
    for(let i=0;i<p.count;i++)low=Math.min(low,new THREE.Vector3().fromBufferAttribute(p,i).applyMatrix4(o.matrixWorld).y);
   });
   assert.equal(shoes,2);
   assert.ok(Math.abs(low-.205)<1e-6,`${person} sole height ${low} at ${frame/120}`);
   if(previousY!==undefined)assert.ok(Math.abs(cast.actor.position.y-previousY)<.012,'vertical pose snapped');
   previousY=cast.actor.position.y;
  }
 }
});


test('villa entrance landing and stair treads have consistent 20cm risers',async()=>{
 const b=await model('villa');b.updateMatrixWorld(true);
 const ray=new THREE.Raycaster();
 for(const [z,height] of [[-5.5,.6],[-6.2,.6],[-6.6,.4],[-7,.2],[-3,.6]]){
  ray.set(new THREE.Vector3(0,1,z),new THREE.Vector3(0,-1,0));
  const hit=ray.intersectObject(b,true)[0];
  assert.ok(hit);assert.ok(Math.abs(hit.point.y-height)<1e-5,`tread ${z}: ${hit.point.y}`);
 }
});


test('stair leg poses reach independent ankle targets without tilting the soles',async()=>{
 const {CityLegs}=await import('../.runtime/frontend-test/city3dLegs.js');
 for(const name of ['person','woman']){
  const actor=await model(name),legs=new CityLegs(actor);
  actor.position.set(30,.5,20);actor.rotation.y=1.1;
  for(const side of [-1,1])for(const y of [.14,.24,.34])for(const z of [-.2,0,.2]){
   const target=new THREE.Vector3(side*.12,y,z);
   assert.equal(legs.place(side,target),true);
   const ankle=actor.getObjectByName(`ankle${side}`);
   assert.ok(actor.worldToLocal(ankle.getWorldPosition(new THREE.Vector3())).distanceTo(target)<1e-6);
   const up=new THREE.Vector3(0,1,0).transformDirection(ankle.matrixWorld);
   assert.ok(up.distanceTo(new THREE.Vector3(0,1,0))<1e-6);
   let sole=Infinity;
   ankle.traverse(o=>{if(!(o instanceof THREE.Mesh))return;const vertices=o.geometry.attributes.position;for(let i=0;i<vertices.count;i++)sole=Math.min(sole,new THREE.Vector3().fromBufferAttribute(vertices,i).applyMatrix4(o.matrixWorld).y);});
   assert.ok(Math.abs(sole-(.5+y-.095))<1e-6,`authored sole height ${sole}`);
  }
  const before=actor.getObjectByName('leg1').quaternion.clone();
  assert.equal(legs.place(1,new THREE.Vector3(0,-10,0)),false);
  assert.ok(before.equals(actor.getObjectByName('leg1').quaternion));
 }
});


test('stair descent keeps support planted and actor geometry above the authored treads',async()=>{
 const {CityStairDescent}=await import('../.runtime/frontend-test/city3dStairs.js');
 const building=await model('villa');building.updateMatrixWorld(true);
 for(const name of ['person','woman']){
  const cast=new CityStairDescent(await model(name));cast.root.position.set(0,.6,-6.23);
  let last;
  for(let frame=0;frame<=Math.ceil(cast.duration*120);frame++){
   const t=frame/120,pose=cast.update(t);assert.ok(pose.reached.every(Boolean),`${name} unreachable at ${t}`);
   const support=1-pose.moving;
   if(last&&last.moving===pose.moving)assert.ok(last.feet[support].distanceTo(pose.feet[support])<1e-8,'support foot slides');
   cast.actor.traverse(o=>{
    if(!(o instanceof THREE.Mesh))return;
    const p=o.geometry.attributes.position;
    for(let i=0;i<p.count;i++){
     const v=new THREE.Vector3().fromBufferAttribute(p,i).applyMatrix4(o.matrixWorld);
     const surface=v.z> -6.4?.6:v.z> -6.8?.4:v.z> -7.2?.2:0;
     assert.ok(v.y>=surface-.00001,`${name} geometry penetrates tread at ${t}: ${v.toArray()}`);
    }
   });
   if(last)assert.ok(cast.actor.position.distanceTo(last.body)<.025,'body position snapped');
   last={...pose,body:cast.actor.position.clone()};
  }
  const finish=cast.update(cast.duration);assert.ok(finish.done);assert.ok(Math.abs(cast.actor.position.y+.605)<1e-6,'did not stand upright');assert.ok(finish.feet.every(f=>Math.abs(f.y+.5)<1e-6&&Math.abs(f.z+1.2)<1e-6));
 }
});

test('villa approach joins stair descent with no body jump or tread penetration',async()=>{
 const {CityVillaExit}=await import('../.runtime/frontend-test/city3dVillaExit.js');
 for(const name of ['person','woman']){
  const cast=new CityVillaExit(await model(name));cast.root.position.set(0,.6,-5.12);
  const {planterReservation}=await import('../.runtime/frontend-test/city3dEvents.js');
  const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
  const slot=planterReservation({x:0,z:-5.12}),size=trafficSize(slot.model);
  let previous;
  for(let frame=0;frame<=1548;frame++){
   const t=frame/120,p=cast.update(t);assert.ok(p.reached);
   cast.root.updateMatrixWorld(true);
   const body=cast.actor.getWorldPosition(new THREE.Vector3());
   if(previous)assert.ok(body.distanceTo(previous)<.025,`${name} body jumps at ${t}`);
   previous=body;
   cast.actor.traverse(o=>{
    if(!(o instanceof THREE.Mesh))return;const vertices=o.geometry.attributes.position;
    for(let i=0;i<vertices.count;i++){
     const v=new THREE.Vector3().fromBufferAttribute(vertices,i).applyMatrix4(o.matrixWorld);
     assert.ok(Math.abs(v.x-slot.pose.x)<=size.width/2&&Math.abs(v.z-slot.pose.z)<=size.length/2,'villa actor leaves reserved footprint');
     const surface=v.z> -6.4?.6:v.z> -6.8?.4:v.z> -7.2?.2:0;
     assert.ok(v.y>=surface-1e-5,`${name} intersects landing/tread at ${t}`);
    }
   });
  }
  assert.ok(cast.update(cast.duration).done);assert.equal(cast.update(cast.duration).door,0);
  // Replays can seek back across the rig handoff without retaining the stair pose.
  cast.update(0);assert.ok(cast.actor.getWorldPosition(new THREE.Vector3()).z>-3);
 }
});

test('accident fall keeps both actual rigs above ground and preserves the recorded fatal outcome',async()=>{
 const {CityAccident}=await import('../.runtime/frontend-test/city3dAccident.js');
 for(const name of ['person','woman'])for(const fatal of [false,true]){
  const cast=new CityAccident(await model(name),fatal);
  cast.root.position.set(30,.2,40);cast.root.rotation.y=1.1;
  let previous;
  for(let frame=0;frame<=480;frame++){
   const state=cast.update(frame/120),bounds=new THREE.Box3().setFromObject(cast.actor,true);
   assert.ok(Math.abs(bounds.min.y-.205)<1e-6,`${name} lost ground contact at ${frame/120}`);
   assert.equal(state.fatal,fatal);if(fatal)assert.equal(state.recovering,false);
   assert.ok(Number.isFinite(bounds.max.y));
   if(previous)assert.ok(cast.actor.position.distanceTo(previous)<.06,'discontinuous fall');
   previous=cast.actor.position.clone();
  }
  assert.ok(fatal?cast.actor.rotation.x<-1.5:cast.actor.rotation.x>-.8);
  assert.equal(cast.update(4).done,true);
  cast.update(0);assert.equal(cast.actor.rotation.x,0);
 }
});

test('surviving accident recovery plants both palms before raising the body',async()=>{
 const {CityAccident}=await import('../.runtime/frontend-test/city3dAccident.js');
 for(const name of ['person','woman']){
  const cast=new CityAccident(await model(name),false);cast.root.position.set(12,.2,9);cast.root.rotation.y=.7;
  for(let frame=240;frame<=480;frame++){
   cast.update(frame/120);
   for(const side of [-1,1]){
    const wrist=cast.actor.getObjectByName('accident-wrist'+side);
    const at=cast.root.worldToLocal(wrist.getWorldPosition(new THREE.Vector3()));
    assert.ok(at.distanceTo(new THREE.Vector3(side*.43,.084,-.6))<1e-5,`${name} palm slides at ${frame/120}: ${at.toArray()}`);
    const hand=wrist.children[0],bounds=new THREE.Box3().setFromObject(hand,true);
    assert.ok(bounds.min.y>=.2-1e-6,'palm penetrates floor');
   }
  }
  cast.update(0);const reset=cast.actor.getObjectByName('arm1').quaternion;
  assert.ok(reset.angleTo(new THREE.Quaternion())<1e-8,'replay retains arm brace');
 }
});

test('accident reservations enclose the whole posed cast and yield occupied frontage slots',async()=>{
 const {CityAccident}=await import('../.runtime/frontend-test/city3dAccident.js');
 const {accidentReservation,sceneSlots,availableSceneSlot}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const lot={id:'review',x:48,z:48,row:1,col:1,model:'shop'};
 const slots=sceneSlots(lot,'accident');
 assert.equal(slots.length,3);
 assert.notDeepEqual(availableSceneSlot(lot,'accident',[slots[0]]),slots[0]);
 assert.equal(availableSceneSlot(lot,'accident',slots),undefined);
 for(const name of ['person','woman'])for(const fatal of [false,true])for(const heading of [0,Math.PI/2,Math.PI,1.1]){
  const cast=new CityAccident(await model(name),fatal),slot=accidentReservation(slots[0].root,heading),size=trafficSize(slot.model);
  cast.root.position.set(slot.root.x,.2,slot.root.z);cast.root.rotation.y=heading;
  const inverse=new THREE.Matrix4().compose(new THREE.Vector3(slot.pose.x,0,slot.pose.z),new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,1,0),heading),new THREE.Vector3(1,1,1)).invert();
  for(let frame=0;frame<=120;frame++){
   cast.update(frame/30);
   cast.actor.traverse(part=>{
    if(!(part instanceof THREE.Mesh))return;
    const vertices=part.geometry.attributes.position,p=new THREE.Vector3();
    for(let i=0;i<vertices.count;i++){
     p.fromBufferAttribute(vertices,i).applyMatrix4(part.matrixWorld).applyMatrix4(inverse);
     assert.ok(Math.abs(p.x)<=size.width/2&&Math.abs(p.z)<=size.length/2,`${name}/${fatal} leaves reserved footprint at ${frame/30}: ${p.toArray()}`);
    }
   });
  }
 }
 // The reserved rectangle itself stays on this frontage's pavement, outside
 // both the road (z<=36) and the conservative building envelope (z>=39.5).
 for(const slot of slots){
  const size=trafficSize(slot.model);
  assert.ok(slot.pose.z-size.width/2>36);
  assert.ok(slot.pose.z+size.width/2<39.5);
 }
});

test('fatal accident aftermath retains the final rig pose and occupied slot until cleanup',async()=>{
 const {CityAccident}=await import('../.runtime/frontend-test/city3dAccident.js');
 const {CityAftermath}=await import('../.runtime/frontend-test/city3dAftermath.js');
 const {sceneSlots}=await import('../.runtime/frontend-test/city3dEvents.js');
 const lot={id:'bar',x:80,z:48,row:1,col:2},lots=new Map([['bar',lot]]);
 for(const [name,face] of [['person',1],['woman',3]]){
  const source=await model(name),cast=new CityAccident(source.clone(true),true),city=new CityAftermath();
  const slot=sceneSlots(lot,'accident')[1],id='player:1';
  cast.root.position.set(slot.root.x,.2,slot.root.z);cast.root.rotation.y=slot.pose.heading;cast.update(4);
  const records=[{id:'fatal',target:'bar',victim:{id,name:'Recorded victim'},cause:'charge-accident',face,minute:600,police_at:605,cleanup_at:780}];
  const models=new Map([[name,source],['police',new THREE.Group()],['police-officer',source]]);
  city.rememberBody(id,slot,0);
  city.update(records,600,lots,models,()=>{throw Error('looked up current player instead of recorded appearance');},[],new Set([id]));
  assert.equal(city.slots().length,0);
  city.update(records,600,lots,models,()=>'',[],new Set());
  const body=city.object('aftermath:fatal:body'),actor=body.children.at(-1);
  body.updateMatrixWorld(true);
  cast.actor.traverse(part=>{
   if(!(part instanceof THREE.Mesh))return;
   const held=actor.getObjectByName(part.name);assert.ok(held,part.name);
   assert.ok(part.matrixWorld.elements.every((n,i)=>Math.abs(n-held.matrixWorld.elements[i])<1e-7),`${name} ${part.name} jumped at handoff`);
  });
  assert.deepEqual(city.slots()[0],slot);
  city.update(records,605,lots,models,()=>'',[],new Set());
  assert.equal(city.inspect().filter(e=>e.id.includes('police')).length,1);
  city.update(records,780,lots,models,()=>'',[],new Set());assert.equal(city.slots().length,0);city.dispose();
 }
});

test('strike victims stay grounded throughout every fall and its retained aftermath',async()=>{
 const {CityAssassination}=await import('../.runtime/frontend-test/city3dAssassination.js');
 for(const name of ['person','woman'])for(const variant of ['back-of-head','close-shot','burst','close-quarters']){
  const victim=await model(name),cast=new CityAssassination(await model('person'),victim,variant==='close-quarters'?undefined:new THREE.Group(),variant);
  cast.root.position.set(11,.2,23);cast.root.rotation.y=1.1;
  let previous;
  for(let frame=0;frame<=cast.duration*30;frame++){
   cast.update(frame/30);
   let lowest=Infinity;
   victim.traverseVisible(part=>{
    if(!(part instanceof THREE.Mesh))return;
    const vertices=part.geometry.attributes.position,p=new THREE.Vector3();
    for(let i=0;i<vertices.count;i++)lowest=Math.min(lowest,p.fromBufferAttribute(vertices,i).applyMatrix4(part.matrixWorld).y);
   });
   assert.ok(Math.abs(lowest-.205)<1e-6,`${name}/${variant} at ${frame/30} has floor ${lowest}`);
   if(previous!==undefined)assert.ok(Math.abs(victim.position.y-previous)<.09,`${name}/${variant} ground correction jumped`);
   previous=victim.position.y;
  }
 }
});
