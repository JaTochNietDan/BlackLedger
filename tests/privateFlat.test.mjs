import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
import {placementsForInterior} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
test('private apartment never stages building occupants and both player rigs have a clear entry',async()=>{
 assert.equal(placementsForInterior('flat',[{id:'visitor'},{id:'concierge',role:'Concierge'}]).size,0);
 const room=await load('interior-flat');room.updateMatrixWorld(true);
 for(const name of ['person','woman']){
  const actor=await load(name),arrival=new InteriorArrival(actor,'flat');
  for(let i=0;i<=40;i++){
   arrival.pose(arrival.duration*i/40);const {x,z}=actor.position;
   const floor=new THREE.Raycaster(new THREE.Vector3(x,.1,z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(floor&&floor.point.y>=.009&&floor.point.y<=.019,`entry ${i}: unsupported floor`);
   for(const y of [.3,.9,1.5])for(let r=0;r<16;r++){
    const dir=new THREE.Vector3(Math.cos(r*Math.PI/8),0,Math.sin(r*Math.PI/8));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(x,y,z),dir,0,.48).intersectObject(room,true).length,0,`entry ${i}: blocked`);
   }
  }
 }
});
