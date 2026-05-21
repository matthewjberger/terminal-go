// Package world holds the game state and the queries over it.
//
// The world is a single [World] value with parallel-array tables
// keyed by typed [RoomID] / [ItemID] indices. Mutation lives in
// the free functions in this package and in package command;
// methods are not used for game logic.
package world

// RoomID indexes every Rooms.* slice. Two negative sentinels are
// reserved: InvalidRoom means "no room", InventoryRoom means "the
// player is carrying this item".
type RoomID int32

// ItemID indexes every Items.* slice. InvalidItem means "no item".
type ItemID int32

const (
	InvalidRoom   RoomID = -1
	InventoryRoom RoomID = -2

	InvalidItem ItemID = -1
)

// Direction is one of the six exit directions.
type Direction int8

const (
	North Direction = iota
	South
	East
	West
	Up
	Down
)

// Rooms is the room table in struct-of-arrays form. Every slice
// has the same length; a RoomID is a valid index into all of them.
type Rooms struct {
	Name        []string
	Description []string
	Dark        []bool
}

// Items is the item table in struct-of-arrays form. Location[id]
// is the RoomID the item currently sits in; InventoryRoom means
// it is carried by the player. Aliases[id] lists alternative
// noun phrases the parser will accept. Lit[id] is true when the
// item dispels darkness while carried.
type Items struct {
	Name        []string
	Description []string
	Aliases     [][]string
	Location    []RoomID
	Takeable    []bool
	Readable    []bool
	ReadText    []string
	Lit         []bool
}

// Exit is one directed edge in the room graph. Locked exits open
// when the player presents KeyItem.
type Exit struct {
	From    RoomID
	Dir     Direction
	To      RoomID
	Locked  bool
	KeyItem ItemID
}

// Player holds the per-player state. Inventory is derived from
// Items.Location == InventoryRoom, not stored here.
type Player struct {
	Room RoomID
}

// World is the single owner of all game state.
type World struct {
	Rooms  Rooms
	Items  Items
	Exits  []Exit
	Player Player

	GoalRoom RoomID
	GoalItem ItemID

	Quit bool
	Won  bool
}

// DirectionFromWord parses a direction word or single-letter
// abbreviation. Recognized: north/n, south/s, east/e, west/w,
// up/u, down/d.
func DirectionFromWord(word string) (Direction, bool) {
	switch word {
	case "n", "north":
		return North, true
	case "s", "south":
		return South, true
	case "e", "east":
		return East, true
	case "w", "west":
		return West, true
	case "u", "up":
		return Up, true
	case "d", "down":
		return Down, true
	}
	return 0, false
}

// DirectionName returns the canonical lowercase name for d.
func DirectionName(d Direction) string {
	switch d {
	case North:
		return "north"
	case South:
		return "south"
	case East:
		return "east"
	case West:
		return "west"
	case Up:
		return "up"
	case Down:
		return "down"
	}
	return "?"
}

// Opposite returns the reverse direction.
func Opposite(d Direction) Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	case Up:
		return Down
	case Down:
		return Up
	}
	return d
}

// FindItemInRoom returns the first item in room whose name or any
// alias equals needle (assumed already lowercased and trimmed),
// or InvalidItem.
func FindItemInRoom(w *World, room RoomID, needle string) ItemID {
	for index := range w.Items.Name {
		id := ItemID(index)
		if w.Items.Location[id] != room {
			continue
		}
		if itemMatches(w, id, needle) {
			return id
		}
	}
	return InvalidItem
}

// FindItemInInventory returns the first carried item matching
// needle, or InvalidItem.
func FindItemInInventory(w *World, needle string) ItemID {
	return FindItemInRoom(w, InventoryRoom, needle)
}

func itemMatches(w *World, id ItemID, needle string) bool {
	if w.Items.Name[id] == needle {
		return true
	}
	for _, alias := range w.Items.Aliases[id] {
		if alias == needle {
			return true
		}
	}
	return false
}

// ExitFrom returns the exit out of room in direction dir, along
// with its index in w.Exits so callers can flip Locked in place.
func ExitFrom(w *World, room RoomID, dir Direction) (Exit, int, bool) {
	for index, exit := range w.Exits {
		if exit.From == room && exit.Dir == dir {
			return exit, index, true
		}
	}
	return Exit{}, -1, false
}

// ExitsFrom returns every exit out of room.
func ExitsFrom(w *World, room RoomID) []Exit {
	var out []Exit
	for _, exit := range w.Exits {
		if exit.From == room {
			out = append(out, exit)
		}
	}
	return out
}

// IsDark reports whether room is dark and the player carries no
// item flagged as a light source.
func IsDark(w *World, room RoomID) bool {
	if !w.Rooms.Dark[room] {
		return false
	}
	for index := range w.Items.Name {
		id := ItemID(index)
		if w.Items.Location[id] != InventoryRoom {
			continue
		}
		if w.Items.Lit[id] {
			return false
		}
	}
	return true
}

// IsCarrying returns true if the named item is in the player's
// inventory.
func IsCarrying(w *World, needle string) bool {
	return FindItemInInventory(w, needle) != InvalidItem
}
