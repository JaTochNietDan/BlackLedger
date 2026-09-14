import * as THREE from 'three';
import {type Card,knownCard,pipOf,isRedSuit} from './cards';

export function cardTexture(card?:Card){
 const canvas=document.createElement('canvas');canvas.width=256;canvas.height=368;
 const c=canvas.getContext('2d')!;
 c.fillStyle='#eee6cf';c.fillRect(0,0,256,368);
 if(!card||!knownCard(card)){
  c.fillStyle='#702920';c.fillRect(12,12,232,344);c.strokeStyle='#d1ad77';c.lineWidth=1;
  for(let y=-240;y<600;y+=14){c.beginPath();c.moveTo(18,y);c.lineTo(238,y+220);c.stroke();c.beginPath();c.moveTo(238,y);c.lineTo(18,y+220);c.stroke();}
  c.strokeStyle='#eee6cf';c.lineWidth=5;c.strokeRect(20,20,216,328);
 }else{
  c.fillStyle=isRedSuit(card.suit)?'#a32922':'#1f231f';
  for(let i=0;i<2;i++){c.save();if(i){c.translate(256,368);c.rotate(Math.PI);}c.font='bold 60px Georgia';c.fillText(card.rank,14,63);c.font='48px Georgia';c.fillText(pipOf(card.suit),16,110);c.restore();}
  c.textAlign='center';
  const n=Number(card.rank);
  const rows:Record<number,number[][]>={2:[[0,-1],[0,1]],3:[[0,-1],[0,0],[0,1]],4:[[-1,-1],[1,-1],[-1,1],[1,1]],5:[[-1,-1],[1,-1],[0,0],[-1,1],[1,1]],6:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1]],7:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1],[0,-.5]],8:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1],[0,-.5],[0,.5]],9:[[-1,-1],[1,-1],[-1,-.33],[1,-.33],[-1,.33],[1,.33],[-1,1],[1,1],[0,0]],10:[[-1,-1],[1,-1],[-1,-.33],[1,-.33],[-1,.33],[1,.33],[-1,1],[1,1],[0,-.67],[0,.67]]};
  if(rows[n]){c.font='43px Georgia';for(const [x,y] of rows[n])c.fillText(pipOf(card.suit),128+x*39,195+y*87);}
  else {c.font='bold 76px Georgia';c.fillText(card.rank,128,176);c.font='64px Georgia';c.fillText(pipOf(card.suit),128,246);}
 }
 const texture=new THREE.CanvasTexture(canvas);texture.colorSpace=THREE.SRGBColorSpace;texture.flipY=false;texture.anisotropy=4;return texture;
}

