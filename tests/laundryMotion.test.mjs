import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {LaundryMotion,laundryRunningMachines} from '../.runtime/frontend-test/laundryMotion.js';
const healthy={trading:1,staff:3,supply:20,condition:100,trouble:false};
test('washer workload reflects public operation without guessing missing state',()=>{
 assert.equal(laundryRunningMachines(healthy),3);assert.equal(laundryRunningMachines({...healthy,trading:.4}),2);
 for(const patch of [{trading:0},{trading:NaN},{trading:undefined},{staff:0},{supply:0},{condition:39},{trouble:true}])assert.equal(laundryRunningMachines({...healthy,...patch}),0);
 assert.equal(laundryRunningMachines(),0);assert.equal(laundryRunningMachines({...healthy,trading:1.4}),3);
});
test('exported washer pivots turn within fixed seals and pause without catching up',async()=>{
 const b=readFileSync(new URL('../public/art/models/interior-laundry.glb',import.meta.url)),loader=new GLTFLoader();
 loader.register(()=>({name:'geometry-only',loadMaterial(){return Promise.resolve(new THREE.MeshBasicMaterial());}}));
 const room=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
 const motion=new LaundryMotion(room);assert.equal(motion.count,3);
 const drums=[0,1,2].map(i=>room.getObjectByName(`laundry-drum-${i}`));
 const fixed=[];room.traverse(o=>{if(o.isMesh&&!drums.some(d=>{let p=o;while(p){if(p===d)return true;p=p.parent;}return false;}))fixed.push([o,o.matrixWorld.clone()]);});
 const first=drums[0].quaternion.clone(),last=drums[2].quaternion.clone();
 for(let frame=0;frame<720;frame++){
  assert.equal(motion.step(1/60,2,true,false),true);room.updateMatrixWorld(true);
  for(let i=0;i<3;i++){
   const center=drums[i].getWorldPosition(new THREE.Vector3());
   assert.ok(Math.abs(center.x-[-3.35,-1.35,.65][i])<.001&&Math.abs(center.y-.94)<.001&&Math.abs(center.z+3.815)<.001);
   drums[i].traverse(o=>{if(!o.isMesh)return;const points=o.geometry.attributes.position;const v=new THREE.Vector3();for(let j=0;j<points.count;j++){
    v.fromBufferAttribute(points,j).applyMatrix4(o.matrixWorld).sub(center);
    assert.ok(Math.hypot(v.x,v.y)<.402,'laundry clips fixed rubber seal');
    assert.ok(Math.abs(v.z)<.05,'drum moves through cabinet or door');
   }});
  }
 }
 assert.ok(first.angleTo(drums[0].quaternion)>.1);assert.ok(last.angleTo(drums[2].quaternion)<1e-7,'idle washer rotates');
 for(const [o,matrix] of fixed)assert.deepEqual(o.matrixWorld.elements,matrix.elements,'static furniture moves');
 const seconds=motion.seconds,pose=drums[0].quaternion.clone();
 for(const args of [[1,3,false,false],[1,3,true,true],[1,0,true,false],[NaN,3,true,false]])assert.equal(motion.step(...args),false);
 assert.equal(motion.seconds,seconds);assert.ok(pose.angleTo(drums[0].quaternion)<1e-7);
 motion.step(120,3,true,false);assert.ok(Math.abs(motion.seconds-seconds-.05)<1e-8,'resume jumps after background delay');
});
