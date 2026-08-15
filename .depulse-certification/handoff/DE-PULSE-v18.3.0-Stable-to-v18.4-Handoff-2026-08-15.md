# DE.PULSE — v18.3.0 Stable → v18.4 Handoff

**Closure date:** 2026-08-15  
**Certified release:** `v18.3.0` STABLE  
**Stable tag:** `v18.3.0-stable`  
**Certified Stable source:** `adc3adc88e92b53071f0fc0f23d00ce34f091e42`  
**Stable source fingerprint:** `535d6c854874201501e62330f26683cdb252a97935fdc17c9816498146c71c5f`  
**Stable source archive SHA-256:** `c70ffba7eaf97cc33762db4016714497932c7f6802c3776d0d2a92897b4194b1`  
**Build ID:** `v18.3.0-stable-postgresql-hosted-shared-state-20260815`

## 1. Release truth

v18.3.0 is fully certified Stable. `main` was fast-forwarded to the exact certified Stable source with no merge rewrite or force, then the annotated `v18.3.0-stable` tag was created on that exact source. Only after the tag existed did G16 fingerprint-excluded metadata begin advancing `main`.

The Stable identity promotion from the TEST RC changed exactly 15 identity/documentation/gate files. PostgreSQL repository/runtime implementation, provider/scanner/intelligence owners, and protected deterministic Day/Swing/Long formula owners were unchanged by Stable promotion.

## 2. G0–G16 evidence summary

| Evidence | Result | Authority |
|---|---|---|
| G0–G8 final TEST candidate | PASS | run `31912832420` |
| G9–G10 final TEST candidate | PASS | run `31912832422`; TEST fingerprint `1b2251231d4d4a140ad034bc022d62c5ffd68f5c1ff0a2b4ed926089bc025de4` |
| TEST G11–G15 | PASS | run `31913540543`; `STABLE_PROMOTION_READY` |
| Stable source identity promotion | PASS | `adc3adc88e92b53071f0fc0f23d00ce34f091e42` |
| Stable G11 | PASS | run `31914647719`; immutable source artifact `9254586144` |
| Stable G12 | PASS | run `31914647719`; full race + PostgreSQL 17 + hosted readiness + UI/runtime; artifact `9254667117` |
| Stable G13/G14 macOS Apple Silicon | PASS | retry run `31915196766`; artifact `9254722138`; package SHA-256 `3c2ff8332e4579a2a8165a188e3e934ec17e48de90b37465c9fbdf2d8e3ff8be` |
| Stable G13/G14 Windows x64 | PASS | retry run `31915196766`; artifact `9254726071`; package SHA-256 `7539443a71cb156d8f9a73a40de827d4e4274d59e66d021764f0b9d948963675` |
| Stable G15 | PASS | retry run `31915196766`; artifact `9254728203`; `STABLE_CERTIFIED` |
| Stable tag | PASS | run `31915319399`; tag resolves to `adc3adc88e92b53071f0fc0f23d00ce34f091e42` |
| v18.3 temporary branch cleanup | PASS | run `31915391884`; zero `v18.3*` working branches remain |
| G16 | PASS | this handoff + resume checkpoints |

## 3. What v18.3 delivered

- PostgreSQL 17 repository parity beneath the canonical `PersistenceBackend`.
- Desktop macOS/Windows remain canonical SQLite/local by default; hosted PostgreSQL selection is explicit and fail-closed.
- Dedicated-session migration locking, bounded pool/contention behavior and database diagnostics.
- Integrity-protected persistence archive, backup/restore and SQLite→PostgreSQL migration continuity.
- Database-aware readiness with bounded outage/recovery hysteresis.
- Hosted liveness/readiness and actual hosted runtime certification.
- Per-user workspace isolation while retaining one shared deduplicated market/provider/scanner/intelligence universe.
- Source-responsibility hardening: hosted health/readiness HTTP ownership separated from the oversized API owner without behavior change.
- Required macOS Apple Silicon and Windows x64 actual packaged-runtime certification.
- Protected deterministic Day/Swing/Long formulas and the permanent No Execution Boundary remained unchanged.

