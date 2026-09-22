import {packager} from '@electron/packager';
import fs from 'node:fs/promises';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const here = path.dirname(fileURLToPath(import.meta.url));
const [game, destination, platform, arch] = process.argv.slice(2);
const meta = JSON.parse(await fs.readFile(path.join(here, 'package.json'), 'utf8'));
const source = path.join(path.dirname(game), 'desktop-source');
await fs.mkdir(source);
await fs.writeFile(path.join(source, 'package.json'), JSON.stringify({name: meta.name, version: meta.version,
  description: meta.description, license: meta.license, main: 'main.cjs'}));
for (const name of ['main.cjs', 'backend.cjs']) await fs.copyFile(path.join(here, name), path.join(source, name));
const [result] = await packager({dir: source, out: destination, name: 'Black Ledger', executableName: 'BlackLedger',
  platform, arch, electronVersion: meta.devDependencies.electron, asar: true, prune: false,
  appBundleId: 'com.jatochnietdan.blackledger', appCategoryType: 'public.app-category.games',
  appCopyright: 'Copyright 2026 JaTochNietDan and Black Ledger contributors',
  extraResource: [game], quiet: true});
await fs.writeFile(path.join(destination, 'package-path.txt'), result);
