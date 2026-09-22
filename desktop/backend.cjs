const {spawn} = require('node:child_process');
const {setTimeout: delay} = require('node:timers/promises');

// The pipe is a lifetime lease: Go exits when this process closes or crashes.
function startBackend(executable, args, options = {}) {
  const {originPattern = /http:\/\/127\.0\.0\.1:\d+/, ...spawnOptions} = options;
  const child = spawn(executable, args, {windowsHide: true, stdio: ['pipe', 'pipe', 'pipe'], ...spawnOptions});
  let output = '', origin, exited = false, failure;
  child.stdin.on('error', () => {});
  for (const stream of [child.stdout, child.stderr]) stream.on('data', chunk => {
    output = (output + chunk).slice(-16000);
    origin ||= output.match(originPattern)?.[0]?.match(/http:\/\/127\.0\.0\.1:\d+/)?.[0];
  });
  child.on('error', error => { failure = error; });
  const closed = new Promise(resolve => child.once('close', () => { exited = true; resolve(); }));
  return {
    child,
    closed,
    async ready(timeout = 45000) {
      const deadline = Date.now() + timeout;
      while (Date.now() < deadline) {
        if (failure || exited) throw new Error(`Game server could not start. ${failure?.message || output}`);
        if (origin) {
          try {
            const response = await fetch(origin + '/api/health', {signal: AbortSignal.timeout(1000)});
            const health = await response.json();
            if (response.ok && health.ok && health.game === 'Black Ledger') return origin;
          } catch {}
        }
        await delay(100);
      }
      throw new Error(`The local service did not become ready within ${Math.round(timeout/1000)} seconds. ` + output);
    },
    async stop() {
      if (exited) return;
      child.stdin.end();
      await Promise.race([closed, delay(6000)]);
      if (!exited) { child.kill('SIGKILL'); await closed; }
    },
    get output() { return output; },
  };
}
module.exports = {startBackend};
