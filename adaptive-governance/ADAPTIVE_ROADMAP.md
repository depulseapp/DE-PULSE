# DE.PULSE — Adaptive Roadmap

## Permanent Build Recoverability Contract

The Adaptive Roadmap permanently includes **GitHub-backed build recoverability** from v18.2 onward.

Every roadmap release must be resumable after interruption from the last trustworthy GitHub/CI evidence without relying on conversation history. Each release therefore owns an active development/release branch, a durable build checkpoint, fingerprint-bound gate evidence, and immutable release artifacts as they become available.

Roadmap sequencing remains unchanged by an interruption. Recovery determines the last trustworthy PASS, resumes the current approved release at the next required step, completes G0–G16, then proceeds to the next approved roadmap release.

A source change invalidates affected downstream evidence. A conversation/tool interruption by itself does not invalidate unchanged-source evidence and must not force unnecessary reruns.

At G16, the current checkpoint is archived and the next approved roadmap branch/checkpoint is seeded from the exact promoted Stable commit/tag.

## Permanent Adaptive CI learning direction

From v18.2 onward the build system itself follows the DE.PULSE adaptive principle: **observe → classify → diagnose → generalize → prevent → execute → measure → learn**.

Future roadmap releases inherit:

- one logical Build Coordinator inside canonical G0–G16;
- mandatory CI failure classification so harness/infra/no-op/superseded events are not mislabeled as product defects;
- a preventative-learning preflight driven by accumulated release lessons;
- a GitHub-reconciled Build State Ledger / Checkpoint v2;
- idempotent/path-isolated CI to reduce duplicate and self-triggered workflow noise;
- exact-source partial resume rather than unnecessary full reruns;
- mandatory Mac/Windows native user-delivery closure for completed native releases;
- G16 consolidation of obsolete/duplicate workflow machinery while preserving assurance.

The canonical rules are defined in:
- `adaptive-governance/BUILD_RESUME_PROTOCOL.md`
- `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`

They apply to all future roadmap families unless replaced by a stricter approved contract.

## Permanent Role-Aware UI Composition & Capability Governance

From v18.2 onward every roadmap release must preserve a deliberate, role-aware product experience across the complete application shell and every current or newly introduced tab.

The governing principle is **role/capability changes composition, not information hierarchy**. DE.PULSE must not design one OWNER page and merely hide cards for everyone else. Each role receives an intentional composition built from the same canonical information hierarchy, authorized capabilities and reusable responsive primitives.

Permanent roadmap direction:

- **SUPER_OWNER / OWNER** retain full authorized governance and implementation/operational visibility.
- **ADMIN is capability-based, not automatically full-power.** Selected admins may receive deep/full operational rights when explicitly delegated; other admins receive only their assigned administrative capabilities.
- **USER / DEMO receive no implementation machinery** such as provider plumbing, queue/cache/database/scheduler internals, subscription capacity, routing/circuit/fallback internals, secrets or policy/model implementation controls.
- role visibility is never security; matching backend/API authorization and server-side sensitive-data suppression are mandatory.
- removing unauthorized content must trigger natural layout recomposition: no blank grid cells, orphan headings, dead space, placeholder geometry, hierarchy inversion, clipping or awkward page endings.
- privileged diagnostics/administration remain contextually placed; they must not move above higher-value market intelligence merely because the viewer is privileged.
- a new tab is allowed only when it materially improves utility, clarity, security or workflow separation. It must not be created merely to relocate content. A dedicated **Administration** tab is permitted/expected when it cleanly serves delegated administrative capabilities; its navigation and contents remain capability-scoped.
- every new tab or major surface automatically enters the role/capability audit contract.
- protected deterministic Day/Swing/Long logic and the No Execution Boundary are role-invariant.

This roadmap contract is enforced inside existing G0–G16, especially G1/G2/G3/G9/G10/G12/G16. It creates no new gate.

Canonical role-aware rules: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`.
