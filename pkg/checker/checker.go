package checker

import (
	"log"
	"regexp"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func New() Checker {
	val := validator.New()

	// Register custom tags
	for _, cc := range customCheckers {
		val.RegisterValidation(cc.Tag, cc.Validator)
	}

	return Checker{val}
}

type Checker struct {
	val *validator.Validate
}

func (c *Checker) Struct(input any) error {
	return c.val.Struct(input)
}

func (c *Checker) Var(field any, tag string) error {
	return c.val.Var(field, tag)
}

type Validator interface {
	Struct(input any) error
	Var(field any, tag string) error
}

var customCheckers = []struct {
	Tag       string
	Validator func(fl validator.FieldLevel) bool
}{
	{TagUsername, func(fl validator.FieldLevel) bool {
		return Username(fl.Field().String())
	}},
	{TagPassword, func(fl validator.FieldLevel) bool {
		return Password(fl.Field().String())
	}},
}

const (
	CharDigit     = `\d`
	CharLowerCase = `[a-z]`
	CharUpperCase = `[A-Z]`
	CharSpecial   = `[~!@#$%^&*()_+\x60\-={}|[\]\\":<>?,./]`
	CharUsername  = `^[a-zA-Z][a-zA-Z0-9]*$`
	MinPassword   = 8
	MaxPassword   = 20
	MinUsername   = 6
	MaxUsername   = 16
	TagUsername   = "username"
	TagPassword   = "password"
)

// Pre-compiled Regexes,
// Panics when it failed to parsed the expressions.
var (
	RegexDigit     = regexp.MustCompile(CharDigit)
	RegexLowerCase = regexp.MustCompile(CharLowerCase)
	RegexUpperCase = regexp.MustCompile(CharUpperCase)
	RegexSpecial   = regexp.MustCompile(CharSpecial)
	RegexUsername  = regexp.MustCompile(CharUsername)
)

// Verify Password format with:
//
//	8-20 length,
//	a lower case [a-z],
//	an upper case [A-Z],
//	a special character [~!@#$%^&*()_+\x60\-={}|[\]\\":<>?,./].
func Password(input string) bool {
	if !verifyLen(input, MinPassword, MaxPassword) {
		return false
	}

	hasLower := RegexLowerCase.MatchString(input)
	if !hasLower {
		return false
	}

	hasUpper := RegexUpperCase.MatchString(input)
	if !hasUpper {
		return false
	}

	hasDigit := RegexDigit.MatchString(input)
	if !hasDigit {
		return false
	}

	hasSpecial := RegexSpecial.MatchString(input)
	if !hasSpecial {
		return false
	}

	return true
}

// Verify Username format with:
//
//	6-16 length,
//	letters and numbers `^[a-zA-Z][a-zA-Z0-9]*$`,
//	username can't be only numbers,
//	numbers can't be at front.
func Username(input string) bool {
	if !verifyLen(input, MinUsername, MaxUsername) {
		return false
	}
	return RegexUsername.MatchString(input)
}

func verifyLen(input string, min, max int) bool {
	length := utf8.RuneCountInString(input)
	if length < min || length > max {
		return false
	}
	return true
}

func reportPassword(kind string, err error) {
	log.Printf("Err: Validate password has %s: %s", kind, err)
}
