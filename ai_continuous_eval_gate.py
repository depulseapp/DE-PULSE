#!/usr/bin/env python3
"""Offline v18.6 AI hardening/rights evaluation contract."""

from __future__ import annotations

import json
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).resolve().parent
POLICY = ROOT / "ai_eval_policy.json"
RIGHTS = ROOT / "provider_dataset_ai_rights_registry.json"
HARDENING = ROOT / "ai_hardening_v18_6.go"
CLIENTS = ROOT / "ai_clients.go"
TESTS = ROOT / "ai_hardening_test.go"
DOC_IMPACT = ROOT / "release" / "v18.6.0" / "DOCUMENTATION-IMPACT.md"

REQUIRED_LANES = {
    "golden",
    "citation",
    "contradiction",
    "missing-evidence",
    "injection-adversarial",
    "bounded-context",
    "cache-identity-ttl",
    "strict-schema-abstention",
    "rights-approved-denied",
    "structured-output-capability",
    "cost-latency-telemetry",
}

REQUIRED_TESTS = {
    "TestV186AIGoldenStructuredOutput",
    "TestV186AICitationStrictness",
    "TestV186AIContradictionMissingEvidenceContext",
    "TestV186AIInjectionAdversarialBoundary",
    "TestV186AIBoundedSemanticContext",
    "TestV186AICacheIdentityAndTTL",
    "TestV186AIStrictSchemaSafeAbstention",
    "TestV186AIRightsApprovedDeniedFixtures",
    "TestV186AIStructuredOutputProviderPayloads",
    "TestV186AICostLatencyTelemetry",
}


def fail(message: str) -> None:
    raise AssertionError(message)


def load(path: pathlib.Path) -> dict:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} invalid JSON: {exc}")


def source_int(source: str, name: str) -> int:
    match = re.search(rf"\b{name}\s*=\s*(\d+)", source)
    if not match:
        fail(f"missing integer constant {name}")
    return int(match.group(1))


def main() -> int:
    policy = load(POLICY)
    rights = load(RIGHTS)
    hardening = HARDENING.read_text(encoding="utf-8")
    clients = CLIENTS.read_text(encoding="utf-8")
    tests = TESTS.read_text(encoding="utf-8")
    docs = DOC_IMPACT.read_text(encoding="utf-8")

    if policy.get("schema") != "DE.PULSE-AI-CONTINUOUS-EVAL-1" or policy.get("policyVersion") != "v18.6-ai-eval-v1":
        fail("AI eval policy schema/version mismatch")
    lanes = set(policy.get("requiredLanes", []))
    missing_lanes = REQUIRED_LANES - lanes
    if missing_lanes:
        fail(f"AI eval policy missing lanes: {sorted(missing_lanes)}")
    budgets = policy.get("budgets", {})
    if budgets.get("liveProviderCallsInCI") != 0:
        fail("normal CI AI eval must remain offline with zero live provider calls")
    if budgets.get("maxPromptBytes") != source_int(hardening, "aiMaxPromptBytes"):
        fail("AI prompt byte budget drift")
    if budgets.get("maxPromptTokenUpperBound") != source_int(hardening, "aiMaxPromptTokenUpperBound"):
        fail("AI prompt token-upper-bound drift")
    if budgets.get("maxCacheEntries") != source_int(hardening, "aiCacheMaxEntries"):
        fail("AI cache bound policy drift")
    if budgets.get("cacheTTLSeconds") != 900 or "15 * time.Minute" not in hardening:
        fail("AI cache TTL policy drift")

    required_prod_markers = (
        "filterAIResearchPackageForEgress(pkg)",
        "buildBoundedAIUserPrompt(task, pkg)",
        "aiInferenceCacheKey(req, task, pkg, routing, settings)",
        "loadAICacheV186(cacheKey, time.Now())",
        "parseAIStructuredPayloadStrict(text, pkg)",
        "recordAIInferenceTelemetry",
    )
    for marker in required_prod_markers:
        if marker not in clients:
            fail(f"production AI path missing hardening marker: {marker}")
    generate_body = clients.split("func (a *Application) GenerateAIForUser", 1)[1].split("type OpenRouterConfig", 1)[0]
    if "parseAIStructuredPayload(text)" in generate_body:
        fail("production inference still accepts the historical lenient schema parser")
    for marker in ('"response_format"', '"require_parameters": true', '"responseMimeType": "application/json"', '"text": map[string]any{'):
        if marker not in clients:
            fail(f"structured-output capability marker missing: {marker}")

    if rights.get("schema") != "DE.PULSE-PROVIDER-DATASET-AI-RIGHTS-1":
        fail("provider/dataset AI rights schema mismatch")
    if rights.get("policyVersion") != "v18.6.0-ai-egress-rights-1":
        fail("provider/dataset AI rights policy version mismatch")
    if str(rights.get("defaultDecision", "")).upper() != "BLOCK":
        fail("provider/dataset AI rights default must fail closed")
    rows = rights.get("records")
    if not isinstance(rows, list) or not rows:
        fail("provider/dataset AI rights registry is empty")
    keys: set[tuple[str, str]] = set()
    for row in rows:
        provider = str(row.get("provider", "")).strip().lower()
        dataset = str(row.get("dataset", "")).strip().lower()
        if not provider or not dataset:
            fail("rights row has empty provider/dataset")
        key = (provider, dataset)
        if key in keys:
            fail(f"duplicate rights row {provider}/{dataset}")
        keys.add(key)
        decision = str(row.get("decision", "")).strip().upper()
        if decision not in {"ALLOW", "BLOCK", "DENY"}:
            fail(f"unsupported AI egress decision {decision} for {provider}/{dataset}")
        if decision == "ALLOW":
            if not row.get("evidenceBound") or any(str(row.get(field, "")).upper() != "APPROVED" for field in ("commercialUse", "redistribution", "aiUse")):
                fail(f"ALLOW without evidence-bound commercial/redistribution/AI approval for {provider}/{dataset}")
    if ("unknown", "unknown") not in keys:
        fail("rights registry lacks explicit unknown/unknown fail-closed record")

    for test_name in REQUIRED_TESTS:
        if f"func {test_name}(" not in tests:
            fail(f"AI eval lane missing test {test_name}")
    if "http.DefaultClient = &http.Client{Transport:" not in tests:
        fail("structured-output provider tests are not network-isolated")

    if "AUDIT-18-AI-001" not in docs or "AUDIT-18-AI-RIGHTS-001" not in docs:
        fail("Documentation Impact Manifest lacks AI hardening/rights dispositions")
    if "AI egress" not in docs or "continuous" not in docs.lower():
        fail("Documentation Impact Manifest lacks AI egress/eval proof requirements")

    for path in (POLICY, RIGHTS):
        text = path.read_text(encoding="utf-8").lower()
        for secret_marker in ("sk-", "api_key=", "api-key=", "bearer ey"):
            if secret_marker in text:
                fail(f"possible secret-like value in {path.name}")

    print(
        "PASS: v18.6 AI hardening offline eval contract; "
        f"{len(REQUIRED_LANES)} lanes, {len(rows)} provider/dataset rights rows, zero live-provider CI calls"
    )
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except AssertionError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
