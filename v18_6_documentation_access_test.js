'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');

const source=fs.readFileSync('renderer/documentation-access-v18.6.js','utf8');
const index=fs.readFileSync('renderer/index.html','utf8');
const context=vm.createContext({console});
vm.runInContext(`
let authPrincipal={role:'USER'};
let documentationTab='developer';
let hydrateCalls=0;
async function hydrateDocumentation(){hydrateCalls++;return documentationTab}
function renderDocumentation(){return '<div class="doc-tabs"><button class="doc-tab" data-doc-tab="user">User Documentation</button><button class="doc-tab active" data-doc-tab="developer">Developer Documentation</button><button class="doc-tab" data-doc-tab="limitations">Capabilities & Limitations</button></div><article>'+documentationTab+'</article>'}
`,context);
vm.runInContext(source,context,{filename:'documentation-access-v18.6.js'});

const userHTML=vm.runInContext('renderDocumentation()',context);
assert(!userHTML.includes('data-doc-tab="developer"'),'USER must not receive developer tab');
assert(userHTML.includes('<article>user</article>'),'forbidden programmatic developer selection must normalize to user docs');
assert.strictEqual(vm.runInContext('__v186DocumentationAccess.canViewDeveloper()',context),false);

vm.runInContext(`authPrincipal={role:'DEMO'};documentationTab='developer'`,context);
assert(!vm.runInContext('renderDocumentation()',context).includes('Developer Documentation'),'DEMO must not receive developer docs');

for(const role of ['ADMIN','OWNER','SUPER_OWNER']){
  vm.runInContext(`authPrincipal={role:'${role}'};documentationTab='developer'`,context);
  const html=vm.runInContext('renderDocumentation()',context);
  assert(html.includes('data-doc-tab="developer"'),`${role} should receive developer tab`);
  assert(html.includes('<article>developer</article>'),`${role} developer selection should remain intact`);
}

vm.runInContext(`authPrincipal={role:'USER'};documentationTab='developer'`,context);
vm.runInContext('hydrateDocumentation()',context).then(()=>{
  assert.strictEqual(vm.runInContext('documentationTab',context),'user','hydrate must not fetch forbidden developer docs');
  assert(index.includes('<script src="documentation-access-v18.6.js?v=18.5.2"></script>'),'index must load v18.6 documentation access extension');
  assert(!source.includes('/api/'),'composition extension must not create a parallel authorization API');
  console.log('v18.6 documentation access regression PASS');
}).catch(err=>{console.error(err);process.exitCode=1});
