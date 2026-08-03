# Lab 05: Goroutines and Channels

## Objectives

- Build a small pipeline with channel stages and sync.WaitGroup.
- Fan out work to multiple worker goroutines.
- Use select with a timeout to bound waiting.

## Setup

1. cd labs/05-goroutines-channels
2. Run go run . (output may stop early due to the demo timeout).

## Exercises

1. Add a third stage that filters even squares only.
2. Replace the timeout demo with context.WithTimeout (preview of lab 07).
3. Measure whether fan-out with 1 vs 4 workers changes output order; explain why.
4. Fix a deliberate bug: what happens if out is never closed? Trace the program.

## What to take away

- Close channels when producers are done; consumers range until close.
- select multiplexes channel operations; default or timer cases avoid blocking forever.

## Cleanup

None required.

## Related Modules

- [Goroutines and scheduler](../../modules/10-goroutines-scheduler.md)
- [Channels and select](../../modules/11-channels-select.md)