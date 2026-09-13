import * as THREE from 'three';
import type {Snapshot} from './types';
import type {CityAftermath} from './city3dAftermath';

export function waterArc(from:THREE.Vector3,to:THREE.Vector3,t:number){
 const u=Math.max(0,Math.min(1,t));return from.clone().lerp(to,u).add(new THREE.Vector3(0,.65*4*u*(1-u),0));
}
type Line={root:THREE.Group;nozzle:THREE.Group;hose:THREE.Mesh;water:THREE.Points;crew:THREE.Group;from:THREE.Vector3;to:THREE.Vector3};
export class CitySuppression {
 readonly root=new THREE.Group();
 private lines=new Map<string,Line>();
 private hoseMaterial=new THREE.MeshStandardMaterial({color:0x9b8961,roughness:.95});
 private waterMaterial=new THREE.PointsMaterial({color:0xc9e7ee,size:2,sizeAttenuation:false,transparent:true,opacity:.8,depthWrite:false});
 private clock=0;
 update(fires:NonNullable<Snapshot['building_fires']>,minute:number,buildings:Map<string,THREE.Group>,models:Map<string,THREE.Group>,aftermath:CityAftermath,dt:number,motion:boolean){
  if(motion)this.clock+=Math.min(100,Math.max(0,dt))/1000;
  const desired=new Set<string>();
  for(const fire of fires){
   if(minute<fire.brigade_at||minute>=fire.extinguished_at)continue;
   const engine=aftermath.object(`aftermath:${fire.id}:fire-engine`),building=buildings.get(fire.target),source=models.get('fire-nozzle');
   if(!engine?.visible||!building||!source)continue;
   const vents:THREE.Object3D[]=[];building.traverse(o=>{if(o.name.startsWith('fire-window-'))vents.push(o);});
   const upper=vents.filter(v=>v.name.startsWith('fire-window-1-'));const targets=(upper.length?upper:vents).slice(0,4).map(v=>v.getWorldPosition(new THREE.Vector3()));
   if(!targets.length)continue;
   for(const suffix of ['a','b']){
    const id=`${fire.id}:${suffix}`,crew=aftermath.object(`aftermath:${fire.id}:firefighter-${suffix}`);
    if(!crew?.visible)continue;desired.add(id);
    let line=this.lines.get(id);
    if(!line){
     const actor=crew.children[0] as THREE.Group;
     const to=targets.reduce((a,b)=>Math.abs(b.x-crew.position.x)<Math.abs(a.x-crew.position.x)?b:a).clone();
     actor.rotation.y=Math.atan2(to.x-crew.position.x,to.z-crew.position.z);
     for(const side of [1,-1]){const arm=actor.getObjectByName(`arm${side}`);if(arm){arm.rotation.x=-1.3;arm.rotation.z=-side*.55;}}
     const forward=new THREE.Vector3(Math.sin(actor.rotation.y),0,Math.cos(actor.rotation.y));
     const grip=crew.position.clone().addScaledVector(forward,.5).add(new THREE.Vector3(0,1.21,0));
     const nozzle=source.clone(true);nozzle.position.copy(grip);
     nozzle.quaternion.setFromUnitVectors(new THREE.Vector3(0,0,1),to.clone().sub(grip).normalize());nozzle.updateMatrixWorld(true);
     const from=nozzle.getObjectByName('water-outlet')!.getWorldPosition(new THREE.Vector3());
     const side=engine.position.x>building.position.x?1:-1;
     const lane=suffix==='a'?0:.25;
     const port=engine.position.clone().add(new THREE.Vector3(side*1.125,.85,-.27-lane));
     const outer=building.position.x+side*(11.5+lane*.5);
     const behind=crew.position.z-.8-lane*.5;
     const points=[port,new THREE.Vector3(outer,.23,port.z),new THREE.Vector3(outer,.23,behind),new THREE.Vector3(crew.position.x+side*.65,.23,behind),new THREE.Vector3(crew.position.x+side*.65,.23,grip.z),grip];
     const path=new THREE.CurvePath<THREE.Vector3>();for(let i=1;i<points.length;i++)path.add(new THREE.LineCurve3(points[i-1],points[i]));
     const hose=new THREE.Mesh(new THREE.TubeGeometry(path,100,.045,6,false),this.hoseMaterial);
     const geometry=new THREE.BufferGeometry();geometry.setAttribute('position',new THREE.BufferAttribute(new Float32Array(144),3));
     const water=new THREE.Points(geometry,this.waterMaterial);water.frustumCulled=false;
     const root=new THREE.Group();root.add(nozzle,hose,water);this.root.add(root);
     line={root,nozzle,hose,water,crew,from,to};this.lines.set(id,line);
    }
    const positions=line.water.geometry.getAttribute('position') as THREE.BufferAttribute;
    for(let i=0;i<48;i++){const t=(i/48+this.clock*1.7)%1,p=waterArc(line.from,line.to,t);positions.setXYZ(i,p.x,p.y,p.z);}
    positions.needsUpdate=true;
   }
  }
  for(const [id,line] of this.lines)if(!desired.has(id))this.remove(id,line);
 }
 private remove(id:string,line:Line){
  const actor=line.crew.children[0];for(const name of ['arm1','arm-1']){const arm=actor?.getObjectByName(name);if(arm)arm.rotation.set(0,0,0);}
  this.root.remove(line.root);line.hose.geometry.dispose();line.water.geometry.dispose();this.lines.delete(id);
 }
 inspect(){return [...this.lines].map(([id,l])=>({id,from:l.from.toArray(),to:l.to.toArray()}));}
 dispose(){for(const [id,line] of this.lines)this.remove(id,line);this.hoseMaterial.dispose();this.waterMaterial.dispose();}
}
