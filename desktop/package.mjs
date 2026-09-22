import {packager} from '@electron/packager';
import fs from 'node:fs/promises';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const here = path.dirname(fileURLToPath(import.meta.url));
const [game, destination, platform, arch] = process.argv.slice(2);
if(platform!==process.platform||arch!==process.arch)throw new Error('Package on the target OS and architecture so native AI dependencies match.');
const meta = JSON.parse(await fs.readFile(path.join(here, 'package.json'), 'utf8'));
const source = path.join(path.dirname(game), 'desktop-source');
await fs.mkdir(source);
await fs.writeFile(path.join(source, 'package.json'), JSON.stringify({name: meta.name, version: meta.version,
  description: meta.description, license: meta.license, main: 'main.cjs', dependencies:meta.dependencies}));
for (const name of ['main.cjs', 'backend.cjs']) await fs.copyFile(path.join(here, name), path.join(source, name));
await fs.cp(path.join(here,'ai'),path.join(source,'ai'),{recursive:true});
const lock=JSON.parse(await fs.readFile(path.join(here,'package-lock.json'),'utf8'));
for(const [name,entry] of Object.entries(lock.packages)) {
  if(!name||entry.dev||entry.devOptional)continue;
  const from=path.join(here,name),to=path.join(source,name);
  try {await fs.access(from);} catch {if(entry.optional)continue;throw new Error('Missing desktop dependency '+name);}
  await fs.cp(from,to,{recursive:true,verbatimSymlinks:true});
}
const [result] = await packager({dir: source, out: destination, name: 'Black Ledger', executableName: 'BlackLedger',
  platform, arch, electronVersion: meta.devDependencies.electron, asar: {unpackDir:'node_modules'}, prune: false,
  appBundleId: 'com.jatochnietdan.blackledger', appCategoryType: 'public.app-category.games',
  appCopyright: 'Copyright 2026 JaTochNietDan and Black Ledger contributors',
  extraResource: [game], quiet: true});
await fs.writeFile(path.join(destination, 'package-path.txt'), result);
