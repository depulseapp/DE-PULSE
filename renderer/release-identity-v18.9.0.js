'use strict';
/* DE.PULSE v18.9.0 release identity overlay.
   TradeInsight adds adaptive SHADOW/fallback evidence through existing canonical
   owners without changing protected trading semantics or the No Execution boundary. */
var DEPULSE_RELEASE_VERSION = '18.9.0';
var DEPULSE_RELEASE_BUILD_ID = 'v18.9.0-stable-20260821';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-21',
  status: 'STABLE',
  summary: 'TradeInsight adaptive provider integration: SHADOW Congressional evidence, adjusted daily OHLCV fallback/backfill, bounded admission-controlled history fan-out, and shared provider telemetry with direct SEC authority preserved.',
  file: 'v18.9.0.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.9.0/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1890(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.9.0 renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1890 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1890(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1890();
  };
}
