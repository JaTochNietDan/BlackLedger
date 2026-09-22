const test=require('node:test');const assert=require('node:assert/strict');
const http=require('node:http');const fs=require('node:fs/promises');const path=require('node:path');const os=require('node:os');const {createHash}=require('node:crypto');
const {download}=require('./ai/download.cjs');
async function fixture(t,serve){const root=await fs.mkdtemp(path.join(os.tmpdir(),'ledger-download-'));const server=http.createServer(serve);await new Promise(r=>server.listen(0,'127.0.0.1',r));t.after(async()=>{server.closeAllConnections();await new Promise(r=>server.close(r));await fs.rm(root,{recursive:true,force:true});});return {target:path.join(root,'model'),url:'http://127.0.0.1:'+server.address().port};}
const data=Buffer.from('a verified model file for restart testing');const sha256=createHash('sha256').update(data).digest('hex');
test('setup resumes a partial download and verifies its complete checksum',async t=>{
 let range;const f=await fixture(t,(req,res)=>{range=req.headers.range;res.writeHead(206,{'Content-Range':`bytes 8-${data.length-1}/${data.length}`});res.end(data.subarray(8));});
 await fs.writeFile(f.target+'.part',data.subarray(0,8));await download({...f,size:data.length,sha256},f.target);
 assert.equal(range,'bytes=8-');assert.deepEqual(await fs.readFile(f.target),data);
});
test('setup restarts when a server ignores Range and rejects a corrupt download',async t=>{
 const f=await fixture(t,(_req,res)=>res.end(data));await fs.writeFile(f.target+'.part',data.subarray(0,8));
 await download({...f,size:data.length,sha256},f.target);assert.deepEqual(await fs.readFile(f.target),data);
 await fs.rm(f.target);await assert.rejects(download({...f,size:data.length,sha256:'0'.repeat(64)},f.target),/verification failed/);
 await assert.rejects(fs.stat(f.target));await assert.rejects(fs.stat(f.target+'.part'));
});
test('cancelling setup retains completed partial data for retry',async t=>{
 const f=await fixture(t,(_req,res)=>{res.writeHead(200);res.write(data.subarray(0,8));});const controller=new AbortController();
 await assert.rejects(download({...f,size:data.length,sha256},f.target,()=>controller.abort(),controller.signal));
 assert.equal((await fs.stat(f.target+'.part')).size,8);await assert.rejects(fs.stat(f.target));
});
