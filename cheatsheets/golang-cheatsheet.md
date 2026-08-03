# Go Quick Reference Cheatsheet

Syntax, types, concurrency primitives, tooling, and popular libraries. For deeper explanations see the [module guides](../README.md#curriculum).

---

## Package & Imports

```go
package main

import (
    "fmt"
    "errors"

    "github.com/gin-gonic/gin"          // third-party
    mypkg "example.com/proj/internal/x" // alias
)
```

| Rule | Detail |
|------|--------|
| Exported | Name starts with **uppercase** |
| Unexported | Name starts with **lowercase** |
| `init()` | Runs before `main`, order = import order |
| `go.mod` | Module path + Go version + dependencies |

---

## Basic Types & Zero Values

| Type | Zero value | Notes |
|------|------------|-------|
| `bool` | `false` | |
| `int`, `int8`…`int64` | `0` | `int`/`uint` size = platform word |
| `float32`, `float64` | `0.0` | |
| `string` | `""` | Immutable UTF-8 byte sequence |
| `byte` | `0` | Alias for `uint8` |
| `rune` | `0` | Alias for `int32` (Unicode code point) |
| pointer | `nil` | |
| slice, map, chan, func, interface | `nil` | `nil` slice is valid (len=0) |
| struct | all fields zeroed | |
| array | all elements zeroed | Value type, fixed size |

```go
var x int           // zero value
y := 42             // short declaration (function body only)
const Pi = 3.14
const (
    StatusOK = 200
    StatusNotFound = 404
)
```

---

## Control Flow

```go
// if — no parens; optional init statement
if n := len(s); n > 0 {
    // ...
} else if n == 0 {
    // ...
}

// for — only loop keyword
for i := 0; i < 10; i++ { }
for i < 10 { }          // while-style
for { }                 // infinite
for i, v := range slice { }
for k, v := range m { }
for v := range ch { }   // channel receive

// switch
switch x {
case 1, 2:
    fallthrough // explicit only
case 3:
default:
}

switch {                // tagless switch
case x > 0:
}

// defer — LIFO at function return
defer f.Close()
defer func() { recover() }()
```

---

## Functions

```go
func add(a, b int) int { return a + b }

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Variadic
func sum(nums ...int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}

// Named returns
func split(s string) (first, rest string) {
    i := strings.Index(s, " ")
    first, rest = s[:i], s[i+1:]
    return // naked return
}
```

| Pattern | Example |
|---------|---------|
| Method (value receiver) | `func (p Point) Dist() float64` |
| Method (pointer receiver) | `func (p *Point) Move(dx, dy int)` |
| Function type | `type Handler func(http.ResponseWriter, *http.Request)` |
| Closures | Capture variables by reference |

---

## Arrays vs Slices

```go
var arr [3]int = [3]int{1, 2, 3}   // fixed size, value type
s := []int{1, 2, 3}                 // slice header: ptr, len, cap
s = make([]int, 5)                  // len=5, cap=5
s = make([]int, 0, 10)              // len=0, cap=10
s = append(s, 4, 5)
sub := s[1:3]                       // shares backing array
copy(dst, src)
```

---

## Maps

```go
m := make(map[string]int)
m["key"] = 42
v, ok := m["key"]   // ok=false if missing
delete(m, "key")
for k, v := range m { }
```

---

## Structs & Embedding

```go
type Person struct {
    Name string `json:"name" db:"name"`
    Age  int    `json:"age,omitempty"`
}

type Employee struct {
    Person          // anonymous embedding — promotes fields
    ID       int
    Dept     string
}

e := Employee{Person: Person{Name: "Ada"}, ID: 1}
fmt.Println(e.Name) // promoted
```

### Common Struct Tags

| Tag | Used by | Example |
|-----|---------|---------|
| `json:"field"` | `encoding/json` | `json:"user_id"` |
| `json:"-"` | skip field | |
| `json:",omitempty"` | omit if zero | |
| `db:"column"` | sqlx | `db:"created_at"` |
| `gorm:"primaryKey"` | GORM | `gorm:"column:id"` |
| `validate:"required"` | go-playground/validator | |
| `xml:"name"` | `encoding/xml` | |
| `yaml:"key"` | gopkg.in/yaml.v3 | |

---

## Interfaces

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Implicit satisfaction — no "implements" keyword
type MyReader struct{}
func (r MyReader) Read(p []byte) (int, error) { return 0, io.EOF }

var r Reader = MyReader{} // OK

// Type assertion
v, ok := i.(string)
v := i.(string) // panics if wrong type

// Type switch
switch v := i.(type) {
case string:
case int:
default:
}

// Empty interface (any since Go 1.18)
var x any = 42
```

---

## Pointers

```go
x := 10
p := &x       // address-of
*p = 20       // dereference
new(int)      // *int on heap, zeroed
```

| Use pointer receiver when | Use value receiver when |
|---------------------------|-------------------------|
| Method mutates receiver | Small, immutable struct |
| Struct is large (copy cost) | Value semantics desired |
| Consistency (some methods need `*`) | No mutation needed |

---

## Error Handling

```go
if err != nil {
    return fmt.Errorf("load config: %w", err) // wrap
}

// Custom error
type NotFoundError struct {
    ID string
}
func (e NotFoundError) Error() string {
    return fmt.Sprintf("not found: %s", e.ID)
}

// errors.Is / errors.As
if errors.Is(err, os.ErrNotExist) { }
var nfe NotFoundError
if errors.As(err, &nfe) { }

// panic / recover — only for truly unrecoverable or init
defer func() {
    if r := recover(); r != nil { log.Println(r) }
}()
```

---

## Concurrency Primitives

### Goroutines

```go
go func() { fmt.Println("async") }()
```

### Channels

```go
ch := make(chan int)       // unbuffered — sync handoff
buf := make(chan int, 10)  // buffered — async up to cap

ch <- 1       // send
v := <-ch     // receive
close(ch)     // no more sends; receivers drain

for v := range ch { } // until closed
```

### select

```go
select {
case msg := <-ch1:
case ch2 <- val:
case <-time.After(5 * time.Second):
case <-ctx.Done():
default: // non-blocking
}
```

### sync Package

```go
var mu sync.Mutex
mu.Lock(); defer mu.Unlock()

var wg sync.WaitGroup
wg.Add(1)
go func() { defer wg.Done(); work() }()
wg.Wait()

var once sync.Once
once.Do(func() { init() })

var m sync.Map // concurrent map (niche use)
```

### context

```go
ctx := context.Background()
ctx, cancel := context.WithCancel(ctx)
defer cancel()

ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
defer cancel()

ctx = context.WithValue(ctx, key, value) // sparingly

select {
case <-ctx.Done():
    return ctx.Err() // context.Canceled or DeadlineExceeded
}
```

### Common Patterns

```go
// Worker pool
jobs := make(chan Job, 100)
for w := 0; w < numWorkers; w++ {
    go worker(jobs)
}

// Fan-out / fan-in
// errgroup
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return task1(ctx) })
if err := g.Wait(); err != nil { }
```

---

## Generics (Go 1.18+)

```go
func Map[T, U any](s []T, f func(T) U) []U {
    out := make([]U, len(s))
    for i, v := range s { out[i] = f(v) }
    return out
}

