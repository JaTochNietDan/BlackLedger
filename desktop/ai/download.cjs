const fs = require('node:fs');
const fsp = require('node:fs/promises');
const path = require('node:path');
const {createHash} = require('node:crypto');

async function digest(file, signal) {
 const hash = createHash('sha256');
 for await (const chunk of fs.createReadStream(file)) {signal?.throwIfAborted();hash.update(chunk);}
 return hash.digest('hex');
}
async function download(asset, target, progress = () => {}, signal) {
 await fsp.mkdir(path.dirname(target), {recursive:true});
 const valid = async file => (await fsp.stat(file).catch(()=>null))?.size === asset.size && await digest(file, signal) === asset.sha256;
 if (await valid(target)) { progress(asset.size, asset.size); return; }
 const partial = target + '.part';
 let offset = (await fsp.stat(partial).catch(()=>null))?.size || 0;
 if (offset > asset.size) { await fsp.unlink(partial); offset = 0; }
 if (offset === asset.size) {
  if (await valid(partial)) { await fsp.rename(partial,target); return; }
  await fsp.unlink(partial); offset = 0;
 }
 const response = await fetch(asset.url, {headers: offset ? {Range:`bytes=${offset}-`} : {}, signal});
 if (!response.ok) throw new Error(`Download failed (${response.status}). Please retry when your connection is available.`);
 if (offset && response.status !== 206) offset = 0;
 if (response.status === 206 && !response.headers.get('content-range')?.startsWith(`bytes ${offset}-`)) throw new Error('The download server returned an invalid resume position.');
 const file = await fsp.open(partial, offset ? 'a' : 'w');
 let count = offset;
 try {
  for await (const chunk of response.body) {
   signal?.throwIfAborted();
   count += chunk.length;
   if (count > asset.size) throw new Error('Downloaded file is larger than the verified release.');
   let written=0;
   while(written<chunk.length){const result=await file.write(chunk,written,chunk.length-written);if(!result.bytesWritten)throw new Error('Unable to write AI download');written+=result.bytesWritten;}
   progress(count, asset.size);
  }
 } finally { await file.close(); }
 if (count !== asset.size || await digest(partial, signal) !== asset.sha256) {
  await fsp.unlink(partial);
  throw new Error('Download verification failed. Retry to get an intact copy.');
 }
 await fsp.rename(partial, target);
}
module.exports = {download, digest};
