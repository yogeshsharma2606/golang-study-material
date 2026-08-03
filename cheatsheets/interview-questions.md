# Go Interview Question Bank

120+ questions organized by topic with model answers and follow-ups. Cross-reference modules for deeper study.

**Format:** numbered question → **A:** model answer → ↳ follow-up probe.

---

## Golang Fundamentals (M01–M04)

*See [M01](../modules/01-go-fundamentals.md), [M02](../modules/02-types-control-flow.md), [M03](../modules/03-functions-methods.md), [M04](../modules/04-pointers-value-semantics.md)*

1. What is the zero value of each basic type in Go?
**A:** `false` for bool, `0` for numeric types, `""` for string, `nil` for pointers, slices, maps, channels, functions, and interfaces. Structs zero all fields; arrays zero all elements.
↳ Why does Go have zero values instead of uninitialized variables? It eliminates a class of undefined-behavior bugs and makes constructors optional.

2. How do Go modules differ from GOPATH?
**A:** GOPATH was a single global workspace under `$GOPATH/src`. Go Modules (since 1.11, default 1.16+) give each project its own `go.mod` with explicit versioned dependencies, enabling reproducible builds from any directory.
↳ How do you consume a private module? Set `GOPRIVATE` and configure git credentials or a private proxy.

3. What is the difference between an array and a slice?
**A:** An array has a fixed size known at compile time and is a value type (copying copies all elements). A slice is a dynamic view over a backing array with a header (pointer, len, cap); copying a slice copies the header, not the backing array.
↳ What happens when `append` exceeds capacity? A new backing array is allocated (typically 2× growth), elements are copied, and the slice header points to the new array.

4. Can you compare two slices with `==`?
**A:** Only to `nil`. Two non-nil slices are never equal with `==`; use `bytes.Equal` for `[]byte` or `reflect.DeepEqual` (or element-by-element comparison) for other element types.
↳ Why? Slices are headers; `==` would not compare backing array contents and would be ambiguous for overlapping slices.

5. What is a pointer and when should you use one?
**A:** A pointer holds the memory address of a value. Use pointers to mutate a value across function boundaries, avoid copying large structs, or share mutable state. Use values for small, immutable data.
↳ What does `new(T)` return? A `*T` pointing to a zero-initialized value on the heap (escape analysis may still stack-allocate in some cases).

6. Explain pass-by-value semantics in Go.
**A:** Go always passes arguments by value. For slices, maps, and channels, the value copied is a header or internal pointer — mutations through that copy are visible to other holders of the same underlying data.
↳ Give an example where a function cannot modify a caller's slice length. `append` that triggers reallocation updates only the local slice header unless you return the new slice.

7. How does `defer` work?
**A:** `defer` schedules a function call to run when the surrounding function returns, in LIFO order. Arguments are evaluated immediately; the deferred function runs at return time (after return values are set but before the function exits to the caller).
↳ What is a common defer pitfall with loops? Deferring inside a loop without a closure captures loop variables incorrectly — wrap in `func() { defer f() }()`.

8. What is the difference between `panic` and returning an `error`?
**A:** `error` is the idiomatic, expected failure path. `panic` unwinds the stack and should be reserved for programmer bugs or truly unrecoverable states. Use `recover` only in deferred functions at API boundaries (e.g., HTTP middleware).
↳ Should libraries panic? Generally no — return errors. `json.Unmarshal` panics only on programmer error (invalid type passed to `Unmarshal`).

9. What are variadic functions?
**A:** Functions declared as `func f(nums ...int)` accept zero or more arguments of type `int`, available inside as a slice `nums`. Callers can pass `f(1, 2, 3)` or `f(slice...)`.
↳ What is the type of `nums` inside the function? `[]int` — a slice, not a special type.

10. How does type safety work in Go?
**A:** Go is statically typed with no implicit numeric conversions. Assignments and operations require compatible types; use explicit conversion `T(v)`. Interfaces provide polymorphism but require explicit type assertions for concrete access.
↳ Is `interface{}` type-safe? It defers type checks to runtime via assertions; misuse causes panics. Prefer generics or concrete types where possible.

11. What is an interface in Go and how is it satisfied?
**A:** An interface is a set of method signatures. A type satisfies an interface implicitly by implementing all methods — no `implements` keyword. Interface values hold a (type, value) pair.
↳ What is the nil interface gotcha? An interface is nil only if both type and value are nil. `var p *T = nil; var i interface{} = p` — `i != nil` because the type is `*T`.

12. Explain type assertion and type switch.
**A:** `v, ok := i.(T)` extracts concrete type `T` from interface `i`; without `ok`, wrong type panics. `switch v := i.(type)` dispatches on dynamic type.
↳ When would you use a type switch over polymorphism? Interop with `encoding/json` unmarshaling, plugin systems, or handling multiple concrete types in a single handler.

13. What is the difference between a named struct type and an anonymous struct?
**A:** A named struct (`type User struct{...}`) can have methods, be reused across packages, and documents intent. An anonymous struct is defined inline — useful for one-off JSON shapes or test fixtures without polluting the type namespace.
↳ Can anonymous structs have methods? No — methods require a named receiver type.

14. What is struct embedding and how does it differ from inheritance?
**A:** Embedding promotes the embedded type's fields and methods to the outer struct. It is composition, not inheritance — no subtype polymorphism unless the outer type satisfies interfaces through promoted methods.
↳ What if outer and embedded types both have a method with the same name? The outer method shadows; access embedded via `outer.Embedded.Method()`.

15. How do you create custom errors in Go?
**A:** Implement the `error` interface with `Error() string`. For sentinel errors use `var ErrNotFound = errors.New("not found")`. For structured errors use custom types and check with `errors.Is`/`errors.As`. Wrap with `fmt.Errorf("context: %w", err)`.
↳ When should you export error variables? Export when callers need `errors.Is`; keep internal errors unexported.

16. What is `errors.Is` vs `errors.As`?
**A:** `errors.Is(err, target)` walks the error chain checking equality to a sentinel. `errors.As(err, &target)` finds the first error in the chain assignable to `target`'s type and assigns it.
↳ Does `==` work on wrapped errors? No — use `errors.Is` to traverse `%w` chains.

