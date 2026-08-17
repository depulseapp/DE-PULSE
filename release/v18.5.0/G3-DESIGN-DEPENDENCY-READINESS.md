# DE.PULSE v18.5.0 — G3 Design / Dependency Readiness

Status: PASS — execution plan frozen; later evidence may block release.

v18.5 uses a **measure → isolate → correct → remeasure → certify** closure design. Existing certified v18 machinery is reused first; new code is justified only by a demonstrated closure gap or defect.

## Closure execution order

1. Reconstruct v18.0–v18.4 approved scope/provenance and protected contracts.
2. Run static architecture/source/data-utility/security review and existing regression suites.
3. Run inherited v16/v17 performance, persistence, priority/backpressure, readiness and lineage harnesses as baseline evidence.
4. Add a v18.5 ADR-GDI scenario matrix around the canonical runtime, covering provider, freshness/disagreement, queue/backpressure, database, restart, fan-out, background-work, blast-radius, hysteresis and truthful UX semantics.
5. Profile realistic active-market workload and attribute degradation to provider/external, database, local-runtime, queue/backpressure or evidence-quality causes rather than reporting one unexplained broad state.
6. Correct only proven product/runtime defects; do not redesign mature certified subsystems merely to create new closure code.
7. Re-run focused + full/race/randomized/Extreme-30/HTTP/responsive/native/hosted qualification on the exact candidate source.
8. Freeze immutable RC, perform G11–G15 exact-source release certification, then bounded Stable identity promotion and fresh Stable recertification before G16.

## Dependency rules

- Real PostgreSQL 17 is required for hosted pressure/recovery evidence; desktop SQLite remains default and must retain truthful continuity.
- Provider failure/rate-limit scenarios must use deterministic test doubles/fault injection where possible; certification must not depend on consuming paid quota or deliberately disrupting external services.
- Actual native package audits are mandatory for macOS Apple Silicon and Windows x64.
- Mature ASBI/TDTI/AODR/adaptive-13F work is not a dependency of v18.5 closure; only existing dependency-compatible foundations are validated.
- Legal/provider data-rights determinations are not fabricated by closure tests; readiness remains fail closed where evidence is absent.

Any closure test that reveals local overload delaying or misrepresenting decision-critical live/current evidence becomes implementation scope automatically because it is a frozen G1 release blocker, not optional feature work.
