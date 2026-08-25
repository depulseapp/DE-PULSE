#!/usr/bin/env python3
"""Fail-closed T9 packaged cross-platform runtime/release/provenance assurance."""
from __future__ import annotations

import argparse
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
T9 = PROGRAM / "T9_PACKAGED_CROSS_PLATFORM_RELEASE_ASSURANCE.json"
T8 = PROGRAM / "T8_PERFORMANCE_LOAD_SOAK_CONCURRENCY_RESOURCE_ASSURANCE.json"
CURRENT = ROOT / "governance" / "current-state.json"
WORK_SLICE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "work-slice.json"
CLOSURE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "closure.json"
IDENTITY = ROOT / "release_identity.json"
VERSION = ROOT / "VERSION.txt"
BOOTSTRAP = ROOT / "app_bootstrap.go"
RENDERER_IDENTITY = ROOT / "renderer" / "release-identity.js"
RENDERER_INDEX = ROOT / "renderer" / "index.html"
RELEASE_CONTRACT = ROOT / "release" / "v18.10.0" / "release_contract.json"
CERT_MANIFEST = ROOT / "release" / "v18.10.0" / "certification-manifest.json"
PREVIOUS_STABLE = ROOT / "release" / "v18.9.1" / "stable-evidence-manifest.json"
MAC = ROOT / "tools" / "release" / "native_macos.sh"
WINDOWS = ROOT / "tools" / "release" / "native_windows.ps1"
G15 = ROOT / "tools" / "release" / "g15_assurance.py"
PLANNER = ROOT / "tools" / "ci" / "impact_plan.py"
PLANNER_TEST = ROOT / "tools" / "ci" / "impact_plan_self_test.py"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
CI_QUALIFIED = ROOT / ".github" / "workflows" / "ci-qualified.yml"
RELEASE_WORKFLOW = ROOT / ".github" / "workflows" / "release.yml"

EXPECTED_IDENTITY = {
    "version": "18.10.0",
    "display_version": "DE.PULSE v18.10.0",
    "build_id": "v18.10.0-stable-20260825",
    "channel": "STABLE",
    "stable_baseline": "v18.9.1",
    "previous_stable": "v18.9.1",
    "bundle_version": "181000",
    "runtime_config": "PersonalMarketTerminal",
    "application_bundle": "De-Pulse.app",
    "renderer_identity_asset": "release-identity.js",
}
EXPECTED_PREVIOUS = {
    "release": "v18.9.1",
    "status": "STABLE_PUBLISHED",
    "buildId": "v18.9.1-stable-20260821",
    "stableTag": "v18.9.1-stable",
    "certifiedCandidate": "e55d8d25b15cec2ffb0f5411bc358bc40b359cf9",
    "sourceFingerprint": "0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff",
}
REQUIRED_CONTROL_AREAS = {
    "candidate identity and monotonic platform build",
    "exact source SHA and Git-object fingerprint",
    "macOS Apple Silicon clean package extraction and launch",
    "macOS Apple Silicon fresh and warm native-window lifecycle",
    "macOS previous-Stable profile upgrade preservation",
    "Windows x64 clean package extraction and launch",
    "Windows x64 fresh and warm runtime lifecycle",
    "Windows previous-Stable profile upgrade preservation",
    "SQLite migrations/integrity/state retention after upgrade",
    "renderer contracts and role/responsive behavior",
    "Chrome primary behavior",
    "WebKit primary compatibility",
    "backend full/race/randomized qualification",
    "persistence/DB integration",
    "security/data-rights contracts",
    "artifact SHA-256 and source provenance",
    "G15 same-candidate cross-platform binding",
    "no-rebuild Stable publication semantics",
}
REQUIRED_RULES = {
    "sameSourceHeadAcrossAllT9Qualification",
    "sameGitObjectFingerprintAcrossPlatforms",
    "actualPackagedRuntimeRequired",
    "freshInstallRequired",
    "warmRelaunchRequired",
    "previousStableUpgradeRequired",
    "persistenceIntegrityRequired",
    "rendererChromeWebkitRequired",
    "backendRaceRandomizedRequired",
    "securityDataRightsRequired",
    "artifactDigestAndProvenanceRequired",
    "noRebuildPublicationSemanticsMustRemainIntact",
    "t9MayNotPublishStable",
    "t9MayNotStartT10",
    "noParallelPackagingOrReleaseSubsystem",
}


