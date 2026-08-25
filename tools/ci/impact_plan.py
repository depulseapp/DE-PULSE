#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import subprocess

PROCESS_ONLY_PREFIXES = (
    ".github/workflows/",
    "tools/ci/",
    "tools/release/",
    "adaptive-governance/",
    "governance/",
    "handoff/",
    ".depulse-certification/resume/",
)
PROCESS_ONLY_EXACT = {
    "tools/release/source_fingerprint.py",
    "README.md",
    "AGENTS.md",
    "CLAUDE.md",
}
STABLE_EVIDENCE_RE = re.compile(r"^release/v[^/]+/stable-evidence-manifest\.json$")
ACTIVE_RELEASE_BROWSER_TEST_RE = re.compile(r"^release/v[^/]+/browser_[^/]+_test\.py$")

FAILURE_TAXONOMY = (
    "PRODUCT_FAIL",
    "GATE_TEST_FAIL",
    "CI_HARNESS_FAIL",
    "INFRA_FAIL",
    "EXPECTED_NOOP",
    "SUPERSEDED",
)

CHANGE_CLASSES = (
    "CI_HARNESS",
    "RELEASE_TOOLING",
    "BACKEND",
    "RENDERER_UI",
    "AUTH_SECURITY",
    "PROVIDER_ROUTER",
    "DATA_RIGHTS",
    "PERSISTENCE",
    "RELIABILITY_PERFORMANCE",
    "CERTIFICATION_GOVERNANCE",
)

RENDERER_EVIDENCE_FILES = {
    "tools/ci/deterministic_equivalence_test.js",
}
WEBKIT_EVIDENCE_FILES = {
    "tools/ci/webkit_browser_test.py",
    "tools/ci/webkit_targeted_test.py",
    "tools/ci/browser_risk_routing_gate.py",
}
CAPABILITY_EVIDENCE_FILES = RENDERER_EVIDENCE_FILES | {
    "tools/ci/webkit_browser_test.py",
}
# Active cross-layer assurance controls are not ordinary process-only files.
# A change to their ownership/closure semantics must re-qualify the executable
# backend + renderer evidence graph they certify, without inventing a new lane.
T2_CROSS_LAYER_ASSURANCE_FILES = {
    "tools/ci/v18_t2_unit_contract_assurance_gate.py",
    "governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T2_UNIT_CONTRACT_ASSURANCE.json",
}
NATIVE_MACOS_FILES = {
    "tools/release/native_macos.sh",
    "desktop_lifecycle.go",
    "desktop_lifecycle_regression_test.go",
}
NATIVE_WINDOWS_FILES = {
    "tools/release/native_windows.ps1",
}
SHARED_NATIVE_RELEASE_FILES = {
    ".github/workflows/release.yml",
    ".github/workflows/ci-qualified.yml",
    "release_identity.json",
    "tools/release/release_identity.py",
    "tools/release/source_fingerprint.py",
    "tools/release/g15_assurance.py",
    "tools/release/verify_promotion_evidence.py",
}


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), check=check, text=True, capture_output=True)


def resolve_commit(ref: str) -> str:
    candidate = ref.strip()
    if not candidate:
        return ""
    probes = (candidate, f"origin/{candidate}") if not candidate.startswith("origin/") else (candidate,)
    for probe in probes:
        result = git("rev-parse", "--verify", f"{probe}^{{commit}}", check=False)
        if result.returncode == 0:
            return result.stdout.strip()
    return ""


def resolve_base(base: str, head: str, target_ref: str) -> tuple[str, str]:
    explicit = resolve_commit(base)
    if explicit:
        merge = git("merge-base", explicit, head, check=False)
        if merge.returncode != 0 or not merge.stdout.strip():
            raise SystemExit(f"unable to derive merge-base for explicit base {base!r} and head {head}")
        return merge.stdout.strip(), "EXPLICIT_BASE_MERGE_BASE"

    target = resolve_commit(target_ref)
    if not target:
        raise SystemExit(
            "trustworthy CI base unavailable: provide --base or a resolvable --target-ref; HEAD^ fallback is prohibited"
        )
    merge = git("merge-base", target, head, check=False)
    if merge.returncode != 0 or not merge.stdout.strip():
        raise SystemExit(f"unable to derive merge-base for target {target_ref!r} and head {head}")
    return merge.stdout.strip(), "TARGET_REF_MERGE_BASE"


def is_process_only(path: str) -> bool:
    if path in CAPABILITY_EVIDENCE_FILES or path in T2_CROSS_LAYER_ASSURANCE_FILES:
        return False
    if STABLE_EVIDENCE_RE.fullmatch(path):
        return True
    return path in PROCESS_ONLY_EXACT or path.startswith(PROCESS_ONLY_PREFIXES)


