import test from 'node:test';
import assert from 'node:assert/strict';
import {planCards,cardPose} from '../.runtime/frontend-test/blackjackPresentation.js';
const card=(rank,suit='clubs')=>({rank,suit,value:Number(rank)||10});
const hand={playing:true,mine:[card('5'),card('5','hearts')],theirs:[card('8','diamonds')]};
test('dealing preserves public cards, conceals the hole card and finishes above the cloth',()=>{
 const plan=planCards({playing:false},hand,true);
 assert.equal(plan.cards.length,4);assert.equal(plan.cards.filter(c=>!c.card).length,1);
 assert.ok(plan.duration>1000);
 assert.deepEqual([...plan.cards].sort((a,b)=>a.delay-b.delay).map(c=>[c.row,c.index]),[[1,0],[0,0],[1,1],[0,1]]);
 for(const move of plan.cards){
  if(move.delay)assert.equal(cardPose(move,move.delay-1).visible,false);
  for(let ms=0;ms<=plan.duration;ms+=10)assert.ok(cardPose(move,ms).position[1]>=.875-1e-9);
  assert.deepEqual(cardPose(move,plan.duration).position,move.to);
 }
});
test('hit moves existing cards without redealing and stand only reveals authoritative dealer draws',()=>{
 const hit={...hand,mine:[...hand.mine,card('K','diamonds')]};
 const p=planCards(hand,hit,true);assert.equal(p.duration,480);
 assert.equal(p.cards.filter(c=>c.duration===480).length,1);
 const settled={...hit,playing:false,settled:true,theirs:[...hit.theirs,card('2'),card('3'),card('5','spades')]};
 const q=planCards(hit,settled,true);
 assert.equal(q.cards.filter(c=>!c.card).length,0);
 assert.equal(q.cards.filter(c=>c.duration===480).length,3);
 assert.deepEqual(q.cards.filter(c=>c.row===0).map(c=>c.card),settled.theirs);
 for(const state of [hand,hit,settled]){
  const immediate=planCards(state,state,false);assert.equal(immediate.duration,0);
  for(const move of immediate.cards)assert.deepEqual(cardPose(move,0).position,move.to);
 }
});
