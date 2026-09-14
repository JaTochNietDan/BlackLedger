import test from 'node:test';
import assert from 'node:assert/strict';
import {policeCast,availableSceneSlot,officerApproach,raidEntryPose,policeSceneSeconds,BlastAudio} from '../.runtime/frontend-test/city3dEvents.js';
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
 assert.equal(cast.length,4);assert.deepEqual(cast.find(c=>c.kind==='detainee').actors,[{id:'player',name:'Alex'}]);
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
   assert.ok(Math.abs(slot.root.z+pose.travelled-slot.pose.z)+.7<=3.300001,'actor exceeds reserved sweep');
   previous=pose.travelled;
  }
  assert.equal(officerApproach(3,distance,officer).travelled,distance);
  assert.equal(officerApproach(3,distance,officer).walking,false);
 }
 assert.equal(officerApproach(.5,2.4,3).travelled,0);
 assert.ok(officerApproach(.5,2.4,0).travelled>0);
});

test('raid entry waits for the door to clear and stays inside its reserved corridor',()=>{
 const slot=availableSceneSlot({id:'bar',x:80,z:48,row:1,col:2},'raid-officer',[]);
 for(const distance of [2.4,2.92,3])for(let t=0;t<=7;t+=.01){
  const pose=raidEntryPose(t,distance);
  assert.ok(Math.abs(slot.root.z+pose.travelled-slot.pose.z)+.7<=3.300001);
  if(pose.travelled>distance)assert.equal(pose.door,1,'entry preceded door opening');
 }
 assert.equal(raidEntryPose(7,2.92).walking,false);
 for(const cue of policeCast({id:'raid',kind:'raid',target:'bar'}))assert.equal(policeSceneSeconds(cue.kind),10);
});

test('raid approach reservation leaves the public entrance clear',()=>{
 const slot=availableSceneSlot({id:'bar',x:80,z:48,row:1,col:2},'raid-officer',[]);
 assert.equal(trafficOverlap(slot.pose,slot.model,{x:80,z:36.65,heading:0},'person'),false);
});

test('delayed breach audio fires only at visible contact and consumes muted or missed beats',()=>{
 let played=0,stopped=0;const make=()=>new BlastAudio(()=>{played++;return()=>stopped++;});
 const audio=make();audio.update(-2);audio.update(-.01);assert.equal(played,0);
 audio.update(.02);audio.update(.1);assert.equal(played,1);audio.dispose();assert.equal(stopped,1);
 const muted=make();muted.update(-1);muted.update(.01,false);muted.update(.03,true);assert.equal(played,1);
 const late=make();late.update(.4);late.update(.5);assert.equal(played,1);
});

test('a casualty in the central forecourt leaves a shorter door-aligned breach corridor',()=>{
 const lot={id:'bar',x:80,z:48,row:1,col:2},entry={x:80,z:41.92};
 const occupied=[{pose:{x:80.8,z:38.35,heading:0},model:'casualty'},
  {pose:{x:80,z:36.65,heading:0},model:'person'}];
 const slot=availableSceneSlot(lot,'raid-officer',occupied,entry);
 assert.ok(slot);assert.equal(slot.root.x,entry.x);
 assert.ok(slot.root.z>38.35);
 for(const other of occupied)assert.equal(trafficOverlap(slot.pose,slot.model,other.pose,other.model),false);
 const distance=entry.z-slot.root.z-.65;
 for(let t=0;t<=10;t+=.01){
  const pose=raidEntryPose(t,distance);
  assert.ok(Math.abs(slot.root.z+pose.travelled-slot.pose.z)+.7<=3.300001);
 }
});

test('shallower buildings stage close enough to reach their doorway',()=>{
 const lot={id:'burlesque',x:80,z:48,row:1,col:2},entry={x:80,z:42.92};
 const slot=availableSceneSlot(lot,'raid-officer',[],entry);
 assert.equal(slot.root.x,entry.x);assert.ok(entry.z-slot.root.z-.65<=3);
});
