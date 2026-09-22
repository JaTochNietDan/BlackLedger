const test=require('node:test');const assert=require('node:assert/strict');
const {profileFor,wav}=require('./ai/voice.cjs');
const profiles=require('./ai/voice-profiles.json');
test('saved voice profiles keep their blends, rates and narrator',()=>{
 assert.equal(Object.keys(profiles).length,48);
 assert.deepEqual(profileFor({voice:'warm',profile:'cast2-01'}),{voices:['af_alloy','af_aoede'],weight:.7,speed:.96,language:'a'});
 assert.deepEqual(profileFor({voice:'narrator'}).voices,['bm_george']);
 assert.deepEqual(profileFor({voice:'warm',speaker:' Mara   Vale '}),profileFor({voice:'warm',speaker:'mara vale'}));
 assert.throws(()=>profileFor({voice:'warm',profile:'../../anything'}));
});
test('voice output is bounded mono PCM WAV with finite samples',()=>{
 const bytes=wav(new Float32Array([0,1,-1,.5]));assert.equal(bytes.toString('ascii',0,4),'RIFF');assert.equal(bytes.readUInt32LE(24),24000);assert.equal(bytes.length,52);
 assert.equal(bytes.readInt16LE(46),32767);assert.equal(bytes.readInt16LE(48),-32767);assert.throws(()=>wav([NaN]));
});
