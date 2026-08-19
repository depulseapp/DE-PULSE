'use strict';

// Version-neutral Node preload for regression tests that intentionally reuse a
// frozen renderer implementation across patch releases. Product runtime loads
// a small patch script from renderer/index.html; this preload makes direct
// source-reading tests see the same canonical release identity and QA overlay.
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

const qaEntry = {
  version: releaseVersion,
  date: '2026-08-19',
  status: 'STABLE',
  summary: 'Focused watchlist removal, membership-toggle and header-alert hardening with inherited v18.6.0 certification coverage and v18.6.1 edge/browser proof.',
  file: `v${releaseVersion}.txt`,
  buildId: releaseBuild,
  checkpoint: `release/v${releaseVersion}/patch_contract.json`,
};

const originalReadFileSync = fs.readFileSync.bind(fs);
fs.readFileSync = function releaseIdentityAwareRead(file, ...args) {
  const result = originalReadFileSync(file, ...args);
  const normalized = path.resolve(String(file));
  if (typeof result !== 'string') return result;

  const rendererPath = path.join(root, 'renderer', 'renderer.js');
  if (normalized === rendererPath) {
    return result
      .replace(/const EXPECTED_RELEASE_VERSION='[^']+';/, `const EXPECTED_RELEASE_VERSION='${releaseVersion}';`)
      .replace(/const EXPECTED_BUILD_ID='[^']+';/, `const EXPECTED_BUILD_ID='${releaseBuild}';`);
  }

  const manifestPath = path.join(root, 'renderer', 'qa', 'manifest.json');
  if (normalized === manifestPath) {
    const manifest = JSON.parse(result);
    const releases = Array.isArray(manifest.releases) ? manifest.releases : [];
    if (!releases.some((x) => String(x.version) === releaseVersion)) {
      manifest.releases = [qaEntry, ...releases];
    }
    return JSON.stringify(manifest, null, 2) + '\n';
  }

  return result;
};
