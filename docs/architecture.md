# Architecture Overview

## Current Product Shape

Pet Everyone is an online petting zoo. Users can create pets, upload images, open a pet page, pet the image in real time, and chat with other active petters. Guests can participate through a generated guest identity.

The project is moving from a browser-first Go web app toward an API-first backend that can support web, iOS, and Android clients.

## Runtime Components

### Go API / Web Server

Entry point: `cmd/web/main.go`

Owns:

- user accounts and sessions
- guest identities
- pets and pet ownership
- user-visible pet image records
- pet counts
- WebSocket/SSE client-facing realtime APIs
- HTML templates for the current browser UI

Current dependencies:

- PostgreSQL
- Redis connection, mostly unused today
- local filesystem assets
- per-pet WebSocket hubs

### PostgreSQL

Main product database.

Owns product state:

- `RegisteredUser`
- `Visitor`
- `SessionTokens`
- `Pet`
- `PetImage`
- `UserPetsClickCount`
- `DisplayName`

### Image Processor Service

Separate service under development in `pet-everyone-image-processor`.

Expected ownership:

- internal image processing jobs
- worker claims/locks/retries
- processing attempts and operational errors
- intermediate processing artifacts

It should not own users, pets, active pet image selection, or final product-visible image state.

### Browser Client

Current web client is server-rendered HTML plus small JavaScript files under `app/`.

Current browser behavior:

- cookie auth
- pet tapping over WebSocket
- chat over WebSocket
- personal count loaded through HTTP
- hard-coded leaderboard placeholders

### Future Mobile Clients

Mobile clients should consume JSON APIs and realtime protocols directly. They should not depend on server-rendered pages or browser-only cookie behavior.

## Current Package Boundaries

### `cmd/web/application`

Application facade and dependency container. It currently owns JSON helpers, asset path helpers, signup orchestration, pet listing/detail DTO construction, and create-pet image processing.

This package is useful but overloaded.

### `cmd/web/handler`

HTTP routes, page rendering, API handlers, and WebSocket upgrade handlers.

Handlers currently call models through `application.Config` directly in several places. This is acceptable for recovery, but service boundaries should become clearer before the API grows.

### `internal/data/model`

SQL data access layer. Models execute SQL directly.

### `internal/service`

Domain-ish services for realtime pet counts, chat, display-name resolution, and an unfinished leaderboard service.

### `internal/registry`

Manages one lazily-created WebSocket hub per pet.

### `internal/websocket`

Owns WebSocket connection pumps, hub registration/unregistration, and broadcast fanout.

### `internal/transport`

Routes incoming WebSocket commands to services and serializes service events back to hub broadcasts.

### `internal/wal`

Simple local WAL for failed pet count updates. Recovery is not wired into startup yet.

## Important Runtime Flows

### App Startup

1. Load `.env`.
2. Connect to PostgreSQL.
3. Connect to Redis and ping it.
4. Build `application.Config`.
5. Create hub registry.
6. Ensure asset directory exists.
7. Register routes and start HTTP server.

### Pet Page

1. `GET /pet/{pet_id}`.
2. Ensure guest cookie if no valid session.
3. Load pet metadata from `PetModel`.
4. Render `app/templates/pet.html`.
5. Browser JS loads personal count and opens WebSocket.

### WebSocket Pet Room

1. Client opens `/pet/{pet_id}/ws`.
2. Server resolves session user or guest identity.
3. Registry gets or creates hub for pet.
4. Client joins hub.
5. Incoming messages flow through router to pet/chat services.
6. Events are serialized and broadcast back through the hub.

### Current Image Upload

1. Registered user submits multipart form to `POST /api/create`.
2. Go server validates media type.
3. Go server writes temp file.
4. Go server invokes local Python background-removal script.
5. Go server writes processed image to assets directory.
6. Go server creates `PetImage` and `Pet` rows.

This flow is planned to be replaced by the image processor service.

## Known Architectural Drift

- The app is still browser-first, while the product direction is API-first plus mobile.
- OpenAPI create-pet contract does not fully match implementation.
- Redis is connected but does not yet have a firm responsibility.
- Leaderboard service and UI are placeholders.
- WAL recovery exists but is not connected to startup.
- The image processor service boundary needs to be formalized before new upload work continues.

## Recovery Priority

1. Document current contracts.
2. Align implementation and OpenAPI.
3. Introduce image job/image status model.
4. Replace synchronous local image processing with processor-service submission.
5. Keep browser UI working as one client, but make JSON APIs the canonical client contract.
