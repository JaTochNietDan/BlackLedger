import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {CitySuppression,waterArc,impactSpray} from '../.runtime/frontend-test/city3dSuppression.js';

test('water arcs meet the nozzle and window without overshooting the facade',()=>{
 const a=new THREE.Vector3(77,1.5,38),b=new THREE.Vector3(80,5,42);
 assert.deepEqual(waterArc(a,b,0).toArray(),a.toArray());assert.deepEqual(waterArc(a,b,1).toArray(),b.toArray());
 for(let i=0;i<=100;i++){const p=waterArc(a,b,i/100);assert.ok(p.z>=38&&p.z<=42&&p.y>=1.5&&p.y<=5.65);}
});
test('suppression requires an attending crew, freezes motion and releases hose geometry at extinguishing',()=>{
 const effect=new CitySuppression(),building=new THREE.Group();building.position.set(80,0,48);
 const vent=new THREE.Object3D();vent.name='fire-window-1-0';vent.position.set(0,5,-6);building.add(vent);building.updateMatrixWorld(true);
 const engine=new THREE.Group();engine.position.set(90,.2,48);
 const crew=()=>{const g=new THREE.Group(),actor=new THREE.Group();for(const name of ['arm1','arm-1']){const arm=new THREE.Group();arm.name=name;actor.add(arm);}g.add(actor);return g;};
 const a=crew(),b=crew();a.position.set(77,.2,38.35);b.position.set(83,.2,38.35);
 const objects=new Map([['aftermath:f:fire-engine',engine],['aftermath:f:firefighter-a',a],['aftermath:f:firefighter-b',b]]);
 const nozzle=new THREE.Group(),tip=new THREE.Object3D();tip.name='water-outlet';tip.position.z=.29;nozzle.add(tip);
 const fires=[{id:'f',target:'bar',minute:480,brigade_at:490,extinguished_at:660,cleanup_at:720}];
 const update=(minute,motion=true)=>effect.update(fires,minute,new Map([['bar',building]]),new Map([['fire-nozzle',nozzle]]),{object:id=>objects.get(id)},100,motion);
 update(489);assert.equal(effect.inspect().length,0);engine.visible=false;update(490);assert.equal(effect.inspect().length,0);
 engine.visible=true;update(490);assert.equal(effect.inspect().length,2);
 const particles=effect.root.children[0].children.filter(o=>o instanceof THREE.Points);
 const before=particles.map(o=>Array.from(o.geometry.attributes.position.array));
 update(500,false);assert.deepEqual(particles.map(o=>Array.from(o.geometry.attributes.position.array)),before);
 const jet=effect.root.children[0].children.find(o=>o instanceof THREE.Mesh&&o.geometry.parameters.radius===.023);
 const arc=jet.geometry.parameters.path;
 for(let i=0;i<=20;i++)assert.ok(arc.getPoint(i/20).distanceTo(waterArc(arc.v0,arc.v2,i/20))<1e-10);
 let disposed=0;effect.root.traverse(o=>{if(o.geometry)o.geometry.addEventListener('dispose',()=>disposed++);});
 update(660);assert.equal(effect.inspect().length,0);assert.equal(disposed,8);assert.equal(a.getObjectByName('arm1').rotation.x,0);
 effect.dispose();
});

test('window impact spray remains outside the facade and close to its impact',()=>{
 const to=new THREE.Vector3(80,5,42);
 for(let i=0;i<24;i++)for(let t=0;t<2;t+=.01){
  const p=impactSpray(to,i,t);assert.ok(p.z<to.z&&p.z>=to.z-.56);
  assert.ok(Math.abs(p.x-to.x)<.6&&p.y>=to.y-1.1&&p.y<=to.y+.1);
 }
});
