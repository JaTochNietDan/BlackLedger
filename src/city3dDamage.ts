import * as THREE from 'three';

type Wear = {amount: {value: number}; origin: {value: THREE.Vector3}; base: THREE.Color; cutaway: {value: number}; window: {value: THREE.Vector3}};
/** Condition describes damage, not its cause; stains imply no ongoing fire. */
export function buildingCondition(material: THREE.MeshStandardMaterial, condition: number, origin: THREE.Vector3) {
  let wear = material.userData.cityWear as Wear | undefined;
  if (!wear) {
    wear = {amount: {value: 0}, origin: {value: origin.clone()}, base: material.color.clone(), cutaway: {value: 0}, window: {value: new THREE.Vector3()}};
    material.userData.cityWear = wear;
    const uniforms = wear;
    material.onBeforeCompile = shader => {
      shader.uniforms.cityWear = uniforms.amount;
      shader.uniforms.cityCutaway = uniforms.cutaway;
      shader.uniforms.cityWindow = uniforms.window;
      shader.uniforms.cityWearOrigin = uniforms.origin;
      shader.vertexShader = 'varying vec3 citySurface;\n' + shader.vertexShader;
      shader.vertexShader = shader.vertexShader.replace('#include <worldpos_vertex>', `#include <worldpos_vertex>
        citySurface = (modelMatrix * vec4(transformed, 1.0)).xyz;`);
      shader.fragmentShader = `
        varying vec3 citySurface;
        uniform float cityWear;
        uniform float cityCutaway;
        uniform vec3 cityWindow;
        uniform vec3 cityWearOrigin;
        float cityHash(vec3 p) { return fract(sin(dot(p,vec3(127.1,311.7,74.7)))*43758.5453); }
        float cityNoise(vec3 p) {
          vec3 cell=floor(p), f=fract(p); f=f*f*(3.0-2.0*f);
          return mix(mix(mix(cityHash(cell),cityHash(cell+vec3(1,0,0)),f.x),
                         mix(cityHash(cell+vec3(0,1,0)),cityHash(cell+vec3(1,1,0)),f.x),f.y),
                     mix(mix(cityHash(cell+vec3(0,0,1)),cityHash(cell+vec3(1,0,1)),f.x),
                         mix(cityHash(cell+vec3(0,1,1)),cityHash(cell+vec3(1,1,1)),f.x),f.y),f.z);
        }
      ` + shader.fragmentShader;
      shader.fragmentShader = shader.fragmentShader.replace('#include <color_fragment>', `#include <color_fragment>
        if (cityCutaway > 0.001) {
          float opening = 1.0 - smoothstep(cityWindow.z * .7, cityWindow.z, distance(gl_FragCoord.xy, cityWindow.xy));
          float grain = fract(sin(dot(floor(gl_FragCoord.xy), vec2(12.9898,78.233))) * 43758.5453);
          if (grain < opening * cityCutaway * .98) discard;
        }
        if (cityWear > 0.0) {
        vec3 wearPoint=(citySurface-cityWearOrigin)*vec3(.8,.36,.8);
        float stain=cityNoise(wearPoint)*.72+cityNoise(wearPoint*3.1)*.28;
        float coverage=smoothstep(.68-cityWear*.35,.88-cityWear*.20,stain);
        diffuseColor.rgb *= 1.0-cityWear*(.16+coverage*.70);
        }`);
    };
    material.customProgramCacheKey = () => 'city-condition-cutaway-v2';
    material.needsUpdate = true;
  }
  wear.amount.value = 1 - Math.max(0, Math.min(100, Number.isFinite(condition) ? condition : 100)) / 100;
  wear.origin.value.copy(origin);
  material.color.copy(wear.base);
}

/** Private building materials share no visibility state with their source asset. */
export function buildingCutaway(material: THREE.MeshStandardMaterial, amount: number, window: THREE.Vector3) {
  const wear = material.userData.cityWear as Wear | undefined;
  if (!wear) return;
  wear.cutaway.value = Math.max(0, Math.min(1, amount));
  wear.window.value.copy(window);
}

/** Broken glazing remains a condition state until repairs restore the facade. */
export function buildingGlazing(building:THREE.Group,condition:number){
 const broken=Number.isFinite(condition)&&condition<60;
 const intact=building.getObjectByName('window-intact'),damaged=building.getObjectByName('window-broken');
 if(intact)intact.visible=!broken;if(damaged)damaged.visible=broken;
}

export const GLASS_BREAK_AT=.09;
export function glazingDuringBlast(condition:number,before:number,seconds:number,preview:boolean){
 return seconds<GLASS_BREAK_AT?before:preview?0:condition;
}
