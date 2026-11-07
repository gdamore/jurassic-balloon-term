# Project guidance

Jurassic Balloon Term is a Go/tcell terminal game.

- Keep game rules and state in `internal/game`; they must not depend on tcell.
- Keep terminal input and rendering in `internal/ui`.
- Preserve the fixed-timestep update model used by the initial app loop.
- Add focused unit tests for game logic.
- Run `go test ./...` and `go vet ./...` before committing.

`proxima5` is an architectural reference, not a source tree to copy wholesale.
Prefer current tcell v2 APIs and small, explicit components.

