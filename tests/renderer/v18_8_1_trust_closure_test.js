'use strict';
const {execFileSync}=require('child_process');

require('./research_information_architecture_test.js');
require('./symbol_desk_correctness_test.js');
require('./research_correctness_closure_test.js');

for(const gate of [
  'tools/ci/v18_8_1_readiness_contract.py',
  'tools/ci/v18_8_1_freshness_contract.py',
  'tools/ci/v18_8_1_zero_miss_reconciliation_gate.py'
]){
  execFileSync('python3',[gate],{cwd:process.cwd(),stdio:'inherit'});
}
console.log('v18.8.1 user-trust closure PASS');
console.log('Research IA + Symbol/Desk + Readiness + Freshness + Research correctness + Zero-Miss: PASS');
