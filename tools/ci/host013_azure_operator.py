#!/usr/bin/env python3
"""Governed HOST-013..014 Azure AKS operator.

This is the single repository owner for the future manual Azure infrastructure run.
It assumes GitHub/Azure OIDC authentication is already established and Terraform is
installed. It never accepts a client secret, storage key, SAS token or kubeconfig.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
AZURE_DIR = ROOT / "governance" / "hosted-infrastructure" / "azure"
RENDERER = ROOT / "tools" / "hosted" / "render_kubernetes_trust.py"
LIVE_EVIDENCE = ROOT / "tools" / "ci" / "host013_azure_live_evidence.py"
TRAFFIC_PROBE = ROOT / "tools" / "ci" / "host013_azure_traffic_probe.py"
CONFIRMATION = "HOST013_AZURE_AKS_OPERATOR_DRILL"
ISTIO_RE = re.compile(r"^asm-(\d+)-(\d+)$")


def fail(message: str) -> None:
    raise SystemExit("HOST-013/014 Azure operator: " + message)


def run(args: list[str], *, env: dict[str, str] | None = None, capture: bool = False) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        args,
        cwd=ROOT,
        env=env,
        text=True,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.STDOUT if capture else None,
        check=False,
    )
    if result.returncode != 0:
        detail = (result.stdout or "").strip() if capture else ""
        fail(f"command failed ({result.returncode}): {' '.join(args)}" + ("\n" + detail if detail else ""))
    return result


def run_json(args: list[str]) -> object:
    result = run(args, capture=True)
    try:
        return json.loads(result.stdout or "null")
    except json.JSONDecodeError as exc:
        fail(f"non-JSON command response from {' '.join(args)}: {exc}")


def revision_names(value: object) -> list[str]:
    names: set[str] = set()

    def walk(node: object) -> None:
        if isinstance(node, str) and ISTIO_RE.fullmatch(node):
            names.add(node)
        elif isinstance(node, list):
            for item in node:
                walk(item)
        elif isinstance(node, dict):
            for item in node.values():
                walk(item)

    walk(value)
    return sorted(names, key=lambda item: tuple(int(v) for v in ISTIO_RE.fullmatch(item).groups()))


def choose_revision(value: object) -> str:
    names = revision_names(value)
    if not names:
        fail("Azure returned no supported managed Istio asm-X-Y revision")
    return names[-1]


def terraform_env(args: argparse.Namespace, istio_revision: str) -> dict[str, str]:
    env = os.environ.copy()
    env.update({
        "ARM_USE_OIDC": "true",
        "ARM_USE_AZUREAD": "true",
        "ARM_CLIENT_ID": args.client_id,
        "ARM_TENANT_ID": args.tenant_id,
        "ARM_SUBSCRIPTION_ID": args.subscription_id,
        "TF_VAR_subscription_id": args.subscription_id,
        "TF_VAR_tenant_id": args.tenant_id,
        "TF_VAR_environment": args.environment,
        "TF_VAR_location": args.location,
        "TF_VAR_istio_revisions": json.dumps([istio_revision]),
        "TF_IN_AUTOMATION": "true",
    })
    return env


def terraform(args: argparse.Namespace, env: dict[str, str], *extra: str, capture: bool = False) -> subprocess.CompletedProcess[str]:
    return run(["terraform", f"-chdir={AZURE_DIR}", *extra], env=env, capture=capture)


def terraform_outputs(args: argparse.Namespace, env: dict[str, str]) -> dict[str, object]:
    result = terraform(args, env, "output", "-json", capture=True)
    try:
        payload = json.loads(result.stdout or "{}")
    except json.JSONDecodeError as exc:
        fail("Terraform output was not valid JSON: " + str(exc))
    if not isinstance(payload, dict):
        fail("Terraform output JSON must be an object")
    return payload


def output_value(outputs: dict[str, object], name: str) -> str:
    row = outputs.get(name)
    if not isinstance(row, dict):
        fail("missing Terraform output: " + name)
    value = str(row.get("value") or "").strip()
    if not value:
        fail("empty Terraform output: " + name)
    return value


def fingerprint(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()[:16]


def require_evidence(path: Path, expected_schema: str, expected_status: str) -> dict[str, object]:
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        fail(f"evidence unreadable {path.name}: {exc}")
    if not isinstance(payload, dict):
        fail(f"evidence must be an object: {path.name}")
    if payload.get("schema") != expected_schema or payload.get("status") != expected_status:
        fail(f"evidence did not reach required state: {path.name}")
    if payload.get("containsSecrets") is not False:
        fail(f"evidence must explicitly be secret-free: {path.name}")
    return payload


def self_test() -> None:
    fixture = {
        "revisions": [
            {"name": "asm-1-26"},
            {"revision": "asm-1-27"},
            {"note": "not-a-revision"},
        ]
    }
    if choose_revision(fixture) != "asm-1-27":
        fail("managed Istio revision selection self-test failed")
    try:
        choose_revision({"revisions": []})
    except SystemExit:
        pass
    else:
        fail("empty managed Istio revision set did not fail closed")
    if output_value({"x": {"value": "abc"}}, "x") != "abc":
        fail("Terraform output extraction self-test failed")
    try:
        output_value({}, "missing")
    except SystemExit:
        pass
    else:
        fail("missing Terraform output did not fail closed")
    print("HOST-013/014 Azure operator self-test: PASS")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--confirmation", default="")
    parser.add_argument("--candidate-sha", default="")
    parser.add_argument("--client-id", default="")
    parser.add_argument("--tenant-id", default="")
    parser.add_argument("--subscription-id", default="")
    parser.add_argument("--state-resource-group", default="")
    parser.add_argument("--state-storage-account", default="")
    parser.add_argument("--state-container", default="")
    parser.add_argument("--environment", default="dev", choices=["dev"])
    parser.add_argument("--location", default="canadacentral")
    parser.add_argument("--evidence-dir", default=".depulse-host013-azure")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return 0
    if args.confirmation != CONFIRMATION:
        fail(f"explicit confirmation must equal {CONFIRMATION}")
    required = (
        "candidate_sha", "client_id", "tenant_id", "subscription_id",
        "state_resource_group", "state_storage_account", "state_container",
    )
    for name in required:
        if not str(getattr(args, name)).strip():
            fail("missing required non-secret identifier --" + name.replace("_", "-"))
    head = run(["git", "rev-parse", "HEAD"], capture=True).stdout.strip()
    if head != args.candidate_sha:
        fail(f"exact-head mismatch: repository={head} requested={args.candidate_sha}")

    account = run_json(["az", "account", "show", "-o", "json"])
    if not isinstance(account, dict) or str(account.get("id") or "") != args.subscription_id:
        fail("active Azure subscription does not match requested subscription")

    revision_payload = run_json(["az", "aks", "mesh", "get-revisions", "--location", args.location, "-o", "json"])
    istio_revision = choose_revision(revision_payload)
    env = terraform_env(args, istio_revision)
    state_key = f"host013/{args.environment}/terraform.tfstate"

    terraform(args, env, "fmt", "-check", "-diff")
    terraform(
        args, env, "init", "-input=false", "-reconfigure",
        f"-backend-config=resource_group_name={args.state_resource_group}",
        f"-backend-config=storage_account_name={args.state_storage_account}",
        f"-backend-config=container_name={args.state_container}",
        f"-backend-config=key={state_key}",
    )
    terraform(args, env, "validate")

    evidence_dir = ROOT / args.evidence_dir
    evidence_dir.mkdir(parents=True, exist_ok=True)
    plan_path = evidence_dir / "host013-azure.tfplan"
    terraform(args, env, "plan", "-input=false", "-out=" + str(plan_path))
    terraform(args, env, "apply", "-input=false", "-auto-approve", str(plan_path))

    outputs = terraform_outputs(args, env)
    workload_client_id = output_value(outputs, "workload_identity_client_id")
    workload_subject = output_value(outputs, "workload_identity_subject")
    expected_subject = f"system:serviceaccount:depulse-{args.environment}:depulse-web-{args.environment}"
    if workload_subject != expected_subject:
        fail(f"workload identity subject drift: expected={expected_subject} actual={workload_subject}")

    rg = f"rg-depulse-{args.environment}"
    cluster = f"aks-depulse-{args.environment}"
    manifest_path = evidence_dir / "host013-aks-managed-trust.yaml"
    run([
        sys.executable, str(RENDERER),
        "--environment", args.environment,
        "--mesh-profile", "aks-managed",
        "--istio-revision", istio_revision,
        "--workload-identity-client-id", workload_client_id,
        "--output", str(manifest_path),
    ])
    run([
        "az", "aks", "command", "invoke", "--resource-group", rg,
        "--name", cluster, "--command", "kubectl apply -f host013-aks-managed-trust.yaml",
        "--file", str(manifest_path), "-o", "json",
    ], capture=True)

    evidence_path = evidence_dir / "host013-azure-live-evidence.json"
    run([
        sys.executable, str(LIVE_EVIDENCE),
        "--subscription-id", args.subscription_id,
        "--resource-group", rg,
        "--cluster-name", cluster,
        "--environment", args.environment,
        "--location", args.location,
        "--candidate-sha", args.candidate_sha,
        "--expected-istio-revision", istio_revision,
        "--workload-identity-client-id", workload_client_id,
        "--output", str(evidence_path),
    ])
    require_evidence(evidence_path, "DE.PULSE-HOST013-AZURE-LIVE-EVIDENCE-1", "PASS_CONFIGURATION_AND_IDENTITY")

    traffic_path = evidence_dir / "host013-azure-traffic-evidence.json"
    run([
        sys.executable, str(TRAFFIC_PROBE),
        "--resource-group", rg,
        "--cluster-name", cluster,
        "--environment", args.environment,
        "--candidate-sha", args.candidate_sha,
        "--istio-revision", istio_revision,
        "--output", str(traffic_path),
    ])
    traffic = require_evidence(traffic_path, "DE.PULSE-HOST013-AZURE-TRAFFIC-EVIDENCE-1", "PASS")
    traffic_checks = traffic.get("checks")
    if not isinstance(traffic_checks, dict) or not traffic_checks or not all(value is True for value in traffic_checks.values()):
        fail("traffic probe evidence is incomplete")

    drift = subprocess.run(
        ["terraform", f"-chdir={AZURE_DIR}", "plan", "-input=false", "-detailed-exitcode"],
        cwd=ROOT, env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=False,
    )
    if drift.returncode != 0:
        fail(f"post-verification Terraform drift check returned {drift.returncode}; expected 0")

    summary = {
        "schema": "DE.PULSE-HOST013-AZURE-OPERATOR-1",
        "candidateSha": args.candidate_sha,
        "environment": args.environment,
        "location": args.location,
        "managedIstioRevision": istio_revision,
        "workloadIdentityClientFingerprint": fingerprint(workload_client_id),
        "workloadIdentitySubject": workload_subject,
        "stateKey": state_key,
        "terraformPostApplyDriftExitCode": drift.returncode,
        "configurationIdentityEvidence": evidence_path.name,
        "trafficEvidence": traffic_path.name,
        "containsSecrets": False,
        "status": "PASS",
    }
    (evidence_dir / "host013-azure-operator-summary.json").write_text(
        json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    plan_path.unlink(missing_ok=True)
    manifest_path.unlink(missing_ok=True)
    print("HOST-013/014 Azure operator: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
