import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {cityPlan,route} from '../.runtime/frontend-test/city3dPlan.js';
const plan=cityPlan(JSON.parse(readFileSync(new URL('../core/locations.json',import.meta.url))));
async function load(name){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));const loader=new GLTFLoader();loader.register(p=>({name:'harbour-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));return(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;}

test('authored landing reaches the quay, stands in the water and stays clear of all street routes',async()=>{
 const dock=plan.lots.find(l=>l.id==='docks');assert.equal(dock.col,0);
 const pier=await load('harbour-pier');pier.position.set(-24,0,dock.z);pier.updateMatrixWorld(true);
 const bounds=new THREE.Box3().setFromObject(pier,true);
 assert.ok(bounds.min.y<-.9,'piles must extend below the water surface');
 assert.ok(bounds.max.x>=-17&&bounds.max.x<=-15.9,'landing must meet the promenade');
 assert.ok(bounds.min.x>=-32.1&&bounds.max.z<=dock.z+4&&bounds.min.z>=dock.z-4);
 for(const driving of [false,true])for(const a of plan.lots)for(const b of plan.lots)for(const point of route(a,b,driving))
  assert.ok(point.x-(driving?2.9:.7)>bounds.max.x,'harbour structures cannot intrude into street routes');
 const ray=new THREE.Raycaster(new THREE.Vector3(-23.75,4,dock.z),new THREE.Vector3(0,-1,0));
 assert.ok(Math.abs(ray.intersectObject(pier,true)[0].point.y-.2)<.001,'deck boards need a level upper surface');
});

test('quay coping joins the raised promenade while retaining wall extends underwater',async()=>{
 const quay=await load('quay-section');quay.updateMatrixWorld(true);const bounds=new THREE.Box3().setFromObject(quay,true);
 assert.ok(bounds.min.y< -2.8&&bounds.max.y<=.251&&bounds.max.y>=.249);
 assert.ok(bounds.min.z>=-4.001&&bounds.max.z<=4.001,'eight metre sections join without overlap');
 assert.ok(bounds.min.x>=-.376&&bounds.max.x<=.376);
});
