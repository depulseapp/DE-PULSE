package hostedenv

import (
	"strings"
	"testing"
)

func validDeclaration(t *testing.T, environment string) RuntimeDeclaration {
	t.Helper()
	version, desired, err := DesiredStateFor(environment)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := ExpectedDigest(environment)
	if err != nil {
		t.Fatal(err)
	}
	return RuntimeDeclaration{
		Environment:         environment,
		DesiredStateVersion: version,
		DesiredStateSHA256:  digest,
		IsolationID:         desired.IsolationID,
		ServiceIdentity:     desired.ServiceIdentity,
		IngressPolicy:       desired.IngressPolicy,
		EgressPolicy:        desired.EgressPolicy,
		NetworkPolicy:       desired.NetworkPolicy,
		TLSPolicy:           desired.TLSPolicy,
		InternalMTLS:        desired.InternalMTLS,
		TrustedProxyHeaders: true,
		PublicOrigin:        "https://" + environment + ".depulse.example",
	}
}

func TestDesiredStateDefinesFourIsolatedEnvironments(t *testing.T) {
	seenIsolation := map[string]bool{}
	seenService := map[string]bool{}
	for _, environment := range canonicalEnvironments {
		version, desired, err := DesiredStateFor(environment)
		if err != nil {
			t.Fatalf("%s: %v", environment, err)
		}
		if version != "v1" {
			t.Fatalf("%s version = %q", environment, version)
		}
		if seenIsolation[desired.IsolationID] || seenService[desired.ServiceIdentity] {
			t.Fatalf("%s is not isolated: %+v", environment, desired)
		}
		seenIsolation[desired.IsolationID] = true
		seenService[desired.ServiceIdentity] = true
		if !desired.InternalMTLS || !strings.Contains(desired.EgressPolicy, "default-deny") || !strings.Contains(desired.NetworkPolicy, "default-deny") {
			t.Fatalf("%s lacks least-privilege service trust: %+v", environment, desired)
		}
	}
	if _, _, err := DesiredStateFor("production"); err == nil {
		t.Fatal("non-canonical environment alias must fail closed")
	}
}

func TestValidateAcceptsExactDesiredState(t *testing.T) {
	for _, environment := range canonicalEnvironments {
		if err := Validate(validDeclaration(t, environment)); err != nil {
			t.Fatalf("%s exact declaration rejected: %v", environment, err)
		}
	}
}

func TestValidateRejectsEnvironmentAndServiceTrustDrift(t *testing.T) {
	base := validDeclaration(t, "stage")
	tests := []struct {
		name   string
		mutate func(*RuntimeDeclaration)
	}{
		{"version", func(d *RuntimeDeclaration) { d.DesiredStateVersion = "v0" }},
		{"digest", func(d *RuntimeDeclaration) { d.DesiredStateSHA256 = strings.Repeat("0", 64) }},
		{"isolation", func(d *RuntimeDeclaration) { d.IsolationID = "depulse-prod" }},
		{"service identity", func(d *RuntimeDeclaration) { d.ServiceIdentity = "depulse-web-prod" }},
		{"ingress", func(d *RuntimeDeclaration) { d.IngressPolicy = "allow-all:*" }},
		{"egress", func(d *RuntimeDeclaration) { d.EgressPolicy = "allow-all:*" }},
		{"network", func(d *RuntimeDeclaration) { d.NetworkPolicy = "allow-all" }},
		{"tls", func(d *RuntimeDeclaration) { d.TLSPolicy = "tls1.0" }},
		{"mtls", func(d *RuntimeDeclaration) { d.InternalMTLS = false }},
		{"proxy trust", func(d *RuntimeDeclaration) { d.TrustedProxyHeaders = false }},
		{"http origin", func(d *RuntimeDeclaration) { d.PublicOrigin = "http://stage.depulse.example" }},
		{"wildcard origin", func(d *RuntimeDeclaration) { d.PublicOrigin = "https://*.depulse.example" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			declaration := base
			tc.mutate(&declaration)
			if err := Validate(declaration); err == nil {
				t.Fatalf("%s drift was accepted", tc.name)
			}
		})
	}
}

func TestExpectedDigestIsEnvironmentBound(t *testing.T) {
	digests := map[string]bool{}
	for _, environment := range canonicalEnvironments {
		digest, err := ExpectedDigest(environment)
		if err != nil {
			t.Fatal(err)
		}
		if len(digest) != 64 || digests[digest] {
			t.Fatalf("invalid or duplicate %s digest %q", environment, digest)
		}
		digests[digest] = true
	}
}
