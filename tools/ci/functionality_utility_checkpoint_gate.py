#!/usr/bin/env python3
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
REGISTRY = ROOT / "functionality_utility_registry.json"
REMEDIATION = ROOT / "functionality_utility_remediation.json"
INDEX = ROOT / "renderer" / "index.html"
CONTRACT = ROOT / "adaptive-governance" / "FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md"
DATA_REGISTRY = ROOT / "data_utility_registry.json"
PROVIDER_MODE_REGISTRY = ROOT / "provider_market_mode_integration_registry.json"
PROVIDER_ROUTER = ROOT / "provider_router.go"
PROVIDER_CAPABILITIES = ROOT / "provider_capabilities.go"

errors = []

for path, label in (
    (REGISTRY, "functionality_utility_registry.json"),
    (REMEDIATION, "functionality_utility_remediation.json"),
    (CONTRACT, "FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md"),
    (DATA_REGISTRY, "data_utility_registry.json"),
    (PROVIDER_MODE_REGISTRY, "provider_market_mode_integration_registry.json"),
    (PROVIDER_ROUTER, "provider_router.go"),
    (PROVIDER_CAPABILITIES, "provider_capabilities.go"),
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

# Every current or planned provider capability must receive a Market Mode
# disposition. This turns the adaptive correlation method into an executable
# conservation rule: adding a provider to the Router/capability registry without
# assessing its Market Mode role fails G2/G10 instead of becoming silent scope.
provider_mode = json.loads(PROVIDER_MODE_REGISTRY.read_text())
if provider_mode.get("schema") != "DE.PULSE-PROVIDER-MARKET-MODE-INTEGRATION-1":
    errors.append("unexpected provider/Market Mode registry schema")
if str(provider_mode.get("version", "")).strip() != "18.5.2":
    errors.append("provider/Market Mode registry version drift")

allowed_dispositions = {
    "INTEGRATED",
    "CONTEXTUAL_ONLY",
    "NOT_RELEVANT",
    "INTENTIONALLY_HIDDEN",
}
allowed_lifecycles = {
    "NOT_IMPLEMENTED",
    "SHADOW",
    "VALIDATED",
    "APPROVED",
    "PRODUCTION",
}
allowed_mode_scopes = {
    "MARKET_WIDE",
    "DAY",
    "SWING",
    "LONG",
    "SECTOR",
    "INDUSTRY",
    "NONE",
}
assessments = provider_mode.get("assessments", [])
if not assessments:
    errors.append("provider/Market Mode registry contains no assessments")
if set(provider_mode.get("allowedDispositions", [])) != allowed_dispositions:
    errors.append("provider/Market Mode disposition vocabulary drift")
if set(provider_mode.get("allowedLifecycles", [])) != allowed_lifecycles:
    errors.append("provider/Market Mode lifecycle vocabulary drift")
if set(provider_mode.get("allowedMarketModeScopes", [])) != allowed_mode_scopes:
    errors.append("provider/Market Mode scope vocabulary drift")
if provider_mode.get("promotion") != "SHADOW -> VALIDATED -> APPROVED -> PRODUCTION":
    errors.append("provider/Market Mode promotion lifecycle drift")

assessment_ids = set()
assessed_providers = set()
for assessment in assessments:
    ident = str(assessment.get("id", "")).strip()
    provider = str(assessment.get("provider", "")).strip()
    capability = str(assessment.get("capability", "")).strip()
    if not ident:
        errors.append("provider/Market Mode assessment missing id")
        continue
    if ident in assessment_ids:
        errors.append(f"duplicate provider/Market Mode assessment id: {ident}")
    assessment_ids.add(ident)
    if not provider:
        errors.append(f"{ident}: provider missing")
    else:
        assessed_providers.add(provider)
    for field in (
        "capability",
        "requirementId",
        "routerRole",
        "marketModeDisposition",
        "lifecycle",
        "productionInfluence",
        "canonicalOwner",
        "evidencePolicy",
        "acceptance",
    ):
        if not str(assessment.get(field, "")).strip():
            errors.append(f"{ident}: {field} missing")
    if not assessment.get("consumers"):
        errors.append(f"{ident}: consumers missing")

    disposition = str(assessment.get("marketModeDisposition", ""))
    lifecycle = str(assessment.get("lifecycle", ""))
    influence = str(assessment.get("productionInfluence", ""))
    scopes = assessment.get("marketModeScopes", [])
    if disposition not in allowed_dispositions:
        errors.append(f"{ident}: invalid Market Mode disposition {disposition}")
    if lifecycle not in allowed_lifecycles:
        errors.append(f"{ident}: invalid lifecycle {lifecycle}")
    if not isinstance(scopes, list) or not scopes:
        errors.append(f"{ident}: marketModeScopes missing")
        scopes = []
    invalid_scopes = sorted(set(scopes) - allowed_mode_scopes)
    if invalid_scopes:
        errors.append(f"{ident}: invalid Market Mode scopes {invalid_scopes}")
    if "NONE" in scopes and len(scopes) != 1:
        errors.append(f"{ident}: NONE cannot be combined with active Market Mode scopes")
    if disposition == "INTEGRATED" and "NONE" in scopes:
        errors.append(f"{ident}: INTEGRATED requires at least one active Market Mode scope")
    if disposition in {"NOT_RELEVANT", "INTENTIONALLY_HIDDEN"} and scopes != ["NONE"]:
        errors.append(f"{ident}: {disposition} must use the single NONE scope")
    if lifecycle in {"NOT_IMPLEMENTED", "SHADOW", "VALIDATED"} and influence != "NONE":
        errors.append(f"{ident}: pre-approval lifecycle cannot influence production Market Modes")

router_text = PROVIDER_ROUTER.read_text()
route_start = router_text.find("func routeChains()")
route_end = router_text.find("\n}\n", route_start)
if route_start < 0 or route_end < 0:
    errors.append("cannot locate canonical provider routeChains")
    route_block = ""
else:
    route_block = router_text[route_start:route_end]
route_providers = set()
for values in re.findall(r":\s*\{([^}]*)\}", route_block):
    route_providers.update(re.findall(r'"([^"]+)"', values))
capability_providers = set(
    re.findall(r'Provider:\s*"([^"]+)"', PROVIDER_CAPABILITIES.read_text())
)
for provider in sorted(route_providers | capability_providers):
    if provider not in assessed_providers:
        errors.append(f"provider lacks Market Mode assessment: {provider}")

tradeinsight_expected = {
    "Congressional Trading Intelligence",
    "SEC Form 4 Enrichment",
    "Historical OHLCV Backfill",
    "Corporate Actions",
    "Symbol Metadata / Search",
    "Top Movers",
    "Controlled AI/MCP Research",
}
tradeinsight_rows = [row for row in assessments if row.get("provider") == "TradeInsight"]
tradeinsight_capabilities = {row.get("capability") for row in tradeinsight_rows}
if tradeinsight_capabilities != tradeinsight_expected:
    errors.append(
        "TradeInsight Market Mode capability coverage drift: "
        + repr(sorted(tradeinsight_capabilities))
    )
for row in tradeinsight_rows:
    ident = row.get("id", "TradeInsight")
    if row.get("requirementId") != "IMPL-18-TRADEINSIGHT-001":
        errors.append(f"{ident}: TradeInsight assessment must bind IMPL-18-TRADEINSIGHT-001")
    if row.get("lifecycle") != "NOT_IMPLEMENTED":
        errors.append(f"{ident}: current source cannot claim implemented TradeInsight lifecycle")
    if row.get("productionInfluence") != "NONE":
        errors.append(f"{ident}: TradeInsight cannot influence production before implementation/approval")

contract_text = CONTRACT.read_text()
for phrase in (
    "Need → reuse → correlate → consolidate → place → measure.",
    "one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.",
    "G9 — Cross-Module / UI / UX",
    "G10 — Pre-Freeze Qualification",
    "Provider → Market Mode assessment",
    "SHADOW → VALIDATED → APPROVED → PRODUCTION",
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
    f"{len(remediation_by_name)} durable remediation plans · {len(required_data)} v18.2 identity/admin data classes governed · "
    f"{len(assessments)} provider/Market Mode assessments across {len(assessed_providers)} providers"
)

