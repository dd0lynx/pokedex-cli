package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "  hello  world  ",
			want:  []string{"hello", "world"},
		},
		{
			input: "  heLLoWorLd  ",
			want:  []string{"helloworld"},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			want:  []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		got := CleanInput(c.input)
		if len(got) != len(c.want) {
			t.Errorf("CleanInput: Unexpected lengeth for input: '%s'\nexpeceted length: %d %#v\nactual length: %d %#v", c.input, len(c.want), c.want, len(got), got)
			continue
		}

		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("CleanInput: Values don't match for input: '%s'\nExpected: %#v\nActual: %#v", c.input, c.want, got)
				break
			}
		}
	}
}
