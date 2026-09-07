import test from 'node:test';
import assert from 'node:assert/strict';
import {pickBuilding} from '../public/art/hit-test.js';
const mask=(alpha)=>({width:2,height:1,data:new Uint8ClampedArray([0,0,0,alpha[0],0,0,0,alpha[1]])});
test('transparent front sprite padding does not steal rear building selection',()=>{
 const buildings=[{id:'front',x:0,y:0,w:100,h:50,depth:90},{id:'rear',x:0,y:0,w:100,h:50,depth:20}];
 const masks={front:mask([0,255]),rear:mask([255,255])};
 assert.equal(pickBuilding(buildings,masks,20,20).id,'rear');
 assert.equal(pickBuilding(buildings,masks,80,20).id,'front');
 assert.equal(pickBuilding(buildings,masks,100,20),undefined);
 assert.equal(pickBuilding(buildings,masks,20,50),undefined);
 assert.equal(pickBuilding(buildings,masks,-1,20),undefined);
 assert.equal(buildings[0].id,'front');
});
test('unloaded or transparent artwork does not acquire a click target',()=>{
 const b=[{id:'market',x:0,y:0,w:100,h:50,depth:10}];
 assert.equal(pickBuilding(b,{},25,25),undefined);
 assert.equal(pickBuilding(b,{market:mask([0,0])},25,25),undefined);
});
