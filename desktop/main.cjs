const {app, BrowserWindow, dialog, Menu} = require('electron');
const path = require('node:path');
const fs = require('node:fs');
const {startBackend} = require('./backend.cjs');

const smoke = process.argv.includes('--smoke-test') && process.env.BLACK_LEDGER_SMOKE_DIR;
const smokeDir = smoke && path.resolve(process.env.BLACK_LEDGER_SMOKE_DIR);
if (smokeDir) {
  fs.mkdirSync(smokeDir, {recursive: true});
  app.setPath('userData', path.join(smokeDir, 'desktop-profile'));
}
let backend, window, quitting = false, cleaned = false;
let exitCode = 0;
function failure(error) {
  exitCode = 1;
  if (smokeDir) fs.writeFileSync(path.join(smokeDir, 'error.txt'), String(error.stack || error));
  else dialog.showErrorBox('Black Ledger could not continue', error.message + '\nYour saved campaign is kept separately from the app.');
  app.quit();
}
if (!app.requestSingleInstanceLock()) {
  app.quit();
} else {
  app.on('second-instance', () => {
    if (window) { if (window.isMinimized()) window.restore(); window.show(); window.focus(); }
  });
  app.on('window-all-closed', () => app.quit());
  app.on('before-quit', event => {
    quitting = true;
    if (!cleaned && backend) {
      event.preventDefault();
      backend.stop().finally(() => { cleaned = true; app.exit(exitCode); });
    }
  });
  app.whenReady().then(async () => {
    const game = path.join(process.resourcesPath, 'game');
    const save = smokeDir ? path.join(smokeDir, 'campaign.sqlite3') :
      process.env.BLACK_LEDGER_DB || path.join(app.getPath('appData'), 'BlackLedger', 'campaign.sqlite3');
    const executable = path.join(game, process.platform === 'win32' ? 'blackledger.exe' : 'blackledger');
    const env = {...process.env};
    delete env.BLACK_LEDGER_WEB;
    backend = startBackend(executable,
      ['-desktop', '-browser=false', '-addr', ':0', '-db', save, '-parent-stdio'], {env, cwd: game});
    backend.closed.then(() => { if (!quitting) failure(new Error('The game server stopped unexpectedly.\n' + backend.output)); });
    const origin = await backend.ready();
    if (quitting) return;
    window = new BrowserWindow({
      title: 'Black Ledger', width: 1440, height: 960, minWidth: 900, minHeight: 640,
      show: false, backgroundColor: '#17140f',
      webPreferences: {nodeIntegration: false, contextIsolation: true, sandbox: true},
    });
    window.webContents.setWindowOpenHandler(() => ({action: 'deny'}));
    window.webContents.on('will-navigate', (event, url) => {
      if (new URL(url).origin !== origin) event.preventDefault();
    });
    window.webContents.on('will-attach-webview', event => event.preventDefault());
    window.webContents.session.setPermissionRequestHandler((_web, _permission, callback) => callback(false));
    window.webContents.session.setPermissionCheckHandler(() => false);
    window.webContents.on('render-process-gone', (_event, details) => {
      if (!quitting) failure(new Error('The game window stopped: ' + details.reason));
    });
    Menu.setApplicationMenu(Menu.buildFromTemplate([
      ...(process.platform === 'darwin' ? [{label: 'Black Ledger', submenu: [{role: 'about'}, {type: 'separator'}, {role: 'quit'}]}] : []),
      {label: 'Game', submenu: [{role: 'reload'}, {role: 'togglefullscreen'}, {type: 'separator'}, {role: 'quit'}]},
      {role: 'editMenu'},
    ]));
    await window.loadURL(origin);
    if (smokeDir) {
      // Test the packaged window and server, then exercise the normal close path.
      const response = await fetch(origin + '/api/state');
      if (!response.ok) throw new Error('The packaged game did not return a campaign');
      const state = await response.json();
      fs.writeFileSync(path.join(smokeDir, 'ready.json'), JSON.stringify({origin, backendPID: backend.child.pid,
        url: window.webContents.getURL(), revision: state.revision, player: state.player}));
      window.close();
    } else window.show();
  }).catch(failure);
}
