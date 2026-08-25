#!/usr/bin/env python3
from __future__ import annotations

import json
import os
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
CLOSED_PROCESS_STATES = {"COMPLETE", "COMPLETED", "CLOSED", "DELIVERED"}
ACTIVE_PRODUCT_STATES = {"ACTIVE", "IN_PROGRESS"}


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception:
        return {}


def git(*args: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["git", *args],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )


def branch_name() -> str:
    env = (
        os.environ.get("GITHUB_HEAD_REF")
        or os.environ.get("GITHUB_REF_NAME")
        or ""
    ).strip()
    if env:
        return env.removeprefix("refs/heads/")
    return git("branch", "--show-current").stdout.strip()


def clean_version(value: object) -> str:
    return str(value or "").strip().lstrip("v")


def version_tuple(value: object) -> tuple[int, ...]:
    cleaned = clean_version(value)
    try:
        return tuple(int(part) for part in cleaned.split("."))
    except Exception:
        return ()


def fail(items: list[str]) -> int:
    print("DE.PULSE post-Stable continuity: FAIL", file=sys.stderr)
    for item in items:
        print(f" - {item}", file=sys.stderr)
    return 1


def stable_documentation_errors(stable_tag: str) -> list[str]:
    errors: list[str] = []
    handoff_path = ROOT / "handoff" / "CURRENT.md"
    handoff = (
        handoff_path.read_text(encoding="utf-8", errors="ignore")
        if handoff_path.is_file()
        else ""
    )
    if f"**Certified Stable:** `{stable_tag}`" not in handoff:
        errors.append("handoff/CURRENT.md does not identify the aligned Stable tag")

    current_docs = [
        ROOT / "adaptive-governance" / "CURRENT_ADAPTIVE_ROADMAP.md",
        ROOT / "adaptive-governance" / "CURRENT_ADAPTIVE_BUILD_PLAN.md",
        ROOT / "adaptive-governance" / "CURRENT_ADAPTIVE_BUILD_PROCESS.md",
        ROOT / "adaptive-governance" / "CURRENT_ADAPTIVE_DELIVERY_PROCESS.md",
    ]
    for path in current_docs:
        text = (
            path.read_text(encoding="utf-8", errors="ignore")
            if path.is_file()
            else ""
        )
        if stable_tag not in text:
            errors.append(f"{path.relative_to(ROOT)} is not aligned to {stable_tag}")
    return errors


