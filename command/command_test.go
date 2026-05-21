package command_test

import (
	"bytes"
	"strings"
	"testing"

	"terminal-go/command"
	"terminal-go/parse"
	"terminal-go/world"
)

func run(w *world.World, line string) string {
	var buf bytes.Buffer
	command.Run(w, parse.Tokenize(line), &buf)
	return buf.String()
}

func TestLockedExitRequiresKey(t *testing.T) {
	w := world.NewDemo()
	out := run(w, "down")
	if !strings.Contains(out, "locked") {
		t.Fatalf("going down without key should mention being locked, got:\n%s", out)
	}
	if w.Player.Room != w.GoalRoom {
		t.Fatalf("player should not have moved; room = %d", w.Player.Room)
	}
}

func TestWalkthroughWins(t *testing.T) {
	w := world.NewDemo()
	steps := []string{
		"north",
		"take brass key",
		"south",
		"west",
		"take lantern",
		"east",
		"down",
		"take gold coin",
		"up",
	}
	var buf bytes.Buffer
	for _, step := range steps {
		buf.Reset()
		command.Run(w, parse.Tokenize(step), &buf)
		if w.Won {
			break
		}
	}
	if !w.Won {
		t.Fatalf("expected to win after walkthrough, last output:\n%s", buf.String())
	}
}

func TestDarkCellarBlocksTake(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "take brass key")
	run(w, "south")
	run(w, "down")
	out := run(w, "take gold coin")
	if !strings.Contains(out, "dark") {
		t.Fatalf("taking coin in unlit cellar should mention darkness, got:\n%s", out)
	}
	if world.IsCarrying(w, "gold coin") {
		t.Fatal("coin should not have been taken in the dark")
	}
}

func TestExamineUsesAliases(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	out := run(w, "examine key")
	if !strings.Contains(out, "brass") {
		t.Fatalf("examine key should describe the brass key, got:\n%s", out)
	}
}
