# Go Study Guide

A broad, **interview-prep-depth** curriculum for the Go programming language — from language fundamentals through concurrency, runtime internals, networking, databases, microservices, and production tooling. Hands-on labs reinforce every major topic.

**Estimated study time:** ~60 hours over 8 weeks.

Every module follows the same interview-prep template: **TL;DR → Concept → How It Really Works (internals + diagram) → Why/When/Trade-offs → Worked Scenario → Gotchas → Interview Q&A → Verify → Further Reading**.

---

## Prerequisites

| Tool | Purpose | Install (Windows) |
|------|---------|-------------------|
| Go 1.22+ | Compiler, toolchain, stdlib | [go.dev/dl](https://go.dev/dl/) or `choco install golang` |
| Git | Clone repos, module proxy auth | `choco install git` |
| Docker Desktop | Container labs (gRPC, microservices, K8s) | [docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop/) |
| VS Code / Cursor + Go extension | IDE, debugging, `gopls` | Extension: `golang.go` |
| `protoc` + plugins (optional) | gRPC lab code generation | `choco install protobuf` |
| PostgreSQL client (optional) | Database lab inspection | Bundled in Docker image or `choco install postgresql` |

You should be comfortable with at least one programming language, basic command-line usage, and HTTP concepts. No prior Go experience required.

---

## Quick Start

Verify your toolchain:

```bash
cd labs/00-setup/hello
go run .
```

Expected output: `Hello, Go!`

Then explore modules starting with [Module 1 — Go Fundamentals](modules/01-go-fundamentals.md).

---

## Curriculum (28 modules)

### Group 1 — Language Fundamentals

| # | Module | Topics |
|---|--------|--------|
| 1 | [Go Fundamentals & Tooling](modules/01-go-fundamentals.md) | Go modules, GOPATH, project layout, `go` CLI, `gofmt`, `go vet` |
| 2 | [Types, Variables & Control Flow](modules/02-types-control-flow.md) | Basic types, zero values, constants, `if`/`for`/`switch`, type conversions |
| 3 | [Functions, Methods & defer/panic/recover](modules/03-functions-methods.md) | Functions, methods, variadics, closures, defer LIFO, panic/recover |
| 4 | [Pointers & Value Semantics](modules/04-pointers-value-semantics.md) | Pointers, `new`, pass-by-value vs reference, escape analysis intro |

### Group 2 — Data Structures & Abstractions

| # | Module | Topics |
|---|--------|--------|
| 5 | [Structs & Composition](modules/05-structs-composition.md) | Structs, embedding, promoted fields, anonymous vs named types |
| 6 | [Interfaces & Polymorphism](modules/06-interfaces-polymorphism.md) | Implicit interfaces, type assertions, empty interface, `io` patterns |
| 7 | [Slices & Maps Internals](modules/07-slices-maps-internals.md) | Slice header, capacity, append, map internals, iteration |
| 8 | [Strings, Runes & Encoding](modules/08-strings-runes-encoding.md) | UTF-8, runes, `range` over strings, `strings`/`bytes` packages |
| 9 | [Error Handling](modules/09-error-handling.md) | `error` interface, wrapping, `errors.Is`/`As`, custom errors |

### Group 3 — Concurrency

| # | Module | Topics |
|---|--------|--------|
| 10 | [Goroutines & the Scheduler](modules/10-goroutines-scheduler.md) | Goroutines vs threads, GMP model, `GOMAXPROCS`, scheduler internals |
| 11 | [Channels & select](modules/11-channels-select.md) | Buffered vs unbuffered, close semantics, `select`, pipelines |
| 12 | [Sync Primitives & Memory Model](modules/12-sync-memory-model.md) | Mutex, RWMutex, WaitGroup, `sync.Once`, atomic, happens-before |
| 13 | [context Package](modules/13-context-package.md) | Cancellation, deadlines, values, propagation, HTTP integration |
| 14 | [Concurrency Patterns](modules/14-concurrency-patterns.md) | Worker pools, fan-in/fan-out, errgroup, semaphores, leaks |

### Group 4 — Advanced Language Features

| # | Module | Topics |
|---|--------|--------|
| 15 | [Generics](modules/15-generics.md) | Type parameters, constraints, generic data structures |
| 16 | [Reflection & Code Generation](modules/16-reflection-codegen.md) | `reflect` package, struct tags, `go:generate`, protobuf |
| 17 | [CGO, FFI & Unsafe](modules/17-cgo-ffi.md) | C interop, `unsafe` pointer rules, when to avoid cgo |
| 18 | [Testing, Benchmarks & Fuzzing](modules/18-testing-benchmarks-fuzzing.md) | Table tests, mocks, benchmarks, fuzzing, coverage |

### Group 5 — Networking & Data

| # | Module | Topics |
|---|--------|--------|
| 19 | [Networking & HTTP](modules/19-networking-http.md) | TCP, HTTP servers, middleware, WebSockets, HTTP/2, TLS, proxies |
| 20 | [Databases & ORM](modules/20-databases-orm.md) | `database/sql`, sqlx, GORM, transactions, pooling, migrations |
| 21 | [gRPC & Protobuf](modules/21-grpc-protobuf.md) | Protobuf, gRPC vs REST, streaming, interceptors, codegen |

### Group 6 — Runtime & Performance

| # | Module | Topics |
|---|--------|--------|
| 22 | [GC, Runtime & Performance](modules/22-gc-runtime-performance.md) | GC tri-color mark-sweep, escape analysis, pprof, optimization |

### Group 7 — Architecture & Design

| # | Module | Topics |
|---|--------|--------|
| 23 | [Design Patterns in Go](modules/23-design-patterns-go.md) | Functional options, DI, repository, strategy, Go-idiomatic patterns |
| 24 | [Observability](modules/24-observability.md) | Structured logging, metrics, tracing, OpenTelemetry |
| 25 | [Project Architecture](modules/25-project-architecture.md) | Layered/clean/hexagonal, config, module versioning, API design |

### Group 8 — Production & Operations

| # | Module | Topics |
|---|--------|--------|
| 26 | [Microservices in Go](modules/26-microservices-go.md) | Service boundaries, Saga, CQRS, service mesh, resilience |
| 27 | [Security & Distributed Systems](modules/27-security-distributed.md) | Auth, TLS, secrets, CAP, consensus, distributed locks |
| 28 | [Tooling, CI/CD & Containers](modules/28-tooling-cicd-containers.md) | Docker, Kubernetes, GitHub Actions, linting, release |

---

## Hands-On Labs

Complete in order; each lab has its own `README.md` with instructions.

| Lab | Folder | Topic |
|-----|--------|-------|
| 00 | [labs/00-setup](labs/00-setup/) | Toolchain verification (`go run`, modules) |
| 01 | [labs/01-basics](labs/01-basics/) | Types, pointers, defer, variadic functions |
| 02 | [labs/02-slices-maps](labs/02-slices-maps/) | Slices, maps, internals |
| 03 | [labs/03-interfaces](labs/03-interfaces/) | Interfaces, polymorphism, type assertions |
| 04 | [labs/04-errors](labs/04-errors/) | Custom errors, wrapping, `errors.Is`/`As` |
| 05 | [labs/05-goroutines-channels](labs/05-goroutines-channels/) | Goroutines, channels, `select` |
| 06 | [labs/06-sync-memory-model](labs/06-sync-memory-model/) | Mutex, WaitGroup, race conditions |
| 07 | [labs/07-context](labs/07-context/) | Context cancellation and timeouts |
| 08 | [labs/08-testing-benchmarks](labs/08-testing-benchmarks/) | Unit tests, table tests, benchmarks |
| 09 | [labs/09-http-api](labs/09-http-api/) | HTTP server, routing, middleware |
| 10 | [labs/10-databases](labs/10-databases/) | `database/sql`, queries, transactions |
| 11 | [labs/11-profiling](labs/11-profiling/) | CPU/memory profiling with pprof |
| 12 | [labs/12-concurrency-patterns](labs/12-concurrency-patterns/) | Worker pools, fan-in/fan-out |
| 13 | [labs/13-design-patterns](labs/13-design-patterns/) | Functional options, repository pattern |
| 14 | [labs/14-microservices](labs/14-microservices/) | Multi-service HTTP communication |
| 15 | [labs/15-grpc](labs/15-grpc/) | gRPC server/client with Protobuf |
| 16 | [labs/16-docker](labs/16-docker/) | Docker, docker-compose, Kubernetes manifests |

---

## Cheatsheets

- [Interview question bank](cheatsheets/interview-questions.md) — 120+ Q&A organized by topic with follow-ups
- [Go cheatsheet](cheatsheets/golang-cheatsheet.md) — syntax, types, concurrency, CLI, popular libraries
- [Decision cheatsheet](cheatsheets/decision-cheatsheet.md) — channels vs mutex, ORM picker, gRPC vs REST, deployment checklist

---

## Suggested Schedule

| Week | Modules | Labs | Hours |
|------|---------|------|-------|
| 1 | 1–4 | 00–01 | ~8 |
| 2 | 5–9 | 02–04 | ~8 |
| 3 | 10–12 | 05–06 | ~8 |
| 4 | 13–15 | 07–08 | ~8 |
| 5 | 16–18 | 09 | ~7 |
| 6 | 19–21 | 10, 15 | ~8 |
| 7 | 22–25 | 11–13 | ~8 |
| 8 | 26–28 + interview bank | 14–16 | ~8 |

---

## How Each Module Is Structured

1. **TL;DR** — the mental model in a few lines
2. **Concept** — plain-language explanation
3. **How It Really Works** — internals with a diagram
4. **Why / When / Trade-offs** — senior decision-making
5. **Worked Scenario** — a realistic situation end-to-end
6. **Gotchas & Failure Modes** — what bites people in production
7. **Interview Q&A** — sharp answers + follow-ups
8. **Verify** — commands and code to run yourself
9. **Further Reading** — authoritative sources

---

## Official References

- [The Go Programming Language Specification](https://go.dev/ref/spec)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Blog](https://go.dev/blog/)
- [pkg.go.dev](https://pkg.go.dev/) — standard library and module documentation
- *The Go Programming Language* — Donovan & Kernighan
- *Concurrency in Go* — Katherine Cox-Buday
- *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi
