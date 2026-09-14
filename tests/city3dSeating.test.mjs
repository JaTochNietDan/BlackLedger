import test from 'node:test';import assert from 'node:assert/strict';import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';import {readFileSync} from 'node:fs';
import {seatDriver} from '../.runtime/frontend-test/city3dSeating.js';
async function load(name){const b=readFileSync(`public/art/models/${name}.glb`),l=new GLTFLoader();l.register(()=>({name:'seating-geometry',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial());}}));return(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
for(const name of ['ford','hudson','packard'])test(`${name} seats full-size driver rigs above the floor and below the roof`,async()=>{
 const car=await load(name);
 for(const rig of ['person','woman']){
  const actor=await load(rig),seat=car.getObjectByName('seat-front-left');seatDriver(actor,seat.position,['left','right'].map(side=>car.getObjectByName('seat-driver-grip-'+side).position));car.add(actor);car.updateMatrixWorld(true);
  const bounds=new THREE.Box3();actor.traverseVisible(o=>{if(o.isMesh)bounds.union(new THREE.Box3().setFromObject(o,true));});
  assert.ok(bounds.min.y>=.15,`${rig} enters the floor: ${bounds.min.y}`);
  assert.ok(bounds.max.y<=1.43,`${rig} enters roof: ${bounds.max.y}`);
  assert.ok(bounds.min.x>=-.745&&bounds.max.x<=.745,'driver enters a door');
  const wheelZ=car.getObjectByName('seat-driver-grip-left').position.z;
  assert.ok(bounds.max.z<=wheelZ+.20,`${rig} extends beyond the wheel into the dashboard: ${bounds.max.z}`);
  assert.equal(actor.scale.x,1);assert.equal(actor.scale.y,1);assert.equal(actor.scale.z,1);
  const pelvis=new THREE.Vector3(0,.86,0).applyMatrix4(actor.matrixWorld);assert.ok(pelvis.distanceTo(seat.position)<1e-6);
  for(const side of [-1,1]){
   const elbow=actor.getObjectByName(`elbow${side}`);
   const hand=elbow.localToWorld(new THREE.Vector3(0,-.295,0));
   const grip=car.getObjectByName('seat-driver-grip-'+(side<0?'left':'right')).getWorldPosition(new THREE.Vector3());
   assert.ok(hand.distanceTo(grip)<.025,`${rig} hand misses steering wheel by ${hand.distanceTo(grip)}`);
  }
  car.remove(actor);
 }
});
