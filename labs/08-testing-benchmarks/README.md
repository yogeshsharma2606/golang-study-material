# Lab 08: Testing and Benchmarks

## Objectives

- Write table-driven tests with 	.Run subtests.
- Measure performance with benchmarks and .ReportAllocs().
- Explore fuzzing with 	esting.F (Go 1.18+).

## Setup

1. cd labs/08-testing-benchmarks
2. Run go test -v
3. Run go test -bench=. -benchmem
4. Run go test -fuzz=FuzzIsPalindrome -fuzztime=10s

## Exercises

1. Add table cases for Unicode palindromes (e.g. "Ã©tÃ©").
2. Add a failing subtest on purpose and read the failure output; fix it.
3. Compare benchmark ns/op before and after a micro-optimization in IsPalindrome.
4. Extend the fuzz target to reject strings longer than 1 KiB with 	.Skip().

## What to take away

- Table-driven tests keep cases declarative and easy to extend.
- Benchmarks answer "how fast" and "how many allocations"; fuzzing searches for crashes and invariant breaks.

## Cleanup

Remove 	estdata or uzz corpora only if you created extra fuzz outputs you do not want to keep.

## Related Modules

- [Testing, benchmarks, fuzzing](../../modules/18-testing-benchmarks-fuzzing.md)