#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path

from release_identity import asset_cache_token

ROOT = Path(__file__).resolve().parents[2]
I = json.loads((ROOT / "release_identity.json").read_text())
VERSION = str(I["version"])
BUILD = str(I["build_id"])
PREV = str(I["previous_stable"])
STABLE = str(I["stable_baseline"])
CHANNEL = str(I["channel"])
DISPLAY = str(I["display_version"])
errs: list[str] = []
LEGACY_REGISTRY_NAMES = {
    "certification": "certification_" + "plan.json",
    "ci_pipeline": "ci_pipeline_" + "plan.json",
}


def need(ok: bool, msg: str) -> None:
    if not ok:
        errs.append(msg)


def legacy_registry_path(kind: str, version: str) -> Path:
    name = LEGACY_REGISTRY_NAMES.get(kind)
    if not name:
        raise SystemExit(f"Version consistency: unknown legacy registry kind: {kind}")
    version = str(version or "").strip()
    if not version:
        raise SystemExit(f"Version consistency: legacy registry version is required: {kind}")
    archived = ROOT / "release" / "history" / f"v{version}" / "legacy-ci" / name
    if archived.is_file():
        return archived
    raise SystemExit(f"Version consistency: legacy registry missing: kind={kind} version={version}")


boot = (ROOT / "app_bootstrap.go").read_text()
renderer = (ROOT / "renderer" / "renderer.js").read_text()
version = (ROOT / "VERSION.txt").read_text()
index = (ROOT / "renderer" / "index.html").read_text()

patch_contract_path = ROOT / "release" / f"v{VERSION}" / "patch_contract.json"
patch_contract = json.loads(patch_contract_path.read_text()) if patch_contract_path.exists() else None
release_contract_path = ROOT / "release" / f"v{VERSION}" / "release_contract.json"
release_contract = json.loads(release_contract_path.read_text()) if release_contract_path.exists() else None

if patch_contract:
    cert_version = str(patch_contract.get("inherited_certification_plan", "")).strip()
    ci_version = str(patch_contract.get("inherited_ci_plan", "")).strip()
elif release_contract:
    cert_version = str(release_contract.get("legacy_certification_plan_version", "")).strip() or VERSION
    ci_version = str(release_contract.get("legacy_ci_plan_version", "")).strip() or VERSION
else:
    cert_version = VERSION
    ci_version = VERSION

ci = json.loads(legacy_registry_path("ci_pipeline", ci_version).read_text())
cert = json.loads(legacy_registry_path("certification", cert_version).read_text())

patch_path = ROOT / "renderer" / f"watchlist-desk-contract-v{VERSION}.js"
patch = patch_path.read_text() if patch_path.exists() else ""

overlay_name = str((release_contract or {}).get("identity_asset", "")).strip()
overlay_path = ROOT / "renderer" / overlay_name if overlay_name else None
overlay = overlay_path.read_text() if overlay_path and overlay_path.is_file() else ""

need(f'const appVersion = "{VERSION}"' in boot, "appVersion mismatch")
need(f'const buildID = "{BUILD}"' in boot, "buildID mismatch")
need(
    DISPLAY in version
    and BUILD in version
    and f"Stable baseline: {STABLE}" in version
    and f"Previous Stable: {PREV}" in version,
    "VERSION.txt mismatch",
)

renderer_identity = (
    f"const EXPECTED_RELEASE_VERSION='{VERSION}';" in renderer
    and f"const EXPECTED_BUILD_ID='{BUILD}';" in renderer
)
patch_identity = (
    f"DEPULSE_PATCH_VERSION = '{VERSION}'" in patch
    and f"DEPULSE_PATCH_BUILD_ID = '{BUILD}'" in patch
    and f"watchlist-desk-contract-v{VERSION}.js?v={VERSION}" in index
)
overlay_identity = bool(overlay_name) and (
    f"DEPULSE_RELEASE_VERSION = '{VERSION}'" in overlay
    and f"DEPULSE_RELEASE_BUILD_ID = '{BUILD}'" in overlay
    and "DEPULSE_RELEASE_QA_ENTRY" in overlay
    and f"{overlay_name}?v={asset_cache_token(overlay_name, canonical_git=True)}" in index
)
need(renderer_identity or patch_identity or overlay_identity, "renderer/patch/release-overlay version mismatch")

