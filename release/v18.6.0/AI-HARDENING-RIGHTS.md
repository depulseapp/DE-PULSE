# AI Bounded-Context / Cache / Schema / Eval / Rights Hardening — v18.6.0

**Requirements:** `AUDIT-18-AI-001`, `AUDIT-18-AI-RIGHTS-001`  
**Gate model:** existing `G0–G16` only  
**Eval policy:** `ai_eval_policy.json`  
**Canonical provider × dataset AI-rights source:** `provider_dataset_ai_rights_registry.json`

## Production-path changes

The v18.6 AI path now performs rights filtering before prompt construction and before any external model request. A provider/dataset evidence item is eligible for egress only when commercial-use, redistribution and AI-use rights are all explicitly `APPROVED`, the evidence decision is bound, and the row decision is `ALLOW`. Missing, unknown, blocked or denied records fail closed. A working API key or endpoint response is not rights evidence.

The current registry deliberately contains no inferred approvals. Unapproved evidence is removed from the outbound package, role/citation ID lists are rebuilt, a missing-evidence/safety warning is added, and diagnostics expose only counts/policy state rather than blocked source text, URLs or credentials.

Context construction is semantic/materiality-aware. Risk, event, catalyst, bull/bear and market evidence receive deterministic priority before low-value base context. The outbound JSON envelope is repeatedly compacted while retaining valid JSON until both the hard byte bound and conservative provider-independent token upper bound are satisfied. Oversized user task text is bounded separately.

Inference cache identity now includes evidence/egress package identity, evidence snapshot, request kind/task/scope/ticker, route policy, provider/model/fallback configuration, output-token route limit, prompt fingerprint, safety policy, schema policy/fingerprint, context policy and rights-egress policy. Cache entries expire after 15 minutes and clock-skew/future entries fail invalidation rather than being trusted indefinitely. Existing bounded 256-entry storage remains in force.

Production output no longer treats malformed/unstructured model text as a successful informational response. Required fields, unknown fields, enums, list sizes, string sizes and evidence citations are strictly validated. Unknown/duplicate citation IDs fail. Invalid structured output may fall through to another configured provider; if no candidate produces valid output the path returns a safe failure and does not cache the raw response.

## Provider structured-output capability

The production request shape asks capable providers for schema-constrained JSON:

- Groq/OpenAI-compatible Responses uses `text.format` JSON schema with strict mode.
- OpenRouter uses `response_format=json_schema` with strict mode and requires provider parameter support so routing does not silently choose a backend that ignores the schema request.
- Gemini `generateContent` requests JSON response MIME type plus response schema.

Local strict validation remains authoritative even when a provider advertises structured output.

## Continuous offline eval

`v18_6_ai_hardening_test.go` and `ai_continuous_eval_gate.py` provide the mandatory normal-CI lane with **zero live provider calls**. Provider request-shape tests replace `http.DefaultClient` with an isolated in-memory transport.

Required eval lanes:

1. golden structured output;
2. citation validation;
3. contradiction preservation;
4. missing-evidence preservation;
5. injection/adversarial boundary;
6. semantic bounded-context behavior;
7. complete cache identity and TTL invalidation;
8. strict-schema safe abstention;
9. approved/denied/unknown rights fixtures and redacted diagnostics;
10. provider structured-output request capability;
11. bounded aggregate cost/latency/token telemetry.

The runtime telemetry is intentionally aggregate/bounded: counters, cumulative token/cost totals, last/max latency and last outcome only. It does not accumulate prompts, provider source text, URLs or an unbounded event history.

## Protected boundaries

This slice does not alter Day/Swing/Long formulas, Market Mode, Smart Provider Router ownership, desk membership, execution behavior, or the No Execution boundary. AI remains a research-analysis second opinion. It cannot self-promote, modify protected production decision logic, or infer missing provider agreement/rights.

## G10 entry condition

`AUDIT-18-AI-001` and `AUDIT-18-AI-RIGHTS-001` may be presented as current implementation evidence only after fresh Go regression, continuous-eval gate, CI Fast and CI Qualified evidence pass on the same candidate head. Documentation/source markers alone are insufficient.
