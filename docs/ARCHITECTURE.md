# Architecture

terminal-go is a tiny data-oriented text-adventure runtime. The
game state is one struct of parallel arrays; every verb is a free
function that mutates that struct and writes to an `io.Writer`.

## Principles

- **Data, not objects.** The world is one [`world.World`] value.
  Typed `RoomID` and `ItemID` index struct-of-array tables.
  Mutation lives in free functions.
- **Tags, not columns.** Binary attributes go into a `[]RoomTag`
  / `[]ItemTag` bitset, not a new bool slice per flag. Adding an
  attribute is one bit constant.
- **Flat aliases.** Item aliases are a flat `[]string` plus a
  `[]int32` start table. No per-item slice allocation.
- **Free-function commands.** Each verb is
  `func(*world.World, ..., io.Writer)`. The dispatcher in
  `command.Run` is a `switch`. No `Command` interface.
- **One allocation, save and load.** `World` is plain data with
  no func or interface fields. `world.Encode` / `world.Decode`
  is one `gob` call.

## Module layout

| Package      | Purpose                                                |
|--------------|--------------------------------------------------------|
| `world/`     | Rooms / Items / Exits / Player tables, ID and tag types, queries, save+load, demo seeder. |
| `parse/`     | Tokenizer that produces a normalized verb + args. |
| `command/`   | Verb handlers and the dispatcher. |
| `cmd/zork/`  | Entry point: build the world, drive the input loop. |

## Tables

`Rooms` and `Items` are struct-of-arrays. A `RoomID` is an `int32`
that indexes every `Rooms.*` slice; `ItemID` indexes every
`Items.*` slice. `Items.Location[id]` is the `RoomID` the item
sits in; the `InventoryRoom` sentinel means "carried by the
player".

`Items.AliasFlat` is a single string slice with every alias
concatenated. `Items.AliasStart` has length `len(Items.Name)+1`,
so item `id`'s aliases are
`AliasFlat[AliasStart[id]:AliasStart[id+1]]`. The `+1` element is
the running total; `addItem` maintains it.

Exits are a flat `[]Exit` scanned linearly by `(from, dir)`. The
exit count is small enough that a linear scan beats a map every
time and keeps `gob` save/load trivial.

## Tags

`RoomTag` is `uint8`; `ItemTag` is `uint16`. Each constant is a
single bit: `RoomDark`, `ItemTakeable`, `ItemReadable`,
`ItemLit`. Combine with `|`, test with `world.RoomHasTag` /
`world.ItemHasTag`. Both helpers bound-check the index and
return false for sentinel or out-of-range IDs.

To add a new flag: add a `1 << iota` constant in the tag block,
set it on items that need it in `setup.go`, query it where it
matters. No new column on the table.

## Frame

Each input line:

1. Read a line from stdin.
2. `parse.Tokenize(line)` returns `{Verb, Args, Object}`.
3. `command.Run(world, token, out)` dispatches, mutates `world`,
   writes the response.
4. Loop until `world.Quit` or `world.Won`.

There is no game loop, no delta time. The world only advances
when the player types a verb.

## The single-appender rule

`addItem` is the only function that appends into the `Items`
table. It grows every column (`Name`, `Description`, `Location`,
`Tags`, `ReadText`, `AliasFlat`, `AliasStart`) in one pass.
`addRoom` does the same for `Rooms`. Outside those two functions,
the SoA invariant only holds because nothing else appends. Break
that and `w.Items.Name[id]` and `w.Items.Tags[id]` drift apart.

## Locked exits

`addLockedExit(from, dir, to, key)` is the only path to a locked
edge. `addExit(from, dir, to)` and `addExitPair(a, dir, b)` build
open edges and set `KeyItem: InvalidItem`. There is no way to
express a locked exit without naming its key item, which closes
the zero-value-collision footgun where `ItemID(0)` would
accidentally become a valid key.

## Light

A room with `RoomDark` is dark unless `world.IsDark(w, room)`
finds an item with `ItemLit` either in the player's inventory or
in the room itself. The lantern and silver candle both light the
cellar; either path through the demo works.

## Save and load

`world.Encode(w, io.Writer)` and `world.Decode(io.Reader)` use
`encoding/gob` over the whole `World` value. Every field is
exported and plain data (slices of primitives, slices of structs
of primitives), so the round-trip is one call each way. The
`save` and `load` verbs wrap a file open around those two
functions.
