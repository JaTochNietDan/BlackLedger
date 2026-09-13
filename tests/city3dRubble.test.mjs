import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {CityRubble} from '../.runtime/frontend-test/city3dRubble.js';
import {windowDebris} from '../.runtime/frontend-test/city3dBlast.js';
import {surfaceHeight} from '../.runtime/frontend-test/city3dPlan.js';

test('saved rubble matches settled fragments, survives extinguishing and yields to replay until cleanup',()=>{
 const rubble=new CityRubble(),building=new THREE.Group(),source=new THREE.Group();
 source.add(new THREE.Mesh(new THREE.BoxGeometry(.18,.08,.09),new THREE.MeshStandardMaterial()));
 const windows=[-3,-1,1,3].map(x=>new THREE.Vector3(80+x,5.13,42.28));building.userData.blastWindows=windows;
 const fire={id:'f',target:'bar',minute:480,brigade_at:490,extinguished_at:660,cleanup_at:720};
 const update=(minute,playing=new Set())=>rubble.update([fire],minute,new Map([['bar',building]]),source,playing);
 update(479);assert.equal(rubble.inspect().length,0);
 update(480);const mesh=rubble.root.children[0];assert.equal(mesh.count,12);
 for(let i=0;i<12;i++){
  const matrix=new THREE.Matrix4();mesh.getMatrixAt(i,matrix);const p=new THREE.Vector3().setFromMatrixPosition(matrix),w=windows[i%4];
  const end=windowDebris(i,2.5,w,surfaceHeight({x:w.x,z:w.z-3}));
  assert.ok(Math.abs(p.x-end.x)<1e-5&&Math.abs(p.z-end.z)<1e-5);
  assert.ok(Math.abs(p.y-(surfaceHeight(p)+.04*end.scale+.006))<1e-6);
 }
 const before=Array.from(mesh.instanceMatrix.array);update(600);assert.deepEqual(Array.from(mesh.instanceMatrix.array),before);
 update(600,new Set(['bar']));assert.equal(mesh.visible,false);update(660);assert.equal(mesh.visible,true);assert.equal(rubble.root.children[0],mesh);
 let instancesDisposed=0,geometryDisposed=0,sourceDisposed=0;
 mesh.addEventListener('dispose',()=>instancesDisposed++);mesh.geometry.addEventListener('dispose',()=>geometryDisposed++);source.children[0].geometry.addEventListener('dispose',()=>sourceDisposed++);
 update(720);assert.equal(rubble.inspect().length,0);assert.equal(instancesDisposed,1);assert.equal(geometryDisposed,0);
 rubble.dispose();assert.equal(geometryDisposed,1);assert.equal(sourceDisposed,0);
});
