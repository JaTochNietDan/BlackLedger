import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';

test('Saint Agnes interior has supported clear standing bays and bounded room geometry',async()=>{
 const b=readFileSync(new URL('../public/art/models/interior-saint-agnes.glb',import.meta.url));
 const loader=new GLTFLoader();loader.register(p=>({name:'interior-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 const room=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));
 assert.ok(room.getObjectByName('interior-wall-back'));
 const bounds=new THREE.Box3().setFromObject(room,true);
 assert.ok(bounds.min.x>=-6.11&&bounds.max.x<=6.01&&bounds.max.y<=3.51);
 for(let i=0;i<9;i++){
  const x=-1+(i%3)*2,z=-Math.floor(i/3)*2;
  // Sample under the shoe footprint, away from the central grout crossing.
  const floor=new THREE.Raycaster(new THREE.Vector3(x+.12,.1,z+.12),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
  assert.ok(floor&&Math.abs(floor.point.y-.028)<.002,'standing bay must meet floor');
  for(const y of [.2,.6,1,1.5,1.85])for(let angle=0;angle<16;angle++){
   const direction=new THREE.Vector3(Math.cos(angle*Math.PI/8),0,Math.sin(angle*Math.PI/8));
   assert.equal(new THREE.Raycaster(new THREE.Vector3(x,y,z),direction,0,.6).intersectObject(room,true).length,0,`furniture intrudes into bay ${i} at height ${y}`);
  }
 }
});
