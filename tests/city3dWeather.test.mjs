import test from 'node:test';
import assert from 'node:assert/strict';
import {cityNightAmount,cityWeather,rainVertices} from '../.runtime/frontend-test/city3dWeather.js';

test('weather follows public rain and wetness independently, including drying streets',()=>{
 const dry=cityWeather({kind:'clear',wet:0},false),rain=cityWeather({kind:'rain',wet:1},false),drying=cityWeather({kind:'clear',wet:.45},false);
 assert.equal(dry.rain,false);assert.equal(rain.rain,true);assert.equal(drying.rain,false);
 assert.ok(rain.roadRoughness<drying.roadRoughness&&drying.roadRoughness<dry.roadRoughness);
 assert.ok(rain.roadTone<drying.roadTone&&drying.roadTone<dry.roadTone);
 assert.ok(rain.sun<dry.sun);
 assert.equal(cityWeather({kind:'fog',wet:.3},false).fogFar,440);
 assert.equal(cityWeather({kind:'clear',wet:NaN},false).wet,0);
 assert.equal(cityWeather(undefined,false).wet,0);
 assert.equal(cityWeather({kind:'rain',wet:2},true).wet,1);
});

test('rain fills a finite bounded reusable buffer and follows elapsed seconds at any frame rate',()=>{
 const vertices=new Float32Array(1800*6);rainVertices(vertices,0,192,160);const initial=vertices.slice();
 for(const seconds of [0,.1,2,10,10000]){
  rainVertices(vertices,seconds,192,160);
  for(let i=0;i<vertices.length;i+=6){
   const [x,y,z,tx,ty,tz]=vertices.subarray(i,i+6);
   assert.ok([x,y,z,tx,ty,tz].every(Number.isFinite));
   assert.ok(y>=.19&&ty<=32.21&&ty>=y&&ty-y<=1.41);
   assert.ok(x>=0&&x<=195&&z>=0&&z<=160);assert.equal(z,tz);
  }
 }
 for(const fps of [30,60,144]){
  let elapsed=0;for(let i=0;i<fps*2;i++)elapsed+=1/fps;
  rainVertices(vertices,elapsed,192,160);const expected=new Float32Array(vertices.length);rainVertices(expected,2,192,160);
  for(let i=0;i<vertices.length;i++)assert.ok(Math.abs(vertices[i]-expected[i])<1e-5);
 }
 assert.notDeepEqual(vertices,initial);
});


test('twilight is periodic and continuous at dawn, dusk and midnight',()=>{
 for(const [minute,amount] of [[0,1],[330,1],[360,.5],[390,0],[720,0],[1170,0],[1200,.5],[1230,1],[1440,1]]){
  assert.equal(cityNightAmount(minute),amount);
  assert.equal(cityNightAmount(minute+1440*7),amount);
  assert.equal(cityNightAmount(minute-1440),amount);
 }
 for(const minute of [0,330,360,390,1170,1200,1230,1440])
  assert.ok(Math.abs(cityNightAmount(minute-.001)-cityNightAmount(minute+.001))<.00006);
 for(const start of [330,1170]){
  let previous=cityNightAmount(start);
  for(let i=1;i<=3600;i++){
   const value=cityNightAmount(start+i/60);
   assert.ok(value>=0&&value<=1);
   assert.ok(Math.abs(value-previous)<.00042);
   assert.ok(start===330?value<=previous:value>=previous);
   previous=value;
  }
 }
});

test('all weather conditions blend their lighting without changing weather or paused time',()=>{
 for(const kind of ['clear','rain','fog','overcast']){
  const sky={kind,wet:.45},day=cityWeather(sky,false),night=cityWeather(sky,true);
  assert.deepEqual(cityWeather(sky,0),day);assert.deepEqual(cityWeather(sky,1),night);
  const middle=cityWeather(sky,cityNightAmount(1200));
  assert.ok(Math.abs(middle.sun-(day.sun+night.sun)/2)<1e-12);
  assert.ok(Math.abs(middle.ambient-(day.ambient+night.ambient)/2)<1e-12);
  for(const key of ['wet','rain','fogFar','roadRoughness','pavementRoughness','roadTone','pavementTone'])assert.equal(middle[key],day[key]);
  for(let frame=0;frame<144;frame++)assert.deepEqual(cityWeather(sky,cityNightAmount(1200)),middle);
  let previous=day;
  for(let minute=1171;minute<=1230;minute++){
   const next=cityWeather(sky,cityNightAmount(minute));
   assert.ok(next.sun<=previous.sun&&previous.sun-next.sun<.072);
   for(const shift of [0,8,16])assert.ok(Math.abs(((next.background>>shift)&255)-((previous.background>>shift)&255))<=3);
   previous=next;
  }
 }
});
