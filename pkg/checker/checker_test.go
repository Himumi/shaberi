package checker

import (
	"fmt"
	"strings"
	"testing"
)

func TestPassword(t *testing.T) {
	type Case struct {
		name   string
		input  string
		result bool
	}

	cases := []Case{
		{fmt.Sprintf("less than %d chars", MinPassword), "less", false},
		{fmt.Sprintf("more than %d chars", MaxPassword), "morethantwentymorethantwentymorethantwentymorethantwenty", false},
		{"not have a lower case", "NOTHAVEALOWERCASE", false},
		{"not have an upper case", "nothaveanuppercase", false},
		{"not have a digit", "NOTHAVEadigit", false},
		{"not have a special char", "NOTHAVE0specialchar", false},
		{"correct password format", "SatifiedRegex10?", true},
	}

	for _, c := range cases {
		result := Password(c.input)
		t.Run(c.name, func(t *testing.T) {
			reportTestPassword(t, c.input, result, c.result)
		})
	}
}

func reportTestPassword(t testing.TB, input string, got bool, want bool) {
	t.Helper()
	if got != want {
		const format = "Failed validating password: input %s: got %t, want %t\n"
		t.Fatalf(format, input, got, want)
	}
}

func TestUsername(t *testing.T) {
	type Case struct {
		name   string
		input  string
		result bool
	}

	cases := []Case{
		{fmt.Sprintf("less than %d chars", MinUsername), "less", false},
		{fmt.Sprintf("more than %d chars", MaxUsername), "morethansixteenMORETHANSIXTEEN", false},
		{"Contain not allowed chars (underscores)", "contain_", false},
		{"Contain not allowed chars (spaces)", "user name", false},
		{"contain only numbers", "1234567", false},
		{"contain numbers at front", "123user", false},
		{"contain upper case", "USERNAME", true},
		{"contain lower case", "username", true},
		{"contain numbers", "username123", true},
		{"contain numbers at the middle", "user123name", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := Username(c.input)
			reportTestUsername(t, c.input, result, c.result)
		})
	}
}

func reportTestUsername(t testing.TB, input string, got, want bool) {
	t.Helper()
	if got != want {
		const format = "Failed validating username: input %s: got %t, want %t"
		t.Fatalf(format, input, got, want)
	}
}

func TestVarErrors(t *testing.T) {
	cases := []struct {
		input any
		tag   string
		msgs  []string
	}{
		{"123", "username", []string{"This field must be a valid username."}},
		{"123", "password", []string{"This field must be a valid password."}},
		{"", "required", []string{"This field is required."}},
		{"invalid@email", "email", []string{"This field must be a valid email address."}},
		{9, "min=10", []string{"This field must be at least 10."}},
		{11, "max=10", []string{"This field must not exceed 10."}},
		{9, "gte=10", []string{"This field must be greater than or equal to 10."}},
		{"red", "oneof=yellow blue", []string{"This field must be one of [yellow, blue]."}},
	}

	ch := New()
	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			err := ch.Var(c.input, c.tag)
			reportTestVar(t, err, c.tag, c.msgs)
		})
	}
}

func reportTestVar(t testing.TB, err error, tag string, wantMsgs []string) {
	t.Helper()
	if err == nil {
		t.Fatalf("Must return an error for %s tag\n", tag)
	}
	msgs := FormatVarErrors(err, "")

	gotLen := len(msgs)
	wantLen := len(wantMsgs)

	if gotLen != wantLen {
		t.Fatalf("Invalid messages length: got %d, want %d\n", gotLen, wantLen)
	}
	for i, gotMsg := range msgs {
		wantMsg := wantMsgs[i]
		if !strings.EqualFold(gotMsg, wantMsg) {
			t.Fatalf("Invalid message: got %v, want %v\n", gotMsg, wantMsg)
		}
	}
}
