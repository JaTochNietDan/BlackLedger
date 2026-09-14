import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
import {placementsForInterior,interiorPlayerSpot,poseInteriorOccupant} from '../.runtime/frontend-test/interiorStaging.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshStandardMaterial());}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('every interior entrance keeps both real rigs on the floor, clear of occupants and furniture, and settles without replay',async()=>{
 for(const name of ['person','woman'])for(const place of ['bar','room','mercercourt','laundry']){
  const source=await load(name),actor=source.clone(true),arrival=new InteriorArrival(actor,place),end=interiorPlayerSpot(place);
  const people=[{id:'worker',role:place==='bar'?'Barman':place==='room'?'Landlady':place==='laundry'?'Laundress':'Superintendent'},...Array.from({length:20},(_,i)=>({id:`guest-${i}`}))];
  const occupied=[...placementsForInterior(place,people)].map(([id,s])=>{const npc=source.clone(true);poseInteriorOccupant(npc,s);return {id,box:new THREE.Box3().setFromObject(npc,true)};});
  const obstacle=place==='laundry'?new THREE.Box3(new THREE.Vector3(1.97,0,-1.67),new THREE.Vector3(4.73,1.5,-.53)):place==='room'?new THREE.Box3(new THREE.Vector3(2.37,0,-5.7),new THREE.Vector3(4.85,4.2,.05)):place==='mercercourt'?new THREE.Box3(new THREE.Vector3(2.85,0,-6.31),new THREE.Vector3(5.65,4.2,.02)):new THREE.Box3(new THREE.Vector3(-5,0,2.0),new THREE.Vector3(5.3,1.3,3.3));
  let previous;let swung=false;
  for(let frame=0;frame<=180;frame++){
   const seconds=arrival.duration*frame/180,active=arrival.pose(seconds),box=new THREE.Box3().setFromObject(actor,true);
   assert.equal(active,frame<180);
   assert.ok(box.min.y>=.0175,`${place}/${name} penetrates floor`);
   assert.ok(box.min.x> -5.7&&box.max.x<5.8&&box.min.z> -4.8&&box.max.z<(place==='room'?4:5),`${place}/${name} leaves floor`);
   assert.equal(box.intersectsBox(obstacle),false,`${place}/${name} intersects furniture`);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${place}/${name} crosses ${other.id}`);
   const mat=place==='room'?{x:1.1,z0:2.225,z1:3.575,top:.0555}:place==='mercercourt'?{x:1.4,z0:2.65,z1:4.15,top:.0575}:undefined;
   actor.traverse(o=>{if(o.isMesh&&o.name.startsWith('shoe')){const b=new THREE.Box3().setFromObject(o,true);if(mat&&b.min.x<mat.x&&b.max.x> -mat.x&&b.min.z<mat.z1&&b.max.z>mat.z0)assert.ok(b.min.y>=mat.top,`${place}/${name} shoe penetrates mat`);}});
   if(previous)assert.ok(actor.position.distanceTo(previous)<.07,'entrance jumps between frames');
   previous=actor.position.clone();swung ||= Math.abs(actor.getObjectByName('leg1').rotation.x)>.1;
  }
  assert.ok(swung,'entrance slides without moving legs');assert.ok(Math.abs(actor.position.x-end.x)<1e-9&&Math.abs(actor.position.z-end.z)<1e-9);
  assert.ok(Math.abs(Math.sin(actor.rotation.y-end.yaw))<1e-9);
  for(const side of [-1,1])for(const part of ['leg','knee','arm','elbow'])assert.ok(Math.abs(actor.getObjectByName(part+side).rotation.x)<1e-9,'final pose is not at rest');
  const final=actor.position.clone();assert.equal(arrival.pose(arrival.duration+100),false);assert.ok(actor.position.distanceTo(final)<1e-9);
 }
});
