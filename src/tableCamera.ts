import {KeyboardPan,bindKeyboardPan,cameraCommand} from './city3dControls.js';
import * as THREE from 'three';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';

// Local camera only: no world commands or document-wide keyboard capture.
export class TableCamera {
 readonly controls:OrbitControls;
 private distance:number;
 private aspectScale=1;
 private held=new KeyboardPan();
 private unbind:()=>void;
 private frame=0;
 private lastTime=0;
 constructor(private camera:THREE.PerspectiveCamera,private canvas:HTMLCanvasElement,private changed:()=>void,private center:THREE.Vector3,distance:number,private azimuth=0,private fitAspect=.95){
  this.distance=distance;
  this.controls=new OrbitControls(camera,canvas);
  this.controls.minDistance=1.8;this.controls.maxDistance=9;
  this.controls.minPolarAngle=.15;this.controls.maxPolarAngle=Math.PI*.38;
  this.controls.enableDamping=false;this.controls.screenSpacePanning=false;
  this.controls.addEventListener('change',this.change);
  canvas.tabIndex=0;canvas.addEventListener('keydown',this.key);canvas.addEventListener('pointerdown',this.focus);canvas.addEventListener('table-reset',this.reset);
  this.unbind=bindKeyboardPan(canvas,this.held);
  this.reset();this.lastTime=performance.now();this.frame=requestAnimationFrame(this.tick);
 }
 private change=()=>{const t=this.controls.target;t.x=THREE.MathUtils.clamp(t.x,-1.5,1.5);t.z=THREE.MathUtils.clamp(t.z,-1.3,1.3);t.y=this.center.y;this.changed();};
 private focus=()=>this.canvas.focus({preventScroll:true});
 private reset=()=>{this.held.clear();this.controls.target.copy(this.center);this.aspectScale=Math.max(1,this.fitAspect/this.camera.aspect);const d=this.distance*this.aspectScale;this.camera.position.copy(this.center).add(new THREE.Vector3(0,d*.84,d*.54).applyAxisAngle(new THREE.Vector3(0,1,0),this.azimuth));this.controls.update();this.changed();};
 private tick=(now:number)=>{
  const seconds=(now-this.lastTime)/1000;this.lastTime=now;
  if(!document.hidden){
   const delta=this.held.step(this.camera.position,this.controls.target,seconds,1.5),turn=this.held.rotation(seconds);
   if(delta.x||delta.z||turn){
    const offset=this.camera.position.clone().sub(this.controls.target).applyAxisAngle(new THREE.Vector3(0,1,0),turn);
    this.controls.target.x=THREE.MathUtils.clamp(this.controls.target.x+delta.x,-1.5,1.5);
    this.controls.target.z=THREE.MathUtils.clamp(this.controls.target.z+delta.z,-1.3,1.3);
    this.camera.position.copy(this.controls.target).add(offset);this.controls.update();this.changed();
   }
  }
  this.frame=requestAnimationFrame(this.tick);
 };
 private key=(e:KeyboardEvent)=>{
  const command=cameraCommand(e);if(!command)return;
  e.preventDefault();e.stopPropagation();
  if(command==='reset'){this.reset();return;}
  if(this.held.press(e))return;
  const offset=this.camera.position.clone().sub(this.controls.target);
  offset.multiplyScalar(command==='zoom-out'?1.08:.92).clampLength(this.controls.minDistance,this.controls.maxDistance);
  this.camera.position.copy(this.controls.target).add(offset);this.controls.update();this.changed();
 };
 frameView(center:THREE.Vector3,distance:number,azimuth=this.azimuth){this.center.copy(center);this.distance=distance;this.azimuth=azimuth;this.reset();}
 resize(){const scale=Math.max(1,this.fitAspect/this.camera.aspect);this.camera.position.sub(this.controls.target).multiplyScalar(scale/this.aspectScale).add(this.controls.target);this.aspectScale=scale;this.controls.update();this.changed();}
 dispose(){cancelAnimationFrame(this.frame);this.unbind();this.controls.removeEventListener('change',this.change);this.controls.dispose();this.canvas.removeEventListener('keydown',this.key);this.canvas.removeEventListener('pointerdown',this.focus);this.canvas.removeEventListener('table-reset',this.reset);}
}
