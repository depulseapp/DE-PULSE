#!/usr/bin/env python3
"""Render the HOST-017/018 Azure Key Vault CSI managed-secret contract.

The renderer accepts only non-secret version references. Secret values are never
accepted, printed, stored, or serialized by this contract.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import re
from pathlib import Path

ENVIRONMENTS = ("dev", "test", "stage", "prod")
UUID_RE = re.compile(r"^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
KEY_VAULT_RE = re.compile(r"^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$")
OBJECT_VERSION_RE = re.compile(r"^[0-9a-fA-F]{32}$")
REFERENCE_SCHEMA = "DE.PULSE-HOSTED-MANAGED-SECRET-REFERENCES-1"
SECRET_CONTRACT = "v1"
SECRET_ALIASES = (
    "finnhub",
    "tradeinsight",
    "alpaca-key",
    "alpaca-secret",
    "groq",
    "openrouter",
    "gemini",
    "fred",
    "bls",
    "eia",
    "twelvedata",
    "marketaux",
)


def q(value: str) -> str:
    return json.dumps(value)


def validate_uuid(name: str, value: str) -> str:
    if not UUID_RE.fullmatch(value):
        raise SystemExit(f"{name} must be a concrete UUID")
    return value


def load_reference_manifest(path: Path, environment: str) -> list[tuple[str, str]]:
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        raise SystemExit(f"managed-secret reference manifest unreadable: {exc}") from exc
    if not isinstance(payload, dict) or payload.get("schema") != REFERENCE_SCHEMA:
        raise SystemExit("unsupported managed-secret reference manifest schema")
    if set(payload) != {"schema", "environment", "references"}:
        raise SystemExit("managed-secret reference manifest may contain reference metadata only")
    if payload.get("environment") != environment:
        raise SystemExit("managed-secret reference manifest environment mismatch")
    references = payload.get("references")
    if not isinstance(references, dict) or not references:
        raise SystemExit("managed-secret reference manifest must contain at least one reference")
    unknown = sorted(set(references) - set(SECRET_ALIASES))
    if unknown:
        raise SystemExit("managed-secret reference manifest contains unsupported aliases: " + ",".join(unknown))

    ordered: list[tuple[str, str]] = []
    for alias in SECRET_ALIASES:
        if alias not in references:
            continue
        version = references[alias]
        if not isinstance(version, str) or not OBJECT_VERSION_RE.fullmatch(version):
            raise SystemExit(f"managed-secret reference {alias!r} must be an exact 32-character Key Vault object version")
        ordered.append((alias, version.lower()))
    return ordered


def reference_generation(environment: str, references: list[tuple[str, str]]) -> str:
    payload = {
        "schema": REFERENCE_SCHEMA,
        "environment": environment,
        "references": [{"alias": alias, "objectVersion": version} for alias, version in references],
    }
    encoded = json.dumps(payload, separators=(",", ":"), sort_keys=True).encode("utf-8")
    return hashlib.sha256(encoded).hexdigest()


def required_aliases(references: list[tuple[str, str]]) -> str:
    return ",".join(alias for alias, _ in references)


def secret_provider_class_name(generation: str) -> str:
    return "depulse-managed-secrets-" + generation[:12]


def render(
    environment: str,
    key_vault_name: str,
    tenant_id: str,
    workload_identity_client_id: str,
    references: list[tuple[str, str]],
) -> str:
    if environment not in ENVIRONMENTS:
        raise SystemExit("unsupported environment")
    if not KEY_VAULT_RE.fullmatch(key_vault_name):
        raise SystemExit("key vault name must satisfy Azure Key Vault naming rules")
    validate_uuid("tenant id", tenant_id)
    validate_uuid("workload identity client id", workload_identity_client_id)
    if not references:
        raise SystemExit("managed-secret references must not be empty")

    generation = reference_generation(environment, references)
    namespace = f"depulse-{environment}"
    provider_class = secret_provider_class_name(generation)
    objects = []
    for alias, version in references:
        objects.append(
            "        - |\n"
            f"          objectName: {alias}\n"
            "          objectType: secret\n"
            f"          objectVersion: {q(version)}\n"
            f"          objectAlias: {alias}"
        )
    object_block = "\n".join(objects)
    return f"""apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: {provider_class}
  namespace: {namespace}
  labels:
    depulse.io/secret-contract: {q(SECRET_CONTRACT)}
  annotations:
    depulse.io/secret-generation-sha256: {q(generation)}
spec:
  provider: azure
  parameters:
    usePodIdentity: {q('false')}
    clientID: {q(workload_identity_client_id)}
    keyvaultName: {q(key_vault_name)}
    cloudName: {q('')}
    objects: |
      array:
{object_block}
    tenantId: {q(tenant_id)}
"""


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--environment", required=True, choices=ENVIRONMENTS)
    parser.add_argument("--key-vault-name", required=True)
    parser.add_argument("--tenant-id", required=True)
    parser.add_argument("--workload-identity-client-id", required=True)
    parser.add_argument("--reference-manifest", required=True, type=Path)
    args = parser.parse_args()
    references = load_reference_manifest(args.reference_manifest, args.environment)
    print(render(args.environment, args.key_vault_name, args.tenant_id, args.workload_identity_client_id, references), end="")


if __name__ == "__main__":
    main()
