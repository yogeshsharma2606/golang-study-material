package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price_cents"`
}

type store struct {
	mu    sync.RWMutex
	items map[int]item
	next  int
}

func newStore() *store {
	s := &store{items: make(map[int]item), next: 1}
	s.items[1] = item{ID: 1, Name: "Notebook", Price: 499}
	return s
}

func (s *store) list() []item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	return out
}

func (s *store) get(id int) (item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	return it, ok
}

func (s *store) create(name string, price int) item {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	it := item{ID: id, Name: name, Price: price}
	s.items[id] = it
	return it
}

type ctxKey string

const requestIDKey ctxKey = "request_id"

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = strconv.FormatInt(time.Now().UnixNano(), 36)
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rid, _ := r.Context().Value(requestIDKey).(string)
		lw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf("request_id=%s method=%s path=%s status=%d duration=%s",
			rid, r.Method, r.URL.Path, lw.status, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

type api struct {
	store *store
}

func (a *api) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.handleHealth)
	mux.HandleFunc("GET /api/items", a.handleListItems)
	mux.HandleFunc("GET /api/items/{id}", a.handleGetItem)
	mux.HandleFunc("POST /api/items", a.handleCreateItem)
	var h http.Handler = mux
	h = requestIDMiddleware(h)
	h = loggingMiddleware(h)
	return h
}

func (a *api) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *api) handleListItems(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.store.list())
}

func (a *api) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	it, ok := a.store.get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (a *api) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name       string `json:"name"`
		PriceCents int    `json:"price_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Name == "" || body.PriceCents < 0 {
		writeError(w, http.StatusBadRequest, "name required and price_cents >= 0")
		return
	}
	it := a.store.create(body.Name, body.PriceCents)
	writeJSON(w, http.StatusCreated, it)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      (&api{store: newStore()}).routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
