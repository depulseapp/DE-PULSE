#!/usr/bin/env python3
"""Materialize a deterministic G1 placement/ownership draft for all v18 ledger rows.

The output is a review artifact until it is committed as the G1 placement map.
It conserves every tracked ID and never marks implementation or evidence complete.
"""
from __future__ import annotations

import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LEDGER = ROOT / "release" / "v18.5.1" / "V17-V18-IMPLEMENTATION-RECONCILIATION.json"
OUT = ROOT / "evidence" / "v18.5.1-G1-placement-draft.json"
CURRENT_IDS = {
    "COPY-18.5.1-001", "SYMBOL-18.5.1-001", "SYMBOL-18.5.1-002", "NAV-18.5.1-001",
    "RESEARCH-v15.1.0-17-19-REOPENED", "VERSION-18.5.1-002", "HOVER-18.5.1-001",
    "AUDIT-18-UI-001", "AUDIT-18-CI-001", "AUDIT-18-QA-001",
}
FINAL = {"FRESH_PASS", "INTENTIONALLY_SUPERSEDED", "NOT_APPLICABLE"}


def collect(node, path="$"):
    rows = []
    if isinstance(node, dict):
        if isinstance(node.get("id"), str) and isinstance(node.get("status"), str):
            rows.append((path, node))
        for k, v in node.items():
            rows.extend(collect(v, f"{path}.{k}"))
    elif isinstance(node, list):
        for i, v in enumerate(node):
            rows.extend(collect(v, f"{path}[{i}]"))
    return rows


def source_hints(row):
    hints = []
    for ev in row.get("historicalEvidence", []) if isinstance(row.get("historicalEvidence"), list) else []:
        if isinstance(ev, dict) and ev.get("file"):
            hints.append(ev["file"])
    for key in ("source", "file", "path"):
        if isinstance(row.get(key), str):
            hints.append(row[key])
    return list(dict.fromkeys(hints))[:8]


def owner_lane(row, hints):
    text = " ".join([row.get("id", ""), row.get("title", ""), row.get("name", ""), *hints]).lower()
    if ".github/" in text or " ci" in text or "workflow" in text or "delivery" in text:
        return "CI / Delivery", "CI_DELIVERY"
    if "renderer" in text or "css" in text or "responsive" in text or "header" in text or "hover" in text or "research" in text:
        return "Renderer / UX", "RENDERER_UI"
    if "security" in text or "credential" in text or "auth" in text or "secret" in text:
        return "Security / Platform", "SECURITY_PLATFORM"
    if "provider" in text or "finnhub" in text or "alpaca" in text or "tradeinsight" in text or "data" in text:
        return "Data / Provider Intelligence", "DATA_PROVIDER"
    if any(h.endswith(".go") for h in hints) or " api" in text or "persistence" in text or "state" in text:
        return "Core / Backend", "CORE_BACKEND"
    return "DE.PULSE Core", "CROSS_CUTTING_RECONCILIATION"


def release_for(row):
    rid = row["id"]
    if rid in CURRENT_IDS:
        return "v18.5.1", "CURRENT_RECOVERY"
    if row.get("candidateRelease"):
        return row["candidateRelease"], row.get("completionLane") or "EXPLICIT_LEDGER_PLACEMENT"
    if row.get("placement"):
        return row["placement"], "ROADMAP_FUTURE"
    if row.get("status") in FINAL:
        return "NONE", "FINAL_DISPOSITION"
    return "v18.8+", "V18_ZERO_GAP_RECONCILIATION"


def evidence_contract(row):
    explicit = row.get("requiredCurrentEvidence") or row.get("requiredClosure") or []
    if not isinstance(explicit, list):
        explicit = [str(explicit)]
    base = [
        "current source/behavior owner resolved before implementation closure",
        "focused executable regression for affected behavior",
        "behavior-first runtime/browser proof when user-facing",
        "exact macOS Apple Silicon and Windows x64 package proof when affected",
    ]
    return list(dict.fromkeys([str(x) for x in explicit] + base))


def main():
    raw = LEDGER.read_bytes()
    data = json.loads(raw)
    rows = collect(data)
    seen = set()
    placements = []
    for path, row in rows:
        rid = row["id"]
        if rid in seen:
            raise SystemExit(f"duplicate id: {rid}")
        seen.add(rid)
        hints = source_hints(row)
        owner, lane = owner_lane(row, hints)
        release, release_lane = release_for(row)
        placements.append({
            "id": rid,
            "title": row.get("title") or row.get("name") or rid,
            "ledgerPath": path,
            "ledgerStatus": row["status"],
            "scopeDisposition": "CURRENT_V18.5.1" if rid in CURRENT_IDS else ("FINAL" if row["status"] in FINAL else "PLACED_NEXT_OR_FUTURE"),
            "candidateRelease": release,
            "releaseLane": release_lane,
            "accountableOwner": owner,
            "buildLane": lane,
            "sourceHints": hints,
            "dependencyIds": [x for x in row.get("links", []) if isinstance(x, str)] if isinstance(row.get("links"), list) else [],
            "regressionId": "REG-" + "".join(ch if ch.isalnum() else "-" for ch in rid).strip("-"),
            "userImpact": row.get("title") or row.get("name") or rid,
            "placementReason": (
                "Audited urgent-recovery scope frozen for v18.5.1" if rid in CURRENT_IDS
                else "Preserve explicit ledger/roadmap release placement" if row.get("candidateRelease") or row.get("placement")
                else "Conserved into later v18 zero-gap reconciliation; may move earlier only through a future G1 replan"
            ),
            "evidenceContract": evidence_contract(row),
            "closureState": "OPEN_UNTIL_CURRENT_EVIDENCE",
        })
    declared = data.get("baseline", {}).get("totalTrackedRows")
    if len(placements) != declared:
        raise SystemExit(f"placement count {len(placements)} != declared {declared}")
    out = {
        "schema": "DE.PULSE-v18.5.1-G1-PLACEMENT-MAP-DRAFT-1",
        "release": "v18.5.1",
        "sourceLedger": str(LEDGER.relative_to(ROOT)),
        "sourceLedgerSha256": hashlib.sha256(raw).hexdigest(),
        "declaredTrackedRows": declared,
        "materializedRows": len(placements),
        "currentSliceIds": sorted(CURRENT_IDS),
        "policy": "Placement/ownership metadata only. This artifact cannot mark implementation, regression, behavior, package or closure evidence complete.",
        "placements": sorted(placements, key=lambda x: x["id"]),
    }
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(out, indent=2) + "\n")
    print(f"G1 placement draft: {len(placements)} rows -> {OUT.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
