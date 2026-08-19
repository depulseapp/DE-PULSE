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
    "source_fingerprint.py",
    "README.md",
    "AGENTS.md",
    "CLAUDE.md",
}
STABLE_EVIDENCE_RE = re.compile(r"^release/v[^/]+/stable-evidence-manifest\.json$")

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

WEBKIT_EVIDENCE_FILES = {
    "tools/ci/webkit_targeted_test.py",
    "tools/ci/browser_risk_routing_gate.py",
}


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), check=check, text=True, capture_output=True)


def resolve_base(base: str, head: str) -> str:
    candidate = base.strip()
    if candidate and set(candidate) != {"0"} and git("cat-file", "-e", f"{candidate}^{{commit}}", check=False).returncode == 0:
        return candidate
    parent = git("rev-parse", f"{head}^", check=False)
    if parent.returncode == 0:
        return parent.stdout.strip()
    return head


def is_process_only(path: str) -> bool:
    if STABLE_EVIDENCE_RE.fullmatch(path):
        return True
    return path in PROCESS_ONLY_EXACT or path.startswith(PROCESS_ONLY_PREFIXES)


def classify_path(path: str) -> set[str]:
    p = path.lower()
    classes: set[str] = set()

    if path.startswith(".github/workflows/") or path.startswith("tools/ci/"):
        classes.add("CI_HARNESS")
    if (
        path.startswith("tools/release/")
        or path == ".github/workflows/release.yml"
        or path.startswith("release/")
        or path == "release_identity.py"
        or path == "release_identity.json"
        or path == "version_consistency_test.py"
    ):
        classes.add("RELEASE_TOOLING")
    if path.startswith(("adaptive-governance/", "governance/", "handoff/", ".depulse-certification/")):
        classes.add("CERTIFICATION_GOVERNANCE")

    if path.startswith(("renderer/", "tests/renderer/", "tests/browser/")) or path.endswith((".html", ".css")):
        classes.add("RENDERER_UI")
    if path.endswith(".go") or path in {"go.mod", "go.sum"}:
        classes.add("BACKEND")

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

    # Unknown non-process files are product-affecting by default. This fail-closed
    # behavior keeps the full Qualified lane whenever classification is uncertain.
    if not classes and not is_process_only(path):
        classes.add("BACKEND")
    return classes


def analyze_changed_paths(changed: list[str]) -> dict[str, object]:
    process_only = bool(changed) and all(is_process_only(path) for path in changed)
    go_required = any(path.endswith(".go") or path in {"go.mod", "go.sum"} for path in changed)
    node_required = any(path.endswith((".js", ".mjs", ".cjs")) for path in changed)
    classes = sorted({c for path in changed for c in classify_path(path)})
    qualified_lane = "ci-harness" if process_only else "full"
    release_rehearsal_required = bool({"CI_HARNESS", "RELEASE_TOOLING"} & set(classes))
    webkit_required = "RENDERER_UI" in classes or any(path in WEBKIT_EVIDENCE_FILES for path in changed)

    return {
        "processOnly": process_only,
        "goRequired": go_required,
        "nodeRequired": node_required,
        "qualifiedLane": qualified_lane,
        "changeClasses": classes,
        "releaseRehearsalRequired": release_rehearsal_required,
        "webkitRequired": webkit_required,
        "failureTaxonomyVersion": 1,
        "failureTaxonomy": list(FAILURE_TAXONOMY),
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Plan DE.PULSE CI lanes from a deterministic Git diff")
    parser.add_argument("--base", default="")
    parser.add_argument("--head", default="HEAD")
    parser.add_argument("--github-output")
    parser.add_argument("--json-out")
    args = parser.parse_args()

    head = git("rev-parse", args.head).stdout.strip()
    base = resolve_base(args.base, head)
    raw = git("diff", "--name-only", base, head).stdout
    changed = sorted({line.strip() for line in raw.splitlines() if line.strip()})

    analysis = analyze_changed_paths(changed)
    plan = {
        "schema": "DE.PULSE-CI-IMPACT-PLAN-2",
        "baseSha": base,
        "headSha": head,
        "changedPaths": changed,
        **analysis,
        "reason": (
            "Only canonical CI/release/governance/handoff tooling changed; run harness + portability only."
            if analysis["processOnly"]
            else "Product, test, dependency, release identity, or other non-process content changed; run full qualified coverage."
        ),
    }
    text = json.dumps(plan, indent=2, sort_keys=True) + "\n"
    print(text, end="")

    if args.json_out:
        out = Path(args.json_out)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(text, encoding="utf-8")
    if args.github_output:
        out = Path(args.github_output)
        with out.open("a", encoding="utf-8") as f:
            f.write(f"qualified_lane={analysis['qualifiedLane']}\n")
            f.write(f"go_required={str(analysis['goRequired']).lower()}\n")
            f.write(f"node_required={str(analysis['nodeRequired']).lower()}\n")
            f.write(f"process_only={str(analysis['processOnly']).lower()}\n")
            f.write(f"change_classes={','.join(analysis['changeClasses'])}\n")
            f.write(f"release_rehearsal_required={str(analysis['releaseRehearsalRequired']).lower()}\n")
            f.write(f"webkit_required={str(analysis['webkitRequired']).lower()}\n")
            f.write("impact_plan_schema=DE.PULSE-CI-IMPACT-PLAN-2\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
