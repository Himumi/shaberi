package token

import (
	"errors"
	"testing"
	"time"
	"uuid"
)

func TestGenerate(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		secret := "secret"
		sub := uuid.NewV4()
		exp := time.Now().Add(time.Hour).Unix()

		tokenStr, err := Generate(secret, sub, exp)
		if err != nil {
			t.Fatalf("Failed generating token: %v\n", err)
		}

		if len(tokenStr) == 0 {
			t.Fatalf("Failed generating token string: %s\n", tokenStr)
		}
	})
	t.Run("Fails", func(t *testing.T) {
		cases := []struct {
			secret string
			sub    uuid.UUID
			err    error
		}{
			{"", uuid.NewV4(), ErrEmptySecret},
			{"secret", uuid.Nil(), ErrEmptySubject},
		}

		for _, c := range cases {
			t.Run(c.err.Error(), func(t *testing.T) {
				_, err := Generate(c.secret, c.sub, 0)
				if !errors.Is(err, c.err) {
					t.Fatalf("Invalid error: got %v, want %v\n", err, c.err)
				}
			})
		}
	})
}

func TestParse(t *testing.T) {
	secret := "my-super-secret-key"
	wrongSecret := "wrong-secret-key"
	testUUID := uuid.New()

	// Generate a valid token (expires 1 hour from now)
	validExp := time.Now().Add(time.Hour).Unix()
	validTokenStr, err := Generate(secret, testUUID, validExp)
	if err != nil {
		t.Fatalf("failed to generate valid test token: %v", err)
	}

	// Generate an expired token (expired 1 hour ago)
	expiredExp := time.Now().Add(-1 * time.Hour).Unix()
	expiredTokenStr, err := Generate(secret, testUUID, expiredExp)
	if err != nil {
		t.Fatalf("failed to generate expired test token: %v", err)
	}

	// Generate a token with a different secret
	wrongSecretTokenStr, err := Generate(wrongSecret, testUUID, validExp)
	if err != nil {
		t.Fatalf("failed to generate wrong-secret test token: %v", err)
	}

	cases := []struct {
		name      string
		tokenStr  string
		secret    string
		wantSub   uuid.UUID
		wantError error
	}{
		{
			name:      "Error: Empty Secret",
			tokenStr:  validTokenStr,
			secret:    "",
			wantSub:   uuid.Nil(),
			wantError: ErrEmptySecret,
		},
		{
			name:      "Error: Empty Token String",
			tokenStr:  "",
			secret:    secret,
			wantSub:   uuid.Nil(),
			wantError: ErrEmptyTokenStr,
		},
		{
			name:      "Success: Valid Token",
			tokenStr:  validTokenStr,
			secret:    secret,
			wantSub:   testUUID,
			wantError: nil,
		},
		{
			name:      "Error: Expired Token",
			tokenStr:  expiredTokenStr,
			secret:    secret,
			wantSub:   uuid.Nil(),
			wantError: ErrExpiredToken,
		},
		{
			name:      "Error: Invalid Signature (Wrong Secret)",
			tokenStr:  wrongSecretTokenStr,
			secret:    secret,
			wantSub:   uuid.Nil(),
			wantError: nil, // jwt parser returns a signature error, not nil
		},
		{
			name:      "Error: Malformed Token String",
			tokenStr:  "not.a.real.jwt.token",
			secret:    secret,
			wantSub:   uuid.Nil(),
			wantError: nil, // jwt parser returns a malformed error
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			claims, err := Parse(c.tokenStr, c.secret)

			// Check error expectations
			if c.wantError != nil {
				if !errors.Is(err, c.wantError) {
					t.Fatalf("expected error %v, got %v", c.wantError, err)
				}
			} else {
				if c.tokenStr == validTokenStr && err != nil {
					t.Fatalf("unexpected error for valid token: %v", err)
				}
				if c.tokenStr != validTokenStr && err == nil {
					t.Fatalf("expected an error for invalid/malformed token, got nil")
				}
			}

			// Verify Claims.Subject
			if c.wantError == nil && c.tokenStr == validTokenStr {
				if claims.Subject != c.wantSub {
					t.Errorf("expected subject %v, got %v", c.wantSub, claims.Subject)
				}
			} else {
				if claims != nil {
					t.Errorf("expected nil claims on error, got %+v", claims)
				}
			}
		})
	}
}
