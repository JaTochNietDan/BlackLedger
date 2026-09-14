import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityIncendiary,incendiaryFlight,INCENDIARY_RELEASE,INCENDIARY_IMPACT} from '../.runtime/frontend-test/city3dIncendiary.js';
import {clearBlastWindows} from '../.runtime/frontend-test/city3dFire.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'collision-geometry',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 const root=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;root.updateMatrixWorld(true);return root;
}
test('incendiary bottle centre reaches facade without passing through earlier masonry',async()=>{
 const actor=await model('person'),bottle=await model('incendiary-bottle');
 for(const name of ['tavern','monarch','tenement','shop','civic','casino','warehouse','bluehour','goldenlily','papermoon','mariner','mercer-court']){
  const building=await model(name);building.position.set(0,.18,16);building.updateMatrixWorld(true);
  for(const x of [3,0,-3]){
   const origin=new THREE.Vector3(x,.2,6.35),windows=clearBlastWindows(building).sort((a,b)=>Math.abs(a.x-x)-Math.abs(b.x-x));
   assert.ok(windows.length,name);const target=windows[0].clone().sub(origin);
   const cast=new CityIncendiary(actor.clone(true),bottle.clone(true),target);cast.root.position.copy(origin);
   const flight=incendiaryFlight(cast.root.localToWorld(cast.release.clone()),windows,building);
   assert.ok(flight,`${name} x=${x} has no clear throw`);
   cast.target.copy(flight.target).sub(origin);cast.loft=flight.loft;
   let previous;
   for(let frame=0;frame<80;frame++){
    cast.update(INCENDIARY_RELEASE+(INCENDIARY_IMPACT-INCENDIARY_RELEASE)*frame/80);
    const point=cast.bottle.getWorldPosition(new THREE.Vector3());
    if(previous){const direction=point.clone().sub(previous);const hits=new THREE.Raycaster(previous,direction.clone().normalize(),0,direction.length()).intersectObject(building,true);
     assert.equal(hits.length,0,`${name} x=${x} hits ${hits[0]?.object.name} before impact at frame ${frame}`);
    }previous=point;
   }
  }
 }
});