def explicit_classes(path: str) -> set[str]:
    p = path.lower()
    classes: set[str] = set()

    if path.startswith(".github/workflows/") or path.startswith("tools/ci/"):
        classes.add("CI_HARNESS")
    if (
        path.startswith("tools/release/")
        or path == ".github/workflows/release.yml"
        or path.startswith("release/")
        or path == "tools/release/release_identity.py"
        or path == "release_identity.json"
        or path == "tools/release/version_consistency_test.py"
    ):
        classes.add("RELEASE_TOOLING")
    if path.startswith(("adaptive-governance/", "governance/", "handoff/", ".depulse-certification/")):
        classes.add("CERTIFICATION_GOVERNANCE")

    if (
        path.startswith(("renderer/", "tests/renderer/", "tests/browser/"))
        or path.endswith((".html", ".css"))
        or path in RENDERER_EVIDENCE_FILES
        or ACTIVE_RELEASE_BROWSER_TEST_RE.fullmatch(path)
    ):
        classes.add("RENDERER_UI")
    if path.endswith(".go") or path in {"go.mod", "go.sum"}:
        classes.add("BACKEND")

    if path in T2_CROSS_LAYER_ASSURANCE_FILES:
        classes.update({"BACKEND", "RENDERER_UI"})

    if any(token in p for token in ("auth", "login", "rbac", "security", "secret", "credential", "permission")):
        classes.add("AUTH_SECURITY")
    if any(token in p for token in ("provider", "router", "finnhub", "alpaca", "tradeinsight", "twelve", "fred", "bls", "eia")):
        classes.add("PROVIDER_ROUTER")
    if any(token in p for token in ("license", "licence", "entitlement", "data_right", "data-right", "redistribution", "ai_use", "ai-use")):
        classes.add("DATA_RIGHTS")
    if any(token in p for token in ("sqlite", "migration", "persist", "storage", "cache", "canonical_state")):
        classes.add("PERSISTENCE")
    if any(token in p for token in ("performance", "load", "runtime", "backpressure", "circuit", "retry", "latency", "stability", "reliability")):
        classes.add("RELIABILITY_PERFORMANCE")
    return classes


def classify_path(path: str) -> set[str]:
    classes = explicit_classes(path)
    if not classes and not is_process_only(path):
        classes.add("BACKEND")
    return classes


def adaptive_requirements(changed: list[str]) -> dict[str, object]:
    process_only = bool(changed) and all(is_process_only(path) for path in changed)
    explicit_by_path = {path: explicit_classes(path) for path in changed}
    unknown_paths = sorted(path for path, classes in explicit_by_path.items() if not classes and not is_process_only(path))
    classes = sorted({c for path in changed for c in classify_path(path)})
    class_set = set(classes)

    fail_closed_full = bool(unknown_paths)
    backend_required = bool(class_set & {"BACKEND", "PROVIDER_ROUTER", "PERSISTENCE", "RELIABILITY_PERFORMANCE", "AUTH_SECURITY"})
    renderer_required = "RENDERER_UI" in class_set
    chrome_required = renderer_required
    webkit_required = renderer_required or any(path in WEBKIT_EVIDENCE_FILES for path in changed)
    security_rights_required = bool(class_set & {"AUTH_SECURITY", "DATA_RIGHTS"})
    db_integration_required = "PERSISTENCE" in class_set
    portability_required = process_only and bool(class_set & {"CI_HARNESS", "RELEASE_TOOLING", "CERTIFICATION_GOVERNANCE"})

    macos_specific = any(path in NATIVE_MACOS_FILES or "macos" in path.lower() for path in changed)
    windows_specific = any(path in NATIVE_WINDOWS_FILES or "windows" in path.lower() for path in changed)
    shared_native = any(path in SHARED_NATIVE_RELEASE_FILES for path in changed)
    native_macos_required = macos_specific or shared_native
    native_windows_required = windows_specific or shared_native
    release_rehearsal_required = native_macos_required or native_windows_required

    if fail_closed_full:
        backend_required = renderer_required = chrome_required = webkit_required = True
        security_rights_required = db_integration_required = True
        native_macos_required = native_windows_required = True
        release_rehearsal_required = True

    if process_only:
        qualified_lane = "ci-harness"
    elif renderer_required and not backend_required and not security_rights_required and not db_integration_required:
        qualified_lane = "renderer"
    elif (chrome_required or webkit_required) and not backend_required and not security_rights_required and not db_integration_required:
        qualified_lane = "browser"
    elif backend_required and not renderer_required and not chrome_required and not webkit_required:
        qualified_lane = "backend"
    else:
        qualified_lane = "full"

    return {
        "processOnly": process_only,
        "goRequired": any(path.endswith(".go") or path in {"go.mod", "go.sum"} for path in changed),
        "nodeRequired": any(path.endswith((".js", ".mjs", ".cjs")) for path in changed),
        "qualifiedLane": qualified_lane,
        "changeClasses": classes,
        "unknownPaths": unknown_paths,
        "failClosedFull": fail_closed_full,
        "backendRequired": backend_required,
        "rendererRequired": renderer_required,
        "chromeRequired": chrome_required,
        "webkitRequired": webkit_required,
        "securityRightsRequired": security_rights_required,
        "dbIntegrationRequired": db_integration_required,
        "portabilityRequired": portability_required,
        "nativeMacosRequired": native_macos_required,
        "nativeWindowsRequired": native_windows_required,
        "releaseRehearsalRequired": release_rehearsal_required,
        "failureTaxonomyVersion": 1,
        "failureTaxonomy": list(FAILURE_TAXONOMY),
    }


