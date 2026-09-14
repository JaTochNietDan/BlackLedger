import test from 'node:test';
import assert from 'node:assert/strict';
import {poolAimContact} from '../.runtime/frontend-test/billiardsAim.js';
const ball=(id,x,y,pocket=-1)=>({id,position:[x,y,.028575],pocket});
const table=(balls)=>({width:1.27,length:2.54,radius:.028575,balls});
test('aim stops one ball diameter before the first object ball',()=>{
 const p=table([ball(0,.4,.5),ball(2,.4,1.5),ball(1,.4,1)]),c=poolAimContact(p,Math.PI/2);
 assert.equal(c.ball,1);assert.equal(c.kind,'ball');assert.ok(Math.abs(c.y-(1-2*p.radius))<1e-10);
 assert.deepEqual(poolAimContact({...p,balls:[...p.balls].reverse()},Math.PI/2),c);
});
test('aim handles cuts, misses, pocketed balls and a cushion before a ball',()=>{
 const p=table([ball(0,.4,.5),ball(1,.43,1),ball(2,.4,.7,0)]),c=poolAimContact(p,Math.PI/2);
 assert.equal(c.ball,1);assert.ok(Math.abs(Math.hypot(c.x-.43,c.y-1)-2*p.radius)<1e-10);
 const miss=poolAimContact(table([ball(0,.4,.5),ball(1,.48,1)]),Math.PI/2);assert.equal(miss.kind,'cushion');assert.ok(Math.abs(miss.y-(2.54-p.radius))<1e-10);
 const rail=poolAimContact(table([ball(0,.4,.5),ball(1,1.5,.5)]),0);assert.equal(rail.kind,'cushion');assert.ok(Math.abs(rail.x-(1.27-p.radius))<1e-10);
});
test('pocket openings do not act like a continuous cushion',()=>{
 const p=table([ball(0,.4,1.27)]);const c=poolAimContact(p,Math.PI);assert.equal(c.kind,'edge');assert.ok(c.x<0);
 const jaw=poolAimContact(table([ball(0,.4,1.27-.075)]),Math.PI);assert.equal(jaw.kind,'cushion');assert.ok(jaw.x>=0);
});
test('touching balls allow aiming away and reads leave the table unchanged',()=>{
 const p=table([ball(0,.4,.5),ball(1,.4+2*.028575,.5)]),before=JSON.stringify(p);
 assert.equal(poolAimContact(p,0).distance,0);assert.equal(poolAimContact(p,Math.PI).kind,'cushion');assert.equal(JSON.stringify(p),before);
 assert.equal(poolAimContact(p,NaN),null);
});
