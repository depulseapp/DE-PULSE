#!/usr/bin/env python3
"""Materialize the v18.5.1 G1 placement map and immutable freeze overlay.

The audit reconciliation ledger remains the immutable G0/G1 input snapshot. G1
adds a separate placement map so every tracked ID has an accountable owner,
build lane, next release/final disposition, regression identity and evidence
contract without rewriting historical audit evidence.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LEDGER = ROOT / "release" / "v18.5.1" / "V17-V18-IMPLEMENTATION-RECONCILIATION.json"
DRAFT = ROOT / "evidence" / "v18.5.1-G1-placement-draft.json"
MAP = ROOT / "release" / "v18.5.1" / "V18.5.1-G1-PLACEMENT-MAP.json"
FREEZE = ROOT / "release" / "v18.5.1" / "G1-SCOPE-FREEZE.json"
EXPECTED_ROWS = 295
EXPECTED_IDS = [
    "COPY-18.5.1-001", "SYMBOL-18.5.1-001", "SYMBOL-18.5.1-002", "NAV-18.5.1-001",
    "RESEARCH-v15.1.0-17-19-REOPENED", "VERSION-18.5.1-002", "HOVER-18.5.1-001",
    "AUDIT-18-UI-001", "AUDIT-18-CI-001", "AUDIT-18-QA-001",
]
CURRENT_IDS = set(EXPECTED_IDS)
FINAL = {"FRESH_PASS", "INTENTIONALLY_SUPERSEDED", "NOT_APPLICABLE"}
CURRENT_OWNERS = {
    "COPY-18.5.1-001": ("Renderer / UX", "RENDERER_UI"),
    "SYMBOL-18.5.1-001": ("Core / Backend", "CORE_BACKEND"),
    "SYMBOL-18.5.1-002": ("Core / Backend", "CORE_BACKEND"),
    "NAV-18.5.1-001": ("Renderer / UX", "RENDERER_UI"),
    "RESEARCH-v15.1.0-17-19-REOPENED": ("Renderer / UX", "RENDERER_UI"),
    "VERSION-18.5.1-002": ("Release Identity / UX", "RELEASE_IDENTITY_UI"),
    "HOVER-18.5.1-001": ("Renderer / UX", "RENDERER_UI"),
    "AUDIT-18-UI-001": ("Renderer / UX", "RENDERER_UI"),
    "AUDIT-18-CI-001": ("CI / Delivery", "CI_DELIVERY"),
    "AUDIT-18-QA-001": ("QA / Release Assurance", "QA_BEHAVIOR"),
}


def collect(node, path="$"):
    rows = []
    if isinstance(node, dict):
        if isinstance(node.get("id"), str) and isinstance(node.get("status"), str):
            rows.append((path, node))
        for key, value in node.items():
            rows.extend(collect(value, f"{path}.{key}"))
    elif isinstance(node, list):
        for idx, value in enumerate(node):
            rows.extend(collect(value, f"{path}[{idx}]"))
    return rows


def source_hints(row):
    hints = []
    historical = row.get("historicalEvidence")
    if isinstance(historical, list):
        for evidence in historical:
            if isinstance(evidence, dict) and isinstance(evidence.get("file"), str):
                hints.append(evidence["file"])
    for key in ("source", "file", "path"):
        if isinstance(row.get(key), str):
            hints.append(row[key])
    return list(dict.fromkeys(hints))[:2]


def owner_lane(row, hints):
    rid = row["id"]
    if rid in CURRENT_OWNERS:
        return CURRENT_OWNERS[rid]
    text = " ".join([rid, row.get("title", ""), row.get("name", ""), *hints]).lower()
    if ".github/" in text or "workflow" in text or " ci" in text or "delivery" in text:
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
        return str(row["candidateRelease"]), "EXPLICIT_PLACEMENT"
    if row.get("placement"):
        return str(row["placement"]), "EXPLICIT_PLACEMENT"
    if row.get("status") in FINAL:
        return "NONE", "FINAL"
    return "v18.8+", "V18_ZERO_GAP"


def build_map(data, raw):
    rows = collect(data)
    ids = [row["id"] for _, row in rows]
    if len(ids) != EXPECTED_ROWS or len(set(ids)) != EXPECTED_ROWS:
        raise SystemExit(f"ledger conservation failed: rows={len(ids)} unique={len(set(ids))} expected={EXPECTED_ROWS}")
    assigned = data.get("adaptiveReleaseTrain", {}).get("currentSlice", {}).get("assignedIds", [])
    if assigned != EXPECTED_IDS:
        raise SystemExit("audited ten-ID current slice drifted before G1 freeze")
    placements = []
    for path, row in rows:
        hints = source_hints(row)
        owner, lane = owner_lane(row, hints)
        release, reason = release_for(row)
        rid = row["id"]
        placements.append({
            "id": rid,
            "path": path,
            "status": row["status"],
            "disposition": "CURRENT_V18.5.1" if rid in CURRENT_IDS else ("FINAL" if row["status"] in FINAL else "PLACED_NEXT_OR_FUTURE"),
            "release": release,
            "owner": owner,
            "lane": lane,
            "sourceHints": hints,
            "deps": [x for x in row.get("links", []) if isinstance(x, str)] if isinstance(row.get("links"), list) else [],
            "regression": "REG-" + "".join(ch if ch.isalnum() else "-" for ch in rid).strip("-"),
            "impact": row.get("title") or row.get("name") or rid,
            "reason": reason,
        })
    return {
        "schema": "DE.PULSE-v18.5.1-G1-PLACEMENT-MAP-1",
        "release": "v18.5.1",
        "sourceLedger": str(LEDGER.relative_to(ROOT)),
        "sourceLedgerSha256": hashlib.sha256(raw).hexdigest(),
        "trackedRows": EXPECTED_ROWS,
        "ownerSemantics": "G1 accountable owner/lane; G2/G4 must resolve exact current source owner before implementation closure.",
        "evidenceProfile": {
            "required": [
                "current source/behavior owner",
                "focused executable regression",
                "behavior-first runtime/browser proof when user-facing",
                "exact macOS Apple Silicon and Windows x64 package proof when affected",
            ],
            "rule": "Placement never marks implementation, regression, behavior, package or closure evidence complete.",
        },
        "placements": sorted(placements, key=lambda item: item["id"]),
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--apply", action="store_true", help="write canonical release placement map and freeze manifest")
    args = parser.parse_args()
    raw = LEDGER.read_bytes()
    data = json.loads(raw)
    placement = build_map(data, raw)
    target = MAP if args.apply else DRAFT
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(json.dumps(placement, indent=2) + "\n")
    map_sha = hashlib.sha256(target.read_bytes()).hexdigest()
    if args.apply:
        manifest = {
            "schema": "DE.PULSE-v18.5.1-G1-SCOPE-FREEZE-1",
            "release": "v18.5.1",
            "gate": "G1",
            "freezeState": "FROZEN",
            "decision": "G1_SCOPE_FROZEN",
            "date": "2026-08-17",
            "incomingStable": {"tag": "v18.5.0-stable", "commit": "0d37ca35f5fc3ad89cebed506cc5a4c2d6a7a680"},
            "auditIntegrationCommit": "664418e143a969b63cd2169616278dd54e501d6b",
            "sourceLedger": {"path": str(LEDGER.relative_to(ROOT)), "sha256": placement["sourceLedgerSha256"], "trackedRows": EXPECTED_ROWS},
            "placementMap": {"path": str(MAP.relative_to(ROOT)), "sha256": map_sha, "rows": EXPECTED_ROWS},
            "currentSliceIds": EXPECTED_IDS,
            "laterCapacity": ["v18.6", "v18.7", "v18.8+", "FINAL_V18_X", "v19-v20 where already explicit"],
            "rules": [
                "The audit reconciliation ledger is preserved as the G0/G1 input snapshot; this manifest is the binding G1 freeze overlay.",
                "No tracked ID may disappear. Any addition, removal, current-slice reassignment or placement-map hash change reopens G1.",
                "Current v18.5.1 scope is exactly the ten IDs listed here; linked dependencies do not become independent scope IDs unless G1 is reopened.",
                "Placement is not implementation. No row closes without executable behavior and applicable exact-package evidence.",
                "Future placements may move only through a later release G1 replan with requirement conservation.",
            ],
        }
        FREEZE.write_text(json.dumps(manifest, indent=2) + "\n")
        print(f"G1 scope materialized: {EXPECTED_ROWS} rows · mapSha256={map_sha}")
    else:
        print(f"G1 placement draft: {EXPECTED_ROWS} rows -> {target.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
