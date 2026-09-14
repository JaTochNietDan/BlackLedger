import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const l=new GLTFLoader();l.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));return(await l.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}
for(const [name,x,y,height] of [['interior-flat',-1.25,2.65,.78],['interior-lodging-room',-.67,2.48,.77],['interior-cypress',2.5,1.7,.74]])test(`${name} drawer has a real tray that opens out of its case`,async()=>{
 const room=await load(name),drawer=room.getObjectByName('burglary-drawer');assert.ok(drawer);room.updateMatrixWorld(true);
 const closed=new THREE.Box3().setFromObject(drawer,true);drawer.position.z=.34;room.updateMatrixWorld(true);
 const open=new THREE.Box3().setFromObject(drawer,true);assert.ok(Math.abs(open.min.z-closed.min.z-.34)<.001);assert.ok(Math.abs(open.min.x-closed.min.x)<.001);
 const point=new THREE.Vector3(x,1,-y+.53),hit=new THREE.Raycaster(point,new THREE.Vector3(0,-1,0),0,1).intersectObject(room,true)[0];
 assert.ok(hit&&Math.abs(hit.point.y-(height*.49+.0125))<.002,'open tray floor is not accessible from above');
 drawer.position.z=0;room.updateMatrixWorld(true);const restored=new THREE.Box3().setFromObject(drawer,true);assert.ok(restored.min.distanceTo(closed.min)<.001);
});
