'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');
const crypto=require('crypto');

const ownerSource=fs.readFileSync('renderer/market-header-ui.js','utf8');
const index=fs.readFileSync('renderer/index.html','utf8');

function gitBlobToken(path){
  const data=fs.readFileSync(path);
  const header=Buffer.from(`blob ${data.length}\0`,'utf8');
  return crypto.createHash('sha1').update(header).update(data).digest('hex').slice(0,16);
}

class Element {
  constructor(tagName='div') {
    this.tagName=String(tagName).toUpperCase();
    this.id='';
    this.className='';
    this.parentElement=null;
    this.children=[];
    this.attributes={};
    this.textContent='';
    this.title='';
  }
  get firstChild(){ return this.children[0]||null; }
  setAttribute(name,value){ this.attributes[name]=String(value); }
  appendChild(child){
    if(child.parentElement) child.parentElement.removeChild(child);
    this.children.push(child);
    child.parentElement=this;
    return child;
  }
  insertBefore(child,reference){
    if(child.parentElement) child.parentElement.removeChild(child);
    const index=this.children.indexOf(reference);
    if(index<0) return this.appendChild(child);
    this.children.splice(index,0,child);
    child.parentElement=this;
    return child;
  }
  removeChild(child){
    const index=this.children.indexOf(child);
    if(index>=0) this.children.splice(index,1);
    child.parentElement=null;
    return child;
  }
  insertAdjacentElement(position,child){
    assert.strictEqual(position,'afterend');
    assert(this.parentElement,'afterend requires a parent');
    const siblings=this.parentElement.children;
    const selfIndex=siblings.indexOf(this);
    if(child.parentElement) child.parentElement.removeChild(child);
    siblings.splice(selfIndex+1,0,child);
    child.parentElement=this.parentElement;
    return child;
  }
  querySelector(selector){ return querySelectorFrom(this,selector,false); }
}

function matches(element,selector){
  if(selector.startsWith('#')) return element.id===selector.slice(1);
  if(selector.startsWith('.')) return String(element.className).split(/\s+/).includes(selector.slice(1));
  return element.tagName.toLowerCase()===selector.toLowerCase();
}
function walk(root,visitor,includeRoot=true){
  if(includeRoot && visitor(root)) return root;
  for(const child of root.children){
    const match=walk(child,visitor,true);
    if(match) return match;
  }
  return null;
}
function querySelectorFrom(root,selector,includeRoot=true){
  return walk(root,element=>matches(element,selector),includeRoot);
}

const root=new Element('body');
const topbar=root.appendChild(new Element('header'));
topbar.className='topbar';
function attach(id,className=''){
  const element=topbar.appendChild(new Element('div'));
  element.id=id;
  element.className=className;
  return element;
}
const session=attach('market-session-context');
const clocks=attach('','market-clocks');
const status=attach('runtime-status');
const toggle=attach('runtime-toggle');

const document={
  createElement:tagName=>new Element(tagName),
  getElementById:id=>walk(root,element=>element.id===id,true),
  querySelector:selector=>querySelectorFrom(root,selector,true),
};

let baseCalls=0;
let lastChangedSymbol='';
let healthCalls=0;
function baseUpdateChrome(changedSymbol=''){
  baseCalls++;
  lastChangedSymbol=changedSymbol;
}
function headerDataHealth(){
  healthCalls++;
  return {label:'LIVE',detail:'Fresh · canonical market evidence'};
}

const context=vm.createContext({console,document,updateChrome:baseUpdateChrome,headerDataHealth});
vm.runInContext(ownerSource,context,{filename:'market-header-ui.js'});

const owner=vm.runInContext('__DE_PULSE_MARKET_HEADER_UI__',context);
assert.strictEqual(owner.owner,'renderer/market-header-ui.js');
assert.strictEqual(owner.version,1);
assert.strictEqual(owner.registry().state,'ACTIVE_OWNER_WITH_COMPAT_ALIAS');
assert(owner.registry().responsibilities.includes('market-pulse-ribbon'));
assert(owner.registry().responsibilities.includes('market-clocks'));
assert(owner.registry().responsibilities.includes('data-runtime-control'));
assert.strictEqual(owner.registry().legacyCompatibilityFile,'renderer/header-v18.5.1.js');
assert.strictEqual(vm.runInContext('__v1851HeaderContracts===__DE_PULSE_MARKET_HEADER_UI__',context),true);

const bar=document.getElementById('market-status-bar');
assert(bar,'Market Pulse Ribbon must be created');
assert.strictEqual(root.children[root.children.indexOf(topbar)+1],bar,'ribbon must remain directly below topbar');
const content=bar.querySelector('.market-status-content');
assert(content,'ribbon content container must exist');
const summary=document.getElementById('market-data-summary');
const detail=document.getElementById('market-data-detail');
assert(summary&&detail,'market data coverage summary/detail must exist');
assert.deepStrictEqual(content.children,[session,summary,clocks,toggle],'ribbon order must remain session -> coverage -> clocks -> runtime control');
assert.strictEqual(summary.children[0],status,'runtime status must be owned by market-data summary');
assert.strictEqual(detail.textContent,'Market data state unavailable');

assert.strictEqual(owner.ensureSecondaryMarketStatus(),bar);
assert.strictEqual(owner.ensureSecondaryMarketStatus(),bar);
assert.strictEqual(root.children.filter(element=>element.id==='market-status-bar').length,1,'ensure must be idempotent');
assert.deepStrictEqual(content.children,[session,summary,clocks,toggle],'repeated ensure must not duplicate/reorder capability children');

vm.runInContext("updateChrome('SPY')",context);
assert.strictEqual(baseCalls,1,'canonical chrome owner must be invoked exactly once per update');
assert.strictEqual(lastChangedSymbol,'SPY');
assert.strictEqual(healthCalls,1,'header health must be evaluated exactly once per update');
assert.strictEqual(detail.textContent,'Fresh · canonical market evidence');
assert.strictEqual(detail.title,'Fresh · canonical market evidence');
vm.runInContext("updateChrome('QQQ')",context);
assert.strictEqual(baseCalls,2,'wrapper must not multiply base update calls');
assert.strictEqual(healthCalls,2,'wrapper must not multiply health evaluation');

const ownerTag=`<script src="market-header-ui.js?v=${gitBlobToken('renderer/market-header-ui.js')}"></script>`;
assert(index.includes(ownerTag),'index must load stable Market Header owner with content-derived cache identity');
assert(!index.includes('<script src="header-v18.5.1.js'),'version-stacked header must remain inactive runtime evidence');

console.log('Market Header capability owner regression PASS');
