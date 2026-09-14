export type PoolPose=[number,number,number,number,number,number,number,number,number];
export interface PoolInput {x?:number;y?:number;angle?:number;speed?:number;top?:number;side?:number;ball?:number;pocket?:number;safety?:boolean}
export interface PoolState {
 place:string;opponent:string;opponent_name:string;stake:number;pot:number;settled:boolean;voided:boolean;
 turn:number;groups:[number,number];breaking:boolean;ball_in_hand:boolean;behind_head_string:boolean;winner:number;shots:number;
 balls:{id:number;position:[number,number,number];rotation:[number,number,number,number];pocket:number}[];
 legal_balls:number[];break_choices:{id:string;label:string}[];outcome:string;foul:boolean;unavailable:string;
 width:number;length:number;radius:number;replay:string;
 stroke:null|{shooter:number;intent:PoolInput;placement?:[number,number,number];decision:string};
}
export interface PoolOpponent {id:string;name:string;max_stake:number;unavailable:string}
export interface PoolReplay {v:1;duration:number;frames:{t:number;balls:PoolPose[]}[];events:{Time:number;Kind:string;Ball:number;Other:number;Speed:number}[]}
export async function decodePoolReplay(encoded:string):Promise<PoolReplay>{
 if(encoded.length>2*1024*1024)throw new Error('Billiards replay exceeds its size limit.');
 const bytes=Uint8Array.from(atob(encoded),c=>c.charCodeAt(0));
 const reader=new Blob([bytes]).stream().pipeThrough(new DecompressionStream('deflate')).getReader();
 const chunks:Uint8Array[]=[];let size=0;
 try{for(;;){const {done,value}=await reader.read();if(done)break;size+=value.length;if(size>8*1024*1024)throw new Error('Billiards replay expands beyond its size limit.');chunks.push(value);}}
 finally{await reader.cancel();reader.releaseLock();}
 const raw=new Uint8Array(size);let offset=0;for(const c of chunks){raw.set(c,offset);offset+=c.length;}
 const r=JSON.parse(new TextDecoder().decode(raw)) as PoolReplay;
 if(r.v!==1||!Number.isFinite(r.duration)||r.duration<0||r.duration>60||!Array.isArray(r.frames)||!r.frames.length||!Array.isArray(r.events))throw new Error('Invalid billiards replay.');
 let time=-1;for(const f of r.frames){
  if(!Number.isFinite(f.t)||f.t<time||f.t>r.duration+.001||!Array.isArray(f.balls)||f.balls.length!==16)throw new Error('Invalid billiards frame.');time=f.t;
  const ids=new Set<number>();for(const b of f.balls){if(b.length!==9||!b.every(Number.isFinite)||!Number.isInteger(b[0])||b[0]<0||b[0]>15||ids.has(b[0])||!Number.isInteger(b[8])||b[8]<0||b[8]>6)throw new Error('Invalid billiards ball.');ids.add(b[0]);}
 }
 return r;
}
// Binary search preserves the actual adjacent impact samples; never resample a
// whole stroke into a straight start/end path through cushions or other balls.
export function poolFramePair(r:PoolReplay,time:number){
 let lo=0,hi=r.frames.length-1;
 while(lo<hi){const mid=Math.ceil((lo+hi)/2);if(r.frames[mid].t<=time)lo=mid;else hi=mid-1;}
 const a=r.frames[lo],b=r.frames[Math.min(lo+1,r.frames.length-1)];
 return {a,b,mix:b.t>a.t?Math.max(0,Math.min(1,(time-a.t)/(b.t-a.t))):0};
}
export const poolPockets=['Left near','Left far','Right near','Right far','Left middle','Right middle'];
