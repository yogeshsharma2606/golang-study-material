# Lab 12: Concurrency Patterns

Practice three common patterns: worker pool, bounded buffer, and `errgroup` for coordinated goroutines.

## Objectives

- Limit parallelism with a fixed worker pool.
- Back-pressure with a bounded channel buffer.
- Cancel sibling work when one goroutine fails using `errgroup`.

## Setup

```bash
cd labs/12-concurrency-patterns
go mod tidy
go run .
```

## Exercises

1. Reduce buffer capacity to 1 and add a third `Put` without a consumer; observe blocking until timeout.
2. Change worker count and measure runtime for a larger job batch.
3. Force `billing` to always fail and confirm other fetches stop early via shared context.
4. Replace `rand` failures with explicit error injection for deterministic tests.

## What to take away

- Channels encode ownership and back-pressure.
- Worker pools bound resource usage (CPU, DB connections).
- `errgroup` simplifies “run these in parallel, fail fast” orchestration.

## Cleanup

No files or services to tear down.

## Related Modules

- Goroutines, channels, and `context` in the course modules.
- `golang.org/x/sync` extended primitives.
