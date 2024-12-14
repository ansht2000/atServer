package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	cases := []struct{
		input string
		expected error
	}{
		{
			input: "password",
			expected: nil,
		},
		{
			input: "better_password",
			expected: nil,
		},
		{
			input: "betterpassword",
			expected: nil,
		},
		{
			input: "better_word",
			expected: nil,
		},
		{
			input: "bettord",
			expected: nil,
		},
		{
			input: "better_paord",
			expected: nil,
		},
		{
			input: "tter_password",
			expected: nil,
		},
	}

	for _, c := range cases {
		hashedPass, _ := HashPassword(c.input)
		actual := CheckPasswordHash(c.input, hashedPass)
		if actual != c.expected {
			t.Errorf("Test failed for password: %v\n", c.input)
			t.Fail()
		}
	}
}