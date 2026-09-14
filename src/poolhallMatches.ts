import * as THREE from 'three';
import type {PoolTournamentState} from './billiards';
import type {InteriorSpot} from './interiorStaging';
import {ballTexture} from './billiardsBallTexture.js';

export function poolhallTableOrigin(number:number):[number,number]|null {
 if(!Number.isInteger(number)||number<1||number>6)return null;
 return [number<=3?-3.7:3.7,[4.25,0,-4.25][(number-1)%3]];
}
export function tournamentHallSpots(t:PoolTournamentState|null|undefined){
 const spots=new Map<string,InteriorSpot>();
 if(!t||t.settled)return spots;
 for(const g of t.games){const origin=poolhallTableOrigin(g.table_number);if(!origin||g.resolved||!g.table)continue;
  g.players.forEach((id,seat)=>{const side=seat===0?-1:1;spots.set(id===t.player_id?'player':id,{id:`match-${g.index}-seat-${seat}`,x:origin[0]+side*1.35,z:origin[1],yaw:side<0?Math.PI/2:-Math.PI/2});});
 }
 return spots;
}

// The hall shows the authoritative resting positions. Stroke playback stays in
// the close table view; never interpolate an unwatched rack through geometry.
export class PoolhallMatches {
 group=new THREE.Group();
 private geometry=new THREE.SphereGeometry(.028575,16,12);
 private materials=Array.from({length:16},(_,i)=>new THREE.MeshStandardMaterial({map:ballTexture(i),roughness:.22}));
 private decorations:THREE.Object3D[]=[];
 constructor(private room:THREE.Object3D){room.updateMatrixWorld(true);room.traverse(o=>{if(o.name.startsWith('billiard_ball')||o.name.startsWith('billiard ball'))this.decorations.push(o)});}
 update(t:PoolTournamentState|null|undefined){
  this.group.clear();this.decorations.forEach(o=>o.visible=true);
  const tables=new Map<number,NonNullable<PoolTournamentState['games'][number]['table']>>();
  for(const g of t?.games??[])if(g.table&&poolhallTableOrigin(g.table_number))tables.set(g.table_number,g.table);
  for(const [number,table] of tables){const [x,z]=poolhallTableOrigin(number)!;
   for(const o of this.decorations){const p=o.getWorldPosition(new THREE.Vector3());if(Math.abs(p.x-x)<.8&&Math.abs(p.z-z)<1.5)o.visible=false;}
   for(const ball of table.balls){if(ball.pocket>=0)continue;const m=new THREE.Mesh(this.geometry,this.materials[ball.id]);m.position.set(x+ball.position[0]-table.width/2,.78+ball.position[2],z+table.length/2-ball.position[1]);const [qx,qy,qz,qw]=ball.rotation;m.quaternion.setFromAxisAngle(new THREE.Vector3(1,0,0),-Math.PI/2).multiply(new THREE.Quaternion(qx,qy,qz,qw));m.castShadow=true;m.receiveShadow=true;this.group.add(m);}
  }
 }
 dispose(){this.group.removeFromParent();this.group.clear();this.geometry.dispose();for(const m of this.materials){m.map?.dispose();m.dispose();}}
}

export function poolhallGameAt(t:PoolTournamentState|null|undefined,x:number,z:number):number|null {
 if(!Number.isFinite(x)||!Number.isFinite(z))return null;
 for(const g of [...(t?.games??[])].reverse()){
  const origin=poolhallTableOrigin(g.table_number);
  if(g.table&&origin&&Math.abs(x-origin[0])<=g.table.width/2+.17&&Math.abs(z-origin[1])<=g.table.length/2+.17)return g.index;
 }
 return null;
}
