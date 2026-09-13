import * as THREE from 'three';
import type {Snapshot} from './types';
import {windowDebris} from './city3dBlast.js';
import {surfaceHeight} from './city3dPlan.js';

/** Saved fire cleanup owns rubble lifetime; scene playback temporarily owns its fragments. */
export class CityRubble {
 readonly root=new THREE.Group();
 private entries=new Map<string,THREE.InstancedMesh>();
 private geometry?:THREE.BufferGeometry;
 private material?:THREE.Material|THREE.Material[];
 update(records:NonNullable<Snapshot['building_fires']>,minute:number,buildings:Map<string,THREE.Group>,source:THREE.Group|undefined,playing:Set<string>){
  if(!this.geometry&&source){
   source.updateMatrixWorld(true);
   source.traverse(o=>{if(o instanceof THREE.Mesh&&!this.geometry){this.geometry=o.geometry.clone().applyMatrix4(o.matrixWorld);this.material=o.material;}});
  }
  const active=new Set<string>();
  for(const fire of records){
   if(minute<fire.minute||minute>=fire.cleanup_at)continue;
   const windows=buildings.get(fire.target)?.userData.blastWindows as THREE.Vector3[]|undefined;
   if(!windows?.length||!this.geometry||!this.material)continue;
   active.add(fire.id);
   let mesh=this.entries.get(fire.id);
   if(!mesh){
    mesh=new THREE.InstancedMesh(this.geometry,this.material,12);mesh.frustumCulled=false;
    const transform=new THREE.Object3D();
    for(let i=0;i<12;i++){
     const window=windows[i%windows.length];
     const p=windowDebris(i,2.5,window,surfaceHeight({x:window.x,z:window.z-3}),windows.length);
     transform.rotation.set(p.rx,p.ry,p.rz);transform.scale.setScalar(p.scale);
     transform.position.set(p.x,surfaceHeight({x:p.x,z:p.z})+.04*p.scale+.006,p.z);
     transform.updateMatrix();mesh.setMatrixAt(i,transform.matrix);
    }
    mesh.instanceMatrix.needsUpdate=true;this.root.add(mesh);this.entries.set(fire.id,mesh);
   }
   mesh.visible=!playing.has(fire.target);
  }
  for(const [id,mesh] of this.entries)if(!active.has(id)){this.root.remove(mesh);mesh.dispose();this.entries.delete(id);}
 }
 inspect(){return [...this.entries].map(([id,mesh])=>({id,visible:mesh.visible,fragments:mesh.count}));}
 dispose(){for(const mesh of this.entries.values())mesh.dispose();this.entries.clear();this.root.clear();this.geometry?.dispose();this.geometry=undefined;}
}
