import test from 'node:test';
import assert from 'node:assert/strict';
import {previewScene, previewScenes} from '../.runtime/frontend-test/city3dPreview.js';

test('building drive-by preview captures distinct cast and equipment without damaging campaign',()=>{
 const state={id:'qa',minute:600,player:{name:'Alex'},locations:[{id:'bar',condition:100}]};
 const before=JSON.stringify(state),preview=previewScene(state,'bar','Building drive-by','drive');
 assert.equal(JSON.stringify(state),before);
 assert.equal(preview.cue.kind,'driveby-building');
 assert.equal(preview.cue.attacker.weapon,3);
 assert.notEqual(preview.cue.attacker.id,preview.cue.drive_by.driver.id);
 assert.equal(preview.cue.drive_by.vehicle_tier,3);
 assert.equal(preview.state.locations[0].condition,100);
 assert.equal(preview.state.last_result.elapsed,0);
 const damaged={...state,locations:[{id:'bar',condition:43}]};
 const replay=previewScene(damaged,'bar','Building drive-by','damaged');
 assert.equal(replay.cue.drive_by.condition_before,43);
 assert.equal(replay.cue.drive_by.condition_after,15);
 assert.equal(damaged.locations[0].condition,43);
});
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
   name.startsWith('Assassination')||name==='Explosion · casualty'?2:1);
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

test('ordinary gunfire previews use each public gun tier without inventing a fatality',()=>{
 const state=Object.freeze({id:'campaign',minute:600,revision:10});
 for(const [name,weapon] of [['Gunfight',1],['Gunfight · shotgun',2],['Gunfight · Thompson',3]]){
  assert.ok(previewScenes.includes(name));
  const {state:preview}=previewScene(state,'laundry',name,name);
  const cues=preview.last_result.cues;assert.equal(cues.length,1);
  assert.equal(cues[0].kind,weapon?'gunfight':'attack');assert.equal(cues[0].attacker.weapon,weapon);
  assert.equal(cues[0].strike,undefined);assert.deepEqual(cues[0].actors,[]);
  assert.equal(preview.revision,state.revision);assert.equal(preview.minute,state.minute);
 }
});

test('incendiary preview uses an unarmed attacker and an isolated target fire',()=>{
 const existing={id:'other',target:'bar',minute:500,brigade_at:510,extinguished_at:680,cleanup_at:740};
 const state={id:'campaign',revision:10,minute:600,building_fires:[existing]};
 const before=JSON.stringify(state),{state:preview,cue}=previewScene(state,'club','Incendiary','fire');
 assert.equal(cue.kind,'incendiary');assert.equal(cue.attacker.weapon,0);assert.deepEqual(cue.actors,[]);
 assert.equal(preview.building_fires.length,2);assert.equal(preview.building_fires[0],existing);
 assert.equal(preview.building_fires[1].target,'club');assert.ok(preview.building_fires[1].brigade_at>state.minute);
 assert.equal(JSON.stringify(state),before);
});

test('explicit demolition outcome overrides an unrelated same-minute fire',async()=>{
 const {internalDetonation}=await import('../.runtime/frontend-test/city3dBlast.js');
 const cue={id:'accident',target:'club',minute:600,detonation:'premature'};
 assert.equal(internalDetonation(cue,[{target:'club',minute:600}]),false);
 assert.equal(internalDetonation({...cue,id:'preview:accident'},[]),false);
 assert.equal(internalDetonation({...cue,detonation:'planted'},[]),true,'recorded planted charge lost its origin after fire cleanup');
 assert.equal(internalDetonation({id:'old',target:'club',minute:600},[{target:'club',minute:600}]),true);
});

test('premature explosion preview retains its accident outcome without inventing a fire',()=>{
 const state={id:'campaign',minute:600};
 const {state:preview,cue}=previewScene(state,'club','Explosion · premature','accident');
 assert.equal(cue.kind,'explosion');assert.equal(cue.detonation,'premature');
 assert.deepEqual(cue.accident,{health_lost:30,fatal:false});
 const fatal=previewScene(state,'club','Explosion · fatal accident','fatal');
 assert.equal(fatal.cue.detonation,'premature');assert.equal(fatal.cue.attacker.id,'preview-planter');
 assert.deepEqual(fatal.cue.accident,{health_lost:40,fatal:true});assert.equal(fatal.state.building_fires,undefined);
 assert.equal(state.accident,undefined);
 assert.equal(cue.attacker.id,'preview-planter');assert.equal(preview.building_fires,undefined);
 assert.equal(previewScene(state,'club','Explosion','planted').cue.detonation,'planted');
});

test('combined explosion preview pairs the casualty with a recorded player planter',()=>{
 const state={id:'campaign',minute:600,player:{name:'Alex Varga'}};
 const {state:preview}=previewScene(state,'club','Explosion · casualty','paired');
 const [victim,blast]=preview.last_result.cues;
 assert.equal(victim.kind,'killing');assert.equal(victim.actors[0].id,'preview-victim');assert.equal(victim.attacker,undefined);
 assert.equal(blast.kind,'explosion');assert.equal(blast.detonation,'planted');assert.equal(blast.attacker.id,'player');assert.equal(blast.attacker.name,'Alex Varga');
 assert.equal(victim.target,blast.target);assert.equal(victim.minute,blast.minute);assert.equal(preview.building_fires.length,1);
});