## 4. G16 Adaptive Process Learning Closure

Permanent/reusable lessons from v18.3:

1. **Regression-harness ownership:** inherited version-specific tests must verify their historical behavior contract, not own the current release title/version. Current release identity belongs to canonical identity/version gates.
2. **Hosted environment scoping:** hosted runtime environment variables belong to the launched hosted process; they must not contaminate neutral PostgreSQL test/build steps.
3. **Fingerprint authority:** canonical fingerprints come from the actual Git source transformation/bytes. An independently reconstructed equivalent is not authoritative.
4. **Markdown-aware formatting:** generic whitespace checks must not reject intentional Markdown hard-break syntax.
5. **Native G14 ownership:** G14 validates actual packaged runtime responsibilities. It must not duplicate Smart Router/Rapid Move active-state assertions already owned by authoritative G12 regression.
6. **Windows audit semantics:** force UTF-8 when reading UTF-8 repository documents on Windows; POSIX permission-bit assertions are not the Windows security contract.
7. **Immutable-source discipline:** tooling-only branches/scripts may evolve, but immutable RC/Stable source may not.
8. **Patch/build staging:** new files must be represented correctly and staging integrity should bind to repository Git source/blob identity.
9. **Checkpoint ordering:** create/freeze the immutable Stable tag first; only then may fingerprint-excluded G16 checkpoints/handoffs advance `main`.
10. **Delta-first requalification:** harness/tooling-only corrections inherit unaffected trustworthy evidence and rerun only changed/affected gates plus final reconciliation.

## 5. Permanent product/process boundaries to carry forward

- G0–G16 remains the only release-gate model; do not create G17+.
- DE.PULSE remains research/intelligence/decision support only: no paper trading, execution, order routing, P&L, portfolio or journal features.
- US equities / approved US ETFs only under the established product contract.
- Preserve Adaptive Intelligence North Star, ASBI behavioral-state intelligence, Data Utility/Correlation rules, canonical fetch-once/calculate-once ownership and evidence reuse.
- Preserve macOS Apple Silicon and Windows x64 as required release targets.
- Preserve shared canonical market/provider/scanner/intelligence computation instead of multiplying market-wide work per user.
- Preserve production logic validation boundaries; adaptive learning must not silently self-modify production decision logic.
- Preserve branch/release hygiene and fingerprint-excluded post-tag G16 metadata discipline.

## 6. v18.4 start rule

**Do not infer or invent v18.4 product scope from this closure.** The authoritative v18.4 roadmap/scope must be recovered and reconciled before G1 is frozen.

Start v18.4 with:

1. G0 exact baseline = `v18.3.0-stable` @ `adc3adc88e92b53071f0fc0f23d00ce34f091e42`, fingerprint `535d6c854874201501e62330f26683cdb252a97935fdc17c9816498146c71c5f`.
2. Create `v18.4-development` from the final v18.3 G16 `main` metadata head so closure evidence is inherited while the certified source/tag remains immutable.
3. Recover the authoritative v18.4 roadmap and outstanding Defect/Improvement Register items.
4. Run Delta AIPLC / architecture-data-utility reconciliation before freezing G1.
5. Freeze v18.4 scope only after the roadmap, dependencies, evidence plan and protected-owner blast radius agree.
6. Continue autonomous execution under the existing authorization through v18.5; routine implementation/testing/certification/promotion should not require intermediate manual approval. Stop only for genuinely unavoidable secrets/credentials, new financial commitments, legal/licensing/data-rights decisions, irreversible/high-impact external actions, or genuinely new material product decisions.

## 7. Resume instruction

If a new chat resumes this project, treat `v18.3.0-stable` as the immutable product baseline and the final G16 `main` head as process/handoff metadata truth. Reconcile actual Git state first, then begin v18.4 G0 and roadmap recovery. Do not reopen v18.3 implementation unless a new defect is proven against the certified Stable source/artifacts.
