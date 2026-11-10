package phash

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerate(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		pass := "validPassword"

		got, err := Generate(pass)
		if err != nil {
			t.Fatalf("Failed generating password: %v\n", err)
		}

		if err := bcrypt.CompareHashAndPassword(got, []byte(pass)); err != nil {
			t.Fatalf("Invalid Password: %v\n", err)
		}
	})
	t.Run("Fails", func(t *testing.T) {
		cases := []struct {
			input string
			err   error
		}{
			{"short", ErrTooShort},
			{"tooLongTooLongtooLongtooLongtooLong", ErrTooLong},
		}
		for _, c := range cases {
			_, err := Generate(c.input)
			if err == nil {
				t.Fatalf("Must return an error\n")
			}

			if !errors.Is(err, c.err) {
				t.Fatalf("Invalid error: got %v, want %v\n", err, c.err)
			}
		}
	})
}
