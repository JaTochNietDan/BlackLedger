import * as THREE from 'three';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';

// Local camera only: no world commands or document-wide keyboard capture.
export class TableCamera {
 readonly controls:OrbitControls;
 private distance:number;
 private aspectScale=1;
 constructor(private camera:THREE.PerspectiveCamera,private canvas:HTMLCanvasElement,private changed:()=>void,private center:THREE.Vector3,distance:number){
  this.distance=distance;
  this.controls=new OrbitControls(camera,canvas);
  this.controls.minDistance=1.8;this.controls.maxDistance=9;
  this.controls.minPolarAngle=.15;this.controls.maxPolarAngle=Math.PI*.38;
  this.controls.enableDamping=false;this.controls.screenSpacePanning=false;
  this.controls.addEventListener('change',this.change);
  canvas.tabIndex=0;canvas.addEventListener('keydown',this.key);canvas.addEventListener('pointerdown',this.focus);canvas.addEventListener('table-reset',this.reset);
  this.reset();
 }
 private change=()=>{const t=this.controls.target;t.x=THREE.MathUtils.clamp(t.x,-1.5,1.5);t.z=THREE.MathUtils.clamp(t.z,-1.3,1.3);t.y=this.center.y;this.changed();};
 private focus=()=>this.canvas.focus({preventScroll:true});
 private reset=()=>{this.controls.target.copy(this.center);this.aspectScale=Math.max(1,.95/this.camera.aspect);const d=this.distance*this.aspectScale;this.camera.position.copy(this.center).add(new THREE.Vector3(0,d*.84,d*.54));this.controls.update();this.changed();};
 private key=(e:KeyboardEvent)=>{
  if(e.ctrlKey||e.metaKey||e.altKey)return;
  const k=e.key.toLowerCase();
  if(!['w','a','s','d','arrowup','arrowdown','arrowleft','arrowright','q','e','+','=','-','home'].includes(k))return;
  e.preventDefault();e.stopPropagation();
  if(k==='home'){this.reset();return;}
  const offset=this.camera.position.clone().sub(this.controls.target);
  if(k==='q'||k==='e')offset.applyAxisAngle(new THREE.Vector3(0,1,0),k==='q'?.09:-.09);
  else if(['+','=','-'].includes(k))offset.multiplyScalar(k==='-'?1.08:.92).clampLength(this.controls.minDistance,this.controls.maxDistance);
  else {
   const right=new THREE.Vector3().setFromMatrixColumn(this.camera.matrixWorld,0);right.y=0;right.normalize();
   const forward=new THREE.Vector3(-right.z,0,right.x);
   const delta=['a','arrowleft','d','arrowright'].includes(k)?right.multiplyScalar(['a','arrowleft'].includes(k)?-.09:.09):forward.multiplyScalar(['w','arrowup'].includes(k)?-.09:.09);
   this.controls.target.add(delta);this.change();
  }
  this.camera.position.copy(this.controls.target).add(offset);this.controls.update();this.changed();
 };
 resize(){const scale=Math.max(1,.95/this.camera.aspect);this.camera.position.sub(this.controls.target).multiplyScalar(scale/this.aspectScale).add(this.controls.target);this.aspectScale=scale;this.controls.update();this.changed();}
 dispose(){this.controls.removeEventListener('change',this.change);this.controls.dispose();this.canvas.removeEventListener('keydown',this.key);this.canvas.removeEventListener('pointerdown',this.focus);this.canvas.removeEventListener('table-reset',this.reset);}
}