for doc in ("user.md", "developer.md", "limitations.md"):
    text = (ROOT / "renderer" / "docs" / doc).read_text()
    expected = {
        "user.md": "# DE.PULSE — User documentation",
        "developer.md": "# DE.PULSE — Developer documentation",
        "limitations.md": "# DE.PULSE — Capabilities & Limitations",
    }[doc]
    need(text.splitlines()[0] == expected, f"{doc} canonical first heading mismatch")

if patch_contract:
    need(str(patch_contract.get("release")) == VERSION, "patch contract release mismatch")
    need(
        str(cert.get("version")) == str(patch_contract.get("inherited_certification_plan")),
        "inherited certification plan mismatch",
    )
    need(
        str(ci.get("version")) == str(patch_contract.get("inherited_ci_plan")),
        "inherited CI plan mismatch",
    )
    need("no gate" in str(patch_contract.get("quality_rule", "")), "patch quality inheritance rule missing")
elif release_contract and overlay_name:
    need(str(release_contract.get("release")) == VERSION, "release contract version mismatch")
    need(bool(overlay_path and overlay_path.is_file()), "release identity overlay missing")
    legacy_cert = str(release_contract.get("legacy_certification_plan_version", "")).strip()
    legacy_ci = str(release_contract.get("legacy_ci_plan_version", "")).strip()
    need(str(cert.get("version")) == (legacy_cert or VERSION), "release contract certification plan inheritance mismatch")
    need(str(ci.get("version")) == (legacy_ci or VERSION), "release contract CI plan inheritance mismatch")
    legacy_plan_rule = (
        "Historical "
        + LEGACY_REGISTRY_NAMES["certification"]
        + " and "
        + LEGACY_REGISTRY_NAMES["ci_pipeline"]
        + " remain conserved legacy registries"
    )
    need(
        legacy_plan_rule in str(release_contract.get("legacyPlanRule", "")),
        "release contract legacy-plan rule missing",
    )
    need(
        "release identity overlay" in str(release_contract.get("documentationIdentityRule", "")).lower(),
        "release contract documentation identity rule missing",
    )
    readme = (ROOT / "README.md").read_text()
    for marker in (
        "## Current project truth",
        "AGENTS.md",
        "CLAUDE.md",
        "governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md",
        "handoff/CURRENT.md",
        "GitHub is the durable source of truth",
    ):
        need(marker in readme, f"README portability/resume marker missing: {marker}")
else:
    readme = (ROOT / "README.md").read_text()
    need(readme.startswith(f"# DE.PULSE v{VERSION}"), "README title mismatch")
    need(BUILD in readme and f"Current Stable baseline:** {PREV}" in readme, "README immediate Stable predecessor mismatch")
    need(f"**Channel:** {CHANNEL}" in "\n".join(readme.splitlines()[:8]), "README current channel mismatch")
    for doc in ("user.md", "developer.md", "limitations.md"):
        text = (ROOT / "renderer" / "docs" / doc).read_text()
        need(f"v{VERSION} {CHANNEL}" in "\n".join(text.splitlines()[:30]), f"{doc} current build section missing")
    need(str(ci.get("version")) == VERSION, "CI plan current-release identity mismatch")
    need(ci.get("policy", {}).get("baseline") == PREV + " Stable", "CI plan predecessor mismatch")
    need(str(cert.get("version")) == VERSION, "certification plan current-release identity mismatch")

if errs:
    print("Version consistency: FAIL")
    for error in errs:
        print(" -", error)
    raise SystemExit(2)

mode = "patch" if patch_contract else ("release-overlay" if overlay_identity else "monolith")
print(
    f"Version consistency: PASS · v{VERSION} canonical release identity aligned · "
    f"mode={mode} · historical long-form docs/plans conserved only when explicitly declared by the release contract."
)
