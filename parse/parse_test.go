package parse_test

import (
	"reflect"
	"testing"

	"terminal-go/parse"
)

func TestTokenizeEmpty(t *testing.T) {
	got := parse.Tokenize("   ")
	if got.Verb != "" || len(got.Args) != 0 || got.Object != "" {
		t.Fatalf("empty input produced %+v", got)
	}
}

func TestTokenizeAliasesAndFiller(t *testing.T) {
	cases := []struct {
		input string
		verb  string
		args  []string
	}{
		{"n", "north", nil},
		{"go to the north", "go", []string{"north"}},
		{"TAKE THE BRASS KEY", "take", []string{"brass", "key"}},
		{"x plaque", "examine", []string{"plaque"}},
		{"i", "inventory", nil},
		{"look", "look", nil},
		{"get lamp", "take", []string{"lamp"}},
		{"  drop  the  coin  ", "drop", []string{"coin"}},
		{"?", "help", nil},
	}
	for _, c := range cases {
		got := parse.Tokenize(c.input)
		if got.Verb != c.verb {
			t.Errorf("Tokenize(%q).Verb = %q, want %q", c.input, got.Verb, c.verb)
		}
		wantArgs := c.args
		if wantArgs == nil {
			wantArgs = []string{}
		}
		if len(got.Args) == 0 && len(wantArgs) == 0 {
			continue
		}
		if !reflect.DeepEqual(got.Args, wantArgs) {
			t.Errorf("Tokenize(%q).Args = %#v, want %#v", c.input, got.Args, wantArgs)
		}
	}
}

func TestTokenizeObjectJoinsArgs(t *testing.T) {
	got := parse.Tokenize("take the brass key")
	if got.Object != "brass key" {
		t.Errorf("Object = %q, want %q", got.Object, "brass key")
	}
}
