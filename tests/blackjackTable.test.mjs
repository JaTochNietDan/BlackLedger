import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
async function model(name){
 const b=readFileSync(new URL(`../public/art/models/${name}.glb`,import.meta.url)),loader=new GLTFLoader();
 loader.register(p=>({name:'geometry-only',loadMaterial(i){return Promise.resolve(new THREE.MeshBasicMaterial({name:p.json.materials[i].name}));}}));
 const root=(await loader.parseAsync(b.buffer.slice(b.byteOffset,b.byteOffset+b.byteLength),'')).scene;root.updateMatrixWorld(true);return root;
}
test('blackjack cards rest above the felt and expose upright printed faces',async()=>{
 const table=await model('blackjack-table'),card=await model('playing-card');
 for(const z of [-.38,.48])for(const x of [-.85,0,.85]){
  card.scale.set(1.4,1,1.4);card.position.set(x,.875,z);card.updateMatrixWorld(true);
  const box=new THREE.Box3().setFromObject(card,true);
  for(const cx of [box.min.x,box.max.x])for(const cz of [box.min.z,box.max.z]){
   const hit=new THREE.Raycaster(new THREE.Vector3(cx,2,cz),new THREE.Vector3(0,-1,0)).intersectObject(table,true)[0];
   assert.ok(hit);assert.equal(hit.object.material.name,'blackjack green baize');assert.ok(box.min.y>hit.point.y,'card intersects cloth');
  }
  const hit=new THREE.Raycaster(new THREE.Vector3(x,2,z-.1),new THREE.Vector3(0,-1,0)).intersectObject(card,true)[0];
  assert.equal(hit.object.material.name,'card printed face');assert.ok(hit.uv.y<.5,'card face is upside down');
 }
});

test('revealed hole card turns from its existing seat with its full stock above the felt',async()=>{
 const {planCards,cardPose}=await import('../.runtime/frontend-test/blackjackPresentation.js');
 const first={rank:'8',suit:'diamonds',value:8},hole={rank:'9',suit:'hearts',value:9};
 const before={playing:true,mine:[first],theirs:[first]},after={playing:false,settled:true,mine:[first],theirs:[first,hole]};
 const plan=planCards(before,after,true),move=plan.cards.find(c=>c.flip);
 assert.ok(move);assert.equal(move.row,0);assert.equal(move.index,1);assert.deepEqual(move.from,move.to);
 const card=await model('playing-card');card.scale.set(1.4,1,1.4);
 let face;card.traverse(o=>{if(o instanceof THREE.Mesh&&o.material.name==='card printed face')face=o;});
 const reverse=face.clone();reverse.rotation.z=Math.PI;card.add(reverse);
 for(let frame=0;frame<=120;frame++){
  const pose=cardPose(move,frame/120*move.duration);card.position.fromArray(pose.position);card.rotation.z=pose.rotation;card.updateMatrixWorld(true);
  assert.ok(new THREE.Box3().setFromObject(card,true).min.y>.87,'turning card intersects cloth');
  const normal=new THREE.Vector3(0,1,0).transformDirection(card.matrixWorld);
  if(frame===0)assert.ok(normal.y<-.9999,'face is exposed before turnover');
  if(frame===120)assert.ok(normal.y>.9999,'card does not finish face up');
 }
 assert.equal(planCards(before,after,false).cards.some(c=>c.flip),false);
});

test('only public dealers are staged and both rigs clear the table',async()=>{
 const {blackjackDealer,poseBlackjackDealer}=await import('../.runtime/frontend-test/blackjackDealer.js');
 const patron={id:'patron',name:'Patron',standing:'',role:'Bellandi Family'};
 assert.equal(blackjackDealer([patron]),undefined);
 const dealer={...patron,id:'dealer',name:'Dante',role:'Croupier'};
 assert.equal(blackjackDealer([patron,dealer]),dealer);
 for(const name of ['person','woman']){
  const actor=await model(name);poseBlackjackDealer(actor);
  actor.traverse(o=>{
   if(!(o instanceof THREE.Mesh))return;
   const pos=o.geometry.getAttribute('position');
   for(let i=0;i<pos.count;i++){
    const v=new THREE.Vector3().fromBufferAttribute(pos,i).applyMatrix4(o.matrixWorld);
    assert.ok(!(v.y>.60&&v.y<.875&&v.x*v.x/(1.67*1.67)+v.z*v.z/(1.18*1.18)<1),'dealer enters table apron or rail');
   }
  });
 }
});
