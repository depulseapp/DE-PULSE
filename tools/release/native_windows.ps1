Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

foreach($name in 'DEPULSE_CANDIDATE_SHA','DEPULSE_SOURCE_FINGERPRINT','DEPULSE_VERSION','DEPULSE_BUILD_ID') {
  if(-not (Get-Item "Env:$name" -ErrorAction SilentlyContinue).Value) { throw "$name is required" }
}

$root = if($env:GITHUB_WORKSPACE){$env:GITHUB_WORKSPACE}else{(Get-Location).Path}
$out = if($env:DEPULSE_OUTPUT_DIR){$env:DEPULSE_OUTPUT_DIR}else{Join-Path $root 'out/windows-x64'}
$stage = Join-Path $env:RUNNER_TEMP 'depulse-native-windows-stage'
$clean = Join-Path $env:RUNNER_TEMP 'depulse-native-windows-clean'
$profile = Join-Path $env:RUNNER_TEMP 'depulse-native-windows-profile'
Remove-Item -Recurse -Force $out,$stage,$clean,$profile -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $out,$stage,$clean,$profile | Out-Null

Push-Location $root
try {
  if((git rev-parse HEAD).Trim() -ne $env:DEPULSE_CANDIDATE_SHA){ throw 'candidate SHA mismatch' }
  $gitFp=(python tools/release/source_fingerprint.py --mode git --commit $env:DEPULSE_CANDIDATE_SHA).Trim()
  if($gitFp -ne $env:DEPULSE_SOURCE_FINGERPRINT){ throw "Git-object fingerprint mismatch: $gitFp" }
  python tools/release/release_identity.py --verify
  if($LASTEXITCODE -ne 0){ throw 'release identity verification failed' }
  $rid=Get-Content -Raw release_identity.json | ConvertFrom-Json
  if($rid.version -ne $env:DEPULSE_VERSION){ throw 'release version mismatch' }
  if($rid.build_id -ne $env:DEPULSE_BUILD_ID){ throw 'build ID mismatch' }
  $runtimeConfig=[string]$rid.runtime_config

  $exe=Join-Path $stage 'DePulse.exe'
  $env:CGO_ENABLED='0'; $env:GOOS='windows'; $env:GOARCH='amd64'
  go build -trimpath -ldflags '-H=windowsgui' -o $exe .
  if($LASTEXITCODE -ne 0){ throw 'Windows native build failed' }
  Copy-Item 'platform-icons/DePulse.ico' (Join-Path $stage 'DePulse.ico') -Force
  @("DE.PULSE v$($env:DEPULSE_VERSION) STABLE","Build ID: $($env:DEPULSE_BUILD_ID)","Certified source: $($env:DEPULSE_CANDIDATE_SHA)","Certified fingerprint: $($env:DEPULSE_SOURCE_FINGERPRINT)") | Set-Content -Encoding UTF8 (Join-Path $stage 'VERSION.txt')

  $bytes=[IO.File]::ReadAllBytes($exe)
  if($bytes[0] -ne 0x4d -or $bytes[1] -ne 0x5a){ throw 'missing MZ header' }
  $pe=[BitConverter]::ToInt32($bytes,0x3c)
  if([Text.Encoding]::ASCII.GetString($bytes,$pe,4) -ne "PE`0`0"){ throw 'missing PE header' }
  $machine=[BitConverter]::ToUInt16($bytes,$pe+4)
  if($machine -ne 0x8664){ throw ('Windows machine mismatch: 0x{0:x4}' -f $machine) }
  ('PE32+ x86-64 machine=0x{0:x4}' -f $machine) | Set-Content (Join-Path $out 'binary-format.txt')

  $package="De-Pulse-v$($env:DEPULSE_VERSION)-Stable-Windows-x64.zip"
  $zip=Join-Path $out $package
  Compress-Archive -Path (Join-Path $stage '*') -DestinationPath $zip -CompressionLevel Optimal -Force
  Expand-Archive -LiteralPath $zip -DestinationPath $clean -Force
  $pkgExe=Join-Path $clean 'DePulse.exe'
  if(-not (Test-Path $pkgExe)){ throw 'clean-extracted executable missing' }
  $pkgBytes=[IO.File]::ReadAllBytes($pkgExe)
  $pkgPe=[BitConverter]::ToInt32($pkgBytes,0x3c)
  if([BitConverter]::ToUInt16($pkgBytes,$pkgPe+4) -ne 0x8664){ throw 'clean-extracted architecture mismatch' }

  function Invoke-PackagedCycle {
    param(
      [Parameter(Mandatory=$true)][string]$Executable,
      [Parameter(Mandatory=$true)][string]$ExpectedVersion,
      [Parameter(Mandatory=$true)][string]$ExpectedBuild,
      [Parameter(Mandatory=$true)][string]$CycleProfile,
      [Parameter(Mandatory=$true)][string]$Phase,
      [Parameter(Mandatory=$true)][string]$Snapshot
    )
    New-Item -ItemType Directory -Force -Path $CycleProfile | Out-Null
    $env:APPDATA=$CycleProfile
    $env:DEPULSE_HEADLESS='1'
    $env:PMT_NO_BROWSER='1'
    $instance=Join-Path $CycleProfile "$runtimeConfig/instance.json"
    Remove-Item -Force $instance -ErrorAction SilentlyContinue
    $stdout=Join-Path $out "$Phase.stdout.log"
    $stderr=Join-Path $out "$Phase.stderr.log"
    Remove-Item -Force $stdout,$stderr -ErrorAction SilentlyContinue
    $proc=Start-Process -FilePath $Executable -PassThru -RedirectStandardOutput $stdout -RedirectStandardError $stderr
    try {
      for($i=0;$i -lt 180 -and -not (Test-Path $instance);$i++){ Start-Sleep -Milliseconds 100 }
      if(-not (Test-Path $instance)){ throw "packaged Windows app did not create instance.json ($Phase)" }
      $inst=Get-Content -Raw $instance | ConvertFrom-Json
      $base=$inst.url.TrimEnd('/')
      $health=$null
      for($i=0;$i -lt 50 -and $null -eq $health;$i++){
        try { $health=Invoke-RestMethod -Uri "$base/api/health" -TimeoutSec 2 } catch { Start-Sleep -Milliseconds 150 }
      }
      if($null -eq $health){ throw "health endpoint unavailable ($Phase)" }
      if($health.version -ne $ExpectedVersion -or $health.buildId -ne $ExpectedBuild){ throw "runtime identity mismatch ($Phase)" }
      $ready=$null
      for($i=0;$i -lt 50 -and $null -eq $ready;$i++){
        try { $ready=Invoke-RestMethod -Uri "$base/api/ready" -TimeoutSec 2 } catch { Start-Sleep -Milliseconds 150 }
      }
      if($null -eq $ready -or -not $ready.ok -or $ready.persistence.backend -ne 'sqlite' -or -not $ready.persistence.ready){ throw "SQLite readiness failed ($Phase)" }
      $rootResponse=Invoke-WebRequest -Uri "$base/" -TimeoutSec 2 -UseBasicParsing
      if($rootResponse.StatusCode -ne 200 -or $rootResponse.Content.Length -eq 0){ throw "root surface failed ($Phase)" }
      $db=Join-Path $CycleProfile "$runtimeConfig/depulse-v17.db"
      if(-not (Test-Path $db)){ throw "SQLite database missing ($Phase)" }
      python -c "import json,sqlite3,sys; db,out,phase=sys.argv[1:]; c=sqlite3.connect(db); assert c.execute('pragma integrity_check').fetchone()[0]=='ok'; v=[r[0] for r in c.execute('select version from schema_migrations order by version')]; s=c.execute('select count(*) from symbol_registry').fetchone()[0]; i=c.execute('select count(*) from identity_state').fetchone()[0]; assert v and v==sorted(set(v)),v; assert s>0 and i>=1,(s,i); json.dump({'phase':phase,'migrations':v,'symbols':s,'identities':i},open(out,'w'),sort_keys=True); c.close()" $db $Snapshot $Phase
      if($LASTEXITCODE -ne 0){ throw "SQLite audit failed ($Phase)" }
      Write-Host "Packaged Windows runtime health/readiness/root/SQLite: PASS ($Phase)"
    }
    finally {
      if($proc -and -not $proc.HasExited){ Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue }
      Remove-Item -Force $instance -ErrorAction SilentlyContinue
    }
  }

  # Current-package fresh + warm relaunch on one profile.
  Invoke-PackagedCycle -Executable $pkgExe -ExpectedVersion $env:DEPULSE_VERSION -ExpectedBuild $env:DEPULSE_BUILD_ID -CycleProfile $profile -Phase 'current-fresh' -Snapshot (Join-Path $out 'current-fresh.json')
  Invoke-PackagedCycle -Executable $pkgExe -ExpectedVersion $env:DEPULSE_VERSION -ExpectedBuild $env:DEPULSE_BUILD_ID -CycleProfile $profile -Phase 'current-warm' -Snapshot (Join-Path $out 'current-warm.json')

  # T9 previous-Stable upgrade. Invoke the immutable v18.9.1 Stable tag's own
  # certified Windows native harness in a separate PowerShell process rooted at
  # its isolated worktree, then launch the exact current package on that profile.
  $previousTag='v18.9.1-stable'
  $previousSha='e55d8d25b15cec2ffb0f5411bc358bc40b359cf9'
  $previousFp='0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff'
  $previousVersion='18.9.1'
  $previousBuild='v18.9.1-stable-20260821'
  if((git rev-list -n 1 $previousTag).Trim() -ne $previousSha){ throw 'previous Stable tag/candidate mismatch' }
  if((python tools/release/source_fingerprint.py --mode git --commit $previousSha).Trim() -ne $previousFp){ throw 'previous Stable fingerprint mismatch' }

  $previousWorktree=Join-Path $env:RUNNER_TEMP 'depulse-v18.9.1-stable-worktree'
  $previousRunnerTemp=Join-Path $env:RUNNER_TEMP 'depulse-v18.9.1-stable-runner'
  $previousOut=Join-Path $env:RUNNER_TEMP 'depulse-v18.9.1-stable-output'
  $previousClean=Join-Path $env:RUNNER_TEMP 'depulse-v18.9.1-stable-clean'
  $upgradeProfile=Join-Path $env:RUNNER_TEMP 'depulse-native-windows-upgrade-profile'
  Remove-Item -Recurse -Force $previousWorktree,$previousRunnerTemp,$previousOut,$previousClean,$upgradeProfile -ErrorAction SilentlyContinue
  git worktree add --detach $previousWorktree $previousTag | Out-Null
  try {
    New-Item -ItemType Directory -Force -Path $previousRunnerTemp,$previousOut,$previousClean,$upgradeProfile | Out-Null
    $savedEnv=@{}
    foreach($name in 'GITHUB_WORKSPACE','RUNNER_TEMP','DEPULSE_OUTPUT_DIR','DEPULSE_CANDIDATE_SHA','DEPULSE_SOURCE_FINGERPRINT','DEPULSE_VERSION','DEPULSE_BUILD_ID') {
      $savedEnv[$name]=[Environment]::GetEnvironmentVariable($name,'Process')
    }
    try {
      $env:GITHUB_WORKSPACE=$previousWorktree
      $env:RUNNER_TEMP=$previousRunnerTemp
      $env:DEPULSE_OUTPUT_DIR=$previousOut
      $env:DEPULSE_CANDIDATE_SHA=$previousSha
      $env:DEPULSE_SOURCE_FINGERPRINT=$previousFp
      $env:DEPULSE_VERSION=$previousVersion
      $env:DEPULSE_BUILD_ID=$previousBuild
      $pwsh=Join-Path $PSHOME 'pwsh.exe'
      & $pwsh -NoProfile -File (Join-Path $previousWorktree 'tools/release/native_windows.ps1')
      if($LASTEXITCODE -ne 0){ throw 'previous Stable certified Windows harness failed' }
    }
    finally {
      foreach($name in $savedEnv.Keys) {
        [Environment]::SetEnvironmentVariable($name,$savedEnv[$name],'Process')
      }
    }

    $previousZip=Join-Path $previousOut 'De-Pulse-v18.9.1-Stable-Windows-x64.zip'
    if(-not (Test-Path $previousZip)){ throw 'previous Stable certified Windows package missing' }
    Expand-Archive -LiteralPath $previousZip -DestinationPath $previousClean -Force
    $previousPkgExe=Join-Path $previousClean 'DePulse.exe'
    if(-not (Test-Path $previousPkgExe)){ throw 'previous Stable clean-extracted executable missing' }

    $before=Join-Path $out 'upgrade-before.json'
    $after=Join-Path $out 'upgrade-after.json'
    Invoke-PackagedCycle -Executable $previousPkgExe -ExpectedVersion $previousVersion -ExpectedBuild $previousBuild -CycleProfile $upgradeProfile -Phase 'previous-stable' -Snapshot $before
    Invoke-PackagedCycle -Executable $pkgExe -ExpectedVersion $env:DEPULSE_VERSION -ExpectedBuild $env:DEPULSE_BUILD_ID -CycleProfile $upgradeProfile -Phase 'upgrade-current' -Snapshot $after
    python -c "import json,sys; b=json.load(open(sys.argv[1])); a=json.load(open(sys.argv[2])); assert set(b['migrations']).issubset(set(a['migrations'])),(b,a); assert a['symbols']>=b['symbols'],(b,a); assert a['identities']>=b['identities'],(b,a); print('Windows previous-Stable upgrade: PASS')" $before $after
    if($LASTEXITCODE -ne 0){ throw 'previous-Stable upgrade state preservation failed' }
  }
  finally {
    git worktree remove --force $previousWorktree 2>$null | Out-Null
  }

  $sha=(Get-FileHash -Algorithm SHA256 $zip).Hash.ToLowerInvariant()
  $evidence=[ordered]@{
    schema='DE.PULSE-G13-G14-NATIVE-3'; release="v$($env:DEPULSE_VERSION)"; platform='Windows x64'; status='PASS';
    certifiedSourceSha=$env:DEPULSE_CANDIDATE_SHA; sourceFingerprint=$env:DEPULSE_SOURCE_FINGERPRINT; buildId=$env:DEPULSE_BUILD_ID;
    previousStable=[ordered]@{tag='v18.9.1-stable';candidateSha='e55d8d25b15cec2ffb0f5411bc358bc40b359cf9';buildId='v18.9.1-stable-20260821'};
    artifact=$package; artifactSha256=$sha;
    host=[ordered]@{os=[Environment]::OSVersion.VersionString; arch=$env:PROCESSOR_ARCHITECTURE};
    checks=[ordered]@{exactGitObjectFingerprint='PASS';nativeBuild='PASS';peX64Format='PASS';cleanExtraction='PASS';actualPackagedLaunch='PASS';warmPackagedRelaunch='PASS';healthIdentity='PASS';readySQLite='PASS';sqliteMigrations='PASS';sqliteIntegrity='PASS';rootSurface='PASS';previousStableCertifiedHarness='PASS';previousStableTagCandidateBinding='PASS';previousStableProfileUpgrade='PASS';upgradeSQLiteIntegrity='PASS';upgradeStatePreserved='PASS'};
    generatedAt=[DateTime]::UtcNow.ToString('o')
  }
  $evidence | ConvertTo-Json -Depth 6 | Set-Content -Encoding UTF8 (Join-Path $out 'G13-G14-Windows-x64.json')
  Write-Host "PASS: G13/G14 Windows x64 v$($env:DEPULSE_VERSION) (fresh + warm + previous-Stable upgrade)"
}
finally {
  Pop-Location
}
