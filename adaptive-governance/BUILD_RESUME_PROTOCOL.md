# DE.PULSE — Permanent Build Resume Protocol

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process  
Effective: v18.2 and all later releases

## Purpose

DE.PULSE builds must be recoverable after chat interruption, tool/session loss, CI interruption, runner timeout, or handoff to a new conversation without depending on chat history or repeating already-proven work.

The protocol is an efficiency and recoverability mechanism only. It does **not** add a gate beyond G0–G16 and never permits weakening, bypassing, inventing, or reusing invalid release evidence.

## Authoritative recovery hierarchy

On every resume, use this order of authority:

1. immutable incoming Stable tag/release and its provenance;
2. active development/release branch and actual Git commit SHA;
3. committed build checkpoint under `.depulse-certification/resume/`;
4. GitHub Actions run state, CI evidence, retained artifacts and immutable RC/native hashes;
5. chat/handoff narrative as advisory context only.

A checkpoint is never trusted by itself. It must be reconciled with the branch SHA, canonical release identity, source fingerprint where applicable, CI runs and retained artifacts.

## Mandatory resume algorithm

At the start of a continuation or after an interruption:

1. Detect the active release branch.
2. Verify the immutable incoming Stable tag/commit/release.
3. Read `.depulse-certification/resume/build-checkpoint.json`.
4. Read the actual active-branch HEAD and canonical `release_identity.json`.
5. Inspect the latest relevant CI runs and retained G10/G11/G12/G14/G15 artifacts.
6. Compare the checkpoint source identity/fingerprint with the evidence source identity/fingerprint.
7. Determine the **last trustworthy PASS**, not merely the last reported PASS.
8. Resume at the next required step or earliest invalidated gate.
9. Update the checkpoint after meaningful state transitions.

## Fingerprint and evidence-reuse rule

PASS evidence may be reused only when all evidence-relevant source identity is unchanged and the evidence is bound to the same source fingerprint or immutable RC/artifact identity.

If product source changes, downstream evidence is invalidated from the earliest affected gate onward. The build resumes from that gate; it does not restart earlier unrelated work and it does not retain later PASS labels from the superseded fingerprint.

Checkpoint-only metadata commits under `.depulse-certification/` are intentionally excluded from the canonical product-source fingerprint. They may advance branch HEAD without invalidating product evidence, provided every intervening commit changes only fingerprint-excluded checkpoint/evidence metadata.

Release-channel transformations such as TEST → STABLE require their own prescribed Stable identity/certification evidence. TEST certification is an input, not a substitute for Stable certification.

## Mandatory checkpoint cadence

Update durable recovery state at least after:

- G0–G3 scope/architecture/readiness freeze;
- meaningful code commits that change the candidate fingerprint;
- G4 Development Exit;
- G10 Pre-Freeze Qualification;
- G11 Immutable RC freeze;
- G12 Full Certification;
- each native G14 platform result;
- G15 Release Assurance / promotion decision;
- G16 retrospective, archive, cleanup and next-release seed.

No meaningful uncommitted local work is considered durable or resumable. Before leaving a material implementation phase, persist the work to the active GitHub branch and bind test/evidence state to that commit/fingerprint.

## Two-level checkpoint model

### Build checkpoint

`.depulse-certification/resume/build-checkpoint.json` records current release/branch/baseline, source commit/fingerprint state, gate states, last trustworthy PASS, next required step, blocking issue and update time.

### Release evidence checkpoint

`.depulse-certification/resume/release-evidence-checkpoint.json` records immutable evidence as it becomes available: G10 fingerprint, G11 RC/source SHA, G12 evidence, native macOS/Windows package/runtime evidence, G15 assurance, Stable tag/release identity and G16 archive state.

The evidence checkpoint may reference GitHub run/artifact IDs, hashes and immutable release identities; it must never claim an artifact PASS that does not exist.

## Partial native/certification resume

If an interruption occurs after one platform or certification lane passed, exact-source evidence may be reused. Example: G12 PASS + macOS G14 PASS + Windows G14 pending resumes at Windows G14/G15 when the immutable source/artifact identity is unchanged. It must not rerun unrelated expensive lanes merely because the conversation stopped.

If source or packaging identity changed, affected native evidence is invalidated and rerun.

## G16 closure and next build

After Stable promotion:

1. bind `main`, immutable Stable tag and permanent release to the certified source/artifacts;
2. archive the final checkpoint/evidence state;
3. clean obsolete working branches/artifacts according to repository hygiene rules;
4. create the next approved development branch from the exact Stable commit/tag;
5. seed a new checkpoint for the next release at G0.

## User-interruption policy

Normal recovery must require no routine manual work from the user. Ask only for an unavoidable external credential/account authorization, paid-service decision, legal/data-rights decision, or genuinely new material product decision that cannot be resolved from approved DE.PULSE contracts.

## Permanent safety rule

**Resume means continue from the last trustworthy evidence, not from the last optimistic status.** GitHub branch state, immutable source identity and CI/artifact evidence control the decision.