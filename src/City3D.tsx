import {useEffect, useRef, useState} from 'react';
import * as THREE from 'three';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import type {Snapshot, VisualCue} from './types';
import {
  cityPlan,
  entrance,
  route,
  PITCH,
  STREET_WIDTH,
  streetsidePosition,
  parkingSpot,
  lampPositions,
} from './city3dPlan';
import type {Lot, Point} from './city3dPlan';
import type {Journey} from './TravelPresentation';
import './city3d.css';
import {CityCueQueue, availableSceneSlot, casualtyFall, gunfightPose, casualtySceneStart} from './city3dEvents';
import type {SceneSlot} from './city3dEvents';
import {StreetTraffic} from './city3dTraffic';
import {pedestrianModel, isPedestrian} from './city3dCast';

type Props = {
  state: Snapshot;
  selected: string;
  onSelect: (id: string) => void;
  onTravel: (id: string) => void;
  onEnter: () => void;
  motion: boolean;
  journey: Journey | null;
  busy: boolean;
  activeCue: VisualCue | null;
  onSkipCue: () => void;
  onSkipJourney: () => void;
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
  arrived?: boolean;
};
type Effect = {
  cue: VisualCue;
  since: number;
  mesh: THREE.InstancedMesh;
  light: THREE.PointLight;
  extra?: THREE.Group;
  slot?: SceneSlot;
  gunArm?: THREE.Object3D;
  muzzle?: THREE.Object3D;
};
const modelNames = [
  'tenement',
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
  'filling',
  'garage',
  'dealer',
  'chapel',
  'docks',
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
  const latest = useRef(props);
  latest.current = props;
  const controlsRef = useRef<OrbitControls | null>(null);
  const focus = useRef<(id?: string) => void>(() => {});
  const [expanded, setExpanded] = useState(false);
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
      '3D Bellwether city. Drag to rotate, right drag to pan, scroll to zoom.',
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
    controls.maxZoom = 12;
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
    const textures: THREE.Texture[] = [];
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
    const selection = new THREE.Mesh(
      new THREE.RingGeometry(9.5, 9.8, 64),
      new THREE.MeshBasicMaterial({color: 0xf3ce83, side: THREE.DoubleSide, depthWrite: false}),
    );
    selection.rotation.x = -Math.PI / 2;
    selection.position.y = 0.22;
    scene.add(selection);
    const playerRing = new THREE.Mesh(
      new THREE.RingGeometry(1, 1.25, 32),
      new THREE.MeshBasicMaterial({color: 0xffd885, side: THREE.DoubleSide}),
    );
    playerRing.rotation.x = -Math.PI / 2;
    scene.add(playerRing);
    const models = new Map<string, THREE.Group>();
    const buildings = new Map<string, THREE.Group>();
    const labels = new Map<string, THREE.Sprite>();
    const actors = new Map<string, Actor>();
    const effects: Effect[] = [];
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
      const hits = raycaster.intersectObjects([...buildings.values()], true);
      if (hits.length) {
        let ob: THREE.Object3D | null = hits[0].object;
        while (ob && !ob.userData.place) ob = ob.parent;
        if (ob) latest.current.onSelect(ob.userData.place);
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
      const hits = raycaster.intersectObjects([...buildings.values()], true);
      let ob: THREE.Object3D | null = hits[0]?.object || null;
      while (ob && !ob.userData.place) ob = ob.parent;
      hovered = ob?.userData.place || '';
      canvas.style.cursor = hovered ? 'pointer' : 'grab';
    };
    canvas.addEventListener('pointermove', pointerMove);
    canvas.addEventListener('pointerdown', pointerDown);
    canvas.addEventListener('pointerup', pointerUp);
    const keys = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setExpanded(false);
      if (['+', '=', '-'].includes(e.key)) {
        e.preventDefault();
        camera.zoom = THREE.MathUtils.clamp(
          camera.zoom * (e.key === '-' ? 0.9 : 1.1),
          controls.minZoom,
          controls.maxZoom,
        );
        camera.updateProjectionMatrix();
      }
      if (['q', 'e'].includes(e.key.toLowerCase())) {
        e.preventDefault();
        const offset = camera.position.clone().sub(controls.target);
        offset.applyAxisAngle(
          new THREE.Vector3(0, 1, 0),
          e.key.toLowerCase() === 'q' ? 0.12 : -0.12,
        );
        camera.position.copy(controls.target).add(offset);
        controls.update();
      }
      if (e.key === 'Home') {
        e.preventDefault();
        reset();
      }
      if (['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(e.key)) {
        e.preventDefault();
        const move = new THREE.Vector3(
          e.key === 'ArrowLeft' ? -5 : e.key === 'ArrowRight' ? 5 : 0,
          0,
          e.key === 'ArrowUp' ? -5 : e.key === 'ArrowDown' ? 5 : 0,
        );
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
    let movementClock = performance.now();
    const addActor = (id: string, model: string): Actor => {
      const object = models.get(model)!.clone(true);
      scene.add(object);
      const limbs: THREE.Object3D[] = [];
      object.traverse(o => {
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
        limbs,
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
      if (a && a.model !== model) {
        scene.remove(a.object);
        actors.delete(id);
        a = undefined;
      }
      if (!a) a = addActor(id, model);
      // Arrivals hide their outdoor actor; a later journey must show it again.
      a.object.visible = true;
      a.arrived = false;
      a.points = points;
      a.start = start;
      a.end = end;
      a.since = movementClock;
      a.realSince = now;
      a.duration = duration;
      a.walking = isPedestrian(model) && start !== end;
    };
    const loader = new GLTFLoader();
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
          disposeTree(gltf.scene);
          return;
        }
        models.set(name, gltf.scene);
      }),
    )
      .then(() => {
        if (dead) return;
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
          const box = new THREE.Box3().setFromObject(model);
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
        motion = p.motion && !reduced.matches;
      const playbackStarted = !!p.activeCue && lastActive !== p.activeCue.id;
      if (lastActive && !p.activeCue) {
        for (const effect of effects) {
          scene.remove(effect.mesh, effect.light);
          if (effect.extra) scene.remove(effect.extra);
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
          for (const a of actors.values()) scene.remove(a.object);
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
            scene.remove(a.object);
            actors.delete(id);
          }
        const here = lots.get(w.player.location);
        if (here && w.player.alive) {
          if (!actors.has('player')) assign('player', personModel('player'), [entrance(here)], 0, 0, 0, now);
        } else {
          const a = actors.get('player');
          if (a) scene.remove(a.object);
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
              map: particleTexture,
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
          let extra: THREE.Group | undefined, gunArm: THREE.Object3D | undefined, muzzle: THREE.Object3D | undefined;
          if (['killing', 'gunfight', 'raid', 'arrest'].includes(cue.kind)) {
            const model = cue.kind === 'killing' ? personModel(cue.actors?.[0]?.id || '')
              : cue.kind === 'gunfight' ? 'person' : 'police';
            extra = models.get(model)!.clone(true);
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
          effects.push({cue, since: now, mesh, light, extra, gunArm, muzzle});
          if (p.activeCue?.id === cue.id) focus.current(cue.target);
        }
        // Damage is a persistent scorch state, not evidence of a continuing fire.
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
                  for (const m of Array.isArray(o.material) ? o.material : [o.material]) {
                    if (m instanceof THREE.MeshStandardMaterial)
                      m.userData.baseColor = m.color.clone();
                  }
                }
                for (const m of Array.isArray(o.material) ? o.material : [o.material])
                  if (m instanceof THREE.MeshStandardMaterial)
                    m.color
                      .copy(m.userData.baseColor as THREE.Color)
                      .multiplyScalar(
                        0.42 + (0.58 * Math.max(0, Math.min(100, place.condition))) / 100,
                      );
              }
            });
        }
        const hour = (w.minute % 1440) / 60;
        const night = hour < 6 || hour >= 20;
        sky.intensity = night ? 0.9 : 2.1;
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
        sun.intensity = night ? 0.35 : 3.2;
        scene.background = new THREE.Color(night ? 0x17232c : 0x657477);
        scene.fog = new THREE.Fog(scene.background, 260, w.sky?.kind === 'fog' ? 440 : 850);
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
          scene.remove(parked.object);
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
        if (!p.journey && actors.has('player') && actors.get('player')!.model !== personModel('player')) {
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
                model: a.model,
                pose: {x: a.object.position.x, z: a.object.position.z, heading: a.object.rotation.y},
              })),
            ...effects.flatMap(other => other.slot ? [other.slot] : []),
          ];
          e.slot = availableSceneSlot(lots.get(e.cue.target)!, e.cue.kind, occupied);
          if (e.slot) {
            e.extra.position.set(e.slot.root.x, 0.2, e.slot.root.z);
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
                model: a.model,
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
          a.object.position.set(at.x, 0.2, at.z);
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
          if (id === 'player') playerRing.position.set(at.x, 0.23, at.z);
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
            e.mesh.dispose();
            (e.mesh.material as THREE.Material).dispose();
            effects.splice(i, 1);
            continue;
          }
          const at = e.slot?.root || entrance(lot);
          const blast = e.cue.kind === 'explosion';
          const shot = e.cue.kind === 'gunfight';
          const firing = gunfightPose(t * 3);
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
            const smoke = blast && t > 0.12 + (j % 8) * 0.055;
            const r = blast
              ? Math.sin((Math.min(1, t * 2) * Math.PI) / 2) * (2 + (j % 5))
              : shot
                ? 0.35
                : 1.8;
            tmp.position.set(
              at.x + Math.cos(a) * r,
              1 + (smoke ? t * 9 : Math.sin(a) * r * 0.4),
              at.z + Math.sin(a) * r,
            );
            tmp.scale.setScalar(
              blast
                ? smoke
                  ? 3.5 * t + 1.3
                  : (1 - t) * 4.1
                : shot
                  ? j < 3 && Math.floor(t * 22) % 3 === 0
                    ? 0.45
                    : 0.001
                  : casualty
                    ? 0.001
                    : 0.25,
            );
            if (shot) {
              tmp.position.copy(muzzlePosition);
              const smoke = j === 1 && firing.smoke > 0;
              if (smoke) tmp.position.y += (1 - firing.smoke) * 0.3;
              tmp.scale.setScalar(j === 0 && firing.flash ? 0.28 : smoke ? 0.15 + (1 - firing.smoke) * 0.3 : 0.001);
            }
            if (police) {
              // A period rotating red roof beacon, rather than sparks around the car.
              tmp.position.set(at.x, 1.96, at.z);
              tmp.scale.setScalar(j === 0 ? 0.45 + 0.35 * Math.max(0, Math.sin(t * 38)) : 0.001);
            }
            tmp.quaternion.copy(camera.quaternion);
            tmp.updateMatrix();
            e.mesh.setMatrixAt(j, tmp.matrix);
            e.mesh.setColorAt(
              j,
              new THREE.Color(
                smoke || (shot && j === 1)
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
          e.mesh.instanceMatrix.needsUpdate = true;
          if (e.mesh.instanceColor) e.mesh.instanceColor.needsUpdate = true;
          (e.mesh.material as THREE.MeshBasicMaterial).opacity = 1 - t;
          e.light.intensity = blast
            ? 100 * (1 - t) ** 3
            : shot && firing.flash
              ? 30
              : police
                ? 8 * Math.max(0, Math.sin(t * 38))
                : 0;
          if (police) {
            e.light.color.setHex(0xe53220);
            e.light.position.y = 1.96;
          }
        }
      }
      if (ready) lastActive = p.activeCue?.id || null;
      motionWas = motion;
      controls.enableDamping = motion;
      controls.update();
      // Keep pan within the city plus its waterfront margin.
      const before = controls.target.clone();
      controls.target.x = THREE.MathUtils.clamp(controls.target.x, -15, plan.width + 15);
      controls.target.z = THREE.MathUtils.clamp(controls.target.z, -15, plan.depth + 15);
      camera.position.add(controls.target.clone().sub(before));
      renderer.render(scene, camera);
      if (ready)
        canvas.dataset.presentation = JSON.stringify({
          revision: w.revision,
          minute: w.minute,
          playbackRate: playback.current,
          actors: [...actors]
            .filter(([, a]) => a.object.visible)
            .map(([id, a]) => ({
              id,
              model: a.model,
              x: a.object.position.x,
              z: a.object.position.z,
            })),
          waiting: [...actors].filter(([, a]) => !a.arrived && !a.object.visible).map(([id]) => id),
          effects: effects.map(e => ({
            id: e.cue.id, kind: e.cue.kind, target: e.cue.target,
            staged: !e.extra || e.extra.visible, x: e.slot?.root.x, z: e.slot?.root.z,
            arm: e.gunArm?.rotation.x,
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
      disposeTree(scene);
      for (const m of models.values()) disposeTree(m);
      textures.forEach(t => t.dispose());
      effectGeometry.dispose();
      renderer.dispose();
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
        if (!expanded) return;
        if (e.key === 'Escape') {
          setExpanded(false);
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
          } else if (!e.shiftKey && at === items.length - 1) {
            e.preventDefault();
            items[0]?.focus();
          }
        }
      }}
    >
      <div className="city3d-canvas" ref={host} />
      <div className="city3d-heading">
        <small>BELLWETHER · 1950</small>
        <span>Drag to rotate · Right drag to pan · Scroll to zoom</span>
      </div>
      <div className="city3d-tools">
        <button aria-pressed={expanded} onClick={() => setExpanded(v => !v)}>
          {expanded ? 'Return to game' : 'Expand city'}
        </button>
        <button onClick={() => focus.current()}>Whole city</button>
        <button onClick={() => focus.current(props.state.player.location)}>Find me</button>
        <button onClick={() => focus.current(props.selected)}>Focus address</button>
        <button
          onClick={() => setPlaybackRate(rate => (rate === 1 ? 4 : 1))}
          aria-label={`Travel playback speed: ${playbackRate} times`}
        >
          Travel {playbackRate}×
        </button>
      </div>
      {(status || failure) && (
        <p className="city3d-status" role="status">
          {failure || status}
        </p>
      )}
      {expanded && props.activeCue && (
        <div className="city3d-event" role="status">
          <strong>{props.activeCue.caption}</strong>
          <button onClick={props.onSkipCue}>Skip scene →</button>
        </div>
      )}
      {waiting > 0 && (
        <p className="city3d-status" role="status">
          {waiting} travellers waiting for space at their departures.
        </p>
      )}
      {expanded && props.journey && (
        <div className="city3d-event" role="status">
          <strong>Travelling to {props.journey.to.name}</strong>
          <button onClick={props.onSkipJourney}>Skip journey →</button>
        </div>
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
function disposeTree(root: THREE.Object3D) {
  const geometry = new Set<THREE.BufferGeometry>(),
    materials = new Set<THREE.Material>();
  root.traverse(o => {
    if (o instanceof THREE.Mesh) {
      geometry.add(o.geometry);
      for (const m of Array.isArray(o.material) ? o.material : [o.material]) materials.add(m);
    } else if (o instanceof THREE.Sprite) materials.add(o.material);
  });
  geometry.forEach(g => g.dispose());
  materials.forEach(m => {
    for (const v of Object.values(m)) if (v instanceof THREE.Texture) v.dispose();
    m.dispose();
  });
}
