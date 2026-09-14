import test from 'node:test';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {butcherPlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Butcher customers and worker occupy clear supported floor around the actual display',async()=>{
 const room=await load('interior-butcher');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const people=[{id:'a-worker',role:'Butcher'},...Array.from({length:16},(_,i)=>({id:`customer-${i}`}))];
 const places=butcherPlacements(people);assert.equal(places.size,10);assert.equal(places.get('a-worker').id,'butcher-service');
 assert.deepEqual([...butcherPlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('butcher')];
 for(const name of ['person','woman']){
  const source=await load(name),occupied=[];
  for(const spot of spots){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${name}: ${spot.id} overlaps ${other.id}`);
   occupied.push({id:spot.id,box});
   const floor=new THREE.Raycaster(new THREE.Vector3(spot.x+.09,.1,spot.z+.09),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(floor&&Math.abs(floor.point.y-.0175)<.001,`${spot.id}: unsupported floor`);
   for(const y of [.2,.65,1.1,1.65])for(let i=0;i<24;i++){
    const dir=new THREE.Vector3(Math.cos(i*Math.PI/12),0,Math.sin(i*Math.PI/12));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(spot.x,y,spot.z),dir,0,.48).intersectObject(room,true).length,0,`${spot.id}: furniture at ${y}`);
   }
  }
 }
});

test('Butcher entrance stays on the shop floor and clears the fullest public roster',async()=>{
 const room=await load('interior-butcher');room.updateMatrixWorld(true);
 for(const name of ['person','woman']){
  const source=await load(name),actor=source.clone(true),arrival=new InteriorArrival(actor,'butcher');
  const boxes=[...butcherPlacements(Array.from({length:20},(_,i)=>({id:`p-${i}`}))).values()].map(s=>{
   const o=source.clone(true);poseInteriorOccupant(o,s);return new THREE.Box3().setFromObject(o,true);
  });
  for(let frame=0;frame<=60;frame++){
   arrival.pose(arrival.duration*frame/60);const box=new THREE.Box3().setFromObject(actor,true);
   assert.ok(box.min.y>=.0175&&box.max.z<4.5,'feet leave the entrance floor');
   for(const other of boxes)assert.equal(box.intersectsBox(other),false,'entrance crosses a customer');
  }
 }
});
