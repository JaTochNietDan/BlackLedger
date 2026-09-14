import test from 'node:test';
import assert from 'node:assert/strict';
import {previewScene, previewScenes} from '../.runtime/frontend-test/city3dPreview.js';
import {CityCueQueue} from '../.runtime/frontend-test/city3dEvents.js';

test('every debug scene uses a private world and leaves campaign results intact',()=>{
 const state=Object.freeze({id:'campaign',life:1,revision:42,minute:600,
  last_result:Object.freeze({cues:Object.freeze([{id:'real',kind:'arrest',target:'bar'}])})});
 const before=JSON.stringify(state),queue=new CityCueQueue();
 for(const name of previewScenes){
  const preview=previewScene(state,'bar',name,name);
  assert.notEqual(preview.state.id,state.id);
  assert.equal(preview.state.revision,42);assert.equal(preview.state.minute,600);
  assert.ok(preview.state.last_result.cues.every(c=>c.target==='bar'&&c.id.startsWith('preview:')));
  assert.equal(queue.take(preview.state.id,preview.state.last_result.cues,preview.cue,true).length,
   name.startsWith('Assassination')?2:1);
 }
 assert.equal(JSON.stringify(state),before);
 assert.deepEqual(queue.take(state.id,state.last_result.cues,null,true),[],'leaving debug must not replay saved events');
});


test('assassination previews cover recorded variants with the actual weapon and victim',()=>{
 const state={id:'campaign',minute:600,last_result:{elapsed:45,cash:-100,kind:'travel'}};
 const expected=[['back-of-head',1],['close-shot',1],['close-shot',2],['burst',3],['close-quarters',0]];
 const scenarios=previewScenes.filter(name=>name.startsWith('Assassination'));
 assert.equal(scenarios.length,expected.length);
 scenarios.forEach((name,i)=>{
  const {state:preview}=previewScene(state,'bar',name,String(i));
  const [victim,attack]=preview.last_result.cues;
  assert.equal(attack.strike.variant,expected[i][0]);
  assert.equal(attack.attacker.weapon,expected[i][1]);
  assert.equal(attack.kind,expected[i][1]?'gunfight':'attack');
  assert.equal(attack.strike.victim.id,victim.actors[0].id);
  assert.deepEqual(victim.strike,attack.strike);
  assert.equal(preview.last_result.elapsed,0);
  assert.equal(preview.last_result.cash,0);
  assert.equal(preview.last_result.kind,undefined,'preview must not inherit a campaign action');
 });
});
