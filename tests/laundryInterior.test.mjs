import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {laundryPlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Bluebird public occupants fit the authored furniture and both rigs remain separately selectable',async()=>{
 const room=await load('interior-laundry');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const bounds=new THREE.Box3().setFromObject(room,true);assert.ok(bounds.min.x>=-5.11&&bounds.max.x<=5.01&&bounds.max.y<=4.21);
 const people=[{id:'a-worker',role:'Laundress'},...Array.from({length:24},(_,i)=>({id:`customer-${i}`}))];
 const places=laundryPlacements(people);assert.equal(places.size,12);assert.equal(places.get('a-worker').id,'laundry-counter');
 assert.deepEqual([...laundryPlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('laundry')];
 for(const name of ['person','woman']){
  const source=await load(name),occupied=[];
  for(const spot of spots){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${name}: ${spot.id} overlaps ${other.id}`);
   occupied.push({id:spot.id,box});
   if(spot.seat!==undefined)continue;
   const floor=new THREE.Raycaster(new THREE.Vector3(spot.x+.09,.1,spot.z+.09),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(floor&&Math.abs(floor.point.y-.0175)<.002,`${spot.id}: unsupported floor`);
   for(const y of [.2,.65,1.1,1.65])for(let i=0;i<24;i++){
    const dir=new THREE.Vector3(Math.cos(i*Math.PI/12),0,Math.sin(i*Math.PI/12));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(spot.x,y,spot.z),dir,0,spot.id==='laundry-counter'?.28:.48).intersectObject(room,true).length,0,`${spot.id}: furniture at ${y}`);
   }
  }
 }
});
