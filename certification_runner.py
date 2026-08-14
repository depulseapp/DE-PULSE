#!/usr/bin/env python3
"""Checkpointed DE.PULSE G0-G16 certification orchestrator.

Design goals:
- each expensive certification item runs as an independent process;
- PASS checkpoints survive interruption and are reused only for the exact same source fingerprint;
- runner/infrastructure failures are distinguished from product/test failures;
- blocked release-host prerequisites do not masquerade as product defects;
- all logs/results are durable and resumable outside the source tree.
"""
from __future__ import annotations
import argparse, datetime as dt, fcntl, hashlib, json, os, shutil, signal, subprocess, sys, time
from pathlib import Path
from source_fingerprint import canonical_source_fingerprint

ROOT = Path(__file__).resolve().parent
DEFAULT_PLAN = ROOT / "certification_plan.json"
STATE_DIR = Path(os.environ.get("DEPULSE_CERT_DIR", str(ROOT.parent / f".depulse-certification-{ROOT.name}"))).resolve()
STATE_FILE = STATE_DIR / "state.json"
SUMMARY_FILE = STATE_DIR / "summary.json"
LOG_DIR = STATE_DIR / "logs"
RESULT_DIR = STATE_DIR / "results"
LOCK_FILE = STATE_DIR / "runner.lock"

EXCLUDED_DIRS = {".depulse-certification", "__pycache__", ".git"}
EXCLUDED_SUFFIXES = {".log", ".out", ".tmp", ".test", ".exe"}
INFRA_PATTERNS = (
    "resource temporarily unavailable",
    "no space left on device",
    "broken pipe",
    "epipe",
    "browser executable doesn't exist",
    "executable doesn't exist",
    "cannot find module 'playwright'",
    "killed",
)


def utc_now() -> str:
    return dt.datetime.now(dt.timezone.utc).isoformat()




