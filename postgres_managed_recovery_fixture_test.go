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
