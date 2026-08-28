#!/usr/bin/env python3
"""Fail-closed HOST-014 external egress conservation gate.

Every literal external http/https/wss destination in production Go source must be
classified in the canonical hosted egress inventory. Runtime authorization is
separate: only inventory rows with runtimeAllowed=true and an allowed protocol
(https/wss) may be emitted to the hosted trust renderer.
"""
from __future__ import annotations

import json
import re
from pathlib import Path
from urllib.parse import urlsplit

ROOT = Path(__file__).resolve().parents[2]
INVENTORY = ROOT / "governance" / "hosted-infrastructure" / "external-egress-v1.json"
URL_RE = re.compile(r"(?P<scheme>https?|wss)://(?P<host>[A-Za-z0-9.-]+)(?::\d+)?(?:[/\"'?]|$)", re.IGNORECASE)
IGNORED_HOSTS = {"localhost", "127.0.0.1", "::1", "example.invalid"}


def fail(message: str) -> None:
    print("FAIL:", message)
    raise SystemExit(1)


def normalize_host(value: str) -> str:
    host = value.strip().lower().rstrip(".")
    if not host or "*" in host or "/" in host or " " in host:
        fail("invalid egress host: " + value)
    if host.startswith(".") or host.endswith("."):
        fail("invalid egress host boundary: " + value)
    return host


def production_go_files() -> list[Path]:
    return sorted(
        p for p in ROOT.rglob("*.go")
        if not p.name.endswith("_test.go")
        and ".git" not in p.parts
        and "vendor" not in p.parts
    )


def discovered_urls() -> dict[tuple[str, str], set[str]]:
    found: dict[tuple[str, str], set[str]] = {}
    for path in production_go_files():
        text = path.read_text(encoding="utf-8", errors="strict")
        rel = path.relative_to(ROOT).as_posix()
        for match in URL_RE.finditer(text):
            scheme = match.group("scheme").lower()
            host = normalize_host(match.group("host"))
            if host in IGNORED_HOSTS or host.endswith(".example.invalid"):
                continue
            found.setdefault((host, scheme), set()).add(rel)
    return found


def load_inventory() -> tuple[dict, dict[str, dict]]:
    try:
        payload = json.loads(INVENTORY.read_text(encoding="utf-8"))
    except Exception as exc:
        fail("egress inventory unreadable: " + str(exc))
    if payload.get("schema") != "DE.PULSE-HOSTED-EXTERNAL-EGRESS-1":
        fail("unexpected egress inventory schema")
    if payload.get("version") != "v1":
        fail("unexpected egress inventory version")
    if str(payload.get("defaultPolicy", "")).upper() != "DENY":
        fail("hosted external egress must default DENY")
    allowed = payload.get("allowedProtocols")
    if allowed != ["https", "wss"]:
        fail("allowedProtocols must remain exactly https,wss")
    rules = payload.get("rules")
    if not isinstance(rules, list) or not rules:
        fail("egress inventory rules missing")
    by_host: dict[str, dict] = {}
    for row in rules:
        if not isinstance(row, dict):
            fail("egress rule must be an object")
        host = normalize_host(str(row.get("host", "")))
        if host in by_host:
            fail("duplicate egress host: " + host)
        protocols = row.get("protocols")
        if not isinstance(protocols, list) or not protocols:
            fail("egress protocols missing for " + host)
        protocols = [str(p).lower() for p in protocols]
        if len(set(protocols)) != len(protocols):
            fail("duplicate protocol for " + host)
        runtime_allowed = row.get("runtimeAllowed") is True
        if runtime_allowed:
            if any(p not in allowed for p in protocols):
                fail("runtime-authorized host uses forbidden protocol: " + host)
            if not str(row.get("consumer", "")).strip():
                fail("runtime-authorized host missing consumer: " + host)
        else:
            if not str(row.get("blockReason", "")).strip():
                fail("blocked host missing blockReason: " + host)
        source_paths = row.get("sourcePaths")
        if not isinstance(source_paths, list) or not source_paths:
            fail("egress sourcePaths missing for " + host)
        for rel in source_paths:
            if not (ROOT / str(rel)).is_file():
                fail(f"egress trace path missing for {host}: {rel}")
        row["host"] = host
        row["protocols"] = protocols
        by_host[host] = row
    return payload, by_host


def main() -> int:
    payload, by_host = load_inventory()
    found = discovered_urls()
    missing: list[str] = []
    blocked_seen: list[str] = []
    for (host, scheme), paths in sorted(found.items()):
        row = by_host.get(host)
        if row is None or scheme not in row.get("protocols", []):
            missing.append(f"{scheme}://{host} <- {', '.join(sorted(paths))}")
            continue
        if row.get("runtimeAllowed") is not True:
            blocked_seen.append(f"{scheme}://{host} <- {', '.join(sorted(paths))}")
    if missing:
        fail("unclassified production external destinations:\n  " + "\n  ".join(missing))

    runtime_hosts = sorted(
        host for host, row in by_host.items()
        if row.get("runtimeAllowed") is True
        and any(proto in payload["allowedProtocols"] for proto in row.get("protocols", []))
    )
    if not runtime_hosts:
        fail("runtime egress allowlist is empty")
    if any(host in {"0.0.0.0", "::"} for host in runtime_hosts):
        fail("broad network destination present in runtime allowlist")

    print(f"HOST-014 external egress conservation: PASS · {len(found)} source destinations classified · {len(runtime_hosts)} runtime hosts")
    if blocked_seen:
        print("KNOWN BLOCKED LEGACY DESTINATIONS (not authorized for hosted runtime):")
        for item in blocked_seen:
            print("-", item)
    print("RUNTIME HOSTS:")
    for host in runtime_hosts:
        print(host)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
