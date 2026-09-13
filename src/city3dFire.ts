import * as THREE from 'three';
import type {Snapshot} from './types';

export class CityFire {
 readonly root=new THREE.Group();
 private entries=new Map<string,{mesh:THREE.InstancedMesh;vents:THREE.Vector3[]}>();
 private geometry=new THREE.PlaneGeometry(1,1);
 private material:THREE.MeshBasicMaterial;
 private transform=new THREE.Object3D();
 private clock=0;
 constructor(texture:THREE.Texture){
  this.material=new THREE.MeshBasicMaterial({map:texture,transparent:true,opacity:.72,depthWrite:false,toneMapped:false});
 }
 update(records:NonNullable<Snapshot['building_fires']>,minute:number,buildings:Map<string,THREE.Group>,camera:THREE.Camera,dt:number,motion:boolean){
  if(motion)this.clock+=Math.min(100,Math.max(0,dt))/1000;
  const active=new Set<string>();
  for(const fire of records){
   if(minute<fire.minute||minute>=fire.extinguished_at)continue;
   active.add(fire.id);
   let entry=this.entries.get(fire.id);
   if(!entry){
    const building=buildings.get(fire.target);if(!building)continue;
    const vents:THREE.Vector3[]=[];
    building.updateMatrixWorld(true);
    const anchors:THREE.Object3D[]=[];
    building.traverse(o=>{if(o.name.startsWith('fire-window-'))anchors.push(o);});
    anchors.sort((a,b)=>Number(b.name.startsWith('fire-window-1-'))-Number(a.name.startsWith('fire-window-1-')));
    for(const anchor of anchors.slice(0,4))vents.push(anchor.getWorldPosition(new THREE.Vector3()));
    if(!vents.length)continue;
    const mesh=new THREE.InstancedMesh(this.geometry,this.material,vents.length*12);mesh.frustumCulled=false;
    this.root.add(mesh);entry={mesh,vents};this.entries.set(fire.id,entry);
   }
   let index=0;
   for(const [ventIndex,vent] of entry.vents.entries())for(let i=0;i<12;i++){
    const smoke=i>=4,cycle=smoke?3.8:.85;
    const age=(this.clock+i*.317+ventIndex*.53)%cycle;
    const life=age/cycle,fade=Math.sin(Math.PI*life);
    this.transform.position.copy(vent).add(new THREE.Vector3(
     Math.sin(i*2.4+age*1.3)*(smoke?.25+age*.15:.18),
     age*(smoke?1.4:1.6),-age*(smoke?.3:.48)));
    this.transform.quaternion.copy(camera.quaternion);
    const size=(smoke?.55+age*.48:.42+life*.8)*Math.sqrt(Math.max(.001,fade));
    this.transform.scale.set(size*(smoke?1:.65),size*(smoke?1:1.65),1);
    this.transform.updateMatrix();entry.mesh.setMatrixAt(index,this.transform.matrix);
    entry.mesh.setColorAt(index++,new THREE.Color(smoke?0x514f49:i%2?0xff751c:0xffc44d));
   }
   entry.mesh.instanceMatrix.needsUpdate=true;
   if(entry.mesh.instanceColor)entry.mesh.instanceColor.needsUpdate=true;
  }
  for(const [id,entry] of this.entries)if(!active.has(id)){this.root.remove(entry.mesh);entry.mesh.dispose();this.entries.delete(id);}
 }
 inspect(){return {clock:this.clock,scenes:[...this.entries].map(([id,e])=>({id,vents:e.vents.map(v=>({x:v.x,y:v.y,z:v.z})),particles:e.mesh.count}))};}
 dispose(){for(const e of this.entries.values())e.mesh.dispose();this.entries.clear();this.root.clear();this.geometry.dispose();this.material.dispose();}
}

/** Keep the exit corridor clear of projecting signs and facade ornament. */
export function clearBlastWindows(building:THREE.Group){
 building.updateMatrixWorld(true);
 const vents:THREE.Object3D[]=[];building.traverse(o=>{if(o.name.startsWith('fire-window-'))vents.push(o);});
 const upper=vents.filter(o=>o.name.startsWith('fire-window-1-'));
 return (upper.length?upper:vents).map(o=>o.getWorldPosition(new THREE.Vector3())).filter(p=>{
  for(const x of [-.55,0,.55]){
   const ray=new THREE.Raycaster(p.clone().add(new THREE.Vector3(x,0,-.06)),new THREE.Vector3(0,0,-1),0,3.4);
   if(ray.intersectObject(building,true).length)return false;
  }return true;
 }).slice(0,4);
}
