import test from 'node:test';
import assert from 'node:assert/strict';
import * as THREE from 'three';
import {disposeCityResources} from '../.runtime/frontend-test/city3dResources.js';

test('city teardown releases shared instances, shadow targets, textures and bitmaps once',()=>{
 const original=globalThis.ImageBitmap;
 class Bitmap {closed=0;close(){this.closed++;}}
 globalThis.ImageBitmap=Bitmap;
 try{
  const bitmap=new Bitmap(),texture=new THREE.Texture(bitmap),cube=new THREE.CubeTexture(Array(6).fill(bitmap));
  const geometry=new THREE.BoxGeometry(),material=new THREE.MeshStandardMaterial({map:texture,normalMap:texture,envMap:cube});
  const other=material.clone(),spriteMaterial=new THREE.SpriteMaterial({map:texture});
  const instance=new THREE.InstancedMesh(geometry,material,4);
  const scene=new THREE.Scene();scene.add(instance,new THREE.Mesh(geometry,other),new THREE.Sprite(spriteMaterial));
  const prototype=new THREE.Group();prototype.add(new THREE.Mesh(geometry,material));
  const light=new THREE.DirectionalLight();light.shadow.map=new THREE.WebGLRenderTarget(4,4);light.shadow.mapPass=new THREE.WebGLRenderTarget(4,4);scene.add(light);
  const watched=[instance,geometry,material,other,spriteMaterial,texture,cube,light.shadow.map,light.shadow.mapPass];
  const counts=new Map(watched.map(o=>[o,0]));
  for(const object of watched)object.addEventListener('dispose',()=>counts.set(object,counts.get(object)+1));
  disposeCityResources([scene,prototype],{geometries:[geometry],materials:[material],textures:[texture,cube]});
  for(const object of watched)assert.equal(counts.get(object),1,object.type||'shadow target');
  assert.equal(bitmap.closed,1,'shared decoded pixels must be closed only once');
 }finally{globalThis.ImageBitmap=original;}
});

test('a late-loading model can release resources without a renderer or browser bitmap implementation',()=>{
 const geometry=new THREE.BoxGeometry(),texture=new THREE.Texture(),material=new THREE.MeshBasicMaterial({map:texture});
 let disposed=0;for(const object of [geometry,texture,material])object.addEventListener('dispose',()=>disposed++);
 disposeCityResources([new THREE.Mesh(geometry,material)]);
 assert.equal(disposed,3);
});