type Stack[T any] struct { items []T }
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
```

---

## Testing Commands

```bash
go test ./...                    # all packages
go test -v ./pkg                 # verbose
go test -run TestName            # single test
go test -count=1                 # disable cache
go test -race ./...              # race detector
go test -cover ./...             # coverage
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out

# Benchmarks
go test -bench=. -benchmem ./...
go test -bench=BenchmarkFoo -benchtime=5s

# Fuzzing (Go 1.18+)
go test -fuzz=FuzzReverse -fuzztime=30s
```

### Test Skeleton

```go
func TestAdd(t *testing.T) {
    got := add(2, 3)
    if got != 5 {
        t.Fatalf("got %d, want 5", got)
    }
}

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        add(2, 3)
    }
}
```

---

## go CLI

| Command | Purpose |
|---------|---------|
| `go mod init <path>` | Create module |
| `go mod tidy` | Add/remove deps |
| `go get pkg@v1.2.3` | Add/update dependency |
| `go build -o bin/app ./cmd/server` | Compile |
| `go run .` | Build + run |
| `go install ./cmd/...` | Install to `$GOBIN` |
| `go vet ./...` | Static analysis |
| `gofmt -w .` | Format |
| `goimports -w .` | Format + fix imports |
| `go doc fmt.Println` | Documentation |
| `go list -m all` | List module versions |
| `go work init` | Workspace (monorepo) |

### Profiling

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.
go tool pprof cpu.prof
go tool pprof -http=:8080 mem.prof

# Runtime HTTP endpoint
import _ "net/http/pprof"
# GET /debug/pprof/
```

