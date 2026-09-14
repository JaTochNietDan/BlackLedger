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

const {PoolTap,poolAimAngle,poolPlacementHint,poolPocketCenters}=await import('../.runtime/frontend-test/billiards.js');
test('table taps never reinterpret an orbit or multi-pointer gesture as cue input',()=>{
 const tap=new PoolTap();tap.begin(1,100,100,true,0);assert.equal(tap.end(1,103,102),true);
 tap.begin(1,100,100,true,0);tap.move(1,140,100);assert.equal(tap.end(1,100,100),false);
 tap.begin(1,100,100,true,0);tap.begin(2,110,100,false,0);assert.equal(tap.end(1,100,100),false);
 tap.begin(1,100,100,true,2);assert.equal(tap.end(1,100,100),false);
 tap.begin(1,100,100,true,0);tap.cancel();assert.equal(tap.end(1,100,100),false);
});
test('cue aim follows cloth axes and placement preview respects occupied space/head string',()=>{
 assert.equal(poolAimAngle([.5,.5],[.5,1]),Math.PI/2);
 assert.equal(poolAimAngle([.5,.5],[0,.5]),Math.PI);
 assert.equal(poolAimAngle([.5,.5],[.5,.5]),null);
 assert.equal(poolAimAngle([.5,.5],[NaN,.5]),null);
 const p={width:1.27,length:2.54,radius:.028575,behind_head_string:true,balls:[{id:0,position:[.5,.3,0],pocket:-1},{id:1,position:[.7,.3,0],pocket:-1}]};
 assert.equal(poolPlacementHint(p,.5,.3),'');
 assert.match(poolPlacementHint(p,.7,.3),/room/);
 assert.match(poolPlacementHint(p,.5,.635),/head string/);
 assert.match(poolPlacementHint(p,0,.3),/cushions/);
 assert.equal(poolPlacementHint({...p,behind_head_string:false},.5,1),'');
 const pockets=poolPocketCenters(1.27,2.54);assert.equal(pockets.length,6);assert.deepEqual(pockets[4],[-.045,1.27]);assert.deepEqual(pockets[5],[1.315,1.27]);
});

const {poolCueStroke,POOL_CUE_CONTACT,POOL_CUE_END}=await import('../.runtime/frontend-test/billiards.js');
test('cue contacts the sphere before physical replay advances, including off-centre spin',()=>{
 const radius=.028575,top=.012,side=-.01;
 for(const speed of [.1,4,8]){
  const start=poolCueStroke(0,speed,radius,top,side),back=poolCueStroke(.42,speed,radius,top,side),contact=poolCueStroke(POOL_CUE_CONTACT,speed,radius,top,side);
  assert.ok(back.front<start.front);assert.equal(back.ballTime,0);
  assert.ok(Math.abs(contact.front**2+top**2+side**2-radius**2)<1e-12);assert.equal(contact.ballTime,0);assert.equal(contact.contact,true);
  assert.ok(poolCueStroke(.719,speed,radius).front < -radius);
  assert.ok(Math.abs(poolCueStroke(.82,speed,radius).ballTime-.1)<1e-12);
  assert.equal(poolCueStroke(POOL_CUE_END,speed,radius).visible,false);
 }
});
test('cue motion is continuous at phase boundaries',()=>{
 for(const t of [.42,POOL_CUE_CONTACT,.88,POOL_CUE_END])assert.ok(Math.abs(poolCueStroke(t-1e-8,7,.028575).front-poolCueStroke(t+1e-8,7,.028575).front)<1e-6);
});
