import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {blockingBuildings} from '../.runtime/frontend-test/city3dOcclusion.js';
import {buildingCondition, buildingCutaway} from '../.runtime/frontend-test/city3dDamage.js';

test('orthographic cutaways select actual intervening geometry, including off-centre sight lines',()=>{
 const camera=new THREE.OrthographicCamera(-10,10,10,-10,.1,3000);
 camera.position.set(0,1,20);camera.lookAt(0,1,0);camera.updateMatrixWorld();
 const make=(x,z)=>{const g=new THREE.Group();g.add(new THREE.Mesh(new THREE.BoxGeometry(2,3,2),new THREE.MeshStandardMaterial()));g.position.set(x,1,z);g.updateMatrixWorld(true);g.userData.sightBounds=new THREE.Box3().setFromObject(g);return g;};
 const buildings=new Map([['front',make(5,5)],['behind',make(5,-5)],['beside',make(0,5)]]);
 assert.deepEqual([...blockingBuildings(camera,new THREE.Vector3(5,1,0),buildings)],['front']);
 assert.equal(blockingBuildings(camera,new THREE.Vector3(8,1,0),buildings).size,0);
 camera.position.set(5,1,-20);camera.lookAt(5,1,0);camera.updateMatrixWorld();
 assert.deepEqual([...blockingBuildings(camera,new THREE.Vector3(5,1,0),buildings)],['behind']);
});

test('cutaways update private shader uniforms and fully restore without changing material opacity',()=>{
 const source=new THREE.MeshStandardMaterial(),a=source.clone(),b=source.clone();
 for(const m of [a,b])buildingCondition(m,100,new THREE.Vector3());
 const shader={uniforms:{},vertexShader:THREE.ShaderLib.standard.vertexShader,fragmentShader:THREE.ShaderLib.standard.fragmentShader};
 a.onBeforeCompile(shader,{});const version=a.version;
 buildingCutaway(a,1,new THREE.Vector3(300,200,120));
 assert.equal(shader.uniforms.cityCutaway.value,1);
 assert.deepEqual(shader.uniforms.cityWindow.value.toArray(),[300,200,120]);
 assert.equal(b.userData.cityWear.cutaway.value,0);
 buildingCutaway(a,0,new THREE.Vector3());
 assert.equal(shader.uniforms.cityCutaway.value,0);assert.equal(a.version,version);
 assert.equal(a.opacity,1);assert.equal(a.transparent,false);assert.equal(a.depthWrite,true);
});
