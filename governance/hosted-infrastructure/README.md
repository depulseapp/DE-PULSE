# DE.PULSE Hosted Infrastructure Owner

This directory is the repository-side infrastructure owner for v19 HOST-013..014. It **projects** the canonical application desired state from `internal/hostedenv/desired_state_v1.json`; it does not replace or fork that authority.

## Current implementation boundary

`tools/hosted/render_kubernetes_trust.py` renders deployable Kubernetes + Istio service-trust resources for exactly one canonical environment (`dev`, `test`, `stage`, `prod`). The renderer enforces:

- one environment namespace equal to the canonical `isolationId`;
- one Kubernetes ServiceAccount equal to the canonical `serviceIdentity`;
- namespace-wide default-deny ingress and egress;
- ingress to the DE.PULSE workload only from the managed mesh ingress namespace on application port 8080;
- DNS egress required for service discovery;
- HTTPS/WSS transport egress on TCP/443, while Istio outbound traffic remains `REGISTRY_ONLY` and only explicitly supplied external hosts are registered;
- Istio `PeerAuthentication` in `STRICT` mode for internal mTLS;
- an `AuthorizationPolicy` binding inbound service traffic to the managed ingress service-account principal;
- fail-closed rendering when no explicit external-host inventory is supplied;
- annotations binding every rendered object to the canonical desired-state version and SHA-256 digest.

The renderer deliberately does **not** create cloud resources, clusters, DNS, certificates, load balancers, workload-identity federation, KMS/secrets or provider-specific firewalls. Those require a selected managed hosting platform and real operator evidence. Repository rendering therefore advances HOST-013..014 implementation but does not by itself make the band VERIFIED.

## Render

```bash
python3 tools/hosted/render_kubernetes_trust.py \
  --environment dev \
  --egress-hosts-file /secure/operator/canonical-egress-hosts.txt \
  --output hosted-dev.yaml
```

The egress file is newline-delimited. Blank lines and `#` comments are ignored. Entries must be concrete DNS hostnames; wildcards, URLs, IP literals and ports are rejected. Keep credentials and secrets out of this file.

## Verification contract

`python3 tools/ci/hosted_infrastructure_contract_gate.py` is the repository gate. It proves deterministic projection, fail-closed egress handling, unique environment/service identities, strict mTLS policy, default-deny network posture, managed-ingress authorization and canonical desired-state digest binding.

A future HOST-013..014 VERIFIED transition additionally requires real managed-environment evidence showing the rendered controls are actually applied and drift detection fails closed in the selected hosting platform.
