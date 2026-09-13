import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {frameScene} from '../.runtime/frontend-test/city3dFraming.js';

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
