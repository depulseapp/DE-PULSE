#!/usr/bin/env python3
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import json
from pathlib import Path
from typing import Any

SCHEMA = "DE.PULSE-CI-TELEMETRY-2"
FAST_NAME = "DE.PULSE | CI Fast"
QUALIFIED_NAME = "DE.PULSE | CI Qualified"
RELEASE_NAME = "DE.PULSE | Release G11-G16"
EVIDENCE_JOB = "Qualified evidence summary"


def parse_time(value: Any) -> datetime | None:
    if not value:
        return None
    text = str(value).replace("Z", "+00:00")
    try:
        return datetime.fromisoformat(text)
    except ValueError:
        return None


def seconds_between(start: Any, end: Any) -> int | None:
    a = parse_time(start)
    b = parse_time(end)
    if not a or not b:
        return None
    return max(0, int((b - a).total_seconds()))


def platform_for(job: dict[str, Any]) -> str:
    labels = " ".join(str(x).lower() for x in job.get("labels", []) or [])
    name = str(job.get("name", "")).lower()
    combined = f"{labels} {name}"
    if "macos" in combined:
        return "macos"
    if "windows" in combined:
        return "windows"
    if "ubuntu" in combined or "linux" in combined:
        return "linux"
    return "unknown"


def normalize_bool(value: str | None) -> bool | None:
    if value is None or value == "":
        return None
    low = value.strip().lower()
    if low == "true":
        return True
    if low == "false":
        return False
    return None


def normalize_int(value: str | None) -> int | None:
    if value is None or value == "":
        return None
    try:
        return max(0, int(value))
    except ValueError:
        return None


def job_rows(payload: dict[str, Any]) -> tuple[list[dict[str, Any]], dict[str, int]]:
    rows: list[dict[str, Any]] = []
    totals = {"linux": 0, "macos": 0, "windows": 0, "unknown": 0}
    for job in payload.get("jobs", []) or []:
        name = str(job.get("name", ""))
        # The evidence job is still running while telemetry is collected, so exclude
        # it from runtime totals rather than recording a misleading partial duration.
        if name == EVIDENCE_JOB:
            continue
        platform = platform_for(job)
        queue_seconds = seconds_between(job.get("created_at"), job.get("started_at"))
        run_seconds = seconds_between(job.get("started_at"), job.get("completed_at"))
        if run_seconds is not None:
            totals[platform] += run_seconds
        rows.append({
            "name": name,
            "status": job.get("status"),
            "conclusion": job.get("conclusion"),
            "platform": platform,
            "queueSeconds": queue_seconds,
            "runSeconds": run_seconds,
        })
    rows.sort(key=lambda row: row["name"])
    return rows, totals


def amplification(runs_payload: dict[str, Any] | None, pr_created_at: str | None) -> dict[str, Any]:
    counts = {"fast": 0, "qualified": 0, "release": 0}
    if runs_payload:
        cutoff = parse_time(pr_created_at)
        for run in runs_payload.get("workflow_runs", []) or []:
            created = parse_time(run.get("created_at"))
            if cutoff and created and created < cutoff:
                continue
            name = run.get("name")
            if name == FAST_NAME:
                counts["fast"] += 1
            elif name == QUALIFIED_NAME:
                counts["qualified"] += 1
            elif name == RELEASE_NAME:
                counts["release"] += 1
    thresholds = {"fastWarningAbove": 3, "qualifiedWarningAbove": 2, "releaseWarningAbove": 1}
    warnings: list[str] = []
    if counts["fast"] > thresholds["fastWarningAbove"]:
        warnings.append(f"Fast run amplification: {counts['fast']} runs on this PR branch")
    if counts["qualified"] > thresholds["qualifiedWarningAbove"]:
        warnings.append(f"Qualified run amplification: {counts['qualified']} runs on this PR branch")
    if counts["release"] > thresholds["releaseWarningAbove"]:
        warnings.append(f"Release amplification: {counts['release']} runs on this PR branch")
    return {"counts": counts, "thresholds": thresholds, "warnings": warnings, "status": "WARN" if warnings else "OK"}


def trustworthy_evidence_accounting(
    jobs: list[dict[str, Any]],
    runner_minutes: dict[str, float],
    browser_setup: dict[str, dict[str, Any]],
) -> dict[str, Any]:
    successful = [row["name"] for row in jobs if row.get("status") == "completed" and row.get("conclusion") == "success"]
    skipped = [row["name"] for row in jobs if row.get("conclusion") == "skipped"]
    failed = [
        row["name"]
        for row in jobs
        if row.get("status") == "completed" and row.get("conclusion") not in {None, "success", "skipped"}
    ]
    total_minutes = round(sum(float(value or 0) for value in runner_minutes.values()), 2)
    evidence_units = len(successful)
    per_unit = round(total_minutes / evidence_units, 3) if evidence_units else None

    skipped_browser_setups = 0
    for name in skipped:
        low = name.lower()
        if "chrome" in low or "webkit" in low or "browser" in low:
            skipped_browser_setups += 1

    observed_setup_seconds = sum(
        int(entry["setupSeconds"])
        for entry in browser_setup.values()
        if entry.get("setupSeconds") is not None
    )
    cache_hits = sum(1 for entry in browser_setup.values() if entry.get("pipCacheHit") is True)

    return {
        "definition": "One trustworthy evidence unit is one completed-success non-summary qualification job recorded by GitHub Actions.",
        "successfulEvidenceUnits": evidence_units,
        "successfulJobs": successful,
        "failedJobs": failed,
        "runnerMinutesConsumed": total_minutes,
        "runnerMinutesPerTrustworthyEvidence": per_unit,
        "avoidedWork": {
            "impactRoutedJobRunsSkipped": len(skipped),
            "skippedJobs": skipped,
            "browserSetupOperationsSkipped": skipped_browser_setups,
            "observedBrowserSetupSeconds": observed_setup_seconds,
            "browserSetupCacheHits": cache_hits,
            "runnerMinutesAvoided": None,
            "runnerMinutesAvoidedReason": (
                "Skipped jobs have no observed counterfactual runtime. DE.PULSE records the avoided job/setup count "
                "and refuses to fabricate avoided minutes without a comparable baseline."
            ),
            "setupSecondsAvoided": None,
            "setupSecondsAvoidedReason": (
                "A cache hit proves reuse but not the cold-install duration that would have occurred. "
                "Observed setup seconds and cache-hit counts are retained; avoided seconds require an external baseline."
            ),
        },
    }


