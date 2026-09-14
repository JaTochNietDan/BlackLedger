import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CounterWipe} from '../.runtime/frontend-test/interiorService.js';
import {poseInteriorOccupant} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();loader.register(()=>({name:'no-textures',loadMaterial(){return Promise.resolve(new THREE.MeshStandardMaterial());}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('male and female bartender hands remain on the cloth and clear the counter fixtures',async()=>{
 const cloth=await load('bar-cloth');
 for(const name of ['person','woman']){
  const actor=await load(name);poseInteriorOccupant(actor,{id:'service',x:0,z:3.45,yaw:Math.PI});
  const service=new CounterWipe(actor);assert.equal(service.available,true);
  let hand;actor.getObjectByName('elbow1').traverse(o=>{if(o.isMesh&&o.name.startsWith('hand'))hand=o;});
  const positions=[];
  for(let frame=0;frame<360;frame++){
   service.pose(frame/30);cloth.position.copy(service.clothPosition);cloth.updateMatrixWorld(true);
   const handBox=new THREE.Box3().setFromObject(hand,true),clothBox=new THREE.Box3().setFromObject(cloth,true);
   assert.ok(Math.abs(handBox.min.y-1.201)<.002,`${name} hand loses contact at ${frame}: ${handBox.min.y}`);
   assert.ok(clothBox.min.y>1.185&&clothBox.max.y<1.206);
   assert.ok(clothBox.min.x>-.7&&clothBox.max.x<.1&&clothBox.min.z>2.8&&clothBox.max.z<3.225,'cloth leaves clear rear countertop strip');
   for(const side of [-1,1]){const hip=actor.getObjectByName('leg'+side);assert.equal(hip.rotation.x,0,'wiping moved feet');}
   assert.ok(actor.position.distanceTo(new THREE.Vector3(0,.03,3.45))<1e-9);
   positions.push(service.clothPosition.x);
  }
  assert.ok(Math.max(...positions)-Math.min(...positions)>.2,'wipe does not actually move');
 }
});

test('missing service rig is safely unavailable',()=>{assert.equal(new CounterWipe(new THREE.Group()).available,false);});
