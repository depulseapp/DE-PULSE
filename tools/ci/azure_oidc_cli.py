#!/usr/bin/env python3
"""Refresh Azure CLI authentication from GitHub Actions OIDC without secret material.

Long-running managed-infrastructure operations can outlive the short GitHub OIDC
assertion used by an earlier azure/login step. This helper requests a fresh assertion
from the GitHub Actions OIDC endpoint and exchanges it through `az login` immediately
before Azure CLI work. The assertion and request token are never printed or persisted.
"""
from __future__ import annotations

import json
import os
import subprocess
from urllib.parse import parse_qsl, urlencode, urlsplit, urlunsplit
from urllib.request import Request, urlopen

AUDIENCE = "api://AzureADTokenExchange"


def oidc_request_url(base_url: str, audience: str = AUDIENCE) -> str:
    parts = urlsplit(base_url)
    query = dict(parse_qsl(parts.query, keep_blank_values=True))
    query["audience"] = audience
    return urlunsplit((parts.scheme, parts.netloc, parts.path, urlencode(query), parts.fragment))


def extract_oidc_value(payload: object) -> str:
    if not isinstance(payload, dict):
        raise RuntimeError("GitHub OIDC response must be an object")
    value = str(payload.get("value") or "").strip()
    if not value:
        raise RuntimeError("GitHub OIDC response did not contain an assertion")
    return value


def _required_env(name: str) -> str:
    value = str(os.environ.get(name) or "").strip()
    if not value:
        raise RuntimeError(f"required GitHub/Azure OIDC environment value is missing: {name}")
    return value


def _fresh_github_assertion() -> str:
    request_url = oidc_request_url(_required_env("ACTIONS_ID_TOKEN_REQUEST_URL"))
    request_token = _required_env("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
    request = Request(
        request_url,
        headers={"Authorization": "Bearer " + request_token, "Accept": "application/json"},
        method="GET",
    )
    try:
        with urlopen(request, timeout=20) as response:
            payload = json.load(response)
    except Exception as exc:
        raise RuntimeError("failed to request a fresh GitHub Actions OIDC assertion") from exc
    return extract_oidc_value(payload)


def refresh_azure_cli_oidc() -> None:
    """Renew the Azure CLI session using a just-in-time GitHub OIDC assertion."""
    client_id = _required_env("AZURE_CLIENT_ID")
    tenant_id = _required_env("AZURE_TENANT_ID")
    subscription_id = _required_env("AZURE_SUBSCRIPTION_ID")
    assertion = _fresh_github_assertion()

    login = subprocess.run(
        [
            "az", "login",
            "--service-principal",
            "--username", client_id,
            "--tenant", tenant_id,
            "--federated-token", assertion,
            "--only-show-errors",
            "--output", "none",
        ],
        text=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.PIPE,
        check=False,
    )
    assertion = ""  # Drop the local reference as soon as the exchange completes.
    if login.returncode != 0:
        raise RuntimeError(f"Azure CLI OIDC refresh failed with exit code {login.returncode}")

    selected = subprocess.run(
        ["az", "account", "set", "--subscription", subscription_id, "--only-show-errors"],
        text=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.PIPE,
        check=False,
    )
    if selected.returncode != 0:
        raise RuntimeError(f"Azure CLI subscription selection failed with exit code {selected.returncode}")


def self_test() -> None:
    url = oidc_request_url("https://example.invalid/oidc?api-version=1")
    if "api-version=1" not in url or "audience=api%3A%2F%2FAzureADTokenExchange" not in url:
        raise AssertionError("OIDC request URL construction failed")
    if extract_oidc_value({"value": "header.payload.signature"}) != "header.payload.signature":
        raise AssertionError("OIDC response extraction failed")
    try:
        extract_oidc_value({})
    except RuntimeError:
        pass
    else:
        raise AssertionError("empty OIDC response did not fail closed")
    print("Azure GitHub OIDC CLI refresh self-test: PASS")


if __name__ == "__main__":
    self_test()