def process_work_slice_errors(
    *,
    identity: dict,
    identity_version: str,
    checkpoint_version: str,
    branch: str,
) -> list[str]:
    errors: list[str] = []
    state = load_json(ROOT / "governance" / "current-state.json")
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    active = (
        state.get("activeWorkSlice", {})
        if isinstance(state.get("activeWorkSlice"), dict)
        else {}
    )
    capability_gate = (
        state.get("productCapabilityGate", {})
        if isinstance(state.get("productCapabilityGate"), dict)
        else {}
    )

    work_slice_id = str(active.get("workSliceId", "")).strip()
    work_slice_path = (
        ROOT / "governance" / "work-slices" / work_slice_id / "work-slice.json"
        if work_slice_id
        else ROOT / "governance" / "work-slices" / "<missing>" / "work-slice.json"
    )
    work_slice = load_json(work_slice_path)

    stable_tag = str(stable.get("tag", "")).strip()
    stable_candidate = str(stable.get("candidateSha", "")).strip()
    stable_fingerprint = str(stable.get("sourceFingerprint", "")).strip()
    expected_tag = f"v{identity_version}-stable"
    previous_stable = clean_version(identity.get("previous_stable"))

    registered_branch = str(active.get("branch", "")).strip()
    work_registered_branch = str(work_slice.get("branch", "")).strip() if work_slice else ""
    closure_branch = str(active.get("closureBranch", "")).strip()
    work_closure_branch = str(work_slice.get("closureBranch", "")).strip() if work_slice else ""
    active_status = str(active.get("status", "")).strip().upper()
    work_status = str(work_slice.get("status", "")).strip().upper() if work_slice else ""
    is_closure_phase = (
        bool(closure_branch)
        and branch == closure_branch
        and work_closure_branch == closure_branch
        and active_status in CLOSED_PROCESS_STATES
        and work_status in CLOSED_PROCESS_STATES
    )

    if identity.get("channel") != "STABLE":
        errors.append("process work slice requires unchanged STABLE release identity")
    if clean_version(stable.get("productVersion")) != identity_version:
        errors.append("current-state Stable productVersion / release identity drift")
    if stable_tag != expected_tag:
        errors.append(f"current-state Stable tag mismatch: {stable_tag or '<missing>'} != {expected_tag}")
    if str(stable.get("buildId", "")) != str(identity.get("build_id", "")):
        errors.append("current-state Stable buildId / release identity drift")
    if str(stable.get("platformBuildNumber", "")) != str(identity.get("bundle_version", "")):
        errors.append("current-state platform build number / release identity drift")
    if stable.get("publication") != "PASS_NO_REBUILD":
        errors.append("current-state Stable publication must remain PASS_NO_REBUILD")
    if not stable_candidate or not stable_fingerprint:
        errors.append("current-state Stable candidate/fingerprint missing")
    if not str(stable.get("qualifiedSourceSha", "")).strip():
        errors.append("current-state Stable qualified source SHA missing")
    try:
        if int(stable.get("releaseRunId", 0)) <= 0:
            errors.append("current-state Stable release run id missing")
    except Exception:
        errors.append("current-state Stable release run id invalid")

    if not work_slice_id:
        errors.append("registered active process workSliceId missing")
    if not work_slice:
        errors.append(f"registered process work-slice metadata missing: {work_slice_path.relative_to(ROOT)}")
    else:
        if str(work_slice.get("workSliceId", "")) != work_slice_id:
            errors.append("current-state/work-slice id drift")
        if work_slice.get("type") != "PROCESS_RELEASE_ENGINEERING":
            errors.append("registered work slice is not PROCESS_RELEASE_ENGINEERING")
        if work_slice.get("publicProductVersion") is not None:
            errors.append("process work slice must not consume a public product version")
        if work_slice.get("productBehaviorChange") is not False:
            errors.append("process work slice must explicitly declare productBehaviorChange=false")
        if clean_version(work_slice.get("stableProductVersionAtStart")) != identity_version:
            errors.append("process work-slice Stable-at-start version / release identity drift")
        if str(work_slice.get("baselineCandidateSha", "")) != stable_candidate:
            errors.append("process work-slice baseline candidate / current Stable candidate drift")
        if str(work_slice.get("baselineSourceFingerprint", "")) != stable_fingerprint:
            errors.append("process work-slice baseline fingerprint / current Stable fingerprint drift")
        if str(work_slice.get("baselineBuildId", "")) != str(identity.get("build_id", "")):
            errors.append("process work-slice baseline build / release identity drift")
        if work_registered_branch != registered_branch:
            errors.append("current-state/work-slice registered implementation branch drift")

        if is_closure_phase:
            if work_slice.get("blocksNextProductCapability") is not False:
                errors.append("completed process work slice must unblock next product capability")
            final_evidence = str(work_slice.get("finalQualificationEvidence", "")).strip()
            if not final_evidence or not (ROOT / final_evidence).is_file():
                errors.append("completed process closure requires retained finalQualificationEvidence")
            if str(work_slice.get("mergedCommitSha", "")).strip() != str(active.get("mergedCommitSha", "")).strip():
                errors.append("completed process merged-commit binding drift")
        else:
            if branch != work_registered_branch:
                errors.append("process work-slice registered branch / current branch drift")
            if work_slice.get("blocksNextProductCapability") is not True:
                errors.append("in-flight process work slice must keep next product capability blocked")
        if active.get("issue") != work_slice.get("issue"):
            errors.append("current-state/work-slice issue drift")

    if is_closure_phase:
        if capability_gate.get("blocked") is not False:
            errors.append("completed process closure must unblock product capability gate")
        if capability_gate.get("blockedByIssue") is not None:
            errors.append("completed process closure must clear product capability blocker")
    else:
        if registered_branch != branch:
            errors.append("current-state active process branch / current branch drift")
        if capability_gate.get("blocked") is not True:
            errors.append("next product capability must remain blocked during in-flight process work slice")
        if capability_gate.get("blockedByIssue") != active.get("issue"):
            errors.append("product capability blocker must be the active process issue")

    if active.get("type") != "PROCESS_RELEASE_ENGINEERING":
        errors.append("current-state active work slice is not PROCESS_RELEASE_ENGINEERING")
    if active.get("publicProductVersion") is not None:
        errors.append("current-state process work slice must not have a public product version")
    if active.get("productBehaviorChange") is not False:
        errors.append("current-state process work slice must declare productBehaviorChange=false")

    if checkpoint_version != identity_version and checkpoint_version != previous_stable:
        errors.append(
            "process work-slice checkpoint may only equal the current Stable or its immediate predecessor: "
            f"checkpoint=v{checkpoint_version}, current=v{identity_version}, predecessor=v{previous_stable or '<missing>'}"
        )

    if stable_tag and stable_candidate:
        tag = git("rev-parse", "--verify", f"refs/tags/{stable_tag}^{{commit}}")
        if tag.returncode != 0:
            errors.append(f"registered Stable tag is not available in Git history: {stable_tag}")
        elif tag.stdout.strip() != stable_candidate:
            errors.append(
                f"registered Stable tag/candidate mismatch: {stable_tag} -> {tag.stdout.strip()} != {stable_candidate}"
            )

    if stable_candidate:
        head = git("rev-parse", "HEAD")
        if head.returncode != 0:
            errors.append("cannot resolve current HEAD for process-work ancestry check")
        else:
            ancestry = git("merge-base", "--is-ancestor", stable_candidate, head.stdout.strip())
            if ancestry.returncode != 0:
                errors.append("process work branch does not descend from the recorded Stable candidate")

    return errors


