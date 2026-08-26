package hostedpersistence

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

const policySchema = "DE.PULSE-HOSTED-POSTGRES-POLICY-1"

var canonicalEnvironments = []string{"dev", "test", "stage", "prod"}

//go:embed postgres_policy_v1.json
var policyJSON []byte

type EnvironmentPolicy struct {
	Schema                   string `json:"schema"`
	MaxOpenConnections       int    `json:"maxOpenConnections"`
	HARequired               bool   `json:"haRequired"`
	BackupEncryptionRequired bool   `json:"backupEncryptionRequired"`
	PITRRequired             bool   `json:"pitrRequired"`
	MaxRPO                   string `json:"maxRPO"`
	MaxRTO                   string `json:"maxRTO"`
	RestoreDrillMaxAge       string `json:"restoreDrillMaxAge"`
}

type policyManifest struct {
	Schema            string                       `json:"schema"`
	Version           string                       `json:"version"`
	RLSDisposition    string                       `json:"rlsDisposition"`
	MigrationStrategy string                       `json:"migrationStrategy"`
	RollbackPolicy    string                       `json:"rollbackPolicy"`
	Environments      map[string]EnvironmentPolicy `json:"environments"`
}

type RuntimeDeclaration struct {
	Environment        string
	DatabaseURL        string
	MaxOpenConnections int
	PolicyVersion      string
	RLSDisposition     string
	MigrationStrategy  string
	HAReady            bool
	BackupEncrypted    bool
	PITRReady          bool
	RPO                string
	RTO                string
	RestoreDrillAt     string
	RestoreDrillID     string
	RollbackPolicy     string
	RollbackVerified   bool
}

var (
	manifestOnce sync.Once
	manifestData policyManifest
	manifestErr  error
)

func loadManifest() (policyManifest, error) {
	manifestOnce.Do(func() {
		if err := json.Unmarshal(policyJSON, &manifestData); err != nil {
			manifestErr = fmt.Errorf("decode hosted PostgreSQL policy: %w", err)
			return
		}
		manifestErr = validateManifest(manifestData)
	})
	return manifestData, manifestErr
}

func validateManifest(m policyManifest) error {
	if m.Schema != policySchema {
		return fmt.Errorf("unsupported hosted PostgreSQL policy schema %q", m.Schema)
	}
	if strings.TrimSpace(m.Version) == "" {
		return errors.New("hosted PostgreSQL policy version is empty")
	}
	if m.RLSDisposition != "application-authorization-canonical-identity" {
		return errors.New("hosted PostgreSQL policy must explicitly retain tenant authorization in the canonical identity owner until RLS is adopted")
	}
	if m.MigrationStrategy != "expand-contract" {
		return errors.New("hosted PostgreSQL policy must use expand-contract migrations")
	}
	if strings.TrimSpace(m.RollbackPolicy) == "" {
		return errors.New("hosted PostgreSQL rollback policy is empty")
	}
	if len(m.Environments) != len(canonicalEnvironments) {
		return errors.New("hosted PostgreSQL policy must define exactly dev/test/stage/prod")
	}
	seenSchemas := map[string]struct{}{}
	for _, environment := range canonicalEnvironments {
		policy, ok := m.Environments[environment]
		if !ok {
			return fmt.Errorf("hosted PostgreSQL policy missing %s", environment)
		}
		if policy.Schema != "depulse_"+environment {
			return fmt.Errorf("hosted PostgreSQL %s schema must be depulse_%s", environment, environment)
		}
		if _, exists := seenSchemas[policy.Schema]; exists {
			return fmt.Errorf("hosted PostgreSQL schema %q is reused across environments", policy.Schema)
		}
		seenSchemas[policy.Schema] = struct{}{}
		if policy.MaxOpenConnections < 2 || policy.MaxOpenConnections > 128 {
			return fmt.Errorf("hosted PostgreSQL %s pool capacity is outside supported bounds", environment)
		}
		if !policy.BackupEncryptionRequired {
			return fmt.Errorf("hosted PostgreSQL %s must require encrypted backups", environment)
		}
		for name, raw := range map[string]string{"maxRPO": policy.MaxRPO, "maxRTO": policy.MaxRTO, "restoreDrillMaxAge": policy.RestoreDrillMaxAge} {
			if _, err := positiveDuration(raw); err != nil {
				return fmt.Errorf("hosted PostgreSQL %s %s: %w", environment, name, err)
			}
		}
	}
	return nil
}

func PolicyFor(environment string) (string, string, string, string, EnvironmentPolicy, error) {
	m, err := loadManifest()
	if err != nil {
		return "", "", "", "", EnvironmentPolicy{}, err
	}
	environment = strings.ToLower(strings.TrimSpace(environment))
	policy, ok := m.Environments[environment]
	if !ok {
		return "", "", "", "", EnvironmentPolicy{}, fmt.Errorf("unsupported hosted PostgreSQL environment %q; expected dev, test, stage, or prod", environment)
	}
	return m.Version, m.RLSDisposition, m.MigrationStrategy, m.RollbackPolicy, policy, nil
}