def build(args: argparse.Namespace) -> dict[str, Any]:
    jobs_payload = json.loads(Path(args.jobs_json).read_text(encoding="utf-8"))
    runs_payload = None
    if args.runs_json and Path(args.runs_json).is_file():
        runs_payload = json.loads(Path(args.runs_json).read_text(encoding="utf-8"))
    jobs, totals = job_rows(jobs_payload)
    runner_minutes = {key: round(value / 60.0, 2) for key, value in totals.items()}
    browser_setup = {
        "chrome": {
            "pipCacheHit": normalize_bool(args.chrome_cache_hit),
            "setupSeconds": normalize_int(args.chrome_setup_seconds),
        },
        "webkit": {
            "pipCacheHit": normalize_bool(args.webkit_cache_hit),
            "setupSeconds": normalize_int(args.webkit_setup_seconds),
        },
    }
    evidence = trustworthy_evidence_accounting(jobs, runner_minutes, browser_setup)
    return {
        "schema": SCHEMA,
        "generatedAt": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        "workflow": args.workflow,
        "runId": int(args.run_id),
        "candidateSha": args.candidate_sha,
        "lane": args.lane,
        "webkitRequired": normalize_bool(args.webkit_required),
        "jobs": jobs,
        "runnerSecondsByPlatform": totals,
        "runnerMinutesByPlatform": runner_minutes,
        "browserDependencySetup": browser_setup,
        "workflowAmplification": amplification(runs_payload, args.pr_created_at),
        "trustworthyEvidenceAccounting": evidence,
        "costAccounting": {
            "actualCurrencyCost": None,
            "currencyCostPerTrustworthyEvidence": None,
            "runnerMinutesPerTrustworthyEvidence": evidence["runnerMinutesPerTrustworthyEvidence"],
            "reason": (
                "GitHub billing remains authoritative for currency cost. DE.PULSE reports measured runner consumption "
                "per trustworthy evidence unit and never invents billing rates or counterfactual savings."
            ),
        },
    }


def render_summary(data: dict[str, Any]) -> str:
    minutes = data["runnerMinutesByPlatform"]
    amp = data["workflowAmplification"]
    cache = data["browserDependencySetup"]
    evidence = data["trustworthyEvidenceAccounting"]
    avoided = evidence["avoidedWork"]
    lines = [
        "## DE.PULSE CI Telemetry",
        f"- lane: `{data['lane']}`",
        f"- exact candidate: `{data['candidateSha']}`",
        f"- runner minutes: Linux {minutes['linux']}, macOS {minutes['macos']}, Windows {minutes['windows']}, unknown {minutes['unknown']}",
        f"- trustworthy evidence units: {evidence['successfulEvidenceUnits']}",
        f"- runner minutes / trustworthy evidence: {evidence['runnerMinutesPerTrustworthyEvidence']}",
        f"- impact-routed job runs skipped: {avoided['impactRoutedJobRunsSkipped']}",
        f"- browser setup operations skipped: {avoided['browserSetupOperationsSkipped']}",
        f"- observed browser setup seconds: {avoided['observedBrowserSetupSeconds']}; cache hits: {avoided['browserSetupCacheHits']}",
        f"- PR workflow counts: Fast {amp['counts']['fast']}, Qualified {amp['counts']['qualified']}, Release {amp['counts']['release']} — {amp['status']}",
        f"- Chrome dependency setup: cache={cache['chrome']['pipCacheHit']} seconds={cache['chrome']['setupSeconds']}",
        f"- WebKit dependency setup: cache={cache['webkit']['pipCacheHit']} seconds={cache['webkit']['setupSeconds']}",
        "- avoided runner/setup minutes: intentionally not fabricated without a counterfactual baseline",
        "- currency cost: intentionally not estimated; use GitHub billing for actual currency",
    ]
    for warning in amp["warnings"]:
        lines.append(f"- ⚠ {warning}")
    return "\n".join(lines) + "\n"


def parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description="Build deterministic DE.PULSE CI telemetry from GitHub Actions API payloads")
    p.add_argument("--jobs-json", required=True)
    p.add_argument("--runs-json")
    p.add_argument("--workflow", required=True)
    p.add_argument("--run-id", required=True)
    p.add_argument("--candidate-sha", required=True)
    p.add_argument("--lane", required=True)
    p.add_argument("--webkit-required", default="")
    p.add_argument("--pr-created-at", default="")
    p.add_argument("--chrome-cache-hit", default="")
    p.add_argument("--chrome-setup-seconds", default="")
    p.add_argument("--webkit-cache-hit", default="")
    p.add_argument("--webkit-setup-seconds", default="")
    p.add_argument("--output", required=True)
    p.add_argument("--summary-output")
    return p


def main() -> int:
    args = parser().parse_args()
    data = build(args)
    Path(args.output).write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    summary = render_summary(data)
    print(summary, end="")
    if args.summary_output:
        Path(args.summary_output).write_text(summary, encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
