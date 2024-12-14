package main

import (
	"testing"
)

func TestProfanityFilter(t *testing.T) {
	cases := []struct{
		input string
		expected string
	}{
		{
			input: "hello kerfuffle",
			expected: "hello ****",
		},
		{
			input: "hello sharbert",
			expected: "hello ****",
		},
		{
			input: "hello fornax",
			expected: "hello ****",
		},
		{
			input: "hello KERFUFFLE",
			expected: "hello ****",
		},
		{
			input: "hello SHARBERT",
			expected: "hello ****",
		},
		{
			input: "hello FORNAX",
			expected: "hello ****",
		},
		{
			input: "hello Kerfuffle",
			expected: "hello ****",
		},
		{
			input: "hello Sharbert",
			expected: "hello ****",
		},
		{
			input: "hello Fornax",
			expected: "hello ****",
		},
	}

	for _, c := range cases {
		actual := profanityFilter(c.input)
		if actual != c.expected {
			t.Errorf("Test failed for message: %v\n", c.input)
			t.Fail()
		}
	}
}