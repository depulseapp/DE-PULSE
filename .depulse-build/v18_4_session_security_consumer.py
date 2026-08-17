#!/usr/bin/env python3
from pathlib import Path

p=Path('identity_model.go')
s=p.read_text()
old='''func (s *IdentityService) status(token string) map[string]any {
\tout := map[string]any{"authenticated": false, "bootstrapRequired": s.ownerNeedsPassword()}
\tif p, err := s.resolve(token, false); err == nil {
\t\tout["authenticated"] = true
\t\tout["principal"] = p
\t}
\treturn out
}'''
new='''func (s *IdentityService) status(token string) map[string]any {
\tout := map[string]any{"authenticated": false, "bootstrapRequired": s.ownerNeedsPassword(), "recentAuthentication": false}
\tif p, err := s.resolve(token, false); err == nil {
\t\tout["authenticated"] = true
\t\tout["principal"] = p
\t\tout["recentAuthentication"] = s.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL)
\t}
\treturn out
}'''
if s.count(old) != 1:
    raise SystemExit(f'status consumer: expected exactly one match, got {s.count(old)}')
p.write_text(s.replace(old,new,1))
