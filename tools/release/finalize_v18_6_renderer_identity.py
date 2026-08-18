#!/usr/bin/env python3
from __future__ import annotations

import json
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
IDENTITY = ROOT / "release_identity.json"
MANIFEST = ROOT / "renderer" / "qa" / "manifest.json"


def main() -> None:
    identity = json.loads(IDENTITY.read_text(encoding="utf-8"))
    version = identity["version"]
    build_id = identity["build_id"]
    if version != "18.6.0" or build_id != "v18.6.0-stable-20260818":
        raise SystemExit(f"unexpected release identity: {version} / {build_id}")

    # Re-run the canonical sync after the release-coupled asset list was hardened.
    subprocess.run(["python3", "release_identity.py", "--sync"], cwd=ROOT, check=True)

    manifest = json.loads(MANIFEST.read_text(encoding="utf-8"))
    releases = manifest.get("releases")
    if not isinstance(releases, list):
        raise SystemExit("renderer QA manifest releases must be a list")
    releases = [entry for entry in releases if str(entry.get("version")) != version]
    releases.insert(0, {
        "version": version,
        "date": "2026-08-18",
        "status": "STABLE",
        "summary": "STABLE adaptive utility and intelligence hardening: shared Scanner/Radar snapshot acquisition, serialized Session Intelligence Coordinator, supporting-surface and legacy-route consolidation, role-aware documentation, explicit dependency/provider readiness, bounded AI context/cache/schema/evaluation, and rights-aware fail-closed AI egress; deterministic Day/Swing/Long formulas, Smart Provider Router ownership, U.S. equities scope and No Execution remain unchanged.",
        "file": "v18.6.0.txt",
        "buildId": build_id,
        "checkpoint": "release/v18.6.0/G1-IMMUTABLE-SCOPE.md",
    })
    manifest["releases"] = releases
    MANIFEST.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")

    qa_text = ROOT / "renderer" / "qa" / "v18.6.0.txt"
    qa_text.write_text(
        "DE.PULSE v18.6.0 — Adaptive Utility & Intelligence Hardening\n"
        f"Build: {build_id}\n"
        "Release target: STABLE; promotion remains gated by G11–G16 certification.\n\n"
        "Implemented and pre-freeze qualified scope:\n"
        "- Shared Scanner / Opportunity Radar broad-snapshot acquisition with bounded freshness-aware reuse, partial-miss fetching and in-flight coalescing.\n"
        "- One serialized Session Intelligence Coordinator for Pre-Market Prep and Market Open Prep.\n"
        "- Market Activity demoted to supporting drill-down and legacy evidence routes consolidated.\n"
        "- Role-aware documentation with server-authoritative privileged access.\n"
        "- Explicit external dependency/provider readiness and durable User Action Required governance.\n"
        "- Bounded AI semantic context, complete cache identity + TTL, strict structured output/citation validation, safe abstention and continuous-evaluation telemetry.\n"
        "- Provider×dataset rights-aware AI egress; unknown or denied content is withheld before external model calls.\n\n"
        "Protected invariants:\n"
        "- Deterministic Day/Swing/Long formulas unchanged.\n"
        "- Smart Provider Router remains sole provider-routing owner; provider count cannot change Market Mode.\n"
        "- U.S. Equities Processing Boundary preserved.\n"
        "- Permanent No Execution Boundary preserved.\n"
        "- Desktop SQLite / hosted PostgreSQL architecture preserved.\n\n"
        "Certification state at QA-history creation: implementation complete; canonical G10 exact-identity qualification in progress; G11–G16 not yet claimed.\n",
        encoding="utf-8",
    )

    subprocess.run(["python3", "release_identity.py", "--verify"], cwd=ROOT, check=True)
    subprocess.run(["python3", "version_consistency_test.py"], cwd=ROOT, check=True)
    print("PASS: v18.6 renderer/cache/QA release identity finalized")


if __name__ == "__main__":
    main()
