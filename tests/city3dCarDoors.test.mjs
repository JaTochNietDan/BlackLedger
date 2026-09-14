import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
for(const name of ['ford','hudson','packard','police'])test(`${name} has four real door openings and seat anchors`,async()=>{
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(parser=>({name:'door-test-materials',loadMaterial(index){return Promise.resolve(new THREE.MeshBasicMaterial({name:parser.json.materials[index].name,side:THREE.DoubleSide}));}}));
 const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 for(const row of ['front','rear'])for(const side of ['left','right']){
  const hinge=model.getObjectByName(`car-door-${row}-${side}`),seat=model.getObjectByName(`seat-${row}-${side}`);
  assert.ok(hinge&&seat,'door and seat must have stable independent nodes');
  model.updateMatrixWorld(true);
  const center=seat.getWorldPosition(new THREE.Vector3()),sign=side==='left'?-1:1;
  const ray=new THREE.Raycaster(new THREE.Vector3(sign*2,.95,center.z),new THREE.Vector3(-sign,0,0),0,1.45);
  assert.ok(ray.intersectObject(model,true).length,'closed door must close the cabin');
  hinge.rotation.y=-sign*Math.PI/2;model.updateMatrixWorld(true);
  assert.equal(ray.intersectObject(model,true).length,0,'opening must not retain solid chassis or glazing across entry');
  hinge.rotation.y=0;
 }
});
