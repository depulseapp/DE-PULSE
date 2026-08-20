'use strict';
const fs=require('fs');
const assert=require('assert');
const renderer=fs.readFileSync('renderer/renderer.js','utf8');
const liveDom=fs.readFileSync('renderer/live-dom-reconcile.js','utf8');

const required=[
  "function v151ResearchReadyForAI(sym)",
  "String(pkg.state||'').toUpperCase()!=='FRESH'",
  "pkg.evidenceSnapshotId||runtime?.evidenceSnapshot?.id",
  "function v151AIStatus(sym)",
  "Research Evidence Incomplete",
  "AI is optional and never changes deterministic Action/Score.",
  "Research freshness",
  "return context preserved.",
  "maxlength=\"8\"",
  "['overview','Overview']",
  "['earnings','Earnings']",
  "['fundamentals','Fundamentals']",
  "['sec','SEC & Ownership']",
  "['catalysts','Catalysts']",
  "['technical','Technical Context']",
  "function v151SelectResearchTarget(sym)",
  "/^[A-Z][A-Z0-9.\\-]{0,7}$/",
  "api('/api/research/refresh',{symbol:v})",
  "v151ResearchHydrationError=err.message||'Research refresh failed'",
  "globalThis.DePulseLogic",
  "renderResearch,researchOverview,researchFundamentals,researchCatalysts,researchTechnicals,researchSECAndOwnership,researchEarningsV2"
];
for(const token of required)assert(renderer.includes(token),`Research IA contract missing: ${token}`);
assert(renderer.includes("Trading decisions remain in the horizon-specific desks."),'Research must remain decision-support, not execution');
assert(renderer.includes("data-research-back")&&renderer.includes("['dashboard','discovery','day','swing','long'].includes(researchOrigin)"),'Research return-origin containment missing');
for(const token of ['selectionStart','selectionEnd','window.scrollX','window.scrollY','main.scrollTop','active.focus({preventScroll:true})','setSelectionRange(selectionStart, selectionEnd)']){
  assert(liveDom.includes(token),`live DOM interaction preservation missing: ${token}`);
}
assert(liveDom.includes("if (!symbol) return originalScheduleLiveRender(changedSymbol)"),'structural events must retain authoritative full-render path');
console.log('Research information architecture regression PASS');
console.log('Research target/freshness/recovery/containment: PASS');
console.log('Research focus/selection/scroll preservation: PASS');
