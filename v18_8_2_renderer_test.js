'use strict';
const assert = require('assert');
const fs = require('fs');
const path = require('path');
const truth = require('./renderer/market-intelligence-truth.js');

assert.strictEqual(truth.tradeabilityScoreLabel('DATA DEGRADED', 0), 'UNAVAILABLE', 'degraded evidence must not render as a meaningful 0/100');
assert.strictEqual(truth.tradeabilityScoreLabel('UNAVAILABLE', 0), 'UNAVAILABLE', 'unavailable evidence must remain unavailable');
assert.strictEqual(truth.tradeabilityScoreLabel('WAIT', 0), '0/100', 'an evaluated zero score remains numeric when evidence is current');
assert.strictEqual(truth.tradeabilityScoreLabel('SELECTIVE', 62), '62/100', 'evaluated tradeability keeps its numeric score');
assert.strictEqual(truth.tradeabilityScoreLabel('TRADE NORMALLY', 83), '83/100', 'normal evaluated tradeability keeps its numeric score');

const source = fs.readFileSync(path.join(__dirname, 'renderer', 'market-intelligence-truth.js'), 'utf8');
assert(source.includes("label.textContent.trim() !== 'Tradeability'"), 'dashboard Tradeability row must be reconciled');
assert(source.includes("label.textContent.trim() !== 'Market Tradeability'"), 'Market Intelligence Tradeability card must be reconciled');
assert(source.includes('DATA DEGRADED'), 'truth layer must explicitly recognize DATA DEGRADED');
assert(source.includes('UNAVAILABLE'), 'truth layer must explicitly render unavailable score truth');

const index = fs.readFileSync(path.join(__dirname, 'renderer', 'index.html'), 'utf8');
assert(index.includes('<script src="market-intelligence-truth.js?v=18.8.2"></script>'), 'renderer index must load the v18.8.2 Market Intelligence truth layer');
assert(index.indexOf('market-intelligence-truth.js?v=18.8.2') > index.indexOf('renderer.js?v=18.8.1'), 'truth layer must load after the primary renderer');

console.log('v18.8.2 Market Intelligence renderer truth contract: PASS');
