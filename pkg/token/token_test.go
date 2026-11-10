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
