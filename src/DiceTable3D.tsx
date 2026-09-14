import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {PresentedDie,DICE_ROLL_MS} from './dicePresentation';
import {disposeCityResources} from './city3dResources';
import './diceTable3d.css';
export function DiceTable3D(props:{faces:number[];rolling:boolean;turn:number}){
 const host=useRef<HTMLDivElement>(null),latest=useRef(props);latest.current=props;
 const [status,setStatus]=useState('Opening the dice tray…');
 useEffect(()=>{
  const el=host.current!;let dead=false,frame=0,dirty=true,seen='',start=0,wasRolling=false,rendered=0;
  const scene=new THREE.Scene();scene.background=new THREE.Color(0x173128);
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;renderer.toneMappingExposure=1.2;
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','3D dice tray. The two upper faces show the committed roll once the dice settle.');el.append(canvas);
  const camera=new THREE.PerspectiveCamera(36,1,.01,20);camera.position.set(1.45,2.65,2.7);camera.lookAt(0,0,0);
  scene.add(new THREE.HemisphereLight(0xffedcd,0x20352a,2.6));
  const lamp=new THREE.DirectionalLight(0xffe4bc,3);lamp.position.set(-1,4,2);lamp.castShadow=true;lamp.shadow.mapSize.set(1024,1024);Object.assign(lamp.shadow.camera,{left:-2,right:2,top:2,bottom:-2,near:.1,far:10});lamp.shadow.bias=-.0001;scene.add(lamp);
  const resize=()=>{renderer.setSize(el.clientWidth,el.clientHeight);camera.aspect=el.clientWidth/el.clientHeight;camera.updateProjectionMatrix();dirty=true;};const observer=new ResizeObserver(resize);observer.observe(el);resize();
  const models:THREE.Group[]=[];let dice:PresentedDie[]=[];
  const loader=new GLTFLoader();Promise.all(['dice-tray','gaming-die'].map(async name=>{const g=await loader.loadAsync(`/art/models/${name}.glb`);if(dead){disposeCityResources([g.scene]);return undefined;}models.push(g.scene);return g.scene;})).then(([tray,die])=>{
   if(dead||!tray||!die)return;tray.traverse(o=>{if(o instanceof THREE.Mesh)o.receiveShadow=true;});scene.add(tray);
   dice=[0,1].map(i=>{const object=die.clone(true);object.traverse(o=>{if(o instanceof THREE.Mesh)o.castShadow=true;});scene.add(object);return new PresentedDie(object,i);});dirty=true;setStatus('');
  }).catch(()=>{if(!dead)setStatus('The dice tray could not load. The committed roll remains available in the result below.');});
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current;
   if(p.rolling&&!wasRolling)start=now;wasRolling=p.rolling;
   const key=JSON.stringify([p.faces,p.rolling,p.turn]);if(key!==seen){seen=key;dirty=true;}
   if(p.rolling)dirty=true;if(!dirty||document.hidden)return;
   const progress=p.rolling?Math.min(1,(now-start)/DICE_ROLL_MS):1;
   dice.forEach((die,i)=>die.pose(p.faces[i],progress));renderer.render(scene,camera);rendered++;
   canvas.dataset.dice=JSON.stringify({faces:p.faces,rolling:p.rolling,progress,rendered,drawCalls:renderer.info.render.calls,triangles:renderer.info.render.triangles});dirty=false;
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();disposeCityResources([scene,...models]);renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[]);
 return <div className="dice3d"><div ref={host}/>{status&&<p role="status">{status}</p>}{!props.rolling&&props.faces.length===2&&<span className="dice3d-caption">{props.faces.join(' + ')} = {props.faces.reduce((a,b)=>a+b,0)}</span>}</div>;
}
