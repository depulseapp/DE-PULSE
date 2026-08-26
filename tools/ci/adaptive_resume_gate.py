#!/usr/bin/env python3
"""Validate the permanent GitHub-backed DE.PULSE Build Resume Protocol."""
from __future__ import annotations

import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
errors: list[str] = []
CLOSED_STATES = {"COMPLETE", "COMPLETED", "CLOSED", "DELIVERED", "VERIFIED"}
ACTIVE_PRODUCT_STATES = {"ACTIVE", "IN_PROGRESS"}


def need(ok: bool, msg: str) -> None:
    if not ok:
        errors.append(msg)


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception as exc:
        errors.append(f"unreadable JSON {path.relative_to(ROOT)}: {exc}")
        return {}


def semver(value: object) -> tuple[int, int, int] | None:
    match = re.fullmatch(r"(?:v)?(\d+)\.(\d+)\.(\d+)(?:-stable)?", str(value or "").strip())
    return tuple(int(part) for part in match.groups()) if match else None


def clean_version(value: object) -> str:
    raw = str(value or "").strip()
    if raw.startswith("v"):
        raw = raw[1:]
    if raw.endswith("-stable"):
        raw = raw[:-7]
    return raw


def stable_tag_for(version: object) -> str:
    policy = load_json(ROOT / "governance" / "versioning-policy.json")
    cleaned = clean_version(version)
    current = semver(cleaned)
    cutover = semver(policy.get("effectiveAfterProductVersion"))
    need(bool(current and cutover), "versioning policy / Stable version is invalid")
    if not current or not cutover:
        return f"v{cleaned}"
    pattern = policy.get("legacyStableTagPattern") if current <= cutover else policy.get("futureStableTagPattern")
    return str(pattern).format(productVersion=cleaned)


required_docs = [
    ROOT / "AGENTS.md",
    ROOT / "CLAUDE.md",
    ROOT / "governance" / "AI-ASSISTANT-PORTABILITY-CONTRACT.md",
    ROOT / "handoff" / "CURRENT.md",
    ROOT / "adaptive-governance" / "BUILD_RESUME_PROTOCOL.md",
    ROOT / "adaptive-governance" / "ADAPTIVE_ROADMAP.md",
    ROOT / "adaptive-governance" / "ADAPTIVE_BUILD_PLAN.md",
    ROOT / "adaptive-governance" / "ADAPTIVE_BUILD_PROCESS.md",
    ROOT / "adaptive-governance" / "ADAPTIVE_DELIVERY_PROCESS.md",
]
for path in required_docs:
    need(path.exists(), f"missing permanent resume governance document: {path.relative_to(ROOT)}")

handoff = ""
if all(path.exists() for path in required_docs):
    corpus = "\n".join(path.read_text(errors="ignore") for path in required_docs)
    for term in (
        "last trustworthy PASS", "source fingerprint", "G0–G16", "GitHub", "G16",
        "No Execution", "handoff/CURRENT.md", "AI-ASSISTANT-PORTABILITY-CONTRACT.md", "Claude",
    ):
        need(term in corpus, f"permanent resume governance missing required contract term: {term}")
    contract = (ROOT / "governance" / "AI-ASSISTANT-PORTABILITY-CONTRACT.md").read_text(errors="ignore")
    agents = (ROOT / "AGENTS.md").read_text(errors="ignore")
    claude = (ROOT / "CLAUDE.md").read_text(errors="ignore")
    handoff = (ROOT / "handoff" / "CURRENT.md").read_text(errors="ignore")
    for adapter, text in (("AGENTS.md", agents), ("CLAUDE.md", claude)):
        need("governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md" in text, f"{adapter} must point to the vendor-neutral portability contract")
        need("handoff/CURRENT.md" in text, f"{adapter} must point to the current authoritative handoff")
        need("tools/ci/adaptive_resume_gate.py" in text, f"{adapter} must require the canonical resume gate")
    for term in (
        "GitHub source-of-truth hierarchy", "Mandatory fresh-session algorithm", "Durable handoff rule",
        "Secrets and account independence", "No upload of an old chat handoff is required",
    ):
        need(term in contract, f"portability contract missing required section/term: {term}")
    need("SUPERSEDES ALL PRIOR CHAT HANDOFFS" in handoff, "handoff/CURRENT.md must be the single current handoff authority")
    need("Exactly one next action" in handoff, "handoff/CURRENT.md must name exactly one next action")

