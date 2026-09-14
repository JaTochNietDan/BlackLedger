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
