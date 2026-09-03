//go:build cgo || windows

package main

func (b *sqlitePersistenceBackend) Capabilities() []string {
	return []string{"global-symbol-registry", "canonical-quotes", "quote-history", "evidence-records", "point-in-time-evidence", "decision-lineage", "outcome-history", "derived-feature-store", "user-workspaces"}
}
