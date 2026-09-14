import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {readFileSync} from 'node:fs';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {InteriorCastBatch} from '../.runtime/frontend-test/interiorCastBatch.js';
import {dressPedestrian,wardrobe} from '../.runtime/frontend-test/city3dWardrobe.js';
import {interiorPlacements,poseInteriorOccupant} from '../.runtime/frontend-test/interiorStaging.js';

test('batch retains joint transforms, individual colors, hidden parts and pick identity',()=>{
 const geometry=new THREE.BoxGeometry(.2,.2,.2),source=new THREE.MeshStandardMaterial();
 const actors=new Map(),joints=[];
 for(let i=0;i<3;i++){
  const actor=new THREE.Group(),joint=new THREE.Group(),material=source.clone();material.userData.castSource=source.uuid;material.color.setRGB(.2+i*.2,.3,.4);
  const part=new THREE.Mesh(geometry,material);joint.add(part);actor.add(joint);actor.position.set(i*2,1,0);joint.position.y=.5;actors.set('person-'+i,actor);joints.push(joint);
 }
 joints[2].visible=false;
 const batch=new InteriorCastBatch(actors),mesh=batch.root.children[0];
 assert.equal(batch.root.children.length,1);assert.equal(mesh.count,2);assert.deepEqual(mesh.userData.people,['person-0','person-1']);
 const matrix=new THREE.Matrix4(),color=new THREE.Color();mesh.getMatrixAt(1,matrix);assert.equal(matrix.elements[12],2);assert.equal(matrix.elements[13],1.5);
 mesh.getColorAt(1,color);assert.ok(Math.abs(color.r-.4)<1e-6);assert.equal(mesh.material.color.r,1);
 joints[1].rotation.z=.6;joints[2].visible=true;batch.update();assert.equal(mesh.count,3);
 mesh.getMatrixAt(1,matrix);assert.ok(matrix.elements.every((v,i)=>Math.abs(v-joints[1].children[0].matrixWorld.elements[i])<1e-6));
 mesh.updateMatrixWorld(true);const ray=new THREE.Raycaster(new THREE.Vector3(4,1.5,3),new THREE.Vector3(0,0,-1));
 const hit=ray.intersectObject(mesh)[0];assert.ok(hit);assert.equal(mesh.userData.people[hit.instanceId],'person-2');
 let geometryDisposed=0,sourceDisposed=0,copyDisposed=0;geometry.addEventListener('dispose',()=>geometryDisposed++);source.addEventListener('dispose',()=>sourceDisposed++);mesh.material.addEventListener('dispose',()=>copyDisposed++);
 batch.dispose();assert.equal(geometryDisposed,0);assert.equal(sourceDisposed,0);assert.equal(copyDisposed,1);assert.equal(joints[0].children[0].visible,true);assert.equal(batch.root.children.length,0);
});

test('authored seated cast batches without losing geometry or material tint',async()=>{
 const loader=new GLTFLoader();loader.register(()=>({name:'geometry-only',loadMaterial(index){const m=new THREE.MeshStandardMaterial();m.name=['skin','wool suit','ivory shirt'][index%3];return Promise.resolve(m);}}));
 const sources=new Map();for(const name of ['person','woman']){const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url));sources.set(name,(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene);}
 const people=Array.from({length:6},(_,i)=>({id:`guest-${i}`})),spots=interiorPlacements(people),actors=new Map(),expected=[];
 people.forEach((p,i)=>{const name=i%2?'woman':'person',actor=sources.get(name).clone(true);dressPedestrian(actor,name,wardrobe(p.id,i+1));poseInteriorOccupant(actor,spots.get(p.id));actors.set(p.id,actor);
  actor.traverseVisible(o=>{if(o.isMesh)expected.push({id:p.id,geometry:o.geometry,matrix:o.matrixWorld.clone(),color:o.material.color.clone()});});
 });
 const batch=new InteriorCastBatch(actors);let count=0,draws=0;
 for(const mesh of batch.root.children){if(!mesh.visible)continue;draws++;
  for(let i=0;i<mesh.count;i++){count++;const matrix=new THREE.Matrix4(),color=new THREE.Color();mesh.getMatrixAt(i,matrix);mesh.getColorAt(i,color);
   assert.ok(expected.some(p=>p.id===mesh.userData.people[i]&&p.geometry===mesh.geometry&&p.matrix.elements.every((v,k)=>Math.abs(v-matrix.elements[k])<1e-5)&&['r','g','b'].every(k=>Math.abs(p.color[k]-color[k])<1e-5)));
  }
 }
 assert.equal(count,expected.length);assert.ok(draws<expected.length*.65,`${draws} draws for ${expected.length} original parts`);batch.dispose();
});

test('unsupported materials stay visible and available for ordinary picking',()=>{
 const actor=new THREE.Group(),mesh=new THREE.Mesh(new THREE.BoxGeometry(),new THREE.MeshBasicMaterial());actor.add(mesh);
 const batch=new InteriorCastBatch(new Map([['guest',actor]]));
 assert.equal(mesh.visible,true);assert.deepEqual(batch.unbatched,[mesh]);assert.equal(batch.root.children.length,0);batch.dispose();assert.equal(mesh.visible,true);
});
