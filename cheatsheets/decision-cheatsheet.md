# Go Decision Cheatsheet

Quick decision matrices for architecture, concurrency, data access, and deployment. Pair with the [interview question bank](interview-questions.md) and [module guides](../README.md#curriculum).

---

## Channels vs Mutex

| Factor | Prefer **Channels** | Prefer **Mutex** |
|--------|---------------------|------------------|
| Primary goal | Coordinate goroutines, pass ownership of work/data | Protect shared mutable state |
| Ownership model | "Don't communicate by sharing memory" | "Share memory, communicate by locking" |
| Complexity | Fan-in/fan-out, pipelines, select-based cancellation | Simple counter, cache, map guarded by one lock |
| Performance | Overhead of send/recv; great for orchestration | Lower overhead for short critical sections |
| Deadlock risk | Blocked sends/receives, forgotten close | Lock ordering, forgotten unlock (use `defer`) |
| Testing | Harder to reason about ordering | Straightforward with `-race` |
| API surface | Channel as part of public API (careful) | Hide lock inside struct methods |

**Rule of thumb:** Use a channel when goroutines need to signal, hand off work, or shut down. Use a mutex when multiple goroutines read/write the same in-memory structure.

```go
// Mutex — shared cache
type Cache struct {
    mu sync.RWMutex
    m  map[string]string
}

// Channel — worker pool
jobs := make(chan Job, 100)
results := make(chan Result, 100)
```

---

## Buffered vs Unbuffered Channels

| Factor | **Unbuffered** `make(chan T)` | **Buffered** `make(chan T, n)` |
|--------|-------------------------------|--------------------------------|
| Synchronization | Sender blocks until receiver ready (rendezvous) | Sender blocks only when buffer full |
| Guarantee | Strong handoff — work accepted = being processed | Decouples producer/consumer rates |
| Use when | Synchronous pipeline, backpressure by default | Bursty producers, known max in-flight |
| Risk | Deadlock if no receiver | Hidden backlog, memory if consumer stalls |
| Close semantics | Same — receivers drain, then get zero value | Same |
| Size selection | N/A | Match worker count, rate limits, or batch size |

**Rule of thumb:** Start unbuffered. Add buffering only when profiling shows producer blocking and you understand the backlog bound.

---

## database/sql vs sqlx vs GORM

| Factor | **database/sql** | **sqlx** | **GORM** |
|--------|------------------|----------|----------|
| Abstraction | Minimal — you write SQL | Thin — struct scanning, named queries | Full ORM — models, associations, hooks |
| Learning curve | Low (stdlib) | Low–medium | Medium–high |
| SQL control | Full | Full | Partial (raw SQL escape hatch) |
| Migrations | External tool (golang-migrate) | External tool | Built-in AutoMigrate (limited) |
| Performance | Best (no magic) | Near-best | Watch N+1, reflection overhead |
| Type safety | Manual `Scan` | Struct tags + `Get`/`Select` | Struct tags + conventions |
| Transactions | `BeginTx` | `BeginTxx` | `db.Transaction(func(tx *gorm.DB) error)` |
| Best for | Libraries, tight SQL, max control | Services wanting ergonomics + SQL | CRUD-heavy apps, rapid prototyping |
| Pitfalls | Verbose boilerplate | Still manual relation loading | Hidden queries, magic defaults, soft-delete surprises |

**Decision tree:**
1. Need full SQL control or building a library → **database/sql**
2. Want struct scanning without ORM magic → **sqlx**
3. CRUD-heavy domain, team prefers convention → **GORM** (with query logging in dev)

---

## gRPC vs REST vs GraphQL

| Factor | **REST** (HTTP/JSON) | **gRPC** (HTTP/2 + Protobuf) | **GraphQL** |
|--------|----------------------|------------------------------|-------------|
| Contract | OpenAPI optional | `.proto` files (strong contract) | Schema (SDL) |
| Payload | JSON (human-readable, larger) | Binary Protobuf (compact, fast) | JSON (flexible queries) |
| Browser support | Native | Needs grpc-web proxy | Native |
| Streaming | SSE, chunked (limited) | Bidirectional streaming built-in | Subscriptions |
| Versioning | URL/header conventions | Package + field evolution rules | Schema deprecation |
| Tooling | Universal | Strong codegen (Go excels) | Apollo, gqlgen |
| Best for | Public APIs, simple CRUD, caching | Internal microservices, high throughput | Flexible client-driven queries |
| Go libraries | `net/http`, chi, gin | `google.golang.org/grpc` | `gqlgen`, `graphql-go` |

**Rule of thumb:** REST for public/browser APIs; gRPC for service-to-service; GraphQL when clients need flexible field selection (BFF layer).

---

## Caching: In-Memory vs Redis

| Factor | **In-Memory** (sync.Map, map+RWMutex, ristretto) | **Redis** |
|--------|---------------------------------------------------|-----------|
| Latency | Microseconds (same process) | Sub-ms over network |
| Consistency | Per-instance (stale across replicas) | Shared across all instances |
| Eviction | Manual or library (LRU/TTL) | Built-in TTL, LRU policies |
| Persistence | Lost on restart | Optional RDB/AOF |
| Use when | Single instance, hot read path, computed values | Multi-instance, shared sessions, rate limits |
| Pitfalls | Memory pressure, no cluster coherence | Network failures, cache stampede, key design |

**Pattern:** L1 in-process cache + L2 Redis for multi-instance deployments. Always define TTL and invalidation strategy.

---

## Resiliency Patterns Picker

| Problem | Pattern | Go approach |
|---------|---------|-------------|
| Downstream service slow/failing | **Circuit breaker** | `sony/gobreaker`, `afex/hystrix-go` |
| Retry transient failures | **Retry with backoff** | Exponential backoff + jitter; cap attempts |
| Prevent cascade overload | **Bulkhead** | Separate connection pools, worker pools per dependency |
| Graceful degradation | **Fallback** | Return cached/default response |
| Request deadline exceeded | **Timeout** | `context.WithTimeout` on every outbound call |
| Thundering herd on recovery | **Jittered backoff** | Randomize retry delay |
| Duplicate processing | **Idempotency key** | Store key in DB/Redis before side effects |
| Rate limiting | **Token bucket / leaky bucket** | Middleware + Redis or in-memory for single node |
| Partial failure in distributed tx | **Saga** | Choreography (events) or orchestration (coordinator) |
| Read/write scale separation | **CQRS** | Separate read models, eventual consistency |

**Minimum production checklist per outbound call:** context timeout → retry (idempotent only) → circuit breaker → structured logging + metrics.

---

## Value vs Pointer Receiver Rules

| Use **value receiver** `(t T)` | Use **pointer receiver** `(t *T)` |
|--------------------------------|-----------------------------------|
| Small struct (< few pointers worth) | Method mutates receiver |
| Immutable semantics | Large struct (avoid copy) |
| Basic types wrapped in struct | Consistency: if any method needs `*T`, use `*T` for all |
| `sync.Mutex` embedded — **never** value receiver | Slices/maps inside that method reslices or mutates |
| Interface satisfaction where copy is OK | `json.Unmarshaler` and similar need pointer |

**Interface gotcha:** `var r Reader = MyType{}` works if methods have value receivers. If any method has pointer receiver, you need `&MyType{}`.

**Nil receiver:** Pointer-receiver methods can be called on nil receiver if they guard against nil (rare, document clearly).

---

## Performance Optimization Checklist

Work through in order — measure before and after each step.

### 1. Measure
- [ ] Reproduce with realistic load (not micro-benchmark alone)
- [ ] `go test -bench` for hot functions
- [ ] CPU profile: `go test -cpuprofile` or `/debug/pprof`
- [ ] Memory profile: `-memprofile` or heap profile
- [ ] Trace: `go tool trace` for latency gaps
- [ ] `go test -race` in CI

### 2. Algorithm & I/O
- [ ] Fix O(n²) loops, unnecessary allocations in loops
- [ ] Batch DB/network calls; eliminate N+1 queries
- [ ] Right-size connection pools (`SetMaxOpenConns`, idle, lifetime)
- [ ] Use `io.Copy`, buffered I/O, streaming for large payloads
- [ ] Add indexes for slow queries (EXPLAIN ANALYZE)

### 3. Concurrency
- [ ] Right-size worker pool (don't spawn unbounded goroutines)
- [ ] Check for goroutine leaks (`pprof/goroutine`)
- [ ] Reduce lock contention (finer locks, `sync.RWMutex`, sharding)
- [ ] Avoid holding locks during I/O

### 4. Memory & GC
- [ ] Reduce allocations in hot paths (`-benchmem`)
- [ ] Reuse buffers (`sync.Pool` — profile first)
- [ ] Prefer value types / stack allocation where escape analysis allows
- [ ] Check `GOGC` only after profiling (default 100 is usually fine)

### 5. Serialization & Caching
- [ ] Faster JSON (`jsoniter`, or Protobuf for internal APIs)
- [ ] Cache immutable/hot reads with TTL
- [ ] Compress large payloads (gzip middleware)

### Red flags
- Premature `sync.Pool` without allocation proof
- `runtime.GC()` in request path
- Global mutable state without benchmarks
- Cgo in hot path without measurement

---

## Microservices Patterns

| Pattern | When to use | Trade-off |
|---------|-------------|-----------|
| **Saga** | Distributed transactions across services | Eventual consistency; compensating actions required |
| **CQRS** | Read/write load profiles differ greatly | Complexity, eventual read model lag |
| **Event Sourcing** | Full audit trail, temporal queries | Storage growth, replay complexity |
| **API Gateway** | Single entry, auth, routing, rate limit | Potential bottleneck — scale horizontally |
| **BFF** | Different client needs (web vs mobile) | Extra service to maintain |
| **Strangler Fig** | Gradual monolith migration | Long transition period |
| **Outbox** | Reliable event publish with DB write | Requires poller/relay infrastructure |
| **Sidecar** | Cross-cutting concerns (proxy, mesh) | Operational overhead |

### Saga styles

| Choreography | Orchestration |
|--------------|---------------|
| Services react to events | Central coordinator drives steps |
| Looser coupling | Clearer flow visibility |
| Harder to trace | Single point of failure (mitigate with HA) |

---

## Docker / Kubernetes Deployment Checklist

### Dockerfile
- [ ] Multi-stage build (`golang:alpine` builder → `distroless`/`alpine` runtime)
- [ ] `CGO_ENABLED=0` for static binary (unless cgo required)
- [ ] Copy `go.mod`/`go.sum` first for layer cache
- [ ] Run as non-root user
- [ ] Single binary in final image — no source code
- [ ] `.dockerignore` excludes `vendor/`, `.git`, test files
- [ ] HEALTHCHECK or rely on K8s probes

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /server /server
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```

### Kubernetes
- [ ] **Liveness** probe — restart unhealthy pods
- [ ] **Readiness** probe — remove from service endpoints during startup/drain
- [ ] **Startup** probe — slow-starting apps
- [ ] Resource requests & limits (CPU/memory)
- [ ] `replicas >= 2` for HA
- [ ] Config via ConfigMap/Secret, not baked in image
- [ ] Graceful shutdown: handle `SIGTERM`, drain in-flight requests (`http.Server.Shutdown`)
- [ ] `preStop` hook + `terminationGracePeriodSeconds`
- [ ] Horizontal Pod Autoscaler on CPU/custom metrics
- [ ] NetworkPolicy for least-privilege
- [ ] Ingress + TLS termination
- [ ] Structured logs to stdout (JSON) for aggregation

### Go-specific runtime
- [ ] Set `GOMAXPROCS` (automatic since Go 1.5, verify in containers)
- [ ] Expose `/metrics` (Prometheus) and `/health`
- [ ] OpenTelemetry traces propagated via headers
- [ ] Version info endpoint or build-time `-ldflags` injection

---

## Quick Reference Links

| Topic | Module |
|-------|--------|
| Concurrency | [M10–M14](../modules/10-goroutines-scheduler.md) |
| GC & performance | [M22](../modules/22-gc-runtime-performance.md) |
| Architecture | [M25](../modules/25-project-architecture.md) |
| Databases | [M20](../modules/20-databases-orm.md) |
| gRPC | [M21](../modules/21-grpc-protobuf.md) |
| Microservices | [M26](../modules/26-microservices-go.md) |
| Security | [M27](../modules/27-security-distributed.md) |
| Tooling & containers | [M28](../modules/28-tooling-cicd-containers.md) |
