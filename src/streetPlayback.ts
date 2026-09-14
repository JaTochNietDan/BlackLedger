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
