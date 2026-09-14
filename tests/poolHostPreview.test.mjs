import test from 'node:test';
import assert from 'node:assert/strict';
import {poolHostPreview} from '../.runtime/frontend-test/billiards.js';
test('Owner fee preview selects funded field and exactly splits gross entries',()=>{
 const offers=Array.from({length:8},(_,i)=>({id:String(8-i),name:`Player ${8-i}`,max_fee:i<4?100:25}));
 let p=poolHostPreview(offers,25,17,false,0);
 assert.equal(p.reason,'');assert.equal(p.count,8);assert.equal(p.gross,200);assert.equal(p.cut,34);assert.equal(p.prize,166);
 p=poolHostPreview(offers,50,20,false,0);assert.equal(p.count,4);assert.equal(p.prize,160);assert.equal(p.cut,40);assert.deepEqual(p.entrants.map(n=>n.id),['5','6','7','8']);
 assert.match(poolHostPreview(offers,101,20,false,1000).reason,/Need 4/);
 assert.match(poolHostPreview(offers,50,20,true,49).reason,/cannot cover/);
 p=poolHostPreview(offers,50,20,true,50);assert.equal(p.count,4);assert.equal(p.entrants.length,3);
 assert.match(poolHostPreview(offers,NaN,20,false,0).reason,/Set/);
});
