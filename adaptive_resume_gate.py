#!/usr/bin/env python3
"""Validate the permanent GitHub-backed DE.PULSE Build Resume Protocol.

This is an owning-check inside existing G0/G2/G10 qualification, not a new gate.
Checkpoint metadata is deliberately fingerprint-excluded so recording progress
cannot mutate the product candidate it describes.

The resume model intentionally distinguishes two truths during an in-flight
release slice:
1) the immutable last certified Stable checkpoint/release evidence; and
2) the current canonical release candidate identity + engineering branch.
Preparing a new candidate must never rewrite prior Stable PASS evidence merely
so version strings match.
"""
from pathlib import Path
import json, re, subprocess, sys

ROOT = Path(__file__).resolve().parent
errors = []


def need(ok, msg):
    if not ok:
        errors.append(msg)


required_docs = [
    ROOT / 'AGENTS.md',
    ROOT / 'CLAUDE.md',
    ROOT / 'governance' / 'AI-ASSISTANT-PORTABILITY-CONTRACT.md',
    ROOT / 'handoff' / 'CURRENT.md',
    ROOT / 'adaptive-governance' / 'BUILD_RESUME_PROTOCOL.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_ROADMAP.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_BUILD_PLAN.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_BUILD_PROCESS.md',
    ROOT / 'adaptive-governance' / 'ADAPTIVE_DELIVERY_PROCESS.md',
]
for p in required_docs:
    need(p.exists(), f'missing permanent resume governance document: {p.relative_to(ROOT)}')

handoff = ''
if all(p.exists() for p in required_docs):
    corpus = '\n'.join(p.read_text(errors='ignore') for p in required_docs)
    for term in [
        'last trustworthy PASS',
        'source fingerprint',
        'G0–G16',
        'GitHub',
        'G16',
        'No Execution',
        'handoff/CURRENT.md',
        'AI-ASSISTANT-PORTABILITY-CONTRACT.md',
        'Claude',
    ]:
        need(term in corpus, f'permanent resume governance missing required contract term: {term}')

    contract = (ROOT / 'governance' / 'AI-ASSISTANT-PORTABILITY-CONTRACT.md').read_text(errors='ignore')
    agents = (ROOT / 'AGENTS.md').read_text(errors='ignore')
    claude = (ROOT / 'CLAUDE.md').read_text(errors='ignore')
    handoff = (ROOT / 'handoff' / 'CURRENT.md').read_text(errors='ignore')
    for adapter, text in [('AGENTS.md', agents), ('CLAUDE.md', claude)]:
        need('governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md' in text, f'{adapter} must point to the vendor-neutral portability contract')
        need('handoff/CURRENT.md' in text, f'{adapter} must point to the current authoritative handoff')
        need('adaptive_resume_gate.py' in text, f'{adapter} must require the owning resume gate')
    for term in [
        'GitHub source-of-truth hierarchy',
        'Mandatory fresh-session algorithm',
        'Durable handoff rule',
        'Secrets and account independence',
        'No upload of an old chat handoff is required',
    ]:
        need(term in contract, f'portability contract missing required section/term: {term}')
    need('SUPERSEDES ALL PRIOR CHAT HANDOFFS' in handoff, 'handoff/CURRENT.md must be the single current handoff authority')
    need('Exactly one next action' in handoff, 'handoff/CURRENT.md must name exactly one next action')

identity_path = ROOT / 'release_identity.json'
checkpoint_path = ROOT / '.depulse-certification' / 'resume' / 'build-checkpoint.json'
evidence_path = ROOT / '.depulse-certification' / 'resume' / 'release-evidence-checkpoint.json'
plan_path = ROOT / 'ci_pipeline_plan.json'

try:
    ident = json.loads(identity_path.read_text())
except Exception as exc:
    ident = {}
    errors.append(f'release identity unreadable: {exc}')

identity_release = str(ident.get('version', '')).lstrip('v')
identity_build = str(ident.get('build_id', ''))
identity_previous_stable = str(ident.get('previous_stable', '')).lstrip('v')
identity_stable_baseline = str(ident.get('stable_baseline', '')).lstrip('v')

