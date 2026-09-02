//go:build postgres

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	host012ManagedRecoveryFixtureAck    = "HOST012_MANAGED_RECOVERY_LIFECYCLE_FIXTURE"
	host012ManagedRecoveryFixtureSchema = "DE.PULSE-HOST012-MANAGED-RECOVERY-LIFECYCLE-FIXTURE-1"
)

type host012ManagedRecoveryFixtureArtifact struct {
	Schema                    string `json:"schema"`
	CandidateSHA              string `json:"candidateSha"`
	EnvironmentClass          string `json:"environmentClass"`
	Phase                     string `json:"phase"`
	DatabaseIdentity          string `json:"databaseIdentity"`
	UserID                    string `json:"userId"`
	UserIDSHA256              string `json:"userIdSha256"`
	TenantID                  string `json:"tenantId"`
	AccountStatus             string `json:"accountStatus"`
	TombstoneReason           string `json:"tombstoneReason,omitempty"`
	TombstoneDeletedAt        int64  `json:"tombstoneDeletedAt,omitempty"`
	DeviceCount               int    `json:"deviceCount"`
	SessionCount              int    `json:"sessionCount"`
	PersonalWorkspacePresent  bool   `json:"personalWorkspacePresent"`
	CanonicalLifecycleOwner   string `json:"canonicalLifecycleOwner"`
	ProviderBackupPITRClaimed bool   `json:"providerBackupPitrClaimed"`
	CreatedAt                 string `json:"createdAt"`
}

func host012ManagedRecoveryFixtureEnvironment(t *testing.T) (string, string, string, string) {
	t.Helper()
	if got := host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_FIXTURE_ACK"); got != host012ManagedRecoveryFixtureAck {
		t.Fatalf("DEPULSE_MANAGED_RECOVERY_FIXTURE_ACK must equal %q", host012ManagedRecoveryFixtureAck)
	}
	environmentClass := strings.ToLower(host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_ENV_CLASS"))
	switch environmentClass {
	case "development", "test", "stage":
	case "production":
		t.Fatal("managed recovery lifecycle fixture refuses production")
	default:
		t.Fatalf("unsupported DEPULSE_MANAGED_RECOVERY_ENV_CLASS %q", environmentClass)
	}
	candidateSHA := strings.ToLower(host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA"))
	if len(candidateSHA) != 40 {
		t.Fatalf("DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA must be a full 40-character git SHA")
	}
	for _, r := range candidateSHA {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			t.Fatalf("DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA must be hexadecimal")
		}
	}
	phase := strings.ToLower(host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_FIXTURE_PHASE"))
	if phase != "seed" && phase != "delete" {
		t.Fatalf("DEPULSE_MANAGED_RECOVERY_FIXTURE_PHASE must be seed or delete")
	}
	artifactPath := host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_FIXTURE_ARTIFACT_PATH")
	return environmentClass, candidateSHA, phase, artifactPath
}

func host012ManagedRecoveryFixtureManager(t *testing.T, databaseURL string) *PersistenceManager {
	t.Helper()
	backend := newPostgresPersistenceBackend(postgresPersistenceConfig{
		DatabaseURL:     databaseURL,
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: time.Minute,
	})
	manager := newPersistenceManagerWithBackend(backend)
	if diag := manager.Diagnostics(); !diag.Ready {
		_ = manager.Close()
		t.Fatalf("managed recovery lifecycle fixture persistence unavailable: %s", diag.LastError)
	}
	return manager
}

func host012ManagedRecoveryFixtureOwner(identity *IdentityService) (Principal, error) {
	if identity == nil {
		return Principal{}, errors.New("identity unavailable")
	}
	identity.mu.Lock()
	defer identity.mu.Unlock()
	for _, user := range identity.state.Users {
		if user.Status != UserActive || isAccountDeletionTombstone(user) {
			continue
		}
		if user.Role != RoleOwner && user.Role != RoleSuperOwner {
			continue
		}
		return Principal{
			TenantID:    normalizedTenantID(user.TenantID),
			UserID:      user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
		}, nil
	}
	return Principal{}, errors.New("an active owner or super owner is required to create the recovery fixture")
}

func host012ManagedRecoveryFixtureApplication(t *testing.T, manager *PersistenceManager, identity *IdentityService) *Application {
	t.Helper()
	loaded, err := manager.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("load managed recovery fixture workspaces: %v", err)
	}
	workspaces := make(map[string]UserWorkspace, len(loaded))
	for _, workspace := range loaded {
		workspace = normalizeUserWorkspace(workspace)
		if strings.TrimSpace(workspace.UserID) != "" {
			workspaces[workspace.UserID] = workspace
		}
	}
	return &Application{
		state:       defaultEmptyWorkspaceState(),
		persistence: manager,
		identity:    identity,
		workspaces:  workspaces,
	}
}

