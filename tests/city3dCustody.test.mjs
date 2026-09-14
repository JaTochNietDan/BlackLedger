import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityCustody} from '../.runtime/frontend-test/city3dCustody.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(()=>({name:'custody-geometry',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial());}}));
 return(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('paired custody keeps the prisoner still, approaches before contact and fits its complete reservation',async()=>{
 for(const name of ['person','woman']){
  const detainee=await model(name),officer=await model('police-officer'),cuffs=await model('handcuffs');cuffs.name='custody-restraint';detainee.add(cuffs);
  const scene=new CityCustody(detainee,officer);scene.root.position.y=.2;
  let previous=0;
  for(let frame=0;frame<=360;frame++){
   const seconds=frame/60,p=scene.update(seconds),bounds=new THREE.Box3().setFromObject(scene.root,true);
   assert.ok(p.distance>=previous&&p.distance-previous<=1.25/60+.000001);previous=p.distance;
   assert.equal(detainee.position.x,3);assert.ok(officer.position.x<=2.42);
   assert.ok(bounds.min.x>=-.7&&bounds.max.x<=3.7&&bounds.min.z>=-.7&&bounds.max.z<=.7,'cast exceeds reserved corridor');
   assert.ok(bounds.min.y>=.17,'cast enters pavement');
   if(seconds<2.79)assert.equal(cuffs.visible,false);
   if(seconds>=2.81)for(const side of [-1,1]){
    const hand=officer.getObjectByName(`elbow${side}`).localToWorld(new THREE.Vector3(0,-.295,0));
    const wrist=detainee.getObjectByName(`elbow${side}`).localToWorld(new THREE.Vector3(0,-.295,0));
    assert.ok(hand.distanceTo(wrist)<.035,`officer misses wrist: ${hand.distanceTo(wrist)}`);
   }
  }
 }
});