def product_work_slice_errors(
    *,
    identity: dict,
    identity_version: str,
    checkpoint_version: str,
    branch: str,
) -> list[str]:
    errors: list[str] = []
    state = load_json(ROOT / "governance" / "current-state.json")
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    active = (
        state.get("activeWorkSlice", {})
        if isinstance(state.get("activeWorkSlice"), dict)
        else {}
    )
    capability_gate = (
        state.get("productCapabilityGate", {})
        if isinstance(state.get("productCapabilityGate"), dict)
        else {}
    )

    reserved_id = str(capability_gate.get("reservedWorkSliceId", "")).strip()
    reserved_branch = str(capability_gate.get("reservedBranch", "")).strip()
    reserved_status = str(capability_gate.get("reservationStatus", "")).strip().upper()
    reserved_issue = capability_gate.get("reservedIssue")
    work_slice_rel = str(capability_gate.get("workSlicePath", "")).strip()
    expected_rel = (
        f"governance/work-slices/{reserved_id}/work-slice.json"
        if reserved_id
        else ""
    )
    work_slice_path = ROOT / work_slice_rel if work_slice_rel else ROOT / "<missing>"
    work_slice = load_json(work_slice_path) if work_slice_rel else {}

    stable_tag = str(stable.get("tag", "")).strip()
    stable_candidate = str(stable.get("candidateSha", "")).strip()
    stable_fingerprint = str(stable.get("sourceFingerprint", "")).strip()
    previous_stable = clean_version(identity.get("previous_stable"))
    work_type = str(work_slice.get("type", "")).strip()
    is_release_closure = work_type == "PRODUCT_RELEASE_CLOSURE"
    published_stable_version = previous_stable if is_release_closure else identity_version
    expected_tag = f"v{published_stable_version}-stable"
    completed_process_id = str(active.get("workSliceId", "")).strip()
    completed_process_path = (
        ROOT / "governance" / "work-slices" / completed_process_id / "work-slice.json"
        if completed_process_id
        else ROOT / "<missing>"
    )
    completed_process = load_json(completed_process_path) if completed_process_id else {}

    if identity.get("channel") != "STABLE":
        errors.append("product work slice requires STABLE candidate release identity")
    if is_release_closure:
        if not previous_stable:
            errors.append("release closure candidate must declare previous_stable")
        if clean_version(identity.get("stable_baseline")) != previous_stable:
            errors.append("release closure stable_baseline / previous_stable drift")
        if not version_tuple(identity_version) or not version_tuple(previous_stable) or version_tuple(identity_version) <= version_tuple(previous_stable):
            errors.append("release closure candidate version must be strictly newer than published Stable")
        if clean_version(stable.get("productVersion")) != previous_stable:
            errors.append("release closure must preserve previous Stable productVersion until publication")
        if stable_tag != expected_tag:
            errors.append(f"release closure must preserve published Stable tag: {stable_tag or '<missing>'} != {expected_tag}")
        if str(stable.get("buildId", "")) != str(work_slice.get("baselineBuildId", "")):
            errors.append("release closure published Stable buildId / baseline build drift")
        try:
            stable_build = int(stable.get("platformBuildNumber", 0))
            candidate_build = int(identity.get("bundle_version", 0))
            if stable_build <= 0 or candidate_build <= stable_build:
                errors.append("release closure platform build must advance monotonically beyond published Stable")
        except Exception:
            errors.append("release closure platform build numbers are invalid")
    else:
        if clean_version(stable.get("productVersion")) != identity_version:
            errors.append("current-state Stable productVersion / release identity drift")
        if stable_tag != expected_tag:
            errors.append(f"current-state Stable tag mismatch: {stable_tag or '<missing>'} != {expected_tag}")
        if str(stable.get("buildId", "")) != str(identity.get("build_id", "")):
            errors.append("current-state Stable buildId / release identity drift")
        if str(stable.get("platformBuildNumber", "")) != str(identity.get("bundle_version", "")):
            errors.append("current-state platform build number / release identity drift")
    if stable.get("publication") != "PASS_NO_REBUILD":
        errors.append("current-state Stable publication must remain PASS_NO_REBUILD")
    if not stable_candidate or not stable_fingerprint:
        errors.append("current-state Stable candidate/fingerprint missing")
    if not str(stable.get("qualifiedSourceSha", "")).strip():
        errors.append("current-state Stable qualified source SHA missing")
    try:
        if int(stable.get("releaseRunId", 0)) <= 0:
            errors.append("current-state Stable release run id missing")
    except Exception:
        errors.append("current-state Stable release run id invalid")

    if capability_gate.get("blocked") is not False:
        errors.append("registered product work requires an unblocked product capability gate")
    if capability_gate.get("blockedByIssue") is not None:
        errors.append("registered product work must not retain a process capability blocker")
    if not str(capability_gate.get("nextReservedCapability", "")).strip():
        errors.append("registered product capability name missing")
    if not reserved_id:
        errors.append("registered product workSliceId missing")
    if not isinstance(reserved_issue, int) or reserved_issue <= 0:
        errors.append("registered product issue missing or invalid")
    if not reserved_branch:
        errors.append("registered product branch missing")
    elif reserved_branch != branch:
        errors.append("registered product branch / current branch drift")
    if reserved_status not in ACTIVE_PRODUCT_STATES:
        errors.append("registered product reservationStatus must be ACTIVE or IN_PROGRESS")
    if work_slice_rel != expected_rel:
        errors.append(
            f"registered product workSlicePath mismatch: {work_slice_rel or '<missing>'} != {expected_rel or '<missing>'}"
        )
    if not work_slice:
        errors.append(
            f"registered product work-slice metadata missing: {work_slice_rel or '<missing>'}"
        )
    else:
        if work_slice.get("schema") != "DE.PULSE-WORK-SLICE-1":
            errors.append("registered product work-slice schema mismatch")
        if str(work_slice.get("workSliceId", "")).strip() != reserved_id:
            errors.append("product capability/work-slice id drift")
        if work_slice.get("issue") != reserved_issue:
            errors.append("product capability/work-slice issue drift")
        if work_type not in {"PRODUCT_CAPABILITY", "PRODUCT_RELEASE_CLOSURE"}:
            errors.append("registered product work slice must be PRODUCT_CAPABILITY or PRODUCT_RELEASE_CLOSURE")
        if str(work_slice.get("status", "")).strip().upper() not in ACTIVE_PRODUCT_STATES:
            errors.append("registered product work-slice status is not active")
        if str(work_slice.get("branch", "")).strip() != reserved_branch:
            errors.append("product capability/work-slice branch drift")
        if is_release_closure:
            if clean_version(work_slice.get("publicProductVersion")) != identity_version:
                errors.append("release closure publicProductVersion / candidate identity drift")
            if clean_version(work_slice.get("stableProductVersionAtStart")) != previous_stable:
                errors.append("release closure Stable-at-start must equal previous Stable")
            if str(work_slice.get("targetStable", "")).strip() != f"v{identity_version}-stable":
                errors.append("release closure targetStable / candidate identity drift")
            if not str(work_slice.get("releaseClaim", "")).strip():
                errors.append("release closure releaseClaim missing")
            closure_rel = str(work_slice.get("closureLedger", "")).strip()
            if not closure_rel or not (ROOT / closure_rel).is_file():
                errors.append("release closure requires a retained closure ledger")
        elif work_slice.get("publicProductVersion") is not None:
            errors.append("product capability identity must remain separate from public SemVer")
        if work_slice.get("productBehaviorChange") is not True:
            errors.append("product work slice must declare productBehaviorChange=true")
        if work_slice.get("blocksNextProductCapability") is not True:
            errors.append("active product work must block starting the next product capability")
        if not is_release_closure and clean_version(work_slice.get("stableProductVersionAtStart")) != identity_version:
            errors.append("product work-slice Stable-at-start version / release identity drift")
        if str(work_slice.get("baselineCandidateSha", "")).strip() != stable_candidate:
            errors.append("product work-slice baseline candidate / current Stable candidate drift")
        if str(work_slice.get("baselineSourceFingerprint", "")).strip() != stable_fingerprint:
            errors.append("product work-slice baseline fingerprint / current Stable fingerprint drift")
        expected_baseline_build = str(stable.get("buildId", "")) if is_release_closure else str(identity.get("build_id", ""))
        if str(work_slice.get("baselineBuildId", "")).strip() != expected_baseline_build:
            errors.append("product work-slice baseline build / governing Stable identity drift")

    if not completed_process:
        errors.append("completed process work-slice metadata missing while product work is active")
    else:
        if active.get("type") != "PROCESS_RELEASE_ENGINEERING":
            errors.append("product work requires retained completed process-work authority")
        if str(active.get("status", "")).strip().upper() not in CLOSED_PROCESS_STATES:
            errors.append("product work requires the retained process work slice to be complete")
        if str(completed_process.get("status", "")).strip().upper() not in CLOSED_PROCESS_STATES:
            errors.append("retained process work-slice metadata is not complete")
        if completed_process.get("blocksNextProductCapability") is not False:
            errors.append("retained completed process work slice has not unblocked product work")
        if capability_gate.get("unblockedByCompletedWorkSlice") != completed_process_id:
            errors.append("product capability gate / completed process unblock binding drift")
        final_evidence = str(completed_process.get("finalQualificationEvidence", "")).strip()
        if not final_evidence or not (ROOT / final_evidence).is_file():
            errors.append("retained completed process work slice lacks final qualification evidence")

    if is_release_closure:
        if checkpoint_version != previous_stable:
            errors.append(
                "release closure checkpoint must remain the immediate published Stable until publication: "
                f"checkpoint=v{checkpoint_version}, previous=v{previous_stable or '<missing>'}"
            )
    elif checkpoint_version != identity_version and checkpoint_version != previous_stable:
        errors.append(
            "product work-slice checkpoint may only equal the current Stable or its immediate predecessor: "
            f"checkpoint=v{checkpoint_version}, current=v{identity_version}, predecessor=v{previous_stable or '<missing>'}"
        )

    if stable_tag and stable_candidate:
        tag = git("rev-parse", "--verify", f"refs/tags/{stable_tag}^{{commit}}")
        if tag.returncode != 0:
            errors.append(f"registered Stable tag is not available in Git history: {stable_tag}")
        elif tag.stdout.strip() != stable_candidate:
            errors.append(
                f"registered Stable tag/candidate mismatch: {stable_tag} -> {tag.stdout.strip()} != {stable_candidate}"
            )

    if stable_candidate:
        head = git("rev-parse", "HEAD")
        if head.returncode != 0:
            errors.append("cannot resolve current HEAD for product-work ancestry check")
        else:
            ancestry = git("merge-base", "--is-ancestor", stable_candidate, head.stdout.strip())
            if ancestry.returncode != 0:
                errors.append("product work branch does not descend from the recorded Stable candidate")

    return errors


