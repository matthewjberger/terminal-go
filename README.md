# terminal-go

A small text adventure in Go, written data-oriented from the
ground up. The whole game state is one struct of parallel arrays
indexed by typed IDs. Verbs are free functions that mutate that
struct and write to an `io.Writer`. Save and load are a single
`encoding/gob` call each way.

About 840 lines of Go across four packages. Six rooms, six
items, a dozen verbs, one lit lantern.

Architecture notes live in
[`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Quickstart

```
just run        # zork
just run zork   # same
```

`just --list` for the rest.

## License

Dual-licensed under [MIT](LICENSE-MIT) or
[Apache-2.0](LICENSE-APACHE) at your option.
