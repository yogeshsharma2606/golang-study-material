# Lab 11: Profiling

Run a CPU-heavy workload alongside the `net/http/pprof` endpoints to practice capturing profiles.

## Objectives

- Expose pprof over HTTP.
- Generate CPU load intentionally.
- Capture and inspect a CPU profile with the `go tool pprof` CLI.

## Setup

```bash
cd labs/11-profiling
go run .
```

Optional: `WORKERS=2 PPROF_ADDR=:6060 go run .`

## Exercises

1. While the program runs, capture 30 seconds of CPU profile:
   ```bash
   go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
   ```
2. In the flame graph, find `slowFib` and discuss why it dominates.
3. Try `go tool pprof http://localhost:6060/debug/pprof/heap` and compare heap vs CPU.
4. Fix the hotspot (memoization or iterative fib) and re-profile to validate improvement.

## What to take away

- Profile before optimizing; measure real bottlenecks.
- pprof integrates with the standard library with minimal code.
- Profiling overhead is low enough for many production services (often behind auth).

## Cleanup

Stop the process with `Ctrl+C`.

## Related Modules

- Performance tuning and observability modules.
- Concurrency (worker goroutines driving load).
