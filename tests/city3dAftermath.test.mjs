import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {CityAftermath} from '../.runtime/frontend-test/city3dAftermath.js';
import {cityPlan} from '../.runtime/frontend-test/city3dPlan.js';
import {trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';

test('aftermath persists, reserves separate body/police bays and clears on cleanup',()=>{
 const city=new CityAftermath();
 const lot={id:'bar',x:80,z:48,row:1,col:2};const lots=new Map([['bar',lot]]);
 const person=new THREE.Group();person.add(new THREE.Mesh(new THREE.BoxGeometry(.7,1.8,.4),new THREE.MeshStandardMaterial()));
 const models=new Map([['person',person],['police',new THREE.Group()],['police-officer',person]]);
 const records=[{id:'death:mara',target:'bar',victim:{id:'mara',name:'Mara'},minute:480,police_at:485,cleanup_at:660}];
 const update=(minute,animating=new Set())=>city.update(records,minute,lots,models,()=> 'person',[],animating);
 update(480);assert.equal(city.reservations().length,1);
 const root=city.root.children[0];update(484);assert.equal(city.root.children[0],root);
 update(485);assert.equal(city.reservations().length,4);
 for(const [i,a] of city.slots().entries())for(const b of city.slots().slice(i+1))assert.equal(trafficOverlap(a.pose,a.model,b.pose,b.model),false);
 update(486,new Set(['mara']));assert.equal(city.reservations().length,0,'replay hides the corpse and its later police response');
 update(487);assert.equal(city.reservations().length,4);
 update(660);assert.equal(city.reservations().length,0);assert.equal(city.root.children.length,0);
 city.dispose();
});

test('raid cordon persists independently, yields to playback, and releases all reservations at cleanup',()=>{
 const city=new CityAftermath(), lot={id:'bar',x:80,z:48,row:1,col:2};
 const models=new Map([['police',new THREE.Group()],['police-officer',new THREE.Group()]]);
 const presence=[{id:'raid-one',target:'bar',minute:480,cleanup_at:600}];
 const update=(minute,active=new Set())=>city.update([],minute,new Map([['bar',lot]]),models,()=> 'person',[],new Set(),presence,active);
 update(480);assert.equal(city.reservations().length,7);
 assert.equal(city.inspect().some(e=>e.id.endsWith(':body')),false,'raid invents no casualty');
 const object=city.root.children[0];update(599);assert.equal(city.root.children[0],object);
 for(const [i,a] of city.slots().entries())for(const b of city.slots().slice(i+1))assert.equal(trafficOverlap(a.pose,a.model,b.pose,b.model),false);
 update(599,new Set(['bar']));assert.equal(city.reservations().length,0,'playback takes ownership of scene');
 update(599);assert.equal(city.reservations().length,7);
 update(600);assert.equal(city.reservations().length,0);assert.equal(city.root.children.length,0);
 city.dispose();
});

test('execution hands off the body at its actual fall position and yaw',()=>{
 const city=new CityAftermath(),lot={id:'bar',x:80,z:48,row:1,col:2};
 const models=new Map([['person',new THREE.Group()]]),lots=new Map([['bar',lot]]);
 const records=[{id:'death:mara',target:'bar',victim:{id:'mara'},minute:480,police_at:485,cleanup_at:660}];
 city.update(records,480,lots,models,()=> 'person',[],new Set());
 city.suppressVictims(new Set(['mara']));assert.equal(city.slots().length,0,'old body must not block shared scene');
 const slot={root:{x:81,z:38.35},pose:{x:81.8,z:38.35,heading:0},model:'casualty'};
 city.rememberBody('mara',slot,Math.PI/2);
 city.update(records,480,lots,models,()=> 'person',[],new Set(['mara']));assert.equal(city.slots().length,0);
 city.update(records,480,lots,models,()=> 'person',[],new Set());
 assert.equal(city.inspect()[0].x,81);assert.equal(city.inspect()[0].z,38.35);
 const object=city.root.children[0].children.at(-1);
 const expected=new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,1,0),Math.PI/2)
  .premultiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,0,1),-Math.PI/2));
 assert.ok(object.quaternion.angleTo(expected)<1e-7);
 city.update(records,660,lots,models,()=> 'person',[],new Set());assert.equal(city.slots().length,0);
 city.dispose();
});

test('a replay yields only its own later response and retains an older nearby crime scene',()=>{
 const city=new CityAftermath(),lot={id:'casino',x:112,z:112,row:3,col:3};
 const lots=new Map([['casino',lot]]),models=new Map([['person',new THREE.Group()],['police',new THREE.Group()],['police-officer',new THREE.Group()]]);
 const records=[{id:'older',target:'casino',victim:{id:'piet'},minute:1605,police_at:1610,cleanup_at:1785},
  {id:'current',target:'casino',victim:{id:'zoltan'},minute:1665,police_at:1670,cleanup_at:1845}];
 const update=animating=>city.update(records,1725,lots,models,()=> 'person',[],animating);
 update(new Set());const older=city.inspect().filter(e=>e.id.startsWith('aftermath:older:'));
 assert.equal(older.length,4);
 city.suppressVictims(new Set(['zoltan']));
 assert.deepEqual(city.inspect(),older,'unrelated attendance changed');
 update(new Set(['zoltan']));assert.deepEqual(city.inspect(),older,'response returned during replay');
 update(new Set());assert.ok(city.inspect().some(e=>e.id.startsWith('aftermath:current:')),'response did not return after playback');
 city.dispose();
});
