# DEC-2026-08-17-001 — Adaptive CI Cost & Qualification Contract

**Status:** APPROVED / PERMANENT  
**Date:** 2026-08-17  
**Affects:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, G0–G16 execution, GitHub Actions, qualification checkpoints, native certification, AIPLC and release evidence

## Decision

Adopt `governance/ADAPTIVE-CI-QUALIFICATION-CONTRACT.md` as a permanent DE.PULSE operating contract.

DE.PULSE will optimize **when required tests execute**, never **which required tests exist**.

Permanent execution model:

**DEVELOP → ACCUMULATE COHERENT FIXES → DEVELOPMENT CHECKPOINT → G0–G5 FAST QUALIFICATION → BATCH FIX/RECHECKPOINT IF NEEDED → G6–G10 QUALIFICATION → IMMUTABLE RC → G11–G15 FULL/NATIVE CERTIFICATION → G16**

Ordinary development commits must not automatically launch medium/full/browser-intensive/database-intensive/performance/native/release certification. Source and regression tests are implemented during development, then qualified at a deliberate exact-SHA checkpoint.

G0–G5 must run cheap/fail-fast first. Only after FAST is green may G6–G9 run, preferably in bounded parallel when independent. G10 joins the exact checkpoint evidence. G11–G15 remain exhaustive and run only against an immutable RC.

Failed qualification is repaired as a coherent batch and requalified at a new checkpoint. Re-run only failed jobs when the source SHA, workflow/test definition, toolchain/environment and inputs remain equivalent. CI-generated evidence must not recursively trigger the same qualification pipeline.

Quality invariant:

> **CI efficiency may eliminate redundant execution, never required evidence. No cost optimization may weaken scope, behavioral test depth, security, data/rights validation, performance/capacity/stability testing, professional acceptance, supported-platform certification, provenance, or actual packaged-runtime proof.**

Any material source/tooling change after checkpoint qualification or RC invalidates affected evidence and requires requalification of the new fingerprint.

## Immediate v18.5.1 application

`v18.5.1-development` uses `release/v18.5.1/QUALIFICATION_CHECKPOINT.json` as the deliberate qualification trigger. Normal source commits are quiet. The qualification workflow is read-only/non-mutating, exact-SHA bound, cheap-first, and retains the complete G0–G10 evidence chain. G11–G15 remain a later immutable-RC release phase.

## Non-goals

- no G17+;
- no removal of G0–G16 responsibilities;
- no reduction in macOS Apple Silicon or Windows x64 certification requirements;
- no evidence reuse across materially different fingerprints;
- no documentation-only closure of product defects;
- no use of cost as justification for false PASS.

## Canonical contract

`governance/ADAPTIVE-CI-QUALIFICATION-CONTRACT.md`
