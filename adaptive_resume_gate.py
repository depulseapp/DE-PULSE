#!/usr/bin/env python3
"""Validate the permanent GitHub-backed DE.PULSE Build Resume Protocol.

This is an owning-check inside existing G0/G2/G10 qualification, not a new gate.
Checkpoint metadata is deliberately fingerprint-excluded so recording progress
cannot mutate the product candidate it describes.
"""
from pathlib import Path
import json, re, sys

ROOT = Path(__file__).resolve().parent
errors = []

def need(ok, msg):
    if not ok:
        errors.append(msg)

required_docs = [
    ROOT / 'adaptive-governance' / 'BUILD_RESUME_PROTOCOL.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_ROADMAP.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_BUILD_PLAN.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_BUILD_PROCESS.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_DELIVERY_PROCESS.md',
]
for p in required_docs:
    need(p.exists(), f'missing permanent resume governance document: {p.relative_to(ROOT)}')

if all(p.exists() for p in required_docs):
    corpus = '\n'.join(p.read_text(errors='ignore') for p in required_docs)
    for term in [
        'last trustworthy PASS',
        'source fingerprint',
        'G0–G16',
        'GitHub',
        'G16',
        'No Execution',
    ]:
        need(term in corpus, f'permanent resume governance missing required contract term: {term}')

identity_path = ROOT / 'release_identity.json'
checkpoint_path = ROOT / '.depulse-certification' / 'resume' / 'build-checkpoint.json'
evidence_path = ROOT / '.depulse-certification' / 'resume' / 'release-evidence-checkpoint.json'
plan_path = ROOT / 'ci_pipeline_plan.json'

try:
    ident = json.loads(identity_path.read_text())
except Exception as exc:
    ident = {}
    errors.append(f'release identity unreadable: {exc}')

try:
    plan = json.loads(plan_path.read_text())
    policy = plan.get('policy', {})
    need(policy.get('build_resume_protocol') == 'adaptive-governance/BUILD_RESUME_PROTOCOL.md', 'CI plan missing canonical build_resume_protocol')
    need(policy.get('resume_checkpoint') == '.depulse-certification/resume/build-checkpoint.json', 'CI plan missing canonical resume_checkpoint')
    need(policy.get('release_evidence_checkpoint') == '.depulse-certification/resume/release-evidence-checkpoint.json', 'CI plan missing release_evidence_checkpoint')
    need(policy.get('resume_reconciliation_required') is True, 'CI plan must require resume reconciliation')
    need(policy.get('checkpoint_metadata_excluded_from_product_fingerprint') is True, 'CI plan must mark checkpoint metadata fingerprint-excluded')
except Exception as exc:
    errors.append(f'CI plan unreadable: {exc}')

try:
    ci_text = (ROOT / 'ci_pipeline.py').read_text(errors='ignore')
    need("'.depulse-certification'" in ci_text, 'ci_pipeline.py must exclude .depulse-certification from product fingerprint')
    pf_text = (ROOT / 'prefreeze_qualification.py').read_text(errors='ignore')
    need("'.depulse-certification'" in pf_text, 'prefreeze fingerprint must exclude .depulse-certification metadata')
except Exception as exc:
    errors.append(f'fingerprint owner validation failed: {exc}')

try:
    cp = json.loads(checkpoint_path.read_text())
    need(cp.get('schemaVersion') == 1, 'build checkpoint schemaVersion must be 1')
    release = str(cp.get('release', '')).lstrip('v')
    need(release == str(ident.get('version', '')).lstrip('v'), 'checkpoint release must match canonical release identity')
    need(cp.get('branch'), 'checkpoint branch missing')
    need(cp.get('sourceCommit'), 'checkpoint sourceCommit missing')
    fp = cp.get('sourceFingerprint')
    fp_state = cp.get('sourceFingerprintState')
    need((isinstance(fp, str) and re.fullmatch(r'[0-9a-f]{64}', fp)) or (fp is None and fp_state in {'PENDING_REQUALIFICATION','NOT_FROZEN'}), 'checkpoint source fingerprint must be verified SHA-256 or explicitly pending')
    gates = cp.get('gates', {})
    need(all(f'G{i}' in gates for i in range(17)), 'checkpoint must contain G0-G16 states')
    need(cp.get('nextStep'), 'checkpoint nextStep missing')
    need('updatedAt' in cp, 'checkpoint updatedAt missing')
except Exception as exc:
    errors.append(f'build checkpoint invalid/unreadable: {exc}')

try:
    ev = json.loads(evidence_path.read_text())
    need(ev.get('schemaVersion') == 1, 'release evidence checkpoint schemaVersion must be 1')
    need(str(ev.get('release', '')).lstrip('v') == str(ident.get('version', '')).lstrip('v'), 'release evidence checkpoint release mismatch')
    need(isinstance(ev.get('evidence'), dict), 'release evidence checkpoint evidence map missing')
except Exception as exc:
    errors.append(f'release evidence checkpoint invalid/unreadable: {exc}')

if errors:
    print('Adaptive Build Resume Contract: FAIL')
    for e in errors:
        print(' -', e)
    sys.exit(1)

print('Adaptive Build Resume Contract: PASS · roadmap/build-plan/build-process/delivery-process integrated · checkpoint + evidence reconciliation enforced · metadata fingerprint-excluded')