func Validate(declaration RuntimeDeclaration, now time.Time) error {
	environment := strings.ToLower(strings.TrimSpace(declaration.Environment))
	version, rlsDisposition, migrationStrategy, rollbackPolicy, policy, err := PolicyFor(environment)
	if err != nil {
		return err
	}
	if strings.TrimSpace(declaration.PolicyVersion) != version {
		return fmt.Errorf("hosted PostgreSQL policy drift: version does not match %s", environment)
	}
	if strings.TrimSpace(declaration.RLSDisposition) != rlsDisposition {
		return fmt.Errorf("hosted PostgreSQL policy drift: RLS disposition does not match %s", environment)
	}
	if strings.TrimSpace(declaration.MigrationStrategy) != migrationStrategy {
		return fmt.Errorf("hosted PostgreSQL policy drift: migration strategy does not match %s", environment)
	}
	if strings.TrimSpace(declaration.RollbackPolicy) != rollbackPolicy {
		return fmt.Errorf("hosted PostgreSQL policy drift: rollback policy does not match %s", environment)
	}
	if !declaration.RollbackVerified {
		return errors.New("hosted PostgreSQL rollback readiness is not verified")
	}
	if declaration.MaxOpenConnections < 2 || declaration.MaxOpenConnections > policy.MaxOpenConnections {
		return fmt.Errorf("hosted PostgreSQL pool exceeds %s capacity policy", environment)
	}
	if err := validateDatabaseURL(declaration.DatabaseURL, policy.Schema); err != nil {
		return err
	}
	if policy.HARequired && !declaration.HAReady {
		return fmt.Errorf("hosted PostgreSQL %s requires HA/failover readiness", environment)
	}
	if policy.BackupEncryptionRequired && !declaration.BackupEncrypted {
		return fmt.Errorf("hosted PostgreSQL %s requires encrypted backups", environment)
	}
	if policy.PITRRequired && !declaration.PITRReady {
		return fmt.Errorf("hosted PostgreSQL %s requires point-in-time recovery readiness", environment)
	}
	if err := validateRecoveryDuration("RPO", declaration.RPO, policy.MaxRPO); err != nil {
		return err
	}
	if err := validateRecoveryDuration("RTO", declaration.RTO, policy.MaxRTO); err != nil {
		return err
	}
	if err := validateRestoreDrill(declaration.RestoreDrillAt, declaration.RestoreDrillID, policy.RestoreDrillMaxAge, now); err != nil {
		return err
	}
	return nil
}

func validateDatabaseURL(raw, expectedSchema string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u == nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Host == "" {
		return errors.New("hosted PostgreSQL database URL must be an explicit postgres/postgresql URL")
	}
	query := u.Query()
	if query.Get("sslmode") != "verify-full" {
		return errors.New("hosted PostgreSQL database URL must require sslmode=verify-full")
	}
	if query.Get("search_path") != expectedSchema {
		return fmt.Errorf("hosted PostgreSQL database URL must bind search_path to %s", expectedSchema)
	}
	return nil
}

func validateRecoveryDuration(name, actualRaw, maximumRaw string) error {
	actual, err := positiveDuration(actualRaw)
	if err != nil {
		return fmt.Errorf("hosted PostgreSQL %s: %w", name, err)
	}
	maximum, err := positiveDuration(maximumRaw)
	if err != nil {
		return fmt.Errorf("hosted PostgreSQL policy %s: %w", name, err)
	}
	if actual > maximum {
		return fmt.Errorf("hosted PostgreSQL %s %s exceeds policy maximum %s", name, actual, maximum)
	}
	return nil
}

func positiveDuration(raw string) (time.Duration, error) {
	value, err := time.ParseDuration(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("duration %q must be positive", strings.TrimSpace(raw))
	}
	return value, nil
}

func validateRestoreDrill(rawAt, drillID, maxAgeRaw string, now time.Time) error {
	if strings.TrimSpace(drillID) == "" {
		return errors.New("hosted PostgreSQL restore drill evidence id is required")
	}
	drillAt, err := time.Parse(time.RFC3339, strings.TrimSpace(rawAt))
	if err != nil {
		return errors.New("hosted PostgreSQL restore drill timestamp must be RFC3339")
	}
	if drillAt.After(now.Add(5 * time.Minute)) {
		return errors.New("hosted PostgreSQL restore drill timestamp is in the future")
	}
	maxAge, err := positiveDuration(maxAgeRaw)
	if err != nil {
		return fmt.Errorf("hosted PostgreSQL restore drill policy: %w", err)
	}
	if now.Sub(drillAt) > maxAge {
		return fmt.Errorf("hosted PostgreSQL restore drill evidence is older than %s", maxAge)
	}
	return nil
}
