import test from 'node:test';import assert from 'node:assert/strict';
import {streetAt} from '../.runtime/frontend-test/streetPlayback.js';
const leg={id:'mara',name:'Mara',from_id:'bar',to_id:'market',from:'Bar',to:'Market',because:'Work',yours:false,minutes:10,progress:0,from_minute:105,to_minute:115,end_progress:1};
test('NPC departures and completed arrivals use the shared minute',()=>{
 assert.equal(streetAt([leg],100).size,0);
 assert.equal(streetAt([leg],105).get('mara').progress,0);
 assert.equal(streetAt([leg],110).get('mara').progress,.5);
 assert.equal(streetAt([leg],115).size,0);
});
test('partway journeys end at their committed progress without inventing arrival',()=>{
 const partial={...leg,progress:.25,end_progress:.75};
 assert.equal(streetAt([partial],110).get('mara').progress,.5);
 assert.equal(streetAt([partial],115).get('mara').progress,.75);
 assert.equal(streetAt([partial],116).size,0);
});
test('one character can make consecutive legs without duplicate actors',()=>{
 const second={...leg,from_id:'market',to_id:'club',from_minute:115,to_minute:125};
 const samples=streetAt([leg,second],115);assert.equal(samples.size,1);assert.equal(samples.get('mara').segment.to_id,'club');assert.equal(samples.get('mara').progress,0);
});

test('journey clock preserves time for collision-delayed walking without stopping the world',async()=>{
 const {advanceJourneyClock}=await import('../.runtime/frontend-test/streetPlayback.js');
 const {StreetTraffic,trafficOverlap}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const finishTimes=[];
 for(const fps of [30,60,144])for(const rate of [1,4]){
  const traffic=new StreetTraffic();let clock=0,playerProgress=0,waitFrames=0,elapsed=0;
  const player={id:'player',model:'person',points:[{x:80,z:36.65},{x:80,z:27.35},{x:59.35,z:27.35},{x:59.35,z:68.65},{x:112,z:68.65}],progress:0};
  const routeSeconds=123.9/1.8;
  for(let frame=0;frame<fps*180/rate;frame++){
   const step=rate/fps;elapsed+=step;
   const previous=clock;
   clock=advanceJourneyClock(clock,playerProgress,step,routeSeconds,routeSeconds);
   assert.ok(clock>=previous&&clock<1,'the clock must remain short of arrival while the player is en route');
   player.progress=Math.min(1,elapsed/routeSeconds);
   const car={id:'car',model:'packard',points:[{x:48,z:33.6},{x:84,z:33.6}],progress:Math.min(1,clock/.65)};
   const requests=clock<.65?[player,car]:[player];
   const placements=traffic.update(requests,1/fps,rate),at=placements.get('player');
   if(at.blockedBy){waitFrames++;assert.ok(clock>previous,'world time must advance so the obstruction can clear');}
   if(placements.has('car')&&!placements.get('car').waiting)assert.equal(trafficOverlap(at.pose,'person',placements.get('car').pose,'packard'),false);
   playerProgress=at.progress;
   if(playerProgress>=1){
    clock=advanceJourneyClock(clock,playerProgress,0,routeSeconds,routeSeconds);
    break;
   }
  }
  assert.ok(waitFrames>0);assert.equal(playerProgress,1);assert.equal(clock,1);
  finishTimes.push(elapsed);
 }
 assert.ok(Math.max(...finishTimes)-Math.min(...finishTimes)<1,'frame rate and speed settings must not materially alter simulated travel duration');
});

test('unobstructed travel retains its pace and paused playback cannot spend time',async()=>{
 const {advanceJourneyClock}=await import('../.runtime/frontend-test/streetPlayback.js');
 let clock=0;
 for(let step=0;step<600;step++){
  clock=advanceJourneyClock(clock,step/600,1/60,10,10);
  assert.ok(Math.abs(clock-(step+1)/600)<1e-6);
 }
 assert.equal(advanceJourneyClock(.45,.25,0,10,10),.45);
 assert.equal(advanceJourneyClock(.45,.25,-1,10,10),.45);
 assert.equal(advanceJourneyClock(.99,1,0,10,10),1);
});

test('delayed NPCs finish their current leg before a later recorded departure',async()=>{
 const {StreetPlayback,streetLegKey}=await import('../.runtime/frontend-test/streetPlayback.js');
 const second={...leg,from_id:'market',to_id:'club',from_minute:116,to_minute:126};
 const timeline=new StreetPlayback([second,leg]);let arrived='';
 const sample=minute=>timeline.sample(minute,(id,key)=>key===arrived);
 assert.equal(sample(100).size,0);
 assert.equal(sample(110).get('mara').progress,.5);
 assert.equal(sample(119).get('mara').segment,leg,'an overdue first leg must not disappear or jump to its successor');
 assert.equal(sample(119).get('mara').progress,1);
 arrived=streetLegKey(leg);
 assert.equal(sample(120).get('mara').segment,second);
 assert.equal(sample(120).get('mara').progress,.4);
 assert.equal(sample(130).get('mara').segment,second);
 arrived=streetLegKey(second);assert.equal(sample(130).size,0);
 assert.equal(sample(140).size,0,'finished legs must not reappear');
});

test('partial observed legs stop at their saved fraction without inventing an arrival',async()=>{
 const {StreetPlayback}=await import('../.runtime/frontend-test/streetPlayback.js');
 const partial={...leg,progress:.25,end_progress:.75};
 const timeline=new StreetPlayback([partial]);
 assert.equal(timeline.sample(120,()=>true).get('mara').progress,.75);
 const completed=new StreetPlayback([leg]);
 assert.equal(completed.sample(114,()=>true).get('mara').progress,.9,'rendering ahead cannot erase a leg before its recorded arrival');
 assert.equal(completed.sample(115,()=>true).size,0);
});

test('a queued NPC leg enters at the door after physical collision recovery',async()=>{
 const {StreetPlayback,streetLegKey}=await import('../.runtime/frontend-test/streetPlayback.js');
 const {StreetTraffic}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const second={...leg,from_id:'market',to_id:'club',from_minute:116,to_minute:126};
 const timeline=new StreetPlayback([leg,second]),traffic=new StreetTraffic();
 let currentKey='',arrived=false,previousPose,finished=0;
 for(let frame=0;frame<1800;frame++){
  const samples=timeline.sample(130,(id,key)=>key===currentKey&&arrived),sample=samples.get('mara');
  if(!sample){assert.equal(finished,2);break;}
  const key=streetLegKey(sample.segment),starting=key!==currentKey;
  const points=sample.segment===leg?[{x:100,z:16},{x:120,z:16}]:[{x:120,z:16},{x:140,z:16}];
  const requests=[{id:'mara',model:'person',points,progress:starting?sample.segment.progress:sample.progress}];
  if(frame<180)requests.push({id:'obstacle',model:'parked-ford',points:[{x:107,z:16}],progress:0});
  const placed=traffic.update(requests,1/60).get('mara');
  if(previousPose)assert.ok(Math.hypot(placed.pose.x-previousPose.x,placed.pose.z-previousPose.z)<=1.8/60+1e-8,'late successor must not teleport down its route');
  previousPose={...placed.pose};currentKey=key;arrived=placed.progress>=1;
  if(arrived)finished++;
  if(frame===179){assert.equal(sample.segment,leg);assert.equal(placed.blockedBy,'obstacle');}
 }
 assert.equal(finished,2);
});
