# Jurassic Balloon Term

A terminal-native version of Jurassic Balloon, built in Go with
[tcell](https://github.com/gdamore/tcell).

This repository currently contains the initial playable workspace: a tcell
screen, a fixed-timestep game loop, queued terminal events, resize handling,
and a small model that can be exercised independently of the terminal. It also
includes a terminal-independent sprite engine with animation, velocity,
layering, transparent cells, and cell-accurate collision detection.

## Run it

```sh
go run ./cmd/jurassic-balloon-term
```

Use the arrow keys or `WASD` to move the balloon. Press `Space` to pause and
`Q` or `Esc` to quit.

## Development

```sh
make test
make build
```

The project borrows its initial loop architecture from
[Escape from Proxima 5](https://github.com/gdamore/proxima5): input is queued,
simulation updates happen on a fixed cadence, and rendering is kept separate
from game state. Jurassic Balloon Term uses current tcell v2 APIs and keeps
the model testable without a real terminal.

### About Proxima 5

Proxima 5 was the first game Garrett D'Amore created, apart from a number
guessing game he wrote as a child around 1980. It was built mostly over a
single weekend, after he realized that tcell—a framework he had created for a
business application—could also power a terminal game.

That history makes Proxima 5 an important source of ideas and working game
mechanics, but not a design blueprint. Some of its choices reflect the speed
and exploratory nature of that first game. Jurassic Balloon Term should reuse
the good concepts while reconsidering APIs, ownership, timing, testing, and
package boundaries for a maintainable modern implementation.

## Layout

- `cmd/jurassic-balloon-term`: executable entry point
- `internal/game`: terminal-independent game state and updates
- `internal/sprite`: animated sprites, scenes, layers, and collisions
- `internal/ui`: tcell event loop and renderer

Licensed under the Apache License, Version 2.0.

## Design documents

- [Jurassic Balloon game-design handoff](docs/JURASSIC_BALLOON_DESIGN_HANDOFF.md)
- [Source discussion transcript](docs/JURASSIC_BALLOON_SOURCE_TRANSCRIPT.md)

