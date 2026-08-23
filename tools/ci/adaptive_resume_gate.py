#!/usr/bin/env python3
"""Validate the permanent GitHub-backed DE.PULSE Build Resume Protocol.

#70 convergence note: current resume truth is owned by governance/current-state.json,
the registered work-slice metadata, the three canonical workflows, the current
release identity, and immutable Stable checkpoints. Retired v18.6 CI plans are
historical evidence only and are deliberately not runtime dependencies here.
"""
from pathlib import Path
import json
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
errors: list[str] = []


def need(ok: bool, msg: str) -> None:
    if not ok:
        errors.append(msg)


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
        "last trustworthy PASS",
        "source fingerprint",
        "G0–G16",
        "GitHub",
        "G16",
        "No Execution",
        "handoff/CURRENT.md",
        "AI-ASSISTANT-PORTABILITY-CONTRACT.md",
        "Claude",
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
        "GitHub source-of-truth hierarchy",
        "Mandatory fresh-session algorithm",
        "Durable handoff rule",
        "Secrets and account independence",
        "No upload of an old chat handoff is required",
    ):
        need(term in contract, f"portability contract missing required section/term: {term}")
    need("SUPERSEDES ALL PRIOR CHAT HANDOFFS" in handoff, "handoff/CURRENT.md must be the single current handoff authority")
    need("Exactly one next action" in handoff, "handoff/CURRENT.md must name exactly one next action")

identity_path = ROOT / "release_identity.json"
state_path = ROOT / "governance" / "current-state.json"
checkpoint_path = ROOT / ".depulse-certification" / "resume" / "build-checkpoint.json"
evidence_path = ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json"

try:
    ident = json.loads(identity_path.read_text())
except Exception as exc:
    ident = {}
    errors.append(f"release identity unreadable: {exc}")

identity_release = str(ident.get("version", "")).lstrip("v")
identity_build = str(ident.get("build_id", ""))
identity_previous_stable = str(ident.get("previous_stable", "")).lstrip("v")
identity_stable_baseline = str(ident.get("stable_baseline", "")).lstrip("v")

try:
    state = json.loads(state_path.read_text())
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
    gate = state.get("productCapabilityGate", {}) if isinstance(state.get("productCapabilityGate"), dict) else {}

    need(stable.get("productVersion") == identity_release, "current-state Stable productVersion / release identity drift")
    need(stable.get("buildId") == identity_build, "current-state Stable buildId / release identity drift")
    need(stable.get("platformBuildNumber") == str(ident.get("bundle_version", "")), "current-state platform build number / release identity drift")
    need(stable.get("tag") == f"v{identity_release}-stable", "current-state Stable tag / release identity drift")
    need(isinstance(stable.get("candidateSha"), str) and re.fullmatch(r"[0-9a-f]{40}", stable.get("candidateSha", "")), "current-state Stable candidateSha invalid")
    need(isinstance(stable.get("sourceFingerprint"), str) and re.fullmatch(r"[0-9a-f]{64}", stable.get("sourceFingerprint", "")), "current-state Stable sourceFingerprint invalid")
    need(stable.get("publication") == "PASS_NO_REBUILD", "current-state Stable publication must remain PASS_NO_REBUILD")

    work_slice_id = str(active.get("workSliceId", "")).strip()
    need(bool(work_slice_id), "current-state active work slice missing")
    need(active.get("publicProductVersion") is None, "process work slice must not consume a public product version")
    need(active.get("productBehaviorChange") is False, "#70 process work slice must remain product-behavior neutral")
    need(gate.get("blocked") is True and gate.get("blockedByIssue") == 70, "next product capability must remain blocked by issue #70")

    work_slice_path = ROOT / "governance" / "work-slices" / work_slice_id / "work-slice.json"
    work_slice = json.loads(work_slice_path.read_text())
    need(work_slice.get("workSliceId") == work_slice_id, "current-state/work-slice ID mismatch")
    need(work_slice.get("issue") == active.get("issue"), "current-state/work-slice issue mismatch")
    need(work_slice.get("branch") == active.get("branch"), "current-state/work-slice branch mismatch")
    need(work_slice.get("publicProductVersion") is None, "registered process work slice consumed a public product version")
    need(work_slice.get("productBehaviorChange") is False, "registered process work slice must remain product-behavior neutral")
    need(work_slice.get("baselineCandidateSha") == stable.get("candidateSha"), "work-slice baseline candidate / Stable candidate drift")
    need(work_slice.get("baselineSourceFingerprint") == stable.get("sourceFingerprint"), "work-slice baseline fingerprint / Stable fingerprint drift")
    need(work_slice.get("baselineBuildId") == stable.get("buildId"), "work-slice baseline build / Stable build drift")
    need(work_slice.get("blocksNextProductCapability") is True, "work slice must block next product capability until closure")

    workflow_dir = ROOT / ".github" / "workflows"
    routine = sorted(path.name for path in workflow_dir.glob("*.yml"))
    need(routine == ["ci-fast.yml", "ci-qualified.yml", "release.yml"], f"routine workflow set drift: {routine}")

    prefreeze_text = (ROOT / "tools" / "release" / "prefreeze_qualification.py").read_text(errors="ignore")
    need("'.depulse-certification'" in prefreeze_text or '".depulse-certification"' in prefreeze_text, "prefreeze fingerprint must exclude .depulse-certification metadata")
except Exception as exc:
    errors.append(f"canonical current-state/work-slice validation failed: {exc}")

