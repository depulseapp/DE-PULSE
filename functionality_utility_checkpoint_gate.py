#!/usr/bin/env python3
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent
REGISTRY = ROOT / "functionality_utility_registry.json"
REMEDIATION = ROOT / "functionality_utility_remediation.json"
INDEX = ROOT / "renderer" / "index.html"
CONTRACT = ROOT / "adaptive-governance" / "FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md"
DATA_REGISTRY = ROOT / "data_utility_registry.json"

errors = []

for path, label in (
    (REGISTRY, "functionality_utility_registry.json"),
    (REMEDIATION, "functionality_utility_remediation.json"),
    (CONTRACT, "FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md"),
    (DATA_REGISTRY, "data_utility_registry.json"),
):
    if not path.exists():
        errors.append(label + " missing")

if errors:
    print("Functionality Utility & Integration Checkpoint: FAIL")
    for err in errors:
        print(" -", err)
    raise SystemExit(2)

registry = json.loads(REGISTRY.read_text())
items = registry.get("items", [])
if registry.get("schema") != "DE.PULSE-FUNCTIONALITY-UTILITY-1":
    errors.append("unexpected registry schema")
if not items:
    errors.append("registry contains no functionality items")

remediation = json.loads(REMEDIATION.read_text())
if remediation.get("schema") != "DE.PULSE-FUNCTIONALITY-REMEDIATION-1":
    errors.append("unexpected remediation schema")
if not str(remediation.get("targetRelease", "")).strip():
    errors.append("remediation targetRelease missing")
remediation_by_name = {
    str(x.get("name", "")).strip(): x for x in remediation.get("items", []) if str(x.get("name", "")).strip()
}

allowed_ui = {
    "EXISTING_SURFACE",
    "EMBED_EXISTING_SURFACE",
    "SUMMARY_PLUS_MAINTENANCE_DETAIL",
    "INTERNAL_OR_DRILLDOWN",
    "REMOVE_OR_REDIRECT",
    "NEW_SURFACE_JUSTIFIED",
}
allowed_dup = {
    "CONSOLIDATE_PRESENTATION",
    "CONSOLIDATE_ACQUISITION",
    "RETAIN_SEPARATE",
    "CANONICAL_DEEP_HOME",
    "EXTEND_EXISTING_COORDINATOR",
    "EXTEND_EXISTING_OWNER",
    "REMOVE_PROMINENT_SURFACE",
    "REMOVE_OR_RETIRE",
}
allowed_status = {"ACTIVE", "CONSOLIDATE", "RETIRE", "DEFER"}
actionable_dup = {
    "CONSOLIDATE_PRESENTATION",
    "CONSOLIDATE_ACQUISITION",
    "EXTEND_EXISTING_COORDINATOR",
    "EXTEND_EXISTING_OWNER",
    "REMOVE_PROMINENT_SURFACE",
    "REMOVE_OR_RETIRE",
}

names = set()
surfaces = {}
for item in items:
    name = str(item.get("name", "")).strip()
    if not name:
        errors.append("item missing name")
        continue
    if name in names:
        errors.append(f"duplicate item name: {name}")
    names.add(name)

    for field in ("kind", "owner", "purpose", "reusePolicy", "uiDecision", "duplicationDecision", "separationReason", "status"):
        if not str(item.get(field, "")).strip():
            errors.append(f"{name}: {field} missing")
    if not item.get("consumers"):
        errors.append(f"{name}: consumers missing")
    if not item.get("correlationTargets"):
        errors.append(f"{name}: correlationTargets missing")

    ui = str(item.get("uiDecision", ""))
    dup = str(item.get("duplicationDecision", ""))
    status = str(item.get("status", ""))
    if ui not in allowed_ui:
        errors.append(f"{name}: invalid uiDecision {ui}")
    if dup not in allowed_dup:
        errors.append(f"{name}: invalid duplicationDecision {dup}")
    if status not in allowed_status:
        errors.append(f"{name}: invalid status {status}")
    if ui == "NEW_SURFACE_JUSTIFIED" and not str(item.get("newTabJustification", "")).strip():
        errors.append(f"{name}: NEW_SURFACE_JUSTIFIED requires newTabJustification")

    # Any audit result that calls for consolidation, extension or removal must
    # have a durable next-release remediation record. This prevents audit debt
    # from being acknowledged once and then silently forgotten.
    if dup in actionable_dup:
        plan = remediation_by_name.get(name)
        if not plan:
            errors.append(f"{name}: actionable {dup} finding lacks durable remediation plan")
        else:
            if str(plan.get("action", "")) != dup:
                errors.append(f"{name}: remediation action does not match registry disposition {dup}")
            for field in ("priority", "plan", "acceptance"):
                if not str(plan.get(field, "")).strip():
                    errors.append(f"{name}: remediation {field} missing")

    sid = str(item.get("surfaceId", "")).strip()
    if sid:
        if sid in surfaces:
            errors.append(f"duplicate surfaceId {sid}: {surfaces[sid]} and {name}")
        surfaces[sid] = name

