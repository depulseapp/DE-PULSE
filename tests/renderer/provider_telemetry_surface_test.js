'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');

class El{
  constructor(tag='div'){this.tag=tag;this.children=[];this.dataset={};this.className='';this.textContent='';this.removed=false;}
  append(...xs){this.children.push(...xs)}
  appendChild(x){this.children.push(x);return x}
  remove(){this.removed=true}
  querySelector(sel){
    if(sel==='span'||sel==='b'||sel==='small')return this.children.find(x=>x.tag===sel)||null;
    return null;
  }
}
function card(label){const e=new El();const s=new El('span'),b=new El('b'),sm=new El('small');s.textContent=label;e.append(s,b,sm);return e}
const grid=new El();grid.append(card('Finnhub Requests'));
let semanticCards=[];
const document={
  createElement(tag){return new El(tag)},
  querySelector(sel){return sel==='.performance-observability .maintenance-health-grid'?grid:null},
  querySelectorAll(sel){return sel==='.provider-usefulness-card'?semanticCards.filter(x=>!x.removed):[]}
};
const context={console,globalThis:null,window:{addEventListener(){}},document,requestAnimationFrame(fn){fn()},afterRender(){},authPrincipal:{role:'ADMIN'},runtime:{runtimeLoad:{providerRequests:[{provider:'Finnhub',requestsLastMinute:7,inFlight:1,successes:9,errors:1,successPct:90,p50LatencyMs:35,p95LatencyMs:90,averageLatencyMs:44,rateLimited:1,peakInFlight:2}],providerUsefulness:[{provider:'Finnhub',state:'OBSERVING',agreementPct:80,eligibleSamples:10,crossSourceSamples:8,agreementSamples:6,conflictSamples:2,singleSourceSamples:2,canonicalSelections:7,excludedSamples:1}]}}};
context.globalThis=context;
const originalCreate=document.createElement.bind(document);
document.createElement=(tag)=>{const e=originalCreate(tag);if(tag==='div')semanticCards.push(e);return e};
vm.createContext(context);
vm.runInContext(fs.readFileSync('renderer/provider-telemetry-ui.js','utf8'),context,{filename:'renderer/provider-telemetry-ui.js'});
vm.runInContext('afterRender()',context);
const reqCard=grid.children[0];
assert.equal(reqCard.dataset.providerTransportReliability,'true');
assert(reqCard.querySelector('b').textContent.includes('7/min'));
assert(reqCard.querySelector('small').textContent.includes('Transport reliability'));
const semantic=grid.children.find(x=>x.className==='provider-usefulness-card');
assert(semantic,'privileged role must receive semantic usefulness card');
assert(semantic.querySelector('b').textContent.includes('80.0% cross-source agreement'));
assert(semantic.querySelector('small').textContent.includes('ADVISORY ONLY'));
assert(semantic.querySelector('small').textContent.includes('no routing effect'));

context.authPrincipal={role:'USER'};
vm.runInContext('afterRender()',context);
assert(semantic.removed,'non-privileged role must remove provider usefulness projection');
console.log('Provider telemetry/usefulness privileged surface functional regression PASS');
