import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {PresentedDie} from '../.runtime/frontend-test/dicePresentation.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial());}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('every committed dice pair lands on the authored pips, with no tray penetration or pair overlap',async()=>{
 const model=await load('gaming-die');let one=0;
 for(let face=1;face<=6;face++){
  const group=model.getObjectByName(`die-face-${face}`);assert.ok(group);let vertices=0;
  group.traverse(o=>{if(o.isMesh)vertices+=o.geometry.attributes.position.count;});
  if(face===1)one=vertices;assert.equal(vertices,one*face,'pip count differs from face label');
 }
 const objects=[model.clone(true),model.clone(true)],dice=objects.map((o,i)=>new PresentedDie(o,i));
 for(let a=1;a<=6;a++)for(let b=1;b<=6;b++){
  for(let frame=0;frame<=120;frame++){
   dice[0].pose(a,frame/120);dice[1].pose(b,frame/120);
   const boxes=objects.map(o=>new THREE.Box3().setFromObject(o,true));
   for(const box of boxes){assert.ok(box.min.y>=.0259,'die passes through baize');assert.ok(box.min.x> -1.205&&box.max.x<1.205&&box.min.z>-.70&&box.max.z<.70,'die leaves inner tray');}
   assert.equal(boxes[0].intersectsBox(boxes[1]),false,'dice overlap');
   if(frame===120)for(let i=0;i<2;i++){
    assert.ok(Math.abs(boxes[i].min.y-.026)<1e-6,'settled die floats');
    const face=i===0?a:b,normal=objects[i].getObjectByName(`die-face-${face}`).getWorldPosition(new THREE.Vector3()).sub(objects[i].position).normalize();
    assert.ok(normal.y>.99999,`wrong upper face for ${face}`);
   }
  }
 }
 dice[0].pose(0,1);assert.equal(objects[0].visible,false);dice[0].pose(7,1);assert.equal(objects[0].visible,false);
});
