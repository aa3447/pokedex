package main

import (
	"testing"
)

type testCase struct {
		input string
		expected []string
}

func TestCleanInput(t *testing.T){

	cases := []testCase{
		{
		input: "Hello World",
		expected: []string{"hello", "world"},
		},
		{
		input: "hello world",
		expected: []string{"hello", "world"},
		},
		{
		input: "Hello World     ",
		expected: []string{"hello", "world"},
		},
		{
		input: "     Hello World",
		expected: []string{"hello", "world"},
		},
	}

	for _, test := range cases {
		result := cleanInput(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("Expected length %d, got %d", len(test.expected), len(result))
		}
		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("Expected %s, got %s", test.expected[i], v)
			}
		}
	}
}
