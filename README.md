# terminal-go

A small collection of text-based games in Go. Data-oriented from
the ground up: room, item, and exit tables are flat parallel arrays
keyed by typed IDs, all behavior lives in free functions, no game
state is hidden behind methods.

Architecture notes live in [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Quickstart

```
just run        # zork
just run zork   # same
```

`just --list` for the rest.

## License

Dual-licensed under [MIT](LICENSE-MIT) or [Apache-2.0](LICENSE-APACHE) at your option.
