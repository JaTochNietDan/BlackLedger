import test from 'node:test';
import assert from 'node:assert/strict';
import {deflateSync} from 'node:zlib';
import {decodePoolReplay,poolFramePair} from '../.runtime/frontend-test/billiards.js';
const balls=Array.from({length:16},(_,id)=>[id,.5,1,.028575,0,0,0,1,0]);
const tape={v:1,duration:1,frames:[{t:0,balls},{t:.39,balls},{t:.4,balls},{t:1,balls}],events:[]};
const encode=r=>deflateSync(JSON.stringify(r)).toString('base64');
test('billiards replay decodes server zlib format and retains collision timestamps',async()=>{
 const r=await decodePoolReplay(encode(tape));assert.deepEqual(r,tape);
 const before=poolFramePair(r,.395);assert.equal(before.a.t,.39);assert.equal(before.b.t,.4);assert.ok(Math.abs(before.mix-.5)<1e-10);
 const hit=poolFramePair(r,.4);assert.equal(hit.a.t,.4);assert.equal(hit.mix,0);
 assert.equal(poolFramePair(r,9).a.t,1);assert.equal(poolFramePair(r,-1).mix,0);
});
test('billiards replay rejects malformed ball identities and frame order',async()=>{
 const bad=structuredClone(tape);bad.frames[0].balls[1][0]=0;await assert.rejects(decodePoolReplay(encode(bad)),/ball/);
 const reverse=structuredClone(tape);reverse.frames[2].t=.1;await assert.rejects(decodePoolReplay(encode(reverse)),/frame/);
 await assert.rejects(decodePoolReplay(encode({...tape,v:2})),/Invalid/);
 await assert.rejects(decodePoolReplay('invalid replay'));
});
test('billiards replay bounds compressed and decompressed sizes',async()=>{
 await assert.rejects(decodePoolReplay('x'.repeat(2*1024*1024+1)),/size limit/);
 const bomb=deflateSync(' '.repeat(8*1024*1024+1)).toString('base64');await assert.rejects(decodePoolReplay(bomb),/size limit/);
});