17. What are runes and how do they relate to strings?
**A:** A `rune` is an alias for `int32` representing a Unicode code point. Go strings are immutable UTF-8 byte sequences. Indexing a string yields bytes, not characters — use `range` for runes or `utf8.DecodeRuneInString`.
↳ Why is `len("é")` sometimes 2? `é` can be two UTF-8 bytes; `len` counts bytes, not runes.

18. What is the difference between `make` and `new`?
**A:** `new(T)` allocates zeroed memory and returns `*T`. `make` initializes slices, maps, and channels and returns the type itself (not a pointer). Only `make` can specify slice capacity and channel buffer size.
↳ Can you `make` a map with initial capacity? `make(map[K]V, hint)` pre-allocates buckets for ~hint entries.

19. What are constants in Go and what is `iota`?
**A:** Constants are compile-time values. `iota` in a const block auto-increments, enabling enumerated constants and bitmasks.
↳ Can a const be a slice? No — constants must be compile-time computable; slices are runtime values.

20. How do you export identifiers in Go?
**A:** Names starting with an uppercase letter are exported (public to other packages). Lowercase names are package-private.
↳ Can you export a struct field selectively? Yes — mixed exported/unexported fields in the same struct.

21. What is the `internal` directory convention?
**A:** Packages under `internal/` can only be imported by code within the parent tree of `internal`. The compiler enforces this — it is stronger than unexported names.
↳ Can another module import your `internal` package? No — compile error regardless of visibility keywords.

22. Explain short variable declaration `:=`.
**A:** `:=` declares and initializes at least one new variable in the current scope. It can only be used inside functions. At least one variable on the left must be new.
↳ What happens with `x, err := f()` when `err` was already declared? `err` is redeclared in the inner block if in a new scope (e.g., `if`), otherwise compile error.

23. What is a method and how do value vs pointer receivers differ?
**A:** A method is a function with a receiver. Value receivers operate on a copy; pointer receivers can mutate and avoid copy cost. If any method needs a pointer receiver, use pointer receivers consistently for that type.
↳ Can you call pointer-receiver methods on values? Yes — Go takes the address automatically (`v.Method()` → `(&v).Method()`).

24. What is the blank identifier `_` used for?
**A:** Discard unwanted values: `_, err := f()`. Import for side effects: `_ "net/http/pprof"`. Interface satisfaction compile check: `var _ io.Reader = (*MyType)(nil)`.
↳ Why import for side effects? Register HTTP handlers, database drivers (`_ "github.com/lib/pq"`).

25. What are type aliases vs type definitions?
**A:** `type New = Old` creates an alias — identical type. `type New Old` creates a distinct type with the same underlying type — no implicit conversion between them.
↳ When use a definition over an alias? When you want distinct methods or to prevent accidental interchange.

26. How does Go handle unused imports and variables?
**A:** Compile error — unused imports and variables are not allowed. Use `_` to explicitly discard.
↳ Why is this strict? Forces clean code and catches dead code early.

---

## Concurrency (M10–M14)

*See [M10](../modules/10-goroutines-scheduler.md), [M11](../modules/11-channels-select.md), [M12](../modules/12-sync-memory-model.md), [M13](../modules/13-context-package.md), [M14](../modules/14-concurrency-patterns.md)*

27. What is a goroutine and how does it differ from an OS thread?
**A:** A goroutine is a lightweight concurrent function scheduled by the Go runtime. Thousands fit in memory (≈2 KB initial stack, grows as needed) vs ~1 MB per OS thread. The runtime multiplexes goroutines onto `GOMAXPROCS` kernel threads (GMP model).
↳ What is the GMP model? G = goroutine, M = OS thread, P = logical processor (context for scheduling).

28. How do you start a goroutine?
**A:** Prefix a function call with `go`: `go func() { work() }()`. The caller does not wait — synchronization requires channels, `sync` primitives, or `WaitGroup`.
↳ What happens if the main goroutine exits before others finish? The program terminates — all goroutines are killed.

29. What is a race condition and how do you detect it?
**A:** A race occurs when two goroutines access shared memory concurrently and at least one access is a write, without synchronization. Run `go test -race` or `go run -race` to detect.
↳ Is a race a compile-time error? No — it is undefined behavior detected at runtime by the race detector (instrumentation overhead ~2–10×).

30. What is `sync.WaitGroup`?
**A:** A counter for waiting on a group of goroutines. Call `Add(n)` before spawning, `Done()` in each goroutine (typically `defer wg.Done()`), and `Wait()` in the parent to block until the counter reaches zero.
↳ What happens if `Add` is called after `Wait`? Panic — `Add` must complete before `Wait` returns.

31. What is the difference between buffered and unbuffered channels?
**A:** Unbuffered channels synchronize sender and receiver (rendezvous) — send blocks until receive. Buffered channels (`make(chan T, n)`) allow up to `n` sends without a receiver; send blocks when full.
↳ When should you use buffer size 0 vs N? Default unbuffered for backpressure; buffer when you have a known bound on in-flight work and measured producer blocking.

32. What happens when you close a channel?
**A:** No more sends allowed (panic if sent after close). Receivers drain remaining values, then receive zero value and `ok == false`. Closing is only the sender's responsibility.
↳ Who should close a channel? The sender — never the receiver. Use `sync.Once` or a dedicated goroutine if multiple senders.

33. How does `select` work?
**A:** `select` blocks until one of its cases can proceed (send/receive on channels) or `default` runs immediately (non-blocking). If multiple cases are ready, one is chosen pseudo-randomly.
↳ How do you implement a timeout with select? `case <-time.After(d):` or prefer `context.WithTimeout`.

34. What is the difference between `sync.Mutex` and `sync.RWMutex`?
**A:** `Mutex` provides exclusive lock. `RWMutex` allows multiple concurrent readers OR one writer — better for read-heavy workloads.
↳ When does RWMutex hurt performance? When writes are frequent — reader lock overhead exceeds benefit.

35. What is `sync.Once`?
**A:** Ensures a function runs exactly once, even with concurrent calls — used for singleton initialization (e.g., `sync.Once` around `db.Open`).
↳ Can `Once` be reset? No — create a new `Once` if you need re-initialization.

36. Explain the Go memory model happens-before rules.
**A:** Synchronization events establish ordering: channel send happens-before receive, `Unlock` happens-before subsequent `Lock`, `Once.Do` completion happens-before any return from `Do`, goroutine creation happens-before goroutine start.
↳ Why can't you rely on write visibility without sync? CPU caches and compiler reordering — unsynchronized reads may see stale values.

