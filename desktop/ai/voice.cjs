const fs = require('node:fs/promises');
const path = require('node:path');
const {createHash} = require('node:crypto');
const profiles = require('./voice-profiles.json');
function profileFor(payload) {
 if(payload.voice === 'narrator') return {voices:['bm_george'],weight:1,speed:1};
 if(payload.voice !== 'warm') throw new Error('Unknown voice');
 let id = payload.profile;
 if(!id) {
  if(typeof payload.speaker !== 'string' || !payload.speaker.trim() || payload.speaker.length>120) throw new Error('Invalid speaker');
  const identity=payload.speaker.toLowerCase().trim().replace(/\s+/g,' ');
  id='cast2-'+String(createHash('sha256').update(identity).digest().readUInt32BE(0)%Object.keys(profiles).length).padStart(2,'0');
 }
 if(!profiles[id]) throw new Error('Unknown voice profile');
 return profiles[id];
}
function wav(samples, rate=24000) {
 const data=Buffer.alloc(44+samples.length*2);
 data.write('RIFF');data.writeUInt32LE(data.length-8,4);data.write('WAVEfmt ',8);
 data.writeUInt32LE(16,16);data.writeUInt16LE(1,20);data.writeUInt16LE(1,22);
 data.writeUInt32LE(rate,24);data.writeUInt32LE(rate*2,28);data.writeUInt16LE(2,32);data.writeUInt16LE(16,34);
 data.write('data',36);data.writeUInt32LE(samples.length*2,40);
 for(let i=0;i<samples.length;i++){
  if(!Number.isFinite(samples[i])) throw new Error('Invalid audio sample');
  data.writeInt16LE(Math.round(Math.max(-1,Math.min(1,samples[i]))*32767),44+i*2);
 }
 return data;
}
async function createVoice(root) {
 const {env,Tensor,RawAudio}=await import('@huggingface/transformers');
 env.allowRemoteModels=false;env.allowLocalModels=true;
 env.backends.onnx.wasm.numThreads=1;
 const {KokoroTTS,TextSplitterStream}=await import('kokoro-js');
 const tts=await KokoroTTS.from_pretrained(path.join(root,'voice'),{dtype:'q8',device:'cpu'});
 const voiceRoot=path.join(path.dirname(require.resolve('kokoro-js')),'../voices');
 const cache=new Map();let active;
 const load=async id=>{
  if(!cache.has(id)){
   const bytes=await fs.readFile(path.join(voiceRoot,id+'.bin'));
   cache.set(id,new Float32Array(bytes.buffer.slice(bytes.byteOffset,bytes.byteOffset+bytes.byteLength)));
  }
  return cache.get(id);
 };
 // Keep the established cast2 blends and speaking rates across platforms.
 tts.generate_from_ids=async(input_ids,{speed})=>{
  const offset=256*Math.min(Math.max(input_ids.dims.at(-1)-2,0),509);
  const a=await load(active.voices[0]),b=active.voices[1]?await load(active.voices[1]):a;
  const style=new Float32Array(256);
  for(let i=0;i<256;i++)style[i]=a[offset+i]*active.weight+b[offset+i]*(1-active.weight);
  const {waveform}=await tts.model({input_ids,style:new Tensor('float32',style,[1,256]),speed:new Tensor('float32',[speed],[1])});
  return new RawAudio(waveform.data,24000);
 };
 let busy=false;
 return async payload=>{
  if(busy){const e=new Error('Speech busy');e.status=429;throw e;}
  if(typeof payload.text!=='string'||!payload.text.trim()||payload.text.length>1000)throw new Error('Speech text must be 1–1000 characters');
  const text=payload.text.replace(/\[\[.*?\]\]/gs,'').replace(/[\[\]\x00-\x1f]/g,' ').replace(/\s+/g,' ').trim();
  active=profileFor(payload);busy=true;
  try{
   const chunks=[];let length=0;
   const input=new TextSplitterStream();input.push(text);input.close();
   for await(const {audio} of tts.stream(input,{voice:active.voices[0],speed:active.speed})){
    chunks.push(audio.audio);length+=audio.audio.length;
    if(length>24000*150)throw new Error('Speech exceeded duration limit');
   }
   if(!length)throw new Error('No speech produced');
   const samples=new Float32Array(length);let offset=0;
   for(const chunk of chunks){samples.set(chunk,offset);offset+=chunk.length;}
   return wav(samples);
  }finally{busy=false;}
 };
}
module.exports={createVoice,profileFor,wav};