identity_path = ROOT / "release_identity.json"
state_path = ROOT / "governance" / "current-state.json"
checkpoint_path = ROOT / ".depulse-certification" / "resume" / "build-checkpoint.json"
evidence_path = ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json"
ident = load_json(identity_path)
state = load_json(state_path)
identity_release = clean_version(ident.get("version"))
identity_build = str(ident.get("build_id", ""))
identity_previous_stable = clean_version(ident.get("previous_stable"))
identity_stable_baseline = clean_version(ident.get("stable_baseline"))
stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
gate = state.get("productCapabilityGate", {}) if isinstance(state.get("productCapabilityGate"), dict) else {}
published_stable_version = clean_version(stable.get("productVersion"))
product_work_rel = str(gate.get("workSlicePath", "")).strip()
product_work = load_json(ROOT / product_work_rel) if product_work_rel and (ROOT / product_work_rel).is_file() else {}
product_status = str(gate.get("reservationStatus", "")).strip().upper()
product_is_release_closure = product_work.get("type") == "PRODUCT_RELEASE_CLOSURE"
release_closure = product_is_release_closure and product_status in ACTIVE_PRODUCT_STATES
completed_release_closure = product_is_release_closure and product_status in CLOSED_STATES

need(ident.get("channel") == "STABLE", "release identity channel must remain STABLE")
if release_closure:
    need(identity_previous_stable == published_stable_version, "release-closure previous_stable must equal published Stable")
    need(identity_stable_baseline == published_stable_version, "release-closure stable_baseline must equal published Stable")
    current_semver = semver(identity_release)
    previous_semver = semver(published_stable_version)
    need(bool(current_semver and previous_semver and current_semver > previous_semver), "release-closure candidate SemVer must be newer than published Stable")
    need(stable.get("tag") == stable_tag_for(published_stable_version), "release-closure published Stable tag drift")
    need(product_work.get("workSliceId") == gate.get("reservedWorkSliceId"), "release-closure reserved workSliceId drift")
    need(product_work.get("issue") == gate.get("reservedIssue"), "release-closure reserved issue drift")
    need(product_work.get("branch") == gate.get("reservedBranch"), "release-closure reserved branch drift")
    need(clean_version(product_work.get("publicProductVersion")) == identity_release, "release-closure publicProductVersion / candidate identity drift")
    need(clean_version(product_work.get("stableProductVersionAtStart")) == published_stable_version, "release-closure Stable-at-start drift")
    need(product_work.get("baselineCandidateSha") == stable.get("candidateSha"), "release-closure baseline candidate / published Stable drift")
    need(product_work.get("baselineSourceFingerprint") == stable.get("sourceFingerprint"), "release-closure baseline fingerprint / published Stable drift")
    need(product_work.get("baselineBuildId") == stable.get("buildId"), "release-closure baseline build / published Stable drift")
    need(product_work.get("targetStable") == stable_tag_for(identity_release), "release-closure target Stable drift")
    need(product_work.get("productBehaviorChange") is True, "release-closure must declare productBehaviorChange=true")
    need(product_work.get("blocksNextProductCapability") is True, "release-closure must block subsequent product capability work")
else:
    need(published_stable_version == identity_release, "current-state Stable productVersion / release identity drift")
    need(stable.get("buildId") == identity_build, "current-state Stable buildId / release identity drift")
    need(str(stable.get("platformBuildNumber", "")) == str(ident.get("bundle_version", "")), "current-state platform build number / release identity drift")
    need(stable.get("tag") == stable_tag_for(identity_release), "current-state Stable tag / versioning-policy drift")

if completed_release_closure:
    need(str(product_work.get("status", "")).upper() in CLOSED_STATES, "completed release-closure work-slice status drift")
    need(clean_version(product_work.get("publicProductVersion")) == identity_release, "completed release closure publicProductVersion / Stable drift")
    need(product_work.get("targetStable") == stable_tag_for(identity_release), "completed release closure targetStable / Stable tag drift")
    need(gate.get("postClosureSourceOverlapAuditRequired") is True, "completed final-v18 closure must retain post-closure source-overlap audit requirement")
    need(str(gate.get("postClosureSourceOverlapAuditStatus", "")).upper() in {"PENDING", "PASS"}, "post-closure source-overlap audit status invalid")

need(isinstance(stable.get("candidateSha"), str) and re.fullmatch(r"[0-9a-f]{40}", stable.get("candidateSha", "")), "current-state Stable candidateSha invalid")
need(isinstance(stable.get("sourceFingerprint"), str) and re.fullmatch(r"[0-9a-f]{64}", stable.get("sourceFingerprint", "")), "current-state Stable sourceFingerprint invalid")
need(stable.get("publication") == "PASS_NO_REBUILD", "current-state Stable publication must remain PASS_NO_REBUILD")

