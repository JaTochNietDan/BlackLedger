import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';

test('authored slot paper presents each exact payline through its window and lever clears cabinet',async()=>{
 const b=readFileSync(new URL('../public/art/models/slot-cabinet.glb',import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const room=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;room.updateMatrixWorld(true);
 const bound=new THREE.Box3().setFromObject(room,true);assert.ok(bound.min.x>=-.52&&bound.max.x<.72&&bound.max.y<1.72);
 for(const [i,x] of [-.28,0,.28].entries()){
  for(const y of [.975,1.055,1.135]){
   const hit=new THREE.Raycaster(new THREE.Vector3(x,y,2),new THREE.Vector3(0,0,-1),0,3).intersectObject(room,true)[0];
   assert.ok(hit,'window is missing');assert.equal(hit.object.material.name,`slot-reel-${i}`,'case obscures reel');
   if(y===1.055)assert.ok(Math.abs(hit.uv.y-.5)<1e-5&&Math.abs(hit.uv.x-.5)<1e-5,'payline misses middle printed symbol');
   if(y>1.055)assert.ok(hit.uv.y<.5,'reel print is upside down');
  }
 }
 const lever=room.getObjectByName('slot-lever');assert.ok(lever);
 for(let frame=0;frame<=90;frame++){
  lever.rotation.x=.75*Math.sin(frame/90*Math.PI);room.updateMatrixWorld(true);
  const box=new THREE.Box3().setFromObject(lever,true);assert.ok(box.min.x>.5,'lever crosses cabinet side');
 }
 const coins=room.getObjectByName('slot-coins');assert.ok(coins);
 const tray=new THREE.Box3().setFromObject(coins,true);
 assert.ok(tray.min.x>-.44&&tray.max.x<.24&&tray.min.z>.33&&tray.max.z<.585&&tray.min.y>.228&&tray.max.y<.29,'coins leave the tray');
});
