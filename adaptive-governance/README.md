# DE.PULSE Adaptive Governance — Canonical Baseline

**Status:** AUTHORITATIVE DOCUMENT INDEX  
**Rebaselined:** 2026-08-28  
**Decision:** DEC-2026-08-28-001

This index prevents the repository from accumulating multiple documents that all appear current. Code/runtime/package evidence defines implementation truth; the files below define approved intent and execution rules.

## Canonical narrative authorities

| Concern | Single authority |
|---|---|
| Permanent adaptive operating rules | `governance/ADAPTIVE-OPERATING-CONTRACT.md` |
| Product sequencing and time horizons | `governance/ROADMAP.md` |
| Build planning | `adaptive-governance/ADAPTIVE_BUILD_PLAN.md` |
| Build execution and qualification | `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md` |
| Packaging, delivery and release | `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md` |
| Exact current machine status | `governance/current-state.json` + active closure ledger |
| One current human resume | `handoff/CURRENT.md` |

No other file may independently redefine those concerns.

## Audit-rebaseline evidence

- `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md` — approved audit synthesis and execution rebaseline;
- `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md` — full section/appendix conservation;
- `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json` — all Executive findings and audit-wide risks;
- `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json` — eleven maturity-domain targets;
- `governance/programs/V19-V20-REBASELINE/product-audit-section-coverage.json` — machine coverage;
- `governance/programs/V19-V20-REBASELINE/product-audit-third-pass-delta.json` — third-pass residual check;
- `governance/V19_V20_REBASELINE.md` and machine maps — version conservation;
- `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md` — additive provider-registry placement.

## Current projections

The following paths remain for compatibility with resume/current-state gates and historical references. They project Stable/work-slice/issue/branch/ledger identity only and never own scope, rules or next action:

- `CURRENT_ADAPTIVE_ROADMAP.md`;
- `CURRENT_ADAPTIVE_BUILD_PLAN.md`;
- `CURRENT_ADAPTIVE_BUILD_PROCESS.md`;
- `CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`;
- `CURRENT_ADAPTIVE_CI_CONVERGENCE.md`;
- `CURRENT_ADAPTIVE_GAP_CLOSURE.md`.

`ADAPTIVE_ROADMAP.md` is a compatibility pointer to `governance/ROADMAP.md`.

## Permanent specialized contracts

These remain active only within the canonical hierarchy and may not override it:

| File | Disposition |
|---|---|
| `ADAPTIVE_CI_OPERATING_CONTRACT.md` | KEEP; CI specialization |
| `ADAPTIVE_CI_V19_REBASELINE.md` | KEEP as v19+ CI addendum; later merge into CI contract |
| `ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md` | KEEP; provider onboarding specialization |
| `ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md` | KEEP; decomposition specialization |
| `BUILD_RESUME_PROTOCOL.md` | KEEP; portability/resume specialization |
| `FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md` | KEEP; G2/G9 utility checkpoint |
| `GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md` | KEEP; governed-to-delivered traceability |
| `NAMING_AND_IDENTITY_CONTRACT.md` | KEEP; naming specialization |
| `PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md` | KEEP; persistence/session data specialization |
| `ROLE_AWARE_SESSION_SECURITY_CONTRACT.md` | KEEP; identity/session specialization |
| `ROLE_AWARE_UI_COMPOSITION_CONTRACT.md` | KEEP; UI/RBAC specialization |
| `SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md` | NEEDS UPDATE; semantic principle retained, obsolete v18.3–v18.5 placement superseded by roadmap |
| `SQLITE_POSTGRES_SYNC_CONTRACT.md` | NEEDS UPDATE/MERGE; compatibility principles retained, future authority is Postgres v2/outbox/versioned sync |

If a specialized contract conflicts with the audit-rebaselined canonical five, classify and record a new decision before changing behavior. Version-specific text in a permanent specialized contract does not change the canonical roadmap.

## Historical evidence — do not use as current authority

- `LEGACY_TEST_GATE_CLEANUP_PLAN.md`;
- `LEGACY_TEST_GATE_INVENTORY.md`;
- `V18_6_1_CURRENT_RECONCILIATION.md`;
- `V18_6_7_CURRENT_RECONCILIATION.md`;
- `V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`;
- `V18.8.1-ZERO-MISS-RECONCILIATION.json`.

These remain immutable traceability inputs until a later repository-layout/archive change moves them. They do not describe the current product or build.

## Update rule

1. Change approved intent in the applicable canonical authority.
2. Record a material Decision Log entry when scope, authority or sequencing changes.
3. Update machine maps/closure evidence for execution status.
4. Update `handoff/CURRENT.md` with exactly one next action.
5. Keep projection files thin; never paste the canonical narrative into them.
6. Run adaptive resume, current-state projection, Data Health, workflow-policy and requirement-conservation gates.

This rebaseline preserves history; it does not claim planned audit gaps are implemented.
