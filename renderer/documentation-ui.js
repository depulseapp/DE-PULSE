'use strict';

(function establishDocumentationOwner(){
  const OWNER='renderer/documentation-ui.js';
  const OWNER_VERSION=2;
  const legacyFallback={
    renderMarkdown:typeof renderMarkdown==='function',
    hydrateDocumentation:typeof hydrateDocumentation==='function',
    renderDocumentation:typeof renderDocumentation==='function'
  };
  const architectureOwner=globalThis.__DE_PULSE_DOCUMENTATION_ARCHITECTURE__?.owner||'';
  if(!architectureOwner||typeof architectureDiagram!=='function'){
    throw new Error('Canonical documentation architecture owner is required before Documentation UI initialization.');
  }

  function documentationMarkdown(md){
    const lines=String(md||'').split(/\r?\n/),out=[];
    let list=false,code=false,buf=[];
    const flushList=()=>{if(list){out.push('</ul>');list=false}};
    const flushCode=()=>{if(code){out.push(`<pre class="doc-code">${esc(buf.join('\n'))}</pre>`);buf=[];code=false}};
    for(const raw of lines){
      const line=raw.trimEnd();
      if(line.startsWith('```')){flushList();if(code)flushCode();else{code=true;buf=[]}continue}
      if(code){buf.push(raw);continue}
      const diagram=line.match(/^\[\[diagram:([a-z-]+)\]\]$/i);
      if(diagram){flushList();out.push(architectureDiagram(diagram[1].toLowerCase()));continue}
      if(/^###\s+/.test(line)){flushList();out.push(`<h3>${docInline(line.replace(/^###\s+/,''))}</h3>`)}
      else if(/^##\s+/.test(line)){flushList();out.push(`<h2>${docInline(line.replace(/^##\s+/,''))}</h2>`)}
      else if(/^#\s+/.test(line)){flushList();out.push(`<h1>${docInline(line.replace(/^#\s+/,''))}</h1>`)}
      else if(/^[-*]\s+/.test(line)){if(!list){out.push('<ul>');list=true}out.push(`<li>${docInline(line.replace(/^[-*]\s+/,''))}</li>`)}
      else if(!line.trim()){flushList()}
      else{flushList();out.push(`<p>${docInline(line)}</p>`)}
    }
    flushList();flushCode();return out.join('');
  }

  async function documentationHydrate(){
    const kind=documentationTab;
    if(docCache[kind]!=null)return;
    docCache[kind]='';
    try{
      const response=await fetch(`/docs/${kind}.md`,{cache:'no-store'});
      if(!response.ok)throw new Error(`Documentation unavailable (${response.status})`);
      docCache[kind]=await response.text();
    }catch(error){
      docCache[kind]=`# Documentation unavailable\n${error.message}`;
    }
    if(page==='documentation')render();
  }

  function documentationRender(){
    const content=docCache[documentationTab];
    const head=`<header class="page-head no-regime"><div class="page-head-copy"><h1>Documentation</h1><p>User, developer, and current data-capability references bundled with this ${DOC_BRAND} build.</p></div></header>`;
    return `<div class="page documentation-page" data-render-owner="documentation-ui"><section class="documentation-owner-marker" hidden data-documentation-owner="${OWNER}" data-architecture-owner="${architectureOwner}"></section>${head}<section class="card documentation-shell"><div class="doc-tabs" role="tablist"><button class="btn ${documentationTab==='user'?'primary':''}" data-doc-tab="user" role="tab">User Documentation</button><button class="btn ${documentationTab==='developer'?'primary':''}" data-doc-tab="developer" role="tab">Developer Documentation</button><button class="btn ${documentationTab==='limitations'?'primary':''}" data-doc-tab="limitations" role="tab">Capabilities &amp; Limitations</button></div><article id="doc-content" class="doc-content">${content?documentationMarkdown(content):'<div class="empty">Loading documentation…</div>'}</article></section></div>`;
  }

  renderMarkdown=documentationMarkdown;
  hydrateDocumentation=documentationHydrate;
  renderDocumentation=documentationRender;

  const registry=globalThis.__DE_PULSE_RENDERER_OWNERS__||(globalThis.__DE_PULSE_RENDERER_OWNERS__={});
  registry.documentation={
    owner:OWNER,
    version:OWNER_VERSION,
    state:'ACTIVE_CANONICAL_OWNER_WITH_LEGACY_RENDER_FALLBACK',
    responsibilities:['markdown','documentation-hydration','documentation-view'],
    dependencies:['shared-ui-state','shared-escaping-branding','documentation-architecture'],
    architectureOwner,
    legacyFallbackPresent:legacyFallback,
    deletionGate:'Remove legacy monolith Documentation render/hydration implementations only after direct source extraction/equivalence evidence proves no fallback consumer remains.'
  };
  globalThis.__DE_PULSE_DOCUMENTATION_UI__={
    owner:OWNER,
    version:OWNER_VERSION,
    renderMarkdown:documentationMarkdown,
    hydrate:documentationHydrate,
    render:documentationRender,
    registry:()=>registry.documentation
  };
})();
