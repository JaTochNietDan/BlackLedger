import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {sceneWeapon,poseLongGun,weaponShots,pumpOffset} from '../.runtime/frontend-test/city3dWeapons.js';
import {gunfightPose,GunfireAudio} from '../.runtime/frontend-test/city3dEvents.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'weapon-geometry',loadMaterial(){return Promise.resolve(new THREE.MeshStandardMaterial());}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('recorded weapon tiers select equipment and preserve anonymous legacy gunfire',()=>{
 assert.equal(sceneWeapon(), 'revolver');assert.equal(sceneWeapon(0),null);assert.equal(sceneWeapon(1),'revolver');assert.equal(sceneWeapon(2),'shotgun');assert.equal(sceneWeapon(3),'thompson');assert.equal(sceneWeapon(99),null);
});
test('authored long guns have muzzle/grip anchors and both rigid arms reach them throughout playback',async()=>{
 for(const personName of ['person','woman'])for(const gunName of ['shotgun','thompson']){
  const person=await model(personName),gun=await model(gunName);
  gun.updateMatrixWorld(true);const bounds=new THREE.Box3().setFromObject(gun,true),size=bounds.getSize(new THREE.Vector3());
  assert.ok(size.x<.15&&size.y<.4&&size.z<1.2,`${gunName} dimensions`);
  assert.ok(gun.getObjectByName('muzzle')&&gun.getObjectByName('support-grip'));
  person.add(gun);person.rotation.y=Math.PI/2;person.position.set(80,.2,38.35);
  for(let frame=0;frame<=180;frame++){
   const pump=gun.getObjectByName('pump-slide');if(pump)pump.position.z=pumpOffset(frame/60);
   poseLongGun(person,gun,gunfightPose(frame/60,weaponShots(gunName)).arm);
   for(const side of [1,-1]){
    const elbow=person.getObjectByName(`elbow${side}`);assert.ok(elbow);
    const hand=elbow.localToWorld(new THREE.Vector3(0,-.295,0));
    const grip=(side===1?gun:gun.getObjectByName('support-grip')).getWorldPosition(new THREE.Vector3());
    assert.ok(hand.distanceTo(grip)<.015,`${personName}/${gunName} hand ${side} gap ${hand.distanceTo(grip)} frame ${frame}`);
   }
   const box=new THREE.Box3().setFromObject(gun,true);
   const cast=new THREE.Box3().setFromObject(person,true);
   assert.ok(cast.min.x>=79.5&&cast.max.x<=82.1&&cast.min.z>=37.65&&cast.max.z<=39.05,`${personName} aiming pose exceeds scene clearance`);
   assert.ok(box.min.y>.17,`${gunName} penetrates pavement`);
   assert.ok(box.min.x>79.5&&box.max.x<82.1&&box.min.z>37.65&&box.max.z<39.05,`${gunName} exceeds scene clearance`);
  }
 }
});

test('weapon cadence uses two pumped shots or two automatic bursts with one audio onset per beat',async()=>{
 const {GunfireAudio}=await import('../.runtime/frontend-test/city3dEvents.js');
 for(const kind of ['revolver','shotgun','thompson']){
  const beats=weaponShots(kind);let played=0;const audio=new GunfireAudio(()=>{played++;return ()=>{};},beats);
  for(let frame=0;frame<180;frame++)audio.update(frame/60);
  assert.equal(played,beats.length);audio.dispose();
 }
 assert.equal(weaponShots('shotgun').length,2);assert.equal(weaponShots('thompson').length,6);
 assert.equal(pumpOffset(.7),0);assert.ok(pumpOffset(1.1)<-.08);assert.equal(pumpOffset(1.7),0);
});


test('recorded back-of-head strikes fire exactly once, including late playback',()=>{
 const beats=weaponShots('revolver','back-of-head');assert.deepEqual(beats,[3.65]);
 let count=0;const audio=new GunfireAudio(()=>{count++;return ()=>{};},beats);
 for(let i=0;i<390;i++)audio.update(i/60);
 assert.equal(count,1);audio.dispose();
 let late=0;const skipped=new GunfireAudio(()=>{late++;return ()=>{};},beats);
 skipped.update(4.5);skipped.update(5);assert.equal(late,0);skipped.dispose();
});
