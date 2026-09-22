import * as THREE from 'three';
import {CityAccident} from './city3dAccident.js';
import {pedestrianModel} from './city3dCast.js';
import {groundCharacter} from './city3dGround.js';
import {availableSceneSlot} from './city3dEvents.js';
import type {SceneSlot} from './city3dEvents.js';
import type {Lot, Point} from './city3dPlan.js';
import {vehicleRootHeight, PITCH, LANE} from './city3dPlan.js';
import type {TrafficPlacement} from './city3dTraffic.js';
import {dressPedestrian, wardrobe} from './city3dWardrobe.js';
import type {Snapshot} from './types';

export type BodyJoints = Record<string,[number,number,number,number]>;
const bodyJointNames=['arm1','arm-1','elbow1','elbow-1','leg1','leg-1','knee1','knee-1'];
/** Copy transforms only: the retained pose never owns a model or material. */
export function captureBodyJoints(actor:THREE.Group):BodyJoints {
  const pose:BodyJoints={};
  for(const name of bodyJointNames){const joint=actor.getObjectByName(name);if(joint)pose[name]=joint.quaternion.toArray();}
  return pose;
}
type Entry = {group: THREE.Group; slot: SceneSlot; owned: THREE.Material[]; victim?:string;
  approach?:Point[]; approachStarted?:boolean; waitingFor?:string; wheelDistance?:number};