---

## net/http Essentials

```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "ok")
})
log.Fatal(http.ListenAndServe(":8080", nil))

// Handler interface
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Middleware pattern
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

---

## database/sql Essentials

```go
db, err := sql.Open("postgres", dsn) // pool, not a connection
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

row := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id=$1", id)
err = row.Scan(&name)

tx, err := db.BeginTx(ctx, nil)
_, err = tx.ExecContext(ctx, "INSERT ...")
err = tx.Commit() // or tx.Rollback()
```

---

## Popular Libraries

| Library | Use case | Import |
|---------|----------|--------|
| **chi** | Lightweight HTTP router + middleware | `github.com/go-chi/chi/v5` |
| **gin** | Fast HTTP framework with binding/validation | `github.com/gin-gonic/gin` |
| **gRPC** | RPC over HTTP/2 + Protobuf | `google.golang.org/grpc` |
| **GORM** | ORM for SQL databases | `gorm.io/gorm` |
| **sqlx** | Extensions over database/sql | `github.com/jmoiron/sqlx` |
| **viper** | Configuration (files, env) | `github.com/spf13/viper` |
| **zap** / **zerolog** | Structured logging | `go.uber.org/zap` |
| **testify** | Test assertions & mocks | `github.com/stretchr/testify` |
| **otel** | OpenTelemetry tracing/metrics | `go.opentelemetry.io/otel` |
| **redis** | Redis client | `github.com/redis/go-redis/v9` |
| **migrate** | DB migrations | `github.com/golang-migrate/migrate/v4` |

### chi Example

```go
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Get("/users/{id}", getUser)
r.Route("/api", func(r chi.Router) {
    r.Post("/items", createItem)
})
```

### gin Example

```go
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
r.Run(":8080")
```

### gRPC Server Sketch

```go
lis, _ := net.Listen("tcp", ":50051")
grpcServer := grpc.NewServer()
pb.RegisterGreeterServer(grpcServer, &server{})
grpcServer.Serve(lis)
```

### GORM Sketch

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
db.AutoMigrate(&User{})
db.Create(&User{Name: "Ada"})
db.Where("name = ?", "Ada").First(&user)
```

---

## Useful Standard Library Packages

| Package | Purpose |
|---------|---------|
| `context` | Cancellation, deadlines, request-scoped values |
| `encoding/json` | JSON marshal/unmarshal |
| `fmt` | Formatted I/O |
| `io` / `io/fs` | Reader/Writer abstractions |
| `net/http` | HTTP client & server |
| `os` / `os/exec` | OS interface, subprocess |
| `path/filepath` | Cross-platform paths |
| `strconv` | String ↔ number conversion |
| `strings` / `bytes` | String/byte utilities |
| `sync` / `sync/atomic` | Concurrency primitives |
| `time` | Time and timers |
| `log/slog` | Structured logging (Go 1.21+) |

---

## Build Tags & Conditional Compilation

```go
//go:build linux
// +build linux

package mypkg
```

---

## Quick Memory Model Reminders

- Goroutines are cheap (~2 KB stack, grows as needed); OS threads are expensive.
- Passing a slice copies the **header**, not the backing array.
- Maps are reference types (internally a pointer); not safe for concurrent writes without sync.
- `range` over slice copies the element value (use index for large structs).
- `defer` has small overhead; avoid in tight hot loops.

---

## Further Reading

- [Go spec](https://go.dev/ref/spec)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go memory model](https://go.dev/ref/mem)
- [Module reference](https://go.dev/ref/mod)
