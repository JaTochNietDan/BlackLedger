const test = require('node:test');
const assert = require('node:assert/strict');
const {startBackend} = require('./backend.cjs');

test('desktop parent waits for its own healthy server and stops it through stdin', async () => {
 const script = `const http=require('node:http');const s=http.createServer((q,r)=>r.end(JSON.stringify({ok:true,game:'Black Ledger'})));s.listen(0,'127.0.0.1',()=>console.log('http://127.0.0.1:'+s.address().port));process.stdin.resume();process.stdin.on('end',()=>s.close());`;
 const backend = startBackend(process.execPath, ['-e', script]);
 try {
  const origin = await backend.ready(5000);
  assert.match(origin, /^http:\/\/127\.0\.0\.1:\d+$/);
 } finally { await backend.stop(); }
 assert.equal(backend.child.exitCode, 0);
});
test('a failed game process is reported rather than opening an unrelated server', async () => {
 const backend = startBackend(process.execPath, ['-e', `console.error('save unavailable');process.exit(2);`]);
 await assert.rejects(backend.ready(5000), /save unavailable/);
 await backend.stop();
});
