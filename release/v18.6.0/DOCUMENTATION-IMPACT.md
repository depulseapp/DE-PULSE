# Documentation Impact Manifest — v18.6.0

**Scope owner:** G0–G16 Adaptive Delivery Process  
**Requirement:** `IMPL-18-DOC-001` / `CONVO-V17-005`  
**Development branch:** `v18.6-development`  
**Rule:** every material product/process change receives an explicit documentation disposition; silence is not an acceptable disposition.

## Audience policy

| Audience | User Documentation | Capabilities & Limitations | Developer Documentation |
|---|---|---|---|
| SUPER_OWNER | Allow | Allow | Allow |
| OWNER | Allow | Allow | Allow |
| delegated ADMIN | Allow | Allow | Allow |
| USER | Allow | Allow | Deny |
| DEMO | Allow | Allow | Deny |

The server is authoritative. Renderer composition is a usability layer only and cannot grant access. Direct `/docs/*` requests require a valid server-side session when IdentityService is active. `/docs/developer.md` additionally requires the canonical `ADMIN` role threshold, which naturally includes OWNER and SUPER_OWNER through the existing role hierarchy. No new documentation-only role system exists.

Identity-disabled local desktop operation retains its existing embedded-document compatibility behavior; it has no multi-user role boundary to enforce.

## v18.6 material-change dispositions

| Requirement / change | Audience impact | Documentation disposition | Evidence / rationale |
|---|---|---|---|
| Watchlist membership remediation | User workflow semantics clarified; no new product capability | `NO_TEXT_CHANGE` | Restores canonical persisted desk membership and removes misleading current-desk semantics; existing user intent is unchanged. Regression/browser proof is the primary evidence. |
| Shared Scanner/Radar broad snapshot broker (`IMPL-18-UTILITY-001`) | Internal architecture/performance | `DEVELOPER_EVIDENCE` | No user workflow or ranking semantics change. Architecture, freshness reuse, coalescing and diagnostics are captured in source/tests/release evidence rather than adding end-user instructions. |
| Session Intelligence Coordinator (`IMPL-18-UTILITY-002`) | Internal scheduler ownership; checkpoint semantics preserved | `DEVELOPER_EVIDENCE` | Pre-Market Prep and Market Open Prep timing/meaning remain unchanged. Coordinator ownership and catch-up behavior are source/test/release evidence. |
| Market Activity surface demotion (`IMPL-18-UTILITY-003`) | Discovery presentation | `SELF_DESCRIBING_UI` | Market Activity Seeds remain visible only as a labeled supporting-input drill-down; no user action or interpretation contract is added. |
| Legacy evidence-route retirement (`IMPL-18-UTILITY-004`) | Navigation consolidation | `SELF_DESCRIBING_UI` | Legacy market-wide evidence resolves to Market Intelligence; ticker evidence resolves to Research subviews. Navigation/browser proof is authoritative. |
| Role-aware Documentation (`IMPL-18-DOC-001`) | Documentation audience and direct-path access | `MANIFEST_AND_POLICY` | This manifest records the audience contract. UI hides developer machinery from USER/DEMO and server direct-path enforcement independently denies it. |
| External Dependency & Provider Readiness (`IMPL-17-DEPS-001`) | Owner/release-engineering setup, entitlement, rights and deployment readiness | `MANIFEST_AND_REGISTRY` | Canonical dependency/readiness and durable User Action Required registries record owner, capability, blocker, user action, rights/entitlement and evidence. CI binds the contract to existing G0–G16 gates only. |

## Direct-path and composition proof required before G10

- Anonymous IdentityService-mode `/docs/*` request: `401`.
- USER/DEMO `/docs/user.md` and `/docs/limitations.md`: allowed.
- USER/DEMO `/docs/developer.md`: `403` and developer tab absent.
- ADMIN/OWNER/SUPER_OWNER `/docs/developer.md`: allowed and tab present.
- Programmatic/UI attempt to select `developer` while unauthorized normalizes to user documentation and never triggers a privileged document fetch.
- Server checks reuse canonical session resolution and role hierarchy; renderer visibility is never accepted as authorization evidence.

## Dependency/readiness proof required before G10

- `dependency_readiness_gate.py` passes in CI Fast and CI Qualified.
- Provider operational entitlement remains separate from contractual data rights.
- All provider-rights evidence remains fail closed until explicitly bound.
- Desktop SQLite and conditional hosted PostgreSQL remain the only canonical persistence architecture.
- Durable user actions contain no secret values and bind only G0–G16.
- CI/release runtime-policy versions have no source-of-truth drift.

## Manifest maintenance contract

Each later v18.6 material slice must append or update a row before its G10 evidence is accepted. `AUDIT-18-AI-001` and `AUDIT-18-AI-RIGHTS-001` therefore receive their final documentation dispositions in this manifest as those slices are implemented. G10 must fail review if a material v18.6 change has no disposition.
