package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	persistenceArchiveSchemaVersion = 1
	persistenceRestoreModeEmpty     = "empty"
	persistenceRestoreModeReplace   = "replace"
	persistenceExportPathEnv        = "DEPULSE_PERSISTENCE_EXPORT_PATH"
	persistenceRestorePathEnv       = "DEPULSE_PERSISTENCE_RESTORE_PATH"
	persistenceRestoreModeEnv       = "DEPULSE_PERSISTENCE_RESTORE_MODE"
)

type PersistenceCanonicalQuoteRecord struct {
	Symbol            string `json:"symbol"`
	Quote             Quote  `json:"quote"`
	ProviderTimestamp int64  `json:"providerTimestamp,omitempty"`
	ReceivedTimestamp int64  `json:"receivedTimestamp,omitempty"`
	PersistedAt       int64  `json:"persistedAt"`
	Source            string `json:"source,omitempty"`
	DataState         string `json:"dataState,omitempty"`
}

type PersistenceQuoteHistoryRecord struct {
	Symbol            string  `json:"symbol"`
	Bucket            int64   `json:"bucket"`
	ProviderTimestamp int64   `json:"providerTimestamp,omitempty"`
	ReceivedTimestamp int64   `json:"receivedTimestamp,omitempty"`
	Price             float64 `json:"price"`
	Bid               float64 `json:"bid,omitempty"`
	Ask               float64 `json:"ask,omitempty"`
	Volume            float64 `json:"volume,omitempty"`
	Source            string  `json:"source,omitempty"`
	DataState         string  `json:"dataState,omitempty"`
}

type PersistenceArchive struct {
	SchemaVersion     int                               `json:"schemaVersion"`
	SourceBackend     string                            `json:"sourceBackend"`
	SourceStoreSchema int                               `json:"sourceStoreSchema"`
	ExportedAt        int64                             `json:"exportedAt"`
	Symbols           []SymbolRegistryRecord            `json:"symbols,omitempty"`
	CanonicalQuotes   []PersistenceCanonicalQuoteRecord `json:"canonicalQuotes,omitempty"`
	QuoteHistory      []PersistenceQuoteHistoryRecord   `json:"quoteHistory,omitempty"`
	Evidence          []EvidenceRecord                  `json:"evidence,omitempty"`
	Decisions         []DecisionLineageRecord           `json:"decisions,omitempty"`
	Outcomes          []OutcomeHistoryRecord            `json:"outcomes,omitempty"`
	Features          []DerivedFeatureRecord            `json:"features,omitempty"`
	HasIdentity       bool                              `json:"hasIdentity"`
	Identity          IdentityPersistentState           `json:"identity,omitempty"`
	UserWorkspaces    []UserWorkspace                   `json:"userWorkspaces,omitempty"`
}

type persistenceArchiveEnvelope struct {
	SchemaVersion int                `json:"schemaVersion"`
	SHA256        string             `json:"sha256"`
	Archive       PersistenceArchive `json:"archive"`
}

type persistenceArchiveBackend interface {
	ExportPersistenceArchive(context.Context) (PersistenceArchive, error)
	RestorePersistenceArchive(context.Context, PersistenceArchive, string) error
}

func normalizePersistenceRestoreMode(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", persistenceRestoreModeEmpty:
		return persistenceRestoreModeEmpty, nil
	case persistenceRestoreModeReplace:
		return persistenceRestoreModeReplace, nil
	default:
		return "", fmt.Errorf("unsupported persistence restore mode %q", raw)
	}
}

func persistenceArchiveDigest(archive PersistenceArchive) (string, []byte, error) {
	raw, err := json.Marshal(archive)
	if err != nil {
		return "", nil, err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), raw, nil
}

func writePersistenceArchiveFile(path string, archive PersistenceArchive) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("persistence archive path is required")
	}
	digest, _, err := persistenceArchiveDigest(archive)
	if err != nil {
		return err
	}
	envelope := persistenceArchiveEnvelope{SchemaVersion: persistenceArchiveSchemaVersion, SHA256: digest, Archive: archive}
	raw, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, append(raw, '\n'), 0600)
}

