package hostedenv

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

const desiredStateSchema = "DE.PULSE-HOSTED-DESIRED-STATE-1"

var canonicalEnvironments = []string{"dev", "test", "stage", "prod"}

//go:embed desired_state_v1.json
var desiredStateJSON []byte

type DesiredState struct {
	IsolationID      string `json:"isolationId"`
	ServiceIdentity string `json:"serviceIdentity"`
	IngressPolicy   string `json:"ingressPolicy"`
	EgressPolicy    string `json:"egressPolicy"`
	NetworkPolicy   string `json:"networkPolicy"`
	TLSPolicy       string `json:"tlsPolicy"`
	InternalMTLS    bool   `json:"internalMTLS"`
}

type manifest struct {
	Schema       string                  `json:"schema"`
	Version      string                  `json:"version"`
	Environments map[string]DesiredState `json:"environments"`
}

type RuntimeDeclaration struct {
	Environment         string
	DesiredStateVersion string
	DesiredStateSHA256  string
	IsolationID         string
	ServiceIdentity     string
	IngressPolicy       string
	EgressPolicy        string
	NetworkPolicy       string
	TLSPolicy           string
	InternalMTLS        bool
	TrustedProxyHeaders bool
	PublicOrigin        string
}

var (
	manifestOnce sync.Once
	manifestData manifest
	manifestErr  error
)

func loadManifest() (manifest, error) {
	manifestOnce.Do(func() {
		if err := json.Unmarshal(desiredStateJSON, &manifestData); err != nil {
			manifestErr = fmt.Errorf("decode desired state: %w", err)
			return
		}
		manifestErr = validateManifest(manifestData)
	})
	return manifestData, manifestErr
}

func validateManifest(m manifest) error {
	if m.Schema != desiredStateSchema {
		return fmt.Errorf("unsupported desired-state schema %q", m.Schema)
	}
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("desired-state version is empty")
	}
	if len(m.Environments) != len(canonicalEnvironments) {
		return fmt.Errorf("desired state must define exactly dev/test/stage/prod")
	}
	isolationIDs := map[string]struct{}{}
	serviceIDs := map[string]struct{}{}
	for _, environment := range canonicalEnvironments {
		state, ok := m.Environments[environment]
		if !ok {
			return fmt.Errorf("desired state missing %s environment", environment)
		}
		if err := validateDesiredState(environment, state); err != nil {
			return err
		}
		if _, exists := isolationIDs[state.IsolationID]; exists {
			return fmt.Errorf("duplicate isolation id %q", state.IsolationID)
		}
		isolationIDs[state.IsolationID] = struct{}{}
		if _, exists := serviceIDs[state.ServiceIdentity]; exists {
			return fmt.Errorf("duplicate service identity %q", state.ServiceIdentity)
		}
		serviceIDs[state.ServiceIdentity] = struct{}{}
	}
	return nil
}

func validateDesiredState(environment string, state DesiredState) error {
	fields := map[string]string{
		"isolationId":      state.IsolationID,
		"serviceIdentity": state.ServiceIdentity,
		"ingressPolicy":   state.IngressPolicy,
		"egressPolicy":    state.EgressPolicy,
		"networkPolicy":   state.NetworkPolicy,
		"tlsPolicy":       state.TLSPolicy,
	}
	for name, value := range fields {
		value = strings.TrimSpace(value)
		if value == "" {
			return fmt.Errorf("%s desired state %s is empty", environment, name)
		}
		lower := strings.ToLower(value)
		if strings.Contains(value, "*") || strings.Contains(lower, "allow-all") || strings.Contains(lower, "allow_any") {
			return fmt.Errorf("%s desired state %s is not least privilege", environment, name)
		}
	}
	if !state.InternalMTLS {
		return fmt.Errorf("%s desired state must require internal mTLS", environment)
	}
	if !strings.Contains(strings.ToLower(state.IngressPolicy), "https:443") {
		return fmt.Errorf("%s ingress policy must expose HTTPS only", environment)
	}
	if !strings.Contains(strings.ToLower(state.EgressPolicy), "default-deny") {
		return fmt.Errorf("%s egress policy must default deny", environment)
	}
	if !strings.Contains(strings.ToLower(state.NetworkPolicy), "default-deny") {
		return fmt.Errorf("%s network policy must default deny", environment)
	}
	if !strings.Contains(strings.ToLower(state.TLSPolicy), "tls1.2") || !strings.Contains(strings.ToLower(state.TLSPolicy), "mtls") {
		return fmt.Errorf("%s TLS policy must require TLS 1.2+ and internal mTLS", environment)
	}
	return nil
}

func DesiredStateFor(environment string) (string, DesiredState, error) {
	m, err := loadManifest()
	if err != nil {
		return "", DesiredState{}, err
	}
	environment = strings.ToLower(strings.TrimSpace(environment))
	state, ok := m.Environments[environment]
	if !ok {
		return "", DesiredState{}, fmt.Errorf("unsupported hosted environment %q; expected dev, test, stage, or prod", environment)
	}
	return m.Version, state, nil
}

func ExpectedDigest(environment string) (string, error) {
	version, state, err := DesiredStateFor(environment)
	if err != nil {
		return "", err
	}
	canonical := struct {
		Schema      string       `json:"schema"`
		Version     string       `json:"version"`
		Environment string       `json:"environment"`
		State       DesiredState `json:"state"`
	}{desiredStateSchema, version, strings.ToLower(strings.TrimSpace(environment)), state}
	payload, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("encode desired state: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func Validate(declaration RuntimeDeclaration) error {
	environment := strings.ToLower(strings.TrimSpace(declaration.Environment))
	version, desired, err := DesiredStateFor(environment)
	if err != nil {
		return err
	}
	digest, err := ExpectedDigest(environment)
	if err != nil {
		return err
	}
	checks := []struct {
		name string
		got  string
		want string
	}{
		{"desired-state version", declaration.DesiredStateVersion, version},
		{"desired-state sha256", strings.ToLower(declaration.DesiredStateSHA256), digest},
		{"isolation id", declaration.IsolationID, desired.IsolationID},
		{"service identity", declaration.ServiceIdentity, desired.ServiceIdentity},
		{"ingress policy", declaration.IngressPolicy, desired.IngressPolicy},
		{"egress policy", declaration.EgressPolicy, desired.EgressPolicy},
		{"network policy", declaration.NetworkPolicy, desired.NetworkPolicy},
		{"TLS policy", declaration.TLSPolicy, desired.TLSPolicy},
	}
	for _, check := range checks {
		if strings.TrimSpace(check.got) != check.want {
			return fmt.Errorf("hosted desired-state drift: %s does not match %s", check.name, environment)
		}
	}
	if declaration.InternalMTLS != desired.InternalMTLS {
		return fmt.Errorf("hosted desired-state drift: internal mTLS does not match %s", environment)
	}
	if !declaration.TrustedProxyHeaders {
		return fmt.Errorf("hosted service trust requires the explicitly trusted managed proxy path")
	}
	if err := validatePublicOrigin(declaration.PublicOrigin); err != nil {
		return err
	}
	return nil
}

func validatePublicOrigin(origin string) error {
	origin = strings.TrimSpace(origin)
	u, err := url.Parse(origin)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" || strings.Contains(u.Host, "*") {
		return fmt.Errorf("hosted public origin must be one explicit HTTPS origin")
	}
	return nil
}
