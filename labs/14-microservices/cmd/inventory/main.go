package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type stock struct {
	SKU   string `json:"sku"`
	Count int    `json:"count"`
}

var (
	mu    sync.RWMutex
	items = map[string]int{"go-book": 5, "mug": 12}
)

func main() {
	addr := ":8081"
	if v := os.Getenv("INVENTORY_ADDR"); v != "" {
		addr = v
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"service": "inventory", "status": "ok"})
	})
	mux.HandleFunc("GET /stock/{sku}", handleGetStock)
	log.Printf("inventory listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleGetStock(w http.ResponseWriter, r *http.Request) {
	sku := r.PathValue("sku")
	mu.RLock()
	count, ok := items[sku]
	mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, stock{SKU: sku, Count: count})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