def acquire_runner_lock():
    """Prevent overlapping certification sessions on the same source tree."""
    STATE_DIR.mkdir(exist_ok=True)
    fh = LOCK_FILE.open("a+")
    try:
        fcntl.flock(fh.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
    except BlockingIOError:
        fh.seek(0)
        owner = fh.read().strip() or "unknown owner"
        print(f"Certification runner: BLOCKED - another session owns this candidate ({owner})")
        raise SystemExit(3)
    fh.seek(0); fh.truncate()
    fh.write(f"pid={os.getpid()} started_at={utc_now()}\n"); fh.flush()
    return fh


def source_fingerprint() -> str:
    def exclude_build_output(p: Path, rel) -> bool:
        return p.name.startswith("De-Pulse-") and p.suffix.lower() in {".zip", ".app"}
    return canonical_source_fingerprint(
        ROOT,
        excluded_dirs=EXCLUDED_DIRS,
        excluded_suffixes=EXCLUDED_SUFFIXES,
        exclude_file=exclude_build_output,
    )


def load_plan(path: Path) -> dict:
    plan = json.loads(path.read_text())
    ids = [c["id"] for c in plan.get("checks", [])]
    if len(ids) != len(set(ids)):
        raise SystemExit("certification plan contains duplicate check IDs")
    return plan


def load_state(fp: str) -> dict:
    STATE_DIR.mkdir(exist_ok=True); LOG_DIR.mkdir(exist_ok=True); RESULT_DIR.mkdir(exist_ok=True)
    if STATE_FILE.exists():
        try:
            state = json.loads(STATE_FILE.read_text())
        except Exception:
            state = {}
    else:
        state = {}
    if state.get("source_fingerprint") != fp:
        if state:
            archive = STATE_DIR / f"state-invalidated-{int(time.time())}.json"
            archive.write_text(json.dumps(state, indent=2) + "\n")
        state = {
            "schema": "DE.PULSE-CERTIFICATION-STATE-1",
            "source_fingerprint": fp,
            "created_at": utc_now(),
            "updated_at": utc_now(),
            "results": {},
        }
        save_state(state)
    return state


def save_state(state: dict) -> None:
    state["updated_at"] = utc_now()
    tmp = STATE_FILE.with_suffix(".tmp")
    tmp.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
    tmp.replace(STATE_FILE)


def classify(returncode: int | None, timed_out: bool, output_tail: str, missing: bool = False) -> str:
    if missing:
        return "BLOCKED"
    if timed_out:
        return "INFRA FAIL"
    if returncode == 0:
        return "PASS"
    if returncode == 3:
        return "BLOCKED"
    if returncode is not None and (returncode < 0 or returncode in {137, 143}):
        return "INFRA FAIL"
    lower = output_tail.lower()
    if any(p in lower for p in INFRA_PATTERNS):
        return "INFRA FAIL"
    return "PRODUCT FAIL"


def render_command(raw: list[str]) -> list[str]:
    mapping = {
        "{python}": sys.executable,
        "{root}": str(ROOT),
        "{cert_dir}": str(STATE_DIR),
    }
    return [mapping.get(x, x.replace("{root}", str(ROOT)).replace("{cert_dir}", str(STATE_DIR))) for x in raw]


def run_check(check: dict, state: dict, force: bool = False) -> dict:
    cid = check["id"]
    prior = state["results"].get(cid)
    if prior and prior.get("status") == "PASS" and not force:
        print(f"[RESUME] {cid}: PASS · {prior.get('duration_seconds', 0):.1f}s")
        return prior

    cmd = render_command(check["command"])
    timeout = int(check.get("timeout_seconds", 120))
    log_path = LOG_DIR / f"{cid}.log"
    started = time.time()
    started_at = utc_now()
    timed_out = False
    missing = False
    rc: int | None = None
    print(f"[RUN] {check.get('gate','?')} {cid} — {check.get('label', cid)}", flush=True)
    print("      " + " ".join(cmd), flush=True)

    with log_path.open("w", encoding="utf-8", errors="replace") as log:
        log.write(f"DE.PULSE certification check: {cid}\n")
        log.write(f"gate={check.get('gate')}\nstarted_at={started_at}\nsource_fingerprint={state['source_fingerprint']}\n")
        log.write("command=" + " ".join(cmd) + "\n\n")
        log.flush()
        try:
            p = subprocess.Popen(cmd, cwd=ROOT, stdout=log, stderr=subprocess.STDOUT, text=True, start_new_session=True)
            try:
                rc = p.wait(timeout=timeout)
            except subprocess.TimeoutExpired:
                timed_out = True
                try:
                    os.killpg(p.pid, signal.SIGTERM)
                    p.wait(timeout=5)
                except Exception:
                    try: os.killpg(p.pid, signal.SIGKILL)
                    except Exception: pass
                rc = p.returncode
                log.write(f"\nCERTIFICATION_RUNNER_TIMEOUT after {timeout}s\n")
        except FileNotFoundError as exc:
            missing = True
            log.write(f"CERTIFICATION_RUNNER_MISSING_EXECUTABLE: {exc}\n")
        except Exception as exc:
            timed_out = True
            log.write(f"CERTIFICATION_RUNNER_INFRA_EXCEPTION: {type(exc).__name__}: {exc}\n")

    duration = time.time() - started
    tail = log_path.read_text(errors="replace")[-12000:]
    status = classify(rc, timed_out, tail, missing)
    result = {
        "id": cid,
        "gate": check.get("gate"),
        "label": check.get("label", cid),
        "status": status,
        "blocking": bool(check.get("blocking", True)),
        "started_at": started_at,
        "finished_at": utc_now(),
        "duration_seconds": round(duration, 3),
        "returncode": rc,
        "timeout_seconds": timeout,
        "log": str(log_path),
        "source_fingerprint": state["source_fingerprint"],
    }
    state["results"][cid] = result
    save_state(state)
    (RESULT_DIR / f"{cid}.json").write_text(json.dumps(result, indent=2, sort_keys=True) + "\n")
    print(f"[{status}] {cid} · {duration:.1f}s", flush=True)
    if status != "PASS":
        print("      log: " + str(log_path), flush=True)
    return result


def summarize(plan: dict, state: dict, selected_ids: set[str] | None = None) -> dict:
    rows = []
    for c in plan.get("checks", []):
        if selected_ids is not None and c["id"] not in selected_ids:
            continue
        r = state["results"].get(c["id"])
        rows.append({
            "id": c["id"], "gate": c.get("gate"), "label": c.get("label", c["id"]),
            "blocking": bool(c.get("blocking", True)),
            "status": r.get("status") if r else "NOT RUN",
            "duration_seconds": r.get("duration_seconds") if r else None,
            "log": r.get("log") if r else None,
        })
    statuses = [r["status"] for r in rows if r["blocking"]]
    if any(x == "PRODUCT FAIL" for x in statuses): overall = "PRODUCT FAIL"
    elif any(x in {"BLOCKED", "INFRA FAIL", "NOT RUN"} for x in statuses): overall = "BLOCKED"
    else: overall = "PASS"
    summary = {
        "schema": "DE.PULSE-CERTIFICATION-SUMMARY-1",
        "generated_at": utc_now(),
        "source_fingerprint": state["source_fingerprint"],
        "overall": overall,
        "checks": rows,
    }
    SUMMARY_FILE.write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
    return summary


def print_summary(summary: dict) -> None:
    print("\nDE.PULSE Certification Summary")
    print("Source fingerprint:", summary["source_fingerprint"])
    for r in summary["checks"]:
        dur = "" if r["duration_seconds"] is None else f" · {r['duration_seconds']:.1f}s"
        print(f"{r['gate']:>3}  {r['status']:<12}  {r['id']}{dur}")
    print("OVERALL:", summary["overall"])


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--plan", default=str(DEFAULT_PLAN))
    ap.add_argument("--phase", action="append", help="run one gate/phase, e.g. G12")
    ap.add_argument("--check", action="append", help="run one check ID; repeatable")
    ap.add_argument("--all", action="store_true", help="run all plan checks")
    ap.add_argument("--force", action="store_true", help="rerun PASS checkpoints")
    ap.add_argument("--status", action="store_true", help="show checkpoint status without executing")
    ap.add_argument("--list", action="store_true", help="list available checks")
    ap.add_argument("--reset", action="store_true", help="delete current certification checkpoints")
    ap.add_argument("--fail-fast", action="store_true")
    args = ap.parse_args()

    if args.reset:
        shutil.rmtree(STATE_DIR, ignore_errors=True)
        print("Certification checkpoints reset")
        return 0

    plan = load_plan(Path(args.plan))
    if args.list:
        for c in plan.get("checks", []):
            print(f"{c.get('gate','?')}\t{c['id']}\t{c.get('label', '')}")
        return 0

    lock_handle = acquire_runner_lock()
    fp = source_fingerprint()
    state = load_state(fp)
    checks = plan.get("checks", [])
    selected = []
    if args.check:
        wanted = set(args.check)
        missing = wanted - {c["id"] for c in checks}
        if missing:
            raise SystemExit("unknown check IDs: " + ", ".join(sorted(missing)))
        selected = [c for c in checks if c["id"] in wanted]
    elif args.phase:
        phases = set(args.phase)
        selected = [c for c in checks if c.get("gate") in phases]
    elif args.all:
        selected = checks
    elif args.status:
        selected = checks
    else:
        ap.error("choose --all, --phase, --check, --status, --list, or --reset")

    selected_ids = {c["id"] for c in selected}
    if not args.status:
        for c in selected:
            r = run_check(c, state, force=args.force)
            if args.fail_fast and r["status"] != "PASS":
                break
    summary = summarize(plan, state, selected_ids if not args.status else None)
    print_summary(summary)
    if summary["overall"] == "PASS": return 0
    if summary["overall"] == "PRODUCT FAIL": return 2
    return 3

if __name__ == "__main__":
    raise SystemExit(main())
