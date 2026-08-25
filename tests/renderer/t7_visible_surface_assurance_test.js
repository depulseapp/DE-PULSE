'use strict';

const fs = require('fs');
const assert = require('assert');

// Reuse existing canonical renderer contracts instead of duplicating their
// implementation logic inside a closure-only test.
require('./provider_telemetry_surface_test.js');
require('./symbol_desk_correctness_test.js');
require('./surface_consolidation_test.js');

const renderer = fs.readFileSync('renderer/renderer.js', 'utf8');
const css = fs.readFileSync('renderer/styles.css', 'utf8');

// Opportunity Radar is the canonical visible merge point for unusual
// participation / dislocation intelligence. It is decision support only.
assert(renderer.includes('opportunity-radar-head'), 'Opportunity Radar visible table header missing');
assert(renderer.includes('<h2>Unusual volume, volatility & rapid dislocation</h2>'), 'Opportunity Radar must explain the merged rapid-dislocation responsibility');
assert(renderer.includes('Decision support only; no execution.'), 'Opportunity Radar must preserve No Execution copy');
assert(css.includes('.opportunity-radar-head,.opportunity-radar-row'), 'Opportunity Radar responsive row ownership missing');
assert(css.includes('.opportunity-radar-row{grid-template-columns:repeat(2,minmax(0,1fr))'), 'Opportunity Radar narrow-layout contract missing');

// Rapid Move remains merged into the existing notification / Opportunity Radar
// presentation rather than becoming a new navigation page or execution surface.
assert(renderer.includes("if(d.type==='rapid-move'&&d.event)"), 'Rapid Move SSE event must have an explicit user-facing projection');
assert(renderer.includes('window.__lastRapidMoveToast!==toastKey'), 'Rapid Move user notification must dedupe stable event/state transitions');
assert(renderer.includes('toast(`${ev.symbol} · RAPID MOVE'), 'Rapid Move alert must remain clearly identified in the existing notification surface');
assert(renderer.includes("ev.catalystState==='CONFIRMED'"), 'Rapid Move alert must surface confirmed catalyst context when available');
assert(renderer.includes("ev.catalystState==='UNEXPLAINED'?' · Catalyst validating':''"), 'Rapid Move alert must truthfully show catalyst validation instead of inventing a cause');

// Identity chrome is a visible shell responsibility. The role-responsive owner
// proves placement/visibility; these assertions close the action semantics.
assert(renderer.includes("$('#identity-signout')?.addEventListener('click'"), 'authenticated shell must bind Sign Out action');
assert(renderer.includes("await api('/api/auth/logout',{})"), 'Sign Out must use the canonical logout route');
assert(renderer.includes("toast('Sign Out Failed',err.message,'error')"), 'Sign Out failure must remain visible and recoverable');

console.log('T7 visible-surface composition regression PASS');
console.log('Opportunity Radar + Rapid Move + provider semantics + Master Symbols + identity Sign Out: PASS');
