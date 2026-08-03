# Module 2 — Types, Variables & Control Flow

## TL;DR

Go has a small set of built-in types with **zero values** (0, `""`, `nil`, `false`). Variables are statically typed; type inference uses `:=`. Control flow is limited to `for`, `if`, `switch`, and `select` — no `while` or `do-while`.

## Concept

Every variable declared without an explicit value gets its **zero value**:

| Type | Zero Value |
|------|------------|
| `int`, `float` | `0` |
| `string` | `""` |
| `bool` | `false` |
| pointer, slice, map, channel, func, interface | `nil` |

```go
var count int        // 0
name := "Go"         // type inferred as string
const MaxRetries = 3
```

**Arrays vs slices** (preview): arrays have fixed size (`[5]int`), slices are dynamic views (`[]int`). Slices are used 99% of the time.

## How It Really Works (Internals)

- **Type safety**: No implicit numeric conversions — `int32 + int64` is a compile error.
- **iota**: Pre-declared identifier that increments in const blocks — useful for enums.
- **Short declaration `:=`**: Only valid inside functions.
- **switch**: Cases don't fall through by default.

## Why / When / Trade-offs

- Zero values eliminate uninitialized variable bugs.
- Explicit conversions prevent silent precision loss.
- Prefer `for range` over index loops when you don't need the index.

## Worked Scenario

```go
type StatusCategory int

const (
    CategoryInfo StatusCategory = iota
    CategorySuccess
    CategoryClientError
    CategoryServerError
)

func categorize(code int) StatusCategory {
    switch {
    case code >= 500:
        return CategoryServerError
    case code >= 400:
        return CategoryClientError
    case code >= 200:
        return CategorySuccess
    default:
        return CategoryInfo
    }
}
```

## Gotchas & Failure Modes

- `:=` can shadow outer variables.
- `nil` slice vs empty slice: both `len == 0`, but JSON serialization differs.
- Structs with slice/map fields are not comparable.

## Interview Q&A

**Q: What are zero values in Go?**
A: Every declared variable is initialized to a type-specific default: 0 for numbers, `""` for strings, `false` for bool, `nil` for pointers/slices/maps/channels/interfaces/functions.
↳ How does this differ from C? C stack variables are uninitialized (garbage values).

**Q: What are the key differences between an array and a slice?**
A: Arrays have fixed compile-time size and are value types. Slices are dynamic-length views over an array, passed by a header (pointer, len, cap).

**Q: How does Go's type system ensure type safety?**
A: Static typing, no implicit conversions, explicit type assertions, compile-time generics checking.

## Verify

```bash
cd labs/01-basics && go run .
```

## Further Reading

- [Go Spec — Types](https://go.dev/ref/spec#Types)
- [Go Blog — Constants](https://go.dev/blog/constants)
