#!/usr/bin/env python3
"""Executable all-provider Data Health classification contract.

#80 extends the pre-existing Adaptive Data Health gate instead of creating a new
CI family. The gate fails closed when provider registry/Router membership or a
production external network host appears without durable classification.
"""
from __future__ import annotations

import json
from pathlib import Path
import re
import sys
from urllib.parse import urlsplit

ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = ROOT / "governance/policies/data_health_policy.json"
MATRIX_PATH = ROOT / "governance/data-health/provider-capability-matrix.json"
SLO_PATH = ROOT / "governance/data-health/data-health-slo.json"
FETCH_PATH = ROOT / "governance/data-health/provider-fetch-paths.json"
ADAPTIVE_CURRENT_CONTRACTS = {
    ROOT / "adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md": (
        "#80", "#81", "#82", "#83", "#78", "#84",
        "Smart Provider Router v2", "PARTIAL COVERAGE", "DATA DEGRADED",
    ),
    ROOT / "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md": (
        "provider-capability-matrix.json", "data-health-slo.json",
        "provider-fetch-paths.json", "Adaptive Roadmap", "Build Plan",
        "Build Process", "Delivery Process",
    ),
    ROOT / "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md": (
        "Smart Provider Router v2", "fail closed", "canonical freshness",
        "#81/#82/#83/#78/#84",
    ),
    ROOT / "adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md": (
        "canonical Fast exact-head PASS", "Qualified exact-head PASS",
        "Direct SEC/EDGAR", "No Execution",
    ),
}

EXPECTED_SCHEMAS = {
    MATRIX_PATH: "DE.PULSE-DATA-HEALTH-PROVIDER-CAPABILITY-MATRIX-1",
    SLO_PATH: "DE.PULSE-DATA-HEALTH-SLO-1",
    FETCH_PATH: "DE.PULSE-DATA-HEALTH-PROVIDER-FETCH-PATHS-1",
}
ALLOWED_DISPOSITIONS = {"MIGRATE", "DIRECT_AUTHORITY", "N/A"}
REQUIRED_METRICS = {
    "healthyCoverage", "freshnessCompliance", "fallbackSuccess",
    "degradationFrequency", "degradationDuration", "recoveryTime",
    "quotaPressure", "consumerImpact", "providerCallsAvoided",
}
EXCLUDED_SOURCE_PREFIXES = (
    "release/", "tools/", "tests/", "governance/", "adaptive-governance/",
    "handoff/", ".github/", "assets/", "docs/", "internal/vendorcrypto/",
)
# Browser CSP allowlists describe destinations the renderer may contact; they are
# not backend provider transports and must not create false provider-bypass debt.
NON_BACKEND_NETWORK_LITERAL_FILES = {"http_security.go"}
LOCAL_OR_EXAMPLE_HOSTS = {"localhost", "127.0.0.1", "0.0.0.0", "example.com"}
# Router v2 historically names the SEC capability member "SEC" while the durable
# provider identity and direct authority is "SEC EDGAR". Normalize the alias
# rather than creating a false duplicate provider row.
ROUTER_PROVIDER_ALIASES = {"SEC": "SEC EDGAR"}


def load_json(path: Path, errors: list[str]) -> dict:
    if not path.is_file():
        errors.append(f"missing contract: {path.relative_to(ROOT)}")
        return {}
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        errors.append(f"invalid JSON {path.relative_to(ROOT)}: {exc}")
        return {}
    if value.get("schema") != EXPECTED_SCHEMAS.get(path, value.get("schema")):
        errors.append(f"unsupported schema in {path.relative_to(ROOT)}: {value.get('schema')!r}")
    return value


def nonempty_strings(value: object) -> bool:
    return isinstance(value, list) and bool(value) and all(isinstance(x, str) and x.strip() for x in value)


def registry_providers() -> set[str]:
    src = (ROOT / "provider_capabilities.go").read_text(encoding="utf-8")
    providers = set(re.findall(r'Provider:\s*"([^"]+)"', src))
    if "Provider: tradeInsightProviderName" in src:
        providers.add("TradeInsight")
    return providers


