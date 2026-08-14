#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRCZIP="$ROOT/input/De-Pulse-v18.0.4-TEST-Source.zip"
EXPECTED_SHA="69f87c3fd2f94e9adc8fe5e3fb273843bd5b07cc019f62eeaf1e4ad688adb66f"
EXPECTED_FP="a1ff742baf04176d6122da338fc2360c0d21c7220977c6901299b755ce2cfc5b"
BUILD_ID="v18.0.4-test-native-windows-lifecycle-g14-harness-hardening-20260813"
OUT="$ROOT/out/macos-arm64"
WORK="$ROOT/.work/macos-arm64"
APP="$OUT/De-Pulse-v18.0.4-TEST.app"
ZIP="$OUT/De-Pulse-v18.0.4-TEST-macOS-Apple-Silicon.zip"
[[ "$(uname -s)" == "Darwin" ]] || { echo 'ERROR: macOS runner required'; exit 1; }
[[ "$(uname -m)" == "arm64" ]] || { echo 'ERROR: Apple Silicon arm64 runner required'; exit 1; }
for x in go clang codesign otool curl python3 shasum ditto; do command -v "$x" >/dev/null; done
ACTUAL_SHA="$(shasum -a 256 "$SRCZIP" | awk '{print $1}')"
[[ "$ACTUAL_SHA" == "$EXPECTED_SHA" ]] || { echo "ERROR: source SHA mismatch: $ACTUAL_SHA"; exit 1; }
rm -rf "$OUT" "$WORK"; mkdir -p "$OUT" "$WORK/src" "$APP/Contents/MacOS" "$APP/Contents/Resources"
unzip -q "$SRCZIP" -d "$WORK/src"
cd "$WORK/src"
[[ "$(python3 ci_pipeline.py --fingerprint)" == "$EXPECTED_FP" ]] || { echo 'ERROR: source fingerprint mismatch'; exit 1; }
CGO_ENABLED=1 go test -count=1 ./...
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags='-s -w' -o "$APP/Contents/MacOS/DePulse-arm64" .
cp platform-icons/DePulse.icns "$APP/Contents/Resources/DePulse.icns"
cat > "$APP/Contents/MacOS/DePulseLauncher" <<'LAUNCHER'
#!/bin/sh
APP_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec "$APP_DIR/DePulse-arm64" "$@"
LAUNCHER
chmod +x "$APP/Contents/MacOS/DePulseLauncher" "$APP/Contents/MacOS/DePulse-arm64"
cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleDevelopmentRegion</key><string>en</string>
<key>CFBundleDisplayName</key><string>DE.PULSE v18.0.4 TEST</string>
<key>CFBundleExecutable</key><string>DePulseLauncher</string>
<key>CFBundleIconFile</key><string>DePulse</string>
<key>CFBundleIdentifier</key><string>com.deivaram.depulse.v1804test</string>
<key>CFBundleInfoDictionaryVersion</key><string>6.0</string>
<key>CFBundleName</key><string>DE.PULSE v18.0.4 TEST</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleShortVersionString</key><string>18.0.4</string>
<key>CFBundleVersion</key><string>18004</string>
<key>LSMinimumSystemVersion</key><string>12.0</string>
<key>NSHighResolutionCapable</key><true/>
</dict></plist>
PLIST
codesign --force --deep --sign - "$APP"
codesign --verify --deep --strict "$APP"
file "$APP/Contents/MacOS/DePulse-arm64" | tee "$OUT/macos-filetype.txt"
grep -q 'arm64' "$OUT/macos-filetype.txt"
otool -L "$APP/Contents/MacOS/DePulse-arm64" | tee "$OUT/macos-linkage.txt"
grep -qi 'libsqlite3' "$OUT/macos-linkage.txt"
go version -m "$APP/Contents/MacOS/DePulse-arm64" | tee "$OUT/macos-go-build-metadata.txt"
RHOME="$WORK/runtime-home"
BASE="$RHOME/Library/Application Support"
STABLE="$BASE/PersonalMarketTerminal"
TARGET="$BASE/PersonalMarketTerminal-v18.0.4-TEST"
mkdir -p "$STABLE"
printf 'DE.PULSE-G14-STABLE-SENTINEL\n' > "$STABLE/g14-stable-sentinel.txt"
cat > "$STABLE/secrets.json" <<'SECRETS'
{"finnhub":"g14-finnhub","alpacaKey":"g14-alpaca-key","alpacaSecret":"g14-alpaca-secret","groq":"g14-groq","openrouter":"g14-openrouter","gemini":"g14-gemini","fred":"g14-fred","bls":"g14-bls","eia":"g14-eia","twelveData":"g14-twelve","marketaux":"g14-marketaux"}
SECRETS
chmod 600 "$STABLE/secrets.json"
STABLE_BEFORE="$(shasum -a 256 "$STABLE/g14-stable-sentinel.txt" | awk '{print $1}')"
STABLE_SECRETS_BEFORE="$(shasum -a 256 "$STABLE/secrets.json" | awk '{print $1}')"
PASSWORD="G14-$(python3 - <<'PY'
import secrets
print(secrets.token_urlsafe(24))
PY
)"
start_app() {
  local phase="$1"
  HOME="$RHOME" DEPULSE_HEADLESS=1 "$APP/Contents/MacOS/DePulseLauncher" >"$OUT/runtime-${phase}-stdout.log" 2>"$OUT/runtime-${phase}-stderr.log" &
  APP_PID=$!
  INSTANCE="$TARGET/instance.json"
  for _ in {1..80}; do [[ -s "$INSTANCE" ]] && break; sleep 0.25; done
  [[ -s "$INSTANCE" ]] || { echo "ERROR: instance.json not created ($phase)"; exit 1; }
  URL="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["url"])' "$INSTANCE")"
}
stop_app() {
  kill "$APP_PID" 2>/dev/null || true
  wait "$APP_PID" 2>/dev/null || true
}
trap 'stop_app 2>/dev/null || true' EXIT
start_app first
curl -fsS "${URL}api/health" | tee "$OUT/runtime-first-health.json"
python3 - "$OUT/runtime-first-health.json" "$BUILD_ID" <<'PY'
import json,sys
j=json.load(open(sys.argv[1])); assert j.get('version')=='18.0.4'; assert j.get('buildId')==sys.argv[2]
PY
COOKIE1="$OUT/runtime-first-cookies.txt"
curl -fsS -c "$COOKIE1" "$URL" > "$OUT/runtime-first-root.html"
curl -fsS -b "$COOKIE1" "${URL}api/auth/status" | tee "$OUT/runtime-first-auth-status.json"
python3 - "$OUT/runtime-first-auth-status.json" <<'PY'
import json,sys
j=json.load(open(sys.argv[1])); assert j.get('authenticated') is True; assert j.get('bootstrapRequired') is True
p=j.get('principal') or {}; assert p.get('userId')=='bootstrap-owner'; assert p.get('role')=='OWNER'
PY
[[ -f "$TARGET/.v18.0.4-test-profile-migration.json" ]]
[[ "$(cat "$TARGET/g14-stable-sentinel.txt")" == 'DE.PULSE-G14-STABLE-SENTINEL' ]]
[[ -f "$TARGET/secrets.json" ]]
[[ "$(shasum -a 256 "$TARGET/secrets.json" | awk '{print $1}')" == "$STABLE_SECRETS_BEFORE" ]]
[[ "$(stat -f '%Lp' "$TARGET/secrets.json")" == '600' ]]
STABLE_AFTER="$(shasum -a 256 "$STABLE/g14-stable-sentinel.txt" | awk '{print $1}')"
STABLE_SECRETS_AFTER="$(shasum -a 256 "$STABLE/secrets.json" | awk '{print $1}')"
[[ "$STABLE_BEFORE" == "$STABLE_AFTER" ]]
[[ "$STABLE_SECRETS_BEFORE" == "$STABLE_SECRETS_AFTER" ]]
DB="$TARGET/depulse-v17.db"; [[ -s "$DB" ]]
CSRF="$(awk '$6=="depulse_csrf" {print $7}' "$COOKIE1" | tail -1)"; [[ -n "$CSRF" ]]
python3 - "$PASSWORD" > "$WORK/password.json" <<'PY'
import json,sys
print(json.dumps({'password':sys.argv[1]}))
PY
curl -fsS -b "$COOKIE1" -c "$COOKIE1" -H 'Content-Type: application/json' -H "X-DE-PULSE-CSRF: $CSRF" --data-binary @"$WORK/password.json" "${URL}api/auth/set-password" > "$OUT/runtime-password-setup.json"
curl -fsS -b "$COOKIE1" "${URL}api/auth/status" > "$OUT/runtime-post-password-auth-status.json"
python3 - "$OUT/runtime-post-password-auth-status.json" <<'PY'
import json,sys
j=json.load(open(sys.argv[1])); assert j.get('authenticated') is True; assert j.get('bootstrapRequired') is False
PY
stop_app
sleep 0.4
start_app restart
curl -fsS "${URL}api/health" > "$OUT/runtime-restart-health.json"
COOKIE2="$OUT/runtime-restart-cookies.txt"
python3 - "$PASSWORD" > "$WORK/login.json" <<'PY'
import json,sys
print(json.dumps({'username':'owner','password':sys.argv[1]}))
PY
curl -fsS -c "$COOKIE2" -H 'Content-Type: application/json' --data-binary @"$WORK/login.json" "${URL}api/auth/login" > "$OUT/runtime-restart-login.json"
curl -fsS -b "$COOKIE2" "${URL}api/auth/status" > "$OUT/runtime-restart-auth-status.json"
curl -fsS -b "$COOKIE2" "${URL}api/bootstrap" > "$OUT/runtime-bootstrap.json"
python3 - "$OUT/runtime-restart-auth-status.json" "$OUT/runtime-bootstrap.json" <<'PY'
import json,sys
a=json.load(open(sys.argv[1])); assert a.get('authenticated') is True; assert a.get('bootstrapRequired') is False
p=a.get('principal') or {}; assert p.get('role')=='OWNER'
b=json.load(open(sys.argv[2])); r=b['runtime']; assert r['providerRouter']['policyVersion']=='smart-router-v2.0.0'; assert r['rapidMove']['policyVersion']=='rapid-move-v1.0.0'; assert r['rapidMove']['coverage']['state']=='TIERED_PARTIAL'
PY
[[ -s "$DB" ]]
[[ -f "$TARGET/.v18.0.4-test-profile-migration.json" ]]
STABLE_FINAL="$(shasum -a 256 "$STABLE/g14-stable-sentinel.txt" | awk '{print $1}')"; [[ "$STABLE_BEFORE" == "$STABLE_FINAL" ]]
STABLE_SECRETS_FINAL="$(shasum -a 256 "$STABLE/secrets.json" | awk '{print $1}')"; [[ "$STABLE_SECRETS_BEFORE" == "$STABLE_SECRETS_FINAL" ]]
stop_app; trap - EXIT
/usr/bin/ditto -c -k --sequesterRsrc --keepParent "$APP" "$ZIP"
ZIP_SHA="$(shasum -a 256 "$ZIP" | awk '{print $1}')"
printf '%s  %s\n' "$ZIP_SHA" "$(basename "$ZIP")" > "$OUT/De-Pulse-v18.0.4-TEST-macOS-Apple-Silicon.zip.sha256"
printf '%s  De-Pulse-v18.0.4-TEST-Source.zip\n' "$EXPECTED_SHA" > "$OUT/verified-source.sha256"
python3 - "$OUT/g14-evidence.json" "$ZIP_SHA" "${ImageOS:-macOS}" "${ImageVersion:-unknown}" <<PY
import json,sys,datetime
out,zipsha,imageos,imagever=sys.argv[1:]
e={
 'schemaVersion':1,'gate':'G14','status':'PASS','release':'v18.0.4 TEST','buildId':'$BUILD_ID',
 'sourceSha256':'$EXPECTED_SHA','sourceFingerprint':'$EXPECTED_FP','platform':'macos-arm64','nativeArchitecture':'arm64',
 'runner':{'os':imageos,'imageVersion':imagever},
 'checks':{k:True for k in ['sourceSha','sourceFingerprint','nativeCompile','nativeSQLite','packageLauncherExecution','releaseIdentity','stableProfileMigration','stableSecretsMigration','stableIsolation','bootstrapOwner','csrfPasswordSetup','restartPersistence','credentialLoginAfterRestart','smartRouterV2Runtime','rapidMoveRuntime','coverageTruth']},
 'artifact':{'name':'De-Pulse-v18.0.4-TEST-macOS-Apple-Silicon.zip','sha256':zipsha},
 'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat()
}
json.dump(e,open(out,'w'),indent=2); open(out,'a').write('\n')
PY
echo 'PASS: G14 Apple Silicon native package + migration + auth restart + Router/Rapid Move runtime audit'
