# DE.PULSE Claude Entry Instructions

Claude must use the same GitHub-backed contract as every other assistant. Do not infer project state from Claude memory, a prior Claude Project, chat history or an uploaded summary.

Before planning, editing or claiming status:

1. Read `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`.
2. Read `AGENTS.md`; its repository constraints apply equally to Claude.
3. Read `governance/README.md` and `handoff/CURRENT.md`.
4. Inspect the actual GitHub default branch, active branch, open PR, current HEAD, latest Stable tag/release, checks and artifacts.
5. Read `release_identity.json`, `.depulse-certification/resume/build-checkpoint.json` and `.depulse-certification/resume/release-evidence-checkpoint.json`.
6. Read the canonical Roadmap plus the Adaptive Build Plan, Build Process and Delivery Process.
7. Run or inspect `python3 adaptive_resume_gate.py`.
8. Reconcile GitHub/checkpoint disagreements and resume the smallest safe step from the last trustworthy PASS.

GitHub evidence outranks this file and all assistant prose. Use G0–G16 only; preserve the U.S. Equities Processing Boundary, No Execution Boundary, deterministic Market Mode ownership, secrets policy and governed promotion lifecycles. Never ask the user to recreate context already present in the repository.

Before ending meaningful work, commit durable changes, refresh `handoff/CURRENT.md`, update the machine checkpoint after the candidate commit, and leave exactly one next action for any future assistant/account.

