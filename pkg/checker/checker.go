package checker

import (
	"log"
	"regexp"
	"unicode/utf8"
)

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
