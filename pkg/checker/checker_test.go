package checker

import (
	"fmt"
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
