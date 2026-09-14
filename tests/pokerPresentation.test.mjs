import test from 'node:test';
import assert from 'node:assert/strict';
import {planPoker,pokerCardPose,pokerCards} from '../.runtime/frontend-test/pokerPresentation.js';
const card=value=>({value,rank:String(value),suit:'spades'});
const hand={hands:1,mine:[card(1),card(12)],board:[],seats:[{who:'a',folded:false},{who:'b',folded:false}]};
test('poker public plan deals round-robin and never invents opponent cards',()=>{
 const p=planPoker(null,hand,true);
 assert.deepEqual(p.cards.map(c=>c.id),['mine-0','a-0','b-0','mine-1','a-1','b-1']);
 assert.equal(p.cards.filter(c=>c.card!==undefined).length,2);
 assert.equal(pokerCardPose(p.cards[5],0).visible,false);
 for(const c of p.cards){const pose=pokerCardPose(c,p.duration);assert.equal(pose.x,c.x);assert.equal(pose.z,c.z);assert.equal(pose.turn,0);}
});
test('flop moves only three public board cards; showdown flips only newly exposed hands',()=>{
 const flop={...hand,board:[card(2),card(3),card(4)]};
 const p=planPoker(hand,flop,true);assert.equal(p.cards.filter(c=>c.duration).length,3);
 const exposed={...flop,seats:[{who:'a',folded:false,cards:[card(8),card(9)]},{who:'b',folded:true}]};
 const reveal=planPoker(flop,exposed,true);
 assert.equal(reveal.cards.filter(c=>c.flip).length,2);
 assert.equal(reveal.cards.some(c=>c.id.startsWith('b-')),false);
 const move=reveal.cards.find(c=>c.id==='a-0');assert.equal(pokerCardPose(move,0).turn,Math.PI);
 assert.equal(pokerCardPose(move,reveal.duration).turn,0);
});
test('restored/reduced-motion hands settle immediately and a new hand deals again',()=>{
 assert.equal(planPoker(null,hand,false).duration,0);
 assert.equal(planPoker(hand,hand,true).duration,0);
 assert.ok(planPoker(hand,{...hand,hands:2},true).duration>0);
 assert.equal(pokerCards(null).length,0);
});
