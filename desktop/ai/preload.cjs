const {contextBridge,ipcRenderer}=require('electron');
contextBridge.exposeInMainWorld('aiSetup',{
 action:value=>{if(['install','skip'].includes(value))ipcRenderer.send('ai-setup-action',value);},
 onProgress:callback=>ipcRenderer.on('ai-setup-progress',(_event,value)=>callback(value)),
});
