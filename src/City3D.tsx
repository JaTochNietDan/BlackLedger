import {CityBuildingDriveBy} from './city3dBuildingDriveBy';
import {CityVillaExit} from './city3dVillaExit';
import {CityAccident} from './city3dAccident';
import {CityPlanter} from './city3dPlanter';
import {CityIncendiary, incendiaryStagingFlight, incendiaryShard, incendiaryShardObstructed, INCENDIARY_IMPACT} from './city3dIncendiary';
import {StreetPlayback,streetLegKey,responseRecords,advanceJourneyClock} from './streetPlayback';
import {CityCustody} from './city3dCustody';
import {CityAssassination, assassinationBatch, isStagedStrike, MELEE_IMPACTS, ASSASSINATION_SHOT, ASSASSINATION_VICTIM_X, executionSpatter} from './city3dAssassination';
import {poseCustody,sceneWeapon,poseLongGun,weaponShots,pumpOffset} from './city3dWeapons';
import {CityRubble} from './city3dRubble';
import {CitySuppression} from './city3dSuppression';
import {CityFire, clearBlastWindows} from './city3dFire';
import {CityAftermath,captureBodyJoints} from './city3dAftermath';
import {groundCharacter} from './city3dGround';
import {previewScenes, previewScene, type PreviewScene} from './city3dPreview';
import {frameScene, stagedSceneBounds,ScenePullback, impactPulse, renderImpact} from './city3dFraming';
import {seatDriver} from './city3dSeating';
import {wardrobe, dressPedestrian} from './city3dWardrobe';
import {headlightAlpha, headlightCentre} from './city3dHeadlights';
import {cityNightAmount, cityWeather, rainVertices} from './city3dWeather';
import {blockingBuildings, characterSightPoints} from './city3dOcclusion';
import {buildingCondition, buildingGlazing, glazingDuringBlast, GLASS_BREAK_AT, buildingCutaway} from './city3dDamage';
import {disposeCityResources} from './city3dResources';
import {useEffect, useRef, useState, type ReactNode} from 'react';
import * as THREE from 'three';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import type {Snapshot, VisualCue} from './types';
import {
  cityPlan,
  crossingPaint,
  entrance,
  route,
  onRoute,
  PITCH,
  STREET_WIDTH,
  streetsidePosition,
  parkingSpot,
  lampPositions,
  vehicleRootHeight,
  pedestrianRootHeight,
  surfaceHeight,
} from './city3dPlan';
import type {Lot, Point} from './city3dPlan';
import type {Journey} from './TravelPresentation';
import './city3d.css';
import {CityCueQueue, gunVictim, gunCastReady, raidEntryPose, policeSceneSeconds, officerApproach, policeCast, sceneSlots, availableSceneSlot, planterReservation, casualtyFall, gunfightPose, casualtySceneStart, GunfireAudio, BlastAudio} from './city3dEvents';
import type {SceneSlot} from './city3dEvents';
import {StreetTraffic, trafficSpeed, trafficSize, trafficModel, advanceWheel, wheelSteering, advanceSteering, frontWheelSteering} from './city3dTraffic';
import {pedestrianModel, isPedestrian} from './city3dCast';
import {playCityGunshot, preloadCityGunshots, cityGunshotStatus, preloadCityEffects, CityAmbientAudio, cityEffectStatus, playRecordedEffect, playMoment, soundOn} from './sound';
import {cameraCommand, KeyboardPan, bindKeyboardPan} from './city3dControls';
import {blastParticle, windowBurst, internalDetonation, windowDebris, blastLight, blastOpacity, billowAlpha, debrisPose, fragmentBlocked} from './city3dBlast';

const pathLength=(points:Point[])=>points.slice(1).reduce((sum,p,i)=>sum+Math.hypot(p.x-points[i].x,p.z-points[i].z),0);
type Props = {
  immersive?:boolean;
  state: Snapshot;
  beforeConditions?:Record<string,number>;
  replaySerial?:number;
  overlay?: ReactNode;
  selected: string;
  onSelect: (id: string) => void;
  onTravel: (id: string) => void;
  onEnter: () => void;
  motion: boolean;
  journey: Journey | null;
  busy: boolean;
  activeCue: VisualCue | null;
  onJourneyDone: () => void;
  onJourneyProgress?: (progress: number) => void;
  onJourneyBlocked?: (blocked: boolean) => void;
  onSceneDone?: (id: string) => void;
};
type Actor = {
  object: THREE.Group;
  model: string;
  points: Point[];
  start: number;
  end: number;
  since: number;
  duration: number;
  walking: boolean;
  phase: number;
  realSince: number;
  limbs: THREE.Object3D[];
  wheels: THREE.Object3D[];
  wheelPhase: number;
  steering: number;
  wheelPlaced: boolean;
  lamps: THREE.MeshStandardMaterial[];
  wardrobe: THREE.MeshStandardMaterial[];
  wardrobeKey: string;
  driver?: THREE.Group;
  arrived?: boolean;
  timelineProgress?:number;
  legKey?:string;
};
type Effect = {
  cue: VisualCue;
  since: number;
  mesh: THREE.InstancedMesh;
  light: THREE.PointLight;
  debris?: THREE.InstancedMesh;
  extra?: THREE.Group;
  wardrobe?: THREE.MeshStandardMaterial[];
  slot?: SceneSlot;
  glassBlocked?: Set<number>;
  gunArm?: THREE.Object3D;
  weapon?:THREE.Group;
  weaponModel?:string;
  muzzle?: THREE.Object3D;
  audio?: GunfireAudio | BlastAudio;
  glassAudio?: BlastAudio;
  reactionAudio?:BlastAudio;
  departureAudio?:BlastAudio;
  glazingBefore?:number;
  assassination?:CityAssassination;
  driveBy?:CityBuildingDriveBy;
  custody?:CityCustody;
  incendiary?:CityIncendiary;
  planter?:CityPlanter|CityVillaExit;
  accident?:CityAccident;
  pullback?:{move:ScenePullback;intent:number};
};
const modelNames = [
  'tenement',
  'mercer-court',
  'mariner',
  'tavern',
  'casino',
  'warehouse',
  'civic',
  'shop',
  'villa',
  'fire-nozzle',
  'firefighter',
  'fire-engine',
  'ford',
  'hudson',
  'packard',
  'person',
  'woman',
  'revolver',
  'handcuffs',
  'incendiary-bottle',
  'bottle-shard',
  'shotgun',
  'thompson',
  'blast-fragment',
  'street-bed',
  'vacant-lot',
  'filling',
  'garage',
  'dealer',
  'undertaker',
  'docks',
  'harbour-pier',
  'quay-section',
  'haulage',
  'police',
  'police-officer',
  'streetside',
  'monarch',
  'bluehour',
  'goldenlily',
  'papermoon',
];
const carModel = (name = '') =>
  /packard/i.test(name) ? 'packard' : /hudson/i.test(name) ? 'hudson' : 'ford';
const tmp = new THREE.Object3D();

