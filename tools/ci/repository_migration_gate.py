#!/usr/bin/env python3
"""Canonical repository migration/root/current-state convergence gate."""
from __future__ import annotations
import json
from pathlib import Path
import tempfile
import current_state_projection_gate, repository_migration_gate_impl, root_convergence_gate, root_ownership_gate
ROOT=Path(__file__).resolve().parents[2]
POLICY_PATH=ROOT/"governance"/"root-layout-policy.json";MIGRATIONS_PATH=ROOT/"governance"/"repository-migrations.json"
def synthesized_migration_policy():
    policy=json.loads(POLICY_PATH.read_text());migrations=json.loads(MIGRATIONS_PATH.read_text());transitional={}
    for row in migrations.get("moves",[]) if isinstance(migrations.get("moves",[]),list) else []:
        if not isinstance(row,dict):continue
        new=str(row.get("newPath","")).strip();owner=str(row.get("owner","")).strip();reason=str(row.get("reason","")).strip();removal=str(row.get("removalCondition","")).strip()
        if new and "/" not in new and all((owner,reason,removal)):transitional[new]={"owner":owner,"reason":reason,"expiry":"REGISTERED_MIGRATION_TARGET","removalCondition":removal}
    final=policy.get("finalRootEvidenceFiles",{})
    if isinstance(final,dict):
        for path,meta in final.items():
            if not isinstance(meta,dict) or "/" in str(path):continue
            owner=str(meta.get("owner","")).strip();reason=str(meta.get("reason","")).strip()
            if owner and reason:transitional[str(path)]={"owner":owner,"reason":reason,"expiry":"FINAL_PACKAGE_MAIN_EVIDENCE","removalCondition":"MOVE_WITH_CAPABILITY_ONLY_WHEN_PRIVATE_PACKAGE_MAIN_ACCESS_CAN_BE_PRESERVED_WITHOUT_TEST_ONLY_PRODUCTION_EXPORTS"}
    compatibility=dict(policy);compatibility["transitionalRootFiles"]=transitional;return compatibility
def main():
    if root_ownership_gate.main()!=0:return 1
    if root_convergence_gate.main()!=0:return 1
    if current_state_projection_gate.main()!=0:return 1
    compatibility=synthesized_migration_policy();original=repository_migration_gate_impl.POLICY_PATH
    with tempfile.NamedTemporaryFile(mode="w",encoding="utf-8",suffix=".json",delete=True) as tmp:
        json.dump(compatibility,tmp,indent=2,sort_keys=True);tmp.write("\n");tmp.flush();repository_migration_gate_impl.POLICY_PATH=Path(tmp.name)
        try:result=repository_migration_gate_impl.main()
        finally:repository_migration_gate_impl.POLICY_PATH=original
    if result!=0:return result
    print("DE.PULSE canonical repository/root/current-state convergence: PASS");return 0
if __name__=="__main__":raise SystemExit(main())
