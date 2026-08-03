# Lab 13: Design Patterns

Explore three Go-friendly patterns: functional options, decorator-style middleware, and strategy via interfaces.

## Objectives

- Configure structs with `Option` functions instead of large constructors.
- Wrap `http.Handler` to add behavior (logging, recovery).
- Swap algorithms with a small interface (`Pricer`).

## Setup

```bash
cd labs/13-design-patterns
go run .
```

## Exercises

1. Add a `WithReadHeaderTimeout` option and apply it to a real `http.Server`.
2. Start the decorated handler on `:9090` and hit `/hi`; trigger a panic in a handler and confirm recovery.
3. Implement `WholesalePricer` with tiered discounts.
4. Compare functional options vs a dedicated `Config` struct—when is each clearer?

## What to take away

- Functional options scale configuration APIs without breaking callers.
- Middleware is the decorator pattern applied to HTTP handlers.
- Interfaces keep pricing/shipping/tax rules pluggable and testable.

## Cleanup

Stop any HTTP server you started for experiments.

## Related Modules

- Interfaces and composition modules.
- HTTP middleware (Lab 09).
