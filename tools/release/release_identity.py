#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
IDENTITY = ROOT / "release_identity.json"
RELEASE_COUPLED_ASSETS = (
    "renderer.js",
    "documentation-ui.js",
    "live-dom-reconcile.js",
    "watchlist-ui.js",
    "watchlist-desk.css",
    "market-header-ui.js",
    "ui-layout-contracts.css",
    "surface-consolidation.js",
    "surface-consolidation.css",
    "documentation-access.js",
)
LEGACY_REGISTRY_NAMES = {
    "certification": "certification_" + "plan.json",
    "ci_pipeline": "ci_pipeline_" + "plan.json",
}


def load():
    x = json.loads(IDENTITY.read_text())
    for key in (
        "version", "display_version", "build_id", "channel", "stable_baseline",
        "previous_stable", "scope", "bundle_version", "runtime_config", "application_bundle",
    ):
        if not str(x.get(key, "")).strip():
            raise SystemExit(f"release identity missing {key}")
    asset = str(x.get("renderer_identity_asset", "")).strip()
    if asset and ("/" in asset or "\\" in asset or asset.startswith(".")):
        raise SystemExit("release identity renderer asset must be a renderer-local filename")
    return x


def patch_contract(x):
    path = ROOT / "release" / f"v{x['version']}" / "patch_contract.json"
    if not path.exists():
        return None
    data = json.loads(path.read_text())
    if str(data.get("release")) != str(x["version"]):
        raise SystemExit("patch contract release mismatch")
    return data


def release_contract(x):
    path = ROOT / "release" / f"v{x['version']}" / "release_contract.json"
    if not path.exists():
        return None
    data = json.loads(path.read_text())
    if str(data.get("release")) != str(x["version"]):
        raise SystemExit("release contract version mismatch")
    return data


def renderer_identity_asset_name(x, contract):
    asset = str(x.get("renderer_identity_asset", "")).strip()
    if not asset and contract:
        asset = str(contract.get("identity_asset", "")).strip()
    if not asset:
        return ""
    if "/" in asset or "\\" in asset or asset.startswith("."):
        raise SystemExit("release identity asset must be a renderer-local filename")
    return asset


def overlay_asset_path(x, contract):
    asset = renderer_identity_asset_name(x, contract)
    return ROOT / "renderer" / asset if asset else None