37. What is `context.Context` and why use it?
**A:** `context` carries deadlines, cancellation signals, and request-scoped values across API boundaries. When a context is canceled, all derived operations should stop promptly.
↳ Should you store contexts in structs? No — pass as the first function parameter: `func(ctx context.Context, ...)`.

38. What are `context.WithCancel`, `WithTimeout`, and `WithDeadline`?
**A:** `WithCancel` returns a cancel function. `WithTimeout`/`WithDeadline` auto-cancel after duration/time. Always call `cancel()` to release resources (use `defer`).
↳ What error does a canceled context return? `context.Canceled` or `context.DeadlineExceeded` via `ctx.Err()`.

39. When should you use `context.WithValue`?
**A:** Sparingly — for request-scoped data (trace ID, auth claims) that crosses API boundaries. Not for optional parameters — use function arguments instead.
↳ Why is `WithValue` controversial? Untyped keys, no compile-time safety, encourages god-context anti-pattern.

40. What is a worker pool pattern?
**A:** Fixed number of goroutine workers consume jobs from a channel, limiting concurrency and resource usage. Producers send to `jobs` channel; workers process and optionally send to `results`.
↳ How do you shut down a worker pool gracefully? Close `jobs` after sending all work, `WaitGroup` on workers, or cancel via context.

41. Explain fan-out and fan-in.
**A:** Fan-out: multiple goroutines read from one channel (parallelize work). Fan-in: multiplex multiple channels into one (merge results). Often combined in pipelines.
↳ How do you fan-in without leaking goroutines? Use `sync.WaitGroup` + one merger goroutine, or `errgroup`.

42. What is `golang.org/x/sync/errgroup`?
**A:** Runs goroutines with shared context — first error cancels siblings. `g.Go(func() error {...})` then `g.Wait()` returns the first error.
↳ How does errgroup differ from WaitGroup? Built-in error propagation and context cancellation.

43. What is a semaphore in Go?
**A:** Limit concurrent goroutines — buffered channel of empty structs (`sem := make(chan struct{}, n)`), acquire by send, release by receive. Or use `golang.org/x/sync/semaphore`.
↳ Why channel over counting semaphore library? Idiomatic, no extra dependency for simple cases.

44. How do you prevent goroutine leaks?
**A:** Ensure every goroutine has an exit path: context cancellation, channel close, or timeout. Avoid blocked sends/receives when no counterpart exists.
↳ How do you find leaks in production? `pprof` goroutine profile — look for growing stacks blocked on channels.

45. Can maps be used concurrently without synchronization?
**A:** No — concurrent read/write or write/write on maps panics. Use `sync.Mutex`/`RWMutex` around the map, `sync.Map` for specific read-heavy patterns, or shard maps.
↳ When is `sync.Map` appropriate? Many reads, infrequent writes, stable key sets (e.g., cache).

46. What does `GOMAXPROCS` control?
**A:** Maximum OS threads executing Go code simultaneously. Defaults to `runtime.NumCPU()`. Increase rarely helps CPU-bound; I/O-bound benefits from many goroutines regardless.
↳ What about container CPU limits? Go 1.19+ respects cgroup quotas via `automaxprocs` libraries or runtime updates.

47. What is the difference between concurrency and parallelism?
**A:** Concurrency is structure — dealing with many things at once. Parallelism is execution — doing many things simultaneously on multiple cores. Go enables concurrent design; parallelism requires `GOMAXPROCS > 1` and CPU-bound work.
↳ Can you have concurrency without parallelism? Yes — single-core machine with goroutines interleaving on I/O waits.

48. How do atomics compare to mutexes?
**A:** `sync/atomic` provides lock-free operations on integers and pointers — lower overhead for simple counters/flags. Mutexes protect compound invariants across multiple fields.
↳ When is atomic insufficient? When you need to read-modify-write multiple related fields atomically — use a mutex.

49. What is a pipeline pattern in Go?
**A:** Stages connected by channels — each stage is a group of goroutines running the same function, receiving from inbound and sending to outbound channels.
↳ How do you handle stage errors in pipelines? Use `errgroup`, pass `context`, or dedicated error channel merged at the end.

50. Explain the "share memory by communicating" principle.
**A:** Prefer passing data ownership via channels over protecting shared state with locks. Channels encode synchronization and ownership transfer explicitly.
↳ When should you ignore this principle? Hot-path shared caches, simple counters — mutex may be simpler and faster.

51. What happens with `select` and a nil channel?
**A:** Operations on nil channels block forever — a `case` with nil channel is never selected, useful for disabling cases dynamically.
↳ Give a use case. Fan-in merger disabling an exhausted input channel by setting it to nil.

---

## Memory Management (M22)

*See [M22](../modules/22-gc-runtime-performance.md)*

52. How does the Go garbage collector work?
**A:** Go uses a concurrent tri-color mark-and-sweep collector with write barriers. It runs alongside the application (with brief STW phases for termination and sweep setup). Target: sub-millisecond STW pauses.
↳ What are the three colors? White (unreachable), gray (reachable, children not scanned), black (reachable, children scanned).

53. What is escape analysis?
**A:** Compiler analysis determining whether a variable's lifetime exceeds its stack frame. If it "escapes" (returned pointer, stored in heap structure, closure capture), it is allocated on the heap.
↳ How do you see escape decisions? `go build -gcflags='-m'` shows escape messages.

54. What is the difference between stack and heap allocation?
**A:** Stack: fast, scoped to function, no GC pressure. Heap: survives function return, GC-managed, slower allocation. Go stacks grow/shrink dynamically (starting ~2 KB).
↳ Can you force stack allocation? No explicit control — write code that doesn't escape; compiler decides.

55. What does `runtime.GC()` do and when should you call it?
**A:** Forces a garbage collection cycle. Rarely needed in application code — the runtime triggers GC based on `GOGC` (default 100 = collect when heap doubles). May help in benchmarks or memory-sensitive batch jobs.
↳ What is `GOGC`? Percentage heap growth trigger; `GOGC=off` disables GC (testing only).

56. How do you profile memory in Go?
**A:** `go test -memprofile=mem.prof`, `go tool pprof -http=:8080 mem.prof`, or HTTP `/debug/pprof/heap`. Look at `inuse_space`, `alloc_space`, and top allocators.
↳ What is the difference between inuse and alloc? Inuse = currently live; alloc = cumulative allocations (find churn).

