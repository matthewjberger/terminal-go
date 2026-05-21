package world

// NewDemo builds the bundled adventure and returns a fresh World
// positioned in the starting room.
func NewDemo() *World {
	w := &World{}

	foyer := addRoom(w, "Foyer",
		"A vaulted entry hall floored in cracked marble. Dust drifts in a single shaft of light from somewhere far above. A wide stone stair descends through a hole in the floor.",
		false)
	library := addRoom(w, "Library",
		"Floor to ceiling shelves, mostly empty. A few warped books lean against each other in the corner. A reading desk squats in the middle of the room.",
		false)
	kitchen := addRoom(w, "Kitchen",
		"A long flat-stone counter runs along one wall. The hearth is cold and lined with ash. Bare hooks dangle where utensils used to hang.",
		false)
	cellar := addRoom(w, "Cellar",
		"Roughly hewn walls and a packed earth floor. The air smells of cold iron and damp stone.",
		true)

	brassKey := addItem(w, itemSpec{
		Name:        "brass key",
		Aliases:     []string{"key", "brass"},
		Description: "A small brass key, worn smooth by use.",
		Location:    library,
		Takeable:    true,
	})
	addItem(w, itemSpec{
		Name:        "lantern",
		Aliases:     []string{"lamp", "light"},
		Description: "An oil lantern, its wick burning steady and low.",
		Location:    kitchen,
		Takeable:    true,
	})
	goldCoin := addItem(w, itemSpec{
		Name:        "gold coin",
		Aliases:     []string{"coin", "gold"},
		Description: "A single gold coin. Heavier than it looks.",
		Location:    cellar,
		Takeable:    true,
	})
	addItem(w, itemSpec{
		Name:        "plaque",
		Aliases:     []string{"sign", "bronze plaque"},
		Description: "A worn bronze plaque set into the wall beside the stair.",
		Location:    foyer,
		Readable:    true,
		ReadText:    "Take what you came for and walk back out the way you came.",
	})

	addExitPair(w, foyer, North, library)
	addExitPair(w, foyer, West, kitchen)
	addExit(w, foyer, Down, cellar, Exit{Locked: true, KeyItem: brassKey})
	addExit(w, cellar, Up, foyer, Exit{KeyItem: InvalidItem})

	w.Player.Room = foyer
	w.GoalRoom = foyer
	w.GoalItem = goldCoin
	return w
}

func addRoom(w *World, name, description string, dark bool) RoomID {
	id := RoomID(len(w.Rooms.Name))
	w.Rooms.Name = append(w.Rooms.Name, name)
	w.Rooms.Description = append(w.Rooms.Description, description)
	w.Rooms.Dark = append(w.Rooms.Dark, dark)
	w.Rooms.Visited = append(w.Rooms.Visited, false)
	return id
}

type itemSpec struct {
	Name        string
	Aliases     []string
	Description string
	Location    RoomID
	Takeable    bool
	Readable    bool
	ReadText    string
}

func addItem(w *World, spec itemSpec) ItemID {
	id := ItemID(len(w.Items.Name))
	w.Items.Name = append(w.Items.Name, spec.Name)
	w.Items.Description = append(w.Items.Description, spec.Description)
	w.Items.Aliases = append(w.Items.Aliases, spec.Aliases)
	w.Items.Location = append(w.Items.Location, spec.Location)
	w.Items.Takeable = append(w.Items.Takeable, spec.Takeable)
	w.Items.Readable = append(w.Items.Readable, spec.Readable)
	w.Items.ReadText = append(w.Items.ReadText, spec.ReadText)
	return id
}

func addExit(w *World, from RoomID, dir Direction, to RoomID, modifiers Exit) {
	modifiers.From = from
	modifiers.Dir = dir
	modifiers.To = to
	if !modifiers.Locked {
		modifiers.KeyItem = InvalidItem
	}
	w.Exits = append(w.Exits, modifiers)
}

func addExitPair(w *World, a RoomID, dir Direction, b RoomID) {
	addExit(w, a, dir, b, Exit{})
	addExit(w, b, Opposite(dir), a, Exit{})
}
