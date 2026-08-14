#!/usr/bin/env python3
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent
REGISTRY = ROOT / "functionality_utility_registry.json"
INDEX = ROOT / "renderer" / "index.html"
CONTRACT = ROOT / "adaptive-governance" / "FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md"
DATA_REGISTRY = ROOT / "data_utility_registry.json"

errors = []

if not REGISTRY.exists():
    errors.append("functionality_utility_registry.json missing")
if not CONTRACT.exists():
    errors.append("FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md missing")
if not DATA_REGISTRY.exists():
    errors.append("data_utility_registry.json missing")

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

    sid = str(item.get("surfaceId", "")).strip()
    if sid:
        if sid in surfaces:
            errors.append(f"duplicate surfaceId {sid}: {surfaces[sid]} and {name}")
        surfaces[sid] = name

# Every primary navigation tab must automatically enter the audit registry.
html = INDEX.read_text()
nav_surfaces = set(re.findall(r'data-page="([^"]+)"', html))
for sid in sorted(nav_surfaces):
    if sid not in surfaces:
        errors.append(f"primary navigation surface not audited: {sid}")

# Explicitly require the current canonical high-value engines/checkpoints so
# they cannot silently disappear from the overlap/reuse review.
required_items = {
    "Discovery Scanner",
    "Opportunity Radar",
    "Rapid Move / Market Shock",
    "Pre-Market Prep",
    "Market Open Prep",
    "Earnings & Material Catalyst Reaction Watch",
    "Research",
    "Market Intelligence",
    "Maintenance",
    "Settings",
}
for required in sorted(required_items):
    if required not in names:
        errors.append(f"required audited capability missing: {required}")

# Data-utility remains a companion blocking contract. Functionality may not
# invent a dataset without the data registry carrying purpose/consumer truth.
data_registry = json.loads(DATA_REGISTRY.read_text())
if not data_registry.get("datasets"):
    errors.append("data utility registry has no datasets")

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
    f"{len(items)} audited items · {len(nav_surfaces)} primary navigation surfaces covered"
)