57. What causes memory leaks in Go?
**A:** GC only collects unreachable objects. Leaks occur when references are unintentionally retained: global slices/maps growing, goroutines blocked forever, `time.Ticker` not stopped, HTTP response bodies not closed, finalizer chains.
↳ How do goroutines cause leaks? Blocked on channel send/receive with no counterpart — goroutine and its stack vars stay reachable.

58. What are pointer best practices in Go?
**A:** Prefer value semantics; use pointers for mutation, large structs, or optional fields (`*string` for SQL NULL). Avoid unnecessary pointers in slices of small structs (cache unfriendly). Don't use `unsafe` without extreme care.
↳ Should slice elements be pointers? Use `[]*T` when elements are large or need in-place mutation shared across references; `[]T` when small and immutable.

59. What is `sync.Pool`?
**A:** A pool of temporary objects to reduce allocation pressure. Objects may be GC'd at any time — no guarantee of retrieval. Use for frequently allocated short-lived buffers after profiling.
↳ When should you NOT use sync.Pool? When object lifecycle or correctness depends on persistence across GC cycles.

60. How does Go handle stack growth?
**A:** Goroutine stacks start small and grow by copying to larger stacks when needed (historically segmented stacks, now contiguous copy). This enables cheap goroutine creation.
↳ What is stack copying cost? Proportional to stack size — deep recursion can trigger multiple copies.

61. What is write barrier overhead?
**A:** During GC mark phase, the runtime intercepts pointer writes to maintain the tri-color invariant. Most application code has negligible impact; pointer-heavy write paths during GC may see slight cost.
↳ Is Go's GC generational? Not traditionally — Go 1.19+ has soft generational hints (page spans) but is primarily mark-sweep.

62. How do you reduce GC pressure?
**A:** Reduce allocations in hot paths, reuse buffers (`sync.Pool`), prefer value types, preallocate slices with known capacity, avoid unnecessary string concatenation (use `strings.Builder`).
↳ Does fewer pointers help GC? Yes — fewer scanned objects and smaller live heap.

63. What is the cost of interface values for GC?
**A:** Interface values are two words (type, data). If data is a pointer, GC must scan it. Passing large values as interfaces may cause boxing and heap allocation.
↳ What is boxing? Storing a value type in an interface may allocate on heap if value doesn't fit in pointer word.

---

## Architecture & Design (M23, M25)

*See [M23](../modules/23-design-patterns-go.md), [M25](../modules/25-project-architecture.md)*

64. What is the standard Go project layout?
**A:** Common convention: `cmd/` for entry points, `internal/` for private packages, `pkg/` for public libraries (optional), `api/` for protos/OpenAPI, `configs/`, `scripts/`. Not enforced by compiler but widely adopted.
↳ What goes in `cmd/server/main.go`? Wiring only — parse config, construct dependencies, start server.

65. How do you implement dependency injection in Go?
**A:** Constructor injection — `func NewService(repo Repository, log Logger) *Service`. Manual wiring in `main` or use wire/fx for codegen/DI frameworks. Interfaces define boundaries for testing.
↳ Why avoid global singletons? Hard to test, hidden dependencies, init order issues.

66. What is the functional options pattern?
**A:** `func NewServer(opts ...Option) *Server` where each `Option` is `func(*Server)`. Provides optional, readable configuration without telescoping constructors.
↳ When prefer explicit struct config? When options are required and validation is complex — struct is clearer.

67. How do you manage configuration in Go services?
**A:** Layered: defaults → config file (YAML/JSON) → environment variables → flags. Libraries: `viper`, `envconfig`. Validate at startup; fail fast on missing required config.
↳ How do you handle secrets? Environment variables or secret managers (Vault, AWS SM) — never commit secrets.

68. What is middleware in HTTP handlers?
**A:** Functions wrapping `http.Handler` to add cross-cutting concerns: logging, auth, recovery, tracing. Chain with `handler = middleware1(middleware2(actualHandler))`.
↳ How does middleware access request context? `r.Context()` — add values via context or attach to custom request struct.

69. What HTTP patterns are idiomatic in Go?
**A:** `net/http` with `http.Handler` interface, chi/gin routers for routing groups, handler per resource, `context.Context` per request, structured JSON responses, explicit status codes, graceful shutdown with `Server.Shutdown`.
↳ REST vs RPC-style in Go? REST for public APIs; internal can use gRPC. Avoid `POST /getUser` RPC-on-HTTP smell in public APIs.

70. How does distributed tracing integrate with Go?
**A:** OpenTelemetry SDK — propagate trace context via HTTP/gRPC headers (`traceparent`), create spans around handlers and outbound calls, export to Jaeger/Tempo/Datadog.
↳ Where do you start/end spans? Middleware starts server span; each outbound call gets a child span.

71. When do you choose gRPC over REST?
**A:** gRPC for internal service-to-service: strong contracts (protobuf), binary efficiency, streaming, codegen. REST for browser-facing, public APIs, caching, human debugging.
↳ How do browsers call gRPC? grpc-web with Envoy proxy, or expose REST gateway (gRPC-Gateway).

72. How does semantic versioning work in Go modules?
**A:** v0/v1: module path unchanged. v2+: must include `/v2` in module path and import paths. Tag releases `v1.2.3`; `go get module@v1.2.3` pins version.
↳ What is a breaking change in Go modules? Any change requiring importers to update code — exported API removal, behavior change. Major version bump required.

73. What is the repository pattern in Go?
**A:** Interface defining data access (`GetUser`, `SaveOrder`); implementation wraps `database/sql` or ORM. Business logic depends on interface — easy to mock in tests.
↳ Where does the interface live? Consumer package (service) defines the interface it needs — Go idiom, not provider-defined.

74. What is clean/hexagonal architecture in Go?
**A:** Domain at center, ports (interfaces) define boundaries, adapters implement HTTP/DB/external APIs. `internal/domain`, `internal/ports`, `internal/adapters` — dependency points inward.
↳ Is this overkill for small services? Yes — start simple, extract when complexity grows.

75. How do you version APIs?
**A:** URL prefix (`/v1/`), header negotiation, or separate modules. Deprecate with headers and documentation; maintain backward compatibility within major version.
↳ How does this relate to Go module versions? API version (HTTP) is independent of Go module semver.

