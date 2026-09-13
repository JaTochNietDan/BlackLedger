import test from 'node:test';
import assert from 'node:assert/strict';
import {policeCast,availableSceneSlot,officerApproach} from '../.runtime/frontend-test/city3dEvents.js';
import {trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';

test('raid stages three vehicles and four officers without shared bays',()=>{
 const cast=policeCast({id:'raid',kind:'raid',target:'bar',caption:''});
 assert.equal(cast.length,7);assert.equal(cast.filter(c=>c.kind==='raid-officer').length,4);
 const lot={id:'bar',x:80,z:48,row:1,col:2},slots=[];
 for(const c of cast){const slot=availableSceneSlot(lot,c.kind,slots);assert.ok(slot,`no space for ${c.kind}`);slots.push(slot);}
 for(let i=0;i<slots.length;i++)for(const b of slots.slice(i+1))assert.equal(trafficOverlap(slots[i].pose,slots[i].model,b.pose,b.model),false);
});
test('arrest uses explicit detainee identity and never substitutes its detective actor',()=>{
 const cue={id:'arrest',kind:'arrest',target:'precinct',caption:'',actors:[{id:'detective',name:'Detective'}]};
 assert.equal(policeCast(cue).filter(c=>c.kind==='detainee').length,0);
 const cast=policeCast({...cue,detainee:{id:'player',name:'Alex'}});
 assert.equal(cast.length,5);assert.deepEqual(cast.find(c=>c.kind==='detainee').actors,[{id:'player',name:'Alex'}]);
 assert.equal(new Set(cast.map(c=>c.id)).size,cast.length);
});

test('officer approach is staggered, speed bounded and stays inside its swept reservation',()=>{
 const lot={id:'bar',x:80,z:48,row:1,col:2};const slot=availableSceneSlot(lot,'raid-officer',[]);
 for(const distance of [0,.15,1.6,2.4])for(let officer=0;officer<4;officer++){
  let previous=0;
  for(let time=0;time<=3;time+=.01){
   const pose=officerApproach(time,distance,officer);
   assert.ok(pose.travelled>=previous&&pose.travelled<=distance);
   assert.ok(pose.travelled-previous<=.014001,'approach exceeds walking speed');
   assert.ok(Math.abs(slot.root.z+pose.travelled-slot.pose.z)+.7<=2,'actor exceeds reserved sweep');
   previous=pose.travelled;
  }
  assert.equal(officerApproach(3,distance,officer).travelled,distance);
  assert.equal(officerApproach(3,distance,officer).walking,false);
 }
 assert.equal(officerApproach(.5,2.4,3).travelled,0);
 assert.ok(officerApproach(.5,2.4,0).travelled>0);
});
