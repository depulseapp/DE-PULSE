const fs=require('fs'),assert=require('assert');
const s=fs.readFileSync('renderer/renderer.js','utf8');
for(const token of ['Context & Alternative Intelligence','Sentiment Composite','Market / Sector Heat Map','Gamma Exposure Context','UNTRUSTED COMMUNITY INTELLIGENCE','Oil / Energy','data-community-add','data-community-delete','deterministic impact NONE']) assert(s.includes(token),`renderer missing ${token}`);
assert(s.includes('alternativeIntelligence:runtime?.alternativeIntelligence||{}'),'AI client context missing alternative intelligence');
assert(!s.includes('data-page="community"')&&!s.includes('data-page="gex"'),'v16.5 must not create top-level Community/GEX pages');
console.log('v16.5 renderer context/alternative integration tests: PASS');
