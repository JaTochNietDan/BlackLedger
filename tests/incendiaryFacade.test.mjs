import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityIncendiary,incendiaryStagingFlight,incendiaryCorridorClear,incendiaryFlight,INCENDIARY_RELEASE,INCENDIARY_IMPACT} from '../.runtime/frontend-test/city3dIncendiary.js';
import {clearBlastWindows} from '../.runtime/frontend-test/city3dFire.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'collision-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 const root=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;root.updateMatrixWorld(true);return root;
}
test('incendiary bottle centre reaches facade without passing through earlier masonry',async()=>{
 const actor=await model('person'),bottle=await model('incendiary-bottle');
 for(const name of ['tavern','monarch','tenement','shop','civic','casino','warehouse','bluehour','goldenlily','papermoon','mariner','mercer-court','filling','garage','dealer','docks','haulage','villa','undertaker']){
  const building=await model(name);building.position.set(0,.18,16);if(['filling','garage','dealer','docks','haulage'].includes(name))building.rotation.y=Math.PI;building.updateMatrixWorld(true);
  for(const x of [3,0,-3]){
   assert.ok(incendiaryCorridorClear({x,z:6.35},building),`${name} x=${x} actor corridor overlaps authored geometry`);
   const origin=new THREE.Vector3(x,.2,6.35),windows=clearBlastWindows(building).sort((a,b)=>Math.abs(a.x-x)-Math.abs(b.x-x));
   assert.ok(windows.length,name);const target=windows[0].clone().sub(origin);
   const cast=new CityIncendiary(actor.clone(true),bottle.clone(true),target);cast.root.position.copy(origin);
   const flight=incendiaryFlight(cast.root.localToWorld(cast.release.clone()),windows,building);
   assert.ok(flight,`${name} x=${x} has no clear throw`);
   cast.target.copy(flight.target).sub(origin);cast.loft=flight.loft;
   let previous;
   for(let frame=0;frame<80;frame++){
    cast.update(INCENDIARY_RELEASE+(INCENDIARY_IMPACT-INCENDIARY_RELEASE)*frame/80);
    const point=cast.bottle.getWorldPosition(new THREE.Vector3());
    if(previous){const direction=point.clone().sub(previous);const hits=new THREE.Raycaster(previous,direction.clone().normalize(),0,direction.length()).intersectObject(building,true);
     assert.equal(hits.length,0,`${name} x=${x} hits ${hits[0]?.object.name} before impact at frame ${frame}`);
    }previous=point;
   }
  }
 }
});

test('industrial fire markers face the street and have real glazing behind them',async()=>{
 for(const name of ['filling','garage','dealer','docks','haulage']){
  const building=await model(name);building.position.set(0,.18,16);building.rotation.y=Math.PI;building.updateMatrixWorld(true);
  const windows=clearBlastWindows(building);assert.equal(windows.length,3,name);
  for(const point of windows){
   const hit=new THREE.Raycaster(point,new THREE.Vector3(0,0,1),0,.25).intersectObject(building,true)[0];
   assert.ok(hit,`${name} marker is not attached to a facade`);
   assert.match(hit.object.material.name,/glazing/,`${name} marker points at solid wall`);
  }
 }
});

test('special facade fire markers sit in front of upper-window glazing',async()=>{
 for(const [name,count] of [['villa',3],['undertaker',4]]){
  const building=await model(name),windows=clearBlastWindows(building);assert.equal(windows.length,count);
  for(const point of windows){
   assert.ok(point.y>4);
   const hit=new THREE.Raycaster(point,new THREE.Vector3(0,0,1),0,.25).intersectObject(building,true)[0];
   assert.ok(hit,`${name} marker lacks backing`);assert.match(hit.object.material.name,/glazing/);
  }
 }
});

