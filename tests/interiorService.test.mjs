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

test('laundry hands smooth linen without penetrating the counter or stretching the rig',async()=>{
 const {LinenPress}=await import('../.runtime/frontend-test/interiorService.js');
 for(const name of ['person','woman']){
  const actor=await load(name);poseInteriorOccupant(actor,{id:'laundry-counter',x:3.35,z:-2,yaw:0});
  const service=new LinenPress(actor),hands=[];assert.equal(service.available,true);
  for(const side of [-1,1])actor.getObjectByName('elbow'+side).traverse(o=>{if(o.isMesh&&o.name.startsWith('hand'))hands.push(o);});
  let min=Infinity,max=-Infinity;
  for(let frame=0;frame<360;frame++){
   assert.equal(service.step(1/30,true,true,false),true);
   const boxes=hands.map(hand=>new THREE.Box3().setFromObject(hand,true));
   for(const box of boxes){
    assert.ok(Math.abs(box.min.y-1.109)<.003,`${name} loses linen contact: ${box.min.y}`);
    assert.ok(box.min.x>3.02&&box.max.x<3.68&&box.min.z> -1.66&&box.max.z< -1.32,`${name} leaves linen: ${JSON.stringify(box)}`);
   }
   assert.equal(boxes[0].intersectsBox(boxes[1]),false,'hands cross');
   min=Math.min(min,boxes[0].min.x);max=Math.max(max,boxes[0].min.x);
   for(const side of [-1,1]){
    assert.equal(actor.getObjectByName('leg'+side).rotation.x,0);
    assert.deepEqual(actor.getObjectByName('arm'+side).scale.toArray(),[1,1,1]);
    // The lower sleeve must remain above the solid countertop wherever it
    // crosses the rear edge, not merely put its hand on the correct height.
    actor.getObjectByName('arm'+side).traverse(o=>{if(!o.isMesh)return;const attr=o.geometry.attributes.position,v=new THREE.Vector3();for(let i=0;i<attr.count;i++){v.fromBufferAttribute(attr,i).applyMatrix4(o.matrixWorld);if(v.z>=-1.67)assert.ok(v.y>=1.089,`${name} sleeve penetrates counter`);}});
   }
  }
  assert.ok(max-min>.07,'smoothing is motionless');
  const seconds=service.seconds;
  assert.equal(service.step(2,true,false,false),false);assert.equal(service.step(2,true,true,true),false);assert.equal(service.seconds,seconds);
  assert.equal(service.step(1,false,true,false),true);assert.equal(service.active,false);
  for(const side of [-1,1])for(const part of ['arm','elbow'])assert.ok(actor.getObjectByName(part+side).quaternion.angleTo(new THREE.Quaternion())<1e-7);
  assert.equal(service.step(1,false,true,false),false);
  assert.equal(service.step(120,true,true,false),true);assert.ok(Math.abs(service.seconds-seconds-.05)<1e-8);
 }
 assert.equal(new LinenPress(new THREE.Group()).available,false);
});
