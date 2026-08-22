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
test "$(python3 source_fingerprint.py --mode git --commit "$DEPULSE_CANDIDATE_SHA")" = "$DEPULSE_SOURCE_FINGERPRINT"
python3 release_identity.py --verify

eval "$(python3 - <<'PY'
import json,shlex
r=json.load(open('release_identity.json'))
for key,name in [('version','RID_VERSION'),('build_id','RID_BUILD'),('runtime_config','RUNTIME_CONFIG'),('application_bundle','APP_BUNDLE'),('bundle_version','BUNDLE_VERSION')]:
    print(f"{name}={shlex.quote(str(r[key]))}")
PY
)"
test "$RID_VERSION" = "$DEPULSE_VERSION"
test "$RID_BUILD" = "$DEPULSE_BUILD_ID"

APP="$STAGE/$APP_BUNDLE"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp platform-icons/DePulse.icns "$APP/Contents/Resources/DePulse.icns"
SDKROOT="$(xcrun --sdk macosx --show-sdk-path)"
CC="$(xcrun -f clang)"
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
/usr/bin/codesign --verify --deep --strict "$PKG_APP"
PKG_BIN="$PKG_APP/Contents/MacOS/DePulse-arm64"
otool -L "$PKG_BIN" | grep -qi libsqlite3

# The certified package must have one canonical executable contract. A legacy
# Contents/MacOS/De-Pulse binary would indicate a stale/merged/non-canonical
# bundle and must never be accepted as the exact G13/G14 artifact.
test -x "$PKG_APP/Contents/MacOS/DePulseLauncher"
test -x "$PKG_BIN"
test ! -e "$PKG_APP/Contents/MacOS/De-Pulse"
/usr/libexec/PlistBuddy -c 'Print :CFBundleExecutable' "$PKG_APP/Contents/Info.plist" | tee "$OUT/bundle-executable.txt" | grep -qx 'DePulseLauncher'

# Existing packaged-backend/SQLite proof. This deliberately suppresses the
# native window, so it is followed by a real Cocoa/WebKit window lifecycle
# probe below rather than being mislabeled as complete desktop launch proof.
HOME="$HOME_DIR" DEPULSE_HEADLESS=1 PMT_NO_BROWSER=1 "$PKG_APP/Contents/MacOS/DePulseLauncher" >"$OUT/runtime.log" 2>&1 &
PID=$!
cleanup(){ kill -INT "$PID" 2>/dev/null || true; wait "$PID" 2>/dev/null || true; }
trap cleanup EXIT
INSTANCE="$HOME_DIR/Library/Application Support/$RUNTIME_CONFIG/instance.json"
for _ in $(seq 1 180); do [ -s "$INSTANCE" ] && break; sleep 0.1; done
test -s "$INSTANCE"
python3 - "$INSTANCE" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" <<'PY'
import json,sys,time,urllib.request
inst=json.load(open(sys.argv[1])); version=sys.argv[2]; build=sys.argv[3]; base=inst['url'].rstrip('/')
def get(path):
    last=None
    for _ in range(50):
        try:
            with urllib.request.urlopen(base+path,timeout=2) as r:
                return r.status,json.loads(r.read().decode())
        except Exception as e:
            last=e; time.sleep(.15)
    raise SystemExit(f'{path} unavailable: {last}')
status,health=get('/api/health')
assert status==200 and health.get('version')==version and health.get('buildId')==build,health
status,ready=get('/api/ready')
assert status==200 and ready.get('ok') is True,ready
p=ready.get('persistence',{})
assert p.get('backend')=='sqlite' and p.get('ready') is True,p
with urllib.request.urlopen(base+'/',timeout=2) as r:
    assert r.status==200 and len(r.read(256))>0
print('Packaged macOS backend health/readiness/root: PASS')
PY

