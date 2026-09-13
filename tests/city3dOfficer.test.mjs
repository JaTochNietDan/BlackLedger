import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';

test('authored uniformed officer fits the reserved standing footprint at every heading',async()=>{
 const b=readFileSync(new URL('../public/art/models/police-officer.glb',import.meta.url));
 const loader=new GLTFLoader();loader.register(p=>({name:'officer-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 for(let i=0;i<32;i++){
  model.rotation.y=i*Math.PI/16;model.updateMatrixWorld(true);const box=new THREE.Box3().setFromObject(model,true);
  assert.ok(box.min.x>=-.5&&box.max.x<=.5&&box.min.z>=-.5&&box.max.z<=.5,'uniform or equipment exceeds standing clearance');
  assert.ok(box.min.y>=-.001&&box.max.y<2,'feet must meet ground and cap remain at human height');
 }
});