try:
    plan = json.loads(plan_path.read_text())
    policy = plan.get('policy', {})
    need(policy.get('build_resume_protocol') == 'adaptive-governance/BUILD_RESUME_PROTOCOL.md', 'CI plan missing canonical build_resume_protocol')
    need(policy.get('resume_checkpoint') == '.depulse-certification/resume/build-checkpoint.json', 'CI plan missing canonical resume_checkpoint')
    need(policy.get('release_evidence_checkpoint') == '.depulse-certification/resume/release-evidence-checkpoint.json', 'CI plan missing release_evidence_checkpoint')
    need(policy.get('resume_reconciliation_required') is True, 'CI plan must require resume reconciliation')
    need(policy.get('checkpoint_metadata_excluded_from_product_fingerprint') is True, 'CI plan must mark checkpoint metadata fingerprint-excluded')
    need(policy.get('ai_assistant_portability_contract') == 'governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md', 'CI plan missing canonical AI assistant portability contract')
    need(policy.get('authoritative_handoff') == 'handoff/CURRENT.md', 'CI plan missing authoritative current handoff')
    need(policy.get('assistant_account_independent') is True, 'CI plan must require assistant/account independence')
    need(policy.get('ai_entrypoints') == ['AGENTS.md', 'CLAUDE.md'], 'CI plan AI entrypoints drift')
except Exception as exc:
    errors.append(f'CI plan unreadable: {exc}')

try:
    ci_text = (ROOT / 'ci_pipeline.py').read_text(errors='ignore')
    need("'.depulse-certification'" in ci_text, 'ci_pipeline.py must exclude .depulse-certification from product fingerprint')
    pf_text = (ROOT / 'prefreeze_qualification.py').read_text(errors='ignore')
    need("'.depulse-certification'" in pf_text, 'prefreeze fingerprint must exclude .depulse-certification metadata')
except Exception as exc:
    errors.append(f'fingerprint owner validation failed: {exc}')

stable_release = ''
try:
    cp = json.loads(checkpoint_path.read_text())
    schema = cp.get('schemaVersion')
    need(isinstance(schema, int) and schema >= 2, 'build checkpoint schemaVersion must be v2 or later')
    stable_release = str(cp.get('release', '')).lstrip('v')
    need(bool(stable_release), 'build checkpoint Stable release missing')
    need(cp.get('channel') == 'STABLE', 'build checkpoint must describe certified Stable evidence')
    need(cp.get('branch'), 'checkpoint branch missing')

    certified = cp.get('certifiedStable', {}) if isinstance(cp.get('certifiedStable'), dict) else {}
    certified_release = str(certified.get('version', '')).lstrip('v')
    need(certified_release == stable_release, 'checkpoint release must match certifiedStable version')

    if identity_release != stable_release:
        need(identity_previous_stable == stable_release, 'in-flight candidate previous_stable must match immutable Stable checkpoint')
        need(identity_stable_baseline == stable_release, 'in-flight candidate stable_baseline must match immutable Stable checkpoint')
    else:
        need(bool(identity_previous_stable), 'promoted Stable previous_stable must identify its predecessor')
        need(bool(identity_stable_baseline), 'promoted Stable stable_baseline must identify its certified predecessor baseline')
        need(identity_previous_stable == identity_stable_baseline, 'promoted Stable previous_stable/stable_baseline mismatch')

    candidate = cp.get('candidateSourceCommit') or certified.get('candidateSourceCommit') or certified.get('certifiedSourceCheckout')
    need(isinstance(candidate, str) and re.fullmatch(r'[0-9a-f]{40}', candidate), 'checkpoint candidate source commit must be a Git SHA')

    metadata_rule = cp.get('metadataHeadRule') or cp.get('stableIdentityRule')
    if not metadata_rule:
        post_release = cp.get('postReleaseOperationalMetadata', {})
        if isinstance(post_release, dict):
            metadata_rule = post_release.get('rule')
    need(bool(metadata_rule), 'checkpoint metadata/stable identity rule missing')

    fp = cp.get('sourceFingerprint') or certified.get('sourceFingerprint')
    fp_state = cp.get('sourceFingerprintState')
    need((isinstance(fp, str) and re.fullmatch(r'[0-9a-f]{64}', fp)) or (fp is None and fp_state in {'PENDING_REQUALIFICATION','NOT_FROZEN'}), 'checkpoint source fingerprint must be verified SHA-256 or explicitly pending')
    gates = cp.get('gates', {})
    need(all(f'G{i}' in gates for i in range(17)), 'checkpoint must contain G0-G16 states')
    need(cp.get('nextStep'), 'checkpoint nextStep missing')
    need('updatedAt' in cp, 'checkpoint updatedAt missing')
    portability = cp.get('assistantPortability', {})
    need(portability.get('status') == 'ENFORCED', 'checkpoint assistantPortability must be ENFORCED')
    need(portability.get('authoritativeHandoff') == 'handoff/CURRENT.md', 'checkpoint authoritative handoff drift')
    need(portability.get('entrypoints') == ['AGENTS.md', 'CLAUDE.md'], 'checkpoint assistant entrypoints drift')

    stable_handoff_ok = (
        f'**Certified Stable:** `v{stable_release}-stable`' in handoff or
        f'**Release:** `v{stable_release}`' in handoff
    )
    need(stable_handoff_ok, 'current handoff must identify the immutable Stable checkpoint release')

    if identity_release != stable_release:
        expected_engineering_branch = f'v{identity_release}-development'
        need(
            f'**Engineering branch:** `{expected_engineering_branch}`' in handoff,
            'current handoff engineering branch must match in-flight canonical release identity',
        )
        need(
            f'**Candidate package identity:** `{identity_release}` / `{identity_build}`' in handoff,
            'current handoff candidate identity/build must match canonical release identity',
        )
        need(
            'checkpoints intentionally remain anchored' in handoff or 'checkpoints intentionally remain anchored'.lower() in handoff.lower(),
            'current handoff must explain why Stable checkpoints remain immutable during candidate development',
        )
    else:
        old_branch_ok = f"**Active branch:** `{cp.get('branch')}`" in handoff
        need(old_branch_ok or stable_handoff_ok, 'current handoff must reconcile active Stable branch/release')
