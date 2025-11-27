package main

import (
	"testing"
)

func TestValidateMessage(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input:    "Thanks to sharbert I'll do nothing today",
			expected: "Thanks to **** I'll do nothing today",
		},
		{
			input:    "Hakuna matata!",
			expected: "Hakuna matata!",
		},
		{
			input:    "Sharbert! Don't do this",
			expected: "Sharbert! Don't do this",
		},
	}

	for _, c := range cases {
		actual := validateMessage(c.input)

		if actual != c.expected {
			t.Errorf("tested word is different than the expected word")
		}
	}
}
