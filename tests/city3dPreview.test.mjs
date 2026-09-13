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
   name==='Assassination'?2:1);
 }
 assert.equal(JSON.stringify(state),before);
 assert.deepEqual(queue.take(state.id,state.last_result.cues,null,true),[],'leaving debug must not replay saved events');
});
