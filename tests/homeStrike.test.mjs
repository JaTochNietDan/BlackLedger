import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {isHomeStrike,homeStrikeFor,homeStrikeRoom} from '../.runtime/frontend-test/homeStrike.js';
import {CityAssassination} from '../.runtime/frontend-test/city3dAssassination.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const l=new GLTFLoader();l.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
const hit={id:'hit',kind:'gunfight',target:'mercercourt',minute:10,attacker:{id:'player',weapon:1},strike:{setting:'home',variant:'back-of-head',victim:{id:'mara'}}};
test('private room playback requires an explicit supported home setting and linked victim',()=>{
 assert.equal(isHomeStrike(hit),true);
 assert.equal(isHomeStrike({...hit,target:'estate'}),true);
 assert.equal(homeStrikeRoom('estate').model,'interior-cypress');
 assert.equal(isHomeStrike({...hit,strike:{...hit.strike,setting:undefined}}),false);
 assert.equal(isHomeStrike({...hit,target:'room'}),false);
 const death={id:'death',kind:'killing',target:'mercercourt',minute:10,actors:[{id:'mara'}]};
 assert.equal(homeStrikeFor(death,[hit]),hit);
 assert.equal(homeStrikeFor({...death,actors:[{id:'leo'}]},[hit]),undefined);
 assert.equal(homeStrikeFor({...death,minute:11},[hit]),undefined);
});
test('full-sized home strike cast stays inside furnished homes for approach, impact and fall',async()=>{
 for(const target of ['mercercourt','estate']){
 const settings=homeStrikeRoom(target),room=await load(settings.model);room.updateMatrixWorld(true);
 for(const [model,variant,gun] of [['person','back-of-head','revolver'],['woman','close-shot','shotgun'],['person','burst','thompson'],['woman','close-quarters',undefined]]){
  const a=await load(model),v=await load(model),weapon=gun?await load(gun):undefined;
  const cast=new CityAssassination(a,v,weapon,variant,gun);cast.root.position.set(settings.origin.x,0,settings.origin.z);
  for(let i=0;i<=80;i++){
   cast.update(cast.duration*i/80);
   for(const who of [a,v]){
    const box=new THREE.Box3().setFromObject(who,true);
    assert.ok(box.min.x> -settings.bounds.x&&box.max.x<settings.bounds.x&&box.min.z> -settings.bounds.z&&box.max.z<settings.bounds.z,`${target}/${variant} crosses room boundary at ${i}`);
    const center=who.getWorldPosition(new THREE.Vector3());
    for(const y of [.4,1,1.6])for(const dir of [[1,0,0],[-1,0,0],[0,0,1],[0,0,-1]]){
     assert.equal(new THREE.Raycaster(new THREE.Vector3(center.x,y,center.z),new THREE.Vector3(...dir),0,.28).intersectObject(room,true).length,0,`${variant} intersects furniture`);
    }
   }
  }
 }
 }
});
