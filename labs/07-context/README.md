# Lab 07: Context

## Objectives

- Propagate cancellation and deadlines with context.Context.
- Simulate an HTTP handler that times out long work.
- Cancel work manually with context.WithCancel.

## Setup

1. cd labs/07-context
2. Run go run . â€” expect timeout on the simulated /api/users request and cancellation on the manual demo.

## Exercises

1. Switch WithTimeout to WithDeadline using 	ime.Now().Add(...).
2. Pass a value with context.WithValue (use a private key type) and read it inside simulateSlowWork.
3. Start a real http.Server on :8080 with the same handler; curl with and without slow backend.
4. Ensure goroutines exit when context is canceled (no leaks after timeout).

## What to take away

- Always pass ctx as the first parameter to I/O-bound functions.
- Derive child contexts with cancel functions; call defer cancel() to release timer resources.

## Cleanup

Stop any server you started for exercise 3.

## Related Modules

- [Context package](../../modules/13-context-package.md)