---

## Performance Tuning (M22)

*See [M22](../modules/22-gc-runtime-performance.md), [M18](../modules/18-testing-benchmarks-fuzzing.md)*

76. What is the difference between benchmarking and profiling?
**A:** Benchmarks (`go test -bench`) measure function/package performance in isolation with controlled inputs. Profiling (CPU, memory, trace) samples running programs to find where time/memory is actually spent in realistic workloads.
↳ When to use each? Benchmark for regression detection; profile for finding real bottlenecks.

77. How do you write a Go benchmark?
**A:** `func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { X() } }`. Run with `go test -bench=. -benchmem`. Use `b.ResetTimer()`, `b.StopTimer()` for setup.
↳ What is `b.N`? Auto-tuned iteration count for statistical stability.

78. What tools profile Go programs?
**A:** `pprof` (CPU, heap, goroutine, mutex, block), `go tool trace` (scheduler/latency), `runtime/pprof`, HTTP `/debug/pprof/`, `go test -cpuprofile`, Datadog/Pyroscope continuous profiling.
↳ What does the flame graph show? Inclusive time per call stack — wide bars are hot paths.

79. How do you optimize CPU-bound Go code?
**A:** Profile first. Reduce allocations, improve algorithms, use `sync.Pool` for buffers, parallelize with worker pools (`GOMAXPROCS` cores), avoid unnecessary interface boxing and reflection.
↳ When does parallelization hurt? Small work items — goroutine overhead exceeds benefit.

80. How do you detect goroutine leaks?
**A:** `go test -race` won't catch leaks. Use goroutine profile: compare count over time, look for blocked stacks on channels. Test with `goleak` library in integration tests.
↳ What does a leaked goroutine stack show? Blocked on `chan send` or `chan receive` indefinitely.

81. What is lock contention and how do you reduce it?
**A:** Multiple goroutines competing for the same mutex — profile with `mutex` profile. Fix: shrink critical section, `RWMutex` for reads, shard locks, lock-free atomics for simple counters.
↳ What does mutex profile show? Cumulative time spent waiting for locks.

82. How do you find performance bottlenecks systematically?
**A:** Measure end-to-end latency → CPU profile hot paths → memory profile allocations → trace for scheduling gaps → check DB/network with spans. Optimize the largest slice first.
↳ What is the biggest mistake? Optimizing without measuring — micro-benchmarking the wrong function.

83. How does GC impact performance?
**A:** GC consumes CPU for marking/sweeping and causes STW pauses (brief). High allocation rate increases GC frequency. Reduce live heap and allocation rate to reduce GC overhead.
↳ What latency percentile does GC affect? Tail latency (p99) — occasional STW and assist marking.

84. How do you use pprof in production?
**A:** Import `net/http/pprof`, protect endpoint with auth/network policy, or use continuous profiler. Sample at low rate to minimize overhead.
↳ Security concern? `/debug/pprof` exposes internals — never expose publicly.

85. How do you optimize I/O-bound Go services?
**A:** Enough goroutines to saturate I/O (not one per request on server — use goroutine per request for simplicity up to scale, then worker pools). Connection pooling, batch requests, context timeouts, keep-alive, appropriate buffer sizes.
↳ What is the right number of DB connections? `max_open ≈ (cores × 2) + effective_spindle_count` rule of thumb — tune with load test.

86. What is `benchstat`?
**A:** Tool comparing benchmark results across commits — reports statistical significance of performance changes. Essential for CI performance regression detection.
↳ How many benchmark runs? At least 5–10 with `-count=10` for stable comparison.

87. When should you use `unsafe` for performance?
**A:** Almost never in application code — zero-copy string/byte conversions (`unsafe.String`, Go 1.20+) in hot paths only after profiling proves benefit. Breaks GC assumptions if misused.
↳ Example legitimate use? Interop with C via cgo, or stdlib-level optimizations.

---

## Networking (M19, M21)

*See [M19](../modules/19-networking-http.md), [M21](../modules/21-grpc-protobuf.md)*

88. How does TCP differ from HTTP in Go's networking stack?
**A:** TCP (`net` package) is byte-stream transport — you manage framing/protocol. HTTP (`net/http`) is application-layer on top of TCP with request/response semantics, headers, and convenience APIs.
↳ When use raw TCP? Custom protocols, gaming, high-performance internal RPC before HTTP overhead matters.

89. How does `net/http` handle concurrency?
**A:** `ListenAndServe` spawns a goroutine per connection (HTTP/1.1) or per request. Handlers must be thread-safe — no shared mutable state without sync.
↳ Is `DefaultServeMux` safe? Yes for registration at init; handlers must be safe for concurrent calls.

90. What is graceful shutdown for HTTP servers?
**A:** On SIGTERM, call `server.Shutdown(ctx)` — stops accepting new connections, waits for in-flight requests up to context deadline, then closes.
↳ What happens to long-running requests? They get context canceled if you wire `r.Context()`; otherwise wait until deadline.

91. What are WebSockets and how do you use them in Go?
**A:** Full-duplex communication over HTTP upgrade. Library: `gorilla/websocket` or `nhooyr.io/websocket`. Server upgrades connection; both sides send frames.
↳ WebSockets vs SSE? SSE is server→client one-way over HTTP; simpler for feeds. WebSockets for bidirectional.

92. What is HTTP/2 and does Go support it?
**A:** Multiplexed streams over single TCP connection, header compression, server push (rarely used). Go enables HTTP/2 by default for HTTPS in `net/http`.
↳ How do you force HTTP/2 in development? Use TLS — h2c (HTTP/2 cleartext) requires explicit configuration.

93. What is connection pooling in HTTP clients?
**A:** `http.Client` reuses TCP connections via `Transport`'s idle connection pool (`MaxIdleConns`, `MaxIdleConnsPerHost`). Creating new `Client` per request defeats pooling.
↳ What happens without pooling? TCP + TLS handshake per request — high latency.

94. What is the difference between synchronous and asynchronous HTTP handling?
**A:** Sync: handler blocks until work completes (standard `ServeHTTP`). Async: handler returns immediately, processes in background goroutine — must handle client disconnect and backpressure carefully.
↳ When is async handler appropriate? Long-polling, fire-and-forget with webhook callback, SSE streams.

