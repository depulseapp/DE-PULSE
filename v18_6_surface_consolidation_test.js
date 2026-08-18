'use strict';

const fs = require('fs');
const vm = require('vm');

const extension = fs.readFileSync('renderer/surface-consolidation-v18.6.js', 'utf8');
const css = fs.readFileSync('renderer/surface-consolidation-v18.6.css', 'utf8');
const index = fs.readFileSync('renderer/index.html', 'utf8');
const renderer = fs.readFileSync('renderer/renderer.js', 'utf8');

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

console.log('PASS: v18.6 demotes Market Activity to an optional supporting-input drilldown and retires standalone News/Earnings/Filings navigation into canonical Research/Market Intelligence destinations.');
