# Lab 01: Basics â€” Types, Pointers, Defer, Variadics

## Objectives

- Work with basic types, zero values, and short variable declarations.
- Pass and modify values through pointers.
- Use defer and understand LIFO cleanup order.
- Write and call variadic functions.

## Setup

1. cd labs/01-basics
2. Run go run . to see the guided demo output.

## Exercises

1. Add a function swap(a, b *int) that exchanges two integers; call it from main and print before/after.
2. Write a variadic max(nums ...int) (int, error) that returns an error if no arguments are passed.
3. Add a defer in main that prints "main exiting" and observe order relative to deferDemo.
4. Change describe to guard against a nil pointer (print a message instead of dereferencing).

## What to take away

- Pointers let callees mutate caller state; & takes address, * dereferences.
- defer schedules work at function return; multiple defers run last-in-first-out.
- ...T in the parameter list collects arguments into a slice inside the function.

## Cleanup

No generated artifacts beyond optional binaries from go build.

## Related Modules

- [Go fundamentals](../../modules/01-go-fundamentals.md)
- [Types and control flow](../../modules/02-types-control-flow.md)
- [Functions and methods](../../modules/03-functions-methods.md)
- [Pointers and value semantics](../../modules/04-pointers-value-semantics.md)