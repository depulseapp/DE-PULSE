#!/bin/bash
set -euo pipefail

: "${DEPULSE_CANDIDATE_SHA:?}"
: "${DEPULSE_SOURCE_FINGERPRINT:?}"
: "${DEPULSE_VERSION:?}"
: "${DEPULSE_BUILD_ID:?}"

ROOT="${GITHUB_WORKSPACE:-$(pwd)}"
OUT="${DEPULSE_OUTPUT_DIR:-$ROOT/out/macos-arm64}"
STAGE="${RUNNER_TEMP:-/tmp}/depulse-native-macos-stage"
CLEAN="${RUNNER_TEMP:-/tmp}/depulse-native-macos-clean"
HOME_DIR="${RUNNER_TEMP:-/tmp}/depulse-native-macos-home"
GUI_HOME="${RUNNER_TEMP:-/tmp}/depulse-native-macos-gui-home"
rm -rf "$OUT" "$STAGE" "$CLEAN" "$HOME_DIR" "$GUI_HOME"
mkdir -p "$OUT" "$STAGE" "$CLEAN" "$HOME_DIR" "$GUI_HOME"

cd "$ROOT"
test "$(git rev-parse HEAD)" = "$DEPULSE_CANDIDATE_SHA"
test "$(python3 tools/release/source_fingerprint.py --mode git --commit "$DEPULSE_CANDIDATE_SHA")" = "$DEPULSE_SOURCE_FINGERPRINT"
python3 tools/release/release_identity.py --verify

eval "$(python3 - <<'PY'
import json,shlex
r=json.load(open('release_identity.json'))
for key,name in [('version','RID_VERSION'),('build_id','RID_BUILD'),('runtime_config','RUNTIME_CONFIG'),('application_bundle','APP_BUNDLE'),('bundle_version','BUNDLE_VERSION')]:
    print(f"{name}={shlex.quote(str(r[key]))}")
PY
)"
test "$RID_VERSION" = "$DEPULSE_VERSION"
test "$RID_BUILD" = "$DEPULSE_BUILD_ID"

SDKROOT="$(xcrun --sdk macosx --show-sdk-path)"
CC="$(xcrun -f clang)"
APP="$STAGE/$APP_BUNDLE"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp platform-icons/DePulse.icns "$APP/Contents/Resources/DePulse.icns"
MACOSX_DEPLOYMENT_TARGET=11.0 CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="$CC" \
  CGO_CFLAGS="-isysroot $SDKROOT -arch arm64" \
  CGO_LDFLAGS="-isysroot $SDKROOT -arch arm64" \
  go build -trimpath -o "$APP/Contents/MacOS/DePulse-arm64" .
cat > "$APP/Contents/MacOS/DePulseLauncher" <<'LAUNCH'
#!/bin/sh
APP_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
case "$(/usr/bin/uname -m)" in
  arm64) exec "$APP_DIR/DePulse-arm64" "$@" ;;
  *) /usr/bin/osascript -e 'display dialog "DE.PULSE requires Apple Silicon for this package." buttons {"OK"} default button "OK" with icon stop' >/dev/null 2>&1; exit 1 ;;
esac
LAUNCH
chmod +x "$APP/Contents/MacOS/DePulseLauncher" "$APP/Contents/MacOS/DePulse-arm64"
cat > "$APP/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleDevelopmentRegion</key><string>en</string>
<key>CFBundleDisplayName</key><string>DE.PULSE</string>
<key>CFBundleExecutable</key><string>DePulseLauncher</string>
<key>CFBundleIconFile</key><string>DePulse</string>
<key>CFBundleIdentifier</key><string>com.deivaram.depulse</string>
<key>CFBundleName</key><string>DE.PULSE</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleShortVersionString</key><string>$DEPULSE_VERSION</string>
<key>CFBundleVersion</key><string>$BUNDLE_VERSION</string>
<key>LSMinimumSystemVersion</key><string>11.0</string>
<key>NSHighResolutionCapable</key><true/>
</dict></plist>
PLIST
file "$APP/Contents/MacOS/DePulse-arm64" | tee "$OUT/binary-format.txt" | grep -q 'Mach-O 64-bit executable arm64'
otool -L "$APP/Contents/MacOS/DePulse-arm64" | tee "$OUT/sqlite-linkage.txt" | grep -qi libsqlite3
/usr/bin/codesign --force --deep --sign - "$APP"
/usr/bin/codesign --verify --deep --strict "$APP"

