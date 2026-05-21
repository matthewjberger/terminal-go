# Architecture

terminal-go is a tiny data-oriented text-adventure runtime. The
game state is one struct of parallel arrays keyed by typed IDs;
every verb is a free function that mutates that struct and writes
output to an `io.Writer`. No interfaces, no method dispatch on
game objects, no per-room or per-item types.

## Principles

- **Data, not objects.** The world is one [`world.World`] value
  with struct-of-arrays tables (`Rooms`, `Items`, plus a flat
  `Exits` slice). Typed IDs (`RoomID`, `ItemID`) index those
  tables. Mutation lives in free functions that take `*World`.
- **Free-function commands.** Each verb (`look`, `go`, `take`,
  `examine`, ...) is a `func(*world.World, ..., io.Writer)`. The
  dispatcher in `command` is a `switch` on the parsed verb. There
  is no `Command` interface and no command struct.
- **One world, one allocation.** All game state lives in a single
  `World` value. Save and load is one `encoding/gob` round-trip
  away if you want it.
- **Content is data.** Rooms, items, exits, and the win goal are
  populated by `world.NewDemo` appending to slices. Authoring a
  new adventure means writing a new `NewDemo` style seeder, not
  subclassing anything.

## Module layout

| Package      | Purpose                                                |
|--------------|--------------------------------------------------------|
| `world/`     | Rooms / Items / Exits / Player tables, ID types, direction helpers, demo seeder. |
| `parse/`     | Tokenizer that turns a raw input line into a normalized verb plus arguments. |
| `command/`   | Free-function verb handlers and the dispatcher. |
| `cmd/zork/`  | Entry point: build the world, drive the input loop, print output. |

No package depends on a higher layer. `command` depends on
`world` and `parse`; `cmd/zork` ties them together.

## Tables

Rooms and items are stored as struct-of-arrays. A `RoomID` is a
typed `int32` that indexes every `Rooms.*` slice; an `ItemID`
indexes every `Items.*` slice. `Items.Location[id]` holds the
`RoomID` the item currently sits in, with the sentinel
`InventoryRoom` meaning "carried by the player".

Exits are a flat `[]Exit` scanned linearly by `(from, dir)`. The
exit count is small enough that a linear scan beats a map every
time, and a flat slice keeps save/load trivial.

## Frame

Each input line:

1. Read a line from stdin.
2. `parse.Tokenize(line)` returns the normalized verb, the raw
   args (post filler stripping), and the joined object string.
3. `command.Run(world, token, out)` dispatches on the verb,
   mutates `world`, and writes the player-facing response.
4. Loop until `world.Quit` or `world.Won`.

There is no game loop, no delta time, no schedule. The world only
advances when the player types a verb.

## Dark rooms and locked exits

Two cross-cutting state checks live in `world` as pure functions:

- `IsDark(w, room)` returns true when `Rooms.Dark[room]` is set
  and no item the player is carrying is configured as a light
  source (today: any item named `lantern`).
- `ExitFrom(w, room, dir)` returns the exit plus its slice index
  so the mover can flip `Locked` to false in place once the key
  is presented.

The locked exit holds an `ItemID` for the key; if the player is
carrying that item, the mover clears `Locked` on the exit (and
its reverse, if any) and lets the move through.

## Why this shape

A text adventure is the smallest non-trivial DOD playground:

- Iteration order is one verb per tick, so cache layout doesn't
  matter, but the discipline of "indices, not pointers; tables,
  not objects" does. The same code shape extends to thousands of
  entities without restructuring.
- Save/load is `gob.Encode(world)` because the `World` is plain
  data with no `interface{}` field, no `func` field, and no
  pointer graph.
- The dispatcher pattern (one switch in `command.Run`) means
  adding a verb is one case plus one handler, never touching a
  type hierarchy.
