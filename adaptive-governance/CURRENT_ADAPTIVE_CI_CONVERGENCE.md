# CURRENT Adaptive CI Convergence

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Canonical retained process closure ledger:** `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Completed foundation:** #80 / Fast #859 / Qualified #182 / PR #86  
**Completed Router adoption:** #81 / Fast #878 / Qualified #184 / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Active product work:** #82 / `ADAPT-DATAHEALTH-RUNTIME-001` / `adapt-datahealth-runtime-001`

#70/#73 remain the completed CI/repository control-plane foundation. #82 uses only the existing three canonical workflows: `.github/workflows/ci-fast.yml`, `.github/workflows/ci-qualified.yml`, and `.github/workflows/release.yml`. No temporary workflow family, redundant permanent gate family, gate weakening, direct-main development, or fabricated qualification evidence is permitted.

The existing Data Health/source-health and runtime regression owners must evolve with #82 so scoped degradation, warm-state reuse, fallback/recovery, hysteresis, backpressure/load shedding and false-global-degradation prevention are executable properties rather than documentation. Reuse existing test/gate owners; do not create a G17+ family.

Final #82 closure requires exact-head Fast and impact-selected Qualified. Release is not invoked merely for Data Health runtime work; #84 owns final zero-gap release-quality closure.
