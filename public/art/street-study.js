// Presentation-only study: no API commands, economy, outcomes, or authoritative clock.
const embedded=new URLSearchParams(location.search).get('embed')==='1';
if(embedded)document.body.classList.add('embedded');
const canvas=document.querySelector('#city'),ctx=canvas.getContext('2d'),light=document.querySelector('#light'),motion=document.querySelector('#motion'),status=document.querySelector('#status');
const reduced=matchMedia('(prefers-reduced-motion: reduce)');motion.checked=!reduced.matches;
const assets={};let t=0,last=0,frame=0,selected='',arrival=null;
const buildings=[{id:'cafe',name:'Saint Agnes',file:'cafe-noir-v1.png',x:70,y:35,w:630,h:420,depth:445,door:[446,443],info:'Your neighborhood contact keeps a table here. Building selection will connect to the existing action panel.'},{id:'casino',name:'The Monarch',file:'casino-noir-v1.png',x:490,y:230,w:630,h:420,depth:625,door:[900,625],info:'Bellandi territory. Its lights are an independent layer over the building image.'}];
const bulbs=[[794,699],[818,714],[849,727],[880,738],[913,747],[947,754],[982,757],[1017,755],[1051,750],[1084,742],[1115,730],[1144,716],[1168,701],[1188,683]];
function line(x1,y1,x2,y2,width,color){ctx.beginPath();ctx.moveTo(x1,y1);ctx.lineTo(x2,y2);ctx.lineWidth=width;ctx.strokeStyle=color;ctx.stroke()}
function ellipse(x,y,rx,ry,color){ctx.fillStyle=color;ctx.beginPath();ctx.ellipse(x,y,rx,ry,0,0,Math.PI*2);ctx.fill()}
function ground(n){ctx.fillStyle=n>.5?'#253333':'#536055';ctx.fillRect(0,0,1280,820);
 // Street and sidewalks share the building assets' projected ground axes.
 line(-200,230,1480,1070,190,'#656960');line(-200,230,1480,1070,142,'#383f3e');
 line(-200,230,1480,1070,1,'#858574');
 line(300,640,1300,140,150,'#656960');line(300,640,1300,140,110,'#3a4140');
 // Stable surface variation avoids animated noise.
 for(let i=0;i<480;i++){const x=(i*191)%1280,y=(i*137)%820;ctx.fillStyle=i%3?'#00000007':'#efddba0b';ctx.fillRect(x,y,2+(i%5),1)}
 for(let x=-100;x<1300;x+=44){line(x,.5*x+242,x-15,.5*x+257,1,'#292f2b55');line(x,.5*x+406,x-15,.5*x+420,1,'#292f2b55')}
 if(n>0){ctx.fillStyle=`rgba(9,18,29,${n*.28})`;ctx.fillRect(0,0,1280,820)}}
function building(b,n){ctx.save();ctx.filter=`brightness(${1-n*.48})`;ctx.drawImage(assets[b.id],b.x,b.y,b.w,b.h);ctx.restore();
 if(b.id==='casino'&&n>0){bulbs.forEach(([x,y],i)=>{const bright=(Math.floor(t*5)+i)%4===0?.25:1;ctx.save();ctx.globalAlpha=n*bright;ctx.shadowColor='#ffd387';ctx.shadowBlur=7;ellipse(b.x+x*b.w/1536,b.y+y*b.h/1024,1.25,1.25,'#ffe4a6');ctx.restore()});glow(b.door[0],b.door[1]-10,70,n*.3)}
 if(b.id==='cafe'&&n>0){glow(b.door[0],b.door[1]-24,38,n*.2)}
 if(selected===b.id){ctx.strokeStyle='#eac88b';ctx.lineWidth=2;ctx.beginPath();ctx.ellipse(...b.door,26,10,0,0,Math.PI*2);ctx.stroke()}}