def apply_lane_override(plan: dict[str, object], requested: str) -> dict[str, object]:
    requested = (requested or "adaptive").strip()
    valid = {"adaptive", "full", "backend", "renderer", "browser", "ci-harness"}
    if requested not in valid:
        raise SystemExit(f"invalid qualification lane: {requested}")
    if requested == "adaptive":
        plan["qualificationMode"] = "adaptive"
        return plan

    plan["qualificationMode"] = requested
    plan["qualifiedLane"] = requested
    plan["backendRequired"] = requested in {"full", "backend"}
    plan["rendererRequired"] = requested in {"full", "renderer"}
    plan["chromeRequired"] = requested in {"full", "browser"}
    plan["webkitRequired"] = requested in {"full", "browser"} or bool(plan["webkitRequired"])
    plan["portabilityRequired"] = requested == "ci-harness"
    # Explicit lane overrides may broaden ordinary qualification, but they never
    # suppress native rehearsals selected by the changed-path dependency graph.
    return plan


def analyze_changed_paths(changed: list[str], requested_lane: str = "adaptive") -> dict[str, object]:
    return apply_lane_override(adaptive_requirements(changed), requested_lane)


def main() -> int:
    parser = argparse.ArgumentParser(description="Plan DE.PULSE CI jobs from a trustworthy full Git delta")
    parser.add_argument("--base", default="")
    parser.add_argument("--target-ref", default="main")
    parser.add_argument("--head", default="HEAD")
    parser.add_argument("--requested-lane", default="adaptive")
    parser.add_argument("--github-output")
    parser.add_argument("--json-out")
    args = parser.parse_args()

    head = git("rev-parse", args.head).stdout.strip()
    base, base_resolution = resolve_base(args.base, head, args.target_ref)
    raw = git("diff", "--name-only", base, head).stdout
    changed = sorted({line.strip() for line in raw.splitlines() if line.strip()})

    analysis = analyze_changed_paths(changed, args.requested_lane)
    selected = [
        key
        for key, field in (
            ("ci-harness", None),
            ("portability", "portabilityRequired"),
            ("backend", "backendRequired"),
            ("renderer", "rendererRequired"),
            ("chrome", "chromeRequired"),
            ("webkit", "webkitRequired"),
            ("security-rights", "securityRightsRequired"),
            ("db-integration", "dbIntegrationRequired"),
            ("native-macos", "nativeMacosRequired"),
            ("native-windows", "nativeWindowsRequired"),
        )
        if field is None or bool(analysis[field])
    ]
    reason = (
        "Unknown/unclassified changed path detected; fail-closed full evidence plus native rehearsals selected."
        if analysis["failClosedFull"]
        else "Deterministic dependency-aware job selection from the full merge-base delta."
    )
    plan = {
        "schema": "DE.PULSE-CI-IMPACT-PLAN-3",
        "baseSha": base,
        "baseResolution": base_resolution,
        "targetRef": args.target_ref,
        "headSha": head,
        "changedPaths": changed,
        **analysis,
        "selectedJobs": selected,
        "reason": reason,
    }
    text = json.dumps(plan, indent=2, sort_keys=True) + "\n"
    print(text, end="")

    if args.json_out:
        out = Path(args.json_out)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(text, encoding="utf-8")
    if args.github_output:
        out = Path(args.github_output)
        mapping = {
            "qualified_lane": "qualifiedLane",
            "go_required": "goRequired",
            "node_required": "nodeRequired",
            "process_only": "processOnly",
            "backend_required": "backendRequired",
            "renderer_required": "rendererRequired",
            "chrome_required": "chromeRequired",
            "webkit_required": "webkitRequired",
            "security_rights_required": "securityRightsRequired",
            "db_integration_required": "dbIntegrationRequired",
            "portability_required": "portabilityRequired",
            "native_macos_required": "nativeMacosRequired",
            "native_windows_required": "nativeWindowsRequired",
            "release_rehearsal_required": "releaseRehearsalRequired",
            "fail_closed_full": "failClosedFull",
        }
        with out.open("a", encoding="utf-8") as f:
            for output, field in mapping.items():
                value = analysis[field]
                if isinstance(value, bool):
                    value = str(value).lower()
                f.write(f"{output}={value}\n")
            f.write(f"change_classes={','.join(analysis['changeClasses'])}\n")
            f.write(f"selected_jobs={','.join(selected)}\n")
            f.write(f"base_sha={base}\n")
            f.write(f"base_resolution={base_resolution}\n")
            f.write("impact_plan_schema=DE.PULSE-CI-IMPACT-PLAN-3\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
