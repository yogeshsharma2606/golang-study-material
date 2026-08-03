# Lab 09: HTTP API

Build a small REST API using only the Go standard library (`net/http`), with middleware and graceful shutdown.

## Objectives

- Register routes with Go 1.22+ `ServeMux` method patterns.
- Compose middleware for logging and request IDs.
- Return JSON responses with appropriate status codes.
- Shut down the server cleanly on `SIGINT` / `SIGTERM`.

## Setup

```bash
cd labs/09-http-api
go run .
```

The server listens on `:8080` (override with `PORT`).

## Exercises

1. Call `GET http://localhost:8080/health` and confirm `{"status":"ok"}`.
2. List items: `GET http://localhost:8080/api/items`.
3. Create an item:
   ```bash
   curl -X POST http://localhost:8080/api/items \
     -H "Content-Type: application/json" \
     -d "{\"name\":\"Pen\",\"price_cents\":199}"
   ```
4. Fetch by id: `GET http://localhost:8080/api/items/1`.
5. Send a custom `X-Request-ID` header and verify it appears in logs and the response.
6. Press `Ctrl+C` and confirm the server logs shutdown without killing in-flight requests abruptly (try a slow client if you extend the lab).

## What to take away

- Middleware wraps `http.Handler` to add cross-cutting concerns.
- `http.Server` timeouts protect against slow clients.
- `Shutdown` with a context deadline is the idiomatic graceful stop.

## Cleanup

Stop the server with `Ctrl+C`. No external services or files are required.

## Related Modules

- HTTP fundamentals and routing in the course modules on `net/http`.
- Concurrency and context (used during shutdown).
