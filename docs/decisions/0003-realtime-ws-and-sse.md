# 0003: Split Realtime Responsibilities Between WebSocket and SSE

## Status

Draft accepted direction

## Context

The current app uses WebSockets for pet taps, pet count requests, and chat. The product may add mobile clients and server-driven status updates such as leaderboards and image processing status.

Mobile clients can use WebSockets while active in the foreground, but they should not depend on long-lived sockets while backgrounded. Some updates are naturally one-way from server to client.

## Decision

Use WebSockets for active pet-room interactions and SSE for one-way status streams.

WebSocket responsibilities:

- active pet tapping
- chat
- live count deltas
- room presence, if added

SSE responsibilities:

- leaderboard updates
- image processing status
- other low-frequency server-to-client updates

HTTP remains the transport for one-shot commands and queries.

## Consequences

- The existing per-pet WebSocket hub architecture can be preserved.
- SSE can be introduced in small vertical slices without replacing WebSockets.
- Mobile clients need reconnect handling.
- Background mobile notifications remain a future push notification concern, not a WebSocket/SSE concern.

## Open Questions

- Should leaderboard updates be produced from Postgres, Redis, or a future leaderboard service?
- Should SSE endpoints require session auth, guest auth, or both?
- How should clients resume after missed SSE events?