def clean_head_blob_bytes(relative_path: str) -> bytes | None:
    """Return canonical HEAD bytes only when the tracked worktree path is unchanged.

    Windows checkouts may materialize CRLF while Git stores LF. Cache identity is
    a source identity, so verification must use repository blob bytes when the
    logical tracked file is unchanged. If the file is intentionally staged or
    modified (for example during --sync --verify), callers fall back to the
    working-tree bytes that are being prepared for the next commit.
    """
    for diff_args in (
        ("git", "diff", "--quiet", "--", relative_path),
        ("git", "diff", "--cached", "--quiet", "--", relative_path),
    ):
        result = subprocess.run(
            diff_args,
            cwd=ROOT,
            check=False,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        if result.returncode != 0:
            return None
    result = subprocess.run(
        ("git", "cat-file", "blob", f"HEAD:{relative_path}"),
        cwd=ROOT,
        check=False,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
    )
    return result.stdout if result.returncode == 0 else None


def asset_cache_token(name: str, *, canonical_git: bool = False) -> str:
    path = ROOT / "renderer" / name
    if not path.is_file():
        raise SystemExit(f"renderer cache identity asset missing: {path}")
    relative_path = f"renderer/{name}"
    data = clean_head_blob_bytes(relative_path) if canonical_git else None
    if data is None:
        data = path.read_bytes()
    # Match Git's blob object identity so the token is deterministic from source
    # bytes, independent of product version/build/work-slice and OS materialization.
    header = f"blob {len(data)}\0".encode("utf-8")
    return hashlib.sha1(header + data).hexdigest()[:16]


def sync_asset_cache_identity(index: str, name: str) -> str:
    token = asset_cache_token(name)
    return re.sub(rf"{re.escape(name)}\?v=[A-Za-z0-9._-]+", f"{name}?v={token}", index)


def legacy_registry_path(kind: str, version: str) -> Path:
    name = LEGACY_REGISTRY_NAMES.get(kind)
    if not name:
        raise SystemExit(f"unknown legacy registry kind: {kind}")
    version = str(version or "").strip()
    if not version:
        raise SystemExit(f"legacy registry version is required: {kind}")
    archived = ROOT / "release" / "history" / f"v{version}" / "legacy-ci" / name
    if archived.is_file():
        return archived
    raise SystemExit(f"legacy registry missing: kind={kind} version={version}")


def legacy_registry_versions(x, patch, contract) -> tuple[str, str]:
    if patch:
        return (
            str(patch.get("inherited_certification_plan", "")).strip(),
            str(patch.get("inherited_ci_plan", "")).strip(),
        )
    if contract:
        return (
            str(contract.get("legacy_certification_plan_version", "")).strip(),
            str(contract.get("legacy_ci_plan_version", "")).strip(),
        )
    return str(x["version"]), str(x["version"])


def sync(x):
    (ROOT / "VERSION.txt").write_text(
        f"{x['display_version']}\n"
        f"Build: {x['build_id']}\n"
        f"Channel: {x['channel']}\n"
        f"Stable baseline: {x['stable_baseline']}\n"
        f"Previous Stable: {x['previous_stable']}\n"
        f"Scope: {x['scope']}\n"
        f"Application bundle: {x['application_bundle']}\n"
        "Deterministic Day/Swing/Long formulas: unchanged / v14.3.7-compatible\n"
        f"Runtime/config continuity: {x['runtime_config']} (preserves compatible prior Stable settings and API keys)\n"
        "Release learning: adaptive G0-G16 Release Learning Registry active · streamlined one-branch/one-PR exact-head Fast + Qualified + one Release G11-G16 lifecycle\n"
        "Adaptive contracts: US Equities Processing Boundary · Data Utility/Evidence Value · Performance/Scalability · Data Health/Freshness/Cache · Testing · Intelligence/Learning preserved\n"
    )
    path = ROOT / "app_bootstrap.go"
    text = path.read_text()
    text = re.sub(r'const appVersion = "[^"]+"', f'const appVersion = "{x["version"]}"', text)
    text = re.sub(r'const buildID = "[^"]+"', f'const buildID = "{x["build_id"]}"', text)
    text = re.sub(r'const releaseChannel = "[^"]+"', f'const releaseChannel = "{x["channel"]}"', text)
    path.write_text(text)

    patch = patch_contract(x)
    contract = release_contract(x)
    overlay = overlay_asset_path(x, contract)
    if patch:
        path = ROOT / "renderer" / "watchlist-desk-contract.js"
        if not path.exists():
            raise SystemExit(f"patch identity asset missing: {path}")
        text = path.read_text()
        text = re.sub(r"DEPULSE_PATCH_VERSION = '[^']+'", f"DEPULSE_PATCH_VERSION = '{x['version']}'", text)
        text = re.sub(r"DEPULSE_PATCH_BUILD_ID = '[^']+'", f"DEPULSE_PATCH_BUILD_ID = '{x['build_id']}'", text)
        path.write_text(text)
    elif overlay:
        if not overlay.exists():
            raise SystemExit(f"release identity overlay missing: {overlay}")
        text = overlay.read_text()
        text = re.sub(r"DEPULSE_RELEASE_VERSION = '[^']+'", f"DEPULSE_RELEASE_VERSION = '{x['version']}'", text)
        text = re.sub(r"DEPULSE_RELEASE_BUILD_ID = '[^']+'", f"DEPULSE_RELEASE_BUILD_ID = '{x['build_id']}'", text)
        overlay.write_text(text)
    else:
        path = ROOT / "renderer" / "renderer.js"
        text = path.read_text()
        text = re.sub(r"const EXPECTED_RELEASE_VERSION='[^']+';", f"const EXPECTED_RELEASE_VERSION='{x['version']}';", text)
        text = re.sub(r"const EXPECTED_BUILD_ID='[^']+';", f"const EXPECTED_BUILD_ID='{x['build_id']}';", text)
        path.write_text(text)

    path = ROOT / "renderer" / "index.html"
    text = path.read_text()
    text = re.sub(r"<title>DE\.PULSE v[^<]+</title>", f"<title>DE.PULSE v{x['version']}</title>", text)
    for asset in RELEASE_COUPLED_ASSETS:
        text = sync_asset_cache_identity(text, asset)
    if overlay:
        text = sync_asset_cache_identity(text, overlay.name)
    if patch:
        text = sync_asset_cache_identity(text, "watchlist-desk-contract.js")
    path.write_text(text)

    legacy_cert, legacy_ci = legacy_registry_versions(x, patch, contract)
    if not patch and not legacy_cert and not legacy_ci:
        for kind in ("certification", "ci_pipeline"):
            target = legacy_registry_path(kind, x["version"])
            data = json.loads(target.read_text())
            data["version"] = x["version"]
            if kind == "ci_pipeline":
                data.setdefault("policy", {})["baseline"] = x["previous_stable"] + " Stable"
                data["policy"]["release_channel"] = x["channel"]
                data["policy"]["canonical_release_identity"] = "release_identity.json"
                data["policy"]["pre_freeze_qualification"] = True
                data["policy"]["unique_test_evidence"] = True
            target.write_text(json.dumps(data, indent=2) + "\n")


def verify_cache_identity(index: str, name: str) -> bool:
    return f"{name}?v={asset_cache_token(name, canonical_git=True)}" in index


def verify(x):
    errs = []
    version_text = (ROOT / "VERSION.txt").read_text()
    boot = (ROOT / "app_bootstrap.go").read_text()
    renderer = (ROOT / "renderer" / "renderer.js").read_text()
    index = (ROOT / "renderer" / "index.html").read_text()
    patch = patch_contract(x)
    contract = release_contract(x)
    overlay = overlay_asset_path(x, contract)
    legacy_cert, legacy_ci = legacy_registry_versions(x, patch, contract)
    cert = json.loads(legacy_registry_path("certification", legacy_cert or x["version"]).read_text())
    ci = json.loads(legacy_registry_path("ci_pipeline", legacy_ci or x["version"]).read_text())

    renderer_ok = (
        f"const EXPECTED_RELEASE_VERSION='{x['version']}';" in renderer
        and f"const EXPECTED_BUILD_ID='{x['build_id']}';" in renderer
    )
    legacy_plans = False
    if patch:
        patch_name = "watchlist-desk-contract.js"
        patch_asset = (ROOT / "renderer" / patch_name).read_text()
        renderer_ok = (
            f"DEPULSE_PATCH_VERSION = '{x['version']}'" in patch_asset
            and f"DEPULSE_PATCH_BUILD_ID = '{x['build_id']}'" in patch_asset
            and verify_cache_identity(index, patch_name)
        )
        cert_ok = str(cert.get("version")) == str(patch.get("inherited_certification_plan"))
        ci_ok = str(ci.get("version")) == str(patch.get("inherited_ci_plan"))
        legacy_plans = True
    elif overlay:
        if not overlay.exists():
            raise SystemExit(f"release identity overlay missing: {overlay}")
        overlay_text = overlay.read_text()
        renderer_ok = (
            f"DEPULSE_RELEASE_VERSION = '{x['version']}'" in overlay_text
            and f"DEPULSE_RELEASE_BUILD_ID = '{x['build_id']}'" in overlay_text
            and verify_cache_identity(index, overlay.name)
        )
        cert_ok = str(cert.get("version")) == (legacy_cert or str(x["version"]))
        ci_ok = str(ci.get("version")) == (legacy_ci or str(x["version"]))
        legacy_plans = bool(legacy_cert or legacy_ci)
    else:
        cert_ok = str(cert.get("version")) == str(x["version"])
        ci_ok = str(ci.get("version")) == str(x["version"])

    checks = [
        (x["display_version"] in version_text, "VERSION display"),
        (x["build_id"] in version_text, "VERSION build"),
        (f'Previous Stable: {x["previous_stable"]}' in version_text, "VERSION predecessor"),
        (f'const appVersion = "{x["version"]}"' in boot, "appVersion"),
        (f'const buildID = "{x["build_id"]}"' in boot, "buildID"),
        (f'const releaseChannel = "{x["channel"]}"' in boot, "release channel"),
        (renderer_ok, "renderer/overlay identity"),
        (f"<title>DE.PULSE v{x['version']}</title>" in index, "HTML title"),
        (cert_ok, "certification plan inheritance/version"),
        (ci_ok, "CI plan inheritance/version"),
    ]
    if not patch and not legacy_plans:
        checks.append((ci.get("policy", {}).get("baseline") == x["previous_stable"] + " Stable", "CI baseline"))
    checks.extend(
        (verify_cache_identity(index, asset), f"{asset} content-derived cache identity")
        for asset in RELEASE_COUPLED_ASSETS
    )
    if f"?v={x['version']}" in index:
        checks.append((False, "renderer cache identity must not use public version"))
    errs.extend(label for ok, label in checks if not ok)
    if errs:
        raise SystemExit("Release identity: FAIL · " + ", ".join(errs))
    mode = "overlay" if overlay else ("patch" if patch else "monolith")
    print(f"Release identity: PASS · {x['version']} · {x['build_id']} · renderer identity={mode} · cache identity=content-derived")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--sync", action="store_true")
    parser.add_argument("--verify", action="store_true")
    args = parser.parse_args()
    identity = load()
    if args.sync:
        sync(identity)
    if args.verify or not args.sync:
        verify(identity)


if __name__ == "__main__":
    main()
