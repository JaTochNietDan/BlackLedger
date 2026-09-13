import test from 'node:test';
import assert from 'node:assert/strict';
import {StreetTraffic,trafficOverlap} from '../.runtime/frontend-test/city3dTraffic.js';
const line=[{x:0,z:0},{x:200,z:0}];
function clear(requests,poses){
 for(const a of requests)for(const b of requests){if(a.id>=b.id)continue;
 const ap=poses.get(a.id),bp=poses.get(b.id);
 if(!ap.waiting&&!bp.waiting)assert.equal(trafficOverlap(ap.pose,a.model,bp.pose,b.model),false,`${a.id} intersects ${b.id}`);
 }
}
test('twelve cars sharing committed progress queue without overlap or backwards movement',()=>{
 const traffic=new StreetTraffic();
 const requests=Array.from({length:12},(_,i)=>({id:`car-${i}`,model:i%2?'packard':'ford',points:line,progress:.7}));
 let previous=traffic.update(requests,1/60);clear(requests,previous);
 assert.equal([...previous.values()].filter(p=>p.waiting).length,0);
 for(let frame=0;frame<180;frame++){
  for(const r of requests)r.progress=.7+.2*Math.min(1,frame/100);
  const poses=traffic.update(requests,1/60);clear(requests,poses);
  for(const r of requests){assert.ok(poses.get(r.id).progress>=previous.get(r.id).progress);assert.ok(poses.get(r.id).progress<=r.progress+1e-8);}
  previous=poses;
 }
});
test('crowded departures remain represented in the waiting set and enter when space opens',()=>{
 const traffic=new StreetTraffic();const requests=[0,1].map(i=>({id:`car-${i}`,model:'ford',points:line,progress:0}));
 let poses=traffic.update(requests,1/60);assert.equal([...poses.values()].filter(p=>p.waiting).length,1);
 for(let frame=0;frame<90;frame++){requests.forEach(r=>r.progress=.5);poses=traffic.update(requests,1/60);clear(requests,poses);}
 assert.ok([...poses.values()].every(p=>!p.waiting&&p.progress>0));
});
test('opposing lanes do not block each other',()=>{
 const traffic=new StreetTraffic();const requests=[{id:'east',model:'packard',points:[{x:0,z:1.6},{x:100,z:1.6}],progress:.5},{id:'west',model:'packard',points:[{x:100,z:-1.6},{x:0,z:-1.6}],progress:.5}];
 const poses=traffic.update(requests,1/60);clear(requests,poses);assert.ok([...poses.values()].every(p=>p.progress===.5&&!p.waiting));
});
test('crossing traffic cannot tunnel through occupied vehicles between frames',()=>{
 const traffic=new StreetTraffic();const requests=[{id:'east',model:'ford',points:[{x:0,z:0},{x:100,z:0}],progress:.4},{id:'north',model:'packard',points:[{x:50,z:-50},{x:50,z:50}],progress:.3}];
 traffic.update(requests,1/60);
 for(let frame=0;frame<180;frame++){requests.forEach(r=>r.progress=.9);clear(requests,traffic.update(requests,1/30));}
 for(const p of traffic.update(requests,1/30).values())assert.ok(p.progress>.89,'crossing traffic must finish');
});
test('a newly restarted journey on the same route never inherits old progress',()=>{
 const traffic=new StreetTraffic();const request={id:'returning',model:'ford',points:line,progress:.8};
 traffic.update([request],1/60);request.progress=.1;
 assert.equal(traffic.update([request],1/60).get(request.id).progress,.1);
 assert.equal(traffic.update([],1/60).size,0);
});
test('occupancy dimensions enclose the authored car meshes',async()=>{
 const {readFileSync}=await import('node:fs');const {trafficSize}=await import('../.runtime/frontend-test/city3dTraffic.js');
 const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
 for(const model of ['ford','hudson','packard','police']){
  const [min,max]=manifest[model].bounds_blender,size=trafficSize(model);
  assert.ok(size.width>=max[0]-min[0],`${model} width`);assert.ok(size.length>=max[1]-min[1],`${model} length`);
 }
});
