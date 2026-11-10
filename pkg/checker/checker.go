package checker

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
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

func FormatStructErrors(err error) []string {
	fieldErrors, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return []string{err.Error()}
	}

	var messages []string
	for _, fe := range fieldErrors {
		messages = append(messages, FormatError(fe, fe.Field()))
	}
	return messages
}

func FormatVarErrors(err error, fieldName string) []string {
	if len(fieldName) == 0 {
		fieldName = "This field"
	}

	fieldErrors, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return []string{err.Error()}
	}

	var messages []string
	for _, fe := range fieldErrors {
		messages = append(messages, FormatError(fe, fieldName))
	}
	return messages
}

func FormatError(fe validator.FieldError, fieldName string) string {
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "uuid", "uuidv4", "uuidv7":
		return fmt.Sprintf("%s must be a valid UUID.", fieldName)
	case TagUsername:
		return fmt.Sprintf("%s must be a valid username.", fieldName)
	case TagPassword:
		return fmt.Sprintf("%s must be a valid password.", fieldName)
	case "required":
		return fmt.Sprintf("%s is required.", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address.", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s.", fieldName, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s.", fieldName, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s.", fieldName, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s.", fieldName, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s].", fieldName, strings.ReplaceAll(param, " ", ", "))
	default:
		return fmt.Sprintf("%s failed the '%s' validation rule.", fieldName, tag)
	}
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
