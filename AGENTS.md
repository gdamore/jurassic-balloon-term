# Project guidance

Jurassic Balloon Term is a Go/tcell terminal game.

Read `docs/JURASSIC_BALLOON_DESIGN_HANDOFF.md` before making gameplay or scope
decisions. It distinguishes confirmed decisions from tentative and superseded
ideas. Consult the linked source transcript when the summary is ambiguous.

- Keep game rules and state in `internal/game`; they must not depend on tcell.
- Keep reusable animation and collision behavior in `internal/sprite`.
- Keep terminal input and rendering in `internal/ui`.
- Preserve the fixed-timestep update model used by the initial app loop.
- Add focused unit tests for game logic.
- Run `go test ./...` and `go vet ./...` before committing.

`proxima5` was Garrett's first substantial game and was created mostly over a
weekend as an exploration of using tcell for games. It is a source of useful
ideas and behavior, not an architectural blueprint or a source tree to copy
wholesale. Evaluate its design choices before carrying them forward. Prefer
current tcell v2 APIs, deterministic timing, tests, and small explicit
components.

