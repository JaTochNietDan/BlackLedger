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
