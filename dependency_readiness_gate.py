#!/usr/bin/env python3
"""DE.PULSE v18.6 external dependency/provider readiness contract."""

from __future__ import annotations

import json
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).resolve().parent
REGISTRY_PATH = ROOT / "dependency_readiness_registry.json"
ACTION_PATH = ROOT / "user_action_required_registry.json"
DOC_IMPACT_PATH = ROOT / "release" / "v18.6.0" / "DOCUMENTATION-IMPACT.md"
ENV_POLICY_PATH = ROOT / "release_environment_policy.json"
FAST_PATH = ROOT / ".github" / "workflows" / "ci-fast.yml"
QUALIFIED_PATH = ROOT / ".github" / "workflows" / "ci-qualified.yml"

REQUIRED_CATEGORIES = {"PROVIDER", "DATABASE", "PACKAGE_RUNTIME", "CREDENTIAL_CONFIG"}
ALLOWED_STATUSES = {
    "READY",
    "READY_NO_CREDENTIAL",
    "RUNTIME_EVALUATED",
    "CONDITIONAL",
    "CI_ENFORCED",
    "RELEASE_ENFORCED",
}
ALLOWED_ACTION_STATES = {"OPEN", "CONDITIONAL", "SATISFIED"}
REQUIRED_EVIDENCE = {
    "provider_capabilities.go",
    "provider_data_rights.go",
    "persistence_backend_select.go",
    "persistence_backend_postgres.go",
    "release_environment_policy.json",
    "app_model.go",
}
MARKET_PROVIDER_OWNER = "Smart Provider Router"


def fail(message: str) -> None:
    raise AssertionError(message)


def load_json(path: pathlib.Path) -> dict:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} is not valid JSON: {exc}")


def validate_gate_ids(values: list[str], label: str) -> None:
    if not values:
        fail(f"{label} must bind at least one existing G0-G16 gate")
    for gate in values:
        if not re.fullmatch(r"G(\d|1[0-6])", str(gate)):
            fail(f"{label} contains invalid/new top-level gate {gate!r}; G0-G16 is immutable")


def evidence_path(raw: str) -> pathlib.Path:
    return ROOT / raw.split("#", 1)[0]


def workflow_go_version(text: str, label: str) -> str:
    match = re.search(r"^\s*GO_VERSION:\s*['\"]([^'\"]+)['\"]\s*$", text, re.MULTILINE)
    if not match:
        fail(f"{label} has no pinned GO_VERSION")
    return match.group(1)


