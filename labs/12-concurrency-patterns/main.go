package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// --- Bounded buffer (channel-backed) ---

type BoundedBuffer[T any] struct {
	ch chan T
}

func NewBoundedBuffer[T any](capacity int) *BoundedBuffer[T] {
	return &BoundedBuffer[T]{ch: make(chan T, capacity)}
}

func (b *BoundedBuffer[T]) Put(ctx context.Context, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.ch <- v:
		return nil
	}
}

func (b *BoundedBuffer[T]) Get(ctx context.Context) (T, error) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case v, ok := <-b.ch:
		if !ok {
			return zero, errors.New("buffer closed")
		}
		return v, nil
	}
}

func (b *BoundedBuffer[T]) Close() { close(b.ch) }

// --- Worker pool ---

func workerPool(jobs <-chan int, results chan<- int, workers int) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				results <- j * j
			}
		}(i)
	}
	return &wg
}

// --- errgroup demo ---

func fetch(ctx context.Context, name string, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		if name == "billing" && rand.Float32() < 0.15 {
			return fmt.Errorf("%s service unavailable", name)
		}
		fmt.Printf("  fetched %s\n", name)
		return nil
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("== Bounded buffer ==")
	buf := NewBoundedBuffer[string](2)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = buf.Put(ctx, "a")
	_ = buf.Put(ctx, "b")
	v, err := buf.Get(ctx)
	fmt.Println("got", v, err)
	buf.Close()

	fmt.Println("\n== Worker pool ==")
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	wg := workerPool(jobs, results, 3)
	for i := 1; i <= 8; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	close(results)
	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println("sum of squares:", sum)

	fmt.Println("\n== errgroup ==")
	g, gctx := errgroup.WithContext(context.Background())
	services := []struct {
		name string
		d    time.Duration
	}{
		{"catalog", 100 * time.Millisecond},
		{"inventory", 150 * time.Millisecond},
		{"billing", 120 * time.Millisecond},
	}
	for _, s := range services {
		s := s
		g.Go(func() error {
			return fetch(gctx, s.name, s.d)
		})
	}
	if err := g.Wait(); err != nil {
		log.Printf("errgroup failed: %v", err)
	} else {
		fmt.Println("all services OK")
	}
}
