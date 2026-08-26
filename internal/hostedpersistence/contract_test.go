package hostedpersistence

import (
	"strings"
	"testing"
	"time"
)

func validProdDeclaration(now time.Time) RuntimeDeclaration {
	return RuntimeDeclaration{
		Environment:        "prod",
		DatabaseURL:        "postgresql://svc:super-secret@db.example.invalid/depulse?sslmode=verify-full&search_path=depulse_prod",
		MaxOpenConnections: 32,
		PolicyVersion:      "v1",
		RLSDisposition:     "application-authorization-canonical-identity",
		MigrationStrategy:  "expand-contract",
		HAReady:            true,
		BackupEncrypted:    true,
		PITRReady:          true,
		RPO:                "10m",
		RTO:                "45m",
		RestoreDrillAt:     now.Add(-24 * time.Hour).UTC().Format(time.RFC3339),
		RestoreDrillID:     "restore-drill-prod-2026-08",
		RollbackPolicy:     "forward-fix-or-tested-restore",
		RollbackVerified:   true,
	}
}

func TestHostedPostgresPolicyUsesUniqueEnvironmentSchemas(t *testing.T) {
	seen := map[string]string{}
	for _, environment := range canonicalEnvironments {
		version, rls, migration, rollback, policy, err := PolicyFor(environment)
		if err != nil {
			t.Fatal(err)
		}
		if version != "v1" || rls != "application-authorization-canonical-identity" || migration != "expand-contract" || rollback != "forward-fix-or-tested-restore" {
			t.Fatalf("unexpected canonical policy for %s", environment)
		}
		if prior, ok := seen[policy.Schema]; ok {
			t.Fatalf("environment schemas are not isolated: %s and %s share %s", prior, environment, policy.Schema)
		}
		seen[policy.Schema] = environment
	}
}

func TestHostedPostgresProductionDeclarationPasses(t *testing.T) {
	now := time.Date(2026, time.August, 26, 20, 30, 0, 0, time.UTC)
	if err := Validate(validProdDeclaration(now), now); err != nil {
		t.Fatalf("valid production declaration rejected: %v", err)
	}
}

func TestHostedPostgresPolicyFailsClosedOnTrustTenancyAndRecoveryDrift(t *testing.T) {
	now := time.Date(2026, time.August, 26, 20, 30, 0, 0, time.UTC)
	tests := map[string]func(*RuntimeDeclaration){
		"environment schema": func(d *RuntimeDeclaration) {
			d.DatabaseURL = "postgresql://svc:super-secret@db.example.invalid/depulse?sslmode=verify-full&search_path=depulse_stage"
		},
		"tls": func(d *RuntimeDeclaration) {
			d.DatabaseURL = "postgresql://svc:super-secret@db.example.invalid/depulse?sslmode=require&search_path=depulse_prod"
		},
		"rls disposition":   func(d *RuntimeDeclaration) { d.RLSDisposition = "rls-enabled" },
		"migration":         func(d *RuntimeDeclaration) { d.MigrationStrategy = "destructive" },
		"pool capacity":     func(d *RuntimeDeclaration) { d.MaxOpenConnections = 65 },
		"ha":                func(d *RuntimeDeclaration) { d.HAReady = false },
		"backup encryption": func(d *RuntimeDeclaration) { d.BackupEncrypted = false },
		"pitr":              func(d *RuntimeDeclaration) { d.PITRReady = false },
		"rpo":               func(d *RuntimeDeclaration) { d.RPO = "16m" },
		"rto":               func(d *RuntimeDeclaration) { d.RTO = "61m" },
		"restore drill":     func(d *RuntimeDeclaration) { d.RestoreDrillAt = now.Add(-31 * 24 * time.Hour).Format(time.RFC3339) },
		"rollback policy":   func(d *RuntimeDeclaration) { d.RollbackPolicy = "hope" },
		"rollback evidence": func(d *RuntimeDeclaration) { d.RollbackVerified = false },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			declaration := validProdDeclaration(now)
			mutate(&declaration)
			err := Validate(declaration, now)
			if err == nil {
				t.Fatal("drift declaration unexpectedly passed")
			}
			if strings.Contains(err.Error(), "super-secret") {
				t.Fatalf("validation leaked database credentials: %v", err)
			}
		})
	}
}