def routed_providers() -> set[str]:
    src = (ROOT / "provider_router.go").read_text(encoding="utf-8")
    start = src.find("func routeChains")
    if start < 0:
        return set()
    tail = src[start:]
    end = tail.find("\n}\n")
    block = tail if end < 0 else tail[: end + 3]
    providers: set[str] = set()
    for match in re.finditer(r':\s*\{([^}]*)\}', block):
        providers.update(re.findall(r'"([^"]+)"', match.group(1)))
        if "tradeInsightProviderName" in match.group(1):
            providers.add("TradeInsight")
    return {ROUTER_PROVIDER_ALIASES.get(provider, provider) for provider in providers}


def production_go_files() -> list[Path]:
    files: list[Path] = []
    for path in ROOT.rglob("*.go"):
        rel = path.relative_to(ROOT).as_posix()
        if path.name.endswith("_test.go") or rel.startswith(EXCLUDED_SOURCE_PREFIXES):
            continue
        files.append(path)
    return sorted(files)


def production_hosts() -> dict[str, set[str]]:
    hosts: dict[str, set[str]] = {}
    rx = re.compile(r'(?P<url>(?:https?|wss?)://[A-Za-z0-9._:-]+)')
    for path in production_go_files():
        rel = path.relative_to(ROOT).as_posix()
        if rel in NON_BACKEND_NETWORK_LITERAL_FILES:
            continue
        text = path.read_text(encoding="utf-8", errors="replace")
        for match in rx.finditer(text):
            raw = match.group("url")
            host = (urlsplit(raw).hostname or "").lower().strip(".")
            if not host or host in LOCAL_OR_EXAMPLE_HOSTS:
                continue
            hosts.setdefault(host, set()).add(rel)
    return hosts


def validate_matrix(matrix: dict, errors: list[str]) -> tuple[set[str], set[str]]:
    if matrix.get("scopeId") != "ADAPT-DATAHEALTH-BASELINE-001" or matrix.get("issue") != 80:
        errors.append("provider matrix scope/issue mismatch")
    if matrix.get("generalRoutingAuthority") != "Smart Provider Router v2":
        errors.append("provider matrix must retain Smart Provider Router v2 general authority")
    rows = matrix.get("providers")
    if not isinstance(rows, list) or not rows:
        errors.append("provider matrix providers must be a non-empty array")
        return set(), set()
    names: set[str] = set()
    capability_keys: set[str] = set()
    required_cap_fields = {
        "capability", "canonicalOwner", "consumers", "authorityClass",
        "freshnessPolicy", "cachePersistenceOwner", "fallbackCandidates",
        "materiality", "degradedImpact", "routerDataset", "lifecycle",
    }
    for i, row in enumerate(rows):
        if not isinstance(row, dict):
            errors.append(f"provider[{i}] must be an object")
            continue
        name = str(row.get("provider", "")).strip()
        pid = str(row.get("id", "")).strip()
        if not name or not pid or not str(row.get("class", "")).strip():
            errors.append(f"provider[{i}] missing id/provider/class")
            continue
        if name in names:
            errors.append(f"duplicate provider matrix provider: {name}")
        names.add(name)
        caps = row.get("capabilities")
        if not isinstance(caps, list) or not caps:
            errors.append(f"{name}: capabilities must be non-empty")
            continue
        for j, cap in enumerate(caps):
            if not isinstance(cap, dict):
                errors.append(f"{name} capability[{j}] must be object")
                continue
            missing = [field for field in required_cap_fields if field not in cap]
            if missing:
                errors.append(f"{name} capability[{j}] missing fields: {', '.join(sorted(missing))}")
                continue
            if not str(cap.get("capability", "")).strip() or not str(cap.get("canonicalOwner", "")).strip():
                errors.append(f"{name} capability[{j}] missing capability/canonical owner")
            if not isinstance(cap.get("consumers"), list) or not cap.get("consumers"):
                errors.append(f"{name} capability[{j}] consumers must be non-empty")
            if not isinstance(cap.get("fallbackCandidates"), list):
                errors.append(f"{name} capability[{j}] fallbackCandidates must be an array")
            capability_keys.add(name + "::" + str(cap.get("capability", "")).strip())
    return names, capability_keys


