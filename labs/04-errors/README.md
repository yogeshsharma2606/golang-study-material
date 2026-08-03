# Lab 04: Errors

## Objectives

- Define sentinel errors and custom error types.
- Use errors.Is for sentinel comparison through wrapped errors.
- Use errors.As to extract typed errors from a chain.
- Wrap errors with %w for context.

## Setup

1. cd labs/04-errors
2. Run go run .

## Exercises

1. Add ErrPermissionDenied and return it from loadUser when id == 999.
2. Implement Unwrap() error on ValidationError if you nest a cause (optional pattern).
3. Write describe(err error) string that uses errors.As to format validation errors differently.
4. Compare err == ErrNotFound vs errors.Is when the error is wrapped twice.

## What to take away

- Return errors as values; add context with mt.Errorf("...: %w", err).
- errors.Is checks identity through wrapping; errors.As finds a type in the chain.

## Cleanup

None required.

## Related Modules

- [Error handling](../../modules/09-error-handling.md)