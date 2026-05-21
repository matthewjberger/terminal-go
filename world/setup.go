package world

func NewDemo() *World {
	w := &World{Version: SaveVersion}

	entryway := addRoom(w, "Entryway",
		"A high entry hall paved in chipped slate. A shaft of light filters down through a cracked skylight far above. A wide arch opens to the north.",
		0)
	hallway := addRoom(w, "Hallway",
		"A long corridor lined with faded portraits. Doors stand to the east and west, and a heavier door is set into the north wall.",
		0)
	library := addRoom(w, "Library",
		"Floor to ceiling shelves lean inward over a single reading desk. The air is thick with old paper.",
		0)
	kitchen := addRoom(w, "Kitchen",
		"A long flat-stone counter runs along one wall. The hearth is cold and lined with ash. Bare hooks dangle where utensils used to hang.",
		0)
	study := addRoom(w, "Study",
		"A panelled study dominated by a heavy oak desk. A square trapdoor is set into the floorboards beside it.",
		0)
	cellar := addRoom(w, "Cellar",
		"Roughly hewn walls and a packed earth floor. The air smells of cold iron and damp stone.",
		RoomDark)

	addItem(w, itemSpec{
		Name:        "plaque",
		Aliases:     []string{"sign", "wooden plaque"},
		Description: "A weathered wooden plaque, fixed beside the arch.",
		Location:    entryway,
		Tags:        ItemReadable,
		ReadText:    "Bring back what shines, from the dark below.",
	})
	brassKey := addItem(w, itemSpec{
		Name:        "brass key",
		Aliases:     []string{"key", "brass"},
		Description: "A small brass key, worn smooth by use.",
		Location:    library,
		Tags:        ItemTakeable,
	})
	addItem(w, itemSpec{
		Name:        "torn note",
		Aliases:     []string{"note", "paper"},
		Description: "A scrap of paper, torn at one edge.",
		Location:    library,
		Tags:        ItemTakeable | ItemReadable,
		ReadText:    "The door north of the hall takes brass.",
	})
	addItem(w, itemSpec{
		Name:        "lantern",
		Aliases:     []string{"lamp", "light"},
		Description: "An oil lantern, its wick burning steady and low.",
		Location:    kitchen,
		Tags:        ItemTakeable | ItemLit,
	})
	addItem(w, itemSpec{
		Name:        "silver candle",
		Aliases:     []string{"candle", "silver"},
		Description: "A short silver candle, half burnt, smelling faintly of beeswax.",
		Location:    study,
		Tags:        ItemTakeable | ItemLit,
	})
	goldCoin := addItem(w, itemSpec{
		Name:        "gold coin",
		Aliases:     []string{"coin", "gold"},
		Description: "A single gold coin. Heavier than it looks.",
		Location:    cellar,
		Tags:        ItemTakeable,
	})

	addExitPair(w, entryway, North, hallway)
	addExitPair(w, hallway, East, library)
	addExitPair(w, hallway, West, kitchen)
	addLockedExit(w, hallway, North, study, brassKey)
	addExit(w, study, South, hallway)
	addExitPair(w, study, Down, cellar)

	w.PlayerRoom = entryway
	w.GoalRoom = entryway
	w.GoalItem = goldCoin
	return w
}

func addRoom(w *World, name, description string, tags RoomTag) RoomID {
	id := RoomID(len(w.Rooms.Name))
	w.Rooms.Name = append(w.Rooms.Name, name)
	w.Rooms.Description = append(w.Rooms.Description, description)
	w.Rooms.Tags = append(w.Rooms.Tags, tags)
	return id
}

type itemSpec struct {
	Name        string
	Aliases     []string
	Description string
	Location    RoomID
	Tags        ItemTag
	ReadText    string
}

func addItem(w *World, spec itemSpec) ItemID {
	id := ItemID(len(w.Items.Name))
	w.Items.Name = append(w.Items.Name, spec.Name)
	w.Items.Description = append(w.Items.Description, spec.Description)
	w.Items.Location = append(w.Items.Location, spec.Location)
	w.Items.Tags = append(w.Items.Tags, spec.Tags)
	w.Items.ReadText = append(w.Items.ReadText, spec.ReadText)

	if len(w.Items.AliasStart) == 0 {
		w.Items.AliasStart = append(w.Items.AliasStart, 0)
	}
	w.Items.AliasFlat = append(w.Items.AliasFlat, spec.Aliases...)
	w.Items.AliasStart = append(w.Items.AliasStart, int32(len(w.Items.AliasFlat)))
	return id
}

func addExit(w *World, from RoomID, dir Direction, to RoomID) {
	w.Exits = append(w.Exits, Exit{
		From:    from,
		Dir:     dir,
		To:      to,
		KeyItem: InvalidItem,
	})
}

func addLockedExit(w *World, from RoomID, dir Direction, to RoomID, key ItemID) {
	w.Exits = append(w.Exits, Exit{
		From:    from,
		Dir:     dir,
		To:      to,
		Locked:  true,
		KeyItem: key,
	})
}

func addExitPair(w *World, a RoomID, dir Direction, b RoomID) {
	addExit(w, a, dir, b)
	addExit(w, b, Opposite(dir), a)
}
