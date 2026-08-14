const fs=require('fs'),assert=require('assert');
const s=fs.readFileSync('renderer/renderer.js','utf8');
for(const token of ['Bull case','Base case','Bear case','Not win probability','Evidence package','Canonical snapshot','External-content safety','untrusted data','Routing & Cost Policy','Material Evidence Package caching','AI cannot change Action/Score']) assert(s.includes(token),`renderer missing ${token}`);
for(const mode of ['manual','efficient','balanced','deep']) assert(s.includes(`value="${mode}"`),`routing mode missing ${mode}`);
assert(s.includes('evidencePackageId')&&s.includes('evidenceSnapshotId')&&s.includes('evidenceIds'),'evidence grounding fields missing');
console.log('v16.4 renderer AI evidence/safety/routing tests: PASS');