PACKAGE="De-Pulse-v${DEPULSE_VERSION}-Stable-macOS-Apple-Silicon.zip"
ZIP="$OUT/$PACKAGE"
/usr/bin/ditto -c -k --sequesterRsrc --keepParent "$APP" "$ZIP"
/usr/bin/ditto -x -k "$ZIP" "$CLEAN"
PKG_APP="$CLEAN/$APP_BUNDLE"
PKG_BIN="$PKG_APP/Contents/MacOS/DePulse-arm64"
/usr/bin/codesign --verify --deep --strict "$PKG_APP"
otool -L "$PKG_BIN" | grep -qi libsqlite3
test -x "$PKG_APP/Contents/MacOS/DePulseLauncher"
test -x "$PKG_BIN"
test ! -e "$PKG_APP/Contents/MacOS/De-Pulse"
/usr/libexec/PlistBuddy -c 'Print :CFBundleExecutable' "$PKG_APP/Contents/Info.plist" | tee "$OUT/bundle-executable.txt" | grep -qx 'DePulseLauncher'

snapshot_db() {
  db="$1"; out="$2"; phase="$3"
  python3 - "$db" "$out" "$phase" <<'PY'
import json,sqlite3,sys
path,out,phase=sys.argv[1:]
con=sqlite3.connect(path)
assert con.execute('pragma integrity_check').fetchone()[0]=='ok'
versions=[r[0] for r in con.execute('select version from schema_migrations order by version')]
symbols=con.execute('select count(*) from symbol_registry').fetchone()[0]
identities=con.execute('select count(*) from identity_state').fetchone()[0]
assert versions and versions==sorted(set(versions)),versions
assert symbols>0 and identities>=1,(symbols,identities)
json.dump({'phase':phase,'migrations':versions,'symbols':symbols,'identities':identities},open(out,'w'),sort_keys=True)
con.close()
PY
}

run_headless_cycle() {
  launcher="$1"; expected_version="$2"; expected_build="$3"; home="$4"; phase="$5"; snapshot="$6"
  instance="$home/Library/Application Support/$RUNTIME_CONFIG/instance.json"
  db="$home/Library/Application Support/$RUNTIME_CONFIG/depulse-v17.db"
  mkdir -p "$home"
  rm -f "$instance"
  HOME="$home" DEPULSE_HEADLESS=1 PMT_NO_BROWSER=1 "$launcher" >"$OUT/${phase}.log" 2>&1 &
  pid=$!
  for _ in $(seq 1 180); do [ -s "$instance" ] && break; sleep 0.1; done
  test -s "$instance"
  python3 - "$instance" "$expected_version" "$expected_build" <<'PY'
import json,sys,time,urllib.request
inst=json.load(open(sys.argv[1])); version=sys.argv[2]; build=sys.argv[3]; base=inst['url'].rstrip('/')
last=None
for _ in range(50):
    try:
        with urllib.request.urlopen(base+'/api/health',timeout=2) as r: health=json.loads(r.read().decode())
        with urllib.request.urlopen(base+'/api/ready',timeout=2) as r: ready=json.loads(r.read().decode())
        with urllib.request.urlopen(base+'/',timeout=2) as r: root=r.read(256)
        assert health.get('version')==version and health.get('buildId')==build,health
        assert ready.get('ok') is True and ready.get('persistence',{}).get('backend')=='sqlite' and ready.get('persistence',{}).get('ready') is True,ready
        assert root
        break
    except Exception as exc:
        last=exc; time.sleep(.15)
else:
    raise SystemExit(f'packaged runtime unavailable: {last}')
PY
  kill -INT "$pid" 2>/dev/null || true
  wait "$pid" 2>/dev/null || true
  test -s "$db"
  snapshot_db "$db" "$snapshot" "$phase"
  rm -f "$instance"
}

run_headless_cycle "$PKG_APP/Contents/MacOS/DePulseLauncher" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" "$HOME_DIR" current-fresh "$OUT/current-fresh.json"
run_headless_cycle "$PKG_APP/Contents/MacOS/DePulseLauncher" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" "$HOME_DIR" current-warm "$OUT/current-warm.json"