def validate_slo(slo: dict, errors: list[str]) -> None:
    if slo.get("scopeId") != "ADAPT-DATAHEALTH-BASELINE-001" or slo.get("issue") != 80:
        errors.append("Data Health SLO scope/issue mismatch")
    metrics = slo.get("metrics", {})
    if not isinstance(metrics, dict):
        errors.append("Data Health SLO metrics must be an object")
    else:
        missing = sorted(REQUIRED_METRICS - set(metrics))
        if missing:
            errors.append("Data Health SLO metrics missing: " + ", ".join(missing))
    time_semantics = slo.get("timeSemantics", {})
    if "provider observation" not in str(time_semantics.get("authoritativeEvidenceTime", "")).lower():
        errors.append("SLO must make provider observation/evidence time authoritative")
    if "bookkeeping" not in str(time_semantics.get("retrievalAndCacheTime", "")).lower():
        errors.append("SLO must treat retrieval/cache time as bookkeeping")
    if "unknown" not in str(time_semantics.get("unknownObservationTime", "")).lower():
        errors.append("SLO must preserve unknown observation time")
    shared = set(time_semantics.get("sharedConsumers", [])) if isinstance(time_semantics.get("sharedConsumers"), list) else set()
    if not {"Scanner", "Opportunity Radar"}.issubset(shared):
        errors.append("SLO must bind Scanner and Opportunity Radar to shared freshness semantics")
    degradation = slo.get("degradation", {})
    if not isinstance(degradation, dict) or not nonempty_strings(degradation.get("partialCoverageReasons")) or not nonempty_strings(degradation.get("dataDegradedReasons")):
        errors.append("SLO must define truthful PARTIAL COVERAGE and DATA DEGRADED reasons")
    recovery = slo.get("recovery", {})
    for key in ("automatic", "hysteresisRequired", "clearStaleBannerWhenRecovered", "preserveAuthorityRules"):
        if not isinstance(recovery, dict) or recovery.get(key) is not True:
            errors.append(f"SLO recovery contract requires {key}=true")


def validate_fetch_paths(fetch: dict, matrix_names: set[str], errors: list[str]) -> tuple[set[str], dict[str, set[str]]]:
    if fetch.get("scopeId") != "ADAPT-DATAHEALTH-BASELINE-001" or fetch.get("issue") != 80:
        errors.append("fetch-path contract scope/issue mismatch")
    entries = fetch.get("entries")
    if not isinstance(entries, list) or not entries:
        errors.append("fetch-path entries must be a non-empty array")
        return set(), {}
    registered_hosts: set[str] = set()
    host_owners: dict[str, set[str]] = {}
    seen_ids: set[str] = set()
    for i, row in enumerate(entries):
        if not isinstance(row, dict):
            errors.append(f"fetch entry[{i}] must be object")
            continue
        fid = str(row.get("id", "")).strip()
        provider = str(row.get("provider", "")).strip()
        host = str(row.get("host", "")).lower().strip().strip(".")
        disposition = str(row.get("disposition", "")).strip()
        owners = row.get("sourceOwners")
        if not fid or fid in seen_ids:
            errors.append(f"fetch entry[{i}] missing/duplicate id: {fid!r}")
        seen_ids.add(fid)
        if provider not in matrix_names:
            errors.append(f"{fid}: provider missing from matrix: {provider}")
        if not host:
            errors.append(f"{fid}: host missing")
        else:
            registered_hosts.add(host)
        if disposition not in ALLOWED_DISPOSITIONS:
            errors.append(f"{fid}: invalid disposition {disposition!r}")
        if not nonempty_strings(owners):
            errors.append(f"{fid}: sourceOwners must be non-empty")
        else:
            host_owners.setdefault(host, set()).update(str(x) for x in owners)
            missing_owners = [str(x) for x in owners if not (ROOT / str(x)).is_file()]
            if missing_owners:
                errors.append(f"{fid}: source owner paths do not exist: {', '.join(missing_owners)}")
        for field in ("capability", "reason", "justification"):
            if not str(row.get(field, "")).strip():
                errors.append(f"{fid}: {field} missing")
        if not isinstance(row.get("generalRouterApplicable"), bool) or not isinstance(row.get("bypass"), bool):
            errors.append(f"{fid}: generalRouterApplicable/bypass must be booleans")
        if disposition == "MIGRATE" and row.get("followupIssue") not in (81, 82):
            errors.append(f"{fid}: MIGRATE requires followupIssue #81 or #82")
        if disposition == "DIRECT_AUTHORITY" and row.get("bypass") is not True:
            errors.append(f"{fid}: DIRECT_AUTHORITY must explicitly acknowledge direct bypass=true")
        if disposition == "N/A" and not str(row.get("reason", "")).strip():
            errors.append(f"{fid}: N/A requires explicit reason")
    return registered_hosts, host_owners