func host012WriteManagedRecoveryFixtureArtifact(t *testing.T, path string, artifact host012ManagedRecoveryFixtureArtifact) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create managed recovery fixture artifact directory: %v", err)
	}
	payload, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		t.Fatalf("marshal managed recovery fixture artifact: %v", err)
	}
	payload = append(payload, '\n')
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write managed recovery fixture artifact: %v", err)
	}
}

func host012ManagedRecoveryFixtureSeed(t *testing.T, manager *PersistenceManager, identity *IdentityService, candidateSHA, environmentClass, databaseIdentity, artifactPath string) {
	t.Helper()
	owner, err := host012ManagedRecoveryFixtureOwner(identity)
	if err != nil {
		t.Fatal(err)
	}
	suffix := time.Now().UTC().UnixNano()
	username := fmt.Sprintf("host012-%s-%d", candidateSHA[:8], suffix)
	temporaryPassword := fmt.Sprintf("Tmp#A1%s%d", candidateSHA[:12], suffix)
	finalPassword := fmt.Sprintf("Fixture#B2%s%d", candidateSHA[12:24], suffix)
	created, err := identity.adminCreateUser(owner, username, "HOST-012 Recovery Fixture", RoleUser, temporaryPassword)
	if err != nil {
		t.Fatalf("create managed recovery fixture user: %v", err)
	}
	_, principal, err := identity.setPassword(created.ID, finalPassword)
	if err != nil {
		t.Fatalf("credential managed recovery fixture user: %v", err)
	}
	fingerprint := host012UserHash(fmt.Sprintf("host012-device:%s:%d", created.ID, suffix))
	device, err := identity.registerHostedDevice(principal, "HOST-012 Recovery Fixture", fingerprint)
	if err != nil {
		t.Fatalf("register managed recovery fixture device: %v", err)
	}
	if err := identity.bindHostedDeviceToSession(principal, device.ID); err != nil {
		t.Fatalf("bind managed recovery fixture device: %v", err)
	}

	workspaceState := defaultEmptyWorkspaceState()
	setTrackedSymbolLocked(&workspaceState, "NVDA", true)
	setTrackedSymbolLocked(&workspaceState, "AMD", true)
	workspaceState.UI = UIState{ScopeType: "watchlist", WatchlistID: "swing", SelectedTicker: "NVDA"}
	workspace := workspaceFromState(created.ID, workspaceState, time.Now().UTC())
	if err := manager.SaveUserWorkspace(context.Background(), workspace); err != nil {
		t.Fatalf("save managed recovery fixture workspace: %v", err)
	}

	profile, hosted, err := identity.accountPrivacyExport(created.ID)
	if err != nil {
		t.Fatalf("export seeded recovery fixture account: %v", err)
	}
	if profile.Status != UserActive || profile.ID != created.ID {
		t.Fatalf("seeded recovery fixture profile is not active")
	}
	if len(hosted.Devices) == 0 || len(hosted.Sessions) == 0 {
		t.Fatalf("seeded recovery fixture must retain at least one device and session")
	}
	persistedWorkspaces, err := manager.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("load seeded recovery fixture workspaces: %v", err)
	}
	persistedWorkspace, ok := host012WorkspaceForUser(persistedWorkspaces, created.ID)
	if !ok || !host012WorkspaceContainsPersonalData(persistedWorkspace) {
		t.Fatalf("seeded recovery fixture must contain personal workspace state")
	}

	host012WriteManagedRecoveryFixtureArtifact(t, artifactPath, host012ManagedRecoveryFixtureArtifact{
		Schema:                    host012ManagedRecoveryFixtureSchema,
		CandidateSHA:              candidateSHA,
		EnvironmentClass:          environmentClass,
		Phase:                     "seed",
		DatabaseIdentity:          databaseIdentity,
		UserID:                    created.ID,
		UserIDSHA256:              host012UserHash(created.ID),
		TenantID:                  normalizedTenantID(profile.TenantID),
		AccountStatus:             string(UserActive),
		DeviceCount:               len(hosted.Devices),
		SessionCount:              len(hosted.Sessions),
		PersonalWorkspacePresent:  true,
		CanonicalLifecycleOwner:   "IdentityService.adminCreateUser+setPassword+registerHostedDevice+bindHostedDeviceToSession;PersistenceManager.SaveUserWorkspace",
		ProviderBackupPITRClaimed: false,
		CreatedAt:                 time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func host012ManagedRecoveryFixtureDelete(t *testing.T, manager *PersistenceManager, identity *IdentityService, candidateSHA, environmentClass, databaseIdentity, artifactPath string) {
	t.Helper()
	userID := host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_USER_ID")
	before, err := manager.LoadIdentityState(context.Background())
	if err != nil {
		t.Fatalf("load managed recovery fixture source identity: %v", err)
	}
	beforeUser, ok := host012FindUser(before, userID)
	if !ok || beforeUser.Status != UserActive || isAccountDeletionTombstone(beforeUser) {
		t.Fatalf("delete phase requires the active seeded recovery fixture account")
	}
	beforeWorkspaces, err := manager.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("load managed recovery fixture source workspaces: %v", err)
	}
	beforeWorkspace, ok := host012WorkspaceForUser(beforeWorkspaces, userID)
	if !ok || !host012WorkspaceContainsPersonalData(beforeWorkspace) {
		t.Fatalf("delete phase requires seeded personal workspace state")
	}

	app := host012ManagedRecoveryFixtureApplication(t, manager, identity)
	tombstone, err := app.deleteAccountData(userID, accountDeletionReasonUserRequest)
	if err != nil {
		t.Fatalf("delete managed recovery fixture account: %v", err)
	}
	if tombstone.UserID != userID || tombstone.Reason != accountDeletionReasonUserRequest || tombstone.DeletedAt <= 0 {
		t.Fatalf("managed recovery fixture deletion did not produce the canonical tombstone")
	}

	after, err := manager.LoadIdentityState(context.Background())
	if err != nil {
		t.Fatalf("reload deleted managed recovery fixture identity: %v", err)
	}
	deletedUser, ok := host012FindUser(after, userID)
	if !ok || !isAccountDeletionTombstone(deletedUser) {
		t.Fatalf("managed recovery fixture source is missing the canonical deletion tombstone")
	}
	for _, device := range after.Devices {
		if device.UserID == userID {
			t.Fatalf("managed recovery fixture deletion retained a device")
		}
	}
	for _, session := range after.Sessions {
		if session.UserID == userID {
			t.Fatalf("managed recovery fixture deletion retained a session")
		}
	}
	afterWorkspaces, err := manager.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("reload deleted managed recovery fixture workspaces: %v", err)
	}
	afterWorkspace, ok := host012WorkspaceForUser(afterWorkspaces, userID)
	if !ok || host012WorkspaceContainsPersonalData(afterWorkspace) {
		t.Fatalf("managed recovery fixture deletion did not leave a canonical privacy-blank workspace")
	}
	if _, _, err := identity.accountPrivacyExport(userID); err == nil {
		t.Fatalf("deleted managed recovery fixture still exposes account privacy data")
	}

	host012WriteManagedRecoveryFixtureArtifact(t, artifactPath, host012ManagedRecoveryFixtureArtifact{
		Schema:                    host012ManagedRecoveryFixtureSchema,
		CandidateSHA:              candidateSHA,
		EnvironmentClass:          environmentClass,
		Phase:                     "delete",
		DatabaseIdentity:          databaseIdentity,
		UserID:                    userID,
		UserIDSHA256:              host012UserHash(userID),
		TenantID:                  normalizedTenantID(tombstone.TenantID),
		AccountStatus:             string(UserDisabled),
		TombstoneReason:           tombstone.Reason,
		TombstoneDeletedAt:        tombstone.DeletedAt,
		DeviceCount:               0,
		SessionCount:              0,
		PersonalWorkspacePresent:  false,
		CanonicalLifecycleOwner:   "Application.deleteAccountData->IdentityService.deleteAccount+privacyBlankWorkspace",
		ProviderBackupPITRClaimed: false,
		CreatedAt:                 time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func TestHOST012ManagedRecoveryLifecycleFixture(t *testing.T) {
	environmentClass, candidateSHA, phase, artifactPath := host012ManagedRecoveryFixtureEnvironment(t)
	sourceURL := host012ManagedRecoveryDatabaseURL(t, "DEPULSE_MANAGED_RECOVERY_SOURCE_URL")
	databaseIdentity := host012ManagedRecoveryDatabaseIdentity(t, sourceURL)
	manager := host012ManagedRecoveryFixtureManager(t, sourceURL)
	defer func() {
		if err := manager.Close(); err != nil {
			t.Errorf("close managed recovery fixture persistence: %v", err)
		}
	}()
	identity, err := NewIdentityService(manager)
	if err != nil {
		t.Fatalf("initialize managed recovery fixture identity: %v", err)
	}

	switch phase {
	case "seed":
		host012ManagedRecoveryFixtureSeed(t, manager, identity, candidateSHA, environmentClass, databaseIdentity, artifactPath)
	case "delete":
		host012ManagedRecoveryFixtureDelete(t, manager, identity, candidateSHA, environmentClass, databaseIdentity, artifactPath)
	default:
		t.Fatalf("unsupported managed recovery fixture phase %q", phase)
	}
}

func host015PostgresConfig(t *testing.T) postgresPersistenceConfig {
	t.Helper()
	return postgresPersistenceConfig{
		DatabaseURL:     v183PostgresURL(t),
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}
}

func host015PreparePostgresTenancy(t *testing.T) postgresPersistenceConfig {
	t.Helper()
	config := host015PostgresConfig(t)
	raw := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := raw.Init(ctx); err != nil {
		t.Fatalf("prepare HOST-015 postgres: %v", err)
	}
	defer raw.Close()
	if _, err := raw.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_user_workspaces`); err != nil {
		t.Fatal(err)
	}
	if _, err := raw.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_identity_state`); err != nil {
		t.Fatal(err)
	}
	if _, err := raw.db.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version=$1`, hostedTenantPostgresSchemaVersion); err != nil {
		t.Fatal(err)
	}
	if _, err := raw.db.ExecContext(ctx, `TRUNCATE TABLE user_workspaces, identity_state`); err != nil {
		t.Fatal(err)
	}
	return config
}

func host015OpenPostgresTenancy(t *testing.T, config postgresPersistenceConfig) *hostedTenantPostgresBackend {
	t.Helper()
	t.Setenv(runtimeModeEnv, "hosted")
	raw := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	wrapped, ok := wrapHostedTenantPostgresBackend(raw).(*hostedTenantPostgresBackend)
	if !ok {
		t.Fatalf("hosted postgres tenant wrapper not selected: %T", wrapHostedTenantPostgresBackend(raw))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := wrapped.Init(ctx); err != nil {
		t.Fatalf("open HOST-015 postgres tenancy backend: %v", err)
	}
	return wrapped
}

func host015CleanupPostgresTenancy(t *testing.T, backend *hostedTenantPostgresBackend) {
	t.Helper()
	if backend == nil || backend.pg == nil || backend.pg.db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = backend.pg.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_user_workspaces`)
	_, _ = backend.pg.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_identity_state`)
	_, _ = backend.pg.db.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version=$1`, hostedTenantPostgresSchemaVersion)
	_, _ = backend.pg.db.ExecContext(ctx, `TRUNCATE TABLE user_workspaces, identity_state`)
	_ = backend.Close()
}