def load(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8-sig"))


def text(path: Path) -> str:
    return path.read_text(encoding="utf-8", errors="replace")


def contains(path: Path, *needles: str) -> bool:
    body = text(path)
    return all(needle in body for needle in needles)


def closure_row(closure: dict, row_id: str) -> dict | None:
    return next((row for row in closure.get("gaps", []) if isinstance(row, dict) and row.get("id") == row_id), None)


def governed_status(product: dict, track: str) -> str:
    return next((str(row.get("status") or "") for row in product.get("nextGovernedTracks", []) if isinstance(row, dict) and row.get("track") == track), "")


def git_hash_object(path: Path) -> str:
    result = subprocess.run(
        ("git", "hash-object", str(path.relative_to(ROOT))),
        cwd=ROOT,
        check=True,
        text=True,
        capture_output=True,
    )
    return result.stdout.strip()


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()
    errors: list[str] = []

    required_paths = (
        T9, T8, CURRENT, WORK_SLICE, CLOSURE, IDENTITY, VERSION, BOOTSTRAP,
        RENDERER_IDENTITY, RENDERER_INDEX, RELEASE_CONTRACT, CERT_MANIFEST,
        PREVIOUS_STABLE, MAC, WINDOWS, G15, PLANNER, PLANNER_TEST, CI_FAST,
        CI_QUALIFIED, RELEASE_WORKFLOW,
    )
    for path in required_paths:
        if not path.is_file():
            errors.append(f"required T9 owner missing: {path.relative_to(ROOT)}")
    if errors:
        print("V18 T9 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1

    t9 = load(T9)
    t8 = load(T8)
    current = load(CURRENT)
    work_slice = load(WORK_SLICE)
    closure = load(CLOSURE)
    identity = load(IDENTITY)
    contract = load(RELEASE_CONTRACT)
    manifest = load(CERT_MANIFEST)
    previous = load(PREVIOUS_STABLE)
    product = current.get("productCapabilityGate") or {}
    state = str(t9.get("state") or "")

    if t9.get("schema") != "DE.PULSE-V18-T9-PACKAGED-CROSS-PLATFORM-RELEASE-ASSURANCE-1" or t9.get("programIssue") != 113 or t9.get("trackIssue") != 122 or t9.get("targetVersion") != "v18.10.0":
        errors.append("T9 assurance identity/schema mismatch")
    if set(t9.get("requiredControlAreas") or []) != REQUIRED_CONTROL_AREAS:
        errors.append("T9 required control-area matrix drifted")
    rules = t9.get("rules") or {}
    for key in sorted(REQUIRED_RULES):
        if rules.get(key) is not True:
            errors.append(f"T9 rule must remain true: {key}")
    if t9.get("stablePublishedByThisArtifact") is not False or t9.get("t10StartedByThisArtifact") is not False:
        errors.append("T9 artifact must not publish Stable or start T10")
    if t9.get("requiredPlatforms") != ["macOS Apple Silicon", "Windows x64"]:
        errors.append("T9 required packaged platforms drifted")
    if t9.get("requiredBrowserSurfaces") != ["Chrome", "WebKit"]:
        errors.append("T9 primary browser surfaces drifted")
    if t9.get("linuxScope") != "CI_TEST_ONLY_NOT_RELEASE_TARGET" or t9.get("hostedWebGA") != "NOT_CLAIMED_BY_V18":
        errors.append("T9 Linux/hosted-Web applicability boundary drifted")

    if t8.get("state") != "COMPLETE" or t8.get("knownCoverageGaps"):
        errors.append("T9 requires zero-gap COMPLETE T8")
    for prior_id in (
        "T1-FEATURE-TRACEABILITY", "T2-UNIT-CONTRACT-PROPERTY", "T3-FUNCTIONAL-E2E",
        "T4-EDGE-FAILURE-DATA-TRUTH", "T5-PERSISTENCE-LIFECYCLE-RECOVERY",
        "T6-SECURITY-ROLES-RIGHTS", "T7-UI-UX-IA-CONTENT", "T8-PERFORMANCE-CONCURRENCY-SOAK",
    ):
        row = closure_row(closure, prior_id)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T9 requires prior durable VERIFIED closure: {prior_id}")

    for key, expected in EXPECTED_IDENTITY.items():
        if identity.get(key) != expected:
            errors.append(f"v18.10.0 release identity mismatch: {key}={identity.get(key)!r}, expected {expected!r}")
    if not contains(VERSION, "DE.PULSE v18.10.0", "Build: v18.10.0-stable-20260825", "Previous Stable: v18.9.1"):
        errors.append("VERSION.txt is not synchronized to v18.10.0 candidate identity")
    if not contains(BOOTSTRAP, 'const appVersion = "18.10.0"', 'const buildID = "v18.10.0-stable-20260825"', 'const releaseChannel = "STABLE"'):
        errors.append("app_bootstrap.go is not synchronized to v18.10.0 candidate identity")
    if not contains(RENDERER_IDENTITY, "DEPULSE_RELEASE_VERSION = '18.10.0'", "DEPULSE_RELEASE_BUILD_ID = 'v18.10.0-stable-20260825'", "release/v18.10.0/release_contract.json"):
        errors.append("renderer release identity is not synchronized to v18.10.0 candidate")
    renderer_hash = git_hash_object(RENDERER_IDENTITY)
    if f"release-identity.js?v={renderer_hash[:16]}" not in text(RENDERER_INDEX):
        errors.append("renderer/index.html release-identity cache token does not match current Git blob")
    if "<title>DE.PULSE v18.10.0</title>" not in text(RENDERER_INDEX):
        errors.append("renderer/index.html title is not v18.10.0")

    if contract.get("schema") != "DE.PULSE-V18.10.0-RELEASE-CONTRACT-1" or contract.get("release") != "18.10.0" or contract.get("certifiedStableBaseline") != "v18.9.1-stable":
        errors.append("v18.10.0 release contract identity/baseline mismatch")
    if contract.get("identity_asset") != "release-identity.js" or "ADAPT-V18-FINAL-CLOSURE-10-10-001" not in str(contract.get("scope") or ""):
        errors.append("v18.10.0 release contract scope/identity owner mismatch")
    if manifest.get("schema") != "DE.PULSE-G12-EVIDENCE-MANIFEST-1" or manifest.get("productVersion") != "18.10.0" or manifest.get("workSliceId") != "ADAPT-V18-FINAL-CLOSURE-10-10-001" or manifest.get("evidenceSchemaVersion") != 1:
        errors.append("v18.10.0 certification manifest identity mismatch")
    python_gates = {tuple(row) for row in manifest.get("pythonGates", []) if isinstance(row, list)}
    if ("python3", "tools/ci/v18_t9_packaged_cross_platform_release_assurance_gate.py") not in python_gates:
        errors.append("v18.10.0 certification manifest does not bind the T9 assurance gate")

    for key, expected in EXPECTED_PREVIOUS.items():
        if previous.get(key) != expected:
            errors.append(f"previous Stable evidence mismatch: {key}={previous.get(key)!r}, expected {expected!r}")
    if (previous.get("gates") or {}).get("StablePublication") != "PASS_NO_REBUILD":
        errors.append("previous Stable must retain PASS_NO_REBUILD publication evidence")
    previous_t9 = t9.get("previousStable") or {}
    if previous_t9.get("tag") != EXPECTED_PREVIOUS["stableTag"] or previous_t9.get("candidateSha") != EXPECTED_PREVIOUS["certifiedCandidate"] or previous_t9.get("sourceFingerprint") != EXPECTED_PREVIOUS["sourceFingerprint"] or previous_t9.get("buildId") != EXPECTED_PREVIOUS["buildId"]:
        errors.append("T9 previous-Stable identity does not match immutable Stable evidence")

    common_previous_markers = (
        "v18.9.1-stable",
        "e55d8d25b15cec2ffb0f5411bc358bc40b359cf9",
        "0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff",
        "v18.9.1-stable-20260821",
        "previous-Stable upgrade",
        "DE.PULSE-G13-G14-NATIVE-3",
        "previousStableCertifiedHarness",
        "upgradeStatePreserved",
    )
    for path, label in ((MAC, "macOS"), (WINDOWS, "Windows")):
        for marker in common_previous_markers:
            if marker not in text(path):
                errors.append(f"{label} T9 native harness missing marker: {marker}")
    if not contains(MAC, 'GITHUB_WORKSPACE="$PREVIOUS_WORKTREE"', "git worktree add --detach", "current-fresh", "current-warm", "run_native_window_cycle fresh", "run_native_window_cycle warm", "upgrade-before.json", "upgrade-after.json"):
        errors.append("macOS native harness lost fresh/warm/native/previous-Stable worktree upgrade evidence")
    if not contains(WINDOWS, "$env:GITHUB_WORKSPACE=$previousWorktree", "git worktree add --detach", "current-fresh", "current-warm", "previous-stable", "upgrade-current", "upgrade-before.json", "upgrade-after.json"):
        errors.append("Windows native harness lost fresh/warm/previous-Stable worktree upgrade evidence")
    if "De-Pulse-v${DEPULSE_VERSION}-Stable-macOS-Apple-Silicon.zip" not in text(MAC):
        errors.append("macOS canonical package name drifted")
    if 'De-Pulse-v$($env:DEPULSE_VERSION)-Stable-Windows-x64.zip' not in text(WINDOWS):
        errors.append("Windows canonical package name drifted")

    planner = text(PLANNER)
    for marker in (
        "T9_CROSS_LAYER_ASSURANCE_FILES",
        "T9_PACKAGED_CROSS_PLATFORM_RELEASE_ASSURANCE.json",
        "release/v18.10.0/release_contract.json",
        "release/v18.10.0/certification-manifest.json",
        '"RENDERER_UI"', '"AUTH_SECURITY"', '"DATA_RIGHTS"', '"PERSISTENCE"',
        "path in T9_CROSS_LAYER_ASSURANCE_FILES",
    ):
        if marker not in planner:
            errors.append(f"Planner v3 lost T9 full-evidence routing marker: {marker}")
    planner_test = text(PLANNER_TEST)
    for marker in (
        "T9 assurance state must not be treated as process-only",
        "T9 assurance must select full cross-layer qualification",
        "T9 gate must select both native platforms",
        "T9 manifest must select both native packages",
    ):
        if marker not in planner_test:
            errors.append(f"Planner v3 self-test lost T9 invariant: {marker}")

    qualified = text(CI_QUALIFIED)
    for marker in (
        "Full Go suite", "go test -count=1 ./...", "Race detector", "go test -race -count=1 ./...",
        "Randomized package order", "go test -shuffle=on -count=1 ./...",
        "Qualified persistence / DB integration", "Persistence-focused regression",
        "Qualified renderer contracts", "Primary Chrome behavior", "Primary WebKit browser compatibility",
        "Qualified security / data-rights contracts", "Qualified macOS native lifecycle rehearsal",
        "Qualified Windows native runtime rehearsal", "tools/release/native_macos.sh", "tools/release/native_windows.ps1",
    ):
        if marker not in qualified:
            errors.append(f"Qualified CI lost required T9 lane/command: {marker}")

    fast = text(CI_FAST)
    if "python3 tools/ci/v18_t9_packaged_cross_platform_release_assurance_gate.py" not in fast:
        errors.append("T9 assurance gate is not bound into canonical CI Fast")

    g15 = text(G15)
    for marker in (
        "certifiedSourceSha", "sourceFingerprint", "buildId", "artifactSha256",
        "assert actual == data['artifactSha256']", "macOS Apple Silicon", "Windows x64",
        "DE.PULSE-G15-ASSURANCE-2", "noExecutionBoundary",
    ):
        if marker not in g15:
            errors.append(f"G15 same-candidate/artifact binding drifted: {marker}")

    release = text(RELEASE_WORKFLOW)
    for marker in (
        "DE.PULSE/fast-head", "DE.PULSE/qualified-head",
        "tools/release/run_full_certification.py", "G13/G14 macOS Apple Silicon", "G13/G14 Windows x64",
        "tools/release/native_macos.sh", "tools/release/native_windows.ps1",
        "G15 Release Assurance", "Publish exact same-run certified artifacts",
        "Publish Stable release without rebuild or overwrite",
        "immutable Stable asset conflict for $name",
        "existing Stable asset byte-identical; reuse",
        "published/reused $tag from this exact certified run without rebuilding or overwriting differing bytes",
        "depulse-stable-release",
    ):
        if marker not in release:
            errors.append(f"canonical release no-rebuild/provenance contract drifted: {marker}")

    if (product.get("downstreamAssuranceState") or {}).get("T10Started") is not False:
        errors.append("T10 must remain not started during T9 closure handoff")

    gaps = [row for row in t9.get("knownCoverageGaps", []) if isinstance(row, dict)]
    open_gap_ids = {str(row.get("id")) for row in gaps if row.get("status") == "OPEN"}
    t9_closure = closure_row(closure, "T9-PACKAGED-CROSS-PLATFORM-RELEASE")
    completed = product.get("completedChildTracks") or []
    completed_t9 = next((row for row in completed if isinstance(row, dict) and row.get("track") == "T9" and row.get("issue") == 122), None)

    if state == "IN_PROGRESS":
        if work_slice.get("nextTrack") != {"track": "T9", "issue": 122}:
            errors.append("IN_PROGRESS T9 requires work-slice nextTrack T9/#122")
        if product.get("nextChildIssue") != 122 or product.get("nextChildTrack") != "T9":
            errors.append("IN_PROGRESS T9 requires current-state active child T9/#122")
        if product.get("nextCompanionChildIssue") != 123 or product.get("nextCompanionChildTrack") != "T10":
            errors.append("IN_PROGRESS T9 requires T10/#123 companion")
        if governed_status(product, "T9") != "IN_PROGRESS" or governed_status(product, "T10") != "NOT_STARTED":
            errors.append("IN_PROGRESS T9 requires T9 IN_PROGRESS and T10 NOT_STARTED in current-state")
        if not isinstance(t9_closure, dict) or t9_closure.get("status") != "IMPLEMENTED_UNVERIFIED":
            errors.append("IN_PROGRESS T9 closure row must be IMPLEMENTED_UNVERIFIED")
        expected_gaps = {"T9-EXACT-HEAD-FAST", "T9-IDENTICAL-HEAD-QUALIFIED", "T9-NATIVE-EVIDENCE-BINDING"}
        if open_gap_ids != expected_gaps:
            errors.append(f"IN_PROGRESS T9 open gaps must be exactly {sorted(expected_gaps)}")
        if completed_t9 is not None:
            errors.append("IN_PROGRESS T9 must not already appear in completedChildTracks")
    elif state == "COMPLETE":
        if work_slice.get("nextTrack") != {"track": "T10", "issue": 123}:
            errors.append("COMPLETE T9 requires work-slice handoff to T10/#123")
        if work_slice.get("nextCompanionTrack") not in (None, {}):
            errors.append("COMPLETE T9 must not retain a companion track beside sole remaining T10")
        if not isinstance(completed_t9, dict) or completed_t9.get("status") != "COMPLETE" or not str(completed_t9.get("frozenHeadSha") or "") or not completed_t9.get("fastRunId") or not completed_t9.get("qualifiedRunId"):
            errors.append("COMPLETE T9 requires durable exact-head Fast/Qualified completedChildTracks evidence")
        if completed_t9.get("frozenHeadSha") != t9.get("finalSourceSha") or completed_t9.get("fastRunId") != t9.get("fastRunId") or completed_t9.get("qualifiedRunId") != t9.get("qualifiedRunId"):
            errors.append("COMPLETE T9 current-state evidence must match T9 assurance artifact")
        if not isinstance(t9_closure, dict) or t9_closure.get("status") != "VERIFIED":
            errors.append("COMPLETE T9 closure row must be VERIFIED")
        if gaps:
            errors.append("COMPLETE T9 cannot retain coverage gaps")
        if product.get("nextChildIssue") != 123 or product.get("nextChildTrack") != "T10":
            errors.append("COMPLETE T9 must hand off next child to T10/#123")
        if product.get("nextCompanionChildIssue") is not None or product.get("nextCompanionChildTrack") is not None:
            errors.append("COMPLETE T9 must not retain a companion child beside sole remaining T10")
        if governed_status(product, "T10") != "NOT_STARTED":
            errors.append("COMPLETE T9 must not silently start T10")
    else:
        errors.append(f"unsupported T9 state: {state!r}")

    strict = args.strict or state == "COMPLETE"
    if strict and gaps:
        errors.append("strict T9 closure cannot retain coverage gaps")

    print("V18 T9 PACKAGED / CROSS-PLATFORM / RELEASE / PROVENANCE ASSURANCE")
    print(f"state: {state}")
    print(f"candidate: v{identity.get('version')} / {identity.get('build_id')} / platform-build {identity.get('bundle_version')}")
    print(f"previous Stable: {previous.get('stableTag')} @ {previous.get('certifiedCandidate')}")
    print(f"required control areas: {len(REQUIRED_CONTROL_AREAS)}")
    print(f"open gaps: {len(open_gap_ids)}")
    print("Qualified graph: backend/race/randomized + DB + renderer + Chrome + WebKit + security/data-rights + macOS + Windows")
    print("native upgrade: immutable v18.9.1 Stable harness/profile -> exact v18.10.0 package on both required platforms")
    print("Linux: CI/test only; Hosted Web GA: not claimed by v18")
    print("T9 publication: PROHIBITED; T10 authorization: NOT IMPLIED")

    if errors:
        print("V18 T9 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T9 ASSURANCE GATE: PASS (strict={strict})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
