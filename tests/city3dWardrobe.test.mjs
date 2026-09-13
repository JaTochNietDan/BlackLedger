import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {wardrobe,dressPedestrian} from '../.runtime/frontend-test/city3dWardrobe.js';

test('cast palettes follow explicit portraits and remain stable for generated identities',()=>{
 for(let face=1;face<=24;face++)assert.deepEqual(wardrobe('player-a',face,true),wardrobe('player-b',face,true));
 assert.ok(new Set(Array.from({length:24},(_,i)=>wardrobe('cast',i+1).coat)).size>=20);
 for(const id of ['mara','Élodie','constructor','npc-a']){
  assert.deepEqual(wardrobe(id),wardrobe(id,NaN));
  assert.deepEqual(wardrobe(id),wardrobe(id,2.5));
  assert.match(wardrobe(id).skin,/^#[a-f0-9]{6}$/);
 }
 assert.notDeepEqual(wardrobe('Alex Varga',5,true),wardrobe('Alex Varga',0,true));
});

test('cast material copies preserve weave maps, geometry and source materials across actors',()=>{
 const coat=new THREE.MeshStandardMaterial({name:'wool suit',map:new THREE.Texture(),normalMap:new THREE.Texture()});
 const geometry=new THREE.BoxGeometry(),source=new THREE.Group();
 for(let i=0;i<5;i++)source.add(new THREE.Mesh(geometry,coat));
 const a=source.clone(true),b=source.clone(true);
 const ownedA=dressPedestrian(a,'person',wardrobe('cast',1));
 const ownedB=dressPedestrian(b,'person',wardrobe('cast',8));
 assert.equal(ownedA.length,1,'all suit pieces share one private material');
 assert.equal(new Set(a.children.map(m=>m.material)).size,1);
 assert.notEqual(ownedA[0],ownedB[0]);assert.ok(!ownedA[0].color.equals(ownedB[0].color));
 assert.equal(ownedA[0].map,coat.map);assert.equal(ownedA[0].normalMap,coat.normalMap);
 assert.equal(a.children[0].geometry,geometry);assert.equal(source.children[0].material,coat);
 assert.ok(coat.color.equals(new THREE.Color(0xffffff)));
 let textureDisposals=0;coat.map.addEventListener('dispose',()=>textureDisposals++);
 ownedA.forEach(m=>m.dispose());assert.equal(textureDisposals,0);
});
