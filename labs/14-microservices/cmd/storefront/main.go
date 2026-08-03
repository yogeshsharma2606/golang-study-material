package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	inventoryURL := os.Getenv("INVENTORY_URL")
	if inventoryURL == "" {
		inventoryURL = "http://localhost:8081"
	}
	addr := ":8080"
	if v := os.Getenv("STOREFRONT_ADDR"); v != "" {
		addr = v
	}

	client := &http.Client{Timeout: 3 * time.Second}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"service": "storefront", "status": "ok"})
	})
	mux.HandleFunc("GET /products/{sku}", func(w http.ResponseWriter, r *http.Request) {
		sku := r.PathValue("sku")
		url := fmt.Sprintf("%s/stock/%s", inventoryURL, sku)
		resp, err := client.Get(url)
		if err != nil {
			http.Error(w, "inventory unreachable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "inventory error", resp.StatusCode)
			return
		}
		var stock struct {
			SKU   string `json:"sku"`
			Count int    `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
			http.Error(w, "bad upstream json", http.StatusBadGateway)
			return
		}
		writeJSON(w, map[string]any{
			"sku":       stock.SKU,
			"available": stock.Count > 0,
			"count":     stock.Count,
		})
	})

	log.Printf("storefront listening on %s (inventory %s)", addr, inventoryURL)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
