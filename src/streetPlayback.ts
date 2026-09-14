import type {StreetSegment} from './types';
/** Retain scenes which expire mid-journey; newer saved records win by ID. */
export function responseRecords<T extends {id:string}>(before:readonly T[] = [],after:readonly T[] = []):T[]{
 return [...new Map([...before,...after].map(record=>[record.id,record])).values()];
}
/** Sample only travel actually observed within the committed command. */
export function streetAt(segments:readonly StreetSegment[],minute:number){
 const active=new Map<string,{segment:StreetSegment;progress:number}>();
 for(const segment of segments){
  if(minute<segment.from_minute||minute>segment.to_minute||minute===segment.to_minute&&segment.end_progress>=1)continue;
  const t=Math.max(0,Math.min(1,(minute-segment.from_minute)/Math.max(1,segment.to_minute-segment.from_minute)));
  active.set(segment.id,{segment,progress:segment.progress+(segment.end_progress-segment.progress)*t});
 }
 return active;
}

/** Spend the remaining displayed time over the physical travel still to render.
 * Traffic can keep moving while the player yields, but the clock cannot consume
 * its final minute before that player reaches the destination. */
export function advanceJourneyClock(progress:number, playerProgress:number, seconds:number, duration:number, routeSeconds:number) {
 if(playerProgress>=1)return 1;
 const step=Math.max(0,seconds);
 const remaining=Math.max(step,Math.max(0,1-playerProgress)*routeSeconds);
 const rate=Math.min(1/Math.max(.001,duration),(1-progress)/Math.max(.001,remaining));
 return Math.min(1-1e-7,Math.max(progress,progress+step*rate));
}

export const streetLegKey=(segment:StreetSegment)=>`${segment.from_id}:${segment.to_id}:${segment.from_minute}:${segment.vehicle||''}`;
/** Keep observed legs in order when cosmetic traffic delays an arrival. */
export class StreetPlayback {
 private legs=new Map<string,StreetSegment[]>();
 private next=new Map<string,number>();
 constructor(segments:readonly StreetSegment[]){
  for(const segment of segments){
   const legs=this.legs.get(segment.id)||[];legs.push(segment);this.legs.set(segment.id,legs);
  }
  for(const legs of this.legs.values())legs.sort((a,b)=>a.from_minute-b.from_minute);
 }
 sample(minute:number,arrived:(id:string,key:string)=>boolean){
  const samples=new Map<string,{segment:StreetSegment;progress:number}>();
  for(const [id,legs] of this.legs){
   let index=this.next.get(id)||0;
   while(index<legs.length){
    const segment=legs[index];
    if(minute<segment.from_minute)break;
    if(segment.end_progress>=1&&minute>=segment.to_minute&&arrived(id,streetLegKey(segment))){index++;continue;}
    const fraction=Math.max(0,Math.min(1,(minute-segment.from_minute)/Math.max(1,segment.to_minute-segment.from_minute)));
    samples.set(id,{segment,progress:segment.progress+(segment.end_progress-segment.progress)*fraction});
    break;
   }
   this.next.set(id,index);
  }
  return samples;
 }
}
