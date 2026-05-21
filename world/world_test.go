package world_test

import (
	"bytes"
	"testing"

	"terminal-go/world"
)

func TestNewDemoShape(t *testing.T) {
	w := world.NewDemo()

	if got := len(w.Rooms.Name); got != 6 {
		t.Fatalf("rooms = %d, want 6", got)
	}
	if got := len(w.Items.Name); got != 6 {
		t.Fatalf("items = %d, want 6", got)
	}
	if w.Player.Room != w.GoalRoom {
		t.Fatal("player should start in the goal room")
	}

	for index, location := range w.Items.Location {
		if location == world.InventoryRoom {
			t.Fatalf("item %d starts in inventory; demo expects all items placed in rooms", index)
		}
		if int(location) < 0 || int(location) >= len(w.Rooms.Name) {
			t.Fatalf("item %d has invalid location %d", index, location)
		}
	}

	if want := len(w.Items.Name) + 1; len(w.Items.AliasStart) != want {
		t.Fatalf("AliasStart length = %d, want %d", len(w.Items.AliasStart), want)
	}

	for index, exit := range w.Exits {
		if int(exit.From) < 0 || int(exit.From) >= len(w.Rooms.Name) {
			t.Fatalf("exit %d has invalid From %d", index, exit.From)
		}
		if int(exit.To) < 0 || int(exit.To) >= len(w.Rooms.Name) {
			t.Fatalf("exit %d has invalid To %d", index, exit.To)
		}
	}
}

func TestItemTagsCombine(t *testing.T) {
	w := world.NewDemo()
	lantern := world.FindItemInRoom(w, roomByName(t, w, "Kitchen"), "lantern")
	if lantern == world.InvalidItem {
		t.Fatal("demo is missing the lantern item")
	}
	if !world.ItemHasTag(w, lantern, world.ItemTakeable) {
		t.Error("lantern should be takeable")
	}
	if !world.ItemHasTag(w, lantern, world.ItemLit) {
		t.Error("lantern should be lit")
	}
	if world.ItemHasTag(w, lantern, world.ItemReadable) {
		t.Error("lantern should not be readable")
	}
}

func TestItemAliasesLookup(t *testing.T) {
	w := world.NewDemo()
	for index, name := range w.Items.Name {
		id := world.ItemID(index)
		aliases := world.ItemAliases(w, id)
		if name == "plaque" {
			if len(aliases) != 2 || aliases[0] != "sign" || aliases[1] != "wooden plaque" {
				t.Errorf("plaque aliases = %v, want [sign, wooden plaque]", aliases)
			}
		}
		if name == "brass key" {
			if len(aliases) != 2 || aliases[0] != "key" || aliases[1] != "brass" {
				t.Errorf("brass key aliases = %v, want [key, brass]", aliases)
			}
		}
	}
}

func TestIsDarkOutOfBoundsRoom(t *testing.T) {
	w := world.NewDemo()
	if world.IsDark(w, world.InventoryRoom) {
		t.Fatal("IsDark on InventoryRoom should return false, not panic")
	}
	if world.IsDark(w, world.RoomID(999)) {
		t.Fatal("IsDark on out-of-bounds room should return false")
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	w := world.NewDemo()

	lantern := world.FindItemInRoom(w, roomByName(t, w, "Kitchen"), "lantern")
	w.Items.Location[lantern] = world.InventoryRoom
	w.Player.Room = roomByName(t, w, "Cellar")

	var buf bytes.Buffer
	if err := world.Encode(w, &buf); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	restored, err := world.Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if restored.Version != world.SaveVersion {
		t.Errorf("Version = %d, want %d", restored.Version, world.SaveVersion)
	}
	if restored.Player.Room != w.Player.Room {
		t.Errorf("Player.Room = %d, want %d", restored.Player.Room, w.Player.Room)
	}
	if restored.Items.Location[lantern] != world.InventoryRoom {
		t.Error("lantern location lost in save/load")
	}
	if !world.ItemHasTag(restored, lantern, world.ItemLit) {
		t.Error("lantern lit tag lost in save/load")
	}
	if len(restored.Items.AliasStart) != len(w.Items.AliasStart) {
		t.Errorf("AliasStart length = %d, want %d", len(restored.Items.AliasStart), len(w.Items.AliasStart))
	}
}

func TestDecodeRejectsVersionMismatch(t *testing.T) {
	w := world.NewDemo()
	w.Version = world.SaveVersion + 99

	var buf bytes.Buffer
	if err := world.Encode(w, &buf); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if _, err := world.Decode(&buf); err == nil {
		t.Fatal("Decode should reject a wrong-version blob")
	}
}

func roomByName(t *testing.T, w *world.World, name string) world.RoomID {
	t.Helper()
	for index, n := range w.Rooms.Name {
		if n == name {
			return world.RoomID(index)
		}
	}
	t.Fatalf("room %q not found", name)
	return -1
}