def validate_existing_policy(errors: list[str]) -> None:
    p = json.loads(POLICY_PATH.read_text(encoding="utf-8"))
    policy = p.get("policy", {})
    for key in (
        "session_aware", "provider_vs_cache_timestamps", "stale_while_revalidate",
        "fallback_recovery", "material_change_priority", "selected_symbol_priority",
        "adaptive_policy_changes_require_shadow_validation",
    ):
        if policy.get(key) is not True:
            errors.append("existing Data Health policy missing " + key)
    if policy.get("market_critical_priority", [])[:2] != ["SPY", "QQQ"]:
        errors.append("existing Data Health policy lost SPY/QQQ critical priority")
    freshness = (ROOT / "data_freshness.go").read_text(encoding="utf-8")
    for token in ("ProviderTimestamp", "CacheAt", "CheckAgeMs", "DataAgeMs", "FreshLimitMs", "StaleLimitMs", "Fallback", "Reason"):
        if token not in freshness:
            errors.append("freshness diagnostic missing " + token)
    live = (ROOT / "live_subscription_manager.go").read_text(encoding="utf-8")
    if 'marketCriticalLiveSymbols = []string{"SPY", "QQQ"}' not in live:
        errors.append("live priority source lost SPY/QQQ")


def validate_adaptive_current_contracts(errors: list[str]) -> None:
    for path, tokens in ADAPTIVE_CURRENT_CONTRACTS.items():
        if not path.is_file():
            errors.append(f"missing CURRENT Adaptive Data Health contract: {path.relative_to(ROOT)}")
            continue
        text = path.read_text(encoding="utf-8")
        for token in tokens:
            if token not in text:
                errors.append(
                    f"CURRENT Adaptive Data Health contract drift in {path.relative_to(ROOT)}: missing {token!r}"
                )


def main() -> int:
    errors: list[str] = []
    validate_existing_policy(errors)
    validate_adaptive_current_contracts(errors)
    matrix = load_json(MATRIX_PATH, errors)
    slo = load_json(SLO_PATH, errors)
    fetch = load_json(FETCH_PATH, errors)
    matrix_names, capability_keys = validate_matrix(matrix, errors)
    validate_slo(slo, errors)
    registered_hosts, _ = validate_fetch_paths(fetch, matrix_names, errors)

    registry = registry_providers()
    missing_registry = sorted(registry - matrix_names)
    if missing_registry:
        errors.append("active provider registry rows missing from matrix: " + ", ".join(missing_registry))

    routed = routed_providers()
    missing_routed = sorted(routed - matrix_names)
    if missing_routed:
        errors.append("Smart Router v2 providers missing from matrix: " + ", ".join(missing_routed))

    discovered = production_hosts()
    missing_hosts = sorted(set(discovered) - registered_hosts)
    if missing_hosts:
        for host in missing_hosts:
            errors.append(f"unclassified production external host {host}: {', '.join(sorted(discovered[host]))}")

    print("DE.PULSE Adaptive Data Health provider contract")
    print(f"provider matrix rows: {len(matrix_names)}")
    print(f"capability rows: {len(capability_keys)}")
    print(f"provider registry covered: {len(registry)}/{len(registry)}" if not missing_registry else f"provider registry missing: {len(missing_registry)}")
    print(f"Router v2 provider members classified: {len(routed)}")
    print(f"production external hosts discovered/classified: {len(discovered)}/{len(registered_hosts)} registered")
    if errors:
        print("DE.PULSE Adaptive Data Health: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 2
    print("provider/capability classification: PASS")
    print("Data Health SLO/degradation/recovery contract: PASS")
    print("CURRENT Adaptive Roadmap/Build Plan/Build Process/Delivery Process: PASS")
    print("runtime external-host recurrence protection: PASS")
    print("DE.PULSE Adaptive Data Health: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
