# DE.PULSE — AI Assistant & Account Portability Contract

**Status:** PERMANENT / GOVERNING  
**Applies to:** ChatGPT, Codex, Claude, GitHub Copilot, future AI assistants, human developers, and all user accounts  
**Authority:** GitHub repository `depulseapp/DE-PULSE`

## Purpose

DE.PULSE must never depend on one AI assistant, one conversation, one ChatGPT account, one Claude account, one temporary workspace, or one person’s memory. A newly authorized assistant must be able to reconstruct the last trustworthy project state and resume the smallest correct next step from GitHub alone.

Conversation memory and external handoff copies are convenience only. They may help locate GitHub, but they are never required release evidence or approved-intent authority.

This contract strengthens the existing G0/G2/G10/G16 resume controls. It does not create a gate beyond G0–G16.

## GitHub source-of-truth hierarchy

On every fresh account, assistant change, model change, handoff, or interrupted build, use this order:

1. immutable Stable tag and GitHub Release, including certified artifacts and hashes;
2. actual active development/release branch and current Git commit SHA;
3. open pull request state, diff, checks, reviews and retained artifacts;
4. `release_identity.json` and `VERSION.txt`;
5. `.depulse-certification/resume/build-checkpoint.json` and `release-evidence-checkpoint.json`, reconciled against actual GitHub state;
6. `governance/README.md`, approved scope, permanent contracts, canonical roadmap and decision log;
7. `handoff/CURRENT.md`, which is the one current human-readable continuation report;
8. release-specific evidence under `release/<version>/`;
9. chat memory, uploaded handoffs, local notes or temporary workspaces as advisory context only.

No lower source may override contradictory higher evidence. A documentation claim never proves implementation or package delivery.

## Repository entry points

- `AGENTS.md` is the concise entry adapter for ChatGPT/Codex and other agents that load repository agent instructions.
- `CLAUDE.md` is the concise entry adapter for Claude.
- Both adapters point to this vendor-neutral contract and the same canonical files. They may not create different product rules.
- `handoff/CURRENT.md` is the single current human continuation record and supersedes earlier handoff narratives.
- The Build State Ledger is the current machine-readable state. It must be reconciled with the actual branch and PR before trust.

An assistant that cannot access the private GitHub repository must stop and request repository connection/access. It must not reconstruct project truth from model memory.

## Mandatory fresh-session algorithm

Before proposing or changing work, every new assistant/account must:

1. confirm repository `depulseapp/DE-PULSE` is accessible;
2. read the applicable root adapter (`AGENTS.md` or `CLAUDE.md`);
3. read this contract and `governance/README.md`;
4. read `handoff/CURRENT.md`;
5. inspect the actual default branch, active branch, open PR, current HEAD, latest Stable tag/release and available checks/artifacts;
6. read `release_identity.json`, the Build State Ledger and release-evidence checkpoint;
7. read the canonical Roadmap, Adaptive Build Plan, Adaptive Build Process and Adaptive Delivery Process;
8. run or inspect `tools/ci/adaptive_resume_gate.py` and classify any disagreement;
9. determine the last trustworthy PASS and earliest required resume gate;
10. continue exactly one smallest safe next action, without asking the user to restate repository-resident context.

If GitHub and the checkpoint disagree, correct the checkpoint truthfully. Never change GitHub history or product intent merely to make a stale checkpoint pass.

## Durable handoff rule

After every meaningful build/checkpoint and at G16:

- commit all meaningful source, governance, test and evidence changes to the active branch;
- update the machine-readable checkpoint after the product/tooling/documentation commit so `candidateSourceCommit` points to the exact candidate it describes;
- update `handoff/CURRENT.md` with what changed, what passed, what remains, exactly one next action and a provider-neutral continuation instruction;
- update Roadmap, Build Plan, Build Process and Delivery Process only when their durable state or contract changed;
- preserve older certified history in immutable tags, Releases and release evidence; do not accumulate several files all claiming to be current;
- do not leave essential evidence only in ChatGPT Library, Claude Projects, a local Mac, an AI sandbox, email or chat.

Metadata-only checkpoint commits may follow the candidate commit under the existing fingerprint-exclusion rule. A new assistant must verify intervening commits are truly metadata-only before reusing evidence.

## Assistant-neutral behavior

Every assistant must:

- use the canonical G0–G16 model; never invent G17+;
- follow LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE;
- prefer REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD;
- preserve the U.S. Equities Processing Boundary and permanent No Execution Boundary;
- distinguish approved intent, implemented source, executable evidence, actual package evidence and future roadmap placement;
- keep open requirements visible until executable closure evidence exists;
- preserve deterministic/statistical ownership of price truth and numeric Market Mode calculation;
- never treat its own prior prose, confidence or memory as evidence;
- never silently weaken gates because an assistant lacks a tool, a paid CI budget, a credential or prior conversation access.

## Secrets and account independence

Repository continuity files must contain no API keys, passwords, cookies, tokens, private signing materials or account-specific session data. They may record that a credential is required, where it is configured, and the evidence needed after configuration.

Connecting the same GitHub repository from another ChatGPT/Claude account must be sufficient to recover project state. The new account must not require access to the previous account’s chat history or file storage.

## Validation

`tools/ci/adaptive_resume_gate.py` is the owning machine check inside existing G0/G2/G10 qualification. It verifies:

- both assistant adapters exist and point to the same canonical contract;
- the canonical handoff and machine checkpoints exist;
- release/branch/PR identity is consistent across current continuity artifacts;
- CI policy declares GitHub authority and account/assistant independence;
- all four adaptive layers preserve the portability requirement.

The free exact-source runner and normal pre-freeze qualification must execute this check. A missing/stale portability artifact is a blocking continuity failure, not a reason to guess.

## New account quick start

1. Connect GitHub and grant access to `depulseapp/DE-PULSE`.
2. Open the repository’s active PR/branch.
3. Tell the assistant: “Read the repository instructions and resume DE.PULSE from GitHub source of truth. Do not rely on chat memory.”
4. The assistant must execute the mandatory fresh-session algorithm above.

No upload of an old chat handoff is required when GitHub is accessible.

