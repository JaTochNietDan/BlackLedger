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
 update(486,new Set(['mara']));assert.equal(city.reservations().length,3,'replay hides duplicate corpse');
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