func TestHOST015PostgresTenantPersistenceSurvivesRestart(t *testing.T) {
	config := host015PreparePostgresTenancy(t)
	backend := host015OpenPostgresTenancy(t, config)
	ctx := context.Background()
	state := host015TenantFixture()
	if err := backend.SaveIdentityState(ctx, state); err != nil {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal(err)
	}
	for _, userID := range []string{"user-a", "user-b"} {
		workspace := defaultUserWorkspace(userID)
		workspace.UpdatedAt = 10
		if err := backend.SaveUserWorkspace(ctx, workspace); err != nil {
			host015CleanupPostgresTenancy(t, backend)
			t.Fatal(err)
		}
	}
	stats, err := backend.Stats(ctx)
	if err != nil || stats.SchemaVersion != hostedTenantPostgresSchemaVersion || stats.UserCount != 2 || stats.SessionCount != 2 {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatalf("tenant postgres stats mismatch: stats=%+v err=%v", stats, err)
	}
	if err := backend.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := host015OpenPostgresTenancy(t, config)
	defer host015CleanupPostgresTenancy(t, reopened)
	loaded, err := reopened.LoadIdentityState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	index, err := hostedTenantUserIndex(loaded)
	if err != nil {
		t.Fatal(err)
	}
	if index["user-a"] != "tenant-a" || index["user-b"] != "tenant-b" {
		t.Fatalf("tenant identity ownership did not survive restart: %+v", index)
	}
	workspaces, err := reopened.LoadUserWorkspaces(ctx)
	if err != nil || len(workspaces) != 2 {
		t.Fatalf("tenant workspaces did not survive restart: workspaces=%+v err=%v", workspaces, err)
	}
	var legacyIdentity, legacyWorkspaces int
	if err := reopened.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_state`).Scan(&legacyIdentity); err != nil {
		t.Fatal(err)
	}
	if err := reopened.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user_workspaces`).Scan(&legacyWorkspaces); err != nil {
		t.Fatal(err)
	}
	if legacyIdentity != 0 || legacyWorkspaces != 0 {
		t.Fatalf("legacy cross-tenant authorities remain: identity=%d workspaces=%d", legacyIdentity, legacyWorkspaces)
	}
}

