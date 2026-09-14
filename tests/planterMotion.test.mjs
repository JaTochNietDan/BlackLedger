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


test('stair descent keeps support planted and every shoe above the authored treads',async()=>{
 const {CityStairDescent}=await import('../.runtime/frontend-test/city3dStairs.js');
 const building=await model('villa');building.updateMatrixWorld(true);
 for(const name of ['person','woman']){
  const cast=new CityStairDescent(await model(name));cast.root.position.set(0,.6,-6.23);
  let last;
  for(let frame=0;frame<=468;frame++){
   const t=frame/120,pose=cast.update(t);assert.ok(pose.reached.every(Boolean),`${name} unreachable at ${t}`);
   const support=1-pose.moving;
   if(last&&last.moving===pose.moving)assert.ok(last.feet[support].distanceTo(pose.feet[support])<1e-8,'support foot slides');
   cast.actor.traverse(o=>{
    if(!(o instanceof THREE.Mesh)||!o.name.startsWith('shoe'))return;
    const p=o.geometry.attributes.position;
    for(let i=0;i<p.count;i++){
     const v=new THREE.Vector3().fromBufferAttribute(p,i).applyMatrix4(o.matrixWorld);
     const surface=v.z> -6.4?.6:v.z> -6.8?.4:v.z> -7.2?.2:0;
     assert.ok(v.y>=surface-.00001,`${name} shoe penetrates tread at ${t}: ${v.toArray()}`);
    }
   });
   last=pose;
  }
  const finish=cast.update(cast.duration);assert.ok(finish.done);assert.ok(finish.feet.every(f=>Math.abs(f.y+.5)<1e-6&&Math.abs(f.z+1.2)<1e-6));
 }
});