# Retained process authority may predate the current Stable after a later product
# release has completed. Validate its own immutable metadata instead of falsely
# requiring it to have been rebased onto the newer Stable.
work_slice_id = str(active.get("workSliceId", "")).strip()
active_status = str(active.get("status", "")).strip().upper()
need(bool(work_slice_id), "current-state retained work slice missing")
work_slice_path = ROOT / "governance" / "work-slices" / work_slice_id / "work-slice.json"
work_slice = load_json(work_slice_path) if work_slice_path.is_file() else {}
need(bool(work_slice), "retained work-slice metadata missing")
if work_slice:
    need(work_slice.get("workSliceId") == work_slice_id, "current-state/work-slice ID mismatch")
    need(work_slice.get("issue") == active.get("issue"), "current-state/work-slice issue mismatch")
    need(work_slice.get("branch") == active.get("branch"), "current-state/work-slice branch mismatch")
    need(str(work_slice.get("status", "")).strip().upper() == active_status, "current-state/work-slice status mismatch")

    if active.get("type") == "PROCESS_RELEASE_ENGINEERING":
        need(active.get("publicProductVersion") is None, "process work slice must not consume a public product version")
        need(active.get("productBehaviorChange") is False, "process work slice must remain product-behavior neutral")
        need(work_slice.get("publicProductVersion") is None, "registered process work slice consumed a public product version")
        need(work_slice.get("productBehaviorChange") is False, "registered process work slice must remain product-behavior neutral")
        need(active_status in CLOSED_STATES, "retained process authority must be complete once final-v18 release is complete")
        need(work_slice.get("blocksNextProductCapability") is False, "retained completed process work must remain unblocked")
        final_evidence = str(work_slice.get("finalQualificationEvidence", "")).strip()
        need(bool(final_evidence) and (ROOT / final_evidence).is_file(), "retained completed process work lacks final qualification evidence")
        need(active.get("finalQualificationEvidence") == final_evidence, "retained process final qualification evidence drift")
        need(active.get("mergedCommitSha") == work_slice.get("mergedCommitSha"), "retained process merged commit drift")
        if not completed_release_closure:
            need(work_slice.get("baselineCandidateSha") == stable.get("candidateSha"), "process work-slice baseline candidate / Stable drift")
            need(work_slice.get("baselineSourceFingerprint") == stable.get("sourceFingerprint"), "process work-slice baseline fingerprint / Stable drift")
            need(work_slice.get("baselineBuildId") == stable.get("buildId"), "process work-slice baseline build / Stable drift")
        else:
            need(gate.get("unblockedByCompletedWorkSlice") == work_slice_id, "completed release must preserve retained process unblock authority")
    elif active.get("type") == "PRODUCT_RELEASE_CLOSURE":
        need(active_status in CLOSED_STATES, "retained product release closure must be complete")
        need(clean_version(active.get("publicProductVersion")) == identity_release, "retained product release closure public version drift")
    else:
        need(False, f"unsupported retained work-slice type: {active.get('type')!r}")

workflow_dir = ROOT / ".github" / "workflows"
routine = sorted(path.name for path in workflow_dir.glob("*.yml"))
need(routine == ["ci-fast.yml", "ci-qualified.yml", "release.yml"], f"routine workflow set drift: {routine}")
prefreeze_text = (ROOT / "tools" / "release" / "prefreeze_qualification.py").read_text(errors="ignore")
need("'.depulse-certification'" in prefreeze_text or '".depulse-certification"' in prefreeze_text, "prefreeze fingerprint must exclude .depulse-certification metadata")

