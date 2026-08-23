#!/usr/bin/env python3
"""Executable full-root disposition and #73 convergence guard."""
from __future__ import annotations
import argparse, json, re, subprocess, sys
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parents[2]
POLICY=ROOT/"governance"/"work-slices"/"ADAPT-ROOT-CONVERGENCE-001"/"root-convergence-policy.json"
ROOT_POLICY=ROOT/"governance"/"root-layout-policy.json"
STATE=ROOT/"governance"/"current-state.json"
VERSIONED=re.compile(r"^v(?:17|18)(?:[_\.-]|$)",re.I)
CLOSED={"READY_FOR_CLOSURE","COMPLETE","COMPLETED","CLOSED","DELIVERED"}

def tracked_root_files():
    r=subprocess.run(("git","ls-files"),cwd=ROOT,check=True,text=True,capture_output=True)
    return sorted(x.strip() for x in r.stdout.splitlines() if x.strip() and "/" not in x.strip())

def classify(name,canonical,remaining):
    if name in canonical:return "KEEP","Canonical steady-state project root owner."
    if name.endswith(".go"):
        if VERSIONED.match(name):return "CONSOLIDATE","Active version-named package-main Go owner; rename or cohesively extract."
        return "KEEP","Current package-main Go source/test owner; relocation is architectural, not cosmetic."
    if VERSIONED.match(name):return "MOVE","Historical/version-scoped non-Go evidence belongs under governed release/history ownership."
    if Path(name).suffix.lower() in {".py",".js",".sh",".ps1"}:return "MOVE","Reusable root tooling/evidence should converge to tools/ci, tools/release or tests."
    if Path(name).suffix.lower() in {".json",".md",".txt"}:return "MOVE","Non-canonical root metadata should converge to governance/release/test ownership."
    return "MOVE","Non-canonical root file requires an explicit canonical owner outside repository root."

def main():
    p=argparse.ArgumentParser();p.add_argument("--json-out");a=p.parse_args();errors=[]
    if not POLICY.is_file() or not ROOT_POLICY.is_file() or not STATE.is_file():
        print("DE.PULSE root convergence: FAIL",file=sys.stderr);return 1
    policy=json.loads(POLICY.read_text());rp=json.loads(ROOT_POLICY.read_text());state=json.loads(STATE.read_text())
    if policy.get("schema")!="DE.PULSE-ROOT-CONVERGENCE-POLICY-1":errors.append("unsupported root-convergence policy schema")
    if policy.get("workSliceId")!="ADAPT-ROOT-CONVERGENCE-001" or policy.get("issue")!=73:errors.append("root-convergence policy identity drift")
    canonical={str(x) for x in rp.get("canonicalRootFiles",[])};remaining={str(x) for x in policy.get("remainingVersionedGo",[])};files=tracked_root_files();rows=[];counts=Counter()
    for name in files:
        disp,reason=classify(name,canonical,remaining);rows.append({"path":name,"disposition":disp,"reason":reason});counts[disp]+=1
    actual={n for n in files if n.endswith(".go") and VERSIONED.match(n)}
    unregistered=sorted(actual-remaining);stale=sorted(remaining-actual)
    if unregistered:errors.append("unregistered version-named Go root owners: "+", ".join(unregistered))
    phase=str(policy.get("phase","")).upper();hist=sorted(n for n in files if VERSIONED.match(n) and not n.endswith(".go"));tools=sorted(n for n in files if Path(n).suffix.lower() in {".py",".js",".sh",".ps1"} and n not in canonical and not VERSIONED.match(n));moves=[r["path"] for r in rows if r["disposition"]=="MOVE"]
    if phase in {"HISTORICAL_NON_GO_CLEAN","TOOLING_CLEAN","GO_CONVERGED","FINAL"} and hist:errors.append("historical v17/v18 non-Go root files remain: "+", ".join(hist[:20]))
    if phase in {"TOOLING_CLEAN","GO_CONVERGED","FINAL"} and tools:errors.append("non-canonical root tooling remains: "+", ".join(tools[:20]))
    if phase in {"GO_CONVERGED","FINAL"} and actual:errors.append("version-named Go root owners remain: "+", ".join(sorted(actual)))
    active=state.get("activeWorkSlice",{});status=str(active.get("status","")).upper()
    if phase=="FINAL" or status in CLOSED:
        if moves:errors.append("closure has unresolved MOVE dispositions: "+", ".join(moves[:20]))
        if stale:errors.append("closure policy still registers converged versioned Go paths: "+", ".join(stale))
    report={"schema":"DE.PULSE-ROOT-DISPOSITION-INVENTORY-1","workSliceId":"ADAPT-ROOT-CONVERGENCE-001","issue":73,"phase":phase,"status":"FAIL" if errors else "PASS","rootFileCount":len(files),"dispositionCounts":dict(sorted(counts.items())),"historicalVersionedNonGoRootCount":len(hist),"versionedGoRootCount":len(actual),"nonCanonicalRootToolCount":len(tools),"unclassifiedCount":0,"rows":rows,"errors":errors}
    if a.json_out:
        out=Path(a.json_out);out.parent.mkdir(parents=True,exist_ok=True);out.write_text(json.dumps(report,indent=2,sort_keys=True)+"\n")
    print("DE.PULSE #73 root convergence");print(f"root files: {len(files)}");print("dispositions: "+", ".join(f"{k}={counts[k]}" for k in sorted(counts)));print(f"historical versioned non-Go root: {len(hist)}");print(f"versioned Go root: {len(actual)}");print(f"non-canonical root tooling: {len(tools)}");print("unclassified root files: 0")
    if errors:
        print("DE.PULSE root convergence: FAIL",file=sys.stderr)
        for e in errors:print(" - "+e,file=sys.stderr)
        return 1
    print("DE.PULSE root convergence: PASS");return 0
if __name__=="__main__":raise SystemExit(main())
