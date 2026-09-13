import test from 'node:test';
import assert from 'node:assert/strict';
import {crossingPaint,PITCH,FOOTWAY,surfaceHeight} from '../.runtime/frontend-test/city3dPlan.js';

test('crossing paint remains on roads and connects existing pavement islands',()=>{
 const marks=crossingPaint(6,5);assert.equal(marks.length,196);
 for(const mark of marks){
  for(const x of [mark.x-mark.width/2,mark.x+mark.width/2])for(const z of [mark.z-mark.depth/2,mark.z+mark.depth/2])assert.equal(surfaceHeight({x,z}),-.1);
  const horizontal=mark.width>mark.depth;
  for(const side of [-1,1]){
   const end={x:mark.x+(horizontal?side*4.3:0),z:mark.z+(horizontal?0:side*4.3)};
   assert.equal(surfaceHeight(end),.17,'crossing must connect two raised pavements');
   assert.ok(end.x>0&&end.x<6*PITCH&&end.z>0&&end.z<5*PITCH,'no paint toward nonexistent exterior pavement');
  }
 }
 assert.equal(new Set(marks.map(m=>`${m.x}:${m.z}`)).size,marks.length);
});

test('every interior junction has a pair of lines around each pedestrian lane',()=>{
 const marks=crossingPaint(6,5);
 for(let col=1;col<6;col++)for(let row=1;row<5;row++)for(const side of [-1,1]){
  const horizontal=marks.filter(m=>m.x===col*PITCH&&Math.abs(m.z-(row*PITCH+side*FOOTWAY))<.6);
  const vertical=marks.filter(m=>m.z===row*PITCH&&Math.abs(m.x-(col*PITCH+side*FOOTWAY))<.6);
  assert.equal(horizontal.length,2);assert.equal(vertical.length,2);
 }
});