test('flight cache avoids repeated raycasts and invalidates changed obstacles',()=>{
 const building=new THREE.Group(),wall=new THREE.Mesh(new THREE.BoxGeometry(30,30,.2),new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));
 wall.position.set(40,3,2);building.add(wall);
 const ring=new THREE.Mesh(new THREE.TorusGeometry(10,.1,8,24),wall.material);ring.position.set(0,2,2);building.add(ring);
 let rays=0;for(const mesh of [wall,ring]){const original=mesh.raycast;mesh.raycast=function(ray,hits){rays++;original.call(this,ray,hits);};}
 const start=new THREE.Vector3(0,1.7,0),windows=[new THREE.Vector3(0,2,5)];
 const first=incendiaryFlight(start,windows,building);assert.ok(first);const initial=rays;assert.ok(initial>0);
 const expected=first.target.clone();first.target.set(999,999,999);
 assert.deepEqual(incendiaryFlight(start,windows,building).target,expected);assert.equal(rays,initial,'repeat recalculated the same flight');
 wall.position.x=0;assert.equal(incendiaryFlight(start,windows,building),null,'moved wall reused clear path');assert.ok(rays>initial);
 const blocked=rays;assert.equal(incendiaryFlight(start,windows,building),null);assert.equal(rays,blocked,'blocked result was not cached');
 wall.geometry.translate(40,0,0);assert.ok(incendiaryFlight(start,windows,building),'changed geometry reused blocked result');assert.ok(rays>blocked);
 const fresh=wall.geometry.attributes.position.array.slice();
 wall.geometry.setAttribute('position',new THREE.BufferAttribute(fresh,3));incendiaryFlight(start,windows,building);
 const replacement=fresh.slice();for(let i=0;i<replacement.length;i+=3)replacement[i]-=40;
 wall.geometry.setAttribute('position',new THREE.BufferAttribute(replacement,3));wall.geometry.computeBoundingBox();wall.geometry.computeBoundingSphere();
 assert.equal(incendiaryFlight(start,windows,building),null,'replacement attribute reused old geometry at the same version');
 wall.geometry.dispose();ring.geometry.dispose();wall.material.dispose();
});

test('short flight segments reject distant geometry but retain nearer blockers',()=>{
 const building=new THREE.Group(),wall=new THREE.Mesh(new THREE.BoxGeometry(30,30,.1),new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));
 building.position.set(12,.3,-6);building.rotation.y=.6;wall.position.set(0,2,20);building.add(wall);building.updateMatrixWorld(true);
 let rays=0;const original=wall.raycast;wall.raycast=function(ray,hits){rays++;original.call(this,ray,hits);};
 const start=building.localToWorld(new THREE.Vector3(0,1.7,0)),windows=[building.localToWorld(new THREE.Vector3(0,2,5))];
 assert.ok(incendiaryFlight(start,windows,building));assert.equal(rays,0,'triangles beyond every segment were still tested');
 wall.position.z=2;assert.equal(incendiaryFlight(start,windows,building),null,'near blocker was rejected with distant geometry');assert.ok(rays>0);
 wall.geometry.dispose();wall.material.dispose();
});

test('staging chooses another corridor when a facade obstacle blocks the first',async()=>{
 const {availableSceneSlot}=await import('../.runtime/frontend-test/city3dEvents.js');
 const lot={x:0,z:16,row:0},building=new THREE.Group();
 const fence=new THREE.Mesh(new THREE.BoxGeometry(.12,2,2),new THREE.MeshBasicMaterial());fence.position.set(4,1,6.35);building.add(fence);
 const slot=availableSceneSlot(lot,'incendiary',[],undefined,s=>incendiaryCorridorClear(s.root,building));
 assert.ok(slot);assert.equal(slot.root.x,0,'selected the corridor through the fence');
 assert.equal(incendiaryCorridorClear({x:3,z:6.35},building),false);
 fence.position.x=20;assert.equal(incendiaryCorridorClear({x:3,z:6.35},building),true);
 fence.geometry.dispose();fence.material.dispose();
});

test('staging rejects a blocked flight and reserves an alternate clear throw',async()=>{
 const {availableSceneSlot}=await import('../.runtime/frontend-test/city3dEvents.js');
 const building=new THREE.Group(),wall=new THREE.Mesh(new THREE.BoxGeometry(2,10,.1),new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));
 wall.position.set(3,5,8);building.add(wall);
 const lot={x:0,z:16,row:0},release=new THREE.Vector3(.4,1.8,.4),windows=[new THREE.Vector3(3,4,12)];
 assert.ok(incendiaryCorridorClear({x:3,z:6.35},building),'fixture blocks the bottle, not the walking corridor');
 assert.equal(incendiaryStagingFlight({x:3,z:6.35},release,windows,building),null);
 const slot=availableSceneSlot(lot,'incendiary',[],undefined,s=>!!incendiaryStagingFlight(s.root,release,windows,building));
 assert.equal(slot.root.x,0);
 assert.deepEqual(release.toArray(),[.4,1.8,.4],'planning mutated the cast release');
 assert.equal(incendiaryStagingFlight(slot.root,release,[],building),null,'missing windows invented a target');
 wall.scale.x=20;
 assert.equal(availableSceneSlot(lot,'incendiary',[],undefined,s=>!!incendiaryStagingFlight(s.root,release,windows,building)),undefined,'all blocked flights invented a slot');
 wall.geometry.dispose();wall.material.dispose();
});
