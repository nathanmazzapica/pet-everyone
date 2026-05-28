# Realtime Protocol Contract

## Status

Draft. Current implementation uses WebSockets for pet rooms. SSE is planned for server-to-client status streams such as leaderboard and image processing status.

## Transport Split

Use HTTP for one-shot actions:

- login/signup/logout
- create pet
- upload image
- fetch pet details
- fetch counts

Use WebSocket for active pet-room interactions:

- pet taps
- chat messages
- live count deltas
- active room presence, if added

Use SSE for server-to-client updates:

- leaderboard updates
- image processing status
- other low-frequency status streams

Use mobile push notifications later for background notifications. Do not rely on WebSockets or SSE while a mobile app is backgrounded or closed.

## WebSocket Endpoint

Current:

```http
GET /pet/{pet_id}/ws
```

Future API-oriented form:

```http
GET /v1/pets/{pet_id}/ws
```

## WebSocket Auth

Browser:

- `session_token` cookie for registered users
- `guest_id` cookie for guests

Mobile:

- expected future support for `Authorization: Bearer <session-token>`

Avoid auth tokens in query strings.

## Current Client Commands

Pet tap:

```json
{
  "type": "pet",
  "data": null
}
```

Request full pet count:

```json
{
  "type": "petcount",
  "data": null
}
```

Chat:

```json
{
  "type": "chat",
  "data": {
    "msg": "hello"
  }
}
```

## Current Server Events

Pet delta:

```json
{
  "type": "pet",
  "data": {
    "c": 1
  }
}
```

Full pet count:

```json
{
  "type": "petcount",
  "data": 1234
}
```

Chat:

```json
{
  "type": "chat",
  "data": {
    "msg": "hello",
    "author": "displayName1234"
  }
}
```

## Future Protocol Direction

Future messages should add protocol versioning and request/event IDs:

```json
{
  "v": 1,
  "type": "pet.tap",
  "request_id": "uuid",
  "data": {}
}
```

```json
{
  "v": 1,
  "type": "pet.count_delta",
  "event_id": "uuid-or-sequence",
  "pet_id": "uuid",
  "data": {
    "delta": 1
  }
}
```

## SSE Candidate: Leaderboard

Endpoint:

```http
GET /v1/pets/{pet_id}/leaderboard/events
Accept: text/event-stream
```

Event:

```text
event: leaderboard
data: {"pet_id":"uuid","leaders":[{"display_name":"Name1234","count":42}]}
```

## SSE Candidate: Image Processing Status

Endpoint:

```http
GET /v1/images/{image_id}/events
Accept: text/event-stream
```

Event:

```text
event: image_status
data: {"image_id":"uuid","status":"processing"}
```

Ready event:

```text
event: image_status
data: {"image_id":"uuid","status":"ready","processed_image_url":"https://..."}
```

Failure event:

```text
event: image_status
data: {"image_id":"uuid","status":"failed","error_code":"background_removal_failed"}
```
