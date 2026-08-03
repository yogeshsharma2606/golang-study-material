# Lab 14: Microservices

Two small HTTP services: **inventory** (data) and **storefront** (calls inventory). Both expose health checks.

## Objectives

- Run multiple services locally with clear ports and env configuration.
- Call another service over HTTP with timeouts.
- Propagate failures as `502`/`504`-style responses.

## Setup

Terminal 1:

```bash
cd labs/14-microservices
go run ./cmd/inventory
```

Terminal 2:

```bash
cd labs/14-microservices
INVENTORY_URL=http://localhost:8081 go run ./cmd/storefront
```

## Exercises

1. `curl http://localhost:8081/health` and `curl http://localhost:8080/health`.
2. `curl http://localhost:8080/products/go-book` — compare with unknown SKU.
3. Stop inventory and call storefront again; observe graceful degradation messaging.
4. Add a shared `X-Request-ID` header from storefront to inventory.

## What to take away

- Health endpoints are the minimum operability contract for orchestrators.
- Always set `http.Client` timeouts for service-to-service calls.
- Keep services loosely coupled via HTTP + JSON for learning; production may use gRPC or events.

## Cleanup

Stop both processes with `Ctrl+C`.

## Related Modules

- HTTP APIs (Lab 09).
- Distributed systems and resilience topics in the course.
