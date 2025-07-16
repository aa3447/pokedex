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

	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) != len(c.expected) {
			t.Errorf("Expected length %d, got %d", len(c.expected), len(result))
		}
		for i, v := range result {
			if v != c.expected[i] {
				t.Errorf("Expected %s, got %s", c.expected[i], v)
			}
		}
	}
}