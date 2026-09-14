import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityIncendiary,incendiaryShard,incendiaryStride,INCENDIARY_RELEASE,INCENDIARY_IMPACT} from '../.runtime/frontend-test/city3dIncendiary.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('bottle stays at the authored hand, releases continuously and reaches the impact point',async()=>{
 for(const name of ['person','woman']){
  const actor=await model(name),bottle=await model('incendiary-bottle'),target=new THREE.Vector3(0,2.5,4);
  const cast=new CityIncendiary(actor,bottle,target);cast.root.position.set(13,.2,-8);cast.root.rotation.y=.7;
  for(let t=0;t<INCENDIARY_RELEASE;t+=.025){
   cast.update(t);
   const hand=actor.getObjectByName('elbow1').localToWorld(new THREE.Vector3(.01,-.295,.005));
   assert.ok(hand.distanceTo(bottle.getWorldPosition(new THREE.Vector3()))<1e-6,'bottle slips from hand');
   // The substantial body must extend beyond the hand, away from the forearm.
   const body=bottle.localToWorld(new THREE.Vector3(0,-.16,0));
   const elbow=actor.getObjectByName('elbow1').getWorldPosition(new THREE.Vector3());
   const nearest=new THREE.Line3(elbow,hand).closestPointToPoint(body,true,new THREE.Vector3());
   assert.ok(body.distanceTo(nearest)>.14,'body embeds in the forearm');
  }
  cast.update(INCENDIARY_RELEASE-1e-6);const position=bottle.position.clone(),rotation=bottle.quaternion.clone();
  cast.update(INCENDIARY_RELEASE);assert.ok(bottle.position.distanceTo(position)<1e-5);assert.ok(rotation.angleTo(bottle.quaternion)<1e-5);
  for(let t=INCENDIARY_RELEASE;t<INCENDIARY_IMPACT;t+=.01){cast.update(t);assert.equal(bottle.visible,true);assert.ok(bottle.position.y>1.5,'flight drops below the release area');}
  cast.update(INCENDIARY_IMPACT);assert.equal(bottle.visible,false);assert.ok(bottle.position.distanceTo(target)<1e-6);
  cast.update(cast.duration);assert.ok(actor.position.x<=-4);assert.equal(bottle.visible,false);
 }
});

test('the actor sweep fits the reserved city forecourt footprint',async()=>{
 const {sceneSlots}=await import('../.runtime/frontend-test/city3dEvents.js');
 const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const slot=sceneSlots({x:40,z:40,row:1},'incendiary')[0],size=trafficSize(slot.model);
 for(const name of ['person','woman']){
  const actor=await model(name),cast=new CityIncendiary(actor,await model('incendiary-bottle'),new THREE.Vector3(0,4,4));
  cast.root.position.set(slot.root.x,.2,slot.root.z);
  for(let frame=0;frame<=180;frame++){
   cast.update(frame/180*cast.duration);const bounds=new THREE.Box3().setFromObject(actor,true);
   assert.ok(bounds.min.x>=slot.pose.x-size.width/2&&bounds.max.x<=slot.pose.x+size.width/2,'actor escapes lateral reservation');
   assert.ok(bounds.min.z>=slot.pose.z-size.length/2&&bounds.max.z<=slot.pose.z+size.length/2,'actor escapes forecourt depth');
  }
 }
});

test('approach and escape brake continuously without a limb reset',async()=>{
 for(const [distance,duration,ramp] of [[2,2/1.2,.22],[4,4/2.1,.28]]){
  assert.equal(incendiaryStride(-1,distance,duration,ramp).distance,0);
  assert.equal(incendiaryStride(duration+1,distance,duration,ramp).distance,distance);
  let previous=0;
  for(let i=0;i<=240;i++){
   const pose=incendiaryStride(duration*i/240,distance,duration,ramp);
   assert.ok(pose.distance>=previous&&pose.distance<=distance);previous=pose.distance;
  }
  for(const t of [0,ramp,duration-ramp,duration]){
   const left=incendiaryStride(t-1e-5,distance,duration,ramp),right=incendiaryStride(t+1e-5,distance,duration,ramp);
   assert.ok(Math.abs(left.weight-right.weight)<.001,'speed jumps at ramp boundary');
  }
 }
 for(const name of ['person','woman']){
  const actor=await model(name),cast=new CityIncendiary(actor,await model('incendiary-bottle'),new THREE.Vector3(0,4,4));
  const joints=['arm1','elbow1','arm-1','leg1','leg-1','knee1','knee-1'].map(name=>actor.getObjectByName(name));
  for(const t of [.2,.2+2/1.2,3.75,4,4.25,4.25+4/2.1]){
   cast.update(t-1e-5);const before=joints.map(j=>j.quaternion.clone()),position=actor.position.clone();
   cast.update(t+1e-5);
   assert.ok(position.distanceTo(actor.position)<.001,'root jumps at a movement boundary');
   joints.forEach((j,i)=>assert.ok(before[i].angleTo(j.quaternion)<.001,`${name} ${j.name} snaps at ${t}`));
  }
 }
});

test('glass scatter begins at impact and settles above the pavement',async()=>{
 const shard=await model('bottle-shard'),bounds=new THREE.Box3().setFromObject(shard,true);
 assert.ok(bounds.max.x-bounds.min.x<.13);assert.ok(bounds.max.y-bounds.min.y<.02,'glass is a thick masonry fragment');
 for(const height of [1.5,3,5])for(let j=0;j<12;j++){
  assert.equal(incendiaryShard(j,-.01,height).scale,0);
  const origin=incendiaryShard(j,0,height);assert.equal(Math.abs(origin.x),0);assert.equal(origin.y,0);assert.equal(Math.abs(origin.z),0);
  for(let age=0;age<2.4;age+=.01){const p=incendiaryShard(j,age,height);assert.ok(p.y>=-height);assert.ok(p.z<=0,'shard travels into facade');assert.ok(p.scale>=0);}
  assert.equal(incendiaryShard(j,2.4,height).scale,0);
 }
});
