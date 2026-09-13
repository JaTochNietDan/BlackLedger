import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {trafficSize,trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';
import {CityAftermath} from '../.runtime/frontend-test/city3dAftermath.js';

test('exported fire engine fits its reserved footprint',async()=>{
 const b=readFileSync(new URL('../public/art/models/fire-engine.glb',import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'engine-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 const box=new THREE.Box3().setFromObject(model,true),size=trafficSize('parked-fire-engine');
 assert.ok(box.max.x-box.min.x<=size.width);assert.ok(box.max.z-box.min.z<=size.length);
 assert.ok(box.min.y>=0&&box.min.y<.03);assert.ok(box.max.y>2.3&&box.max.y<2.8);
});
test('brigade parks separately from death-scene police and stays after extinguishing',()=>{
 const aftermath=new CityAftermath(),lots=new Map([['bar',{id:'bar',x:80,z:48,row:1,col:2}]]);
 const models=new Map(['person','police','police-officer','fire-engine','firefighter'].map(n=>[n,new THREE.Group()]));
 const deaths=[{id:'death',target:'bar',victim:{id:'mara',name:'Mara'},minute:480,police_at:485,cleanup_at:660}];
 const fires=[{id:'fire',target:'bar',minute:480,brigade_at:490,extinguished_at:660,cleanup_at:720}];
 const update=minute=>aftermath.update(deaths,minute,lots,models,()=> 'person',[{model:'parked-hudson',root:{x:89.6,z:48},pose:{x:89.6,z:48,heading:0}}],new Set(),[],new Set(),fires);
 update(489);assert.equal(aftermath.inspect().filter(e=>e.id.endsWith('fire-engine')).length,0);
 update(490);for(const slot of aftermath.slots())assert.equal(trafficOverlap(slot.pose,slot.model,{x:89.6,z:48,heading:0},'parked-hudson'),false);assert.equal(aftermath.inspect().filter(e=>e.id.includes('firefighter')).length,2);assert.equal(aftermath.inspect().filter(e=>e.id.endsWith('fire-engine')).length,1);
 for(const [i,a] of aftermath.slots().entries())for(const b of aftermath.slots().slice(i+1))assert.equal(trafficOverlap(a.pose,a.model,b.pose,b.model),false);
 update(660);assert.equal(aftermath.inspect().filter(e=>e.id.endsWith('fire-engine')).length,1);
 update(720);assert.equal(aftermath.reservations().length,0);
 aftermath.dispose();
});

test('firefighter export retains an articulated rig inside standing clearance',async()=>{
 const b=readFileSync(new URL('../public/art/models/firefighter.glb',import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'crew-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 for(const name of ['leg1','leg-1','arm1','arm-1'])assert.ok(model.getObjectByName(name));
 const box=new THREE.Box3().setFromObject(model,true);
 assert.ok(box.min.x>=-.5&&box.max.x<=.5&&box.min.z>=-.5&&box.max.z<=.5);
 assert.ok(box.min.y>=-.01&&box.max.y<=2.05);
 assert.equal(model.getObjectByName('headwear-fedora'),undefined);
});
