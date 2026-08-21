'use strict';

const fs = require('fs');
const vm = require('vm');
const marketTruth = require('../../renderer/market-intelligence-truth.js');

const extension = fs.readFileSync('renderer/surface-consolidation-v18.6.js', 'utf8');
const css = fs.readFileSync('renderer/surface-consolidation-v18.6.css', 'utf8');
const index = fs.readFileSync('renderer/index.html', 'utf8');
const renderer = fs.readFileSync('renderer/renderer.js', 'utf8');
const marketTruthSource = fs.readFileSync('renderer/market-intelligence-truth.js', 'utf8');

function assert(condition, message) {
  if (!condition) throw new Error(message);
}

assert(index.includes('surface-consolidation-v18.6.js'), 'v18.6 consolidation script is not loaded');
assert(index.includes('surface-consolidation-v18.6.css'), 'v18.6 consolidation CSS is not loaded');
assert(extension.includes('discovery-supporting-input'), 'supporting Market Activity drilldown missing');
assert(extension.includes("legacyEvidencePages = new Set(['news', 'earnings', 'filings'])"), 'legacy route set missing');
assert(extension.includes("baseSetPage('market-intelligence')"), 'legacy market-wide redirect missing');
assert(extension.includes("kind === 'filings' || kind === 'sec' ? 'filings'"), 'ticker filing route must target Research Filings');
assert(extension.includes("kind === 'news' || kind === 'earnings' ? 'catalysts'"), 'ticker News/Earnings route must target Research Catalysts');
assert(!extension.includes('/api/'), 'presentation consolidation must not create a new data acquisition path');
assert(css.includes('.discovery-supporting-input > summary'), 'Market Activity must be collapsed behind an optional details summary');
assert(renderer.includes('discovery-market-activity full-width'), 'base historical renderer marker unexpectedly disappeared; extension contract must consciously demote it');

// v18.8.2 issue #57: DATA DEGRADED is an evidence state, not a calculated
// zero. Keep evaluated scores numeric, but never present unavailable market
// evidence as a meaningful 0/100 on Dashboard or Market Intelligence.
assert(index.includes('<script src="market-intelligence-truth.js?v=18.8.2"></script>'), 'v18.8.2 Market Intelligence truth layer is not loaded');
assert(index.indexOf('market-intelligence-truth.js?v=18.8.2') > index.indexOf('renderer.js?v=18.8.1'), 'Market Intelligence truth layer must load after the primary renderer');
assert(marketTruthSource.includes("label.textContent.trim() !== 'Tradeability'"), 'Dashboard Tradeability row truth reconciliation missing');
assert(marketTruthSource.includes("label.textContent.trim() !== 'Market Tradeability'"), 'Market Intelligence Tradeability card truth reconciliation missing');
assert(marketTruth.tradeabilityScoreLabel('DATA DEGRADED', 0) === 'UNAVAILABLE', 'DATA DEGRADED must not render as 0/100');
assert(marketTruth.tradeabilityScoreLabel('UNAVAILABLE', 0) === 'UNAVAILABLE', 'UNAVAILABLE must remain non-numeric');
assert(marketTruth.tradeabilityScoreLabel('WAIT', 0) === '0/100', 'an evaluated zero remains numeric when evidence is current');
assert(marketTruth.tradeabilityScoreLabel('SELECTIVE', 62) === '62/100', 'evaluated tradeability score must remain numeric');

const calls = [];
const sandbox = {
  console,
  window: {},
  page: 'discovery',
  state: { ui: { selectedTicker: 'AAPL' } },
  researchSymbol: 'AAPL',
  researchTab: 'overview',
  runtime: {
    marketActivity: {
      status: 'AVAILABLE',
      mostActive: [{ symbol: 'AAPL', changePercent: 1.2 }, { symbol: 'MSFT', changePercent: -0.4 }],
      gainers: [{ symbol: 'AAPL', changePercent: 1.2 }, { symbol: 'NVDA', changePercent: 2.5 }],
      losers: [{ symbol: 'TSLA', changePercent: -3.1 }]
    }
  },
  renderDiscovery: () => 'before<section class="card discovery-market-activity full-width"><div>prominent legacy surface</div></section>after',
  setPage: target => { calls.push(['setPage', target]); return target; },
  bindDynamic: () => { calls.push(['bind']); },
  render: () => { calls.push(['render']); },
  openResearch: (symbol, origin) => { calls.push(['openResearch', symbol, origin]); sandbox.researchSymbol = symbol; sandbox.page = 'research'; },
  v151SelectResearchTarget: undefined,
  $$: () => [],
  esc: value => String(value),
  num: value => Number(value) || 0,
  pctSigned: value => `${Number(value) >= 0 ? '+' : ''}${Number(value).toFixed(2)}%`
};
vm.createContext(sandbox);
vm.runInContext(extension, sandbox, { filename: 'surface-consolidation-v18.6.js' });

const discovery = sandbox.renderDiscovery();
assert(!discovery.includes('discovery-market-activity full-width'), 'prominent Market Activity card survived v18.6 rendering');
assert(discovery.includes('<details class="discovery-supporting-input"'), 'Market Activity was not demoted to optional drilldown');
assert(discovery.includes('Market Activity Seeds'), 'supporting seed explanation missing');
assert(discovery.includes('Scanner / Opportunity Radar seeds'), 'supporting-input purpose is not explained');

sandbox.setPage('news');
assert(calls.at(-1)[0] === 'setPage' && calls.at(-1)[1] === 'market-intelligence', 'direct News route did not redirect to Market Intelligence');
sandbox.setPage('earnings');
assert(calls.at(-1)[1] === 'market-intelligence', 'direct Earnings route did not redirect to Market Intelligence');
sandbox.setPage('filings');
assert(calls.at(-1)[1] === 'market-intelligence', 'direct Filings route did not redirect to Market Intelligence');
sandbox.setPage('day');
assert(calls.at(-1)[1] === 'day', 'non-legacy navigation was changed');

assert(sandbox.window.__v186SurfaceConsolidation.evidenceResearchTab('sec') === 'filings', 'SEC ticker evidence must resolve to Research Filings');
assert(sandbox.window.__v186SurfaceConsolidation.evidenceResearchTab('filings') === 'filings', 'Filings ticker evidence must resolve to Research Filings');
assert(sandbox.window.__v186SurfaceConsolidation.evidenceResearchTab('earnings') === 'catalysts', 'Earnings ticker evidence must resolve to Research Catalysts');
assert(sandbox.window.__v186SurfaceConsolidation.evidenceResearchTab('news') === 'catalysts', 'News ticker evidence must resolve to Research Catalysts');

console.log('PASS: v18.6 surface consolidation + v18.8.2 Market Intelligence degraded-score truth contracts.');