func TestHOST015PostgresRejectsCrossTenantWorkspaceRow(t *testing.T) {
	config := host015PreparePostgresTenancy(t)
	backend := host015OpenPostgresTenancy(t, config)
	defer host015CleanupPostgresTenancy(t, backend)
	ctx := context.Background()
	if err := backend.SaveIdentityState(ctx, host015TenantFixture()); err != nil {
		t.Fatal(err)
	}
	workspace := defaultUserWorkspace("user-a")
	workspace.UpdatedAt = 10
	if err := backend.SaveUserWorkspace(ctx, workspace); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.pg.db.ExecContext(ctx, `UPDATE tenant_user_workspaces SET tenant_id='tenant-b' WHERE user_id='user-a'`); err != nil {
		t.Fatal(err)
	}
	if err := backend.SaveUserWorkspace(ctx, workspace); err == nil || !strings.Contains(err.Error(), "conflicting tenant") {
		t.Fatalf("tampered cross-tenant workspace was silently reassigned: %v", err)
	}
	var storedTenant string
	if err := backend.pg.db.QueryRowContext(ctx, `SELECT tenant_id FROM tenant_user_workspaces WHERE user_id='user-a'`).Scan(&storedTenant); err != nil {
		t.Fatal(err)
	}
	if storedTenant != "tenant-b" {
		t.Fatalf("failed save mutated tampered ownership row: tenant=%q", storedTenant)
	}
	if _, err := backend.LoadUserWorkspaces(ctx); err == nil || !strings.Contains(err.Error(), "crosses tenant boundary") {
		t.Fatalf("tampered cross-tenant workspace unexpectedly loaded: %v", err)
	}
}