GUI_INSTANCE="$GUI_HOME/Library/Application Support/$RUNTIME_CONFIG/instance.json"
GUI_WINDOW_LOG="$GUI_HOME/Library/Application Support/$RUNTIME_CONFIG/native-window.log"
GUI_DB="$GUI_HOME/Library/Application Support/$RUNTIME_CONFIG/depulse-v17.db"
run_native_window_cycle() {
  cycle="$1"; runtime_log="$OUT/native-window-runtime-${cycle}.log"
  rm -f "$GUI_INSTANCE"
  HOME="$GUI_HOME" DEPULSE_HEADLESS= PMT_NO_BROWSER= "$PKG_APP/Contents/MacOS/DePulseLauncher" >"$runtime_log" 2>&1 &
  gui_pid=$!; window_pid=0
  stop_cycle() {
    if [ "$window_pid" -gt 0 ]; then
      kill -TERM "$window_pid" 2>/dev/null || true
      for _ in $(seq 1 30); do kill -0 "$window_pid" 2>/dev/null || break; sleep 0.1; done
      kill -KILL "$window_pid" 2>/dev/null || true
    fi
    kill -INT "$gui_pid" 2>/dev/null || true
    for _ in $(seq 1 30); do kill -0 "$gui_pid" 2>/dev/null || break; sleep 0.1; done
    kill -KILL "$gui_pid" 2>/dev/null || true
    wait "$gui_pid" 2>/dev/null || true
    rm -f "$GUI_INSTANCE"
  }
  for _ in $(seq 1 180); do
    kill -0 "$gui_pid" 2>/dev/null || { cat "$runtime_log" >&2 || true; return 1; }
    [ -s "$GUI_INSTANCE" ] && break; sleep 0.1
  done
  test -s "$GUI_INSTANCE"
  native_identity="$(python3 - "$GUI_INSTANCE" <<'PY'
import json,sys
inst=json.load(open(sys.argv[1])); print(inst.get('url','')); print(int(inst.get('windowPid') or 0))
PY
)"
  base_url="$(printf '%s\n' "$native_identity" | sed -n '1p')"
  window_pid="$(printf '%s\n' "$native_identity" | sed -n '2p')"
  test -n "$base_url"; test "$window_pid" -gt 0; kill -0 "$window_pid"
  python3 - "$base_url" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" <<'PY'
import json,sys,time,urllib.request
base=sys.argv[1].rstrip('/'); version=sys.argv[2]; build=sys.argv[3]; last=None
for _ in range(40):
    try:
        with urllib.request.urlopen(base+'/api/health',timeout=2) as r: health=json.loads(r.read().decode())
        with urllib.request.urlopen(base+'/',timeout=2) as r: root=r.read(256)
        assert health.get('version')==version and health.get('buildId')==build,health
        assert root
        break
    except Exception as exc: last=exc; time.sleep(.15)
else: raise SystemExit(f'native window backend unavailable: {last}')
PY
  sleep 3
  kill -0 "$gui_pid"; kill -0 "$window_pid"
  stop_cycle
}
run_native_window_cycle fresh
snapshot_db "$GUI_DB" "$OUT/gui-fresh.json" gui-fresh
run_native_window_cycle warm
snapshot_db "$GUI_DB" "$OUT/gui-warm.json" gui-warm
test -s "$GUI_WINDOW_LOG"
test "$(grep -c -- "--- DE.PULSE $DEPULSE_VERSION window start" "$GUI_WINDOW_LOG")" -ge 2
! grep -q -- 'protocol does not exist' "$GUI_WINDOW_LOG"

