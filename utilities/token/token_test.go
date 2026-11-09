package token

import (
	"testing"
	"time"
)

func TestDetectTypeRoundTripsServiceTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		issuer   string
		clientID string
		audience string
		subject  string
	}{
		{name: "client", issuer: "aegis", clientID: "sender", audience: "chaos"},
		{name: "plain subject", issuer: "aegis", clientID: "sender", audience: "chaos", subject: "sender"},
		{name: "subject without cli", issuer: "aegis", audience: "chaos", subject: "sender"},
		{name: "self audience", issuer: "sender", clientID: "sender", audience: "sender"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := NewClaimsBuilder().Issuer(tt.issuer).ClientID(tt.clientID).
				Audience(tt.audience).Subject(tt.subject).ExpiresIn(time.Minute).
				Build(NewServiceAccessTokenBuilder().Scope("mail:send"))
			wire, err := Build(input)
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			kind := DetectType(wire)
			if kind != TokenTypeSAT {
				t.Fatalf("DetectType() = %q, want %q", kind, TokenTypeSAT)
			}
			parsed, err := ParseToken(wire, kind)
			if err != nil {
				t.Fatalf("ParseToken() error = %v", err)
			}
			sat, ok := parsed.(*ServiceAccessToken)
			if !ok {
				t.Fatalf("ParseToken() type = %T, want *ServiceAccessToken", parsed)
			}
			if sat.Subject() != tt.subject || sat.ClientID() != tt.clientID {
				t.Fatalf("parsed subject/client = %q/%q, want %q/%q", sat.Subject(), sat.ClientID(), tt.subject, tt.clientID)
			}
			if _, ok := sat.Scopes()["mail:send"]; !ok {
				t.Fatalf("Scopes() = %v, want mail:send", sat.Scopes())
			}
		})
	}
}

func TestDetectTypePreservesChallengeBusinessType(t *testing.T) {
	t.Parallel()
	input := NewClaimsBuilder().Issuer("aegis").ClientID("sender").Audience("chaos").
		ExpiresIn(time.Minute).Build(NewChallengeTokenBuilder().Subject("user@example.com").Type("sat"))
	wire, err := Build(input)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if got := DetectType(wire); got != TokenTypeChallenge {
		t.Fatalf("DetectType() = %q, want %q", got, TokenTypeChallenge)
	}
}