func TestHOST015PostgresExpandMigratesLegacyIdentityAndWorkspace(t *testing.T) {
	config := host015PreparePostgresTenancy(t)
	ctx := context.Background()
	raw := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	if err := raw.Init(ctx); err != nil {
		t.Fatal(err)
	}
	legacy := IdentityPersistentState{
		Version:   1,
		Users:     []UserRecord{{ID: "legacy-user", Username: "legacy", Role: RoleOwner, Status: UserActive, CreatedAt: 1, UpdatedAt: 2}},
		Sessions:  []SessionRecord{{ID: "legacy-session", UserID: "legacy-user", CreatedAt: 1, LastSeenAt: 2, IdleExpiresAt: 100, AbsoluteExpiresAt: 200}},
		UpdatedAt: 2,
	}
	if err := raw.SaveIdentityState(ctx, legacy); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	workspace := defaultUserWorkspace("legacy-user")
	workspace.UpdatedAt = 2
	if err := raw.SaveUserWorkspace(ctx, workspace); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	if err := raw.Close(); err != nil {
		t.Fatal(err)
	}

	backend := host015OpenPostgresTenancy(t, config)
	defer host015CleanupPostgresTenancy(t, backend)
	loaded, err := backend.LoadIdentityState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	index, err := hostedTenantUserIndex(loaded)
	if err != nil {
		t.Fatal(err)
	}
	if index["legacy-user"] != localTenantID {
		t.Fatalf("legacy user was not deterministically bound to local tenant: %+v", index)
	}
	workspaces, err := backend.LoadUserWorkspaces(ctx)
	if err != nil || len(workspaces) != 1 || workspaces[0].UserID != "legacy-user" {
		t.Fatalf("legacy workspace migration failed: workspaces=%+v err=%v", workspaces, err)
	}
	var legacyIdentity, legacyWorkspaces int
	_ = backend.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_state`).Scan(&legacyIdentity)
	_ = backend.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user_workspaces`).Scan(&legacyWorkspaces)
	if legacyIdentity != 0 || legacyWorkspaces != 0 {
		t.Fatalf("legacy authority survived expand migration: identity=%d workspaces=%d", legacyIdentity, legacyWorkspaces)
	}
}

