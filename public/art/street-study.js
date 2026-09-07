// Presentation-only study: no API commands, economy, outcomes, or authoritative clock.
const embedded=new URLSearchParams(location.search).get('embed')==='1';
if(embedded)document.body.classList.add('embedded');
const canvas=document.querySelector('#city'),ctx=canvas.getContext('2d'),light=document.querySelector('#light'),motion=document.querySelector('#motion'),status=document.querySelector('#status');
const reduced=matchMedia('(prefers-reduced-motion: reduce)');motion.checked=!reduced.matches;
const assets={};let t=0,last=0,frame=0,selected='',arrival=null,playerLocation='',journey=null,journeyKey='',propertyState={};
const buildings=await fetch('/art/buildings.json').then(r=>{if(!r.ok)throw Error('Building manifest unavailable');return r.json()});
const bulbs=[[794,699],[818,714],[849,727],[880,738],[913,747],[947,754],[982,757],[1017,755],[1051,750],[1084,742],[1115,730],[1144,716],[1168,701],[1188,683]];
function line(x1,y1,x2,y2,width,color){ctx.beginPath();ctx.moveTo(x1,y1);ctx.lineTo(x2,y2);ctx.lineWidth=width;ctx.strokeStyle=color;ctx.stroke()}
function ellipse(x,y,rx,ry,color){ctx.fillStyle=color;ctx.beginPath();ctx.ellipse(x,y,rx,ry,0,0,Math.PI*2);ctx.fill()}
let asphalt=null,groundCache=null,groundLight=-1;
function ground(n){
 if(!groundCache||groundLight!==n){
  const layer=document.createElement('canvas');layer.width=1280;layer.height=820;
  const g=layer.getContext('2d');
  g.fillStyle='#444b43';g.fillRect(0,0,1280,820);
  if(asphalt){g.globalAlpha=.16;g.fillStyle=asphalt;g.fillRect(0,0,1280,820);g.globalAlpha=1}
  const stroke=(a,b,w,color)=>{g.beginPath();g.moveTo(...a);g.lineTo(...b);g.lineWidth=w;g.strokeStyle=color;g.stroke()};
  const streets=[{a:[-240,210],b:[1520,1090],sidewalk:190,road:142},{a:[-240,910],b:[1520,30],sidewalk:150,road:110}];
  // Both pavements are drawn before either roadway, keeping the junction open.
  for(const road of streets)stroke(road.a,road.b,road.sidewalk,'#777668');
  for(const road of streets){stroke(road.a,road.b,road.road+4,'#2a302c');stroke(road.a,road.b,road.road,asphalt||'#383f3e')}
  // Restrained sidewalk joints, outside the open junction.
  for(let x=-100;x<1300;x+=32){if(x>330&&x<640)continue;stroke([x,.5*x+243],[x-13,.5*x+257],1,'#30352c55');stroke([x,.5*x+403],[x-13,.5*x+417],1,'#30352c55')}
  for(const [x,y] of [[265,462],[720,690],[965,310]]){g.save();g.translate(x,y);g.scale(1,.5);g.fillStyle='#313531';g.beginPath();g.arc(0,0,9,0,Math.PI*2);g.fill();g.strokeStyle='#7d7c6755';g.lineWidth=1;g.stroke();for(let k=-4;k<=4;k+=4)stroke([-5,k],[5,k],1,'#77786866');g.restore()}
  // Ground receives the same gradual night tint as the architecture.
  if(n>0){g.fillStyle=`rgba(9,18,29,${n*.55})`;g.fillRect(0,0,1280,820)}
  groundCache=layer;groundLight=n;
 }
 ctx.drawImage(groundCache,0,0);
}
function building(b,n){ctx.save();ctx.filter=`brightness(${1-n*.48})`;ctx.drawImage(b.damage&&(propertyState[b.id]?.condition??100)<b.damage.below?assets[b.id+':damaged']:assets[b.id],b.x,b.y,b.w,b.h);ctx.restore();
 if(b.id==='club'&&n>0){bulbs.forEach(([x,y],i)=>{const bright=(Math.floor(t*5)+i)%4===0?.25:1;ctx.save();ctx.globalAlpha=n*bright;ctx.shadowColor='#ffd387';ctx.shadowBlur=7;ellipse(b.x+x*b.w/1536,b.y+y*b.h/1024,1.25,1.25,'#ffe4a6');ctx.restore()});glow(b.door[0],b.door[1]-10,70,n*.3)}
 if(b.id==='bar'&&n>0){glow(b.door[0],b.door[1]-24,38,n*.2)}
 if(selected===b.id){ctx.strokeStyle='#eac88b';ctx.lineWidth=2;ctx.beginPath();ctx.ellipse(...b.door,26,10,0,0,Math.PI*2);ctx.stroke()}}