export function City3D(props: Props) {
  const host = useRef<HTMLDivElement>(null);
  const expandButton = useRef<HTMLButtonElement>(null);
  const latest = useRef(props);
  const debug = new URLSearchParams(location.search).has('city-debug');
  const [preview, setPreview] = useState<ReturnType<typeof previewScene> | null>(null);
  const [previewKind, setPreviewKind] = useState<PreviewScene>('Gunfight');
  const canPreview = debug && !props.busy && !props.journey && !props.activeCue && props.motion;
  const shownPreview = canPreview ? preview : null;
  latest.current = shownPreview ? {...props, state: shownPreview.state, activeCue: shownPreview.cue} : props;
  useEffect(() => { setPreview(null); }, [props.state.id, props.state.revision, props.activeCue?.id]);
  const controlsRef = useRef<OrbitControls | null>(null);
  const focus = useRef<(id?: string) => void>(() => {});
  const [expanded, setExpanded] = useState(false);
  const [following, setFollowing] = useState(true);
  const followPlayer = useRef(true);
  const setFollow = (value: boolean) => { followPlayer.current = value; setFollowing(value); };
  const [status, setStatus] = useState('Loading Bellwether…');
  const [failure, setFailure] = useState('');
  const [fps, setFps] = useState('');
  const [waiting, setWaiting] = useState(0);
  const [playbackRate, setPlaybackRate] = useState(1);
  const playback = useRef(playbackRate);
  playback.current = playbackRate;
  useEffect(() => {
    if (failure && props.journey) props.onJourneyDone();
  }, [failure, props.journey, props.onJourneyDone]);
  useEffect(() => {
    const element = host.current!;
    let dead = false,
      frame = 0,
      renderer: THREE.WebGLRenderer;
    try {
      renderer = new THREE.WebGLRenderer({
        antialias: true,
        alpha: false,
        powerPreference: 'high-performance',
      });
    } catch {
      setFailure('The city needs WebGL. Your browser could not start the 3D renderer.');
      return;
    }
    renderer.setPixelRatio(Math.min(devicePixelRatio, 1.5));
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFShadowMap;
    renderer.shadowMap.autoUpdate = false;
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.ACESFilmicToneMapping;
    renderer.toneMappingExposure = 1.25;
    const canvas = renderer.domElement;
    canvas.setAttribute(
      'aria-label',
      '3D Bellwether city. Drag to rotate, right drag to pan, scroll to zoom. Keyboard: WASD or arrows pan, Q and E rotate, plus and minus zoom, Home resets, Escape leaves the expanded city.',
    );
    canvas.tabIndex = 0;
    element.appendChild(canvas);
    const scene = new THREE.Scene();
    scene.background = new THREE.Color('#333e40');
    const camera = new THREE.OrthographicCamera(-100, 100, 80, -80, 0.1, 1500);
    const controls = new OrbitControls(camera, canvas);
    controlsRef.current = controls;
    controls.enableDamping = true;
    controls.dampingFactor = 0.12;
    controls.minPolarAngle = 0.28;
    controls.maxPolarAngle = Math.PI * 0.44;
    controls.minZoom = 0.6;
    controls.maxZoom = 32;
    controls.screenSpacePanning = false;
    controls.mouseButtons = {
      LEFT: THREE.MOUSE.ROTATE,
      MIDDLE: THREE.MOUSE.DOLLY,
      RIGHT: THREE.MOUSE.PAN,
    };
    const plan = cityPlan(latest.current.state.locations);
    const lots = new Map(plan.lots.map(l => [l.id, l]));
    const home = new THREE.Vector3(plan.width / 2, 0, plan.depth / 2);
    let cameraIntent=0;
    const reset = () => {
      cameraIntent++;
      controls.target.copy(home);
      camera.position.copy(home).add(new THREE.Vector3(180, 200, -240));
      camera.zoom = 1;
      camera.updateProjectionMatrix();
      controls.update();
    };
    focus.current = id => {
      cameraIntent++;
      setFollow(false);
      const lot = id ? lots.get(id) : undefined;
      if (!lot) {
        reset();
        return;
      }
      const building = buildings.get(lot.id);
      const envelope = building?.userData.sightBounds as THREE.Box3 | undefined;
      frameScene(camera, controls.target, envelope?.clone().expandByScalar(1) ??
        new THREE.Box3(new THREE.Vector3(lot.x - 9, 0, lot.z - 9), new THREE.Vector3(lot.x + 9, 22, lot.z + 9)));
      controls.update();
    };
    reset();
    camera.zoom = 1.65;
    camera.updateProjectionMatrix();
    const resize = () => {
      const w = element.clientWidth,
        h = element.clientHeight;
      renderer.setSize(w, h);
      const span = Math.max(plan.width, plan.depth) * 0.6;
      camera.left = (-span * w) / h;
      camera.right = (span * w) / h;
      camera.top = span;
      camera.bottom = -span;
      camera.updateProjectionMatrix();
    };
    const observer = new ResizeObserver(resize);
    observer.observe(element);
    resize();
    const sky = new THREE.HemisphereLight(0xc5d8e0, 0x605342, 2.1);
    scene.add(sky);
    const sun = new THREE.DirectionalLight(0xffdbad, 3.2);
    sun.position.set(-80, 160, 90);
    sun.target.position.copy(home);
    scene.add(sun, sun.target);
    sun.castShadow = true;
    sun.shadow.mapSize.set(2048, 2048);
    sun.shadow.bias = -0.0002;
    sun.shadow.normalBias = 0.16;
    Object.assign(sun.shadow.camera, {
      left: -220,
      right: 220,
      top: 220,
      bottom: -220,
      near: 1,
      far: 600,
    });
    sun.shadow.camera.updateProjectionMatrix();
    const groundMat = new THREE.MeshStandardMaterial({color: 0x646460, roughness: 0.92});
    const pavementMat = new THREE.MeshStandardMaterial({color: 0xaaa18b, roughness: 0.9});
    const geo = new THREE.BoxGeometry(1, 1, 1);
    const ground = new THREE.Mesh(geo, groundMat);
    ground.scale.set(plan.width + 34, 0.6, plan.depth + 34);
    ground.position.set(plan.width / 2, -0.4, plan.depth / 2);
    ground.receiveShadow = true;
    scene.add(ground);
    const rainGeometry = new THREE.BufferGeometry();
    const rainPositions = new Float32Array(1800 * 6);
    const rainAttribute = new THREE.BufferAttribute(rainPositions, 3).setUsage(THREE.DynamicDrawUsage);
    rainGeometry.setAttribute('position', rainAttribute);
    const rainMaterial = new THREE.LineBasicMaterial({color: 0xb9c8cf, transparent: true, opacity: .22, depthWrite: false});
    const rainfall = new THREE.LineSegments(rainGeometry, rainMaterial);
    rainfall.frustumCulled = false; rainfall.visible = false; scene.add(rainfall);
    let rainClock = 0;
    const textures: THREE.Texture[] = [];
    const harbourLot = plan.lots.find(lot => lot.id === 'docks' && lot.col === 0);
    let waterNormal: THREE.CanvasTexture | undefined, waterClock = 0;
    if (harbourLot) {
      const waterCanvas = document.createElement('canvas'); waterCanvas.width = waterCanvas.height = 128;
      const context = waterCanvas.getContext('2d')!, pixels = context.createImageData(128, 128);
      for (let y = 0; y < 128; y++) for (let x = 0; x < 128; x++) {
        const u = x * Math.PI * 2 / 128, v = y * Math.PI * 2 / 128;
        const normal = new THREE.Vector3(.35 * Math.cos(u * 3 + v * 2) + .12 * Math.cos(u * 7 - v * 5),
          .23 * Math.cos(u * 3 + v * 2) - .09 * Math.cos(u * 7 - v * 5), 1).normalize();
        pixels.data.set([Math.round((normal.x * .5 + .5) * 255), Math.round((normal.y * .5 + .5) * 255), Math.round((normal.z * .5 + .5) * 255), 255], (y * 128 + x) * 4);
      }
      context.putImageData(pixels, 0, 0); waterNormal = new THREE.CanvasTexture(waterCanvas);
      waterNormal.wrapS = waterNormal.wrapT = THREE.RepeatWrapping;
      waterNormal.repeat.set(100, (plan.depth + 1600) / 10); textures.push(waterNormal);
      const water = new THREE.Mesh(new THREE.PlaneGeometry(1000, plan.depth + 1600),
        new THREE.MeshStandardMaterial({color: 0x344e50, metalness: .28, roughness: .34,
          normalMap: waterNormal, normalScale: new THREE.Vector2(.65, .65)}));
      water.rotation.x = -Math.PI / 2; water.position.set(-517, -.9, plan.depth / 2); water.receiveShadow = true;
      scene.add(water);
      const promenadeGeometry = new THREE.BoxGeometry(13, .28, plan.depth + 40);
      const positions = promenadeGeometry.attributes.position, normals = promenadeGeometry.attributes.normal;
      const uv = promenadeGeometry.attributes.uv;
      for (let i = 0; i < positions.count; i++) {
        const x = positions.getX(i), y = positions.getY(i), z = positions.getZ(i);
        uv.setXY(i, (Math.abs(normals.getX(i)) > .5 ? z : x) / 24,
          (Math.abs(normals.getY(i)) > .5 ? z : y) / 24);
      }
      const promenade = new THREE.Mesh(promenadeGeometry, pavementMat);
      promenade.position.set(-10.5, .03, plan.depth / 2);
      promenade.receiveShadow = true; scene.add(promenade);
    }
    const textureLoader = new THREE.TextureLoader();
    for (const [path, mat, repeat] of [
      ['asphalt', groundMat, 40],
      ['pavement', pavementMat, 4],
    ] as const) {
      textureLoader.load(`/art/ground/${path}.jpg`, t => {
        if (dead) {
          t.dispose();
          return;
        }
        t.colorSpace = THREE.SRGBColorSpace;
        t.wrapS = t.wrapT = THREE.RepeatWrapping;
        t.repeat.set(repeat, repeat);
        t.anisotropy = Math.min(8, renderer.capabilities.getMaxAnisotropy());
        mat.map = t;
        mat.needsUpdate = true;
        textures.push(t);
      });
    }
    const islands = new THREE.InstancedMesh(geo, pavementMat, plan.cols * plan.rows);
    let n = 0;
    for (let c = 0; c < plan.cols; c++)
      for (let r = 0; r < plan.rows; r++) {
        tmp.position.set((c + 0.5) * PITCH, 0.03, (r + 0.5) * PITCH);
        tmp.scale.set(PITCH - STREET_WIDTH, 0.28, PITCH - STREET_WIDTH);
        tmp.rotation.set(0, 0, 0);
        tmp.updateMatrix();
        islands.setMatrixAt(n++, tmp.matrix);
      }
    islands.receiveShadow = true;
    scene.add(islands);
    const stripeMat = new THREE.MeshStandardMaterial({color: 0xc4b998, roughness: 1});
    const stripes: THREE.Matrix4[] = [];
    const stripe = (x: number, z: number, w: number, d: number) => {
      tmp.position.set(x, -0.085, z);
      tmp.scale.set(w, 0.015, d);
      tmp.updateMatrix();
      stripes.push(tmp.matrix.clone());
    };
    for (let c = 0; c <= plan.cols; c++)
      for (let z = 7; z < plan.depth; z += 8)
        if (z % PITCH > 6 && z % PITCH < PITCH - 6) stripe(c * PITCH, z, 0.16, 3);
    for (let r = 0; r <= plan.rows; r++)
      for (let x = 7; x < plan.width; x += 8)
        if (x % PITCH > 6 && x % PITCH < PITCH - 6) stripe(x, r * PITCH, 3, 0.16);
    for (const mark of crossingPaint(plan.cols, plan.rows)) stripe(mark.x, mark.z, mark.width, mark.depth);
    const lines = new THREE.InstancedMesh(geo, stripeMat, stripes.length);
    stripes.forEach((m, i) => lines.setMatrixAt(i, m));
    scene.add(lines);
    // Repeated street furniture in one draw, with lamps above pedestrian clearance.
    const iron = new THREE.MeshStandardMaterial({color: 0x29332f, metalness: 0.5, roughness: 0.55});
    const poles = new THREE.InstancedMesh(
      new THREE.CylinderGeometry(0.1, 0.16, 5, 8),
      iron,
      plan.lots.length * 2,
    );
    const bulbs = new THREE.InstancedMesh(
      new THREE.SphereGeometry(0.25, 8, 6),
      new THREE.MeshStandardMaterial({color: 0xf7cf85, emissive: 0xffa742, emissiveIntensity: 1}),
      plan.lots.length * 2,
    );
    n = 0;
    for (const l of plan.lots)
      for (const at of lampPositions(l)) {
        tmp.scale.set(1, 1, 1);
        tmp.position.set(at.x, 2.6, at.z);
        tmp.updateMatrix();
        poles.setMatrixAt(n, tmp.matrix);
        tmp.position.y = 5.2;
        tmp.updateMatrix();
        bulbs.setMatrixAt(n++, tmp.matrix);
      }
    scene.add(poles, bulbs);
    // Mark parcel corners without drawing a broad circle through the architecture.
    const selectionMaterial = new THREE.MeshBasicMaterial({color: 0xb99452,
      transparent: true, opacity: .88, depthWrite: false, toneMapped: false});
    const cornerGeometry = new THREE.BoxGeometry(2.25, .008, .09);
    const selection = new THREE.InstancedMesh(cornerGeometry, selectionMaterial, 8);
    const markerTransform = new THREE.Object3D();
    let markerIndex = 0;
    for (const x of [-1, 1]) for (const z of [-1, 1]) {
      markerTransform.rotation.y = 0;
      markerTransform.position.set(x * (9.7 - 1.125), 0, z * 9.7);
      markerTransform.updateMatrix(); selection.setMatrixAt(markerIndex++, markerTransform.matrix);
      markerTransform.rotation.y = Math.PI / 2;
      markerTransform.position.set(x * 9.7, 0, z * (9.7 - 1.125));
      markerTransform.updateMatrix(); selection.setMatrixAt(markerIndex++, markerTransform.matrix);
    }
    selection.instanceMatrix.needsUpdate = true;
    selection.position.y = .22; scene.add(selection);
    const playerRing = new THREE.Mesh(
      new THREE.RingGeometry(.92, 1, 48),
      new THREE.MeshBasicMaterial({color: 0xd0bb87, side: THREE.DoubleSide,
        transparent: true, opacity: .9, depthWrite: false, toneMapped: false}),
    );
    playerRing.rotation.order = 'YXZ';
    playerRing.rotation.x = -Math.PI / 2;
    scene.add(playerRing);
    const models = new Map<string, THREE.Group>();
    const buildings = new Map<string, THREE.Group>();
    const clockHands:THREE.Object3D[]=[];
    let lightingKey="", lastNightAmount=-1;
    const luminousMaterials=new Set<THREE.MeshStandardMaterial>();
    const landings: THREE.Group[] = [];
    const labels = new Map<string, THREE.Sprite>();
    let blockers = new Set<string>(), lastSightCheck = -Infinity;
    const cutawayAmounts = new Map<string, number>();
    const cutawayWindow = new THREE.Vector3();
    const streetCutawayMaterials=new Set<THREE.MeshStandardMaterial>();
    const actors = new Map<string, Actor>();
    const effects: Effect[] = [];
    let completedScene = "";
    const suppression=new CitySuppression();scene.add(suppression.root);
    const rubble=new CityRubble();scene.add(rubble.root);
    const aftermath = new CityAftermath();
    scene.add(aftermath.root);
    const disposeDebris = (effect: Effect) => {
      if (!effect.debris) return;
      scene.remove(effect.debris);
      effect.debris.geometry.dispose();
      (effect.debris.material as THREE.Material).dispose();
      effect.debris.dispose();
    };
    const cueQueue = new CityCueQueue();
    const traffic = new StreetTraffic();
    const effectGeometry = new THREE.PlaneGeometry(1, 1);
    const particleCanvas = document.createElement('canvas');
    particleCanvas.width = particleCanvas.height = 64;
    const particleContext = particleCanvas.getContext('2d')!;
    const gradient = particleContext.createRadialGradient(32, 32, 0, 32, 32, 32);
    gradient.addColorStop(0, 'rgba(255,255,255,1)');
    gradient.addColorStop(0.3, 'rgba(255,255,255,.8)');
    gradient.addColorStop(1, 'rgba(255,255,255,0)');
    particleContext.fillStyle = gradient;
    particleContext.fillRect(0, 0, 64, 64);
    const particleTexture = new THREE.CanvasTexture(particleCanvas);
    textures.push(particleTexture);
    const billowCanvas = document.createElement('canvas');
    billowCanvas.width = billowCanvas.height = 128;
    const billowContext = billowCanvas.getContext('2d')!;
    const billowImage = billowContext.createImageData(128, 128);
    for (let y = 0; y < 128; y++) for (let x = 0; x < 128; x++) {
      const index = (y * 128 + x) * 4;
      billowImage.data.set([255, 255, 255, Math.round(255 * billowAlpha(x / 63.5 - 1, y / 63.5 - 1))], index);
    }
    billowContext.putImageData(billowImage, 0, 0);
    const billowTexture = new THREE.CanvasTexture(billowCanvas);
    textures.push(billowTexture);
    const buildingFire=new CityFire(billowTexture);scene.add(buildingFire.root);
    const contactGeometry = new THREE.PlaneGeometry(1, 1);
    const contactMaterial = new THREE.MeshBasicMaterial({
      color: 0x080b09, map: particleTexture, transparent: true,
      opacity: 0.48, depthWrite: false,
    });
    const headlightCanvas = document.createElement('canvas');
    headlightCanvas.width = headlightCanvas.height = 128;
    const headlightContext = headlightCanvas.getContext('2d')!;
    const headlightImage = headlightContext.createImageData(128, 128);
    for (let y = 0; y < 128; y++) for (let x = 0; x < 128; x++)
      headlightImage.data.set([255, 255, 255, Math.round(255 * headlightAlpha(x / 63.5 - 1, y / 127))], (y * 128 + x) * 4);
    headlightContext.putImageData(headlightImage, 0, 0);
    const headlightTexture = new THREE.CanvasTexture(headlightCanvas); textures.push(headlightTexture);
    const headlightGeometry = new THREE.PlaneGeometry(3, 7);
    const headlightMaterial = new THREE.MeshBasicMaterial({map: headlightTexture, color: 0xffdf9b,
      transparent: true, opacity: .28, depthWrite: false, blending: THREE.AdditiveBlending});
    const headlightTransform = new THREE.Object3D();
    let headlightPools = new THREE.InstancedMesh(headlightGeometry, headlightMaterial, 64);
    headlightPools.frustumCulled = false; headlightPools.count = 0; scene.add(headlightPools);
    const addVehicleShadow = (object: THREE.Group, model: string) => {
      const size = trafficSize(model);
      const shadow = new THREE.Mesh(contactGeometry, contactMaterial);
      shadow.name = 'vehicle-contact-shadow';
      shadow.rotation.x = -Math.PI / 2;
      shadow.position.y = 0.017;
      shadow.scale.set(size.width * 1.3, size.length, 1);
      object.add(shadow);
    };
    const pools = new THREE.InstancedMesh(
      new THREE.PlaneGeometry(1, 1),
      new THREE.MeshBasicMaterial({
        color: 0xffbf63,
        map: particleTexture,
        transparent: true,
        opacity: 0.24,
        depthWrite: false,
        blending: THREE.AdditiveBlending,
      }),
      plan.lots.length * 2,
    );
    n = 0;
    for (const lot of plan.lots)
      for (const at of lampPositions(lot)) {
        tmp.position.set(at.x, 0.185, at.z);
        tmp.rotation.set(-Math.PI / 2, 0, 0);
        tmp.scale.set(8, 8, 1);
        tmp.updateMatrix();
        pools.setMatrixAt(n++, tmp.matrix);
      }
    pools.visible = false;
    scene.add(pools);
    const raycaster = new THREE.Raycaster();
    const pointer = new THREE.Vector2();
    let hovered = '';
    let hoveredAt = 0;
    let down = {x: 0, y: 0};
    const pointerDown = (e: PointerEvent) => {
      cameraIntent++;
      setFollow(false);
      down = {x: e.clientX, y: e.clientY};
    };
    const pointerUp = (e: PointerEvent) => {
      if (e.button !== 0 || Math.hypot(e.clientX - down.x, e.clientY - down.y) > 5) return;
      const rect = canvas.getBoundingClientRect();
      pointer.set(
        ((e.clientX - rect.left) / rect.width) * 2 - 1,
        (-(e.clientY - rect.top) / rect.height) * 2 + 1,
      );
      raycaster.setFromCamera(pointer, camera);
      const hits = raycaster.intersectObjects([...buildings.values(), ...landings], true);
      if (hits.length) {
        let ob: THREE.Object3D | null = hits[0].object;
        while (ob && !ob.userData.place) ob = ob.parent;
        if (ob) { setFollow(false); latest.current.onSelect(ob.userData.place); }
      }
    };
    const pointerMove = (e: PointerEvent) => {
      if (e.buttons || performance.now() - hoveredAt < 70) return;
      hoveredAt = performance.now();
      const rect = canvas.getBoundingClientRect();
      pointer.set(
        ((e.clientX - rect.left) / rect.width) * 2 - 1,
        (-(e.clientY - rect.top) / rect.height) * 2 + 1,
      );
      raycaster.setFromCamera(pointer, camera);
      const hits = raycaster.intersectObjects([...buildings.values(), ...landings], true);
      let ob: THREE.Object3D | null = hits[0]?.object || null;
      while (ob && !ob.userData.place) ob = ob.parent;
      hovered = ob?.userData.place || '';
      canvas.style.cursor = hovered ? 'pointer' : 'grab';
    };
    canvas.addEventListener('pointermove', pointerMove);
    canvas.addEventListener('pointerdown', pointerDown);
    canvas.addEventListener('pointerup', pointerUp);
    const manualZoom=()=>{cameraIntent++;setFollow(false);};
    canvas.addEventListener('wheel',manualZoom,{passive:true});
    const keyboardPan = new KeyboardPan();
    const unbindPan = bindKeyboardPan(canvas, keyboardPan);
    let panTime = performance.now();
    const keys = (e: KeyboardEvent) => {
      const command = cameraCommand(e);
      if (!command) return;
      cameraIntent++;
      e.preventDefault();
      if (command === 'reset' || command.startsWith('pan-') || command.startsWith('zoom-')) setFollow(false);
      if (command === 'zoom-in' || command === 'zoom-out') {
        camera.zoom = THREE.MathUtils.clamp(
          camera.zoom * (command === 'zoom-out' ? 0.9 : 1.1),
          controls.minZoom,
          controls.maxZoom,
        );
        camera.updateProjectionMatrix();
      }
      if (command === 'reset') {
        reset();
      }
      keyboardPan.press(e);
    };
    canvas.addEventListener('keydown', keys);
    void preloadCityGunshots();
    void preloadCityEffects();
    const makeLabel = (name: string) => {
      const c = document.createElement('canvas');
      c.width = 512;
      c.height = 80;
      const ctx = c.getContext('2d')!;
      ctx.fillStyle = '#242924dd';
      ctx.fillRect(0, 0, 512, 80);
      ctx.strokeStyle = '#ae9970';
      ctx.strokeRect(1, 1, 510, 78);
      ctx.fillStyle = '#eadbb8';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.font = '500 28px Georgia';
      ctx.fillText(name, 256, 40, 490);
      const t = new THREE.CanvasTexture(c);
      t.colorSpace = THREE.SRGBColorSpace;
      textures.push(t);
      const s = new THREE.Sprite(
        new THREE.SpriteMaterial({map: t, depthTest: true, transparent: true, toneMapped: false}),
      );
      s.scale.set(17, 2.65, 1);
      return s;
    };
    const personModel = (id: string) => {
      const w = latest.current.state;
      return id === 'player'
        ? pedestrianModel(w.player.name, w.player.face, true)
        : pedestrianModel(id, w.everyone?.find(person => person.id === id)?.face);
    };
    const personWardrobe = (id: string) => {
      const w = latest.current.state;
      return id === 'player' ? wardrobe(w.player.name, w.player.face, true)
        : wardrobe(id, w.everyone?.find(person => person.id === id)?.face);
    };
    let movementClock = performance.now();
    const releaseActor = (actor: Actor) => {
      scene.remove(actor.object);
      actor.lamps.forEach(material => material.dispose());
      actor.wardrobe.forEach(material => material.dispose());
    };
    const addActor = (id: string, model: string): Actor => {
      const object = models.get(model)!.clone(true);
      if (!isPedestrian(model)) addVehicleShadow(object, model);
      scene.add(object);
      const limbs: THREE.Object3D[] = [];
      const wheels: THREE.Object3D[] = [];
      const lamps: THREE.MeshStandardMaterial[] = [];
      object.traverse(o => {
        if (o instanceof THREE.Mesh) {
          const cloneLamp = (m: THREE.Material) => {
            if (!(m instanceof THREE.MeshStandardMaterial) || !['headlamps', 'tail lamps'].includes(m.name)) return m;
            const own = m.clone(); own.emissiveIntensity = 0; lamps.push(own); return own;
          };
          o.material = Array.isArray(o.material) ? o.material.map(cloneLamp) : cloneLamp(o.material);
        }
        if (o.name.startsWith('wheel-roll-')) { o.rotation.order = 'YXZ'; wheels.push(o); }
        if (o.name.startsWith('leg') || o.name.startsWith('arm') || o.name.startsWith('knee'))
          limbs.push(o);
      });
      const costume=isPedestrian(model)?dressPedestrian(object,model,personWardrobe(id)):[];
      let driver:THREE.Group|undefined;
      if(['ford','hudson','packard'].includes(model)){
        const seat=object.getObjectByName('seat-front-left');
        if(seat){
          const driverModel=personModel(id);driver=models.get(driverModel)!.clone(true);
          costume.push(...dressPedestrian(driver,driverModel,personWardrobe(id)));
          seatDriver(driver,seat.position,['left','right'].map(side=>object.getObjectByName('seat-driver-grip-'+side)!.position));driver.name='vehicle-driver';driver.visible=false;
          const windows=new Map<THREE.Material,THREE.MeshStandardMaterial>();
          object.traverse(part=>{if(part instanceof THREE.Mesh){
            const glaze=(m:THREE.Material)=>{if(!(m instanceof THREE.MeshStandardMaterial)||m.name!=='car glass')return m;
              let own=windows.get(m);if(!own){own=m.clone();own.transparent=true;own.opacity=.38;own.depthWrite=false;windows.set(m,own);costume.push(own);}return own;};
            part.material=Array.isArray(part.material)?part.material.map(glaze):glaze(part.material);
          }});
          object.add(driver);
        }
      }
      const actor = {
        object, driver,
        model,
        points: [{x: 0, z: 0}],
        start: 0,
        end: 0,
        since: 0,
        duration: 0,
        walking: false,
        phase: 0,
        realSince: 0,
        wardrobe: costume,
        wardrobeKey: JSON.stringify(personWardrobe(id)),
        limbs, wheels, lamps, wheelPhase: 0, steering: 0, wheelPlaced: false,
      };
      actors.set(id, actor);
      return actor;
    };
    const assign = (
      id: string,
      model: string,
      points: Point[],
      start: number,
      end: number,
      duration: number,
      now: number,
    ) => {
      let a = actors.get(id);
      if (a && (a.model !== model || ((isPedestrian(model)||a.driver) && a.wardrobeKey !== JSON.stringify(personWardrobe(id))))) {
        releaseActor(a);
        actors.delete(id);
        a = undefined;
      }
      if (!a) a = addActor(id, model);
      // Arrivals hide their outdoor actor; a later journey must show it again.
      a.object.visible = true;
      a.arrived = false;
      a.timelineProgress=undefined;a.legKey=undefined;
      a.wheelPlaced = false;
      a.points = points;
      a.start = start;
      a.end = end;
      a.since = movementClock;
      a.realSince = now;
      a.duration = duration;
      a.walking = isPedestrian(model) && start !== end;
      if (start === end) {
        a.steering = 0;
        for (const wheel of a.wheels) wheel.rotation.y = 0;
      }
    };
    const loader = new GLTFLoader();
    let graphicsLost = false;
    let ready = false,
      revision = -1,
      worldID = '',
      previous: Snapshot | null = null,
      journeyKey = '',
      reportedJourneyProgress = -1,
      reportedJourneyBlocked = false,
      playedJourneyProgress = 0,
      renderedJourneyProgress = 0,
      streetPlayback:StreetPlayback|null = null,
      streetPlaybackOwner = '',
      settlingStreet = false,
      wasFollowing = false,
      followZoom: number | null = null,
      motionWas = true;
    let lastActive: string | null = null;
    let completedJourney = '';
    Promise.all(
      modelNames.map(async name => {
        const gltf = await loader.loadAsync(`/art/models/${name}.glb`);
        if (dead) {
          disposeCityResources([gltf.scene]);
          return;
        }
        models.set(name, gltf.scene);
      }),
    )
      .then(() => {
        if (dead) return;
        const vacantLot = models.get('vacant-lot')!;
        vacantLot.updateMatrixWorld(true);
        if (plan.vacant.length) vacantLot.traverse(part => {
          if (!(part instanceof THREE.Mesh)) return;
          const instances = new THREE.InstancedMesh(part.geometry, part.material, plan.vacant.length);
          plan.vacant.forEach((lot, index) => {
            const transform = new THREE.Matrix4().makeTranslation(lot.x, .18, lot.z);
            instances.setMatrixAt(index, transform.multiply(part.matrixWorld));
          });
          instances.castShadow = instances.receiveShadow = true;
          scene.add(instances);
        });
        const streetBed = models.get('street-bed')!;
        const parcels = [...plan.lots, ...plan.vacant];
        streetBed.updateMatrixWorld(true);
        streetBed.traverse(part => {
          if (!(part instanceof THREE.Mesh)) return;
          const instances = new THREE.InstancedMesh(part.geometry, part.material, parcels.length);
          parcels.forEach((lot, index) => {
            const transform = new THREE.Matrix4().makeTranslation(lot.x, 0, lot.z);
            instances.setMatrixAt(index, transform.multiply(part.matrixWorld));
          });
          instances.receiveShadow = true;
          scene.add(instances);
        });
        const furniture = models.get('streetside')!;
        furniture.updateMatrixWorld(true);
        furniture.traverse(part => {
          if (!(part instanceof THREE.Mesh)) return;
          const cutawayMaterial=(source:THREE.Material)=>{
            if(!(source instanceof THREE.MeshStandardMaterial))return source;
            const material=source.clone();buildingCondition(material,100,new THREE.Vector3());
            streetCutawayMaterials.add(material);return material;
          };
          const material=Array.isArray(part.material)?part.material.map(cutawayMaterial):cutawayMaterial(part.material);
          const instances = new THREE.InstancedMesh(part.geometry, material, plan.lots.length);
          plan.lots.forEach((lot, index) => {
            const at = streetsidePosition(lot);
            const transform = new THREE.Matrix4().makeTranslation(at.x, 0.18, at.z);
            instances.setMatrixAt(index, transform.multiply(part.matrixWorld));
          });
          instances.castShadow = instances.receiveShadow = true;
          scene.add(instances);
        });
        if (harbourLot) {
          const pier = models.get('harbour-pier')!.clone(true);
          pier.position.set(-24, 0, harbourLot.z); pier.userData.place = harbourLot.id;
          pier.traverse(part => {if (part instanceof THREE.Mesh) part.castShadow = part.receiveShadow = true;});
          scene.add(pier); landings.push(pier);
          const quay = models.get('quay-section')!; quay.updateMatrixWorld(true);
          const sections: number[] = [];
          for (let z = -16; z <= plan.depth + 16; z += 8) if (Math.abs(z - harbourLot.z) >= 4) sections.push(z);
          quay.traverse(part => {
            if (!(part instanceof THREE.Mesh)) return;
            const instances = new THREE.InstancedMesh(part.geometry, part.material, sections.length);
            sections.forEach((z, i) => instances.setMatrixAt(i, new THREE.Matrix4().makeTranslation(-17.15, 0, z).multiply(part.matrixWorld)));
            instances.castShadow = instances.receiveShadow = true; scene.add(instances);
          });
        }
        for (const lot of plan.lots) {
          const model = models.get(lot.model)!.clone(true);
          model.position.set(lot.x, 0.18, lot.z);
          model.userData.place = lot.id;
          if (['filling', 'garage', 'dealer', 'chapel', 'docks', 'haulage'].includes(lot.model))
            model.rotation.y = Math.PI;
          model.traverse(o => {
            if (o instanceof THREE.Mesh) {
              o.castShadow = true;
              o.receiveShadow = true;
            }
          });
          scene.add(model);
          buildings.set(lot.id, model);
          model.traverse(o=>{if(o.name==='clock-hand-minute'||o.name==='clock-hand-hour')clockHands.push(o);});
          const label = makeLabel(latest.current.state.locations.find(p => p.id === lot.id)!.name);
          const box = new THREE.Box3().setFromObject(model, true);
          model.userData.sightBounds = box.clone();
          model.userData.debrisOrigin = {x: lot.x, y: 0.25, z: box.min.z - 0.15};
          const blastWindows=clearBlastWindows(model);
          model.userData.blastWindows=blastWindows.slice(0,4);
          model.userData.blastOrigin = blastWindows.length
            ? {x:lot.x,y:blastWindows[0].y,z:blastWindows[0].z+.45}
            : {x:lot.x,y:Math.min(2,box.max.y*.5),z:lot.z};
          label.position.set(lot.x, box.max.y + 2, lot.z);
          scene.add(label);
          const anchor = model.getObjectByName('sign-anchor');
          if (anchor) {
            const sign = new THREE.Mesh(
              new THREE.PlaneGeometry(7.5, 1.05),
              new THREE.MeshBasicMaterial({
                map: (label.material as THREE.SpriteMaterial).map,
                toneMapped: false,
              }),
            );
            if (!['filling', 'garage', 'dealer', 'chapel', 'docks', 'haulage'].includes(lot.model))
              sign.rotation.y = Math.PI;
            anchor.add(sign);
          }
          labels.set(lot.id, label);
        }
        renderer.shadowMap.needsUpdate = true;
        ready = true;
        setStatus('');
      })
      .catch(err => {
        if (!dead) setFailure(`City assets could not load: ${String(err.message || err)}`);
      });
    const ambientAudio=new CityAmbientAudio();
    const vehicleSounds=new Map<string,{key:string;stop?:()=>void}>();
    const silenceAmbient=()=>{ambientAudio.dispose();for(const voice of vehicleSounds.values())voice.stop?.();};
    const hiddenAudio=()=>{if(document.hidden)silenceAmbient();};
    document.addEventListener('visibilitychange',hiddenAudio);
    const reduced = matchMedia('(prefers-reduced-motion: reduce)');
    let sampleStart = performance.now(),
      samples: number[] = [],
      last = sampleStart;
    const tick = (now: number) => {
      if (dead) return;
      frame = requestAnimationFrame(tick);
      const inputSeconds = (now - panTime) / 1000;
      const turn = keyboardPan.rotation(inputSeconds);
      if (turn) {
        cameraIntent++;
        const offset = camera.position.clone().sub(controls.target);
        offset.applyAxisAngle(new THREE.Vector3(0, 1, 0), turn);
        camera.position.copy(controls.target).add(offset);
      }
      const pan = keyboardPan.step(camera.position, controls.target, inputSeconds, 110 / camera.zoom);
      panTime = now;
      if (pan.x || pan.z) {
        cameraIntent++;
        setFollow(false);
        camera.position.x += pan.x; camera.position.z += pan.z;
        controls.target.x += pan.x; controls.target.z += pan.z;
      }
      if (document.hidden) {
        last = now;
        return;
      }
      const dt = now - last;
      movementClock += Math.min(100, dt) * playback.current;
      last = now;
      const p = latest.current,
        w = p.state,
        motion = p.motion && !reduced.matches && !graphicsLost;
      let presentationMinute=w.minute;
      const aftermathRecords=p.journey?responseRecords(p.journey.beforeResponse?.aftermath,w.aftermath):w.aftermath||[];
      const policeRecords=p.journey?responseRecords(p.journey.beforeResponse?.police_presence,w.police_presence):w.police_presence||[];
      const fireRecords=p.journey?responseRecords(p.journey.beforeResponse?.building_fires,w.building_fires):w.building_fires||[];
      const nearbyIdle:string[]=[];
      const impact={x:0,y:0};
      const addImpact=(age:number,strength:number)=>{const pulse=impactPulse(age,strength);impact.x+=pulse.x;impact.y+=pulse.y;};
      const activePlayback=p.activeCue?`${p.activeCue.id}:${p.replaySerial??0}`:null;
      const playbackStarted = !!activePlayback && lastActive !== activePlayback;
      const handOffSurvivor=(effect:Effect)=>{
        if(!effect.accident||effect.accident.fatal||!effect.slot||effect.cue.id.startsWith('preview:')||
          effect.cue.attacker?.id!=='player'||!w.player.alive||w.player.location!==effect.cue.target)return;
        const actor=actors.get('player');if(!actor||actor.start!==actor.end||!isPedestrian(actor.model))return;
        effect.accident.update(effect.accident.duration);
        const at=effect.accident.actor.getWorldPosition(new THREE.Vector3());
        const pose={x:at.x,z:at.z,heading:effect.slot.pose.heading};
        if(traffic.adoptFrontage({id:'player',model:actor.model,points:actor.points,progress:actor.end},pose)){
          actor.object.position.set(at.x,at.y,at.z);actor.object.rotation.set(0,pose.heading,0);
          actor.limbs.forEach(limb=>limb.rotation.x=0);actor.object.visible=true;
        }
      };
      if ((lastActive && (!p.activeCue||playbackStarted)) || worldID !== `${w.id}:${w.life}`) {
        for (const effect of effects) {
          if(!playbackStarted&&worldID===`${w.id}:${w.life}`)handOffSurvivor(effect);
          scene.remove(effect.mesh, effect.light);
          if (effect.extra) scene.remove(effect.extra);
          effect.wardrobe?.forEach(material => material.dispose());
          effect.audio?.dispose(); effect.glassAudio?.dispose(); effect.reactionAudio?.dispose(); effect.departureAudio?.dispose();
          effect.driveBy?.dispose(); disposeDebris(effect);
          effect.mesh.dispose();
          (effect.mesh.material as THREE.Material).dispose();
        }
        effects.length = 0;
      }
      const streetOwner=`${w.id}:${w.life}:${w.revision}`;
      if(settlingStreet&&(streetPlaybackOwner!==streetOwner||!motion||playbackStarted)){
        // A new authoritative state or explicit presentation change replaces the
        // remainder. Disabling motion reconciles it immediately as well.
        if(streetPlaybackOwner===streetOwner)revision=-1;
        settlingStreet=false;streetPlayback=null;streetPlaybackOwner='';
      }
      if (
        ready &&
        (revision !== w.revision || worldID !== `${w.id}:${w.life}` || playbackStarted)
      ) {
        const first = worldID !== `${w.id}:${w.life}` || revision < 0;
        if (first) {
          for (const a of actors.values()) releaseActor(a);
          actors.clear();
          traffic.clear();
          previous = null;
        }
        const seen = new Set<string>();
        for (const j of p.journey?.street ? [] : w.street || []) {
          const from = lots.get(j.from_id),
            to = lots.get(j.to_id);
          if (!from || !to) continue;
          seen.add(j.id);
          const model = j.vehicle ? carModel(j.vehicle) : personModel(j.id);
          const before = previous?.street?.find(
            b => b.id === j.id && b.from_id === j.from_id && b.to_id === j.to_id,
          );
          assign(
            j.id,
            model,
            route(from, to, !!j.vehicle),
            motion && before ? before.progress : j.progress,
            j.progress,
            motion ? 1600 : 0,
            now,
          );
        }
        for (const person of p.journey?.street ? [] : w.everyone || []) {
          if (seen.has(person.id) || !person.where_id) continue;
          const lot = lots.get(person.where_id);
          if (!lot) continue;
          // Only observed arrivals get outdoor actors; an indoor person stays indoors.
          const before = previous?.street?.find(j => j.id === person.id && j.to_id === lot.id);
          if (before && motion) {
            const from = lots.get(before.from_id);
            if (from) {
              seen.add(person.id);
              assign(
                person.id,
                before.vehicle ? carModel(before.vehicle) : personModel(person.id),
                route(from, lot, !!before.vehicle),
                before.progress,
                1,
                1600,
                now,
              );
            }
          }
        }
        for (const [id, a] of actors)
          if (!id.startsWith('player') && !seen.has(id)) {
            releaseActor(a);
            actors.delete(id);
          }
        const here = lots.get(w.player.location);
        if (here && w.player.alive) {
          if (!actors.has('player')) assign('player', personModel('player'), [entrance(here)], 0, 0, 0, now);
        } else {
          const a = actors.get('player');
          if (a) releaseActor(a);
          actors.delete('player');
        }
        for (const cue of assassinationBatch(p.journey
          ? []
          : cueQueue.take(
              `${w.id}:${w.life}`,
              w.last_result?.cues || [],
              p.activeCue,
              first,
              !first && revision === w.revision && playbackStarted,
            ).flatMap(policeCast))) {
          if (!motion) continue;
          const lot = lots.get(cue.target);
          if (!lot) continue;
          const mesh = new THREE.InstancedMesh(
            effectGeometry,
            new THREE.MeshBasicMaterial({
              color: 0xffffff,
              map: cue.kind === 'explosion' ? billowTexture : particleTexture,
              transparent: true,
              opacity: 0.85,
              depthWrite: false,
            }),
            32,
          );
          mesh.frustumCulled = false;
          scene.add(mesh);
          const light = new THREE.PointLight(0xffa64d, 0, 25);
          const at = entrance(lot);
          light.position.set(at.x, 3, at.z);
          scene.add(light);
          let costume: THREE.MeshStandardMaterial[] | undefined;
          let weapon:THREE.Group|undefined;
          let assassination:CityAssassination|undefined;
          let driveBy:CityBuildingDriveBy|undefined;
          let custody:CityCustody|undefined;
          let incendiary:CityIncendiary|undefined;
          let planter:CityPlanter|CityVillaExit|undefined;
          let accident:CityAccident|undefined;
          const hasAccident=cue.kind==='explosion'&&cue.detonation==='premature'&&!!cue.accident&&!!cue.attacker;
          const planterBuilding=buildings.get(cue.target);
          const hasPlanter=cue.kind==='explosion'&&cue.detonation==='planted'&&!!cue.attacker&&
            !!(planterBuilding?.getObjectByName('entrance-threshold')||planterBuilding?.getObjectByName('entrance-landing'))&&!!planterBuilding.getObjectByName('entrance-door-hinge');
          const weaponModel=sceneWeapon(cue.attacker?.weapon);
          let extra: THREE.Group | undefined, gunArm: THREE.Object3D | undefined, muzzle: THREE.Object3D | undefined;
          if (hasAccident || hasPlanter || ['incendiary', 'attack', 'killing', 'gunfight', 'raid', 'arrest','raid-unit','police-unit','officer','detainee','raid-officer'].includes(cue.kind)) {
            const model = ['killing','detainee'].includes(cue.kind) ? personModel(cue.actors?.[0]?.id || '')
              : hasAccident || hasPlanter || ['gunfight','attack','incendiary'].includes(cue.kind) ? personModel(cue.attacker?.id||'anonymous-shooter') : ['officer','raid-officer'].includes(cue.kind)?'police-officer':'police';
            extra = models.get(model)!.clone(true);
            if (isPedestrian(model)) costume = dressPedestrian(extra, model, personWardrobe(['killing','detainee'].includes(cue.kind) ? cue.actors?.[0]?.id || '' : cue.attacker?.id||'anonymous-shooter'));
            if (model === 'police') addVehicleShadow(extra, model);
            if(hasAccident){accident=new CityAccident(extra,cue.accident!.fatal);extra=accident.root;}
            if(hasPlanter){planter=planterBuilding?.getObjectByName('entrance-landing')?new CityVillaExit(extra):new CityPlanter(extra);extra=planter.root;}
            if(cue.kind==='incendiary'){
              incendiary=new CityIncendiary(extra,models.get('incendiary-bottle')!.clone(true),new THREE.Vector3(0,3,4));extra=incendiary.root;
            }
            if(cue.kind==='detainee'){
              const cuffs=models.get('handcuffs')!.clone(true);cuffs.name='custody-restraint';cuffs.visible=false;extra.add(cuffs);
              custody=new CityCustody(extra,models.get('police-officer')!.clone(true));extra=custody.root;
            }
            if (cue.kind === 'gunfight') {
              extra.rotation.y = Math.PI / 2;
              gunArm = extra.getObjectByName('arm1');
              if(weaponModel){
                weapon = models.get(weaponModel)!.clone(true);
                if(weaponModel==='revolver'){weapon.position.set(0,-.58,0);weapon.rotation.x=Math.PI/2;gunArm?.add(weapon);}
                else extra.add(weapon);
                muzzle=weapon.getObjectByName('muzzle');
              }
            }
            if(isStagedStrike(cue)){
              const victimID=cue.strike!.victim.id, victimModel=personModel(victimID);
              const victim=models.get(victimModel)!.clone(true);
              costume=[...(costume||[]),...dressPedestrian(victim,victimModel,personWardrobe(victimID))];
              assassination=new CityAssassination(extra,victim,weapon,cue.strike!.variant,weaponModel||'revolver');extra=assassination.root;
            }
            extra.visible = false;
            mesh.visible = false;
            scene.add(extra);
          }
          if(cue.kind==='driveby-building'&&cue.drive_by&&cue.attacker&&weaponModel){
            const vehicleModel=carModel(cue.drive_by.vehicle),car=models.get(vehicleModel)!.clone(true);
            const shooterModel=personModel(cue.attacker.id),driverModel=personModel(cue.drive_by.driver.id);
            const shooter=models.get(shooterModel)!.clone(true),driver=models.get(driverModel)!.clone(true);
            costume=[...dressPedestrian(shooter,shooterModel,personWardrobe(cue.attacker.id)),...dressPedestrian(driver,driverModel,personWardrobe(cue.drive_by.driver.id))];
            weapon=models.get(weaponModel)!.clone(true);gunArm=shooter.getObjectByName('arm1');muzzle=weapon.getObjectByName('muzzle');
            addVehicleShadow(car,vehicleModel);
            driveBy=new CityBuildingDriveBy(car,driver,shooter,weapon,weaponModel);extra=driveBy.root;
            extra.visible=false;mesh.visible=false;scene.add(extra);
          }
          let debris: THREE.InstancedMesh | undefined;
          if (cue.kind === 'explosion' || incendiary) {
            const model = models.get(incendiary?'bottle-shard':'blast-fragment')!;
            model.updateMatrixWorld(true);
            model.traverse(part => {
              if (!(part instanceof THREE.Mesh) || debris) return;
              const material = (part.material as THREE.MeshStandardMaterial).clone();
              material.transparent = true;
              debris = new THREE.InstancedMesh(part.geometry.clone().applyMatrix4(part.matrixWorld), material, 12);
              debris.frustumCulled = false;
              if(incendiary||planter||accident)debris.visible=false;
              scene.add(debris);
            });
          }
          effects.push({cue, accident, driveBy, assassination, custody, incendiary, planter, glassBlocked:incendiary?new Set():undefined, since: now, mesh, light, debris, extra, wardrobe: costume, gunArm, muzzle, weapon, weaponModel:weaponModel||undefined,
            glazingBefore:(cue.id.startsWith('preview:')?undefined:p.beforeConditions?.[cue.target]) ?? buildings.get(cue.target)?.userData.condition ?? w.locations.find(p=>p.id===cue.target)?.condition ?? 100,
            departureAudio:driveBy?new BlastAudio(()=>playRecordedEffect('drive-away')):undefined,
            reactionAudio:driveBy?new BlastAudio(()=>playRecordedEffect('vehicle-approach')):assassination?new BlastAudio(()=>playRecordedEffect('pain')):cue.kind==='explosion'?new BlastAudio(()=>playRecordedEffect('panic'))
              :cue.kind==='killing'&&(w.last_result?.cues||[]).some(gun=>gunVictim(gun,cue))?new BlastAudio(()=>playRecordedEffect('pain')):undefined,
            glassAudio:['explosion','incendiary'].includes(cue.kind)?new BlastAudio(()=>playMoment('glass-break')):undefined,
            audio: driveBy?new GunfireAudio(()=>playCityGunshot(weaponModel!),driveBy.shots,true):cue.kind==='attack'&&assassination?new GunfireAudio(()=>playMoment('body-hit'),MELEE_IMPACTS,true):cue.kind === 'gunfight' && weaponModel ? new GunfireAudio(()=>playCityGunshot(weaponModel),weaponShots(weaponModel||undefined,cue.strike?.variant),true)
              : cue.kind === 'raid-officer' ? new BlastAudio(() => playMoment('door-breach'))
              : ['explosion','raid','arrest'].includes(cue.kind) ? new BlastAudio(() => playMoment(cue.kind)) : undefined});
          if (p.activeCue?.id === cue.id || (assassination && p.activeCue?.strike?.victim.id===cue.strike?.victim.id)) {
            setFollow(false);
            const envelope = new THREE.Box3();
            const siblings = (w.last_result?.cues || []).filter(other => other.target === cue.target);
            for (const other of [...siblings, cue].flatMap(policeCast)) {
              for (const slot of sceneSlots(lot, isStagedStrike(other)?'assassination':other.kind)) {
                envelope.expandByPoint(new THREE.Vector3(slot.root.x - 3.5, 0, slot.root.z - 3.5));
                envelope.expandByPoint(new THREE.Vector3(slot.root.x + 3.5, 3, slot.root.z + 3.5));
              }
              if (other.kind === 'explosion') {
                const building = buildings.get(lot.id);
                if (building) envelope.union(new THREE.Box3().setFromObject(building));
                const origin = building?.userData.blastOrigin || {x: lot.x, y: 0, z: at.z};
                envelope.expandByPoint(new THREE.Vector3(origin.x - 8, 0, origin.z - 8));
                envelope.expandByPoint(new THREE.Vector3(origin.x + 8, 15, origin.z + 8));
              }
            }
            if (envelope.isEmpty()) envelope.set(
              new THREE.Vector3(lot.x - 10, 0, at.z - 4),
              new THREE.Vector3(lot.x + 10, 5, at.z + 8));
            frameScene(camera, controls.target, envelope);
            controls.update();
          }
        }
        // Condition-driven surface stains imply no ongoing fire or invented collapse.
        for (const place of w.locations) {
          const b = buildings.get(place.id);
          if (b) {b.userData.condition=place.condition;buildingGlazing(b,place.condition);}
          if (b)
            b.traverse(o => {
              if (o instanceof THREE.Mesh) {
                if (!o.userData.ownMaterial) {
                  o.material = Array.isArray(o.material)
                    ? o.material.map(m => m.clone())
                    : o.material.clone();
                  o.userData.ownMaterial = true;
                }
                for (const m of Array.isArray(o.material) ? o.material : [o.material])
                  if (m instanceof THREE.MeshStandardMaterial)
                    buildingCondition(m, place.condition, b.position);
              }
            });
        }
        lightingKey="";
        revision = w.revision;
        worldID = `${w.id}:${w.life}`;
        previous = w;
      }
      if (ready) {
        const hereForCar = lots.get(w.player.location);
        const parked = actors.get('player-car');
        // Death does not delete the saved car; retain its empty parked model
        // until the authoritative vehicle/location changes (including new life).
        if (
          !p.journey &&
          hereForCar &&
          /Ford|Hudson|Packard/i.test(w.vehicle?.car || '')
        ) {
          if (!parked || parked.model !== carModel(w.vehicle?.car))
            assign('player-car', carModel(w.vehicle?.car), [parkingSpot(hereForCar)], 0, 0, 0, now);
          const car = actors.get('player-car')!;
          car.points = [parkingSpot(hereForCar)];
        } else if (parked) {
          releaseActor(parked);
          actors.delete('player-car');
        }
        const key = p.journey
          ? `${w.id}:${w.life}:${p.journey.from.id}:${p.journey.to.id}:${w.revision}`
          : '';
        if (journeyKey && !key) {
          settlingStreet=!!streetPlayback&&streetPlaybackOwner===streetOwner&&completedJourney===journeyKey&&motion&&!playbackStarted&&!p.activeCue;
          if(!settlingStreet)revision=-1;
        }
        if (key !== journeyKey || motion !== motionWas) {
          journeyKey = key;
          reportedJourneyProgress = -1;
          reportedJourneyBlocked=false;p.onJourneyBlocked?.(false);
          playedJourneyProgress = 0;renderedJourneyProgress=0;
          if(p.journey?.street){
            streetPlayback=new StreetPlayback(p.journey.street);streetPlaybackOwner=streetOwner;settlingStreet=false;
          }else if(!settlingStreet){streetPlayback=null;streetPlaybackOwner='';}
          if (p.journey) { setFollow(true); followZoom = 8; }
          const here = lots.get(w.player.location);
          if (here && w.player.alive) {
            const from = p.journey ? lots.get(p.journey.from.id) : undefined;
            const driving = p.journey?.driving ?? !!w.vehicle?.running;
            if (from && motion)
              assign(
                'player',
                driving ? carModel(p.journey?.vehicle || w.vehicle?.car) : personModel('player'),
                route(from, here, driving),
                0,
                1,
                p.journey?.street ? Math.max(2400,
                  pathLength(route(from,here,driving))/trafficSpeed(driving?carModel(p.journey?.vehicle):personModel('player'))*1000,
                  ...p.journey.street.map(segment=>{
                    const a=lots.get(segment.from_id),b=lots.get(segment.to_id);
                    if(!a||!b)return 0;
                    return pathLength(route(a,b,!!segment.vehicle))*(segment.end_progress-segment.progress)
                      /trafficSpeed(segment.vehicle?carModel(segment.vehicle):personModel(segment.id))
                      *p.journey!.minutes/Math.max(1,segment.to_minute-segment.from_minute)*1000;
                  })) : 2400,
                now,
              );
            else assign('player', personModel('player'), [entrance(here)], 0, 0, 0, now);
          }
        }
        // Non-travel commands can also change a life/location without a journey.
        if (!p.journey && actors.has('player') && (actors.get('player')!.model !== personModel('player')
          || actors.get('player')!.wardrobeKey !== JSON.stringify(personWardrobe('player')))) {
          const here = lots.get(w.player.location);
          if (here) assign('player', personModel('player'), [entrance(here)], 0, 0, 0, now);
        }
        const player = actors.get('player');
        if (player && !p.journey) {
          const here = lots.get(w.player.location);
          if (here) {
            player.points = [entrance(here)];
            player.start = player.end = 0;
          }
        }
        if (p.journey?.street || settlingStreet) {
          const traveller=actors.get('player');
          if(traveller&&p.journey)playedJourneyProgress=motion?advanceJourneyClock(
            playedJourneyProgress,renderedJourneyProgress,Math.min(100,Math.max(0,dt))/1000*playback.current,
            traveller.duration/1000,pathLength(traveller.points)/trafficSpeed(traveller.model)):1;
          const minute=p.journey?(p.journey.fromMinute??w.minute-p.journey.minutes)+p.journey.minutes*playedJourneyProgress:w.minute;
          const samples=streetPlayback!.sample(minute,(id,key)=>actors.get(id)?.legKey===key&&!!actors.get(id)?.arrived);
          if(settlingStreet&&!samples.size){settlingStreet=false;streetPlayback=null;streetPlaybackOwner='';}
          for(const [id,a] of actors)if(!id.startsWith('player')&&!samples.has(id)){releaseActor(a);actors.delete(id);}
          for(const [id,{segment,progress}] of samples){
            const from=lots.get(segment.from_id),to=lots.get(segment.to_id);if(!from||!to)continue;
            const legKey=streetLegKey(segment);
            const startingLeg=actors.get(id)?.legKey!==legKey;
            if(startingLeg){
              assign(id,segment.vehicle?carModel(segment.vehicle):personModel(id),route(from,to,!!segment.vehicle),segment.progress,segment.end_progress,1,now);
              actors.get(id)!.legKey=legKey;
            }
            // Admit a delayed successor at its observed start, never halfway
            // down the next street merely because its scheduled time has passed.
            actors.get(id)!.timelineProgress=startingLeg?segment.progress:progress;
          }
        }
        presentationMinute=p.journey?(p.journey.fromMinute??w.minute-p.journey.minutes)+p.journey.minutes*playedJourneyProgress:w.minute;
        if (!motion) traffic.clear();
        // Stationary actors own their known destination even before their first
        // visible frame. Otherwise response vehicles can steal a parked bay.
        const replacedActors=new Set(effects.flatMap(e=>e.driveBy?[e.cue.attacker?.id,...(isPedestrian(actors.get(e.cue.drive_by?.driver.id||'')?.model||'')?[e.cue.drive_by?.driver.id]:[]),...(e.cue.attacker?.id==='player'?['player-car']:[])]:e.incendiary||e.planter||e.accident?[e.cue.attacker?.id]:[]));
        const actorSpaces=[...actors].filter(([id])=>!replacedActors.has(id)).filter(([,a])=>a.object.visible||a.start===a.end).map(([id,a])=>{
          const pose=a.start===a.end?(traffic.placement(id)?.pose||onRoute(a.points,1)):{x:a.object.position.x,z:a.object.position.z,heading:a.object.rotation.y};
          return {model:trafficModel(a.model,a.start===a.end),root:{x:pose.x,z:pose.z},pose};
        });
        const animatingVictims=new Set(effects.flatMap(e=>e.accident?.fatal&&!e.cue.id.startsWith('preview:')?[`player:${w.life}`]:e.assassination?[e.cue.strike!.victim.id]:e.cue.kind==='killing'?e.cue.actors?.map(a=>a.id)||[]:[]));
        aftermath.suppressVictims(animatingVictims);
        // Reserve the planter/shooter before associated casualties, so the
        // casualties cannot occupy the path their triggering scene needs.
        const stagingOrder = [...effects].sort((a, b) =>
          (b.planter?2:Number(b.cue.kind === 'gunfight')) - (a.planter?2:Number(a.cue.kind === 'gunfight')));
        let policeStaged=false;
        for (const e of stagingOrder) {
          if (!e.extra || e.slot) continue;
          const occupied = [
            ...actorSpaces,
            ...effects.flatMap(other => other.slot ? [other.slot] : []),
            ...aftermath.slots(),
          ];
          const entryBuilding=buildings.get(e.cue.target);
          const entry=(entryBuilding?.getObjectByName('entrance-threshold')||(e.planter?entryBuilding?.getObjectByName('entrance-landing'):undefined))?.getWorldPosition(new THREE.Vector3());
          const fireBuilding=e.incendiary?buildings.get(e.cue.target):undefined;
          const fireWindows=fireBuilding?clearBlastWindows(fireBuilding):[];
          let stagedFlight:ReturnType<typeof incendiaryStagingFlight>=null;
          e.slot = availableSceneSlot(lots.get(e.cue.target)!, e.accident?'accident':e.planter?'planter':e.assassination?'assassination':e.custody?'custody':e.cue.kind, occupied,entry,e.incendiary?slot=>{
            stagedFlight=fireBuilding?incendiaryStagingFlight(slot.root,e.incendiary!.release,fireWindows,fireBuilding):null;
            return stagedFlight!==null;
          }:undefined);
          if (e.slot) {
            e.extra.position.set(e.slot.root.x, (e.slot.model === 'parked-police'||e.driveBy) ? vehicleRootHeight(e.slot.root) : 0.2, e.slot.root.z);
            e.light.position.set(e.slot.root.x, 3, e.slot.root.z);
            e.since = now;
            if(e.driveBy){
              const facade=entryBuilding?.userData.sightBounds as THREE.Box3|undefined;
              e.driveBy.target.set(0,1.8,(facade?.min.z??e.slot.root.z+8)-e.slot.root.z);
              frameScene(camera,controls.target,new THREE.Box3(new THREE.Vector3(e.slot.root.x-13,0,e.slot.root.z-2),new THREE.Vector3(e.slot.root.x+14,3,e.slot.root.z+e.driveBy.target.z+1)));controls.update();
            }
            if(e.accident){
              e.extra.rotation.y=e.slot.pose.heading;
              if(e.accident.fatal&&!e.cue.id.startsWith('preview:'))aftermath.rememberBody(`player:${w.life}`,e.slot,0);
              frameScene(camera,controls.target,new THREE.Box3(new THREE.Vector3(e.slot.root.x-3,0,e.slot.root.z-2),new THREE.Vector3(e.slot.root.x+3,3,e.slot.root.z+2)));
              controls.update();
            }
            if(e.planter&&entry){
              e.extra.position.y=entry.y;
              e.since=now+e.planter.duration*1000;
              const bounds=new THREE.Box3().setFromPoints([
                new THREE.Vector3(entry.x-3,0,entry.z-3.2),new THREE.Vector3(entry.x+1,3,entry.z+2.8)]);
              frameScene(camera,controls.target,bounds);controls.update();
              const wide=new THREE.Box3(new THREE.Vector3(entry.x-8,0,entry.z-8),new THREE.Vector3(entry.x+8,15,entry.z+5));
              const facade=buildings.get(e.cue.target)?.userData.sightBounds as THREE.Box3|undefined;
              if(facade)wide.union(facade);
              e.pullback={move:new ScenePullback(camera,controls.target,wide),intent:cameraIntent};
            }
            if(e.incendiary){
              // A slot is accepted only with a checked flight. Never invent a
              // fallback through a facade when a window has no clear approach.
              const flight=stagedFlight!;
              const target=flight.target;
              e.incendiary.loft=flight.loft;
              e.incendiary.target.copy(target).sub(e.extra.position);
              frameScene(camera,controls.target,new THREE.Box3().setFromPoints([
                new THREE.Vector3(e.slot.root.x-5,0,e.slot.root.z-1.2),
                new THREE.Vector3(e.slot.root.x+1.5,2.5,e.slot.root.z+1.2),target.clone().add(new THREE.Vector3(1,1,0))]));
              controls.update();
            }
            if(p.activeCue && ['arrest','raid'].includes(p.activeCue.kind) && (e.cue.id===p.activeCue.id||e.cue.id.startsWith(p.activeCue.id+':')))policeStaged=true;
            if(e.assassination){
              frameScene(camera,controls.target,new THREE.Box3(
                new THREE.Vector3(e.slot.root.x-.8,0,e.slot.root.z-1.2),
                new THREE.Vector3(e.slot.root.x+7.5,2.4,e.slot.root.z+1.2)));
              controls.update();
            }
            if(e.assassination&&!e.cue.id.startsWith('preview:')){
              const root={x:e.slot.root.x+ASSASSINATION_VICTIM_X,z:e.slot.root.z};
              e.assassination.update(e.assassination.duration);
              const joints=captureBodyJoints(e.assassination.victim);
              e.assassination.update(0);
              aftermath.rememberBody(e.cue.strike!.victim.id,{root,pose:{x:root.x+.8,z:root.z,heading:0},model:'casualty'},e.assassination.victimYaw,joints);
            }
          }
        }
        if(policeStaged&&p.activeCue){
          const group=effects.filter(e=>e.cue.id===p.activeCue!.id||e.cue.id.startsWith(p.activeCue!.id+':'));
          const bounds=stagedSceneBounds(group.flatMap(e=>e.slot?[e.slot]:[]));
          frameScene(camera,controls.target,bounds);controls.update();
        }
        aftermath.update(aftermathRecords, presentationMinute, lots, models, personModel,
          [...effects.flatMap(e=>e.slot?[e.slot]:[]), ...actorSpaces], animatingVictims, policeRecords,
          new Set(effects.filter(e=>['raid','raid-officer','raid-unit'].includes(e.cue.kind)).map(e=>e.cue.target)),fireRecords.filter(f=>!effects.some(e=>(e.incendiary||e.planter)&&e.cue.target===f.target)));
        const placements = traffic.update(
          [...actors]
            .filter(([id, a]) => !a.arrived&&!replacedActors.has(id))
            .map(([id, a]) => {
              const t =
                motion && a.duration ? Math.min(1, (movementClock - a.since) / a.duration) : 1;
              return {
                id,
                model: trafficModel(a.model,a.start === a.end),
                points: a.points,
                progress: a.timelineProgress ?? a.start + (a.end - a.start) * t,
              };
            })
            .concat(effects.flatMap(e => e.slot ? [{
              id: `scene:${e.cue.id}`,
              model: e.slot.model,
              points: [e.slot.pose],
              progress: 0,
            }] : [])).concat(aftermath.reservations()),
          dt / 1000,
          playback.current,
          effects.flatMap(e=>{
            if(e.slot)return [];
            if(e.driveBy)return sceneSlots(lots.get(e.cue.target)!,'driveby-building');
            if(!e.planter)return [];
            const building=buildings.get(e.cue.target);
            const entry=(building?.getObjectByName('entrance-threshold')||building?.getObjectByName('entrance-landing'))?.getWorldPosition(new THREE.Vector3());
            return entry?[planterReservation(entry)]:[];
          }),
        );
        aftermath.show(placements);
        for (const [id, a] of actors) {
          const placement = placements.get(id);
          a.object.visible = !!placement && !placement.waiting && !(id==='player'&&effects.some(e=>e.cue.kind==='detainee'&&e.cue.actors?.[0]?.id==='player'&&e.slot));
          if (!placement || placement.waiting) continue;
          const at = placement.pose;
          const distance = Math.hypot(a.object.position.x - at.x, a.object.position.z - at.z);
          const moved = distance > 0.0001;
          const walking = a.walking || !!placement.yielding;
          if(!isPedestrian(a.model)&&a.object.visible&&Math.hypot(at.x-controls.target.x,at.z-controls.target.z)<22){
            const key=`${a.since}:${a.start}:${a.end}`,previous=vehicleSounds.get(id);
            if(motion&&moved&&a.start!==a.end&&previous?.key!==key){
              previous?.stop?.();vehicleSounds.set(id,{key,stop:playRecordedEffect('vehicle-approach')});
            }else if(!moved)nearbyIdle.push('engine-idle');
          }
          if (walking && moved)
            a.phase = (a.phase + (distance / 1.15) * Math.PI * 2) % (Math.PI * 2);
          if (a.wheels.length && a.wheelPlaced && motion && a.start !== a.end && moved) {
            a.wheelPhase = advanceWheel(a.wheelPhase, distance);
            a.steering = advanceSteering(a.steering, wheelSteering(a.points,placement.progress,a.model),distance);
            for (const wheel of a.wheels) {
              wheel.rotation.x = a.wheelPhase;
              wheel.rotation.y = wheel.name.includes('-front-') ? frontWheelSteering(a.steering, a.model, wheel.position.x) : 0;
            }
          }
          a.wheelPlaced = true;
          a.object.position.set(at.x, isPedestrian(a.model) ? pedestrianRootHeight(at) : vehicleRootHeight(at), at.z);
          a.object.rotation.y = at.heading;
          if(a.driver)a.driver.visible=a.start!==a.end&&!a.arrived&&camera.zoom>=6;
          for (const limb of a.limbs) {
            const side = limb.name.endsWith('-1') ? 0 : Math.PI;
            const phase = a.phase + side;
            limb.rotation.x =
              !walking || !moved || !motion
                ? 0
                : limb.name.startsWith('knee')
                  ? Math.max(0, Math.sin(phase + 0.7)) * 0.65
                  : Math.sin(phase + (limb.name.startsWith('arm') ? Math.PI : 0)) *
                    (limb.name.startsWith('arm') ? 0.23 : 0.35);
          }
          if (walking && moved && motion)
            a.object.position.y += 0.05 + Math.abs(Math.sin(a.phase)) * 0.015;
          if (id === 'player') {
            playerRing.position.set(at.x, isPedestrian(a.model) ? surfaceHeight(at) + .045 : vehicleRootHeight(at) + .04, at.z);
            const pixel = (camera.top - camera.bottom) / (Math.max(1, element.clientHeight) * camera.zoom);
            const radius = Math.max(.58, Math.min(2.2, pixel * 7));
            if (isPedestrian(a.model)) playerRing.scale.set(radius, radius, 1);
            else {
              const size = trafficSize(a.model);
              playerRing.scale.set(Math.max(radius, size.width / 2 + .2), Math.max(radius, size.length / 2 + .25), 1);
            }
            playerRing.rotation.y = at.heading;
          }
          if (id !== 'player' && a.end === 1 && placement.progress >= 1) {
            a.arrived = true;
            a.object.visible = false;
          }
        }
        playerRing.visible = !!actors.get('player')?.object.visible;
        const arrival = placements.get('player');
        const journeyBlocked=!!p.journey&&!!arrival?.blockedBy;
        if(journeyBlocked!==reportedJourneyBlocked){reportedJourneyBlocked=journeyBlocked;p.onJourneyBlocked?.(journeyBlocked);}
        if (p.journey && arrival) {
          renderedJourneyProgress=arrival.progress;
          if(arrival.progress>=1)playedJourneyProgress=1;
          if(!p.journey.street)playedJourneyProgress = Math.max(playedJourneyProgress, arrival.progress);
          const progress = Math.floor(THREE.MathUtils.clamp(playedJourneyProgress, 0, 1) * Math.max(1, p.journey.minutes)) / Math.max(1, p.journey.minutes);
          if (progress !== reportedJourneyProgress) { reportedJourneyProgress = progress; p.onJourneyProgress?.(progress); }
        }
        if (
          p.journey &&
          completedJourney !== journeyKey &&
          (!motion || (arrival && !arrival.waiting && arrival.progress >= 1))
        ) {
          completedJourney = journeyKey;
          canvas.dataset.arrival = JSON.stringify({
            revision: w.revision,
            minute: w.minute,
            progress: arrival?.progress,
            elapsedMs: Math.round(now - (actors.get('player')?.realSince ?? now)),
            reducedMotion: !motion,
          });
          p.onJourneyDone();
        }
        const lot = lots.get(p.selected);
        selection.visible = !!lot;
        if (lot) selection.position.set(lot.x, 0.22, lot.z);
        for (const [id, label] of labels) {
          label.visible =
            camera.zoom > 2.7 || id === p.selected || id === w.player.location || id === hovered;
          const labelScale = Math.max(1, camera.zoom / 2.7);
          label.scale.set(17 / labelScale, 2.65 / labelScale, 1);
        }
        for(const [id,building] of buildings){
          const door=building.getObjectByName('entrance-door-hinge');if(!door)continue;
          const playing=effects.some(e=>e.cue.target===id&&e.cue.kind==='raid-officer');
          door.rotation.y=!playing&&policeRecords.some(p=>p.target===id&&presentationMinute>=p.minute&&presentationMinute<p.cleanup_at)?-Math.PI/2:0;
        }
        for (let i = effects.length - 1; i >= 0; i--) {
          const e = effects[i];
          if (e.extra && motion) {
            const castReady=gunCastReady(e.cue,effects.map(other=>other.cue),id=>{
              const victim=effects.find(other=>other.cue.id===id);
              const placement=placements.get(`scene:${id}`);
              return !!victim?.slot&&!!placement&&!placement.waiting;
            });
            const staged = !!e.slot && !placements.get(`scene:${e.cue.id}`)?.waiting && castReady;
            e.extra.visible = staged;
            e.mesh.visible = staged&&(!e.planter||now>=e.since);
            if((e.incendiary||e.planter||e.accident)&&e.debris)e.debris.visible=staged&&(!e.planter||now>=e.since);
            if (!staged) {
              e.since += dt;
              continue;
            }
          }
          if (e.cue.kind === 'killing') {
            const preamble=effects.find(other=>other.planter&&other.cue.target===e.cue.target&&other.cue.minute===e.cue.minute&&(!other.slot||now<other.since));
            if(preamble){e.since=now;e.mesh.visible=false;if(e.extra)e.extra.visible=false;continue;}
            const gunScene = effects.find(other => gunVictim(other.cue,e.cue));
            e.since = casualtySceneStart(e.since, now, gunScene?.since);
          }
          const t = (now - e.since) / 3000,
            lot = lots.get(e.cue.target)!;
          if (t * 3 >= (e.driveBy?.duration??e.incendiary?.duration??e.assassination?.duration??policeSceneSeconds(e.cue.kind)) || !motion) {
            handOffSurvivor(e);
            scene.remove(e.mesh, e.light);
            if (e.extra) scene.remove(e.extra);
            e.driveBy?.dispose();
            e.wardrobe?.forEach(material => material.dispose());
            e.audio?.dispose(); e.glassAudio?.dispose(); e.reactionAudio?.dispose(); e.departureAudio?.dispose();
            disposeDebris(e);
            e.mesh.dispose();
            (e.mesh.material as THREE.Material).dispose();
            effects.splice(i, 1);
            continue;
          }
          const at = e.slot?.root || entrance(lot);
          const blast = e.cue.kind === 'explosion';
          const blastBuilding=buildings.get(lot.id);
          const internal=internalDetonation(e.cue,w.building_fires||[]);
          const blastOrigin = e.accident?{x:at.x,y:.95,z:at.z}:internal?blastBuilding?.userData.blastOrigin:blastBuilding?.userData.debrisOrigin;
          const blastWindows = (internal?blastBuilding?.userData.blastWindows || []:[]) as THREE.Vector3[];
          const debrisOrigin=e.accident?{x:at.x,y:.2,z:at.z}:blastBuilding?.userData.debrisOrigin;
          if (blast && blastOrigin) e.light.position.set(blastOrigin.x, blastOrigin.y, blastOrigin.z);
          const shot = (e.cue.kind === 'gunfight'||!!e.driveBy) && !!e.weapon;
          if(e.planter){
            const pose=e.planter.update(t*3+e.planter.duration);
            if(e.pullback){
              if(cameraIntent!==e.pullback.intent||followPlayer.current)e.pullback.move.cancel();
              if(e.pullback.move.update(camera,controls.target,(t*3+1.6)/1.4))controls.update();
            }
            const door=blastBuilding?.getObjectByName('entrance-door-hinge');if(door)door.rotation.y=-Math.PI/2*pose.door;
          }
          e.accident?.update(t*3);
          e.assassination?.update(t*3);
          e.driveBy?.update(t*3);
          e.departureAudio?.update(t*3-3.8,soundOn());
          e.custody?.update(t*3);
          e.incendiary?.update(t*3);
          if(e.incendiary){e.glassAudio?.update(t*3-INCENDIARY_IMPACT,soundOn());addImpact(t*3-INCENDIARY_IMPACT,1.4);}
          const firing = gunfightPose(t * 3,e.driveBy?.shots??weaponShots(e.weaponModel,e.cue.strike?.variant));
          if(blast)addImpact(t*3,11);
          if(shot)for(const beat of e.driveBy?.shots??weaponShots(e.weaponModel,e.cue.strike?.variant))addImpact(t*3-beat,3);
          if(e.cue.kind!=='raid-officer')e.audio?.update(t * 3, soundOn());
          e.reactionAudio?.update(t*3-(e.cue.kind==='attack'?MELEE_IMPACTS[2]:e.assassination?ASSASSINATION_SHOT:blast?1:.06),soundOn());
          if(e.cue.kind==='attack'&&e.assassination)for(const beat of MELEE_IMPACTS)addImpact(t*3-beat,1.2);
          const muzzlePosition = new THREE.Vector3(at.x, 1.4, at.z);
          if (shot && e.gunArm && e.muzzle) {
            if(e.assassination||e.driveBy){ /* the shared cast owns its arm and weapon rig */ }
            else if(e.weaponModel==='revolver')e.gunArm.rotation.x=firing.arm;
            else {const pump=e.weapon?.getObjectByName('pump-slide');if(pump)pump.position.z=pumpOffset(t*3);poseLongGun(e.extra!,e.weapon!,firing.arm);}
            e.extra!.updateMatrixWorld(true);
            e.muzzle.getWorldPosition(muzzlePosition);
            e.light.position.copy(muzzlePosition);
          }
          const casualty = e.cue.kind === 'killing';
          if (casualty && e.extra) {
            const fall = casualtyFall(t);
            e.extra.rotation.z = fall.rotation;
            e.extra.position.y = fall.height;
            groundCharacter(e.extra,.205);
          }
          const police = ['raid', 'arrest','raid-unit','police-unit'].includes(e.cue.kind);
          const personnel=['officer','detainee','raid-officer'].includes(e.cue.kind);
          if(personnel&&e.extra){
            if(e.cue.kind==='raid-officer'){
              const front=(buildings.get(lot.id)?.userData.sightBounds as THREE.Box3|undefined)?.min.z ?? at.z;
              const distance=THREE.MathUtils.clamp(front-at.z-.7,0,2.4);
              const building=buildings.get(lot.id);
              const threshold=building?.getObjectByName('entrance-threshold');
              const door=building?.getObjectByName('entrance-door-hinge');
              const entry=threshold?.getWorldPosition(new THREE.Vector3());
              const approach=entry?entry.z-at.z-.65:0;
              const breaching=!!door&&!!entry&&Math.abs(entry.x-at.x)<.05&&approach>0&&approach<=3;
              const breach=breaching?raidEntryPose(t*3,approach):null;
              const walk=breach || officerApproach(t*3,distance,Number(e.cue.id.split(':').at(-1))||0);
              if(breach&&door){
                door.rotation.y=-Math.PI/2*breach.door;
                const impactAge=t*3-(.35+approach/1.4+.3);
                e.audio?.update(impactAge,soundOn());
                addImpact(impactAge,5);
              }
              e.extra.position.z=at.z+walk.travelled;e.extra.rotation.y=0;
              e.extra.position.y=.2+(walk.walking?Math.abs(Math.sin(walk.phase))*.025:0);
              for(const name of ['leg1','leg-1','knee1','knee-1','arm1','arm-1']){
                const limb=e.extra.getObjectByName(name);if(!limb)continue;
                const phase=walk.phase+(name.endsWith('-1')?0:Math.PI);
                limb.rotation.x=!walk.walking?0:name.startsWith('knee')?Math.max(0,Math.sin(phase+.7))*.65
                  :Math.sin(phase+(name.startsWith('arm')?Math.PI:0))*(name.startsWith('arm')?.23:.35);
              }
              if(breach){
                const leg=e.extra.getObjectByName('leg1'), knee=e.extra.getObjectByName('knee1');
                if(leg)leg.rotation.x-=breach.kick*.9;
                if(knee)knee.rotation.x+=breach.kick*.5;
              }
            }else if(e.cue.kind==='detainee'){ /* Shared custody rig owns both participants. */ }
            else e.extra.rotation.y=Math.atan2(lot.x-at.x,lot.row*PITCH+6.35-at.z);
          }
          for (let j = 0; j < 32; j++) {
            const a = j * 2.399;
            const r = shot ? 0.35 : 1.8;
            tmp.position.set(at.x + Math.cos(a) * r, 1 + Math.sin(a) * r * 0.4, at.z + Math.sin(a) * r);
            // Saved window fire owns the incendiary effect; do not add the
            // generic incident's floating ring of particles at the doorway.
            tmp.scale.setScalar(casualty || shot || personnel || e.assassination || e.cue.kind === 'incendiary' ? 0.001 : 0.25);
            const burst = blast ? blastParticle(j, t * 3) : null;
            if (burst && blastOrigin) {
              const window=blastWindows[j%blastWindows.length];
              const vent=window?windowBurst(j,t*3,window):null;
              tmp.position.set(vent?.x ?? blastOrigin.x+burst.x,vent?.y ?? blastOrigin.y+burst.y,vent?.z ?? blastOrigin.z+burst.z);
              tmp.scale.setScalar(Math.max(.001,vent?.size ?? burst.size));
            }
            if(e.incendiary&&e.incendiary.bottle.visible&&j<2){
              e.incendiary.bottle.getObjectByName('bottle-flame')?.getWorldPosition(tmp.position);
              tmp.position.y+=j*.05;
              const flicker=.075+.025*Math.sin(t*93+j*2);
              tmp.scale.set(flicker,flicker*1.7,1);
              if(j===0)e.light.position.copy(tmp.position);
            }
            if (shot) {
              tmp.position.copy(muzzlePosition);
              const smoke = j === 1 && firing.smoke > 0;
              if (smoke) tmp.position.y += (1 - firing.smoke) * 0.3;
              tmp.scale.setScalar(j === 0 && firing.flash ? 0.28 : smoke ? 0.15 + (1 - firing.smoke) * 0.3 : 0.001);
              if(e.assassination&&j>=2){const drop=executionSpatter(j-2,t*3,e.cue.strike?.variant==='back-of-head');
                tmp.position.set(at.x+drop.x,.2+drop.y,at.z+drop.z);tmp.scale.setScalar(Math.max(.001,drop.size));}

            }
            if(e.driveBy&&j>=2){
              const beat=[...e.driveBy.shots].reverse().find(at=>t*3>=at),age=beat===undefined?-1:t*3-beat;
              const hit=e.driveBy.root.localToWorld(e.driveBy.target.clone());
              const drift=Math.max(0,age);
              tmp.position.set(hit.x+Math.cos(a)*drift*.7,hit.y+Math.sin(a)*drift*.45+.25*drift,hit.z-.12-drift*(.35+(j%4)*.1));
              tmp.scale.setScalar(age>=0&&age<.48?(.07+(j%3)*.025)*(1-age/.48):.001);
            }
            if (police) {
              // A period rotating red roof beacon, rather than sparks around the car.
              tmp.position.set(at.x, (e.extra?.position.y ?? 0.2) + 1.76, at.z);
              tmp.scale.setScalar(j === 0 ? 0.45 + 0.35 * Math.max(0, Math.sin(t * 38)) : 0.001);
            }
            tmp.quaternion.copy(camera.quaternion);
            tmp.updateMatrix();
            e.mesh.setMatrixAt(j, tmp.matrix);
            e.mesh.setColorAt(
              j,
              new THREE.Color(
                e.driveBy&&j>=2?0xa99d85:e.incendiary ? (j?0xff7b23:0xffd37b) : e.assassination&&j>=2 ? 0x720c12 : burst ? burst.color : (shot && j === 1)
                  ? 0x55534e
                  : police
                    ? 0xde3426
                    : !blast && !shot
                      ? 0xbacbd0
                      : j % 3
                        ? 0xff9c2f
                        : 0xffefb8,
              ),
            );
          }
          if(e.debris&&e.incendiary){
            const impact=e.incendiary.root.localToWorld(e.incendiary.target.clone());
            const age=t*3-INCENDIARY_IMPACT;
            for(let j=0;j<12;j++){
              const floor=surfaceHeight({x:impact.x,z:impact.z-1})+.07;
              const height=impact.y-floor;
              const piece=incendiaryShard(j,age,height);
              const to=impact.clone().add(new THREE.Vector3(piece.x,piece.y,piece.z));
              if(piece.scale>0&&!e.glassBlocked!.has(j)&&blastBuilding&&
                incendiaryShardObstructed(j,age-dt/1000,age,height,impact,blastBuilding))e.glassBlocked!.add(j);
              tmp.position.copy(to);tmp.rotation.set(piece.rx,piece.ry,piece.rz);
              tmp.scale.setScalar(e.glassBlocked!.has(j)?0:piece.scale);tmp.updateMatrix();e.debris.setMatrixAt(j,tmp.matrix);
            }
            e.debris.instanceMatrix.needsUpdate=true;
          }
          if (e.debris && debrisOrigin && !e.incendiary) {
            for (let j = 0; j < 12; j++) {
              const window=blastWindows[j%blastWindows.length];
              const airborne=window?windowDebris(j,t*3,window,surfaceHeight({x:window.x,z:window.z-3}),blastWindows.length) : null;
              const fragment = airborne || debrisPose(j, t * 3);
              const x = airborne?fragment.x:debrisOrigin.x+fragment.x, z = airborne?fragment.z:debrisOrigin.z+fragment.z;
              let blocked = false;
              for (const actor of actors.values()) {
                if (!actor.object.visible) continue;
                const size = trafficSize(actor.model);
                if (fragmentBlocked(x,z,{x:actor.object.position.x,z:actor.object.position.z,heading:actor.object.rotation.y},size.width,size.length)) blocked = true;
              }
              for (const other of effects) {
                if (!other.slot || !other.extra?.visible) continue;
                const size = trafficSize(other.slot.model);
                if (fragmentBlocked(x,z,other.slot.pose,size.width,size.length)) blocked = true;
              }
              tmp.position.set(x, 0, z);
              tmp.rotation.set(fragment.rx, fragment.ry, fragment.rz);
              tmp.scale.setScalar(blocked ? .001 : Math.max(.001, fragment.scale));
              tmp.updateMatrix();
              // Support the rotated fragment's exported half extents on the surface.
              const basis = tmp.matrix.elements;
              const support = .09*Math.abs(basis[1]) + .04*Math.abs(basis[5]) + .045*Math.abs(basis[9]);
              tmp.position.y = surfaceHeight({x,z}) + support + .006 + fragment.height;
              tmp.updateMatrix();
              e.debris.setMatrixAt(j,tmp.matrix);
            }
            e.debris.instanceMatrix.needsUpdate = true;
            (e.debris.material as THREE.MeshStandardMaterial).opacity = blastWindows.length&&w.building_fires?.some(f=>f.target===lot.id&&w.minute<f.cleanup_at)?1:blastOpacity(t*3)/.8;
          }
          e.mesh.instanceMatrix.needsUpdate = true;
          if (e.mesh.instanceColor) e.mesh.instanceColor.needsUpdate = true;
          (e.mesh.material as THREE.MeshBasicMaterial).opacity = blast ? blastOpacity(t * 3) : e.assassination||e.incendiary||e.driveBy?1:Math.max(0,1 - t*3/policeSceneSeconds(e.cue.kind));
          e.light.intensity = e.planter&&t<0?0:e.incendiary ? (e.incendiary.bottle.visible?1.5:0) : blast
            ? blastLight(t * 3)
            : shot && firing.flash
              ? 30
              : police
                ? 8 * Math.max(0, Math.sin(t * 38))
                : 0;
          if (police) {
            e.light.color.setHex(0xe53220);
            e.light.position.y = (e.extra?.position.y ?? 0.2) + 1.76;
          }
        }
      }
      const stagedAttackers=new Set(effects.filter(e=>e.cue.attacker&&e.extra?.visible).map(e=>e.cue.attacker!.id));
      for(const e of effects)if(e.cue.attacker&&e.extra?.visible){
        const actor=actors.get(e.cue.attacker.id);if(actor)actor.object.visible=false;
        if(e.cue.attacker.id==='player')playerRing.visible=false;
      }
      if (!p.activeCue) completedScene = '';
      if (ready && p.activeCue && effects.length===0 && completedScene!==activePlayback) {
        completedScene=activePlayback!;p.onSceneDone?.(p.activeCue.id);
      }
      if (ready) lastActive = activePlayback;
      motionWas = motion;
      controls.enableDamping = motion;
      controls.update();
      const followed = actors.get('player');
      if (followPlayer.current && !wasFollowing) followZoom = 8;
      wasFollowing = followPlayer.current;
      if (followPlayer.current && ready && followed?.object.visible) {
        const ease = motion ? 1 - Math.exp(-Math.min(dt, 100) / 170) : 1;
        const dx = (followed.object.position.x - controls.target.x) * ease;
        const dz = (followed.object.position.z - controls.target.z) * ease;
        if (followZoom !== null) {
          camera.zoom += (followZoom - camera.zoom) * ease;
          if (Math.abs(followZoom - camera.zoom) < .005) followZoom = null;
          camera.updateProjectionMatrix();
        }
        camera.position.x += dx; camera.position.z += dz;
        controls.target.x += dx; controls.target.z += dz;
      }
      // Keep pan within the city plus its waterfront margin.
      const before = controls.target.clone();
      controls.target.x = THREE.MathUtils.clamp(controls.target.x, harbourLot ? -45 : -15, plan.width + 15);
      controls.target.z = THREE.MathUtils.clamp(controls.target.z, -15, plan.depth + 15);
      camera.position.add(controls.target.clone().sub(before));
      // Keep the followed player or the active staged cast readable through buildings.
      camera.updateMatrixWorld();
      const eventTarget = p.activeCue?.target || effects.find(e => e.extra?.visible)?.cue.target;
      const sightPoints = ready && followPlayer.current && followed?.object.visible
        ? [followed.object.position.clone().add(new THREE.Vector3(0, .9, 0))]
        : ready ? effects.filter(e => e.extra?.visible && e.slot && e.cue.target === eventTarget)
          .flatMap(e => e.driveBy?[e.driveBy.driver,e.driveBy.shooter].map(a=>a.getWorldPosition(new THREE.Vector3()).add(new THREE.Vector3(0,1.35,0))):e.planter?[e.planter.actor.getWorldPosition(new THREE.Vector3()).add(new THREE.Vector3(0,1.35,0))]:e.incendiary?[e.incendiary.actor.getWorldPosition(new THREE.Vector3()).add(new THREE.Vector3(0,1.35,0))]:e.custody?[e.custody.officer,e.custody.detainee].map(a=>a.getWorldPosition(new THREE.Vector3()).add(new THREE.Vector3(0,.9,0))):e.assassination?[e.assassination.attacker,e.assassination.victim].map(a=>a.getWorldPosition(new THREE.Vector3()).add(new THREE.Vector3(0,.9,0))):[e.extra!.position.clone().add(new THREE.Vector3(0, .9, 0))]) : [];
      let cutawayDepth=0;
      if (sightPoints.length) {
        const sightEnvelope=sightPoints.flatMap(sight => characterSightPoints(camera,sight));
        if (now - lastSightCheck >= 100) {
          blockers = new Set(sightEnvelope
            .flatMap(sight => [...blockingBuildings(camera, sight, buildings)]));
          // An authored doorway should occlude an entering officer naturally;
          // dissolving the target would erase the door being breached.
          for(const e of effects)if(e.cue.kind==='raid-officer'&&buildings.get(e.cue.target)?.getObjectByName('entrance-door-hinge'))blockers.delete(e.cue.target);
          lastSightCheck = now;
        }
        const size = renderer.getDrawingBufferSize(new THREE.Vector2());
        const projected = sightEnvelope.map(sight => {
          const screen = sight.project(camera);
          cutawayDepth=Math.max(cutawayDepth,(screen.z+1)/2);
          return new THREE.Vector2((screen.x + 1) * size.x / 2, (screen.y + 1) * size.y / 2);
        });
        const centre = projected.reduce((sum, point) => sum.add(point), new THREE.Vector2()).divideScalar(projected.length);
        const radius = Math.max(...projected.map(point => point.distanceTo(centre)))
          + (65 + camera.zoom * 2) * renderer.getPixelRatio();
        cutawayWindow.set(centre.x, centre.y, radius);
      } else { blockers.clear(); lastSightCheck = -Infinity; }
      // Repeated trees and street furniture keep their instanced draw calls.
      // Only fragments in front of the watched cast dissolve in the opening.
      for(const material of streetCutawayMaterials)
        buildingCutaway(material,sightPoints.length?1:0,cutawayWindow,cutawayDepth);
      for (const [id, building] of buildings) {
        const desired = blockers.has(id) ? 1 : 0, previous = cutawayAmounts.get(id) || 0;
        const amount = motion ? THREE.MathUtils.lerp(previous, desired, 1 - Math.exp(-Math.min(dt, 100) / 90)) : desired;
        const settled = Math.abs(amount - desired) < .001 ? desired : amount;
        cutawayAmounts.set(id, settled);
        if (!settled && !previous) continue;
        building.traverse(part => {
          if (!(part instanceof THREE.Mesh)) return;
          for (const material of Array.isArray(part.material) ? part.material : [part.material])
            if (material instanceof THREE.MeshStandardMaterial) buildingCutaway(material, settled, cutawayWindow);
        });
      }
      if (actors.size * 2 > headlightPools.instanceMatrix.count) {
        scene.remove(headlightPools); headlightPools.dispose();
        headlightPools = new THREE.InstancedMesh(headlightGeometry, headlightMaterial, actors.size * 4);
        headlightPools.frustumCulled = false; scene.add(headlightPools);
      }
      headlightPools.count = 0;
      const nightAmount=cityNightAmount(presentationMinute), lampsOn=nightAmount>0;
      const nextLightingKey=JSON.stringify(w.sky);
      const lightingChanged=lightingKey!==nextLightingKey;
      if(ready&&(lightingChanged||lastNightAmount!==nightAmount)){
        // Snapshot changes can replace damage/cutaway materials. Cache their new
        // references once; twilight itself must not traverse every building.
        if(lightingChanged){
          luminousMaterials.clear();
          for(const b of buildings.values())b.traverse(o=>{
            if(o instanceof THREE.Mesh)for(const m of Array.isArray(o.material)?o.material:[o.material])
              if(m instanceof THREE.MeshStandardMaterial&&m.emissive.getHex()!==0)luminousMaterials.add(m);
          });
        }
        lightingKey=nextLightingKey;lastNightAmount=nightAmount;
        const weather=cityWeather(w.sky,nightAmount);
        sky.intensity=weather.ambient;sun.intensity=weather.sun;
        groundMat.color.setHex(0x646460).multiplyScalar(weather.roadTone);
        pavementMat.color.setHex(0xaaa18b).multiplyScalar(weather.pavementTone);
        groundMat.roughness=weather.roadRoughness;pavementMat.roughness=weather.pavementRoughness;
        pools.visible=lampsOn;pools.material.opacity=.24*nightAmount;
        headlightMaterial.opacity=.28*nightAmount;
        for(const material of luminousMaterials)material.emissiveIntensity=.2+nightAmount;
        if(scene.background instanceof THREE.Color)scene.background.setHex(weather.background);
        else scene.background=new THREE.Color(weather.background);
        if(scene.fog instanceof THREE.Fog){scene.fog.color.setHex(weather.background);scene.fog.far=weather.fogFar;}
        else scene.fog=new THREE.Fog(weather.background,260,weather.fogFar);
        if(lightingChanged)renderer.shadowMap.needsUpdate=true;
      }
      for(const hand of clockHands){
        const period=hand.name==='clock-hand-minute'?60:720;
        hand.rotation.z=(presentationMinute%period)/period*Math.PI*2;
      }
      for (const actor of actors.values()) {
        const lit = ready && lampsOn && actor.object.visible && actor.points.length > 1;
        for (const material of actor.lamps) material.emissiveIntensity = lit ? (material.name === 'headlamps' ? 1.6 : .8)*nightAmount : 0;
        if (!lit || !actor.lamps.length) continue;
        for (const side of [-1, 1]) {
          const at = headlightCentre(actor.object.position.x, actor.object.position.z, actor.object.rotation.y, trafficSize(actor.model).length, side);
          headlightTransform.position.set(at.x, actor.object.position.y + .02, at.z);
          headlightTransform.rotation.set(-Math.PI / 2, actor.object.rotation.y, 0, 'YXZ'); headlightTransform.updateMatrix();
          headlightPools.setMatrixAt(headlightPools.count++, headlightTransform.matrix);
        }
      }
      if (headlightPools.count) headlightPools.instanceMatrix.needsUpdate = true;
      rainfall.visible = ready && motion && w.sky?.kind === 'rain';
      if (rainfall.visible) {
        rainClock += Math.min(dt, 100) / 1000;
        rainVertices(rainPositions, rainClock, plan.width, plan.depth);
        rainAttribute.needsUpdate = true;
      }
      if (ready && motion && waterNormal) {
        waterClock += Math.min(dt, 100) / 1000;
        waterNormal.offset.set((waterClock * .006) % 1, (waterClock * .003) % 1);
      }
      for(const [id,b] of buildings){
        const condition=b.userData.condition??100;
        const blast=effects.find(e=>e.cue.kind==='explosion'&&e.cue.target===id&&internalDetonation(e.cue,w.building_fires||[]));
        const before=blast?.glazingBefore??condition,age=blast?(now-blast.since)/1000:0;
        const preview=!!blast?.cue.id.startsWith('preview:');
        buildingGlazing(b,blast?glazingDuringBlast(condition,before,age,preview):condition);
        if(blast&&before>=60&&(preview||condition<60)&&b.getObjectByName('window-broken'))blast.glassAudio?.update(age-GLASS_BREAK_AT,soundOn());
      }
      const revealedFires=fireRecords.filter(f=>!effects.some(e=>e.cue.target===f.target&&((e.incendiary&&(!e.slot||now-e.since<INCENDIARY_IMPACT*1000))||(e.planter&&(!e.slot||now<e.since)))));
      const responseFires=revealedFires.filter(f=>!effects.some(e=>(e.incendiary||e.planter)&&e.cue.target===f.target));
      rubble.update(revealedFires,presentationMinute,buildings,models.get('blast-fragment'),new Set(effects.filter(e=>e.cue.kind==='explosion'&&now-e.since<3000).map(e=>e.cue.target)));
      suppression.update(responseFires,presentationMinute,buildings,models,aftermath,dt,motion);
      buildingFire.update(revealedFires,presentationMinute,buildings,camera,dt,motion);
      const nearFire=revealedFires.some(f=>{
        const b=buildings.get(f.target);
        return b&&presentationMinute>=f.minute&&presentationMinute<f.extinguished_at&&Math.hypot(b.position.x-controls.target.x,b.position.z-controls.target.z)<28;
      });
      ambientAudio.update(graphicsLost?[]:[...(nearFire?['fire']:[]),...nearbyIdle]);
      for(const [id,voice] of vehicleSounds)if(!actors.has(id)){voice.stop?.();vehicleSounds.delete(id);}
      renderImpact(camera,motion?impact.x:0,motion?impact.y:0,canvas.clientWidth,canvas.clientHeight,()=>renderer.render(scene,camera));
      if (ready)
        canvas.dataset.presentation = JSON.stringify({
          revision: w.revision,
          minute: w.minute,
          playbackRate: playback.current,
          gunAudio:cityGunshotStatus(),
          effectAudio:cityEffectStatus(),
          impact:motion?impact:{x:0,y:0},
          suppression:suppression.inspect(),
          rubble:rubble.inspect(),
          glazing:[...buildings].flatMap(([id,b])=>{const broken=b.getObjectByName('window-broken');return broken?[{id,broken:broken.visible}]:[];}),
          buildingFire:buildingFire.inspect(),
          headlightPools: headlightPools.count,
          harbour: {visible: !!harbourLot, waterClock},
          followingPlayer: followPlayer.current,
          journeyProgress: {clock:playedJourneyProgress,player:renderedJourneyProgress},
          settlingStreet,
          trafficBlocks: [...actors.keys()].flatMap(id=>{const by=traffic.placement(id)?.blockedBy;return by?[{id,by}]:[]}),
          streetMinute:p.journey?(p.journey.fromMinute??w.minute-p.journey.minutes)+p.journey.minutes*playedJourneyProgress:w.minute,
          cutawayBuildings: [...blockers],
          weather: {night:nightAmount>=.5,nightAmount,sunIntensity:sun.intensity,ambientIntensity:sky.intensity,lightingMinute:presentationMinute,clockHands:clockHands.map(h=>({name:h.name,angle:h.rotation.z})),kind: w.sky?.kind || 'clear', wet: w.sky?.wet || 0, rainVisible: rainfall.visible, rainClock},
          camera: {zoom: camera.zoom, x: camera.position.x, z: camera.position.z,
            targetX: controls.target.x, targetZ: controls.target.z},
          actors: [...actors]
            .filter(([, a]) => a.object.visible)
            .map(([id, a]) => ({
              id,
              model: a.model,
              x: a.object.position.x,
              z: a.object.position.z,
              y: a.object.position.y,
              driverVisible:a.driver?.visible,
              wheelPhase: a.wheels.length ? a.wheelPhase : undefined,
              steering: a.wheels.length ? a.steering : undefined,
              lampsOn: a.lamps.length ? a.lamps.some(m => m.emissiveIntensity > 0) : undefined,
              frontWheels: a.wheels.length ? a.wheels.filter(w => w.name.includes('-front-')).map(w => ({x: w.position.x, angle: w.rotation.y})) : undefined,
            })),
          waiting: [...actors].filter(([id, a]) => !a.arrived && !a.object.visible&&!stagedAttackers.has(id)).map(([id]) => id),
          aftermath: aftermath.inspect(),
          doors:[...buildings].flatMap(([id,b])=>{const door=b.getObjectByName('entrance-door-hinge');return door?[{id,angle:door.rotation.y}]:[];}),
          effects: effects.map(e => ({
            id: e.cue.id, kind: e.cue.kind, target: e.cue.target,
            driveBy:e.driveBy?{seconds:(now-e.since)/1000,car:e.driveBy.car.getWorldPosition(new THREE.Vector3()),shots:e.audio?.started}:undefined,
            accident:e.accident?{fatal:e.accident.fatal,seconds:(now-e.since)/1000,rotation:e.accident.actor.rotation.x}:undefined,
            staged: !e.extra || e.extra.visible, x: e.slot?.root.x, z: e.slot?.root.z,
            reservation:e.slot?{model:e.slot.model,authored:e.slot.pose,admitted:traffic.placement(`scene:${e.cue.id}`)?.pose}:undefined,
            approach: e.cue.kind==='raid-officer'&&e.extra?{x:e.extra.position.x,z:e.extra.position.z,leg:e.extra.getObjectByName('leg1')?.rotation.x}:undefined,
            planter:e.planter?{seconds:(now-e.since)/1000+e.planter.duration,blastSeconds:(now-e.since)/1000,actor:e.planter.actor.getWorldPosition(new THREE.Vector3())}:undefined,
            debris: e.debris?.count,
            glassShards:e.incendiary&&e.debris?{visible:e.debris.visible,blocked:e.glassBlocked?.size,
              shown:Array.from({length:12},(_,i)=>Math.hypot(...Array.from(e.debris!.instanceMatrix.array.slice(i*16,i*16+3)))>.001).filter(Boolean).length}:undefined,
            blastOrigin: e.accident&&e.slot?{x:e.slot.root.x,y:.95,z:e.slot.root.z}:e.cue.kind === 'explosion' ? buildings.get(e.cue.target)?.userData[internalDetonation(e.cue,w.building_fires||[])?'blastOrigin':'debrisOrigin'] : undefined,
            blastWindows: e.cue.kind === 'explosion' && internalDetonation(e.cue,w.building_fires||[]) ? buildings.get(e.cue.target)?.userData.blastWindows : undefined,
            arm: e.gunArm?.rotation.x,
            incendiary:e.incendiary?{seconds:(now-e.since)/1000,actor:e.incendiary.actor.getWorldPosition(new THREE.Vector3()),bottle:e.incendiary.bottle.getWorldPosition(new THREE.Vector3()),held:(now-e.since)<3100,visible:e.incendiary.bottle.visible,target:e.incendiary.target,loft:e.incendiary.loft}:undefined,
            custody:e.custody?{officer:e.custody.officer.getWorldPosition(new THREE.Vector3()),detainee:e.custody.detainee.getWorldPosition(new THREE.Vector3()),restrained:e.custody.detainee.getObjectByName('custody-restraint')?.visible}:undefined,
            execution:e.assassination?{seconds:(now-e.since)/1000,attacker:e.assassination.attacker.getWorldPosition(new THREE.Vector3()),victim:e.assassination.victim.getWorldPosition(new THREE.Vector3()),muzzle:e.muzzle?.getWorldPosition(new THREE.Vector3())}:undefined,
            attacker:e.cue.attacker,weapon:e.cue.kind==='gunfight'?e.weaponModel:undefined,
            audioShots: e.cue.kind === 'gunfight' ? e.audio?.started : undefined,
            audioBreaches: e.cue.kind === 'raid-officer' ? e.audio?.started : undefined,
            glassBreaks:e.glassAudio?.started,
            audioBlasts: e.cue.kind === 'explosion' ? e.audio?.started : undefined,
            fall: e.cue.kind === 'killing' ? e.extra?.rotation.z : undefined,
          })),
        });
      if (ready && dt > 0 && dt < 250) samples.push(dt);
      if (now - sampleStart > 2500 && samples.length) {
        samples.sort((a, b) => a - b);
        const median = samples[Math.floor(samples.length / 2)];
        const p95 = samples[Math.floor(samples.length * 0.95)];
        const metrics = {
          fps: Math.round(1000 / median),
          p95ms: Math.round(p95 * 10) / 10,
          drawCalls: renderer.info.render.calls,
          triangles: renderer.info.render.triangles,
          actors: actors.size,
          buildings: buildings.size,
          geometries: renderer.info.memory.geometries,
          textures: renderer.info.memory.textures,
        };
        canvas.dataset.metrics = JSON.stringify(metrics);
        setFps(`${metrics.fps} FPS · ${metrics.drawCalls} draws`);
        setWaiting([...actors].filter(([id,a]) => !a.arrived && !a.object.visible&&!stagedAttackers.has(id)).length);
        samples = [];
        sampleStart = now;
      }
    };
    frame = requestAnimationFrame(tick);
    const lost = (e: Event) => {
      e.preventDefault();
      graphicsLost = true;
      silenceAmbient();
      effects.forEach(effect => {effect.audio?.dispose();effect.glassAudio?.dispose(); effect.reactionAudio?.dispose(); effect.departureAudio?.dispose();});
      setFailure('The graphics context was interrupted. Reload the city to restore it.');
    };
    canvas.addEventListener('webglcontextlost', lost);
    return () => {
      dead = true;
      silenceAmbient();
      document.removeEventListener('visibilitychange',hiddenAudio);
      cancelAnimationFrame(frame);
      observer.disconnect();
      controls.dispose();
      controlsRef.current = null;
      canvas.removeEventListener('pointerdown', pointerDown);
      canvas.removeEventListener('pointermove', pointerMove);
      canvas.removeEventListener('pointerup', pointerUp);
      canvas.removeEventListener('wheel',manualZoom);
      unbindPan();
      canvas.removeEventListener('keydown', keys);
      canvas.removeEventListener('webglcontextlost', lost);
      effects.forEach(effect => { effect.audio?.dispose(); effect.glassAudio?.dispose(); effect.reactionAudio?.dispose(); effect.departureAudio?.dispose(); effect.driveBy?.dispose(); disposeDebris(effect); });
      rubble.dispose();
      suppression.dispose();
      buildingFire.dispose();
      aftermath.dispose();
      disposeCityResources([scene, ...models.values()], {
        textures, geometries: [effectGeometry, contactGeometry], materials: [contactMaterial],
      });
      renderer.dispose();
      renderer.forceContextLoss();
      canvas.remove();
    };
  }, []);
  const place = props.state.locations.find(l => l.id === props.selected);
  const travel = place?.actions.find(a => a.id === 'travel');
  return (
    <section
      className={'city3d' + (props.immersive?' city3d-immersive':expanded ? ' city3d-expanded' : '')}
      aria-label="Bellwether city"
      role={expanded&&!props.immersive ? 'dialog' : 'region'}
      aria-modal={expanded&&!props.immersive || undefined}
      onKeyDown={e => {
        if (props.immersive || !expanded || e.ctrlKey || e.metaKey || e.altKey || e.nativeEvent.isComposing) return;
        if (e.key === 'Escape') {
          e.preventDefault();
          e.stopPropagation();
          setExpanded(false);
          expandButton.current?.focus();
          return;
        }
        if (e.key === 'Tab') {
          const items = Array.from(
            e.currentTarget.querySelectorAll<HTMLElement>('button:not(:disabled),select,canvas'),
          );
          const at = items.indexOf(document.activeElement as HTMLElement);
          if (e.shiftKey && at <= 0) {
            e.preventDefault();
            items.at(-1)?.focus();
          } else if (!e.shiftKey && (at < 0 || at === items.length - 1)) {
            e.preventDefault();
            items[0]?.focus();
          }
        }
      }}
    >
      <div className="city3d-canvas" ref={host} />
      <div className="city3d-top">
        <div className="city3d-heading">
          <small>BELLWETHER · 1950</small>
          <span>Drag to rotate · Right drag to pan · Scroll to zoom</span>
          <span>WASD / Arrows: pan · Q/E: rotate · +/−: zoom · Home: reset · Esc: return</span>
        </div>
        <div className="city3d-tools">
          {!props.immersive&&<button ref={expandButton} aria-pressed={expanded} onClick={() => {
            setExpanded(!expanded);
            if (!expanded) host.current?.querySelector('canvas')?.focus();
          }}>
            {expanded ? 'Return to game' : 'Expand city'}
          </button>}
          <button data-shortcut="o" aria-keyshortcuts="O" title="Whole city (O)" onClick={() => focus.current()}>Whole city · O</button>
          <button data-shortcut="f" aria-keyshortcuts="F" aria-pressed={following} title={following ? "Stop following your character" : "Follow your character or car"}
            onClick={() => setFollow(!followPlayer.current)}>{following ? "Stop following" : "Find me"}</button>
          <button data-shortcut="z" aria-keyshortcuts="Z" title="Focus address (Z)" onClick={() => focus.current(props.selected)}>Focus address · Z</button>
          <button
            data-shortcut="t" aria-keyshortcuts="T"
            onClick={() => setPlaybackRate(rate => (rate === 1 ? 4 : 1))}
            aria-label={`Travel playback speed: ${playbackRate} times`}
          >
            Travel {playbackRate}×
          </button>
        </div>
        {debug && <div className="city3d-preview" aria-label="Scene preview controls">
          <label>Debug scene <select aria-label="Debug scene" value={previewKind}
            onChange={e => setPreviewKind(e.target.value as PreviewScene)}>
            {previewScenes.map(kind => <option key={kind}>{kind}</option>)}
          </select></label>
          <button disabled={!canPreview || !!status || !!failure} onClick={() => {
            setPreview(previewScene(props.state, props.selected, previewKind, crypto.randomUUID()));
          }}>Play preview</button>
          <button disabled={!preview} onClick={() => setPreview(null)}>Stop preview</button>
          <small>{shownPreview ? 'VISUAL PREVIEW · Campaign unchanged' : !props.motion ? 'Enable City animation in Settings' : 'At selected address'}</small>
        </div>}
        {props.overlay && <div className="city3d-story">{props.overlay}</div>}
      </div>
      {(status || failure) && (
        <p className="city3d-status" role="status">
          {failure || status}
        </p>
      )}
      {waiting > 0 && (
        <p className="city3d-status" role="status">
          {waiting} travellers waiting for space at their departures.
        </p>
      )}
      {place && (
        <div className="city3d-address">
          <small>
            {place.owned
              ? 'YOUR PREMISES'
              : place.locked
                ? 'BEYOND YOUR REACH'
                : 'SELECTED ADDRESS'}
          </small>
          <strong>{place.name}</strong>
          <span>{place.blurb}</span>
          {place.id === props.state.player.location ? (
            <button data-shortcut="g" aria-keyshortcuts="G" title="Step inside (G)" disabled={!!shownPreview||props.busy||!!props.journey} onClick={props.onEnter}>Step inside · G →</button>
          ) : (
            <button
              data-shortcut="g" aria-keyshortcuts="G"
              disabled={!travel || travel.disabled || props.busy || !!props.journey || !!shownPreview}
              title={travel?.reason}
              onClick={() => props.onTravel(place.id)}
            >
              {travel?.disabled ? travel.reason : `Go there · ${travel?.minutes ?? '—'} min →`}
            </button>
          )}
        </div>
      )}
      <span
        className="city3d-metrics"
        hidden={!debug}
        aria-hidden="true"
      >
        {fps}
      </span>
      <label className="city3d-directory">
        Find an address
        <select
          data-shortcut="j" aria-keyshortcuts="J" aria-label="Find an address"
          value={props.selected}
          onChange={e => {
            props.onSelect(e.target.value);
            focus.current(e.target.value);
          }}
        >
          {props.state.locations.map(l => (
            <option key={l.id} value={l.id}>
              {l.name}
            </option>
          ))}
        </select>
      </label>
    </section>
  );
}
