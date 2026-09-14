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
