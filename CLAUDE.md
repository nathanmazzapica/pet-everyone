# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pet Everyone is an online petting zoo where registered users can upload images of their pets and invite their friends to pet them. The application is built with Go, using WebSockets for real-time interactions, PostgreSQL for persistent storage, and Redis for caching.

## Development Commands

### Running the Application
```bash
# Start dependencies (PostgreSQL and Redis)
docker-compose up -d

# Run the web server
go run cmd/web/main.go
```

The server will be available at `http://localhost:<PORT>/app` (PORT from .env file).

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/service
go test ./internal/wal

# Run a specific test
go test -run TestFunctionName ./path/to/package
```

### Database
PostgreSQL runs on port 5433 (mapped from 5432) with credentials:
- User: `devuser`
- Password: `devpass`
- Database: `devdb`

Migrations are located in `internal/data/db/migrations/`.

### Environment Variables
Required in `.env`:
- `FILEPATH_ROOT`: Root directory for file storage
- `ASSETS_ROOT`: Directory for static assets
- `PORT`: Server port
- Database and Redis connection settings

## Architecture

### Real-Time WebSocket System

The application uses a sophisticated multi-layer architecture for WebSocket-based real-time interactions:

**Hub Registry Pattern** (`internal/registry/registry.go`):
- `HubRegistry` maintains a thread-safe map of active Hubs, one per pet
- Each Hub is lazily created when the first client connects to view a pet
- Hubs automatically shut down when the last client disconnects (via context cancellation)
- This ensures resources are only allocated for pets being actively viewed

**Hub Lifecycle** (`internal/websocket/hub.go`):
- Each `Hub` manages WebSocket clients for a specific pet
- Contains channels for `register`, `unregister`, and `broadcast`
- The Hub's `Run()` loop handles client registration/unregistration and message broadcasting
- When the last client unregisters, the Hub starts a 30-second shutdown timer (prevents unnecessary teardown on page refreshes)
- If a client reconnects during the delay, the shutdown is cancelled
- If the timer expires, the Hub calls its `cancel()` function to trigger cleanup

**Service Layer** (`internal/service/`):
- `PetService`: Manages pet interaction counts with in-memory buffering and periodic DB writes
  - Buffers clicks in `dbQueue` and flushes to DB every 10 seconds
  - Implements retry logic with `pendingUpdates` for failed DB writes
  - Maintains a live `petCount` for instant UI updates
- `ChatService`: Handles chat message broadcasting
- Services publish `Event` structs to an output channel

**Transport Router** (`internal/transport/router.go`):
- Routes incoming commands to appropriate services based on command type
- Commands flow: WebSocket → Hub → Router → Service → Events → Serializer → Broadcast
- `Envelope` wraps commands with sender metadata
- Supports command types: "pet", "petcount", "chat"

**Serializer** (`internal/transport/serializer.go`):
- `JSONSerializer` subscribes to service events and converts them to JSON
- Publishes serialized messages to the Hub's broadcast channel
- Decouples service layer from WebSocket message format

**Data Flow**:
1. Client sends WebSocket message → `Client.readPump()`
2. Client forwards to Hub's command channel via Router
3. Router unmarshals and routes to appropriate Service
4. Service processes and emits Event
5. Serializer converts Event to JSON
6. Hub broadcasts JSON to all connected clients

### Persistence and Reliability

**Write-Ahead Log (WAL)** (`internal/wal/`):
- `SimpleWAL` provides crash recovery for pet counts
- Format: `userID,count\n` for each entry
- `Recover()` rebuilds state by scanning the entire log
- File kept open during service lifetime for performance

**Database Layer** (`internal/data/`):
- `model/`: Database models (User, Pet, SessionToken)
- `db/`: PostgreSQL connection using pgx
- `cache/`: Redis connection for session storage
- Services use `PetDatabase` interface for testing

### Web Application Structure

**Main Entry** (`cmd/web/main.go`):
- Loads `.env` configuration
- Connects to PostgreSQL and Redis
- Initializes `HubRegistry` with `PetModel`
- Creates `application.Config` with all dependencies
- Sets up HTTP server with routes

**Application Layer** (`cmd/web/application/`):
- `Config`: Holds DB, cache, registry, and path configuration
- Asset management for uploaded pet images
- JSON helpers for API responses

**Handlers** (`cmd/web/handler/`):
- HTTP request handlers for signup, login, pet management
- WebSocket upgrade endpoint connects clients to appropriate Hub
- Routes defined in `handler/` (exact file to be determined)

**Middleware** (`cmd/web/middleware/`):
- `auth.go`: Session authentication
- `logging.go`: Request logging

### Code Standards (from .github/instructions)

When reviewing or implementing code in this repository, apply senior-level engineering standards:
- Evaluate architectural decisions: why this approach vs alternatives?
- Identify concurrency risks, race conditions, and edge cases
- Ensure code is testable through loose coupling and clear interfaces
- Challenge tight coupling and missing abstractions that reduce testability
- Communicate intent through clear naming and focused comments
- Use idiomatic Go patterns (proper error handling, context usage, goroutine management)
