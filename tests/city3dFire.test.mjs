import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {CityFire} from '../.runtime/frontend-test/city3dFire.js';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';

test('saved fire persists on paused game time, freezes motion and extinguishes exactly',()=>{
 const texture=new THREE.Texture(),fire=new CityFire(texture),camera=new THREE.PerspectiveCamera();
 const building=new THREE.Group(),vent=new THREE.Object3D();vent.name='fire-window-0-0';vent.position.set(0,1.8,-6.22);building.add(vent);
 const buildings=new Map([['bar',building]]),records=[{id:'fire',target:'bar',minute:480,brigade_at:490,extinguished_at:525,cleanup_at:570}];
 const update=(minute,dt=100,motion=true)=>fire.update(records,minute,buildings,camera,dt,motion);
 update(480);const mesh=fire.root.children[0];assert.equal(mesh.count,12);
 for(let i=0;i<100;i++)update(480);
 assert.equal(fire.root.children[0],mesh,'paused decisions cannot expire a saved fire');
 const before=Array.from(mesh.instanceMatrix.array);update(480,100,false);assert.deepEqual(Array.from(mesh.instanceMatrix.array),before);
 update(524);assert.equal(fire.inspect().scenes.length,1);
 update(525);assert.equal(fire.inspect().scenes.length,0);assert.equal(fire.root.children.length,0);
 fire.dispose();texture.dispose();
});
test('exported fire anchors are on the authored front windows',async()=>{
 for(const name of ['tavern','monarch','tenement','shop','civic','casino','warehouse','bluehour','goldenlily','papermoon']){
  const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
  loader.register(p=>({name:'fire-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
  const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;model.updateMatrixWorld(true);
  const anchors=[];model.traverse(o=>{if(o.name.startsWith('fire-window-'))anchors.push(o.getWorldPosition(new THREE.Vector3()));});
  assert.ok(anchors.length>=4,name);
  for(const v of anchors){
   assert.ok(v.y>1&&v.z<0);
   const ray=new THREE.Raycaster(v,new THREE.Vector3(0,0,1),0,.1);
   assert.ok(ray.intersectObject(model,true).some(h=>/glass|windows/.test(h.object.material.name)),`${name} anchor lacks window immediately behind it`);
  }
 }
});
