package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// --- Functional options ---

type Server struct {
	addr    string
	timeout time.Duration
	tls     bool
}

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) { s.addr = addr }
}

func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

func WithTLS(enabled bool) Option {
	return func(s *Server) { s.tls = enabled }
}

func NewServer(opts ...Option) *Server {
	s := &Server{addr: ":8080", timeout: 5 * time.Second}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// --- Decorator middleware ---

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) { f(w, r) }

func WithLogging(next http.Handler) http.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[log] %s %s %s\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func WithRecovery(next http.Handler) http.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				fmt.Println("[recover]", rec)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// --- Strategy ---

type Pricer interface {
	Price(baseCents int) int
}

type RetailPricer struct{}

func (RetailPricer) Price(base int) int { return base }

type MemberPricer struct{ Discount float64 }

func (m MemberPricer) Price(base int) int {
	return int(float64(base) * (1 - m.Discount))
}

type PromoPricer struct{ OffCents int }

func (p PromoPricer) Price(base int) int {
	out := base - p.OffCents
	if out < 0 {
		return 0
	}
	return out
}

func checkout(label string, p Pricer, base int) {
	fmt.Printf("  %s: %d cents -> %d cents\n", label, base, p.Price(base))
}

func main() {
	srv := NewServer(WithAddr(":9090"), WithTLS(true), WithTimeout(2*time.Second))
	fmt.Printf("Server config: addr=%s timeout=%s tls=%v\n", srv.addr, srv.timeout, srv.tls)

	mux := http.NewServeMux()
	mux.HandleFunc("/hi", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})
	handler := WithRecovery(WithLogging(mux))
	fmt.Println("Decorator chain: Recovery -> Logging -> mux (not starting server in demo)")

	base := 1000
	fmt.Println("Strategy pricing (base 1000 cents):")
	checkout("retail", RetailPricer{}, base)
	checkout("member 10%", MemberPricer{Discount: 0.10}, base)
	checkout("promo $3 off", PromoPricer{OffCents: 300}, base)

	_ = handler
	_ = strings.Builder{}
}