95. How does TLS work in Go servers?
**A:** `http.ListenAndServeTLS(addr, certFile, keyFile, handler)` or custom `tls.Config` on `http.Server`. Minimum version, cipher suites, and cert rotation via config or sidecar (cert-manager).
↳ How do you skip TLS verify in dev? `tls.Config{InsecureSkipVerify: true}` — never in production.

96. What is a reverse proxy in Go?
**A:** `httputil.ReverseProxy` forwards requests to upstream servers, handles hop-by-hop headers, can modify requests/responses. Used in API gateways, load balancing layers.
↳ How do you propagate tracing through proxy? Copy `traceparent` headers to upstream request.

97. How do gRPC interceptors work?
**A:** Middleware for gRPC — unary/stream interceptors for logging, auth, tracing, recovery. Chain with `grpc.ChainUnaryInterceptor(a, b, c)`.
↳ Equivalent in HTTP? `http.Handler` middleware wrapping pattern.

98. What are gRPC streaming modes?
**A:** Unary (request-response), server streaming, client streaming, bidirectional streaming. Go generates `Send`/`Recv` methods on streams.
↳ When use streaming? Large datasets, real-time feeds, upload chunking.

99. How do you set timeouts on HTTP clients?
**A:** `http.Client{Timeout: 30 * time.Second}` covers entire exchange. Finer control: `context.WithTimeout` on `http.NewRequestWithContext`, or `Transport.ResponseHeaderTimeout`.
↳ Difference between Client.Timeout and context? Client.Timeout is coarse overall; context cancels at any phase and propagates to server if supported.

---

## Databases & ORM (M20)

*See [M20](../modules/20-databases-orm.md)*

100. What does `sql.Open` actually do?
**A:** Creates a `*sql.DB` connection **pool** — it does not connect immediately. First operation establishes connections. Configure pool with `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`.
↳ What if you forget to close DB? Connections leak until process exit — call `db.Close()` on shutdown.

101. How do database transactions work in Go?
**A:** `tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})`. All operations use `tx`. `Commit()` or `Rollback()`. Always rollback on error (defer).
↳ What happens if you don't commit? Connection returned to pool with open transaction — blocks and leaks.

102. When do you use `Query` vs `QueryRow` vs `Exec`?
**A:** `Query` — multiple rows (`rows.Next()` loop). `QueryRow` — single row expected, returns on `Scan`. `Exec` — INSERT/UPDATE/DELETE without row data, check `RowsAffected`.
↳ What if `QueryRow` finds no rows? `Scan` returns `sql.ErrNoRows`.

103. What are the pros and cons of GORM?
**A:** Pros: rapid CRUD, associations, hooks, migrations (AutoMigrate), conventions. Cons: hidden SQL (N+1), magic defaults, harder to optimize complex queries, reflection overhead, soft-delete surprises.
↳ How do you debug GORM SQL? `db.Debug()` or logger mode in config.

104. How does sqlx differ from database/sql?
**A:** sqlx embeds `*sql.DB`, adds struct scanning (`Get`, `Select`), named queries, `In` clause expansion. Still write SQL — no ORM magic.
↳ When choose sqlx over GORM? SQL-heavy services wanting ergonomics without ORM abstraction.

105. What is the N+1 query problem?
**A:** Loading N parent records then one query per child in a loop. Fix with JOIN, batch `WHERE id IN (...)`, or ORM preloading (`Preload` in GORM).
↳ How do you detect N+1? Enable query logging, count queries per request in tests.

106. How do you optimize SQL queries in Go services?
**A:** EXPLAIN ANALYZE, proper indexes, avoid `SELECT *`, pagination with keyset not offset for large tables, prepared statements (automatic in `database/sql`), batch inserts, read replicas for reads.
↳ What is keyset pagination? `WHERE id > $last_id ORDER BY id LIMIT 20` — stable vs large OFFSET.

107. How do connection pools affect database performance?
**A:** Too few — goroutines wait. Too many — overwhelm DB. Tune `MaxOpenConns` ≤ DB limit, `ConnMaxLifetime` for load balancer failover, `ConnMaxIdleTime` to release unused.
↳ Should MaxOpenConns equal number of goroutines? No — goroutines multiplex; pool size ≈ concurrent queries.

108. How do you handle database migrations?
**A:** Tools: `golang-migrate`, `goose`, `atlas`. Versioned SQL files, run on deploy. GORM `AutoMigrate` for prototyping only — not for production schema management alone.
↳ How do you do zero-downtime migrations? Expand-contract pattern — add column nullable, backfill, add constraint.

109. What are common database/sql pitfalls?
**A:** Not closing `rows.Close()`, ignoring `rows.Err()` after loop, not using context variants, SQL injection via string concat (use parameters), connection pool misconfiguration, long transactions holding locks.
↳ Why close rows? Unclosed rows hold connections from pool — pool exhaustion.

110. What are soft deletes and their trade-offs?
**A:** `deleted_at` timestamp instead of DELETE. GORM supports via `gorm.DeletedAt`. Pros: audit, recovery. Cons: unique constraints complicated, queries must filter `deleted_at IS NULL`, table bloat.
↳ How do unique constraints work with soft delete? Composite unique on `(email, deleted_at)` or partial index.

111. How do you handle long-running queries?
**A:** Pass `context` with timeout, use `QueryContext`, cancel on client disconnect. For reports, use read replica, materialized views, or async job with result polling.
↳ What does DB see on context cancel? Driver sends cancel request (PostgreSQL `pg_cancel_backend`).

112. How do you test database code in Go?
**A:** Interface-based repository with mock, `sqlmock` for SQL expectations, testcontainers for integration tests with real Postgres/MySQL.
↳ When is sqlmock enough? Unit testing SQL construction and scan logic without real DB.

---

## Security & Distributed (M27)

*See [M27](../modules/27-security-distributed.md)*

113. How do you handle authentication in Go APIs?
**A:** JWT middleware (validate signature, expiry, claims), OAuth2/OIDC for delegated auth, API keys for service-to-service. Store secrets in env/Vault, never in code.
↳ Where validate JWT? Middleware before handlers — attach claims to context.

114. What is the principle of least privilege for Go services?
**A:** DB user with minimal permissions, IAM roles not static keys, network policies in K8s, read-only filesystem in containers, non-root user.
↳ How do you prevent SQL injection? Parameterized queries always — never `fmt.Sprintf` for SQL values.