# T9 previous-Stable upgrade. Invoke the immutable v18.9.1 Stable tag's own
# certified native harness in an isolated worktree/temp root, then run the exact
# current package on the same prior-Stable profile and compare persisted state.
PREVIOUS_TAG="v18.9.1-stable"
PREVIOUS_SHA="e55d8d25b15cec2ffb0f5411bc358bc40b359cf9"
PREVIOUS_FP="0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff"
PREVIOUS_VERSION="18.9.1"
PREVIOUS_BUILD_ID="v18.9.1-stable-20260821"
PREVIOUS_WORKTREE="${RUNNER_TEMP:-/tmp}/depulse-v18.9.1-stable-worktree"
PREVIOUS_RUNNER_TEMP="${RUNNER_TEMP:-/tmp}/depulse-v18.9.1-stable-runner"
PREVIOUS_OUT="${RUNNER_TEMP:-/tmp}/depulse-v18.9.1-stable-output"
PREVIOUS_CLEAN="${RUNNER_TEMP:-/tmp}/depulse-v18.9.1-stable-clean"
UPGRADE_HOME="${RUNNER_TEMP:-/tmp}/depulse-native-macos-upgrade-home"
rm -rf "$PREVIOUS_WORKTREE" "$PREVIOUS_RUNNER_TEMP" "$PREVIOUS_OUT" "$PREVIOUS_CLEAN" "$UPGRADE_HOME"
test "$(git rev-list -n 1 "$PREVIOUS_TAG")" = "$PREVIOUS_SHA"
test "$(python3 tools/release/source_fingerprint.py --mode git --commit "$PREVIOUS_SHA")" = "$PREVIOUS_FP"
git worktree add --detach "$PREVIOUS_WORKTREE" "$PREVIOUS_TAG" >/dev/null
(
  cd "$PREVIOUS_WORKTREE"
  GITHUB_WORKSPACE="$PREVIOUS_WORKTREE" \
  RUNNER_TEMP="$PREVIOUS_RUNNER_TEMP" \
  DEPULSE_OUTPUT_DIR="$PREVIOUS_OUT" \
  DEPULSE_CANDIDATE_SHA="$PREVIOUS_SHA" \
  DEPULSE_SOURCE_FINGERPRINT="$PREVIOUS_FP" \
  DEPULSE_VERSION="$PREVIOUS_VERSION" \
  DEPULSE_BUILD_ID="$PREVIOUS_BUILD_ID" \
  bash tools/release/native_macos.sh
)
PREVIOUS_ZIP="$PREVIOUS_OUT/De-Pulse-v18.9.1-Stable-macOS-Apple-Silicon.zip"
test -s "$PREVIOUS_ZIP"
mkdir -p "$PREVIOUS_CLEAN"
/usr/bin/ditto -x -k "$PREVIOUS_ZIP" "$PREVIOUS_CLEAN"
PREVIOUS_PKG_APP="$PREVIOUS_CLEAN/$APP_BUNDLE"
/usr/bin/codesign --verify --deep --strict "$PREVIOUS_PKG_APP"
run_headless_cycle "$PREVIOUS_PKG_APP/Contents/MacOS/DePulseLauncher" "$PREVIOUS_VERSION" "$PREVIOUS_BUILD_ID" "$UPGRADE_HOME" previous-stable "$OUT/upgrade-before.json"
run_headless_cycle "$PKG_APP/Contents/MacOS/DePulseLauncher" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" "$UPGRADE_HOME" upgrade-current "$OUT/upgrade-after.json"
python3 - "$OUT/upgrade-before.json" "$OUT/upgrade-after.json" <<'PY'
import json,sys
before=json.load(open(sys.argv[1])); after=json.load(open(sys.argv[2]))
assert set(before['migrations']).issubset(set(after['migrations'])),(before,after)
assert after['symbols'] >= before['symbols'],(before,after)
assert after['identities'] >= before['identities'],(before,after)
print('macOS previous-Stable upgrade: PASS')
PY
git worktree remove --force "$PREVIOUS_WORKTREE" >/dev/null

ART_SHA="$(shasum -a 256 "$ZIP" | awk '{print $1}')"
python3 - "$OUT/G13-G14-macOS-Apple-Silicon.json" "$PACKAGE" "$ART_SHA" <<'PY'
import json,os,sys,datetime,platform
path,artifact,sha=sys.argv[1:]
data={
 'schema':'DE.PULSE-G13-G14-NATIVE-3','release':'v'+os.environ['DEPULSE_VERSION'],'platform':'macOS Apple Silicon','status':'PASS',
 'certifiedSourceSha':os.environ['DEPULSE_CANDIDATE_SHA'],'sourceFingerprint':os.environ['DEPULSE_SOURCE_FINGERPRINT'],'buildId':os.environ['DEPULSE_BUILD_ID'],
 'previousStable':{'tag':'v18.9.1-stable','candidateSha':'e55d8d25b15cec2ffb0f5411bc358bc40b359cf9','buildId':'v18.9.1-stable-20260821'},
 'artifact':artifact,'artifactSha256':sha,'host':{'os':platform.platform(),'arch':platform.machine()},
 'checks':{
   'exactGitObjectFingerprint':'PASS','nativeBuild':'PASS','arm64Format':'PASS','sqliteLinkage':'PASS','codeSign':'PASS','cleanExtraction':'PASS','canonicalBundleExecutable':'PASS','legacyExecutableAbsent':'PASS',
   'actualPackagedBackendLaunch':'PASS','warmPackagedBackendRelaunch':'PASS','healthIdentity':'PASS','readySQLite':'PASS','sqliteMigrations':'PASS','sqliteIntegrity':'PASS','rootSurface':'PASS',
   'actualPackagedNativeWindowLaunch':'PASS','nativeWindowStartupDwell':'PASS','warmNativeWindowRelaunch':'PASS','warmProfileSQLiteReuse':'PASS','nativeWindowProtocolResolution':'PASS','nativeWindowCleanup':'PASS',
   'previousStableCertifiedHarness':'PASS','previousStableTagCandidateBinding':'PASS','previousStableProfileUpgrade':'PASS','upgradeSQLiteIntegrity':'PASS','upgradeStatePreserved':'PASS'
 },
 'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat()
}
open(path,'w').write(json.dumps(data,indent=2,sort_keys=True)+'\n')
PY
echo "PASS: G13/G14 macOS Apple Silicon v$DEPULSE_VERSION (fresh + warm native + previous-Stable upgrade)"
