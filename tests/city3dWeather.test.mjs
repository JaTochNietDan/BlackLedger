import test from 'node:test';
import assert from 'node:assert/strict';
import {cityWeather,rainVertices} from '../.runtime/frontend-test/city3dWeather.js';

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
