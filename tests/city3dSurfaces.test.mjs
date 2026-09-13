import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import {cityPlan,route,onRoute,surfaceHeight,vehicleRootHeight,parkingSpot,PITCH} from '../.runtime/frontend-test/city3dPlan.js';
const places=JSON.parse(readFileSync(new URL('../core/locations.json',import.meta.url)));
const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
const plan=cityPlan(places);
test('moving tyres sit on asphalt and parked tyres sit on the raised pavement',()=>{
 for(const model of ['ford','hudson','packard','police']){
  const tyre=manifest[model].bounds_blender[0][2];
  for(const lot of plan.lots){
   const parked=parkingSpot(lot);
   assert.equal(surfaceHeight(parked),.17);
   assert.ok(Math.abs(vehicleRootHeight(parked)+tyre-.175)<1e-8);
  }
  for(const from of plan.lots)for(const to of plan.lots){
   const points=route(from,to,true);
   for(let i=0;i<=100;i++){
    const at=onRoute(points,i/100);
    assert.equal(surfaceHeight(at),-.1,`${from.id}/${to.id}`);
    assert.ok(Math.abs(vehicleRootHeight(at)+tyre+.095)<1e-8);
   }
  }
 }
});
test('street-bed kerbs stay below pedestrian soles and road furniture clears tyre tracks',()=>{
 const bed=manifest['street-bed'];assert.ok(bed);
 assert.ok(bed.bounds_blender[1][2]<=.18+1e-6);
 for(const model of ['person','woman'])assert.ok(manifest[model].motion_bounds_blender[0][2]+.25>.18);
 // All cars share a .87m wheel track offset and .22m tyre width. The longest
 // chassis supplies the most demanding leading/trailing contacts on turns.
 for(const from of plan.lots)for(const to of plan.lots){
  const points=route(from,to,true);
  for(let i=0;i<=150;i++){
   const at=onRoute(points,i/150),sin=Math.sin(at.heading),cos=Math.cos(at.heading);
   for(const side of [-1,1])for(const fore of [-1,1]){
    const x=at.x+side*.87*cos+fore*1.74*sin;
    const z=at.z-side*.87*sin+fore*1.74*cos;
    assert.equal(surfaceHeight({x,z}),-.1,`tyre crosses pavement ${from.id}/${to.id}`);
    for(const lot of plan.lots){
     assert.ok(Math.hypot(x-lot.x,z-lot.row*PITCH)>.60,'manhole under a tyre');
     for(const drain of [-8.5,8.5])assert.ok(Math.abs(x-lot.x-drain)>.52||Math.abs(z-(lot.z-12.4))>.32,'drain under a tyre');
    }
   }
  }
 }
});
test('the exported chapel entry is visible in front of the tower wall',async()=>{
 const THREE=await import('three');
 const {GLTFLoader}=await import('three/addons/loaders/GLTFLoader.js');
 const bytes=readFileSync(new URL('../public/art/models/chapel.glb',import.meta.url));
 const loader=new GLTFLoader();
 // Geometry/raycast verification does not need DOM image decoding in Node.
 loader.register(parser=>({name:'geometry-test-materials',loadMaterial(index){
  return Promise.resolve(new THREE.MeshBasicMaterial({name:parser.json.materials[index].name}));
 }}));
 const model=(await loader.parseAsync(bytes.buffer.slice(bytes.byteOffset,bytes.byteOffset+bytes.byteLength),'')).scene;
 model.rotation.y=Math.PI;model.updateMatrixWorld(true);
 for(const x of [-.48,0,.48]){
  const ray=new THREE.Raycaster(new THREE.Vector3(x,1.7,-10),new THREE.Vector3(0,0,1));
  const hit=ray.intersectObject(model,true)[0];assert.ok(hit);
  assert.match(hit.object.material.name,/chapel oak door/,'door was buried by tower masonry');
  assert.ok(hit.point.z<-5.5);
 }
});
