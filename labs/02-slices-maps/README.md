# Lab 02: Slices and Maps

## Objectives

- See how slice headers share backing arrays (aliasing).
- Learn when ppend reuses capacity vs allocates a new array.
- Iterate maps and understand non-deterministic key order.

## Setup

1. cd labs/02-slices-maps
2. Run go run . several times and watch map iteration order.

## Exercises

1. Write uniqueInts(in []int) []int preserving first-seen order without extra packages.
2. Given s := make([]int, 0, 2), append three elements and explain when s and a returned slice might alias.
3. Build a map[string][]string of tags per user; safely append to a slice value in the map (copy pattern if needed).
4. Delete a key while ranging (use a two-step collect-then-delete pattern).

## What to take away

- Slices are descriptors (pointer, len, cap); copying the slice copies the descriptor, not necessarily the array.
- ppend may grow into a new array; always treat overlapping slices carefully.
- Map iteration order is randomized; do not depend on it for logic.

## Cleanup

None required.

## Related Modules

- [Slices and maps internals](../../modules/07-slices-maps-internals.md)