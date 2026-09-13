import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import {castLooks,castFace,pedestrianModel} from '../.runtime/frontend-test/city3dCast.js';
const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
test('3D appearances follow every backend cast portrait and painted identity',()=>{
 const core=readFileSync(new URL('../core/voices.go',import.meta.url),'utf8');
 const sheet=core.match(/var castLooks[^\{]+\{([^}]+)/)[1].match(/looks(?:Man|Woman)/g).map(s=>s==='looksWoman'?'f':'m').join('');
 assert.equal(castLooks,sheet);
 for(let i=1;i<=24;i++)assert.equal(pedestrianModel('player-selected',i,true),sheet[i-1]==='f'?'woman':'person');
 const painted=core.match(/var paintedLooks[^\{]+\{([^}]+)/)[1];
 for(const match of painted.matchAll(/"([^"]+)": looks(Man|Woman)/g))assert.equal(pedestrianModel(match[1]),match[2]==='Woman'?'woman':'person');
 assert.equal(pedestrianModel('Alex Varga',3,true),'woman');
 assert.equal(pedestrianModel('Alex Varga',0,true),'person');
});
test('generated identities remain stable when a public face is unavailable',()=>{
 assert.equal(castFace('hello'),Number(0x4f9f2cabn%24n));
 assert.equal(castFace('foobar'),Number(0xbf9cf968n%24n));
 for(const id of ['npc-1','npc-2','mara','Élodie']){
  const face=castFace(id)+1;
  assert.equal(pedestrianModel(id),pedestrianModel(id,face));
  assert.equal(pedestrianModel(id),pedestrianModel(id,NaN));
 }
});
test('both articulated models fit walking and full casualty reservations',async()=>{
 const {trafficSize,trafficSpeed}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const {casualtyFall}=await import('../.runtime/frontend-test/city3dEvents.js');
 for(const model of ['person','woman']){
  const [min,max]=manifest[model].motion_bounds_blender,size=trafficSize(model);
  assert.equal(trafficSpeed(model),1.8);
  assert.ok(size.width/2>=Math.max(-min[0],max[0]));
  assert.ok(size.length/2>=Math.max(-min[1],max[1]));
  assert.ok(min[2]+.25>=.17);
  const [lo,hi]=manifest[model].bounds_blender,fallSize=trafficSize('casualty');
  for(let i=0;i<=120;i++){
   const {rotation,height}=casualtyFall(i/120);
   for(const x of [lo[0],hi[0]])for(const y of [lo[2],hi[2]]){
    assert.ok(height+x*Math.sin(rotation)+y*Math.cos(rotation)>=.17);
    assert.ok(Math.abs(x*Math.cos(rotation)-y*Math.sin(rotation)-.8)<=fallSize.width/2);
   }
  }
 }
});

test('exported models retain both arms, hips and knees as articulated joints',()=>{
 for(const model of ['person','woman']){
  const glb=readFileSync(new URL(`../public/art/models/${model}.glb`,import.meta.url));
  const json=JSON.parse(glb.subarray(20,20+glb.readUInt32LE(12)).toString());
  for(const joint of ['arm1','arm-1','leg1','leg-1','knee1','knee-1']){
   const node=json.nodes.find(n=>n.name===joint);assert.ok(node?.children?.length,`${model}/${joint}`);
  }
 }
});
