# 0001: Recover the Existing Go Server

## Status

Accepted

## Context

The project was paused for several months. The current Go server contains working ideas for users, guests, pets, WebSockets, display names, and pet count persistence, but also includes incomplete systems such as Redis-backed leaderboard work, WAL recovery, and the old local Python image processing flow.

The product direction is expanding from browser-first to API-first with future mobile clients.

## Decision

Keep and recover the existing Go server instead of starting over.

Recovery means:

- document current architecture
- preserve the current rate-limiting-era state on a `recovery` branch
- clarify API and service boundaries before more feature work
- make JSON APIs the canonical client contract over time
- keep the browser UI working as one client, not the center of the architecture

## Consequences

- Existing code and product decisions are preserved.
- Recovery work must separate stable systems from unfinished experiments.
- Some packages may remain messy while docs and contracts catch up.
- Large rewrites are deferred unless a boundary is clearly blocking progress.

## Open Questions

- Which stale systems should be deleted versus documented as planned?
- How much of the current OpenAPI spec should be corrected before mobile work begins?
- What minimum API surface does the first mobile client need?