DB="$HOME_DIR/Library/Application Support/$RUNTIME_CONFIG/depulse-v17.db"
test -s "$DB"
python3 - "$DB" <<'PY'
import sqlite3,sys
con=sqlite3.connect(sys.argv[1])
assert con.execute('pragma integrity_check').fetchone()[0]=='ok'
versions=[r[0] for r in con.execute('select version from schema_migrations order by version')]
assert versions and versions==sorted(set(versions)),versions
assert con.execute('select count(*) from symbol_registry').fetchone()[0] > 0
assert con.execute('select count(*) from identity_state').fetchone()[0] >= 1
print(f'Packaged macOS SQLite: PASS migrations={versions}')
con.close()
PY
cleanup
trap - EXIT

# Real native-window proof. This executes the exact packaged launcher without
# DEPULSE_HEADLESS/PMT_NO_BROWSER so desktop_lifecycle.go must start its JXA
# Cocoa + WKWebView process. We require both backend and window processes to
# stay alive long enough to serve health/root, then repeat against the same
# persisted profile to cover warm relaunch. Cleanup termination is test-owned
# and is not interpreted as a product crash.
GUI_INSTANCE="$GUI_HOME/Library/Application Support/$RUNTIME_CONFIG/instance.json"
GUI_WINDOW_LOG="$GUI_HOME/Library/Application Support/$RUNTIME_CONFIG/native-window.log"

run_native_window_cycle() {
  cycle="$1"
  runtime_log="$OUT/native-window-runtime-${cycle}.log"
  rm -f "$GUI_INSTANCE"
  HOME="$GUI_HOME" DEPULSE_HEADLESS= PMT_NO_BROWSER= "$PKG_APP/Contents/MacOS/DePulseLauncher" >"$runtime_log" 2>&1 &
  gui_pid=$!
  window_pid=0

  dump_native_logs() {
    echo "--- native window runtime ($cycle) ---" >&2
    cat "$runtime_log" >&2 2>/dev/null || true
    echo "--- native-window.log ($cycle) ---" >&2
    cat "$GUI_WINDOW_LOG" >&2 2>/dev/null || true
  }

  for _ in $(seq 1 180); do
    if ! kill -0 "$gui_pid" 2>/dev/null; then
      echo "ERROR: packaged backend exited before native window became ready (cycle=$cycle)" >&2
      dump_native_logs
      wait "$gui_pid" 2>/dev/null || true
      return 1
    fi
    [ -s "$GUI_INSTANCE" ] && break
    sleep 0.1
  done
  if [ ! -s "$GUI_INSTANCE" ]; then
    echo "ERROR: native-window instance evidence missing (cycle=$cycle)" >&2
    dump_native_logs
    kill -INT "$gui_pid" 2>/dev/null || true
    wait "$gui_pid" 2>/dev/null || true
    return 1
  fi

  base_url="$(python3 - "$GUI_INSTANCE" <<'PY'
import json,sys
inst=json.load(open(sys.argv[1]))
print(inst.get('url',''))
PY
)"
  window_pid="$(python3 - "$GUI_INSTANCE" <<'PY'
import json,sys
inst=json.load(open(sys.argv[1]))
print(int(inst.get('windowPid') or 0))
PY
)"
  if [ -z "$base_url" ] || [ "$window_pid" -le 0 ]; then
    echo "ERROR: native window identity incomplete (cycle=$cycle url=$base_url pid=$window_pid)" >&2
    dump_native_logs
    kill -INT "$gui_pid" 2>/dev/null || true
    wait "$gui_pid" 2>/dev/null || true
    return 1
  fi
  if ! kill -0 "$window_pid" 2>/dev/null; then
    echo "ERROR: JXA/Cocoa window process already exited (cycle=$cycle pid=$window_pid)" >&2
    dump_native_logs
    kill -INT "$gui_pid" 2>/dev/null || true
    wait "$gui_pid" 2>/dev/null || true
    return 1
  fi

  python3 - "$base_url" "$DEPULSE_VERSION" "$DEPULSE_BUILD_ID" <<'PY'
