package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
)

func slowFib(n int) int {
	if n <= 1 {
		return n
	}
	return slowFib(n-1) + slowFib(n-2)
}

func workLoop(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			_ = slowFib(35)
		}
	}
}

func main() {
	addr := ":6060"
	if v := os.Getenv("PPROF_ADDR"); v != "" {
		addr = v
	}

	stop := make(chan struct{})
	defer close(stop)
	workers := runtime.NumCPU()
	if w := os.Getenv("WORKERS"); w != "" {
		if n, err := parsePositiveInt(w); err == nil {
			workers = n
		}
	}
	for i := 0; i < workers; i++ {
		go workLoop(stop)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	log.Printf("pprof on http://localhost%s/debug/pprof/", addr)
	log.Printf("CPU load: %d workers computing fib(35)", workers)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, os.ErrInvalid
		}
		n = n*10 + int(c-'0')
	}
	if n <= 0 {
		return 0, os.ErrInvalid
	}
	return n, nil
}
