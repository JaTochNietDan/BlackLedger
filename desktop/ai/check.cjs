// Run under the packaged Electron Node runtime on each native CI runner.
(async()=>{
 const {phonemize}=await import('phonemizer');
 const phones=await phonemize('The city remembers.','en-us');
 if(!phones?.length)throw new Error('Speech phonemizer produced no output');
 const ort=require('onnxruntime-node');
 if(!ort.InferenceSession)throw new Error('ONNX native runtime unavailable');
 await import('kokoro-js');
 if(!require('node:zlib').createZstdDecompress)throw new Error('AI archive decompression unavailable');
 if(process.argv[2]){
  const root=process.argv[2],path=require('node:path'),fs=require('node:fs/promises');
  const {download}=require('./download.cjs'),{voice}=require('./manifest.json');
  const signal=AbortSignal.timeout(240000);
  for(const file of voice.files)await download(file,path.join(root,'voice',file.path),undefined,signal);
  const render=await require('./voice.cjs').createVoice(root);
  const audio=await render({voice:'narrator',text:'The city remembers every promise.'});
  if(audio.length<=44||audio.toString('ascii',0,4)!=='RIFF')throw new Error('Packaged speech generation failed');
  await fs.writeFile(path.join(root,'smoke-voice.wav'),audio);
  console.log('AI speech generation ready');
 }
 console.log('AI native dependencies ready');
})().catch(error=>{console.error(error);process.exitCode=1;});
