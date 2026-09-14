import type {Card} from './cards';
type CardsState={hands:number;mine:Card[];board:Card[];seats:{who:string;folded:boolean;cards?:Card[]}[]};
export const pokerSeatX=(index:number,count:number)=>(index-(count-1)/2)*Math.min(1.05,2.65/Math.max(1,count-1));
export type PokerCard={id:string;card?:Card;x:number;z:number};
export function pokerCards(p:CardsState|null):PokerCard[]{
 if(!p)return [];
 const hands=[p.mine.map((card,i)=>({id:`mine-${i}`,card,x:(i-.5)*.30,z:.78})),
  ...p.seats.map((s,i)=>s.folded?[]:[0,1].map(j=>({id:`${s.who}-${j}`,card:s.cards?.[j],x:pokerSeatX(i,p.seats.length)+(j-.5)*.29,z:-.78})))];
 return [
  ...[0,1].flatMap(i=>hands.flatMap(hand=>hand[i]?[hand[i]]:[])),
  ...p.board.map((card,i)=>({id:`board-${i}`,card,x:(i-2)*.31,z:0})),
 ];
}
export function planPoker(before:CardsState|null,after:CardsState|null,animate:boolean){
 const prior=new Map(pokerCards(before).map(c=>[c.id,c]));
 const fresh=!before||before.hands!==after?.hands;
 let sequence=0;
 const cards=pokerCards(after).map(card=>{
  const old=fresh?undefined:prior.get(card.id);
  const changed=!old||JSON.stringify(old.card)!==JSON.stringify(card.card);
  const flip=!!old&&!old.card&&!!card.card;
  return {...card,deal:!old,flip,delay:animate&&changed?sequence++*130:0,duration:animate&&changed?430:0};
 });
 return {cards,duration:Math.max(0,...cards.map(c=>c.delay+c.duration))};
}
export function pokerCardPose(move:ReturnType<typeof planPoker>['cards'][number],elapsed:number){
 const t=move.duration?Math.max(0,Math.min(1,(elapsed-move.delay)/move.duration)):1;
 const smooth=t*t*(3-2*t),turn=move.flip?Math.PI*(1-smooth):0;
 return {visible:!move.deal||elapsed>=move.delay,turn,x:move.deal&&t<1?1.55+(move.x-1.55)*smooth:move.x,
  y:.924+.14*Math.sin(t*Math.PI)+(move.flip?.13*Math.sin(turn):0),z:move.deal&&t<1?-.3+(move.z+.3)*smooth:move.z};
}
