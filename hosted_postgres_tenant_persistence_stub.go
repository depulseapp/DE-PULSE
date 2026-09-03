//go:build !postgres

package main

func wrapHostedTenantPostgresBackend(inner PersistenceBackend) PersistenceBackend {
	return inner
}