115. How do you manage secrets in Go?
**A:** Environment variables at runtime, HashiCorp Vault, cloud secret managers, K8s Secrets (encrypted at rest). Use `os.Getenv` or dedicated config loader — no secrets in `go.mod` or git.
↳ What about embedding secrets in binaries? Reversible — anyone can extract.

116. What is mutual TLS (mTLS)?
**A:** Both client and server present certificates — common in service meshes (Istio) and internal gRPC. Go: `tls.Config{ClientAuth: tls.RequireAndVerifyClientCert}`.
↳ When is TLS alone enough? Public APIs where only server identity matters.

117. Explain CAP theorem in practical terms.
**A:** In a partition, choose Consistency (CP — reject requests) or Availability (AP — serve possibly stale data). Go services must handle retries, idempotency, and conflict resolution based on chosen trade-off.
↳ Is PostgreSQL CP or AP? CP by default — synchronous replication tunable.

118. What is distributed consensus and where is it used?
**A:** Agreement among nodes despite failures — Raft (etcd, Consul), Paxos. Used for leader election, distributed locks, config management.
↳ Go client for etcd? `go.etcd.io/etcd/client/v3` for distributed locking with leases.

119. How do you implement rate limiting?
**A:** Token bucket (`golang.org/x/time/rate`) in middleware, or Redis-backed for distributed rate limiting across instances.
↳ What headers to return? `X-RateLimit-Limit`, `Retry-After` on 429.

120. What are common security headers for HTTP APIs?
**A:** `Strict-Transport-Security`, `Content-Security-Policy`, `X-Content-Type-Options: nosniff`, `X-Frame-Options`, CORS configured explicitly (not `*` with credentials).
↳ How set in Go? Middleware setting `w.Header().Set(...)`.

121. How do you prevent timing attacks in Go?
**A:** Use `crypto/subtle.ConstantTimeCompare` for secret comparison, `crypto/hmac` for MAC verification. Avoid early return on first byte mismatch for secrets.
↳ Is `==` on strings safe for passwords? No — use `bcrypt` or `argon2` via `golang.org/x/crypto`.

122. What is input validation best practice?
**A:** Validate at API boundary — type, length, format (use `validator` tags or custom). Sanitize output for XSS in HTML contexts. Reject unknown JSON fields if strictness needed.
↳ Library? `github.com/go-playground/validator/v10` with struct tags.

---

## Microservices (M26)

*See [M26](../modules/26-microservices-go.md)*

123. How do you define service boundaries in Go microservices?
**A:** Align with business capabilities (bounded contexts), not technical layers. Each service owns its data, deploys independently, communicates via APIs/events. Start modular monolith, split when scaling teams or domains justify it.
↳ What is the strangler fig pattern? Gradually replace monolith routes with new services behind a proxy.

124. What is the Saga pattern?
**A:** Distributed transaction as a sequence of local transactions with compensating actions on failure. Choreography (events) or orchestration (central coordinator).
↳ Example compensation? `CreateOrder` fails after `ReserveInventory` → run `ReleaseInventory`.

125. What is CQRS?
**A:** Command Query Responsibility Segregation — separate write model (commands) from read model (optimized queries). Often paired with event sourcing but not required.
↳ When is CQRS worth it? Read/write load profiles differ dramatically; complex domain with many query shapes.

126. What is event sourcing?
**A:** Store state changes as append-only events; current state is replay/reduce. Enables audit trail and temporal queries. Complexity: snapshots, schema evolution, eventual consistency.
↳ Go libraries? Custom event store, or Kafka/NATS as log with consumer projections.

127. How do services discover each other?
**A:** DNS/K8s services, Consul/etcd, or service mesh (Istio/Linkerd). Client-side load balancing with gRPC resolver or round-robin HTTP.
↳ What is the sidecar pattern? Proxy container alongside app — handles mTLS, retries, metrics.

128. How do you handle partial failures in microservices?
**A:** Timeouts, retries (idempotent ops only), circuit breakers, fallbacks, bulkheads (isolate connection pools per dependency).
↳ What is a bulkhead? Separate resource pools so one failing dependency doesn't exhaust all connections.

129. What is an API gateway's role?
**A:** Single entry point — routing, auth, rate limiting, SSL termination, request aggregation. Implemented in Go with custom server or Kong/Envoy.
↳ When skip gateway? Internal mesh with mTLS and few clients — direct service calls may suffice.

130. How do you implement health checks?
**A:** Liveness (`/healthz` — process up) and readiness (`/readyz` — can serve traffic, DB connected). K8s probes use these. Return 503 when not ready.
↳ Should health check the database? Readiness yes; liveness no — avoid restart loops on DB outage.

131. What is the outbox pattern?
**A:** Write business data and outbound event in same DB transaction (outbox table). Separate poller publishes to message broker — guarantees at-least-once delivery without dual-write problem.
↳ Go implementation? Transaction inserts row; background goroutine polls and publishes.

132. How do you version microservice APIs?
**A:** URL versioning, header versioning, or protobuf package versioning for gRPC. Maintain backward compatibility; deprecate with metrics on old version usage.
↳ gRPC breaking change? New `package v2` and new service definitions.

---

## Tooling (M28)

*See [M28](../modules/28-tooling-cicd-containers.md)*

133. What linters should Go projects use?
**A:** `go vet` (stdlib), `staticcheck` (comprehensive), `golangci-lint` (aggregator), `go fmt`/`goimports` (formatting). Run in CI on every PR.
↳ What does staticcheck catch beyond vet? Unused code, simplifications, API misuse, deprecated functions.

134. How do you structure a CI pipeline for Go?
**A:** Steps: checkout → setup Go (cache modules) → `go vet ./...` → `golangci-lint` → `go test -race -cover ./...` → build binary → (optional) integration tests → docker build → deploy.
↳ How cache Go modules in CI? `actions/cache` on `~/go/pkg/mod` and build cache.

135. What makes a good Go Dockerfile?
**A:** Multi-stage build: compile with `golang` image, copy static binary to `distroless` or `alpine`. `CGO_ENABLED=0` for static binary. Non-root user. Minimal attack surface.
↳ Why distroless? No shell — smaller image, fewer CVEs.

