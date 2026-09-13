import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {CityFire,clearBlastWindows} from '../.runtime/frontend-test/city3dFire.js';
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

test('window debris clears exported facades and canopies before settling',async()=>{
 const {windowDebris}=await import('../.runtime/frontend-test/city3dBlast.js');
 for(const name of ['tavern','monarch','tenement','shop','civic','casino','warehouse','bluehour','goldenlily','papermoon']){
  const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
  loader.register(p=>({name:'debris-geometry',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));}}));
  const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;model.updateMatrixWorld(true);
  const windows=clearBlastWindows(model);
  for(let i=0;i<12;i++){
   const window=windows[i%windows.length];assert.ok(window,name);
   let prev;
   for(let t=0;t<1.8;t+=.01){
    const p=windowDebris(i,t,window,0,windows.length),point=new THREE.Vector3(p.x,p.height+.05,p.z);
    if(prev){const step=point.clone().sub(prev);if(step.length()>1e-6){
     const ray=new THREE.Raycaster(prev,step.clone().normalize(),0,step.length());
     const hits=ray.intersectObject(model,true);assert.equal(hits.length,0,`${name} fragment ${i} collides at ${t}: ${hits[0]?.object.name} ${JSON.stringify(hits[0]?.point)}`);
    }}prev=point;
    assert.ok(p.z<window.z&&p.height>=0);
   }
   const settled=windowDebris(i,2,window,0,windows.length);assert.ok(settled.height<1e-10);assert.equal(settled.rx,0);assert.equal(settled.rz,0);
  }
 }
});

test('authored broken glazing has a missing centre and retains separate intact panes',async()=>{
 for(const name of ['tavern','monarch','tenement','shop','civic','casino','warehouse','bluehour','goldenlily','papermoon']){
  const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
  loader.register(p=>({name:'glazing-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({side:THREE.DoubleSide,name:p.json.materials[i].name}));}}));
  const model=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;model.updateMatrixWorld(true);
  const intact=model.getObjectByName('window-intact'),broken=model.getObjectByName('window-broken');assert.ok(intact&&broken,name);
  const window=model.getObjectByName('fire-window-1-1')||model.getObjectByName('fire-window-0-1');
  const p=window.getWorldPosition(new THREE.Vector3());p.x+=.12;
  const ray=new THREE.Raycaster(p,new THREE.Vector3(0,0,1),0,.2);
  assert.ok(ray.intersectObject(intact,true).length,`${name} intact pane missing`);
  const hits=ray.intersectObject(broken,true);assert.ok(hits.length,`${name} recess missing`);
  assert.equal(hits[0].object.material.name,'unlit window recess');
 }
});
