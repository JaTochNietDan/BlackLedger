import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const root=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;root.updateMatrixWorld(true);return root;
}
test('blackjack cards rest above the felt and expose upright printed faces',async()=>{
 const table=await model('blackjack-table'),card=await model('playing-card');
 for(const z of [-.38,.48])for(const x of [-.85,0,.85]){
  card.scale.set(1.4,1,1.4);card.position.set(x,.875,z);card.updateMatrixWorld(true);
  const box=new THREE.Box3().setFromObject(card,true);
  for(const cx of [box.min.x,box.max.x])for(const cz of [box.min.z,box.max.z]){
   const hit=new THREE.Raycaster(new THREE.Vector3(cx,2,cz),new THREE.Vector3(0,-1,0)).intersectObject(table,true)[0];
   assert.ok(hit);assert.equal(hit.object.material.name,'blackjack green baize');assert.ok(box.min.y>hit.point.y,'card intersects cloth');
  }
  const hit=new THREE.Raycaster(new THREE.Vector3(x,2,z-.1),new THREE.Vector3(0,-1,0)).intersectObject(card,true)[0];
  assert.equal(hit.object.material.name,'card printed face');assert.ok(hit.uv.y<.5,'card face is upside down');
 }
});
