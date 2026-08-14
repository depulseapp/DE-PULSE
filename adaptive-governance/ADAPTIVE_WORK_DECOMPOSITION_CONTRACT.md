# DE.PULSE — Adaptive Work Decomposition & Gate Evolution Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Adaptive CI Operating Contract  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE engineering work must follow the same adaptive-intelligence principle as the product itself. Large tasks, tests, certifications and delivery activities should not be forced through one monolithic execution unit when smaller independently verifiable units would improve reliability, speed, diagnosis, recovery, reuse or learning.

Permanent execution loop:

**Understand → decompose → map dependencies → reuse → execute → checkpoint → evaluate evidence → adapt next work → integrate → certify → learn.**

## 1. Adaptive decomposition is the default for heavy work

Before a heavy task, gate or release phase begins, the Build Coordinator must decide whether it should remain one unit or be decomposed into smaller work packages, checkpoints, sub-stages, shards or parallel lanes.

Decomposition is preferred when it materially improves one or more of:

- fault isolation and root-cause diagnosis;
- exact-source evidence reuse after a partial failure;
- safe parallel execution of independent work;
- CI/runtime resource efficiency;
- interruption recovery and resumability;
- ownership clarity;
- test evidence quality and traceability;
- prevention of one slow/flaky responsibility blocking unrelated work;
- reuse of already-qualified evidence;
- ability to learn from intermediate outcomes and adapt the next action.

Do not decompose merely for ceremony. A split that adds coordination cost without better evidence, reliability, recovery, speed or clarity should be rejected.

## 2. Dependency-aware execution

Every decomposed work package must have:

- one canonical owner/responsibility;
- explicit inputs and dependencies;
- source/fingerprint/artifact identity where relevant;
- clear PASS/FAIL/BLOCKED or completion criteria;
- evidence location;
- downstream consumers;
- invalidation rules when its inputs change.

Independent packages may execute in parallel. Dependent packages remain ordered. The Build Coordinator owns the overall dependency graph and prevents multiple lanes from mutating the same canonical release state concurrently.

Parallelism must be bounded by runner/provider/runtime capacity and must not create duplicate expensive work, duplicate provider acquisition, conflicting writes or misleading evidence.

## 3. Checkpoint-first recovery

A heavy phase should expose meaningful intermediate checkpoints whenever a later failure would otherwise force unnecessary repetition.

A checkpoint is reusable only when its relevant source/fingerprint/artifact inputs remain unchanged and its dependencies still hold. A checkpoint is not a substitute for a required release gate; it is a smaller evidence boundary inside or between gates.

On failure or interruption:

**inspect actual state → identify the smallest affected work package → preserve unrelated PASS evidence → repair → rerun only affected/dependent work → continue.**

## 4. AI/LLM-style adaptive execution

The build process should not behave like a fixed blind script when evidence shows a better next action. After each meaningful checkpoint, the coordinator may adapt sequencing, parallelism, test focus, consolidation work or remediation based on actual evidence while preserving immutable scope, product contracts and release assurance.

Examples:

- a performance hotspot discovered early may cause profiling/remediation to run before broader certification;
- a platform-specific failure may rerun only that native lane when shared certified source is unchanged;
- a repeated CI harness failure may become a preventative preflight and remove redundant downstream retries;
- a new feature audit may redirect work toward reuse/consolidation rather than adding another tab/engine;
- a large all-tab UI audit may be sharded by role × viewport × surface family and then recombined into one G9 conclusion.

Adaptive execution never means silent production-policy self-modification, arbitrary scope growth, weakened tests or bypassed evidence.

## 5. Gate evolution rule

The current G0–G16 model remains the canonical release-gate map by default, but the gate model itself is allowed to evolve when evidence proves that the existing structure no longer provides a clean, non-duplicative assurance boundary.

Before adding a new G-gate, the proposal must pass a **Gate Utility Test**:

1. **Distinct risk/responsibility** — it protects a materially different release concern that is not already owned cleanly by an existing gate.
2. **Independent evidence** — it has clear inputs, PASS/FAIL criteria and durable evidence.
3. **Non-duplication** — the need cannot be solved more cleanly as a checkpoint, sub-stage, parallel lane or strengthened existing gate.
4. **Material value** — it improves correctness, safety, traceability, recovery or delivery assurance enough to justify additional process complexity.
5. **Canonical ownership** — exactly one gate owns the responsibility after the change; overlapping old checks are consolidated or removed.
6. **Process-wide update** — the Roadmap, Build Plan, Build Process, Delivery Process, CI contract, checkpoint/ledger schema and relevant automation are updated together.
7. **Migration clarity** — in-flight and future releases know which gate map applies; historical Stable evidence is never retroactively rewritten.
8. **G16 review** — the new gate is reviewed after use to confirm it reduced risk/duplication rather than adding ceremony.

No workflow may invent an ad-hoc `G17`, `G18`, etc. in isolation. A new release gate exists only after the canonical gate map is deliberately revised under this rule.

## 6. Preferred hierarchy of process change

When a heavy responsibility needs improvement, prefer in this order:

**reuse existing evidence → split into checkpoints → shard/parallelize independent lanes → strengthen/reassign an existing gate → add a new canonical gate only when materially justified.**

This mirrors DE.PULSE product development:

**reuse → correlate → consolidate → add only when needed.**

## 7. G0–G16 integration

Under the current canonical map:

- **G1/G2/G3:** plan decomposition, dependencies, ownership, evidence and expected parallel lanes before heavy implementation/qualification.
- **G4:** implementation work may be split into independently testable work packages while preserving immutable scope.
- **G5/G6:** use fast/medium qualification as early checkpoints rather than waiting for one giant full-suite result.
- **G7/G8/G9:** shard data/security/adaptive, performance/capacity and cross-module/UI work where independence is real, then aggregate to one gate conclusion.
- **G10:** verify every required sub-checkpoint/lane for the candidate is complete before freeze.
- **G11/G12:** immutable RC identity binds downstream certification evidence; full certification may internally decompose but produces one authoritative G12 conclusion.
- **G13/G14:** package/runtime work is naturally platform-separated; macOS and Windows retain independent artifact evidence.
- **G15:** promotion consumes the complete dependency graph rather than assuming one monolithic prior task.
- **G16:** audit whether decomposition/parallelism improved cycle time and reliability, whether any packages were redundant, and whether gate evolution is warranted.

## 8. Performance and load discipline

Breaking a heavy task into many jobs must not accidentally create a heavier system. The coordinator must consider total CPU, memory, runner minutes, provider/API calls, browser instances, artifact size, storage, database load and external rate limits.

Prefer shared setup/artifacts, canonical data acquisition, cached dependencies, reusable evidence and bounded parallelism over duplicate work in multiple shards.

## 9. Permanent rule

**DE.PULSE should solve large engineering problems the way an adaptive intelligence system reasons: break complexity into meaningful parts, understand dependencies, reuse what is already known, execute independent work efficiently, inspect evidence, adapt the next step, integrate the result, and learn. The release process may evolve—including its gate structure—when measured evidence shows a simpler, safer or more reliable model, but it must never grow process machinery without proving utility.**
