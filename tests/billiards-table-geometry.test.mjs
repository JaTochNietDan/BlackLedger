import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {poolRails,poolCushionGeometry} from '../.runtime/frontend-test/billiardsTableGeometry.js';
const W=1.27,L=2.54,R=.028575,H=.78;
test('every cushion face stays behind the physical contact line, including jaws',()=>{
 const rails=poolRails(W,L);assert.equal(rails.length,18);
 for(const rail of rails){const [ax,ay,bx,by]=rail,d=Math.hypot(bx-ax,by-ay);let nx=-(by-ay)/d,ny=(bx-ax)/d;if(nx*(W/2-(ax+bx)/2)+ny*(L/2-(ay+by)/2)>0){nx=-nx;ny=-ny;}
  const g=poolCushionGeometry(rail,W,L,R,H),p=g.getAttribute('position');let noses=0;
  for(let i=0;i<p.count;i++){const x=p.getX(i)+W/2,y=L/2-p.getZ(i),out=(x-ax)*nx+(y-ay)*ny;assert.ok(out>=-1e-6,'visual cushion intrudes onto cloth');if(Math.abs(out)<1e-6){assert.ok(Math.abs(p.getY(i)-H-R)<1e-6);noses++;}}
  assert.ok(noses>=2);g.dispose();
 }
});
test('rendered end cushion is met at the same tangent point as the solver',()=>{
 const g=poolCushionGeometry(poolRails(W,L)[0],W,L,R,H),mesh=new THREE.Mesh(g,new THREE.MeshBasicMaterial({side:THREE.DoubleSide}));mesh.updateMatrixWorld();
 const ray=new THREE.Raycaster(new THREE.Vector3(0,H+R,0),new THREE.Vector3(0,0,1));
 const hits=ray.intersectObject(mesh);assert.ok(hits.length);assert.ok(Math.abs(hits[0].point.z-L/2)<1e-6);
 g.dispose();mesh.material.dispose();
});
