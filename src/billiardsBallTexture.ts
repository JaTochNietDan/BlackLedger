import * as THREE from 'three';
const colours=['#f2ead8','#e6ad23','#2248a0','#b7312b','#613774','#c66322','#236548','#761e26','#151615'];
export function ballTexture(id:number){
 const canvas=document.createElement('canvas');canvas.width=512;canvas.height=256;const c=canvas.getContext('2d')!;
 c.fillStyle=id>8?'#f2ead8':colours[id];c.fillRect(0,0,512,256);
 if(id>8){c.fillStyle=colours[id-8];c.fillRect(0,77,512,102);}
 if(id)for(const x of [128,384]){c.fillStyle='#f7efda';c.beginPath();c.arc(x,128,36,0,Math.PI*2);c.fill();c.fillStyle='#151515';c.font='bold 49px Georgia';c.textAlign='center';c.textBaseline='middle';c.fillText(String(id),x,130);}
 const t=new THREE.CanvasTexture(canvas);t.colorSpace=THREE.SRGBColorSpace;return t;
}