stable_release = ""
cp = load_json(checkpoint_path)
if cp:
    schema = cp.get("schemaVersion")
    need(isinstance(schema, int) and schema >= 2, "build checkpoint schemaVersion must be v2 or later")
    stable_release = clean_version(cp.get("release"))
    need(bool(stable_release), "build checkpoint Stable release missing")
    need(cp.get("channel") == "STABLE", "build checkpoint must describe certified Stable evidence")
    need(cp.get("branch"), "checkpoint branch missing")
    certified = cp.get("certifiedStable", {}) if isinstance(cp.get("certifiedStable"), dict) else {}
    need(clean_version(certified.get("version")) == stable_release, "checkpoint release must match certifiedStable version")
    need(certified.get("tag") == stable_tag_for(stable_release), "checkpoint certified Stable tag / versioning-policy drift")
    expected_checkpoint = published_stable_version
    need(stable_release == expected_checkpoint, f"build checkpoint must align to published Stable v{expected_checkpoint}")
    if identity_release != stable_release:
        need(identity_previous_stable == stable_release, "candidate previous_stable must match immutable predecessor checkpoint")
        need(identity_stable_baseline == stable_release, "candidate stable_baseline must match immutable predecessor checkpoint")
    else:
        need(bool(identity_previous_stable), "promoted Stable previous_stable must identify its predecessor")
        need(bool(identity_stable_baseline), "promoted Stable stable_baseline must identify its certified predecessor baseline")
        need(identity_previous_stable == identity_stable_baseline, "promoted Stable previous_stable/stable_baseline mismatch")
    candidate = cp.get("candidateSourceCommit") or certified.get("candidateSourceCommit") or certified.get("certifiedSourceCheckout")
    need(isinstance(candidate, str) and re.fullmatch(r"[0-9a-f]{40}", candidate), "checkpoint candidate source commit must be a Git SHA")
    need(candidate == stable.get("candidateSha"), "checkpoint candidate / current Stable candidate drift")
    fp = cp.get("sourceFingerprint") or certified.get("sourceFingerprint")
    need(isinstance(fp, str) and re.fullmatch(r"[0-9a-f]{64}", fp or ""), "checkpoint source fingerprint must be verified SHA-256")
    need(fp == stable.get("sourceFingerprint"), "checkpoint fingerprint / current Stable fingerprint drift")
    metadata_rule = cp.get("metadataHeadRule") or cp.get("stableIdentityRule")
    if not metadata_rule and isinstance(cp.get("postReleaseOperationalMetadata"), dict):
        metadata_rule = cp["postReleaseOperationalMetadata"].get("rule")
    need(bool(metadata_rule), "checkpoint metadata/stable identity rule missing")
    gates = cp.get("gates", {})
    need(all(f"G{i}" in gates for i in range(17)), "checkpoint must contain G0-G16 states")
    need(cp.get("nextStep"), "checkpoint nextStep missing")
    need("updatedAt" in cp, "checkpoint updatedAt missing")
    portability = cp.get("assistantPortability", {})
    need(portability.get("status") == "ENFORCED", "checkpoint assistantPortability must be ENFORCED")
    need(portability.get("authoritativeHandoff") == "handoff/CURRENT.md", "checkpoint authoritative handoff drift")
    need(portability.get("entrypoints") == ["AGENTS.md", "CLAUDE.md"], "checkpoint assistant entrypoints drift")
    expected_handoff_tag = stable_tag_for(stable_release)
    need(f"**Certified Stable:** `{expected_handoff_tag}`" in handoff, "current handoff must identify the immutable Stable checkpoint release")

_ev = load_json(evidence_path)
if _ev:
    schema = _ev.get("schemaVersion")
    need(isinstance(schema, int) and schema >= 2, "release evidence checkpoint schemaVersion must be v2 or later")
    evidence_release = clean_version(_ev.get("release"))
    need(evidence_release == stable_release, "release evidence checkpoint must match immutable Stable build checkpoint")
    need(_ev.get("channel") == "STABLE", "release evidence checkpoint must describe certified Stable evidence")
    stable_ev = _ev.get("stable", {}) if isinstance(_ev.get("stable"), dict) else {}
    need(stable_ev.get("tag") == stable_tag_for(stable_release), "release evidence Stable tag mismatch")
    need(stable_ev.get("promotionCommit") == stable.get("candidateSha"), "release evidence promotion commit / current Stable drift")
    need(stable_ev.get("sourceFingerprint") == stable.get("sourceFingerprint"), "release evidence fingerprint / current Stable drift")
    need(isinstance(_ev.get("evidence"), dict), "release evidence checkpoint evidence map missing")
    evidence_metadata_rule = _ev.get("metadataHeadRule")
    if not evidence_metadata_rule and isinstance(_ev.get("postReleaseOperationalMetadata"), dict):
        evidence_metadata_rule = _ev["postReleaseOperationalMetadata"].get("rule")
    need(bool(evidence_metadata_rule), "release evidence checkpoint metadata/stable identity rule missing")

try:
    continuity = subprocess.run(
        [sys.executable, str(ROOT / "tools" / "ci" / "post_stable_continuity_gate.py")],
        cwd=ROOT, text=True, capture_output=True, check=False,
    )
    need(continuity.returncode == 0, "post-Stable continuity contract failed: " + (continuity.stdout + continuity.stderr).strip())
except Exception as exc:
    errors.append(f"post-Stable continuity contract unreadable: {exc}")

if errors:
    print("Adaptive Build Resume Contract: FAIL")
    for error in errors:
        print(" -", error)
    sys.exit(1)

mode = "PRODUCT_RELEASE_CLOSURE" if release_closure else "STABLE_ALIGNED"
print(
    "Adaptive Build Resume Contract: PASS · "
    f"mode={mode} · published Stable=v{stable_release} · candidate=v{identity_release} · "
    "GitHub-only ChatGPT/Codex/Claude portability enforced · current-state/work-slice truth bound · "
    "three-workflow control plane enforced · post-Stable continuity enforced · metadata fingerprint-excluded"
)
