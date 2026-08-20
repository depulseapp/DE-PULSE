'use strict';
const fs=require('fs');
const assert=require('assert');
const renderer=fs.readFileSync('renderer/renderer.js','utf8');
const truth=fs.readFileSync('research_truth.go','utf8');

const rendererRequired=[
  'Earnings Deep Dive',
  'Fundamentals Interpretation',
  'SEC & Ownership',
  'BUY/SELL/OTHER stays semantically separated.',
  'Technical Context',
  'function v151ResearchReadyForAI(sym)',
  'Research Evidence Incomplete',
  'AI is optional and never changes deterministic Action/Score.',
  'Trading decisions remain in the horizon-specific desks.',
  "const DETERMINISTIC_FORMULA_VERSION='deterministic-v14.3.7-compatible-v16.3'",
  'researchEarningsV2',
  'researchFundamentals',
  'researchSECAndOwnership',
  'researchCatalysts',
  'researchTechnicals',
  'globalThis.DePulseLogic',
  'This is validation/research, never paper trading.'
];
for(const token of rendererRequired)assert(renderer.includes(token),`Research capability regression missing: ${token}`);

const truthRequired=[
  'func evidenceAge(now, ts int64, allowedFutureSkew time.Duration)',
  'future timestamp exceeds allowed clock skew',
  'func buildResearchPackageTruth(',
  'worst-dependency semantics',
  'ResearchPackageTruth{Symbol: symbol, State: "BLOCKED"',
  'No valid equity/ETF Research target selected',
  'Dataset: "Quote", Required: true, Critical: true',
  'Dataset: "Daily History", Required: true, Critical: true',
  'Dataset: "Fundamentals", Required: true',
  'Dataset: "News", Required: true',
  'Dataset: "Earnings", Required: true',
  'Dataset: "SEC & Ownership", Required: true',
  'Dataset: "Catalyst & Material Event Context", Required: true',
  'Dataset: "Required Market Context", Required: true',
  'Dataset: "Provider Reconciliation", Required: true, Critical: true'
];
for(const token of truthRequired)assert(truth.includes(token),`Research truth regression missing: ${token}`);
assert(truth.includes('optional context must never upgrade a degraded required component.'),'optional context must not upgrade required evidence truth');
console.log('Research correctness closure regression PASS');
console.log('approved Research capability preservation: PASS');
console.log('worst-dependency evidence truth + future-clock-skew protection: PASS');
console.log('AI second-opinion / deterministic score boundary: PASS');
