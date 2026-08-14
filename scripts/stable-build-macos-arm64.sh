#!/bin/bash
set -euo pipefail
R="$(cd "$(dirname "$0")/.."&&pwd)"; S="$R/input/De-Pulse-v18.0.4-STABLE-Source.zip"; SHA='a7fc226577a23ce89ad214183dfad73f11cde1e43b59eba08eedb200ddc0da2a'; FP='9fa0e285cc90d7848aa7676a6169091a540601942f18422a31243b9ba6a8bee8'; B='v18.0.4-stable-native-cross-platform-closure-20260813'; O="$R/out/macos-arm64"; W="$R/.work/macos-arm64"; A="$O/De-Pulse.app"; Z="$O/De-Pulse-v18.0.4-STABLE-macOS-Apple-Silicon.zip"
[[ "$(uname -s)" == Darwin && "$(uname -m)" == arm64 ]]; [[ "$(shasum -a 256 "$S"|awk '{print $1}')" == "$SHA" ]]
rm -rf "$O" "$W"; mkdir -p "$O" "$W/src" "$A/Contents/MacOS" "$A/Contents/Resources"; unzip -q "$S" -d "$W/src"; cd "$W/src"; [[ "$(python3 ci_pipeline.py --fingerprint)" == "$FP" ]]
python3 v18_0_4_stable_promotion_gate.py; python3 v18_0_4_scope_gate.py; python3 v18_baseline_gate.py; CGO_ENABLED=1 go test -count=1 ./...; CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags='-s -w' -o "$A/Contents/MacOS/DePulse-arm64" .; cp platform-icons/DePulse.icns "$A/Contents/Resources/DePulse.icns"
printf '#!/bin/sh\nD=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)\nexec "$D/DePulse-arm64" "$@"\n' > "$A/Contents/MacOS/DePulseLauncher"; chmod +x "$A/Contents/MacOS/DePulseLauncher" "$A/Contents/MacOS/DePulse-arm64"
cat > "$A/Contents/Info.plist" <<'P'
<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"><plist version="1.0"><dict><key>CFBundleDisplayName</key><string>DE.PULSE</string><key>CFBundleExecutable</key><string>DePulseLauncher</string><key>CFBundleIconFile</key><string>DePulse</string><key>CFBundleIdentifier</key><string>com.deivaram.depulse</string><key>CFBundleName</key><string>DE.PULSE</string><key>CFBundlePackageType</key><string>APPL</string><key>CFBundleShortVersionString</key><string>18.0.4</string><key>CFBundleVersion</key><string>18004</string><key>LSMinimumSystemVersion</key><string>12.0</string><key>NSHighResolutionCapable</key><true/></dict></plist>
P
codesign --force --deep --sign - "$A" >/dev/null; codesign --verify --deep --strict "$A"; file "$A/Contents/MacOS/DePulse-arm64"|tee "$O/macos-filetype.txt"|grep -q arm64; otool -L "$A/Contents/MacOS/DePulse-arm64"|tee "$O/macos-linkage.txt"|grep -qi libsqlite3
H="$W/home"; D="$H/Library/Application Support/PersonalMarketTerminal"; mkdir -p "$D"; printf 'DE.PULSE-G14-STABLE-CONTINUITY\n' > "$D/g14-stable-sentinel.txt"; printf '%s' '{"finnhub":"g14-finnhub","alpacaKey":"g14-alpaca-key","alpacaSecret":"g14-alpaca-secret"}' > "$D/secrets.json"; chmod 600 "$D/secrets.json"; SS="$(shasum -a 256 "$D/g14-stable-sentinel.txt"|awk '{print $1}')"; KS="$(shasum -a 256 "$D/secrets.json"|awk '{print $1}')"; PWD="G14-$(python3 -c 'import secrets;print(secrets.token_urlsafe(24))')"
start(){ p="$1"; rm -f "$D/instance.json"; HOME="$H" DEPULSE_HEADLESS=1 "$A/Contents/MacOS/DePulseLauncher" >"$O/$p.out" 2>"$O/$p.err" & PID=$!; for _ in {1..120}; do if [[ -s "$D/instance.json" ]]; then U="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["url"])' "$D/instance.json")"; if curl -fsS "${U}api/health" > "$O/$p-health.json" 2>/dev/null; then python3 - "$O/$p-health.json" "$B" <<'PY'
import json,sys
j=json.load(open(sys.argv[1])); assert j.get('version')=='18.0.4' and j.get('buildId')==sys.argv[2]
PY
return; fi; fi; kill -0 "$PID" 2>/dev/null || { cat "$O/$p.err"; exit 1; }; sleep .25; done; cat "$O/$p.err"; exit 1; }
stop(){ kill "$PID" 2>/dev/null||true; wait "$PID" 2>/dev/null||true; }
trap 'stop 2>/dev/null||true' EXIT; start first; [[ "$(shasum -a 256 "$D/g14-stable-sentinel.txt"|awk '{print $1}')" == "$SS" && "$(shasum -a 256 "$D/secrets.json"|awk '{print $1}')" == "$KS" ]]; C="$O/cookies"; curl -fsS -c "$C" "$U" >/dev/null; curl -fsS -b "$C" "${U}api/auth/status" > "$O/auth1.json"; python3 - "$O/auth1.json" <<'PY'
import json,sys
j=json.load(open(sys.argv[1])); assert j.get('authenticated') is True and j.get('bootstrapRequired') is True and (j.get('principal') or {}).get('role')=='OWNER'
PY
CSRF="$(awk '$6=="depulse_csrf"{print $7}' "$C"|tail -1)"; python3 - "$PWD" > "$W/p.json" <<'PY'
import json,sys;print(json.dumps({'password':sys.argv[1]}))
PY
curl -fsS -b "$C" -c "$C" -H 'Content-Type: application/json' -H "X-DE-PULSE-CSRF: $CSRF" --data-binary @"$W/p.json" "${U}api/auth/set-password" >/dev/null; stop; sleep .4; start restart; C2="$O/cookies2"; python3 - "$PWD" > "$W/l.json" <<'PY'
import json,sys;print(json.dumps({'username':'owner','password':sys.argv[1]}))
PY
curl -fsS -c "$C2" -H 'Content-Type: application/json' --data-binary @"$W/l.json" "${U}api/auth/login" >/dev/null; curl -fsS -b "$C2" "${U}api/auth/status" > "$O/auth2.json"; curl -fsS -b "$C2" "${U}api/bootstrap" > "$O/boot.json"; python3 - "$O/auth2.json" "$O/boot.json" <<'PY'
import json,sys
a=json.load(open(sys.argv[1])); assert a.get('authenticated') is True and a.get('bootstrapRequired') is False and (a.get('principal') or {}).get('role')=='OWNER'; r=json.load(open(sys.argv[2]))['runtime']; assert r['providerRouter']['policyVersion']=='smart-router-v2.0.0' and r['rapidMove']['policyVersion']=='rapid-move-v1.0.0' and r['rapidMove']['coverage']['state']=='TIERED_PARTIAL'
PY
[[ -s "$D/depulse-v17.db" && "$(shasum -a 256 "$D/g14-stable-sentinel.txt"|awk '{print $1}')" == "$SS" && "$(shasum -a 256 "$D/secrets.json"|awk '{print $1}')" == "$KS" ]]; stop; trap - EXIT; /usr/bin/ditto -c -k --sequesterRsrc --keepParent "$A" "$Z"; ZS="$(shasum -a 256 "$Z"|awk '{print $1}')"
python3 - "$O/g14-evidence.json" "$ZS" <<PY
import json,sys,datetime
K=['sourceSha','sourceFingerprint','nativeCompile','nativeSQLite','packageLauncherExecution','releaseIdentity','stableRuntimeContinuity','stableSecretsPreserved','bootstrapOwner','csrfPasswordSetup','restartPersistence','credentialLoginAfterRestart','smartRouterV2Runtime','rapidMoveRuntime','coverageTruth']; json.dump({'schemaVersion':1,'gate':'G14','status':'PASS','release':'v18.0.4 STABLE','buildId':'$B','sourceSha256':'$SHA','sourceFingerprint':'$FP','platform':'macos-arm64','nativeArchitecture':'arm64','checks':{k:True for k in K},'artifact':{'name':'$(basename "$Z")','sha256':sys.argv[2]},'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat()},open(sys.argv[1],'w'),indent=2)
PY
echo 'PASS: G14 macOS Apple Silicon v18.0.4 STABLE'
