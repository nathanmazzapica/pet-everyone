# 0002: Use a Separate Image Processor Service

## Status

Draft accepted direction

## Context

The original Go server processes uploaded images synchronously by invoking a local Python `withoutbg` script. This blocks the request lifecycle, couples Go app runtime to Python environment setup, and makes retries/status handling difficult.

A separate image processor service is being developed to own image processing and queue execution.

The local script approach also has unfavorable runtime characteristics:

- the background-removal model has a high memory footprint
- invoking the script from the Go request path causes an expensive cold start for each upload
- model loading/setup cost is paid repeatedly instead of being amortized by a warm worker process
- request latency and server memory pressure become tied to image processing workload

## Decision

The Go server will submit image processing jobs to a separate image processor service over an internal API authenticated with a shared secret.

The Go server remains the source of truth for product state. The image processor may have its own database, but that database is internal to the processor and stores job execution state only.

## Consequences

- Image processing becomes asynchronous.
- Clients need image status states.
- The Go server needs product image records before processing finishes.
- The image processor needs stable job IDs and idempotent submission behavior.
- The services communicate through HTTP APIs/webhooks, not direct database access.
- Image processing workers can keep expensive model/runtime state warm instead of paying cold-start cost per request.
- Memory-heavy processing is isolated from the Go web server.

## Ownership Boundary

Go server owns:

- users
- pets
- image ownership
- user-visible image status
- active pet image selection
- final processed image URL/storage key

Image processor owns:

- processing jobs
- attempts
- worker locks
- retries
- operational errors
- intermediate artifacts

## Open Questions

- Should the processor callback the Go server, or should the Go server poll?
- Where should originals and processed images be stored?
- What queue semantics does the processor database provide?
- How does the Go server expose image status to browser/mobile clients?
- What is the retry/dead-letter policy?
