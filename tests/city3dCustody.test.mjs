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
test('custody approaches, cuffs and escorts both cast members within the reserved path',async()=>{
 for(const name of ['person','woman']){
  const detainee=await model(name),officer=await model('police-officer'),cuffs=await model('handcuffs');cuffs.name='custody-restraint';detainee.add(cuffs);
  const scene=new CityCustody(detainee,officer);scene.root.position.y=.2;
  let previous=0;
  for(let frame=0;frame<=540;frame++){
   const seconds=frame/60,p=scene.update(seconds),bounds=new THREE.Box3().setFromObject(scene.root,true);
   assert.ok(p.distance>=previous&&p.distance-previous<=1.25/60+.000001);previous=p.distance;
   if(seconds<5.5)assert.equal(detainee.position.x,3);
   assert.ok(officer.position.x<=3);
   if(seconds>=5.5){assert.ok(Math.abs(detainee.position.x-officer.position.x)<.001);assert.ok(Math.abs(officer.position.z-.75)<.001);assert.ok(detainee.rotation.y< -1.56);}
   assert.ok(bounds.min.x>=-.7&&bounds.max.x<=3.7&&bounds.min.z>=-.7&&bounds.max.z<=1.3,'cast exceeds reserved corridor');
   assert.ok(bounds.min.y>=.17,'cast enters pavement');
   if(seconds<2.79)assert.equal(cuffs.visible,false);
   if(seconds>=2.81&&seconds<=2.95)for(const side of [-1,1]){
    const hand=officer.getObjectByName(`elbow${side}`).localToWorld(new THREE.Vector3(0,-.295,0));
    const wrist=detainee.getObjectByName(`elbow${side}`).localToWorld(new THREE.Vector3(0,-.295,0));
    assert.ok(hand.distanceTo(wrist)<.035,`officer misses wrist: ${hand.distanceTo(wrist)}`);
   }
   if(seconds>=5.5){
    const hand=officer.getObjectByName('elbow-1').localToWorld(new THREE.Vector3(0,-.295,0));
    const sleeve=detainee.getObjectByName('arm1').localToWorld(new THREE.Vector3(0,-.08,0));
    assert.ok(hand.distanceTo(sleeve)<.035,'escort loses upper-arm contact');
   }
  }
  assert.ok(detainee.position.x<.71,'prisoner was not escorted away');
 }
});
