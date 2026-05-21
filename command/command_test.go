package command_test

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestLockedDoorRequiresKey(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	out := run(w, "north")
	if !strings.Contains(out, "locked") {
		t.Fatalf("going north of hallway without key should mention locked, got:\n%s", out)
	}
	if w.Rooms.Name[w.PlayerRoom] != "Hallway" {
		t.Fatalf("player should still be in hallway, got %s", w.Rooms.Name[w.PlayerRoom])
	}
}

func TestWalkthroughWins(t *testing.T) {
	w := world.NewDemo()
	steps := []string{
		"north",
		"east",
		"take brass key",
		"west",
		"west",
		"take lantern",
		"east",
		"north",
		"down",
		"take gold coin",
		"up",
		"south",
		"south",
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
	run(w, "east")
	run(w, "take brass key")
	run(w, "west")
	run(w, "north")
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
	run(w, "east")
	out := run(w, "examine key")
	if !strings.Contains(out, "brass") {
		t.Fatalf("examine key should describe the brass key, got:\n%s", out)
	}
}

func TestExamineInventoryItemInTheDark(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "east")
	run(w, "take brass key")
	run(w, "west")
	run(w, "north")
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

func TestSilverCandleLightsTheCellar(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "east")
	run(w, "take brass key")
	run(w, "west")
	run(w, "north")
	run(w, "take silver candle")
	run(w, "down")
	out := run(w, "take gold coin")
	if !world.IsCarrying(w, "gold coin") {
		t.Fatalf("silver candle should light the cellar; output:\n%s", out)
	}
}

func TestDroppedLanternStillLightsTheRoom(t *testing.T) {
	w := world.NewDemo()
	run(w, "north")
	run(w, "east")
	run(w, "take brass key")
	run(w, "west")
	run(w, "west")
	run(w, "take lantern")
	run(w, "east")
	run(w, "north")
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
	run(w, "east")
	run(w, "pick up brass key")
	if !world.IsCarrying(w, "brass key") {
		t.Fatal("pick up brass key should have taken the brass key")
	}
	run(w, "put down brass key")
	if world.IsCarrying(w, "brass key") {
		t.Fatal("put down brass key should have dropped the brass key")
	}
}

func TestSaveAndLoadCommands(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.save")

	w := world.NewDemo()
	run(w, "north")
	run(w, "east")
	run(w, "take brass key")
	out := run(w, "save "+path)
	if !strings.Contains(out, "saved") {
		t.Fatalf("save should report success, got:\n%s", out)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("save file should exist: %v", err)
	}

	w2 := world.NewDemo()
	out = run(w2, "load "+path)
	if !strings.Contains(out, "restored") {
		t.Fatalf("load should report success, got:\n%s", out)
	}
	if !world.IsCarrying(w2, "brass key") {
		t.Fatal("loaded world should have the brass key in inventory")
	}
	if w2.Rooms.Name[w2.PlayerRoom] != "Library" {
		t.Fatalf("loaded world should place player in library, got %s", w2.Rooms.Name[w2.PlayerRoom])
	}
}

func TestSavePreservesFilenameCase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "MixedCase.Save")

	w := world.NewDemo()
	out := run(w, "save "+path)
	if !strings.Contains(out, path) {
		t.Errorf("save output should include the original-case path %q, got:\n%s", path, out)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range matches {
		if filepath.Base(m) == "MixedCase.Save" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected file named MixedCase.Save, found: %v", matches)
	}
}

func TestLoadOfCompletedRunReportsIt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "winning.save")

	w := world.NewDemo()
	w.Won = true
	if err := writeSave(w, path); err != nil {
		t.Fatal(err)
	}

	w2 := world.NewDemo()
	out := run(w2, "load "+path)
	if !strings.Contains(out, "completed") {
		t.Fatalf("loading a completed save should acknowledge it, got:\n%s", out)
	}
}

func writeSave(w *world.World, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return world.Encode(w, file)
}