# Every statically declared primary navigation tab must automatically enter the
# audit registry. Conditional/dynamic surfaces are additionally required below.
html = INDEX.read_text()
nav_surfaces = set(re.findall(r'data-page="([^"]+)"', html))
for sid in sorted(nav_surfaces):
    if sid not in surfaces:
        errors.append(f"primary navigation surface not audited: {sid}")

# Explicitly require canonical high-value engines/checkpoints and the v18.2
# conditional Administration surface so they cannot silently disappear from
# overlap/reuse/placement review merely because they are dynamically rendered.
required_items = {
    "Discovery Scanner",
    "Opportunity Radar",
    "Rapid Move / Market Shock",
    "Pre-Market Prep",
    "Market Open Prep",
    "Earnings & Material Catalyst Reaction Watch",
    "Research",
    "Market Intelligence",
    "Administration",
    "Maintenance",
    "Settings",
}
for required in sorted(required_items):
    if required not in names:
        errors.append(f"required audited capability missing: {required}")

# Data Utility is a companion blocking contract. Current-release identity,
# session, presence and delegated-capability data must be explicitly owned and
# governed; administrative audit evidence is strategic until its durable audit
# repository is activated.
data_registry = json.loads(DATA_REGISTRY.read_text())
datasets = data_registry.get("datasets", [])
if not datasets:
    errors.append("data utility registry has no datasets")
if str(data_registry.get("version", "")) != "18.2.0":
    errors.append("data utility registry must be reconciled to v18.2.0")

data_by_name = {str(d.get("dataset", "")).strip(): d for d in datasets}
required_data = {
    "Identity Accounts",
    "Authenticated Sessions",
    "Presence State",
    "Delegated Administrative Capabilities",
    "Administrative Audit Evidence",
}
for dataset in sorted(required_data):
    item = data_by_name.get(dataset)
    if not item:
        errors.append(f"required v18.2 data-utility entry missing: {dataset}")
        continue
    for field in ("owner", "purpose", "retention", "utility", "sensitivity", "authorization"):
        if not str(item.get(field, "")).strip():
            errors.append(f"{dataset}: {field} missing")
    if not item.get("consumers"):
        errors.append(f"{dataset}: consumers missing")

contract_text = CONTRACT.read_text()
for phrase in (
    "Need → reuse → correlate → consolidate → place → measure.",
    "one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.",
    "G9 — Cross-Module / UI / UX",
    "G10 — Pre-Freeze Qualification",
):
    if phrase not in contract_text:
        errors.append(f"checkpoint contract drift: missing phrase {phrase}")

if errors:
    print("Functionality Utility & Integration Checkpoint: FAIL")
    for err in errors:
        print(" -", err)
    raise SystemExit(2)

print(
    "Functionality Utility & Integration Checkpoint: PASS · "
    f"{len(items)} audited items · {len(nav_surfaces)} static primary navigation surfaces covered · "
    f"{len(remediation_by_name)} durable remediation plans · {len(required_data)} v18.2 identity/admin data classes governed"
)
