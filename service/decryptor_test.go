package service

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/heliantheon/aegis-go/utilities/key"
	pasetokit "github.com/heliantheon/aegis-go/utilities/paseto"
	tokendef "github.com/heliantheon/aegis-go/utilities/token"
)

func TestInterpretSignedServiceToken(t *testing.T) {
	t.Parallel()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	signingKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(privateKey)
	if err != nil {
		t.Fatalf("NewV4AsymmetricSecretKeyFromBytes() error = %v", err)
	}
	pid, err := pasetokit.ComputePID(signingKey.Public())
	if err != nil {
		t.Fatalf("ComputePID() error = %v", err)
	}
	footer, err := pasetokit.NewFooter(pid).Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	input := tokendef.NewClaimsBuilder().Issuer("aegis").ClientID("sender").
		Audience("chaos").Subject("service-principal").ExpiresIn(time.Minute).
		Build(tokendef.NewServiceAccessTokenBuilder().Scope("mail:send"))
	wire, err := tokendef.Build(input)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	wire.SetFooter(footer)
	signed := wire.V4Sign(signingKey, nil)
	provider := key.SingleOf(func(_ context.Context, id string) ([]byte, error) {
		if id != "sender" {
			return nil, key.ErrNotFound
		}
		return publicKey, nil
	})

	for _, tt := range []struct {
		name     string
		audience string
		wantErr  bool
	}{
		{name: "matching audience", audience: "chaos"},
		{name: "wrong audience", audience: "other", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parsed, err := NewDecryptor(tt.audience, nil, provider).Interpret(t.Context(), signed)
			if tt.wantErr {
				if err == nil {
					t.Fatal("Interpret() accepted a token for another audience")
				}
				return
			}
			if err != nil {
				t.Fatalf("Interpret() error = %v", err)
			}
			if _, ok := parsed.(*tokendef.ServiceAccessToken); !ok {
				t.Fatalf("Interpret() type = %T, want *ServiceAccessToken", parsed)
			}
			if parsed.Subject() != "service-principal" || parsed.ClientID() != "sender" {
				t.Fatalf("Interpret() subject/client = %q/%q, want service-principal/sender", parsed.Subject(), parsed.ClientID())
			}
		})
	}
}
