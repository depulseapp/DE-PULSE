$ErrorActionPreference='Stop'
$Root=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$SrcZip=Join-Path $Root 'input\De-Pulse-v18.0.4-TEST-Source.zip'
$ExpectedSha='69f87c3fd2f94e9adc8fe5e3fb273843bd5b07cc019f62eeaf1e4ad688adb66f'
$ExpectedFingerprint='a1ff742baf04176d6122da338fc2360c0d21c7220977c6901299b755ce2cfc5b'
$BuildId='v18.0.4-test-native-windows-lifecycle-g14-harness-hardening-20260813'
$Out=Join-Path $Root 'out\windows-x64'; $Work=Join-Path $Root '.work\windows-x64'; $Exe=Join-Path $Out 'De-Pulse-v18.0.4-TEST.exe'; $Zip=Join-Path $Out 'De-Pulse-v18.0.4-TEST-Windows-x64.zip'
$ActualSha=(Get-FileHash -Algorithm SHA256 $SrcZip).Hash.ToLowerInvariant(); if($ActualSha -ne $ExpectedSha){throw "source SHA mismatch: $ActualSha"}
Remove-Item -Recurse -Force $Out,$Work -ErrorAction SilentlyContinue; New-Item -ItemType Directory -Force -Path $Out,(Join-Path $Work 'src') | Out-Null
Expand-Archive -LiteralPath $SrcZip -DestinationPath (Join-Path $Work 'src') -Force; Set-Location (Join-Path $Work 'src')
$fp=(python ci_pipeline.py --fingerprint).Trim(); if($fp -ne $ExpectedFingerprint){throw "fingerprint mismatch: $fp"}
go test -count=1 ./...
$env:CGO_ENABLED='0'; $env:GOOS='windows'; $env:GOARCH='amd64'; go build -trimpath -ldflags='-s -w' -o $Exe .
go version -m $Exe | Tee-Object -FilePath (Join-Path $Out 'windows-go-build-metadata.txt')
$Meta=Get-Content (Join-Path $Out 'windows-go-build-metadata.txt') -Raw; if($Meta -notmatch 'GOOS=windows' -or $Meta -notmatch 'GOARCH=amd64'){throw 'Windows metadata mismatch'}
$env:APPDATA=Join-Path $Work 'appdata'; New-Item -ItemType Directory -Force -Path $env:APPDATA | Out-Null
$Stable=Join-Path $env:APPDATA 'PersonalMarketTerminal'; $Target=Join-Path $env:APPDATA 'PersonalMarketTerminal-v18.0.4-TEST'; New-Item -ItemType Directory -Force -Path $Stable | Out-Null
$Sentinel=Join-Path $Stable 'g14-stable-sentinel.txt'; Set-Content -LiteralPath $Sentinel -Value "DE.PULSE-G14-STABLE-SENTINEL`n" -Encoding ascii -NoNewline
$Secrets=Join-Path $Stable 'secrets.json'
$SecretJson='{"finnhub":"g14-finnhub","alpacaKey":"g14-alpaca-key","alpacaSecret":"g14-alpaca-secret","groq":"g14-groq","openrouter":"g14-openrouter","gemini":"g14-gemini","fred":"g14-fred","bls":"g14-bls","eia":"g14-eia","twelveData":"g14-twelve","marketaux":"g14-marketaux"}'
Set-Content -LiteralPath $Secrets -Value $SecretJson -Encoding utf8 -NoNewline
$StableBefore=(Get-FileHash -Algorithm SHA256 $Sentinel).Hash.ToLowerInvariant()
$StableSecretsBefore=(Get-FileHash -Algorithm SHA256 $Secrets).Hash.ToLowerInvariant()
$Password='G14-'+[Convert]::ToBase64String([Security.Cryptography.RandomNumberGenerator]::GetBytes(24)).Replace('/','A').Replace('+','B').Replace('=','C')
$env:DEPULSE_HEADLESS='1'
function Start-Depulse([string]$Phase) {
  $script:Proc=Start-Process -FilePath $Exe -PassThru -RedirectStandardOutput (Join-Path $Out "runtime-$Phase-stdout.log") -RedirectStandardError (Join-Path $Out "runtime-$Phase-stderr.log")
  $script:Instance=Join-Path $Target 'instance.json'
  for($i=0;$i -lt 80 -and -not (Test-Path $script:Instance);$i++){Start-Sleep -Milliseconds 250}
  if(-not (Test-Path $script:Instance)){throw "instance.json not created ($Phase)"}
  $script:Url=(Get-Content $script:Instance -Raw | ConvertFrom-Json).url
}
function Stop-Depulse { if($script:Proc){Stop-Process -Id $script:Proc.Id -Force -ErrorAction SilentlyContinue; Wait-Process -Id $script:Proc.Id -ErrorAction SilentlyContinue; $script:Proc=$null} }
try {
  Start-Depulse 'first'
  $Health=Invoke-RestMethod -Uri ($Url+'api/health'); $Health | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-first-health.json')
  if($Health.version -ne '18.0.4' -or $Health.buildId -ne $BuildId){throw 'runtime release identity mismatch'}
  $Session=New-Object Microsoft.PowerShell.Commands.WebRequestSession; Invoke-WebRequest -Uri $Url -WebSession $Session -UseBasicParsing | Out-File (Join-Path $Out 'runtime-first-root.html')
  $Auth=Invoke-RestMethod -Uri ($Url+'api/auth/status') -WebSession $Session; $Auth | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-first-auth-status.json')
  if(-not $Auth.authenticated -or -not $Auth.bootstrapRequired -or $Auth.principal.userId -ne 'bootstrap-owner' -or $Auth.principal.role -ne 'OWNER'){throw 'bootstrap owner audit failed'}
  $Migration=Join-Path $Target '.v18.0.4-test-profile-migration.json'; if(-not (Test-Path $Migration)){throw 'migration marker missing'}
  if((Get-Content (Join-Path $Target 'g14-stable-sentinel.txt') -Raw).Trim() -ne 'DE.PULSE-G14-STABLE-SENTINEL'){throw 'stable sentinel not cloned'}
  $TargetSecrets=Join-Path $Target 'secrets.json'; if(-not (Test-Path $TargetSecrets)){throw 'Stable secrets not cloned'}
  $TargetSecretsHash=(Get-FileHash -Algorithm SHA256 $TargetSecrets).Hash.ToLowerInvariant(); if($TargetSecretsHash -ne $StableSecretsBefore){throw 'migrated Stable secrets mismatch'}
  $DangerousAcl=(Get-Acl $TargetSecrets).Access | Where-Object { $_.AccessControlType -eq 'Allow' -and $_.IdentityReference.Value -match '(?i)Everyone|ANONYMOUS LOGON' -and $_.FileSystemRights.ToString() -match '(?i)Write|Modify|FullControl' }
  if($DangerousAcl){throw 'migrated secrets ACL grants broad write access'}
  $StableAfter=(Get-FileHash -Algorithm SHA256 $Sentinel).Hash.ToLowerInvariant(); if($StableAfter -ne $StableBefore){throw 'Stable profile was modified'}
  $StableSecretsAfter=(Get-FileHash -Algorithm SHA256 $Secrets).Hash.ToLowerInvariant(); if($StableSecretsAfter -ne $StableSecretsBefore){throw 'Stable secrets were modified'}
  $Db=Join-Path $Target 'depulse-v17.db'; if(-not (Test-Path $Db)){throw 'SQLite database not created'}
  $Csrf=$Session.Cookies.GetCookies([uri]$Url)['depulse_csrf'].Value; if([string]::IsNullOrWhiteSpace($Csrf)){throw 'CSRF cookie missing'}
  $Body=@{password=$Password}|ConvertTo-Json -Compress
  Invoke-RestMethod -Method Post -Uri ($Url+'api/auth/set-password') -WebSession $Session -Headers @{'X-DE-PULSE-CSRF'=$Csrf} -ContentType 'application/json' -Body $Body | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-password-setup.json')
  $PostAuth=Invoke-RestMethod -Uri ($Url+'api/auth/status') -WebSession $Session; if(-not $PostAuth.authenticated -or $PostAuth.bootstrapRequired){throw 'password setup did not complete'}
  Stop-Depulse; Start-Sleep -Milliseconds 400
  Start-Depulse 'restart'
  $Session2=New-Object Microsoft.PowerShell.Commands.WebRequestSession
  $LoginBody=@{username='owner';password=$Password}|ConvertTo-Json -Compress
  $Login=Invoke-RestMethod -Method Post -Uri ($Url+'api/auth/login') -WebSession $Session2 -ContentType 'application/json' -Body $LoginBody; $Login | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-restart-login.json')
  $Auth2=Invoke-RestMethod -Uri ($Url+'api/auth/status') -WebSession $Session2; $Auth2 | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-restart-auth-status.json')
  if(-not $Auth2.authenticated -or $Auth2.bootstrapRequired -or $Auth2.principal.role -ne 'OWNER'){throw 'credential login after restart failed'}
  $Boot=Invoke-RestMethod -Uri ($Url+'api/bootstrap') -WebSession $Session2; $Boot | ConvertTo-Json -Depth 30 | Set-Content -Encoding utf8 (Join-Path $Out 'runtime-bootstrap.json')
  if($Boot.runtime.providerRouter.policyVersion -ne 'smart-router-v2.0.0'){throw 'Smart Router v2 runtime policy missing'}
  if($Boot.runtime.rapidMove.policyVersion -ne 'rapid-move-v1.0.0'){throw 'Rapid Move runtime policy missing'}
  if($Boot.runtime.rapidMove.coverage.state -ne 'TIERED_PARTIAL'){throw 'Rapid Move Coverage Truth mismatch'}
  if(-not (Test-Path $Db)){throw 'SQLite database missing after restart'}
  $StableFinal=(Get-FileHash -Algorithm SHA256 $Sentinel).Hash.ToLowerInvariant(); if($StableFinal -ne $StableBefore){throw 'Stable profile changed after restart'}
  $StableSecretsFinal=(Get-FileHash -Algorithm SHA256 $Secrets).Hash.ToLowerInvariant(); if($StableSecretsFinal -ne $StableSecretsBefore){throw 'Stable secrets changed after restart'}
} finally { Stop-Depulse }
@"
DE.PULSE v18.0.4 TEST - Windows x64
Build ID: $BuildId
G14 native runtime: migration + isolation + secure password + restart login + Router/Rapid Move PASS.
"@ | Set-Content -Encoding utf8 (Join-Path $Out 'README.txt')
Compress-Archive -Path $Exe,(Join-Path $Out 'README.txt') -DestinationPath $Zip -Force
$ZipSha=(Get-FileHash -Algorithm SHA256 $Zip).Hash.ToLowerInvariant(); "$ZipSha  De-Pulse-v18.0.4-TEST-Windows-x64.zip" | Set-Content -Encoding ascii (Join-Path $Out 'De-Pulse-v18.0.4-TEST-Windows-x64.zip.sha256')
"$ExpectedSha  De-Pulse-v18.0.4-TEST-Source.zip" | Set-Content -Encoding ascii (Join-Path $Out 'verified-source.sha256')
$Evidence=[ordered]@{
 schemaVersion=1; gate='G14'; status='PASS'; release='v18.0.4 TEST'; buildId=$BuildId; sourceSha256=$ExpectedSha; sourceFingerprint=$ExpectedFingerprint; platform='windows-x64'; nativeArchitecture='amd64'; runner=@{os=$env:ImageOS;imageVersion=$env:ImageVersion};
 checks=[ordered]@{sourceSha=$true;sourceFingerprint=$true;nativeCompile=$true;nativeSQLite=$true;packageLauncherExecution=$true;releaseIdentity=$true;stableProfileMigration=$true;stableSecretsMigration=$true;stableIsolation=$true;bootstrapOwner=$true;csrfPasswordSetup=$true;restartPersistence=$true;credentialLoginAfterRestart=$true;smartRouterV2Runtime=$true;rapidMoveRuntime=$true;coverageTruth=$true};
 artifact=@{name='De-Pulse-v18.0.4-TEST-Windows-x64.zip';sha256=$ZipSha}; generatedAt=[DateTime]::UtcNow.ToString('o')
}
$Evidence|ConvertTo-Json -Depth 10 | Set-Content -Encoding utf8 (Join-Path $Out 'g14-evidence.json')
Write-Host 'PASS: G14 Windows x64 native package + migration + auth restart + Router/Rapid Move runtime audit'