stable_release = ""
try:
    cp = json.loads(checkpoint_path.read_text())
    schema = cp.get("schemaVersion")
    need(isinstance(schema, int) and schema >= 2, "build checkpoint schemaVersion must be v2 or later")
    stable_release = str(cp.get("release", "")).lstrip("v")
    need(bool(stable_release), "build checkpoint Stable release missing")
    need(cp.get("channel") == "STABLE", "build checkpoint must describe certified Stable evidence")
    need(cp.get("branch"), "checkpoint branch missing")

    certified = cp.get("certifiedStable", {}) if isinstance(cp.get("certifiedStable"), dict) else {}
    certified_release = str(certified.get("version", "")).lstrip("v")
    need(certified_release == stable_release, "checkpoint release must match certifiedStable version")

    if identity_release != stable_release:
        need(identity_previous_stable == stable_release, "current Stable previous_stable must match immutable predecessor checkpoint")
        need(identity_stable_baseline == stable_release, "current Stable stable_baseline must match immutable predecessor checkpoint")
    else:
        need(bool(identity_previous_stable), "promoted Stable previous_stable must identify its predecessor")
        need(bool(identity_stable_baseline), "promoted Stable stable_baseline must identify its certified predecessor baseline")
        need(identity_previous_stable == identity_stable_baseline, "promoted Stable previous_stable/stable_baseline mismatch")

    candidate = cp.get("candidateSourceCommit") or certified.get("candidateSourceCommit") or certified.get("certifiedSourceCheckout")
    need(isinstance(candidate, str) and re.fullmatch(r"[0-9a-f]{40}", candidate), "checkpoint candidate source commit must be a Git SHA")

    metadata_rule = cp.get("metadataHeadRule") or cp.get("stableIdentityRule")
    if not metadata_rule:
        post_release = cp.get("postReleaseOperationalMetadata", {})
        if isinstance(post_release, dict):
            metadata_rule = post_release.get("rule")
    need(bool(metadata_rule), "checkpoint metadata/stable identity rule missing")

    fp = cp.get("sourceFingerprint") or certified.get("sourceFingerprint")
    fp_state = cp.get("sourceFingerprintState")
    need((isinstance(fp, str) and re.fullmatch(r"[0-9a-f]{64}", fp)) or (fp is None and fp_state in {"PENDING_REQUALIFICATION", "NOT_FROZEN"}), "checkpoint source fingerprint must be verified SHA-256 or explicitly pending")
    gates = cp.get("gates", {})
    need(all(f"G{i}" in gates for i in range(17)), "checkpoint must contain G0-G16 states")
    need(cp.get("nextStep"), "checkpoint nextStep missing")
    need("updatedAt" in cp, "checkpoint updatedAt missing")
    portability = cp.get("assistantPortability", {})
    need(portability.get("status") == "ENFORCED", "checkpoint assistantPortability must be ENFORCED")
    need(portability.get("authoritativeHandoff") == "handoff/CURRENT.md", "checkpoint authoritative handoff drift")
    need(portability.get("entrypoints") == ["AGENTS.md", "CLAUDE.md"], "checkpoint assistant entrypoints drift")

    stable_handoff_ok = any(
        marker in handoff
        for marker in (
            f"**Certified Stable:** `v{stable_release}-stable`",
            f"**Release:** `v{stable_release}`",
            f"**Immutable predecessor resume checkpoint release:** `v{stable_release}` / `v{stable_release}-stable`",
        )
    )
    need(stable_handoff_ok, "current handoff must identify the immutable Stable checkpoint release")
except Exception as exc:
    errors.append(f"build checkpoint invalid/unreadable: {exc}")

try:
    ev = json.loads(evidence_path.read_text())
    schema = ev.get("schemaVersion")
    need(isinstance(schema, int) and schema >= 2, "release evidence checkpoint schemaVersion must be v2 or later")
    evidence_release = str(ev.get("release", "")).lstrip("v")
    need(evidence_release == stable_release, "release evidence checkpoint must match immutable Stable build checkpoint")
    need(ev.get("channel") == "STABLE", "release evidence checkpoint must describe certified Stable evidence")
    stable = ev.get("stable", {}) if isinstance(ev.get("stable"), dict) else {}
    need(stable.get("tag") == f"v{stable_release}-stable", "release evidence Stable tag mismatch")
    need(isinstance(ev.get("evidence"), dict), "release evidence checkpoint evidence map missing")
    evidence_metadata_rule = ev.get("metadataHeadRule")
    if not evidence_metadata_rule:
        post_release = ev.get("postReleaseOperationalMetadata", {})
        if isinstance(post_release, dict):
            evidence_metadata_rule = post_release.get("rule")
    need(bool(evidence_metadata_rule), "release evidence checkpoint metadata/stable identity rule missing")
except Exception as exc:
    errors.append(f"release evidence checkpoint invalid/unreadable: {exc}")

try:
    continuity = subprocess.run(
        [sys.executable, str(ROOT / "tools" / "ci" / "post_stable_continuity_gate.py")],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    need(continuity.returncode == 0, "post-Stable continuity contract failed: " + (continuity.stdout + continuity.stderr).strip())
except Exception as exc:
    errors.append(f"post-Stable continuity contract unreadable: {exc}")

if errors:
    print("Adaptive Build Resume Contract: FAIL")
    for error in errors:
        print(" -", error)
    sys.exit(1)

print(
    "Adaptive Build Resume Contract: PASS · "
    f"canonical Stable=v{identity_release} · immutable predecessor checkpoint=v{stable_release} · "
    "GitHub-only ChatGPT/Codex/Claude portability enforced · current-state/work-slice truth bound · "
    "three-workflow control plane enforced · post-Stable continuity enforced · metadata fingerprint-excluded"
)