func readPersistenceArchiveFile(path string) (PersistenceArchive, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return PersistenceArchive{}, errors.New("persistence archive path is required")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return PersistenceArchive{}, err
	}
	var envelope persistenceArchiveEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return PersistenceArchive{}, fmt.Errorf("decode persistence archive: %w", err)
	}
	if envelope.SchemaVersion != persistenceArchiveSchemaVersion || envelope.Archive.SchemaVersion != persistenceArchiveSchemaVersion {
		return PersistenceArchive{}, fmt.Errorf("unsupported persistence archive schema envelope=%d archive=%d", envelope.SchemaVersion, envelope.Archive.SchemaVersion)
	}
	digest, _, err := persistenceArchiveDigest(envelope.Archive)
	if err != nil {
		return PersistenceArchive{}, err
	}
	if !strings.EqualFold(strings.TrimSpace(envelope.SHA256), digest) {
		return PersistenceArchive{}, errors.New("persistence archive integrity check failed")
	}
	return envelope.Archive, nil
}

func (p *PersistenceManager) ExportArchiveFile(ctx context.Context, path string) (PersistenceArchive, error) {
	if p == nil || p.backend == nil {
		return PersistenceArchive{}, errors.New("persistence unavailable")
	}
	backend, ok := p.backend.(persistenceArchiveBackend)
	if !ok {
		return PersistenceArchive{}, fmt.Errorf("persistence backend %s does not support archive export", p.backend.Name())
	}
	// Flush the current coalesced batch before taking the backend's consistent snapshot.
	p.flushPending()
	archive, err := backend.ExportPersistenceArchive(ctx)
	if err != nil {
		p.recordPersistenceFailure(err)
		return PersistenceArchive{}, err
	}
	archive.SchemaVersion = persistenceArchiveSchemaVersion
	archive.SourceBackend = p.backend.Name()
	archive.ExportedAt = time.Now().UTC().UnixMilli()
	if archive.SourceStoreSchema == 0 {
		if stats, statsErr := p.backend.Stats(ctx); statsErr == nil {
			archive.SourceStoreSchema = stats.SchemaVersion
		}
	}
	if _, err := EvidenceAsKnownAt(archive.Evidence, 1<<63-1); err != nil {
		return PersistenceArchive{}, fmt.Errorf("archive temporal evidence validation: %w", err)
	}
	if err := writePersistenceArchiveFile(path, archive); err != nil {
		return PersistenceArchive{}, err
	}
	return archive, nil
}

func (p *PersistenceManager) RestoreArchiveFile(ctx context.Context, path, mode string) (PersistenceArchive, error) {
	if p == nil || p.backend == nil {
		return PersistenceArchive{}, errors.New("persistence unavailable")
	}
	backend, ok := p.backend.(persistenceArchiveBackend)
	if !ok {
		return PersistenceArchive{}, fmt.Errorf("persistence backend %s does not support archive restore", p.backend.Name())
	}
	normalizedMode, err := normalizePersistenceRestoreMode(mode)
	if err != nil {
		return PersistenceArchive{}, err
	}
	archive, err := readPersistenceArchiveFile(path)
	if err != nil {
		return PersistenceArchive{}, err
	}
	if _, err := EvidenceAsKnownAt(archive.Evidence, 1<<63-1); err != nil {
		return PersistenceArchive{}, fmt.Errorf("restore temporal evidence validation: %w", err)
	}
	currentIdentity := IdentityPersistentState{}
	if normalizedMode == persistenceRestoreModeReplace {
		currentIdentity, err = p.backend.LoadIdentityState(ctx)
		if err != nil {
			p.recordPersistenceFailure(err)
			return PersistenceArchive{}, fmt.Errorf("load current deletion tombstones before replace restore: %w", err)
		}
	}
	archive = enforceArchiveAccountDeletionPrivacy(archive, currentIdentity)
	if err := backend.RestorePersistenceArchive(ctx, archive, normalizedMode); err != nil {
		p.recordPersistenceFailure(err)
		return PersistenceArchive{}, err
	}
	p.recordPersistenceHealthy()
	p.refreshStoreStats()
	return archive, nil
}

func configureStartupPersistenceRestore(p *PersistenceManager) error {
	path := strings.TrimSpace(os.Getenv(persistenceRestorePathEnv))
	if path == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := p.RestoreArchiveFile(ctx, path, os.Getenv(persistenceRestoreModeEnv))
	return err
}

func configureStartupPersistenceExport(p *PersistenceManager) error {
	path := strings.TrimSpace(os.Getenv(persistenceExportPathEnv))
	if path == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := p.ExportArchiveFile(ctx, path)
	return err
}
