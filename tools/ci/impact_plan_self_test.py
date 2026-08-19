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

    renderer = analyze_changed_paths(["renderer/watchlist-v18.5.1.js"])
    require(renderer["qualifiedLane"] == "full", "renderer change must use full qualification")
    require(renderer["webkitRequired"] is True, "renderer change must require primary WebKit evidence")
    require("RENDERER_UI" in renderer["changeClasses"], "renderer class missing")

    organized_renderer_test = analyze_changed_paths(["tests/renderer/documentation_access_test.js"])
    require(organized_renderer_test["qualifiedLane"] == "full", "organized renderer test must use full qualification")
    require(organized_renderer_test["webkitRequired"] is True, "organized renderer test must retain primary WebKit evidence")
    require("RENDERER_UI" in organized_renderer_test["changeClasses"], "organized renderer test class missing")
    require(organized_renderer_test["nodeRequired"] is True, "organized renderer test must request Node")

    organized_browser_test = analyze_changed_paths(["tests/browser/watchlist_membership_test.py"])
    require(organized_browser_test["qualifiedLane"] == "full", "organized browser test must use full qualification")
    require(organized_browser_test["webkitRequired"] is True, "organized browser test must retain primary WebKit evidence")
    require("RENDERER_UI" in organized_browser_test["changeClasses"], "organized browser test class missing")

    process = analyze_changed_paths([
        ".github/workflows/ci-qualified.yml",
        "tools/ci/workflow_policy.py",
        "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md",
    ])
    require(process["qualifiedLane"] == "ci-harness", "process-only change must use ci-harness")
    require(process["releaseRehearsalRequired"] is True, "CI hardening must require release rehearsal")

    stable_manifest = analyze_changed_paths(["release/v18.6.1/stable-evidence-manifest.json"])
    require(stable_manifest["qualifiedLane"] == "ci-harness", "durable Stable evidence index must stay process-only")
    require("RELEASE_TOOLING" in stable_manifest["changeClasses"], "Stable evidence index must remain release-governed")
    require(stable_manifest["releaseRehearsalRequired"] is True, "Stable evidence changes require release rehearsal")
    require(stable_manifest["webkitRequired"] is False, "Stable evidence-only work must not consume WebKit")

    non_manifest_release = analyze_changed_paths(["release/v18.6.1/run_full_certification.sh"])
    require(non_manifest_release["qualifiedLane"] == "full", "release executable/script changes must remain full qualification")

    webkit_harness = analyze_changed_paths(["tools/ci/webkit_targeted_test.py"])
    require(webkit_harness["qualifiedLane"] == "ci-harness", "WebKit harness-only change must stay process-only")
    require(webkit_harness["webkitRequired"] is True, "WebKit harness changes must execute WebKit proof")

    provider = analyze_changed_paths(["provider_router.go"])
    require("BACKEND" in provider["changeClasses"], "provider Go change must be backend")
    require("PROVIDER_ROUTER" in provider["changeClasses"], "provider router class missing")
    require(provider["webkitRequired"] is False, "backend/provider-only work must not require WebKit")

    rights = analyze_changed_paths(["governance/provider-data-rights-registry.md"])
    require("DATA_RIGHTS" in rights["changeClasses"], "data-rights class missing")
    require("PROVIDER_ROUTER" in rights["changeClasses"], "provider data-rights should include provider-router class")

    mixed = analyze_changed_paths(["tools/ci/workflow_policy.py", "renderer/index.html"])
    require(mixed["qualifiedLane"] == "full", "mixed process/product change must fail closed to full")
    require(mixed["webkitRequired"] is True, "mixed renderer change must retain WebKit requirement")

    expected_taxonomy = (
        "PRODUCT_FAIL",
        "GATE_TEST_FAIL",
        "CI_HARNESS_FAIL",
        "INFRA_FAIL",
        "EXPECTED_NOOP",
        "SUPERSEDED",
    )
    require(FAILURE_TAXONOMY == expected_taxonomy, "failure taxonomy changed unexpectedly")

    print("DE.PULSE CI impact planner v2 self-test: PASS")
    print("process-only routing: PASS")
    print("durable Stable evidence process-only routing: PASS")
    print("release executable/script full-qualification protection: PASS")
    print("mixed/product fail-closed routing: PASS")
    print("Chrome + WebKit primary renderer-risk signal: PASS")
    print("organized renderer/browser test routing: PASS")
    print("WebKit harness self-validation routing: PASS")
    print("backend/provider WebKit suppression: PASS")
    print("failure taxonomy contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
