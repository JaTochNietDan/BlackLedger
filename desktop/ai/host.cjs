// A child Node runtime supplied by Electron; no system Python/Node/Ollama needed.
const {spawn,spawnSync}=require('node:child_process');
const http=require('node:http');
const net=require('node:net');
const path=require('node:path');
const {setTimeout:delay}=require('node:timers/promises');
const {runtimeExecutable,manifest}=require('./install.cjs');
const {createVoice}=require('./voice.cjs');
const root=process.argv[2];
// The verified Ollama Windows archive supplies app-local MSVC runtime DLLs.
// Set the DLL search path before loading ONNX; do not require a system install.
if(process.platform==='win32'){
 const runtime=path.dirname(runtimeExecutable(root));
 process.env.PATH=path.join(runtime,'lib','ollama','cuda_v12')+path.delimiter+(process.env.PATH||'');
}
let ollama,server,stopping=false;
async function stop(code=0){
 if(stopping)return;stopping=true;
 if(server)server.close();
 if(ollama?.pid){
  if(process.platform==='win32')spawnSync('taskkill',['/PID',String(ollama.pid),'/T','/F'],{windowsHide:true,stdio:'ignore'});
  else {try{process.kill(-ollama.pid,'SIGTERM');}catch{} }
 }
 setTimeout(()=>{
  if(ollama?.pid&&process.platform!=='win32'){try{process.kill(-ollama.pid,'SIGKILL');}catch{}}
  process.exit(code);
 },1200);
}
process.stdin.resume();process.stdin.on('end',()=>stop());
process.on('SIGTERM',()=>stop());process.on('SIGINT',()=>stop());
process.on('uncaughtException',error=>{console.error(error);stop(1);});
async function freePort(){const s=net.createServer();await new Promise((resolve,reject)=>{s.once('error',reject);s.listen(0,'127.0.0.1',resolve);});const port=s.address().port;await new Promise(resolve=>s.close(resolve));return port;}
async function main(){
 const port=await freePort(),origin='http://127.0.0.1:'+port;
 if(stopping)return;
 ollama=spawn(runtimeExecutable(root),['serve'],{cwd:path.dirname(runtimeExecutable(root)),windowsHide:true,detached:process.platform!=='win32',
  env:{...process.env,OLLAMA_HOST:'127.0.0.1:'+port,OLLAMA_MODELS:path.join(root,'models'),OLLAMA_NO_CLOUD:'1'},stdio:['ignore','ignore','pipe']});
 ollama.stderr.on('data',chunk=>process.stderr.write(chunk));
 let exited=false;ollama.on('exit',()=>{exited=true;if(!stopping){console.error('AI runtime stopped');stop(1);}});ollama.on('error',error=>{console.error(error);stop(1);});
 let ready=false;
 for(let i=0;i<300&&!stopping;i++){
  if(exited)throw new Error('AI runtime did not start');
  try{const r=await fetch(origin+'/api/tags',{signal:AbortSignal.timeout(1000)});const data=await r.json();if(data.models?.some(m=>m.name===manifest.director.name)){ready=true;break;}}catch{}
  await delay(100);
 }
 if(!ready)throw new Error('Installed story model could not be opened');
 const render=await createVoice(root);
 if(stopping)return;
 server=http.createServer(async(req,res)=>{
  const send=(status,body,type='application/json')=>{res.writeHead(status,{'Content-Type':type});res.end(type==='application/json'?JSON.stringify(body):body);};
  // Only the desktop/Go clients call this service, never a web page directly.
  if(req.headers.origin)return send(403,{error:'Browser origins cannot call local AI directly'});
  if(req.method==='GET'&&req.url==='/api/health')return send(200,{ok:true,game:'Black Ledger',ai:true,model:manifest.director.name});
  if(req.method!=='POST'||!['/speech','/api/chat'].includes(req.url))return send(404,{error:'Not found'});
  try{
   const chunks=[];let bytes=0;for await(const chunk of req){bytes+=chunk.length;if(bytes>2*1024*1024)throw new Error('Request too large');chunks.push(chunk);}
   const body=Buffer.concat(chunks).toString('utf8'),payload=JSON.parse(body);
   if(req.url==='/speech')return send(200,await render(payload),'audio/wav');
   if(payload.model!==manifest.director.name)throw new Error('Unknown installed model');
   const upstream=await fetch(origin+'/api/chat',{method:'POST',headers:{'Content-Type':'application/json'},body,signal:AbortSignal.timeout(190000)});
   res.writeHead(upstream.status,{'Content-Type':'application/json'});res.end(await upstream.text());
  }catch(error){send(error.status||503,{error:error.message});}
 });
 await new Promise(resolve=>server.listen(0,'127.0.0.1',resolve));
 console.log('Local AI ready at http://127.0.0.1:'+server.address().port);
}
main().catch(error=>{console.error(error);stop(1);});
