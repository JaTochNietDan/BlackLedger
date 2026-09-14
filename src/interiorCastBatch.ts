import * as THREE from 'three';

type Part={mesh:THREE.Mesh;actor:THREE.Group;id:string;visible:boolean};
type Batch={mesh:THREE.InstancedMesh;parts:Part[]};

/** Draw shared character parts together while retaining the original joint rigs.
 * Geometry/textures remain owned by the loaded prototypes. Only instance buffers
 * and neutral-color material copies belong to this batch.
 */
export class InteriorCastBatch {
 readonly root=new THREE.Group();
 readonly unbatched:THREE.Mesh[]=[];
 private batches:Batch[]=[];
 private parts:Part[]=[];
 constructor(actors:ReadonlyMap<string,THREE.Group>){
  const groups=new Map<string,{material:THREE.MeshStandardMaterial;parts:Part[]}>();
  for(const [id,actor] of actors)actor.traverse(object=>{
   if(!(object instanceof THREE.Mesh))return;
   if(object instanceof THREE.SkinnedMesh||Array.isArray(object.material)||!(object.material instanceof THREE.MeshStandardMaterial)){this.unbatched.push(object);return;}
   const material=object.material;
   const key=[object.geometry.uuid,material.userData.castSource||material.uuid,object.castShadow,object.receiveShadow].join(':');
   let group=groups.get(key);if(!group){group={material,parts:[]};groups.set(key,group);}
   const part={mesh:object,actor,id,visible:object.visible};group.parts.push(part);this.parts.push(part);
  });
  for(const {material,parts} of groups.values()){
   const neutral=material.clone();neutral.color.setRGB(1,1,1);
   const mesh=new THREE.InstancedMesh(parts[0].mesh.geometry,neutral,parts.length);
   mesh.name='interior-cast-batch';mesh.castShadow=parts[0].mesh.castShadow;mesh.receiveShadow=parts[0].mesh.receiveShadow;
   mesh.instanceMatrix.setUsage(THREE.DynamicDrawUsage);
   mesh.userData.people=[];this.batches.push({mesh,parts});this.root.add(mesh);
  }
  this.update();
  for(const part of this.parts)part.mesh.visible=false;
 }
 update(){
  for(const actor of new Set(this.parts.map(p=>p.actor)))actor.updateWorldMatrix(true,true);
  for(const {mesh,parts} of this.batches){
   let count=0;const people:string[]=[];
   for(const part of parts){
    let visible=part.visible;
    for(let parent:THREE.Object3D|null=part.mesh.parent;parent;parent=parent.parent){if(!parent.visible){visible=false;break;}if(parent===part.actor)break;}
    if(!visible)continue;
    mesh.setMatrixAt(count,part.mesh.matrixWorld);mesh.setColorAt(count,(part.mesh.material as THREE.MeshStandardMaterial).color);
    people.push(part.id);count++;
   }
   mesh.count=count;mesh.visible=count>0;mesh.userData.people=people;
   mesh.instanceMatrix.needsUpdate=true;if(mesh.instanceColor)mesh.instanceColor.needsUpdate=true;
   mesh.computeBoundingBox();mesh.computeBoundingSphere();
  }
 }
 dispose(){
  this.root.removeFromParent();
  for(const {mesh} of this.batches){mesh.dispose();(mesh.material as THREE.Material).dispose();}
  for(const part of this.parts)part.mesh.visible=part.visible;
  this.root.clear();this.batches=[];this.parts=[];this.unbatched.length=0;
 }
}