function glow(x,y,r,a){const g=ctx.createRadialGradient(x,y,0,x,y,r);g.addColorStop(0,`rgba(255,193,91,${a})`);g.addColorStop(1,'rgba(255,193,91,0)');ctx.fillStyle=g;ctx.fillRect(x-r,y-r,r*2,r*2)}
function person(x,y,i,n){ctx.save();ctx.translate(x,y);ctx.scale(1.5,1.5);x=0;y=0;ellipse(x+2,y,5,2,'#00000035');const step=Math.sin(t*5+i)*1.8;line(x-1,y-6,x-2-step,y,1.7,'#191f22');line(x+1,y-6,x+2+step,y,1.7,'#191f22');ctx.fillStyle=['#514b40','#453533','#777468','#34454b'][i%4];ctx.fillRect(x-3,y-13,6,8);ellipse(x,y-16,2.3,3,'#b69e7c');ellipse(x,y-18,4,1.3,'#2a2b27');ctx.restore()}
function car(x,y,n){ctx.save();ellipse(x,y+9,31,10,'#00000035');ctx.filter=`brightness(${1-n*.35})`;ctx.drawImage(assets.car,x-42,y-35,84,56);ctx.restore();if(n>.1)glow(x+29,y+10,28,n*.32)}
// Authored pavement routes are presentation metadata, not simulation pathfinding.
function streetRoute(from,to){
 const a=from?.walkway||[[70,745],[530,514]],b=to.walkway||[to.door];
 return [...a,...[...b].reverse().slice(1)];
}
function routePoint(points,progress){
 const lengths=points.slice(1).map((p,i)=>Math.hypot(p[0]-points[i][0],p[1]-points[i][1]));
 let distance=lengths.reduce((a,b)=>a+b,0)*Math.max(0,Math.min(1,progress));
 for(let i=0;i<lengths.length;i++){if(distance<=lengths[i]){const f=lengths[i]?distance/lengths[i]:0;return [points[i][0]+(points[i+1][0]-points[i][0])*f,points[i][1]+(points[i+1][1]-points[i][1])*f]}distance-=lengths[i]}
 return points[points.length-1];
}
function render(){const n=Number(light.value)/100;ctx.setTransform(canvas.width/1280,0,0,canvas.height/820,0,0);ground(n);const objects=buildings.map(b=>({depth:b.depth,draw:()=>building(b,n)}));
 for(let i=0;i<10;i++){const x=((i*123+t*(i%2?-7:9))+1500)%1500-100,y=.5*x+(i%3===0?415:249);objects.push({depth:y,draw:()=>person(x,y,i,n)})}
 for(let i=0;i<2;i++){let x=((t*40+i*700)%1560)-160,y=.5*x+320;objects.push({depth:y+16,draw:()=>car(x,y,n)})}
 if(arrival){const elapsed=t-arrival.start,x=Math.min(844,-120+elapsed*110),y=.5*x+320-40*Math.max(0,Math.min(1,(x-644)/200));objects.push({depth:y+16,draw:()=>car(x,y,n)});if(x===844&&!arrival.reported){arrival.reported=true;status.textContent='The car has arrived outside The Monarch. This is a presentation preview, not a simulated story outcome.'}}
 let marker=buildings.find(b=>b.id===playerLocation)?.door;
 if(journey){const elapsed=Math.min(1,(performance.now()-journey.start)/2200);marker=routePoint(journey.route,elapsed);objects.push({depth:marker[1]+1,draw:()=>person(marker[0],marker[1],1,n)})}
 objects.sort((a,b)=>a.depth-b.depth).forEach(o=>o.draw());
 if(marker){ctx.strokeStyle='#efcf8b';ctx.lineWidth=2;ctx.beginPath();ctx.ellipse(marker[0],marker[1]+3,11,4,0,0,Math.PI*2);ctx.stroke();if(!journey){ctx.fillStyle='#f1d5a0';ctx.font='11px system-ui';ctx.textAlign='center';ctx.fillText('YOU',marker[0],marker[1]+20)}}
 // Chimney smoke: a few translucent particles, not a physics simulation.
 if(motion.checked){for(let i=0;i<5;i++){const phase=(t*.25+i/5)%1;ellipse(337+phase*12,59-phase*34,3+phase*7,2+phase*5,`rgba(182,179,166,${(1-phase)*.13})`)}}
}
function tick(now){frame=0;if(document.hidden)return;const dt=last?Math.min(.05,(now-last)/1000):0;last=now;if(motion.checked)t+=dt;render();if(motion.checked||journey)frame=requestAnimationFrame(tick)}
function wake(){if(frame)cancelAnimationFrame(frame);last=0;frame=requestAnimationFrame(tick)}
function resize(){canvas.width=Math.round(canvas.clientWidth*Math.min(devicePixelRatio,2));canvas.height=Math.round(canvas.clientWidth*820/1280*Math.min(devicePixelRatio,2));wake()}
function inspect(id){selected=id;status.textContent=buildings.find(b=>b.id===id).info;if(embedded)parent.postMessage({type:'blackledger:inspect',location:id},location.origin);wake()}
document.querySelector('#cafe').onclick=()=>inspect('bar');document.querySelector('#casino').onclick=()=>inspect('club');document.querySelector('#laundry').onclick=()=>inspect('laundry');document.querySelector('#hotel').onclick=()=>inspect('room');light.oninput=wake;motion.onchange=wake;document.querySelector('#arrival').onclick=()=>{if(reduced.matches){status.textContent='Arrival preview skipped because reduced motion is enabled.';return}motion.checked=true;arrival={start:t,reported:false};status.textContent='A car is approaching The Monarch. Preview only; no game time passes.';wake()};document.addEventListener('visibilitychange',wake);reduced.addEventListener('change',()=>{motion.checked=!reduced.matches;wake()});
canvas.addEventListener('click',e=>{const r=canvas.getBoundingClientRect(),x=(e.clientX-r.left)*1280/r.width,y=(e.clientY-r.top)*820/r.height;const b=[...buildings].reverse().find(b=>x>b.x+b.w*.18&&x<b.x+b.w*.82&&y>b.y+b.h*.2&&y<b.y+b.h*.97);if(b)inspect(b.id)});
try{await Promise.all([...buildings.flatMap(b=>[[b.id,b.file],...(b.damage?[[b.id+':damaged',b.damage.file]]:[])]),['car','sedan-noir-v1.png'],['asphalt','asphalt-noir-v1.png']].map(async([id,file])=>{const img=new Image();img.src='/art/'+file;await img.decode();assets[id]=img}));// Reuse each original sprite's alpha as the damage layer's rendering mask.
for(const b of buildings){if(!b.damage)continue;const original=assets[b.id],layer=document.createElement('canvas');layer.width=original.width;layer.height=original.height;const surface=layer.getContext('2d');surface.drawImage(assets[b.id+':damaged'],0,0,layer.width,layer.height);surface.globalCompositeOperation='destination-in';surface.drawImage(original,0,0);assets[b.id+':damaged']=layer}
asphalt=ctx.createPattern(assets.asphalt,'repeat');asphalt.setTransform(new DOMMatrix([.28,.14,-.28,.14,0,0]));new ResizeObserver(resize).observe(canvas);resize()}catch(e){status.className='error';status.textContent='An art asset could not load. Reload this study to retry.'}

if(embedded){window.addEventListener('message',e=>{if(e.source!==parent||e.origin!==location.origin||e.data?.type!=='blackledger:presentation')return;const data=e.data;if(Array.isArray(data.properties)){propertyState={};for(const p of data.properties){if(buildings.some(b=>b.id===p.id)&&Number.isFinite(p.condition))propertyState[p.id]={condition:Math.max(0,Math.min(100,p.condition))}}}playerLocation=typeof data.position==='string'?data.position:'';const key=data.journey?data.journey.from+'>'+data.journey.to:'';if(key!==journeyKey){journeyKey=key;const from=buildings.find(b=>b.id===data.journey?.from),to=buildings.find(b=>b.id===data.journey?.to);journey=key&&to&&!reduced.matches?{route:streetRoute(from,to),start:performance.now()}:null}selected=buildings.some(b=>b.id===data.selected)?data.selected:'';if(Number.isFinite(data.minute)){const hour=(data.minute%1440)/60;light.value=String(Math.round((hour<6||hour>=20?1:hour<8?(8-hour)/2:hour>17?(hour-17)/3:0)*100))}if(typeof data.motion==='boolean')motion.checked=data.motion&&!reduced.matches;wake()});parent.postMessage({type:'blackledger:ready'},location.origin)}
