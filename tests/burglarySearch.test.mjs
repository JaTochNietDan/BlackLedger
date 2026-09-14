import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {BurglarySearch,isIndoorSearch} from '../.runtime/frontend-test/burglarySearch.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const l=new GLTFLoader();l.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
test('search playback does not misrepresent occupied, failed, fatal or legacy burglaries',()=>{
 const cue={kind:'robbery',target:'room',burglary:{success:true,taken:0,resident_present:false,fatal:false}};
 assert.ok(isIndoorSearch(cue));
 for(const change of [{success:false},{resident_present:true},{fatal:true}])assert.equal(isIndoorSearch({...cue,burglary:{...cue.burglary,...change}}),false);
 assert.ok(isIndoorSearch({...cue,target:'estate'}));assert.equal(isIndoorSearch({...cue,target:'restaurant'}),false);assert.equal(isIndoorSearch({...cue,burglary:undefined}),false);
});
test('burglar approaches and searches the real open tray without crossing furniture',async()=>{
 for(const [place,model] of [['room','interior-lodging-room'],['mercercourt','interior-flat'],['estate','interior-cypress']])for(const rig of ['person','woman']){
  const room=await load(model),actor=await load(rig),drawer=room.getObjectByName('burglary-drawer');
  const search=new BurglarySearch(actor,drawer,place,180);room.updateMatrixWorld(true);
  let moved=false,opened=false,sawCash=false;
  for(let i=0;i<=150;i++){
   search.update(search.duration*i/150);room.updateMatrixWorld(true);
   const box=new THREE.Box3().setFromObject(actor,true);assert.ok(box.min.y>=.018,'feet cross the floor');
   for(const y of [.3,.9,1.3])for(const dir of [[1,0,0],[-1,0,0],[0,0,1],[0,0,-1]])assert.equal(new THREE.Raycaster(new THREE.Vector3(actor.position.x,y,actor.position.z),new THREE.Vector3(...dir),0,.25).intersectObject(room,true).length,0,`${place}/${rig} intersects furniture at ${i}`);
   moved ||= Math.abs(actor.getObjectByName('leg1').rotation.x)>.2;opened ||= drawer.position.z>.33;sawCash ||= search.money.visible;
  }
  assert.ok(moved&&opened&&sawCash);assert.ok(drawer.position.z<.001);assert.equal(actor.visible,false);
  search.update(search.arrival+1.5);room.updateMatrixWorld(true);
  let hand;actor.getObjectByName('elbow1').traverse(o=>{if(o.isMesh&&o.name.startsWith('hand'))hand=o;});
  const handBox=new THREE.Box3().setFromObject(hand,true),tray=new THREE.Box3().setFromObject(drawer,true),p=handBox.getCenter(new THREE.Vector3());
  assert.ok(p.x>tray.min.x&&p.x<tray.max.x&&p.z<tray.max.z-.04&&p.z>tray.min.z,`${place}/${rig} hand stops short of the open tray: ${JSON.stringify({hand:p,tray})}`);
  assert.ok(p.y>tray.min.y&&p.y<tray.max.y+.12,`${place}/${rig} hand misses tray height: ${JSON.stringify({hand:p,tray})}`);
  const empty=new BurglarySearch(actor,drawer,place,0);empty.update(empty.arrival+3);assert.equal(empty.money.visible,false);
 }
});

test('search turns continuously at corners, drawer arrival, departure and exit',async()=>{
 for(const place of ['room','mercercourt','estate']){
  const actor=await load('person'),search=new BurglarySearch(actor,new THREE.Group(),place,0);
  let previous=actor.rotation.y;
  for(let frame=1;frame<=Math.ceil(search.duration*60);frame++){
   search.update(frame/60);
   const difference=Math.abs(Math.atan2(Math.sin(actor.rotation.y-previous),Math.cos(actor.rotation.y-previous)));
   assert.ok(difference<.12,`${place}: abrupt heading change at ${frame/60}: ${difference}`);
   previous=actor.rotation.y;
  }
  // Seeking backwards must reconstruct the pose rather than depend on frame history.
  search.update(search.arrival+2);const expected=actor.quaternion.clone();
  search.update(0);search.update(search.arrival+2);
  assert.ok(expected.angleTo(actor.quaternion)<1e-6);
 }
});
