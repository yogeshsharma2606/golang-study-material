# Lab 06: Sync and Memory Model

## Objectives

- Observe a data race and fix it with sync.Mutex.
- Use sync.Once for one-time initialization.
- Know when to run tests and programs with -race.

## Setup

1. cd labs/06-sync-memory-model
2. Run go run . for the safe demos.
3. Uncomment unRace() in main.go and run go run -race . to see the race detector report.

## Exercises

1. Replace the mutex counter with tomic.AddInt64 and compare idioms.
2. Add a sync.RWMutex around a shared map[string]int with concurrent readers and writers.
3. Explain why sync.Once is preferable to a boolean flag guarded by a mutex for lazy init.
4. Run go test -race ./... after adding a small test that spawns goroutines.

## What to take away

- The memory model requires synchronization for shared mutable state; races are undefined behavior.
- Prefer mutexes or atomics; sync.Once guarantees a function runs exactly once.

## Cleanup

None required.

## Related Modules

- [Sync and memory model](../../modules/12-sync-memory-model.md)