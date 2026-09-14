import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {heraldPlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Herald staff, visitors and entrance clear desks and filing cabinets',async()=>{
 const room=await load('interior-herald');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 const people=[{id:'a-worker',role:'Editor'},{id:'b-worker',role:'Reporter'},...Array.from({length:16},(_,i)=>({id:`visitor-${i}`}))];
 const places=heraldPlacements(people);assert.equal(places.size,7);assert.equal(places.get('a-worker').id,'herald-desk-0');
 assert.equal(places.get('b-worker').id,'herald-desk-1');
 assert.deepEqual([...heraldPlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('herald')];
 for(const name of ['person','woman']){
  const source=await load(name),occupied=[];
  for(const spot of spots){
   const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);
   for(const other of occupied)assert.equal(box.intersectsBox(other.box),false,`${name}: ${spot.id} overlaps ${other.id}`);
   occupied.push({id:spot.id,box});
   if(spot.seat!==undefined){
    const support=new THREE.Raycaster(new THREE.Vector3(spot.x,.75,spot.z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
    assert.ok(support&&Math.abs(support.point.y-spot.seat)<.01,`${spot.id}: missing cushion`);
    continue;
   }
   for(const y of [.2,.65,1.1,1.65])for(let i=0;i<24;i++){
    const dir=new THREE.Vector3(Math.cos(i*Math.PI/12),0,Math.sin(i*Math.PI/12));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(spot.x,y,spot.z),dir,0,.48).intersectObject(room,true).length,0,`${spot.id}: furniture at ${y}`);
   }
  }
  const actor=source.clone(true),arrival=new InteriorArrival(actor,'herald');
  for(let i=0;i<=40;i++){
   arrival.pose(arrival.duration*i/40);
   const {x,z}=actor.position;
   const support=new THREE.Raycaster(new THREE.Vector3(x,.1,z),new THREE.Vector3(0,-1,0),0,.2).intersectObject(room,true)[0];
   assert.ok(support&&support.point.y>=.009&&support.point.y<=.019,`entry ${i}: unsupported floor`);
   for(const y of [.3,.9,1.5])for(let r=0;r<16;r++){
    const dir=new THREE.Vector3(Math.cos(r*Math.PI/8),0,Math.sin(r*Math.PI/8));
    assert.equal(new THREE.Raycaster(new THREE.Vector3(x,y,z),dir,0,.48).intersectObject(room,true).length,0,`entry ${i}: blocked`);
   }
  }
 }
});

test('Herald typing keeps both hands above real keys and stops for reduced motion',async()=>{
 const {NewsroomTyping}=await import('../.runtime/frontend-test/newsroomTyping.js');
 for(const name of ['person','woman']){
  const actor=await load(name),spot=heraldPlacements([{id:'editor',role:'Editor'}]).get('editor');
  poseInteriorOccupant(actor,spot);const typing=new NewsroomTyping(actor);
  const hands=[];for(const side of [-1,1])actor.getObjectByName(`elbow${side}`).traverse(o=>{if(o instanceof THREE.Mesh&&o.name.startsWith('hand'))hands.push(o);});
  for(let i=0;i<150;i++){
   typing.step(.05,true,false);
   for(const hand of hands){const box=new THREE.Box3().setFromObject(hand,true);assert.ok(box.min.y>=.969&&box.min.y<1.05,`${name}: hand height ${box.min.y}`);assert.ok(box.min.z> -2.4&&box.max.z< -2.0,`${name}: hand leaves keyboard`);assert.ok(box.min.x>spot.x-.6&&box.max.x<spot.x+.18);}
   for(const side of [-1,1]){assert.equal(actor.getObjectByName(`arm${side}`).scale.y,1);assert.equal(actor.getObjectByName(`elbow${side}`).scale.y,1);}
  }
  const before=typing.seconds;assert.equal(typing.step(1,false,false),false);assert.equal(typing.step(1,true,true),false);assert.equal(typing.seconds,before);
 }
});
