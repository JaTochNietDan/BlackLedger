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
 for(const name of ['monarch','tavern','tenement','mercer-court','casino','civic','shop','warehouse','bluehour','goldenlily','papermoon','mariner','undertaker'])for(const person of ['person','woman']){
  const building=await model(name),actor=await model(person),cast=new CityPlanter(actor);
  building.updateMatrixWorld(true);
  cast.root.position.copy(building.getObjectByName('entrance-threshold').getWorldPosition(new THREE.Vector3()));
  const door=building.getObjectByName('entrance-door-hinge'),triangle=new THREE.Triangle(),bounds=new THREE.Box3();
  for(let frame=0;frame<=186;frame++){
   const pose=cast.update(frame/30);door.rotation.y=-Math.PI/2*pose.door;building.updateMatrixWorld(true);
   const actorBounds=new THREE.Box3().setFromObject(actor,true);actorBounds.min.y+=.12;
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