136. How do you cross-compile Go binaries?
**A:** `GOOS=linux GOARCH=amd64 go build -o app ./cmd/server`. Go supports cross-compilation natively (cgo excepted).
↳ What about ARM (Apple Silicon / Graviton)? `GOARCH=arm64` — same command pattern.

137. What is a Go workspace (`go.work`)?
**A:** Multi-module development — local replace without editing `go.mod`. `go work init ./module-a ./module-b`. Useful for monorepos.
↳ When not use workspaces? Published libraries — use normal module dependencies.

138. How do you inject build version info?
**A:** `-ldflags "-X main.version=1.2.3 -X main.commit=$(git rev-parse HEAD)"` at build time. Expose via `/version` endpoint or CLI flag.
↳ Where define variables? Package-level `var version string` in `main` or `internal/build`.

139. What Kubernetes probes should Go services implement?
**A:** Liveness (restart if deadlocked), readiness (remove from LB if DB down), startup (slow init). Use `http.Server` with short probe-specific handlers.
↳ Graceful shutdown interaction? `preStop` hook + `Shutdown` drains before pod terminates.

140. How do you run Go in Docker with proper signal handling?
**A:** Use exec form `ENTRYPOINT ["/server"]` so PID 1 is your binary. Handle `SIGTERM` in Go, call `server.Shutdown`. Set `STOPSIGNAL SIGTERM` and adequate `terminationGracePeriodSeconds`.
↳ What if you use shell form ENTRYPOINT? Signals may not reach Go process — zombie/orphan issues.

---

## Scenario Prompts

Design and troubleshooting scenarios for senior-level interviews. Practice structuring answers: requirements → trade-offs → design → failure modes.

### Scenario 1: Rate-Limited API Gateway
**Prompt:** Design a Go API gateway that routes to 5 backend services with per-client rate limiting (1000 req/min), JWT auth, and request logging.

**Key points:** chi/gin router, JWT middleware, Redis token bucket for distributed rate limit, `ReverseProxy` per upstream, OpenTelemetry tracing, graceful shutdown, health checks, config via env.

↳ How do you test rate limiting? Integration test with burst requests, assert 429 after threshold.

### Scenario 2: Worker Pool for Image Processing
**Prompt:** Process 10,000 images concurrently without exhausting memory — resize and upload to S3.

**Key points:** Bounded worker pool (e.g., 20 workers), job channel buffered, context for cancellation, `errgroup` for error propagation, stream from disk not load all in memory, semaphore if multiple resources.

↳ What if one image corrupts? Per-job error handling, dead-letter queue, don't fail entire batch.

### Scenario 3: Database Connection Storm
**Prompt:** Service scales from 2 to 50 pods and PostgreSQL starts rejecting connections.

**Key points:** `SetMaxOpenConns` per pod (total < PG `max_connections`), PgBouncer connection pooler, reduce idle conns, `ConnMaxLifetime`, readiness probe checks pool health.

↳ Formula? 50 pods × 10 conns = 500 — exceeds typical PG limit. Target ~20 total via PgBouncer.

### Scenario 4: Goroutine Leak in Production
**Prompt:** Memory grows 50 MB/day, goroutine count climbs to 100k. Diagnose and fix.

**Key points:** `pprof/goroutine` profile, find blocked on channel, trace to handler spawning goroutine without exit on client disconnect, fix with `context.Cancel`, add `goleak` test.

↳ Common cause? SSE/long-poll handler not checking `r.Context().Done()`.

### Scenario 5: Breaking Up a Go Monolith
**Prompt:** E-commerce monolith — extract inventory service. Migration plan?

**Key points:** Strangler fig, define bounded context, extract DB tables, dual-write or CDC, gRPC API, feature flag traffic shift, Saga for order-inventory consistency.

↳ How handle distributed transactions? Saga with `ReserveInventory` / `ReleaseInventory` compensation.

### Scenario 6: Optimizing P99 Latency
**Prompt:** HTTP API p50=20ms, p99=2s. Go service, Postgres backend.

**Key points:** Trace slow requests, check N+1 queries, missing indexes, lock contention, GC pauses (trace), external API timeouts, connection pool wait, add caching for hot reads.

↳ First tool? Distributed tracing waterfall — find the 2s span.

### Scenario 7: Secure Internal gRPC Mesh
**Prompt:** 10 internal Go services — mTLS, auth, observability.

**Key points:** Shared CA, per-service certs, gRPC interceptors for authZ, OpenTelemetry propagators, service mesh optional (Istio), JWT for user context propagation on edge.

↳ Service mesh vs manual TLS? Mesh for large deployments; manual fine for <10 services.

### Scenario 8: Config Hot Reload
**Prompt:** Feature flags change without redeploying Go services.

**Key points:** Poll config service / watch etcd, `atomic.Value` or `RWMutex` for config snapshot, graceful feature toggle, validate before swap, log config changes.

↳ Race during reload? Copy-on-write — build new struct, atomic swap pointer.

### Scenario 9: Handling Duplicate Events
**Prompt:** Payment service receives duplicate webhook deliveries (at-least-once).

**Key points:** Idempotency key in DB (unique constraint), check-before-process, store webhook ID, return 200 on duplicate, transactional outbox for downstream.

↳ Schema? `idempotency_keys(key PRIMARY KEY, processed_at, result)`.

### Scenario 10: CI/CD for Go Monorepo
**Prompt:** 5 Go services in one repo — efficient CI.

**Key points:** Path-filtered workflows (only test changed modules), `go work` or matrix per module, shared golangci-lint config, module cache, build images per `cmd/`, semantic release tags per service.

↳ Tooling? `go list -m` + git diff to detect changed packages; Turborepo-style or custom script.

---

## Quick Topic Index

| Topic | Questions | Module |
|-------|-----------|--------|
| Fundamentals | 1–26 | M01–M04 |
| Concurrency | 27–51 | M10–M14 |
| Memory / GC | 52–63 | M22 |
| Architecture | 64–75 | M23, M25 |
| Performance | 76–87 | M22, M18 |
| Networking | 88–99 | M19, M21 |
| Databases | 100–112 | M20 |
| Security | 113–122 | M27 |
| Microservices | 123–132 | M26 |
| Tooling | 133–140 | M28 |
| Scenarios | S1–S10 | All |

**Total: 140 numbered questions + 10 scenario prompts**
