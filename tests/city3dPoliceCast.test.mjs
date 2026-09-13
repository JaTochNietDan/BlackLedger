import test from 'node:test';
import assert from 'node:assert/strict';
import {policeCast,availableSceneSlot} from '../.runtime/frontend-test/city3dEvents.js';
import {trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';

test('raid stages three vehicles and four officers without shared bays',()=>{
 const cast=policeCast({id:'raid',kind:'raid',target:'bar',caption:''});
 assert.equal(cast.length,7);assert.equal(cast.filter(c=>c.kind==='officer').length,4);
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
