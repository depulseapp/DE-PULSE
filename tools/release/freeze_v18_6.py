#!/usr/bin/env python3
from __future__ import annotations

import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
IDENTITY = ROOT / "release_identity.json"

TARGET = {
    "schema": "DE.PULSE-RELEASE-IDENTITY-1",
    "version": "18.6.0",
    "display_version": "DE.PULSE v18.6.0",
    "build_id": "v18.6.0-stable-20260818",
    "channel": "STABLE",
    "stable_baseline": "v17.5.1",
    "previous_stable": "v18.5.2",
    "scope": "Adaptive utility and intelligence hardening: shared Scanner/Radar snapshot acquisition; serialized Session Intelligence Coordinator; supporting-surface and legacy-route consolidation; role-aware documentation; external dependency/provider readiness; bounded AI context, cache identity, strict schema and continuous evaluation; rights-aware fail-closed AI egress; preserve deterministic Day/Swing/Long formulas, Smart Provider Router ownership, U.S. Equities Processing Boundary and permanent No Execution Boundary",
    "bundle_version": "18600",
    "runtime_config": "PersonalMarketTerminal",
    "patch_predecessor": "v18.5.2",
    "application_bundle": "De-Pulse.app",
}

DOC_SECTION = {
    "user.md": """## v18.6.0 STABLE — Adaptive utility and intelligence hardening\n\nDE.PULSE now shares broad Scanner/Opportunity Radar snapshot acquisition, coordinates Pre-Market Prep and Market Open Prep through one serialized session-intelligence scheduler, and keeps low-value Market Activity / legacy evidence routes in supporting drill-down paths. Documentation is role-aware and privileged developer content is enforced server-side. Provider/dependency readiness is explicit and fail-closed. External AI research uses bounded material context, strict structured outputs, expiring cache identity, continuous-evaluation telemetry, and provider×dataset rights filtering before egress. Deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary, Smart Provider Router ownership, and the permanent No Execution Boundary remain unchanged.\n\n""",
    "developer.md": """## v18.6.0 STABLE — Adaptive utility and intelligence hardening\n\nThe v18.6 release consolidates broad snapshot acquisition behind the shared broker, replaces duplicate prep scheduling with the Session Intelligence Coordinator, demotes supporting Market Activity and retires legacy evidence routes, makes documentation access server-authoritative by role, and adds machine-enforced dependency/provider readiness. AI inference now uses a hard semantic context envelope, full provider/model/prompt/schema/rights-aware cache identity with TTL, strict schema/citation validation and safe abstention, continuous evaluation telemetry, and fail-closed provider×dataset egress rights. No parallel provider router, authentication model, persistence owner, deterministic scoring formula, or execution path is introduced.\n\n""",
    "limitations.md": """## v18.6.0 STABLE — Capabilities & limitations update\n\nAI and adaptive intelligence remain decision-support only. External AI receives only evidence that passes the explicit provider×dataset rights registry; unknown or denied rights are withheld. Model responses must satisfy the strict response schema and evidence-ID contract or DE.PULSE abstains/falls back. Context is materially ranked under a hard prompt envelope and cached inference expires. Provider count never changes Market Mode, deterministic Day/Swing/Long formulas remain canonical, and DE.PULSE still does not place, route, simulate, or execute trades.\n\n""",
}


def replace_line(text: str, prefix: str, replacement: str) -> str:
    pattern = rf"^{re.escape(prefix)}.*$"
    if not re.search(pattern, text, flags=re.MULTILINE):
        raise SystemExit(f"missing expected README line prefix: {prefix}")
    return re.sub(pattern, replacement, text, count=1, flags=re.MULTILINE)


def update_readme() -> None:
    path = ROOT / "README.md"
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()
    if not lines:
        raise SystemExit("README is empty")
    lines[0] = "# DE.PULSE v18.6.0 STABLE — Adaptive Utility & Intelligence Hardening"
    text = "\n".join(lines) + ("\n" if text.endswith("\n") else "")
    text = replace_line(text, "**Build:**", "**Build:** `v18.6.0-stable-20260818`")
    text = replace_line(text, "**Channel:**", "**Channel:** STABLE")
    text = replace_line(text, "**Current Stable baseline:**", "**Current Stable baseline:** v18.5.2")
    text = replace_line(text, "**Major v18 provenance anchor:**", "**Major v18 provenance anchor:** v17.5.1")
    marker = "v18.5 is the mandatory v18 Major Closure before v19."
    summary = (
        "v18.6.0 is the adaptive utility and intelligence hardening release over v18.5.2 STABLE. "
        "It shares broad Scanner/Radar acquisition, serializes session-preparation scheduling, consolidates supporting evidence surfaces, "
        "enforces role-aware documentation and provider/dependency readiness, and hardens AI context, cache, schema, evaluation and rights-aware egress. "
        "Deterministic Day/Swing/Long formulas, Smart Provider Router ownership, the U.S. Equities Processing Boundary and permanent No Execution Boundary are preserved."
    )
    idx = text.find(marker)
    if idx >= 0:
        end = text.find("\n\n", idx)
        if end < 0:
            end = len(text)
        text = text[:idx] + summary + text[end:]
    path.write_text(text, encoding="utf-8")


def update_docs() -> None:
    for name, section in DOC_SECTION.items():
        path = ROOT / "renderer" / "docs" / name
        text = path.read_text(encoding="utf-8")
        if "## v18.6.0 STABLE" in text[:5000]:
            continue
        first_newline = text.find("\n")
        if first_newline < 0:
            raise SystemExit(f"missing heading newline in {path}")
        text = text[: first_newline + 1] + "\n" + section + text[first_newline + 1 :]
        path.write_text(text, encoding="utf-8")


def main() -> None:
    current = json.loads(IDENTITY.read_text(encoding="utf-8"))
    current.update(TARGET)
    IDENTITY.write_text(json.dumps(current, indent=2) + "\n", encoding="utf-8")
    subprocess.run(["python3", "release_identity.py", "--sync"], cwd=ROOT, check=True)
    update_readme()
    update_docs()
    subprocess.run(["python3", "release_identity.py", "--verify"], cwd=ROOT, check=True)
    subprocess.run(["python3", "version_consistency_test.py"], cwd=ROOT, check=True)
    print("PASS: v18.6.0 canonical release identity synchronized")


if __name__ == "__main__":
    main()
