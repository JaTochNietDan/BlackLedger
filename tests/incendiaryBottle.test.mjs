import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
async function bottle(){
 const b=readFileSync(new URL('../public/art/models/incendiary-bottle.glb',import.meta.url));
 const loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const root=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 root.updateMatrixWorld(true);return root;
}
test('authored bottle has a neck-centred grip, clear label and flame attachment',async()=>{
 const root=await bottle(),bounds=new THREE.Box3().setFromObject(root,true);
 assert.ok(bounds.max.x-bounds.min.x>.079&&bounds.max.x-bounds.min.x<.09,'bottle does not fit a hand');
 assert.ok(Math.abs(bounds.min.y+.25)<.001,'base must be 25cm below grip');
 const grip=root.getObjectByName('bottle-grip'),flame=root.getObjectByName('bottle-flame');
 assert.ok(grip&&flame);assert.ok(grip.getWorldPosition(new THREE.Vector3()).length()<1e-6);
 assert.ok(flame.getWorldPosition(new THREE.Vector3()).y>bounds.max.y,'flame starts inside the cloth');
 const hits=new THREE.Raycaster(new THREE.Vector3(0,-.15,.2),new THREE.Vector3(0,0,-1)).intersectObject(root,true);
 assert.ok(hits.length>=2);assert.equal(hits[0].object.material.name,'bottle aged paper','label is inside the glass or facing inward');
 assert.equal(hits[1].object.material.name,'bottle olive glass');assert.ok(hits[0].point.z-hits[1].point.z>.0003);
 let triangles=0;root.traverse(o=>{if(o.isMesh)triangles+=(o.geometry.index?.count??o.geometry.attributes.position.count)/3;});
 assert.ok(triangles<10000,`bottle has ${triangles} triangles`);
});
