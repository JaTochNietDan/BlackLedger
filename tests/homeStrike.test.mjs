import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {isFlatHomeStrike,flatHomeStrikeFor,HOME_STRIKE_ORIGIN} from '../.runtime/frontend-test/homeStrike.js';
import {CityAssassination} from '../.runtime/frontend-test/city3dAssassination.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const l=new GLTFLoader();l.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
const hit={id:'hit',kind:'gunfight',target:'mercercourt',minute:10,attacker:{id:'player',weapon:1},strike:{setting:'home',variant:'back-of-head',victim:{id:'mara'}}};
test('private room playback requires an explicit supported home setting and linked victim',()=>{
 assert.equal(isFlatHomeStrike(hit),true);
 assert.equal(isFlatHomeStrike({...hit,strike:{...hit.strike,setting:undefined}}),false);
 assert.equal(isFlatHomeStrike({...hit,target:'estate'}),false);
 const death={id:'death',kind:'killing',target:'mercercourt',minute:10,actors:[{id:'mara'}]};
 assert.equal(flatHomeStrikeFor(death,[hit]),hit);
 assert.equal(flatHomeStrikeFor({...death,actors:[{id:'leo'}]},[hit]),undefined);
 assert.equal(flatHomeStrikeFor({...death,minute:11},[hit]),undefined);
});
test('full-sized home strike cast stays inside the furnished flat for approach, impact and fall',async()=>{
 const room=await load('interior-flat');room.updateMatrixWorld(true);
 for(const [model,variant,gun] of [['person','back-of-head','revolver'],['woman','close-shot','shotgun'],['person','burst','thompson'],['woman','close-quarters',undefined]]){
  const a=await load(model),v=await load(model),weapon=gun?await load(gun):undefined;
  const cast=new CityAssassination(a,v,weapon,variant,gun);cast.root.position.set(HOME_STRIKE_ORIGIN.x,0,HOME_STRIKE_ORIGIN.z);
  for(let i=0;i<=80;i++){
   cast.update(cast.duration*i/80);
   for(const who of [a,v]){
    const box=new THREE.Box3().setFromObject(who,true);
    assert.ok(box.min.x> -4&&box.max.x<4&&box.min.z> -3.5&&box.max.z<3.5,`${variant} crosses room boundary at ${i}`);
    const center=who.getWorldPosition(new THREE.Vector3());
    for(const y of [.4,1,1.6])for(const dir of [[1,0,0],[-1,0,0],[0,0,1],[0,0,-1]]){
     assert.equal(new THREE.Raycaster(new THREE.Vector3(center.x,y,center.z),new THREE.Vector3(...dir),0,.28).intersectObject(room,true).length,0,`${variant} intersects furniture`);
    }
   }
  }
 }
});
