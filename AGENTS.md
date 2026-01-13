# AGENTS.md

This file guides agentic coding assistants working in the `pet-everyone` repository. Follow these instructions whenever you touch code in this repo.

## Scope
- This file applies to the entire repository unless a more specific AGENTS file is added in a subdirectory.
- When editing files in subdirectories, re-check for nested AGENTS.md files and obey the most specific scope.

## Environment & Services
- Go modules; Go toolchain required.
- Dependencies: PostgreSQL (mapped 5433->5432), Redis (6379), WAL files on local disk.
- Bring up services: `docker-compose up -d`.
- Required env vars (see `.env`): `FILEPATH_ROOT`, `ASSETS_ROOT`, `PORT`, DB creds, Redis creds.
- Main entry: `go run cmd/web/main.go` (expects env + running DB/Redis).

## Build / Test / Lint Commands
- Run all tests: `go test ./...`
- Verbose tests: `go test -v ./...`
- Single package: `go test ./internal/service`
- Single test: `go test -run TestName ./path/to/package`
- Race check (targeted): `go test -race ./internal/...` when feasible
- Format: `go fmt ./...`
- Lint (Makefile pr-check): `make pr-check` (runs go fmt, go vet, staticcheck)
- Staticcheck standalone: `staticcheck ./...`
- Run app locally: `go run cmd/web/main.go`

## Data & Migrations
- Migrations live at `internal/data/db/migrations/`; apply via your preferred tool (e.g., `golang-migrate`) when needed.
- WAL artifacts should **not** be committed; `**/*.log` and `**/*.wal` are ignored in `.gitignore`.

## Code Style (Go)
- Use `gofmt`/`goimports` on all Go files; keep imports grouped standard / third-party / local.
- Prefer explicit types; avoid `interface{}` unless required.
- Naming: exported identifiers are PascalCase with correct acronyms (`ID`, `URL`, `HTTP`); unexported in camelCase.
- Package names: short, lower, no underscores.
- Keep functions small; favor clear, testable units.

