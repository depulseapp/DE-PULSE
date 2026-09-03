#!/usr/bin/env python3
"""Render the HOST-017/018 Azure Key Vault CSI managed-secret contract.

The renderer emits references and aliases only. It never accepts, prints, stores,
or serializes secret values.
"""
from __future__ import annotations

import argparse
import json
import re

ENVIRONMENTS = ("dev", "test", "stage", "prod")
UUID_RE = re.compile(r"^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
KEY_VAULT_RE = re.compile(r"^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$")
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


def render(environment: str, key_vault_name: str, tenant_id: str, workload_identity_client_id: str) -> str:
    if environment not in ENVIRONMENTS:
        raise SystemExit("unsupported environment")
    if not KEY_VAULT_RE.fullmatch(key_vault_name):
        raise SystemExit("key vault name must satisfy Azure Key Vault naming rules")
    validate_uuid("tenant id", tenant_id)
    validate_uuid("workload identity client id", workload_identity_client_id)

    namespace = f"depulse-{environment}"
    objects = []
    for alias in SECRET_ALIASES:
        objects.append(
            "        - |\n"
            f"          objectName: {alias}\n"
            "          objectType: secret\n"
            "          objectVersion: \"\"\n"
            f"          objectAlias: {alias}"
        )
    object_block = "\n".join(objects)
    return f"""apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: depulse-managed-secrets
  namespace: {namespace}
  labels:
    depulse.io/secret-contract: {q('v1')}
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
    args = parser.parse_args()
    print(render(args.environment, args.key_vault_name, args.tenant_id, args.workload_identity_client_id), end="")


if __name__ == "__main__":
    main()
