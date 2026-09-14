import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityIncendiary,INCENDIARY_RELEASE,INCENDIARY_IMPACT} from '../.runtime/frontend-test/city3dIncendiary.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('bottle stays at the authored hand, releases continuously and reaches the impact point',async()=>{
 for(const name of ['person','woman']){
  const actor=await model(name),bottle=await model('incendiary-bottle'),target=new THREE.Vector3(0,2.5,4);
  const cast=new CityIncendiary(actor,bottle,target);cast.root.position.set(13,.2,-8);cast.root.rotation.y=.7;
  for(let t=0;t<INCENDIARY_RELEASE;t+=.025){
   cast.update(t);
   const hand=actor.getObjectByName('elbow1').localToWorld(new THREE.Vector3(.01,-.295,.005));
   assert.ok(hand.distanceTo(bottle.getWorldPosition(new THREE.Vector3()))<1e-6,'bottle slips from hand');
   // The substantial body must extend beyond the hand, away from the forearm.
   const body=bottle.localToWorld(new THREE.Vector3(0,-.16,0));
   const elbow=actor.getObjectByName('elbow1').getWorldPosition(new THREE.Vector3());
   const nearest=new THREE.Line3(elbow,hand).closestPointToPoint(body,true,new THREE.Vector3());
   assert.ok(body.distanceTo(nearest)>.14,'body embeds in the forearm');
  }
  cast.update(INCENDIARY_RELEASE-1e-6);const position=bottle.position.clone(),rotation=bottle.quaternion.clone();
  cast.update(INCENDIARY_RELEASE);assert.ok(bottle.position.distanceTo(position)<1e-5);assert.ok(rotation.angleTo(bottle.quaternion)<1e-5);
  for(let t=INCENDIARY_RELEASE;t<INCENDIARY_IMPACT;t+=.01){cast.update(t);assert.equal(bottle.visible,true);assert.ok(bottle.position.y>1.5,'flight drops below the release area');}
  cast.update(INCENDIARY_IMPACT);assert.equal(bottle.visible,false);assert.ok(bottle.position.distanceTo(target)<1e-6);
  cast.update(cast.duration);assert.ok(actor.position.x<=-4);assert.equal(bottle.visible,false);
 }
});
