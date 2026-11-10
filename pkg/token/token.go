package token

import (
	"errors"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken  = errors.New("Token is expired")
	ErrInvalidToken  = errors.New("Token is invalid")
	ErrEmptySecret   = errors.New("Secret is empty")
	ErrEmptySubject  = errors.New("Subject is empty")
	ErrEmptyTokenStr = errors.New("Token string is empty")
	SigningMethod    = jwt.SigningMethodHS256
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

type Claims struct {
	// Override Subject type
	Subject uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

// Parse token string to a Custom Claims
//
// Possible Errors:
//   - ErrEmptySecret
//   - ErrEmptyTokenStr
//   - ErrExpiredToken
//   - ErrInvalidToken
//   - jwt.ParseWithClaims's Errors
func Parse(tokenStr, secret string) (*Claims, error) {
	if len(secret) == 0 {
		return nil, ErrEmptySecret
	}

	if len(tokenStr) == 0 {
		return nil, ErrEmptyTokenStr
	}

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	if token.Valid {
		return &claims, nil
	}
	return nil, ErrInvalidToken
}
