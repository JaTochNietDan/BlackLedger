const install=document.getElementById('install'),skip=document.getElementById('skip'),status=document.getElementById('status'),progress=document.getElementById('progress');
install.addEventListener('click',()=>{install.disabled=true;skip.textContent='Cancel setup and play';window.aiSetup.action('install');});
skip.addEventListener('click',()=>{skip.disabled=true;window.aiSetup.action('skip');});
window.aiSetup.onProgress(update=>{
 status.textContent=update.label;
 if(update.done){install.disabled=true;skip.disabled=true;}
 if(update.total){progress.max=update.total;progress.value=update.completed;status.textContent+=` — ${(update.completed/1e9).toFixed(1)} / ${(update.total/1e9).toFixed(1)} GB`;}
 if(update.error){install.disabled=false;install.textContent='Retry setup';skip.disabled=false;skip.textContent='Play without AI';}
});