export class CityAftermath {
  readonly root = new THREE.Group();
  private entries = new Map<string, Entry>();
  private minute?:number;
  private bodyPoses=new Map<string,{slot:SceneSlot;yaw:number;joints?:BodyJoints}>();
  private blood = new THREE.MeshStandardMaterial({color: 0x480a0b, roughness: .31, metalness: .05, polygonOffset: true, polygonOffsetFactor: -1});
  private pool: THREE.ShapeGeometry;
  constructor() {
    const shape = new THREE.Shape();
    for (let i=0;i<=48;i++) {
      const angle=i/48*Math.PI*2, r=1+.12*Math.sin(angle*7)+.07*Math.sin(angle*13);
      const x=Math.cos(angle)*.65*r, y=Math.sin(angle)*.46*r;
      if(i===0)shape.moveTo(x,y);else shape.lineTo(x,y);
    }
    this.pool=new THREE.ShapeGeometry(shape);
  }
  rememberBody(id:string,slot:SceneSlot,yaw:number,joints?:BodyJoints){this.bodyPoses.set(id,{slot,yaw,joints: joints?structuredClone(joints):undefined});}
  suppressVictims(ids:Set<string>){
    for(const [key,e] of this.entries)if(e.victim&&ids.has(e.victim)){
      this.root.remove(e.group);e.owned.forEach(m=>m.dispose());this.entries.delete(key);
    }
  }
  update(records: NonNullable<Snapshot['aftermath']>, minute: number, lots: Map<string,Lot>, models: Map<string,THREE.Group>,
    modelFor: (id:string)=>string, occupied: SceneSlot[], animating: Set<string>, presence: NonNullable<Snapshot['police_presence']> = [], activeRaids = new Set<string>(), fires: NonNullable<Snapshot['building_fires']> = []) {
    const desired=new Set<string>();
    const before=this.minute;this.minute=minute;
    for(const id of this.bodyPoses.keys())if(!records.some(r=>r.victim.id===id&&minute<r.cleanup_at))this.bodyPoses.delete(id);
    const scenes = [
      ...records.map(record=>({...record,raid:false,fire:false})),
      ...presence.map(record=>({...record,raid:true,fire:false,police_at:record.minute,victim:{id:'',name:''}})),
      ...fires.map(record=>({...record,raid:false,fire:true,police_at:record.brigade_at,victim:{id:'',name:''}})),
    ];
    for(const record of scenes) {
      if(minute<record.minute || minute>=record.cleanup_at)continue;
      if(record.raid && activeRaids.has(record.target))continue;
      if(!record.raid&&!record.fire&&animating.has(record.victim.id))continue;
      const cast=record.fire ? (minute>=record.police_at?['fire-engine','firefighter-a','firefighter-b']:[]) : record.raid ? ['police','police-b','police-c','officer-a','officer-b','officer-c','officer-d']
        : minute>=record.police_at ? ['body','police','officer-a','officer-b'] : ['body'];
      for(const kind of cast) {
        const vehicle=kind.startsWith('police')||kind==='fire-engine';
        const key=`aftermath:${record.id}:${kind}`;
        desired.add(key);
        if(this.entries.has(key))continue;
        const lot=lots.get(record.target); if(!lot)continue;
        const taken=[...occupied,...[...this.entries.values()].map(e=>e.slot)];
        const remembered=kind==='body'?this.bodyPoses.get(record.victim.id):undefined;
        const accident='cause' in record && record.cause==='charge-accident';
        const face='face' in record&&typeof record.face==='number'?record.face:undefined;
        const slot=remembered?.slot || availableSceneSlot(lot,kind==='fire-engine'?'fire-engine':vehicle?'arrest':kind==='body'&&accident?'accident':'killing',taken);if(!slot)continue;
        const model=kind.startsWith('firefighter')?'firefighter':kind==='fire-engine'?'fire-engine':kind==='body'?(accident?pedestrianModel(record.victim.name,face,true):modelFor(record.victim.id)):vehicle?'police':'police-officer';
        const source=models.get(model);if(!source)continue;
        const group=new THREE.Group(), object=source.clone(true);
        const owned=kind==='body'?dressPedestrian(object,model,accident?wardrobe(record.victim.name,face,true):wardrobe(record.victim.id)):[];
        group.position.set(slot.root.x,kind==='body'?0:vehicle?vehicleRootHeight(slot.root):.2,slot.root.z);
        if(kind==='body') {
          if(accident){
            const fall=new CityAccident(object,true);fall.update(fall.duration);
            group.position.y=.2;group.rotation.y=slot.pose.heading;
          }else{
          for(const [name,rotation] of Object.entries(remembered?.joints||{}))object.getObjectByName(name)?.quaternion.fromArray(rotation);
          object.quaternion.setFromAxisAngle(new THREE.Vector3(0,1,0),remembered?.yaw||0)
            .premultiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,0,1),-Math.PI/2));object.position.y=.6;
          }
          const pool=new THREE.Mesh(this.pool,this.blood);
          pool.rotation.x=-Math.PI/2;pool.position.set(accident?0:1.1,accident?.006:.181,0);group.add(pool);
        }
        if(kind.startsWith('firefighter'))object.rotation.y=Math.atan2(lot.x-slot.root.x,lot.z-slot.root.z);
        if(kind.startsWith('officer')) {
          const body=this.entries.get(`aftermath:${record.id}:body`);
          object.rotation.y=Math.atan2((body?.slot.root.x ?? lot.x)+.8-slot.root.x,(body?.slot.root.z ?? lot.z)-slot.root.z);
        }
        object.traverse(part=>{if(part instanceof THREE.Mesh){part.castShadow=true;part.receiveShadow=true;}});
        group.add(object);
        if(kind==='body'&&!accident)groundCharacter(object,.205);
        const arriving=!record.raid&&!record.fire&&kind==='police'&&before!==undefined&&before<record.police_at&&minute>=record.police_at;
        const approach=arriving?[{x:slot.root.x,z:lot.row*PITCH+LANE},{...slot.root}]:undefined;
        if(approach){group.position.set(approach[0].x,vehicleRootHeight(approach[0]),approach[0].z);group.visible=false;}
        const waitingFor=!record.raid&&!record.fire&&kind.startsWith('officer')?`aftermath:${record.id}:police`:undefined;
        this.root.add(group);this.entries.set(key,{group,slot,owned,victim:record.victim.id||undefined,approach,waitingFor,wheelDistance:0});
      }
    }
    for(const [key,entry] of this.entries) if(!desired.has(key)) {
      this.root.remove(entry.group);entry.owned.forEach(m=>m.dispose());this.entries.delete(key);
    }
  }
  object(id:string){return this.entries.get(id)?.group;}
  reservations() {return [...this.entries].map(([id,e])=>({id,model:e.approach?'police':e.slot.model,points:e.approach||[e.slot.pose],progress:e.approach&&e.approachStarted?1:0}));}
  slots() {return [...this.entries.values()].map(e=>e.slot);}
  show(placements: Map<string,{waiting:boolean}&Partial<TrafficPlacement>>) {
    for(const [id,e] of this.entries){
      const placement=placements.get(id);
      const transport=e.waitingFor?this.entries.get(e.waitingFor):undefined;
      e.group.visible=!!placement&&!placement.waiting&&(!e.waitingFor||!!transport&&!transport.approach&&transport.group.visible);
      if(e.approach&&placement?.pose&&!placement.waiting){
        e.approachStarted=true;
        const at=placement.pose;
        const distance=Math.hypot(at.x-e.group.position.x,at.z-e.group.position.z);
        e.wheelDistance=(e.wheelDistance||0)+distance;
        e.group.position.set(at.x,vehicleRootHeight(at),at.z);e.group.rotation.y=at.heading;
        e.group.traverse(part=>{if(part.name.startsWith('wheel-spin'))part.rotation.x=e.wheelDistance!/.37;});
        if((placement.progress||0)>=1)e.approach=undefined;
      }
    }
  }
  inspect(){return [...this.entries].map(([id,e])=>({id,x:e.group.position.x,z:e.group.position.z,visible:e.group.visible,arriving:!!e.approach}));}
  dispose(){this.bodyPoses.clear();for(const e of this.entries.values())e.owned.forEach(m=>m.dispose());this.entries.clear();this.root.clear();this.pool.dispose();this.blood.dispose();}
}
