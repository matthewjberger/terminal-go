// Package world holds the game state and the queries over it.
package world

// RoomID indexes Rooms.* slices. The sentinel InventoryRoom is
// used as an item location to mean "carried by the player".
type RoomID int32

// ItemID indexes Items.* slices. InvalidItem is the "no item"
// sentinel.
type ItemID int32

const (
	InventoryRoom RoomID = -2

	InvalidItem ItemID = -1
)

type Direction int8

const (
	North Direction = iota
	South
	East
	West
	Up
	Down
)

type RoomTag uint8

const (
	RoomDark RoomTag = 1 << iota
)

type ItemTag uint16

const (
	ItemTakeable ItemTag = 1 << iota
	ItemReadable
	ItemLit
)

type Rooms struct {
	Name        []string
	Description []string
	Tags        []RoomTag
}

// Items holds the item table. AliasStart has length len(Name)+1
// so item id's aliases are AliasFlat[AliasStart[id]:AliasStart[id+1]].
type Items struct {
	Name        []string
	Description []string
	Location    []RoomID
	Tags        []ItemTag
	ReadText    []string

	AliasFlat  []string
	AliasStart []int32
}

type Exit struct {
	From    RoomID
	Dir     Direction
	To      RoomID
	Locked  bool
	KeyItem ItemID
}

type Player struct {
	Room RoomID
}

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

func RoomHasTag(w *World, room RoomID, tag RoomTag) bool {
	if room < 0 || int(room) >= len(w.Rooms.Tags) {
		return false
	}
	return w.Rooms.Tags[room]&tag != 0
}

func ItemHasTag(w *World, item ItemID, tag ItemTag) bool {
	if item < 0 || int(item) >= len(w.Items.Tags) {
		return false
	}
	return w.Items.Tags[item]&tag != 0
}

func ItemAliases(w *World, item ItemID) []string {
	if item < 0 || int(item)+1 >= len(w.Items.AliasStart) {
		return nil
	}
	return w.Items.AliasFlat[w.Items.AliasStart[item]:w.Items.AliasStart[item+1]]
}

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

// FindItemInRoom returns the first item in room matching needle by
// name or alias, or InvalidItem. needle is assumed lowercased.
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

func FindItemInInventory(w *World, needle string) ItemID {
	return FindItemInRoom(w, InventoryRoom, needle)
}

func itemMatches(w *World, id ItemID, needle string) bool {
	if w.Items.Name[id] == needle {
		return true
	}
	for _, alias := range ItemAliases(w, id) {
		if alias == needle {
			return true
		}
	}
	return false
}

// ExitFrom returns the exit out of room in dir along with its
// index in w.Exits so callers can flip Locked in place.
func ExitFrom(w *World, room RoomID, dir Direction) (Exit, int, bool) {
	for index, exit := range w.Exits {
		if exit.From == room && exit.Dir == dir {
			return exit, index, true
		}
	}
	return Exit{}, -1, false
}

func ExitsFrom(w *World, room RoomID) []Exit {
	var out []Exit
	for _, exit := range w.Exits {
		if exit.From == room {
			out = append(out, exit)
		}
	}
	return out
}

// IsDark assumes room is the player's current room; the inventory
// check is meaningless otherwise.
func IsDark(w *World, room RoomID) bool {
	if !RoomHasTag(w, room, RoomDark) {
		return false
	}
	for index := range w.Items.Name {
		id := ItemID(index)
		loc := w.Items.Location[id]
		if loc != InventoryRoom && loc != room {
			continue
		}
		if ItemHasTag(w, id, ItemLit) {
			return false
		}
	}
	return true
}

func IsCarrying(w *World, needle string) bool {
	return FindItemInInventory(w, needle) != InvalidItem
}
