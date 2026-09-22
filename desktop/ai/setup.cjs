const {app,BrowserWindow,ipcMain}=require('electron');
const path=require('node:path');
const fs=require('node:fs/promises');
const {install,installed}=require('./install.cjs');
async function chooseAI(root,{force=false}={}){
 const preference=path.join(root,'preference.json');
 let prior={};try{prior=JSON.parse(await fs.readFile(preference,'utf8'));}catch{}
 if(!force&&prior.enabled===false)return false;
 if(!force&&await installed(root))return true;
 return new Promise((resolve,reject)=>{
  const window=new BrowserWindow({title:'Black Ledger — Local AI',width:780,height:650,resizable:false,
   backgroundColor:'#17140f',webPreferences:{preload:path.join(__dirname,'preload.cjs'),contextIsolation:true,nodeIntegration:false,sandbox:true}});
  const focus=()=>{if(!window.isDestroyed()){if(window.isMinimized())window.restore();window.show();window.focus();}};
  app.on('second-instance',focus);
  window.setMenu(null);window.webContents.setWindowOpenHandler(()=>({action:'deny'}));
  window.webContents.on('will-navigate',event=>event.preventDefault());
  let controller,task,settled=false,last=0;
  const report=update=>{if(!window.isDestroyed()&&(update.error||update.done||Date.now()-last>150)){last=Date.now();window.webContents.send('ai-setup-progress',update);}};
  const finish=async enabled=>{
   if(settled)return;
   await fs.mkdir(root,{recursive:true});await fs.writeFile(preference,JSON.stringify({enabled}));settled=true;
   // Keep this window alive until the game has created its replacement.
   report({label:enabled?'Starting local AI…':'Opening the city…',done:true});ipcMain.removeListener('ai-setup-action',action);
   resolve({enabled,close:()=>{if(!window.isDestroyed())window.close();}});
  };
  const handleAction=async(event,value)=>{
   if(event.sender!==window.webContents||settled)return;
   if(value==='skip'){
    controller?.abort();await task?.catch(()=>{});await finish(false);return;
   }
   if(value!=='install'||task)return;
   controller=new AbortController();
   task=install(root,report,controller.signal);
   try{await task;await finish(true);}catch(error){if(!controller.signal.aborted)report({label:error.message,error:true});}finally{task=null;}
  };
  const action=(event,value)=>{handleAction(event,value).catch(error=>report({label:error.message,error:true}));};
  ipcMain.on('ai-setup-action',action);
  window.on('closed',()=>{app.removeListener('second-instance',focus);controller?.abort();ipcMain.removeListener('ai-setup-action',action);if(!settled)reject(new Error('AI setup closed'));});
  window.loadFile(path.join(__dirname,'setup.html')).catch(reject);
 });
}
module.exports={chooseAI};
