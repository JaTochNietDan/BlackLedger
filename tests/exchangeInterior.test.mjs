import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {exchangePlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Mercer Exchange seats support both rigs and standing visitors and entrance clear the furniture',async()=>{
 const room=await load('interior-exchange');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const people=[{id:'a-worker',role:'Clerk'},{id:'b-worker',role:'Broker'},{id:'c-worker',role:'Teller'},...Array.from({length:16},(_,i)=>({id:`visitor-${i}`}))];
 const places=exchangePlacements(people);assert.equal(places.size,13);assert.equal(places.get('a-worker').id,'exchange-clerk-0');
 assert.deepEqual([...exchangePlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('market')];
 for(const name of ['person','woman']){
  const source=await load(name),occupied=[];
  for(const spot of spots){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${name}: ${spot.id} overlaps ${other.id}`);
   occupied.push({id:spot.id,box});
   if(spot.seat!==undefined){
    const support=new THREE.Raycaster(new THREE.Vector3(spot.x,.75,spot.z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
    assert.ok(support&&Math.abs(support.point.y-spot.seat)<.01,`${spot.id}: missing cushion`);
    continue;
   }
   for(const y of [.2,.65,1.1,1.65])for(let i=0;i<24;i++){
    const dir=new THREE.Vector3(Math.cos(i*Math.PI/12),0,Math.sin(i*Math.PI/12));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(spot.x,y,spot.z),dir,0,.48).intersectObject(room,true).length,0,`${spot.id}: furniture at ${y}`);
   }
  }
  const actor=source.clone(true),arrival=new InteriorArrival(actor,'market');
  for(let i=0;i<=40;i++){
   arrival.pose(arrival.duration*i/40);
   const {x,z}=actor.position;
   const support=new THREE.Raycaster(new THREE.Vector3(x,.1,z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(support&&support.point.y>=.009&&support.point.y<=.019,`entry ${i}: unsupported floor`);
   for(const y of [.3,.9,1.5])for(let r=0;r<16;r++){
    const dir=new THREE.Vector3(Math.cos(r*Math.PI/8),0,Math.sin(r*Math.PI/8));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(x,y,z),dir,0,.48).intersectObject(room,true).length,0,`entry ${i}: blocked`);
   }
  }
 }
});
