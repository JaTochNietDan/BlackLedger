import * as THREE from 'three';
import {availableSceneSlot} from './city3dEvents.js';
import type {SceneSlot} from './city3dEvents.js';
import type {Lot} from './city3dPlan.js';
import {vehicleRootHeight} from './city3dPlan.js';
import {dressPedestrian, wardrobe} from './city3dWardrobe.js';
import type {Snapshot} from './types';

type Entry = {group: THREE.Group; slot: SceneSlot; owned: THREE.Material[]};
export class CityAftermath {
  readonly root = new THREE.Group();
  private entries = new Map<string, Entry>();
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
  update(records: NonNullable<Snapshot['aftermath']>, minute: number, lots: Map<string,Lot>, models: Map<string,THREE.Group>,
    modelFor: (id:string)=>string, occupied: SceneSlot[], animating: Set<string>, presence: NonNullable<Snapshot['police_presence']> = [], activeRaids = new Set<string>(), fires: NonNullable<Snapshot['building_fires']> = []) {
    const desired=new Set<string>();
    const scenes = [
      ...records.map(record=>({...record,raid:false,fire:false})),
      ...presence.map(record=>({...record,raid:true,fire:false,police_at:record.minute,victim:{id:'',name:''}})),
      ...fires.map(record=>({...record,raid:false,fire:true,police_at:record.brigade_at,victim:{id:'',name:''}})),
    ];
    for(const record of scenes) {
      if(minute<record.minute || minute>=record.cleanup_at)continue;
      if(record.raid && activeRaids.has(record.target))continue;
      const cast=record.fire ? (minute>=record.police_at?['fire-engine','firefighter-a','firefighter-b']:[]) : record.raid ? ['police','police-b','police-c','officer-a','officer-b','officer-c','officer-d']
        : minute>=record.police_at ? ['body','police','officer-a','officer-b'] : ['body'];
      for(const kind of cast) {
        const vehicle=kind.startsWith('police')||kind==='fire-engine';
        const key=`aftermath:${record.id}:${kind}`;
        if(kind==='body' && animating.has(record.victim.id))continue;
        desired.add(key);
        if(this.entries.has(key))continue;
        const lot=lots.get(record.target); if(!lot)continue;
        const taken=[...occupied,...[...this.entries.values()].map(e=>e.slot)];
        const slot=availableSceneSlot(lot,kind==='fire-engine'?'fire-engine':vehicle?'arrest':'killing',taken);if(!slot)continue;
        const model=kind.startsWith('firefighter')?'firefighter':kind==='fire-engine'?'fire-engine':kind==='body'?modelFor(record.victim.id):vehicle?'police':'police-officer';
        const source=models.get(model);if(!source)continue;
        const group=new THREE.Group(), object=source.clone(true);
        const owned=kind==='body'?dressPedestrian(object,model,wardrobe(record.victim.id)):[];
        group.position.set(slot.root.x,kind==='body'?0:vehicle?vehicleRootHeight(slot.root):.2,slot.root.z);
        if(kind==='body') {
          object.rotation.z=-Math.PI/2;object.position.y=.6;
          const pool=new THREE.Mesh(this.pool,this.blood);
          pool.rotation.x=-Math.PI/2;pool.position.set(1.1,.181,0);group.add(pool);
        }
        if(kind.startsWith('firefighter'))object.rotation.y=Math.atan2(lot.x-slot.root.x,lot.z-slot.root.z);
        if(kind.startsWith('officer')) {
          const body=this.entries.get(`aftermath:${record.id}:body`);
          object.rotation.y=Math.atan2((body?.slot.root.x ?? lot.x)+.8-slot.root.x,(body?.slot.root.z ?? lot.z)-slot.root.z);
        }
        object.traverse(part=>{if(part instanceof THREE.Mesh){part.castShadow=true;part.receiveShadow=true;}});
        group.add(object);this.root.add(group);this.entries.set(key,{group,slot,owned});
      }
    }
    for(const [key,entry] of this.entries) if(!desired.has(key)) {
      this.root.remove(entry.group);entry.owned.forEach(m=>m.dispose());this.entries.delete(key);
    }
  }
  reservations() {return [...this.entries].map(([id,e])=>({id,model:e.slot.model,points:[e.slot.pose],progress:0}));}
  slots() {return [...this.entries.values()].map(e=>e.slot);}
  show(placements: Map<string,{waiting:boolean}>) {
    for(const [id,e] of this.entries)e.group.visible=!!placements.get(id)&&!placements.get(id)!.waiting;
  }
  inspect(){return [...this.entries].map(([id,e])=>({id,x:e.slot.root.x,z:e.slot.root.z,visible:e.group.visible}));}
  dispose(){for(const e of this.entries.values())e.owned.forEach(m=>m.dispose());this.entries.clear();this.root.clear();this.pool.dispose();this.blood.dispose();}
}
