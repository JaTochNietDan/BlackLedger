import * as THREE from 'three';

export type PoolRail=[number,number,number,number];
export function poolRails(width:number,length:number):PoolRail[]{
 const rails:PoolRail[]=[[.085,0,width-.085,0],[.085,length,width-.085,length]];
 for(const x of [0,width]){const outside=x===0?-.045:width+.045;rails.push([x,.085,x,length/2-.075],[x,length/2+.075,x,length-.085],[x,length/2-.075,outside,length/2-.065],[x,length/2+.075,outside,length/2+.065]);}
 for(const x of [0,width])for(const y of [0,length]){const sx=x===0?1:-1,sy=y===0?1:-1;rails.push([x+sx*.085,y,x+sx*.049,y-sy*.035],[x,y+sy*.085,x-sx*.035,y+sy*.049]);}
 return rails;
}
// The nose sits on the solver segment at ball-centre height. All material lies
// outside that line, unlike a box centred on it which intrudes into the cloth.
export function poolCushionGeometry([ax,ay,bx,by]:PoolRail,width:number,length:number,radius:number,height=.78){
 const distance=Math.hypot(bx-ax,by-ay);let nx=-(by-ay)/distance,ny=(bx-ax)/distance;
 if(nx*(width/2-(ax+bx)/2)+ny*(length/2-(ay+by)/2)>0){nx=-nx;ny=-ny;}
 const profile=[[0,radius],[.012,.052],[.059,.052],[.065,0],[.024,0]];
 const vertices:number[]=[];
 for(const [x,y] of [[ax,ay],[bx,by]])for(const [out,z] of profile)vertices.push(x+nx*out-width/2,height+z,length/2-y-ny*out);
 const indices:number[]=[];
 for(let i=1;i<4;i++){indices.push(0,i+1,i,5,5+i,6+i);}
 for(let i=0;i<5;i++){const j=(i+1)%5;indices.push(i,j,5+j,i,5+j,5+i);}
 // Normalize face winding for either segment direction/outward normal.
 const positions=new THREE.Float32BufferAttribute(vertices,3);let volume=0;
 for(let i=0;i<indices.length;i+=3){const a=new THREE.Vector3().fromBufferAttribute(positions,indices[i]),b=new THREE.Vector3().fromBufferAttribute(positions,indices[i+1]),c=new THREE.Vector3().fromBufferAttribute(positions,indices[i+2]);volume+=a.dot(b.cross(c));}
 if(volume<0)for(let i=0;i<indices.length;i+=3)[indices[i+1],indices[i+2]]=[indices[i+2],indices[i+1]];
 const indexed=new THREE.BufferGeometry();indexed.setAttribute('position',positions);indexed.setIndex(indices);
 const geometry=indexed.toNonIndexed();indexed.dispose();geometry.computeVertexNormals();return geometry;
}
