'use strict';
/* DE.PULSE v18.9.1 release identity overlay.
   Platform-specific macOS native-window reliability corrective; shared product,
   provider-routing, persistence and No Execution semantics remain unchanged. */
var DEPULSE_RELEASE_VERSION = '18.9.1';
var DEPULSE_RELEASE_BUILD_ID = 'v18.9.1-stable-20260821';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-21',
  status: 'STABLE',
  summary: 'macOS native-window reliability corrective: canonical JXA/Cocoa/WKWebView startup, fresh/warm profile reuse, protocol-resolution regression and deterministic cleanup evidence.',
  file: 'v18.9.1.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.9.1/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1891(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.9.1 renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1891 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1891(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1891();
  };
}
