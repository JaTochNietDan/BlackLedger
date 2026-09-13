import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';

test('authored tavern door opens onto a supported person-sized vestibule',async()=>{
 const bytes=readFileSync(new URL('../public/art/models/tavern.glb',import.meta.url));
 const loader=new GLTFLoader();loader.register(p=>({name:'entry-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 const model=(await loader.parseAsync(bytes.buffer.slice(bytes.byteOffset,bytes.byteOffset+bytes.byteLength),'')).scene;
 const hinge=model.getObjectByName('entrance-door-hinge');assert.ok(hinge);
 assert.ok(model.getObjectByName('entrance-threshold'));
 const ray=new THREE.Raycaster();
 const blocked=(x,y)=>{ray.set(new THREE.Vector3(x,y,-7),new THREE.Vector3(0,0,1));ray.far=3.2;return ray.intersectObject(model,true).length>0;};
 model.updateMatrixWorld(true);
 assert.equal(blocked(0,1.5),true,'closed door blocks entry');
 hinge.rotation.y=-Math.PI/2;model.updateMatrixWorld(true);
 for(const x of [-.55,0,.55])for(const y of [.2,.6,1,1.5,1.95])assert.equal(blocked(x,y),false,`open passage blocked at ${x},${y}`);
 for(const z of [-6,-5.5,-5,-4.5]){
   ray.set(new THREE.Vector3(0,1,z),new THREE.Vector3(0,-1,0));ray.far=2;
   const floor=ray.intersectObject(model,true)[0];assert.ok(floor,`floor missing at ${z}`);assert.ok(Math.abs(floor.point.y)<.01);
 }
 const bounds=new THREE.Box3().setFromObject(model,true);
 assert.ok(bounds.max.x-bounds.min.x<=17 && bounds.max.z-bounds.min.z<=17,'open leaf exceeds parcel');
});
