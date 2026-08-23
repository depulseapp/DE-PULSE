'use strict';

function v186DocumentationRole(){return String(authPrincipal?.role||'').toUpperCase()}
function v186CanViewDeveloperDocumentation(){return ['SUPER_OWNER','OWNER','ADMIN'].includes(v186DocumentationRole())}
function v186NormalizeDocumentationTab(){
  if(documentationTab==='developer'&&!v186CanViewDeveloperDocumentation())documentationTab='user';
  return documentationTab;
}

const v186BaseHydrateDocumentation=hydrateDocumentation;
hydrateDocumentation=async function(){
  v186NormalizeDocumentationTab();
  return v186BaseHydrateDocumentation();
};

const v186BaseRenderDocumentation=renderDocumentation;
renderDocumentation=function(){
  v186NormalizeDocumentationTab();
  const html=v186BaseRenderDocumentation();
  if(v186CanViewDeveloperDocumentation())return html;
  return html.replace(/<button\b[^>]*data-doc-tab="developer"[^>]*>[\s\S]*?<\/button>/i,'');
};

const v186DocumentationOwner=globalThis.__DE_PULSE_RENDERER_OWNERS__?.documentation;
if(v186DocumentationOwner){
  const decorators=Array.isArray(v186DocumentationOwner.decorators)?v186DocumentationOwner.decorators:[];
  if(!decorators.includes('renderer/documentation-access.js'))decorators.push('renderer/documentation-access.js');
  v186DocumentationOwner.decorators=decorators;
  v186DocumentationOwner.accessPolicy='developer documentation requires SUPER_OWNER, OWNER, or ADMIN';
}

globalThis.__v186DocumentationAccess={
  role:v186DocumentationRole,
  canViewDeveloper:v186CanViewDeveloperDocumentation,
  normalizeTab:v186NormalizeDocumentationTab
};
