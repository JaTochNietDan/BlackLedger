import {poolRails} from './billiardsTableGeometry.js';
import type {PoolState} from './billiards.js';

// Geometric first contact only. No spin, rebound, pocket or win prediction;
// the authoritative Go stroke remains responsible for every actual outcome.
export function poolAimContact(p:Pick<PoolState,'balls'|'width'|'length'|'radius'>,angle:number){
 const cue=p.balls.find(b=>b.id===0&&b.pocket<0);if(!cue||!Number.isFinite(angle))return null;
 const [x,y]=cue.position,dx=Math.cos(angle),dy=Math.sin(angle),r=p.radius;
 let distance=Math.hypot(p.width,p.length)+.3,ball:number|null=null,kind:'ball'|'cushion'|'edge'='edge';
 const accept=(t:number,next:typeof kind,id:number|null=null)=>{if(t>=-1e-8&&t<distance){distance=Math.max(0,t);kind=next;ball=id;}};
 const circle=(cx:number,cy:number,radius:number)=>{
  const ox=x-cx,oy=y-cy,b=ox*dx+oy*dy,c=ox*ox+oy*oy-radius*radius,disc=b*b-c;
  if(disc<0)return Infinity;
  if(c<=1e-10)return b<0?0:Infinity;
  const t=-b-Math.sqrt(disc);return t>=0?t:Infinity;
 };
 for(const b of p.balls)if(b.id!==0&&b.pocket<0)accept(circle(b.position[0],b.position[1],2*r),'ball',b.id);
 // Sweep a radius-r circle against each physical cushion segment, including
 // jaw tips: its Minkowski boundary is a rectangle and two endpoint circles.
 for(const [ax,ay,bx,by] of poolRails(p.width,p.length)){
  const length=Math.hypot(bx-ax,by-ay),tx=(bx-ax)/length,ty=(by-ay)/length,nx=-ty,ny=tx;
  const side=(x-ax)*nx+(y-ay)*ny,velocity=dx*nx+dy*ny;
  if(Math.abs(velocity)>1e-12)for(const sign of [-1,1]){
   if(sign*velocity>=0)continue;
   const t=(sign*r-side)/velocity,along=(x+dx*t-ax)*tx+(y+dy*t-ay)*ty;
   if(along>=0&&along<=length)accept(t,'cushion');
  }
  accept(circle(ax,ay,r),'cushion');accept(circle(bx,by,r),'cushion');
 }
 // A line through a pocket opening ends just beyond the bed, not across the room.
 for(const [origin,velocity,low,high] of [[x,dx,-.10,p.width+.10],[y,dy,-.10,p.length+.10]]){
  if(Math.abs(velocity)>1e-12)accept(((velocity>0?high:low)-origin)/velocity,'edge');
 }
 return {x:x+dx*distance,y:y+dy*distance,distance,ball,kind};
}
