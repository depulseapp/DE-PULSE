//go:build !postgres

package main

func newPostgresPersistenceBackend(config postgresPersistenceConfig) PersistenceBackend {
	return newUnavailablePersistenceBackend("PostgreSQL persistence is not included in this binary; build the hosted target with -tags postgres")
}
