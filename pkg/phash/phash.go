package phash

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/Himumi/shaberi/pkg/checker"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPassword = checker.MinPassword
	MaxPassword = checker.MaxPassword
	CostInt     = bcrypt.DefaultCost
)

var (
	ErrTooShort        = fmt.Errorf("Password must be at least %d.", MinPassword)
	ErrTooLong         = fmt.Errorf("Password must not exceed %d.", MaxPassword)
	ErrInvalidPassword = errors.New("Password is invalid.")
)

// Generate an encrypted password.
//
// Possible Errors:
//
//   - ErrTooShort
//   - ErrTooLong
//   - bcrypt.GenerateFromPassword's Errors
func Generate(password string) ([]byte, error) {
	passLen := utf8.RuneCountInString(password)
	if passLen < MinPassword {
		return []byte{}, ErrTooShort
	}

	if passLen > MaxPassword {
		return []byte{}, ErrTooLong
	}

	return bcrypt.GenerateFromPassword([]byte(password), CostInt)
}

// Verify an encrypted password with the input
//
// Possible Errors:
//
//   - ErrInvalidPassword
//   - bcrypt.CompareHashAndPassword's Errors
func Verify(hashedPassword, input string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}
