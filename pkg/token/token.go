package token

import (
	"errors"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptySecret  = errors.New("Secret is empty")
	ErrEmptySubject = errors.New("Subject is empty")
	SigningMethod   = jwt.SigningMethodHS256
)

// Generate JWT token string
//
// Possible Errors:
//   - ErrEmptySecret
//   - ErrEmptySubject
//   - jwt.Token.SignedString's Errors
func Generate(secret string, sub uuid.UUID, exp int64) (string, error) {
	if len(secret) == 0 {
		return "", ErrEmptySecret
	}
	if sub == uuid.Nil() {
		return "", ErrEmptySubject
	}

	t := jwt.NewWithClaims(SigningMethod, jwt.MapClaims{
		"sub": sub,
		"exp": exp,
	})

	str, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return str, nil
}