def main() -> int:
    registry = load_json(REGISTRY_PATH)
    actions = load_json(ACTION_PATH)

    if registry.get("schema") != "DE.PULSE-DEPENDENCY-READINESS-1":
        fail("dependency registry schema mismatch")
    if registry.get("policyVersion") != "v18.6.0":
        fail("dependency registry policyVersion must be v18.6.0")
    if "IMPL-17-DEPS-001" not in registry.get("scopeRequirement", ""):
        fail("dependency registry is not bound to IMPL-17-DEPS-001")
    validate_gate_ids(registry.get("gateBindings", []), "dependency registry")

    rows = registry.get("dependencies")
    if not isinstance(rows, list) or not rows:
        fail("dependency registry must contain dependencies")

    ids: set[str] = set()
    categories: set[str] = set()
    all_evidence: set[str] = set()
    counts = {category: 0 for category in REQUIRED_CATEGORIES}

    for row in rows:
        dep_id = str(row.get("id", "")).strip()
        if not dep_id:
            fail("dependency entry missing id")
        if dep_id in ids:
            fail(f"duplicate dependency id {dep_id}")
        ids.add(dep_id)

        for field in ("category", "owner", "capability", "status", "blocker", "userAction", "rightsEntitlement"):
            if not str(row.get(field, "")).strip():
                fail(f"{dep_id} missing required field {field}")

        category = row["category"]
        categories.add(category)
        if category not in REQUIRED_CATEGORIES:
            fail(f"{dep_id} has unsupported category {category}")
        counts[category] += 1
        if row["status"] not in ALLOWED_STATUSES:
            fail(f"{dep_id} has unsupported status {row['status']}")

        evidence = row.get("evidence")
        if not isinstance(evidence, list) or not evidence:
            fail(f"{dep_id} must bind evidence")
        for raw in evidence:
            raw = str(raw).strip()
            if not raw:
                fail(f"{dep_id} has empty evidence path")
            all_evidence.add(raw.split("#", 1)[0])
            if not evidence_path(raw).exists():
                fail(f"{dep_id} evidence does not exist: {raw}")

        if category == "PROVIDER" and row["owner"] == MARKET_PROVIDER_OWNER:
            if "REVIEW_REQUIRED" not in row["rightsEntitlement"]:
                fail(f"{dep_id} market-data rights must remain fail-closed review-required")

    missing_categories = REQUIRED_CATEGORIES - categories
    if missing_categories:
        fail(f"missing dependency categories: {sorted(missing_categories)}")
    if counts["PROVIDER"] < 4:
        fail("provider coverage is unexpectedly narrow")
    if counts["DATABASE"] < 2:
        fail("both desktop SQLite and hosted PostgreSQL must be represented")
    if counts["PACKAGE_RUNTIME"] < 4:
        fail("package/runtime coverage must include build, governance, browser and security tooling")
    if counts["CREDENTIAL_CONFIG"] < 3:
        fail("credential/config coverage is unexpectedly narrow")

    missing_evidence = REQUIRED_EVIDENCE - all_evidence
    if missing_evidence:
        fail(f"canonical evidence coverage missing: {sorted(missing_evidence)}")

    if actions.get("schema") != "DE.PULSE-USER-ACTION-REQUIRED-1":
        fail("user-action registry schema mismatch")
    if actions.get("policyVersion") != "v18.6.0":
        fail("user-action registry policyVersion must be v18.6.0")
    if "IMPL-17-DEPS-001" not in actions.get("scopeRequirement", ""):
        fail("user-action registry is not bound to IMPL-17-DEPS-001")

    action_rows = actions.get("actions")
    if not isinstance(action_rows, list) or not action_rows:
        fail("durable User Action Required register must not be empty")

    action_ids: set[str] = set()
    linked_dependencies: set[str] = set()
    for action in action_rows:
        action_id = str(action.get("id", "")).strip()
        if not action_id:
            fail("user-action entry missing id")
        if action_id in action_ids:
            fail(f"duplicate user-action id {action_id}")
        action_ids.add(action_id)
        if action.get("state") not in ALLOWED_ACTION_STATES:
            fail(f"{action_id} has unsupported state {action.get('state')!r}")
        for field in ("owner", "severity", "condition", "action", "verification"):
            if not str(action.get(field, "")).strip():
                fail(f"{action_id} missing required field {field}")
        dep_ids = action.get("dependencyIds")
        if not isinstance(dep_ids, list) or not dep_ids:
            fail(f"{action_id} must link dependencies")
        for dep_id in dep_ids:
            if dep_id not in ids:
                fail(f"{action_id} links unknown dependency {dep_id}")
            linked_dependencies.add(dep_id)
        validate_gate_ids(action.get("gateBindings", []), action_id)

    if "UA-RIGHTS-001" not in action_ids:
        fail("durable rights-review action is missing")
    if "DEP-DB-POSTGRES" not in linked_dependencies:
        fail("hosted PostgreSQL user action is not durably linked")
    if "DEP-CONFIG-LIVE-EQUITY" not in linked_dependencies:
        fail("live provider credential action is not durably linked")

    docs = DOC_IMPACT_PATH.read_text(encoding="utf-8")
    if "IMPL-17-DEPS-001" not in docs or "MANIFEST_AND_REGISTRY" not in docs:
        fail("Documentation Impact Manifest lacks final dependency-readiness disposition")

    policy = load_json(ENV_POLICY_PATH)
    preferred = str(policy.get("preferred_go_version", "")).strip()
    approved = {str(x).strip() for x in policy.get("approved_go_minor_lines", [])}
    fast_text = FAST_PATH.read_text(encoding="utf-8")
    qualified_text = QUALIFIED_PATH.read_text(encoding="utf-8")
    fast_go = workflow_go_version(fast_text, "CI Fast")
    qualified_go = workflow_go_version(qualified_text, "CI Qualified")
    if not preferred or fast_go != preferred or qualified_go != preferred:
        fail(
            "Go runtime source-of-truth drift: "
            f"policy preferred={preferred!r}, fast={fast_go!r}, qualified={qualified_go!r}"
        )
    minor = ".".join(preferred.split(".")[:2])
    if minor not in approved:
        fail(f"preferred Go version {preferred} is outside approved minor lines {sorted(approved)}")

    for workflow_name, text in (("CI Fast", fast_text), ("CI Qualified", qualified_text)):
        if "dependency_readiness_gate.py" not in text:
            fail(f"{workflow_name} does not bind dependency_readiness_gate.py")

    print(
        "PASS: dependency/provider readiness registry "
        f"{len(rows)} dependencies across {len(categories)} categories; "
        f"{len(action_rows)} durable user actions; G0-G16 bindings only"
    )
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except AssertionError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
