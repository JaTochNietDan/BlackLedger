import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {ashburyPlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Ashbury public occupants fit the authored furniture and both rigs remain separately selectable',async()=>{
 const room=await load('interior-ashbury');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const bounds=new THREE.Box3().setFromObject(room,true);assert.ok(bounds.min.x>=-5.11&&bounds.max.x<=5.01&&bounds.max.y<=4.21);
 const people=[{id:'a-worker',role:'Concierge'},...Array.from({length:24},(_,i)=>({id:`customer-${i}`}))];
 const places=ashburyPlacements(people);assert.equal(places.size,13);assert.equal(places.get('a-worker').id,'ashbury-register');
 assert.deepEqual([...ashburyPlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('apartment')];
 for(const name of ['person','woman']){
  const source=await load(name),occupied=[];
  for(const spot of spots){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${name}: ${spot.id} overlaps ${other.id}`);
   occupied.push({id:spot.id,box});
   if(spot.seat!==undefined){const cushion=new THREE.Raycaster(new THREE.Vector3(spot.x,.85,spot.z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];assert.ok(cushion&&Math.abs(cushion.point.y-spot.seat)<.005);continue;}
   const floor=new THREE.Raycaster(new THREE.Vector3(spot.x+.09,.1,spot.z+.09),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(floor&&floor.point.y>=.009&&floor.point.y<=.019,`${spot.id}: unsupported floor`);
   for(const y of [.2,.65,1.1,1.65])for(let i=0;i<24;i++){
    const dir=new THREE.Vector3(Math.cos(i*Math.PI/12),0,Math.sin(i*Math.PI/12));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(spot.x,y,spot.z),dir,0,spot.id==='ashbury-register'?.28:.48).intersectObject(room,true).length,0,`${spot.id}: furniture at ${y}`);
   }
  }
 }
});

test('Ashbury stairs have supported treads, a landing and an open upper exit',async()=>{
 const room=await load('interior-ashbury');room.updateMatrixWorld(true);
 for(let i=0;i<12;i++){
  const hit=new THREE.Raycaster(new THREE.Vector3(3.53,4,-(.05+i*.38)),new THREE.Vector3(0,-1,0),0,5).intersectObject(room,true)[0];
  assert.ok(hit&&Math.abs(hit.point.y-(i+1)*.22)<.002,`tread ${i}`);
 }
 const landing=new THREE.Raycaster(new THREE.Vector3(3.53,4,-4.71),new THREE.Vector3(0,-1,0),0,5).intersectObject(room,true)[0];
 assert.ok(landing&&Math.abs(landing.point.y-2.64)<.002);
 for(const x of [3.1,3.53,3.9])for(const y of [2.8,3.3,3.9]){
  assert.equal(new THREE.Raycaster(new THREE.Vector3(x,y,-4.6),new THREE.Vector3(0,0,-1),0,1).intersectObject(room,true).length,0,'upper exit blocked');
 }
});
