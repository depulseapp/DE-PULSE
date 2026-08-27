package main

import (
	"os"
	"strings"

	"depulse/internal/hostedenv"
)

const (
	runtimeModeEnv                      = "DEPULSE_RUNTIME_MODE"
	hostedListenAddrEnv                 = "DEPULSE_LISTEN_ADDR"
	hostedConfigDirEnv                  = "DEPULSE_CONFIG_DIR"
	hostedTrustProxyHeadersEnv          = "DEPULSE_TRUST_PROXY_HEADERS"
	hostedPublicOriginEnv               = "DEPULSE_PUBLIC_ORIGIN"
	hostedEnvironmentEnv                = "DEPULSE_HOSTED_ENVIRONMENT"
	hostedDesiredStateVersionEnv        = "DEPULSE_HOSTED_DESIRED_STATE_VERSION"
	hostedDesiredStateSHA256Env         = "DEPULSE_HOSTED_DESIRED_STATE_SHA256"
	hostedIsolationIDEnv                = "DEPULSE_HOSTED_ISOLATION_ID"
	hostedServiceIdentityEnv            = "DEPULSE_HOSTED_SERVICE_IDENTITY"
	hostedIngressPolicyEnv              = "DEPULSE_HOSTED_INGRESS_POLICY"
	hostedEgressPolicyEnv               = "DEPULSE_HOSTED_EGRESS_POLICY"
	hostedNetworkPolicyEnv              = "DEPULSE_HOSTED_NETWORK_POLICY"
	hostedTLSPolicyEnv                  = "DEPULSE_HOSTED_TLS_POLICY"
	hostedInternalMTLSEnv               = "DEPULSE_HOSTED_INTERNAL_MTLS"
	providerRightsEnforcementModeEnv    = "DEPULSE_PROVIDER_RIGHTS_ENFORCEMENT_MODE"
	providerRightsEnforcementPublicMode = "PUBLIC_PRODUCTION"
)

func runtimeMode() string {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(runtimeModeEnv))) {
	case "hosted", "server", "web":
		return "hosted"
	default:
		return "desktop"
	}
}

func isHostedRuntime() bool { return runtimeMode() == "hosted" }

// providerRightsEnforcementActive intentionally stays independent from hosted
// runtime selection. During development and pre-public validation, provider
// rights are evaluated and surfaced as governance/audit truth without reducing
// configured provider capacity. Hard fail-closed routing, fanout, cache,
// persistence and serving activate only when PUBLIC_PRODUCTION is explicitly
// selected for a hosted runtime.
func providerRightsEnforcementActive() bool {
	return isHostedRuntime() && strings.EqualFold(
		strings.TrimSpace(os.Getenv(providerRightsEnforcementModeEnv)),
		providerRightsEnforcementPublicMode,
	)
}

func hostedListenAddress() string {
	if addr := strings.TrimSpace(os.Getenv(hostedListenAddrEnv)); addr != "" {
		return addr
	}
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.Contains(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}

func trustHostedProxyHeaders() bool {
	if !isHostedRuntime() {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(hostedTrustProxyHeadersEnv))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func hostedPublicOrigin() string {
	if !isHostedRuntime() {
		return ""
	}
	return strings.TrimSpace(os.Getenv(hostedPublicOriginEnv))
}

func hostedInternalMTLSEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(hostedInternalMTLSEnv))) {
	case "1", "true", "yes", "on", "required":
		return true
	default:
		return false
	}
}

func validateHostedEnvironment() error {
	if !isHostedRuntime() {
		return nil
	}
	return hostedenv.Validate(hostedenv.RuntimeDeclaration{
		Environment:         strings.TrimSpace(os.Getenv(hostedEnvironmentEnv)),
		DesiredStateVersion: strings.TrimSpace(os.Getenv(hostedDesiredStateVersionEnv)),
		DesiredStateSHA256:  strings.TrimSpace(os.Getenv(hostedDesiredStateSHA256Env)),
		IsolationID:         strings.TrimSpace(os.Getenv(hostedIsolationIDEnv)),
		ServiceIdentity:     strings.TrimSpace(os.Getenv(hostedServiceIdentityEnv)),
		IngressPolicy:       strings.TrimSpace(os.Getenv(hostedIngressPolicyEnv)),
		EgressPolicy:        strings.TrimSpace(os.Getenv(hostedEgressPolicyEnv)),
		NetworkPolicy:       strings.TrimSpace(os.Getenv(hostedNetworkPolicyEnv)),
		TLSPolicy:           strings.TrimSpace(os.Getenv(hostedTLSPolicyEnv)),
		InternalMTLS:        hostedInternalMTLSEnabled(),
		TrustedProxyHeaders: trustHostedProxyHeaders(),
		PublicOrigin:        hostedPublicOrigin(),
	})
}