func TestHOST016PostgresTenantArchiveRecoveryAndRollback(t *testing.T) {
	config := host015PreparePostgresTenancy(t)
	backend := host015OpenPostgresTenancy(t, config)
	ctx := context.Background()
	state := host015TenantFixture()
	if err := backend.SaveIdentityState(ctx, state); err != nil {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal(err)
	}
	for _, userID := range []string{"user-a", "user-b"} {
		workspace := defaultUserWorkspace(userID)
		if userID == "user-b" {
			for i := range workspace.Watchlists {
				if workspace.Watchlists[i].ID == "swing" {
					workspace.Watchlists[i].Symbols = []string{"TSLA"}
				}
			}
			workspace.UI.SelectedTicker = "TSLA"
		}
		workspace.UpdatedAt = 10
		if err := backend.SaveUserWorkspace(ctx, workspace); err != nil {
			host015CleanupPostgresTenancy(t, backend)
			t.Fatal(err)
		}
	}

	archive, err := backend.ExportPersistenceArchive(ctx)
	if err != nil {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal(err)
	}
	archive.SchemaVersion = persistenceArchiveSchemaVersion
	if !archive.HasIdentity || len(archive.Identity.Tenants) != 2 || len(archive.UserWorkspaces) != 2 {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatalf("tenant archive omitted canonical state: %+v", archive)
	}

	if workspace, found := host012WorkspaceForUser(archive.UserWorkspaces, "user-b"); !found || !host012WorkspaceContainsPersonalData(workspace) {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal("recovery fixture did not contain personal workspace data before deletion replay")
	}
	current := archive.Identity
	current.Users = append([]UserRecord(nil), archive.Identity.Users...)
	for i := range current.Users {
		if current.Users[i].ID == "user-b" {
			current.Users[i] = UserRecord{
				ID:          "user-b",
				TenantID:    "tenant-b",
				Username:    accountDeletionTombstoneUsername,
				DisplayName: accountDeletionReasonUserRequest,
				Status:      UserDisabled,
				UpdatedAt:   50,
			}
		}
	}
	current.Devices = nil
	for _, device := range archive.Identity.Devices {
		if device.UserID != "user-b" {
			current.Devices = append(current.Devices, device)
		}
	}
	current.Sessions = nil
	for _, session := range archive.Identity.Sessions {
		if session.UserID != "user-b" {
			current.Sessions = append(current.Sessions, session)
		}
	}
	protected := enforceArchiveAccountDeletionPrivacy(archive, current)
	if _, present := accountDeletionTombstones(protected.Identity)["user-b"]; !present {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal("tenant recovery archive did not retain the authoritative tombstone")
	}
	if workspace, found := host012WorkspaceForUser(protected.UserWorkspaces, "user-b"); !found || host012WorkspaceContainsPersonalData(workspace) {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal("tenant recovery archive did not privacy-blank the deleted workspace")
	}

	tampered := protected
	tampered.Identity.Sessions = append([]SessionRecord(nil), protected.Identity.Sessions...)
	tampered.Identity.Sessions[0].TenantID = "tenant-b"
	if err := backend.RestorePersistenceArchive(ctx, tampered, persistenceRestoreModeReplace); err == nil || !strings.Contains(err.Error(), "crosses tenant boundary") {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatalf("tampered tenant recovery archive did not fail closed: %v", err)
	}
	unchanged, err := backend.LoadIdentityState(ctx)
	if err != nil {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal(err)
	}
	if user, found := host012FindUser(unchanged, "user-b"); !found || user.Status != UserActive {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatal("failed recovery mutated the live tenant state")
	}

	if err := backend.RestorePersistenceArchive(ctx, protected, persistenceRestoreModeReplace); err != nil {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatalf("restore tenant recovery archive: %v", err)
	}
	if err := backend.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := host015OpenPostgresTenancy(t, config)
	defer host015CleanupPostgresTenancy(t, reopened)
	recovered, err := reopened.LoadIdentityState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	index, err := hostedTenantUserIndex(recovered)
	if err != nil {
		t.Fatal(err)
	}
	if index["user-a"] != "tenant-a" || index["user-b"] != "tenant-b" {
		t.Fatalf("recovery changed canonical tenant ownership: %+v", index)
	}
	if _, present := accountDeletionTombstones(recovered)["user-b"]; !present {
		t.Fatal("recovery lost the tenant deletion tombstone")
	}
	for _, device := range recovered.Devices {
		if device.UserID == "user-b" {
			t.Fatal("recovery resurrected a deleted tenant device")
		}
	}
	for _, session := range recovered.Sessions {
		if session.UserID == "user-b" {
			t.Fatal("recovery resurrected a deleted tenant session")
		}
	}
	workspaces, err := reopened.LoadUserWorkspaces(ctx)
	if err != nil || len(workspaces) != 2 {
		t.Fatalf("recovered tenant workspaces invalid: workspaces=%+v err=%v", workspaces, err)
	}
	if workspace, found := host012WorkspaceForUser(workspaces, "user-b"); !found || host012WorkspaceContainsPersonalData(workspace) {
		t.Fatal("recovery resurrected deleted tenant workspace data")
	}
	projection := &IdentityService{state: recovered}
	if _, _, err := projection.accountPrivacyExport("user-b"); err == nil {
		t.Fatal("recovered deleted tenant remained visible through account projection")
	}
	var tenantIdentityRows, legacyIdentityRows int
	if err := reopened.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenant_identity_state`).Scan(&tenantIdentityRows); err != nil {
		t.Fatal(err)
	}
	if err := reopened.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_state`).Scan(&legacyIdentityRows); err != nil {
		t.Fatal(err)
	}
	if tenantIdentityRows != 2 || legacyIdentityRows != 0 {
		t.Fatalf("recovery authority mismatch: tenant=%d legacy=%d", tenantIdentityRows, legacyIdentityRows)
	}
}

