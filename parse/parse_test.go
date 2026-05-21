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

func TestTokenizeTwoWordVerbs(t *testing.T) {
	cases := []struct {
		input  string
		verb   string
		object string
	}{
		{"pick up the brass key", "take", "brass key"},
		{"put down the lantern", "drop", "lantern"},
		{"pick up coin", "take", "coin"},
	}
	for _, c := range cases {
		got := parse.Tokenize(c.input)
		if got.Verb != c.verb {
			t.Errorf("Tokenize(%q).Verb = %q, want %q", c.input, got.Verb, c.verb)
		}
		if got.Object != c.object {
			t.Errorf("Tokenize(%q).Object = %q, want %q", c.input, got.Object, c.object)
		}
	}
}

func TestTokenizeRawObjectPreserveCase(t *testing.T) {
	got := parse.Tokenize("save MyGame.save")
	if got.Verb != "save" {
		t.Errorf("Verb = %q, want %q", got.Verb, "save")
	}
	if got.Object != "mygame.save" {
		t.Errorf("Object = %q, want %q", got.Object, "mygame.save")
	}
	if got.RawObject != "MyGame.save" {
		t.Errorf("RawObject = %q, want %q", got.RawObject, "MyGame.save")
	}
}

func TestTokenizeGrabIsTake(t *testing.T) {
	got := parse.Tokenize("grab the coin")
	if got.Verb != "take" {
		t.Errorf("Verb = %q, want %q", got.Verb, "take")
	}
	if got.Object != "coin" {
		t.Errorf("Object = %q, want %q", got.Object, "coin")
	}
}

func TestTokenizeBarePickAndPutAreNotAliased(t *testing.T) {
	if got := parse.Tokenize("pick coin"); got.Verb != "pick" {
		t.Errorf("'pick coin' Verb = %q, want %q (single-word pick should not normalize)", got.Verb, "pick")
	}
	if got := parse.Tokenize("put coin"); got.Verb != "put" {
		t.Errorf("'put coin' Verb = %q, want %q (single-word put should not normalize)", got.Verb, "put")
	}
}
