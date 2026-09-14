import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {poolhallPlacements,poseInteriorOccupant,interiorPlayerSpot} from '../.runtime/frontend-test/interiorStaging.js';
import {InteriorArrival} from '../.runtime/frontend-test/interiorArrival.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('Six-table billiard hall keeps staff, spectators and entrance clear for both rigs',async()=>{
 const room=await load('interior-poolhall');room.updateMatrixWorld(true);
 assert.ok(room.getObjectByName('interior-wall-left'));assert.ok(room.getObjectByName('interior-wall-back'));
 for(const x of [-3.7,3.7])for(const z of [-4.25,0,4.25]){
  const top=new THREE.Raycaster(new THREE.Vector3(x,.95,z),new THREE.Vector3(0,-1,0),0,.5).intersectObject(room,true)[0];
  assert.ok(top&&Math.abs(top.point.y-.78)<.003,`table ${x},${z}: missing cloth`);
  for(const [dx,dz] of [[-.661,-1.296],[-.661,1.296],[.661,-1.296],[.661,1.296],[-.680,0],[.680,0]]){
   const pocket=new THREE.Raycaster(new THREE.Vector3(x+dx,.95,z+dz),new THREE.Vector3(0,-1,0),0,.5).intersectObject(room,true)[0];
   assert.ok(pocket&&pocket.point.y<.66,`table ${x},${z}: pocket is painted over`);
  }
  for(const sign of [-1,1]){
   const end=new THREE.Raycaster(new THREE.Vector3(x,.78+.028575,z),new THREE.Vector3(0,0,sign),0,1.4).intersectObject(room,true)[0];
   assert.ok(end&&Math.abs(end.distance-1.27)<.002,`table ${x},${z}: end nose differs from physics`);
   const side=new THREE.Raycaster(new THREE.Vector3(x,.78+.028575,z+.9),new THREE.Vector3(sign,0,0),0,.8).intersectObject(room,true)[0];
   assert.ok(side&&Math.abs(side.distance-.635)<.002,`table ${x},${z}: side nose differs from physics`);
  }
  const ball=new THREE.Raycaster(new THREE.Vector3(x-.48,.95,z+.68),new THREE.Vector3(0,-1,0),0,.3).intersectObject(room,true)[0];
  assert.ok(ball&&Math.abs(ball.point.y-(.78+2*.028575))<.001,`table ${x},${z}: ball scale differs from physics`);
 }

 const people=[{id:'a-worker',role:'Marker'},{id:'b-worker',role:'Table hand'},...Array.from({length:16},(_,i)=>({id:`visitor-${i}`}))];
 const places=poolhallPlacements(people);assert.equal(places.size,11);assert.equal(places.get('a-worker').id,'pool-marker');
 assert.equal(places.get('b-worker').id,'pool-table-hand');
 assert.deepEqual([...poolhallPlacements([...people].reverse())],[...places]);
 const spots=[...places.values(),interiorPlayerSpot('poolhall')];
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
  const actor=source.clone(true),arrival=new InteriorArrival(actor,'poolhall');
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

test('Tournament players stand beside their physical table and leave resolved matches',async()=>{
 const {poolhallTableOrigin,tournamentHallSpots}=await import('../.runtime/frontend-test/poolhallMatches.js');
 const games=Array.from({length:6},(_,i)=>({index:i,table_number:i+1,players:[i===0?'p':`a${i}`,`b${i}`],resolved:false,table:{}}));
 const t={player_id:'p',settled:false,games};
 const spots=tournamentHallSpots(t);assert.equal(spots.size,12);assert.ok(spots.has('player'));
 const room=await load('interior-poolhall');room.updateMatrixWorld(true);
 for(let i=1;i<=6;i++){const [x,z]=poolhallTableOrigin(i);const hit=new THREE.Raycaster(new THREE.Vector3(x,.9,z),new THREE.Vector3(0,-1,0),0,.3).intersectObject(room,true)[0];assert.ok(hit&&Math.abs(hit.point.y-.78)<.003);}
 for(const name of ['person','woman']){const source=await load(name),boxes=[];for(const spot of spots.values()){const actor=source.clone(true);poseInteriorOccupant(actor,spot);const box=new THREE.Box3().setFromObject(actor,true);for(const prior of boxes)assert.equal(box.intersectsBox(prior),false);boxes.push(box);assert.ok(box.min.y>-.02);}}
 games[0].resolved=true;assert.equal(tournamentHallSpots(t).has('player'),false);
 assert.equal(tournamentHallSpots({...t,settled:true}).size,0);
 assert.equal(poolhallTableOrigin(0),null);assert.equal(poolhallTableOrigin(7),null);
});

test('Hall table picking follows reused tables and excludes the aisles',async()=>{
 const {poolhallGameAt,poolhallTableOrigin}=await import('../.runtime/frontend-test/poolhallMatches.js');
 const games=[{index:0,table_number:1,table:{width:1.27,length:2.54}},{index:2,table_number:1,table:{width:1.27,length:2.54}}];
 const t={games};const [x,z]=poolhallTableOrigin(1);
 assert.equal(poolhallGameAt(t,x,z),2);
 assert.equal(poolhallGameAt(t,0,z),null);
 assert.equal(poolhallGameAt(t,x,z+2),null);
 assert.equal(poolhallGameAt(t,NaN,z),null);
 assert.equal(poolhallGameAt(undefined,x,z),null);
});

test('Both hall rigs grip the cue without changing arm length or touching the floor',async()=>{
 const {createBilliardsCue,holdBilliardsCue}=await import('../.runtime/frontend-test/billiardsCue.js');
 for(const name of ['person','woman']){
  const actor=await load(name);poseInteriorOccupant(actor,{id:'match',x:0,z:0,yaw:Math.PI/2});
  const {cue}=createBilliardsCue();holdBilliardsCue(actor,cue);
  const palm=actor.getObjectByName('elbow1').localToWorld(new THREE.Vector3(0,-.295,0));
  const shaft=cue.localToWorld(new THREE.Vector3(0,.23,0));
  assert.ok(palm.distanceTo(shaft)<.002,`${name}: cue misses palm`);
  assert.ok(new THREE.Box3().setFromObject(cue,true).min.y>.07,`${name}: cue penetrates floor`);
  assert.equal(actor.getObjectByName('arm1').scale.y,1);assert.equal(actor.getObjectByName('elbow1').scale.y,1);
 }
});