def main() -> int:
    errors: list[str] = []
    identity = load_json(ROOT / "release_identity.json")
    build = load_json(ROOT / ".depulse-certification" / "resume" / "build-checkpoint.json")
    evidence = load_json(
        ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json"
    )

    identity_version = clean_version(identity.get("version"))
    checkpoint_version = clean_version(build.get("release"))
    evidence_version = clean_version(evidence.get("release"))
    branch = branch_name()

    if not identity_version or not checkpoint_version or not evidence_version:
        errors.append("release identity/build checkpoint/release evidence version missing")
        return fail(errors)

    if checkpoint_version != evidence_version:
        errors.append(
            f"build/evidence Stable mismatch: {checkpoint_version} != {evidence_version}"
        )
        return fail(errors)

    mode = ""
    stable_tag = f"v{identity_version}-stable"

    if identity_version == checkpoint_version:
        mode = "STABLE_ALIGNED"
        manifest = ROOT / "release" / f"v{identity_version}" / "stable-evidence-manifest.json"
        if not manifest.is_file():
            errors.append(f"Stable evidence manifest missing: {manifest.relative_to(ROOT)}")
        errors.extend(stable_documentation_errors(stable_tag))
    else:
        expected_candidate_branch = f"v{identity_version}-development"
        state = load_json(ROOT / "governance" / "current-state.json")
        active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
        capability_gate = (
            state.get("productCapabilityGate", {})
            if isinstance(state.get("productCapabilityGate"), dict)
            else {}
        )
        registered_process_branch = str(active.get("branch", "")).strip()
        registered_closure_branch = str(active.get("closureBranch", "")).strip()
        registered_product_branch = str(capability_gate.get("reservedBranch", "")).strip()
        work_slice_rel = str(capability_gate.get("workSlicePath", "")).strip()
        registered_product_work = load_json(ROOT / work_slice_rel) if work_slice_rel else {}
        registered_product_type = str(registered_product_work.get("type", "")).strip()
        is_registered_process = (
            branch in {registered_process_branch, registered_closure_branch}
            and bool(branch)
            and active.get("type") == "PROCESS_RELEASE_ENGINEERING"
        )
        is_registered_product = (
            bool(branch)
            and branch == registered_product_branch
            and capability_gate.get("blocked") is False
        )

        if branch in {"main", "master"}:
            errors.append(
                f"default branch carries STABLE identity v{identity_version} while durable Stable checkpoint is v{checkpoint_version}; post-Stable continuity reconciliation is required"
            )
        elif branch == expected_candidate_branch:
            mode = "IN_FLIGHT_CANDIDATE"
        elif is_registered_process:
            mode = "PROCESS_WORK_SLICE_CLOSURE" if branch == registered_closure_branch else "PROCESS_WORK_SLICE"
            errors.extend(
                process_work_slice_errors(
                    identity=identity,
                    identity_version=identity_version,
                    checkpoint_version=checkpoint_version,
                    branch=branch,
                )
            )
        elif is_registered_product:
            mode = "PRODUCT_RELEASE_CLOSURE" if registered_product_type == "PRODUCT_RELEASE_CLOSURE" else "PRODUCT_WORK_SLICE"
            errors.extend(
                product_work_slice_errors(
                    identity=identity,
                    identity_version=identity_version,
                    checkpoint_version=checkpoint_version,
                    branch=branch,
                )
            )
        else:
            errors.append(
                "identity/checkpoint differ outside an allowed product candidate, registered product work slice, or registered process work slice: "
                f"candidate={expected_candidate_branch}, product={registered_product_branch or '<none>'}, "
                f"process={registered_process_branch or '<none>'}, closure={registered_closure_branch or '<none>'}, "
                f"current={branch or '<detached>'}"
            )

    if errors:
        return fail(errors)

    print(
        "DE.PULSE post-Stable continuity: PASS · "
        f"mode={mode} · branch={branch or '<detached>'} · identity=v{identity_version} · durable checkpoint=v{checkpoint_version}"
    )
    if mode.startswith("PROCESS_WORK_SLICE"):
        print("registered process/closure branch / Stable tag / ancestry / no-product-version invariants: PASS")
        print("prior-Stable immutable resume checkpoint exception: BOUNDED_TO_IMMEDIATE_PREDECESSOR")
    if mode == "PRODUCT_WORK_SLICE":
        print("registered product work-slice / Stable tag / ancestry / SemVer-separation invariants: PASS")
        print("prior-Stable immutable resume checkpoint exception: BOUNDED_TO_IMMEDIATE_PREDECESSOR")
    if mode == "PRODUCT_RELEASE_CLOSURE":
        print("active release-closure candidate / published-Stable separation: PASS")
        print("previous Stable tag/candidate/fingerprint/baseline/ancestry invariants: PASS")
        print("candidate SemVer + platform build monotonicity: PASS")
        print("publication remains deferred until canonical release gates: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
