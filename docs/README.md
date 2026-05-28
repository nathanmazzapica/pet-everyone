# Pet Everyone Docs

This folder is the project memory for Pet Everyone. Keep docs short, current, and tied to code changes.

## Start Here

- [Architecture overview](architecture.md)
- [Image processor API contract](contracts/image-processor-api.md)
- [Realtime protocol contract](contracts/realtime-protocol.md)

## Decision Records

Decision records explain why the project chose a direction. Add one when a change affects service boundaries, APIs, data ownership, auth, storage, or realtime behavior.

- [0001: Recover the existing Go server](decisions/0001-recover-existing-go-server.md)
- [0002: Use a separate image processor service](decisions/0002-image-processor-service.md)
- [0003: Split realtime responsibilities between WebSocket and SSE](decisions/0003-realtime-ws-and-sse.md)

## Drift Rule

Update docs when a change affects:

- an HTTP endpoint or payload
- a WebSocket or SSE event
- a database table or ownership boundary
- an auth/session flow
- an external service contract
- local setup or deployment assumptions

Prefer a short accurate doc over a long aspirational one.
