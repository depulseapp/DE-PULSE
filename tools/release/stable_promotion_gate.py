#!/usr/bin/env python3
from pathlib import Path
import subprocess,sys
R=Path(__file__).resolve().parent
go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
model=(R/'app_model.go').read_text()
errors=[]
if 'filepath.Join(base, "PersonalMarketTerminal")' not in go: errors.append('canonical Stable runtime path missing')
if 'PersonalMarketTerminal-v16-TEST' in go: errors.append('retired TEST runtime path still active')
identity=((R/'VERSION.txt').read_text()+'\n'+(R/'renderer/docs/developer.md').read_text())
if 'De-Pulse.app' not in identity: errors.append('standard Stable app identity missing')
for tok in ['`json:"finnhub"`','`json:"alpacaKey"`','`json:"alpacaSecret"`','`json:"groq"`','`json:"openrouter"`','`json:"gemini"`','`json:"fred"`','`json:"bls,omitempty"`','`json:"eia"`','`json:"twelveData"`','`json:"marketaux,omitempty"`']:
    if tok not in model: errors.append('prior Stable credential schema missing '+tok)
p=subprocess.run(['go','test','-count=1','-run','TestV1630Stable(ApplicationUsesCanonicalStableConfig|PreservesPriorStableSecrets|RuntimeUsesCanonicalConfigPath)$','.'],cwd=R,text=True,capture_output=True)
if p.returncode: errors.append('Stable continuity Go regressions failed: '+(p.stdout+p.stderr)[-1200:])
if errors:
    print('Stable Credential Continuity Gate: FAIL')
    print('\n'.join('- '+e for e in errors));sys.exit(1)
print('Stable Credential Continuity Gate: PASS · canonical Stable runtime + prior secrets schema preserved')
