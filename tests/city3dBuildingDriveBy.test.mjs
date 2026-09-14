import test from 'node:test';import assert from 'node:assert/strict';import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';import {readFileSync} from 'node:fs';
import {CityBuildingDriveBy,buildingDriveByPose,BUILDING_DRIVEBY_SECONDS} from '../.runtime/frontend-test/city3dBuildingDriveBy.js';
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
