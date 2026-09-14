import type {Card} from './cards';
export type PublicHand={playing:boolean;settled?:boolean;mine?:Card[];theirs?:Card[]};
export type CardMove={card?:Card;row:number;index:number;from:number[];to:number[];delay:number;duration:number;flip:boolean};
export const cardPosition=(index:number,count:number,row:number)=>[(index-(count-1)/2)*Math.min(.35,1.7/Math.max(1,count-1)),.875+index*.0001,row===0?-.38:.48];
function rows(hand:PublicHand):(Card|undefined)[][]{
 const dealer: (Card|undefined)[]=[...(hand.theirs??[])];
 if(hand.playing&&dealer.length<2)dealer.push(undefined);
 return [dealer,hand.mine??[]];
}
export function planCards(before:PublicHand,after:PublicHand,animate:boolean){
 const old=rows(before),next=rows(after),fresh=!before.playing;let dealt=0;
 const cards:CardMove[]=[];
 for(let row=0;row<2;row++)for(let i=0;i<next[row].length;i++){
  const card=next[row][i],prior=old[row][i],changed=fresh||i>=old[row].length||JSON.stringify(card)!==JSON.stringify(prior);
  const flip=animate&&!fresh&&i<old[row].length&&!prior&&!!card;
  const to=cardPosition(i,next[row].length,row);
  cards.push({card,row,index:i,to,flip,from:changed&&!flip?[.94,1.12,-.30]:cardPosition(i,old[row].length,row),delay:animate&&changed?(fresh?(i*2+(row===1?0:1)):dealt++)*380:0,duration:animate?(changed?480:240):0});
 }
 return {cards,duration:Math.max(0,...cards.map(c=>c.delay+c.duration))};
}
export function cardPose(move:CardMove,elapsed:number){
 const t=move.duration?Math.max(0,Math.min(1,(elapsed-move.delay)/move.duration)):1;
 const ease=t*t*(3-2*t);
 const rotation=move.flip?Math.PI*(1-ease):0;
 const lift=move.flip?.147*Math.abs(Math.sin(rotation)):0;
 return {rotation,visible:move.flip||elapsed>=move.delay,position:t===1?[...move.to]:move.to.map((v,i)=>move.from[i]+(v-move.from[i])*ease+(i===1?lift+(move.flip?.02:.12)*Math.sin(Math.PI*t):0))};
}
