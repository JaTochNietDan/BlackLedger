import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import {interiorSettings,hasInterior} from '../.runtime/frontend-test/interiorSettings.js';

test('registered rooms reference exported assets and reject unknown destinations',()=>{
 const manifest=JSON.parse(readFileSync(new URL('../public/art/models/manifest.json',import.meta.url)));
 for(const [place,room] of Object.entries(interiorSettings)){
  assert.ok(hasInterior(place));assert.ok(manifest[room.model],`${place} has no exported room`);
 }
 for(const name of ['missing-room-id','constructor','toString','__proto__',''])assert.equal(hasInterior(name),false);
});
