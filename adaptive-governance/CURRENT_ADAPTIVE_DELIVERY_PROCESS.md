# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed process foundation:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`

TradeInsight Settings/API-key UX (#76 / `ADAPT-TRADEINSIGHT-SETTINGS-001`) earned exact-head Fast #848 and Qualified #181 before PR #77 expected-head merge `a171ce2258632bd4bd6aa737176f2d6dffb44689`. The implementation is complete, but this delivery checkpoint does **not** publish a new Stable version; existing `v18.9.1-stable` evidence remains immutable.

## Delivery rule for the provider-production program

Parent #79 / `ADAPT-PROVIDER-PRODUCTION-001` proceeds through #80 → #81/#82 → #83 + #78 → #84. Each child must deliver executable capability evidence before lifecycle/promotion claims. Provider reachability or a configured key is not production admission.

For each governed merge: reconcile with current `main`; exact-head Fast must pass; deliberate Qualified must pass the Planner-selected evidence owners; merge with expected-head protection; then bind immutable candidate/CI/merge evidence into source-of-truth closure where required. Multi-platform/native release evidence is required when the release planner or Stable publication path selects it, not invented for an unrelated slice.

Public Stable publication remains a separate G11–G16 decision. #80–#84 must preserve rollback/degradation truth, authority boundaries, provider lifecycle evidence, scoped recovery proof, and no-secret-leakage. No release is allowed merely because implementation work has merged.
