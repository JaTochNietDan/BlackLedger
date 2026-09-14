import test from 'node:test';import assert from 'node:assert/strict';import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';import {readFileSync} from 'node:fs';
import {CityBuildingDriveBy,buildingDriveByCondition,buildingDriveByShots,buildingDriveByPose,BUILDING_DRIVEBY_SECONDS} from '../.runtime/frontend-test/city3dBuildingDriveBy.js';
import {sceneSlots,availableSceneSlot} from '../.runtime/frontend-test/city3dEvents.js';
import {StreetTraffic,trafficOverlap,trafficSize} from '../.runtime/frontend-test/city3dTraffic.js';
const models=new Map();
async function load(name){if(!models.has(name)){const b=readFileSync(`public/art/models/${name}.glb`),l=new GLTFLoader();l.register(parser=>({name:'driveby-geometry',loadMaterial(index){const m=new THREE.MeshStandardMaterial();m.name=parser.json.materials[index].name;return Promise.resolve(m);}}));models.set(name,(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene);}return models.get(name).clone(true);}
test('drive-by slows for gunfire then accelerates with continuous speed',()=>{
 for(const fps of [30,60,144]){let old=buildingDriveByPose(0);for(let i=1;i<=Math.ceil(BUILDING_DRIVEBY_SECONDS*fps);i++){const p=buildingDriveByPose(i/fps);assert.ok(p.distance>=old.distance);assert.ok((p.distance-old.distance)*fps<=6.001);old=p;}}
 for(const t of [2,3.8]){assert.ok(Math.abs(buildingDriveByPose(t-.00001).speed-buildingDriveByPose(t+.00001).speed)<.0001);}
 assert.equal(buildingDriveByPose(-10).distance,0);assert.equal(buildingDriveByPose(100).distance,buildingDriveByPose(BUILDING_DRIVEBY_SECONDS).distance);
});
for(const carName of ['ford','hudson','packard'])for(const rig of ['person','woman'])for(const gun of ['revolver','shotgun','thompson'])test(`${carName}/${rig}/${gun}: full-size seated cast and firing grips`,async()=>{
 const c=new CityBuildingDriveBy(await load(carName),await load('person'),await load(rig),await load(gun),gun);
 const pane=[];c.car.getObjectByName('car-door-front-right').traverse(o=>{if(o.isMesh&&o.material.name==='car glass')pane.push(o);});
 assert.ok(pane.length>0);assert.ok(pane.every(o=>!o.visible),'passenger fires through a closed window');
 for(const time of [0,...c.shots,6.2]){
  c.update(time);
  const slot=sceneSlots({x:0,row:0},'driveby-building')[0],size=trafficSize(slot.model);
  // Local model space; the reservation's root will be translated onto the lane.
  const bounds=new THREE.Box3();c.root.traverseVisible(m=>{if(m.isMesh){const p=m.geometry.attributes.position;for(let i=0;i<p.count;i++)bounds.expandByPoint(c.root.worldToLocal(m.localToWorld(new THREE.Vector3().fromBufferAttribute(p,i))));}});
  assert.ok(bounds.min.x>=.4-size.length/2&&bounds.max.x<=.4+size.length/2,'car/cast leaves swept road length');
  assert.ok(bounds.min.z>=.4-size.width/2&&bounds.max.z<=.4+size.width/2,'car/cast leaves reserved lane width');
  for(const [actor,seat]of[[c.driver,'seat-front-left'],[c.shooter,'seat-front-right']]){
   const pelvis=c.car.worldToLocal(actor.localToWorld(new THREE.Vector3(0,.86,0)));
   assert.ok(pelvis.distanceTo(c.car.getObjectByName(seat).position)<1e-6,'pelvis leaves cushion');assert.deepEqual(actor.scale.toArray(),[1,1,1]);
   for(const side of[-1,1]){
    const leg=actor.getObjectByName('leg'+side),bounds=new THREE.Box3();
    const joint=c.car.worldToLocal(leg.getWorldPosition(new THREE.Vector3()));
    const seatPoint=c.car.getObjectByName(seat).position.clone().add(new THREE.Vector3(side*.12,0,0));
    assert.ok(joint.distanceTo(seatPoint)<1e-6,'legs turn across one another with the torso');
    leg.traverseVisible(m=>{if(m.isMesh){const p=m.geometry.attributes.position;for(let i=0;i<p.count;i++)bounds.expandByPoint(c.car.worldToLocal(m.localToWorld(new THREE.Vector3().fromBufferAttribute(p,i))));}});
    assert.ok(bounds.min.y>=.15&&bounds.max.x<=.745&&bounds.min.x>=-.745,'legs leave the cabin');
   }
  }
  if(c.shots.includes(time))for(const side of (gun==='revolver'?[1]:[-1,1])){
   const elbow=c.shooter.getObjectByName('elbow'+side),hand=elbow.localToWorld(new THREE.Vector3(0,-.295,0));
   const grip=(side<0?c.weapon.getObjectByName('support-grip'):c.weapon).getWorldPosition(new THREE.Vector3());
   assert.ok(hand.distanceTo(grip)<.025,`hand misses ${side} grip at ${time}: ${hand.distanceTo(grip)}`);
  }
  if(c.shots.includes(time)){
   const grip=c.car.worldToLocal(c.weapon.getWorldPosition(new THREE.Vector3()));
   const muzzle=c.car.worldToLocal(c.weapon.getObjectByName('muzzle').getWorldPosition(new THREE.Vector3()));
   const along=(.745-grip.x)/(muzzle.x-grip.x),window=grip.clone().lerp(muzzle,along);
   assert.ok(muzzle.x>.745,'barrel stays inside the car');
   assert.ok(window.y>1.055&&window.y<1.425,`barrel clips window sill/roof at ${time}: ${window.y}`);
   assert.ok(window.z>-.06&&window.z<.60,`barrel crosses a window pillar at ${time}: ${window.z}`);
  }
 }
 c.dispose();
});

test('drive-by waits for its actual road lane without a pavement fallback',()=>{
 const lot={x:48,row:1,z:48,col:1,id:'club'},slot=sceneSlots(lot,'driveby-building')[0];
 assert.equal(slot.root.z,33.6);
 assert.equal(availableSceneSlot(lot,'driveby-building',[{pose:{x:48,z:33.6,heading:Math.PI/2},model:'packard'}]),undefined);
 assert.ok(availableSceneSlot(lot,'driveby-building',[]));
 // Sidewalk pedestrians and the other road lane remain outside the footprint.
 for(const other of [{pose:{x:48,z:36.65,heading:Math.PI/2},model:'person'},{pose:{x:48,z:30.4,heading:Math.PI/2},model:'packard'}])
  assert.equal(trafficOverlap(slot.pose,slot.model,other.pose,other.model),false);
});

test('pending drive-by lets current cars exit, holds incoming cars, then releases them',()=>{
 const lot={x:48,row:1,z:48,col:1,id:'club'},space=sceneSlots(lot,'driveby-building')[0];
 const traffic=new StreetTraffic(),points=[{x:85,z:33.6},{x:20,z:33.6}];
 const inside={id:'inside',model:'packard',points,progress:.50};
 traffic.update([inside],1/60);
 const behind={id:'behind',model:'ford',points,progress:0};
 traffic.update([inside,behind],1/60,1,[space]);inside.progress=behind.progress=1;
 let poses;
 for(let i=0;i<1200;i++){
  poses=traffic.update([inside,behind],1/60,1,[space]);
  const b=poses.get('behind');if(!b.waiting)assert.equal(trafficOverlap(b.pose,behind.model,space.pose,space.model),false);
 }
 assert.equal(poses.get('inside').progress,1,'existing occupant cannot finish leaving');
 assert.ok(poses.get('behind').progress<.4,'incoming car enters scene');
 assert.ok(availableSceneSlot(lot,'driveby-building',[...poses].filter(([,p])=>!p.waiting).map(([id,p])=>({pose:p.pose,model:id==='inside'?'packard':'ford'}))));
 for(let i=0;i<1200;i++)poses=traffic.update([inside,behind],1/60);
 assert.ok(poses.get('behind').progress>.85,'car fails to resume after release');
});

test('admitted drive-by keeps its swept orientation in traffic',()=>{
 const slot=sceneSlots({x:48,row:1},'driveby-building')[0];
 const traffic=new StreetTraffic();
 const request={id:'scene:driveby',model:slot.model,points:[slot.pose],progress:0};
 let poses=traffic.update([request],1/60);
 assert.equal(poses.get(request.id).pose.heading,slot.pose.heading);
 // These two points distinguish a correct east-west sweep from a rotated one.
 assert.equal(trafficOverlap(poses.get(request.id).pose,slot.model,{x:59,z:33.6,heading:Math.PI/2},'ford'),true);
 assert.equal(trafficOverlap(poses.get(request.id).pose,slot.model,{x:48,z:43,heading:0},'person'),false);
 const moving={id:'incoming',model:'ford',points:[{x:80,z:33.6},{x:20,z:33.6}],progress:0};
 traffic.update([request,moving],1/60);moving.progress=1;
 for(let i=0;i<600;i++){
  poses=traffic.update([request,moving],1/60);
  const car=poses.get(moving.id);
  if(!car.waiting)assert.equal(trafficOverlap(car.pose,moving.model,slot.pose,slot.model),false,'traffic entered the actual swept footprint');
 }
 assert.ok(poses.get(moving.id).progress<.4);
});

for(const gun of ['revolver','shotgun','thompson'])test(`${gun}: damage reveals on shots and respects captured endpoints`,()=>{
 for(const [before,after]of[[100,72],[65,40],[4,0],[0,0]]){
  const beats=buildingDriveByShots(gun);
  assert.equal(buildingDriveByCondition(before,after,-1,gun),before);
  assert.equal(buildingDriveByCondition(before,after,beats[0]-.001,gun),before);
  let last=before;
  for(const beat of beats){const shown=buildingDriveByCondition(before,after,beat,gun);assert.ok(shown<=last&&shown>=after);last=shown;}
  assert.equal(last,after);assert.equal(buildingDriveByCondition(before,after,100,gun),after);
 }
});
