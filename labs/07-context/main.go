package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// simulateSlowWork stands in for DB/API work that respects cancellation.
func simulateSlowWork(ctx context.Context, name string) error {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println(name, "finished")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// handler mimics an HTTP handler that derives a timeout from the request context.
func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	if err := simulateSlowWork(ctx, r.URL.Path); err != nil {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	fmt.Fprintf(w, "ok: %s\n", r.URL.Path)
}

func main() {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	fmt.Println("status:", rec.Code)
	fmt.Println("body:", rec.Body.String())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	if err := simulateSlowWork(ctx, "manual"); err != nil {
		fmt.Println("manual cancel:", err)
	}
}