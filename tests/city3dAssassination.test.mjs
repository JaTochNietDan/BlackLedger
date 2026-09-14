import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityAssassination,assassinationPose,assassinationBatch,ASSASSINATION_SHOT,ASSASSINATION_SECONDS,executionSpatter} from '../.runtime/frontend-test/city3dAssassination.js';
import {sceneSlots,availableSceneSlot} from '../.runtime/frontend-test/city3dEvents.js';
import {trafficSize} from '../.runtime/frontend-test/city3dTraffic.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshStandardMaterial());}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('execution links only its recorded victim and preserves unrelated or legacy cues',()=>{
 const gun={id:'g',kind:'gunfight',target:'bar',minute:480,attacker:{id:'player',weapon:1},strike:{variant:'back-of-head',victim:{id:'mara'}}};
 const death={id:'d',kind:'killing',target:'bar',minute:480,actors:[{id:'mara'}]};
 const stranger={...death,id:'other',actors:[{id:'leo'}]};
 assert.deepEqual(assassinationBatch([death,gun,stranger]),[gun,stranger]);
 assert.deepEqual(assassinationBatch([death,{...gun,minute:500}]).length,2);
 assert.deepEqual(assassinationBatch([death,{...gun,attacker:undefined}]).length,2);
});
test('both cast rigs, muzzle and droplet trajectories stay inside the full scene reserve',async()=>{
 for(const attackerModel of ['person','woman'])for(const victimModel of ['person','woman']){
  const actor=await model(attackerModel),victim=await model(victimModel),gun=await model('revolver');
  const scene=new CityAssassination(actor,victim,gun);scene.root.position.set(76,.2,38.35);
  for(let frame=0;frame<=Math.ceil(ASSASSINATION_SECONDS*60);frame++){
   const seconds=frame/60,p=scene.update(seconds),bounds=new THREE.Box3().setFromObject(scene.root,true);
   assert.ok(bounds.min.x>=75.3&&bounds.max.x<=83.5,`${attackerModel}/${victimModel} x extent ${bounds.min.x}..${bounds.max.x} @${seconds}`);
   assert.ok(bounds.min.z>=37.65&&bounds.max.z<=39.05,`z extent ${bounds.min.z}..${bounds.max.z} @${seconds}`);
   assert.ok(bounds.min.y>=.17,`pavement penetration ${bounds.min.y} @${seconds}`);
   // The alternate left bay passes the lamppost; verify actual cast clearance.
   const leftShift=4.6;
   const shifted=bounds.clone().translate(new THREE.Vector3(-leftShift,0,0));
   const pole=new THREE.Box3(new THREE.Vector3(71.14,.1,37.44),new THREE.Vector3(71.46,5.1,37.76));
   assert.ok(!shifted.intersectsBox(pole),'cast intersects street lamp');
   const elbow=actor.getObjectByName('elbow1'),hand=elbow.localToWorld(new THREE.Vector3(0,-.295,0));
   assert.ok(hand.distanceTo(gun.getWorldPosition(new THREE.Vector3()))<.015,'hand detaches from gun');
   if(seconds<ASSASSINATION_SHOT)assert.ok(p.fall.rotation===0,'victim anticipates shot');
   for(let i=0;i<30;i++){const d=executionSpatter(i,seconds);if(d.size){assert.ok(d.x+d.size<7.5&&Math.abs(d.z)+d.size<.7);assert.ok(d.y>.0);}}
  }
  scene.update(ASSASSINATION_SHOT-.001);
  const muzzle=gun.getObjectByName('muzzle'),origin=muzzle.getWorldPosition(new THREE.Vector3());
  const direction=new THREE.Vector3(0,0,1).applyQuaternion(gun.getWorldQuaternion(new THREE.Quaternion()));
  const hits=new THREE.Raycaster(origin,direction,0,.35).intersectObject(victim,true);
  assert.ok(hits.length,'muzzle does not aim into victim head');
  assert.ok(hits[0].point.y>1.7&&hits[0].point.y<2,'shot misses head height');
  assert.ok(origin.x<81,'muzzle penetrates target before firing');
  assert.ok(assassinationPose(3.5).distance>4&&assassinationPose(3.5).aim===1);
 }
 const slot=sceneSlots({x:80,row:1,z:48,col:2,id:'bar'},'assassination')[0];
 assert.equal(slot.root.x,76);assert.deepEqual(trafficSize(slot.model),{length:1.4,width:8.2});
});


test('a second scene can reserve its approach beside an existing body',()=>{
 const lot={x:80,row:1,z:48,col:2,id:'bar'};
 const body={pose:{x:81.8,z:38.35,heading:0},model:'casualty'};
 const slot=availableSceneSlot(lot,'assassination',[body]);
 assert.ok(slot);assert.equal(slot.root.x,71.4);
 // Whole reservation stays on the forecourt, ahead of the maximum 17m facade.
 for(const s of sceneSlots(lot,'assassination')){
  assert.ok(s.pose.z-.7>32+4.65+.7+.2);
  assert.ok(s.pose.z+.7<48-8.5);
  assert.ok(s.pose.x-4.1>64+4);
 }
});

test('close-quarters cast approaches before blows and falls only after the final impact',async()=>{
 const {meleePose,MELEE_IMPACTS,isStagedStrike}=await import('../.runtime/frontend-test/city3dAssassination.js');
 const cue={id:'attack',kind:'attack',target:'bar',minute:600,attacker:{id:'player',weapon:0},strike:{variant:'close-quarters',victim:{id:'victim'}}};
 const death={id:'death',kind:'killing',target:'bar',minute:600,actors:[{id:'victim'}]};
 assert.ok(isStagedStrike(cue));
 assert.deepEqual(assassinationBatch([death,cue]),[cue]);
 assert.equal(isStagedStrike({...cue,attacker:{weapon:1}}),false);
 for(const name of ['person','woman']){
  const attacker=await model(name),victim=await model('person');
  const cast=new CityAssassination(attacker,victim);cast.root.position.y=.2;
  for(let f=0;f<=390;f++){
   const t=f/60;cast.update(t);
   const b=new THREE.Box3().setFromObject(cast.root,true);
   assert.ok(b.min.y>=.17,`pavement penetration at ${t}: ${b.min.y}`);
   assert.ok(b.min.x>=-.7&&b.max.x<=7.5&&b.min.z>=-.7&&b.max.z<=.7,`cast leaves reserved forecourt at ${t}`);
   if(t<MELEE_IMPACTS[2])assert.ok(meleePose(t).fall.rotation===0);
  }
  for(const at of MELEE_IMPACTS){
   cast.update(at);
   const hand=attacker.getObjectByName('elbow1').localToWorld(new THREE.Vector3(0,-.295,0));
   assert.ok(hand.x>4.75&&hand.x<5.1,`blow misses victim: ${hand.x}`);
   assert.ok(hand.y>1.5&&hand.y<1.9,`blow misses upper body: ${hand.y}`);
  }
  assert.equal(cast.weapon,undefined);
 }
});
