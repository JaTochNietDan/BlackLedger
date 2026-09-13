import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {gunfightPose,sceneSlots,casualtySceneStart} from '../.runtime/frontend-test/city3dEvents.js';
import {trafficSize} from '../.runtime/frontend-test/city3dTraffic.js';
test('gunfire raises, recoils four times and lowers at common frame rates',()=>{
 assert.ok(gunfightPose(0).arm===0);
 assert.ok(gunfightPose(3).arm===0);
 assert.equal(gunfightPose(.5).flash,false);
 assert.ok(gunfightPose(.7).arm<-Math.PI/2);
 for(const fps of [30,60,144]){
  let pulses=0,was=false;
  for(let f=0;f<=fps*3;f++){
   const pose=gunfightPose(f/fps);
   if(pose.flash&&!was)pulses++;
   was=pose.flash;
   assert.ok(Number.isFinite(pose.arm));
   assert.ok(pose.smoke>=0&&pose.smoke<=1);
  }
  assert.equal(pulses,4,`${fps} FPS`);
 }
});
test('actual shooter and revolver geometry stay within the reserved scene space',async()=>{
 const load=async name=>{
  const bytes=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));
  return (await new GLTFLoader().parseAsync(bytes.buffer.slice(bytes.byteOffset,bytes.byteOffset+bytes.byteLength),'')).scene;
 };
 const person=await load('person'),gun=await load('revolver');
 const arm=person.getObjectByName('arm1'),muzzle=gun.getObjectByName('muzzle');
 assert.ok(arm);assert.ok(muzzle);
 arm.add(gun);gun.position.set(0,-.58,0);gun.rotation.x=Math.PI/2;
 person.position.y=.2;person.rotation.y=Math.PI/2;
 const slot=sceneSlots({id:'qa',x:0,z:16,row:0,col:0,model:'shop'},'gunfight')[0];
 const size=trafficSize(slot.model);
 for(let frame=0;frame<=180;frame++){
  const pose=gunfightPose(frame/60);arm.rotation.x=pose.arm;person.updateMatrixWorld(true);
  const box=new THREE.Box3().setFromObject(person,true);
  assert.ok(box.min.y>=.17,`floor ${frame}`);
  assert.ok(box.min.x>=.8-size.width/2&&box.max.x<=.8+size.width/2,`width ${frame}`);
  assert.ok(box.min.z>=-size.length/2&&box.max.z<=size.length/2,`depth ${frame}`);
  if(pose.flash){
   const point=muzzle.getWorldPosition(new THREE.Vector3());
   assert.ok(point.x>.8&&point.y>1.4,'flash is ahead of the raised hand');
  }
 }
});

test('casualties wait for co-located gunfire, including a queued shooter',()=>{
 for(const fps of [30,60,144]){
  let since=0;
  for(let frame=0;frame<fps*2;frame++){
   const now=frame*1000/fps;
   since=casualtySceneStart(since,now,500);
   if(now<1200)assert.equal(since,now);
   else assert.ok(now-since>=now-1200);
  }
  assert.ok(since>=1200-1000/fps&&since<1200);
 }
 assert.equal(casualtySceneStart(100,500),100);
});
test('audio follows visible pulses once, and disposal cancels the current tail',async()=>{
 const {GunfireAudio}=await import('../.runtime/frontend-test/city3dEvents.js');
 for(const fps of [30,60,144]){
  let fired=0,stopped=0;
  const audio=new GunfireAudio(()=>{fired++;return()=>stopped++;});
  for(let f=0;f<=fps*3;f++)audio.update(f/fps);
  assert.equal(fired,4);assert.equal(audio.started,4);assert.equal(stopped,3);
  audio.dispose();audio.dispose();audio.update(4);assert.equal(stopped,4);assert.equal(fired,4);
 }
});
test('queued, skipped, muted and missed gunfire does not play a backlog',async()=>{
 const {GunfireAudio}=await import('../.runtime/frontend-test/city3dEvents.js');
 let fired=0,muted=true;
 const make=()=>new GunfireAudio(()=>{if(muted)return;fired++;return()=>{};});
 const skipped=make();for(let f=0;f<100;f++)skipped.update(0);
 skipped.dispose();skipped.update(.7);assert.equal(fired,0);
 const audio=make();audio.update(.7);muted=false;audio.update(.73);assert.equal(fired,0);
 audio.update(1.05);assert.equal(fired,1);
 audio.update(2.5);audio.update(2.8);assert.equal(fired,1);
 assert.equal(audio.started,1);audio.dispose();
});

test('muting stops a current tail and does not replay it on unmute',async()=>{
 const {GunfireAudio}=await import('../.runtime/frontend-test/city3dEvents.js');
 let stops=0,fired=0;
 const audio=new GunfireAudio(()=>{fired++;return()=>stops++;});
 audio.update(.7);audio.update(.71,false);assert.equal(stops,1);
 audio.update(.72,true);assert.equal(fired,1);
 audio.update(1.05,true);assert.equal(fired,2);audio.dispose();assert.equal(stops,2);
});
