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

func TestExamineInventoryItemInTheDark(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "take brass key")
	run(w, "south")
	run(w, "down")
	out := run(w, "examine brass key")
	if !strings.Contains(out, "brass") {
		t.Fatalf("examining a carried item should work in the dark, got:\n%s", out)
	}
	out = run(w, "examine gold coin")
	if !strings.Contains(out, "dark") {
		t.Fatalf("examining a room item in the dark should be blocked, got:\n%s", out)
	}
}

func TestDroppedLanternStillLightsTheRoom(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "take brass key")
	run(w, "south")
	run(w, "west")
	run(w, "take lantern")
	run(w, "east")
	run(w, "down")
	run(w, "drop lantern")
	out := run(w, "take gold coin")
	if !world.IsCarrying(w, "gold coin") {
		t.Fatalf("taking the coin should still work with the lantern on the floor, got:\n%s", out)
	}
	out = run(w, "take lantern")
	if !world.IsCarrying(w, "lantern") {
		t.Fatalf("taking the lantern back should work; got:\n%s", out)
	}
}

func TestPickUpAndPutDownAliases(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "pick up brass key")
	if !world.IsCarrying(w, "brass key") {
		t.Fatal("pick up brass key should have taken the brass key")
	}
	run(w, "put down brass key")
	if world.IsCarrying(w, "brass key") {
		t.Fatal("put down brass key should have dropped the brass key")
	}
}
