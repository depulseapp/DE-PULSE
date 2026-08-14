#!/usr/bin/env python3
import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent
REG = ROOT / "canonical_naming_registry.json"
FUNC = ROOT / "functionality_utility_registry.json"
ROADMAP = ROOT / "adaptive-governance" / "ADAPTIVE_ROADMAP.md"
PROCESS = ROOT / "adaptive-governance" / "ADAPTIVE_BUILD_PROCESS.md"
CONTRACT = ROOT / "adaptive-governance" / "NAMING_AND_IDENTITY_CONTRACT.md"

errors = []

if not REG.exists():
    errors.append("canonical_naming_registry.json is missing")
else:
    reg = json.loads(REG.read_text())
    if reg.get("product", {}).get("displayName") != "DE.PULSE":
        errors.append("canonical product display name must be DE.PULSE")

    expected_gates = [
        ("G0", "Exact Baseline"),
        ("G1", "Immutable Scope"),
        ("G2", "Architecture & Data Utility"),
        ("G3", "Design & Dependency Readiness"),
        ("G4", "Development Exit"),
        ("G5", "FAST Qualification"),
        ("G6", "Integration & MEDIUM Qualification"),
        ("G7", "Data, Security & Adaptive Intelligence"),
        ("G8", "Performance, Capacity & Stability"),
        ("G9", "Cross-Module UI/UX"),
        ("G10", "Pre-Freeze Qualification"),
        ("G11", "Immutable Release Candidate"),
        ("G12", "Full Certification"),
        ("G13", "Native Packaging & Provenance"),
        ("G14", "Actual Artifact Runtime Audit"),
        ("G15", "Release Assurance & Promotion"),
        ("G16", "Adaptive Retrospective & Handoff"),
    ]
    actual_gates = [(x.get("id"), x.get("name")) for x in reg.get("gateMap", [])]
    if actual_gates != expected_gates:
        errors.append("canonical gate map/name drift detected")

    required_roles = {"SUPER_OWNER", "OWNER", "ADMIN", "USER", "DEMO"}
    if set(reg.get("roles", [])) != required_roles:
        errors.append("canonical role vocabulary drift detected")

    required_states = {"PASS", "FAIL", "PENDING", "BLOCKED", "INVALIDATED", "SUPERSEDED"}
    if set(reg.get("evidenceStates", [])) != required_states:
        errors.append("canonical evidence-state vocabulary drift detected")

    required_failures = {"PRODUCT_FAIL", "GATE_TEST_FAIL", "CI_HARNESS_FAIL", "INFRA_FAIL", "EXPECTED_NOOP", "SUPERSEDED"}
    if set(reg.get("failureClasses", [])) != required_failures:
        errors.append("canonical failure-class vocabulary drift detected")

    required_delivery = {"NOT_READY", "READY", "DELIVERED"}
    if set(reg.get("userDeliveryStates", [])) != required_delivery:
        errors.append("canonical user-delivery vocabulary drift detected")

    surface_ids = {x.get("id") for x in reg.get("surfaces", [])}
    if FUNC.exists():
        func = json.loads(FUNC.read_text())
        registry_surface_ids = {
            x.get("surfaceId") for x in func.get("items", [])
            if x.get("surfaceId")
        }
        missing = registry_surface_ids - surface_ids
        if missing:
            errors.append("functionality surfaces missing from naming registry: " + ", ".join(sorted(missing)))

for p in [ROADMAP, PROCESS, CONTRACT]:
    if not p.exists():
        errors.append(f"missing governing naming source: {p.relative_to(ROOT)}")

if CONTRACT.exists():
    text = CONTRACT.read_text()
    for gate, name in [
        ("G0", "Exact Baseline"), ("G1", "Immutable Scope"),
        ("G2", "Architecture & Data Utility"), ("G3", "Design & Dependency Readiness"),
        ("G4", "Development Exit"), ("G5", "FAST Qualification"),
        ("G6", "Integration & MEDIUM Qualification"),
        ("G7", "Data, Security & Adaptive Intelligence"),
        ("G8", "Performance, Capacity & Stability"),
        ("G9", "Cross-Module UI/UX"), ("G10", "Pre-Freeze Qualification"),
        ("G11", "Immutable Release Candidate"), ("G12", "Full Certification"),
        ("G13", "Native Packaging & Provenance"),
        ("G14", "Actual Artifact Runtime Audit"),
        ("G15", "Release Assurance & Promotion"),
        ("G16", "Adaptive Retrospective & Handoff")
    ]:
        if f"**{gate} — {name}**" not in text:
            errors.append(f"naming contract missing canonical {gate} title")

# Guard machine registries from durable placeholder naming. Documentation may
# intentionally mention these strings as negative examples, so do not scan prose.
if REG.exists():
    txt = REG.read_text().lower()
    for bad in ["final-final", "phase2", "new-check"]:
        if bad in txt:
            errors.append(f"ambiguous durable naming token present in {REG.name}: {bad}")

if errors:
    print("Canonical Naming Gate: FAIL")
    for e in errors:
        print(f"- {e}")
    sys.exit(1)

print("Canonical Naming Gate: PASS")
print("Product, gates, roles, states, surfaces and naming contract are aligned.")
