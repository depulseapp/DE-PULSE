#!/usr/bin/env python3
"""Registration-aware extension of the canonical G2/Data Health recurrence contract.

#95 moved provider route/configuration/diagnostic metadata into provider_registration.go.
The original #80 parser intentionally remains stable for its historical baseline,
while this extension makes the new canonical registration fail closed against the
same provider-capability matrix. It is invoked only by the existing G2 source-health
entrypoint; it does not create a new workflow or architecture gate family.
"""
from __future__ import annotations

import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
REGISTRATION_PATH = ROOT / "provider_registration.go"
MATRIX_PATH = ROOT / "governance" / "data-health" / "provider-capability-matrix.json"


def fail(message: str) -> None:
    print("DE.PULSE provider-registration Data Health: FAIL", file=sys.stderr)
    print(" - " + message, file=sys.stderr)
    raise SystemExit(2)


def matrix_providers() -> set[str]:
    try:
        matrix = json.loads(MATRIX_PATH.read_text(encoding="utf-8"))
    except Exception as exc:
        fail("provider-capability matrix unreadable: " + str(exc))
    rows = matrix.get("providers")
    if not isinstance(rows, list) or not rows:
        fail("provider-capability matrix providers must be a non-empty array")
    names = {
        str(row.get("provider", "")).strip()
        for row in rows
        if isinstance(row, dict) and str(row.get("provider", "")).strip()
    }
    if not names:
        fail("provider-capability matrix resolved zero provider identities")
    return names


def registered_providers(source: str) -> set[str]:
    names = set(re.findall(r'\bName:\s*"([^"]+)"', source))
    if re.search(r'\bName:\s*tradeInsightProviderName\b', source):
        names.add("TradeInsight")
    return names


def routed_registration_providers(source: str) -> set[str]:
    names = set(re.findall(r'inheritedProductionRoute\("([^"]+)"', source))
    if (
        re.search(r'\bName:\s*tradeInsightProviderName\b', source)
        and "tradeInsightProductionHistoryRoute(" in source
    ):
        names.add("TradeInsight")
    return names


def main() -> int:
    if not REGISTRATION_PATH.is_file():
        fail("canonical provider_registration.go is missing")
    source = REGISTRATION_PATH.read_text(encoding="utf-8")
    matrix = matrix_providers()
    registered = registered_providers(source)
    routed = routed_registration_providers(source)

    if not registered:
        fail("canonical provider registration resolved zero providers")
    if not routed:
        fail("canonical provider registration resolved zero routed providers")

    missing_registered = sorted(registered - matrix)
    if missing_registered:
        fail("registered providers missing from provider-capability matrix: " + ", ".join(missing_registered))
    missing_routed = sorted(routed - matrix)
    if missing_routed:
        fail("Router v2 registration providers missing from provider-capability matrix: " + ", ".join(missing_routed))

    # Every provider referenced by a production route must also have a canonical
    # registration identity; otherwise a future helper could bypass onboarding.
    unregistered_routed = sorted(routed - registered)
    if unregistered_routed:
        fail("routed providers missing canonical registration identity: " + ", ".join(unregistered_routed))

    print("DE.PULSE provider-registration Data Health: PASS")
    print(f"canonical provider registrations classified: {len(registered)}/{len(registered)}")
    print(f"registration-driven Router v2 provider members classified: {len(routed)}")
    print("provider registration -> #80 matrix recurrence: FAIL_CLOSED")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
