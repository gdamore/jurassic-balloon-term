# Jurassic Balloon Term

Current version: **0.1**

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

Hold `Space`, `Up`, `W`, or the primary mouse button to fire the burner. Press
`P` to pause. `Q` or `Esc`
returns an active run to the title screen and quits when pressed there. Press
`M` to toggle music, `S` to toggle sound and ambience,
and `-` or `+` to adjust master volume. Tcell v3 uses true press/release events
with Kitty and Win32 keyboard protocols; legacy protocols fall back to repeated
presses with a short grace period. Storm clouds telegraph lightning with rain;
touching a visible bolt pops the balloon. Press `Space`, `R`, or `Enter` to
retry.

Starting from the title first reveals the grounded balloon and a launch prompt;
press and hold Space or the primary mouse button again to lift off.

Press `C` on the title screen for creator, audio-asset, and technology credits.
Full attribution and license details are recorded in
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

## Play in a browser

The same game can be compiled to WebAssembly and hosted as a static site. Build
the complete browser bundle with:

```sh
make web
python3 -m http.server --directory dist/web 8080
```

Then open <http://localhost:8080/>. The first click on the launch card grants
the browser permission to start WebAudio; the in-game title screen follows.
The generated `dist/web` directory includes the Go WebAssembly runtime and
tcell's Ghostty-based browser terminal, and is ignored by Git.

The `pages` GitHub Actions workflow rebuilds and deploys this bundle whenever
changes reach `main`. In the repository's Pages settings, select **GitHub
Actions** as the publishing source if it is not already enabled.

## Development

```sh
make test
make build
```

The project borrows its initial loop architecture from
[Escape from Proxima 5](https://github.com/gdamore/proxima5): input is queued,
simulation updates happen on a fixed cadence, and rendering is kept separate
from game state. Jurassic Balloon Term uses current tcell v3 APIs and keeps
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

