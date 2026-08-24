//go:build !cgo && !windows

package main

import (
	"context"
	"encoding/json"
)

const fallbackInstrumentIdentityFeatureKey = "instrument-identity"
const fallbackInstrumentIdentityFeatureVersion = "v1"

// The unsupported/non-native file fallback reuses its existing derived-feature
// container rather than creating another sidecar store. Supported native
// platforms use the dedicated SQLite table.
func (b *filePersistenceBackend) SaveInstrumentIdentities(ctx context.Context, records []InstrumentIdentityRecord) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if b.data.Features == nil {
		b.data.Features = map[string]DerivedFeatureRecord{}
	}
	written := 0
	for _, record := range canonicalInstrumentIdentities(records) {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		payload, err := json.Marshal(record)
		if err != nil {
			return written, err
		}
		key := record.Symbol + "|" + fallbackInstrumentIdentityFeatureKey + "|" + fallbackInstrumentIdentityFeatureVersion
		b.data.Features[key] = DerivedFeatureRecord{
			Symbol:         record.Symbol,
			FeatureKey:     fallbackInstrumentIdentityFeatureKey,
			FeatureVersion: fallbackInstrumentIdentityFeatureVersion,
			AsOf:           record.ObservedAt,
			SourceHash:     record.Source + "|" + record.ProviderAssetID,
			Payload:        payload,
		}
		written++
	}
	if written == 0 {
		return 0, nil
	}
	return written, b.persistLocked()
}

func (b *filePersistenceBackend) LoadInstrumentIdentities(ctx context.Context) ([]InstrumentIdentityRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := []InstrumentIdentityRecord{}
	for _, feature := range b.data.Features {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if feature.FeatureKey != fallbackInstrumentIdentityFeatureKey || feature.FeatureVersion != fallbackInstrumentIdentityFeatureVersion {
			continue
		}
		var record InstrumentIdentityRecord
		if err := json.Unmarshal(feature.Payload, &record); err != nil {
			continue
		}
		if normalized, ok := normalizeInstrumentIdentity(record); ok {
			out = append(out, normalized)
		}
	}
	return canonicalInstrumentIdentities(out), nil
}
