import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {cardPosition,blackjackCamera} from '../.runtime/frontend-test/blackjackPresentation.js';

test('close blackjack camera contains both long hands at desktop and narrow aspects',()=>{
 for(const [width,height] of [[1280,440],[800,600],[360,440]]){
  const aspect=width/height,camera=new THREE.PerspectiveCamera(38,aspect,.01,20),center=new THREE.Vector3(0,.87,0),distance=blackjackCamera.close*Math.max(1,blackjackCamera.fitAspect/aspect);
  camera.position.copy(center).add(new THREE.Vector3(0,distance*.84,distance*.54));camera.lookAt(center);camera.updateMatrixWorld();
  for(const row of [0,1])for(const count of [2,7,12])for(let i=0;i<count;i++){
   const pos=cardPosition(i,count,row);
   for(const x of [-.15,.15])for(const z of [-.215,.215]){
    const projected=new THREE.Vector3(pos[0]+x,pos[1],pos[2]+z).project(camera);
    assert.ok(Math.abs(projected.x)<.98&&Math.abs(projected.y)<.98,`${width}x${height}: card ${i}/${count} clipped`);
   }
  }
 }
});
