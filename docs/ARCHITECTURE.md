# Architecture

terminal-go is a small text adventure. Six rooms, six items, a
player, and around a dozen verbs, in about 840 lines of Go
across four packages. The point of the project is the shape of
the code, not the gameplay. Everything is data-oriented, and the
whole game state is one struct of parallel arrays that
round-trips through `encoding/gob` in a single call.

## Why one struct of arrays

A more conventional Go layout would give each room and item
their own struct, each with methods. `Room` would have a
`Describe`, `Item` would have a `Take`, the player would hold a
`*Room` pointer, and the world would be a graph of values that
own each other and dispatch through methods. Serializing it
would mean walking that graph and rebuilding the pointers on
load. Adding a new "is the item lit" flag would mean a field on
`Item`, a getter, possibly a setter, and a thought about
encapsulation.

terminal-go flips the layout. Rooms live in a `Rooms` struct
whose fields are parallel slices: `Name []string`,
`Description []string`, `Tags []RoomTag`. Items live in an
`Items` struct with the same shape. A `RoomID` is an `int32`
index into the room slices, and an `ItemID` indexes the item
slices. The player's location is a `RoomID`, not a pointer.
"Item in room" is `w.Items.Location[id]`. Nothing in the world
points at anything else; everything indexes.

At six items the cache argument doesn't pay for itself. The
argument that does pay is that the whole `World` is plain data.
No func fields, no interface fields, no pointers between
entities. `gob` encodes and decodes the value in one call each
way, with no custom marshalers and no rebuilt graphs. Saving the
game is `gob.NewEncoder(file).Encode(w)`. Loading is the inverse
plus a version check. That falls out of the layout, not because
we added it.

## How a table grows

Parallel slices only work if every column grows together.
`Items.Name[3]` and `Items.Tags[3]` describe the same item, and
the moment they don't, the world is broken.

There is exactly one function that appends into `Items`:
`addItem`. It grows every column in one pass: `Name`,
`Description`, `Location`, `Tags`, `ReadText`, plus the two
alias columns described below. `addRoom` does the same for
`Rooms`. Any other code that wanted to append directly to one
column would be a bug, and nothing else does.

## Aliases without per-item allocations

Items have alias lists. The brass key answers to "brass key",
"key", and "brass". The natural Go representation is
`Aliases [][]string`. One inner slice per item, one heap
allocation per item.

Instead, all aliases live in one flat slice with an index table
beside it:

```go
AliasFlat  []string
AliasStart []int32
```

`AliasStart[id]` is where item `id`'s aliases begin and
`AliasStart[id+1]` is where they end. `AliasStart` has length
`len(Items.Name) + 1`, so the lookup
`AliasFlat[AliasStart[id]:AliasStart[id+1]]` works for every
item including the last one without a special case. `addItem`
extends `AliasFlat` and pushes the new running total onto
`AliasStart`. It's the same append-everything-together
discipline that keeps the parallel columns in sync.

## Tags are bits, not columns

A room is dark or it isn't. An item is takeable or not, readable
or not, lit or not. The naive layout is one parallel `[]bool`
per flag, which means every new flag is a new column on the
`Items` struct, a new slice to grow in `addItem`, and a new
field for `gob` to think about. `ItemLit` was added partway
through the project's life, and that cost would compound.

Room flags live in `RoomTag` (a `uint8`) and item flags in
`ItemTag` (a `uint16`). Each flag is `1 << iota`. Combine with
`|`, test with `RoomHasTag` / `ItemHasTag`. Adding a new
attribute is one constant. `addItem` doesn't change, the table
doesn't change, the save format doesn't break.

Two consequences of doing it this way. The tag helpers
bounds-check the index and return false for sentinels and
out-of-range IDs, so `ItemHasTag(w, InvalidItem, ItemTakeable)`
is safe to write. And the dark-room rule "any lit item in the
player's inventory or in the room dispels darkness" is one pass
over `Items.Location` checking `ItemLit`. That's `world.IsDark`,
and the lantern and the silver candle both work.

## Verbs are free functions, not methods

A more OOP approach would have a `Command` interface, a registry
of implementations, and a dispatcher that calls `Execute(world)`
on the matched command. terminal-go has none of that. Every verb
is a free function with the signature
`func(w *world.World, tok parse.Token, out io.Writer)`. The
dispatcher in `command.Run` is a single `switch` on `tok.Verb`.
To add a verb you write a function and add a case.

Directions short-circuit before the switch. If `tok.Verb` parses
as a direction the function calls `move` directly. Two-word
verbs like "pick up" and "put down" are normalized to "take" and
"drop" by `parse.normalizeVerb`, so the dispatcher only ever
sees one token.

## Locked exits force their key

`addExit(from, dir, to)` builds an open edge.
`addLockedExit(from, dir, to, key)` builds a locked one and
requires an `ItemID`. There is no constructor that takes a
locked edge without a key, so the state "this exit is locked
and there is no key for it in the world" can't exist.

The zero value of `ItemID` is 0, which is a real item: the
first one added to the table. If `addLockedExit` accepted a
default-zero key, omitting the argument would silently make the
first-defined item the key to every locked door. Forcing the
parameter rules that out at the type level.

## Save and load is one gob call

`world.Encode(w, out)` is `gob.NewEncoder(out).Encode(w)`.
`world.Decode(in)` is the inverse plus a version check. Every
save carries a `Version int32`, and decode rejects anything
whose version doesn't match the runtime's `SaveVersion`. Adding
a column to `Items` or a new tag bit doesn't break old saves
(gob handles the schema), but reordering fields or changing
the type of one does, and that's what the version is for.

The `save` and `load` verbs are file open, encode, close. There
is no save manager.

## The input loop

`cmd/zork/main.go` reads a line from stdin, calls
`parse.Tokenize`, calls `command.Run`, and loops until
`world.Quit` or `world.Won`. The world only advances when the
player types something. There is no clock.
