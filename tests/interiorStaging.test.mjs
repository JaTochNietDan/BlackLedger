import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {interiorPlacements,poseInteriorOccupant} from '../.runtime/frontend-test/interiorStaging.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshStandardMaterial());}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('room occupancy is stable, exclusive and assigns service only to a public bartender',()=>{
 const people=[{id:'mara'},{id:'worker',role:'Bartender'},...Array.from({length:9},(_,i)=>({id:`guest-${i}`}))];
 const a=interiorPlacements(people),b=interiorPlacements([...people].reverse());
 assert.deepEqual([...a],[...b]);assert.equal(a.size,11);
 assert.equal(a.get('mara').id,'booth-0--1');assert.equal(a.get('worker').id,'service');
 assert.equal(new Set([...a.values()].map(s=>s.id)).size,a.size);
 assert.equal([...interiorPlacements([{id:'visitor',role:'Merchant'}]).values()].some(s=>s.id==='service'),false);
});
test('seated rigs remain above the floor and occupants have separate physical space',async()=>{
 for(const name of ['person','woman']){
  const placements=interiorPlacements(Array.from({length:10},(_,i)=>({id:`guest-${i}`})));
  const source=await model(name),boxes=[];
  for(const [id,spot] of placements){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);
   const box=new THREE.Box3().setFromObject(actor,true);
   assert.ok(box.min.y>=0,`${name}/${spot.id} feet penetrate floor: ${box.min.y}`);
   assert.ok(box.min.x>-5.75&&box.max.x<5.8&&box.min.z>-4.8&&box.max.z<3.7,'occupant leaves room');
   if(spot.seat!==undefined)assert.ok(Math.abs(actor.getObjectByName('leg1').getWorldPosition(new THREE.Vector3()).y-spot.seat)<1e-6,'pelvis is not on its seat');
   for(const other of boxes)assert.equal(box.intersectsBox(other.box),false,`${name}: ${id}/${other.id} overlap ${JSON.stringify({a:box,b:other.box})}`);
   boxes.push({id,box});
  }
 }
});