except Exception as exc:
    errors.append(f'build checkpoint invalid/unreadable: {exc}')

try:
    ev = json.loads(evidence_path.read_text())
    schema = ev.get('schemaVersion')
    need(isinstance(schema, int) and schema >= 2, 'release evidence checkpoint schemaVersion must be v2 or later')
    evidence_release = str(ev.get('release', '')).lstrip('v')
    need(evidence_release == stable_release, 'release evidence checkpoint must match immutable Stable build checkpoint')
    need(ev.get('channel') == 'STABLE', 'release evidence checkpoint must describe certified Stable evidence')
    stable = ev.get('stable', {}) if isinstance(ev.get('stable'), dict) else {}
    need(stable.get('tag') == f'v{stable_release}-stable', 'release evidence Stable tag mismatch')
    need(isinstance(ev.get('evidence'), dict), 'release evidence checkpoint evidence map missing')
    evidence_metadata_rule = ev.get('metadataHeadRule')
    if not evidence_metadata_rule:
        post_release = ev.get('postReleaseOperationalMetadata', {})
        if isinstance(post_release, dict):
            evidence_metadata_rule = post_release.get('rule')
    need(bool(evidence_metadata_rule), 'release evidence checkpoint metadata/stable identity rule missing')
except Exception as exc:
    errors.append(f'release evidence checkpoint invalid/unreadable: {exc}')

try:
    continuity = subprocess.run(
        [sys.executable, str(ROOT / 'tools' / 'ci' / 'post_stable_continuity_gate.py')],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    need(
        continuity.returncode == 0,
        'post-Stable continuity contract failed: ' + (continuity.stdout + continuity.stderr).strip(),
    )
except Exception as exc:
    errors.append(f'post-Stable continuity contract unreadable: {exc}')

if errors:
    print('Adaptive Build Resume Contract: FAIL')
    for e in errors:
        print(' -', e)
    sys.exit(1)

mode = 'IN_FLIGHT_CANDIDATE' if identity_release != stable_release else 'STABLE_ALIGNED'
print(
    'Adaptive Build Resume Contract: PASS · '
    f'mode={mode} · immutable Stable=v{stable_release} · canonical candidate=v{identity_release} · '
    'GitHub-only ChatGPT/Codex/Claude portability enforced · current handoff + Stable checkpoints reconciled · '
    'post-Stable repository continuity enforced · four adaptive layers integrated · metadata fingerprint-excluded'
)
