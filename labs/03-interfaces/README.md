# Lab 03: Interfaces

## Objectives

- Understand implicit interface satisfaction (no implements keyword).
- Use struct embedding to compose behavior.
- Apply a simple mock type that records calls for testing.

## Setup

1. cd labs/03-interfaces
2. Run go run .

## Exercises

1. Add a FileLogger that writes to os.Stderr and use it through Logger.
2. Change 
otify to accept mt.Stringer for the message body; pass a custom type.
3. Embed MockLogger in another struct and override Log once; observe which method runs.
4. Write a function sLogger(v any) (Logger, bool) using a type assertion.

## What to take away

- Interfaces describe behavior; concrete types satisfy them by having the right methods.
- Embedding forwards promoted methods; explicit methods on the outer type override promoted ones.
- Mocks implement the same interface as production code for isolated tests.

## Cleanup

None required.

## Related Modules

- [Structs and composition](../../modules/05-structs-composition.md)
- [Interfaces and polymorphism](../../modules/06-interfaces-polymorphism.md)