#!/usr/bin/env python3
from __future__ import annotations

from impact_plan import FAILURE_TAXONOMY, analyze_changed_paths, classify_path


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def main() -> int:
    release_classes = classify_path(".github/workflows/release.yml")
    require("CI_HARNESS" in release_classes, "release workflow must be CI_HARNESS")
    require("RELEASE_TOOLING" in release_classes, "release workflow must be RELEASE_TOOLING")

    renderer = analyze_changed_paths(["renderer/watchlist-ui.js"])
    require(renderer["qualifiedLane"] == "renderer", "renderer change should use renderer adaptive lane")
    require(renderer["rendererRequired"] is True, "renderer contract job missing")
    require(renderer["chromeRequired"] is True, "renderer change must require Chrome")
    require(renderer["webkitRequired"] is True, "renderer change must require WebKit")
    require(renderer["backendRequired"] is False, "renderer-only change must not force backend")

    organized_browser_test = analyze_changed_paths(["tests/browser/watchlist_membership_test.py"])
    require(organized_browser_test["rendererRequired"] is True, "browser test must retain renderer contracts")
    require(organized_browser_test["chromeRequired"] is True, "browser test must retain Chrome")
    require(organized_browser_test["webkitRequired"] is True, "browser test must retain WebKit")

    deterministic_owner = analyze_changed_paths(["tools/ci/deterministic_equivalence_test.js"])
    require(deterministic_owner["processOnly"] is False, "deterministic capability evidence must not be process-only")
    require(deterministic_owner["qualifiedLane"] == "renderer", "deterministic evidence owner must use renderer lane")
    require(deterministic_owner["rendererRequired"] is True, "deterministic evidence owner must execute renderer contracts")
    require(deterministic_owner["chromeRequired"] is True, "deterministic evidence owner must retain Chrome evidence")
    require(deterministic_owner["webkitRequired"] is True, "deterministic evidence owner must retain WebKit evidence")

    webkit_owner = analyze_changed_paths(["tools/ci/webkit_browser_test.py"])
    require(webkit_owner["processOnly"] is False, "canonical WebKit evidence owner must not be process-only")
    require(webkit_owner["qualifiedLane"] == "browser", "canonical WebKit evidence owner must use browser lane")
    require(webkit_owner["webkitRequired"] is True, "canonical WebKit evidence owner must execute WebKit")
    require(webkit_owner["chromeRequired"] is False, "WebKit-only owner should not consume Chrome by itself")
    require(webkit_owner["nativeMacosRequired"] is False, "WebKit browser owner must not become native-macOS proxy")

    legacy_webkit_helper = analyze_changed_paths(["tools/ci/webkit_targeted_test.py"])
    require(legacy_webkit_helper["webkitRequired"] is True, "legacy WebKit compatibility helper must still select WebKit while referenced")
    require(legacy_webkit_helper["nativeMacosRequired"] is False, "legacy WebKit helper must no longer own native macOS rehearsal")

    active_release_browser = analyze_changed_paths(["release/v18.5.1/browser_live_render_test.py"])
    require(active_release_browser["rendererRequired"] is True, "active release browser test must map to renderer/browser evidence")
    require(active_release_browser["chromeRequired"] is True, "active release browser test must retain Chrome")
    require(active_release_browser["webkitRequired"] is True, "active release browser test must fail safe to WebKit until capability migration")

    process = analyze_changed_paths([
        "tools/ci/workflow_policy.py",
        "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md",
    ])
    require(process["qualifiedLane"] == "ci-harness", "process-only change must use ci-harness")
    require(process["portabilityRequired"] is True, "process-only change must require portability")
    require(process["backendRequired"] is False, "process-only change must not force backend")

    t2_control = analyze_changed_paths([
        "governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T2_UNIT_CONTRACT_ASSURANCE.json"
    ])
    require(t2_control["processOnly"] is False, "active T2 assurance state must not be treated as process-only")
    require(t2_control["qualifiedLane"] == "full", "active T2 assurance state must select backend + renderer evidence")
    require(t2_control["backendRequired"] is True, "active T2 assurance state must require backend unit/package execution")
    require(t2_control["rendererRequired"] is True, "active T2 assurance state must require renderer contract execution")
    require(t2_control["chromeRequired"] is True and t2_control["webkitRequired"] is True, "renderer dependencies must remain intact during T2 qualification")
    require(t2_control["nativeMacosRequired"] is False and t2_control["nativeWindowsRequired"] is False, "T2 must not implicitly start packaged-platform qualification")

    t2_gate = analyze_changed_paths(["tools/ci/v18_t2_unit_contract_assurance_gate.py"])
    require(t2_gate["qualifiedLane"] == "full", "T2 gate changes must re-run their executable backend + renderer evidence graph")
    require(t2_gate["backendRequired"] is True and t2_gate["rendererRequired"] is True, "T2 gate dependency graph incomplete")

    t9_control = analyze_changed_paths([
        "governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T9_PACKAGED_CROSS_PLATFORM_RELEASE_ASSURANCE.json"
    ])
    require(t9_control["processOnly"] is False, "T9 assurance state must not be treated as process-only")
    require(t9_control["qualifiedLane"] == "full", "T9 assurance must select full cross-layer qualification")
    for key in (
        "backendRequired",
        "rendererRequired",
        "chromeRequired",
        "webkitRequired",
        "securityRightsRequired",
        "dbIntegrationRequired",
        "nativeMacosRequired",
        "nativeWindowsRequired",
    ):
        require(t9_control[key] is True, f"T9 assurance missing required evidence selection: {key}")
    require(t9_control["releaseRehearsalRequired"] is True, "T9 assurance must retain release rehearsal")
    require("RELEASE_TOOLING" in t9_control["changeClasses"], "T9 assurance must remain release-governed")

    t9_gate = analyze_changed_paths(["tools/ci/v18_t9_packaged_cross_platform_release_assurance_gate.py"])
    require(t9_gate["qualifiedLane"] == "full", "T9 gate changes must re-run full cross-platform evidence")
    require(t9_gate["nativeMacosRequired"] and t9_gate["nativeWindowsRequired"], "T9 gate must select both native platforms")

    t9_manifest = analyze_changed_paths(["release/v18.10.0/certification-manifest.json"])
    require(t9_manifest["qualifiedLane"] == "full", "v18.10.0 certification manifest must select full T9 evidence")
    require(t9_manifest["chromeRequired"] and t9_manifest["webkitRequired"], "T9 manifest must select both primary browser surfaces")
    require(t9_manifest["dbIntegrationRequired"] and t9_manifest["securityRightsRequired"], "T9 manifest must select DB and security/data-rights")
    require(t9_manifest["nativeMacosRequired"] and t9_manifest["nativeWindowsRequired"], "T9 manifest must select both native packages")

    stable_manifest = analyze_changed_paths(["release/v18.6.1/stable-evidence-manifest.json"])
    require(stable_manifest["qualifiedLane"] == "ci-harness", "durable Stable evidence index must stay process-only")
    require("RELEASE_TOOLING" in stable_manifest["changeClasses"], "Stable evidence index must remain release-governed")
    require(stable_manifest["webkitRequired"] is False, "Stable evidence-only work must not consume WebKit")

    non_manifest_release = analyze_changed_paths(["release/v18.9.1/run_full_certification.sh"])
    require(non_manifest_release["processOnly"] is False, "release executable must not be treated as process-only")
    require(non_manifest_release["qualifiedLane"] == "full", "release executable must stay full qualification")

    mac = analyze_changed_paths(["tools/release/native_macos.sh"])
    require(mac["nativeMacosRequired"] is True, "macOS harness change must require real native macOS rehearsal")
    require(mac["nativeWindowsRequired"] is False, "macOS-only harness change must not consume Windows")
    require(mac["releaseRehearsalRequired"] is True, "macOS harness change must require release rehearsal")

    windows = analyze_changed_paths(["tools/release/native_windows.ps1"])
    require(windows["nativeWindowsRequired"] is True, "Windows harness change must require real native Windows rehearsal")
    require(windows["nativeMacosRequired"] is False, "Windows-only harness change must not consume macOS")

    shared_release = analyze_changed_paths(["release_identity.json"])
    require(
        shared_release["nativeMacosRequired"] is True and shared_release["nativeWindowsRequired"] is True,
        "shared release identity must rehearse both required native platforms",
    )

    provider = analyze_changed_paths(["provider_router.go"])
    require(provider["qualifiedLane"] == "backend", "provider-only Go work should use backend adaptive lane")
    require(provider["backendRequired"] is True, "provider router must require backend")
    require(provider["webkitRequired"] is False, "provider-only work must not require WebKit")

    persistence = analyze_changed_paths(["canonical_state_sqlite.go"])
    require(persistence["backendRequired"] is True, "persistence must require backend")
    require(persistence["dbIntegrationRequired"] is True, "persistence must mark DB integration")

    rights = analyze_changed_paths(["governance/provider-data-rights-registry.md"])
    require(rights["securityRightsRequired"] is True, "data rights must require security/rights evidence")
    require("DATA_RIGHTS" in rights["changeClasses"], "data-rights class missing")

    mixed = analyze_changed_paths(["provider_router.go", "renderer/index.html"])
    require(mixed["backendRequired"] is True and mixed["rendererRequired"] is True, "mixed dependency graph incomplete")
    require(mixed["chromeRequired"] is True and mixed["webkitRequired"] is True, "mixed renderer browser evidence incomplete")

    unknown = analyze_changed_paths(["unclassified.future-format"])
    require(unknown["failClosedFull"] is True, "unknown paths must fail closed")
    require(
        all(
            unknown[key] is True
            for key in (
                "backendRequired",
                "rendererRequired",
                "chromeRequired",
                "webkitRequired",
                "securityRightsRequired",
                "dbIntegrationRequired",
                "nativeMacosRequired",
                "nativeWindowsRequired",
            )
        ),
        "unknown path must select full evidence graph",
    )

    override = analyze_changed_paths(["provider_router.go"], "full")
    require(
        override["backendRequired"] and override["rendererRequired"] and override["chromeRequired"],
        "full override must broaden evidence",
    )
    require(override["webkitRequired"], "full override must require WebKit")

    expected_taxonomy = (
        "PRODUCT_FAIL",
        "GATE_TEST_FAIL",
        "CI_HARNESS_FAIL",
        "INFRA_FAIL",
        "EXPECTED_NOOP",
        "SUPERSEDED",
    )
    require(FAILURE_TAXONOMY == expected_taxonomy, "failure taxonomy changed unexpectedly")

    print("DE.PULSE CI impact planner v3 self-test: PASS")
    print("dependency-aware backend/renderer/Chrome/WebKit selection: PASS")
    print("capability evidence-owner path routing: PASS")
    print("active T2 cross-layer assurance routing: PASS")
    print("T9 full cross-platform assurance routing: PASS")
    print("WebKit/browser vs native-macOS ownership separation: PASS")
    print("active release browser-test fail-safe routing: PASS")
    print("process-only portability selection: PASS")
    print("macOS/Windows targeted native rehearsal routing: PASS")
    print("DB + security/data-rights dependency signals: PASS")
    print("unknown-path fail-closed full fallback: PASS")
    print("explicit lane override broadening without native suppression: PASS")
    print("failure taxonomy contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
