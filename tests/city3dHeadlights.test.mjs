import test from 'node:test';
import assert from 'node:assert/strict';
import {headlightAlpha,headlightCentre} from '../.runtime/frontend-test/city3dHeadlights.js';

test('headlight pools remain ahead of each bumper and rotate symmetrically around the vehicle',()=>{
 for(const length of [4.7,5.1,5.8])for(let heading=-Math.PI;heading<=Math.PI;heading+=.1){
  const left=headlightCentre(17,32,heading,length,-1),right=headlightCentre(17,32,heading,length,1);
  assert.ok(Math.abs(Math.hypot(left.x-right.x,left.z-right.z)-1.22)<1e-10);
  for(const point of [left,right]){
   const forward=(point.x-17)*Math.sin(heading)+(point.z-32)*Math.cos(heading);
   assert.ok(Math.abs(forward-length/2-3.5)<1e-10);
   assert.ok(forward-3.5>=length/2-1e-10,'near edge begins at bumper');
  }
 }
});

test('headlight masks taper to transparent edges without rectangular borders',()=>{
 for(let y=0;y<=1;y+=.01){
  assert.equal(headlightAlpha(-1,y),0);assert.equal(headlightAlpha(1,y),0);
  for(let x=-1;x<=1;x+=.05){
   const alpha=headlightAlpha(x,y);assert.ok(alpha>=0&&alpha<=1);
   assert.ok(Math.abs(alpha-headlightAlpha(-x,y))<1e-12);
   assert.equal(headlightAlpha(x,0),0);assert.equal(headlightAlpha(x,1),0);
  }
 }
 assert.ok(headlightAlpha(0,.2)>headlightAlpha(.8,.2));
 assert.ok(headlightAlpha(0,.2)>headlightAlpha(0,.8));
});
