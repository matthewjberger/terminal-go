package world_test

import (
	"testing"

	"terminal-go/world"
)

func TestNewDemoShape(t *testing.T) {
	w := world.NewDemo()

	if got := len(w.Rooms.Name); got != 4 {
		t.Fatalf("rooms = %d, want 4", got)
	}
	if got := len(w.Items.Name); got != 4 {
		t.Fatalf("items = %d, want 4", got)
	}
	if w.Player.Room != w.GoalRoom {
		t.Fatalf("player should start in the goal room")
	}

	for index, location := range w.Items.Location {
		if location == world.InventoryRoom {
			t.Fatalf("item %d starts in inventory; demo expects all items placed in rooms", index)
		}
		if int(location) < 0 || int(location) >= len(w.Rooms.Name) {
			t.Fatalf("item %d has invalid location %d", index, location)
		}
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

func TestDirectionFromWord(t *testing.T) {
	cases := map[string]world.Direction{
		"n":     world.North,
		"north": world.North,
		"s":     world.South,
		"south": world.South,
		"e":     world.East,
		"east":  world.East,
		"w":     world.West,
		"west":  world.West,
		"u":     world.Up,
		"up":    world.Up,
		"d":     world.Down,
		"down":  world.Down,
	}
	for word, want := range cases {
		got, ok := world.DirectionFromWord(word)
		if !ok {
			t.Errorf("DirectionFromWord(%q) ok = false", word)
			continue
		}
		if got != want {
			t.Errorf("DirectionFromWord(%q) = %v, want %v", word, got, want)
		}
	}
	if _, ok := world.DirectionFromWord("nowhere"); ok {
		t.Errorf("DirectionFromWord(%q) should fail", "nowhere")
	}
}

func TestIsDarkRequiresLantern(t *testing.T) {
	w := world.NewDemo()
	cellar := world.RoomID(-1)
	for index, name := range w.Rooms.Name {
		if name == "Cellar" {
			cellar = world.RoomID(index)
		}
	}
	if cellar < 0 {
		t.Fatal("demo is missing the Cellar room")
	}

	if !world.IsDark(w, cellar) {
		t.Fatal("cellar should be dark without a lantern")
	}

	lantern := world.FindItemInRoom(w, world.RoomID(2), "lantern")
	if lantern == world.InvalidItem {
		for index := range w.Items.Name {
			if w.Items.Name[index] == "lantern" {
				lantern = world.ItemID(index)
			}
		}
	}
	if lantern == world.InvalidItem {
		t.Fatal("demo is missing the lantern item")
	}
	w.Items.Location[lantern] = world.InventoryRoom
	if world.IsDark(w, cellar) {
		t.Fatal("cellar should be lit when carrying the lantern")
	}
}
