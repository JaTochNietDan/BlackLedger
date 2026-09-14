import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {frameScene,stagedSceneBounds} from '../.runtime/frontend-test/city3dFraming.js';

test('event envelope fits at extreme previous zooms, rotations and narrow viewports',()=>{
 for(const aspect of [.35,1,2.4]) for(const zoom of [.6,32]) for(const angle of [0,1,3,5]) {
  const camera=new THREE.OrthographicCamera(-100*aspect,100*aspect,100,-100,.1,3000);
  const target=new THREE.Vector3(180,0,100);
  camera.position.copy(target).add(new THREE.Vector3(Math.sin(angle)*240,200,Math.cos(angle)*240));
  camera.zoom=zoom;camera.lookAt(target);camera.updateProjectionMatrix();
  const box=new THREE.Box3(new THREE.Vector3(5,0,24),new THREE.Vector3(30,22,50));
  frameScene(camera,target,box);
  for(const x of [5,30])for(const y of [0,22])for(const z of [24,50]){
   const screen=new THREE.Vector3(x,y,z).project(camera);
   assert.ok(Math.abs(screen.x)<=.580001 && Math.abs(screen.y)<=.580001);
  }
  assert.ok(camera.zoom>0 && camera.zoom<=32);
 }
});

test('impact projection is bounded and restored even if rendering fails',async()=>{
 const {impactPulse,renderImpact}=await import('../.runtime/frontend-test/city3dFraming.js');
 const camera=new THREE.OrthographicCamera(-10,10,8,-8,.1,1000);camera.updateProjectionMatrix();
 const original=camera.projectionMatrix.clone(),inverse=camera.projectionMatrixInverse.clone();
 assert.deepEqual(impactPulse(-.01,10),{x:0,y:0});assert.deepEqual(impactPulse(.32,10),{x:0,y:0});
 for(let t=0;t<.32;t+=.001){const p=impactPulse(t,11);assert.ok(Math.abs(p.x)<=6.05&&Math.abs(p.y)<=11);}
 assert.throws(()=>renderImpact(camera,100,100,1000,800,()=>{
  assert.ok(Math.abs(camera.projectionMatrix.elements[12]-original.elements[12]-.018)<1e-9);
  assert.ok(Math.abs(camera.projectionMatrix.elements[13]-original.elements[13]-.03)<1e-9);
  throw new Error('render failure');
 }));
 assert.deepEqual(camera.projectionMatrix.elements,original.elements);
 assert.deepEqual(camera.projectionMatrixInverse.elements,inverse.elements);
 renderImpact(camera,0,0,1000,800,()=>assert.deepEqual(camera.projectionMatrix.elements,original.elements));
});


test('arrest framing contains both police cars and every corner of the escort reservation',()=>{
 const slots=[
  {pose:{x:74.5,z:38.65,heading:0},model:'custody'},
  {pose:{x:89.6,z:42,heading:0},model:'parked-police'},
  {pose:{x:89.6,z:48,heading:Math.PI/2},model:'parked-police'},
  {pose:{x:83.8,z:38.35,heading:0},model:'casualty'},
 ];
 const bounds=stagedSceneBounds(slots);
 assert.ok(bounds.containsPoint(new THREE.Vector3(72.3,0,37.65)));
 assert.ok(bounds.distanceToPoint(new THREE.Vector3(91.95,3,49.075))<1e-8);
 for(const aspect of [.55,1.77,2.4])for(const angle of [0,1.4,3.6]){
  const camera=new THREE.OrthographicCamera(-100*aspect,100*aspect,100,-100,.1,3000),target=new THREE.Vector3();
  camera.position.set(Math.sin(angle)*240,200,Math.cos(angle)*240);camera.zoom=32;
  frameScene(camera,target,bounds);
  for(const x of [bounds.min.x,bounds.max.x])for(const y of [0,3])for(const z of [bounds.min.z,bounds.max.z]){
   const p=new THREE.Vector3(x,y,z).project(camera);
   assert.ok(Math.abs(p.x)<=.580001&&Math.abs(p.y)<=.580001);
  }
 }
});
