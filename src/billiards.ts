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

export function poolPocketCenters(width:number,length:number):[number,number][] {
 return [[-.026,-.026],[-.026,length+.026],[width+.026,-.026],[width+.026,length+.026],[-.045,length/2],[width+.045,length/2]];
}
// Local placement feedback only; Go validates the committed position again.
export function poolPlacementHint(p:Pick<PoolState,'width'|'length'|'radius'|'behind_head_string'|'balls'>,x:number,y:number):string {
 if(!Number.isFinite(x)||!Number.isFinite(y)||x<p.radius||x>p.width-p.radius||y<p.radius||y>p.length-p.radius)return 'Keep the cue ball inside the cushions.';
 if(p.behind_head_string&&y>=p.length/4)return 'Place behind the dashed head string.';
 if(p.balls.some(b=>b.id!==0&&b.pocket<0&&Math.hypot(b.position[0]-x,b.position[1]-y)<2*p.radius+1e-5))return 'Leave room around the other balls.';
 return '';
}
export function poolAimAngle(cue:[number,number],point:[number,number]):number|null {
 if(![...cue,...point].every(Number.isFinite)||Math.hypot(point[0]-cue[0],point[1]-cue[1])<.001)return null;
 return (Math.atan2(point[1]-cue[1],point[0]-cue[0])+Math.PI*2)%(Math.PI*2);
}
export class PoolTap {
 private active:{id:number;x:number;y:number;moved:boolean}|null=null;
 begin(id:number,x:number,y:number,primary:boolean,button:number){
  if(this.active||!primary||button!==0){this.active=null;return;}
  this.active={id,x,y,moved:false};
 }
 move(id:number,x:number,y:number){const a=this.active;if(a&&a.id===id&&Math.hypot(x-a.x,y-a.y)>6)a.moved=true;}
 end(id:number,x:number,y:number){this.move(id,x,y);const a=this.active;this.active=null;return !!a&&a.id===id&&!a.moved;}
 cancel(){this.active=null;}
}

export const POOL_CUE_CONTACT=.72;
export const POOL_CUE_END=1.12;
// Cue front is measured along the shot direction from the initial cue-ball
// centre. Only presentation time is delayed; saved physical samples are intact.
export function poolCueStroke(seconds:number,speed:number,radius:number,top=0,side=0){
 const t=Math.max(0,seconds),power=Math.max(0,Math.min(1,speed/8));
 const contact=-Math.sqrt(Math.max(0,radius*radius-top*top-side*side));
 const ready=contact-.045,pull=.08+.18*power;
 let front:number;
 if(t<.42){const u=t/.42;front=ready-pull*(u*u*(3-2*u));}
 else if(t<POOL_CUE_CONTACT){const u=(t-.42)/.30;front=ready-pull+(contact-ready+pull)*u*u;}
 else if(t<.88){const u=(t-POOL_CUE_CONTACT)/.16;front=contact+(.03+.08*power)*(1-(1-u)*(1-u));}
 else {const u=Math.min(1,(t-.88)/.24);front=contact+(.03+.08*power)-.25*u;}
 return {front,opacity:t<=.88?1:Math.max(0,1-(t-.88)/.24),ballTime:Math.max(0,t-POOL_CUE_CONTACT),contact:t>=POOL_CUE_CONTACT,visible:t<POOL_CUE_END};
}

export interface PoolTournamentState {
 fee:number;pot:number;settled:boolean;voided:boolean;finished:boolean;winner:string;player_id:string;withdrawn:boolean;
 names:Record<string,string>;
 games:{index:number;round:number;table_number:number;players:[string,string];player_seat:number;resolved:boolean;winner:string;table:PoolState|null}[];
}

export interface PoolTournamentNoticeState {opens:number;closes:number;fee:number;entrants:{id:string;name:string}[];pot:number;can_enter:boolean;unavailable:string}