import json,sys,time,urllib.request
base=sys.argv[1].rstrip('/'); version=sys.argv[2]; build=sys.argv[3]
last=None
for _ in range(40):
    try:
        with urllib.request.urlopen(base+'/api/health',timeout=2) as r:
            health=json.loads(r.read().decode())
        assert health.get('version')==version and health.get('buildId')==build, health
        with urllib.request.urlopen(base+'/',timeout=2) as r:
            assert r.status==200 and len(r.read(256))>0
        print('Native window backend health/root: PASS')
        break
    except Exception as e:
        last=e; time.sleep(.15)
else:
    raise SystemExit(f'native window backend unavailable: {last}')
PY

  # Immediate framework/JXA aborts normally surface within startup. Keep both
  # processes alive for a bounded dwell so SIGABRT/uncaught-exception exits are
  # observable instead of racing the readiness check.
  sleep 3
  if ! kill -0 "$gui_pid" 2>/dev/null || ! kill -0 "$window_pid" 2>/dev/null; then
    echo "ERROR: native window lifecycle did not remain alive (cycle=$cycle backend=$gui_pid window=$window_pid)" >&2
    dump_native_logs
    kill -TERM "$window_pid" 2>/dev/null || true
    kill -INT "$gui_pid" 2>/dev/null || true
    wait "$gui_pid" 2>/dev/null || true
    return 1
  fi

  # Test-owned cleanup. Terminating the window child causes the existing
  # desktop lifecycle to interrupt the backend; fall back to an explicit
  # backend interrupt if needed. This cleanup exit is intentionally ignored.
  kill -TERM "$window_pid" 2>/dev/null || true
  for _ in $(seq 1 50); do
    kill -0 "$gui_pid" 2>/dev/null || break
    sleep 0.1
  done
  kill -INT "$gui_pid" 2>/dev/null || true
  wait "$gui_pid" 2>/dev/null || true
  rm -f "$GUI_INSTANCE"
  echo "PASS: actual packaged native macOS window cycle=$cycle backendPid=$gui_pid windowPid=$window_pid"
}

run_native_window_cycle fresh
run_native_window_cycle warm

test -s "$GUI_WINDOW_LOG"
grep -q -- "--- DE.PULSE $DEPULSE_VERSION window start" "$GUI_WINDOW_LOG"

ART_SHA="$(shasum -a 256 "$ZIP" | awk '{print $1}')"
python3 - "$OUT/G13-G14-macOS-Apple-Silicon.json" "$PACKAGE" "$ART_SHA" <<'PY'
import json,os,sys,datetime,platform
path,artifact,sha=sys.argv[1:]
data={
 'schema':'DE.PULSE-G13-G14-NATIVE-3','release':'v'+os.environ['DEPULSE_VERSION'],'platform':'macOS Apple Silicon','status':'PASS',
 'certifiedSourceSha':os.environ['DEPULSE_CANDIDATE_SHA'],'sourceFingerprint':os.environ['DEPULSE_SOURCE_FINGERPRINT'],'buildId':os.environ['DEPULSE_BUILD_ID'],
 'artifact':artifact,'artifactSha256':sha,'host':{'os':platform.platform(),'arch':platform.machine()},
 'checks':{
   'exactGitObjectFingerprint':'PASS','nativeBuild':'PASS','arm64Format':'PASS','sqliteLinkage':'PASS','codeSign':'PASS','cleanExtraction':'PASS',
   'canonicalBundleExecutable':'PASS','legacyExecutableAbsent':'PASS',
   'actualPackagedBackendLaunch':'PASS','healthIdentity':'PASS','readySQLite':'PASS','sqliteMigrations':'PASS','sqliteIntegrity':'PASS','rootSurface':'PASS',
   'actualPackagedNativeWindowLaunch':'PASS','nativeWindowStartupDwell':'PASS','warmNativeWindowRelaunch':'PASS'
 },
 'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat()
}
open(path,'w').write(json.dumps(data,indent=2,sort_keys=True)+'\n')
PY
echo "PASS: G13/G14 macOS Apple Silicon v$DEPULSE_VERSION (backend + actual native window + warm relaunch)"
