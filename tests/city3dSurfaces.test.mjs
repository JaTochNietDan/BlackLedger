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
test('the funeral business has a visible entry and a clear rear coach gateway',async()=>{
 const THREE=await import('three');
 const {GLTFLoader}=await import('three/addons/loaders/GLTFLoader.js');
 const bytes=readFileSync(new URL('../public/art/models/undertaker.glb',import.meta.url));
 const loader=new GLTFLoader();
 // Geometry/raycast verification does not need DOM image decoding in Node.
 loader.register(parser=>({name:'geometry-test-materials',loadMaterial(index){
  return Promise.resolve(new THREE.MeshBasicMaterial({name:parser.json.materials[index].name}));
 }}));
 const model=(await loader.parseAsync(bytes.buffer.slice(bytes.byteOffset,bytes.byteOffset+bytes.byteLength),'')).scene;
 model.updateMatrixWorld(true);
 const box=new THREE.Box3().setFromObject(model,true),[min,max]=manifest.undertaker.bounds_blender;
 for(const [actual,expected]of [[box.min.x,min[0]],[box.max.x,max[0]],[box.min.y,min[2]],[box.max.y,max[2]],[box.min.z,-max[1]],[box.max.z,-min[1]]])
  assert.ok(Math.abs(actual-expected)<.001,'manifest must measure actual vertices');
 for(const x of [4.08,4.5,4.92]){
  const ray=new THREE.Raycaster(new THREE.Vector3(x,1.7,-10),new THREE.Vector3(0,0,1));
  const hit=ray.intersectObject(model,true)[0];assert.ok(hit);
  assert.match(hit.object.material.name,/funeral oak/,'entry must face the public street');
  assert.ok(hit.point.z<-7.4);
 }
 for(let x=-1.1;x<=1.1001;x+=.05){
  const ray=new THREE.Raycaster(new THREE.Vector3(x,.8,8),new THREE.Vector3(0,0,-1),0,3);
  assert.equal(ray.intersectObject(model,true).length,0,'rear gates must leave room for the hearse');
 }
 assert.equal(plan.lots.find(l=>l.id==='chapel').model,'undertaker');
});
test('the coach yard uses embedded colour and normal textures',()=>{
 const bytes=readFileSync(new URL('../public/art/models/undertaker.glb',import.meta.url));
 const gltf=JSON.parse(bytes.subarray(20,20+bytes.readUInt32LE(12)).toString());
 const material=gltf.materials.find(m=>m.name==='coach yard gravel');assert.ok(material);
 for(const texture of [material.pbrMetallicRoughness.baseColorTexture,material.normalTexture]){
  assert.ok(texture);const image=gltf.images[gltf.textures[texture.index].source];
  assert.ok(Number.isInteger(image.bufferView));assert.equal(image.mimeType,'image/png');
 }
});