func TestHOST016PostgresLegacyMigrationFailureRollsBackAuthority(t *testing.T) {
	config := host015PreparePostgresTenancy(t)
	ctx := context.Background()
	raw := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	if err := raw.Init(ctx); err != nil {
		t.Fatal(err)
	}
	legacy := IdentityPersistentState{
		Version:   1,
		Users:     []UserRecord{{ID: "legacy-user", Username: "legacy", Role: RoleOwner, Status: UserActive, CreatedAt: 1, UpdatedAt: 2}},
		Sessions:  []SessionRecord{{ID: "legacy-session", UserID: "legacy-user", CreatedAt: 1, LastSeenAt: 2, IdleExpiresAt: 100, AbsoluteExpiresAt: 200}},
		UpdatedAt: 2,
	}
	if err := raw.SaveIdentityState(ctx, legacy); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	validWorkspace := defaultUserWorkspace("legacy-user")
	validWorkspace.UpdatedAt = 2
	if err := raw.SaveUserWorkspace(ctx, validWorkspace); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	orphanWorkspace := defaultUserWorkspace("orphan-user")
	orphanWorkspace.UpdatedAt = 2
	if err := raw.SaveUserWorkspace(ctx, orphanWorkspace); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	if err := raw.Close(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(runtimeModeEnv, "hosted")
	migrationRaw := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	backend := wrapHostedTenantPostgresBackend(migrationRaw).(*hostedTenantPostgresBackend)
	if err := backend.Init(ctx); err == nil || !strings.Contains(err.Error(), "no canonical tenant owner") {
		host015CleanupPostgresTenancy(t, backend)
		t.Fatalf("invalid legacy ownership did not fail migration: %v", err)
	}

	inspector := newPostgresPersistenceBackend(config).(*postgresPersistenceBackend)
	if err := inspector.Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = inspector.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_user_workspaces`)
		_, _ = inspector.db.ExecContext(ctx, `DROP TABLE IF EXISTS tenant_identity_state`)
		_, _ = inspector.db.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version=$1`, hostedTenantPostgresSchemaVersion)
		_, _ = inspector.db.ExecContext(ctx, `TRUNCATE TABLE user_workspaces, identity_state`)
		_ = inspector.Close()
	}()
	var legacyIdentityRows, legacyWorkspaceRows, tenantIdentityRows, tenantWorkspaceRows int
	if err := inspector.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_state`).Scan(&legacyIdentityRows); err != nil {
		t.Fatal(err)
	}
	if err := inspector.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user_workspaces`).Scan(&legacyWorkspaceRows); err != nil {
		t.Fatal(err)
	}
	if err := inspector.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenant_identity_state`).Scan(&tenantIdentityRows); err != nil {
		t.Fatal(err)
	}
	if err := inspector.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenant_user_workspaces`).Scan(&tenantWorkspaceRows); err != nil {
		t.Fatal(err)
	}
	if legacyIdentityRows != 1 || legacyWorkspaceRows != 2 || tenantIdentityRows != 0 || tenantWorkspaceRows != 0 {
		t.Fatalf("failed migration partially contracted authority: legacyIdentity=%d legacyWorkspaces=%d tenantIdentity=%d tenantWorkspaces=%d", legacyIdentityRows, legacyWorkspaceRows, tenantIdentityRows, tenantWorkspaceRows)
	}
}
