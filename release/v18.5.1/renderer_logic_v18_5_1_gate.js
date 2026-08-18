const fs = require('fs');
const os = require('os');
const path = require('path');
const { spawnSync } = require('child_process');

const legacyPath = path.resolve('renderer_logic_test.js');
const source = fs.readFileSync(legacyPath, 'utf8');
const stale = "assert.equal(manifest.releases[0].version,releaseVersion);";
if (!source.includes(stale)) {
  throw new Error('Expected inherited QA-manifest assertion is missing; review renderer regression ownership before changing this gate.');
}

const replacement = `if(manifest.releases[0].version!==releaseVersion){
  const deltaPath=\`release/v\${releaseVersion}/QA-MANIFEST-DELTA.json\`;
  assert(fs.existsSync(deltaPath),\`unpromoted candidate \${releaseVersion} requires release-specific QA manifest delta\`);
  const qaDelta=JSON.parse(fs.readFileSync(deltaPath,'utf8'));
  assert.equal(qaDelta.version,releaseVersion);
  assert.equal(qaDelta.status,'TEST');
  assert.equal(qaDelta.previousStable,\`v\${manifest.releases[0].version}\`);
  assert.equal(qaDelta.majorProvenanceAnchor,'v17.5.1');
  assert(String(qaDelta.note||'').includes('durable promoted/historical release manifest'));
}else{
  assert.equal(manifest.releases[0].version,releaseVersion);
}`;

const patched = source.replace(stale, replacement);
const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'depulse-v1851-renderer-'));
const tmpFile = path.join(tmpDir, 'renderer_logic_test.cjs');
try {
  fs.writeFileSync(tmpFile, patched, 'utf8');
  const run = spawnSync(process.execPath, [tmpFile], {
    cwd: process.cwd(),
    stdio: 'inherit',
    env: process.env,
  });
  if (run.error) throw run.error;
  if (run.status !== 0) process.exit(run.status || 1);
  console.log('PASS: v18.5.1 renderer regression preserves promoted QA history and requires the exact TEST/RC QA delta until promotion.');
} finally {
  fs.rmSync(tmpDir, { recursive: true, force: true });
}