function glow(x,y,r,a){const g=ctx.createRadialGradient(x,y,0,x,y,r);g.addColorStop(0,`rgba(255,193,91,${a})`);g.addColorStop(1,'rgba(255,193,91,0)');ctx.fillStyle=g;ctx.fillRect(x-r,y-r,r*2,r*2)}
function person(x,y,i,n){ctx.save();ellipse(x+2,y,5,2,'#00000035');const step=Math.sin(t*5+i)*1.8;line(x-1,y-6,x-2-step,y,1.7,'#191f22');line(x+1,y-6,x+2+step,y,1.7,'#191f22');ctx.fillStyle=['#514b40','#453533','#777468','#34454b'][i%4];ctx.fillRect(x-3,y-13,6,8);ellipse(x,y-16,2.3,3,'#b69e7c');ellipse(x,y-18,4,1.3,'#2a2b27');ctx.restore()}
function car(x,y,n){ctx.save();ellipse(x,y+9,31,10,'#00000035');ctx.filter=`brightness(${1-n*.35})`;ctx.drawImage(assets.car,x-42,y-35,84,56);ctx.restore();if(n>.1)glow(x+29,y+10,28,n*.32)}
function render(){const n=Number(light.value)/100;ctx.setTransform(canvas.width/1280,0,0,canvas.height/820,0,0);ground(n);const objects=buildings.map(b=>({depth:b.depth,draw:()=>building(b,n)}));
 for(let i=0;i<10;i++){const x=((i*123+t*(i%2?-7:9))+1500)%1500-100,y=.5*x+(i%3===0?415:249);objects.push({depth:y,draw:()=>person(x,y,i,n)})}
 for(let i=0;i<2;i++){let x=((t*40+i*700)%1560)-160,y=.5*x+320;objects.push({depth:y+16,draw:()=>car(x,y,n)})}
 if(arrival){const elapsed=t-arrival.start,x=Math.min(844,-120+elapsed*110),y=.5*x+320-40*Math.max(0,Math.min(1,(x-644)/200));objects.push({depth:y+16,draw:()=>car(x,y,n)});if(x===844&&!arrival.reported){arrival.reported=true;status.textContent='The car has arrived outside The Monarch. This is a presentation preview, not a simulated story outcome.'}}
 objects.sort((a,b)=>a.depth-b.depth).forEach(o=>o.draw());
 // Chimney smoke: a few translucent particles, not a physics simulation.
 if(motion.checked){for(let i=0;i<5;i++){const phase=(t*.25+i/5)%1;ellipse(337+phase*12,59-phase*34,3+phase*7,2+phase*5,`rgba(182,179,166,${(1-phase)*.13})`)}}
}
function tick(now){frame=0;if(document.hidden)return;const dt=last?Math.min(.05,(now-last)/1000):0;last=now;if(motion.checked)t+=dt;render();if(motion.checked)frame=requestAnimationFrame(tick)}
function wake(){if(frame)cancelAnimationFrame(frame);last=0;frame=requestAnimationFrame(tick)}
function resize(){canvas.width=Math.round(canvas.clientWidth*Math.min(devicePixelRatio,2));canvas.height=Math.round(canvas.clientWidth*820/1280*Math.min(devicePixelRatio,2));wake()}
function inspect(id){selected=id;status.textContent=buildings.find(b=>b.id===id).info;if(embedded)parent.postMessage({type:'blackledger:inspect',location:id==='cafe'?'bar':'club'},location.origin);wake()}
document.querySelector('#cafe').onclick=()=>inspect('cafe');document.querySelector('#casino').onclick=()=>inspect('casino');light.oninput=wake;motion.onchange=wake;document.querySelector('#arrival').onclick=()=>{if(reduced.matches){status.textContent='Arrival preview skipped because reduced motion is enabled.';return}motion.checked=true;arrival={start:t,reported:false};status.textContent='A car is approaching The Monarch. Preview only; no game time passes.';wake()};document.addEventListener('visibilitychange',wake);reduced.addEventListener('change',()=>{motion.checked=!reduced.matches;wake()});
canvas.addEventListener('click',e=>{const r=canvas.getBoundingClientRect(),x=(e.clientX-r.left)*1280/r.width,y=(e.clientY-r.top)*820/r.height;const b=[...buildings].reverse().find(b=>x>b.x+b.w*.18&&x<b.x+b.w*.82&&y>b.y+b.h*.2&&y<b.y+b.h*.97);if(b)inspect(b.id)});
try{await Promise.all([...buildings.map(b=>[b.id,b.file]),['car','sedan-noir-v1.png']].map(async([id,file])=>{const img=new Image();img.src='/art/'+file;await img.decode();assets[id]=img}));new ResizeObserver(resize).observe(canvas);resize()}catch(e){status.className='error';status.textContent='An art asset could not load. Reload this study to retry.'}

if(embedded){window.addEventListener('message',e=>{if(e.source!==parent||e.origin!==location.origin||e.data?.type!=='blackledger:presentation')return;const data=e.data;selected=data.selected==='bar'?'cafe':data.selected==='club'?'casino':'';if(Number.isFinite(data.minute)){const hour=(data.minute%1440)/60;light.value=String(Math.round((hour<6||hour>=20?1:hour<8?(8-hour)/2:hour>17?(hour-17)/3:0)*100))}if(typeof data.motion==='boolean')motion.checked=data.motion&&!reduced.matches;wake()});parent.postMessage({type:'blackledger:ready'},location.origin)}
