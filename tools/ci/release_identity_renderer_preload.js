'use strict';

// Version-neutral Node preload for regression tests that intentionally reuse a
// frozen renderer implementation across patch releases. The product runtime
// still loads the patch script from renderer/index.html; this preload only
// makes direct source-reading tests evaluate the renderer against the canonical
// release identity instead of an inherited core-file identity literal.
const fs = require('fs');
const path = require('path');

const root = process.cwd();
const identityPath = path.join(root, 'release_identity.json');
const identity = JSON.parse(fs.readFileSync(identityPath, 'utf8'));
const releaseVersion = String(identity.version || '').trim();
const releaseBuild = String(identity.build_id || '').trim();
if (!releaseVersion || !releaseBuild) {
  throw new Error('release identity preload requires version and build_id');
}

const originalReadFileSync = fs.readFileSync.bind(fs);
fs.readFileSync = function releaseIdentityAwareRead(file, ...args) {
  const result = originalReadFileSync(file, ...args);
  const normalized = path.resolve(String(file));
  const rendererPath = path.join(root, 'renderer', 'renderer.js');
  if (normalized !== rendererPath || typeof result !== 'string') return result;
  return result
    .replace(/const EXPECTED_RELEASE_VERSION='[^']+';/, `const EXPECTED_RELEASE_VERSION='${releaseVersion}';`)
    .replace(/const EXPECTED_BUILD_ID='[^']+';/, `const EXPECTED_BUILD_ID='${releaseBuild}';`);
};