## Error Handling
- Do not `log.Fatal` or `panic` outside `main`/startup; return errors up the stack.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`; use `errors.Is/As` to test.
- Handle `sql.ErrNoRows` explicitly; avoid treating it as fatal.
- When a failure is non-critical, log with context and continue; avoid silent drops unless explicitly documented.

## Logging
- Use structured/contextual logging where available (`app.Logger()`); otherwise `log.Printf` with clear prefixes and identifiers (petID, userID/guestID).
- Avoid logging sensitive data (tokens, passwords, session IDs).

## Concurrency & Contexts
- Always accept `context.Context` in long-running or I/O-bound functions; respect cancellation.
- For goroutines, ensure they exit on context cancel; avoid leaks.
- Protect shared maps/slices with mutexes; prefer RWMutex for read-heavy paths.
- Use buffered channels where backpressure is expected; avoid unbounded growth.
- Clean up timers/tickers (`defer ticker.Stop()`), close channels on shutdown when safe.

## WebSocket Layer
- Upgrader should enforce origin checks (localhost currently) and avoid querystring tokens.
- Authenticate via cookies/session or first-message protocol; do not trust client-provided IDs without verification.
- Hubs are per-pet; registry lazily creates and tears down with context cancellation.
- Ensure unregister paths cannot deadlock; avoid leaked clients by removing unresponsive connections.

## HTTP Handlers & Middleware
- Use existing middleware for auth/logging; keep handlers thin and delegate to services/models.
- Respond via helper methods (`RespondWithError`, etc.) for consistent error shapes.
- Validate inputs; reject missing/invalid IDs with 400s; avoid accepting tokens in URLs.

## Database Access
- Use prepared statements with placeholders; never interpolate SQL.
- Use `COALESCE` to avoid NULL scan failures for aggregates (already used in pet counts).
- Close rows promptly (`defer rows.Close()`); check `rows.Err()` where needed.
- Keep transactions short; handle rollback on errors.

## Redis
- Connect via `internal/data/cache`; validate connectivity with `Ping` at startup; close clients on shutdown.
- Handle Redis unavailability gracefully; prefer fallback paths over crashing.

## WAL (Write-Ahead Log)
- WAL files are local durability artifacts; do not commit.
- Ensure WAL is closed on shutdown; surface failures where recovery is impacted.

## DTOs & Serialization
- DTOs live in `internal/dto`; keep them minimal and typed (no `map[string]interface{}`).
- For JSON, ensure exported fields are properly cased; add tags when needed.

## Testing Guidance
- Favor unit tests with mocks (see `internal/service/mock.go` and interfaces like `PetDatabase`).
- Use `t.Run` subtests for table-driven cases.
- Clean up artifacts with `t.TempDir()`; avoid writing to repo root.
- For WAL tests, ensure files are removed/created under temp dirs.
- Avoid hitting live DB/Redis in unit tests; use fakes/mocks when possible.

## Frontend Templates/Assets
- Templates under `app/templates`; CSS in `app/css`; keep HTML semantic; avoid inline styles when possible.
- When adjusting UI, keep auth forms consistent and responsive; check both login/signup.

## Architecture Pointers
- Registry: `internal/registry` manages per-pet WebSocket hubs with lazy create/teardown.
- Services: `internal/service` (PetService buffers counts, flushes to DB/WAL; ChatService builds chat messages).
- Transport: `internal/transport` routes commands/events; serializer broadcasts JSON.
- Websocket clients: `internal/websocket` manages run loops and broadcast/unregister logic.
- Data layer: `internal/data/model` for DB access; `internal/data/cache` for Redis; migrations in `internal/data/db/migrations`.

## Security Guidelines
- Do not pass auth tokens in URLs; prefer HttpOnly cookies or authenticated WS subprotocols/first-message auth.
- Add CORS/origin checks for WebSockets; restrict to allowed origins (localhost currently).
- Validate user/guest identities server-side; never trust client-supplied IDs.
- Avoid logging secrets or tokens; scrub sensitive fields in errors.

## Performance & Reliability
- Avoid N+1 queries; batch where feasible (e.g., pet counts).
- Guard against unbounded retry queues; surface backpressure/metrics instead of silent drops.
- Use timeouts for external calls (DB, Redis, HTTP) via context.

## Naming & Comments
- Comment exported types/functions with GoDoc style; keep comments concise and intent-focused.
- Use clear, descriptive variable names; avoid one-letter names except loop indices.

## When in Doubt
- Follow Go idioms; keep changes minimal and scoped.
- Prefer clarity over cleverness; choose the simplest thing that works.
- Ask before introducing new dependencies or architectural patterns.

## Quick Checklist Before Commit
- `go fmt ./...`
- `go vet ./...`
- `staticcheck ./...`
- `go test ./...` or targeted `go test -run TestName ./path`
- No credentials/tokens in code, logs, or URLs
- No WAL/log artifacts or assets committed

## Frontend JS Guidelines
- Keep client JS small and event-driven; avoid global state where possible.
- Prefer `fetch` with JSON helpers; centralize API calls if you add more.
- Avoid inlining secrets or tokens; rely on HttpOnly cookies for auth.
- Ensure WebSocket clients send auth via cookies or first-message, not query params.

## Docker / Local Dev
- Use `docker-compose up -d` to start Postgres and Redis before running the app.
- Use `docker-compose logs -f` to debug service startup issues.
- Reset state by stopping containers and removing volumes only if necessary.

## Git Hygiene
- Do not commit generated assets, WAL/log files, or local DB files.
- Keep commits scoped and conventional (`feat:`, `fix:`, `refactor:`, etc.).
- Avoid force pushes to shared branches unless explicitly requested.

## Skills
- Additional internal scripts are documented under `skills/`; prefer using documented commands there (e.g., DSN generator) when relevant.

## Cursor / Copilot Rules
- No additional Cursor or Copilot repo rules are present; follow this AGENTS file for guidance.

## Contact
- If instructions conflict, prioritize more specific AGENTS files (if added later), then this file, then general Go best practices.
