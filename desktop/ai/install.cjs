const fs = require('node:fs');
const fsp = require('node:fs/promises');
const path = require('node:path');
const zlib = require('node:zlib');
const {pipeline} = require('node:stream/promises');
const {createHash} = require('node:crypto');
const {download} = require('./download.cjs');
const manifest = require('./manifest.json');
const manifestID = createHash('sha256').update(JSON.stringify(manifest)).digest('hex');
function runtimeAsset(platform = process.platform, arch = process.arch) {
 const asset = manifest.ollama.platforms[`${platform}-${arch}`];
 if (!asset) throw new Error('Local AI is not available for this system architecture.');
 return asset;
}
function runtimeRoot(root) { return path.join(root, 'runtime-' + manifest.ollama.version); }
function runtimeExecutable(root, platform = process.platform) {
 const runtime = runtimeRoot(root);
 return path.join(runtime, platform === 'win32' ? 'ollama.exe' : platform === 'linux' ? 'bin/ollama' : 'ollama');
}
function modelFiles(root) {
 const m = manifest.director.manifest;
 return [m.config, ...m.layers].map(layer => ({size:layer.size, sha256:layer.digest.slice(7),
  url:`https://registry.ollama.ai/v2/library/qwen3/blobs/${layer.digest}`,
  target:path.join(root, 'models/blobs', layer.digest.replace(':','-'))}));
}
function files(root) { return [...modelFiles(root), ...manifest.voice.files.map(f=>({...f,target:path.join(root,'voice',f.path)}))]; }
async function installed(root) {
 try {
  const ready = JSON.parse(await fsp.readFile(path.join(root,'ready.json'),'utf8'));
  if (ready.manifest !== manifestID || ready.platform !== `${process.platform}-${process.arch}`) return false;
  await fsp.access(runtimeExecutable(root));
  for (const f of files(root)) if ((await fsp.stat(f.target)).size !== f.size) return false;
  await fsp.access(path.join(root,'models/manifests/registry.ollama.ai/library/qwen3/14b'));
  return true;
 } catch { return false; }
}
async function install(root, report, signal) {
 const asset = runtimeAsset();
 await fsp.mkdir(root,{recursive:true});
 const all = [{...asset,target:path.join(root,'downloads',asset.archive)}, ...files(root)];
 const total = all.reduce((sum,f)=>sum+f.size,0);
 const space = await fsp.statfs(root);
 // Allow runtime extraction, partial files and model working space. Existing
 // verified downloads are reused; a retry need not reserve the whole model again.
 let remaining = 6*1024**3;
 for (const f of all) remaining += Math.max(0,f.size - ((await fsp.stat(f.target).catch(()=>null))?.size || (await fsp.stat(f.target+'.part').catch(()=>null))?.size || 0));
 if (Number(space.bavail)*Number(space.bsize) < remaining) throw new Error('Not enough free disk space for local AI. Free additional space and retry.');
 let completed = 0;
 for (const [index,f] of all.entries()) {
  signal?.throwIfAborted();
  const label=index===0?'Downloading AI runtime':index<=modelFiles(root).length?'Downloading story model':'Downloading voices';
  await download(f,f.target,(received)=>report({label,completed:completed+received,total}),signal);
  completed += f.size;
 }
 report({label:'Preparing local AI',completed:total,total});
 const destination = runtimeRoot(root), staging = destination+'.installing';
 await fsp.rm(staging,{recursive:true,force:true});await fsp.mkdir(staging,{recursive:true});
 const archive = all[0].target;
 if (asset.archive.endsWith('.zip')) {
  // Windows 10+ ships bsdtar. It bootstraps before the downloaded MSVC
  // runtime exists, unlike a native Node ZIP addon requiring that runtime.
  const {execFile}=require('node:child_process');
  await new Promise((resolve,reject)=>execFile(path.join(process.env.SystemRoot||'C:\\Windows','System32','tar.exe'),
   ['-xf',archive,'-C',staging],{windowsHide:true,signal},error=>error?reject(error):resolve()));
 }
 else {
  const links = new Map();
  const inside = name => {
   const resolved=path.resolve(staging,name);
   if(resolved!==staging&&!resolved.startsWith(staging+path.sep))throw new Error('Archive link escapes its installation folder');
   return resolved;
  };
  const unpack = require('tar').x({cwd:staging,strict:true,preservePaths:false,filter:(name,entry)=>{
   if(['SymbolicLink','Link'].includes(entry.type)){
    links.set(inside(name),{target:inside(entry.type==='Link'?entry.linkpath:path.join(path.dirname(name),entry.linkpath)),hard:entry.type==='Link'});
    return false;
   }
   return true;
  }});
  if (asset.archive.endsWith('.zst')) {
   if (!zlib.createZstdDecompress) throw new Error('The desktop runtime cannot unpack this AI release. Update Black Ledger.');
   await pipeline(fs.createReadStream(archive),zlib.createZstdDecompress(),unpack);
   } else await pipeline(fs.createReadStream(archive),zlib.createGunzip(),unpack);
  // Official dylibs contain chained links. Resolve those entirely inside the
  // verified archive, after regular files, instead of extracting through links.
  for(const [name,link] of links){
   let target=link.target;const seen=new Set([name]);
   while(links.has(target)){if(seen.has(target))throw new Error('Archive contains a link cycle');seen.add(target);target=links.get(target).target;}
   target=await fsp.realpath(target);inside(target);
   await fsp.mkdir(path.dirname(name),{recursive:true});inside(await fsp.realpath(path.dirname(name)));
   if(link.hard)await fsp.link(target,name);else await fsp.symlink(path.relative(path.dirname(name),target),name);
  }
 }
 signal?.throwIfAborted();
 await fsp.rm(destination,{recursive:true,force:true});await fsp.rename(staging,destination);
 if (process.platform!=='win32') await fsp.chmod(runtimeExecutable(root),0o755);
 const modelPath=path.join(root,'models/manifests/registry.ollama.ai/library/qwen3/14b');
 await fsp.mkdir(path.dirname(modelPath),{recursive:true});
 await fsp.writeFile(modelPath,JSON.stringify(manifest.director.manifest));
 await fsp.mkdir(path.join(root,'licenses'),{recursive:true});
 // fs.cp's native directory traversal bypasses Electron's ASAR filesystem.
 for(const name of await fsp.readdir(path.join(__dirname,'licenses')))
  await fsp.writeFile(path.join(root,'licenses',name),await fsp.readFile(path.join(__dirname,'licenses',name)));
 await fsp.writeFile(path.join(root,'ready.json'),JSON.stringify({manifest:manifestID,platform:`${process.platform}-${process.arch}`}));
 return root;
}
module.exports={install,installed,runtimeExecutable,runtimeAsset,modelFiles,manifest};
