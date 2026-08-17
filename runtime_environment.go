package main

import (
	"os"
	"strings"
)

const (
	runtimeModeEnv              = "DEPULSE_RUNTIME_MODE"
	hostedListenAddrEnv         = "DEPULSE_LISTEN_ADDR"
	hostedConfigDirEnv          = "DEPULSE_CONFIG_DIR"
	hostedTrustProxyHeadersEnv  = "DEPULSE_TRUST_PROXY_HEADERS"
	hostedPublicOriginEnv       = "DEPULSE_PUBLIC_ORIGIN"
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
