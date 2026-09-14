import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {TableCamera} from '../.runtime/frontend-test/tableCamera.js';

test('table camera moves between key repeats and releases on keyup, focus loss and disposal',()=>{
 const old={window:globalThis.window,document:globalThis.document,requestAnimationFrame:globalThis.requestAnimationFrame,cancelAnimationFrame:globalThis.cancelAnimationFrame};
 const doc=new EventTarget(),win=new EventTarget(),canvas=new EventTarget();doc.hidden=false;
 Object.assign(canvas,{style:{},ownerDocument:doc,getRootNode:()=>doc,focus(){},clientWidth:800,clientHeight:600});
 let id=0,now=performance.now();const frames=new Map();
 Object.assign(globalThis,{window:win,document:doc,requestAnimationFrame:fn=>{frames.set(++id,fn);return id;},cancelAnimationFrame:id=>frames.delete(id)});
 const step=()=>{now+=20;const pending=[...frames.values()];frames.clear();pending.forEach(fn=>fn(now));};
 const key=(target,type,key)=>{const e=new Event(type,{cancelable:true});Object.defineProperty(e,'key',{value:key});target.dispatchEvent(e);};
 let view;
 try{
  const camera=new THREE.PerspectiveCamera(38,4/3,.01,30);view=new TableCamera(camera,canvas,()=>{},new THREE.Vector3(0,1,0),4.6);
  key(canvas,'keydown','d');step();const first=view.controls.target.x;step();assert.ok(view.controls.target.x>first,'held key must move without another keydown');
  key(win,'keyup','d');const stopped=camera.position.clone();step();assert.ok(camera.position.distanceTo(stopped)<1e-10);
  key(canvas,'keydown','q');step();const orbit=camera.position.clone();step();assert.ok(camera.position.distanceTo(orbit)>.01);
  canvas.dispatchEvent(new Event('blur'));const blurred=camera.position.clone();step();assert.ok(camera.position.distanceTo(blurred)<1e-10);
  key(canvas,'keydown','w');key(canvas,'keydown','Home');const reset=camera.position.clone();step();assert.ok(camera.position.distanceTo(reset)<1e-10);
  const oldDistance=camera.position.distanceTo(view.controls.target);
  key(canvas,'keydown','d');view.frameView(new THREE.Vector3(0,.924,.20),3.35);
  const framed=camera.position.clone();step();assert.ok(camera.position.distanceTo(framed)<1e-10,'framing must clear held input');
  assert.ok(camera.position.distanceTo(view.controls.target)<oldDistance);assert.equal(view.controls.target.z,.20);
  key(canvas,'keydown','Home');assert.ok(camera.position.distanceTo(framed)<1e-10,'Home preserves selected framing');
  view.dispose();view=undefined;assert.equal(frames.size,0,'unmounted table must cancel its animation loop');
 }finally{view?.dispose();Object.assign(globalThis,old);}
});
