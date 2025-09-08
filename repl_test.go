package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  some people wiLl NeVer undErstand ",
			expected: []string{"some", "people", "will", "never", "understand"},
		},
		{
			input:    "  wHat IS eVeN clEaNInPUT ",
			expected: []string{"what", "is", "even", "cleaninput"},
		},
		{
			input:    "     this is another TEST  ",
			expected: []string{"this", "is", "another", "test"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected slice length: %v does not match Actual slice length: %v", len(c.expected), len(actual))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%v does not match %v", word, expectedWord)
			}
		}
	}
}
