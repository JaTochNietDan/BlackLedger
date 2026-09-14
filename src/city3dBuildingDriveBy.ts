import * as THREE from 'three';
import {seatDriver} from './city3dSeating.js';
import {aimArm,poseLongGun} from './city3dWeapons.js';

export const BUILDING_DRIVEBY_SECONDS=6.2;
export const buildingDriveByShots=(weapon:string):readonly number[]=>weapon==='thompson'
 ? [2.4,2.49,2.58,2.67,3.2,3.29,3.38,3.47]
 : weapon==='shotgun'?[2.4,3.6]:[2.4,2.8,3.2,3.6];
/** Reveal only the recorded total damage, distributed across the firing beats.
 * This is a visual interpolation, never an additional gameplay damage roll. */
export function buildingDriveByCondition(before:number,after:number,seconds:number,weapon:string){
 const shots=buildingDriveByShots(weapon),landed=shots.filter(at=>seconds>=at).length;
 return before+Math.round((after-before)*landed/shots.length);
}
const ease=(v:number)=>{const t=Math.max(0,Math.min(1,v));return t*t*(3-2*t);};

/** First visible authored surface along the shot, not the building's box.
 * Hidden intact/broken variants and invisible materials cannot catch bullets. */
export function buildingDriveByHit(building:THREE.Object3D,origin:THREE.Vector3,toward:THREE.Vector3){
 building.updateWorldMatrix(true,true);
 const delta=toward.clone().sub(origin),distance=delta.length();
 if(distance<.001)return null;
 const meshes:THREE.Mesh[]=[];
 building.traverseVisible(o=>{if(o instanceof THREE.Mesh)meshes.push(o);});
 const ray=new THREE.Raycaster(origin,delta.divideScalar(distance),.001,distance+.05);
 const hit=ray.intersectObjects(meshes,false).find(h=>{
  const m=(h.object as THREE.Mesh).material;
  return (Array.isArray(m)?m[h.face?.materialIndex??0]:m)?.visible;
 });
 return hit?.point.clone()??null;
}

export function buildingDriveByTarget(building:THREE.Object3D,root:THREE.Vector3){
 for(const offset of [0,1.5,-1.5,3,-3]){
  const origin=root.clone().add(new THREE.Vector3(offset,1.8,0));
  const hit=buildingDriveByHit(building,origin,origin.clone().add(new THREE.Vector3(0,0,40)));
  if(hit)return hit;
 }
 return null;
}

/** Integrated continuous velocity: approach, slow firing pass, accelerating exit.
 * Local +Z faces the target facade; the car drives toward -X in the near lane. */
export function buildingDriveByPose(seconds:number){
 const t=Math.max(0,Math.min(BUILDING_DRIVEBY_SECONDS,seconds));
 const distance=t<2?4*t-.5*t*t:t<3.8?6+2*(t-2):9.6+2*(t-3.8)+(t-3.8)**2/1.2;
 const speed=t<2?4-t:t<3.8?2:2+(t-3.8)/.6;
 return {x:10-distance,distance,speed,aim:ease((t-.65)/1.15)*(1-ease((t-3.85)/1.2))};
}

/** No damage is calculated here. The caller supplies captured car/cast/weapon
 * and the facade target, reserves the swept road footprint, and reveals results. */
