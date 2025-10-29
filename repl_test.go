package main
import (
	"testing"
	"fmt"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
	input string
	expected []string
	}{
		{
			input: "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "sunova	doubletabbing		",
			expected: []string{"sunova","doubletabbing"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "Happy Birthday Pikachu!",
			expected: []string{"happy","birthday", "pikachu!"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("FAIL - expected length %d; actual length %d", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("FAIL - expected word at %d = %s; actual word = %s", i, expectedWord, word)
			}
		}

		fmt.Println("PASS")
	}

}
