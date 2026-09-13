import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {buildingCondition} from '../.runtime/frontend-test/city3dDamage.js';

test('condition changes update existing uniforms and repairs restore the original appearance inputs',()=>{
 const material=new THREE.MeshStandardMaterial({color:0xa57755}),base=material.color.clone();
 buildingCondition(material,38,new THREE.Vector3(16,.18,32));
 const shader={uniforms:{},vertexShader:THREE.ShaderLib.standard.vertexShader,fragmentShader:THREE.ShaderLib.standard.fragmentShader};
 material.onBeforeCompile(shader,{});
 const amount=shader.uniforms.cityWear,origin=shader.uniforms.cityWearOrigin,version=material.version;
 assert.equal(amount.value,.62);
 buildingCondition(material,100,new THREE.Vector3(48,.18,32));
 assert.equal(amount.value,0,'repairs must remove all staining');
 assert.deepEqual(origin.value.toArray(),[48,.18,32]);
 assert.equal(material.version,version,'revision updates must not rebuild the shader');
 assert.ok(material.color.equals(base),'texture/base colour must remain recoverable');
 buildingCondition(material,-10,new THREE.Vector3());assert.equal(amount.value,1);
 buildingCondition(material,110,new THREE.Vector3());assert.equal(amount.value,0);
 buildingCondition(material,NaN,new THREE.Vector3());assert.equal(amount.value,0);
});

test('damaged clones do not stain healthy buildings sharing the same source texture',()=>{
 const texture=new THREE.Texture(),source=new THREE.MeshStandardMaterial({map:texture});
 const damaged=source.clone(),healthy=source.clone();
 buildingCondition(damaged,0,new THREE.Vector3());buildingCondition(healthy,100,new THREE.Vector3(32,0,0));
 assert.notEqual(damaged.userData.cityWear.amount,healthy.userData.cityWear.amount);
 assert.equal(damaged.userData.cityWear.amount.value,1);assert.equal(healthy.userData.cityWear.amount.value,0);
 assert.equal(damaged.map,healthy.map);assert.equal(source.userData.cityWear,undefined);
 assert.equal(damaged.customProgramCacheKey(),healthy.customProgramCacheKey(),'material instances can share a program');
});

test('broken glazing follows condition and restores intact panes after repairs',async()=>{
 const {buildingGlazing}=await import('../.runtime/frontend-test/city3dDamage.js');
 const b=new THREE.Group(),intact=new THREE.Group(),broken=new THREE.Group();intact.name='window-intact';broken.name='window-broken';b.add(intact,broken);
 for(const condition of [100,60,NaN]){buildingGlazing(b,condition);assert.equal(intact.visible,true);assert.equal(broken.visible,false);}
 for(const condition of [59,38,0]){buildingGlazing(b,condition);assert.equal(intact.visible,false);assert.equal(broken.visible,true);}
 buildingGlazing(b,80);assert.equal(intact.visible,true);assert.equal(broken.visible,false);
});
