const fs=require('fs'),assert=require('assert');
const s=fs.readFileSync('renderer/renderer.js','utf8');
for(const token of ['data-master-remove-all','/api/master-symbol/remove-all','Tracked Symbols Cleared','User Master Symbol Store','System Market Context','may be empty']) assert(s.includes(token),`v16.6 renderer missing ${token}`);
assert(s.includes('Remove All'));
console.log('v16.6 renderer integration / Master Symbol empty-state tests: PASS');
