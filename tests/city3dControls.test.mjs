import test from 'node:test';
import assert from 'node:assert/strict';
import {cameraCommand, screenPan} from '../.runtime/frontend-test/city3dControls.js';

test('camera shortcuts preserve browser modifiers and composing input', () => {
  for (const key of ['+', '=', '-', 'Q', 'E', 'Home', 'ArrowLeft', 'ArrowUp', 'w', 'a', 's', 'd']) {
    assert.ok(cameraCommand({key}));
    for (const modifier of ['ctrlKey', 'metaKey', 'altKey', 'isComposing'])
      assert.equal(cameraCommand({key, [modifier]: true}), null);
  }
  for (const key of ['Tab', 'Escape', 'Enter', 'h', 'constructor', '__proto__']) assert.equal(cameraCommand({key}), null);
  assert.equal(cameraCommand({key: 'Q'}), 'rotate-left');
  assert.equal(cameraCommand({key: '+'}), 'zoom-in');
});

test('panning stays aligned with the screen through a complete camera orbit', () => {
  const target = {x: 80, z: 64};
  for (let angle = 0; angle < Math.PI * 2; angle += .05) {
    const camera = {x: target.x + 100 * Math.cos(angle), z: target.z + 100 * Math.sin(angle)};
    const up = screenPan(camera, target, 'pan-up');
    const right = screenPan(camera, target, 'pan-right');
    const down = screenPan(camera, target, 'pan-down');
    const left = screenPan(camera, target, 'pan-left');
    assert.ok(Math.abs(Math.hypot(up.x, up.z) - 5) < 1e-10);
    assert.ok(Math.abs(Math.hypot(right.x, right.z) - 5) < 1e-10);
    assert.ok(Math.abs(up.x * right.x + up.z * right.z) < 1e-10);
    assert.ok(up.x * (target.x - camera.x) + up.z * (target.z - camera.z) > 0);
    assert.ok(up.x * right.z - up.z * right.x > 0);
    assert.equal(down.x, -up.x); assert.equal(down.z, -up.z);
    assert.equal(left.x, -right.x); assert.equal(left.z, -right.z);
  }
  assert.deepEqual(screenPan(target, target, 'pan-up'), {x: 0, z: 5});
});

test('WASD matches arrow panning with case-insensitive keys and preserves shortcuts',()=>{
 for(const [letter,arrow] of [['w','ArrowUp'],['a','ArrowLeft'],['s','ArrowDown'],['d','ArrowRight']]){
  assert.equal(cameraCommand({key:letter}),cameraCommand({key:arrow}));
  assert.equal(cameraCommand({key:letter.toUpperCase()}),cameraCommand({key:arrow}));
  assert.equal(cameraCommand({key:letter,metaKey:true}),null);
 }
});

test('held panning is frame-rate independent, normalized and stops on release', async()=>{
 const {KeyboardPan}=await import('../.runtime/frontend-test/city3dControls.js');
 const camera={x:0,z:-10},target={x:0,z:0};
 for(const fps of [30,60,144]){
  const pan=new KeyboardPan();pan.press({key:'w'});
  let z=0;for(let i=0;i<fps;i++){pan.press({key:'w'});z+=pan.step(camera,target,1/fps,10).z;}
  assert.ok(Math.abs(z-10)<1e-10);
  pan.press({key:'ArrowUp'});pan.release('w');
  assert.ok(pan.step(camera,target,.02,10).z>0);
  pan.press({key:'d'});
  const diagonal=pan.step(camera,target,.02,10);assert.ok(Math.abs(Math.hypot(diagonal.x,diagonal.z)-.2)<1e-10);
  pan.press({key:'s'});pan.press({key:'a'});assert.deepEqual(pan.step(camera,target,.02,10),{x:0,z:0});
  pan.clear();assert.deepEqual(pan.step(camera,target,.02,10),{x:0,z:0});
  pan.press({key:'W'});pan.release('w');assert.deepEqual(pan.step(camera,target,.02,10),{x:0,z:0});
  assert.equal(pan.press({key:'w',metaKey:true}),false);
  pan.press({key:'w'});assert.equal(pan.step(camera,target,100,10).z,.5);
 }
});
