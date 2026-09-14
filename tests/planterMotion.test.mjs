import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityPlanter,PLANTER_BLAST} from '../.runtime/frontend-test/city3dPlanter.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name,side:THREE.DoubleSide}));}}));
 return (await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;
}
test('planter exits the authored vestibule without crossing the door or masonry',async()=>{
 for(const name of ['monarch','tavern'])for(const person of ['person','woman']){
  const building=await model(name),actor=await model(person),cast=new CityPlanter(actor);
  building.updateMatrixWorld(true);
  cast.root.position.copy(building.getObjectByName('entrance-threshold').getWorldPosition(new THREE.Vector3()));
  const door=building.getObjectByName('entrance-door-hinge'),triangle=new THREE.Triangle(),bounds=new THREE.Box3();
  for(let frame=0;frame<=186;frame++){
   const pose=cast.update(frame/30);door.rotation.y=-Math.PI/2*pose.door;building.updateMatrixWorld(true);
   const actorBounds=new THREE.Box3().setFromObject(actor,true);actorBounds.min.y+=.12;
   building.traverse(o=>{
    if(!(o instanceof THREE.Mesh))return;
    o.geometry.computeBoundingBox();bounds.copy(o.geometry.boundingBox).applyMatrix4(o.matrixWorld);
    if(!actorBounds.intersectsBox(bounds))return;
    const positions=o.geometry.attributes.position,index=o.geometry.index,count=index?.count??positions.count;
    for(let i=0;i<count;i+=3){
     triangle.a.fromBufferAttribute(positions,index?index.getX(i):i).applyMatrix4(o.matrixWorld);
     triangle.b.fromBufferAttribute(positions,index?index.getX(i+1):i+1).applyMatrix4(o.matrixWorld);
     triangle.c.fromBufferAttribute(positions,index?index.getX(i+2):i+2).applyMatrix4(o.matrixWorld);
     assert.equal(actorBounds.intersectsTriangle(triangle),false,`${name}/${person} intersects ${o.name} at ${frame/30}`);
    }
   });
  }
  const finish=cast.update(PLANTER_BLAST);assert.equal(finish.blast,true);assert.equal(finish.door,0);
  assert.ok(actor.position.z<-2&&actor.position.x<=-2);
  assert.equal(cast.update(PLANTER_BLAST-.001).blast,false);
 }
});
