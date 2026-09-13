import {wardrobe, dressPedestrian} from './city3dWardrobe';
import {headlightAlpha, headlightCentre} from './city3dHeadlights';
import {cityWeather, rainVertices} from './city3dWeather';
import {blockingBuildings} from './city3dOcclusion';
import {buildingCondition, buildingCutaway} from './city3dDamage';
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
import {CityCueQueue, availableSceneSlot, casualtyFall, gunfightPose, casualtySceneStart, GunfireAudio, BlastAudio} from './city3dEvents';
import type {SceneSlot} from './city3dEvents';
import {StreetTraffic, trafficSize, trafficModel, advanceWheel, wheelSteering, advanceSteering, frontWheelSteering} from './city3dTraffic';
import {pedestrianModel, isPedestrian} from './city3dCast';
import {playCityGunshot, playMoment, soundOn} from './sound';
import {cameraCommand, screenPan} from './city3dControls';
import {blastParticle, blastLight, blastOpacity, billowAlpha, debrisPose, fragmentBlocked} from './city3dBlast';

type Props = {
  state: Snapshot;
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
  arrived?: boolean;
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
  gunArm?: THREE.Object3D;
  muzzle?: THREE.Object3D;
  audio?: GunfireAudio | BlastAudio;
};
const modelNames = [
  'tenement',
  'mariner',
  'tavern',
  'casino',
  'warehouse',
  'civic',
  'shop',
  'villa',
  'ford',
  'hudson',
  'packard',
  'person',
  'woman',
  'revolver',
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
  latest.current = props;
  const controlsRef = useRef<OrbitControls | null>(null);
  const focus = useRef<(id?: string) => void>(() => {});
  const [expanded, setExpanded] = useState(false);
  const [following, setFollowing] = useState(false);
  const followPlayer = useRef(false);
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
      '3D Bellwether city. Drag to rotate, right drag to pan, scroll to zoom. Keyboard: arrows pan, Q and E rotate, plus and minus zoom, Home resets, Escape leaves the expanded city.',
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
    const reset = () => {
      controls.target.copy(home);
      camera.position.copy(home).add(new THREE.Vector3(180, 200, -240));
      camera.zoom = 1;
      camera.updateProjectionMatrix();
      controls.update();
    };
    focus.current = id => {
      setFollow(false);
      const lot = id ? lots.get(id) : undefined;
      if (!lot) {
        reset();
        return;
      }
      const offset = camera.position.clone().sub(controls.target);
      controls.target.set(lot.x, 0, lot.z);
      camera.position.copy(controls.target).add(offset);
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
    const landings: THREE.Group[] = [];
    const labels = new Map<string, THREE.Sprite>();
    let blockers = new Set<string>(), lastSightCheck = -Infinity;
    const cutawayAmounts = new Map<string, number>();
    const cutawayWindow = new THREE.Vector3();
    const actors = new Map<string, Actor>();
    const effects: Effect[] = [];
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
      if (e.button === 2 || e.pointerType === 'touch') setFollow(false);
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
    const keys = (e: KeyboardEvent) => {
      const command = cameraCommand(e);
      if (!command) return;
      e.preventDefault();
      if (command === 'reset' || command.startsWith('pan-')) setFollow(false);
      if (command === 'zoom-in' || command === 'zoom-out') {
        camera.zoom = THREE.MathUtils.clamp(
          camera.zoom * (command === 'zoom-out' ? 0.9 : 1.1),
          controls.minZoom,
          controls.maxZoom,
        );
        camera.updateProjectionMatrix();
      }
      if (command === 'rotate-left' || command === 'rotate-right') {
        const offset = camera.position.clone().sub(controls.target);
        offset.applyAxisAngle(
          new THREE.Vector3(0, 1, 0),
          command === 'rotate-left' ? 0.12 : -0.12,
        );
        camera.position.copy(controls.target).add(offset);
        controls.update();
      }
      if (command === 'reset') {
        reset();
      }
      if (command.startsWith('pan-')) {
        const pan = screenPan(camera.position, controls.target, command, 8.25 / camera.zoom);
        const move = new THREE.Vector3(pan.x, 0, pan.z);
        controls.target.add(move);
        camera.position.add(move);
        controls.update();
      }
    };
    canvas.addEventListener('keydown', keys);
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
      const actor = {
        object,
        model,
        points: [{x: 0, z: 0}],
        start: 0,
        end: 0,
        since: 0,
        duration: 0,
        walking: false,
        phase: 0,
        realSince: 0,
        wardrobe: isPedestrian(model) ? dressPedestrian(object, model, personWardrobe(id)) : [],
        wardrobeKey: isPedestrian(model) ? JSON.stringify(personWardrobe(id)) : '',
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
      if (a && (a.model !== model || (isPedestrian(model) && a.wardrobeKey !== JSON.stringify(personWardrobe(id))))) {
        releaseActor(a);
        actors.delete(id);
        a = undefined;
      }
      if (!a) a = addActor(id, model);
      // Arrivals hide their outdoor actor; a later journey must show it again.
      a.object.visible = true;
      a.arrived = false;
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
          const instances = new THREE.InstancedMesh(part.geometry, part.material, plan.lots.length);
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
          const label = makeLabel(latest.current.state.locations.find(p => p.id === lot.id)!.name);
          const box = new THREE.Box3().setFromObject(model, true);
          model.userData.sightBounds = box.clone();
          model.userData.blastOrigin = {x: lot.x, y: 0.25, z: box.min.z - 0.15};
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
    const reduced = matchMedia('(prefers-reduced-motion: reduce)');
    let sampleStart = performance.now(),
      samples: number[] = [],
      last = sampleStart;
    const tick = (now: number) => {
      if (dead) return;
      frame = requestAnimationFrame(tick);
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
      const playbackStarted = !!p.activeCue && lastActive !== p.activeCue.id;
      if ((lastActive && !p.activeCue) || worldID !== `${w.id}:${w.life}`) {
        for (const effect of effects) {
          scene.remove(effect.mesh, effect.light);
          if (effect.extra) scene.remove(effect.extra);
          effect.wardrobe?.forEach(material => material.dispose());
          effect.audio?.dispose();
          disposeDebris(effect);
          effect.mesh.dispose();
          (effect.mesh.material as THREE.Material).dispose();
        }
        effects.length = 0;
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
        for (const j of w.street || []) {
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
        for (const person of w.everyone || []) {
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
        for (const cue of p.journey
          ? []
          : cueQueue.take(
              `${w.id}:${w.life}`,
              w.last_result?.cues || [],
              p.activeCue,
              first,
              !first && revision === w.revision && playbackStarted,
            )) {
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
          let extra: THREE.Group | undefined, gunArm: THREE.Object3D | undefined, muzzle: THREE.Object3D | undefined;
          if (['killing', 'gunfight', 'raid', 'arrest'].includes(cue.kind)) {
            const model = cue.kind === 'killing' ? personModel(cue.actors?.[0]?.id || '')
              : cue.kind === 'gunfight' ? 'person' : 'police';
            extra = models.get(model)!.clone(true);
            if (isPedestrian(model)) costume = dressPedestrian(extra, model, personWardrobe(cue.kind === 'killing' ? cue.actors?.[0]?.id || '' : 'anonymous-shooter'));
            if (model === 'police') addVehicleShadow(extra, model);
            if (cue.kind === 'gunfight') {
              extra.rotation.y = Math.PI / 2;
              gunArm = extra.getObjectByName('arm1');
              const weapon = models.get('revolver')!.clone(true);
              weapon.position.set(0, -0.58, 0);
              weapon.rotation.x = Math.PI / 2;
              gunArm?.add(weapon);
              muzzle = weapon.getObjectByName('muzzle');
            }
            extra.visible = false;
            mesh.visible = false;
            scene.add(extra);
          }
          let debris: THREE.InstancedMesh | undefined;
          if (cue.kind === 'explosion') {
            const model = models.get('blast-fragment')!;
            model.updateMatrixWorld(true);
            model.traverse(part => {
              if (!(part instanceof THREE.Mesh) || debris) return;
              const material = (part.material as THREE.MeshStandardMaterial).clone();
              material.transparent = true;
              debris = new THREE.InstancedMesh(part.geometry.clone().applyMatrix4(part.matrixWorld), material, 12);
              debris.frustumCulled = false;
              scene.add(debris);
            });
          }
          effects.push({cue, since: now, mesh, light, debris, extra, wardrobe: costume, gunArm, muzzle,
            audio: cue.kind === 'gunfight' ? new GunfireAudio(playCityGunshot)
              : cue.kind === 'explosion' ? new BlastAudio(() => playMoment('explosion')) : undefined});
          if (p.activeCue?.id === cue.id) focus.current(cue.target);
        }
        // Condition-driven surface stains imply no ongoing fire or invented collapse.
        for (const place of w.locations) {
          const b = buildings.get(place.id);
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
        const hour = (w.minute % 1440) / 60;
        const night = hour < 6 || hour >= 20;
        const weather = cityWeather(w.sky, night);
        sky.intensity = weather.ambient;
        groundMat.color.setHex(0x646460).multiplyScalar(weather.roadTone);
        pavementMat.color.setHex(0xaaa18b).multiplyScalar(weather.pavementTone);
        groundMat.roughness = weather.roadRoughness;
        pavementMat.roughness = weather.pavementRoughness;
        pools.visible = night;
        for (const b of buildings.values())
          b.traverse(o => {
            if (o.name === 'clock-hand-minute') o.rotation.z = ((w.minute % 60) * Math.PI) / 30;
            if (o.name === 'clock-hand-hour') o.rotation.z = ((w.minute % 720) * Math.PI) / 360;
            if (o instanceof THREE.Mesh)
              for (const m of Array.isArray(o.material) ? o.material : [o.material])
                if (m instanceof THREE.MeshStandardMaterial && m.emissive.getHex() !== 0)
                  m.emissiveIntensity = night ? 1.2 : 0.2;
          });
        sun.intensity = weather.sun;
        scene.background = new THREE.Color(weather.background);
        scene.fog = new THREE.Fog(scene.background, 260, weather.fogFar);
        renderer.shadowMap.needsUpdate = true;
        revision = w.revision;
        worldID = `${w.id}:${w.life}`;
        previous = w;
      }
      if (ready) {
        const hereForCar = lots.get(w.player.location);
        const parked = actors.get('player-car');
        if (
          !p.journey &&
          hereForCar &&
          w.player.alive &&
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
        if (key !== journeyKey || motion !== motionWas) {
          journeyKey = key;
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
                2400,
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
        if (!motion) traffic.clear();
        // Reserve the shooter before associated casualties, so a full batch
        // cannot occupy every slot while waiting for an unstaged first shot.
        const stagingOrder = [...effects].sort((a, b) =>
          Number(b.cue.kind === 'gunfight') - Number(a.cue.kind === 'gunfight'));
        for (const e of stagingOrder) {
          if (!e.extra || e.slot) continue;
          const occupied = [
            ...[...actors.values()]
              .filter(a => a.object.visible)
              .map(a => ({
                model: trafficModel(a.model,a.start === a.end),
                pose: {x: a.object.position.x, z: a.object.position.z, heading: a.object.rotation.y},
              })),
            ...effects.flatMap(other => other.slot ? [other.slot] : []),
          ];
          e.slot = availableSceneSlot(lots.get(e.cue.target)!, e.cue.kind, occupied);
          if (e.slot) {
            e.extra.position.set(e.slot.root.x, e.slot.model === 'parked-police' ? vehicleRootHeight(e.slot.root) : 0.2, e.slot.root.z);
            e.light.position.set(e.slot.root.x, 3, e.slot.root.z);
            e.since = now;
          }
        }
        const placements = traffic.update(
          [...actors]
            .filter(([, a]) => !a.arrived)
            .map(([id, a]) => {
              const t =
                motion && a.duration ? Math.min(1, (movementClock - a.since) / a.duration) : 1;
              return {
                id,
                model: trafficModel(a.model,a.start === a.end),
                points: a.points,
                progress: a.start + (a.end - a.start) * t,
              };
            })
            .concat(effects.flatMap(e => e.slot ? [{
              id: `scene:${e.cue.id}`,
              model: e.slot.model,
              points: [e.slot.pose],
              progress: 0,
            }] : [])),
          dt / 1000,
          playback.current,
        );
        for (const [id, a] of actors) {
          const placement = placements.get(id);
          a.object.visible = !!placement && !placement.waiting;
          if (!placement || placement.waiting) continue;
          const at = placement.pose;
          const distance = Math.hypot(a.object.position.x - at.x, a.object.position.z - at.z);
          const moved = distance > 0.0001;
          if (a.walking && moved)
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
          for (const limb of a.limbs) {
            const side = limb.name.endsWith('-1') ? 0 : Math.PI;
            const phase = a.phase + side;
            limb.rotation.x =
              !a.walking || !moved || !motion
                ? 0
                : limb.name.startsWith('knee')
                  ? Math.max(0, Math.sin(phase + 0.7)) * 0.65
                  : Math.sin(phase + (limb.name.startsWith('arm') ? Math.PI : 0)) *
                    (limb.name.startsWith('arm') ? 0.23 : 0.35);
          }
          if (a.walking && moved && motion)
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
        for (let i = effects.length - 1; i >= 0; i--) {
          const e = effects[i];
          if (e.extra && motion) {
            const staged = !!e.slot && !placements.get(`scene:${e.cue.id}`)?.waiting;
            e.extra.visible = e.mesh.visible = staged;
            if (!staged) {
              e.since += dt;
              continue;
            }
          }
          if (e.cue.kind === 'killing') {
            const gunScene = effects.find(other => other.cue.kind === 'gunfight' && other.cue.target === e.cue.target);
            e.since = casualtySceneStart(e.since, now, gunScene?.since);
          }
          const t = (now - e.since) / 3000,
            lot = lots.get(e.cue.target)!;
          if (t >= 1 || !motion) {
            scene.remove(e.mesh, e.light);
            if (e.extra) scene.remove(e.extra);
            e.wardrobe?.forEach(material => material.dispose());
            e.audio?.dispose();
            disposeDebris(e);
            e.mesh.dispose();
            (e.mesh.material as THREE.Material).dispose();
            effects.splice(i, 1);
            continue;
          }
          const at = e.slot?.root || entrance(lot);
          const blast = e.cue.kind === 'explosion';
          const blastOrigin = buildings.get(lot.id)?.userData.blastOrigin;
          if (blast && blastOrigin) e.light.position.set(blastOrigin.x, 2, blastOrigin.z);
          const shot = e.cue.kind === 'gunfight';
          const firing = gunfightPose(t * 3);
          e.audio?.update(t * 3, soundOn());
          const muzzlePosition = new THREE.Vector3(at.x, 1.4, at.z);
          if (shot && e.gunArm && e.muzzle) {
            e.gunArm.rotation.x = firing.arm;
            e.extra!.updateMatrixWorld(true);
            e.muzzle.getWorldPosition(muzzlePosition);
            e.light.position.copy(muzzlePosition);
          }
          const casualty = e.cue.kind === 'killing';
          if (casualty && e.extra) {
            const fall = casualtyFall(t);
            e.extra.rotation.z = fall.rotation;
            e.extra.position.y = fall.height;
          }
          const police = ['raid', 'arrest'].includes(e.cue.kind);
          for (let j = 0; j < 32; j++) {
            const a = j * 2.399;
            const r = shot ? 0.35 : 1.8;
            tmp.position.set(at.x + Math.cos(a) * r, 1 + Math.sin(a) * r * 0.4, at.z + Math.sin(a) * r);
            tmp.scale.setScalar(casualty || shot ? 0.001 : 0.25);
            const burst = blast ? blastParticle(j, t * 3) : null;
            if (burst && blastOrigin) {
              tmp.position.set(blastOrigin.x + burst.x, blastOrigin.y + burst.y, blastOrigin.z + burst.z);
              tmp.scale.setScalar(Math.max(.001, burst.size));
            }
            if (shot) {
              tmp.position.copy(muzzlePosition);
              const smoke = j === 1 && firing.smoke > 0;
              if (smoke) tmp.position.y += (1 - firing.smoke) * 0.3;
              tmp.scale.setScalar(j === 0 && firing.flash ? 0.28 : smoke ? 0.15 + (1 - firing.smoke) * 0.3 : 0.001);
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
                burst ? burst.color : (shot && j === 1)
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
          if (e.debris && blastOrigin) {
            for (let j = 0; j < 12; j++) {
              const fragment = debrisPose(j, t * 3);
              const x = blastOrigin.x + fragment.x, z = blastOrigin.z + fragment.z;
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
            (e.debris.material as THREE.MeshStandardMaterial).opacity = blastOpacity(t*3)/.8;
          }
          e.mesh.instanceMatrix.needsUpdate = true;
          if (e.mesh.instanceColor) e.mesh.instanceColor.needsUpdate = true;
          (e.mesh.material as THREE.MeshBasicMaterial).opacity = blast ? blastOpacity(t * 3) : 1 - t;
          e.light.intensity = blast
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
      if (ready) lastActive = p.activeCue?.id || null;
      motionWas = motion;
      controls.enableDamping = motion;
      controls.update();
      const followed = actors.get('player');
      if (followPlayer.current && ready && followed?.object.visible) {
        const dx = followed.object.position.x - controls.target.x;
        const dz = followed.object.position.z - controls.target.z;
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
          .map(e => e.extra!.position.clone().add(new THREE.Vector3(0, .9, 0))) : [];
      if (sightPoints.length) {
        if (now - lastSightCheck >= 100) {
          blockers = new Set(sightPoints.flatMap(sight => [...blockingBuildings(camera, sight, buildings)]));
          lastSightCheck = now;
        }
        const size = renderer.getDrawingBufferSize(new THREE.Vector2());
        const projected = sightPoints.map(sight => {
          const screen = sight.project(camera);
          return new THREE.Vector2((screen.x + 1) * size.x / 2, (screen.y + 1) * size.y / 2);
        });
        const centre = projected.reduce((sum, point) => sum.add(point), new THREE.Vector2()).divideScalar(projected.length);
        const radius = Math.max(...projected.map(point => point.distanceTo(centre)))
          + (65 + camera.zoom * 2) * renderer.getPixelRatio();
        cutawayWindow.set(centre.x, centre.y, radius);
      } else { blockers.clear(); lastSightCheck = -Infinity; }
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
      const lampHour = (w.minute % 1440) / 60, lampsOn = lampHour < 6 || lampHour >= 20;
      for (const actor of actors.values()) {
        const lit = ready && lampsOn && actor.object.visible && actor.points.length > 1;
        for (const material of actor.lamps) material.emissiveIntensity = lit ? (material.name === 'headlamps' ? 1.6 : .8) : 0;
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
      renderer.render(scene, camera);
      if (ready)
        canvas.dataset.presentation = JSON.stringify({
          revision: w.revision,
          minute: w.minute,
          playbackRate: playback.current,
          headlightPools: headlightPools.count,
          harbour: {visible: !!harbourLot, waterClock},
          followingPlayer: followPlayer.current,
          cutawayBuildings: [...blockers],
          weather: {kind: w.sky?.kind || 'clear', wet: w.sky?.wet || 0, rainVisible: rainfall.visible, rainClock},
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
              wheelPhase: a.wheels.length ? a.wheelPhase : undefined,
              steering: a.wheels.length ? a.steering : undefined,
              lampsOn: a.lamps.length ? a.lamps.some(m => m.emissiveIntensity > 0) : undefined,
              frontWheels: a.wheels.length ? a.wheels.filter(w => w.name.includes('-front-')).map(w => ({x: w.position.x, angle: w.rotation.y})) : undefined,
            })),
          waiting: [...actors].filter(([, a]) => !a.arrived && !a.object.visible).map(([id]) => id),
          effects: effects.map(e => ({
            id: e.cue.id, kind: e.cue.kind, target: e.cue.target,
            staged: !e.extra || e.extra.visible, x: e.slot?.root.x, z: e.slot?.root.z,
            debris: e.debris?.count,
            blastOrigin: e.cue.kind === 'explosion' ? buildings.get(e.cue.target)?.userData.blastOrigin : undefined,
            arm: e.gunArm?.rotation.x,
            audioShots: e.cue.kind === 'gunfight' ? e.audio?.started : undefined,
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
        setWaiting([...actors.values()].filter(a => !a.arrived && !a.object.visible).length);
        samples = [];
        sampleStart = now;
      }
    };
    frame = requestAnimationFrame(tick);
    const lost = (e: Event) => {
      e.preventDefault();
      graphicsLost = true;
      effects.forEach(effect => effect.audio?.dispose());
      setFailure('The graphics context was interrupted. Reload the city to restore it.');
    };
    canvas.addEventListener('webglcontextlost', lost);
    return () => {
      dead = true;
      cancelAnimationFrame(frame);
      observer.disconnect();
      controls.dispose();
      controlsRef.current = null;
      canvas.removeEventListener('pointerdown', pointerDown);
      canvas.removeEventListener('pointermove', pointerMove);
      canvas.removeEventListener('pointerup', pointerUp);
      canvas.removeEventListener('keydown', keys);
      canvas.removeEventListener('webglcontextlost', lost);
      effects.forEach(effect => { effect.audio?.dispose(); disposeDebris(effect); });
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
      className={'city3d' + (expanded ? ' city3d-expanded' : '')}
      aria-label="Bellwether city"
      role={expanded ? 'dialog' : 'region'}
      aria-modal={expanded || undefined}
      onKeyDown={e => {
        if (!expanded || e.ctrlKey || e.metaKey || e.altKey || e.nativeEvent.isComposing) return;
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
          <span>Arrows: pan · Q/E: rotate · +/−: zoom · Home: reset · Esc: return</span>
        </div>
        <div className="city3d-tools">
          <button ref={expandButton} aria-pressed={expanded} onClick={() => {
            setExpanded(!expanded);
            if (!expanded) host.current?.querySelector('canvas')?.focus();
          }}>
            {expanded ? 'Return to game' : 'Expand city'}
          </button>
          <button onClick={() => focus.current()}>Whole city</button>
          <button aria-pressed={following} title={following ? "Stop following your character" : "Follow your character or car"}
            onClick={() => setFollow(!followPlayer.current)}>{following ? "Stop following" : "Find me"}</button>
          <button onClick={() => focus.current(props.selected)}>Focus address</button>
          <button
            onClick={() => setPlaybackRate(rate => (rate === 1 ? 4 : 1))}
            aria-label={`Travel playback speed: ${playbackRate} times`}
          >
            Travel {playbackRate}×
          </button>
        </div>
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
            <button onClick={props.onEnter}>Step inside →</button>
          ) : (
            <button
              disabled={!travel || travel.disabled || props.busy || !!props.journey}
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
        hidden={!new URLSearchParams(location.search).has('city-debug')}
        aria-hidden="true"
      >
        {fps}
      </span>
      <label className="city3d-directory">
        Find an address
        <select
          aria-label="Find an address"
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
