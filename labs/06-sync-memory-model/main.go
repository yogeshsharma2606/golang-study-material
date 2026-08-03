package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Race demo: uncomment runRace() in main and run with -race to see the detector fire.
func runRace() {
	var counter int
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // data race without synchronization
		}()
	}
	wg.Wait()
	fmt.Println("race demo counter (undefined):", counter)
}

func runMutex() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("mutex counter:", counter)
}

func runOnce() {
	var once sync.Once
	var initialized int32
	var wg sync.WaitGroup
	init := func() {
		atomic.StoreInt32(&initialized, 1)
		fmt.Println("sync.Once init ran")
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(init)
		}()
	}
	wg.Wait()
	fmt.Println("initialized flag:", atomic.LoadInt32(&initialized))
}

func main() {
	// runRace() // enable with: go run -race .
	runMutex()
	runOnce()
}