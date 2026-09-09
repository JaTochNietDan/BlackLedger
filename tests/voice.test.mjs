import {test} from 'node:test';import assert from 'node:assert/strict';
import {VoicePlayer, speaking, IDLE} from '../.runtime/frontend-test/voice.js';
const deferred=()=>{let resolve,reject;const promise=new Promise((a,b)=>{resolve=a;reject=b});return {promise,resolve,reject}};
function setup(load){const handles=[],states=[];let errors=0;const voice=new VoicePlayer({load,audio(){const h={plays:0,pauses:0,disposed:0,play:async()=>{h.plays++},pause:()=>h.pauses++,dispose:()=>h.disposed++,onEnded:f=>h.end=f};handles.push(h);return h},status:s=>states.push(s),unavailable:()=>errors++});return {voice,handles,states,errors:()=>errors}}
test('late voice after decision never creates or plays audio',async()=>{const pending=deferred();let signal;const s=setup((id,sig)=>{signal=sig;return pending.promise});const play=s.voice.speak('old');s.voice.stop();assert.equal(signal.aborted,true);pending.resolve(new Blob());await play;assert.equal(s.handles.length,0);assert.equal(s.states.at(-1),'Read aloud')});
test('old response cannot replace a newer character',async()=>{const old=deferred(),fresh=deferred();const s=setup(id=>id==='old'?old.promise:fresh.promise);const a=s.voice.speak('old'),b=s.voice.speak('new');fresh.resolve(new Blob());await b;old.resolve(new Blob());await a;assert.equal(s.handles.length,1);assert.equal(s.handles[0].plays,1);assert.equal(s.states.at(-1),'Speaking…')});
test('decision releases active audio immediately',async()=>{const s=setup(async()=>new Blob());await s.voice.speak('event');s.voice.stop();assert.equal(s.handles[0].pauses,1);assert.equal(s.handles[0].disposed,1)});
test('completion releases resources and offers replay',async()=>{const s=setup(async()=>new Blob());await s.voice.speak('event');s.handles[0].end();assert.equal(s.handles[0].disposed,1);assert.equal(s.states.at(-1),'Replay voice')});
test('unavailable voice leaves retry state',async()=>{const s=setup(async()=>{throw Error('offline')});await s.voice.speak('event');assert.equal(s.states.at(-1),'Retry voice');assert.equal(s.errors(),1)});
test('the stop control appears only while there is something to stop',async()=>{
 // The scene offered "Stop voice" at all times, beside a button reading "Read
 // aloud". These are the player's own states, so this cannot drift from them.
 const s=setup(async()=>new Blob());
 assert.equal(speaking(IDLE),false,'stop offered before anything was asked for');
 const play=s.voice.speak('event');
 assert.equal(speaking(s.states.at(-1)),true,'no way to stop a reading being prepared');
 await play;
 assert.equal(speaking(s.states.at(-1)),true,'no way to stop a reading in progress');
 s.handles[0].end();
 assert.equal(speaking(s.states.at(-1)),false,'stop still offered after it finished');
 const t=setup(async()=>{throw Error('offline')});
 await t.voice.speak('event');
 assert.equal(speaking(t.states.at(-1)),false,'stop offered after the voice failed');
});
