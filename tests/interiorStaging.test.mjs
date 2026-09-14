import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {interiorPlacements,poseInteriorOccupant,placementsForInterior} from '../.runtime/frontend-test/interiorStaging.js';
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

test('Mercer lobby uses exclusive bench and floor positions, with a superintendent at the register',async()=>{
 const {mercerLobbyPlacements}=await import('../.runtime/frontend-test/interiorStaging.js');
 const people=[{id:'super',role:'Superintendent'},...Array.from({length:14},(_,i)=>({id:`tenant-${i}`}))];
 const placements=mercerLobbyPlacements(people);
 assert.deepEqual([...placements],[...mercerLobbyPlacements([...people].reverse())]);
 assert.equal(placements.size,15);assert.equal(placements.get('super').id,'lobby-register');
 for(const name of ['person','woman']){
  const source=await model(name),boxes=[];
  for(const [id,spot] of placements){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);
   const box=new THREE.Box3().setFromObject(actor,true);
   assert.ok(box.min.y>=.0175,`${name}/${spot.id} below terrazzo`);
   assert.ok(box.min.x> -5.7&&box.max.x<2.8&&box.min.z> -6.3&&box.max.z<4.8,`${name}/${spot.id} clips wall or stair`);
   if(spot.seat!==undefined){
    assert.ok(Math.abs(actor.getObjectByName('leg1').getWorldPosition(new THREE.Vector3()).y-.77)<1e-6);
    assert.ok(box.min.x> -5.25,'seated occupant intersects bench back');
   }
   for(const other of boxes)assert.equal(box.intersectsBox(other.box),false,`${name}/${id} overlaps ${other.id}`);
   boxes.push({id,box});
  }
 }
 const room=await model('interior-mercer-court');
 assert.ok(room.getObjectByName('interior-wall-left'));
 assert.ok(room.getObjectByName('interior-wall-back'));
});

test('player has clear floor space in both rooms without displacing any NPC',async()=>{
 const {interiorPlayerSpot,mercerLobbyPlacements}=await import('../.runtime/frontend-test/interiorStaging.js');
 for(const place of ['bar','mercercourt'])for(const name of ['person','woman']){
  const source=await model(name),actor=source.clone(true),spot=interiorPlayerSpot(place);
  poseInteriorOccupant(actor,spot);const playerBox=new THREE.Box3().setFromObject(actor,true);
  assert.ok(playerBox.min.y>=.0175);
  assert.ok(playerBox.min.x>-5.7&&playerBox.max.x<5.7&&playerBox.min.z>-4.8&&playerBox.max.z<4.9);
  const people=[{id:'worker',role:place==='bar'?'Barman':'Superintendent'},...Array.from({length:14},(_,i)=>({id:`guest-${i}`}))];
  const placements=(place==='bar'?interiorPlacements:mercerLobbyPlacements)(people);
  for(const [id,s] of placements){const npc=source.clone(true);poseInteriorOccupant(npc,s);assert.equal(playerBox.intersectsBox(new THREE.Box3().setFromObject(npc,true)),false,`${place}/${name} overlaps ${id}`);}
  if(place==='mercercourt')assert.ok(playerBox.max.x<2.8&&playerBox.min.z>3.7,'player intrudes into staircase or lobby crowd');
  else assert.ok(playerBox.max.z<.7&&playerBox.min.z>-2,'player intrudes into bar stools or aisle crowd');
 }
});

test('Mariner stages its public landlady at reception and keeps every rig clear of desk, stairs and other occupants',async()=>{
 const {marinerLobbyPlacements,interiorPlayerSpot,placementsForInterior}=await import('../.runtime/frontend-test/interiorStaging.js');
 const people=[{id:'host',role:'Landlady'},...Array.from({length:15},(_,i)=>({id:`boarder-${i}`}))];
 const placements=marinerLobbyPlacements(people);
 assert.deepEqual([...placements],[...marinerLobbyPlacements([...people].reverse())]);
 assert.deepEqual([...placements],[...placementsForInterior('room',people)]);
 assert.equal(placements.size,12);assert.equal(placements.get('host').id,'boarding-reception');
 assert.equal(marinerLobbyPlacements([{id:'guest',role:'Boarder'}]).get('guest').id,'boarding-bench-0');
 assert.equal(new Set([...placements.values()].map(s=>s.id)).size,placements.size);
 const desk=new THREE.Box3(new THREE.Vector3(-4.25,0,-4.85),new THREE.Vector3(-.75,1.34,-3.85));
 const stair=new THREE.Box3(new THREE.Vector3(2.37,0,-5.7),new THREE.Vector3(4.85,4.2,.05));
 for(const name of ['person','woman']){
  const source=await model(name),boxes=[];
  for(const [id,spot] of [...placements,['player',interiorPlayerSpot('room')]]){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   assert.ok(box.min.y>=.0175,`${name}/${id} below floorboards`);
   assert.ok(box.min.x> -4.75&&box.max.x<4.75&&box.min.z> -5.8&&box.max.z<3.95,`${name}/${id} clips wall or leaves floor`);
   assert.equal(box.intersectsBox(desk),false,`${name}/${id} intersects reception`);
   assert.equal(box.intersectsBox(stair),false,`${name}/${id} intersects stairs`);
   if(spot.seat!==undefined){assert.ok(Math.abs(actor.getObjectByName('leg1').getWorldPosition(new THREE.Vector3()).y-.77)<1e-6);assert.ok(box.min.x> -4.59,'occupant penetrates bench back');}
   for(const other of boxes)assert.equal(box.intersectsBox(other.box),false,`${name}/${id} overlaps ${other.id}`);
   boxes.push({id,box});
  }
 }
 const room=await model('interior-mariner');assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const bounds=new THREE.Box3().setFromObject(room,true);assert.ok(bounds.min.x>=-5.12&&bounds.max.x<=5.01&&bounds.min.z>=-6.12&&bounds.max.z<=4.01);
});

test('mortuary registry is reserved for a public attendant and receiving spots stay distinct',()=>{
 const people=[{id:'visitor'}, {id:'worker',role:'Receiving attendant'},...Array.from({length:5},(_,i)=>({id:`guest-${i}`}))];
 const a=placementsForInterior('mortuary',people);
 assert.deepEqual([...a],[...placementsForInterior('mortuary',[...people].reverse())]);
 assert.equal(a.get('worker').id,'registry-clerk');
 assert.equal(new Set([...a.values()].map(s=>s.id)).size,a.size);
 assert.equal([...placementsForInterior('mortuary',[{id:'visitor'}]).values()].some(s=>s.id==='registry-clerk'),false);
});
