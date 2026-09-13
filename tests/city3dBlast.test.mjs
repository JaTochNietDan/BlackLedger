import test from 'node:test';
import assert from 'node:assert/strict';
import {blastParticle,blastLight,blastOpacity,billowAlpha} from '../.runtime/frontend-test/city3dBlast.js';
test('blast fire ends before the rising smoke and all particles finish at three seconds',()=>{
 for(let i=0;i<32;i++){
  for(let frame=0;frame<=180;frame++){
   const p=blastParticle(i,frame/60);
   assert.ok([p.x,p.y,p.z,p.size].every(Number.isFinite));
   assert.ok(p.y>=0&&p.y<9&&p.z<0);
   if(p.size>0)assert.ok(Math.abs(p.x)<4&&Math.abs(p.z)<4);
   assert.ok(p.size>=0&&p.size<4);
  }
  assert.equal(blastParticle(i,3).size,0);
  if(i<12)assert.equal(blastParticle(i,1).size,0);
  else{
   assert.ok(blastParticle(i,1).size>1);
   assert.ok(blastParticle(i,2).y>blastParticle(i,1).y);
  }
 }
 assert.equal(blastOpacity(2),.8);assert.equal(blastOpacity(3),0);
 assert.ok(blastOpacity(2.8)>0&&blastOpacity(2.8)<.8);
 assert.equal(blastLight(0),160);assert.equal(blastLight(.5),0);
});
test('billow texture has a soft transparent boundary and varied interior opacity',()=>{
 for(let i=0;i<=128;i++)for(let j=0;j<=128;j++){
  const x=i/64-1,y=j/64-1,a=billowAlpha(x,y);
  assert.ok(a>=0&&a<=1&&Number.isFinite(a));
  if(Math.hypot(x,y)>1.08)assert.equal(a,0);
 }
 assert.ok(billowAlpha(0,0)>.8);
 assert.notEqual(billowAlpha(.3,.2),billowAlpha(.2,.3));
});

test('debris follows separated facade lanes, settles without continuing to spin, and avoids actors',async()=>{
 const {debrisPose,fragmentBlocked}=await import('../.runtime/frontend-test/city3dBlast.js');
 for(let frame=0;frame<180;frame++){
  const parts=Array.from({length:12},(_,i)=>debrisPose(i,frame/60));
  for(const p of parts){
   assert.ok(Object.values(p).every(Number.isFinite));
   assert.ok(Math.abs(p.x)<=2.81&&p.z>=-2.14&&p.z<=-.25&&p.height>=0);
  }
  for(let i=0;i<parts.length;i++)for(let j=i+1;j<parts.length;j++)
   assert.ok(Math.hypot(parts[i].x-parts[j].x,parts[i].z-parts[j].z)>.26,'fragment envelopes must remain separated');
 }
 for(let i=0;i<12;i++){
  const p=debrisPose(i,1.5),q=debrisPose(i,2.5);
  assert.equal(p.height,0);assert.equal(p.rx,0);assert.equal(p.rz,0);assert.deepEqual(p,q);
  assert.equal(debrisPose(i,3).scale,0);
 }
 for(let heading=0;heading<Math.PI*2;heading+=.1){
  const body={x:10,z:20,heading};
  assert.equal(fragmentBlocked(10,20,body,.85,1.4),true);
  assert.equal(fragmentBlocked(10+Math.sin(heading)*.8,20+Math.cos(heading)*.8,body,.85,1.4),true);
  assert.equal(fragmentBlocked(10+Math.sin(heading)*2,20+Math.cos(heading)*2,body,.85,1.4),false);
 }
});