export class CityBuildingDriveBy {
 readonly root=new THREE.Group();
 readonly duration=BUILDING_DRIVEBY_SECONDS;
 readonly shots:readonly number[];
 private seat:THREE.Vector3;
 private wheels:THREE.Object3D[]=[];
 readonly materials:THREE.Material[]=[];
 constructor(readonly car:THREE.Group,readonly driver:THREE.Group,readonly shooter:THREE.Group,
  readonly weapon:THREE.Group,readonly weaponModel:string,readonly target=new THREE.Vector3(0,1.8,8)){
  this.shots=buildingDriveByShots(weaponModel);
  const driverSeat=car.getObjectByName('seat-front-left');
  const passengerSeat=car.getObjectByName('seat-front-right');
  const grips=['left','right'].map(s=>car.getObjectByName('seat-driver-grip-'+s));
  if(!driverSeat||!passengerSeat||grips.some(g=>!g))throw new Error('Drive-by needs an authored two-seat car cabin');
  this.seat=passengerSeat.position.clone();
  seatDriver(driver,driverSeat.position,grips.map(g=>g!.position));
  for(const name of ['headwear-fedora','headwear-cap']){const hat=shooter.getObjectByName(name);if(hat)hat.visible=false;}
  const glazing=new Map<THREE.Material,THREE.Material>();
  car.traverse(o=>{
   if(o.name.startsWith('wheel-roll-'))this.wheels.push(o);
   if(!(o instanceof THREE.Mesh))return;
   const glaze=(m:THREE.Material)=>{
    if(!(m instanceof THREE.MeshStandardMaterial)||m.name!=='car glass')return m;
    let owned=glazing.get(m);if(!owned){owned=m.clone();owned.transparent=true;owned.opacity=.38;owned.depthWrite=false;glazing.set(m,owned);this.materials.push(owned);}return owned;
   };
   o.material=Array.isArray(o.material)?o.material.map(glaze):glaze(o.material);
  });
  // The passenger's window is wound down; retain its frame and closed door.
  car.getObjectByName('car-door-front-right')?.traverse(o=>{
   if(o instanceof THREE.Mesh&&(Array.isArray(o.material)?o.material:[o.material]).some(m=>m.name==='car glass'))o.visible=false;
  });
  shooter.add(weapon);car.add(driver,shooter);this.root.add(car);
  car.rotation.y=-Math.PI/2;
  this.update(0);
 }
 update(seconds:number){
  const p=buildingDriveByPose(seconds);
  this.car.position.set(p.x,0,0);
  this.root.updateMatrixWorld(true);
  const target=this.car.worldToLocal(this.root.localToWorld(this.target.clone())).sub(this.seat);
  const yaw=Math.atan2(target.x,target.z)*p.aim;
  // Pelvis remains planted in the seat. Legs keep facing the dashboard while
  // the upper body turns; no scaled-down people or feet through the door.
  this.shooter.rotation.set(.38,yaw,0,'YXZ');
  this.shooter.position.copy(this.seat).sub(new THREE.Vector3(0,.86,0).applyQuaternion(this.shooter.quaternion));
  this.root.updateMatrixWorld(true);
  for(const side of [-1,1]){
   const hip=this.shooter.getObjectByName('leg'+side),knee=this.shooter.getObjectByName('knee'+side);
   if(hip){
    hip.position.copy(this.shooter.worldToLocal(this.car.localToWorld(this.seat.clone().add(new THREE.Vector3(side*.12,0,0)))));
    hip.quaternion.copy(this.shooter.quaternion).invert().multiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(1,0,0),-Math.PI/2));
   }
   if(knee)knee.rotation.x=Math.PI/2+1.3;
  }
  const last=[...this.shots].reverse().find(at=>seconds>=at);
  const recoil=last===undefined?0:Math.max(0,1-(seconds-last)/.13);
  if(this.weaponModel==='revolver'){
   this.weapon.position.set(.20-.18*p.aim,.9+.80*p.aim,.10+.24*p.aim-.025*recoil);
   this.weapon.rotation.set(-.38*p.aim-.07*recoil,0,0);
   aimArm(this.shooter,1,this.weapon.position);
   aimArm(this.shooter,-1,new THREE.Vector3(-.14,1.06,.24));
  }else{
   poseLongGun(this.shooter,this.weapon,-Math.PI/2*p.aim-.1*recoil,-.5);
   this.weapon.position.y+=.30*p.aim;
   this.weapon.position.z-=.22*p.aim;
   // Counter the seated lean so the barrel points out through the window.
   this.weapon.rotation.x-=.38*p.aim;
   this.shooter.updateMatrixWorld(true);
   for(const side of [-1,1]){
    const grip=side<0?this.weapon.getObjectByName('support-grip'):this.weapon;
    if(grip)aimArm(this.shooter,side,this.shooter.worldToLocal(grip.getWorldPosition(new THREE.Vector3())));
   }
   const pump=this.weapon.getObjectByName('pump-slide'),age=last===undefined?-1:seconds-last;
   if(pump)pump.position.z=age>.15&&age<.70?-.095*Math.sin((age-.15)/.55*Math.PI):0;
  }
  // The facade is a world target, not a fixed firing direction that would
  // sweep past the building as the car moves. Both hands follow the gun.
  this.root.updateMatrixWorld(true);
  const carry=this.weapon.quaternion.clone();
  this.weapon.lookAt(this.root.localToWorld(this.target.clone()));
  this.weapon.quaternion.copy(carry.slerp(this.weapon.quaternion,p.aim));
  this.root.updateMatrixWorld(true);
  aimArm(this.shooter,1,this.shooter.worldToLocal(this.weapon.getWorldPosition(new THREE.Vector3())));
  if(this.weaponModel!=='revolver'){
   const grip=this.weapon.getObjectByName('support-grip');
   if(grip)aimArm(this.shooter,-1,this.shooter.worldToLocal(grip.getWorldPosition(new THREE.Vector3())));
  }
  for(const wheel of this.wheels)wheel.rotation.x=p.distance/.37;
  this.root.updateMatrixWorld(true);
  return {...p,recoil};
 }
 dispose(){this.materials.forEach(m=>m.dispose());}
}
